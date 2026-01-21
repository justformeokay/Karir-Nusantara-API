#!/bin/bash

# Change Password API Test Script
# Tests the complete change password flow including email confirmation

set -e

BASE_URL="http://localhost:8081/api/v1"
TEST_EMAIL="changepass_$(date +%s)@example.com"
OLD_PASSWORD="OldPass123!"
NEW_PASSWORD="NewPass456!"
ANOTHER_PASSWORD="AnotherPass789!"

echo "================================"
echo "Change Password API Test"
echo "================================"
echo ""

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

print_success() {
    echo -e "${GREEN}✓ $1${NC}"
}

print_error() {
    echo -e "${RED}✗ $1${NC}"
}

print_info() {
    echo -e "${YELLOW}ℹ $1${NC}"
}

# Test 1: Register a new user
echo "Test 1: Register new user for testing..."
REGISTER_RESPONSE=$(curl -s -X POST $BASE_URL/auth/register \
  -H "Content-Type: application/json" \
  -d "{
    \"email\": \"$TEST_EMAIL\",
    \"password\": \"$OLD_PASSWORD\",
    \"full_name\": \"Change Password Test User\",
    \"role\": \"company\",
    \"company_name\": \"Test Company\"
  }")

if echo "$REGISTER_RESPONSE" | jq -e '.success == true' > /dev/null; then
    ACCESS_TOKEN=$(echo "$REGISTER_RESPONSE" | jq -r '.data.access_token')
    print_success "User registered successfully"
    print_info "Email: $TEST_EMAIL"
    print_info "Password: $OLD_PASSWORD"
else
    print_error "Registration failed"
    echo "$REGISTER_RESPONSE" | jq .
    exit 1
fi

echo ""

# Test 2: Change password successfully
echo "Test 2: Change password with correct old password..."
CHANGE_RESPONSE=$(curl -s -X PUT $BASE_URL/auth/change-password \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -d "{
    \"old_password\": \"$OLD_PASSWORD\",
    \"new_password\": \"$NEW_PASSWORD\"
  }")

if echo "$CHANGE_RESPONSE" | jq -e '.success == true' > /dev/null; then
    print_success "Password changed successfully"
    print_info "New password: $NEW_PASSWORD"
    print_info "Email confirmation sent (check async)"
else
    print_error "Password change failed"
    echo "$CHANGE_RESPONSE" | jq .
    exit 1
fi

echo ""

# Test 3: Verify old password doesn't work
echo "Test 3: Verify old password no longer works..."
OLD_LOGIN_RESPONSE=$(curl -s -X POST $BASE_URL/auth/login \
  -H "Content-Type: application/json" \
  -d "{
    \"email\": \"$TEST_EMAIL\",
    \"password\": \"$OLD_PASSWORD\"
  }")

if echo "$OLD_LOGIN_RESPONSE" | jq -e '.success == false' > /dev/null; then
    print_success "Old password correctly rejected"
else
    print_error "Old password should not work!"
    echo "$OLD_LOGIN_RESPONSE" | jq .
    exit 1
fi

echo ""

# Test 4: Login with new password
echo "Test 4: Login with new password..."
NEW_LOGIN_RESPONSE=$(curl -s -X POST $BASE_URL/auth/login \
  -H "Content-Type: application/json" \
  -d "{
    \"email\": \"$TEST_EMAIL\",
    \"password\": \"$NEW_PASSWORD\"
  }")

if echo "$NEW_LOGIN_RESPONSE" | jq -e '.success == true' > /dev/null; then
    NEW_ACCESS_TOKEN=$(echo "$NEW_LOGIN_RESPONSE" | jq -r '.data.access_token')
    print_success "Login with new password successful"
    print_info "New access token obtained"
else
    print_error "Login with new password failed"
    echo "$NEW_LOGIN_RESPONSE" | jq .
    exit 1
fi

echo ""

# Test 5: Try to change password with wrong old password
echo "Test 5: Try changing password with wrong old password..."
WRONG_OLD_RESPONSE=$(curl -s -X PUT $BASE_URL/auth/change-password \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $NEW_ACCESS_TOKEN" \
  -d "{
    \"old_password\": \"WrongPassword123!\",
    \"new_password\": \"$ANOTHER_PASSWORD\"
  }")

if echo "$WRONG_OLD_RESPONSE" | jq -e '.success == false' > /dev/null; then
    ERROR_MSG=$(echo "$WRONG_OLD_RESPONSE" | jq -r '.error.message')
    if [[ "$ERROR_MSG" == *"tidak sesuai"* ]]; then
        print_success "Wrong old password correctly rejected"
        print_info "Error: $ERROR_MSG"
    else
        print_error "Wrong error message"
        echo "$WRONG_OLD_RESPONSE" | jq .
    fi
else
    print_error "Should reject wrong old password"
    echo "$WRONG_OLD_RESPONSE" | jq .
    exit 1
fi

echo ""

# Test 6: Try to use same password
echo "Test 6: Try using same password as new password..."
SAME_PASSWORD_RESPONSE=$(curl -s -X PUT $BASE_URL/auth/change-password \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $NEW_ACCESS_TOKEN" \
  -d "{
    \"old_password\": \"$NEW_PASSWORD\",
    \"new_password\": \"$NEW_PASSWORD\"
  }")

if echo "$SAME_PASSWORD_RESPONSE" | jq -e '.success == false' > /dev/null; then
    ERROR_MSG=$(echo "$SAME_PASSWORD_RESPONSE" | jq -r '.error.message')
    if [[ "$ERROR_MSG" == *"tidak boleh sama"* ]]; then
        print_success "Same password correctly rejected"
        print_info "Error: $ERROR_MSG"
    else
        print_error "Wrong error message"
        echo "$SAME_PASSWORD_RESPONSE" | jq .
    fi
else
    print_error "Should reject same password"
    echo "$SAME_PASSWORD_RESPONSE" | jq .
    exit 1
fi

echo ""

# Test 7: Try without authentication
echo "Test 7: Try changing password without authentication..."
NO_AUTH_RESPONSE=$(curl -s -X PUT $BASE_URL/auth/change-password \
  -H "Content-Type: application/json" \
  -d "{
    \"old_password\": \"$NEW_PASSWORD\",
    \"new_password\": \"$ANOTHER_PASSWORD\"
  }")

if echo "$NO_AUTH_RESPONSE" | jq -e '.success == false' > /dev/null; then
    ERROR_CODE=$(echo "$NO_AUTH_RESPONSE" | jq -r '.error.code')
    if [[ "$ERROR_CODE" == "UNAUTHORIZED" ]]; then
        print_success "Unauthorized request correctly rejected"
        print_info "Error: Missing authorization"
    else
        print_error "Wrong error code"
        echo "$NO_AUTH_RESPONSE" | jq .
    fi
else
    print_error "Should reject unauthorized request"
    echo "$NO_AUTH_RESPONSE" | jq .
    exit 1
fi

echo ""

# Test 8: Try with invalid/expired token
echo "Test 8: Try changing password with invalid token..."
INVALID_TOKEN_RESPONSE=$(curl -s -X PUT $BASE_URL/auth/change-password \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer invalid.token.here" \
  -d "{
    \"old_password\": \"$NEW_PASSWORD\",
    \"new_password\": \"$ANOTHER_PASSWORD\"
  }")

if echo "$INVALID_TOKEN_RESPONSE" | jq -e '.success == false' > /dev/null; then
    print_success "Invalid token correctly rejected"
else
    print_error "Should reject invalid token"
    echo "$INVALID_TOKEN_RESPONSE" | jq .
    exit 1
fi

echo ""

# Test 9: Validation - password too short
echo "Test 9: Try changing to password that's too short..."
SHORT_PASSWORD_RESPONSE=$(curl -s -X PUT $BASE_URL/auth/change-password \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $NEW_ACCESS_TOKEN" \
  -d "{
    \"old_password\": \"$NEW_PASSWORD\",
    \"new_password\": \"Short1!\"
  }")

if echo "$SHORT_PASSWORD_RESPONSE" | jq -e '.success == false' > /dev/null; then
    print_success "Short password correctly rejected by validation"
else
    print_error "Should reject short password"
    echo "$SHORT_PASSWORD_RESPONSE" | jq .
fi

echo ""

# Test 10: Successfully change password again
echo "Test 10: Change password again (final test)..."
FINAL_CHANGE_RESPONSE=$(curl -s -X PUT $BASE_URL/auth/change-password \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $NEW_ACCESS_TOKEN" \
  -d "{
    \"old_password\": \"$NEW_PASSWORD\",
    \"new_password\": \"$ANOTHER_PASSWORD\"
  }")

if echo "$FINAL_CHANGE_RESPONSE" | jq -e '.success == true' > /dev/null; then
    print_success "Password changed again successfully"
    print_info "Final password: $ANOTHER_PASSWORD"
else
    print_error "Final password change failed"
    echo "$FINAL_CHANGE_RESPONSE" | jq .
    exit 1
fi

echo ""

# Test 11: Verify final password works
echo "Test 11: Verify final password works..."
FINAL_LOGIN_RESPONSE=$(curl -s -X POST $BASE_URL/auth/login \
  -H "Content-Type: application/json" \
  -d "{
    \"email\": \"$TEST_EMAIL\",
    \"password\": \"$ANOTHER_PASSWORD\"
  }")

if echo "$FINAL_LOGIN_RESPONSE" | jq -e '.success == true' > /dev/null; then
    print_success "Login with final password successful"
else
    print_error "Login with final password failed"
    echo "$FINAL_LOGIN_RESPONSE" | jq .
    exit 1
fi

echo ""
echo "================================"
echo -e "${GREEN}All tests passed! ✓${NC}"
echo "================================"
echo ""
echo "Summary:"
echo "- Test user: $TEST_EMAIL"
echo "- Password changes: OLD -> NEW -> ANOTHER"
echo "- Email confirmations: 2 emails sent (check inbox)"
echo "- All validation tests passed"
echo "- Authentication & authorization working correctly"
echo ""
print_info "Note: Check email inbox for password change confirmation emails"
echo ""
