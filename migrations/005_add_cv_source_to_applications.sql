-- Migration: Add cv_source field to applications table
-- Purpose: Track which CV type was used (built vs uploaded)
-- Date: 2026-02-15

-- Add cv_source column to applications table
ALTER TABLE `applications` 
ADD COLUMN `cv_source` ENUM('built', 'uploaded') NOT NULL DEFAULT 'built' COMMENT 'Source of CV used for application' AFTER `cv_snapshot_id`,
ADD COLUMN `uploaded_document_id` BIGINT(20) UNSIGNED NULL COMMENT 'Reference to uploaded document if cv_source=uploaded' AFTER `cv_source`;

-- Add index for better query performance
ALTER TABLE `applications` 
ADD INDEX `idx_cv_source` (`cv_source`),
ADD INDEX `idx_uploaded_document_id` (`uploaded_document_id`);

-- Note: We don't add foreign key constraint to uploaded_document_id 
-- because documents can be deleted independently
