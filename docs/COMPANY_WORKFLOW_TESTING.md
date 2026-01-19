# Karir Nusantara - Company Workflow API Testing Guide

Complete testing workflow untuk company dari registrasi hingga mengelola job postings.

## 📋 Daftar Isi

1. [Setup Database](#setup-database)
2. [Company Registration](#company-registration)
3. [Company Login](#company-login)
4. [View Company Information](#view-company-information)
5. [Admin Login](#admin-login)
6. [Admin Verifies Company](#admin-verifies-company)
7. [Create Job Posting (Loker)](#create-job-posting)
8. [View Managed Jobs](#view-managed-jobs)
9. [Publishing & Managing Jobs](#publishing--managing-jobs)

---

## Setup Database

### Prerequisites
- MySQL Server running on `localhost:3306`
- Database: `karir_nusantara`
- API Server: Running on `http://localhost:8081`

### Create Admin User

Jalankan query SQL berikut di MySQL:

```sql
-- Create admin user if not exists
-- Password: admin123
INSERT INTO users (
    email, 
    password_hash, 
    full_name, 
    role, 
    is_active, 
    is_verified, 
    email_verified_at, 
    created_at, 
    updated_at
) VALUES (
    'admin@karirnusantara.com',
    '$2y$10$9t0eT3bFLvFCwZP1.LFbCueNJ.uXsQQQb7vGlpPp5j9lB7Jl6zYwm',
    'Super Admin',
    'admin',
    true,
    true,
    NOW(),
    NOW(),
    NOW()
) ON DUPLICATE KEY UPDATE password_hash = VALUES(password_hash);

-- Create default companies for testing (optional)
INSERT INTO users (
    email, 
    password_hash, 
    full_name, 
    phone,
    company_name,
    role, 
    is_active, 
    is_verified, 
    company_status,
    email_verified_at, 
    created_at, 
    updated_at
) VALUES (
    'test.company@example.com',
    '$2y$10$9t0eT3bFLvFCwZP1.LFbCueNJ.uXsQQQb7vGlpPp5j9lB7Jl6zYwm',
    'Test Company Admin',
    '081234567890',
    'PT Test Company',
    'company',
    true,
    true,
    'verified',
    NOW(),
    NOW(),
    NOW()
) ON DUPLICATE KEY UPDATE password_hash = VALUES(password_hash), company_status = 'verified';
```

---

## Company Registration

### Endpoint
```
POST /api/v1/auth/register
```

### Request

```bash
curl -X POST http://localhost:8081/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "company@example.com",
    "password": "CompanyPass@123",
    "full_name": "Direktur Perusahaan",
    "phone": "081234567890",
    "company_name": "PT Solusi Indonesia",
    "company_description": "Perusahaan teknologi terdepan di Indonesia",
    "company_website": "https://solusi-indonesia.com",
    "role": "company"
  }'
```

### Response (201 Created)

```json
{
  "success": true,
  "message": "Registration successful",
  "data": {
    "user": {
      "id": 9,
      "email": "company@example.com",
      "role": "company",
      "full_name": "Direktur Perusahaan",
      "phone": "081234567890",
      "company_name": "PT Solusi Indonesia",
      "is_active": true,
      "is_verified": false,
      "created_at": "2026-01-19T14:44:42+07:00"
    },
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "uuid-token",
    "expires_in": 900
  }
}
```

**Catatan:**
- Simpan `access_token` untuk request selanjutnya
- Status company awal: `pending` (belum diverifikasi admin)
- Password minimal 8 karakter dengan uppercase, lowercase, dan number

---

## Company Login

### Endpoint
```
POST /api/v1/auth/login
```

### Request

```bash
curl -X POST http://localhost:8081/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "company@example.com",
    "password": "CompanyPass@123"
  }'
```

### Response (200 OK)

```json
{
  "success": true,
  "message": "Login successful",
  "data": {
    "user": {
      "id": 9,
      "email": "company@example.com",
      "role": "company",
      "full_name": "Direktur Perusahaan",
      "phone": "081234567890",
      "company_name": "PT Solusi Indonesia",
      "is_active": true,
      "is_verified": false,
      "created_at": "2026-01-19T14:44:42+07:00"
    },
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "uuid-token",
    "expires_in": 900
  }
}
```

---

## View Company Information

### Endpoint
```
GET /api/v1/auth/me
```

### Request

```bash
COMPANY_TOKEN="your_access_token_here"

curl -X GET http://localhost:8081/api/v1/auth/me \
  -H "Authorization: Bearer $COMPANY_TOKEN"
```

### Response (200 OK)

```json
{
  "success": true,
  "message": "User retrieved",
  "data": {
    "id": 9,
    "email": "company@example.com",
    "role": "company",
    "full_name": "Direktur Perusahaan",
    "phone": "081234567890",
    "company_name": "PT Solusi Indonesia",
    "is_active": true,
    "is_verified": false,
    "created_at": "2026-01-19T14:44:42+07:00"
  }
}
```

---

## Admin Login

### Endpoint
```
POST /api/v1/admin/auth/login
```

### Request

```bash
curl -X POST http://localhost:8081/api/v1/admin/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@karirnusantara.com",
    "password": "admin123"
  }'
```

### Response (200 OK)

```json
{
  "success": true,
  "message": "Login berhasil",
  "data": {
    "admin": {
      "id": 1,
      "email": "admin@karirnusantara.com",
      "full_name": "Super Admin",
      "role": "admin",
      "is_active": true,
      "created_at": "2024-01-01T00:00:00Z"
    },
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "expires_in": 86400
  }
}
```

**Catatan:** Simpan `access_token` untuk admin requests

---

## Admin Verifies Company

### Step 1: List Pending Companies

#### Endpoint
```
GET /api/v1/admin/companies?status=pending
```

#### Request

```bash
ADMIN_TOKEN="your_admin_token_here"

curl -X GET "http://localhost:8081/api/v1/admin/companies?status=pending" \
  -H "Authorization: Bearer $ADMIN_TOKEN"
```

#### Response (200 OK)

```json
{
  "success": true,
  "message": "Daftar perusahaan berhasil diambil",
  "data": [
    {
      "id": 9,
      "email": "company@example.com",
      "full_name": "Direktur Perusahaan",
      "phone": "081234567890",
      "company_name": "PT Solusi Indonesia",
      "company_description": "Perusahaan teknologi terdepan di Indonesia",
      "company_website": "https://solusi-indonesia.com",
      "company_logo_url": null,
      "company_status": "pending",
      "is_active": true,
      "is_verified": false,
      "created_at": "2026-01-19T14:44:42+07:00",
      "jobs_count": 0,
      "active_jobs_count": 0,
      "total_applications": 0
    }
  ],
  "meta": {
    "page": 1,
    "per_page": 10,
    "total_items": 1,
    "total_pages": 1
  }
}
```

### Step 2: Verify/Approve Company

#### Endpoint
```
POST /api/v1/admin/companies/{id}/verify
```

#### Request

```bash
COMPANY_ID=9

curl -X POST "http://localhost:8081/api/v1/admin/companies/$COMPANY_ID/verify" \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "action": "approve",
    "reason": "Dokumen legal telah diverifikasi dan lengkap"
  }'
```

#### Response (200 OK)

```json
{
  "success": true,
  "message": "Perusahaan berhasil diverifikasi",
  "data": {
    "id": 9,
    "company_status": "verified",
    "message": "Company status updated to verified"
  }
}
```

---

## Create Job Posting

### Endpoint
```
POST /api/v1/jobs
```

### Request

```bash
COMPANY_TOKEN="your_company_access_token"

curl -X POST http://localhost:8081/api/v1/jobs \
  -H "Authorization: Bearer $COMPANY_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Senior Backend Engineer",
    "description": "Kami mencari Senior Backend Engineer yang berpengalaman dalam mengembangkan sistem scalable dengan teknologi modern.",
    "requirements": "- Minimal 5 tahun pengalaman backend development\n- Mahir Go, Python, atau Java\n- Pengalaman dengan microservices dan Docker\n- Pengalaman dengan database SQL dan NoSQL\n- Familiar dengan CI/CD dan DevOps practices",
    "responsibilities": "- Merancang dan mengimplementasi RESTful API\n- Melakukan code review dan mentoring\n- Mengoptimalkan performance sistem\n- Berkolaborasi dengan tim frontend dan mobile",
    "benefits": "- Gaji kompetitif 15-25 juta/bulan\n- Asuransi kesehatan dan keluarga\n- Work from home flexibility (2-3 hari per minggu)\n- Training budget untuk development\n- Tunjangan makan dan transport",
    "city": "Jakarta Selatan",
    "province": "DKI Jakarta",
    "is_remote": true,
    "job_type": "full_time",
    "experience_level": "senior",
    "salary_min": 15000000,
    "salary_max": 25000000,
    "salary_currency": "IDR",
    "is_salary_visible": true,
    "skills": ["Go", "PostgreSQL", "Docker", "Kubernetes", "Redis", "API Design"]
  }'
```

### Response (201 Created)

```json
{
  "success": true,
  "message": "Job created successfully",
  "data": {
    "id": 1,
    "company_id": 9,
    "title": "Senior Backend Engineer",
    "slug": "senior-backend-engineer",
    "description": "Kami mencari Senior Backend Engineer...",
    "requirements": "- Minimal 5 tahun pengalaman...",
    "responsibilities": "- Merancang dan mengimplementasi...",
    "benefits": "- Gaji kompetitif...",
    "city": "Jakarta Selatan",
    "province": "DKI Jakarta",
    "is_remote": true,
    "job_type": "full_time",
    "experience_level": "senior",
    "salary_min": 15000000,
    "salary_max": 25000000,
    "salary_currency": "IDR",
    "is_salary_visible": true,
    "skills": ["Go", "PostgreSQL", "Docker", "Kubernetes", "Redis", "API Design"],
    "status": "draft",
    "views_count": 0,
    "applications_count": 0,
    "created_at": "2026-01-19T14:44:42+07:00",
    "updated_at": "2026-01-19T14:44:42+07:00"
  }
}
```

**Catatan:**
- Job awal status: `draft` (belum publish)
- Hanya company yang verified bisa membuat job
- Skills array maksimal 10 skills

---

## View Managed Jobs

### Endpoint
```
GET /api/v1/dashboard/stats
```

### Request

```bash
COMPANY_TOKEN="your_company_access_token"

curl -X GET http://localhost:8081/api/v1/dashboard/stats \
  -H "Authorization: Bearer $COMPANY_TOKEN"
```

### Response (200 OK)

```json
{
  "success": true,
  "message": "Statistics retrieved successfully",
  "data": {
    "jobs_count": 3,
    "active_jobs_count": 2,
    "draft_jobs_count": 1,
    "paused_jobs_count": 0,
    "closed_jobs_count": 0,
    "applications_count": 15,
    "pending_applications_count": 8,
    "accepted_applications_count": 5,
    "rejected_applications_count": 2,
    "total_views": 450,
    "average_views_per_job": 150,
    "last_7_days_applications": 7
  }
}
```

---

## Publishing & Managing Jobs

### Publish a Job

#### Endpoint
```
PATCH /api/v1/jobs/{id}/publish
```

#### Request

```bash
JOB_ID=1

curl -X PATCH "http://localhost:8081/api/v1/jobs/$JOB_ID/publish" \
  -H "Authorization: Bearer $COMPANY_TOKEN"
```

#### Response (200 OK)

```json
{
  "success": true,
  "message": "Job published successfully",
  "data": {
    "id": 1,
    "title": "Senior Backend Engineer",
    "status": "published",
    "published_at": "2026-01-19T14:45:00+07:00"
  }
}
```

### Pause a Job

#### Endpoint
```
PATCH /api/v1/jobs/{id}/pause
```

#### Request

```bash
curl -X PATCH "http://localhost:8081/api/v1/jobs/$JOB_ID/pause" \
  -H "Authorization: Bearer $COMPANY_TOKEN"
```

### Close a Job

#### Endpoint
```
PATCH /api/v1/jobs/{id}/close
```

#### Request

```bash
curl -X PATCH "http://localhost:8081/api/v1/jobs/$JOB_ID/close" \
  -H "Authorization: Bearer $COMPANY_TOKEN"
```

### Reopen a Job

#### Endpoint
```
PATCH /api/v1/jobs/{id}/reopen
```

#### Request

```bash
curl -X PATCH "http://localhost:8081/api/v1/jobs/$JOB_ID/reopen" \
  -H "Authorization: Bearer $COMPANY_TOKEN"
```

---

## Error Responses

### 400 Bad Request
```json
{
  "success": false,
  "error": {
    "code": "BAD_REQUEST",
    "message": "Invalid request body"
  }
}
```

### 401 Unauthorized
```json
{
  "success": false,
  "error": {
    "code": "UNAUTHORIZED",
    "message": "Missing or invalid token"
  }
}
```

### 403 Forbidden
```json
{
  "success": false,
  "error": {
    "code": "FORBIDDEN",
    "message": "Only company users can perform this action"
  }
}
```

### 422 Unprocessable Entity (Validation Error)
```json
{
  "success": false,
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Validation failed",
    "details": {
      "title": ["Title is required"],
      "description": ["Description must be at least 20 characters"],
      "salary_min": ["Salary must be a positive number"]
    }
  }
}
```

---

## Testing Checklist

- [ ] Company registration with valid data
- [ ] Company login with correct credentials
- [ ] View company profile information
- [ ] Admin login with correct credentials
- [ ] Admin can list pending companies
- [ ] Admin can verify/approve company
- [ ] Company can create job posting (after verification)
- [ ] Job posting created in draft status
- [ ] Company dashboard shows job statistics
- [ ] Company can publish job posting
- [ ] Company can pause/close job posting
- [ ] Company can reopen job posting
- [ ] Job statistics update after state changes

---

## Tips untuk Testing

1. **Simpan tokens dalam environment variables:**
   ```bash
   export COMPANY_TOKEN="your_token"
   export ADMIN_TOKEN="your_token"
   ```

2. **Gunakan Postman untuk testing lebih mudah:**
   - Import `postman_collection.json`
   - Set environment variables untuk tokens

3. **Monitor database changes:**
   ```sql
   SELECT * FROM users WHERE role = 'company' ORDER BY created_at DESC LIMIT 1;
   SELECT * FROM jobs WHERE company_id = 9 ORDER BY created_at DESC;
   ```

4. **Check server logs untuk debugging:**
   ```bash
   tail -f server.log
   ```

