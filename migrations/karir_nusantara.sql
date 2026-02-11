-- phpMyAdmin SQL Dump
-- version 5.2.1
-- https://www.phpmyadmin.net/
--
-- Host: localhost
-- Generation Time: Feb 11, 2026 at 04:12 PM
-- Server version: 10.4.28-MariaDB
-- PHP Version: 8.0.28

SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
START TRANSACTION;
SET time_zone = "+00:00";


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;

--
-- Database: `karir_nusantara`
--

-- --------------------------------------------------------

--
-- Table structure for table `announcements`
--

CREATE TABLE `announcements` (
  `id` bigint(20) UNSIGNED NOT NULL,
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
  `updated_at` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp()
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

--
-- Dumping data for table `announcements`
--

INSERT INTO `announcements` (`id`, `title`, `content`, `type`, `target_audience`, `is_active`, `priority`, `start_date`, `end_date`, `created_by`, `updated_by`, `created_at`, `updated_at`) VALUES
(1, 'Pemeliharaan Sistem', 'Sistem akan mengalami maintenance pada tanggal 15 Februari 2026 pukul 00:00 - 04:00 WIB. Mohon maaf atas ketidaknyamanannya.', 'notification', 'all', 1, 10, NULL, NULL, NULL, NULL, '2026-02-03 16:13:25', '2026-02-03 16:13:25'),
(2, 'Promo Khusus', 'Dapatkan diskon 20% untuk pembelian paket premium selama bulan Februari 2026!', 'notification', 'company', 1, 5, NULL, NULL, NULL, NULL, '2026-02-03 16:13:25', '2026-02-03 16:13:25'),
(3, 'Fitur Baru Tersedia', 'Kami telah meluncurkan fitur baru untuk membantu Anda dalam proses rekrutmen. Cek sekarang!', 'notification', 'company', 1, 3, NULL, NULL, NULL, NULL, '2026-02-03 16:13:25', '2026-02-03 16:13:25'),
(4, 'Selamat Datang di Karir Nusantara', 'Platform job portal terbaik di Indonesia. Temukan pekerjaan impian Anda bersama kami.', 'banner', 'all', 1, 10, NULL, NULL, NULL, NULL, '2026-02-03 16:13:25', '2026-02-03 16:13:25'),
(5, 'Update Kebijakan Privasi', 'Kami telah memperbarui kebijakan privasi kami. Silakan baca lebih lanjut untuk informasi lebih detail.', 'information', 'all', 1, 5, NULL, NULL, NULL, NULL, '2026-02-03 16:13:25', '2026-02-03 16:13:25'),
(6, 'Panduan Membuat CV yang Menarik', 'Tips dan trik membuat CV yang menarik perhatian perekrut. Pelajari cara menonjolkan kelebihan Anda.', 'information', 'candidate', 1, 3, NULL, NULL, NULL, NULL, '2026-02-03 16:13:25', '2026-02-03 16:13:25'),
(7, 'Pemeliharaan Sistem', 'Sistem akan mengalami maintenance pada tanggal 15 Februari 2026 pukul 00:00 - 04:00 WIB. Mohon maaf atas ketidaknyamanannya.', 'notification', 'all', 1, 10, NULL, NULL, NULL, NULL, '2026-02-03 16:13:59', '2026-02-03 16:13:59'),
(8, 'Promo Khusus', 'Dapatkan diskon 20% untuk pembelian paket premium selama bulan Februari 2026!', 'notification', 'company', 1, 5, NULL, NULL, NULL, NULL, '2026-02-03 16:13:59', '2026-02-03 16:13:59'),
(9, 'Fitur Baru Tersedia', 'Kami telah meluncurkan fitur baru untuk membantu Anda dalam proses rekrutmen. Cek sekarang!', 'notification', 'company', 1, 3, NULL, NULL, NULL, NULL, '2026-02-03 16:13:59', '2026-02-03 16:13:59'),
(10, 'Selamat Datang di Karir Nusantara', 'Platform job portal terbaik di Indonesia. Temukan pekerjaan impian Anda bersama kami.', 'banner', 'all', 1, 10, NULL, NULL, NULL, NULL, '2026-02-03 16:13:59', '2026-02-03 16:13:59'),
(11, 'Update Kebijakan Privasi', 'Kami telah memperbarui kebijakan privasi kami. Silakan baca lebih lanjut untuk informasi lebih detail.', 'information', 'all', 1, 5, NULL, NULL, NULL, NULL, '2026-02-03 16:13:59', '2026-02-03 16:13:59'),
(12, 'Panduan Membuat CV yang Menarik', 'Tips dan trik membuat CV yang menarik perhatian perekrut. Pelajari cara menonjolkan kelebihan Anda.', 'information', 'candidate', 1, 3, NULL, NULL, NULL, NULL, '2026-02-03 16:13:59', '2026-02-03 16:13:59'),
(13, 'Pemeliharaan Sistem', 'Sistem akan mengalami maintenance pada tanggal 15 Februari 2026 pukul 00:00 - 04:00 WIB. Mohon maaf atas ketidaknyamanannya.', 'notification', 'all', 1, 10, NULL, NULL, NULL, NULL, '2026-02-03 16:14:10', '2026-02-03 16:14:10'),
(14, 'Promo Khusus', 'Dapatkan diskon 20% untuk pembelian paket premium selama bulan Februari 2026!', 'notification', 'company', 1, 5, NULL, NULL, NULL, NULL, '2026-02-03 16:14:10', '2026-02-03 16:14:10'),
(15, 'Fitur Baru Tersedia', 'Kami telah meluncurkan fitur baru untuk membantu Anda dalam proses rekrutmen. Cek sekarang!', 'notification', 'company', 1, 3, NULL, NULL, NULL, NULL, '2026-02-03 16:14:10', '2026-02-03 16:14:10'),
(16, 'Selamat Datang di Karir Nusantara', 'Platform job portal terbaik di Indonesia. Temukan pekerjaan impian Anda bersama kami.', 'banner', 'all', 1, 10, NULL, NULL, NULL, NULL, '2026-02-03 16:14:10', '2026-02-03 16:14:10'),
(17, 'Update Kebijakan Privasi', 'Kami telah memperbarui kebijakan privasi kami. Silakan baca lebih lanjut untuk informasi lebih detail.', 'information', 'all', 1, 5, NULL, NULL, NULL, NULL, '2026-02-03 16:14:10', '2026-02-03 16:14:10'),
(18, 'Panduan Membuat CV yang Menarik', 'Tips dan trik membuat CV yang menarik perhatian perekrut. Pelajari cara menonjolkan kelebihan Anda.', 'information', 'candidate', 1, 3, NULL, NULL, NULL, NULL, '2026-02-03 16:14:10', '2026-02-03 16:14:10');

-- --------------------------------------------------------

--
-- Table structure for table `applicant_documents`
--

CREATE TABLE `applicant_documents` (
  `id` bigint(20) UNSIGNED NOT NULL,
  `user_id` bigint(20) UNSIGNED NOT NULL,
  `document_type` enum('cv_uploaded','cv_generated','certificate','transcript','portfolio','ktp','other') NOT NULL,
  `document_name` varchar(255) NOT NULL COMMENT 'Original filename',
  `document_url` varchar(500) NOT NULL COMMENT 'Path to file',
  `file_size` int(10) UNSIGNED DEFAULT NULL COMMENT 'Size in bytes',
  `mime_type` varchar(100) DEFAULT NULL COMMENT 'e.g., application/pdf',
  `is_primary` tinyint(1) DEFAULT 0 COMMENT 'Is this the primary CV?',
  `description` text DEFAULT NULL COMMENT 'Optional description of the document',
  `uploaded_at` timestamp NOT NULL DEFAULT current_timestamp(),
  `expires_at` timestamp NULL DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

--
-- Dumping data for table `applicant_documents`
--

INSERT INTO `applicant_documents` (`id`, `user_id`, `document_type`, `document_name`, `document_url`, `file_size`, `mime_type`, `is_primary`, `description`, `uploaded_at`, `expires_at`) VALUES
(1, 21, 'cv_uploaded', 'CV Saputra Budianto.pdf', '/docs/applicants/21/cv_uploaded_1769149689.pdf', 41120, 'application/pdf', 1, NULL, '2026-01-23 06:28:09', NULL);

-- --------------------------------------------------------

--
-- Table structure for table `applicant_profiles`
--

CREATE TABLE `applicant_profiles` (
  `id` bigint(20) UNSIGNED NOT NULL,
  `user_id` bigint(20) UNSIGNED NOT NULL,
  `date_of_birth` date DEFAULT NULL,
  `gender` enum('male','female','other','prefer_not_to_say') DEFAULT NULL,
  `nationality` varchar(100) DEFAULT 'Indonesia',
  `marital_status` enum('single','married','divorced','widowed') DEFAULT NULL,
  `nik` varchar(20) DEFAULT NULL COMMENT 'Nomor KTP',
  `address` text DEFAULT NULL,
  `city` varchar(100) DEFAULT NULL,
  `province` varchar(100) DEFAULT NULL,
  `postal_code` varchar(10) DEFAULT NULL,
  `country` varchar(100) DEFAULT 'Indonesia',
  `linkedin_url` varchar(500) DEFAULT NULL,
  `github_url` varchar(500) DEFAULT NULL,
  `portfolio_url` varchar(500) DEFAULT NULL,
  `personal_website` varchar(500) DEFAULT NULL,
  `professional_summary` text DEFAULT NULL,
  `headline` varchar(255) DEFAULT NULL COMMENT 'e.g., Senior Software Engineer',
  `expected_salary_min` bigint(20) UNSIGNED DEFAULT NULL,
  `expected_salary_max` bigint(20) UNSIGNED DEFAULT NULL,
  `preferred_job_types` longtext DEFAULT NULL COMMENT 'JSON array: ["full_time","remote"]',
  `preferred_locations` longtext DEFAULT NULL COMMENT 'JSON array of cities',
  `available_from` date DEFAULT NULL,
  `willing_to_relocate` tinyint(1) DEFAULT 0,
  `profile_completeness` int(10) UNSIGNED DEFAULT 0,
  `created_at` timestamp NOT NULL DEFAULT current_timestamp(),
  `updated_at` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp()
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

--
-- Dumping data for table `applicant_profiles`
--

INSERT INTO `applicant_profiles` (`id`, `user_id`, `date_of_birth`, `gender`, `nationality`, `marital_status`, `nik`, `address`, `city`, `province`, `postal_code`, `country`, `linkedin_url`, `github_url`, `portfolio_url`, `personal_website`, `professional_summary`, `headline`, `expected_salary_min`, `expected_salary_max`, `preferred_job_types`, `preferred_locations`, `available_from`, `willing_to_relocate`, `profile_completeness`, `created_at`, `updated_at`) VALUES
(1, 21, '1999-11-23', 'male', 'Indonesia', 'single', NULL, 'Perumahan Griya Bhayangkara blok i5/07 Desa Masangan Kulon, Kecamatan Sukodono, Kabupaten Sidoarjo, Kode Pos 61258', 'Sidoarjo', 'Jawa Timur', '61258', 'Indonesia', 'https://www.linkedin.com/in/saputra-budianto23/', 'https://github.com/justformeokay23@', NULL, NULL, 'Saya adalah seorang profesional yang berdedikasi di bidang Mobile Apps Developer dengan pengalaman selama 3 tahun dalam mengembangkan cross platform mobile apps. Saya merupakan lulusan dari Universitas Muhammdiyah Sidoarjo, di mana saya mengasah kemampuan analitis dan teknis yang menjadi fondasi karier saya saat ini.\n\nSepanjang perjalanan profesional saya, saya telah berhasil mengembangkan dan mendistribusikan beberapa aplikasi ke marketplace dan web publik. Saya dikenal sebagai pribadi yang adaptif, komunikatif, dan memiliki orientasi kuat terhadap hasil. Saya selalu antusias untuk mempelajari teknologi baru dan berkontribusi dalam proyek-proyek inovatif yang memberikan dampak positif bagi organisasi.', 'Mobile Apps Developer', 5000000, 6000000, '[\"full_time\"]', NULL, NULL, 0, 80, '2026-01-23 06:13:35', '2026-01-23 06:29:06'),
(2, 20, '2003-06-06', 'female', 'Indonesia', 'single', NULL, 'Dusun Sumber Pandan, Desa Bulusari, Kecamatan Gempol, Kabupaten Pasuruan, Provinsi Jawa Timur', 'Pasuruan', 'Jawa Timur', '61832', 'Indonesia', NULL, NULL, NULL, NULL, 'Saya memiliki pengalaman dibidang Sales Marketing Jasa pada salah satu perusahaan asuransi terkenal di Indonesia dan memiliki pengalaman kurang lebih 5 tahun.', 'Marketing Sales', 4500000, 6000000, '[\"full_time\",\"part_time\"]', NULL, '2026-02-23', 1, 75, '2026-01-23 11:34:28', '2026-01-23 11:35:10');

-- --------------------------------------------------------

--
-- Table structure for table `applications`
--

CREATE TABLE `applications` (
  `id` bigint(20) UNSIGNED NOT NULL,
  `user_id` bigint(20) UNSIGNED NOT NULL,
  `job_id` bigint(20) UNSIGNED NOT NULL,
  `cv_snapshot_id` bigint(20) UNSIGNED NOT NULL,
  `cover_letter` text DEFAULT NULL,
  `current_status` enum('submitted','viewed','shortlisted','interview_scheduled','interview_completed','assessment','offer_sent','offer_accepted','hired','rejected','withdrawn') NOT NULL DEFAULT 'submitted',
  `applied_at` timestamp NOT NULL DEFAULT current_timestamp(),
  `last_status_update` timestamp NOT NULL DEFAULT current_timestamp(),
  `created_at` timestamp NOT NULL DEFAULT current_timestamp(),
  `updated_at` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp()
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

--
-- Dumping data for table `applications`
--

INSERT INTO `applications` (`id`, `user_id`, `job_id`, `cv_snapshot_id`, `cover_letter`, `current_status`, `applied_at`, `last_status_update`, `created_at`, `updated_at`) VALUES
(1, 3, 1, 1, 'Dengan pengalaman 4 tahun sebagai Backend Developer, saya yakin dapat memberikan kontribusi signifikan di PT TechCorp Indonesia.', 'hired', '2026-01-17 02:51:28', '2026-01-17 02:55:46', '2026-01-17 02:51:28', '2026-01-17 02:55:46'),
(3, 20, 54, 3, NULL, 'interview_scheduled', '2026-01-23 13:48:42', '2026-01-26 04:47:59', '2026-01-23 13:48:42', '2026-01-26 04:47:59'),
(4, 21, 49, 4, 'Saya tertarik dengan posisi Graphic Designer ini. Saya memiliki pengalaman 3 tahun dalam design grafis.', 'submitted', '2026-01-23 14:33:27', '2026-01-23 14:33:27', '2026-01-23 14:33:27', '2026-01-23 14:33:27'),
(5, 21, 56, 5, 'Saya tertarik dengan posisi Software Engineer ini. Saya memiliki keahlian di Go dan TypeScript.', 'submitted', '2026-01-23 14:36:59', '2026-01-23 14:36:59', '2026-01-23 14:36:59', '2026-01-23 14:36:59'),
(6, 21, 37, 6, NULL, 'submitted', '2026-01-23 14:45:37', '2026-01-23 14:45:37', '2026-01-23 14:45:37', '2026-01-23 14:45:37'),
(7, 21, 38, 7, NULL, 'submitted', '2026-01-23 14:50:24', '2026-01-23 14:50:24', '2026-01-23 14:50:24', '2026-01-23 14:50:24'),
(8, 21, 40, 8, NULL, 'submitted', '2026-01-23 14:55:13', '2026-01-23 14:55:13', '2026-01-23 14:55:13', '2026-01-23 14:55:13'),
(9, 21, 42, 9, NULL, 'interview_scheduled', '2026-01-23 14:58:52', '2026-01-26 05:09:34', '2026-01-23 14:58:52', '2026-01-26 05:09:34'),
(10, 21, 43, 10, NULL, 'interview_scheduled', '2026-01-23 15:04:57', '2026-01-25 21:52:49', '2026-01-23 15:04:57', '2026-01-25 21:52:49'),
(11, 20, 56, 11, NULL, 'interview_scheduled', '2026-01-25 09:07:38', '2026-01-25 21:21:39', '2026-01-25 09:07:38', '2026-01-25 21:21:39'),
(12, 20, 57, 12, NULL, 'submitted', '2026-02-11 14:56:30', '2026-02-11 14:56:30', '2026-02-11 14:56:30', '2026-02-11 14:56:30');

-- --------------------------------------------------------

--
-- Table structure for table `application_timelines`
--

CREATE TABLE `application_timelines` (
  `id` bigint(20) UNSIGNED NOT NULL,
  `application_id` bigint(20) UNSIGNED NOT NULL,
  `status` enum('submitted','viewed','shortlisted','interview_scheduled','interview_completed','assessment','offer_sent','offer_accepted','hired','rejected','withdrawn') NOT NULL,
  `note` text DEFAULT NULL,
  `is_visible_to_applicant` tinyint(1) NOT NULL DEFAULT 1,
  `updated_by_type` enum('system','company','applicant') NOT NULL DEFAULT 'system',
  `updated_by_id` bigint(20) UNSIGNED DEFAULT NULL,
  `scheduled_at` timestamp NULL DEFAULT NULL,
  `scheduled_location` varchar(500) DEFAULT NULL,
  `scheduled_notes` text DEFAULT NULL,
  `interview_type` enum('online','offline','whatsapp_notification') DEFAULT NULL,
  `meeting_link` varchar(500) DEFAULT NULL,
  `meeting_platform` varchar(50) DEFAULT NULL,
  `interview_address` text DEFAULT NULL,
  `contact_person` varchar(255) DEFAULT NULL,
  `contact_phone` varchar(20) DEFAULT NULL,
  `created_at` timestamp NOT NULL DEFAULT current_timestamp()
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

--
-- Dumping data for table `application_timelines`
--

INSERT INTO `application_timelines` (`id`, `application_id`, `status`, `note`, `is_visible_to_applicant`, `updated_by_type`, `updated_by_id`, `scheduled_at`, `scheduled_location`, `scheduled_notes`, `interview_type`, `meeting_link`, `meeting_platform`, `interview_address`, `contact_person`, `contact_phone`, `created_at`) VALUES
(1, 1, 'submitted', 'Lamaran berhasil dikirim', 1, 'system', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, '2026-01-17 02:51:28'),
(2, 1, 'viewed', 'Melihat profil kandidat', 1, 'company', 2, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, '2026-01-17 02:53:22'),
(3, 1, 'shortlisted', 'Kandidat memenuhi kriteria', 1, 'company', 2, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, '2026-01-17 02:53:40'),
(4, 1, 'interview_scheduled', 'Interview tahap 1', 1, 'company', 2, '2026-01-20 03:00:00', 'Kantor Jakarta', NULL, NULL, NULL, NULL, NULL, NULL, NULL, '2026-01-17 02:53:40'),
(5, 1, 'interview_completed', 'Interview berhasil', 1, 'company', 2, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, '2026-01-17 02:54:09'),
(6, 1, 'offer_sent', 'Surat penawaran dikirim', 1, 'company', 2, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, '2026-01-17 02:54:58'),
(7, 1, 'offer_accepted', 'Kandidat menerima penawaran', 1, 'company', 2, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, '2026-01-17 02:55:35'),
(8, 1, 'hired', 'Selamat bergabung di PT TechCorp Indonesia!', 1, 'company', 2, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, '2026-01-17 02:55:46'),
(10, 3, 'submitted', 'Lamaran berhasil dikirim', 1, 'system', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, '2026-01-23 13:48:42'),
(11, 4, 'submitted', 'Lamaran berhasil dikirim', 1, 'system', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, '2026-01-23 14:33:27'),
(12, 5, 'submitted', 'Lamaran berhasil dikirim', 1, 'system', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, '2026-01-23 14:36:59'),
(13, 6, 'submitted', 'Lamaran berhasil dikirim', 1, 'system', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, '2026-01-23 14:45:37'),
(14, 7, 'submitted', 'Lamaran berhasil dikirim', 1, 'system', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, '2026-01-23 14:50:24'),
(15, 8, 'submitted', 'Lamaran berhasil dikirim', 1, 'system', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, '2026-01-23 14:55:13'),
(16, 9, 'submitted', 'Lamaran berhasil dikirim', 1, 'system', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, '2026-01-23 14:58:52'),
(17, 10, 'submitted', 'Lamaran berhasil dikirim', 1, 'system', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, '2026-01-23 15:04:57'),
(18, 11, 'submitted', 'Lamaran berhasil dikirim', 1, 'system', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, '2026-01-25 09:07:38'),
(19, 11, 'viewed', NULL, 1, 'company', 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, '2026-01-25 21:20:36'),
(20, 11, 'shortlisted', NULL, 1, 'company', 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, '2026-01-25 21:21:21'),
(21, 11, 'interview_scheduled', NULL, 1, 'company', 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, '2026-01-25 21:21:39'),
(22, 10, 'viewed', NULL, 1, 'company', 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, '2026-01-25 21:52:47'),
(23, 10, 'shortlisted', NULL, 1, 'company', 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, '2026-01-25 21:52:48'),
(24, 10, 'interview_scheduled', 'Test interview', 1, 'company', 1, '2026-01-30 10:00:00', NULL, NULL, 'offline', NULL, NULL, 'Jl. Test No. 123', 'John Doe', '081234567890', '2026-01-25 21:52:49'),
(25, 3, 'interview_scheduled', NULL, 1, 'company', 1, '2026-01-30 11:30:00', NULL, 'Berpakaian Rapi celana gelap', 'offline', NULL, NULL, 'Perumahan Griya Bhayangkara blok i5/07 Desa Masangan Kulon', 'Saputra Budianto', '0881036480285', '2026-01-26 04:47:59'),
(26, 9, 'interview_scheduled', 'Test Interview', 1, 'company', 1, '2025-01-20 10:00:00', NULL, NULL, 'online', 'https://zoom.us/test', 'Zoom', NULL, NULL, NULL, '2026-01-26 05:09:34'),
(27, 12, 'submitted', 'Lamaran berhasil dikirim', 1, 'system', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, '2026-02-11 14:56:30');

-- --------------------------------------------------------

--
-- Table structure for table `audit_logs`
--

CREATE TABLE `audit_logs` (
  `id` bigint(20) UNSIGNED NOT NULL,
  `user_id` bigint(20) UNSIGNED DEFAULT NULL,
  `action` varchar(100) NOT NULL,
  `entity_type` varchar(50) NOT NULL,
  `entity_id` bigint(20) UNSIGNED DEFAULT NULL,
  `old_values` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL CHECK (json_valid(`old_values`)),
  `new_values` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL CHECK (json_valid(`new_values`)),
  `ip_address` varchar(45) DEFAULT NULL,
  `user_agent` varchar(500) DEFAULT NULL,
  `created_at` timestamp NOT NULL DEFAULT current_timestamp()
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- --------------------------------------------------------

--
-- Table structure for table `chat_messages`
--

CREATE TABLE `chat_messages` (
  `id` bigint(20) UNSIGNED NOT NULL,
  `conversation_id` bigint(20) UNSIGNED NOT NULL,
  `sender_id` bigint(20) UNSIGNED NOT NULL,
  `sender_type` enum('company','admin') NOT NULL,
  `message` text NOT NULL,
  `attachment_url` varchar(500) DEFAULT NULL,
  `attachment_type` enum('image','audio') DEFAULT NULL,
  `attachment_filename` varchar(255) DEFAULT NULL,
  `is_read` tinyint(1) DEFAULT 0,
  `created_at` timestamp NOT NULL DEFAULT current_timestamp()
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

--
-- Dumping data for table `chat_messages`
--

INSERT INTO `chat_messages` (`id`, `conversation_id`, `sender_id`, `sender_type`, `message`, `attachment_url`, `attachment_type`, `attachment_filename`, `is_read`, `created_at`) VALUES
(1, 1, 7, 'company', 'Halo admin, saya coba posting lowongan gagal', NULL, NULL, NULL, 0, '2026-01-21 17:48:41'),
(2, 1, 1, 'admin', 'Halo, saya sudah cek. Mohon pastikan kuota posting Anda masih tersedia. Silakan cek di menu Dashboard.', NULL, NULL, NULL, 1, '2026-01-21 17:49:18'),
(3, 4, 7, 'company', 'Test message from E2E script', NULL, NULL, NULL, 0, '2026-01-22 01:39:48'),
(4, 4, 1, 'admin', 'Admin reply: Percakapan sudah diterima. Terima kasih!', NULL, NULL, NULL, 1, '2026-01-22 01:39:49'),
(5, 4, 7, 'company', 'heyy', NULL, NULL, NULL, 0, '2026-01-22 01:53:38'),
(6, 3, 7, 'company', 'testtt', NULL, NULL, NULL, 0, '2026-01-22 01:59:55'),
(7, 3, 7, 'company', '🎤 Pesan Suara', '/uploads/chat/7_1769048557.webm', 'audio', 'voice_1769048553734.webm', 0, '2026-01-22 02:22:37'),
(8, 3, 7, 'company', 'test bang', '/uploads/chat/7_1769048600.png', 'image', 'Screenshot 2026-01-21 at 6.20.01 PM.png', 0, '2026-01-22 02:23:20'),
(9, 3, 7, 'company', 'sdas', NULL, NULL, NULL, 0, '2026-01-22 02:42:59'),
(10, 4, 7, 'company', 'hello bang', NULL, NULL, NULL, 0, '2026-01-22 03:14:46'),
(11, 5, 7, 'company', 'hello bang', NULL, NULL, NULL, 1, '2026-02-02 13:20:52'),
(12, 6, 7, 'company', 'Permisi bang', NULL, NULL, NULL, 1, '2026-02-02 14:14:23'),
(13, 6, 1, 'admin', 'iya ada yang bisa dibantu?', NULL, NULL, NULL, 1, '2026-02-02 14:14:44'),
(14, 6, 7, 'company', '🎤 Pesan Suara', '/docs/chat/7_1770043583.webm', 'audio', 'voice_1770043581844.webm', 1, '2026-02-02 14:46:23'),
(15, 6, 1, 'admin', 'ini error ya?', '/docs/chat/1_1770043816.jpeg', 'image', 'WhatsApp Image 2026-02-02 at 11.30.32 AM.jpeg', 1, '2026-02-02 14:50:16'),
(16, 6, 1, 'admin', '🎤 Pesan Suara', '/docs/chat/1_1770043836.webm', 'audio', 'voice_1770043835401.webm', 1, '2026-02-02 14:50:36');

-- --------------------------------------------------------

--
-- Table structure for table `companies`
--

CREATE TABLE `companies` (
  `id` bigint(20) UNSIGNED NOT NULL,
  `user_id` bigint(20) UNSIGNED NOT NULL,
  `company_name` varchar(255) NOT NULL,
  `company_description` longtext DEFAULT NULL,
  `company_website` varchar(255) DEFAULT NULL,
  `company_logo_url` varchar(500) DEFAULT NULL,
  `company_industry` varchar(100) DEFAULT NULL,
  `company_size` varchar(50) DEFAULT NULL,
  `company_location` varchar(255) DEFAULT NULL,
  `company_phone` varchar(20) DEFAULT NULL,
  `company_email` varchar(255) DEFAULT NULL,
  `company_address` longtext DEFAULT NULL,
  `company_city` varchar(100) DEFAULT NULL,
  `company_province` varchar(100) DEFAULT NULL,
  `company_postal_code` varchar(20) DEFAULT NULL,
  `established_year` year(4) DEFAULT NULL,
  `employee_count` int(11) DEFAULT NULL,
  `company_status` enum('pending','verified','rejected','suspended') DEFAULT 'pending',
  `ktp_founder_url` varchar(500) DEFAULT NULL,
  `akta_pendirian_url` varchar(500) DEFAULT NULL,
  `npwp_url` varchar(500) DEFAULT NULL,
  `nib_url` varchar(500) DEFAULT NULL,
  `documents_verified_at` timestamp NULL DEFAULT NULL,
  `documents_verified_by` bigint(20) UNSIGNED DEFAULT NULL,
  `verification_notes` longtext DEFAULT NULL,
  `created_at` timestamp NOT NULL DEFAULT current_timestamp(),
  `updated_at` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  `deleted_at` timestamp NULL DEFAULT NULL,
  `referred_by_partner_id` bigint(20) UNSIGNED DEFAULT NULL COMMENT 'referral_partners.id',
  `referral_code_used` varchar(20) DEFAULT NULL COMMENT 'Referral code used at registration'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

--
-- Dumping data for table `companies`
--

INSERT INTO `companies` (`id`, `user_id`, `company_name`, `company_description`, `company_website`, `company_logo_url`, `company_industry`, `company_size`, `company_location`, `company_phone`, `company_email`, `company_address`, `company_city`, `company_province`, `company_postal_code`, `established_year`, `employee_count`, `company_status`, `ktp_founder_url`, `akta_pendirian_url`, `npwp_url`, `nib_url`, `documents_verified_at`, `documents_verified_by`, `verification_notes`, `created_at`, `updated_at`, `deleted_at`, `referred_by_partner_id`, `referral_code_used`) VALUES
(1, 7, 'PT Karya Developer indonesia', 'Perusahaan yang bergerak dibidang industri teknogi informasi', 'https://karyadeveloperindonesia.com', '/docs/companies/1/logo_1768832403.png', 'Teknologi Informasi', '1-10', 'Sidaorjo, Jawa Timur', '+62881036480285', 'info@karyadeveloperindonesia.com', 'Perumahan Griya Bhayangkara blok i5/07 Desa Masangan Kulon, Kecamatan Sukodono, Kabupaten Sidoarjo, Kode Pos 61258', 'Sidoarjo', 'Jawa Timur', '61258', '2025', 8, 'verified', '/docs/companies/1/ktp_1768832313.jpg', '/docs/companies/1/akta_1768832354.pdf', '/docs/companies/1/npwp_1768832371.jpg', '/docs/companies/1/nib_1768832386.pdf', '2026-01-20 04:42:11', NULL, NULL, '2026-01-19 13:58:37', '2026-01-20 04:42:11', NULL, NULL, NULL),
(2, 2, 'PT TechCorp Indonesia', 'Perusahaan teknologi terkemuka yang menyediakan solusi IT', 'https://techcorp.id', NULL, 'Teknologi Informasi', '11-50', 'Jakarta Selatan', '021-123-4567', 'hr@techcorp.id', 'Jl. Sudirman No. 123, Jakarta Selatan', 'Jakarta', 'Jakarta', '12190', '2020', 25, 'verified', NULL, NULL, NULL, NULL, '2026-02-02 13:16:03', NULL, NULL, '2026-01-25 03:15:30', '2026-02-02 13:16:03', NULL, NULL, NULL),
(3, 4, 'CV Baru Startup', 'Startup muda yang mengembangkan aplikasi mobile', 'https://baristartup.com', NULL, 'Teknologi Informasi', '1-10', 'Bandung', '0274-555-6789', 'testcompany@test.com', 'Jl. Gatot Subroto No. 45, Bandung', 'Bandung', 'Jawa Barat', '40271', '2024', 5, 'rejected', NULL, NULL, NULL, NULL, NULL, NULL, 'Dokumen tidak lengkap', '2026-01-22 07:30:00', '2026-01-23 02:00:00', NULL, NULL, NULL),
(4, 5, 'PT Manufacturing Plus', 'Perusahaan manufaktur dengan standar internasional', 'https://manfacturingplus.co.id', NULL, 'Manufaktur', '51-200', 'Surabaya', '031-777-8888', 'company2@test.com', 'Jl. Ahmad Yani No. 888, Surabaya', 'Surabaya', 'Jawa Timur', '60188', '2018', 120, 'suspended', NULL, NULL, NULL, NULL, NULL, NULL, 'Melanggar kebijakan platform', '2026-01-15 01:00:00', '2026-01-24 08:45:00', NULL, NULL, NULL),
(5, 6, 'PT Konsultan HR Jaya', 'Perusahaan konsultasi sumber daya manusia', 'https://hrjaya.com', NULL, 'Konsultasi', '11-50', 'Medan', '061-444-5555', 'company3@test.com', 'Jl. Diponegoro No. 222, Medan', 'Medan', 'Sumatera Utara', '20212', '2019', 35, 'rejected', NULL, NULL, NULL, NULL, NULL, NULL, NULL, '2026-01-26 09:20:00', '2026-02-02 13:18:24', NULL, NULL, NULL),
(7, 26, 'PT Nusa Persada', 'PT Nusa Persada adalah perusahaan yang bergerak di bidang penyediaan peralatan medis yang berdiri sejak tahun 2012 dan berlokasi di Surabaya, tepatnya di kawasan Ngagel. Sejak awal berdirinya, PT Nusa Persada berkomitmen untuk mendukung peningkatan kualitas layanan kesehatan di Indonesia melalui penyediaan produk-produk medis yang berkualitas, aman, dan terpercaya.\n\nKami menyediakan berbagai jenis peralatan dan perlengkapan medis untuk rumah sakit, klinik, puskesmas, laboratorium, serta fasilitas layanan kesehatan lainnya. Seluruh produk yang kami distribusikan telah melalui proses seleksi dan memenuhi standar mutu yang berlaku, sehingga dapat menunjang operasional tenaga medis secara optimal.\n\nDengan pengalaman lebih dari satu dekade, PT Nusa Persada terus mengembangkan jaringan distribusi serta meningkatkan kualitas layanan, baik dari segi kecepatan pengiriman, ketepatan produk, maupun pelayanan purna jual. Didukung oleh tim yang profesional dan berpengalaman, kami senantiasa berupaya menjadi mitra terpercaya dalam memenuhi kebutuhan peralatan medis di wilayah Surabaya dan seluruh Indonesia.\n\nPT Nusa Persada berkomitmen untuk tumbuh bersama pelanggan dengan mengedepankan integritas, profesionalisme, dan kepuasan pelanggan sebagai prioritas utama.', 'https://nusapersada.co.id', '/docs/companies/7/logo_1770811199.png', 'Kesehatan', '51-200', 'Ngagel, Surabaya', '+62881036480285', 'info@nusapersada.com', 'Jl. Kertajaya No.164, Kertajaya, Kec. Gubeng, Surabaya, Jawa Timur 60282', 'Surabaya', 'Jawa Timur', '60282', '2012', 430, 'verified', '/docs/companies/7/ktp_1770811244.jpg', '/docs/companies/7/akta_1770811337.png', '/docs/companies/7/npwp_1770811477.pdf', '/docs/companies/7/nib_1770811465.pdf', '2026-02-11 12:05:28', NULL, NULL, '2026-02-09 13:08:22', '2026-02-11 12:05:28', NULL, 3, 'SAPUC4EB');

-- --------------------------------------------------------

--
-- Table structure for table `company_quotas`
--

CREATE TABLE `company_quotas` (
  `id` bigint(20) UNSIGNED NOT NULL,
  `company_id` bigint(20) UNSIGNED NOT NULL,
  `free_quota_used` int(11) NOT NULL DEFAULT 0,
  `paid_quota` int(11) NOT NULL DEFAULT 0,
  `created_at` timestamp NOT NULL DEFAULT current_timestamp(),
  `updated_at` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp()
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

--
-- Dumping data for table `company_quotas`
--

INSERT INTO `company_quotas` (`id`, `company_id`, `free_quota_used`, `paid_quota`, `created_at`, `updated_at`) VALUES
(1, 4, 0, 0, '2026-01-18 07:15:18', '2026-01-18 07:15:18'),
(3, 1, 10, 31, '2026-01-20 07:33:46', '2026-02-02 04:09:31'),
(4, 2, 0, 0, '2026-02-02 04:06:57', '2026-02-02 04:06:57'),
(5, 5, 0, 0, '2026-02-02 04:06:57', '2026-02-02 04:06:57'),
(6, 3, 0, 3, '2026-02-02 04:06:57', '2026-02-02 04:07:05'),
(7, 7, 0, 0, '2026-02-09 13:08:28', '2026-02-09 13:08:28'),
(8, 26, 1, 0, '2026-02-11 14:54:30', '2026-02-11 14:54:30');

-- --------------------------------------------------------

--
-- Table structure for table `conversations`
--

CREATE TABLE `conversations` (
  `id` bigint(20) UNSIGNED NOT NULL,
  `company_id` bigint(20) UNSIGNED NOT NULL,
  `title` varchar(255) NOT NULL,
  `subject` text NOT NULL,
  `category` enum('complaint','helpdesk','general','urgent') DEFAULT 'general',
  `status` enum('open','in_progress','resolved','closed') DEFAULT 'open',
  `created_at` timestamp NOT NULL DEFAULT current_timestamp(),
  `updated_at` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  `closed_at` timestamp NULL DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

--
-- Dumping data for table `conversations`
--

INSERT INTO `conversations` (`id`, `company_id`, `title`, `subject`, `category`, `status`, `created_at`, `updated_at`, `closed_at`) VALUES
(1, 7, 'Help Desk: Tidak bisa posting lowongan baru', 'Tidak bisa posting lowongan baru', 'helpdesk', 'closed', '2026-01-21 17:42:17', '2026-01-22 03:15:05', '2026-01-22 03:15:05'),
(2, 7, 'Help Desk: Tidak bisa posting lowongan baru', 'Tidak bisa posting lowongan baru', 'helpdesk', 'resolved', '2026-01-21 17:42:27', '2026-01-21 17:49:26', NULL),
(3, 7, 'Pertanyaan: Test dari curl', 'Test dari curl', 'general', 'resolved', '2026-01-22 01:38:43', '2026-02-02 02:12:48', '2026-01-22 03:03:32'),
(4, 7, 'Help Desk: Test E2E Chat', 'Test E2E Chat', 'helpdesk', 'closed', '2026-01-22 01:39:48', '2026-01-22 03:14:59', '2026-01-22 03:14:59'),
(5, 7, 'Help Desk: Pesan Website', 'Pesan Website', 'helpdesk', 'closed', '2026-02-02 13:20:47', '2026-02-02 14:02:44', '2026-02-02 14:02:44'),
(6, 7, 'Komplain: Komplain', 'Komplain', 'complaint', 'closed', '2026-02-02 14:14:16', '2026-02-02 14:51:00', '2026-02-02 14:51:00');

-- --------------------------------------------------------

--
-- Table structure for table `cvs`
--

CREATE TABLE `cvs` (
  `id` bigint(20) UNSIGNED NOT NULL,
  `user_id` bigint(20) UNSIGNED NOT NULL,
  `personal_info` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL CHECK (json_valid(`personal_info`)),
  `education` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL DEFAULT '[]' CHECK (json_valid(`education`)),
  `experience` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL DEFAULT '[]' CHECK (json_valid(`experience`)),
  `skills` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL DEFAULT '[]' CHECK (json_valid(`skills`)),
  `certifications` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL DEFAULT '[]' CHECK (json_valid(`certifications`)),
  `languages` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL DEFAULT '[]' CHECK (json_valid(`languages`)),
  `projects` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL DEFAULT '[]' CHECK (json_valid(`projects`)),
  `last_updated_at` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  `completeness_score` int(10) UNSIGNED NOT NULL DEFAULT 0,
  `created_at` timestamp NOT NULL DEFAULT current_timestamp(),
  `updated_at` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp()
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

--
-- Dumping data for table `cvs`
--

INSERT INTO `cvs` (`id`, `user_id`, `personal_info`, `education`, `experience`, `skills`, `certifications`, `languages`, `projects`, `last_updated_at`, `completeness_score`, `created_at`, `updated_at`) VALUES
(2, 3, '{\"full_name\":\"Budi Santoso\",\"email\":\"budi@gmail.com\",\"phone\":\"+6281234567890\"}', 'null', '[{\"company\":\"PT Software House\",\"position\":\"Backend Developer\",\"start_date\":\"2019-08-01\",\"is_current\":true,\"description\":\"Developing REST APIs\"}]', '[{\"name\":\"Go\",\"level\":\"advanced\"}]', 'null', 'null', 'null', '2026-01-17 03:12:22', 45, '2026-01-17 02:51:16', '2026-01-17 03:12:22'),
(3, 20, '{\"full_name\":\"Jastiska Dwi Wanda Sari\",\"email\":\"jastiska14@gmail.com\",\"phone\":\"08893011438\",\"address\":\"Dusun Sumber Pandan, Desa Bulusari, Kecamatan Gempol, Kabupaten Pasuruan, Provinsi Jawa Timur\",\"summary\":\"Saya memiliki pengalaman dibidang Sales Marketing Jasa pada salah satu perusahaan asuransi terkenal di Indonesia dan memiliki pengalaman kurang lebih 5 tahun.\"}', '[{\"institution\":\"Universitas Islam Malang\",\"degree\":\"S1\",\"field_of_study\":\"PGSD\",\"start_date\":\"2025-01-01\",\"end_date\":\"2026-12-31\"}]', '[{\"company\":\"PT Bank Nasional Indonesia Life\",\"position\":\"Product Marketing\",\"start_date\":\"2025-11-13\",\"is_current\":true,\"description\":\"Melakukan pemasaran produk kepada calon nasabah atau nasabah BNI\"}]', '[{\"name\":\"SQL\",\"level\":\"intermediate\"},{\"name\":\"Node.js\",\"level\":\"intermediate\"},{\"name\":\"React\",\"level\":\"intermediate\"},{\"name\":\"Python\",\"level\":\"intermediate\"},{\"name\":\"Leadership\",\"level\":\"intermediate\"},{\"name\":\"Teamwork\",\"level\":\"intermediate\"},{\"name\":\"Problem Solving\",\"level\":\"intermediate\"}]', '[]', '[]', 'null', '2026-02-11 14:56:25', 80, '2026-01-23 03:55:55', '2026-02-11 14:56:25'),
(4, 21, '{\"full_name\":\"Saputra Budianto\",\"email\":\"craftgirlsssshopping@gmail.com\",\"phone\":\"0881036480285\",\"address\":\"Perumahan Griya Bhayangkara blok i5/07 Desa Masangan Kulon, Kecamatan Sukodono, Kabupaten Sidoarjo, Kode Pos 61258\",\"summary\":\"Saya adalah seorang profesional yang berdedikasi di bidang Mobile Apps Developer dengan pengalaman selama 3 tahun dalam mengembangkan cross platform mobile apps. Saya merupakan lulusan dari Universitas Muhammdiyah Sidoarjo, di mana saya mengasah kemampuan analitis dan teknis yang menjadi fondasi karier saya saat ini.\\n\\nSepanjang perjalanan profesional saya, saya telah berhasil mengembangkan dan mendistribusikan beberapa aplikasi ke marketplace dan web publik. Saya dikenal sebagai pribadi yang adaptif, komunikatif, dan memiliki orientasi kuat terhadap hasil. Saya selalu antusias untuk mempelajari teknologi baru dan berkontribusi dalam proyek-proyek inovatif yang memberikan dampak positif bagi organisasi.\",\"linkedin\":\"https://www.linkedin.com/in/saputra-budianto23/\"}', '[]', '[{\"company\":\"PT Adi Karya Media\",\"position\":\"Mobile Apps Developer\",\"start_date\":\"2023-03-23\",\"end_date\":\"2023-09-30\",\"is_current\":false,\"description\":\"Membuat aplikasi berbasis moble app menggunakan framework flutter dan integrasi API menggunakan Supabase dan Firebase\"},{\"company\":\"PT. All Media Indo\",\"position\":\"Mobile Apps Developer\",\"start_date\":\"2023-11-23\",\"is_current\":true,\"description\":\"Membuat aplikasi berbasis moble app menggunakan framework flutter dan integrasi API menggunakan Supabase dan Firebase, serta membuat API\"}]', '[{\"name\":\"JavaScript\",\"level\":\"intermediate\"},{\"name\":\"React\",\"level\":\"intermediate\"},{\"name\":\"SQL\",\"level\":\"intermediate\"},{\"name\":\"Marketing\",\"level\":\"intermediate\"}]', '[{\"name\":\"AWS Certificate\",\"issuer\":\"Amazon Web Server\",\"issue_date\":\"2025-01-01\",\"credential_id\":\"XWXIS23AO\"}]', '[]', 'null', '2026-01-26 02:55:46', 70, '2026-01-23 03:57:04', '2026-01-26 02:55:46');

-- --------------------------------------------------------

--
-- Table structure for table `cv_snapshots`
--

CREATE TABLE `cv_snapshots` (
  `id` bigint(20) UNSIGNED NOT NULL,
  `cv_id` bigint(20) UNSIGNED NOT NULL,
  `user_id` bigint(20) UNSIGNED NOT NULL,
  `personal_info` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL CHECK (json_valid(`personal_info`)),
  `education` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL CHECK (json_valid(`education`)),
  `experience` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL CHECK (json_valid(`experience`)),
  `skills` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL CHECK (json_valid(`skills`)),
  `certifications` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL CHECK (json_valid(`certifications`)),
  `languages` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL DEFAULT '[]' CHECK (json_valid(`languages`)),
  `projects` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL DEFAULT '[]' CHECK (json_valid(`projects`)),
  `snapshot_hash` varchar(64) NOT NULL,
  `completeness_score` int(10) UNSIGNED NOT NULL DEFAULT 0,
  `created_at` timestamp NOT NULL DEFAULT current_timestamp()
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

--
-- Dumping data for table `cv_snapshots`
--

INSERT INTO `cv_snapshots` (`id`, `cv_id`, `user_id`, `personal_info`, `education`, `experience`, `skills`, `certifications`, `languages`, `projects`, `snapshot_hash`, `completeness_score`, `created_at`) VALUES
(1, 2, 3, '{\"full_name\":\"Budi Santoso\",\"email\":\"budi.kandidat@gmail.com\",\"phone\":\"+6281234567890\"}', '[{\"institution\":\"UI\",\"degree\":\"S1\",\"field_of_study\":\"Informatika\",\"start_date\":\"2015-08-01\",\"end_date\":\"2019-07-01\"}]', 'null', '[{\"name\":\"Go\",\"level\":\"advanced\"}]', 'null', 'null', 'null', 'ff53fa4b736225a2f2cecb569e9455e38fb7a4de93409f39dcf6b96c7d489c74', 50, '2026-01-17 02:51:28'),
(3, 3, 20, '{\"full_name\":\"Jastiska Dwi Wanda Sari\",\"email\":\"jastiska14@gmail.com\",\"phone\":\"08893011438\",\"address\":\"Dusun Sumber Pandan, Desa Bulusari, Kecamatan Gempol, Kabupaten Pasuruan, Provinsi Jawa Timur\",\"summary\":\"Saya memiliki pengalaman dibidang Sales Marketing Jasa pada salah satu perusahaan asuransi terkenal di Indonesia dan memiliki pengalaman kurang lebih 5 tahun.\"}', '[{\"institution\":\"Universitas Islam Malang\",\"degree\":\"S1\",\"field_of_study\":\"PGSD\",\"start_date\":\"2025-01-01\",\"end_date\":\"2026-12-31\"}]', '[{\"company\":\"PT Bank Nasional Indonesia Life\",\"position\":\"Product Marketing\",\"start_date\":\"2025-11-13\",\"is_current\":true,\"description\":\"Melakukan pemasaran produk kepada calon nasabah atau nasabah BNI\"}]', '[{\"name\":\"SQL\",\"level\":\"intermediate\"},{\"name\":\"Node.js\",\"level\":\"intermediate\"},{\"name\":\"React\",\"level\":\"intermediate\"},{\"name\":\"Python\",\"level\":\"intermediate\"},{\"name\":\"Leadership\",\"level\":\"intermediate\"},{\"name\":\"Teamwork\",\"level\":\"intermediate\"},{\"name\":\"Problem Solving\",\"level\":\"intermediate\"}]', '[]', '[]', 'null', 'a370fce1c6f49dc80fe4c9510a0d964e00719940f7f64e81a4db01d60a1b52ce', 80, '2026-01-23 13:48:42'),
(4, 4, 21, '{\"full_name\":\"Saputra Budianto\",\"email\":\"craftgirlsssshopping@gmail.com\",\"phone\":\"0881036480285\",\"address\":\"Perumahan Griya Bhayangkara blok i5/07 Desa Masangan Kulon, Kecamatan Sukodono, Kabupaten Sidoarjo, Kode Pos 61258\",\"summary\":\"Saya adalah seorang profesional yang berdedikasi di bidang Mobile Apps Developer dengan pengalaman selama 3 tahun dalam mengembangkan cross platform mobile apps. Saya merupakan lulusan dari Universitas Muhammdiyah Sidoarjo, di mana saya mengasah kemampuan analitis dan teknis yang menjadi fondasi karier saya saat ini.\\n\\nSepanjang perjalanan profesional saya, saya telah berhasil mengembangkan dan mendistribusikan beberapa aplikasi ke marketplace dan web publik. Saya dikenal sebagai pribadi yang adaptif, komunikatif, dan memiliki orientasi kuat terhadap hasil. Saya selalu antusias untuk mempelajari teknologi baru dan berkontribusi dalam proyek-proyek inovatif yang memberikan dampak positif bagi organisasi.\",\"linkedin\":\"https://www.linkedin.com/in/saputra-budianto23/\"}', '[]', '[{\"company\":\"PT Adi Karya Media\",\"position\":\"Mobile Apps Developer\",\"start_date\":\"2023-03-23\",\"end_date\":\"2023-09-30\",\"is_current\":false,\"description\":\"Membuat aplikasi berbasis moble app menggunakan framework flutter dan integrasi API menggunakan Supabase dan Firebase\"},{\"company\":\"PT. All Media Indo\",\"position\":\"Mobile Apps Developer\",\"start_date\":\"2023-11-23\",\"is_current\":true,\"description\":\"Membuat aplikasi berbasis moble app menggunakan framework flutter dan integrasi API menggunakan Supabase dan Firebase, serta membuat API\"}]', '[{\"name\":\"JavaScript\",\"level\":\"intermediate\"},{\"name\":\"React\",\"level\":\"intermediate\"},{\"name\":\"SQL\",\"level\":\"intermediate\"},{\"name\":\"Marketing\",\"level\":\"intermediate\"}]', '[{\"name\":\"AWS Certificate\",\"issuer\":\"Amazon Web Server\",\"issue_date\":\"2025-01-01\",\"credential_id\":\"XWXIS23AO\"}]', '[]', 'null', 'b137c6afadfa033a8916006e673e571d68bbb867f342cfd7631f12f2ab68d54e', 70, '2026-01-23 14:33:27'),
(5, 4, 21, '{\"full_name\":\"Saputra Budianto\",\"email\":\"craftgirlsssshopping@gmail.com\",\"phone\":\"0881036480285\",\"address\":\"Perumahan Griya Bhayangkara blok i5/07 Desa Masangan Kulon, Kecamatan Sukodono, Kabupaten Sidoarjo, Kode Pos 61258\",\"summary\":\"Saya adalah seorang profesional yang berdedikasi di bidang Mobile Apps Developer dengan pengalaman selama 3 tahun dalam mengembangkan cross platform mobile apps. Saya merupakan lulusan dari Universitas Muhammdiyah Sidoarjo, di mana saya mengasah kemampuan analitis dan teknis yang menjadi fondasi karier saya saat ini.\\n\\nSepanjang perjalanan profesional saya, saya telah berhasil mengembangkan dan mendistribusikan beberapa aplikasi ke marketplace dan web publik. Saya dikenal sebagai pribadi yang adaptif, komunikatif, dan memiliki orientasi kuat terhadap hasil. Saya selalu antusias untuk mempelajari teknologi baru dan berkontribusi dalam proyek-proyek inovatif yang memberikan dampak positif bagi organisasi.\",\"linkedin\":\"https://www.linkedin.com/in/saputra-budianto23/\"}', '[]', '[{\"company\":\"PT Adi Karya Media\",\"position\":\"Mobile Apps Developer\",\"start_date\":\"2023-03-23\",\"end_date\":\"2023-09-30\",\"is_current\":false,\"description\":\"Membuat aplikasi berbasis moble app menggunakan framework flutter dan integrasi API menggunakan Supabase dan Firebase\"},{\"company\":\"PT. All Media Indo\",\"position\":\"Mobile Apps Developer\",\"start_date\":\"2023-11-23\",\"is_current\":true,\"description\":\"Membuat aplikasi berbasis moble app menggunakan framework flutter dan integrasi API menggunakan Supabase dan Firebase, serta membuat API\"}]', '[{\"name\":\"JavaScript\",\"level\":\"intermediate\"},{\"name\":\"React\",\"level\":\"intermediate\"},{\"name\":\"SQL\",\"level\":\"intermediate\"},{\"name\":\"Marketing\",\"level\":\"intermediate\"}]', '[{\"name\":\"AWS Certificate\",\"issuer\":\"Amazon Web Server\",\"issue_date\":\"2025-01-01\",\"credential_id\":\"XWXIS23AO\"}]', '[]', 'null', 'b137c6afadfa033a8916006e673e571d68bbb867f342cfd7631f12f2ab68d54e', 70, '2026-01-23 14:36:59'),
(6, 4, 21, '{\"full_name\":\"Saputra Budianto\",\"email\":\"craftgirlsssshopping@gmail.com\",\"phone\":\"0881036480285\",\"address\":\"Perumahan Griya Bhayangkara blok i5/07 Desa Masangan Kulon, Kecamatan Sukodono, Kabupaten Sidoarjo, Kode Pos 61258\",\"summary\":\"Saya adalah seorang profesional yang berdedikasi di bidang Mobile Apps Developer dengan pengalaman selama 3 tahun dalam mengembangkan cross platform mobile apps. Saya merupakan lulusan dari Universitas Muhammdiyah Sidoarjo, di mana saya mengasah kemampuan analitis dan teknis yang menjadi fondasi karier saya saat ini.\\n\\nSepanjang perjalanan profesional saya, saya telah berhasil mengembangkan dan mendistribusikan beberapa aplikasi ke marketplace dan web publik. Saya dikenal sebagai pribadi yang adaptif, komunikatif, dan memiliki orientasi kuat terhadap hasil. Saya selalu antusias untuk mempelajari teknologi baru dan berkontribusi dalam proyek-proyek inovatif yang memberikan dampak positif bagi organisasi.\",\"linkedin\":\"https://www.linkedin.com/in/saputra-budianto23/\"}', '[]', '[{\"company\":\"PT Adi Karya Media\",\"position\":\"Mobile Apps Developer\",\"start_date\":\"2023-03-23\",\"end_date\":\"2023-09-30\",\"is_current\":false,\"description\":\"Membuat aplikasi berbasis moble app menggunakan framework flutter dan integrasi API menggunakan Supabase dan Firebase\"},{\"company\":\"PT. All Media Indo\",\"position\":\"Mobile Apps Developer\",\"start_date\":\"2023-11-23\",\"is_current\":true,\"description\":\"Membuat aplikasi berbasis moble app menggunakan framework flutter dan integrasi API menggunakan Supabase dan Firebase, serta membuat API\"}]', '[{\"name\":\"JavaScript\",\"level\":\"intermediate\"},{\"name\":\"React\",\"level\":\"intermediate\"},{\"name\":\"SQL\",\"level\":\"intermediate\"},{\"name\":\"Marketing\",\"level\":\"intermediate\"}]', '[{\"name\":\"AWS Certificate\",\"issuer\":\"Amazon Web Server\",\"issue_date\":\"2025-01-01\",\"credential_id\":\"XWXIS23AO\"}]', '[]', 'null', 'b137c6afadfa033a8916006e673e571d68bbb867f342cfd7631f12f2ab68d54e', 70, '2026-01-23 14:45:37'),
(7, 4, 21, '{\"full_name\":\"Saputra Budianto\",\"email\":\"craftgirlsssshopping@gmail.com\",\"phone\":\"0881036480285\",\"address\":\"Perumahan Griya Bhayangkara blok i5/07 Desa Masangan Kulon, Kecamatan Sukodono, Kabupaten Sidoarjo, Kode Pos 61258\",\"summary\":\"Saya adalah seorang profesional yang berdedikasi di bidang Mobile Apps Developer dengan pengalaman selama 3 tahun dalam mengembangkan cross platform mobile apps. Saya merupakan lulusan dari Universitas Muhammdiyah Sidoarjo, di mana saya mengasah kemampuan analitis dan teknis yang menjadi fondasi karier saya saat ini.\\n\\nSepanjang perjalanan profesional saya, saya telah berhasil mengembangkan dan mendistribusikan beberapa aplikasi ke marketplace dan web publik. Saya dikenal sebagai pribadi yang adaptif, komunikatif, dan memiliki orientasi kuat terhadap hasil. Saya selalu antusias untuk mempelajari teknologi baru dan berkontribusi dalam proyek-proyek inovatif yang memberikan dampak positif bagi organisasi.\",\"linkedin\":\"https://www.linkedin.com/in/saputra-budianto23/\"}', '[]', '[{\"company\":\"PT Adi Karya Media\",\"position\":\"Mobile Apps Developer\",\"start_date\":\"2023-03-23\",\"end_date\":\"2023-09-30\",\"is_current\":false,\"description\":\"Membuat aplikasi berbasis moble app menggunakan framework flutter dan integrasi API menggunakan Supabase dan Firebase\"},{\"company\":\"PT. All Media Indo\",\"position\":\"Mobile Apps Developer\",\"start_date\":\"2023-11-23\",\"is_current\":true,\"description\":\"Membuat aplikasi berbasis moble app menggunakan framework flutter dan integrasi API menggunakan Supabase dan Firebase, serta membuat API\"}]', '[{\"name\":\"JavaScript\",\"level\":\"intermediate\"},{\"name\":\"React\",\"level\":\"intermediate\"},{\"name\":\"SQL\",\"level\":\"intermediate\"},{\"name\":\"Marketing\",\"level\":\"intermediate\"}]', '[{\"name\":\"AWS Certificate\",\"issuer\":\"Amazon Web Server\",\"issue_date\":\"2025-01-01\",\"credential_id\":\"XWXIS23AO\"}]', '[]', 'null', 'b137c6afadfa033a8916006e673e571d68bbb867f342cfd7631f12f2ab68d54e', 70, '2026-01-23 14:50:24'),
(8, 4, 21, '{\"full_name\":\"Saputra Budianto\",\"email\":\"craftgirlsssshopping@gmail.com\",\"phone\":\"0881036480285\",\"address\":\"Perumahan Griya Bhayangkara blok i5/07 Desa Masangan Kulon, Kecamatan Sukodono, Kabupaten Sidoarjo, Kode Pos 61258\",\"summary\":\"Saya adalah seorang profesional yang berdedikasi di bidang Mobile Apps Developer dengan pengalaman selama 3 tahun dalam mengembangkan cross platform mobile apps. Saya merupakan lulusan dari Universitas Muhammdiyah Sidoarjo, di mana saya mengasah kemampuan analitis dan teknis yang menjadi fondasi karier saya saat ini.\\n\\nSepanjang perjalanan profesional saya, saya telah berhasil mengembangkan dan mendistribusikan beberapa aplikasi ke marketplace dan web publik. Saya dikenal sebagai pribadi yang adaptif, komunikatif, dan memiliki orientasi kuat terhadap hasil. Saya selalu antusias untuk mempelajari teknologi baru dan berkontribusi dalam proyek-proyek inovatif yang memberikan dampak positif bagi organisasi.\",\"linkedin\":\"https://www.linkedin.com/in/saputra-budianto23/\"}', '[]', '[{\"company\":\"PT Adi Karya Media\",\"position\":\"Mobile Apps Developer\",\"start_date\":\"2023-03-23\",\"end_date\":\"2023-09-30\",\"is_current\":false,\"description\":\"Membuat aplikasi berbasis moble app menggunakan framework flutter dan integrasi API menggunakan Supabase dan Firebase\"},{\"company\":\"PT. All Media Indo\",\"position\":\"Mobile Apps Developer\",\"start_date\":\"2023-11-23\",\"is_current\":true,\"description\":\"Membuat aplikasi berbasis moble app menggunakan framework flutter dan integrasi API menggunakan Supabase dan Firebase, serta membuat API\"}]', '[{\"name\":\"JavaScript\",\"level\":\"intermediate\"},{\"name\":\"React\",\"level\":\"intermediate\"},{\"name\":\"SQL\",\"level\":\"intermediate\"},{\"name\":\"Marketing\",\"level\":\"intermediate\"}]', '[{\"name\":\"AWS Certificate\",\"issuer\":\"Amazon Web Server\",\"issue_date\":\"2025-01-01\",\"credential_id\":\"XWXIS23AO\"}]', '[]', 'null', 'b137c6afadfa033a8916006e673e571d68bbb867f342cfd7631f12f2ab68d54e', 70, '2026-01-23 14:55:13'),
(9, 4, 21, '{\"full_name\":\"Saputra Budianto\",\"email\":\"craftgirlsssshopping@gmail.com\",\"phone\":\"0881036480285\",\"address\":\"Perumahan Griya Bhayangkara blok i5/07 Desa Masangan Kulon, Kecamatan Sukodono, Kabupaten Sidoarjo, Kode Pos 61258\",\"summary\":\"Saya adalah seorang profesional yang berdedikasi di bidang Mobile Apps Developer dengan pengalaman selama 3 tahun dalam mengembangkan cross platform mobile apps. Saya merupakan lulusan dari Universitas Muhammdiyah Sidoarjo, di mana saya mengasah kemampuan analitis dan teknis yang menjadi fondasi karier saya saat ini.\\n\\nSepanjang perjalanan profesional saya, saya telah berhasil mengembangkan dan mendistribusikan beberapa aplikasi ke marketplace dan web publik. Saya dikenal sebagai pribadi yang adaptif, komunikatif, dan memiliki orientasi kuat terhadap hasil. Saya selalu antusias untuk mempelajari teknologi baru dan berkontribusi dalam proyek-proyek inovatif yang memberikan dampak positif bagi organisasi.\",\"linkedin\":\"https://www.linkedin.com/in/saputra-budianto23/\"}', '[]', '[{\"company\":\"PT Adi Karya Media\",\"position\":\"Mobile Apps Developer\",\"start_date\":\"2023-03-23\",\"end_date\":\"2023-09-30\",\"is_current\":false,\"description\":\"Membuat aplikasi berbasis moble app menggunakan framework flutter dan integrasi API menggunakan Supabase dan Firebase\"},{\"company\":\"PT. All Media Indo\",\"position\":\"Mobile Apps Developer\",\"start_date\":\"2023-11-23\",\"is_current\":true,\"description\":\"Membuat aplikasi berbasis moble app menggunakan framework flutter dan integrasi API menggunakan Supabase dan Firebase, serta membuat API\"}]', '[{\"name\":\"JavaScript\",\"level\":\"intermediate\"},{\"name\":\"React\",\"level\":\"intermediate\"},{\"name\":\"SQL\",\"level\":\"intermediate\"},{\"name\":\"Marketing\",\"level\":\"intermediate\"}]', '[{\"name\":\"AWS Certificate\",\"issuer\":\"Amazon Web Server\",\"issue_date\":\"2025-01-01\",\"credential_id\":\"XWXIS23AO\"}]', '[]', 'null', 'b137c6afadfa033a8916006e673e571d68bbb867f342cfd7631f12f2ab68d54e', 70, '2026-01-23 14:58:52'),
(10, 4, 21, '{\"full_name\":\"Saputra Budianto\",\"email\":\"craftgirlsssshopping@gmail.com\",\"phone\":\"0881036480285\",\"address\":\"Perumahan Griya Bhayangkara blok i5/07 Desa Masangan Kulon, Kecamatan Sukodono, Kabupaten Sidoarjo, Kode Pos 61258\",\"summary\":\"Saya adalah seorang profesional yang berdedikasi di bidang Mobile Apps Developer dengan pengalaman selama 3 tahun dalam mengembangkan cross platform mobile apps. Saya merupakan lulusan dari Universitas Muhammdiyah Sidoarjo, di mana saya mengasah kemampuan analitis dan teknis yang menjadi fondasi karier saya saat ini.\\n\\nSepanjang perjalanan profesional saya, saya telah berhasil mengembangkan dan mendistribusikan beberapa aplikasi ke marketplace dan web publik. Saya dikenal sebagai pribadi yang adaptif, komunikatif, dan memiliki orientasi kuat terhadap hasil. Saya selalu antusias untuk mempelajari teknologi baru dan berkontribusi dalam proyek-proyek inovatif yang memberikan dampak positif bagi organisasi.\",\"linkedin\":\"https://www.linkedin.com/in/saputra-budianto23/\"}', '[]', '[{\"company\":\"PT Adi Karya Media\",\"position\":\"Mobile Apps Developer\",\"start_date\":\"2023-03-23\",\"end_date\":\"2023-09-30\",\"is_current\":false,\"description\":\"Membuat aplikasi berbasis moble app menggunakan framework flutter dan integrasi API menggunakan Supabase dan Firebase\"},{\"company\":\"PT. All Media Indo\",\"position\":\"Mobile Apps Developer\",\"start_date\":\"2023-11-23\",\"is_current\":true,\"description\":\"Membuat aplikasi berbasis moble app menggunakan framework flutter dan integrasi API menggunakan Supabase dan Firebase, serta membuat API\"}]', '[{\"name\":\"JavaScript\",\"level\":\"intermediate\"},{\"name\":\"React\",\"level\":\"intermediate\"},{\"name\":\"SQL\",\"level\":\"intermediate\"},{\"name\":\"Marketing\",\"level\":\"intermediate\"}]', '[{\"name\":\"AWS Certificate\",\"issuer\":\"Amazon Web Server\",\"issue_date\":\"2025-01-01\",\"credential_id\":\"XWXIS23AO\"}]', '[]', 'null', 'b137c6afadfa033a8916006e673e571d68bbb867f342cfd7631f12f2ab68d54e', 70, '2026-01-23 15:04:57'),
(11, 3, 20, '{\"full_name\":\"Jastiska Dwi Wanda Sari\",\"email\":\"jastiska14@gmail.com\",\"phone\":\"08893011438\",\"address\":\"Dusun Sumber Pandan, Desa Bulusari, Kecamatan Gempol, Kabupaten Pasuruan, Provinsi Jawa Timur\",\"summary\":\"Saya memiliki pengalaman dibidang Sales Marketing Jasa pada salah satu perusahaan asuransi terkenal di Indonesia dan memiliki pengalaman kurang lebih 5 tahun.\"}', '[{\"institution\":\"Universitas Islam Malang\",\"degree\":\"S1\",\"field_of_study\":\"PGSD\",\"start_date\":\"2025-01-01\",\"end_date\":\"2026-12-31\"}]', '[{\"company\":\"PT Bank Nasional Indonesia Life\",\"position\":\"Product Marketing\",\"start_date\":\"2025-11-13\",\"is_current\":true,\"description\":\"Melakukan pemasaran produk kepada calon nasabah atau nasabah BNI\"}]', '[{\"name\":\"SQL\",\"level\":\"intermediate\"},{\"name\":\"Node.js\",\"level\":\"intermediate\"},{\"name\":\"React\",\"level\":\"intermediate\"},{\"name\":\"Python\",\"level\":\"intermediate\"},{\"name\":\"Leadership\",\"level\":\"intermediate\"},{\"name\":\"Teamwork\",\"level\":\"intermediate\"},{\"name\":\"Problem Solving\",\"level\":\"intermediate\"}]', '[]', '[]', 'null', 'a370fce1c6f49dc80fe4c9510a0d964e00719940f7f64e81a4db01d60a1b52ce', 80, '2026-01-25 09:07:38'),
(12, 3, 20, '{\"full_name\":\"Jastiska Dwi Wanda Sari\",\"email\":\"jastiska14@gmail.com\",\"phone\":\"08893011438\",\"address\":\"Dusun Sumber Pandan, Desa Bulusari, Kecamatan Gempol, Kabupaten Pasuruan, Provinsi Jawa Timur\",\"summary\":\"Saya memiliki pengalaman dibidang Sales Marketing Jasa pada salah satu perusahaan asuransi terkenal di Indonesia dan memiliki pengalaman kurang lebih 5 tahun.\"}', '[{\"institution\":\"Universitas Islam Malang\",\"degree\":\"S1\",\"field_of_study\":\"PGSD\",\"start_date\":\"2025-01-01\",\"end_date\":\"2026-12-31\"}]', '[{\"company\":\"PT Bank Nasional Indonesia Life\",\"position\":\"Product Marketing\",\"start_date\":\"2025-11-13\",\"is_current\":true,\"description\":\"Melakukan pemasaran produk kepada calon nasabah atau nasabah BNI\"}]', '[{\"name\":\"SQL\",\"level\":\"intermediate\"},{\"name\":\"Node.js\",\"level\":\"intermediate\"},{\"name\":\"React\",\"level\":\"intermediate\"},{\"name\":\"Python\",\"level\":\"intermediate\"},{\"name\":\"Leadership\",\"level\":\"intermediate\"},{\"name\":\"Teamwork\",\"level\":\"intermediate\"},{\"name\":\"Problem Solving\",\"level\":\"intermediate\"}]', '[]', '[]', 'null', 'a370fce1c6f49dc80fe4c9510a0d964e00719940f7f64e81a4db01d60a1b52ce', 80, '2026-02-11 14:56:30');

-- --------------------------------------------------------

--
-- Table structure for table `jobs`
--

CREATE TABLE `jobs` (
  `id` bigint(20) UNSIGNED NOT NULL,
  `company_id` bigint(20) UNSIGNED NOT NULL,
  `title` varchar(255) NOT NULL,
  `slug` varchar(300) NOT NULL,
  `description` text NOT NULL,
  `requirements` text DEFAULT NULL,
  `responsibilities` text DEFAULT NULL,
  `benefits` text DEFAULT NULL,
  `city` varchar(100) NOT NULL,
  `province` varchar(100) NOT NULL,
  `is_remote` tinyint(1) NOT NULL DEFAULT 0,
  `job_type` enum('full_time','part_time','contract','internship','freelance') NOT NULL DEFAULT 'full_time',
  `experience_level` enum('entry','junior','mid','senior','lead','executive') NOT NULL DEFAULT 'entry',
  `salary_min` bigint(20) UNSIGNED DEFAULT NULL,
  `salary_max` bigint(20) UNSIGNED DEFAULT NULL,
  `salary_currency` varchar(3) DEFAULT 'IDR',
  `is_salary_visible` tinyint(1) NOT NULL DEFAULT 1,
  `application_deadline` date DEFAULT NULL,
  `max_applications` int(10) UNSIGNED DEFAULT NULL,
  `status` enum('draft','active','paused','closed','filled') NOT NULL DEFAULT 'draft',
  `admin_status` enum('approved','rejected','flagged') DEFAULT NULL,
  `admin_note` text DEFAULT NULL,
  `flag_reason` text DEFAULT NULL,
  `views_count` int(10) UNSIGNED NOT NULL DEFAULT 0,
  `applications_count` int(10) UNSIGNED NOT NULL DEFAULT 0,
  `shares_count` int(10) UNSIGNED NOT NULL DEFAULT 0,
  `published_at` timestamp NULL DEFAULT NULL,
  `created_at` timestamp NOT NULL DEFAULT current_timestamp(),
  `updated_at` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  `deleted_at` timestamp NULL DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

--
-- Dumping data for table `jobs`
--

INSERT INTO `jobs` (`id`, `company_id`, `title`, `slug`, `description`, `requirements`, `responsibilities`, `benefits`, `city`, `province`, `is_remote`, `job_type`, `experience_level`, `salary_min`, `salary_max`, `salary_currency`, `is_salary_visible`, `application_deadline`, `max_applications`, `status`, `admin_status`, `admin_note`, `flag_reason`, `views_count`, `applications_count`, `shares_count`, `published_at`, `created_at`, `updated_at`, `deleted_at`) VALUES
(1, 2, 'Senior Software Engineer', 'senior-software-engineer', 'Kami mencari Senior Software Engineer untuk bergabung dengan tim development kami. Anda akan bekerja dengan teknologi terkini dan tim yang solid. Minimal pengalaman 3 tahun.', 'Minimal 3 tahun pengalaman. Menguasai Go/Python.', NULL, NULL, 'Jakarta Selatan', 'DKI Jakarta', 1, 'full_time', 'senior', 15000000, 25000000, 'IDR', 0, NULL, NULL, 'active', NULL, NULL, NULL, 0, 0, 0, '2026-01-17 02:49:00', '2026-01-17 02:46:11', '2026-01-17 02:49:00', NULL),
(2, 2, 'Test Status Job', 'test-status-job', 'Ini adalah deskripsi panjang untuk testing job status management endpoints yang baru ditambahkan.', NULL, NULL, NULL, 'Jakarta', 'DKI Jakarta', 0, 'full_time', 'mid', NULL, NULL, 'IDR', 0, NULL, NULL, 'closed', NULL, NULL, NULL, 0, 0, 0, '2026-01-17 03:13:52', '2026-01-17 03:13:25', '2026-01-17 03:14:15', NULL),
(3, 4, 'Senior Software Engineer', 'senior-software-engineer-1768720323', 'We are looking for an experienced software engineer', '5+ years experience with Go or Python', 'Build and maintain backend systems', 'Competitive salary, remote work', 'Jakarta Selatan', 'DKI Jakarta', 1, 'full_time', 'senior', 20000000, 35000000, 'IDR', 1, NULL, NULL, 'active', NULL, NULL, NULL, 0, 0, 0, '2026-01-18 07:12:10', '2026-01-18 07:12:03', '2026-01-18 07:12:19', NULL),
(4, 10, 'Senior Backend Engineer', 'senior-backend-engineer', 'Kami mencari Senior Backend Engineer yang berpengalaman dalam mengembangkan sistem scalable', '- Minimal 5 tahun pengalaman backend development\n- Mahir Go, Python, atau Java\n- Pengalaman dengan microservices\n- Pengalaman dengan database SQL dan NoSQL', '- Merancang dan mengimplementasi API\n- Melakukan code review\n- Mengoptimalkan performance sistem\n- Mentoring junior developers', '- Gaji kompetitif 15-25 juta/bulan\n- Asuransi kesehatan\n- Work from home flexibility\n- Training budget', 'Jakarta Selatan', 'DKI Jakarta', 1, 'full_time', 'senior', 15000000, 25000000, 'IDR', 1, NULL, NULL, 'active', NULL, NULL, NULL, 0, 0, 0, '2026-01-19 07:48:04', '2026-01-19 07:48:03', '2026-01-19 07:48:04', NULL),
(5, 10, 'Full Stack Developer', 'full-stack-developer', 'Bergabunglah dengan tim kami sebagai Full Stack Developer', '- Minimal 3 tahun pengalaman\n- React atau Vue.js\n- Node.js atau Python', '- Develop frontend dan backend\n- Collaborate dengan tim design', '- Gaji 8-12 juta/bulan\n- Remote friendly', 'Jakarta Pusat', 'DKI Jakarta', 1, 'full_time', 'mid', 8000000, 12000000, 'IDR', 1, NULL, NULL, 'draft', NULL, NULL, NULL, 0, 0, 0, NULL, '2026-01-19 07:48:03', '2026-01-19 07:48:03', NULL),
(6, 10, 'UI/UX Designer', 'uiux-designer', 'Kami mencari UI/UX Designer untuk mengembangkan product kami', '- 2+ tahun pengalaman UI/UX\n- Figma atau Adobe XD', '- Design interface\n- User research', '- Gaji 6-10 juta/bulan', 'Jakarta Selatan', 'DKI Jakarta', 1, 'full_time', 'junior', 6000000, 10000000, 'IDR', 1, NULL, NULL, 'draft', NULL, NULL, NULL, 0, 0, 0, NULL, '2026-01-19 07:48:03', '2026-01-19 07:48:03', NULL),
(7, 1, 'Senior Backend Developer', 'senior-backend-developer', 'We are looking for an experienced backend developer...', 'Node.js, TypeScript, Docker', 'Design and implement APIs', NULL, 'Jakarta', 'DKI Jakarta', 1, 'full_time', 'senior', NULL, NULL, 'IDR', 0, NULL, NULL, 'active', NULL, NULL, NULL, 0, 0, 0, '2026-01-20 07:33:46', '2026-01-20 03:43:14', '2026-01-20 07:33:46', NULL),
(8, 1, 'Senior Frontend Developer', 'senior-frontend-developer', 'Kami mencari Senior Frontend Developer yang berpengalaman dalam React.js dan TypeScript untuk bergabung dengan tim kami. Posisi ini akan bertanggung jawab untuk mengembangkan fitur-fitur baru.', '- Minimal 3 tahun pengalaman dengan React.js\n- Menguasai TypeScript', NULL, NULL, 'Jakarta Selatan', 'DKI Jakarta', 1, 'full_time', 'senior', 15000000, 25000000, 'IDR', 1, NULL, NULL, 'active', NULL, NULL, NULL, 0, 0, 0, '2026-01-20 05:05:33', '2026-01-20 05:05:33', '2026-01-20 06:11:28', '2026-01-20 06:11:28'),
(9, 1, 'Senior Frontend Developer', 'senior-frontend-developer-1768885756', 'Kami mencari Senior Frontend Developer yang berpengalaman dalam React.js dan TypeScript untuk bergabung dengan tim kami. Posisi ini akan bertanggung jawab untuk mengembangkan fitur-fitur baru.', '- Minimal 3 tahun pengalaman dengan React.js\n- Menguasai TypeScript', NULL, NULL, 'Jakarta Selatan', 'DKI Jakarta', 1, 'full_time', 'senior', 15000000, 25000000, 'IDR', 1, NULL, NULL, 'active', NULL, NULL, NULL, 0, 0, 0, '2026-01-20 05:09:16', '2026-01-20 05:09:16', '2026-01-20 06:11:08', '2026-01-20 06:11:08'),
(10, 1, 'Senior Frontend Developer', 'senior-frontend-developer-1768886367', 'Kami mencari Senior Frontend Developer yang berpengalaman dalam React.js.', 'Minimal 3 tahun pengalaman', NULL, NULL, 'Jakarta Selatan', 'DKI Jakarta', 1, 'full_time', 'senior', 15000000, 25000000, 'IDR', 1, NULL, NULL, 'closed', NULL, NULL, NULL, 0, 0, 0, '2026-01-20 05:19:27', '2026-01-20 05:19:27', '2026-01-20 06:10:33', NULL),
(11, 1, 'Senior Frontend Developer', 'senior-frontend-developer-1768886376', 'Kami mencari Senior Frontend Developer yang berpengalaman dalam React.js.', 'Minimal 3 tahun pengalaman', NULL, NULL, 'Jakarta Selatan', 'DKI Jakarta', 1, 'full_time', 'senior', 15000000, 25000000, 'IDR', 1, NULL, NULL, 'active', NULL, NULL, NULL, 0, 0, 0, '2026-01-20 05:19:36', '2026-01-20 05:19:36', '2026-01-20 06:09:48', '2026-01-20 06:09:48'),
(12, 1, 'Mobile Apps Developer', 'mobile-apps-developer', 'Full job description\nAbout this Role\n\nWe are seeking a highly experienced Flutter Developer to join our mobile development team. In this role, you will take the lead in architecting and building advanced mobile applications, driving technical excellence, and contributing to high-impact projects. You will collaborate with top-tier cross-functional teams to deliver innovative, scalable, and high-performance solutions.\n\nJob Description\n\nAs a Senior Flutter Developer, you will be responsible for designing, developing, and maintaining sophisticated mobile applications in a fast-paced and innovation-driven environment.\n\nWhat It’s Like to Work Here as a Senior Flutter Mobile App DeveloperFull-cycle Technical Ownership\n\nLead the end-to-end development lifecycle of mobile applications using Flutter, from system design and architecture to deployment and maintenance.\nUI/UX Collaboration\n\nTransform complex UI/UX designs into intuitive and polished user experiences, ensuring pixel-perfect implementations and smooth interactions.\nAdvanced Integration Expertise\n\nArchitect and integrate mobile applications with robust backend systems, ensuring high performance, security, and real-time synchronization.\nCode Quality & Engineering Standards\n\nDrive and participate in in-depth code reviews, ensuring clean architecture, maintainability, and adherence to best practices across the team.\nPerformance & Scalability Optimization\n\nIdentify performance bottlenecks, analyze app behavior, and implement advanced optimization techniques across devices and platforms.\nInnovation & Technical Leadership\n\nContribute forward-thinking ideas, mentor junior developers, and help shape technical strategy and direction.\nCross-functional Technical Collaboration\n\nWork closely with product managers, backend developers, and UI/UX teams to deliver seamless, scalable solutions on time.\nContinuous Learning & Research\n\nStay ahead of the latest trends in Flutter, mobile technologies, and tools—recommending and driving adoption of relevant advancements.\nRobust Testing & Quality Assurance\n\nDevelop automated testing, debugging strategies, and quality assurance processes to ensure enterprise-level reliability and security.', 'Requirements\n\nBachelor\'s degree in Computer Science or related field (or equivalent experience).\n5+ years of professional mobile development experience, with 3+ years specifically in Flutter.\nA strong portfolio showcasing complex, high-quality Flutter applications.\nMastery of Flutter, Dart, state management (e.g., Bloc, Riverpod, Provider, GetX), and modular architecture patterns.\nProven experience integrating APIs, working with real-time data, and using third-party libraries efficiently.\nHands-on experience with backend technologies (Node.js, Django, Firebase) is a strong advantage.\nFamiliarity with CI/CD pipelines, automated testing, and modern mobile DevOps practices.\nExperience publishing and maintaining apps on the App Store and Google Play.', NULL, NULL, 'Sidaorjo', 'Jawa Timur', 0, 'full_time', 'junior', 5000000, 5500000, 'IDR', 1, '2026-02-20', NULL, 'closed', NULL, NULL, NULL, 0, 0, 0, '2026-01-20 05:19:42', '2026-01-20 05:19:42', '2026-01-25 03:46:15', NULL),
(13, 1, 'Test Job API Debug', 'test-job-api-debug', 'Testing job creation for debug purposes and verification of company_id field. This is a longer description to pass validation.', 'Testing requirements that need to be longer for validation.', NULL, NULL, 'Jakarta', 'DKI Jakarta', 0, 'full_time', 'senior', NULL, NULL, 'IDR', 0, NULL, NULL, 'active', NULL, NULL, NULL, 0, 0, 0, '2026-01-20 05:30:01', '2026-01-20 05:30:01', '2026-01-20 06:08:42', '2026-01-20 06:08:42'),
(14, 1, 'Mobile Apps Developer', 'mobile-apps-developer-1768890085', 'About this Role\n\nWe are seeking a highly experienced Flutter Developer to join our mobile development team. In this role, you will take the lead in architecting and building advanced mobile applications, driving technical excellence, and contributing to high-impact projects. You will collaborate with top-tier cross-functional teams to deliver innovative, scalable, and high-performance solutions.\n\nJob Description\n\nAs a Senior Flutter Developer, you will be responsible for designing, developing, and maintaining sophisticated mobile applications in a fast-paced and innovation-driven environment.\n\nWhat It’s Like to Work Here as a Senior Flutter Mobile App DeveloperFull-cycle Technical Ownership\n\nLead the end-to-end development lifecycle of mobile applications using Flutter, from system design and architecture to deployment and maintenance.\nUI/UX Collaboration\n\nTransform complex UI/UX designs into intuitive and polished user experiences, ensuring pixel-perfect implementations and smooth interactions.\nAdvanced Integration Expertise\n\nArchitect and integrate mobile applications with robust backend systems, ensuring high performance, security, and real-time synchronization.\nCode Quality & Engineering Standards\n\nDrive and participate in in-depth code reviews, ensuring clean architecture, maintainability, and adherence to best practices across the team.\nPerformance & Scalability Optimization\n\nIdentify performance bottlenecks, analyze app behavior, and implement advanced optimization techniques across devices and platforms.\nInnovation & Technical Leadership\n\nContribute forward-thinking ideas, mentor junior developers, and help shape technical strategy and direction.\nCross-functional Technical Collaboration\n\nWork closely with product managers, backend developers, and UI/UX teams to deliver seamless, scalable solutions on time.\nContinuous Learning & Research\n\nStay ahead of the latest trends in Flutter, mobile technologies, and tools—recommending and driving adoption of relevant advancements.\nRobust Testing & Quality Assurance\n\nDevelop automated testing, debugging strategies, and quality assurance processes to ensure enterprise-level reliability and security.', 'Requirements\n\nBachelor\'s degree in Computer Science or related field (or equivalent experience).\n5+ years of professional mobile development experience, with 3+ years specifically in Flutter.\nA strong portfolio showcasing complex, high-quality Flutter applications.\nMastery of Flutter, Dart, state management (e.g., Bloc, Riverpod, Provider, GetX), and modular architecture patterns.\nProven experience integrating APIs, working with real-time data, and using third-party libraries efficiently.\nHands-on experience with backend technologies (Node.js, Django, Firebase) is a strong advantage.\nFamiliarity with CI/CD pipelines, automated testing, and modern mobile DevOps practices.\nExperience publishing and maintaining apps on the App Store and Google Play.', NULL, NULL, 'Sidaorjo', 'Jawa Timur', 0, 'full_time', 'junior', 5000000, 5500000, 'IDR', 1, '2026-02-20', NULL, 'closed', NULL, NULL, NULL, 0, 0, 0, '2026-01-20 06:21:25', '2026-01-20 06:21:25', '2026-01-20 07:02:12', '2026-01-20 07:02:12'),
(15, 1, 'QA Engineer', 'qa-engineer', 'Looking for experienced QA Engineer to join our team. You will be responsible for testing our applications and ensuring quality standards are met. Must have experience with automated testing tools.', '3+ years experience in QA testing', NULL, NULL, 'Jakarta', 'DKI Jakarta', 0, 'full_time', 'mid', NULL, NULL, 'IDR', 0, NULL, NULL, 'active', NULL, NULL, NULL, 0, 0, 0, '2026-01-20 07:34:26', '2026-01-20 07:34:26', '2026-01-20 07:34:26', NULL),
(16, 1, 'DevOps Engineer', 'devops-engineer', 'Looking for experienced DevOps Engineer to manage our cloud infrastructure. You will be responsible for CI/CD pipelines and infrastructure automation.', '3+ years experience with AWS, Docker, Kubernetes', NULL, NULL, 'Surabaya', 'Jawa Timur', 1, 'full_time', 'mid', NULL, NULL, 'IDR', 0, NULL, NULL, 'active', NULL, NULL, NULL, 0, 0, 0, '2026-01-20 07:39:03', '2026-01-20 07:39:03', '2026-01-20 07:39:03', NULL),
(17, 1, 'Senior Developer - Job 1', 'senior-developer-job-1', 'Kami mencari senior developer yang berpengalaman untuk bergabung dengan tim kami.', 'Pengalaman minimal 5 tahun di bidang development, menguasai TypeScript, React, dan Go', 'Mengembangkan fitur baru, melakukan code review, dan mentoring junior developer', 'Asuransi kesehatan, bonus tahunan, work from home', 'Jakarta', 'DKI Jakarta', 1, 'full_time', 'senior', 15000000, 25000000, 'IDR', 1, '2026-02-20', NULL, 'draft', NULL, NULL, NULL, 0, 0, 0, NULL, '2026-01-20 09:36:05', '2026-01-20 09:36:05', NULL),
(18, 1, 'Senior Developer - Job 2', 'senior-developer-job-2', 'Kami mencari senior developer yang berpengalaman untuk bergabung dengan tim kami.', 'Pengalaman minimal 5 tahun di bidang development, menguasai TypeScript, React, dan Go', 'Mengembangkan fitur baru, melakukan code review, dan mentoring junior developer', 'Asuransi kesehatan, bonus tahunan, work from home', 'Jakarta', 'DKI Jakarta', 1, 'full_time', 'senior', 15000000, 25000000, 'IDR', 1, '2026-02-20', NULL, 'draft', NULL, NULL, NULL, 0, 0, 0, NULL, '2026-01-20 09:36:05', '2026-01-20 09:36:05', NULL),
(19, 1, 'Senior Developer - Job 3', 'senior-developer-job-3', 'Kami mencari senior developer yang berpengalaman untuk bergabung dengan tim kami.', 'Pengalaman minimal 5 tahun di bidang development, menguasai TypeScript, React, dan Go', 'Mengembangkan fitur baru, melakukan code review, dan mentoring junior developer', 'Asuransi kesehatan, bonus tahunan, work from home', 'Jakarta', 'DKI Jakarta', 1, 'full_time', 'senior', 15000000, 25000000, 'IDR', 1, '2026-02-20', NULL, 'draft', NULL, NULL, NULL, 0, 0, 0, NULL, '2026-01-20 09:36:06', '2026-01-20 09:36:06', NULL),
(20, 1, 'Senior Developer - Job 4', 'senior-developer-job-4', 'Kami mencari senior developer yang berpengalaman untuk bergabung dengan tim kami.', 'Pengalaman minimal 5 tahun di bidang development, menguasai TypeScript, React, dan Go', 'Mengembangkan fitur baru, melakukan code review, dan mentoring junior developer', 'Asuransi kesehatan, bonus tahunan, work from home', 'Jakarta', 'DKI Jakarta', 1, 'full_time', 'senior', 15000000, 25000000, 'IDR', 1, '2026-02-20', NULL, 'draft', NULL, NULL, NULL, 0, 0, 0, NULL, '2026-01-20 09:36:06', '2026-01-20 09:36:06', NULL),
(21, 1, 'Senior Developer - Job 5', 'senior-developer-job-5', 'Kami mencari senior developer yang berpengalaman untuk bergabung dengan tim kami.', 'Pengalaman minimal 5 tahun di bidang development, menguasai TypeScript, React, dan Go', 'Mengembangkan fitur baru, melakukan code review, dan mentoring junior developer', 'Asuransi kesehatan, bonus tahunan, work from home', 'Jakarta', 'DKI Jakarta', 1, 'full_time', 'senior', 15000000, 25000000, 'IDR', 1, '2026-02-20', NULL, 'draft', NULL, NULL, NULL, 0, 0, 0, NULL, '2026-01-20 09:36:06', '2026-01-20 09:36:06', NULL),
(22, 1, 'Senior Developer - Job 6', 'senior-developer-job-6', 'Kami mencari senior developer yang berpengalaman untuk bergabung dengan tim kami.', 'Pengalaman minimal 5 tahun di bidang development, menguasai TypeScript, React, dan Go', 'Mengembangkan fitur baru, melakukan code review, dan mentoring junior developer', 'Asuransi kesehatan, bonus tahunan, work from home', 'Jakarta', 'DKI Jakarta', 1, 'full_time', 'senior', 15000000, 25000000, 'IDR', 1, '2026-02-20', NULL, 'draft', NULL, NULL, NULL, 0, 0, 0, NULL, '2026-01-20 09:36:06', '2026-01-20 09:36:06', NULL),
(23, 1, 'Senior Developer - Job 7', 'senior-developer-job-7', 'Kami mencari senior developer yang berpengalaman untuk bergabung dengan tim kami.', 'Pengalaman minimal 5 tahun di bidang development, menguasai TypeScript, React, dan Go', 'Mengembangkan fitur baru, melakukan code review, dan mentoring junior developer', 'Asuransi kesehatan, bonus tahunan, work from home', 'Jakarta', 'DKI Jakarta', 1, 'full_time', 'senior', 15000000, 25000000, 'IDR', 1, '2026-02-20', NULL, 'draft', NULL, NULL, NULL, 0, 0, 0, NULL, '2026-01-20 09:36:06', '2026-01-20 09:36:06', NULL),
(24, 1, 'Senior Developer - Job 8', 'senior-developer-job-8', 'Kami mencari senior developer yang berpengalaman untuk bergabung dengan tim kami.', 'Pengalaman minimal 5 tahun di bidang development, menguasai TypeScript, React, dan Go', 'Mengembangkan fitur baru, melakukan code review, dan mentoring junior developer', 'Asuransi kesehatan, bonus tahunan, work from home', 'Jakarta', 'DKI Jakarta', 1, 'full_time', 'senior', 15000000, 25000000, 'IDR', 1, '2026-02-20', NULL, 'draft', NULL, NULL, NULL, 0, 0, 0, NULL, '2026-01-20 09:36:06', '2026-01-20 09:36:06', NULL),
(25, 1, 'Senior Developer - Job 9', 'senior-developer-job-9', 'Kami mencari senior developer yang berpengalaman untuk bergabung dengan tim kami.', 'Pengalaman minimal 5 tahun di bidang development, menguasai TypeScript, React, dan Go', 'Mengembangkan fitur baru, melakukan code review, dan mentoring junior developer', 'Asuransi kesehatan, bonus tahunan, work from home', 'Jakarta', 'DKI Jakarta', 1, 'full_time', 'senior', 15000000, 25000000, 'IDR', 1, '2026-02-20', NULL, 'draft', NULL, NULL, NULL, 0, 0, 0, NULL, '2026-01-20 09:36:06', '2026-01-20 09:36:06', NULL),
(26, 1, 'Senior Developer - Job 10', 'senior-developer-job-10', 'Kami mencari senior developer yang berpengalaman untuk bergabung dengan tim kami.', 'Pengalaman minimal 5 tahun di bidang development, menguasai TypeScript, React, dan Go', 'Mengembangkan fitur baru, melakukan code review, dan mentoring junior developer', 'Asuransi kesehatan, bonus tahunan, work from home', 'Jakarta', 'DKI Jakarta', 1, 'full_time', 'senior', 15000000, 25000000, 'IDR', 1, '2026-02-20', NULL, 'draft', NULL, NULL, NULL, 0, 0, 0, NULL, '2026-01-20 09:36:06', '2026-01-20 09:36:06', NULL),
(27, 1, 'Senior Developer - Job 11 (Should Fail)', 'senior-developer-job-11-should-fail', 'Ini adalah job ke-11 yang seharusnya gagal karena kuota gratis sudah habis', 'Pengalaman minimal 5 tahun', NULL, NULL, 'Jakarta', 'DKI Jakarta', 0, 'full_time', 'senior', NULL, NULL, 'IDR', 0, NULL, NULL, 'draft', NULL, NULL, NULL, 0, 0, 0, NULL, '2026-01-20 09:36:06', '2026-01-20 09:36:06', NULL),
(28, 1, 'Senior Dev Job 1', 'senior-dev-job-1', 'We are hiring a senior developer with 5+ years experience', '5+ years experience, TypeScript, React, Go', NULL, NULL, 'Jakarta', 'DKI Jakarta', 1, 'full_time', 'senior', NULL, NULL, 'IDR', 0, NULL, NULL, 'active', NULL, NULL, NULL, 0, 0, 0, '2026-01-20 09:37:13', '2026-01-20 09:37:13', '2026-01-20 09:37:13', NULL),
(29, 1, 'Senior Dev Job 2', 'senior-dev-job-2', 'We are hiring a senior developer with 5+ years experience', '5+ years experience, TypeScript, React, Go', NULL, NULL, 'Jakarta', 'DKI Jakarta', 1, 'full_time', 'senior', NULL, NULL, 'IDR', 0, NULL, NULL, 'active', NULL, NULL, NULL, 0, 0, 0, '2026-01-20 09:37:13', '2026-01-20 09:37:13', '2026-01-20 09:37:13', NULL),
(30, 1, 'Senior Dev Job 3', 'senior-dev-job-3', 'We are hiring a senior developer with 5+ years experience', '5+ years experience, TypeScript, React, Go', NULL, NULL, 'Jakarta', 'DKI Jakarta', 1, 'full_time', 'senior', NULL, NULL, 'IDR', 0, NULL, NULL, 'active', NULL, NULL, NULL, 0, 0, 0, '2026-01-20 09:37:13', '2026-01-20 09:37:13', '2026-01-20 09:37:13', NULL),
(31, 1, 'Senior Dev Job 4', 'senior-dev-job-4', 'We are hiring a senior developer with 5+ years experience', '5+ years experience, TypeScript, React, Go', NULL, NULL, 'Jakarta', 'DKI Jakarta', 1, 'full_time', 'senior', NULL, NULL, 'IDR', 0, NULL, NULL, 'active', NULL, NULL, NULL, 0, 0, 0, '2026-01-20 09:37:13', '2026-01-20 09:37:13', '2026-01-20 09:37:13', NULL),
(32, 1, 'Senior Dev Job 5', 'senior-dev-job-5', 'We are hiring a senior developer with 5+ years experience', '5+ years experience, TypeScript, React, Go', NULL, NULL, 'Jakarta', 'DKI Jakarta', 1, 'full_time', 'senior', NULL, NULL, 'IDR', 0, NULL, NULL, 'active', NULL, NULL, NULL, 0, 0, 0, '2026-01-20 09:37:13', '2026-01-20 09:37:13', '2026-01-20 09:37:13', NULL),
(33, 1, 'Senior Dev Job 6', 'senior-dev-job-6', 'We are hiring a senior developer with 5+ years experience', '5+ years experience, TypeScript, React, Go', NULL, NULL, 'Jakarta', 'DKI Jakarta', 1, 'full_time', 'senior', NULL, NULL, 'IDR', 0, NULL, NULL, 'active', NULL, NULL, NULL, 0, 0, 0, '2026-01-20 09:37:13', '2026-01-20 09:37:13', '2026-01-20 09:37:13', NULL),
(34, 1, 'Senior Dev Job 7', 'senior-dev-job-7', 'We are hiring a senior developer with 5+ years experience', '5+ years experience, TypeScript, React, Go', NULL, NULL, 'Jakarta', 'DKI Jakarta', 1, 'full_time', 'senior', NULL, NULL, 'IDR', 0, NULL, NULL, 'active', NULL, NULL, NULL, 0, 0, 0, '2026-01-20 09:37:13', '2026-01-20 09:37:13', '2026-01-20 09:37:13', NULL),
(35, 1, 'Senior Dev Job 8', 'senior-dev-job-8', 'We are hiring a senior developer with 5+ years experience', '5+ years experience, TypeScript, React, Go', NULL, NULL, 'Jakarta', 'DKI Jakarta', 1, 'full_time', 'senior', NULL, NULL, 'IDR', 0, NULL, NULL, 'active', NULL, NULL, NULL, 0, 0, 0, '2026-01-20 09:37:13', '2026-01-20 09:37:13', '2026-01-20 09:37:13', NULL),
(36, 1, 'Senior Dev Job 9', 'senior-dev-job-9', 'We are hiring a senior developer with 5+ years experience', '5+ years experience, TypeScript, React, Go', NULL, NULL, 'Jakarta', 'DKI Jakarta', 1, 'full_time', 'senior', NULL, NULL, 'IDR', 0, NULL, NULL, 'active', NULL, NULL, NULL, 0, 0, 0, '2026-01-20 09:37:13', '2026-01-20 09:37:13', '2026-01-20 09:37:13', NULL),
(37, 1, 'Senior Developer Position 1', 'senior-developer-position-1', 'We are looking for an experienced senior developer with at least 5 years of experience in modern web technologies including TypeScript, React, and backend frameworks.', 'Minimum 5 years of experience with TypeScript, React, Go, and PostgreSQL. Strong problem-solving skills required.', 'Develop and maintain web applications, conduct code reviews, mentor junior developers, and participate in architectural decisions.', 'Competitive salary, health insurance, performance bonus, work from home flexibility, professional development opportunities.', 'Jakarta', 'DKI Jakarta', 1, 'full_time', 'senior', 15000000, 25000000, 'IDR', 0, NULL, NULL, 'active', NULL, NULL, NULL, 0, 0, 0, '2026-01-20 09:38:16', '2026-01-20 09:38:16', '2026-01-20 09:38:16', NULL),
(38, 1, 'Senior Developer Position 2', 'senior-developer-position-2', 'We are looking for an experienced senior developer with at least 5 years of experience in modern web technologies including TypeScript, React, and backend frameworks.', 'Minimum 5 years of experience with TypeScript, React, Go, and PostgreSQL. Strong problem-solving skills required.', 'Develop and maintain web applications, conduct code reviews, mentor junior developers, and participate in architectural decisions.', 'Competitive salary, health insurance, performance bonus, work from home flexibility, professional development opportunities.', 'Jakarta', 'DKI Jakarta', 1, 'full_time', 'senior', 15000000, 25000000, 'IDR', 0, NULL, NULL, 'active', NULL, NULL, NULL, 0, 0, 0, '2026-01-20 09:38:16', '2026-01-20 09:38:16', '2026-01-20 09:38:16', NULL),
(39, 1, 'Senior Developer Position 3', 'senior-developer-position-3', 'We are looking for an experienced senior developer with at least 5 years of experience in modern web technologies including TypeScript, React, and backend frameworks.', 'Minimum 5 years of experience with TypeScript, React, Go, and PostgreSQL. Strong problem-solving skills required.', 'Develop and maintain web applications, conduct code reviews, mentor junior developers, and participate in architectural decisions.', 'Competitive salary, health insurance, performance bonus, work from home flexibility, professional development opportunities.', 'Jakarta', 'DKI Jakarta', 1, 'full_time', 'senior', 15000000, 25000000, 'IDR', 0, NULL, NULL, 'active', NULL, NULL, NULL, 0, 0, 0, '2026-01-20 09:38:16', '2026-01-20 09:38:16', '2026-01-20 09:38:16', NULL),
(40, 1, 'Senior Developer Position 4', 'senior-developer-position-4', 'We are looking for an experienced senior developer with at least 5 years of experience in modern web technologies including TypeScript, React, and backend frameworks.', 'Minimum 5 years of experience with TypeScript, React, Go, and PostgreSQL. Strong problem-solving skills required.', 'Develop and maintain web applications, conduct code reviews, mentor junior developers, and participate in architectural decisions.', 'Competitive salary, health insurance, performance bonus, work from home flexibility, professional development opportunities.', 'Jakarta', 'DKI Jakarta', 1, 'full_time', 'senior', 15000000, 25000000, 'IDR', 0, NULL, NULL, 'active', NULL, NULL, NULL, 0, 0, 0, '2026-01-20 09:38:16', '2026-01-20 09:38:16', '2026-01-20 09:38:16', NULL),
(41, 1, 'Senior Developer Position 5', 'senior-developer-position-5', 'We are looking for an experienced senior developer with at least 5 years of experience in modern web technologies including TypeScript, React, and backend frameworks.', 'Minimum 5 years of experience with TypeScript, React, Go, and PostgreSQL. Strong problem-solving skills required.', 'Develop and maintain web applications, conduct code reviews, mentor junior developers, and participate in architectural decisions.', 'Competitive salary, health insurance, performance bonus, work from home flexibility, professional development opportunities.', 'Jakarta', 'DKI Jakarta', 1, 'full_time', 'senior', 15000000, 25000000, 'IDR', 0, NULL, NULL, 'active', NULL, NULL, NULL, 0, 0, 0, '2026-01-20 09:38:16', '2026-01-20 09:38:16', '2026-01-20 09:38:16', NULL),
(42, 1, 'Senior Developer Position 6', 'senior-developer-position-6', 'We are looking for an experienced senior developer with at least 5 years of experience in modern web technologies including TypeScript, React, and backend frameworks.', 'Minimum 5 years of experience with TypeScript, React, Go, and PostgreSQL. Strong problem-solving skills required.', 'Develop and maintain web applications, conduct code reviews, mentor junior developers, and participate in architectural decisions.', 'Competitive salary, health insurance, performance bonus, work from home flexibility, professional development opportunities.', 'Jakarta', 'DKI Jakarta', 1, 'full_time', 'senior', 15000000, 25000000, 'IDR', 0, NULL, NULL, 'active', NULL, NULL, NULL, 0, 0, 0, '2026-01-20 09:38:16', '2026-01-20 09:38:16', '2026-01-20 09:38:16', NULL),
(43, 1, 'Senior Developer Position 7', 'senior-developer-position-7', 'We are looking for an experienced senior developer with at least 5 years of experience in modern web technologies including TypeScript, React, and backend frameworks.', 'Minimum 5 years of experience with TypeScript, React, Go, and PostgreSQL. Strong problem-solving skills required.', 'Develop and maintain web applications, conduct code reviews, mentor junior developers, and participate in architectural decisions.', 'Competitive salary, health insurance, performance bonus, work from home flexibility, professional development opportunities.', 'Jakarta', 'DKI Jakarta', 1, 'full_time', 'senior', 15000000, 25000000, 'IDR', 0, NULL, NULL, 'active', NULL, NULL, NULL, 1, 1, 1, '2026-01-20 09:38:16', '2026-01-20 09:38:16', '2026-01-23 15:07:00', NULL),
(44, 1, 'Senior Developer Position 8', 'senior-developer-position-8', 'We are looking for an experienced senior developer with at least 5 years of experience in modern web technologies including TypeScript, React, and backend frameworks.', 'Minimum 5 years of experience with TypeScript, React, Go, and PostgreSQL. Strong problem-solving skills required.', 'Develop and maintain web applications, conduct code reviews, mentor junior developers, and participate in architectural decisions.', 'Competitive salary, health insurance, performance bonus, work from home flexibility, professional development opportunities.', 'Jakarta', 'DKI Jakarta', 1, 'full_time', 'senior', 15000000, 25000000, 'IDR', 0, NULL, NULL, 'active', NULL, NULL, NULL, 0, 0, 0, '2026-01-20 09:38:16', '2026-01-20 09:38:16', '2026-01-20 09:38:16', NULL),
(45, 1, 'Senior Developer Position 9', 'senior-developer-position-9', 'We are looking for an experienced senior developer with at least 5 years of experience in modern web technologies including TypeScript, React, and backend frameworks.', 'Minimum 5 years of experience with TypeScript, React, Go, and PostgreSQL. Strong problem-solving skills required.', 'Develop and maintain web applications, conduct code reviews, mentor junior developers, and participate in architectural decisions.', 'Competitive salary, health insurance, performance bonus, work from home flexibility, professional development opportunities.', 'Jakarta', 'DKI Jakarta', 1, 'full_time', 'senior', 15000000, 25000000, 'IDR', 0, NULL, NULL, 'active', NULL, NULL, NULL, 0, 0, 0, '2026-01-20 09:38:16', '2026-01-20 09:38:16', '2026-01-20 09:38:16', NULL),
(46, 1, 'Extra Job Beyond Quota', 'extra-job-beyond-quota', 'This job should fail because we have exhausted all free quota and have not paid for additional quota.', 'Experience required', 'Development work', 'Benefits provided', 'Jakarta', 'DKI Jakarta', 0, 'full_time', 'senior', NULL, NULL, 'IDR', 0, NULL, NULL, 'active', NULL, NULL, NULL, 0, 0, 0, '2026-01-20 09:38:16', '2026-01-20 09:38:16', '2026-01-20 09:38:16', NULL),
(47, 1, 'Frontend Developer - Draft Test', 'frontend-developer-draft-test', 'Ini adalah lowongan test yang disimpan sebagai draft. Deskripsi minimal 50 karakter untuk validasi.', 'Persyaratan test minimal 30 karakter untuk validasi.', NULL, NULL, 'Jakarta', 'DKI Jakarta', 0, 'full_time', 'junior', NULL, NULL, 'IDR', 0, NULL, NULL, 'draft', NULL, NULL, NULL, 0, 0, 0, NULL, '2026-01-20 09:55:43', '2026-01-20 09:55:43', NULL),
(48, 1, 'Sales Marketing', 'sales-marketing', 'Tanggung Jawab Pekerjaan\n\nMencapai target\nMenjaga hubungan dengan pelanggan\nRiset pemasaran\nMenjalankan promosi perusahaan\nSyarat & Keahlian\n\nMin lulusan SMA\nUsia 20 - 45 tahun\nPengalaman tidak diutamakan\nPria/Wanita\nRajin dan bertanggung jawab\n\nPengalaman\n\nTidak diutamakan\nBenefit\n\nKomisi, Piknik', 'Max age 28\nMinimum Bachelor Degree in any major\nMin GPA 3.00\nFresh graduates are welcome to apply\nInterest to work in broadcasting industry\nHave a good looking\nHave a good communication & presentation skill\nHave a good networking skill', NULL, NULL, 'Sidaorjo', 'Jawa Timur', 0, 'contract', 'junior', 4000000, 4000000, 'IDR', 1, '2026-02-12', NULL, 'draft', NULL, NULL, NULL, 0, 0, 0, NULL, '2026-01-20 10:00:06', '2026-01-20 10:00:06', NULL),
(49, 1, 'Graphic Designer', 'graphic-designer', 'This is a full-time remote Senior Graphic Designer role starting in March. You will create polished, globally appealing designs for social media, e-commerce, and marketing campaigns, working closely with an international team. Intermediate English communication is required.\nWork Type: Remote\nEmployment: Full-time\nLevel: Senior\nEnglish: Intermediate (daily team communication)\nSalary: HKD 5,000 / month\n≈ IDR 10,800,000 / month\n\nKey Responsibilities\nCreate high-quality visual designs for social media, ads, campaigns, and e-commerce\nDesign using Canva as the main tool for all primary visuals\nEnsure designs meet international beauty brand standards\nMaintain brand consistency across all platforms\nCollaborate with cross-border teams and respond to feedback\nPrepare final assets for digital publishing', 'Senior-level experience (3+ years) as a Graphic Designer\nExcellent Canva skills (primary tool for all main designs)\nStrong understanding of global / international design aesthetics\nAbility to design for international audiences, not local-only styles\nStrong sense of typography, layout, spacing, and color\nIntermediate English (spoken & written) for team communication\nAble to work independently and meet deadlines\nPortfolio showing modern, clean, internationally relevant designs', NULL, NULL, 'Jakarta', 'DKI Jakarta', 0, 'contract', 'junior', 2400000, 3000000, 'IDR', 1, '2026-02-21', NULL, 'active', NULL, NULL, NULL, 1, 0, 2, '2026-01-21 13:52:32', '2026-01-21 13:52:32', '2026-01-23 14:32:51', NULL),
(50, 1, 'Backend Developer - Email Notification Test', 'backend-developer-email-notification-test', 'Testing email notification feature untuk job posting', 'Menguasai Go, MySQL, dan email integration', NULL, NULL, 'Jakarta', 'DKI Jakarta', 0, 'full_time', 'mid', NULL, NULL, 'IDR', 0, NULL, NULL, 'active', NULL, NULL, NULL, 0, 0, 0, '2026-01-22 03:58:55', '2026-01-22 03:58:55', '2026-01-22 03:58:55', NULL),
(51, 1, 'Frontend Developer - Test Email v2', 'frontend-developer-test-email-v2', 'Testing email notification feature dengan debug logs untuk memastikan email terkirim dengan benar ke alamat company yang terdaftar', NULL, NULL, NULL, 'Bandung', 'Jawa Barat', 0, 'full_time', 'junior', NULL, NULL, 'IDR', 0, NULL, NULL, 'active', NULL, NULL, NULL, 0, 0, 0, '2026-01-22 04:00:54', '2026-01-22 04:00:54', '2026-01-22 04:00:54', NULL),
(52, 1, 'QA Engineer - Email Test With Correct Binary', 'qa-engineer-email-test-with-correct-binary', 'Testing email notification feature dengan binary yang sudah terupdate untuk memastikan email terkirim dengan sempurna', NULL, NULL, NULL, 'Surabaya', 'Jawa Timur', 0, 'full_time', 'mid', NULL, NULL, 'IDR', 0, NULL, NULL, 'active', NULL, NULL, NULL, 0, 0, 0, '2026-01-22 04:02:07', '2026-01-22 04:02:07', '2026-01-22 04:02:07', NULL),
(53, 1, 'DevOps Engineer - Final Email Notification Test', 'devops-engineer-final-email-notification-test', 'Testing email notification dengan background context agar tidak ter-cancel saat request selesai dan email dapat terkirim dengan sempurna', NULL, NULL, NULL, 'Yogyakarta', 'DI Yogyakarta', 0, 'full_time', 'senior', NULL, NULL, 'IDR', 0, NULL, NULL, 'active', NULL, NULL, NULL, 1, 0, 0, '2026-01-22 04:08:45', '2026-01-22 04:08:45', '2026-01-25 10:25:44', NULL),
(54, 1, 'Accounting', 'accounting', 'Tentang Kami Twister Communications adalah creative agency yang fokus pada brand activation, event management, dan design services untuk brand multinasional dan lokal di Indonesia.\n\nDeskripsi Pekerjaan Kami mencari Accounting Staff yang teliti untuk mendukung Divisi Finance, Accounting and Tax dalam mengelola tugas Accounting secara keseluruhan sesuai dengan standar operasional Perusahaan.\n\nTugas & Tanggungjawab\n\nMelakukan pencatatan transaksi keuangan harian secara akurat dan tepat waktu\n\nMenyusun dan menginput jurnal akuntansi ke dalam sistem\n\nMelakukan rekonsiliasi bank serta memastikan kesesuaian saldo kas dan bank\n\nMembantu proses closing laporan keuangan bulanan dan tahunan\n\nMenyusun dan memeriksa laporan keuangan (Laba Rugi, Neraca, Arus Kas)\n\nMenyiapkan data pendukung untuk kebutuhan audit internal maupun eksternal\n\nMembantu pengelolaan dan pelaporan perpajakan (PPN, PPh 21, PPh 23, dll)\n\nMengarsipkan dokumen keuangan dan perpajakan secara rapi dan sistematis\n\nMengoperasikan dan memastikan data pada software akuntansi selalu update\n\nMemastikan seluruh proses akuntansi berjalan sesuai SOP dan kebijakan perusahaan\n\nTugas terkait lainnya sesuai yang diberikan', 'Pendidikan minimal S1 Akuntansi\n\nMempunyai pengalaman 1 tahun sebagai Accounting Staff lebih diutamakan\n\nMengetahui program Accounting (contoh: Accurate, ERP)\n\nMahir rumus excel seperti vlookup & pivot\n\nMenguasai Ms Office\n\nTeliti, Jujur, Cekatan dan bertanggung jawab\n\nMampu bekerja dalam lingkungan yang dinamis\n\nMampu bekerja dalam team\n\nDiutamakan berpengalaman pada bidang Jasa\n\nFresh Graduate are Welcome', NULL, NULL, 'Sidaorjo', 'Jawa Timur', 0, 'full_time', 'junior', NULL, NULL, 'IDR', 1, '2026-02-22', NULL, 'active', NULL, NULL, NULL, 1, 0, 0, '2026-01-22 05:52:04', '2026-01-22 05:52:04', '2026-01-25 12:11:37', NULL),
(55, 1, 'Full Stack Developer - Email Test dengan Logging Detail', 'full-stack-developer-email-test-dengan-logging-detail', 'Kami mencari Full Stack Developer yang berpengalaman untuk bergabung dengan tim kami. Posisi ini akan fokus pada pengembangan aplikasi web modern menggunakan teknologi terkini.', 'Minimal 2 tahun pengalaman, menguasai React, Node.js, dan database', NULL, NULL, 'Jakarta', 'DKI Jakarta', 0, 'full_time', 'mid', NULL, NULL, 'IDR', 0, NULL, NULL, 'active', NULL, NULL, NULL, 0, 0, 0, '2026-01-22 06:17:52', '2026-01-22 06:17:52', '2026-01-22 06:17:52', NULL),
(56, 1, 'Software Engineer SALARY TEST', 'software-engineer-salary-test', 'This is a test job to verify salary display works correctly on the frontend. We are looking for a talented software engineer.', '3 years experience with Go or TypeScript', 'Develop and maintain backend services', 'Health insurance, remote work', 'Jakarta', 'DKI Jakarta', 1, 'full_time', 'mid', 15000000, 25000000, 'IDR', 1, NULL, NULL, 'active', NULL, NULL, NULL, 2, 1, 1, '2026-01-22 06:45:16', '2026-01-22 06:45:16', '2026-01-25 12:11:26', NULL),
(57, 7, 'Senior Flutter Developer', 'senior-flutter-developer', 'Kami mencari Senior Mobile Apps Developer berpengalaman yang memiliki keahlian kuat dalam pengembangan aplikasi mobile menggunakan framework Flutter. Kandidat akan bertanggung jawab dalam merancang, mengembangkan, mengoptimalkan, dan memelihara aplikasi mobile berkualitas tinggi untuk platform Android dan iOS, serta berkolaborasi dengan tim backend, UI/UX, dan product team untuk menghasilkan solusi digital yang scalable dan efisien.\n\nTanggung Jawab:\n\nMengembangkan aplikasi mobile menggunakan Flutter dengan performa tinggi dan UI/UX yang optimal\n\nMendesain arsitektur aplikasi yang scalable, maintainable, dan clean code\n\nMengintegrasikan aplikasi dengan REST API / Web Service\n\nMelakukan debugging, testing, dan optimasi performa aplikasi\n\nMembimbing developer junior dan melakukan code review\n\nBerkolaborasi dengan tim lintas divisi dalam proses pengembangan produk\n\nMengikuti perkembangan teknologi mobile dan best practice industri\n\nKualifikasi:\n\nPengalaman minimal 3–5 tahun dalam pengembangan aplikasi mobile\n\nPengalaman kuat menggunakan Flutter dan Dart\n\nMemahami state management (Bloc, Provider, Riverpod, atau sejenisnya)\n\nBerpengalaman integrasi API, JSON parsing, dan authentication\n\nMemahami prinsip Clean Architecture / MVC / MVVM\n\nTerbiasa menggunakan Git version control\n\nMemiliki kemampuan problem solving dan komunikasi yang baik\n\nNilai tambah jika memiliki pengalaman publikasi aplikasi di Play Store / App Store\n\nNilai Tambah (Optional):\n\nPengalaman dengan Firebase (Auth, Firestore, FCM, Analytics)\n\nPengalaman CI/CD mobile apps\n\nPengalaman native Android (Kotlin/Java) atau iOS (Swift)\n\nPengalaman pengembangan real-time apps (WebSocket / MQTT)', 'Persyaratan:\n\nPendidikan minimal S1 Teknik Informatika, Sistem Informasi, atau bidang terkait (atau pengalaman kerja setara)\n\nPengalaman minimal 3–5 tahun dalam pengembangan aplikasi mobile\n\nPengalaman profesional menggunakan Flutter (Dart) dalam pengembangan production app\n\nMemahami konsep State Management (Bloc, Provider, Riverpod, GetX, atau sejenisnya)\n\nBerpengalaman dalam integrasi REST API / Web Service\n\nMemahami konsep Clean Architecture, SOLID Principle, dan Design Pattern\n\nTerbiasa menggunakan Git Version Control (GitHub / GitLab / Bitbucket)\n\nMemahami proses deployment aplikasi ke Google Play Store dan Apple App Store\n\nMampu melakukan debugging, profiling, dan optimasi performa aplikasi\n\nMemiliki kemampuan analisa dan problem solving yang kuat\n\nMampu bekerja secara tim maupun individu\n\nMemiliki komunikasi yang baik dan mampu bekerja dalam lingkungan agile', NULL, NULL, 'Ngagel', 'Surabaya', 0, 'full_time', 'mid', 4000000, NULL, 'IDR', 1, '2026-03-11', NULL, 'active', NULL, NULL, NULL, 1, 1, 0, '2026-02-11 14:54:30', '2026-02-11 14:54:30', '2026-02-11 14:56:30', NULL);

-- --------------------------------------------------------

--
-- Table structure for table `job_shares`
--

CREATE TABLE `job_shares` (
  `id` bigint(20) UNSIGNED NOT NULL,
  `job_id` bigint(20) UNSIGNED NOT NULL,
  `user_id` bigint(20) UNSIGNED DEFAULT NULL COMMENT 'Applicant user_id who shared (optional)',
  `platform` varchar(50) DEFAULT NULL COMMENT 'Platform: whatsapp, telegram, facebook, twitter, copy_link, etc',
  `shared_at` timestamp NOT NULL DEFAULT current_timestamp()
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

--
-- Dumping data for table `job_shares`
--

INSERT INTO `job_shares` (`id`, `job_id`, `user_id`, `platform`, `shared_at`) VALUES
(1, 49, 21, 'whatsapp', '2026-01-23 14:32:38'),
(2, 49, 21, 'telegram', '2026-01-23 14:32:51'),
(3, 43, 21, '', '2026-01-23 15:07:00'),
(4, 56, 20, 'copy_link', '2026-01-25 09:07:44');

-- --------------------------------------------------------

--
-- Table structure for table `job_skills`
--

CREATE TABLE `job_skills` (
  `id` bigint(20) UNSIGNED NOT NULL,
  `job_id` bigint(20) UNSIGNED NOT NULL,
  `skill_name` varchar(100) NOT NULL,
  `is_required` tinyint(1) NOT NULL DEFAULT 1
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

--
-- Dumping data for table `job_skills`
--

INSERT INTO `job_skills` (`id`, `job_id`, `skill_name`, `is_required`) VALUES
(1, 1, 'Go', 1),
(2, 1, 'Python', 1),
(3, 1, 'Docker', 1),
(4, 3, 'Go', 1),
(5, 3, 'Python', 1),
(6, 3, 'Docker', 1),
(7, 4, 'Go', 1),
(8, 4, 'PostgreSQL', 1),
(9, 4, 'Docker', 1),
(10, 4, 'Kubernetes', 1),
(11, 4, 'Redis', 1),
(12, 5, 'React', 1),
(13, 5, 'Node.js', 1),
(14, 5, 'MongoDB', 1),
(15, 5, 'Docker', 1),
(16, 6, 'Figma', 1),
(17, 6, 'UI Design', 1),
(18, 6, 'UX Research', 1),
(19, 6, 'Prototyping', 1),
(20, 8, 'React', 1),
(21, 8, 'TypeScript', 1),
(22, 8, 'JavaScript', 1),
(23, 9, 'React', 1),
(24, 9, 'TypeScript', 1),
(25, 9, 'JavaScript', 1),
(26, 10, 'React', 1),
(27, 10, 'TypeScript', 1),
(28, 11, 'React', 1),
(29, 11, 'TypeScript', 1),
(30, 17, 'TypeScript', 1),
(31, 17, 'React', 1),
(32, 17, 'Go', 1),
(33, 17, 'PostgreSQL', 1),
(34, 18, 'TypeScript', 1),
(35, 18, 'React', 1),
(36, 18, 'Go', 1),
(37, 18, 'PostgreSQL', 1),
(38, 19, 'TypeScript', 1),
(39, 19, 'React', 1),
(40, 19, 'Go', 1),
(41, 19, 'PostgreSQL', 1),
(42, 20, 'TypeScript', 1),
(43, 20, 'React', 1),
(44, 20, 'Go', 1),
(45, 20, 'PostgreSQL', 1),
(46, 21, 'TypeScript', 1),
(47, 21, 'React', 1),
(48, 21, 'Go', 1),
(49, 21, 'PostgreSQL', 1),
(50, 22, 'TypeScript', 1),
(51, 22, 'React', 1),
(52, 22, 'Go', 1),
(53, 22, 'PostgreSQL', 1),
(54, 23, 'TypeScript', 1),
(55, 23, 'React', 1),
(56, 23, 'Go', 1),
(57, 23, 'PostgreSQL', 1),
(58, 24, 'TypeScript', 1),
(59, 24, 'React', 1),
(60, 24, 'Go', 1),
(61, 24, 'PostgreSQL', 1),
(62, 25, 'TypeScript', 1),
(63, 25, 'React', 1),
(64, 25, 'Go', 1),
(65, 25, 'PostgreSQL', 1),
(66, 26, 'TypeScript', 1),
(67, 26, 'React', 1),
(68, 26, 'Go', 1),
(69, 26, 'PostgreSQL', 1),
(70, 50, 'Golang', 1),
(71, 50, 'MySQL', 1),
(72, 50, 'SMTP', 1),
(73, 51, 'React', 1),
(74, 51, 'TypeScript', 1),
(75, 52, 'Testing', 1),
(76, 52, 'Automation', 1),
(77, 53, 'Docker', 1),
(78, 53, 'Kubernetes', 1),
(79, 53, 'CI/CD', 1),
(80, 55, 'React', 1),
(81, 55, 'Node.js', 1),
(82, 55, 'PostgreSQL', 1),
(83, 55, 'Docker', 1),
(84, 56, 'Go', 1),
(85, 56, 'TypeScript', 1);

-- --------------------------------------------------------

--
-- Table structure for table `job_views`
--

CREATE TABLE `job_views` (
  `id` bigint(20) UNSIGNED NOT NULL,
  `job_id` bigint(20) UNSIGNED NOT NULL,
  `user_id` bigint(20) UNSIGNED NOT NULL COMMENT 'Applicant user_id who viewed',
  `viewed_at` timestamp NOT NULL DEFAULT current_timestamp()
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

--
-- Dumping data for table `job_views`
--

INSERT INTO `job_views` (`id`, `job_id`, `user_id`, `viewed_at`) VALUES
(1, 49, 21, '2026-01-23 14:28:16'),
(3, 43, 21, '2026-01-23 15:06:59'),
(5, 56, 20, '2026-01-25 09:07:19'),
(7, 53, 20, '2026-01-25 10:25:44'),
(12, 56, 21, '2026-01-25 12:11:26'),
(13, 54, 21, '2026-01-25 12:11:37'),
(21, 57, 20, '2026-02-11 14:56:22');

-- --------------------------------------------------------

--
-- Table structure for table `notifications`
--

CREATE TABLE `notifications` (
  `id` bigint(20) UNSIGNED NOT NULL,
  `user_id` bigint(20) UNSIGNED NOT NULL,
  `type` varchar(50) NOT NULL,
  `title` varchar(255) NOT NULL,
  `message` text NOT NULL,
  `data` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL CHECK (json_valid(`data`)),
  `is_read` tinyint(1) NOT NULL DEFAULT 0,
  `read_at` timestamp NULL DEFAULT NULL,
  `created_at` timestamp NOT NULL DEFAULT current_timestamp()
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- --------------------------------------------------------

--
-- Table structure for table `partner_commissions`
--

CREATE TABLE `partner_commissions` (
  `id` bigint(20) UNSIGNED NOT NULL,
  `partner_id` bigint(20) UNSIGNED NOT NULL,
  `referral_id` bigint(20) UNSIGNED NOT NULL COMMENT 'partner_referrals.id',
  `payment_id` bigint(20) UNSIGNED NOT NULL COMMENT 'payments.id',
  `company_id` bigint(20) UNSIGNED NOT NULL,
  `transaction_amount` bigint(20) NOT NULL COMMENT 'Original payment amount (IDR)',
  `commission_rate` decimal(5,2) NOT NULL COMMENT 'Rate at time of transaction (e.g., 40.00)',
  `commission_amount` bigint(20) NOT NULL COMMENT 'Calculated commission (IDR)',
  `job_quota` int(11) NOT NULL COMMENT 'Number of job posts purchased',
  `status` enum('pending','approved','paid','cancelled') NOT NULL DEFAULT 'pending',
  `approved_by` bigint(20) UNSIGNED DEFAULT NULL,
  `approved_at` timestamp NULL DEFAULT NULL,
  `paid_at` timestamp NULL DEFAULT NULL,
  `payout_id` bigint(20) UNSIGNED DEFAULT NULL COMMENT 'Link to partner_payouts when paid',
  `notes` text DEFAULT NULL,
  `created_at` timestamp NOT NULL DEFAULT current_timestamp(),
  `updated_at` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp()
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

--
-- Triggers `partner_commissions`
--
DELIMITER $$
CREATE TRIGGER `after_commission_insert` AFTER INSERT ON `partner_commissions` FOR EACH ROW BEGIN
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
END
$$
DELIMITER ;
DELIMITER $$
CREATE TRIGGER `after_commission_update` AFTER UPDATE ON `partner_commissions` FOR EACH ROW BEGIN
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
END
$$
DELIMITER ;

-- --------------------------------------------------------

--
-- Table structure for table `partner_payouts`
--

CREATE TABLE `partner_payouts` (
  `id` bigint(20) UNSIGNED NOT NULL,
  `partner_id` bigint(20) UNSIGNED NOT NULL,
  `amount` bigint(20) NOT NULL COMMENT 'Payout amount (IDR)',
  `bank_name` varchar(100) NOT NULL,
  `bank_account_number` varchar(50) NOT NULL,
  `bank_account_holder` varchar(255) NOT NULL,
  `status` enum('pending','processing','completed','failed','cancelled') NOT NULL DEFAULT 'pending',
  `transfer_ref` varchar(100) DEFAULT NULL COMMENT 'Bank transfer reference number',
  `failure_reason` text DEFAULT NULL,
  `requested_at` timestamp NOT NULL DEFAULT current_timestamp(),
  `processed_by` bigint(20) UNSIGNED DEFAULT NULL,
  `processed_at` timestamp NULL DEFAULT NULL,
  `completed_at` timestamp NULL DEFAULT NULL,
  `notes` text DEFAULT NULL,
  `created_at` timestamp NOT NULL DEFAULT current_timestamp(),
  `updated_at` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp()
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

--
-- Dumping data for table `partner_payouts`
--

INSERT INTO `partner_payouts` (`id`, `partner_id`, `amount`, `bank_name`, `bank_account_number`, `bank_account_holder`, `status`, `transfer_ref`, `failure_reason`, `requested_at`, `processed_by`, `processed_at`, `completed_at`, `notes`, `created_at`, `updated_at`) VALUES
(1, 1, 5000000, 'Bank Central Asia', '1234567890', 'Ahmad Pratama', 'completed', 'https://example.com/proof/transfer123.jpg', NULL, '2026-02-05 19:28:49', NULL, NULL, '2026-02-05 19:29:05', NULL, '2026-02-05 19:28:49', '2026-02-05 19:29:05');

-- --------------------------------------------------------

--
-- Table structure for table `partner_referrals`
--

CREATE TABLE `partner_referrals` (
  `id` bigint(20) UNSIGNED NOT NULL,
  `partner_id` bigint(20) UNSIGNED NOT NULL COMMENT 'referral_partners.id',
  `company_id` bigint(20) UNSIGNED NOT NULL COMMENT 'companies.id',
  `referral_code_used` varchar(20) NOT NULL COMMENT 'The code used at registration',
  `registered_at` timestamp NOT NULL DEFAULT current_timestamp(),
  `is_verified` tinyint(1) NOT NULL DEFAULT 0 COMMENT 'Company account has been verified',
  `first_payment_at` timestamp NULL DEFAULT NULL COMMENT 'When company made first purchase',
  `notes` text DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

--
-- Dumping data for table `partner_referrals`
--

INSERT INTO `partner_referrals` (`id`, `partner_id`, `company_id`, `referral_code_used`, `registered_at`, `is_verified`, `first_payment_at`, `notes`) VALUES
(2, 3, 7, 'SAPUC4EB', '2026-02-09 13:08:22', 0, NULL, NULL);

--
-- Triggers `partner_referrals`
--
DELIMITER $$
CREATE TRIGGER `after_referral_insert` AFTER INSERT ON `partner_referrals` FOR EACH ROW BEGIN
  UPDATE `referral_partners`
  SET 
    `total_referrals` = `total_referrals` + 1,
    `updated_at` = NOW()
  WHERE `id` = NEW.partner_id;
END
$$
DELIMITER ;

-- --------------------------------------------------------

--
-- Table structure for table `password_resets`
--

CREATE TABLE `password_resets` (
  `id` bigint(20) UNSIGNED NOT NULL,
  `user_id` bigint(20) UNSIGNED NOT NULL,
  `token` varchar(255) NOT NULL,
  `expires_at` datetime NOT NULL,
  `created_at` datetime DEFAULT current_timestamp()
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

--
-- Dumping data for table `password_resets`
--

INSERT INTO `password_resets` (`id`, `user_id`, `token`, `expires_at`, `created_at`) VALUES
(2, 24, '669b9969feaca3acd674ccbe4a7cdf883adc329e9bad1d7a350190f69832e3c2', '2026-02-05 20:38:02', '2026-02-05 19:38:02');

-- --------------------------------------------------------

--
-- Table structure for table `password_reset_tokens`
--

CREATE TABLE `password_reset_tokens` (
  `id` bigint(20) UNSIGNED NOT NULL,
  `email` varchar(255) NOT NULL,
  `token` varchar(255) NOT NULL,
  `expires_at` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  `used_at` timestamp NULL DEFAULT NULL,
  `created_at` timestamp NOT NULL DEFAULT current_timestamp()
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

--
-- Dumping data for table `password_reset_tokens`
--

INSERT INTO `password_reset_tokens` (`id`, `email`, `token`, `expires_at`, `used_at`, `created_at`) VALUES
(1, 'info@karyadeveloperindonesia.com', 'c09200435d7cc982616d4fc6f2de7e74ec300416388a193cdc481bdf35a9b1d4', '2026-01-21 14:30:46', '2026-01-21 14:30:46', '2026-01-21 14:29:59'),
(2, 'info@karyadeveloperindonesia.com', 'cd0fc0faa256ec0bbead7089678b7c04c9ea0dc3fdb0503db9133dbfce83e87e', '2026-01-21 14:33:12', '2026-01-21 14:33:12', '2026-01-21 14:32:52'),
(3, 'info@karyadeveloperindonesia.com', '5fca4dd958bbcc3467118700fe27bcf416e0feb8e79c5d2f0345c8395b53cfc3', '2026-01-21 15:35:36', NULL, '2026-01-21 14:35:36'),
(4, 'info@karyadeveloperindonesia.com', 'd05b1c5e3be84df24506cdc46a07cb2f973a8f2c5efc1df2b2101149a62be8c8', '2026-01-21 15:37:00', NULL, '2026-01-21 14:37:00'),
(5, 'info@karyadeveloperindonesia.com', '80394eb1d1edb7a5d31521137bbd54b0b928f20129e5b5011721cf12be0e7be5', '2026-01-21 16:05:46', NULL, '2026-01-21 15:05:46'),
(6, 'info@karyadeveloperindonesia.com', '8bfdb64aafe8d05c4f74e31d8298840cadfa0f61cc816334a7ef3430ca72750b', '2026-01-21 16:06:37', NULL, '2026-01-21 15:06:37'),
(7, 'info@karyadeveloperindonesia.com', 'a7ee82140991f13b34e41d8e02f03aade50ebfdc706085bc07edab5408695791', '2026-01-21 16:08:19', NULL, '2026-01-21 15:08:19'),
(8, 'craftgirlsssshopping@gmail.com', '59efde66def58e29054f7058d3bcb738b7b38af88d517abb988049ef916de548', '2026-01-26 07:55:44', NULL, '2026-01-26 07:50:44');

-- --------------------------------------------------------

--
-- Table structure for table `payments`
--

CREATE TABLE `payments` (
  `id` bigint(20) UNSIGNED NOT NULL,
  `company_id` bigint(20) UNSIGNED NOT NULL,
  `job_id` bigint(20) UNSIGNED DEFAULT NULL,
  `package_id` varchar(50) DEFAULT NULL,
  `quota_amount` int(11) NOT NULL DEFAULT 1,
  `amount` bigint(20) NOT NULL DEFAULT 15000,
  `proof_image_url` varchar(500) DEFAULT NULL,
  `status` enum('pending','confirmed','rejected') NOT NULL DEFAULT 'pending',
  `note` text DEFAULT NULL,
  `confirmed_by_id` bigint(20) UNSIGNED DEFAULT NULL,
  `submitted_at` timestamp NOT NULL DEFAULT current_timestamp(),
  `confirmed_at` timestamp NULL DEFAULT NULL,
  `created_at` timestamp NOT NULL DEFAULT current_timestamp(),
  `updated_at` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp()
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

--
-- Dumping data for table `payments`
--

INSERT INTO `payments` (`id`, `company_id`, `job_id`, `package_id`, `quota_amount`, `amount`, `proof_image_url`, `status`, `note`, `confirmed_by_id`, `submitted_at`, `confirmed_at`, `created_at`, `updated_at`) VALUES
(1, 1, NULL, NULL, 2, 15000, '/docs/payments/7/proof_1768893950.png', 'confirmed', 'Pembayaran telah diverifikasi', 1, '2026-01-20 07:25:50', '2026-01-20 07:29:01', '2026-01-20 07:25:50', '2026-02-02 04:05:51'),
(2, 1, NULL, NULL, 2, 15000, '/docs/payments/7/proof_1768894335.png', 'pending', NULL, NULL, '2026-01-20 07:32:15', NULL, '2026-01-20 07:32:15', '2026-02-02 04:05:51'),
(3, 1, NULL, NULL, 2, 15000, '/docs/payments/7/proof_1768894343.png', 'confirmed', 'Pembayaran diterima', 1, '2026-01-20 07:32:23', '2026-01-20 07:32:32', '2026-01-20 07:32:23', '2026-02-02 04:05:51'),
(4, 1, NULL, 'pack10', 12, 100000, '/docs/payments/7/proof_1768918594.txt', 'confirmed', 'Pack10 verified - 12 quota will be added', 1, '2026-01-20 14:16:34', '2026-01-20 14:20:48', '2026-01-20 14:16:34', '2026-02-02 03:43:55'),
(5, 1, NULL, 'pack5', 5, 50000, '/docs/payments/7/proof_1769003589.png', 'confirmed', 'Pembayaran telah diverifikasi dan disetujui. Terima kasih!', 1, '2026-01-21 13:53:09', '2026-01-21 16:40:27', '2026-01-21 13:53:09', '2026-02-02 03:43:55'),
(6, 2, NULL, NULL, 5, 50000, NULL, 'pending', NULL, NULL, '2026-02-02 03:44:14', NULL, '2026-02-02 03:44:14', '2026-02-02 04:05:51'),
(7, 3, NULL, NULL, 3, 30000, NULL, 'confirmed', 'Test approval via API', 1, '2026-02-01 03:44:14', '2026-02-02 03:45:53', '2026-02-02 03:44:14', '2026-02-02 04:05:51'),
(8, 4, NULL, NULL, 10, 100000, NULL, 'rejected', 'Payment proof is invalid', 1, '2026-01-31 03:44:14', '2026-02-02 03:46:07', '2026-02-02 03:44:14', '2026-02-02 04:05:51'),
(9, 5, NULL, NULL, 8, 75000, NULL, 'pending', NULL, NULL, '2026-02-02 03:46:46', NULL, '2026-02-02 03:46:46', '2026-02-02 04:05:51'),
(10, 1, NULL, NULL, 20, 200000, NULL, 'confirmed', 'Approved with 20 quota(s)', 1, '2026-02-02 00:46:46', '2026-02-02 03:54:16', '2026-02-02 03:46:46', '2026-02-02 04:05:51'),
(11, 2, NULL, NULL, 4, 40000, NULL, 'pending', NULL, NULL, '2026-02-01 22:46:46', NULL, '2026-02-02 03:46:46', '2026-02-02 04:05:51'),
(12, 1, NULL, 'pack5', 5, 50000, '/docs/payments/1/proof_1770005254.txt', 'pending', NULL, NULL, '2026-02-02 04:07:34', NULL, '2026-02-02 04:07:34', '2026-02-02 04:07:34'),
(13, 1, NULL, 'pack5', 5, 50000, '/docs/payments/1/proof_1770005261.txt', 'confirmed', 'Approved - 5 quota package', 1, '2026-02-02 04:07:41', '2026-02-02 04:08:25', '2026-02-02 04:07:41', '2026-02-02 04:08:25'),
(14, 1, NULL, 'pack10', 12, 100000, '/docs/payments/1/proof_1770005343.txt', 'pending', NULL, NULL, '2026-02-02 04:09:03', NULL, '2026-02-02 04:09:03', '2026-02-02 04:09:03'),
(15, 1, NULL, 'pack10', 12, 100000, '/docs/payments/1/proof_1770005356.txt', 'confirmed', 'Pack 10+2 approved', 1, '2026-02-02 04:09:16', '2026-02-02 04:09:31', '2026-02-02 04:09:16', '2026-02-02 04:09:31');

-- --------------------------------------------------------

--
-- Table structure for table `referral_partners`
--

CREATE TABLE `referral_partners` (
  `id` bigint(20) UNSIGNED NOT NULL,
  `user_id` bigint(20) UNSIGNED NOT NULL COMMENT 'Link to users table',
  `referral_code` varchar(20) NOT NULL COMMENT 'Unique referral code e.g., AHMAD2024',
  `commission_rate` decimal(5,2) NOT NULL DEFAULT 40.00 COMMENT 'Commission percentage (40%)',
  `status` enum('active','inactive','pending','suspended','rejected') NOT NULL DEFAULT 'pending',
  `bank_name` varchar(100) DEFAULT NULL,
  `bank_account_number` varchar(50) DEFAULT NULL,
  `bank_account_holder` varchar(255) DEFAULT NULL,
  `is_bank_verified` tinyint(1) NOT NULL DEFAULT 0,
  `total_referrals` int(11) NOT NULL DEFAULT 0 COMMENT 'Total companies referred',
  `total_commission` bigint(20) NOT NULL DEFAULT 0 COMMENT 'Total commission earned (lifetime)',
  `available_balance` bigint(20) NOT NULL DEFAULT 0 COMMENT 'Balance ready for payout',
  `pending_balance` bigint(20) NOT NULL DEFAULT 0 COMMENT 'Commission pending approval',
  `paid_amount` bigint(20) NOT NULL DEFAULT 0 COMMENT 'Total amount paid out',
  `approved_by` bigint(20) UNSIGNED DEFAULT NULL,
  `approved_at` timestamp NULL DEFAULT NULL,
  `notes` text DEFAULT NULL COMMENT 'Admin notes',
  `created_at` timestamp NOT NULL DEFAULT current_timestamp(),
  `updated_at` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp()
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

--
-- Dumping data for table `referral_partners`
--

INSERT INTO `referral_partners` (`id`, `user_id`, `referral_code`, `commission_rate`, `status`, `bank_name`, `bank_account_number`, `bank_account_holder`, `is_bank_verified`, `total_referrals`, `total_commission`, `available_balance`, `pending_balance`, `paid_amount`, `approved_by`, `approved_at`, `notes`, `created_at`, `updated_at`) VALUES
(1, 22, 'AHMAD2024', 40.00, 'active', 'Bank Central Asia', '1234567890', 'Ahmad Pratama', 1, 0, 47850000, 7500000, 0, 40350000, NULL, '2026-02-04 13:06:58', NULL, '2024-01-14 17:00:00', '2026-02-09 13:27:27'),
(2, 23, 'TESTB403', 35.00, 'active', NULL, NULL, NULL, 0, 0, 0, 0, 0, 0, 1, '2026-02-05 19:28:25', NULL, '2026-02-05 12:26:15', '2026-02-05 19:28:25'),
(3, 24, 'SAPUC4EB', 40.00, 'active', NULL, NULL, NULL, 0, 1, 0, 0, 0, 0, NULL, NULL, NULL, '2026-02-05 12:37:05', '2026-02-09 13:27:27');

-- --------------------------------------------------------

--
-- Table structure for table `refresh_tokens`
--

CREATE TABLE `refresh_tokens` (
  `id` bigint(20) UNSIGNED NOT NULL,
  `user_id` bigint(20) UNSIGNED NOT NULL,
  `token_hash` varchar(255) NOT NULL,
  `expires_at` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  `revoked_at` timestamp NULL DEFAULT NULL,
  `device_info` varchar(500) DEFAULT NULL,
  `ip_address` varchar(45) DEFAULT NULL,
  `created_at` timestamp NOT NULL DEFAULT current_timestamp()
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

--
-- Dumping data for table `refresh_tokens`
--

INSERT INTO `refresh_tokens` (`id`, `user_id`, `token_hash`, `expires_at`, `revoked_at`, `device_info`, `ip_address`, `created_at`) VALUES
(5, 2, 'dedea688d4b52753f79934fa2a885093b3e3c4288cbe507380beb7e5c7b872e3', '2026-01-24 02:44:45', NULL, '', '', '2026-01-17 02:44:45'),
(6, 2, '37b6697c2e3078e4b0f617e3408d89a1b8dbcf069a39e029e5ba8e0afb9590a1', '2026-01-24 02:45:11', NULL, '', '', '2026-01-17 02:45:11'),
(7, 3, '2da7ef5f74dfe0a7bce05c4c928a637802ed9780306206f8b99555db54f087a3', '2026-01-24 02:50:35', NULL, '', '', '2026-01-17 02:50:35'),
(8, 3, '6c6d2187d282e50c7c4e6b4fe9b0f106e36742342320a804f21734ef0e77aa47', '2026-01-24 02:50:57', NULL, '', '', '2026-01-17 02:50:57'),
(9, 3, '575affa07d728513fee1f83c9734b761209e0eba1da5e33d0420c799e6ef5c50', '2026-01-24 02:51:04', NULL, '', '', '2026-01-17 02:51:04'),
(10, 2, 'f0376707835cb2efc1ebe6389e6593deff6a074a53c6d8eeeb01202dcf0e365f', '2026-01-24 02:51:37', NULL, '', '', '2026-01-17 02:51:37'),
(11, 2, '8220bad6ba5100b2eb4036f8dc6e108fb452286f11a118d9af361d675dd54940', '2026-01-24 02:51:45', NULL, '', '', '2026-01-17 02:51:45'),
(12, 2, 'b579617f89af74650dfd38a0142ded500d50fd66eae9f23a882f88eadd9d03da', '2026-01-24 02:53:40', NULL, '', '', '2026-01-17 02:53:40'),
(13, 2, 'bdbfd0c0b672d651473c9344146b870b3eb212bb71521cdf6c99857a4125aaf4', '2026-01-24 02:53:47', NULL, '', '', '2026-01-17 02:53:47'),
(14, 2, '3bc435230debba8098566ee9283a87226f8c94748a7dc7d623946a0cce75f988', '2026-01-24 02:54:24', NULL, '', '', '2026-01-17 02:54:24'),
(15, 3, 'f6f388464142d6c765f22d3712efc1658b6920d19597efbab90f0761e6ed8cac', '2026-01-24 03:12:10', NULL, '', '', '2026-01-17 03:12:10'),
(16, 2, 'baf708972f5296530748fc4853d3c6a0f9d44cd0202ce746f1a8b7731c46a3d4', '2026-01-24 03:12:50', NULL, '', '', '2026-01-17 03:12:50'),
(17, 2, 'f92a70392c8c40d01f6725d23215080c2c534d5d9f867409aebe8bbf547501a8', '2026-01-24 03:12:59', NULL, '', '', '2026-01-17 03:12:59'),
(18, 2, 'a4d1d7fba50ee6bc37b012c81e020dc21bfbbebdaf6b802657f3e3b2ebe0e139', '2026-01-24 03:14:47', NULL, '', '', '2026-01-17 03:14:47'),
(23, 4, 'e684ab0d8a15fdda1bbbc3e506c21e22c7d2d6250996fdbf5a7ce7b647e7d8c4', '2026-01-25 07:08:53', NULL, '', '', '2026-01-18 07:08:53'),
(24, 4, 'db6d36505cab36b3b79e51de0a1694fff9c4b708a2ce8eeeead2bbec141bd049', '2026-01-25 07:09:01', NULL, '', '', '2026-01-18 07:09:01'),
(25, 4, 'f159e74df0470d3c8cb47993f9ce83b67ef87a64243156026c4ad5cd6eea001f', '2026-01-25 07:09:09', NULL, '', '', '2026-01-18 07:09:09'),
(26, 4, 'cc5101943e9c04fba7b82da715de208decf1a75337e6289e72516cb8dde48a12', '2026-01-25 07:09:19', NULL, '', '', '2026-01-18 07:09:19'),
(27, 5, '502b572a406ac15b4f8fbf6e5ea7ba9fc11935e00ca9fdb316e27845aaa03630', '2026-01-25 07:42:02', NULL, '', '', '2026-01-18 07:42:02'),
(28, 6, 'e838d7baf899b2c0a678547a39582ca47b9e9d3a507a8ae197f98df28bcfec26', '2026-01-25 07:43:28', NULL, '', '', '2026-01-18 07:43:28'),
(29, 6, 'e329f75dcbeae4b3aab073e27d1a229bb6716718082e53d5144a3d2455267a94', '2026-01-25 07:43:33', NULL, '', '', '2026-01-18 07:43:33'),
(30, 6, '99afa3c2e6726007e6d2062aa55f3a17e083c60905a910968e492a9b234c08fe', '2026-01-25 07:45:38', NULL, '', '', '2026-01-18 07:45:38'),
(31, 6, '7689a0e6b6bd8e10a0a427a7031f8806615ae7c9f30dc78be7242376255e055d', '2026-01-25 07:45:49', NULL, '', '', '2026-01-18 07:45:49'),
(32, 7, 'a9cf5c2eed090d2b482531c18a2bd20ec17094b080dd0418c59998d942612f3d', '2026-01-21 14:30:46', '2026-01-21 14:30:46', '', '', '2026-01-18 09:13:55'),
(33, 6, '2a001d6285eab9929dd40ed49228792412f4876a85357a5536d6d29dc0dcf23c', '2026-01-25 09:24:26', NULL, '', '', '2026-01-18 09:24:26'),
(34, 6, '73a916bb49ac98ca2209557a1f77ecf37d2e561f70a7c7c581f30e720b9abad8', '2026-01-25 09:24:43', NULL, '', '', '2026-01-18 09:24:43'),
(35, 7, '492b7ac0203728ed653ae262a5ed3bdf0136900d3bb810eff4313aa065367a7c', '2026-01-21 14:30:46', '2026-01-21 14:30:46', '', '', '2026-01-18 09:30:55'),
(36, 7, '3a4494825b610104798cabdb98f802e5ca3deea297d6b1283da67ba6dd908770', '2026-01-21 14:30:46', '2026-01-21 14:30:46', '', '', '2026-01-19 03:19:22'),
(38, 9, 'b84311ea9ecaf56a971324bc225b1282bc4ebcfd23cad0e6d4096d6330efb428', '2026-01-26 07:44:42', NULL, '', '', '2026-01-19 07:44:42'),
(39, 9, 'c422e7b8ce81e16955e97ccbee833d5363a46c8a7833cf90e143e59424c5a1ad', '2026-01-26 07:44:42', NULL, '', '', '2026-01-19 07:44:42'),
(40, 10, '6ba39df159ff2f5ff5ee8a882613018387b0b4c7811d0f7b1e54787cd613a7eb', '2026-01-26 07:48:03', NULL, '', '', '2026-01-19 07:48:03'),
(41, 7, '165952470c425c53003d70e42ebe4064657097c48b5508162c53d44d7427a539', '2026-01-21 14:30:46', '2026-01-21 14:30:46', '', '', '2026-01-19 08:45:03'),
(42, 7, '6d1615471fb01dfe916eac9018717fa29938e503fd6839733bfed8ef81d98b24', '2026-01-21 14:30:46', '2026-01-21 14:30:46', '', '', '2026-01-19 09:08:29'),
(43, 7, '802c7b7e85b969ac06cccb49f5b8fd49f395d59664b85405b9ac8b738c4254b8', '2026-01-21 14:30:46', '2026-01-21 14:30:46', '', '', '2026-01-19 09:17:25'),
(44, 7, '6a099673f30809c847a41958a8112fde7b6f111161775cdca8cc8fee3520add4', '2026-01-21 14:30:46', '2026-01-21 14:30:46', '', '', '2026-01-19 09:23:58'),
(45, 7, '51b304ec574e7e09973cbc26b3ff6e614a2d736e2f87c9f425c1468fa21687b4', '2026-01-21 14:30:46', '2026-01-21 14:30:46', '', '', '2026-01-19 09:26:52'),
(46, 7, '31b5c12979c17e775b5b1ac4e4de03438aaec394f67be38a858c97b1141635b2', '2026-01-21 14:30:46', '2026-01-21 14:30:46', '', '', '2026-01-19 09:28:07'),
(47, 7, 'a7549509411647c8a69f56e7defee2351a08feb15cc6f275925629db077b50c4', '2026-01-21 14:30:46', '2026-01-21 14:30:46', '', '', '2026-01-19 09:28:22'),
(48, 7, 'ca5466a3c7888c45fc1278366f852ba5facee273702e2c51e44db24dba6b9293', '2026-01-21 14:30:46', '2026-01-21 14:30:46', '', '', '2026-01-19 09:28:32'),
(49, 7, '9222dd7ef70b9703e9c7a61b99de032a6425d6e445a40d0d2634d93c8e163a45', '2026-01-21 14:30:46', '2026-01-21 14:30:46', '', '', '2026-01-19 09:29:51'),
(50, 7, '030620c5e4c6e29a9c91edf548e1a8da5bb9f59200d2cbf983b2eca4aef2692b', '2026-01-21 14:30:46', '2026-01-21 14:30:46', '', '', '2026-01-19 09:32:44'),
(51, 7, 'a6e1868c7414376fb9fc23b911ed8e8b5aa948820bc39ed4a3389a03e418f599', '2026-01-21 14:30:46', '2026-01-21 14:30:46', '', '', '2026-01-19 09:33:20'),
(52, 7, '8313fde1d1de58c6b6ec84effd5e547e93e6e65a8536e398b953caf4e67d17e7', '2026-01-21 14:30:46', '2026-01-21 14:30:46', '', '', '2026-01-19 09:50:11'),
(53, 7, '60b94d83767aeedf47e9238e9696c3fb36c18df8222ac0cf2118ca104243643d', '2026-01-21 14:30:46', '2026-01-21 14:30:46', '', '', '2026-01-19 09:52:54'),
(54, 7, '059f7d2ef6d19ba619a43d128429e30e4296652a974422beafe5f0b68cc0d7c0', '2026-01-21 14:30:46', '2026-01-21 14:30:46', '', '', '2026-01-19 11:21:19'),
(55, 7, '2a950819245b52d5645dbbe90d8b406bad4d48cab508e517ac790512b9df0f5f', '2026-01-21 14:30:46', '2026-01-21 14:30:46', '', '', '2026-01-19 11:53:25'),
(56, 7, '8725208f1c198c038722c83254c7c84936a3ad018f244a6740b067fdd0a8e3fa', '2026-01-21 14:30:46', '2026-01-21 14:30:46', '', '', '2026-01-19 12:44:17'),
(57, 7, '6ce5ce65b4ee76f77199c0cbb565f7bd4a05e2af01f127f7e6b18c259bd58edc', '2026-01-21 14:30:46', '2026-01-21 14:30:46', '', '', '2026-01-19 13:00:57'),
(58, 7, 'd6821da179e51d9755f3a5fd45e188dccb59058e838a1967682b4e4f5af9dd9e', '2026-01-21 14:30:46', '2026-01-21 14:30:46', '', '', '2026-01-19 13:07:54'),
(59, 7, '2b87fa6d775286438439d2b8e0c72ec2a15c7f61cef9a406c1bdceeefa22bad0', '2026-01-21 14:30:46', '2026-01-21 14:30:46', '', '', '2026-01-19 13:20:06'),
(60, 7, '277c6b4278d4dbc36231ed26ccacacb5212c19c79c6559231c828a09fd95d74c', '2026-01-21 14:30:46', '2026-01-21 14:30:46', '', '', '2026-01-19 13:49:57'),
(61, 7, '573deb02196412717a9d372689ef4854e26e5d3684ad4542360ff3c305a9d216', '2026-01-21 14:30:46', '2026-01-21 14:30:46', '', '', '2026-01-19 13:57:45'),
(62, 7, 'b712bb56482eb7aa73cbf8673eb651d039583d1ed2623e066b31354a85b209d1', '2026-01-21 14:30:46', '2026-01-21 14:30:46', '', '', '2026-01-19 14:01:13'),
(63, 7, '857b6a3aead65786de0f4e3a8011ecddb3dd1e6f20ee1cf2e81f8d3aa5c98b2a', '2026-01-21 14:30:46', '2026-01-21 14:30:46', '', '', '2026-01-19 14:18:08'),
(64, 7, '54eb77166f52390b05e13640b12451dd8f3c195342cba1eaf0687f93d5663ce8', '2026-01-21 14:30:46', '2026-01-21 14:30:46', '', '', '2026-01-19 14:23:15'),
(65, 7, 'b212a3a0942eb61f5e6ea37fa500b908174b38ad96cbc28fcea9ce3d7badf560', '2026-01-21 14:30:46', '2026-01-21 14:30:46', '', '', '2026-01-19 14:30:53'),
(66, 7, '81acc162faa4528132df14379cd4b3a78087f5603a31b307eefbfdffc365c1de', '2026-01-21 14:30:46', '2026-01-21 14:30:46', '', '', '2026-01-19 14:39:44'),
(67, 7, '0644999b481f7fc85f5d314ffc2c37758b236d1dbfaf460cdbc6d11615685bed', '2026-01-21 14:30:46', '2026-01-21 14:30:46', '', '', '2026-01-19 14:40:09'),
(68, 7, 'e6f46e45b9b2b6a0d62e1122c8ea393efe5c18c557065b850f5d5fd8143899c4', '2026-01-21 14:30:46', '2026-01-21 14:30:46', '', '', '2026-01-19 14:40:43'),
(69, 7, 'b9c1d8f68f2909a7898f3a35b2be0120d3e3c80a9c678f9a4d16ac4a8851742f', '2026-01-21 14:30:46', '2026-01-21 14:30:46', '', '', '2026-01-20 02:08:43'),
(70, 7, '606c5371d404836d3e036c80c137a27672039b3931182f1d08e4f446b8ad31d0', '2026-01-21 14:30:46', '2026-01-21 14:30:46', '', '', '2026-01-20 02:15:37'),
(71, 7, '77b5db128439adb5e41fd311798251a45c8cfa68cddaef77211e34ffa097990b', '2026-01-21 14:30:46', '2026-01-21 14:30:46', '', '', '2026-01-20 02:34:42'),
(72, 7, 'c66865a18e5308a3fb05a4d6dae343b14980abbb30209a83d75c3f314a136c28', '2026-01-21 14:30:46', '2026-01-21 14:30:46', '', '', '2026-01-20 02:42:06'),
(73, 7, '70d498c5070965c45b71c2b2d51154e678b42357ee169a629f4a780023bea7e6', '2026-01-21 14:30:46', '2026-01-21 14:30:46', '', '', '2026-01-20 02:58:00'),
(74, 7, '7692b2d788f3eeb82a8aad9726fd25fb30ca41236b8442f2e7f08d305b95bb4d', '2026-01-21 14:30:46', '2026-01-21 14:30:46', '', '', '2026-01-20 02:58:15'),
(75, 7, '448713cf8a721aab434daa2b07dcef61e9451f8bc538c6206fe3df6505b8a0f0', '2026-01-21 14:30:46', '2026-01-21 14:30:46', '', '', '2026-01-20 03:07:37'),
(76, 7, 'b078bae01a9c5f16236f0dda6301b4a41783c123cd767f49a42c04a1bffed631', '2026-01-21 14:30:46', '2026-01-21 14:30:46', '', '', '2026-01-20 03:08:18'),
(77, 7, '45416cb7e854a0c565dc0565a6e88ae4ae27ca2739643b0f7715b5bd0ed023be', '2026-01-21 14:30:46', '2026-01-21 14:30:46', '', '', '2026-01-20 03:09:40'),
(78, 7, '45875c5981930bee5d3273a54dca1d744a5016c76e4f4aa3f809990078965b14', '2026-01-21 14:30:46', '2026-01-21 14:30:46', '', '', '2026-01-20 03:09:49'),
(79, 7, 'f0fd8265affdf10b594851c7689ec62ce68fbf35134d3be58ee23f6fef47b279', '2026-01-21 14:30:46', '2026-01-21 14:30:46', '', '', '2026-01-20 03:10:16'),
(80, 7, 'df9c221d52f392e6267c40eef51a8ec5c3282df565b2334d74a5961b175e7ac9', '2026-01-21 14:30:46', '2026-01-21 14:30:46', '', '', '2026-01-20 03:10:29'),
(81, 7, '809de13f87c8bd92ae95d2c806ddb320c8e4158a4f7ae6bf261277f9a4284ffe', '2026-01-21 14:30:46', '2026-01-21 14:30:46', '', '', '2026-01-20 03:13:00'),
(82, 7, 'dede9156ce8b831a98528aa5b4554d894f3160caded91ea5e8f5c28081ee00d6', '2026-01-21 14:30:46', '2026-01-21 14:30:46', '', '', '2026-01-20 03:13:08'),
(83, 7, 'ef95e3e43bfb9cf8c0d64cf9e232b7d13495e7efee2f70d93292e348bfb3103b', '2026-01-21 14:30:46', '2026-01-21 14:30:46', '', '', '2026-01-20 03:13:52'),
(84, 7, '509b8d1a4b6ae911fd78899220bf183a78b5893739ffda1c2a73971d78fabbd9', '2026-01-21 14:30:46', '2026-01-21 14:30:46', '', '', '2026-01-20 03:14:23'),
(85, 7, '1f280b56a341a2e87d2184feb7839c5a248ceb2aa4c268ca154b31606d626b40', '2026-01-21 14:30:46', '2026-01-21 14:30:46', '', '', '2026-01-20 03:24:30'),
(86, 7, 'e49a2f7582e40c4faee703e1949dd9a7749500f685ba142c7a50a051ab345ead', '2026-01-21 14:30:46', '2026-01-21 14:30:46', '', '', '2026-01-20 03:35:55'),
(87, 7, '35b89f2b1540359403d0c11c6594a3829dab3ad6427536ec65af58e28e8f8d8d', '2026-01-21 14:30:46', '2026-01-21 14:30:46', '', '', '2026-01-20 03:38:47'),
(88, 1, '757f8097ca62f94e1fdede4c0daa1964c9409d075f690d8b8f2101698a1f6113', '2026-01-27 03:39:10', NULL, '', '', '2026-01-20 03:39:10'),
(89, 1, '5f3d2a21161326c932efeda1b5791ea183ef9cc11b7fa955c7f5c6e073501ab9', '2026-01-27 03:39:29', NULL, '', '', '2026-01-20 03:39:29'),
(90, 1, '4588c3e22ca04ff044c40dda7c5b962c7404575cf9c4e7228a6de603d1fc901e', '2026-01-27 03:42:26', NULL, '', '', '2026-01-20 03:42:26'),
(91, 1, '82bc1046e59836de187fc60cc041b3d91c34bfaf947a2ddc93c8b66a83e7c76d', '2026-01-27 03:43:00', NULL, '', '', '2026-01-20 03:43:00'),
(92, 7, '912696beb110543df2a2576f33042c7bec87f04203442d0ba2acec9e3056798f', '2026-01-21 14:30:46', '2026-01-21 14:30:46', '', '', '2026-01-20 03:43:14'),
(93, 7, '52a7e5200ea85a5f23889d5f274b16653e3dc3d445dbc7a851130dbf54a61dd0', '2026-01-21 14:30:46', '2026-01-21 14:30:46', '', '', '2026-01-20 03:51:22'),
(94, 12, '671b5b9f4eab05ddabd87bdbe0c996b84e45a237194325674f6c87e3e084ddad', '2026-01-27 04:39:36', NULL, '', '', '2026-01-20 04:39:36'),
(95, 7, '41b91ae76fe38c780adfa4fec34f3fcfb682bc5006f88449353ae02e49a89340', '2026-01-21 14:30:46', '2026-01-21 14:30:46', '', '', '2026-01-20 04:45:56'),
(96, 7, '078394421ada52af0f3ce1fc95284550834e730d136d50436cc7b9b1b83810aa', '2026-01-21 14:30:46', '2026-01-21 14:30:46', '', '', '2026-01-20 05:03:43'),
(97, 7, 'a8ea2ad138f34e49ec75e6374d8ae61ef9af45e3eb9c6ce13b1d0bb81a91ae79', '2026-01-21 14:30:46', '2026-01-21 14:30:46', '', '', '2026-01-20 05:03:52'),
(98, 7, 'f4876476ab5ed43b3d48afc3df420aa9edc00932271611fd17c78af43e89ed9f', '2026-01-21 14:30:46', '2026-01-21 14:30:46', '', '', '2026-01-20 05:09:04'),
(99, 7, '38b15225b835bccb57e632f23898c7ce1a642ccf26933b08613dac149b3409d4', '2026-01-21 14:30:46', '2026-01-21 14:30:46', '', '', '2026-01-20 05:11:57'),
(100, 7, 'c00c1dbd409e48265a6b1fcc816f5b60069ce028641fe2fc7ca1f681556eed4b', '2026-01-21 14:30:46', '2026-01-21 14:30:46', '', '', '2026-01-20 05:18:09'),
(101, 7, '44e80ca96e11eca5e24c3f566718a2f4645cc6d7420b3d68db58f7c74ae19883', '2026-01-21 14:30:46', '2026-01-21 14:30:46', '', '', '2026-01-20 05:18:15'),
(102, 7, '3e360f4c576b9d4e422ea73e34be89b628da2f3aebe609e4e9bdaaa15e043c49', '2026-01-21 14:30:46', '2026-01-21 14:30:46', '', '', '2026-01-20 05:27:36'),
(103, 7, 'a483bfa3ffb8c4953ce18df9020b74647859df40a75b02dbebc7164a39c3f9cc', '2026-01-21 14:30:46', '2026-01-21 14:30:46', '', '', '2026-01-20 05:34:03'),
(104, 7, '05299d6089289c5829f03cd3a16612d8f0d07f148642e564a95ffcb79194ba2f', '2026-01-21 14:30:46', '2026-01-21 14:30:46', '', '', '2026-01-20 05:54:36'),
(105, 7, 'e22b19117cf558ecf5027ccfdb2a38ac5d8350b299d1d127296a8b6a06a819b9', '2026-01-21 14:30:46', '2026-01-21 14:30:46', '', '', '2026-01-20 06:09:40'),
(106, 7, 'f78e739c2a6b94e91e075398b0a2e9709a9fe1d6142bfb8c8a6ba7e4b4fc3fd9', '2026-01-21 14:30:46', '2026-01-21 14:30:46', '', '', '2026-01-20 06:26:53'),
(107, 7, '3c45cb8142020c1660fe4f35de367bafe6fff73c512cf25c71ea6373ecc0cee6', '2026-01-21 14:30:46', '2026-01-21 14:30:46', '', '', '2026-01-20 06:27:02'),
(108, 7, '567134ad71cfab1b6396537faad83802f86ebaea5f5c10d274f03fc36c410500', '2026-01-21 14:30:46', '2026-01-21 14:30:46', '', '', '2026-01-20 06:27:09'),
(109, 7, '52c742fe1d8d1d136a768fa947c7f90cb9c045e97a03ba4ef961fc8b0d585490', '2026-01-21 14:30:46', '2026-01-21 14:30:46', '', '', '2026-01-20 06:27:18'),
(110, 7, '4ebb80ffa28db3ad3aca43c203ba5ddd15a297241f2eb2683f3dd2a6b978af87', '2026-01-21 14:30:46', '2026-01-21 14:30:46', '', '', '2026-01-20 06:31:58'),
(111, 7, '2dac3522608bc84455a48efe8e7d794d58cf54c07a373e01d979b65365830673', '2026-01-21 14:30:46', '2026-01-21 14:30:46', '', '', '2026-01-20 06:34:17'),
(112, 7, '04792a89c031d49cf3090a1f1370631150523fa1137a7d16ac2a1534e83d6f3d', '2026-01-21 14:30:46', '2026-01-21 14:30:46', '', '', '2026-01-20 06:35:17'),
(113, 7, 'ebf7cfc30224915ae0e325470f1a0633056d9f4db1229224475d90213333a994', '2026-01-21 14:30:46', '2026-01-21 14:30:46', '', '', '2026-01-20 06:37:57'),
(114, 7, '7092720f84b43092afe6eb57f1615b00e5e705863d0a2089d690d929e97b4a67', '2026-01-21 14:30:46', '2026-01-21 14:30:46', '', '', '2026-01-20 06:43:09'),
(115, 7, '39f955f9ac80739ca001edfbfd5a961af0102e60c9db8e7d72b9fa0eb1618095', '2026-01-21 14:30:46', '2026-01-21 14:30:46', '', '', '2026-01-20 06:43:16'),
(116, 7, '7830a709c10076cf406c4ffa9f3b92853f08883bc80109167e42ebc4d407689c', '2026-01-21 14:30:46', '2026-01-21 14:30:46', '', '', '2026-01-20 06:44:14'),
(117, 7, '92cd1a132e95f06f8b505997a4cc4da83aee2c9647f74984bead31d0f9b52a81', '2026-01-21 14:30:46', '2026-01-21 14:30:46', '', '', '2026-01-20 07:01:11'),
(118, 7, '5e830f1ec835c53822a30451f352791af5baed81dd91637735befeade7ba212b', '2026-01-21 14:30:46', '2026-01-21 14:30:46', '', '', '2026-01-20 07:14:54'),
(119, 1, 'f793b289ec992152ee29fe6a67ee09755930d8f4ff140e2259b5e21c8d6da028', '2026-01-27 07:23:29', NULL, '', '', '2026-01-20 07:23:29'),
(120, 7, '8edbedc91fd5afe59d2bf6a1713d84b7ed87d28fb1dadffc77564a11a8741e2e', '2026-01-21 14:30:46', '2026-01-21 14:30:46', '', '', '2026-01-20 07:23:49'),
(121, 7, 'ddae9b2f1490b3139245b1ee15ae0b6573c28a3645caab214a823933574e92c7', '2026-01-21 14:30:46', '2026-01-21 14:30:46', '', '', '2026-01-20 07:31:40'),
(122, 1, '9278c0609b57d2bdb7b54f29dc36743c2acac55a7dc08745bb01fb7ac6128069', '2026-01-27 07:31:40', NULL, '', '', '2026-01-20 07:31:40'),
(123, 7, '6297cfd594eea06f6ccad1bc355742e217f445906b288b70c8ef1c54fb6357b5', '2026-01-21 14:30:46', '2026-01-21 14:30:46', '', '', '2026-01-20 07:31:47'),
(124, 1, '33f173c2e4512119fa09a9364bac2f041d992111095c89688e98790d515c7a04', '2026-01-27 07:31:55', NULL, '', '', '2026-01-20 07:31:55'),
(125, 7, '88d204fda219c6f34d9a99ce421e68ad9f8f66c1a9cfbf5a3015e1279842f310', '2026-01-21 14:30:46', '2026-01-21 14:30:46', '', '', '2026-01-20 07:38:23'),
(126, 7, '1e1fa7e7e3855209fa2162ff1e7765e8f74aa410362ab29d72ed4f361fcf2fef', '2026-01-21 14:30:46', '2026-01-21 14:30:46', '', '', '2026-01-20 07:38:33'),
(127, 1, '80c99e42686621fa56568b13ce1996d267ec936bec3c910c479be6e79da142ec', '2026-01-27 07:39:26', NULL, '', '', '2026-01-20 07:39:26'),
(128, 7, '61fd0cd20ef183f1acec515956896bb394b1284b17fc5d5f43109ac9007483cd', '2026-01-21 14:30:46', '2026-01-21 14:30:46', '', '', '2026-01-20 07:41:36'),
(129, 7, '059b253aa2a70bd13d6e23759a4cc5cbb5000b009201b4b75da94ff041997c54', '2026-01-21 14:30:46', '2026-01-21 14:30:46', '', '', '2026-01-20 07:46:42'),
(130, 7, 'be538b0bc7b7cf3e00d2c60ba19e0bb7e58814c7a6c880a6ab42502643df0d7e', '2026-01-21 14:30:46', '2026-01-21 14:30:46', '', '', '2026-01-20 08:28:06'),
(131, 7, '796274652bb3cda96459b59a9c60541a88dbab2e03eb5571f77ee1fa6f64f821', '2026-01-21 14:30:46', '2026-01-21 14:30:46', '', '', '2026-01-20 08:40:45'),
(132, 7, '039a8e65e6d016a8fbc3534c6a7bafc1d4c09436b990460276900541cfb90e9f', '2026-01-21 14:30:46', '2026-01-21 14:30:46', '', '', '2026-01-20 09:18:08'),
(133, 7, 'd222e380bd264976f11633225b1d00fcd38fabff89da7e2e44ea3edeef4e7c50', '2026-01-21 14:30:46', '2026-01-21 14:30:46', '', '', '2026-01-20 09:34:29'),
(134, 7, 'a3d5cd3df66637a0dd3df62585f07cbb4cacb96088fdc57cd6376971ad00bd7f', '2026-01-21 14:30:46', '2026-01-21 14:30:46', '', '', '2026-01-20 09:35:21'),
(135, 7, 'f84aae432b1478f5b84b66768cdb25321c9003b945f41b5cf1b50624cdadb926', '2026-01-21 14:30:46', '2026-01-21 14:30:46', '', '', '2026-01-20 09:35:28'),
(136, 7, '30e8fbb61c3621857ecc6942b3e5e8afa5045a72960ef5d7568b0ed3cf6510d2', '2026-01-21 14:30:46', '2026-01-21 14:30:46', '', '', '2026-01-20 09:36:05'),
(137, 7, 'e74f351f6156e3244510490d9fe47a181e7997896d063ce782cb6ff2399ae243', '2026-01-21 14:30:46', '2026-01-21 14:30:46', '', '', '2026-01-20 09:36:40'),
(138, 7, '7d86f63076862b414160744bc0921936783408e7a06e158f3a163dfafd241c40', '2026-01-21 14:30:46', '2026-01-21 14:30:46', '', '', '2026-01-20 09:36:51'),
(139, 7, '2acfc24a8488e55b9311120ca1b42bbe49e7e9386ef20204e6730ee1ad5841c0', '2026-01-21 14:30:46', '2026-01-21 14:30:46', '', '', '2026-01-20 09:37:13'),
(140, 7, 'e5ab349748aba45866f7a437909d8e0a9e2b651cfd259231543ff2c977da2781', '2026-01-21 14:30:46', '2026-01-21 14:30:46', '', '', '2026-01-20 09:38:16'),
(141, 7, '31a41109685d63b50e680968df5ade3e234b033c8132d215807e248094bdf5ca', '2026-01-21 14:30:46', '2026-01-21 14:30:46', '', '', '2026-01-20 09:42:28'),
(142, 7, '2c4397b1410b2882bc8129e390cb187e7bf91960e0af3a5b40875fa80cbd86ba', '2026-01-21 14:30:46', '2026-01-21 14:30:46', '', '', '2026-01-20 09:42:37'),
(143, 7, '8e0aeabafd4726d4c9b66eabd0413d74b40764a2149d4970f67d529b10a63901', '2026-01-21 14:30:46', '2026-01-21 14:30:46', '', '', '2026-01-20 09:42:55'),
(144, 7, 'cb8e3c725aa9aa4b870c2d410ebfbbc3caaa5b83b3c8787a85be92d806b87f04', '2026-01-21 14:30:46', '2026-01-21 14:30:46', '', '', '2026-01-20 09:43:01'),
(145, 7, 'b16e1b6ed90187eb1ee7e4cd54489c8e6b7358cfd0e78b3e16f773e0fbf670f2', '2026-01-21 14:30:46', '2026-01-21 14:30:46', '', '', '2026-01-20 09:44:01'),
(146, 7, '7abf61abb121fdaa0f70527049a7a4a231c516f2b4ab624a96ccd7ba22b8567b', '2026-01-21 14:30:46', '2026-01-21 14:30:46', '', '', '2026-01-20 09:55:26'),
(147, 7, '5620d0c506b4c93bb3b2a6ca5c3cdd3add13bf01c1a79c149a12edf1408b6cb2', '2026-01-21 14:30:46', '2026-01-21 14:30:46', '', '', '2026-01-20 09:59:02'),
(148, 7, '641f872c5eddd729f2e36da02ca976d667d377083686ef77f7aca9ecfb72dde0', '2026-01-21 14:30:46', '2026-01-21 14:30:46', '', '', '2026-01-20 10:53:59'),
(149, 7, 'be94c688626ce3680fbb0982395baca150614f3ccf032cf16b79ad0151969b85', '2026-01-21 14:30:46', '2026-01-21 14:30:46', '', '', '2026-01-20 11:01:19'),
(150, 7, 'b054ad7fe7eceae49b9e0e12b345882f5af4c591f948018c0b666311ded23bb3', '2026-01-21 14:30:46', '2026-01-21 14:30:46', '', '', '2026-01-20 11:43:20'),
(151, 7, '4f34fb29b19fd3c88599355ca780d947c168776ba238d7977878b8dae35df570', '2026-01-21 14:30:46', '2026-01-21 14:30:46', '', '', '2026-01-20 12:44:46'),
(152, 7, '45c3cbbb6e45931a69b0ff0a25b322509a26f9d2e9df34cd7c54e9c00bf9d1ec', '2026-01-21 14:30:46', '2026-01-21 14:30:46', '', '', '2026-01-20 13:09:51'),
(153, 7, '32990e21191f130301dcb381bec1e6429fc750dac989d2fd92e011bad3ad12ce', '2026-01-21 14:30:46', '2026-01-21 14:30:46', '', '', '2026-01-20 13:53:02'),
(154, 7, '8756771edb20eb14e2945d57d618e3b1e260c040fe1a14dfa0c1b5385c15246a', '2026-01-21 14:30:46', '2026-01-21 14:30:46', '', '', '2026-01-20 13:53:10'),
(155, 7, '8361e121ea060ae86e6e549a39a0a7b823ee88ea279b360ca6c35e44c658b0cd', '2026-01-21 14:30:46', '2026-01-21 14:30:46', '', '', '2026-01-20 13:53:46'),
(156, 7, '4dee8624b2dc430b2a6c4cb462d36dd28d7c0199c79b4078bfc1ae6680a14110', '2026-01-21 14:30:46', '2026-01-21 14:30:46', '', '', '2026-01-20 13:56:58'),
(157, 7, 'bec0178b7779f1dec8e1e13da4a3d0cc54954f2da6ce4067ea2474cabb0f8168', '2026-01-21 14:30:46', '2026-01-21 14:30:46', '', '', '2026-01-20 13:59:07'),
(158, 7, '85b103f97cdec44e9daa77e184f4dba8fd8294fe825bad97f6760bfe64d63d82', '2026-01-21 14:30:46', '2026-01-21 14:30:46', '', '', '2026-01-20 13:59:20'),
(159, 7, 'c51abd1d24ad9444526dffb2aba7097d317d3e0ac19c2a4acdfd4fa5cc7e5f65', '2026-01-21 14:30:46', '2026-01-21 14:30:46', '', '', '2026-01-20 14:00:49'),
(160, 7, '775fbc3468567e99ae6a8835ef7c44049e3b479bcd2f0b80cfe3a0edc94ed851', '2026-01-21 14:30:46', '2026-01-21 14:30:46', '', '', '2026-01-20 14:01:55'),
(161, 7, '82c79a13ebc30441d52cc97c82333d2aebd0aff9b42b7a21a935558939eb7b4b', '2026-01-21 14:30:46', '2026-01-21 14:30:46', '', '', '2026-01-20 14:04:31'),
(162, 7, '6773b4a3ec439d2817f125e328772c557f113a347027b861516b1a7e7c5b7266', '2026-01-21 14:30:46', '2026-01-21 14:30:46', '', '', '2026-01-20 14:05:15'),
(163, 7, 'd30fabbe7f287662ea51925ac8c7c2c9727502f90ae2a51a3114c084c2c0409d', '2026-01-21 14:30:46', '2026-01-21 14:30:46', '', '', '2026-01-20 14:06:55'),
(164, 7, 'd1575903b246d68db2c8efc99dea52ab68c4446c93d04105d3d5e6174e17e754', '2026-01-21 14:30:46', '2026-01-21 14:30:46', '', '', '2026-01-20 14:15:32'),
(165, 1, '866236816b8d7f07d3d309556af00032760f5ff0db3e4b78d0b6ab616666024e', '2026-01-27 14:17:05', NULL, '', '', '2026-01-20 14:17:05'),
(166, 7, 'd6de20d1ae171fa94b0b518a30f0912910ddca81969d5fc6aef0b144e1127d66', '2026-01-21 14:30:46', '2026-01-21 14:30:46', '', '', '2026-01-20 14:24:59'),
(167, 7, 'cae3e131609a45b456ec2984128049b98c2755fb0b894aac6c14f316e181eed5', '2026-01-21 14:30:46', '2026-01-21 14:30:46', '', '', '2026-01-21 13:50:19'),
(168, 7, 'e86eba5fb38607eb191e0991cd73344805d9747266242094c137639823dae39d', '2026-01-21 14:33:12', '2026-01-21 14:33:12', '', '', '2026-01-21 14:31:44'),
(169, 7, 'b8aaa042cc65529673ae63463bc6f1bcfa1517ca263c8cb9b38ac16710dbd610', '2026-01-21 15:45:14', '2026-01-21 15:45:14', '', '', '2026-01-21 14:33:24'),
(170, 13, 'c9d2fc0f3e445f98ca7ed67ef23e8eacc7cdcc6100449588f0a3eb80e612f2f4', '2026-01-28 15:06:03', NULL, '', '', '2026-01-21 15:06:03'),
(171, 7, '4d3dbd5eba7dcd38967d0efc8825999eb0f99a3961cb077f2748440d6391bc53', '2026-01-21 15:45:14', '2026-01-21 15:45:14', '', '', '2026-01-21 15:25:43'),
(172, 14, '907af0cfcd5a676f61bd9a8b934805753f52d7e2af7a1fc71bcfde79f65d84a4', '2026-01-21 15:31:30', '2026-01-21 15:31:30', '', '', '2026-01-21 15:31:15'),
(173, 14, '886666d8cdbfddc6ce6f774bdd2c5f4d8447866815aecee186b200353350b93c', '2026-01-28 15:31:50', NULL, '', '', '2026-01-21 15:31:50'),
(174, 15, 'c2a3e52f883769b94489f68fe6f4e2126e093d62f862bb0c4d20ac1ee54f05be', '2026-01-21 15:34:13', '2026-01-21 15:34:13', '', '', '2026-01-21 15:34:13'),
(175, 15, '181693203eecafd8271fb572b54f0b4d8a3513ecbfc617ec0b31ec80a770a872', '2026-01-21 15:34:13', '2026-01-21 15:34:13', '', '', '2026-01-21 15:34:13'),
(176, 15, 'adebb6a3ff24285ec15b9d72b23af6e3df68efaaeddf904768df13f67aae062b', '2026-01-28 15:34:13', NULL, '', '', '2026-01-21 15:34:13'),
(177, 7, 'b0355e3e9be2fd07b732cf0545238fd2db21adcb0a45378de4b70d57646d189a', '2026-01-21 15:45:14', '2026-01-21 15:45:14', '', '', '2026-01-21 15:44:40'),
(178, 7, '3408a60a3f7ce0ae130dcda0d9b17c24692d0974bf1d858ac6e32271d0e5637c', '2026-01-28 15:45:27', NULL, '', '', '2026-01-21 15:45:27'),
(179, 7, '7427e94d8205a9dbf5f460fbace00f279e8ca378cc5c6fe2056c2aa8d61f2648', '2026-01-28 15:48:26', NULL, '', '', '2026-01-21 15:48:26'),
(180, 7, '72e2a7ea88471a8ba9f0058ddf6a6ad594924dcf51bda4dc2bd6c1da214654ea', '2026-01-28 16:25:59', NULL, '', '', '2026-01-21 16:25:59'),
(181, 7, '01181a4335d6442f30a44cbe2cfe945019f05258d1af796e1f4edba2d767c7be', '2026-01-28 16:43:27', NULL, '', '', '2026-01-21 16:43:27'),
(182, 7, '9d5ad9ea763d4dde95bd39eac0093249fd87a0f0188e6a10d70a3f41c12fe926', '2026-01-28 16:59:29', NULL, '', '', '2026-01-21 16:59:29'),
(183, 7, '4d4aa5cf728096d896034d7daaf3bc268de7bded465f4774625c8930487361d7', '2026-01-28 17:00:59', NULL, '', '', '2026-01-21 17:00:59'),
(184, 7, 'a7f6a2d0d947550e975ac39c7c983dc5fef17240637fff48419ec451b7237ac5', '2026-01-28 17:06:28', NULL, '', '', '2026-01-21 17:06:28'),
(185, 7, 'b63821332b0379ffb7f7210d32aa189e1160c19b14a5495623fbbf69abadb66d', '2026-01-28 17:21:32', NULL, '', '', '2026-01-21 17:21:32'),
(186, 7, '07e1fb865f9cc86d965b4eb90a812a005b23738032fb73121e6cc2cf4277871a', '2026-01-28 17:26:43', NULL, '', '', '2026-01-21 17:26:43'),
(187, 7, 'f44c607cf5946770534702e5589da85733022238717f69dc84df8c54a2c2550e', '2026-01-28 17:42:01', NULL, '', '', '2026-01-21 17:42:01'),
(188, 7, 'e210546052cb8a77355ce589d4ba34caa8b049c6cbd93b4d8c22629c1e348a76', '2026-01-28 17:42:07', NULL, '', '', '2026-01-21 17:42:07'),
(189, 7, '9ae015e407d35bd95f6cdf5cb6cf9c564ef65daa5fc3441ccbd7b0d37ac9942a', '2026-01-28 17:43:34', NULL, '', '', '2026-01-21 17:43:34'),
(190, 7, '41a6335ef8d35dd1d01fd1a2b5b595bb3bdb148df746865996a1fa1197942323', '2026-01-28 17:43:41', NULL, '', '', '2026-01-21 17:43:41'),
(191, 7, '468f316b8de9e790cf1fd839f89cc313826af83a7bafc3adb367ac3079739ef5', '2026-01-28 17:48:41', NULL, '', '', '2026-01-21 17:48:41'),
(192, 7, '644556ecbea846cb87ebebdfd95230dcf5c5db2b503e7000d48d3811491d8dab', '2026-01-28 17:48:50', NULL, '', '', '2026-01-21 17:48:50'),
(193, 1, 'dd2f39eedbf9b68865e6b14688671a3063e9cf7b5ec486a8f4fe56508f12b1a2', '2026-01-28 17:49:00', NULL, '', '', '2026-01-21 17:49:00'),
(194, 1, 'c340a51ade96fa04d3579290135cbd2d7ec811b8c7e1b3551236868441030fc1', '2026-01-28 17:49:18', NULL, '', '', '2026-01-21 17:49:18'),
(195, 1, 'a4de93a856ea2213c4034b815baf40e14692157fdf9602278990fabed8ce96c7', '2026-01-28 17:49:26', NULL, '', '', '2026-01-21 17:49:26'),
(196, 7, '5dfb0f476366fee6a245526fc89cdab119e4a2df1478fcf960596564614fe150', '2026-01-28 17:53:05', NULL, '', '', '2026-01-21 17:53:05'),
(197, 7, 'c4a421a85cd63e73256121a0f8a425846a6da0c8758dfcaaee41d4e9d0893fb4', '2026-01-29 01:34:23', NULL, '', '', '2026-01-22 01:34:23'),
(198, 7, '2651c19610b13f8e6b379da9d79cc9c509a03141f327eff79530352ca0895a4d', '2026-01-29 01:38:43', NULL, '', '', '2026-01-22 01:38:43'),
(199, 7, '7271a2a6e9c89365b05143c22ea79cbfbbe5fdbdb46c00686a5d1469569961ae', '2026-01-29 01:39:02', NULL, '', '', '2026-01-22 01:39:02'),
(200, 7, 'b935b60717abc9c35a208e8b3e5d4573e20b856439f6f40af6789d6059936ea8', '2026-01-29 01:39:48', NULL, '', '', '2026-01-22 01:39:48'),
(201, 1, '75aa6395138bcad48ac94afa5f7f59088f9d4c207b6b8a68a60aeb0a7b0b9f2b', '2026-01-29 01:39:49', NULL, '', '', '2026-01-22 01:39:49'),
(202, 7, '5c8883baa4f8013274454d4447020ec2ee1d9315a18a618d41e142caaecd37a4', '2026-01-29 01:53:27', NULL, '', '', '2026-01-22 01:53:27'),
(203, 7, '102c1da2c1d87ed26c58b4023a119de8be117d0104c845f55b9d61cbc5f4f871', '2026-01-29 02:22:03', NULL, '', '', '2026-01-22 02:22:03'),
(204, 7, 'd2b79945bbe5910c3d1fb31100cc7b6fc516a94312cf5c2506bfce67e5a911b8', '2026-01-29 02:41:56', NULL, '', '', '2026-01-22 02:41:56'),
(205, 7, '7ceb88c4f7674949c1c0535fd4bb72c08b5937a43f15f51bf5bbee6bf3d04322', '2026-01-29 02:51:07', NULL, '', '', '2026-01-22 02:51:07'),
(206, 7, '247833d3b3d8560d31780b5c45f7e5d162ba19d0c77e1958085d63073696d967', '2026-01-29 02:51:24', NULL, '', '', '2026-01-22 02:51:24'),
(207, 7, '6289f07e8077fc534549b244bb2989b6b72b13e932e56ed42f451a8ba986b14b', '2026-01-29 02:51:39', NULL, '', '', '2026-01-22 02:51:39'),
(208, 7, '1618899e2e7b472f2b878ab3d5272def795ec4a129cc68274a33139464b02bd7', '2026-01-29 02:57:06', NULL, '', '', '2026-01-22 02:57:06'),
(209, 7, 'ee142b36e58ca27bea31bc07b382064cd2a6dd4cc734fe50bd222bbe4ec037a6', '2026-01-29 02:57:15', NULL, '', '', '2026-01-22 02:57:15'),
(210, 7, '6e79072e5d36b114660f1d9de136d246394e627ff627dc2c3186636752a52820', '2026-01-29 03:03:32', NULL, '', '', '2026-01-22 03:03:32'),
(211, 7, '6a30a6abb4b0922c0d328ec1f59f5e2f33bf23e6bc71acac5d44e7ae8f40f571', '2026-01-29 03:03:41', NULL, '', '', '2026-01-22 03:03:41'),
(212, 7, 'e93f81c3e53bf78ab49466157d406f7280ff4e87ed4dc4d0af7e7d62dd6e2dc5', '2026-01-29 03:03:52', NULL, '', '', '2026-01-22 03:03:52'),
(213, 7, '1399e07a14073cefe9f49cccde1722e352df5c981964b2be68d79feb3c24678b', '2026-01-29 03:04:26', NULL, '', '', '2026-01-22 03:04:26'),
(214, 7, '0fc0a0d2c8d886a0e793e492a2e9d91f9a17cf4c2f0a4544c641a5284264cea0', '2026-01-29 03:04:56', NULL, '', '', '2026-01-22 03:04:56'),
(215, 7, 'ab1c1d9d2068522e18b058f0d223011da553b920b89d510e6b319b49b637ddc1', '2026-01-29 03:13:03', NULL, '', '', '2026-01-22 03:13:03'),
(216, 7, '02bb5cdfec46fa208b1798b9fb684e95176fa82e7124dc0b544965972ea9f04a', '2026-01-29 03:23:24', NULL, '', '', '2026-01-22 03:23:24'),
(217, 7, '6fa36beb607cc78ebe5c683ceac07875518802ae3a979aec3674148053577522', '2026-01-29 03:52:55', NULL, '', '', '2026-01-22 03:52:55'),
(218, 7, '12ac41638b5dd40b70e51fec7278857eba8fa3ba45df6fecbb07af3a10791441', '2026-01-29 03:53:40', NULL, '', '', '2026-01-22 03:53:40'),
(219, 7, '893ec89224aad0149e4d16d472658fa6acfc23a196840c479be58289c819173b', '2026-01-29 03:58:54', NULL, '', '', '2026-01-22 03:58:54'),
(220, 7, '2fbf32a24a6e9eff6c99447d3e22061186f94fdd4c8115e709e12a8fbe88ba7f', '2026-01-29 04:00:38', NULL, '', '', '2026-01-22 04:00:38'),
(221, 7, '44aad3f97669d592b439c0073346350f553b1955d28cb399ec20a09a2b113f68', '2026-01-29 04:00:53', NULL, '', '', '2026-01-22 04:00:53'),
(222, 7, '3cb3f4712d8912034ae40320809bfb361e595adda85117a4081740f0e9df5867', '2026-01-29 04:02:07', NULL, '', '', '2026-01-22 04:02:07'),
(223, 7, '19014b0e639d7bbdce980ad2f9d8dc0405e6dec2b9698aef4f6fb419a9c18041', '2026-01-29 04:08:44', NULL, '', '', '2026-01-22 04:08:44'),
(224, 7, '41948ca0c577557b09f6ba756aba13d3d6be3506da331c3bd45b84b414032366', '2026-01-29 04:15:35', NULL, '', '', '2026-01-22 04:15:35'),
(225, 7, 'b9d82e8d33cda0c3c712c272c679ba15a58cd4d39740eacbe1fd73606ecb61db', '2026-01-29 05:50:08', NULL, '', '', '2026-01-22 05:50:08'),
(226, 7, '9d1a0b07a9ea44ca9dd0c17a50aca4ad0fc72d44ad2a3c4c84e40b5ea44821a5', '2026-01-29 06:02:37', NULL, '', '', '2026-01-22 06:02:37'),
(227, 7, 'f3dc71bb619d23bcac83db9602a21fea5041533d7190d53bfd49f761b4f923ea', '2026-01-29 06:17:15', NULL, '', '', '2026-01-22 06:17:15'),
(228, 7, '5fe0f41c4a0f2826d79626bd21b3cb9c4de2782aa404c9349789904cdd50c044', '2026-01-29 06:17:20', NULL, '', '', '2026-01-22 06:17:20'),
(229, 7, 'a21e7b461c6ab444b1f67cc4d6dd5023c388b4e9876616c8ea2802e821dfa100', '2026-01-29 06:17:51', NULL, '', '', '2026-01-22 06:17:51'),
(230, 7, 'de905dde15c990dad3e1d969968a5c194580d19be65afb03dc05c7e60ff6f317', '2026-01-29 06:18:07', NULL, '', '', '2026-01-22 06:18:07'),
(231, 7, '2f76b6c1549955c7e69c56880e6b2357a4aee9c4608084d328927cdeef5a0222', '2026-01-29 06:40:39', NULL, '', '', '2026-01-22 06:40:39'),
(232, 7, '84fa9c952a2e23378803a72d4763af7b328da18c07e4f1deaf7befddec2d5418', '2026-01-29 06:41:03', NULL, '', '', '2026-01-22 06:41:03'),
(233, 7, '28ad48d5664740257e8d2d31d6285d69ad1115066f5cb2d2039bc142a753c08d', '2026-01-29 06:42:46', NULL, '', '', '2026-01-22 06:42:46'),
(234, 7, '07ea1b56ab19a257faf6cb117ff09a90eb7cba7b598a6e4701626d5c90bad29f', '2026-01-29 06:44:23', NULL, '', '', '2026-01-22 06:44:23'),
(235, 7, 'a97f40a1720c1ad30ce973e3e5b4124544002835a200cb5032c172e57c10fd69', '2026-01-29 06:45:08', NULL, '', '', '2026-01-22 06:45:08'),
(236, 7, '188fe2372d2615ff9f87585b5f3f1b04f7ed52b0db1181f1b66d76dc319ebeaf', '2026-01-29 09:21:32', NULL, '', '', '2026-01-22 09:21:32'),
(237, 7, '2acdf223f70ab0a5195b5a00d47ca00035f5cf7ea6b022a57f0ff0da3bcebd89', '2026-01-30 01:39:17', NULL, '', '', '2026-01-23 01:39:17'),
(238, 7, 'b1f272d8819345afe9a862d5bf7864b23b5ce41454b6d62bfad3c612dadfb51f', '2026-01-30 01:59:53', NULL, '', '', '2026-01-23 01:59:53'),
(239, 7, '8bacf4ebe9e23702c5ef3d81a42e6db7e011e4a98da3b21901d620cdcd534386', '2026-01-30 02:18:03', NULL, '', '', '2026-01-23 02:18:03'),
(240, 16, '33d9c933b3aa3d9c65b156d3233adea4324eb66580bfeb59b7230e34bfae2297', '2026-01-30 03:16:19', NULL, '', '', '2026-01-23 03:16:19'),
(241, 17, '14c70a7b1c6b95964f61fe2f4a7355e46c0181ca59735b90a7a5d12b6b84173f', '2026-01-30 03:16:30', NULL, '', '', '2026-01-23 03:16:30'),
(242, 17, 'f1cabc10d06be21cc68cb227d5319743014d250b12fb1ac2ab31c2e6e79a5e70', '2026-01-30 03:16:45', NULL, '', '', '2026-01-23 03:16:45'),
(243, 18, 'c6514490a9aab0e3246a05a05323de1c94217380485b3bdac41fde7d6c14f69d', '2026-01-30 03:22:05', NULL, '', '', '2026-01-23 03:22:05'),
(244, 19, '1882ea5da7eeab1e91ae7cd7a148488ef3561e56b07946f648ae662850cc3e08', '2026-01-30 03:25:21', NULL, '', '', '2026-01-23 03:25:21'),
(245, 20, '5f10450d93618c6c9a6370f22d294049f6316afca8a01def6dec62ee65315f5e', '2026-01-30 03:55:52', NULL, '', '', '2026-01-23 03:55:52'),
(246, 21, '73911e29c3bb4a63038b004bf43013d0e94b5cc7fd64016e09bf5aa16918dbd3', '2026-01-30 03:57:01', NULL, '', '', '2026-01-23 03:57:01'),
(247, 21, '9bfdbcae8acfa2a3f72f8ee3a13a255db527480ffa375f728a56eb05fc8158e9', '2026-01-30 06:12:17', NULL, '', '', '2026-01-23 06:12:17'),
(248, 21, '8e0e0245062a9c6216f29ba1acf6060d518ee91095eec4ff790ab8458c928445', '2026-01-30 06:27:54', NULL, '', '', '2026-01-23 06:27:54'),
(249, 21, '69f8f9534e58882c3b0dc8cfe0e3135fabddb2a249b4998e865d1cf8bcc2e1fc', '2026-01-30 07:08:18', NULL, '', '', '2026-01-23 07:08:18'),
(250, 21, '83f86411c0ce4091f81ed9a8d08ccdbec47f0b71553fba1a20ab005ecd4df937', '2026-01-30 07:36:41', NULL, '', '', '2026-01-23 07:36:41'),
(251, 21, '861766e692c693d9fc900bc7045f3e8c288d62859565285380d2bc26fae5773b', '2026-01-30 07:53:58', NULL, '', '', '2026-01-23 07:53:58'),
(252, 21, 'f917e55380b1a36aeb304a8d892772efe4c330c434d94a5f6987bd6b75fc3fed', '2026-01-30 08:24:27', NULL, '', '', '2026-01-23 08:24:27'),
(253, 21, 'f73bff41d1d74ef74d2ab37e8af0147fee97cd5b03677be2172db35a42e0be9b', '2026-01-30 08:57:31', NULL, '', '', '2026-01-23 08:57:31'),
(254, 21, '5524fcd331bc97301272b043b6d67fb507c1fb17d85b09f34501d57a03b33598', '2026-01-30 08:59:54', NULL, '', '', '2026-01-23 08:59:54'),
(255, 21, '6a01efdf5f6f28e6fe64dbe1af49031437f9ebe38ad92ca2016f0a99e9723038', '2026-01-30 09:16:42', NULL, '', '', '2026-01-23 09:16:42'),
(256, 21, 'bee82da4e07deed65c0bde5480f9d28af5c1a1a25be4917b8f0f62d2cec5ad12', '2026-01-30 09:21:32', NULL, '', '', '2026-01-23 09:21:32'),
(257, 21, '944aa02ccf3813c27f115ce280efef90d6cf6729f2f2aa3dbf27515e9a7db0c9', '2026-01-30 09:38:27', NULL, '', '', '2026-01-23 09:38:27'),
(258, 21, 'b273ebbaa8181d2fbb495d2acedd863c53614a995c2b050c9f67f09cb854d717', '2026-01-30 11:28:45', NULL, '', '', '2026-01-23 11:28:45'),
(259, 20, 'fa2f5f6045222d7224cc176af9d36cdd2cc104fbe7b2b8526ac522f96c78b6b4', '2026-01-30 11:31:15', NULL, '', '', '2026-01-23 11:31:15'),
(260, 20, '57de748240f983b0eec9d6da80fb3aa836cfbb532fc9d20d77145f0d7cf590d8', '2026-01-30 11:45:13', NULL, '', '', '2026-01-23 11:45:13'),
(261, 20, '0787a7a4b05c094527cdcb441a105449d3fba6fe4dfcdd78caf78e04325db02b', '2026-01-30 11:50:23', NULL, '', '', '2026-01-23 11:50:23'),
(262, 20, '44e9decfa5261618ba2faaa7f58e3338a47f4795ef93156929b601296b25cf48', '2026-01-30 11:59:54', NULL, '', '', '2026-01-23 11:59:54'),
(263, 20, '0247724f6eb228f5598759bbc72b449a70bb6aba284b966bd83115e0e9bbe730', '2026-01-30 12:05:44', NULL, '', '', '2026-01-23 12:05:44'),
(264, 20, '13880ac0a4d01947e7f344778b0f9eaefe37d49491d70fe64b1842102c75c4bc', '2026-01-30 13:28:54', NULL, '', '', '2026-01-23 13:28:54'),
(265, 20, 'fc1e3f3a81952c35e691ac207b3b2782205bf74569d2710f4d644834e584267e', '2026-01-30 13:48:24', NULL, '', '', '2026-01-23 13:48:24'),
(266, 21, 'b7ef78e370aaff9f958e46cb3bf42e99e4c1c82a23be9e96e5f38aa68038e58e', '2026-01-30 14:19:30', NULL, '', '', '2026-01-23 14:19:30'),
(267, 21, 'fb0661dfc82f59d10e2ddb18f0b1c382cda07f561718a519ddb93633a859769c', '2026-01-30 14:19:48', NULL, '', '', '2026-01-23 14:19:48'),
(268, 7, 'f7a6a328f65850b9d6ace64a97bd016756d2e30a41dd4d926be1949f36754d69', '2026-01-30 14:33:00', NULL, '', '', '2026-01-23 14:33:00'),
(269, 21, '2ddb2e4435332eac2c748e9480a33ebf9de730fc01222244ab1a338025a18443', '2026-01-30 14:36:14', NULL, '', '', '2026-01-23 14:36:14'),
(270, 21, '26a0b2d78c17b6dd2ee1b468c616e82f877200b3677630a2233e74b983fadb33', '2026-01-30 14:36:22', NULL, '', '', '2026-01-23 14:36:22'),
(271, 21, '89d57daeda884a7b22b871deca1762a6c14b262005954cd870937271a0917005', '2026-01-30 14:40:33', NULL, '', '', '2026-01-23 14:40:33'),
(272, 21, 'cb962d8942a2c7b56995706c0f0d07a0f87a970ae597e5375a1816aa318e79da', '2026-01-30 14:58:22', NULL, '', '', '2026-01-23 14:58:22'),
(273, 21, 'b0ab531095e6b9529583bb9893a9b297d830052a35988e44b5eb4a1dd8381d18', '2026-01-30 14:58:52', NULL, '', '', '2026-01-23 14:58:52'),
(274, 7, '99b4833802aa218144162da41ee89886e8f943dd817149fffebfc3bfef53e82d', '2026-01-30 15:00:57', NULL, '', '', '2026-01-23 15:00:57'),
(275, 21, '6096947f77168d2a5bedb82b89ada3f9a4802d1902583cd526852457f4d6c284', '2026-01-30 15:04:56', NULL, '', '', '2026-01-23 15:04:56'),
(276, 7, '718cac39ba563ca911466b121460da6d5f6189baa4032e5581aed9ddfeaf52c3', '2026-01-30 15:06:16', NULL, '', '', '2026-01-23 15:06:16'),
(277, 21, '4eba517db6cf19ec8aed14090381f9dde08017f755d2f09ee2d9eb3386af0649', '2026-01-30 15:06:59', NULL, '', '', '2026-01-23 15:06:59'),
(278, 7, 'c2cfee53f4f91b592bc29478d098c1dd74a7a4c81268eb1f00965859c25ffc3d', '2026-01-30 15:06:59', NULL, '', '', '2026-01-23 15:06:59'),
(279, 7, 'da2a5e36b2aa993b9c9337722794b0bbf41f9df7fbb08d0021fff272e89cc10f', '2026-01-31 23:16:08', NULL, '', '', '2026-01-24 23:16:08'),
(280, 7, 'f82c8aea003354831db28cea2d8d861b8d9064861c477ea36bd4ab677142c2b4', '2026-01-31 23:24:30', NULL, '', '', '2026-01-24 23:24:30'),
(281, 7, '8a575f23bd722b21f92a892ee9ee17e501ab76cd2667f066fa0ef8baa3065878', '2026-01-31 23:37:29', NULL, '', '', '2026-01-24 23:37:29'),
(282, 7, 'f81ecc757f175ffae0f0c04d3936af5dff47b4e6747d099a74031139fa6dd1c3', '2026-01-31 23:37:46', NULL, '', '', '2026-01-24 23:37:46'),
(283, 7, '25c7144e8391881835ef51ea92b168d9f46a4f94cc1ef8448c214d09122ea49d', '2026-01-31 23:37:54', NULL, '', '', '2026-01-24 23:37:54'),
(284, 20, 'e741e2860dd9c8ed963f2d61e4a9831eb838932125719cab5a0ab1e9781ed113', '2026-01-31 23:38:11', NULL, '', '', '2026-01-24 23:38:11'),
(285, 7, '61e9709d79758bfe808b2e6d165e36eb07e4f4ecafac7bce7aac17354361d962', '2026-01-31 23:42:12', NULL, '', '', '2026-01-24 23:42:12'),
(286, 7, '423b7ec219df95732c885451002741c62e204811ef86f7cb59e14ca838259638', '2026-01-31 23:45:52', NULL, '', '', '2026-01-24 23:45:52'),
(287, 7, '7a5202a7ae632de482b1ebdd9a8b45d7b53ed517b94ac0f5893e7041bab4fe80', '2026-01-31 23:59:42', NULL, '', '', '2026-01-24 23:59:42'),
(288, 7, 'b8c3833bc4688c0267f84096c32d81d509b6baecfa250c51d49912280ffe9a32', '2026-02-01 00:00:34', NULL, '', '', '2026-01-25 00:00:34'),
(289, 7, 'e30651901378e4bc93a3e306f839ef995559bae7cf5881d2ee148446d275ea87', '2026-02-01 00:08:09', NULL, '', '', '2026-01-25 00:08:09'),
(290, 7, '6dc6a7d588ebecdd270bcb3a4ca7dc67b12d2921cd414edb817faa74c9860bdd', '2026-02-01 00:24:28', NULL, '', '', '2026-01-25 00:24:28'),
(291, 7, '4731906e1e7a5a545c5530f639a0593860053e58870cb6f9c10e98e606a6dee8', '2026-02-01 00:24:43', NULL, '', '', '2026-01-25 00:24:43'),
(292, 20, '52a324d267df2f102f2c5bb68e482ff96cfeefdfc458911ec3f045e0fdee9795', '2026-02-01 00:29:07', NULL, '', '', '2026-01-25 00:29:07'),
(293, 7, '1550fc04be1cd69296eb1cd048c65c774cc90d6b2725686e724a799bbeb463bf', '2026-02-01 00:30:06', NULL, '', '', '2026-01-25 00:30:06'),
(294, 7, '40e333cc11c1b56af66ce50157d3d4890cd0544469685e76a8f9e62b1203848b', '2026-02-01 01:13:08', NULL, '', '', '2026-01-25 01:13:08'),
(295, 7, '3bfdb84c07bf2c76bd4659324d8bbd06f4bd8202be6c00efef8ca323ce781b7c', '2026-02-01 01:27:43', NULL, '', '', '2026-01-25 01:27:43'),
(296, 20, '379013a4ccc60e0476cf5c22cade0c876d4f5a0d517ca3935b3f6592e2ea3dd9', '2026-02-01 01:52:03', NULL, '', '', '2026-01-25 01:52:03'),
(297, 7, 'b3bb609b9f8ec4677c55456a58b06fb7b2411cf27c3995e2e93d7f7726e69d7b', '2026-02-01 01:52:09', NULL, '', '', '2026-01-25 01:52:09'),
(298, 7, 'eb110b3c2ad7589d529b24e5fb7933cb708a480e9a424bc1f7fe84b7b15259ae', '2026-02-01 02:09:40', NULL, '', '', '2026-01-25 02:09:40'),
(299, 7, 'e6bb51c2791a7cc140a57da002d379120cd7e487f3c0b68e07174380340329d8', '2026-02-01 03:37:17', NULL, '', '', '2026-01-25 03:37:17'),
(300, 20, '768a2e7216775abb997f075967fef6c4519c1347cf7674bac7f8da5ff6bef086', '2026-02-01 09:07:15', NULL, '', '', '2026-01-25 09:07:15'),
(301, 7, '8eb30b402d8e308760403de2f5bab96e97e00ebc2ba1b118b6449e9acc2881f6', '2026-02-01 09:07:25', NULL, '', '', '2026-01-25 09:07:25'),
(302, 7, 'cbe5b0d6553b5b415e4e75a21f48e271f459c9fc9672fceb0eba52d5ad14d512', '2026-02-01 09:22:46', NULL, '', '', '2026-01-25 09:22:46'),
(303, 7, 'bc535d4835d803861c5d64d03b1202c7f28a4c1f2234a133742db021c08e8edc', '2026-02-01 09:24:27', NULL, '', '', '2026-01-25 09:24:27'),
(304, 20, '1807f455b996c42afbca2dbd747c9ff634fb13e0cc6590edce8b704c2c7456ea', '2026-02-01 09:24:50', NULL, '', '', '2026-01-25 09:24:50'),
(305, 21, '1a3eda06d040309c20080bcc791e970508b28f12d8ad8946b49493c3953ae760', '2026-02-01 09:45:49', NULL, '', '', '2026-01-25 09:45:49'),
(306, 21, '3931a247777020df9988828a3a4f9a5000a83e74016eea4311dde53b15aabb71', '2026-02-01 09:48:17', NULL, '', '', '2026-01-25 09:48:17'),
(307, 21, '9b119c7afff158b192b1586f4ae43610013a9c5aee2dbf3508c305d7425d1296', '2026-02-01 09:48:26', NULL, '', '', '2026-01-25 09:48:26'),
(308, 20, 'a2a7347d04d390ae5d32080db61c465ad9e1c4b822a05f628b956eefa46e8027', '2026-02-01 09:52:37', NULL, '', '', '2026-01-25 09:52:37'),
(309, 20, '6ac7374e32249c12df06a9b0eb94df9eba0915b76af88bd26d5631b92d4599cb', '2026-02-01 09:54:18', NULL, '', '', '2026-01-25 09:54:18'),
(310, 20, 'edd9384c117a3c8d1fb85f937065b6c50c73347235449d56fd7468de64c3275c', '2026-02-01 10:25:32', NULL, '', '', '2026-01-25 10:25:32'),
(311, 20, 'b0e2fbe9e0017963bf138e3d02ae633dd776e1bc89c0d07fbc231821d6219513', '2026-02-01 10:44:22', NULL, '', '', '2026-01-25 10:44:22'),
(312, 21, '34e2da26bd7d745b14fd4b473b3ab1d19538bfd0f9379e351e666fc8bdd797d0', '2026-02-01 10:52:10', NULL, '', '', '2026-01-25 10:52:10'),
(313, 20, 'b14ad309d3a583c80e263399dfbdea488448b14e2e028b51b8f6a59967f431f4', '2026-02-01 11:01:10', NULL, '', '', '2026-01-25 11:01:10'),
(314, 21, 'cac98a7ab5722a1389ae36d8abedaeb584a0e0768f1b3efefe938695e3c1c1f7', '2026-02-01 11:03:34', NULL, '', '', '2026-01-25 11:03:34'),
(315, 21, 'a4804175f4b80818582aa88013c149c481ddb0b835c600bce4dd668391cac98c', '2026-02-01 11:05:26', NULL, '', '', '2026-01-25 11:05:26'),
(316, 20, 'f53cd8ca86ca5f7c0d8a1328c5581e54423745636f14e16828db596854028c4e', '2026-02-01 11:12:48', NULL, '', '', '2026-01-25 11:12:48'),
(317, 20, 'fc86a3def959f17006acedee78a34f945ab8eb3f4d186b46f8de7208849d9a57', '2026-02-01 11:12:57', NULL, '', '', '2026-01-25 11:12:57'),
(318, 20, 'e7912cba5ba1a6c9477a94b4b140f1dc1130194ab5bc1e64742e7f07ca81b559', '2026-02-01 11:23:24', NULL, '', '', '2026-01-25 11:23:24'),
(319, 20, '05cbf290b708af564f4a6ebd51e1ee16450f654b701e1f62c34e7f9daee2038d', '2026-02-01 11:26:41', NULL, '', '', '2026-01-25 11:26:41'),
(320, 20, 'e712abeb9714c6f3a2d4774f56db2333b9a27592b33182883a7c4daf64c271fb', '2026-02-01 11:27:32', NULL, '', '', '2026-01-25 11:27:32'),
(321, 20, '178f21e9cb9ee7006b8c9a5b7c1f8a4c44d27c7f55f15abd1a3aba8be36df4f8', '2026-02-01 12:04:16', NULL, '', '', '2026-01-25 12:04:16'),
(322, 21, 'c620eb72bdb0aaa00ed2707ffc0d6a680c101295b4ddf38177fc955c56142690', '2026-02-01 12:05:47', NULL, '', '', '2026-01-25 12:05:47'),
(323, 21, '225759c0ef695299de29e97163efe07dbf4c5cf31ba81224b947f332435bb2f1', '2026-02-01 17:29:04', NULL, '', '', '2026-01-25 17:29:04'),
(324, 7, 'fe241bd243ceba94c0ee23dcd2add2f83a83c67680132dd26ea98786d9282272', '2026-02-01 17:30:12', NULL, '', '', '2026-01-25 17:30:12'),
(325, 21, '21620047c8bae3aa6d6a3d319afa8ecae6c831298979ffe4344d61610a628763', '2026-02-01 20:19:21', NULL, '', '', '2026-01-25 20:19:21'),
(326, 21, 'e712ab03dc4ac07e6ba7b4ed61121936522ff3d1a093d4ce1490e4bbd0a22944', '2026-02-01 20:24:32', NULL, '', '', '2026-01-25 20:24:32'),
(327, 7, 'b066d9ca95967a64eecd118df070503aed7ce127e9f4cac7bdc3dc7539731038', '2026-02-01 21:19:04', NULL, '', '', '2026-01-25 21:19:04'),
(328, 21, '54fd7634e1083c87d38fa5017a60c6a1cf761bb92427183e20cfc25c7746f33d', '2026-02-01 21:20:16', NULL, '', '', '2026-01-25 21:20:16'),
(329, 7, '8872cd7f54fe60d0b757c2ae6bcfc0432bdc74c4f965a88f409d0147ac7c5d4e', '2026-02-01 21:37:26', NULL, '', '', '2026-01-25 21:37:26'),
(330, 7, '047b3b006e812b20eb6ac09573f572a7419739709a17595a01b4068744792cac', '2026-02-01 21:38:35', NULL, '', '', '2026-01-25 21:38:35'),
(331, 7, '96baed29cca6bc2abc014cebb158ea1dec608e9bd55b723c9385c680acf22481', '2026-02-01 21:51:23', NULL, '', '', '2026-01-25 21:51:23'),
(332, 7, '4d18ed7a197d062631d208d615bac83b22be70a5e6d01958aa011e68d2c06c11', '2026-02-01 21:51:32', NULL, '', '', '2026-01-25 21:51:32'),
(333, 21, 'd6b2a961aaf92800c638bb7188002419d3405e09b2d799ed42c84f2679bf9c24', '2026-02-02 02:48:39', NULL, '', '', '2026-01-26 02:48:39'),
(334, 7, 'bb81c76720ae5a103401b2e3278310d9842d22cfd6f077cf6c06273377c84607', '2026-02-02 02:59:47', NULL, '', '', '2026-01-26 02:59:47'),
(335, 20, '1b915f33f63b55988327001cc2c6b36c0350e1995e062ede4e546bc6e74775ab', '2026-02-02 03:00:12', NULL, '', '', '2026-01-26 03:00:12'),
(336, 20, 'c71fb63257f457e61f8d48813b47d263f3b52b771706077c358bc677dea62b63', '2026-02-02 03:24:01', NULL, '', '', '2026-01-26 03:24:01'),
(337, 7, '78b95ae4cccbf9b67d2dd5d1ea7eeb249cae73637c466ca83f40a75b38a40744', '2026-02-02 03:28:56', NULL, '', '', '2026-01-26 03:28:56'),
(338, 7, '3bb0257b8a8601c941435977eb4622b1a597a16cc076108a49d6055626a5d1f8', '2026-02-02 03:32:58', NULL, '', '', '2026-01-26 03:32:58'),
(339, 7, 'a4be89529c78d71b2f0ba543abd1e60873f17e8e4bf82961f15b3d8638d29053', '2026-02-02 03:33:06', NULL, '', '', '2026-01-26 03:33:06'),
(340, 7, 'a143dbb159d2fb6e0e75d457ca6b6547f68de0a7a6f7abb5dc1fc49be496cade', '2026-02-02 04:47:41', NULL, '', '', '2026-01-26 04:47:41'),
(341, 7, 'db6a850654fa4f46c79070bbd88d93ff78ffaa41efbd2b21aa73fb047f452ce2', '2026-02-02 04:48:43', NULL, '', '', '2026-01-26 04:48:43'),
(342, 7, 'c64d69dfa4bd2d4172a429b266be3f16839d5ae6b93e788c8af5cb17017cd341', '2026-02-02 04:48:52', NULL, '', '', '2026-01-26 04:48:52'),
(343, 7, '6437e2a37132dc6bdfbf50f9dd784b339c999b77bea28852891366a1ddb56175', '2026-02-02 04:51:59', NULL, '', '', '2026-01-26 04:51:59'),
(344, 7, '5dd2166c1d863eb2534a2fbb21129f471cfaf82477a65a185219d5370cabd99e', '2026-02-02 04:52:27', NULL, '', '', '2026-01-26 04:52:27'),
(345, 7, 'f10c2d671d08abe213691bc495ad34c232d11a1f4bfde5250c39c3a5625d3e66', '2026-02-02 05:08:46', NULL, '', '', '2026-01-26 05:08:46'),
(346, 7, '0959788c09231243905ea5d5b9859457446ce9b7c13d9da5b5c46c6a1726ef1e', '2026-02-02 05:08:57', NULL, '', '', '2026-01-26 05:08:57'),
(347, 7, 'a675f63fdbdcb10d1425339797fc098ccaab7c4baa1402710719d64d48e43357', '2026-02-02 05:09:23', NULL, '', '', '2026-01-26 05:09:23'),
(348, 7, 'c37ba039b1e0caab1e53bd0fbfdb996da249f1cfc616d2a1c9ee9b62a7c5dc3c', '2026-02-02 05:09:34', NULL, '', '', '2026-01-26 05:09:34'),
(349, 20, '1f61165a22b7211a7a91ed61e785f9322c30093b6828ff65fe1c55185efa3ca3', '2026-02-02 05:13:24', NULL, '', '', '2026-01-26 05:13:24'),
(350, 21, '289c2c4a26a1fb7af7e349182cbc4c3ddaa42076fb629571d6de43a26b162201', '2026-02-02 06:26:59', NULL, '', '', '2026-01-26 06:26:59'),
(351, 21, '261132984785ea0d05aa958caab4477f599b62770c8b850a2d4864871f7bb11f', '2026-02-02 06:27:08', NULL, '', '', '2026-01-26 06:27:08'),
(352, 1, 'ed2025aa882f2dca84d95e6a7e9d5ba2d1cd1da902a85993a1b8c153b63e888d', '2026-02-02 06:30:07', NULL, '', '', '2026-01-26 06:30:07'),
(353, 20, 'e6311dd8bbe5fe2e39654eacddeb995857e0681c8684c891f5887e4060ec1cd1', '2026-02-02 06:44:00', NULL, '', '', '2026-01-26 06:44:00'),
(354, 20, '236e5a98fe2850ad41657b90b893238a709c66fe15a8361500190b6ca5b55f2a', '2026-02-02 06:48:29', NULL, '', '', '2026-01-26 06:48:29'),
(355, 20, '345c697e4036a6e108024ecc6e34f69c078bffc8badff1e3f95210736ee90fd5', '2026-02-02 07:10:23', NULL, '', '', '2026-01-26 07:10:23'),
(356, 7, '2546d8fec0d96b91cd776c5c5a68795bbeca75c69785f75374def7daf2572180', '2026-02-02 07:13:09', NULL, '', '', '2026-01-26 07:13:09'),
(357, 7, 'b32e0d2dcf5a75ad4ef8f01646b599eee916b35f7464835c3d49ecc56145c69a', '2026-02-06 04:18:59', NULL, '', '', '2026-01-30 04:18:59'),
(358, 7, '360da15c72433e59aa7698560dc0c78767102ba1447376d3251675b83f2f3f7e', '2026-02-07 01:36:40', NULL, '', '', '2026-01-31 01:36:40'),
(359, 7, 'ee1dcb327902686dd0ec877a65c976d060dc7ddfc20478489b7a46413aa7490b', '2026-02-07 01:59:59', NULL, '', '', '2026-01-31 01:59:59');
INSERT INTO `refresh_tokens` (`id`, `user_id`, `token_hash`, `expires_at`, `revoked_at`, `device_info`, `ip_address`, `created_at`) VALUES
(360, 7, '9d5ba8dae79c95241f4514f4550ab7661bd70b21617f7a8b4dd819c188b3ee57', '2026-02-08 00:38:54', NULL, '', '', '2026-02-01 00:38:54'),
(361, 7, '313f4e207021a86a07b1efe48c2153f5b04c055b35a7a5802627053ffb70ab91', '2026-02-08 00:56:58', NULL, '', '', '2026-02-01 00:56:58'),
(362, 7, '607d6e87c15de351b2faa8ce972c7a8eee47f42e8a955223a724f68696843004', '2026-02-08 01:00:01', NULL, '', '', '2026-02-01 01:00:01'),
(363, 1, '0dab125fe9b95d94c9a56f187f8a3cbec5ca7a33f656a8cead97af2004553660', '2026-02-09 01:55:20', NULL, '', '', '2026-02-02 01:55:20'),
(364, 1, '8134b8ea94b4d42820de99bdfacc22f32b92069d58fbdf19aae774befd5dcb2b', '2026-02-09 01:55:30', NULL, '', '', '2026-02-02 01:55:30'),
(365, 1, '59dc1532ea95b007dff9023bb26d1b827ec764e88c7523ac1459f804ccca5019', '2026-02-09 01:55:35', NULL, '', '', '2026-02-02 01:55:35'),
(366, 1, 'da9440d483a15c781481cc235bd4f9a188675637f7d3ed674f2196361dbfa46c', '2026-02-09 01:55:40', NULL, '', '', '2026-02-02 01:55:40'),
(367, 1, '54994e5ce6b68bb0e78c9aaf852c3abd138a6ad4b526ddf34b29823b31091313', '2026-02-09 01:55:49', NULL, '', '', '2026-02-02 01:55:49'),
(368, 1, 'cfda9284d301b0a005d4742a03a8a0c7e93dd8bfa5e80392f59feb6b9d8bf06b', '2026-02-09 02:08:39', NULL, '', '', '2026-02-02 02:08:39'),
(369, 1, '770edad672790115fa18026ab30f90832bcbffa29a8339200e5b6938aa189b00', '2026-02-09 02:10:51', NULL, '', '', '2026-02-02 02:10:51'),
(370, 1, '12c77e9c99d60425a24b0c6bae536114ab5e3d2245f6b4e33fc87714d5d99990', '2026-02-09 02:12:39', NULL, '', '', '2026-02-02 02:12:39'),
(371, 1, 'dc08f81903066b1a2cb570f4854339498a5bdbe822c1e55e6e1fdda843276d71', '2026-02-09 02:12:48', NULL, '', '', '2026-02-02 02:12:48'),
(372, 1, 'f62a59d2e140a0dd1a0c84e330e727d18956751536a912be273f7a3f87d50d13', '2026-02-09 02:12:57', NULL, '', '', '2026-02-02 02:12:57'),
(373, 1, '77445f6ccaeb2dfb463f62035cfd69156759c873d8d6a78e4f1aa99ee3e27d80', '2026-02-09 02:14:35', NULL, '', '', '2026-02-02 02:14:35'),
(374, 1, '7079815374c3fa589d447644518fbf04af285fb34fb1e3b07fdbfa6e3e66d492', '2026-02-09 03:41:25', NULL, '', '', '2026-02-02 03:41:25'),
(375, 1, '2a57efac4d5f939ba9502ce8cbac16cc1fd0e3b8e94d301d2215117485ee9972', '2026-02-09 03:41:53', NULL, '', '', '2026-02-02 03:41:53'),
(376, 7, 'bdc7d3f6571da6593da7ae4d54fc8912a0335c18125f49288ac433af7439c07b', '2026-02-09 03:57:09', NULL, '', '', '2026-02-02 03:57:09'),
(377, 7, 'efe2cf2f4daf5c58020fb39a5f6d003940c73e04a397ab5ef368ae3893321695', '2026-02-09 04:03:21', NULL, '', '', '2026-02-02 04:03:21'),
(378, 7, 'f59edb9ff94ffbd2335d2a5016bdf1b8f6c5b04f60d892b43578d01d7f41475a', '2026-02-09 04:03:28', NULL, '', '', '2026-02-02 04:03:28'),
(379, 1, 'c7e4fc27446591c59902c5bd57b5eba049f11ba3f62de1a95d26b811b9db184e', '2026-02-09 04:04:27', NULL, '', '', '2026-02-02 04:04:27'),
(380, 1, '68350a6243f729b1215e6284ddd8c7e50ada658447a2193df3769a124d81e389', '2026-02-09 04:04:33', NULL, '', '', '2026-02-02 04:04:33'),
(381, 7, '32778bc04f0ec8134320757c310fbb0e0d1d4e72a242c7583979180019ff9fa5', '2026-02-09 04:04:49', NULL, '', '', '2026-02-02 04:04:49'),
(382, 7, '86056abb50f1a0098a1b695d39675a504b49457d1d258a7f1ea9be21242f3870', '2026-02-09 04:07:21', NULL, '', '', '2026-02-02 04:07:21'),
(383, 7, 'd1f4a5d994bfc80c0d40ddc8d70be98fc2739ba68cab3958a69d4d2694482cbf', '2026-02-09 04:07:26', NULL, '', '', '2026-02-02 04:07:26'),
(384, 1, '0817b12ef0bb71a1558fe6d584dca1fd7ac5aa22eb5803283b41e95bdb89fd12', '2026-02-09 04:07:50', NULL, '', '', '2026-02-02 04:07:50'),
(385, 1, 'fc414666a8a324b14b10aaa345cf2c5760cef549f7a578ff08bb08b04e3e8e21', '2026-02-09 04:07:57', NULL, '', '', '2026-02-02 04:07:57'),
(386, 7, 'dce99fdef5dbf57b637c5322e19fa61a28545e875aad070809da4dadc7883955', '2026-02-09 04:12:25', NULL, '', '', '2026-02-02 04:12:25'),
(387, 7, 'dd014bc644d44b13dc69ff99aeefcdc4662f4d8ffbc2b9da454b71f9ec06d457', '2026-02-09 13:20:31', NULL, '', '', '2026-02-02 13:20:31'),
(388, 7, '55c7a011479c79f08e9053ea286f45d22f6fd2854dea5dccd5ab2d5953434066', '2026-02-09 14:14:01', NULL, '', '', '2026-02-02 14:14:01'),
(389, 7, 'd2fb3aed3e20821b0f54bbb5e3984f7bc7eefc2881677e12756d5608e0e5319b', '2026-02-09 14:45:49', NULL, '', '', '2026-02-02 14:45:49'),
(390, 1, 'dd0a08d9789328573a159c36bf9d59097569e1948bd19b57597e2b3f709c4fe7', '2026-02-10 16:22:03', NULL, '', '', '2026-02-03 16:22:03'),
(391, 20, 'a264abe792e99d466b6ea405cbc19adc8586da0e478ee2116e1e69303a120f06', '2026-02-11 03:34:14', NULL, '', '', '2026-02-04 03:34:14'),
(392, 7, 'ac008b86050f916528231b5d5856fe0ff698fafc7b35730469f48b435be040e8', '2026-02-11 12:44:49', NULL, '', '', '2026-02-04 12:44:49'),
(393, 7, 'f7b1731fe4e9ffdb40f2d1c0c39abcdd10067c544992f54760aab7f9a3cc9019', '2026-02-11 13:35:26', NULL, '', '', '2026-02-04 13:35:26'),
(394, 1, '64f0bb80795e6b91a16af22c04f6b779e09421c5164afca6ab1bac200c9a8b32', '2026-02-15 02:27:49', NULL, '', '', '2026-02-08 02:27:49'),
(395, 1, 'ff6253aa143ac1f39df6aab18f340e3db4f635e1445d72bc77bc9d9097e01085', '2026-02-15 02:27:57', NULL, '', '', '2026-02-08 02:27:57'),
(397, 7, '2ec87638e251b899ee3d4bd9102106043438e818394d2be2f6399a7471c275fc', '2026-02-16 13:05:36', NULL, '', '', '2026-02-09 13:05:36'),
(398, 26, '43754c218c9199dc862b55237f9bc9659c2fc0dc9c43b86a6ef59e0e0d0529bc', '2026-02-16 13:08:22', NULL, '', '', '2026-02-09 13:08:22'),
(399, 1, '3ab612eee2904c8beb05ad34b090a65766c04b4951ba631c4ec72b7c5f7fbb12', '2026-02-16 13:12:30', NULL, '', '', '2026-02-09 13:12:30'),
(400, 1, '58ced0af41798dd1dd9512e511c25b2aa894b7791683e58ed2c3cfb56ccadf2d', '2026-02-16 13:14:33', NULL, '', '', '2026-02-09 13:14:33'),
(401, 1, 'e150ff2b18b34b4c586306ced552fe35dbafbec9d62c003d8665aa3b92726901', '2026-02-16 13:30:05', NULL, '', '', '2026-02-09 13:30:05'),
(402, 26, 'e821c4ed5f82ce0a66343f393f90fed1aef0aa41c959994490fcca642db818f0', '2026-02-18 11:47:32', NULL, '', '', '2026-02-11 11:47:32'),
(403, 26, '5c3400d08407132da6159c3e6f4203d4e134b2ecad4b381a7a8fa4440f22d3d3', '2026-02-18 12:03:42', NULL, '', '', '2026-02-11 12:03:42'),
(404, 26, '1bcf0fe9867b6bcd5c92285b8ccd5197288ab80a7a2894c6c0fd7f1731f8274e', '2026-02-18 14:48:45', NULL, '', '', '2026-02-11 14:48:45'),
(405, 20, '07b5f140fddbb87e90fd033838a9a985eaec9346b22322636b1c287f8ccbb057', '2026-02-18 14:56:22', NULL, '', '', '2026-02-11 14:56:22'),
(406, 26, '8a0ae6b88070879c151078413d895dda9509fcfaf51eaebf3bce444d57b75478', '2026-02-18 14:57:52', NULL, '', '', '2026-02-11 14:57:52'),
(407, 26, '4ce874f456c470911d77469b96eb33d94cd8cb64cab42ec2e0c9c89f7ec3daa6', '2026-02-18 15:11:25', NULL, '', '', '2026-02-11 15:11:25');

-- --------------------------------------------------------

--
-- Table structure for table `saved_jobs`
--

CREATE TABLE `saved_jobs` (
  `id` bigint(20) UNSIGNED NOT NULL,
  `user_id` bigint(20) UNSIGNED NOT NULL,
  `job_id` bigint(20) UNSIGNED NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT current_timestamp()
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

--
-- Dumping data for table `saved_jobs`
--

INSERT INTO `saved_jobs` (`id`, `user_id`, `job_id`, `created_at`) VALUES
(1, 20, 56, '2026-01-25 10:44:38'),
(2, 21, 56, '2026-01-25 17:29:18');

-- --------------------------------------------------------

--
-- Table structure for table `support_tickets`
--

CREATE TABLE `support_tickets` (
  `id` bigint(20) UNSIGNED NOT NULL,
  `user_id` bigint(20) UNSIGNED NOT NULL,
  `title` varchar(100) NOT NULL,
  `description` text NOT NULL,
  `category` varchar(50) NOT NULL DEFAULT 'other',
  `priority` enum('low','medium','high','urgent') NOT NULL DEFAULT 'medium',
  `status` enum('open','in_progress','pending_response','resolved','closed') NOT NULL DEFAULT 'open',
  `email` varchar(255) NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT current_timestamp(),
  `updated_at` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  `resolved_at` timestamp NULL DEFAULT NULL,
  `closed_at` timestamp NULL DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

--
-- Dumping data for table `support_tickets`
--

INSERT INTO `support_tickets` (`id`, `user_id`, `title`, `description`, `category`, `priority`, `status`, `email`, `created_at`, `updated_at`, `resolved_at`, `closed_at`) VALUES
(3, 20, 'Testing Ticket', 'Hello bang', 'cv-builder', 'low', 'open', 'jastiska14@gmail.com', '2026-01-26 06:53:30', '2026-01-26 06:53:30', NULL, NULL);

-- --------------------------------------------------------

--
-- Table structure for table `ticket_responses`
--

CREATE TABLE `ticket_responses` (
  `id` bigint(20) UNSIGNED NOT NULL,
  `ticket_id` bigint(20) UNSIGNED NOT NULL,
  `sender_id` bigint(20) UNSIGNED NOT NULL,
  `sender_type` enum('user','admin') NOT NULL DEFAULT 'user',
  `message` text NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT current_timestamp()
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

--
-- Dumping data for table `ticket_responses`
--

INSERT INTO `ticket_responses` (`id`, `ticket_id`, `sender_id`, `sender_type`, `message`, `created_at`) VALUES
(5, 3, 20, 'user', 'Hello bang', '2026-01-26 06:53:30');

-- --------------------------------------------------------

--
-- Table structure for table `users`
--

CREATE TABLE `users` (
  `id` bigint(20) UNSIGNED NOT NULL,
  `email` varchar(255) NOT NULL,
  `password_hash` varchar(255) NOT NULL,
  `role` enum('job_seeker','company','admin','partner') NOT NULL DEFAULT 'job_seeker',
  `full_name` varchar(255) NOT NULL,
  `phone` varchar(20) DEFAULT NULL,
  `avatar_url` varchar(500) DEFAULT NULL,
  `is_active` tinyint(1) NOT NULL DEFAULT 1,
  `is_verified` tinyint(1) NOT NULL DEFAULT 0,
  `email_verified_at` timestamp NULL DEFAULT NULL,
  `created_at` timestamp NOT NULL DEFAULT current_timestamp(),
  `updated_at` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  `deleted_at` timestamp NULL DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

--
-- Dumping data for table `users`
--

INSERT INTO `users` (`id`, `email`, `password_hash`, `role`, `full_name`, `phone`, `avatar_url`, `is_active`, `is_verified`, `email_verified_at`, `created_at`, `updated_at`, `deleted_at`) VALUES
(1, 'admin@karirnusantara.com', '$2a$10$llmqSju0EZuHEklAzhE7jeubjam0w2DoG252O2cgb1lG73gGL0AEG', 'admin', 'Admin User', NULL, NULL, 1, 0, NULL, '2026-01-19 08:53:05', '2026-01-19 08:53:05', NULL),
(2, 'hr@techcorp.id', '$2a$10$weTcDM1VtUZPmsGF/olunOItpWfuD8rKJ5NcHgsn1WvUuMM0Ls.S.', 'company', 'HR TechCorp', '081234567890', NULL, 1, 0, NULL, '2026-01-17 02:44:45', '2026-01-17 02:44:45', NULL),
(3, 'budi.kandidat@gmail.com', '$2a$10$8jpAPiQoMmKhsyKE1nZJ6u0suj/To00/xPNbx5X94ARBUpUmgJj9a', 'job_seeker', 'Budi Santoso', NULL, NULL, 1, 0, NULL, '2026-01-17 02:50:35', '2026-01-17 02:50:35', NULL),
(4, 'testcompany@test.com', '$2a$10$Owbnmfiz/RUlALkTaI777epAeVCCvSsZmd.bqOu.hjBX1KHGsNYhW', 'company', 'HR Manager', '081234567890', NULL, 1, 0, NULL, '2026-01-18 07:08:53', '2026-01-18 07:08:53', NULL),
(5, 'company2@test.com', '$2a$10$26QNIFt3ErRbLdMEaIyu1uqBfxn72fmYYTK27e1NnK5UaAEsZiEvy', 'company', 'Test Manager', '081234567890', NULL, 1, 0, NULL, '2026-01-18 07:42:02', '2026-01-18 07:42:02', NULL),
(6, 'company3@test.com', '$2a$10$27XxXTvPq9S.tOwT5/kY7u6E7V8ovsBq3jGQiL5txJuCzVm0t0iNa', 'company', 'Recruitment Manager', '081555555555', NULL, 1, 0, NULL, '2026-01-18 07:43:28', '2026-01-18 07:43:28', NULL),
(7, 'info@karyadeveloperindonesia.com', '$2a$10$k52RcefcRAvQkxEcpoe7Eu8/xUA.5tKGIP1gW9kz9Pjklizpc1VEK', 'company', 'Admin', '0881036480285', NULL, 1, 1, NULL, '2026-01-18 09:13:55', '2026-01-21 15:45:14', NULL),
(9, 'company.testing1768808682@karirnusantara.com', '$2a$10$1RpgFQaG58KBk.Fhy8laye2yhBj.C.ogGGaqxkktzmew7TMrlLu.W', 'company', 'Budi Santoso', '081234567890', NULL, 1, 0, NULL, '2026-01-19 07:44:42', '2026-01-19 07:44:42', NULL),
(10, 'company.testing1768808883@karirnusantara.com', '$2a$10$GrXGCrL6oNq.NJAJwmhAg.aQMID5txWxyP7Df/XJuVQ39zSKhioDy', 'company', 'Budi Santoso', '081234567890', NULL, 1, 0, NULL, '2026-01-19 07:48:03', '2026-01-19 07:48:03', NULL),
(12, 'test.company@example.com', '$2a$10$Q7OZlZckA0Tw5RarsLHdHOctnn3n9Sh2J.SM1OwVX58DxtA8Q74le', 'company', 'Test Company Manager', '081234567890', NULL, 1, 0, NULL, '2026-01-20 04:39:36', '2026-01-20 04:39:36', NULL),
(13, 'test@example.com', '$2a$10$PoRfDCFO106yVLpJxqDK.uNcETZcse0VoG.m7efDEvi.KmUzD003y', 'company', 'Test User', NULL, NULL, 1, 0, NULL, '2026-01-21 15:06:03', '2026-01-21 15:06:03', NULL),
(14, 'testchangepass@example.com', '$2a$10$1NXdnUsjN3EENJb4PnSiYuQFANGpmFUEQVB/CR6/aexVSHdlDoDUu', 'company', 'Test Change Password', NULL, NULL, 1, 0, NULL, '2026-01-21 15:31:15', '2026-01-21 15:31:30', NULL),
(15, 'changepass_1769009652@example.com', '$2a$10$AanEN8CC9DSYdlh/2hWjFu4DZXTFFrLtD7N89kpJ5B6R1kkp7YuH6', 'company', 'Change Password Test User', NULL, NULL, 1, 0, NULL, '2026-01-21 15:34:13', '2026-01-21 15:34:13', NULL),
(16, 'test.register@example.com', '$2a$10$JEWbMAD3Vbzv4vkt0PGM1OgGAhkhIINOLSgjoCWSjGOiQ2ZYVmtfK', 'job_seeker', 'Test User', '0881036480285', NULL, 1, 0, NULL, '2026-01-23 03:16:19', '2026-01-23 03:16:19', NULL),
(17, 'test.register123@example.com', '$2a$10$wxnxA3HbaDZSz0Y4gjKYdebDsc5PIghMTWl3BtlkavoqL81VD3H/K', 'job_seeker', 'Test User', '0881036480285', NULL, 1, 0, NULL, '2026-01-23 03:16:30', '2026-01-23 03:16:30', NULL),
(18, 'debug.test@example.com', '$2a$10$GUhOXJgQDIqJxRw.irMOa.D1NLbxjIe9aNgUDacALSB1mCzNBgwoe', 'job_seeker', 'Debug Test', NULL, NULL, 1, 0, NULL, '2026-01-23 03:22:05', '2026-01-23 03:22:05', NULL),
(19, 'emailtest123@example.com', '$2a$10$uCnPcKmTLE2C5nyPUMWmReEUr0o.D1eJKke5iAn7ArvyLeWL1THoW', 'job_seeker', 'Email Test User', NULL, NULL, 1, 0, NULL, '2026-01-23 03:25:21', '2026-01-23 03:25:21', NULL),
(20, 'jastiska14@gmail.com', '$2a$10$aIiXz6C49XHFpRSS9pILOeqIycow/uRlAaeLp6WtXFPF6P7PAgXui', 'job_seeker', 'Jastiska Dwi Wanda Sari', '08893011438', '/docs/avatars/avatar_20_1769167905.png', 1, 0, NULL, '2026-01-23 03:55:52', '2026-01-23 11:31:45', NULL),
(21, 'craftgirlsssshopping@gmail.com', '$2a$10$3NfsAYR/Sm5lhYYnXSCRC.ZG1DIaLC2sD3FBdtioJqqDrLYjEAFES', 'job_seeker', 'Saputra Budianto', '0881036480285', '/docs/avatars/avatar_21_1769148776.jpg', 1, 0, NULL, '2026-01-23 03:57:01', '2026-02-01 15:15:59', NULL),
(22, 'ahmad.pratama@email.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'partner', 'Ahmad Pratama', NULL, NULL, 1, 1, NULL, '2026-02-04 13:06:58', '2026-02-04 13:06:58', NULL),
(23, 'testpartner2@example.com', '$2a$10$dmbFVfwx/WXUQcc1ngOwDOB8IO4mW760dgNUE6g9O7uYl9W.MxoHm', 'partner', 'Test Partner Baru', '081234567890', NULL, 1, 0, NULL, '2026-02-05 12:26:15', '2026-02-05 12:27:43', NULL),
(24, 'localhosting8080@gmail.com', '$2a$10$AhHqSNhjP0q2bTzKwcbktO/6YgyNMc/VGiFobDF4nwCXsUmMI3hv2', 'partner', 'Saputra Budianto', '0881036480285', NULL, 1, 0, NULL, '2026-02-05 12:37:05', '2026-02-08 02:28:54', NULL),
(26, 'localhosting127.0.0.1@gmail.com', '$2a$10$fvhCmmsY6aZ0Wr0TekF/keE5CkztyZBc8xPAQmAOZdlK8yS.fDghq', 'company', 'Admin', '0881036480285', NULL, 1, 0, NULL, '2026-02-09 13:08:22', '2026-02-09 13:08:22', NULL);

-- --------------------------------------------------------

--
-- Stand-in structure for view `v_partner_dashboard_stats`
-- (See below for the actual view)
--
CREATE TABLE `v_partner_dashboard_stats` (
`partner_id` bigint(20) unsigned
,`user_id` bigint(20) unsigned
,`referral_code` varchar(20)
,`total_companies` int(11)
,`total_transactions` bigint(21)
,`total_commission` bigint(20)
,`available_balance` bigint(20)
,`paid_commission` bigint(20)
,`pending_commission` bigint(20)
);

-- --------------------------------------------------------

--
-- Stand-in structure for view `v_partner_monthly_stats`
-- (See below for the actual view)
--
CREATE TABLE `v_partner_monthly_stats` (
`partner_id` bigint(20) unsigned
,`month_year` varchar(7)
,`month_name` varchar(32)
,`total_commission` decimal(41,0)
,`companies_count` bigint(21)
);

-- --------------------------------------------------------

--
-- Structure for view `v_partner_dashboard_stats`
--
DROP TABLE IF EXISTS `v_partner_dashboard_stats`;

CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`localhost` SQL SECURITY DEFINER VIEW `v_partner_dashboard_stats`  AS SELECT `rp`.`id` AS `partner_id`, `rp`.`user_id` AS `user_id`, `rp`.`referral_code` AS `referral_code`, `rp`.`total_referrals` AS `total_companies`, count(distinct `pc`.`id`) AS `total_transactions`, coalesce(`rp`.`total_commission`,0) AS `total_commission`, coalesce(`rp`.`available_balance`,0) AS `available_balance`, coalesce(`rp`.`paid_amount`,0) AS `paid_commission`, coalesce(`rp`.`pending_balance`,0) AS `pending_commission` FROM (`referral_partners` `rp` left join `partner_commissions` `pc` on(`pc`.`partner_id` = `rp`.`id`)) GROUP BY `rp`.`id` ;

-- --------------------------------------------------------

--
-- Structure for view `v_partner_monthly_stats`
--
DROP TABLE IF EXISTS `v_partner_monthly_stats`;

CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`localhost` SQL SECURITY DEFINER VIEW `v_partner_monthly_stats`  AS SELECT `pc`.`partner_id` AS `partner_id`, date_format(`pc`.`created_at`,'%Y-%m') AS `month_year`, date_format(`pc`.`created_at`,'%b') AS `month_name`, sum(`pc`.`commission_amount`) AS `total_commission`, count(distinct `pr`.`company_id`) AS `companies_count` FROM (`partner_commissions` `pc` join `partner_referrals` `pr` on(`pr`.`id` = `pc`.`referral_id`)) WHERE `pc`.`status` in ('approved','paid') GROUP BY `pc`.`partner_id`, date_format(`pc`.`created_at`,'%Y-%m') ORDER BY date_format(`pc`.`created_at`,'%Y-%m') DESC ;

--
-- Indexes for dumped tables
--

--
-- Indexes for table `announcements`
--
ALTER TABLE `announcements`
  ADD PRIMARY KEY (`id`),
  ADD KEY `idx_type` (`type`),
  ADD KEY `idx_target_audience` (`target_audience`),
  ADD KEY `idx_is_active` (`is_active`),
  ADD KEY `idx_created_at` (`created_at`),
  ADD KEY `idx_start_end_date` (`start_date`,`end_date`);

--
-- Indexes for table `applicant_documents`
--
ALTER TABLE `applicant_documents`
  ADD PRIMARY KEY (`id`),
  ADD KEY `idx_user_id` (`user_id`),
  ADD KEY `idx_document_type` (`document_type`),
  ADD KEY `idx_is_primary` (`is_primary`);

--
-- Indexes for table `applicant_profiles`
--
ALTER TABLE `applicant_profiles`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `uk_user_id` (`user_id`),
  ADD KEY `idx_city` (`city`),
  ADD KEY `idx_province` (`province`),
  ADD KEY `idx_profile_completeness` (`profile_completeness`),
  ADD KEY `idx_location` (`city`,`province`),
  ADD KEY `idx_salary_range` (`expected_salary_min`,`expected_salary_max`);

--
-- Indexes for table `applications`
--
ALTER TABLE `applications`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `uk_user_job` (`user_id`,`job_id`),
  ADD KEY `fk_applications_cv_snapshot` (`cv_snapshot_id`),
  ADD KEY `idx_applications_user_id` (`user_id`),
  ADD KEY `idx_applications_job_id` (`job_id`),
  ADD KEY `idx_applications_status` (`current_status`),
  ADD KEY `idx_applications_applied_at` (`applied_at`);

--
-- Indexes for table `application_timelines`
--
ALTER TABLE `application_timelines`
  ADD PRIMARY KEY (`id`),
  ADD KEY `idx_timelines_application_id` (`application_id`),
  ADD KEY `idx_timelines_status` (`status`),
  ADD KEY `idx_timelines_created_at` (`created_at`);

--
-- Indexes for table `audit_logs`
--
ALTER TABLE `audit_logs`
  ADD PRIMARY KEY (`id`),
  ADD KEY `idx_audit_logs_user_id` (`user_id`),
  ADD KEY `idx_audit_logs_entity` (`entity_type`,`entity_id`),
  ADD KEY `idx_audit_logs_action` (`action`),
  ADD KEY `idx_audit_logs_created_at` (`created_at`);

--
-- Indexes for table `chat_messages`
--
ALTER TABLE `chat_messages`
  ADD PRIMARY KEY (`id`),
  ADD KEY `idx_conversation_id` (`conversation_id`),
  ADD KEY `idx_sender_id` (`sender_id`),
  ADD KEY `idx_created_at` (`created_at`),
  ADD KEY `idx_is_read` (`is_read`),
  ADD KEY `idx_attachment_type` (`attachment_type`);

--
-- Indexes for table `companies`
--
ALTER TABLE `companies`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `user_id` (`user_id`),
  ADD KEY `fk_companies_verified_by` (`documents_verified_by`),
  ADD KEY `idx_company_status` (`company_status`),
  ADD KEY `idx_deleted_at` (`deleted_at`),
  ADD KEY `idx_created_at` (`created_at`),
  ADD KEY `idx_referred_by_partner` (`referred_by_partner_id`),
  ADD KEY `idx_companies_referral_code` (`referral_code_used`);

--
-- Indexes for table `company_quotas`
--
ALTER TABLE `company_quotas`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `company_id` (`company_id`),
  ADD KEY `idx_company_quotas_company_id` (`company_id`);

--
-- Indexes for table `conversations`
--
ALTER TABLE `conversations`
  ADD PRIMARY KEY (`id`),
  ADD KEY `idx_company_id` (`company_id`),
  ADD KEY `idx_status` (`status`),
  ADD KEY `idx_category` (`category`),
  ADD KEY `idx_created_at` (`created_at`);

--
-- Indexes for table `cvs`
--
ALTER TABLE `cvs`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `user_id` (`user_id`),
  ADD KEY `idx_cvs_user_id` (`user_id`),
  ADD KEY `idx_cvs_completeness` (`completeness_score`);

--
-- Indexes for table `cv_snapshots`
--
ALTER TABLE `cv_snapshots`
  ADD PRIMARY KEY (`id`),
  ADD KEY `idx_cv_snapshots_cv_id` (`cv_id`),
  ADD KEY `idx_cv_snapshots_user_id` (`user_id`),
  ADD KEY `idx_cv_snapshots_hash` (`snapshot_hash`);

--
-- Indexes for table `jobs`
--
ALTER TABLE `jobs`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `slug` (`slug`),
  ADD KEY `idx_jobs_company_id` (`company_id`),
  ADD KEY `idx_jobs_status` (`status`),
  ADD KEY `idx_jobs_job_type` (`job_type`),
  ADD KEY `idx_jobs_city` (`city`),
  ADD KEY `idx_jobs_province` (`province`),
  ADD KEY `idx_jobs_salary` (`salary_min`,`salary_max`),
  ADD KEY `idx_jobs_published_at` (`published_at`),
  ADD KEY `idx_jobs_deleted_at` (`deleted_at`),
  ADD KEY `idx_jobs_admin_status` (`admin_status`),
  ADD KEY `idx_jobs_status_admin` (`status`,`admin_status`);
ALTER TABLE `jobs` ADD FULLTEXT KEY `idx_jobs_search` (`title`,`description`,`requirements`);

--
-- Indexes for table `job_shares`
--
ALTER TABLE `job_shares`
  ADD PRIMARY KEY (`id`),
  ADD KEY `idx_job_shares_job_id` (`job_id`),
  ADD KEY `idx_job_shares_user_id` (`user_id`);

--
-- Indexes for table `job_skills`
--
ALTER TABLE `job_skills`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `uk_job_skill` (`job_id`,`skill_name`),
  ADD KEY `idx_job_skills_skill_name` (`skill_name`);

--
-- Indexes for table `job_views`
--
ALTER TABLE `job_views`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `unique_job_user_view` (`job_id`,`user_id`),
  ADD KEY `idx_job_views_job_id` (`job_id`),
  ADD KEY `idx_job_views_user_id` (`user_id`);

--
-- Indexes for table `notifications`
--
ALTER TABLE `notifications`
  ADD PRIMARY KEY (`id`),
  ADD KEY `idx_notifications_user_id` (`user_id`),
  ADD KEY `idx_notifications_is_read` (`is_read`),
  ADD KEY `idx_notifications_created_at` (`created_at`);

--
-- Indexes for table `partner_commissions`
--
ALTER TABLE `partner_commissions`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `uk_payment_id` (`payment_id`) COMMENT 'One commission record per payment',
  ADD KEY `idx_partner_id` (`partner_id`),
  ADD KEY `idx_referral_id` (`referral_id`),
  ADD KEY `idx_company_id` (`company_id`),
  ADD KEY `idx_status` (`status`),
  ADD KEY `idx_created_at` (`created_at`),
  ADD KEY `idx_payout_id` (`payout_id`),
  ADD KEY `fk_commission_approved_by` (`approved_by`),
  ADD KEY `idx_commissions_date_status` (`created_at`,`status`);

--
-- Indexes for table `partner_payouts`
--
ALTER TABLE `partner_payouts`
  ADD PRIMARY KEY (`id`),
  ADD KEY `idx_partner_id` (`partner_id`),
  ADD KEY `idx_status` (`status`),
  ADD KEY `idx_created_at` (`created_at`),
  ADD KEY `idx_processed_at` (`processed_at`),
  ADD KEY `fk_payout_processed_by` (`processed_by`),
  ADD KEY `idx_payouts_partner_date` (`partner_id`,`created_at`);

--
-- Indexes for table `partner_referrals`
--
ALTER TABLE `partner_referrals`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `uk_partner_company` (`partner_id`,`company_id`),
  ADD UNIQUE KEY `uk_company_id` (`company_id`) COMMENT 'A company can only have one referrer',
  ADD KEY `idx_partner_id` (`partner_id`),
  ADD KEY `idx_registered_at` (`registered_at`);

--
-- Indexes for table `password_resets`
--
ALTER TABLE `password_resets`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `token` (`token`),
  ADD KEY `idx_token` (`token`),
  ADD KEY `idx_user_id` (`user_id`),
  ADD KEY `idx_expires_at` (`expires_at`);

--
-- Indexes for table `password_reset_tokens`
--
ALTER TABLE `password_reset_tokens`
  ADD PRIMARY KEY (`id`),
  ADD KEY `idx_email` (`email`),
  ADD KEY `idx_token` (`token`),
  ADD KEY `idx_expires_at` (`expires_at`);

--
-- Indexes for table `payments`
--
ALTER TABLE `payments`
  ADD PRIMARY KEY (`id`),
  ADD KEY `job_id` (`job_id`),
  ADD KEY `confirmed_by_id` (`confirmed_by_id`),
  ADD KEY `idx_payments_company_id` (`company_id`),
  ADD KEY `idx_payments_status` (`status`),
  ADD KEY `idx_payments_submitted_at` (`submitted_at`),
  ADD KEY `idx_payments_package_id` (`package_id`);

--
-- Indexes for table `referral_partners`
--
ALTER TABLE `referral_partners`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `uk_user_id` (`user_id`),
  ADD UNIQUE KEY `uk_referral_code` (`referral_code`),
  ADD KEY `idx_status` (`status`),
  ADD KEY `idx_referral_code` (`referral_code`),
  ADD KEY `idx_created_at` (`created_at`),
  ADD KEY `fk_partner_approved_by` (`approved_by`);

--
-- Indexes for table `refresh_tokens`
--
ALTER TABLE `refresh_tokens`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `token_hash` (`token_hash`),
  ADD KEY `idx_refresh_tokens_user_id` (`user_id`),
  ADD KEY `idx_refresh_tokens_expires_at` (`expires_at`);

--
-- Indexes for table `saved_jobs`
--
ALTER TABLE `saved_jobs`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `uk_saved_job` (`user_id`,`job_id`),
  ADD KEY `fk_saved_jobs_job` (`job_id`),
  ADD KEY `idx_saved_jobs_user_id` (`user_id`);

--
-- Indexes for table `support_tickets`
--
ALTER TABLE `support_tickets`
  ADD PRIMARY KEY (`id`),
  ADD KEY `idx_user_id` (`user_id`),
  ADD KEY `idx_status` (`status`),
  ADD KEY `idx_priority` (`priority`),
  ADD KEY `idx_category` (`category`),
  ADD KEY `idx_created_at` (`created_at`);

--
-- Indexes for table `ticket_responses`
--
ALTER TABLE `ticket_responses`
  ADD PRIMARY KEY (`id`),
  ADD KEY `idx_ticket_id` (`ticket_id`),
  ADD KEY `idx_sender_id` (`sender_id`),
  ADD KEY `idx_created_at` (`created_at`);

--
-- Indexes for table `users`
--
ALTER TABLE `users`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `email` (`email`),
  ADD KEY `idx_users_email` (`email`),
  ADD KEY `idx_users_role` (`role`),
  ADD KEY `idx_users_is_active` (`is_active`),
  ADD KEY `idx_users_deleted_at` (`deleted_at`),
  ADD KEY `idx_users_role_status` (`role`,`is_active`);

--
-- AUTO_INCREMENT for dumped tables
--

--
-- AUTO_INCREMENT for table `announcements`
--
ALTER TABLE `announcements`
  MODIFY `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=20;

--
-- AUTO_INCREMENT for table `applicant_documents`
--
ALTER TABLE `applicant_documents`
  MODIFY `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=2;

--
-- AUTO_INCREMENT for table `applicant_profiles`
--
ALTER TABLE `applicant_profiles`
  MODIFY `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=3;

--
-- AUTO_INCREMENT for table `applications`
--
ALTER TABLE `applications`
  MODIFY `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=13;

--
-- AUTO_INCREMENT for table `application_timelines`
--
ALTER TABLE `application_timelines`
  MODIFY `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=28;

--
-- AUTO_INCREMENT for table `audit_logs`
--
ALTER TABLE `audit_logs`
  MODIFY `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `chat_messages`
--
ALTER TABLE `chat_messages`
  MODIFY `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=17;

--
-- AUTO_INCREMENT for table `companies`
--
ALTER TABLE `companies`
  MODIFY `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=8;

--
-- AUTO_INCREMENT for table `company_quotas`
--
ALTER TABLE `company_quotas`
  MODIFY `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=9;

--
-- AUTO_INCREMENT for table `conversations`
--
ALTER TABLE `conversations`
  MODIFY `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=7;

--
-- AUTO_INCREMENT for table `cvs`
--
ALTER TABLE `cvs`
  MODIFY `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=5;

--
-- AUTO_INCREMENT for table `cv_snapshots`
--
ALTER TABLE `cv_snapshots`
  MODIFY `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=13;

--
-- AUTO_INCREMENT for table `jobs`
--
ALTER TABLE `jobs`
  MODIFY `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=58;

--
-- AUTO_INCREMENT for table `job_shares`
--
ALTER TABLE `job_shares`
  MODIFY `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=5;

--
-- AUTO_INCREMENT for table `job_skills`
--
ALTER TABLE `job_skills`
  MODIFY `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=86;

--
-- AUTO_INCREMENT for table `job_views`
--
ALTER TABLE `job_views`
  MODIFY `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=23;

--
-- AUTO_INCREMENT for table `notifications`
--
ALTER TABLE `notifications`
  MODIFY `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `partner_commissions`
--
ALTER TABLE `partner_commissions`
  MODIFY `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `partner_payouts`
--
ALTER TABLE `partner_payouts`
  MODIFY `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=2;

--
-- AUTO_INCREMENT for table `partner_referrals`
--
ALTER TABLE `partner_referrals`
  MODIFY `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=3;

--
-- AUTO_INCREMENT for table `password_resets`
--
ALTER TABLE `password_resets`
  MODIFY `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=3;

--
-- AUTO_INCREMENT for table `password_reset_tokens`
--
ALTER TABLE `password_reset_tokens`
  MODIFY `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=9;

--
-- AUTO_INCREMENT for table `payments`
--
ALTER TABLE `payments`
  MODIFY `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=16;

--
-- AUTO_INCREMENT for table `referral_partners`
--
ALTER TABLE `referral_partners`
  MODIFY `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=4;

--
-- AUTO_INCREMENT for table `refresh_tokens`
--
ALTER TABLE `refresh_tokens`
  MODIFY `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=408;

--
-- AUTO_INCREMENT for table `saved_jobs`
--
ALTER TABLE `saved_jobs`
  MODIFY `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=3;

--
-- AUTO_INCREMENT for table `support_tickets`
--
ALTER TABLE `support_tickets`
  MODIFY `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=4;

--
-- AUTO_INCREMENT for table `ticket_responses`
--
ALTER TABLE `ticket_responses`
  MODIFY `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=6;

--
-- AUTO_INCREMENT for table `users`
--
ALTER TABLE `users`
  MODIFY `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=27;

--
-- Constraints for dumped tables
--

--
-- Constraints for table `applicant_documents`
--
ALTER TABLE `applicant_documents`
  ADD CONSTRAINT `fk_applicant_documents_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE;

--
-- Constraints for table `applicant_profiles`
--
ALTER TABLE `applicant_profiles`
  ADD CONSTRAINT `fk_applicant_profiles_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE;

--
-- Constraints for table `applications`
--
ALTER TABLE `applications`
  ADD CONSTRAINT `fk_applications_cv_snapshot` FOREIGN KEY (`cv_snapshot_id`) REFERENCES `cv_snapshots` (`id`),
  ADD CONSTRAINT `fk_applications_job` FOREIGN KEY (`job_id`) REFERENCES `jobs` (`id`) ON DELETE CASCADE,
  ADD CONSTRAINT `fk_applications_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE;

--
-- Constraints for table `application_timelines`
--
ALTER TABLE `application_timelines`
  ADD CONSTRAINT `fk_timelines_application` FOREIGN KEY (`application_id`) REFERENCES `applications` (`id`) ON DELETE CASCADE;

--
-- Constraints for table `chat_messages`
--
ALTER TABLE `chat_messages`
  ADD CONSTRAINT `chat_messages_ibfk_1` FOREIGN KEY (`conversation_id`) REFERENCES `conversations` (`id`) ON DELETE CASCADE,
  ADD CONSTRAINT `chat_messages_ibfk_2` FOREIGN KEY (`sender_id`) REFERENCES `users` (`id`) ON DELETE CASCADE;

--
-- Constraints for table `companies`
--
ALTER TABLE `companies`
  ADD CONSTRAINT `fk_companies_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE,
  ADD CONSTRAINT `fk_companies_verified_by` FOREIGN KEY (`documents_verified_by`) REFERENCES `users` (`id`) ON DELETE SET NULL,
  ADD CONSTRAINT `fk_company_referrer` FOREIGN KEY (`referred_by_partner_id`) REFERENCES `referral_partners` (`id`) ON DELETE SET NULL;

--
-- Constraints for table `company_quotas`
--
ALTER TABLE `company_quotas`
  ADD CONSTRAINT `company_quotas_ibfk_1` FOREIGN KEY (`company_id`) REFERENCES `users` (`id`) ON DELETE CASCADE;

--
-- Constraints for table `conversations`
--
ALTER TABLE `conversations`
  ADD CONSTRAINT `conversations_ibfk_1` FOREIGN KEY (`company_id`) REFERENCES `users` (`id`) ON DELETE CASCADE;

--
-- Constraints for table `cvs`
--
ALTER TABLE `cvs`
  ADD CONSTRAINT `fk_cvs_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE;

--
-- Constraints for table `cv_snapshots`
--
ALTER TABLE `cv_snapshots`
  ADD CONSTRAINT `fk_cv_snapshots_cv` FOREIGN KEY (`cv_id`) REFERENCES `cvs` (`id`) ON DELETE CASCADE,
  ADD CONSTRAINT `fk_cv_snapshots_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE;

--
-- Constraints for table `jobs`
--
ALTER TABLE `jobs`
  ADD CONSTRAINT `fk_jobs_company` FOREIGN KEY (`company_id`) REFERENCES `users` (`id`) ON DELETE CASCADE;

--
-- Constraints for table `job_shares`
--
ALTER TABLE `job_shares`
  ADD CONSTRAINT `fk_job_shares_job` FOREIGN KEY (`job_id`) REFERENCES `jobs` (`id`) ON DELETE CASCADE,
  ADD CONSTRAINT `fk_job_shares_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE SET NULL;

--
-- Constraints for table `job_skills`
--
ALTER TABLE `job_skills`
  ADD CONSTRAINT `fk_job_skills_job` FOREIGN KEY (`job_id`) REFERENCES `jobs` (`id`) ON DELETE CASCADE;

--
-- Constraints for table `job_views`
--
ALTER TABLE `job_views`
  ADD CONSTRAINT `fk_job_views_job` FOREIGN KEY (`job_id`) REFERENCES `jobs` (`id`) ON DELETE CASCADE,
  ADD CONSTRAINT `fk_job_views_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE;

--
-- Constraints for table `notifications`
--
ALTER TABLE `notifications`
  ADD CONSTRAINT `fk_notifications_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE;

--
-- Constraints for table `partner_commissions`
--
ALTER TABLE `partner_commissions`
  ADD CONSTRAINT `fk_commission_approved_by` FOREIGN KEY (`approved_by`) REFERENCES `users` (`id`) ON DELETE SET NULL,
  ADD CONSTRAINT `fk_commission_company` FOREIGN KEY (`company_id`) REFERENCES `companies` (`id`) ON DELETE CASCADE,
  ADD CONSTRAINT `fk_commission_partner` FOREIGN KEY (`partner_id`) REFERENCES `referral_partners` (`id`) ON DELETE CASCADE,
  ADD CONSTRAINT `fk_commission_payment` FOREIGN KEY (`payment_id`) REFERENCES `payments` (`id`) ON DELETE CASCADE,
  ADD CONSTRAINT `fk_commission_referral` FOREIGN KEY (`referral_id`) REFERENCES `partner_referrals` (`id`) ON DELETE CASCADE;

--
-- Constraints for table `partner_payouts`
--
ALTER TABLE `partner_payouts`
  ADD CONSTRAINT `fk_payout_partner` FOREIGN KEY (`partner_id`) REFERENCES `referral_partners` (`id`) ON DELETE CASCADE,
  ADD CONSTRAINT `fk_payout_processed_by` FOREIGN KEY (`processed_by`) REFERENCES `users` (`id`) ON DELETE SET NULL;

--
-- Constraints for table `partner_referrals`
--
ALTER TABLE `partner_referrals`
  ADD CONSTRAINT `fk_referral_company` FOREIGN KEY (`company_id`) REFERENCES `companies` (`id`) ON DELETE CASCADE,
  ADD CONSTRAINT `fk_referral_partner` FOREIGN KEY (`partner_id`) REFERENCES `referral_partners` (`id`) ON DELETE CASCADE;

--
-- Constraints for table `payments`
--
ALTER TABLE `payments`
  ADD CONSTRAINT `payments_ibfk_1` FOREIGN KEY (`company_id`) REFERENCES `users` (`id`) ON DELETE CASCADE,
  ADD CONSTRAINT `payments_ibfk_2` FOREIGN KEY (`job_id`) REFERENCES `jobs` (`id`) ON DELETE SET NULL,
  ADD CONSTRAINT `payments_ibfk_3` FOREIGN KEY (`confirmed_by_id`) REFERENCES `users` (`id`) ON DELETE SET NULL;

--
-- Constraints for table `referral_partners`
--
ALTER TABLE `referral_partners`
  ADD CONSTRAINT `fk_partner_approved_by` FOREIGN KEY (`approved_by`) REFERENCES `users` (`id`) ON DELETE SET NULL,
  ADD CONSTRAINT `fk_partner_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE;

--
-- Constraints for table `refresh_tokens`
--
ALTER TABLE `refresh_tokens`
  ADD CONSTRAINT `fk_refresh_tokens_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE;

--
-- Constraints for table `saved_jobs`
--
ALTER TABLE `saved_jobs`
  ADD CONSTRAINT `fk_saved_jobs_job` FOREIGN KEY (`job_id`) REFERENCES `jobs` (`id`) ON DELETE CASCADE,
  ADD CONSTRAINT `fk_saved_jobs_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE;

--
-- Constraints for table `support_tickets`
--
ALTER TABLE `support_tickets`
  ADD CONSTRAINT `fk_support_tickets_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE;

--
-- Constraints for table `ticket_responses`
--
ALTER TABLE `ticket_responses`
  ADD CONSTRAINT `fk_ticket_responses_ticket` FOREIGN KEY (`ticket_id`) REFERENCES `support_tickets` (`id`) ON DELETE CASCADE;
COMMIT;

/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
