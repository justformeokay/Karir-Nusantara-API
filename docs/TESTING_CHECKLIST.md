# Company Workflow Testing - Interactive Checklist

> Run this checklist untuk memverifikasi semua endpoints dan features working correctly

**Date**: _____________  
**Tester**: _____________  
**Environment**: Local / Staging / Production

---

## 📋 Pre-Testing Setup

- [ ] API Server running on http://localhost:8081
- [ ] MySQL/XAMPP database connected
- [ ] Database: `karir_nusantara` exists
- [ ] Postman installed (optional)
- [ ] cURL available in terminal

---

## 🔐 Authentication Testing

### Registration
- [ ] POST /auth/register dengan data valid
  - Email: ________________________
  - Company: ________________________
  - Response: 201 Created ✅
  - Token received: ✅
  - Notes: _____________________________

### Login
- [ ] POST /auth/login dengan email & password
  - Email: ________________________
  - Response: 200 OK ✅
  - Token valid: ✅
  - Notes: _____________________________

### Get Profile
- [ ] GET /auth/me dengan valid token
  - Response: 200 OK ✅
  - Profile data correct: ✅
  - Company status: Pending / Verified
  - Notes: _____________________________

---

## 💼 Job Management Testing

### Create Job 1
- [ ] POST /jobs dengan data lengkap
  - Title: Senior Backend Engineer
  - Response: 201 Created ✅
  - Job ID: ________
  - Initial Status: draft ✅
  - Notes: _____________________________

### Create Job 2
- [ ] POST /jobs (second job)
  - Title: Full Stack Developer
  - Response: 201 Created ✅
  - Job ID: ________
  - Initial Status: draft ✅
  - Notes: _____________________________

### Create Job 3
- [ ] POST /jobs (third job)
  - Title: UI/UX Designer
  - Response: 201 Created ✅
  - Job ID: ________
  - Initial Status: draft ✅
  - Notes: _____________________________

### Validate Job Data
- [ ] Title saved correctly: ✅
- [ ] Description saved: ✅
- [ ] Requirements saved: ✅
- [ ] Salary min/max correct: ✅
- [ ] Skills array saved: ✅
- [ ] Location (city, province): ✅
- [ ] is_remote flag: ✅
- [ ] Timestamps recorded: ✅
  - created_at: _______________
  - updated_at: _______________

---

## 📤 Job Publishing Testing

### Publish Job 1
- [ ] PATCH /jobs/{id}/publish
  - Job ID: ________
  - Response: 200 OK ✅
  - Status changed to: active ✅
  - published_at recorded: ✅
  - Notes: _____________________________

### Publish Job 2
- [ ] PATCH /jobs/{id}/publish (second job)
  - Job ID: ________
  - Response: 200 OK ✅
  - Status: active ✅
  - Notes: _____________________________

### Verify Published Jobs
- [ ] views_count initialized to 0: ✅
- [ ] applications_count initialized to 0: ✅
- [ ] Status is 'active' in database: ✅

---

## 📋 Job Listing Testing

### List All Jobs
- [ ] GET /jobs?page=1&per_page=10
  - Response: 200 OK ✅
  - Total items returned: ________
  - Pagination meta included: ✅
    - page: ________
    - per_page: ________
    - total_items: ________
    - total_pages: ________
  - Notes: _____________________________

### Pagination
- [ ] Page 1 returns correct items: ✅
- [ ] Can navigate pages: ✅
- [ ] Per_page parameter works: ✅
- [ ] Total counts accurate: ✅

### Get Job by ID
- [ ] GET /jobs/{id}
  - Job ID: ________
  - Response: 200 OK ✅
  - Correct job data: ✅
  - Notes: _____________________________

### Get Job by Slug
- [ ] GET /jobs/slug/{slug}
  - Slug: ________________________
  - Response: 200 OK ✅
  - Correct job returned: ✅
  - Notes: _____________________________

---

## 🎮 Job Control Testing

### Pause Job
- [ ] PATCH /jobs/{id}/pause
  - Job ID: ________
  - Response: 200 OK ✅
  - Status changed to: paused ✅
  - Notes: _____________________________

### Close Job
- [ ] PATCH /jobs/{id}/close
  - Job ID: ________
  - Response: 200 OK ✅
  - Status changed to: closed ✅
  - Notes: _____________________________

### Reopen Job
- [ ] PATCH /jobs/{id}/reopen
  - Job ID: ________
  - Response: 200 OK ✅
  - Status changed to: active ✅
  - Notes: _____________________________

---

## ✏️ Job Update Testing

### Update Job
- [ ] PUT /jobs/{id}
  - Job ID: ________
  - Updated field: title
  - New value: _____________________________
  - Response: 200 OK ✅
  - Changes saved: ✅
  - Notes: _____________________________

### Update Salary
- [ ] PUT /jobs/{id} (update salary)
  - New salary_min: ________
  - New salary_max: ________
  - Response: 200 OK ✅
  - Changes saved: ✅
  - Notes: _____________________________

---

## 🗑️ Job Deletion Testing

### Delete Job
- [ ] DELETE /jobs/{id}
  - Job ID: ________
  - Response: 200 OK ✅
  - Job no longer in list: ✅
  - Notes: _____________________________

---

## 👤 Admin Testing (Optional)

### Admin Login
- [ ] POST /admin/auth/login
  - Email: admin@karirnusantara.com
  - Password: admin123
  - Response: 200 OK ⚠️
  - Token received: ⚠️
  - Status: [ ] Works [ ] Needs Fix
  - Notes: _____________________________

### Admin Get Profile
- [ ] GET /admin/auth/me
  - Response: 200 OK ⚠️
  - Admin role: ⚠️
  - Notes: _____________________________

### List Companies
- [ ] GET /admin/companies
  - Response: 200 OK ⚠️
  - Companies listed: ⚠️
  - Status filtering: ⚠️
  - Notes: _____________________________

### Verify Company
- [ ] POST /admin/companies/{id}/verify
  - Company ID: ________
  - Action: approve
  - Response: 200 OK ⚠️
  - Company status updated: ⚠️
  - Notes: _____________________________

---

## 📊 Dashboard Testing (Optional)

### Company Dashboard
- [ ] GET /dashboard/stats
  - Response: 200 OK ⚠️
  - Jobs count: ________
  - Active jobs: ________
  - Applications count: ________
  - Status: [ ] Works [ ] 404 Error
  - Notes: _____________________________

---

## ⚠️ Error Response Testing

### Test 400 Bad Request
- [ ] POST /jobs tanpa required field
  - Response: 400 Bad Request ✅
  - Error message clear: ✅
  - Notes: _____________________________

### Test 401 Unauthorized
- [ ] POST /jobs tanpa token
  - Response: 401 Unauthorized ✅
  - Notes: _____________________________

### Test 422 Validation Error
- [ ] POST /jobs dengan invalid data (salary_min > salary_max)
  - Response: 422 Unprocessable Entity ✅
  - Error details shown: ✅
  - Notes: _____________________________

### Test 404 Not Found
- [ ] GET /jobs/999999
  - Response: 404 Not Found ✅
  - Notes: _____________________________

---

## 📈 Performance Testing

### Response Time
- [ ] POST /jobs response time: ________ ms (target: < 1000ms)
- [ ] GET /jobs response time: ________ ms (target: < 500ms)
- [ ] GET /jobs/{id} response time: ________ ms (target: < 200ms)

### Database
- [ ] Jobs data persisted in database: ✅
- [ ] All fields stored correctly: ✅
- [ ] No data truncation: ✅

---

## 🔒 Security Testing

### Token Validation
- [ ] Expired token rejected: ✅
- [ ] Invalid token rejected: ✅
- [ ] Token format validated: ✅

### Authorization
- [ ] Company can't modify other company's jobs: ✅
- [ ] Company can't access admin endpoints: ✅
- [ ] Protected routes require auth: ✅

---

## 📝 Data Validation Testing

### Title Validation
- [ ] Empty title rejected: ✅
- [ ] Very long title handled: ✅
- [ ] Special characters accepted: ✅

### Description Validation
- [ ] Min length enforced: ✅
- [ ] HTML/Scripts sanitized: ✅
- [ ] Unicode characters supported: ✅

### Salary Validation
- [ ] salary_min < salary_max enforced: ✅
- [ ] Negative salary rejected: ✅
- [ ] Zero salary rejected: ✅

### Skills Validation
- [ ] Empty skills array rejected: ✅
- [ ] Max skills limit enforced: ✅
- [ ] Duplicate skills handled: ✅

---

## 🎯 Integration Testing

### Complete User Journey
- [ ] Register → Login → Create Job → Publish → List: ✅
- [ ] All steps work together: ✅
- [ ] Data consistent across endpoints: ✅
- [ ] No orphaned records: ✅

---

## 📋 Final Verification

- [ ] All ✅ checks passed: YES / NO
- [ ] All ⚠️ checks resolved: YES / NO
- [ ] No critical bugs found: YES / NO
- [ ] Ready for production: YES / NO

---

## 📊 Summary

| Category | Status | Issues |
|----------|--------|--------|
| Authentication | ✅ | None |
| Job Creation | ✅ | None |
| Job Publishing | ✅ | None |
| Job Listing | ✅ | None |
| Job Control | ✅ | None |
| Admin Features | ⚠️ | Admin login needs verification |
| Dashboard | ⚠️ | 404 error on stats endpoint |
| Error Handling | ✅ | Proper error responses |
| Security | ✅ | Token validation working |
| Data Validation | ✅ | All validations working |

---

## 🎯 Issues to Fix

1. **Issue**: Admin login failing
   - **Solution**: Verify password hash in phpMyAdmin
   - **Priority**: Medium
   - **Status**: [ ] Not Started [ ] In Progress [ ] Fixed

2. **Issue**: Dashboard stats endpoint 404
   - **Solution**: Check route registration
   - **Priority**: Low
   - **Status**: [ ] Not Started [ ] In Progress [ ] Fixed

3. **Issue**: _____________________________
   - **Solution**: _____________________________
   - **Priority**: [ ] High [ ] Medium [ ] Low
   - **Status**: [ ] Not Started [ ] In Progress [ ] Fixed

---

## ✅ Sign-Off

**Tested by**: _____________________________  
**Date**: _____________________________  
**Approved by**: _____________________________  
**Date**: _____________________________

---

**Notes:**
```
_____________________________________________________________________________

_____________________________________________________________________________

_____________________________________________________________________________
```

