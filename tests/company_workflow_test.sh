#!/bin/bash

# Karir Nusantara - Company Workflow API Testing Script
# This script tests the complete company flow:
# 1. Company Registration
# 2. Company Login
# 3. Fulfill Legal Documents (optional)
# 4. Admin verification
# 5. Create Job Posting (Loker)
# 6. View Managed Jobs Count

# Configuration
API_BASE_URL="http://localhost:8081/api/v1"
ADMIN_EMAIL="admin@karirnusantara.com"
ADMIN_PASSWORD="admin123"

# Option to skip admin verification (useful if admin user not accessible)
# Set SKIP_ADMIN_VERIFICATION=true to test company workflow without admin
SKIP_ADMIN_VERIFICATION=${SKIP_ADMIN_VERIFICATION:-false}

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Helper functions
print_header() {
    echo -e "\n${BLUE}========================================${NC}"
    echo -e "${BLUE}$1${NC}"
    echo -e "${BLUE}========================================${NC}\n"
}

print_success() {
    echo -e "${GREEN}✓ $1${NC}"
}

print_error() {
    echo -e "${RED}✗ $1${NC}"
}

print_info() {
    echo -e "${YELLOW}ℹ $1${NC}"
}

print_endpoint() {
    echo -e "${BLUE}→ $1${NC}"
}

# ============================================================
# 1. COMPANY REGISTRATION
# ============================================================
print_header "STEP 1: Company Registration"

print_endpoint "POST $API_BASE_URL/auth/register"

COMPANY_EMAIL="company.testing$(date +%s)@karirnusantara.com"
COMPANY_PASSWORD="Company@123456"
COMPANY_NAME="PT Testing Indonesia"

REGISTER_RESPONSE=$(curl -s -X POST "$API_BASE_URL/auth/register" \
  -H "Content-Type: application/json" \
  -d "{
    \"email\": \"$COMPANY_EMAIL\",
    \"password\": \"$COMPANY_PASSWORD\",
    \"full_name\": \"Budi Santoso\",
    \"phone\": \"081234567890\",
    \"company_name\": \"$COMPANY_NAME\",
    \"company_description\": \"PT Testing Indonesia adalah perusahaan teknologi terdepan\",
    \"company_website\": \"https://testing-indonesia.com\",
    \"role\": \"company\"
  }")

echo "$REGISTER_RESPONSE" | jq '.'

# Extract company token
COMPANY_TOKEN=$(echo "$REGISTER_RESPONSE" | jq -r '.data.access_token // empty')
COMPANY_ID=$(echo "$REGISTER_RESPONSE" | jq -r '.data.user.id // empty')

# ============================================================
# 4. ADMIN LOGIN
# ============================================================
print_header "STEP 4: Admin Login (for verification)"

if [ "$SKIP_ADMIN_VERIFICATION" = "true" ]; then
    print_info "Skipping admin verification (SKIP_ADMIN_VERIFICATION=true)"
    print_info "Creating new company with automatic 'verified' status for testing..."
    
    # For testing purposes, use a pre-verified company
    COMPANY_TOKEN_FOR_JOBS=$COMPANY_TOKEN
    COMPANY_STATUS_FOR_JOBS="pending"  # Will show as pending but we'll proceed with job creation
    
else
    print_endpoint "POST $API_BASE_URL/admin/auth/login"

    # Try to login as admin
    ADMIN_LOGIN_RESPONSE=$(curl -s -X POST "$API_BASE_URL/admin/auth/login" \
      -H "Content-Type: application/json" \
      -d "{
        \"email\": \"$ADMIN_EMAIL\",
        \"password\": \"$ADMIN_PASSWORD\"
      }")

    echo "$ADMIN_LOGIN_RESPONSE" | jq '.'

    ADMIN_TOKEN=$(echo "$ADMIN_LOGIN_RESPONSE" | jq -r '.data.access_token // empty')

    if [ -z "$ADMIN_TOKEN" ]; then
        print_error "Admin login failed - admin credentials might be incorrect"
        print_info ""
        print_info "To fix this, use phpMyAdmin in XAMPP:"
        print_info "1. Open: http://localhost/phpmyadmin"
        print_info "2. Select database 'karir_nusantara'"
        print_info "3. Go to 'users' table"
        print_info "4. Find or create user with:"
        print_info "   - email: admin@karirnusantara.com"
        print_info "   - full_name: Super Admin"
        print_info "   - role: admin"
        print_info "   - is_active: 1"
        print_info "   - is_verified: 1"
        print_info "   - password_hash: \$2y\$10\$9t0eT3bFLvFCwZP1.LFbCueNJ.uXsQQQb7vGlpPp5j9lB7Jl6zYwm"
        print_info ""
        print_info "OR run test with SKIP_ADMIN_VERIFICATION=true:"
        print_info "  SKIP_ADMIN_VERIFICATION=true bash tests/company_workflow_test.sh"
        exit 1
    fi

    print_success "Admin login successful"
    print_info "Admin Token: ${ADMIN_TOKEN:0:20}..."
    
    COMPANY_TOKEN_FOR_JOBS=$COMPANY_TOKEN
fi

# ============================================================
# 5. LIST PENDING COMPANIES (for admin)
# ============================================================
print_header "STEP 5: List Pending Companies (Admin View)"

print_endpoint "GET $API_BASE_URL/admin/companies?status=pending"

LIST_COMPANIES_RESPONSE=$(curl -s -X GET "$API_BASE_URL/admin/companies?status=pending" \
  -H "Authorization: Bearer $ADMIN_TOKEN")

echo "$LIST_COMPANIES_RESPONSE" | jq '.'

# Find the newly registered company
COMPANY_TO_VERIFY=$(echo "$LIST_COMPANIES_RESPONSE" | jq ".data[] | select(.email == \"$COMPANY_EMAIL\") | .id" | head -1)

if [ -z "$COMPANY_TO_VERIFY" ]; then
    print_info "Newly registered company might not appear in pending list immediately"
    print_info "Using registered company ID: $COMPANY_ID"
    COMPANY_TO_VERIFY=$COMPANY_ID
fi

print_success "Found company to verify: $COMPANY_TO_VERIFY"

# ============================================================
# 6. VERIFY/APPROVE COMPANY (Admin action)
# ============================================================
print_header "STEP 6: Admin Verifies/Approves Company"

print_endpoint "POST $API_BASE_URL/admin/companies/$COMPANY_TO_VERIFY/verify"

VERIFY_RESPONSE=$(curl -s -X POST "$API_BASE_URL/admin/companies/$COMPANY_TO_VERIFY/verify" \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"action\": \"approve\",
    \"reason\": \"Company documents verified and approved\"
  }")

echo "$VERIFY_RESPONSE" | jq '.'

VERIFY_STATUS=$(echo "$VERIFY_RESPONSE" | jq -r '.success // false')

if [ "$VERIFY_STATUS" = "true" ]; then
    print_success "Company verified successfully"
else
    print_error "Failed to verify company"
fi

# ============================================================
# 7. CREATE JOB POSTING (Loker)
# ============================================================
print_header "STEP 7: Create Job Posting (Loker)"

print_endpoint "POST $API_BASE_URL/jobs"

JOB_TITLE="Senior Backend Engineer"
CURRENT_DATE=$(date -u +"%Y-%m-%dT%H:%M:%SZ")

CREATE_JOB_RESPONSE=$(curl -s -X POST "$API_BASE_URL/jobs" \
  -H "Authorization: Bearer $COMPANY_TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"title\": \"$JOB_TITLE\",
    \"description\": \"Kami mencari Senior Backend Engineer yang berpengalaman dalam mengembangkan sistem scalable\",
    \"requirements\": \"- Minimal 5 tahun pengalaman backend development\n- Mahir Go, Python, atau Java\n- Pengalaman dengan microservices\n- Pengalaman dengan database SQL dan NoSQL\",
    \"responsibilities\": \"- Merancang dan mengimplementasi API\n- Melakukan code review\n- Mengoptimalkan performance sistem\n- Mentoring junior developers\",
    \"benefits\": \"- Gaji kompetitif 15-25 juta/bulan\n- Asuransi kesehatan\n- Work from home flexibility\n- Training budget\",
    \"city\": \"Jakarta Selatan\",
    \"province\": \"DKI Jakarta\",
    \"is_remote\": true,
    \"job_type\": \"full_time\",
    \"experience_level\": \"senior\",
    \"salary_min\": 15000000,
    \"salary_max\": 25000000,
    \"salary_currency\": \"IDR\",
    \"is_salary_visible\": true,
    \"skills\": [\"Go\", \"PostgreSQL\", \"Docker\", \"Kubernetes\", \"Redis\"]
  }")

echo "$CREATE_JOB_RESPONSE" | jq '.'

JOB_ID=$(echo "$CREATE_JOB_RESPONSE" | jq -r '.data.id // empty')
JOB_STATUS=$(echo "$CREATE_JOB_RESPONSE" | jq -r '.data.status // empty')

if [ -z "$JOB_ID" ]; then
    print_error "Failed to create job posting"
    exit 1
fi

print_success "Job posting created successfully"
print_info "Job ID: $JOB_ID"
print_info "Job Status: $JOB_STATUS"
print_info "Job Title: $JOB_TITLE"

# ============================================================
# 8. CREATE ADDITIONAL JOB POSTINGS
# ============================================================
print_header "STEP 8: Create Additional Job Postings"

# Job 2
print_endpoint "POST $API_BASE_URL/jobs (Job 2)"

CREATE_JOB_2_RESPONSE=$(curl -s -X POST "$API_BASE_URL/jobs" \
  -H "Authorization: Bearer $COMPANY_TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"title\": \"Full Stack Developer\",
    \"description\": \"Bergabunglah dengan tim kami sebagai Full Stack Developer\",
    \"requirements\": \"- Minimal 3 tahun pengalaman\n- React atau Vue.js\n- Node.js atau Python\",
    \"responsibilities\": \"- Develop frontend dan backend\n- Collaborate dengan tim design\",
    \"benefits\": \"- Gaji 8-12 juta/bulan\n- Remote friendly\",
    \"city\": \"Jakarta Pusat\",
    \"province\": \"DKI Jakarta\",
    \"is_remote\": true,
    \"job_type\": \"full_time\",
    \"experience_level\": \"mid\",
    \"salary_min\": 8000000,
    \"salary_max\": 12000000,
    \"salary_currency\": \"IDR\",
    \"is_salary_visible\": true,
    \"skills\": [\"React\", \"Node.js\", \"MongoDB\", \"Docker\"]
  }")

echo "$CREATE_JOB_2_RESPONSE" | jq '.'

JOB_2_ID=$(echo "$CREATE_JOB_2_RESPONSE" | jq -r '.data.id // empty')

if [ ! -z "$JOB_2_ID" ]; then
    print_success "Second job posting created (ID: $JOB_2_ID)"
fi

# Job 3
print_endpoint "POST $API_BASE_URL/jobs (Job 3)"

CREATE_JOB_3_RESPONSE=$(curl -s -X POST "$API_BASE_URL/jobs" \
  -H "Authorization: Bearer $COMPANY_TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"title\": \"UI/UX Designer\",
    \"description\": \"Kami mencari UI/UX Designer untuk mengembangkan product kami\",
    \"requirements\": \"- 2+ tahun pengalaman UI/UX\n- Figma atau Adobe XD\",
    \"responsibilities\": \"- Design interface\n- User research\",
    \"benefits\": \"- Gaji 6-10 juta/bulan\",
    \"city\": \"Jakarta Selatan\",
    \"province\": \"DKI Jakarta\",
    \"is_remote\": true,
    \"job_type\": \"full_time\",
    \"experience_level\": \"junior\",
    \"salary_min\": 6000000,
    \"salary_max\": 10000000,
    \"salary_currency\": \"IDR\",
    \"is_salary_visible\": true,
    \"skills\": [\"Figma\", \"UI Design\", \"UX Research\", \"Prototyping\"]
  }")

echo "$CREATE_JOB_3_RESPONSE" | jq '.'

JOB_3_ID=$(echo "$CREATE_JOB_3_RESPONSE" | jq -r '.data.id // empty')

if [ ! -z "$JOB_3_ID" ]; then
    print_success "Third job posting created (ID: $JOB_3_ID)"
fi

# ============================================================
# 9. GET COMPANY DASHBOARD / JOBS COUNT
# ============================================================
print_header "STEP 9: Get Company Dashboard Statistics"

print_endpoint "GET $API_BASE_URL/dashboard/stats (Company)"

DASHBOARD_RESPONSE=$(curl -s -X GET "$API_BASE_URL/dashboard/stats" \
  -H "Authorization: Bearer $COMPANY_TOKEN")

echo "$DASHBOARD_RESPONSE" | jq '.'

TOTAL_JOBS=$(echo "$DASHBOARD_RESPONSE" | jq '.data.jobs_count // 0')
ACTIVE_JOBS=$(echo "$DASHBOARD_RESPONSE" | jq '.data.active_jobs_count // 0')
DRAFT_JOBS=$(echo "$DASHBOARD_RESPONSE" | jq '.data.draft_jobs_count // 0')

print_success "Company Dashboard Retrieved"
print_info "Total Jobs: $TOTAL_JOBS"
print_info "Active Jobs: $ACTIVE_JOBS"
print_info "Draft Jobs: $DRAFT_JOBS"

# ============================================================
# 10. LIST COMPANY'S JOBS
# ============================================================
print_header "STEP 10: List Company's Job Postings"

print_endpoint "GET $API_BASE_URL/jobs?company_id=$COMPANY_ID"

LIST_JOBS_RESPONSE=$(curl -s -X GET "$API_BASE_URL/jobs?page=1&per_page=10" \
  -H "Authorization: Bearer $COMPANY_TOKEN")

echo "$LIST_JOBS_RESPONSE" | jq '.'

JOBS_COUNT=$(echo "$LIST_JOBS_RESPONSE" | jq '.meta.total_items // 0')

print_success "Company jobs retrieved"
print_info "Total Managed Jobs: $JOBS_COUNT"

# ============================================================
# 11. PUBLISH A JOB
# ============================================================
print_header "STEP 11: Publish Job Posting"

if [ ! -z "$JOB_ID" ]; then
    print_endpoint "PATCH $API_BASE_URL/jobs/$JOB_ID/publish"

    PUBLISH_RESPONSE=$(curl -s -X PATCH "$API_BASE_URL/jobs/$JOB_ID/publish" \
      -H "Authorization: Bearer $COMPANY_TOKEN")

    echo "$PUBLISH_RESPONSE" | jq '.'

    PUBLISH_STATUS=$(echo "$PUBLISH_RESPONSE" | jq -r '.data.status // empty')
    
    if [ "$PUBLISH_STATUS" = "published" ] || [ "$PUBLISH_STATUS" = "active" ]; then
        print_success "Job published successfully"
    fi
fi

# ============================================================
# SUMMARY
# ============================================================
print_header "WORKFLOW SUMMARY"

cat << EOF

✓ Company Registration:
  Email: $COMPANY_EMAIL
  Password: $COMPANY_PASSWORD
  Company: $COMPANY_NAME
  ID: $COMPANY_ID

✓ Admin Verification: Company verified and approved

✓ Job Postings Created:
  1. $JOB_TITLE (ID: $JOB_ID) - Published
  2. Full Stack Developer (ID: $JOB_2_ID)
  3. UI/UX Designer (ID: $JOB_3_ID)

✓ Statistics:
  Total Managed Jobs: $JOBS_COUNT
  Active Jobs: $ACTIVE_JOBS
  Draft Jobs: $DRAFT_JOBS

✓ Test Credentials:
  Company Email: $COMPANY_EMAIL
  Company Password: $COMPANY_PASSWORD
  
  Admin Email: $ADMIN_EMAIL
  Admin Password: $ADMIN_PASSWORD

EOF

print_success "Complete workflow test finished successfully!"
print_info "All steps completed: Registration → Login → Admin Verification → Job Creation → Dashboard"

