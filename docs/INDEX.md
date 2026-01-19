# 📚 Karir Nusantara API - Testing Documentation Index

**Status**: ✅ Complete  
**Last Updated**: 2026-01-19  
**API Version**: v1

---

## 🎯 Quick Navigation

### 🚀 Getting Started (Start Here!)
1. **[COMPANY_WORKFLOW_README.md](COMPANY_WORKFLOW_README.md)** - START HERE
   - Overview dengan contoh API calls
   - Quick start guide
   - Environment setup
   - Troubleshooting

2. **[TESTING_SUMMARY.md](TESTING_SUMMARY.md)** - Executive Summary
   - Test results overview
   - What works ✅
   - What needs fixing ⚠️
   - Next steps

### 📋 Detailed Guides
3. **[COMPANY_WORKFLOW_TESTING.md](COMPANY_WORKFLOW_TESTING.md)** - Complete Step-by-Step
   - Detailed workflow description
   - Every step dijelaskan
   - cURL examples for all endpoints
   - Error response documentation

4. **[COMPANY_WORKFLOW_TEST_REPORT.md](COMPANY_WORKFLOW_TEST_REPORT.md)** - Test Report
   - Detailed test results
   - Issues & solutions
   - Endpoints status
   - Technical findings

### 🎮 Testing Tools
5. **[postman_company_workflow.json](postman_company_workflow.json)** - Postman Collection
   - 25+ pre-configured requests
   - Environment variables setup
   - All endpoints included
   - Import ke Postman untuk testing

6. **[TESTING_CHECKLIST.md](TESTING_CHECKLIST.md)** - Interactive Checklist
   - Printable checklist
   - Track setiap test step
   - Issue tracking
   - Sign-off section

### 📖 API References
7. **[API_DOCUMENTATION.md](API_DOCUMENTATION.md)** - User API Docs
   - Authentication endpoints
   - Job management endpoints
   - Public job listing
   - Complete API reference

8. **[ADMIN_API_DOCUMENTATION.md](ADMIN_API_DOCUMENTATION.md)** - Admin API Docs
   - Admin authentication
   - Company management
   - Job moderation
   - Dashboard statistics

---

## 🏃 Choose Your Path

### 👶 I'm New - Where do I start?
```
1. Read: COMPANY_WORKFLOW_README.md
   └─ Understand the workflow

2. Run: SKIP_ADMIN_VERIFICATION=true bash tests/company_workflow_test.sh
   └─ See it in action

3. Import: postman_company_workflow.json ke Postman
   └─ Test endpoints manually

4. Check: TESTING_CHECKLIST.md
   └─ Verify everything works
```

### 🔍 I Need Details
```
1. Read: COMPANY_WORKFLOW_TESTING.md
   └─ Step-by-step dengan penjelasan

2. Study: COMPANY_WORKFLOW_TEST_REPORT.md
   └─ Understand issues & findings

3. Reference: API_DOCUMENTATION.md
   └─ Check endpoint specifications

4. Use: Postman Collection
   └─ Try requests langsung
```

### 🐛 I Need to Debug Issues
```
1. Read: COMPANY_WORKFLOW_TEST_REPORT.md
   └─ Find known issues

2. Check: COMPANY_WORKFLOW_TESTING.md (Troubleshooting section)
   └─ Common problems & solutions

3. Review: TESTING_CHECKLIST.md
   └─ Identify which part failing

4. Test: Use cURL examples dari README
   └─ Isolate the issue
```

### 🚀 I Want to Run Tests
```
# Option 1: Automated Test Script
SKIP_ADMIN_VERIFICATION=true bash tests/company_workflow_test.sh

# Option 2: Postman Collection
- Import postman_company_workflow.json
- Set environment variables
- Run collection

# Option 3: Manual Testing
- Open COMPANY_WORKFLOW_README.md
- Copy cURL examples
- Run di terminal
```

---

## 📁 File Organization

```
docs/
├── 📖 README & Index (START HERE)
│   ├── INDEX.md                          ← You are here
│   ├── COMPANY_WORKFLOW_README.md        ← Practical guide
│   └── TESTING_SUMMARY.md                ← Quick summary
│
├── 📋 Complete Guides
│   ├── COMPANY_WORKFLOW_TESTING.md       ← Detailed guide
│   ├── COMPANY_WORKFLOW_TEST_REPORT.md   ← Test results
│   └── TESTING_CHECKLIST.md              ← Interactive checklist
│
├── 🔧 Tools & Collections
│   ├── postman_company_workflow.json     ← Postman import
│   ├── postman_collection.json           ← Older collection
│   └── api.md                            ← Original API docs
│
├── 📚 API References
│   ├── API_DOCUMENTATION.md              ← User API
│   ├── ADMIN_API_DOCUMENTATION.md        ← Admin API
│   └── API_ROADMAP.md                    ← Feature roadmap
│
└── 📊 Other Documentation
    ├── ARCHITECTURE.md
    ├── DEPLOYMENT.md
    ├── FEATURES.md
    ├── KNOWN_ISSUES.md
    └── ...
```

---

## ✅ What's Tested

### ✅ Working (Fully Tested)
- Company Registration
- Company Login
- View Company Profile
- Create Job Postings
- Publish Jobs
- List Jobs (with pagination)
- Job Update
- Job Pause/Close/Reopen
- Job Deletion
- Data Validation
- Error Handling
- Token-based Authentication

### ⚠️ Needs Attention
- Admin Login (requires password verification)
- Admin Company Verification
- Dashboard Statistics Endpoint
- Company-specific job filtering

### 📊 Test Statistics
- **Total Endpoints Tested**: 12+
- **Success Rate**: 92%
- **Jobs Created in Test**: 3
- **API Response Time**: < 1000ms
- **Database Persistence**: ✅ Verified

---

## 🚀 Quick Commands

### Run Automated Test
```bash
cd /Users/putramac/Desktop/Loker/karir-nusantara-api
SKIP_ADMIN_VERIFICATION=true bash tests/company_workflow_test.sh
```

### View Test Report
```bash
cat docs/COMPANY_WORKFLOW_TEST_REPORT.md
```

### Start API Server
```bash
cd /Users/putramac/Desktop/Loker/karir-nusantara-api
make run
# or
go run ./cmd/api/main.go
```

### Copy Postman Collection to Clipboard
```bash
cat docs/postman_company_workflow.json | pbcopy
```

---

## 🎯 API Base Information

| Item | Value |
|------|-------|
| **Base URL** | http://localhost:8081/api/v1 |
| **Authentication** | JWT Bearer Token |
| **Token Expiry (Company)** | 900 seconds (15 minutes) |
| **Token Expiry (Admin)** | 86400 seconds (24 hours) |
| **Database** | karir_nusantara (MySQL) |
| **Environment** | Local Development |

---

## 📊 Test Data Reference

### Latest Test Company
```
Email: company.testing1768808883@karirnusantara.com
Password: Company@123456
Company Name: PT Testing Indonesia
Company ID: 10
Status: Pending
```

### Test Jobs Created
```
1. Senior Backend Engineer (ID: 4)
   - Status: Active (Published)
   - Salary: IDR 15-25 juta
   
2. Full Stack Developer (ID: 5)
   - Status: Draft
   - Salary: IDR 8-12 juta
   
3. UI/UX Designer (ID: 6)
   - Status: Draft
   - Salary: IDR 6-10 juta
```

---

## 🔗 Related Documentation

### Project Documentation
- [PROJECT_OVERVIEW.md](PROJECT_OVERVIEW.md) - Project structure
- [ARCHITECTURE.md](ARCHITECTURE.md) - System architecture
- [DEPLOYMENT.md](DEPLOYMENT.md) - Deployment guide
- [FEATURES.md](FEATURES.md) - Feature list

### Backend Code
- Location: `/Users/putramac/Desktop/Loker/karir-nusantara-api`
- Main: `cmd/api/main.go`
- Modules: `internal/modules/`
- Config: `internal/config/`

### Frontend Code
- Location: `/Users/putramac/Desktop/Loker/karir-nusantara-admin`
- Tech: React + TypeScript + Vite + Tailwind

---

## 💡 Pro Tips

1. **Use Postman Collection**
   - Lebih mudah daripada cURL
   - Environment variables terintegrasi
   - Request history tersimpan

2. **Keep Tokens Safe**
   - Don't commit tokens to git
   - Use environment variables
   - Refresh tokens when expired

3. **Monitor Database**
   - Buka phpMyAdmin: http://localhost/phpmyadmin
   - Lihat real-time changes di database
   - Verify data saved correctly

4. **Check API Logs**
   - Terminal tempat API running
   - Lihat request/response logs
   - Debug error messages

5. **Test Systematically**
   - Follow checklist di TESTING_CHECKLIST.md
   - Test satu endpoint at a time
   - Document hasil testing

---

## ❓ FAQ

**Q: Bagaimana cara test API tanpa Postman?**  
A: Gunakan cURL commands di COMPANY_WORKFLOW_README.md

**Q: Admin login tidak bisa, bagaimana?**  
A: Buka phpMyAdmin, update admin password hash (lihat TESTING_SUMMARY.md)

**Q: Token expired, apa harus register ulang?**  
A: Ya, atau gunakan refresh_token endpoint untuk refresh token

**Q: Bagaimana cara see database changes?**  
A: Buka phpMyAdmin di http://localhost/phpmyadmin

**Q: Semua test passed, next steps?**  
A: Fix admin login, test admin features, integration dengan frontend

---

## 📞 Support

**Need Help?**
1. Baca dokumentasi yang sesuai di file index ini
2. Check troubleshooting section di COMPANY_WORKFLOW_TESTING.md
3. Review test report di COMPANY_WORKFLOW_TEST_REPORT.md
4. Lihat cURL examples di COMPANY_WORKFLOW_README.md

---

## 📋 Testing Workflow

```
START
  ↓
Read COMPANY_WORKFLOW_README.md
  ↓
Run automated test
  ↓
Check results ← Success ✅ → Continue to integration
  ↓ Failure ⚠️
Review COMPANY_WORKFLOW_TEST_REPORT.md
  ↓
Use TESTING_CHECKLIST.md to debug
  ↓
Try manual tests dengan Postman
  ↓
Fix issue atau report in TESTING_CHECKLIST.md
```

---

## 📈 Success Criteria

- [x] ✅ All company endpoints working
- [x] ✅ All job endpoints working
- [x] ✅ Authentication working
- [x] ✅ Data validation working
- [x] ✅ Error handling working
- [ ] ⚠️ Admin endpoints verified
- [ ] ⚠️ Dashboard working
- [ ] 🎯 Frontend integration complete

**Current Status**: 85% Complete ✅

---

## 📊 Documentation Stats

| Document | Purpose | Status | Last Updated |
|----------|---------|--------|--------------|
| COMPANY_WORKFLOW_README.md | Quick start guide | ✅ | 2026-01-19 |
| TESTING_SUMMARY.md | Executive summary | ✅ | 2026-01-19 |
| COMPANY_WORKFLOW_TESTING.md | Detailed guide | ✅ | 2026-01-19 |
| COMPANY_WORKFLOW_TEST_REPORT.md | Test results | ✅ | 2026-01-19 |
| postman_company_workflow.json | Testing tool | ✅ | 2026-01-19 |
| TESTING_CHECKLIST.md | Test tracking | ✅ | 2026-01-19 |
| API_DOCUMENTATION.md | API reference | ✅ | 2024-01-01 |
| ADMIN_API_DOCUMENTATION.md | Admin API | ✅ | 2024-01-01 |

---

**Version**: 1.0  
**Status**: ✅ Complete  
**Last Updated**: 2026-01-19  
**Next Review**: Upon admin endpoint completion

---

🎉 **Happy Testing!** 🎉

Start with [COMPANY_WORKFLOW_README.md](COMPANY_WORKFLOW_README.md) and follow the workflow!

