-- Migration: 007_applicant_profiles.sql
-- Description: Create applicant_profiles and applicant_documents tables for comprehensive applicant data
-- Created: 2026-01-23

-- =====================================================
-- Table: applicant_profiles
-- Stores comprehensive profile information for job seekers
-- =====================================================

CREATE TABLE IF NOT EXISTS `applicant_profiles` (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT,
  `user_id` bigint(20) UNSIGNED NOT NULL,
  
  -- Personal Information
  `date_of_birth` date DEFAULT NULL,
  `gender` enum('male','female','other','prefer_not_to_say') DEFAULT NULL,
  `nationality` varchar(100) DEFAULT 'Indonesia',
  `marital_status` enum('single','married','divorced','widowed') DEFAULT NULL,
  
  -- Identity (optional, for verification)
  `nik` varchar(20) DEFAULT NULL COMMENT 'Nomor KTP',
  
  -- Address
  `address` text DEFAULT NULL,
  `city` varchar(100) DEFAULT NULL,
  `province` varchar(100) DEFAULT NULL,
  `postal_code` varchar(10) DEFAULT NULL,
  `country` varchar(100) DEFAULT 'Indonesia',
  
  -- Professional Links
  `linkedin_url` varchar(500) DEFAULT NULL,
  `github_url` varchar(500) DEFAULT NULL,
  `portfolio_url` varchar(500) DEFAULT NULL,
  `personal_website` varchar(500) DEFAULT NULL,
  
  -- Bio/Summary
  `professional_summary` text DEFAULT NULL,
  `headline` varchar(255) DEFAULT NULL COMMENT 'e.g., Senior Software Engineer',
  
  -- Job Preferences
  `expected_salary_min` bigint(20) UNSIGNED DEFAULT NULL,
  `expected_salary_max` bigint(20) UNSIGNED DEFAULT NULL,
  `preferred_job_types` longtext DEFAULT NULL COMMENT 'JSON array: ["full_time","remote"]',
  `preferred_locations` longtext DEFAULT NULL COMMENT 'JSON array of cities',
  `available_from` date DEFAULT NULL,
  `willing_to_relocate` tinyint(1) DEFAULT 0,
  
  -- Profile Completeness
  `profile_completeness` int(10) UNSIGNED DEFAULT 0,
  
  `created_at` timestamp NOT NULL DEFAULT current_timestamp(),
  `updated_at` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_user_id` (`user_id`),
  KEY `idx_city` (`city`),
  KEY `idx_province` (`province`),
  KEY `idx_profile_completeness` (`profile_completeness`),
  CONSTRAINT `fk_applicant_profiles_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;


-- =====================================================
-- Table: applicant_documents
-- Stores uploaded documents (CV, certificates, etc.)
-- =====================================================

CREATE TABLE IF NOT EXISTS `applicant_documents` (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT,
  `user_id` bigint(20) UNSIGNED NOT NULL,
  
  `document_type` enum('cv_uploaded','cv_generated','certificate','transcript','portfolio','ktp','other') NOT NULL,
  `document_name` varchar(255) NOT NULL COMMENT 'Original filename',
  `document_url` varchar(500) NOT NULL COMMENT 'Path to file',
  `file_size` int(10) UNSIGNED DEFAULT NULL COMMENT 'Size in bytes',
  `mime_type` varchar(100) DEFAULT NULL COMMENT 'e.g., application/pdf',
  `is_primary` tinyint(1) DEFAULT 0 COMMENT 'Is this the primary CV?',
  `description` text DEFAULT NULL COMMENT 'Optional description of the document',
  
  `uploaded_at` timestamp NOT NULL DEFAULT current_timestamp(),
  `expires_at` timestamp NULL DEFAULT NULL,
  
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_document_type` (`document_type`),
  KEY `idx_is_primary` (`is_primary`),
  CONSTRAINT `fk_applicant_documents_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;


-- =====================================================
-- Add indexes for better query performance
-- =====================================================

-- Index for searching applicants by location
ALTER TABLE `applicant_profiles` ADD INDEX `idx_location` (`city`, `province`);

-- Index for filtering by salary expectations
ALTER TABLE `applicant_profiles` ADD INDEX `idx_salary_range` (`expected_salary_min`, `expected_salary_max`);
