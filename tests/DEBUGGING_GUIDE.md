# DEBUGGING GUIDE: Company Verification Issue

## Problem
Admin tidak dapat memverifikasi company account menggunakan API endpoint POST /api/v1/admin/companies/{id}/verify

Error: "UPDATE_FAILED" - Gagal memverifikasi perusahaan

## Diagnosis Steps

### Step 1: Verify Database Query Works
Open XAMPP phpMyAdmin and run these SQL queries:

```sql
-- Check current company status
SELECT id, user_id, company_name, company_status FROM companies WHERE user_id = 7;
```

Expected: 
- id: 1 (company_id)
- user_id: 7
- company_name: PT Karya Developer indonesia
- company_status: pending

```sql
-- Test UPDATE query
UPDATE companies SET company_status = 'verified' WHERE user_id = 7;
SELECT ROW_COUNT() as rows_affected;
```

Expected: rows_affected = 1 (meaning 1 row was updated)

```sql
-- Verify the update worked
SELECT id, user_id, company_name, company_status FROM companies WHERE user_id = 7;
```

Expected: company_status should now be 'verified'

### Step 2: If SQL Query Works But API Still Fails

The issue might be in the Go application code. Here's what could be wrong:

1. **Parameter Binding Issue**: The query uses `?` placeholders which should be replaced with actual values
   - Status value: "verified"
   - User ID value: 7

2. **Response Handling**: The handler might not correctly return success

### Step 3: Debug Output
The updated code has DEBUG output. Look for log messages like:
```
[DEBUG] UpdateCompanyStatus: query=UPDATE companies SET company_status = ? WHERE user_id = ?, status=verified, user_id=7
[DEBUG] Rows affected: 1
```

If rows_affected = 0, it means no company row matches user_id=7

If you see database error, check what error message is shown.

## Quick Fix Options

### Option A: Direct Database Update (Quickest for Testing)
If SQL works but API doesn't, directly update the database in phpMyAdmin:

```sql
UPDATE companies SET company_status = 'verified' WHERE user_id = 7;
```

Then refresh the frontend - it should work.

### Option B: Check API Logs
1. Start API: `cd karir-nusantara-api && go run ./cmd/api/main.go`
2. Run test
3. Check console output for [DEBUG] messages

### Option C: Verify with Updated Company
After either fixing or manually updating, verify in frontend:

1. Login as: info@karyadeveloperindonesia.com / Justformeokay23@
2. Go to Dashboard
3. Check if "Buat Lowongan" button is now ENABLED
4. Try creating a job posting

## Expected Behavior After Verification

Once company_status changes to 'verified':
- Dashboard button "Buat Lowongan" should be ENABLED
- Can access /jobs/new without blocking modal  
- Can create job postings
- All dashboard features available

## Troubleshooting Checklist

- [ ] SQL query works in phpMyAdmin
- [ ] rows_affected = 1 in phpMyAdmin
- [ ] company_status shows 'verified' in phpMyAdmin query result
- [ ] Frontend localStorage cleared (refreshed page)
- [ ] Logged in as company user again
- [ ] Dashboard button shows as enabled
- [ ] Can create new job posting

## Next Steps

1. First, run the SQL queries in phpMyAdmin to confirm the database is accessible and query works
2. If SQL works, manually update via phpMyAdmin and test frontend
3. If SQL doesn't work, check database table structure with:
   ```sql
   DESC companies;
   ```
4. If SQL works but API still fails, we need to debug the Go code further

## Contact Points

**For API debugging:**
- Check `/tmp/api_debug.log` file
- Look for [DEBUG] prefix in logs
- Share any error messages shown

**For Database:**
- Use phpMyAdmin GUI
- Check: http://localhost/phpmyadmin
- Select database: karir_nusantara
- Browse companies table

**Key Credentials:**
- Admin: admin@karirnusantara.com / admin123
- Test Company: info@karyadeveloperindonesia.com / Justformeokay23@
