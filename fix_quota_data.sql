-- =====================================================
-- FIX DATA QUOTA UNTUK PT NUSA PERSADA
-- Jalankan query ini di phpMyAdmin
-- =====================================================

-- Cek data saat ini
SELECT 'DATA SEBELUM FIX:' as info;

SELECT '1. User Info:' as step;
SELECT id, email, role FROM users WHERE email = 'localhosting127.0.0.1@gmail.com';

SELECT '2. Company Info:' as step;
SELECT c.id as company_id, c.user_id, c.company_name, c.company_status
FROM companies c
WHERE c.user_id = 26;

SELECT '3. Company Quotas (semua record):' as step;
SELECT * FROM company_quotas WHERE company_id IN (7, 26);

SELECT '4. Jobs yang dimiliki company_id=7:' as step;
SELECT COUNT(*) as total_jobs, 
       SUM(CASE WHEN status = 'active' THEN 1 ELSE 0 END) as active_jobs,
       SUM(CASE WHEN published_at IS NOT NULL THEN 1 ELSE 0 END) as published_jobs
FROM jobs 
WHERE company_id = 7;

-- =====================================================
-- FIX 1: Hapus record quota yang salah (company_id=26 yang sebenarnya user_id)
-- =====================================================
DELETE FROM company_quotas WHERE company_id = 26;

-- =====================================================
-- FIX 2: Update quota untuk company_id=7 berdasarkan jumlah job yang sudah published
-- =====================================================
UPDATE company_quotas cq
SET cq.free_quota_used = (
    SELECT COUNT(*) FROM jobs j 
    WHERE j.company_id = 7 AND j.published_at IS NOT NULL
),
cq.updated_at = NOW()
WHERE cq.company_id = 7;

-- =====================================================
-- VERIFIKASI SETELAH FIX
-- =====================================================
SELECT 'DATA SETELAH FIX:' as info;

SELECT 'Company Quotas untuk PT Nusa Persada:' as step;
SELECT * FROM company_quotas WHERE company_id = 7;

SELECT 'Sisa Kuota Gratis:' as step;
SELECT 10 - free_quota_used as sisa_kuota_gratis FROM company_quotas WHERE company_id = 7;
