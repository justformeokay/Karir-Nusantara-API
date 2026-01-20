#!/bin/bash

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# API Base URL
API_URL="http://localhost:8081/api/v1"

# Test company email
COMPANY_EMAIL="info@karyadeveloperindonesia.com"

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}ADMIN VERIFICATION TEST SCRIPT${NC}"
echo -e "${BLUE}========================================${NC}"

# Step 1: Admin Login
echo -e "\n${YELLOW}[STEP 1] Admin Login...${NC}"
ADMIN_LOGIN_RESPONSE=$(curl -s -X POST "${API_URL}/admin/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@karircorp.com",
    "password": "admin123456"
  }')

ADMIN_TOKEN=$(echo $ADMIN_LOGIN_RESPONSE | jq -r '.data.token' 2>/dev/null)

if [ -z "$ADMIN_TOKEN" ] || [ "$ADMIN_TOKEN" == "null" ]; then
  echo -e "${RED}❌ Admin login failed!${NC}"
  echo "Response: $ADMIN_LOGIN_RESPONSE"
  exit 1
fi

echo -e "${GREEN}✓ Admin login successful!${NC}"
echo "Admin Token: $ADMIN_TOKEN"

# Step 2: Get Companies List
echo -e "\n${YELLOW}[STEP 2] Fetching companies list...${NC}"
COMPANIES_RESPONSE=$(curl -s -X GET "${API_URL}/admin/companies?status=pending&page=1&page_size=10" \
  -H "Authorization: Bearer $ADMIN_TOKEN")

echo -e "${GREEN}✓ Companies list:${NC}"
echo $COMPANIES_RESPONSE | jq '.' 2>/dev/null || echo $COMPANIES_RESPONSE

# Step 3: Find company by email
echo -e "\n${YELLOW}[STEP 3] Finding company with email: $COMPANY_EMAIL${NC}"
COMPANY_ID=$(echo $COMPANIES_RESPONSE | jq ".data[] | select(.email==\"$COMPANY_EMAIL\") | .id" 2>/dev/null | head -1)

if [ -z "$COMPANY_ID" ] || [ "$COMPANY_ID" == "null" ]; then
  echo -e "${YELLOW}⚠ Company not found in pending list, trying all companies...${NC}"
  COMPANIES_ALL=$(curl -s -X GET "${API_URL}/admin/companies?page=1&page_size=100" \
    -H "Authorization: Bearer $ADMIN_TOKEN")
  COMPANY_ID=$(echo $COMPANIES_ALL | jq ".data[] | select(.email==\"$COMPANY_EMAIL\") | .id" 2>/dev/null | head -1)
fi

if [ -z "$COMPANY_ID" ] || [ "$COMPANY_ID" == "null" ]; then
  echo -e "${RED}❌ Company not found!${NC}"
  exit 1
fi

echo -e "${GREEN}✓ Found company ID: $COMPANY_ID${NC}"

# Step 4: Get Company Details
echo -e "\n${YELLOW}[STEP 4] Getting company details...${NC}"
COMPANY_DETAILS=$(curl -s -X GET "${API_URL}/admin/companies/$COMPANY_ID" \
  -H "Authorization: Bearer $ADMIN_TOKEN")

echo -e "${GREEN}✓ Company details:${NC}"
echo $COMPANY_DETAILS | jq '.' 2>/dev/null || echo $COMPANY_DETAILS

# Step 5: Verify Company (Approve)
echo -e "\n${YELLOW}[STEP 5] Verifying company (APPROVE)...${NC}"
VERIFY_RESPONSE=$(curl -s -X POST "${API_URL}/admin/companies/$COMPANY_ID/verify" \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "action": "approve",
    "notes": "Dokumen dan profil sudah lengkap dan sesuai dengan ketentuan."
  }')

echo $VERIFY_RESPONSE | jq '.' 2>/dev/null || echo $VERIFY_RESPONSE

if echo $VERIFY_RESPONSE | grep -q "berhasil"; then
  echo -e "${GREEN}✓ Company verification successful!${NC}"
else
  echo -e "${YELLOW}⚠ Check verification response above${NC}"
fi

# Step 6: Get Updated Company Details
echo -e "\n${YELLOW}[STEP 6] Getting updated company details...${NC}"
UPDATED_COMPANY=$(curl -s -X GET "${API_URL}/admin/companies/$COMPANY_ID" \
  -H "Authorization: Bearer $ADMIN_TOKEN")

echo -e "${GREEN}✓ Updated company status:${NC}"
echo $UPDATED_COMPANY | jq '.data | {id, email, company_name, company_status, is_verified, created_at}' 2>/dev/null || echo $UPDATED_COMPANY

# Step 7: Company Login with verified account
echo -e "\n${YELLOW}[STEP 7] Testing company login after verification...${NC}"
COMPANY_LOGIN=$(curl -s -X POST "${API_URL}/auth/login" \
  -H "Content-Type: application/json" \
  -d "{
    \"email\": \"$COMPANY_EMAIL\",
    \"password\": \"company123456\"
  }")

COMPANY_TOKEN=$(echo $COMPANY_LOGIN | jq -r '.data.user.token' 2>/dev/null || echo $COMPANY_LOGIN | jq -r '.data.token' 2>/dev/null)

if [ -z "$COMPANY_TOKEN" ] || [ "$COMPANY_TOKEN" == "null" ]; then
  echo -e "${YELLOW}⚠ Company login failed (might be normal if using different password):${NC}"
  echo $COMPANY_LOGIN | jq '.' 2>/dev/null || echo $COMPANY_LOGIN
else
  echo -e "${GREEN}✓ Company login successful!${NC}"
  
  # Get current user info
  echo -e "\n${YELLOW}[STEP 8] Getting verified company profile...${NC}"
  COMPANY_PROFILE=$(curl -s -X GET "${API_URL}/auth/me" \
    -H "Authorization: Bearer $COMPANY_TOKEN")
  
  echo -e "${GREEN}✓ Verified company profile:${NC}"
  echo $COMPANY_PROFILE | jq '.data | {id, email, company_name, company_status, verification_status, is_verified}' 2>/dev/null || echo $COMPANY_PROFILE
fi

echo -e "\n${BLUE}========================================${NC}"
echo -e "${GREEN}✓ VERIFICATION TEST COMPLETED!${NC}"
echo -e "${BLUE}========================================${NC}"

echo -e "\n${YELLOW}Summary:${NC}"
echo "- Company ID: $COMPANY_ID"
echo "- Company Email: $COMPANY_EMAIL"
echo "- Status should be: 'verified'"
echo ""
echo -e "${YELLOW}Next steps:${NC}"
echo "1. Login to company dashboard"
echo "2. Test job creation feature"
echo "3. Test document uploads"
echo "4. Test all dashboard features"
