#!/bin/bash

# Complete Admin Verification Flow
# Admin: admin@karirnusantara.com / admin123
# Company: info@karyadeveloperindonesia.com

API="http://localhost:8081/api/v1"

echo "======================================"
echo "ADMIN VERIFICATION TEST"
echo "======================================"

# Step 1: Admin Login
echo ""
echo "1️⃣  Admin Login (admin@karirnusantara.com)..."
ADMIN_LOGIN=$(curl -s -X POST "$API/admin/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@karirnusantara.com",
    "password": "admin123"
  }')

ADMIN_TOKEN=$(echo $ADMIN_LOGIN | jq -r '.data.access_token // .data.token // empty')
if [ -z "$ADMIN_TOKEN" ]; then
  echo "❌ Admin login failed!"
  echo "$ADMIN_LOGIN" | jq '.'
  exit 1
fi
echo "✅ Admin login success"
echo "Admin Token: ${ADMIN_TOKEN:0:30}..."

# Step 2: Get All Companies
echo ""
echo "2️⃣  Fetching all companies..."
COMPANIES=$(curl -s -X GET "$API/admin/companies" \
  -H "Authorization: Bearer $ADMIN_TOKEN")

# Debug: show first company
echo "$COMPANIES" | jq '.data[0:2]'

# Step 3: Find company by email
echo ""
echo "3️⃣  Finding company: info@karyadeveloperindonesia.com"
COMPANY_ID=$(echo $COMPANIES | jq -r '.data[] | select(.email=="info@karyadeveloperindonesia.com") | .id' | head -1)
if [ -z "$COMPANY_ID" ]; then
  echo "❌ Company not found!"
  echo "Trying with company_name search..."
  COMPANY_ID=$(echo $COMPANIES | jq -r '.data[] | select(.company_name | contains("Karya Developer")) | .id' | head -1)
fi

if [ -z "$COMPANY_ID" ]; then
  echo "❌ Company still not found! Available companies:"
  echo "$COMPANIES" | jq '.data[] | {id, company_name, email, company_status}'
  exit 1
fi
echo "✅ Found company ID: $COMPANY_ID"

# Step 4: Get Company Details Before Verification
echo ""
echo "4️⃣  Company details BEFORE verification:"
BEFORE=$(curl -s -X GET "$API/admin/companies/$COMPANY_ID" \
  -H "Authorization: Bearer $ADMIN_TOKEN")

echo "$BEFORE" | jq '.data | {id, company_name, email, company_status, is_verified, documents_verified_at}'

# Step 5: Verify Company
echo ""
echo "5️⃣  🚀 Verifying company (APPROVE)..."
VERIFY=$(curl -s -X POST "$API/admin/companies/$COMPANY_ID/verify" \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "action": "approve",
    "reason": "Dokumen dan profil lengkap. Verifikasi untuk testing dashboard features."
  }')

echo "$VERIFY" | jq '.'

# Step 6: Get Company Details After Verification
echo ""
echo "6️⃣  Company details AFTER verification:"
AFTER=$(curl -s -X GET "$API/admin/companies/$COMPANY_ID" \
  -H "Authorization: Bearer $ADMIN_TOKEN")

echo "$AFTER" | jq '.data | {id, company_name, email, company_status, is_verified, documents_verified_at}'

# Step 7: Company Login Again to Verify Status
echo ""
echo "7️⃣  Company login to verify new status..."
COMPANY_LOGIN=$(curl -s -X POST "$API/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "info@karyadeveloperindonesia.com",
    "password": "Justformeokay23@"
  }')

COMPANY_TOKEN=$(echo $COMPANY_LOGIN | jq -r '.data.access_token // empty')
if [ -z "$COMPANY_TOKEN" ]; then
  echo "❌ Company login failed!"
  echo "$COMPANY_LOGIN" | jq '.'
  exit 1
fi
echo "✅ Company login success"

# Step 8: Get Company Profile
echo ""
echo "8️⃣  Verified company profile:"
PROFILE=$(curl -s -X GET "$API/auth/me" \
  -H "Authorization: Bearer $COMPANY_TOKEN")

echo "$PROFILE" | jq '.data | {id, email, company_name, company_status, is_verified}'

STATUS=$(echo $PROFILE | jq -r '.data.company_status // empty')

echo ""
echo "======================================"
echo "✅ VERIFICATION COMPLETE!"
echo "======================================"
echo ""
echo "📊 Final Status:"
echo "  - Email: info@karyadeveloperindonesia.com"
echo "  - Company ID: $COMPANY_ID"
echo "  - Status: $STATUS"
echo ""

if [ "$STATUS" == "verified" ]; then
  echo "✅ SUCCESS! Company is now VERIFIED"
  echo ""
  echo "🎉 You can now:"
  echo "  - Create job postings"
  echo "  - Access all dashboard features"
  echo "  - Test company functions"
else
  echo "⚠️  Status is still: $STATUS"
  echo "Please check admin response above"
fi

echo ""
echo "📝 Frontend Testing:"
echo "  1. Refresh browser"
echo "  2. Login with info@karyadeveloperindonesia.com"
echo "  3. Check Dashboard - button should be ENABLED"
echo "  4. Try creating a job posting"
