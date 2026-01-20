# Company Verification Test - Curl Commands

## Prerequisites
- API running on: http://localhost:8081
- Admin account exists with credentials
- Company "Karya Developer Indonesia" exists with email: info@karyadeveloperindonesia.com

## API Endpoints to Test

### 1. Admin Login
```bash
curl -X POST "http://localhost:8081/api/v1/admin/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@karircorp.com",
    "password": "admin123456"
  }'
```
**Response:** Returns `data.token` (JWT token with admin role)

---

### 2. Get All Companies (Admin View)
```bash
curl -X GET "http://localhost:8081/api/v1/admin/companies?page=1&page_size=50" \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN"
```
**Response:** List of companies including pending ones

---

### 3. Find Company ID by Email
```bash
# Get company ID for info@karyadeveloperindonesia.com
curl -s "http://localhost:8081/api/v1/admin/companies?page=1&page_size=50" \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN" | jq '.data[] | select(.email=="info@karyadeveloperindonesia.com") | .id'
```
**Response:** Company ID (UUID or number depending on implementation)

---

### 4. Get Company Details Before Verification
```bash
curl -X GET "http://localhost:8081/api/v1/admin/companies/COMPANY_ID" \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN"
```
**Response:** Company details including status (should be "pending")

---

### 5. Verify Company (APPROVE)
```bash
curl -X POST "http://localhost:8081/api/v1/admin/companies/COMPANY_ID/verify" \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "action": "approve",
    "notes": "Dokumen dan profil sudah lengkap dan sesuai dengan ketentuan."
  }'
```
**Response:** Success message indicating company is now verified

---

### 6. Get Company Details After Verification
```bash
curl -X GET "http://localhost:8081/api/v1/admin/companies/COMPANY_ID" \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN"
```
**Response:** Company status should be "verified"

---

### 7. Company Login with Verified Account
```bash
curl -X POST "http://localhost:8081/api/v1/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "info@karyadeveloperindonesia.com",
    "password": "COMPANY_PASSWORD"
  }'
```
**Response:** Returns JWT token for company user

---

### 8. Get Current User Profile (Verify Company Status)
```bash
curl -X GET "http://localhost:8081/api/v1/auth/me" \
  -H "Authorization: Bearer COMPANY_TOKEN"
```
**Response:** Company profile should show status "verified" and is_verified = true

---

### 9. Test Job Creation (Now Should Work)
```bash
curl -X POST "http://localhost:8081/api/v1/jobs" \
  -H "Authorization: Bearer COMPANY_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "position_title": "Senior Backend Developer",
    "job_description": "Kami mencari Senior Backend Developer berpengalaman...",
    "job_requirements": "Minimum 5 tahun pengalaman...",
    "job_type": "full-time",
    "salary_min": 15000000,
    "salary_max": 25000000,
    "location": "Jakarta",
    "company_id": "COMPANY_ID"
  }'
```
**Response:** Job posting created successfully

---

## Complete Test Flow

### Step 1: Get Admin Token
```bash
ADMIN_TOKEN=$(curl -s -X POST "http://localhost:8081/api/v1/admin/auth/login" \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@karircorp.com","password":"admin123456"}' | jq -r '.data.token')

echo $ADMIN_TOKEN
```

### Step 2: Get Company ID
```bash
COMPANY_ID=$(curl -s -X GET "http://localhost:8081/api/v1/admin/companies?page=1&page_size=50" \
  -H "Authorization: Bearer $ADMIN_TOKEN" | \
  jq -r '.data[] | select(.email=="info@karyadeveloperindonesia.com") | .id' | head -1)

echo $COMPANY_ID
```

### Step 3: Verify Company
```bash
curl -X POST "http://localhost:8081/api/v1/admin/companies/$COMPANY_ID/verify" \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"action":"approve","notes":"Test verification"}'
```

### Step 4: Company Login
```bash
COMPANY_TOKEN=$(curl -s -X POST "http://localhost:8081/api/v1/auth/login" \
  -H "Content-Type: application/json" \
  -d '{"email":"info@karyadeveloperindonesia.com","password":"company123456"}' | jq -r '.data.token')

echo $COMPANY_TOKEN
```

### Step 5: Verify Status
```bash
curl -X GET "http://localhost:8081/api/v1/auth/me" \
  -H "Authorization: Bearer $COMPANY_TOKEN" | jq '.data | {email, company_status, is_verified}'
```

---

## Expected Results After Verification

### Dashboard Should Now Show:
✅ "Buat Lowongan" button ENABLED (not grayed out)  
✅ Green success alert: "Siap Membuat Lowongan"  
✅ Can access /jobs/new without blocking modal  
✅ Can fill and publish job postings  

### Company Status Shows:
- `company_status`: "verified"
- `is_verified`: true
- `documents_verified_at`: Current timestamp
- `documents_verified_by`: Admin user ID

---

## Troubleshooting

### If Admin Login Fails
- Check admin credentials in database
- Verify admin role exists in system
- Check if /api/v1/admin/auth/login endpoint exists

### If Company Not Found
- Verify company with email info@karyadeveloperindonesia.com exists
- Check company_email format matches exactly
- Try listing all companies to confirm it's there

### If Verify Fails
- Check authorization header has valid admin token
- Verify company_id is correct UUID/format
- Ensure "action" field is "approve" or "reject" (not typos)
- Check response for detailed error message

### If Job Creation Still Blocked
- Refresh frontend to load new company status
- Check browser localStorage is cleared
- Verify useCompanyEligibility hook picks up new status
- Check network tab for latest company data

---

## Notes

- Replace `YOUR_ADMIN_TOKEN` with actual token from Step 1
- Replace `COMPANY_ID` with actual company ID from Step 2
- Replace `COMPANY_PASSWORD` with actual company login password
- All timestamps in responses are UTC
- Documents should already be uploaded before verification (from previous steps)
