-- phpMyAdmin SQL Dump
-- version 5.2.1
-- https://www.phpmyadmin.net/
--
-- Host: localhost
-- Generation Time: Jan 19, 2026 at 02:41 PM
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
(1, 3, 1, 1, 'Dengan pengalaman 4 tahun sebagai Backend Developer, saya yakin dapat memberikan kontribusi signifikan di PT TechCorp Indonesia.', 'hired', '2026-01-17 02:51:28', '2026-01-17 02:55:46', '2026-01-17 02:51:28', '2026-01-17 02:55:46');

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
  `created_at` timestamp NOT NULL DEFAULT current_timestamp()
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

--
-- Dumping data for table `application_timelines`
--

INSERT INTO `application_timelines` (`id`, `application_id`, `status`, `note`, `is_visible_to_applicant`, `updated_by_type`, `updated_by_id`, `scheduled_at`, `scheduled_location`, `scheduled_notes`, `created_at`) VALUES
(1, 1, 'submitted', 'Lamaran berhasil dikirim', 1, 'system', NULL, NULL, NULL, NULL, '2026-01-17 02:51:28'),
(2, 1, 'viewed', 'Melihat profil kandidat', 1, 'company', 2, NULL, NULL, NULL, '2026-01-17 02:53:22'),
(3, 1, 'shortlisted', 'Kandidat memenuhi kriteria', 1, 'company', 2, NULL, NULL, NULL, '2026-01-17 02:53:40'),
(4, 1, 'interview_scheduled', 'Interview tahap 1', 1, 'company', 2, '2026-01-20 03:00:00', 'Kantor Jakarta', NULL, '2026-01-17 02:53:40'),
(5, 1, 'interview_completed', 'Interview berhasil', 1, 'company', 2, NULL, NULL, NULL, '2026-01-17 02:54:09'),
(6, 1, 'offer_sent', 'Surat penawaran dikirim', 1, 'company', 2, NULL, NULL, NULL, '2026-01-17 02:54:58'),
(7, 1, 'offer_accepted', 'Kandidat menerima penawaran', 1, 'company', 2, NULL, NULL, NULL, '2026-01-17 02:55:35'),
(8, 1, 'hired', 'Selamat bergabung di PT TechCorp Indonesia!', 1, 'company', 2, NULL, NULL, NULL, '2026-01-17 02:55:46');

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
  `deleted_at` timestamp NULL DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

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
(2, 7, 0, 0, '2026-01-18 09:14:40', '2026-01-18 09:14:40');

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
(2, 3, '{\"full_name\":\"Budi Santoso\",\"email\":\"budi@gmail.com\",\"phone\":\"+6281234567890\"}', 'null', '[{\"company\":\"PT Software House\",\"position\":\"Backend Developer\",\"start_date\":\"2019-08-01\",\"is_current\":true,\"description\":\"Developing REST APIs\"}]', '[{\"name\":\"Go\",\"level\":\"advanced\"}]', 'null', 'null', 'null', '2026-01-17 03:12:22', 45, '2026-01-17 02:51:16', '2026-01-17 03:12:22');

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
(1, 2, 3, '{\"full_name\":\"Budi Santoso\",\"email\":\"budi.kandidat@gmail.com\",\"phone\":\"+6281234567890\"}', '[{\"institution\":\"UI\",\"degree\":\"S1\",\"field_of_study\":\"Informatika\",\"start_date\":\"2015-08-01\",\"end_date\":\"2019-07-01\"}]', 'null', '[{\"name\":\"Go\",\"level\":\"advanced\"}]', 'null', 'null', 'null', 'ff53fa4b736225a2f2cecb569e9455e38fb7a4de93409f39dcf6b96c7d489c74', 50, '2026-01-17 02:51:28');

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
  `published_at` timestamp NULL DEFAULT NULL,
  `created_at` timestamp NOT NULL DEFAULT current_timestamp(),
  `updated_at` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  `deleted_at` timestamp NULL DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

--
-- Dumping data for table `jobs`
--

INSERT INTO `jobs` (`id`, `company_id`, `title`, `slug`, `description`, `requirements`, `responsibilities`, `benefits`, `city`, `province`, `is_remote`, `job_type`, `experience_level`, `salary_min`, `salary_max`, `salary_currency`, `is_salary_visible`, `application_deadline`, `max_applications`, `status`, `admin_status`, `admin_note`, `flag_reason`, `views_count`, `applications_count`, `published_at`, `created_at`, `updated_at`, `deleted_at`) VALUES
(1, 2, 'Senior Software Engineer', 'senior-software-engineer', 'Kami mencari Senior Software Engineer untuk bergabung dengan tim development kami. Anda akan bekerja dengan teknologi terkini dan tim yang solid. Minimal pengalaman 3 tahun.', 'Minimal 3 tahun pengalaman. Menguasai Go/Python.', NULL, NULL, 'Jakarta Selatan', 'DKI Jakarta', 1, 'full_time', 'senior', 15000000, 25000000, 'IDR', 0, NULL, NULL, 'active', NULL, NULL, NULL, 0, 0, '2026-01-17 02:49:00', '2026-01-17 02:46:11', '2026-01-17 02:49:00', NULL),
(2, 2, 'Test Status Job', 'test-status-job', 'Ini adalah deskripsi panjang untuk testing job status management endpoints yang baru ditambahkan.', NULL, NULL, NULL, 'Jakarta', 'DKI Jakarta', 0, 'full_time', 'mid', NULL, NULL, 'IDR', 0, NULL, NULL, 'closed', NULL, NULL, NULL, 0, 0, '2026-01-17 03:13:52', '2026-01-17 03:13:25', '2026-01-17 03:14:15', NULL),
(3, 4, 'Senior Software Engineer', 'senior-software-engineer-1768720323', 'We are looking for an experienced software engineer', '5+ years experience with Go or Python', 'Build and maintain backend systems', 'Competitive salary, remote work', 'Jakarta Selatan', 'DKI Jakarta', 1, 'full_time', 'senior', 20000000, 35000000, 'IDR', 1, NULL, NULL, 'active', NULL, NULL, NULL, 0, 0, '2026-01-18 07:12:10', '2026-01-18 07:12:03', '2026-01-18 07:12:19', NULL),
(4, 10, 'Senior Backend Engineer', 'senior-backend-engineer', 'Kami mencari Senior Backend Engineer yang berpengalaman dalam mengembangkan sistem scalable', '- Minimal 5 tahun pengalaman backend development\n- Mahir Go, Python, atau Java\n- Pengalaman dengan microservices\n- Pengalaman dengan database SQL dan NoSQL', '- Merancang dan mengimplementasi API\n- Melakukan code review\n- Mengoptimalkan performance sistem\n- Mentoring junior developers', '- Gaji kompetitif 15-25 juta/bulan\n- Asuransi kesehatan\n- Work from home flexibility\n- Training budget', 'Jakarta Selatan', 'DKI Jakarta', 1, 'full_time', 'senior', 15000000, 25000000, 'IDR', 1, NULL, NULL, 'active', NULL, NULL, NULL, 0, 0, '2026-01-19 07:48:04', '2026-01-19 07:48:03', '2026-01-19 07:48:04', NULL),
(5, 10, 'Full Stack Developer', 'full-stack-developer', 'Bergabunglah dengan tim kami sebagai Full Stack Developer', '- Minimal 3 tahun pengalaman\n- React atau Vue.js\n- Node.js atau Python', '- Develop frontend dan backend\n- Collaborate dengan tim design', '- Gaji 8-12 juta/bulan\n- Remote friendly', 'Jakarta Pusat', 'DKI Jakarta', 1, 'full_time', 'mid', 8000000, 12000000, 'IDR', 1, NULL, NULL, 'draft', NULL, NULL, NULL, 0, 0, NULL, '2026-01-19 07:48:03', '2026-01-19 07:48:03', NULL),
(6, 10, 'UI/UX Designer', 'uiux-designer', 'Kami mencari UI/UX Designer untuk mengembangkan product kami', '- 2+ tahun pengalaman UI/UX\n- Figma atau Adobe XD', '- Design interface\n- User research', '- Gaji 6-10 juta/bulan', 'Jakarta Selatan', 'DKI Jakarta', 1, 'full_time', 'junior', 6000000, 10000000, 'IDR', 1, NULL, NULL, 'draft', NULL, NULL, NULL, 0, 0, NULL, '2026-01-19 07:48:03', '2026-01-19 07:48:03', NULL);

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
(19, 6, 'Prototyping', 1);

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
-- Table structure for table `payments`
--

CREATE TABLE `payments` (
  `id` bigint(20) UNSIGNED NOT NULL,
  `company_id` bigint(20) UNSIGNED NOT NULL,
  `job_id` bigint(20) UNSIGNED DEFAULT NULL,
  `amount` bigint(20) NOT NULL DEFAULT 30000,
  `proof_image_url` varchar(500) DEFAULT NULL,
  `status` enum('pending','confirmed','rejected') NOT NULL DEFAULT 'pending',
  `note` text DEFAULT NULL,
  `confirmed_by_id` bigint(20) UNSIGNED DEFAULT NULL,
  `submitted_at` timestamp NOT NULL DEFAULT current_timestamp(),
  `confirmed_at` timestamp NULL DEFAULT NULL,
  `created_at` timestamp NOT NULL DEFAULT current_timestamp(),
  `updated_at` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp()
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

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
(32, 7, 'a9cf5c2eed090d2b482531c18a2bd20ec17094b080dd0418c59998d942612f3d', '2026-01-25 09:13:55', NULL, '', '', '2026-01-18 09:13:55'),
(33, 6, '2a001d6285eab9929dd40ed49228792412f4876a85357a5536d6d29dc0dcf23c', '2026-01-25 09:24:26', NULL, '', '', '2026-01-18 09:24:26'),
(34, 6, '73a916bb49ac98ca2209557a1f77ecf37d2e561f70a7c7c581f30e720b9abad8', '2026-01-25 09:24:43', NULL, '', '', '2026-01-18 09:24:43'),
(35, 7, '492b7ac0203728ed653ae262a5ed3bdf0136900d3bb810eff4313aa065367a7c', '2026-01-25 09:30:55', NULL, '', '', '2026-01-18 09:30:55'),
(36, 7, '3a4494825b610104798cabdb98f802e5ca3deea297d6b1283da67ba6dd908770', '2026-01-26 03:19:22', NULL, '', '', '2026-01-19 03:19:22'),
(38, 9, 'b84311ea9ecaf56a971324bc225b1282bc4ebcfd23cad0e6d4096d6330efb428', '2026-01-26 07:44:42', NULL, '', '', '2026-01-19 07:44:42'),
(39, 9, 'c422e7b8ce81e16955e97ccbee833d5363a46c8a7833cf90e143e59424c5a1ad', '2026-01-26 07:44:42', NULL, '', '', '2026-01-19 07:44:42'),
(40, 10, '6ba39df159ff2f5ff5ee8a882613018387b0b4c7811d0f7b1e54787cd613a7eb', '2026-01-26 07:48:03', NULL, '', '', '2026-01-19 07:48:03'),
(41, 7, '165952470c425c53003d70e42ebe4064657097c48b5508162c53d44d7427a539', '2026-01-26 08:45:03', NULL, '', '', '2026-01-19 08:45:03'),
(42, 7, '6d1615471fb01dfe916eac9018717fa29938e503fd6839733bfed8ef81d98b24', '2026-01-26 09:08:29', NULL, '', '', '2026-01-19 09:08:29'),
(43, 7, '802c7b7e85b969ac06cccb49f5b8fd49f395d59664b85405b9ac8b738c4254b8', '2026-01-26 09:17:25', NULL, '', '', '2026-01-19 09:17:25'),
(44, 7, '6a099673f30809c847a41958a8112fde7b6f111161775cdca8cc8fee3520add4', '2026-01-26 09:23:58', NULL, '', '', '2026-01-19 09:23:58'),
(45, 7, '51b304ec574e7e09973cbc26b3ff6e614a2d736e2f87c9f425c1468fa21687b4', '2026-01-26 09:26:52', NULL, '', '', '2026-01-19 09:26:52'),
(46, 7, '31b5c12979c17e775b5b1ac4e4de03438aaec394f67be38a858c97b1141635b2', '2026-01-26 09:28:07', NULL, '', '', '2026-01-19 09:28:07'),
(47, 7, 'a7549509411647c8a69f56e7defee2351a08feb15cc6f275925629db077b50c4', '2026-01-26 09:28:22', NULL, '', '', '2026-01-19 09:28:22'),
(48, 7, 'ca5466a3c7888c45fc1278366f852ba5facee273702e2c51e44db24dba6b9293', '2026-01-26 09:28:32', NULL, '', '', '2026-01-19 09:28:32'),
(49, 7, '9222dd7ef70b9703e9c7a61b99de032a6425d6e445a40d0d2634d93c8e163a45', '2026-01-26 09:29:51', NULL, '', '', '2026-01-19 09:29:51'),
(50, 7, '030620c5e4c6e29a9c91edf548e1a8da5bb9f59200d2cbf983b2eca4aef2692b', '2026-01-26 09:32:44', NULL, '', '', '2026-01-19 09:32:44'),
(51, 7, 'a6e1868c7414376fb9fc23b911ed8e8b5aa948820bc39ed4a3389a03e418f599', '2026-01-26 09:33:20', NULL, '', '', '2026-01-19 09:33:20'),
(52, 7, '8313fde1d1de58c6b6ec84effd5e547e93e6e65a8536e398b953caf4e67d17e7', '2026-01-26 09:50:11', NULL, '', '', '2026-01-19 09:50:11'),
(53, 7, '60b94d83767aeedf47e9238e9696c3fb36c18df8222ac0cf2118ca104243643d', '2026-01-26 09:52:54', NULL, '', '', '2026-01-19 09:52:54'),
(54, 7, '059f7d2ef6d19ba619a43d128429e30e4296652a974422beafe5f0b68cc0d7c0', '2026-01-26 11:21:19', NULL, '', '', '2026-01-19 11:21:19'),
(55, 7, '2a950819245b52d5645dbbe90d8b406bad4d48cab508e517ac790512b9df0f5f', '2026-01-26 11:53:25', NULL, '', '', '2026-01-19 11:53:25'),
(56, 7, '8725208f1c198c038722c83254c7c84936a3ad018f244a6740b067fdd0a8e3fa', '2026-01-26 12:44:17', NULL, '', '', '2026-01-19 12:44:17'),
(57, 7, '6ce5ce65b4ee76f77199c0cbb565f7bd4a05e2af01f127f7e6b18c259bd58edc', '2026-01-26 13:00:57', NULL, '', '', '2026-01-19 13:00:57'),
(58, 7, 'd6821da179e51d9755f3a5fd45e188dccb59058e838a1967682b4e4f5af9dd9e', '2026-01-26 13:07:54', NULL, '', '', '2026-01-19 13:07:54'),
(59, 7, '2b87fa6d775286438439d2b8e0c72ec2a15c7f61cef9a406c1bdceeefa22bad0', '2026-01-26 13:20:06', NULL, '', '', '2026-01-19 13:20:06');

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

-- --------------------------------------------------------

--
-- Table structure for table `users`
--

CREATE TABLE `users` (
  `id` bigint(20) UNSIGNED NOT NULL,
  `email` varchar(255) NOT NULL,
  `password_hash` varchar(255) NOT NULL,
  `role` enum('job_seeker','company','admin') NOT NULL DEFAULT 'job_seeker',
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
(7, 'info@karyadeveloperindonesia.com', '$2a$10$yGx97vEQaIJOEeuDHlIGQONjqgqfdFhWUbBrWWxsJNpmeEXf5reWW', 'company', 'Admin', '0881036480285', NULL, 1, 0, NULL, '2026-01-18 09:13:55', '2026-01-19 13:21:21', NULL),
(9, 'company.testing1768808682@karirnusantara.com', '$2a$10$1RpgFQaG58KBk.Fhy8laye2yhBj.C.ogGGaqxkktzmew7TMrlLu.W', 'company', 'Budi Santoso', '081234567890', NULL, 1, 0, NULL, '2026-01-19 07:44:42', '2026-01-19 07:44:42', NULL),
(10, 'company.testing1768808883@karirnusantara.com', '$2a$10$GrXGCrL6oNq.NJAJwmhAg.aQMID5txWxyP7Df/XJuVQ39zSKhioDy', 'company', 'Budi Santoso', '081234567890', NULL, 1, 0, NULL, '2026-01-19 07:48:03', '2026-01-19 07:48:03', NULL);

--
-- Indexes for dumped tables
--

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
-- Indexes for table `companies`
--
ALTER TABLE `companies`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `user_id` (`user_id`),
  ADD KEY `fk_companies_verified_by` (`documents_verified_by`),
  ADD KEY `idx_company_status` (`company_status`),
  ADD KEY `idx_deleted_at` (`deleted_at`),
  ADD KEY `idx_created_at` (`created_at`);

--
-- Indexes for table `company_quotas`
--
ALTER TABLE `company_quotas`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `company_id` (`company_id`),
  ADD KEY `idx_company_quotas_company_id` (`company_id`);

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
-- Indexes for table `job_skills`
--
ALTER TABLE `job_skills`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `uk_job_skill` (`job_id`,`skill_name`),
  ADD KEY `idx_job_skills_skill_name` (`skill_name`);

--
-- Indexes for table `notifications`
--
ALTER TABLE `notifications`
  ADD PRIMARY KEY (`id`),
  ADD KEY `idx_notifications_user_id` (`user_id`),
  ADD KEY `idx_notifications_is_read` (`is_read`),
  ADD KEY `idx_notifications_created_at` (`created_at`);

--
-- Indexes for table `payments`
--
ALTER TABLE `payments`
  ADD PRIMARY KEY (`id`),
  ADD KEY `job_id` (`job_id`),
  ADD KEY `confirmed_by_id` (`confirmed_by_id`),
  ADD KEY `idx_payments_company_id` (`company_id`),
  ADD KEY `idx_payments_status` (`status`),
  ADD KEY `idx_payments_submitted_at` (`submitted_at`);

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
-- AUTO_INCREMENT for table `applications`
--
ALTER TABLE `applications`
  MODIFY `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=3;

--
-- AUTO_INCREMENT for table `application_timelines`
--
ALTER TABLE `application_timelines`
  MODIFY `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=10;

--
-- AUTO_INCREMENT for table `audit_logs`
--
ALTER TABLE `audit_logs`
  MODIFY `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `companies`
--
ALTER TABLE `companies`
  MODIFY `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `company_quotas`
--
ALTER TABLE `company_quotas`
  MODIFY `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=3;

--
-- AUTO_INCREMENT for table `cvs`
--
ALTER TABLE `cvs`
  MODIFY `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=3;

--
-- AUTO_INCREMENT for table `cv_snapshots`
--
ALTER TABLE `cv_snapshots`
  MODIFY `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=3;

--
-- AUTO_INCREMENT for table `jobs`
--
ALTER TABLE `jobs`
  MODIFY `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=7;

--
-- AUTO_INCREMENT for table `job_skills`
--
ALTER TABLE `job_skills`
  MODIFY `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=20;

--
-- AUTO_INCREMENT for table `notifications`
--
ALTER TABLE `notifications`
  MODIFY `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `payments`
--
ALTER TABLE `payments`
  MODIFY `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `refresh_tokens`
--
ALTER TABLE `refresh_tokens`
  MODIFY `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=60;

--
-- AUTO_INCREMENT for table `saved_jobs`
--
ALTER TABLE `saved_jobs`
  MODIFY `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `users`
--
ALTER TABLE `users`
  MODIFY `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=12;

--
-- Constraints for dumped tables
--

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
-- Constraints for table `companies`
--
ALTER TABLE `companies`
  ADD CONSTRAINT `fk_companies_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE,
  ADD CONSTRAINT `fk_companies_verified_by` FOREIGN KEY (`documents_verified_by`) REFERENCES `users` (`id`) ON DELETE SET NULL;

--
-- Constraints for table `company_quotas`
--
ALTER TABLE `company_quotas`
  ADD CONSTRAINT `company_quotas_ibfk_1` FOREIGN KEY (`company_id`) REFERENCES `users` (`id`) ON DELETE CASCADE;

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
-- Constraints for table `job_skills`
--
ALTER TABLE `job_skills`
  ADD CONSTRAINT `fk_job_skills_job` FOREIGN KEY (`job_id`) REFERENCES `jobs` (`id`) ON DELETE CASCADE;

--
-- Constraints for table `notifications`
--
ALTER TABLE `notifications`
  ADD CONSTRAINT `fk_notifications_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE;

--
-- Constraints for table `payments`
--
ALTER TABLE `payments`
  ADD CONSTRAINT `payments_ibfk_1` FOREIGN KEY (`company_id`) REFERENCES `users` (`id`) ON DELETE CASCADE,
  ADD CONSTRAINT `payments_ibfk_2` FOREIGN KEY (`job_id`) REFERENCES `jobs` (`id`) ON DELETE SET NULL,
  ADD CONSTRAINT `payments_ibfk_3` FOREIGN KEY (`confirmed_by_id`) REFERENCES `users` (`id`) ON DELETE SET NULL;

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
COMMIT;

/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
