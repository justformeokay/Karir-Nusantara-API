# 🔐 Admin Company Verification - Complete Testing Guide

## Overview

API endpoint untuk verifikasi company account sudah **SIAP DIGUNAKAN**:
- **Route**: `POST /api/v1/admin/companies/{id}/verify`
- **Status**: ✅ Implemented
- **Location**: [internal/modules/admin/handler.go](../../internal/modules/admin/handler.go#L151)
- **Protection**: Requires Admin authentication + admin role

---

## Quick Start

### Prerequisites
1. Backend API running: `http://localhost:8081`
2. Admin account credentials
3. Company exists: `info@karyadeveloperindonesia.com`

### One-Command Test

```bash
# 1. Get admin token
ADMIN_TOKEN=$(curl -s -X POST "http://localhost:8081/api/v1/admin/auth/login" \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@karircorp.com","password":"admin123456"}' | jq -r '.data.token')

# 2. Get company ID
COMPANY_ID=$(curl -s "http://localhost:8081/api/v1/admin/companies" \
  -H "Authorization: Bearer $ADMIN_TOKEN" | \
  jq -r '.data[] | select(.email=="info@karyadeveloperindonesia.com") | .id')

# 3. Verify company
curl -X POST "http://localhost:8081/api/v1/admin/companies/$COMPANY_ID/verify" \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"action":"approve","notes":"Test verification"}'
```

---

## Step-by-Step Testing

### Step 1: Admin Login

**Endpoint**: `POST /api/v1/admin/auth/login`

```bash
curl -X POST "http://localhost:8081/api/v1/admin/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@karircorp.com",
    "password": "admin123456"
  }'
```

**Response**:
```json
{
  "message": "Login berhasil",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIs...",
    "admin": {
      "id": 1,
      "email": "admin@karircorp.com",
      "full_name": "Admin Karir",
      "role": "admin"
    }
  }
}
```

**Save token for next steps**:
```bash
ADMIN_TOKEN="your_token_here"
```

---

### Step 2: Get Companies List

**Endpoint**: `GET /api/v1/admin/companies`

```bash
curl -X GET "http://localhost:8081/api/v1/admin/companies?page=1&page_size=50" \
  -H "Authorization: Bearer $ADMIN_TOKEN"
```

**Response** (sample):
```json
{
  "message": "Data berhasil diambil",
  "data": [
    {
      "id": 5,
      "user_id": 23,
      "email": "info@karyadeveloperindonesia.com",
      "company_name": "Karya Developer Indonesia",
      "company_status": "pending",
      "is_verified": false,
      "documents_verified_at": null,
      "ktp_founder_url": "/docs/companies/5/ktp_1234567890.pdf",
      "akta_pendirian_url": "/docs/companies/5/akta_1234567890.pdf",
      "npwp_url": "/docs/companies/5/npwp_1234567890.pdf",
      "nib_url": "/docs/companies/5/nib_1234567890.pdf",
      "created_at": "2024-01-15T10:30:00Z"
    },
    // ... more companies
  ],
  "pagination": {
    "page": 1,
    "page_size": 50,
    "total": 10
  }
}
```

**Extract company ID**:
```bash
COMPANY_ID=$(curl -s "http://localhost:8081/api/v1/admin/companies" \
  -H "Authorization: Bearer $ADMIN_TOKEN" | \
  jq -r '.data[] | select(.email=="info@karyadeveloperindonesia.com") | .id')

echo "Company ID: $COMPANY_ID"
```

---

### Step 3: View Company Details Before Verification

**Endpoint**: `GET /api/v1/admin/companies/{id}`

```bash
curl -X GET "http://localhost:8081/api/v1/admin/companies/$COMPANY_ID" \
  -H "Authorization: Bearer $ADMIN_TOKEN" | jq '.data'
```

**Expected Output**:
```
company_status: "pending"
is_verified: false
documents_verified_at: null
documents_verified_by: null
```

---

### Step 4: Verify Company (APPROVE)

**Endpoint**: `POST /api/v1/admin/companies/{id}/verify`

```bash
curl -X POST "http://localhost:8081/api/v1/admin/companies/$COMPANY_ID/verify" \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "action": "approve",
    "notes": "Dokumen dan profil sudah lengkap sesuai ketentuan. Approved untuk testing."
  }'
```

**Request Body**:
```json
{
  "action": "approve",  // "approve" atau "reject"
  "notes": "Optional notes about verification"
}
```

**Response** (Success):
```json
{
  "message": "Perusahaan berhasil diverifikasi",
  "data": null
}
```

**Response** (Error - Company Not Found):
```json
{
  "code": "NOT_FOUND",
  "message": "Perusahaan tidak ditemukan"
}
```

**Response** (Error - Invalid Action):
```json
{
  "code": "INVALID_ACTION",
  "message": "Action harus 'approve' atau 'reject'"
}
```

---

### Step 5: Verify Company Status After Verification

**Endpoint**: `GET /api/v1/admin/companies/{id}`

```bash
curl -X GET "http://localhost:8081/api/v1/admin/companies/$COMPANY_ID" \
  -H "Authorization: Bearer $ADMIN_TOKEN" | jq '.data | {company_name, company_status, is_verified, documents_verified_at}'
```

**Expected Output**:
```json
{
  "company_name": "Karya Developer Indonesia",
  "company_status": "verified",
  "is_verified": true,
  "documents_verified_at": "2024-01-20T14:30:45Z"
}
```

---

### Step 6: Company Login After Verification

**Endpoint**: `POST /api/v1/auth/login`

```bash
curl -X POST "http://localhost:8081/api/v1/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "info@karyadeveloperindonesia.com",
    "password": "company123456"
  }'
```

**Response**:
```json
{
  "message": "Login berhasil",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIs...",
    "user": {
      "id": 23,
      "email": "info@karyadeveloperindonesia.com",
      "role": "company",
      "company": {
        "id": 5,
        "company_name": "Karya Developer Indonesia",
        "company_status": "verified",
        "is_verified": true
      }
    }
  }
}
```

**Save company token**:
```bash
COMPANY_TOKEN="your_company_token_here"
```

---

### Step 7: Get Current User Profile

**Endpoint**: `GET /api/v1/auth/me`

```bash
curl -X GET "http://localhost:8081/api/v1/auth/me" \
  -H "Authorization: Bearer $COMPANY_TOKEN" | jq '.data'
```

**Expected Output** (Status should be VERIFIED):
```json
{
  "id": 23,
  "email": "info@karyadeveloperindonesia.com",
  "full_name": "Company Admin",
  "role": "company",
  "company": {
    "id": 5,
    "company_name": "Karya Developer Indonesia",
    "company_status": "verified",
    "is_verified": true,
    "ktp_founder_url": "/docs/companies/5/ktp_1234567890.pdf",
    "akta_pendirian_url": "/docs/companies/5/akta_1234567890.pdf",
    "npwp_url": "/docs/companies/5/npwp_1234567890.pdf",
    "nib_url": "/docs/companies/5/nib_1234567890.pdf"
  }
}
```

---

### Step 8: Test Job Creation (Now Should Work)

**Endpoint**: `POST /api/v1/jobs`

```bash
curl -X POST "http://localhost:8081/api/v1/jobs" \
  -H "Authorization: Bearer $COMPANY_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "position_title": "Senior Backend Developer",
    "job_description": "Kami mencari Senior Backend Developer berpengalaman dengan track record yang proven...",
    "job_requirements": "Minimum 5 tahun pengalaman di backend development, familiar dengan Go/Node.js, SQL database...",
    "job_type": "full-time",
    "salary_min": 15000000,
    "salary_max": 25000000,
    "location": "Jakarta",
    "company_id": 5
  }'
```

**Expected Response** (Success):
```json
{
  "message": "Lowongan kerja berhasil dibuat",
  "data": {
    "id": 123,
    "position_title": "Senior Backend Developer",
    "company_id": 5,
    "status": "active",
    "created_at": "2024-01-20T14:35:00Z"
  }
}
```

---

## Frontend Testing Checklist

After verification, login to frontend and verify:

### Dashboard Page
- [ ] "Buat Lowongan" button is **ENABLED** (not grayed out)
- [ ] Hover tooltip shows "Siap membuat lowongan" (or similar green message)
- [ ] Button color is normal (not disabled state)

### Job Form Page (/jobs/new)
- [ ] No blocking modal appears
- [ ] Green success alert shows: "Siap Membuat Lowongan"
- [ ] Form is fully enabled (not grayed out)
- [ ] Can fill all job fields
- [ ] Can click "Publish" button

### Create Job Posting
- [ ] Fill form with job details
- [ ] Click "Publish"
- [ ] Job appears in "Lowongan Aktif" section on Dashboard
- [ ] Job is visible in public job listings

### Additional Checks
- [ ] Profile page shows "Status: Verified"
- [ ] Document URLs are correctly displayed
- [ ] All previous features still work

---

## Error Troubleshooting

| Error | Cause | Solution |
|-------|-------|----------|
| `INVALID_CREDENTIALS` | Admin login failed | Verify admin credentials in database |
| `NOT_FOUND` | Company doesn't exist | Check company email exists exactly |
| `INVALID_ACTION` | Action not "approve"/"reject" | Check spelling in request body |
| `Unauthorized` | Token invalid/expired | Re-login to get fresh token |
| Company still shows "pending" | Verify failed silently | Check API response for errors |
| Job creation still blocked | Frontend cached old data | Clear localStorage/refresh browser |

---

## Implementation Details

### Backend Files Modified/Created

**[internal/modules/admin/handler.go](../../internal/modules/admin/handler.go#L151)**
- `VerifyCompany()` - HTTP handler (lines 151-189)
- Validates input, calls service, returns response

**[internal/modules/admin/service.go](../../internal/modules/admin/service.go#L166)**
- `VerifyCompany()` - Business logic (lines 166-200)
- Updates company status in database
- Logs admin action

**[internal/modules/admin/entity.go](../../internal/modules/admin/entity.go)**
- `CompanyVerificationRequest` - Request DTO
- Fields: `action`, `notes`

**[internal/modules/admin/routes.go](../../internal/modules/admin/routes.go#L53)**
- Route registration (line 53)
- POST `/api/v1/admin/companies/{id}/verify`
- Protected with RequireAdmin middleware

### Database Changes

When verification succeeds:

```sql
UPDATE companies 
SET 
  company_status = 'verified',
  is_verified = true,
  documents_verified_at = NOW(),
  documents_verified_by = :admin_id
WHERE id = :company_id
```

---

## Next Steps

1. ✅ Run backend API server
2. ✅ Run test script: `./test_verify_simple.sh`
3. ✅ Verify company via API
4. ✅ Login to frontend with company account
5. ✅ Test all dashboard features
6. ✅ Create job postings
7. ✅ View applicants
8. ✅ Test other company features

---

## Notes

- All verification is **permanent** - no undo (except by direct database modification)
- Consider adding soft-delete or audit trail in future
- Current implementation overwrites previous verification status
- Admin credentials needed - ensure only authorized admins have access
- Token expiration: Typically 24 hours (check JWT settings)

