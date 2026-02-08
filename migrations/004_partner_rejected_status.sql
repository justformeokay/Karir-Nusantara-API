-- =============================================
-- Migration: Add rejected status for partners
-- Version: 004
-- Date: 2026-02-08
-- Description: Add 'rejected' status to referral_partners status enum
-- =============================================

-- Modify the status enum to include 'rejected'
ALTER TABLE `referral_partners` 
MODIFY `status` enum('active','inactive','pending','suspended','rejected') NOT NULL DEFAULT 'pending';

-- Add index for faster status filtering
CREATE INDEX IF NOT EXISTS `idx_partner_status_created` ON `referral_partners` (`status`, `created_at` DESC);
