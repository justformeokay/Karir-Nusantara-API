-- Migration: Add Interview Scheduling Fields to application_timelines
-- Date: 2026-01-26
-- Description: Add columns for interview type, meeting link, platform, address, and contact info

ALTER TABLE `karir_nusantara`.`application_timelines`
ADD COLUMN `interview_type` ENUM('online', 'offline', 'whatsapp_notification') NULL DEFAULT NULL AFTER `scheduled_notes`,
ADD COLUMN `meeting_link` VARCHAR(500) NULL DEFAULT NULL AFTER `interview_type`,
ADD COLUMN `meeting_platform` VARCHAR(50) NULL DEFAULT NULL AFTER `meeting_link`,
ADD COLUMN `interview_address` TEXT NULL DEFAULT NULL AFTER `meeting_platform`,
ADD COLUMN `contact_person` VARCHAR(255) NULL DEFAULT NULL AFTER `interview_address`,
ADD COLUMN `contact_phone` VARCHAR(20) NULL DEFAULT NULL AFTER `contact_person`;

-- Verify the new columns were added
-- SELECT * FROM `karir_nusantara`.`application_timelines` LIMIT 1;
