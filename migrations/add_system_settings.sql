-- ============================================================
-- Migration: system_settings & quota_packages tables
-- Description: Enables dynamic configuration of platform rules
--              such as free quota limit, price per job, and
--              available quota purchase packages.
-- ============================================================

-- Disable strict mode for this session to allow AUTO_INCREMENT inserts
SET SESSION sql_mode = 'NO_ENGINE_SUBSTITUTION';

-- -----------------------------------------------------------
-- Table: system_settings
-- Key-value store for global platform configuration
-- -----------------------------------------------------------
CREATE TABLE IF NOT EXISTS `system_settings` (
  `id`          bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT,
  `key`         varchar(100)        NOT NULL,
  `value`       longtext            NOT NULL,
  `data_type`   enum('string','integer','decimal','boolean','json') NOT NULL DEFAULT 'string',
  `description` text                DEFAULT NULL,
  `updated_by`  bigint(20) UNSIGNED DEFAULT NULL COMMENT 'Admin user who last updated this setting',
  `created_at`  timestamp           NOT NULL DEFAULT current_timestamp(),
  `updated_at`  timestamp           NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_system_settings_key` (`key`),
  KEY `idx_system_settings_updated_by` (`updated_by`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
  COMMENT='Global platform configuration key-value pairs';

-- Seed initial values (match current hardcoded defaults)
INSERT INTO `system_settings` (`id`, `key`, `value`, `data_type`, `description`) VALUES
  (1, 'free_quota_limit', '3',    'integer', 'Jumlah kuota gratis untuk perusahaan baru yang mendaftar'),
  (2, 'price_per_job',    '20000','integer', 'Harga dasar per posting lowongan kerja (IDR)'),
  (3, 'currency',         'IDR',  'string',  'Mata uang yang digunakan pada platform')
ON DUPLICATE KEY UPDATE updated_at = updated_at;

ALTER TABLE `system_settings` AUTO_INCREMENT = 10;

-- -----------------------------------------------------------
-- Table: quota_packages
-- Available quota purchase bundles shown to companies
-- -----------------------------------------------------------
CREATE TABLE IF NOT EXISTS `quota_packages` (
  `id`            bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT,
  `package_id`    varchar(50)         NOT NULL COMMENT 'Unique string identifier, e.g. single, pack5, pack10',
  `name`          varchar(100)        NOT NULL,
  `quota`         int(11)             NOT NULL COMMENT 'Number of job postings purchased',
  `bonus_quota`   int(11)             NOT NULL DEFAULT 0 COMMENT 'Free bonus postings',
  `price`         bigint(20)          NOT NULL COMMENT 'Price in IDR',
  `description`   text                DEFAULT NULL,
  `is_best_value` tinyint(1)          NOT NULL DEFAULT 0,
  `is_active`     tinyint(1)          NOT NULL DEFAULT 1,
  `display_order` int(11)             NOT NULL DEFAULT 0,
  `created_by`    bigint(20) UNSIGNED DEFAULT NULL,
  `updated_by`    bigint(20) UNSIGNED DEFAULT NULL,
  `created_at`    timestamp           NOT NULL DEFAULT current_timestamp(),
  `updated_at`    timestamp           NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_quota_packages_package_id` (`package_id`),
  KEY `idx_quota_packages_is_active` (`is_active`),
  KEY `idx_quota_packages_display_order` (`display_order`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
  COMMENT='Purchasable quota bundles available to companies';

-- Seed initial packages (match existing hardcoded packages)
INSERT INTO `quota_packages` (`id`, `package_id`, `name`, `quota`, `bonus_quota`, `price`, `description`, `is_best_value`, `is_active`, `display_order`) VALUES
  (1, 'single', '1 Posting',            1,  0, 20000,  'Bayar untuk 1 lowongan',              0, 1, 1),
  (2, 'pack5',  '5 Posting',            5,  0, 100000, 'Hemat waktu, beli 5 sekaligus',       0, 1, 2),
  (3, 'pack10', '10 Posting + 2 GRATIS',10, 2, 200000, 'Beli 10 dapat 12! Hemat Rp 40.000',  1, 1, 3),
  (4, 'pack20', '20 Posting + 5 GRATIS',20, 5, 400000, 'Beli 20 dapat 25! Hemat Rp 100.000', 0, 1, 4)
ON DUPLICATE KEY UPDATE updated_at = updated_at;

ALTER TABLE `quota_packages` AUTO_INCREMENT = 10;
