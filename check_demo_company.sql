-- Check demo company data
SELECT 'User Info:' as info;
SELECT id, email, role FROM users WHERE email = 'localhosting127.0.0.1@gmail.com';

SELECT '\nCompany Info:' as info;
SELECT c.id, c.user_id, c.company_name, c.company_status
FROM companies c
JOIN users u ON c.user_id = u.id
WHERE u.email = 'localhosting127.0.0.1@gmail.com';

SELECT '\nCompany Quota Info:' as info;
SELECT cq.id, cq.company_id, cq.free_quota_used, cq.paid_quota
FROM company_quotas cq
JOIN companies c ON cq.company_id = c.id
JOIN users u ON c.user_id = u.id
WHERE u.email = 'localhosting127.0.0.1@gmail.com';

SELECT '\nJobs Posted:' as info;
SELECT COUNT(*) as total_jobs, 
       SUM(CASE WHEN status = 'active' THEN 1 ELSE 0 END) as active_jobs
FROM jobs j
JOIN companies c ON j.company_id = c.id
JOIN users u ON c.user_id = u.id
WHERE u.email = 'localhosting127.0.0.1@gmail.com';
