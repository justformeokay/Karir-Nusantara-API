-- Test Data Seed for Karir Nusantara API Tests
-- This file creates test data for automated testing

-- ============================================
-- Test Companies
-- ============================================

-- Verified Company (password: TestPassword123!)
INSERT INTO companies (
    id, email, password_hash, full_name, company_name, company_description,
    company_website, company_industry, company_size, company_location,
    is_active, is_verified, verification_status, created_at, updated_at
) VALUES (
    'test-verified-company-id',
    'test-verified@company.com',
    '$2a$10$N9qo8uLOickgx2ZMRZoMye.IjqMOeX.yL9rBvP/9S/1e3pL.JY2Gy', -- TestPassword123!
    'Admin Verified',
    'Test Verified Company',
    'A verified test company for automated testing',
    'https://verified-company.test',
    'Technology',
    '51-200',
    'Jakarta, Indonesia',
    TRUE,
    TRUE,
    'verified',
    NOW(),
    NOW()
) ON DUPLICATE KEY UPDATE email = email;

-- Unverified Company (password: TestPassword123!)
INSERT INTO companies (
    id, email, password_hash, full_name, company_name, company_description,
    company_website, company_industry, company_size, company_location,
    is_active, is_verified, verification_status, created_at, updated_at
) VALUES (
    'test-unverified-company-id',
    'test-unverified@company.com',
    '$2a$10$N9qo8uLOickgx2ZMRZoMye.IjqMOeX.yL9rBvP/9S/1e3pL.JY2Gy', -- TestPassword123!
    'Admin Unverified',
    'Test Unverified Company',
    'An unverified test company',
    'https://unverified-company.test',
    'Technology',
    '11-50',
    'Bandung, Indonesia',
    TRUE,
    FALSE,
    'pending',
    NOW(),
    NOW()
) ON DUPLICATE KEY UPDATE email = email;

-- Suspended Company (password: TestPassword123!)
INSERT INTO companies (
    id, email, password_hash, full_name, company_name, company_description,
    company_website, company_industry, company_size, company_location,
    is_active, is_verified, verification_status, created_at, updated_at
) VALUES (
    'test-suspended-company-id',
    'test-suspended@company.com',
    '$2a$10$N9qo8uLOickgx2ZMRZoMye.IjqMOeX.yL9rBvP/9S/1e3pL.JY2Gy', -- TestPassword123!
    'Admin Suspended',
    'Test Suspended Company',
    'A suspended test company',
    'https://suspended-company.test',
    'Technology',
    '1-10',
    'Surabaya, Indonesia',
    FALSE,
    TRUE,
    'verified',
    NOW(),
    NOW()
) ON DUPLICATE KEY UPDATE email = email;

-- Company with no quota (password: TestPassword123!)
INSERT INTO companies (
    id, email, password_hash, full_name, company_name, company_description,
    company_website, company_industry, company_size, company_location,
    is_active, is_verified, verification_status, created_at, updated_at
) VALUES (
    'test-noquota-company-id',
    'test-noquota@company.com',
    '$2a$10$N9qo8uLOickgx2ZMRZoMye.IjqMOeX.yL9rBvP/9S/1e3pL.JY2Gy', -- TestPassword123!
    'Admin No Quota',
    'Test No Quota Company',
    'A test company with exhausted quota',
    'https://noquota-company.test',
    'Technology',
    '11-50',
    'Yogyakarta, Indonesia',
    TRUE,
    TRUE,
    'verified',
    NOW(),
    NOW()
) ON DUPLICATE KEY UPDATE email = email;

-- ============================================
-- Test Quotas
-- ============================================

INSERT INTO company_quotas (
    id, company_id, free_quota, used_free_quota, paid_quota, price_per_job, created_at, updated_at
) VALUES 
    ('quota-verified', 'test-verified-company-id', 5, 2, 0, 150000, NOW(), NOW()),
    ('quota-unverified', 'test-unverified-company-id', 5, 0, 0, 150000, NOW(), NOW()),
    ('quota-suspended', 'test-suspended-company-id', 5, 1, 0, 150000, NOW(), NOW()),
    ('quota-noquota', 'test-noquota-company-id', 5, 5, 0, 150000, NOW(), NOW())
ON DUPLICATE KEY UPDATE company_id = company_id;

-- ============================================
-- Test Jobs
-- ============================================

-- Draft job
INSERT INTO jobs (
    id, company_id, title, description, requirements, location, type,
    experience_level, salary_min, salary_max, salary_currency, is_salary_visible,
    status, created_at, updated_at
) VALUES (
    'test-job-draft',
    'test-verified-company-id',
    'Draft Software Engineer',
    '<p>Draft job description</p>',
    '<p>Draft requirements</p>',
    'Jakarta, Indonesia',
    'full-time',
    'mid',
    15000000,
    25000000,
    'IDR',
    TRUE,
    'draft',
    NOW(),
    NOW()
) ON DUPLICATE KEY UPDATE id = id;

-- Active job
INSERT INTO jobs (
    id, company_id, title, description, requirements, location, type,
    experience_level, salary_min, salary_max, salary_currency, is_salary_visible,
    status, published_at, created_at, updated_at
) VALUES (
    'test-job-active',
    'test-verified-company-id',
    'Active Software Engineer',
    '<p>Active job description</p>',
    '<p>Active requirements</p>',
    'Jakarta, Indonesia',
    'full-time',
    'senior',
    20000000,
    35000000,
    'IDR',
    TRUE,
    'active',
    NOW(),
    NOW(),
    NOW()
) ON DUPLICATE KEY UPDATE id = id;

-- Closed job
INSERT INTO jobs (
    id, company_id, title, description, requirements, location, type,
    experience_level, salary_min, salary_max, salary_currency, is_salary_visible,
    status, published_at, closed_at, created_at, updated_at
) VALUES (
    'test-job-closed',
    'test-verified-company-id',
    'Closed Software Engineer',
    '<p>Closed job description</p>',
    '<p>Closed requirements</p>',
    'Bandung, Indonesia',
    'full-time',
    'junior',
    8000000,
    15000000,
    'IDR',
    TRUE,
    'closed',
    DATE_SUB(NOW(), INTERVAL 30 DAY),
    NOW(),
    NOW(),
    NOW()
) ON DUPLICATE KEY UPDATE id = id;

-- ============================================
-- Test Applications (for status transition tests)
-- ============================================

-- Application with submitted status
INSERT INTO applications (
    id, job_id, applicant_name, applicant_email, applicant_phone,
    resume_url, cover_letter, status, created_at, updated_at
) VALUES (
    'test-app-submitted',
    'test-job-active',
    'John Submitted',
    'john.submitted@test.com',
    '+6281234567001',
    'https://example.com/resume-submitted.pdf',
    'Cover letter for submitted application',
    'submitted',
    NOW(),
    NOW()
) ON DUPLICATE KEY UPDATE id = id;

-- Application with viewed status
INSERT INTO applications (
    id, job_id, applicant_name, applicant_email, applicant_phone,
    resume_url, cover_letter, status, created_at, updated_at
) VALUES (
    'test-app-viewed',
    'test-job-active',
    'Jane Viewed',
    'jane.viewed@test.com',
    '+6281234567002',
    'https://example.com/resume-viewed.pdf',
    'Cover letter for viewed application',
    'viewed',
    NOW(),
    NOW()
) ON DUPLICATE KEY UPDATE id = id;

-- Application with shortlisted status
INSERT INTO applications (
    id, job_id, applicant_name, applicant_email, applicant_phone,
    resume_url, cover_letter, status, created_at, updated_at
) VALUES (
    'test-app-shortlisted',
    'test-job-active',
    'Bob Shortlisted',
    'bob.shortlisted@test.com',
    '+6281234567003',
    'https://example.com/resume-shortlisted.pdf',
    'Cover letter for shortlisted application',
    'shortlisted',
    NOW(),
    NOW()
) ON DUPLICATE KEY UPDATE id = id;

-- Application with interview_scheduled status
INSERT INTO applications (
    id, job_id, applicant_name, applicant_email, applicant_phone,
    resume_url, cover_letter, status, created_at, updated_at
) VALUES (
    'test-app-interview-scheduled',
    'test-job-active',
    'Alice Interview',
    'alice.interview@test.com',
    '+6281234567004',
    'https://example.com/resume-interview.pdf',
    'Cover letter for interview scheduled application',
    'interview_scheduled',
    NOW(),
    NOW()
) ON DUPLICATE KEY UPDATE id = id;

-- Application with interview_completed status
INSERT INTO applications (
    id, job_id, applicant_name, applicant_email, applicant_phone,
    resume_url, cover_letter, status, created_at, updated_at
) VALUES (
    'test-app-interview-completed',
    'test-job-active',
    'Charlie Completed',
    'charlie.completed@test.com',
    '+6281234567005',
    'https://example.com/resume-completed.pdf',
    'Cover letter for interview completed application',
    'interview_completed',
    NOW(),
    NOW()
) ON DUPLICATE KEY UPDATE id = id;

-- Application with offer_sent status
INSERT INTO applications (
    id, job_id, applicant_name, applicant_email, applicant_phone,
    resume_url, cover_letter, status, created_at, updated_at
) VALUES (
    'test-app-offer-sent',
    'test-job-active',
    'David Offer',
    'david.offer@test.com',
    '+6281234567006',
    'https://example.com/resume-offer.pdf',
    'Cover letter for offer sent application',
    'offer_sent',
    NOW(),
    NOW()
) ON DUPLICATE KEY UPDATE id = id;

-- Application with offer_accepted status
INSERT INTO applications (
    id, job_id, applicant_name, applicant_email, applicant_phone,
    resume_url, cover_letter, status, created_at, updated_at
) VALUES (
    'test-app-offer-accepted',
    'test-job-active',
    'Eve Accepted',
    'eve.accepted@test.com',
    '+6281234567007',
    'https://example.com/resume-accepted.pdf',
    'Cover letter for offer accepted application',
    'offer_accepted',
    NOW(),
    NOW()
) ON DUPLICATE KEY UPDATE id = id;

-- Application with hired status
INSERT INTO applications (
    id, job_id, applicant_name, applicant_email, applicant_phone,
    resume_url, cover_letter, status, created_at, updated_at
) VALUES (
    'test-app-hired',
    'test-job-active',
    'Frank Hired',
    'frank.hired@test.com',
    '+6281234567008',
    'https://example.com/resume-hired.pdf',
    'Cover letter for hired application',
    'hired',
    NOW(),
    NOW()
) ON DUPLICATE KEY UPDATE id = id;

-- Application with rejected status
INSERT INTO applications (
    id, job_id, applicant_name, applicant_email, applicant_phone,
    resume_url, cover_letter, status, created_at, updated_at
) VALUES (
    'test-app-rejected',
    'test-job-active',
    'Grace Rejected',
    'grace.rejected@test.com',
    '+6281234567009',
    'https://example.com/resume-rejected.pdf',
    'Cover letter for rejected application',
    'rejected',
    NOW(),
    NOW()
) ON DUPLICATE KEY UPDATE id = id;

-- Application with withdrawn status
INSERT INTO applications (
    id, job_id, applicant_name, applicant_email, applicant_phone,
    resume_url, cover_letter, status, created_at, updated_at
) VALUES (
    'test-app-withdrawn',
    'test-job-active',
    'Henry Withdrawn',
    'henry.withdrawn@test.com',
    '+6281234567010',
    'https://example.com/resume-withdrawn.pdf',
    'Cover letter for withdrawn application',
    'withdrawn',
    NOW(),
    NOW()
) ON DUPLICATE KEY UPDATE id = id;
