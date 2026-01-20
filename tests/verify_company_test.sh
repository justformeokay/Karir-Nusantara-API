#!/bin/bash

# Company Verification Test with Specific Credentials
# Test company: info@karyadeveloperindonesia.com
# Password: Justformeokay23@

API="http://localhost:8081/api/v1"

echo "======================================"
echo "COMPANY VERIFICATION TEST"
echo "======================================"

# Step 1: Company Login
echo ""
echo "1️⃣  Company Login..."
COMPANY_LOGIN=$(curl -s -X POST "$API/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "info@karyadeveloperindonesia.com",
    "password": "Justformeokay23@"
  }')

COMPANY_TOKEN=$(echo $COMPANY_LOGIN | jq -r '.data.access_token // .data.token // .data.user.token // empty')
if [ -z "$COMPANY_TOKEN" ]; then
  echo "❌ Company login failed!"
  echo "$COMPANY_LOGIN" | jq '.'
  exit 1
fi
echo "✅ Company login success"
echo "Token: $COMPANY_TOKEN"

# Step 2: Get Current User/Company Info
echo ""
echo "2️⃣  Getting company profile..."
COMPANY_PROFILE=$(curl -s -X GET "$API/auth/me" \
  -H "Authorization: Bearer $COMPANY_TOKEN")

echo "$COMPANY_PROFILE" | jq '.'

COMPANY_ID=$(echo $COMPANY_LOGIN | jq -r '.data.user.id // empty')
STATUS=$(echo $COMPANY_LOGIN | jq -r '.data.user.company_status // empty')

echo ""
echo "📊 Company Info:"
echo "  - Company ID: $COMPANY_ID"
echo "  - Status: $STATUS"
echo "  - Token: ${COMPANY_TOKEN:0:20}..."

# Step 3: Try to create job (should fail if not verified)
echo ""
echo "3️⃣  Testing job creation (should fail if not verified)..."
JOB_TEST=$(curl -s -X POST "$API/jobs" \
  -H "Authorization: Bearer $COMPANY_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "position_title": "Test Position",
    "job_description": "Test",
    "job_requirements": "Test",
    "job_type": "full-time",
    "salary_min": 5000000,
    "salary_max": 10000000,
    "location": "Jakarta",
    "company_id": "'$COMPANY_ID'"
  }')

echo "$JOB_TEST" | jq '.'

# Check if job creation is blocked
if echo $JOB_TEST | jq -e '.error' > /dev/null; then
  echo ""
  echo "⚠️  Job creation is blocked (expected if not verified)"
else
  echo ""
  echo "✅ Job creation allowed (company is verified!)"
fi

echo ""
echo "======================================"
echo "✅ TEST COMPLETE!"
echo "======================================"
echo ""
echo "📝 Notes:"
echo "  - Company login successful ✓"
echo "  - Current status: $STATUS"
echo "  - Next: Use admin API to verify company"
echo ""
echo "🔐 For admin verification, need admin credentials"
