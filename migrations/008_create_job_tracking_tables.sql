-- Migration: Create job tracking tables
-- Created: 2026-01-23
-- Purpose: Track job views, shares for analytics

-- --------------------------------------------------------
-- Table structure for `job_views`
-- Tracks unique views per job per applicant (user_id)
-- --------------------------------------------------------

CREATE TABLE IF NOT EXISTS `job_views` (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT,
  `job_id` bigint(20) UNSIGNED NOT NULL,
  `user_id` bigint(20) UNSIGNED NOT NULL COMMENT 'Applicant user_id who viewed',
  `viewed_at` timestamp NOT NULL DEFAULT current_timestamp(),
  PRIMARY KEY (`id`),
  UNIQUE KEY `unique_job_user_view` (`job_id`, `user_id`),
  KEY `idx_job_views_job_id` (`job_id`),
  KEY `idx_job_views_user_id` (`user_id`),
  CONSTRAINT `fk_job_views_job` FOREIGN KEY (`job_id`) REFERENCES `jobs` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_job_views_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- --------------------------------------------------------
-- Table structure for `job_shares`
-- Tracks share count per job (incremental counter)
-- --------------------------------------------------------

CREATE TABLE IF NOT EXISTS `job_shares` (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT,
  `job_id` bigint(20) UNSIGNED NOT NULL,
  `user_id` bigint(20) UNSIGNED DEFAULT NULL COMMENT 'Applicant user_id who shared (optional)',
  `platform` varchar(50) DEFAULT NULL COMMENT 'Platform: whatsapp, telegram, facebook, twitter, copy_link, etc',
  `shared_at` timestamp NOT NULL DEFAULT current_timestamp(),
  PRIMARY KEY (`id`),
  KEY `idx_job_shares_job_id` (`job_id`),
  KEY `idx_job_shares_user_id` (`user_id`),
  CONSTRAINT `fk_job_shares_job` FOREIGN KEY (`job_id`) REFERENCES `jobs` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_job_shares_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- --------------------------------------------------------
-- Add shares_count column to jobs table
-- --------------------------------------------------------

ALTER TABLE `jobs` 
ADD COLUMN `shares_count` int(10) UNSIGNED NOT NULL DEFAULT 0 
AFTER `applications_count`;
