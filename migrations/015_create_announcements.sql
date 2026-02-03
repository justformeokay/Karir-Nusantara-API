-- ========================================================
-- Migration: Create Announcements Table
-- Description: Table for Notifications, Banners, and Information
-- Date: 2026-02-03
-- ========================================================

-- Table for Notifications, Banners, and Information
-- This unified table handles all three types of announcements
CREATE TABLE IF NOT EXISTS `announcements` (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT,
  `title` varchar(255) NOT NULL COMMENT 'Title of the announcement',
  `content` text NOT NULL COMMENT 'Content/body of the announcement',
  `type` enum('notification','banner','information') NOT NULL DEFAULT 'notification' COMMENT 'Type of announcement',
  `target_audience` enum('all','company','candidate','partner') NOT NULL DEFAULT 'all' COMMENT 'Target audience',
  `is_active` tinyint(1) NOT NULL DEFAULT 1 COMMENT 'Whether the announcement is active',
  `priority` int(11) NOT NULL DEFAULT 0 COMMENT 'Priority for ordering (higher = more important)',
  `start_date` timestamp NULL DEFAULT NULL COMMENT 'When to start showing the announcement',
  `end_date` timestamp NULL DEFAULT NULL COMMENT 'When to stop showing the announcement',
  `created_by` bigint(20) UNSIGNED DEFAULT NULL COMMENT 'Admin who created this',
  `updated_by` bigint(20) UNSIGNED DEFAULT NULL COMMENT 'Admin who last updated this',
  `created_at` timestamp NOT NULL DEFAULT current_timestamp(),
  `updated_at` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  PRIMARY KEY (`id`),
  KEY `idx_type` (`type`),
  KEY `idx_target_audience` (`target_audience`),
  KEY `idx_is_active` (`is_active`),
  KEY `idx_created_at` (`created_at`),
  KEY `idx_start_end_date` (`start_date`, `end_date`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Insert sample data
INSERT INTO `announcements` (`title`, `content`, `type`, `target_audience`, `is_active`, `priority`, `created_at`) VALUES
('Pemeliharaan Sistem', 'Sistem akan mengalami maintenance pada tanggal 15 Februari 2026 pukul 00:00 - 04:00 WIB. Mohon maaf atas ketidaknyamanannya.', 'notification', 'all', 1, 10, NOW()),
('Promo Khusus', 'Dapatkan diskon 20% untuk pembelian paket premium selama bulan Februari 2026!', 'notification', 'company', 1, 5, NOW()),
('Fitur Baru Tersedia', 'Kami telah meluncurkan fitur baru untuk membantu Anda dalam proses rekrutmen. Cek sekarang!', 'notification', 'company', 1, 3, NOW()),
('Selamat Datang di Karir Nusantara', 'Platform job portal terbaik di Indonesia. Temukan pekerjaan impian Anda bersama kami.', 'banner', 'all', 1, 10, NOW()),
('Update Kebijakan Privasi', 'Kami telah memperbarui kebijakan privasi kami. Silakan baca lebih lanjut untuk informasi lebih detail.', 'information', 'all', 1, 5, NOW()),
('Panduan Membuat CV yang Menarik', 'Tips dan trik membuat CV yang menarik perhatian perekrut. Pelajari cara menonjolkan kelebihan Anda.', 'information', 'candidate', 1, 3, NOW());
