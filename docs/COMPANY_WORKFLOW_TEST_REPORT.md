# Company Workflow Testing - Complete Report

## ✅ Testing Summary

Workflow testing untuk company dari registrasi hingga membuat dan manage job postings **BERHASIL DIJALANKAN**.

---

## 📊 Test Results

### STEP 1: Company Registration ✅
- **Status**: Success (201 Created)
- **Company Email**: company.testing1768808883@karirnusantara.com
- **Company ID**: 10
- **Company Name**: PT Testing Indonesia
- **Initial Status**: Pending (belum verified admin)
- **Access Token**: Generated successfully

### STEP 2: Company Login ✅
- **Status**: Success (200 OK)
- **Authentication**: Verified with email & password
- **Token Valid**: Yes

### STEP 3: Get Current Company Information ✅
- **Status**: Success (200 OK)
- **Company Status**: Pending
- **Role**: company
- **Is Active**: true

### STEP 4: Admin Login (Skipped) ⏭️
- **Reason**: Using SKIP_ADMIN_VERIFICATION=true untuk testing
- **Note**: Admin akun ada di database, password hash perlu dikonfirmasi via phpMyAdmin XAMPP

### STEP 5: Admin Verifies Company ⏭️
- **Skipped**: Tidak ada admin token, tapi company tetap bisa membuat job

### STEP 6-7: Create Job Postings ✅

#### Job 1: Senior Backend Engineer
- **ID**: 4
- **Status**: Draft → Published ✅
- **Title**: Senior Backend Engineer
- **Salary**: IDR 15-25 juta
- **Location**: Jakarta Selatan, Remote
- **Skills**: Go, PostgreSQL, Docker, Kubernetes, Redis
- **Response**: Success (201 Created)

#### Job 2: Full Stack Developer
- **ID**: 5
- **Status**: Draft
- **Title**: Full Stack Developer
- **Salary**: IDR 8-12 juta
- **Location**: Jakarta Pusat, Remote
- **Skills**: React, Node.js, MongoDB, Docker
- **Response**: Success (201 Created)

#### Job 3: UI/UX Designer
- **ID**: 6
- **Status**: Draft
- **Title**: UI/UX Designer
- **Salary**: IDR 6-10 juta
- **Location**: Jakarta Selatan, Remote
- **Skills**: Figma, UI Design, UX Research, Prototyping
- **Response**: Success (201 Created)

### STEP 8: Get Dashboard Statistics ⚠️
- **Endpoint**: /api/v1/dashboard/stats
- **Status**: 404 Not Found
- **Note**: Route kemungkinan belum di-register dengan benar atau path berbeda
- **Action Needed**: Verify dashboard routes configuration

### STEP 9: List Company Jobs ✅
- **Status**: Success (200 OK)
- **Total Jobs Retrieved**: 2 (dari semua company)
- **Response**: Success dengan pagination
- **Note**: List jobs menampilkan semua jobs dari semua companies, bukan hanya company yang login

### STEP 10: Publish Job ✅
- **Job ID**: 4
- **Status Change**: draft → active ✅
- **Published At**: 2026-01-19T14:48:04+07:00
- **Response**: Success

---

## 🎯 Core Features Tested

| Feature | Status | Notes |
|---------|--------|-------|
| Company Registration | ✅ | Berhasil dengan role='company' |
| Company Login | ✅ | Email & password verification working |
| View Company Profile | ✅ | GET /auth/me working |
| Create Job (Draft) | ✅ | Multiple jobs created successfully |
| Publish Job | ✅ | Status changed from draft to active |
| Job Validation | ✅ | Title, description, salary, skills validated |
| Skills Management | ✅ | Array of skills stored correctly |
| Location/Remote | ✅ | City, province, is_remote fields working |
| Job Salary | ✅ | Min & max salary with currency (IDR) |
| Job Metadata | ✅ | views_count, applications_count tracking |
| Pagination | ✅ | Meta info with page, per_page, total_items |

---

## 🚀 Endpoints Tested & Working

```
✅ POST /api/v1/auth/register        - Company registration
✅ POST /api/v1/auth/login           - Company login
✅ GET  /api/v1/auth/me              - Get current user info
✅ POST /api/v1/jobs                 - Create job posting
✅ PATCH /api/v1/jobs/{id}/publish   - Publish job
✅ GET  /api/v1/jobs                 - List jobs with pagination
⚠️  GET  /api/v1/dashboard/stats     - Dashboard (404 Not Found)
❌ POST /api/v1/admin/auth/login     - Admin login (needs verification)
❌ POST /api/v1/admin/companies/{id}/verify - Company verification (needs admin token)
```

---

## 📝 API Requests Examples

### 1. Company Registration
```bash
curl -X POST http://localhost:8081/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "company@example.com",
    "password": "Company@123456",
    "full_name": "CEO Company",
    "phone": "081234567890",
    "company_name": "PT Company",
    "company_description": "Company description",
    "company_website": "https://company.com",
    "role": "company"
  }'
```

### 2. Create Job Posting
```bash
curl -X POST http://localhost:8081/api/v1/jobs \
  -H "Authorization: Bearer <COMPANY_TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Senior Backend Engineer",
    "description": "Job description...",
    "requirements": "Requirements...",
    "responsibilities": "Responsibilities...",
    "benefits": "Benefits...",
    "city": "Jakarta Selatan",
    "province": "DKI Jakarta",
    "is_remote": true,
    "job_type": "full_time",
    "experience_level": "senior",
    "salary_min": 15000000,
    "salary_max": 25000000,
    "salary_currency": "IDR",
    "is_salary_visible": true,
    "skills": ["Go", "PostgreSQL", "Docker"]
  }'
```

### 3. Publish Job
```bash
curl -X PATCH http://localhost:8081/api/v1/jobs/4/publish \
  -H "Authorization: Bearer <COMPANY_TOKEN>"
```

### 4. List Jobs
```bash
curl -X GET http://localhost:8081/api/v1/jobs?page=1&per_page=10 \
  -H "Authorization: Bearer <COMPANY_TOKEN>"
```

---

## ⚙️ Configuration Tested

- **API Base URL**: http://localhost:8081/api/v1
- **Database**: karir_nusantara (MySQL via XAMPP)
- **Environment**: Local development
- **JWT Token**: Working (900s expiry for company, 86400s for admin)

---

## 🔧 Issues Found & Solutions

### Issue 1: Admin Login Failed ❌
**Problem**: Admin credentials tidak diterima  
**Root Cause**: Password hash di database mungkin tidak sesuai  
**Solution**:  
- Update via phpMyAdmin di XAMPP
- Gunakan bcrypt hash yang benar untuk password 'admin123'

**phpMyAdmin Steps:**
1. Buka: http://localhost/phpmyadmin
2. Database: karir_nusantara → Table: users
3. Find/Create user dengan email: admin@karirnusantara.com
4. Update password_hash dengan bcrypt hash yang valid

### Issue 2: Dashboard Stats Endpoint 404 ⚠️
**Problem**: GET /dashboard/stats returns 404  
**Root Cause**: Route mungkin belum di-register atau path berbeda  
**Solution**:  
- Verify dashboard routes di internal/modules/dashboard/routes.go
- Atau gunakan endpoint alternatif untuk stats

### Issue 3: Job List Shows All Companies 📝
**Problem**: GET /jobs menampilkan jobs dari semua companies  
**Note**: Ini adalah behavior yang mungkin diinginkan (job marketplace public list)  
**Expected**: Company hanya melihat jobs mereka sendiri via dashboard atau filter

---

## ✨ Workflow Complete Flow

```
1. Company Registration
   ↓
2. Company Login (get access token)
   ↓
3. View Company Profile (verify status)
   ↓
4. Create Job Posting (status: draft)
   ↓
5. Publish Job (status: active)
   ↓
6. View Managed Jobs (list dan statistics)
   ↓
7. Edit/Pause/Close Job (manage postings)
```

---

## 📋 Testing Checklist

- [x] Company dapat register dengan role='company'
- [x] Company dapat login dengan email & password
- [x] Company dapat view profile mereka
- [x] Company dapat create job posting
- [x] Company dapat create multiple job postings
- [x] Job posting awalnya dalam status 'draft'
- [x] Job dapat di-publish (status → 'active')
- [x] Job list menampilkan dengan pagination
- [x] Job metadata (skills, salary, location) tersimpan
- [x] Job timestamps (created_at, published_at) recorded
- [ ] Admin dapat verify company
- [ ] Dashboard stats endpoint working
- [ ] Company hanya bisa see jobs mereka (jika needed)
- [ ] Job dapat di-pause
- [ ] Job dapat di-close
- [ ] Job dapat di-reopen

---

## 🚀 Next Steps

1. **Fix Admin Login**: 
   - Verify password hash di phpMyAdmin
   - Generate bcrypt hash yang benar untuk 'admin123'

2. **Fix Dashboard Endpoint**:
   - Check routes registration
   - Verify endpoint path

3. **Add Company Verification Flow**:
   - Admin dapat approve/reject company
   - Company status update ke 'verified'

4. **Add Pagination Filters**:
   - Filter jobs by company_id
   - Filter jobs by status (draft, active, paused, closed)

5. **Add Job Management**:
   - PATCH /jobs/{id}/pause
   - PATCH /jobs/{id}/close
   - PATCH /jobs/{id}/reopen

---

## 🎯 Test Credentials (from latest run)

**Company:**
- Email: company.testing1768808883@karirnusantara.com
- Password: Company@123456
- Company Name: PT Testing Indonesia
- Company ID: 10

**Admin:**
- Email: admin@karirnusantara.com
- Password: admin123
- Status: Needs verification/testing

---

## 📚 Running the Tests

### Without Admin Verification (Skip)
```bash
cd /Users/putramac/Desktop/Loker/karir-nusantara-api
SKIP_ADMIN_VERIFICATION=true bash tests/company_workflow_test.sh
```

### With Admin Verification (After fixing admin login)
```bash
cd /Users/putramac/Desktop/Loker/karir-nusantara-api
bash tests/company_workflow_test.sh
```

### Individual Endpoint Testing
Use the curl examples provided above or import `postman_collection.json` ke Postman.

---

**Generated**: 2026-01-19  
**API Server**: Running on http://localhost:8081  
**Status**: ✅ Core workflow functional, ready for full integration testing

