-- Migration: Add package_id and quota_amount to payments table
-- For Hybrid Top-Up Model

-- Add new columns if they don't exist
ALTER TABLE payments 
ADD COLUMN IF NOT EXISTS package_id VARCHAR(50) NULL AFTER job_id,
ADD COLUMN IF NOT EXISTS quota_amount INT NOT NULL DEFAULT 1 AFTER package_id;

-- Update existing payments to have quota_amount = 1
UPDATE payments SET quota_amount = 1 WHERE quota_amount = 0 OR quota_amount IS NULL;

-- Add index for package_id
CREATE INDEX IF NOT EXISTS idx_payments_package_id ON payments(package_id);
