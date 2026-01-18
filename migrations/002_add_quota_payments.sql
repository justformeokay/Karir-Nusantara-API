-- Migration: Add company quota and payments tables
-- Description: Adds quota tracking and payment system for job postings

-- ============================================
-- COMPANY QUOTAS TABLE
-- ============================================
CREATE TABLE IF NOT EXISTS company_quotas (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    company_id BIGINT UNSIGNED NOT NULL UNIQUE,
    free_quota_used INT NOT NULL DEFAULT 0,
    paid_quota INT NOT NULL DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (company_id) REFERENCES users(id) ON DELETE CASCADE,
    INDEX idx_company_quotas_company_id (company_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================
-- PAYMENTS TABLE
-- ============================================
CREATE TABLE IF NOT EXISTS payments (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    company_id BIGINT UNSIGNED NOT NULL,
    job_id BIGINT UNSIGNED NULL,
    amount BIGINT NOT NULL DEFAULT 30000,
    proof_image_url VARCHAR(500) NULL,
    status ENUM('pending', 'confirmed', 'rejected') NOT NULL DEFAULT 'pending',
    note TEXT NULL,
    confirmed_by_id BIGINT UNSIGNED NULL,
    submitted_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    confirmed_at TIMESTAMP NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (company_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (job_id) REFERENCES jobs(id) ON DELETE SET NULL,
    FOREIGN KEY (confirmed_by_id) REFERENCES users(id) ON DELETE SET NULL,
    INDEX idx_payments_company_id (company_id),
    INDEX idx_payments_status (status),
    INDEX idx_payments_submitted_at (submitted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================
-- UPDATE USERS TABLE FOR COMPANY STATUS
-- ============================================
-- Add company_status column for company verification workflow
ALTER TABLE users 
ADD COLUMN IF NOT EXISTS company_status ENUM('pending', 'verified', 'rejected', 'suspended') 
DEFAULT 'pending' AFTER is_verified;

-- Update existing companies: set verified if is_verified = true
UPDATE users 
SET company_status = 'verified' 
WHERE role = 'company' AND is_verified = 1;

UPDATE users 
SET company_status = 'pending' 
WHERE role = 'company' AND is_verified = 0;
