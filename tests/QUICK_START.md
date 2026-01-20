# QUICK START - KARIR NUSANTARA VERIFICATION TESTING

## 🎯 TUJUAN
Memverifikasi akun perusahaan (info@karyadeveloperindonesia.com) agar dapat menggunakan semua fitur dashboard.

## 🔑 CREDENTIALS

### Admin
```
Email: admin@karirnusantara.com
Password: admin123
```

### Sample Company (Karya Developer Indonesia)
```
Email: info@karyadeveloperindonesia.com
Password: Justformeokay23@
Database user_id: 7
```

---

## 📋 STEP-BY-STEP

### Step 1: Start API Server
```bash
cd /Users/putramac/Desktop/Loker/karir-nusantara-api
go run ./cmd/api/main.go
```
Expected: Server should start on http://localhost:8081

### Step 2: Test API Verification (Option A - Recommended First)
```bash
/Users/putramac/Desktop/Loker/karir-nusantara-api/tests/admin_verify_company.sh
```

Expected output:
- ✅ Admin login success
- ✅ Found company ID: 7  
- ✅ Verifying company (APPROVE)
- ❓ If success, company_status should change to "verified"

### Step 3: If API Test Fails - Use Database

**Option B - Manual Database Update (Quickest)**

1. Open XAMPP Control Panel
2. Click "Admin" next to MySQL → opens phpMyAdmin
3. Select database: `karir_nusantara`
4. Go to: Tools → Query Tool (or SQL tab)
5. Run this query:
   ```sql
   UPDATE companies SET company_status = 'verified' WHERE user_id = 7;
   ```
6. Verify result:
   ```sql
   SELECT id, user_id, company_name, company_status FROM companies WHERE user_id = 7;
   ```

Expected: `company_status` should show `verified`

---

## ✅ VERIFICATION IN FRONTEND

### Step 1: Browser
1. Go to http://localhost:5174 (frontend)
2. Login as: info@karyadeveloperindonesia.com
3. Password: Justformeokay23@

### Step 2: Check Dashboard
- Look at "Buat Lowongan" button
- It should be **ENABLED** (not grayed out)
- Click it - should open job creation form

### Step 3: Check Job Creation Features
1. Go to Dashboard
2. Click "Buat Lowongan" button
3. Form should NOT show blocking modal
4. All form fields should be accessible
5. Try filling form and publishing

### Step 4: Check Job Form Page
1. Navigate to /jobs/new
2. Should show green success alert: "Siap Membuat Lowongan"
3. Form should NOT be disabled/grayed out
4. Can fill all fields

---

## 🐛 TROUBLESHOOTING

### Database Query Works But API Still Fails
- Backend code has issue in parameter binding
- Try manually updating database (Option B)
- Frontend should work after manual update

### Button Still Grayed Out After Update
1. **Clear Browser Cache:**
   - Press Ctrl+Shift+R (hard refresh)
   - Or clear localStorage in developer tools
   
2. **Logout and Login Again**
   - Logout from company account
   - Clear browser cache
   - Login again
   - Check Dashboard

3. **Check Console Errors**
   - Press F12 → Console tab
   - Look for JavaScript errors
   - Check network requests

### API Not Starting
- Check if port 8081 is already in use
- Kill process: `killall go`
- Try again

---

## 📊 EXPECTED BEHAVIOR

### Before Verification
```
Company Status: pending
Dashboard Button: DISABLED (grayed out)
BlockModal: SHOWS (prevents access)
Job Creation: BLOCKED
```

### After Verification
```
Company Status: verified ✅
Dashboard Button: ENABLED (clickable)
BlockModal: HIDDEN
Job Creation: ALLOWED ✅
Semua dashboard features: ACCESSIBLE ✅
```

---

## 🗂️ FILE REFERENCE

**Test Scripts:**
- `admin_verify_company.sh` - Full verification test
- `MANUAL_VERIFICATION_DEBUG.sql` - SQL debugging queries

**Documentation:**
- `DEBUGGING_GUIDE.md` - Detailed troubleshooting
- `IMPLEMENTATION_STATUS.md` - Project status summary
- `VERIFICATION_API_TEST.md` - API endpoint documentation

**Key Code Files:**
- `internal/modules/admin/handler.go` - Verify endpoint
- `internal/modules/admin/service.go` - Business logic
- `internal/modules/admin/repository.go` - Database query
- `src/hooks/useCompanyEligibility.ts` - Frontend validation

---

## ⏱️ EXPECTED TIME

- Option A (API): ~2 minutes
- Option B (Database): ~1 minute
- Frontend testing: ~5 minutes

**Total: ~10 minutes to full verification & testing**

---

## 🎉 SUCCESS CRITERIA

✅ Company email: info@karyadeveloperindonesia.com  
✅ Status: verified (in database)  
✅ Dashboard button enabled  
✅ Can create job posting  
✅ All dashboard features accessible  

---

## NOTES

- Database uses XAMPP MySQL (no command-line access)
- API runs on port 8081
- Frontend runs on port 5174
- Test anytime needed - no data will be lost

---

**Untuk bantuan lebih lanjut, lihat DEBUGGING_GUIDE.md**
