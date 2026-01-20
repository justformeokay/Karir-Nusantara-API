-- ============================================================
-- MANUAL VERIFICATION TEST UNTUK KARIR NUSANTARA
-- ============================================================
-- Jalankan query ini di phpMyAdmin atau MySQL CLI di XAMPP

-- ============================================================
-- STEP 1: Check current company status sebelum verification
-- ============================================================
SELECT id, user_id, company_name, company_status, documents_verified_at 
FROM companies 
WHERE user_id = 7;

-- Expected result: company_id=1, user_id=7, company_status='pending'

-- ============================================================
-- STEP 2: Test UPDATE query - VERIFY COMPANY
-- ============================================================
UPDATE companies 
SET company_status = 'verified' 
WHERE user_id = 7;

-- Check rows affected
SELECT ROW_COUNT() as rows_affected;
-- Expected: 1

-- ============================================================
-- STEP 3: Check updated status AFTER verification
-- ============================================================
SELECT id, user_id, company_name, company_status, documents_verified_at 
FROM companies 
WHERE user_id = 7;

-- Expected result: company_status='verified'

-- ============================================================
-- STEP 4: If still pending, debug by checking:
-- ============================================================

-- Check if company exists with user_id=7
SELECT COUNT(*) as company_count FROM companies WHERE user_id = 7;

-- Check if user_id=7 is actually a company user
SELECT id, email, role FROM users WHERE id = 7;

-- Check all companies in database
SELECT id, user_id, company_name, company_status FROM companies;

-- ============================================================
-- STEP 5: If queries work, revert changes back to pending
-- ============================================================
UPDATE companies 
SET company_status = 'pending' 
WHERE user_id = 7;

-- Verify reverted
SELECT id, user_id, company_name, company_status 
FROM companies 
WHERE user_id = 7;

-- ============================================================
-- NOTES:
-- ============================================================
-- - company_id di companies table adalah PRIMARY KEY (id)
-- - user_id adalah FK ke users.id
-- - Untuk verify company dengan user_id=7, query harus:
--   UPDATE companies SET company_status='verified' WHERE user_id=7
-- - company_status enum values: 'pending','verified','rejected','suspended'
