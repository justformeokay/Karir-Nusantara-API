-- =============================================
-- Migration: Partner Referral System
-- Version: 003
-- Date: 2026-02-04
-- Description: Add tables for partner referral and commission tracking
-- =============================================

-- =============================================
-- STEP 1: Modify users table to add 'partner' role
-- =============================================
ALTER TABLE `users` 
MODIFY `role` enum('job_seeker','company','admin','partner') NOT NULL DEFAULT 'job_seeker';

-- =============================================
-- STEP 2: Create referral_partners table
-- Stores partner account information and bank details
-- =============================================
CREATE TABLE `referral_partners` (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT,
  `user_id` bigint(20) UNSIGNED NOT NULL COMMENT 'Link to users table',
  `referral_code` varchar(20) NOT NULL COMMENT 'Unique referral code e.g., AHMAD2024',
  `commission_rate` decimal(5,2) NOT NULL DEFAULT 40.00 COMMENT 'Commission percentage (40%)',
  `status` enum('active','inactive','pending','suspended') NOT NULL DEFAULT 'pending',
  
  -- Bank account information
  `bank_name` varchar(100) DEFAULT NULL,
  `bank_account_number` varchar(50) DEFAULT NULL,
  `bank_account_holder` varchar(255) DEFAULT NULL,
  `is_bank_verified` tinyint(1) NOT NULL DEFAULT 0,
  
  -- Cached statistics (updated by triggers/application)
  `total_referrals` int(11) NOT NULL DEFAULT 0 COMMENT 'Total companies referred',
  `total_commission` bigint(20) NOT NULL DEFAULT 0 COMMENT 'Total commission earned (lifetime)',
  `available_balance` bigint(20) NOT NULL DEFAULT 0 COMMENT 'Balance ready for payout',
  `pending_balance` bigint(20) NOT NULL DEFAULT 0 COMMENT 'Commission pending approval',
  `paid_amount` bigint(20) NOT NULL DEFAULT 0 COMMENT 'Total amount paid out',
  
  -- Admin approval
  `approved_by` bigint(20) UNSIGNED DEFAULT NULL,
  `approved_at` timestamp NULL DEFAULT NULL,
  `notes` text DEFAULT NULL COMMENT 'Admin notes',
  
  -- Timestamps
  `created_at` timestamp NOT NULL DEFAULT current_timestamp(),
  `updated_at` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_user_id` (`user_id`),
  UNIQUE KEY `uk_referral_code` (`referral_code`),
  KEY `idx_status` (`status`),
  KEY `idx_referral_code` (`referral_code`),
  KEY `idx_created_at` (`created_at`),
  
  CONSTRAINT `fk_partner_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_partner_approved_by` FOREIGN KEY (`approved_by`) REFERENCES `users` (`id`) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- =============================================
-- STEP 3: Create partner_referrals table
-- Links partners to the companies they referred
-- =============================================
CREATE TABLE `partner_referrals` (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT,
  `partner_id` bigint(20) UNSIGNED NOT NULL COMMENT 'referral_partners.id',
  `company_id` bigint(20) UNSIGNED NOT NULL COMMENT 'companies.id',
  `referral_code_used` varchar(20) NOT NULL COMMENT 'The code used at registration',
  `registered_at` timestamp NOT NULL DEFAULT current_timestamp(),
  `is_verified` tinyint(1) NOT NULL DEFAULT 0 COMMENT 'Company account has been verified',
  `first_payment_at` timestamp NULL DEFAULT NULL COMMENT 'When company made first purchase',
  `notes` text DEFAULT NULL,
  
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_partner_company` (`partner_id`, `company_id`),
  UNIQUE KEY `uk_company_id` (`company_id`) COMMENT 'A company can only have one referrer',
  KEY `idx_partner_id` (`partner_id`),
  KEY `idx_registered_at` (`registered_at`),
  
  CONSTRAINT `fk_referral_partner` FOREIGN KEY (`partner_id`) REFERENCES `referral_partners` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_referral_company` FOREIGN KEY (`company_id`) REFERENCES `companies` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- =============================================
-- STEP 4: Create partner_commissions table
-- Tracks commission earned per transaction/payment
-- =============================================
CREATE TABLE `partner_commissions` (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT,
  `partner_id` bigint(20) UNSIGNED NOT NULL,
  `referral_id` bigint(20) UNSIGNED NOT NULL COMMENT 'partner_referrals.id',
  `payment_id` bigint(20) UNSIGNED NOT NULL COMMENT 'payments.id',
  `company_id` bigint(20) UNSIGNED NOT NULL,
  
  -- Transaction details
  `transaction_amount` bigint(20) NOT NULL COMMENT 'Original payment amount (IDR)',
  `commission_rate` decimal(5,2) NOT NULL COMMENT 'Rate at time of transaction (e.g., 40.00)',
  `commission_amount` bigint(20) NOT NULL COMMENT 'Calculated commission (IDR)',
  `job_quota` int(11) NOT NULL COMMENT 'Number of job posts purchased',
  
  -- Status tracking
  `status` enum('pending','approved','paid','cancelled') NOT NULL DEFAULT 'pending',
  `approved_by` bigint(20) UNSIGNED DEFAULT NULL,
  `approved_at` timestamp NULL DEFAULT NULL,
  `paid_at` timestamp NULL DEFAULT NULL,
  `payout_id` bigint(20) UNSIGNED DEFAULT NULL COMMENT 'Link to partner_payouts when paid',
  `notes` text DEFAULT NULL,
  
  -- Timestamps
  `created_at` timestamp NOT NULL DEFAULT current_timestamp(),
  `updated_at` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_payment_id` (`payment_id`) COMMENT 'One commission record per payment',
  KEY `idx_partner_id` (`partner_id`),
  KEY `idx_referral_id` (`referral_id`),
  KEY `idx_company_id` (`company_id`),
  KEY `idx_status` (`status`),
  KEY `idx_created_at` (`created_at`),
  KEY `idx_payout_id` (`payout_id`),
  
  CONSTRAINT `fk_commission_partner` FOREIGN KEY (`partner_id`) REFERENCES `referral_partners` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_commission_referral` FOREIGN KEY (`referral_id`) REFERENCES `partner_referrals` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_commission_payment` FOREIGN KEY (`payment_id`) REFERENCES `payments` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_commission_company` FOREIGN KEY (`company_id`) REFERENCES `companies` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_commission_approved_by` FOREIGN KEY (`approved_by`) REFERENCES `users` (`id`) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- =============================================
-- STEP 5: Create partner_payouts table
-- Tracks payout history and status
-- =============================================
CREATE TABLE `partner_payouts` (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT,
  `partner_id` bigint(20) UNSIGNED NOT NULL,
  `amount` bigint(20) NOT NULL COMMENT 'Payout amount (IDR)',
  
  -- Bank details snapshot at time of payout
  `bank_name` varchar(100) NOT NULL,
  `bank_account_number` varchar(50) NOT NULL,
  `bank_account_holder` varchar(255) NOT NULL,
  
  -- Status tracking
  `status` enum('pending','processing','completed','failed','cancelled') NOT NULL DEFAULT 'pending',
  `transfer_ref` varchar(100) DEFAULT NULL COMMENT 'Bank transfer reference number',
  `failure_reason` text DEFAULT NULL,
  
  -- Processing info
  `requested_at` timestamp NOT NULL DEFAULT current_timestamp(),
  `processed_by` bigint(20) UNSIGNED DEFAULT NULL,
  `processed_at` timestamp NULL DEFAULT NULL,
  `completed_at` timestamp NULL DEFAULT NULL,
  `notes` text DEFAULT NULL,
  
  -- Timestamps
  `created_at` timestamp NOT NULL DEFAULT current_timestamp(),
  `updated_at` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  
  PRIMARY KEY (`id`),
  KEY `idx_partner_id` (`partner_id`),
  KEY `idx_status` (`status`),
  KEY `idx_created_at` (`created_at`),
  KEY `idx_processed_at` (`processed_at`),
  
  CONSTRAINT `fk_payout_partner` FOREIGN KEY (`partner_id`) REFERENCES `referral_partners` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_payout_processed_by` FOREIGN KEY (`processed_by`) REFERENCES `users` (`id`) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- =============================================
-- STEP 6: Modify companies table
-- Add referral tracking columns
-- =============================================
ALTER TABLE `companies`
ADD COLUMN `referred_by_partner_id` bigint(20) UNSIGNED DEFAULT NULL COMMENT 'referral_partners.id',
ADD COLUMN `referral_code_used` varchar(20) DEFAULT NULL COMMENT 'Referral code used at registration',
ADD KEY `idx_referred_by_partner` (`referred_by_partner_id`),
ADD CONSTRAINT `fk_company_referrer` FOREIGN KEY (`referred_by_partner_id`) REFERENCES `referral_partners` (`id`) ON DELETE SET NULL;

-- =============================================
-- STEP 7: Create indexes for performance
-- =============================================

-- Index for finding companies by referral code
CREATE INDEX `idx_companies_referral_code` ON `companies` (`referral_code_used`);

-- Index for commission reporting by date range
CREATE INDEX `idx_commissions_date_status` ON `partner_commissions` (`created_at`, `status`);

-- Index for partner payout history
CREATE INDEX `idx_payouts_partner_date` ON `partner_payouts` (`partner_id`, `created_at` DESC);

-- =============================================
-- STEP 8: Insert sample partner data (for testing)
-- =============================================

-- Create a test partner user
INSERT INTO `users` (`full_name`, `email`, `password_hash`, `role`, `is_verified`, `is_active`, `created_at`, `updated_at`)
VALUES (
  'Ahmad Pratama',
  'ahmad.pratama@email.com',
  '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', -- password: "password"
  'partner',
  1,
  1,
  NOW(),
  NOW()
);

-- Create partner record
INSERT INTO `referral_partners` (
  `user_id`,
  `referral_code`,
  `commission_rate`,
  `status`,
  `bank_name`,
  `bank_account_number`,
  `bank_account_holder`,
  `is_bank_verified`,
  `total_referrals`,
  `total_commission`,
  `available_balance`,
  `pending_balance`,
  `paid_amount`,
  `approved_at`,
  `created_at`
)
SELECT 
  `id`,
  'AHMAD2024',
  40.00,
  'active',
  'Bank Central Asia',
  '1234567890',
  'Ahmad Pratama',
  1,
  8,
  47850000,
  12500000,
  0,
  35350000,
  NOW(),
  '2024-01-15 00:00:00'
FROM `users` 
WHERE `email` = 'ahmad.pratama@email.com';

-- =============================================
-- STEP 9: Create triggers for automatic updates
-- =============================================

-- Trigger: Update partner statistics when commission is added
DELIMITER //
CREATE TRIGGER `after_commission_insert`
AFTER INSERT ON `partner_commissions`
FOR EACH ROW
BEGIN
  UPDATE `referral_partners`
  SET 
    `total_commission` = `total_commission` + NEW.commission_amount,
    `pending_balance` = CASE 
      WHEN NEW.status = 'pending' THEN `pending_balance` + NEW.commission_amount
      ELSE `pending_balance`
    END,
    `available_balance` = CASE 
      WHEN NEW.status = 'approved' THEN `available_balance` + NEW.commission_amount
      ELSE `available_balance`
    END,
    `updated_at` = NOW()
  WHERE `id` = NEW.partner_id;
END//
DELIMITER ;

-- Trigger: Update partner statistics when commission status changes
DELIMITER //
CREATE TRIGGER `after_commission_update`
AFTER UPDATE ON `partner_commissions`
FOR EACH ROW
BEGIN
  -- When status changes from pending to approved
  IF OLD.status = 'pending' AND NEW.status = 'approved' THEN
    UPDATE `referral_partners`
    SET 
      `pending_balance` = `pending_balance` - NEW.commission_amount,
      `available_balance` = `available_balance` + NEW.commission_amount,
      `updated_at` = NOW()
    WHERE `id` = NEW.partner_id;
  END IF;
  
  -- When status changes from approved to paid
  IF OLD.status = 'approved' AND NEW.status = 'paid' THEN
    UPDATE `referral_partners`
    SET 
      `available_balance` = `available_balance` - NEW.commission_amount,
      `paid_amount` = `paid_amount` + NEW.commission_amount,
      `updated_at` = NOW()
    WHERE `id` = NEW.partner_id;
  END IF;
  
  -- When commission is cancelled
  IF NEW.status = 'cancelled' AND OLD.status IN ('pending', 'approved') THEN
    UPDATE `referral_partners`
    SET 
      `total_commission` = `total_commission` - OLD.commission_amount,
      `pending_balance` = CASE 
        WHEN OLD.status = 'pending' THEN `pending_balance` - OLD.commission_amount
        ELSE `pending_balance`
      END,
      `available_balance` = CASE 
        WHEN OLD.status = 'approved' THEN `available_balance` - OLD.commission_amount
        ELSE `available_balance`
      END,
      `updated_at` = NOW()
    WHERE `id` = NEW.partner_id;
  END IF;
END//
DELIMITER ;

-- Trigger: Update referral count when new referral is added
DELIMITER //
CREATE TRIGGER `after_referral_insert`
AFTER INSERT ON `partner_referrals`
FOR EACH ROW
BEGIN
  UPDATE `referral_partners`
  SET 
    `total_referrals` = `total_referrals` + 1,
    `updated_at` = NOW()
  WHERE `id` = NEW.partner_id;
END//
DELIMITER ;

-- =============================================
-- STEP 10: Create views for reporting
-- =============================================

-- View: Partner dashboard statistics
CREATE VIEW `v_partner_dashboard_stats` AS
SELECT 
  rp.id AS partner_id,
  rp.user_id,
  rp.referral_code,
  rp.total_referrals AS total_companies,
  COUNT(DISTINCT pc.id) AS total_transactions,
  COALESCE(rp.total_commission, 0) AS total_commission,
  COALESCE(rp.available_balance, 0) AS available_balance,
  COALESCE(rp.paid_amount, 0) AS paid_commission,
  COALESCE(rp.pending_balance, 0) AS pending_commission
FROM `referral_partners` rp
LEFT JOIN `partner_commissions` pc ON pc.partner_id = rp.id
GROUP BY rp.id;

-- View: Monthly partner statistics (for charts)
CREATE VIEW `v_partner_monthly_stats` AS
SELECT 
  pc.partner_id,
  DATE_FORMAT(pc.created_at, '%Y-%m') AS month_year,
  DATE_FORMAT(pc.created_at, '%b') AS month_name,
  SUM(pc.commission_amount) AS total_commission,
  COUNT(DISTINCT pr.company_id) AS companies_count
FROM `partner_commissions` pc
JOIN `partner_referrals` pr ON pr.id = pc.referral_id
WHERE pc.status IN ('approved', 'paid')
GROUP BY pc.partner_id, DATE_FORMAT(pc.created_at, '%Y-%m')
ORDER BY month_year DESC;

-- View: Partner company details
CREATE VIEW `v_partner_companies` AS
SELECT 
  pr.partner_id,
  c.id AS company_id,
  c.company_name,
  c.created_at AS registration_date,
  CASE WHEN c.company_status = 'verified' THEN 'active' ELSE 'inactive' END AS status,
  COALESCE(jc.job_count, 0) AS total_job_posts,
  COALESCE(pc_sum.total_revenue, 0) AS total_revenue,
  COALESCE(pc_sum.total_commission, 0) AS commission_earned
FROM `partner_referrals` pr
JOIN `companies` c ON c.id = pr.company_id
LEFT JOIN (
  SELECT company_id, COUNT(*) AS job_count
  FROM jobs
  GROUP BY company_id
) jc ON jc.company_id = c.id
LEFT JOIN (
  SELECT 
    company_id,
    SUM(transaction_amount) AS total_revenue,
    SUM(commission_amount) AS total_commission
  FROM partner_commissions
  WHERE status IN ('approved', 'paid')
  GROUP BY company_id
) pc_sum ON pc_sum.company_id = c.id;

-- =============================================
-- MIGRATION COMPLETE
-- =============================================

-- Summary of changes:
-- 1. Added 'partner' role to users.role enum
-- 2. Created referral_partners table
-- 3. Created partner_referrals table
-- 4. Created partner_commissions table
-- 5. Created partner_payouts table
-- 6. Added referral columns to companies table
-- 7. Created indexes for performance
-- 8. Created sample partner for testing
-- 9. Created triggers for automatic balance updates
-- 10. Created views for reporting
