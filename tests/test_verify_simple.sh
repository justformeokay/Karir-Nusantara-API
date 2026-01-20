#!/bin/bash

# Company Verification Test - Simple Version
# Usage: ./test_verify_simple.sh

API="http://localhost:8081/api/v1"

echo "======================================"
echo "ADMIN VERIFICATION TEST"
echo "======================================"

# Store tokens
echo ""
echo "1️⃣  Admin Login..."
ADMIN_RESPONSE=$(curl -s -X POST "$API/admin/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@karircorp.com",
    "password": "admin123456"
  }')

ADMIN_TOKEN=$(echo $ADMIN_RESPONSE | jq -r '.data.token')
if [ "$ADMIN_TOKEN" == "null" ] || [ -z "$ADMIN_TOKEN" ]; then
  echo "❌ Admin login failed!"
  echo "$ADMIN_RESPONSE" | jq '.'
  exit 1
fi
echo "✅ Admin login success"

# Get all companies
echo ""
echo "2️⃣  Fetching companies list..."
COMPANIES=$(curl -s -X GET "$API/admin/companies" \
  -H "Authorization: Bearer $ADMIN_TOKEN")

COMPANY_ID=$(echo $COMPANIES | jq -r '.data[] | select(.email=="info@karyadeveloperindonesia.com") | .id' | head -1)
if [ -z "$COMPANY_ID" ] || [ "$COMPANY_ID" == "null" ]; then
  echo "❌ Company not found!"
  echo "$COMPANIES" | jq '.data[0:3]'
  exit 1
fi
echo "✅ Found company ID: $COMPANY_ID"

# Get company details before
echo ""
echo "3️⃣  Company details BEFORE verification:"
BEFORE=$(curl -s -X GET "$API/admin/companies/$COMPANY_ID" \
  -H "Authorization: Bearer $ADMIN_TOKEN")
echo "$BEFORE" | jq '.data | {id, company_name, email, company_status, is_verified}'

# Verify company
echo ""
echo "4️⃣  Verifying company..."
VERIFY=$(curl -s -X POST "$API/admin/companies/$COMPANY_ID/verify" \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "action": "approve",
    "notes": "Dokumen dan profil lengkap. Verifikasi untuk testing."
  }')

echo "$VERIFY" | jq '.'

# Get company details after
echo ""
echo "5️⃣  Company details AFTER verification:"
AFTER=$(curl -s -X GET "$API/admin/companies/$COMPANY_ID" \
  -H "Authorization: Bearer $ADMIN_TOKEN")
echo "$AFTER" | jq '.data | {id, company_name, email, company_status, is_verified, documents_verified_at}'

echo ""
echo "======================================"
echo "✅ TEST COMPLETE!"
echo "======================================"
echo ""
echo "📋 Summary:"
echo "  - Company Email: info@karyadeveloperindonesia.com"
echo "  - Company ID: $COMPANY_ID"
echo "  - Status should now be: verified ✓"
echo ""
echo "🧪 Next: Login dengan company account di frontend untuk test"
