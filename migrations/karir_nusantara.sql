-- phpMyAdmin SQL Dump
-- version 5.2.1
-- https://www.phpmyadmin.net/
--
-- Host: localhost
-- Generation Time: Jan 20, 2026 at 09:33 AM
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

--
-- Dumping data for table `companies`
--

INSERT INTO `companies` (`id`, `user_id`, `company_name`, `company_description`, `company_website`, `company_logo_url`, `company_industry`, `company_size`, `company_location`, `company_phone`, `company_email`, `company_address`, `company_city`, `company_province`, `company_postal_code`, `established_year`, `employee_count`, `company_status`, `ktp_founder_url`, `akta_pendirian_url`, `npwp_url`, `nib_url`, `documents_verified_at`, `documents_verified_by`, `verification_notes`, `created_at`, `updated_at`, `deleted_at`) VALUES
(1, 7, 'PT Karya Developer indonesia', 'Perusahaan yang bergerak dibidang industri teknogi informasi', 'https://karyadeveloperindonesia.com', '/docs/companies/1/logo_1768832403.png', 'Teknologi Informasi', '1-10', 'Sidaorjo, Jawa Timur', '+62881036480285', 'info@karyadeveloperindonesia.com', 'Perumahan Griya Bhayangkara blok i5/07 Desa Masangan Kulon, Kecamatan Sukodono, Kabupaten Sidoarjo, Kode Pos 61258', 'Sidoarjo', 'Jawa Timur', '61258', '2025', 8, 'verified', '/docs/companies/1/ktp_1768832313.jpg', '/docs/companies/1/akta_1768832354.pdf', '/docs/companies/1/npwp_1768832371.jpg', '/docs/companies/1/nib_1768832386.pdf', '2026-01-20 04:42:11', NULL, NULL, '2026-01-19 13:58:37', '2026-01-20 04:42:11', NULL);

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
(2, 7, 1, 1, '2026-01-18 09:14:40', '2026-01-20 07:39:03'),
(3, 1, 2, 0, '2026-01-20 07:33:46', '2026-01-20 07:34:26');

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
(6, 10, 'UI/UX Designer', 'uiux-designer', 'Kami mencari UI/UX Designer untuk mengembangkan product kami', '- 2+ tahun pengalaman UI/UX\n- Figma atau Adobe XD', '- Design interface\n- User research', '- Gaji 6-10 juta/bulan', 'Jakarta Selatan', 'DKI Jakarta', 1, 'full_time', 'junior', 6000000, 10000000, 'IDR', 1, NULL, NULL, 'draft', NULL, NULL, NULL, 0, 0, NULL, '2026-01-19 07:48:03', '2026-01-19 07:48:03', NULL),
(7, 1, 'Senior Backend Developer', 'senior-backend-developer', 'We are looking for an experienced backend developer...', 'Node.js, TypeScript, Docker', 'Design and implement APIs', NULL, 'Jakarta', 'DKI Jakarta', 1, 'full_time', 'senior', NULL, NULL, 'IDR', 0, NULL, NULL, 'active', NULL, NULL, NULL, 0, 0, '2026-01-20 07:33:46', '2026-01-20 03:43:14', '2026-01-20 07:33:46', NULL),
(8, 1, 'Senior Frontend Developer', 'senior-frontend-developer', 'Kami mencari Senior Frontend Developer yang berpengalaman dalam React.js dan TypeScript untuk bergabung dengan tim kami. Posisi ini akan bertanggung jawab untuk mengembangkan fitur-fitur baru.', '- Minimal 3 tahun pengalaman dengan React.js\n- Menguasai TypeScript', NULL, NULL, 'Jakarta Selatan', 'DKI Jakarta', 1, 'full_time', 'senior', 15000000, 25000000, 'IDR', 1, NULL, NULL, 'active', NULL, NULL, NULL, 0, 0, '2026-01-20 05:05:33', '2026-01-20 05:05:33', '2026-01-20 06:11:28', '2026-01-20 06:11:28'),
(9, 1, 'Senior Frontend Developer', 'senior-frontend-developer-1768885756', 'Kami mencari Senior Frontend Developer yang berpengalaman dalam React.js dan TypeScript untuk bergabung dengan tim kami. Posisi ini akan bertanggung jawab untuk mengembangkan fitur-fitur baru.', '- Minimal 3 tahun pengalaman dengan React.js\n- Menguasai TypeScript', NULL, NULL, 'Jakarta Selatan', 'DKI Jakarta', 1, 'full_time', 'senior', 15000000, 25000000, 'IDR', 1, NULL, NULL, 'active', NULL, NULL, NULL, 0, 0, '2026-01-20 05:09:16', '2026-01-20 05:09:16', '2026-01-20 06:11:08', '2026-01-20 06:11:08'),
(10, 1, 'Senior Frontend Developer', 'senior-frontend-developer-1768886367', 'Kami mencari Senior Frontend Developer yang berpengalaman dalam React.js.', 'Minimal 3 tahun pengalaman', NULL, NULL, 'Jakarta Selatan', 'DKI Jakarta', 1, 'full_time', 'senior', 15000000, 25000000, 'IDR', 1, NULL, NULL, 'closed', NULL, NULL, NULL, 0, 0, '2026-01-20 05:19:27', '2026-01-20 05:19:27', '2026-01-20 06:10:33', NULL),
(11, 1, 'Senior Frontend Developer', 'senior-frontend-developer-1768886376', 'Kami mencari Senior Frontend Developer yang berpengalaman dalam React.js.', 'Minimal 3 tahun pengalaman', NULL, NULL, 'Jakarta Selatan', 'DKI Jakarta', 1, 'full_time', 'senior', 15000000, 25000000, 'IDR', 1, NULL, NULL, 'active', NULL, NULL, NULL, 0, 0, '2026-01-20 05:19:36', '2026-01-20 05:19:36', '2026-01-20 06:09:48', '2026-01-20 06:09:48'),
(12, 1, 'Mobile Apps Developer', 'mobile-apps-developer', 'Full job description\nAbout this Role\n\nWe are seeking a highly experienced Flutter Developer to join our mobile development team. In this role, you will take the lead in architecting and building advanced mobile applications, driving technical excellence, and contributing to high-impact projects. You will collaborate with top-tier cross-functional teams to deliver innovative, scalable, and high-performance solutions.\n\nJob Description\n\nAs a Senior Flutter Developer, you will be responsible for designing, developing, and maintaining sophisticated mobile applications in a fast-paced and innovation-driven environment.\n\nWhat It’s Like to Work Here as a Senior Flutter Mobile App DeveloperFull-cycle Technical Ownership\n\nLead the end-to-end development lifecycle of mobile applications using Flutter, from system design and architecture to deployment and maintenance.\nUI/UX Collaboration\n\nTransform complex UI/UX designs into intuitive and polished user experiences, ensuring pixel-perfect implementations and smooth interactions.\nAdvanced Integration Expertise\n\nArchitect and integrate mobile applications with robust backend systems, ensuring high performance, security, and real-time synchronization.\nCode Quality & Engineering Standards\n\nDrive and participate in in-depth code reviews, ensuring clean architecture, maintainability, and adherence to best practices across the team.\nPerformance & Scalability Optimization\n\nIdentify performance bottlenecks, analyze app behavior, and implement advanced optimization techniques across devices and platforms.\nInnovation & Technical Leadership\n\nContribute forward-thinking ideas, mentor junior developers, and help shape technical strategy and direction.\nCross-functional Technical Collaboration\n\nWork closely with product managers, backend developers, and UI/UX teams to deliver seamless, scalable solutions on time.\nContinuous Learning & Research\n\nStay ahead of the latest trends in Flutter, mobile technologies, and tools—recommending and driving adoption of relevant advancements.\nRobust Testing & Quality Assurance\n\nDevelop automated testing, debugging strategies, and quality assurance processes to ensure enterprise-level reliability and security.', 'Requirements\n\nBachelor\'s degree in Computer Science or related field (or equivalent experience).\n5+ years of professional mobile development experience, with 3+ years specifically in Flutter.\nA strong portfolio showcasing complex, high-quality Flutter applications.\nMastery of Flutter, Dart, state management (e.g., Bloc, Riverpod, Provider, GetX), and modular architecture patterns.\nProven experience integrating APIs, working with real-time data, and using third-party libraries efficiently.\nHands-on experience with backend technologies (Node.js, Django, Firebase) is a strong advantage.\nFamiliarity with CI/CD pipelines, automated testing, and modern mobile DevOps practices.\nExperience publishing and maintaining apps on the App Store and Google Play.', NULL, NULL, 'Sidaorjo', 'Jawa Timur', 0, 'full_time', 'junior', 5000000, 5500000, 'IDR', 1, '2026-02-20', NULL, 'closed', NULL, NULL, NULL, 1, 0, '2026-01-20 05:19:42', '2026-01-20 05:19:42', '2026-01-20 06:12:34', NULL),
(13, 1, 'Test Job API Debug', 'test-job-api-debug', 'Testing job creation for debug purposes and verification of company_id field. This is a longer description to pass validation.', 'Testing requirements that need to be longer for validation.', NULL, NULL, 'Jakarta', 'DKI Jakarta', 0, 'full_time', 'senior', NULL, NULL, 'IDR', 0, NULL, NULL, 'active', NULL, NULL, NULL, 0, 0, '2026-01-20 05:30:01', '2026-01-20 05:30:01', '2026-01-20 06:08:42', '2026-01-20 06:08:42'),
(14, 1, 'Mobile Apps Developer', 'mobile-apps-developer-1768890085', 'About this Role\n\nWe are seeking a highly experienced Flutter Developer to join our mobile development team. In this role, you will take the lead in architecting and building advanced mobile applications, driving technical excellence, and contributing to high-impact projects. You will collaborate with top-tier cross-functional teams to deliver innovative, scalable, and high-performance solutions.\n\nJob Description\n\nAs a Senior Flutter Developer, you will be responsible for designing, developing, and maintaining sophisticated mobile applications in a fast-paced and innovation-driven environment.\n\nWhat It’s Like to Work Here as a Senior Flutter Mobile App DeveloperFull-cycle Technical Ownership\n\nLead the end-to-end development lifecycle of mobile applications using Flutter, from system design and architecture to deployment and maintenance.\nUI/UX Collaboration\n\nTransform complex UI/UX designs into intuitive and polished user experiences, ensuring pixel-perfect implementations and smooth interactions.\nAdvanced Integration Expertise\n\nArchitect and integrate mobile applications with robust backend systems, ensuring high performance, security, and real-time synchronization.\nCode Quality & Engineering Standards\n\nDrive and participate in in-depth code reviews, ensuring clean architecture, maintainability, and adherence to best practices across the team.\nPerformance & Scalability Optimization\n\nIdentify performance bottlenecks, analyze app behavior, and implement advanced optimization techniques across devices and platforms.\nInnovation & Technical Leadership\n\nContribute forward-thinking ideas, mentor junior developers, and help shape technical strategy and direction.\nCross-functional Technical Collaboration\n\nWork closely with product managers, backend developers, and UI/UX teams to deliver seamless, scalable solutions on time.\nContinuous Learning & Research\n\nStay ahead of the latest trends in Flutter, mobile technologies, and tools—recommending and driving adoption of relevant advancements.\nRobust Testing & Quality Assurance\n\nDevelop automated testing, debugging strategies, and quality assurance processes to ensure enterprise-level reliability and security.', 'Requirements\n\nBachelor\'s degree in Computer Science or related field (or equivalent experience).\n5+ years of professional mobile development experience, with 3+ years specifically in Flutter.\nA strong portfolio showcasing complex, high-quality Flutter applications.\nMastery of Flutter, Dart, state management (e.g., Bloc, Riverpod, Provider, GetX), and modular architecture patterns.\nProven experience integrating APIs, working with real-time data, and using third-party libraries efficiently.\nHands-on experience with backend technologies (Node.js, Django, Firebase) is a strong advantage.\nFamiliarity with CI/CD pipelines, automated testing, and modern mobile DevOps practices.\nExperience publishing and maintaining apps on the App Store and Google Play.', NULL, NULL, 'Sidaorjo', 'Jawa Timur', 0, 'full_time', 'junior', 5000000, 5500000, 'IDR', 1, '2026-02-20', NULL, 'closed', NULL, NULL, NULL, 0, 0, '2026-01-20 06:21:25', '2026-01-20 06:21:25', '2026-01-20 07:02:12', '2026-01-20 07:02:12'),
(15, 1, 'QA Engineer', 'qa-engineer', 'Looking for experienced QA Engineer to join our team. You will be responsible for testing our applications and ensuring quality standards are met. Must have experience with automated testing tools.', '3+ years experience in QA testing', NULL, NULL, 'Jakarta', 'DKI Jakarta', 0, 'full_time', 'mid', NULL, NULL, 'IDR', 0, NULL, NULL, 'active', NULL, NULL, NULL, 0, 0, '2026-01-20 07:34:26', '2026-01-20 07:34:26', '2026-01-20 07:34:26', NULL),
(16, 1, 'DevOps Engineer', 'devops-engineer', 'Looking for experienced DevOps Engineer to manage our cloud infrastructure. You will be responsible for CI/CD pipelines and infrastructure automation.', '3+ years experience with AWS, Docker, Kubernetes', NULL, NULL, 'Surabaya', 'Jawa Timur', 1, 'full_time', 'mid', NULL, NULL, 'IDR', 0, NULL, NULL, 'active', NULL, NULL, NULL, 0, 0, '2026-01-20 07:39:03', '2026-01-20 07:39:03', '2026-01-20 07:39:03', NULL);

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
(29, 11, 'TypeScript', 1);

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

INSERT INTO `payments` (`id`, `company_id`, `job_id`, `amount`, `proof_image_url`, `status`, `note`, `confirmed_by_id`, `submitted_at`, `confirmed_at`, `created_at`, `updated_at`) VALUES
(1, 7, NULL, 15000, '/docs/payments/7/proof_1768893950.png', 'confirmed', 'Pembayaran telah diverifikasi', 1, '2026-01-20 07:25:50', '2026-01-20 07:29:01', '2026-01-20 07:25:50', '2026-01-20 07:29:01'),
(2, 7, NULL, 15000, '/docs/payments/7/proof_1768894335.png', 'pending', NULL, NULL, '2026-01-20 07:32:15', NULL, '2026-01-20 07:32:15', '2026-01-20 07:32:15'),
(3, 7, NULL, 15000, '/docs/payments/7/proof_1768894343.png', 'confirmed', 'Pembayaran diterima', 1, '2026-01-20 07:32:23', '2026-01-20 07:32:32', '2026-01-20 07:32:23', '2026-01-20 07:32:32');

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
(59, 7, '2b87fa6d775286438439d2b8e0c72ec2a15c7f61cef9a406c1bdceeefa22bad0', '2026-01-26 13:20:06', NULL, '', '', '2026-01-19 13:20:06'),
(60, 7, '277c6b4278d4dbc36231ed26ccacacb5212c19c79c6559231c828a09fd95d74c', '2026-01-26 13:49:57', NULL, '', '', '2026-01-19 13:49:57'),
(61, 7, '573deb02196412717a9d372689ef4854e26e5d3684ad4542360ff3c305a9d216', '2026-01-26 13:57:45', NULL, '', '', '2026-01-19 13:57:45'),
(62, 7, 'b712bb56482eb7aa73cbf8673eb651d039583d1ed2623e066b31354a85b209d1', '2026-01-26 14:01:13', NULL, '', '', '2026-01-19 14:01:13'),
(63, 7, '857b6a3aead65786de0f4e3a8011ecddb3dd1e6f20ee1cf2e81f8d3aa5c98b2a', '2026-01-26 14:18:08', NULL, '', '', '2026-01-19 14:18:08'),
(64, 7, '54eb77166f52390b05e13640b12451dd8f3c195342cba1eaf0687f93d5663ce8', '2026-01-26 14:23:15', NULL, '', '', '2026-01-19 14:23:15'),
(65, 7, 'b212a3a0942eb61f5e6ea37fa500b908174b38ad96cbc28fcea9ce3d7badf560', '2026-01-26 14:30:53', NULL, '', '', '2026-01-19 14:30:53'),
(66, 7, '81acc162faa4528132df14379cd4b3a78087f5603a31b307eefbfdffc365c1de', '2026-01-26 14:39:44', NULL, '', '', '2026-01-19 14:39:44'),
(67, 7, '0644999b481f7fc85f5d314ffc2c37758b236d1dbfaf460cdbc6d11615685bed', '2026-01-26 14:40:09', NULL, '', '', '2026-01-19 14:40:09'),
(68, 7, 'e6f46e45b9b2b6a0d62e1122c8ea393efe5c18c557065b850f5d5fd8143899c4', '2026-01-26 14:40:43', NULL, '', '', '2026-01-19 14:40:43'),
(69, 7, 'b9c1d8f68f2909a7898f3a35b2be0120d3e3c80a9c678f9a4d16ac4a8851742f', '2026-01-27 02:08:43', NULL, '', '', '2026-01-20 02:08:43'),
(70, 7, '606c5371d404836d3e036c80c137a27672039b3931182f1d08e4f446b8ad31d0', '2026-01-27 02:15:37', NULL, '', '', '2026-01-20 02:15:37'),
(71, 7, '77b5db128439adb5e41fd311798251a45c8cfa68cddaef77211e34ffa097990b', '2026-01-27 02:34:42', NULL, '', '', '2026-01-20 02:34:42'),
(72, 7, 'c66865a18e5308a3fb05a4d6dae343b14980abbb30209a83d75c3f314a136c28', '2026-01-27 02:42:06', NULL, '', '', '2026-01-20 02:42:06'),
(73, 7, '70d498c5070965c45b71c2b2d51154e678b42357ee169a629f4a780023bea7e6', '2026-01-27 02:58:00', NULL, '', '', '2026-01-20 02:58:00'),
(74, 7, '7692b2d788f3eeb82a8aad9726fd25fb30ca41236b8442f2e7f08d305b95bb4d', '2026-01-27 02:58:15', NULL, '', '', '2026-01-20 02:58:15'),
(75, 7, '448713cf8a721aab434daa2b07dcef61e9451f8bc538c6206fe3df6505b8a0f0', '2026-01-27 03:07:37', NULL, '', '', '2026-01-20 03:07:37'),
(76, 7, 'b078bae01a9c5f16236f0dda6301b4a41783c123cd767f49a42c04a1bffed631', '2026-01-27 03:08:18', NULL, '', '', '2026-01-20 03:08:18'),
(77, 7, '45416cb7e854a0c565dc0565a6e88ae4ae27ca2739643b0f7715b5bd0ed023be', '2026-01-27 03:09:40', NULL, '', '', '2026-01-20 03:09:40'),
(78, 7, '45875c5981930bee5d3273a54dca1d744a5016c76e4f4aa3f809990078965b14', '2026-01-27 03:09:49', NULL, '', '', '2026-01-20 03:09:49'),
(79, 7, 'f0fd8265affdf10b594851c7689ec62ce68fbf35134d3be58ee23f6fef47b279', '2026-01-27 03:10:16', NULL, '', '', '2026-01-20 03:10:16'),
(80, 7, 'df9c221d52f392e6267c40eef51a8ec5c3282df565b2334d74a5961b175e7ac9', '2026-01-27 03:10:29', NULL, '', '', '2026-01-20 03:10:29'),
(81, 7, '809de13f87c8bd92ae95d2c806ddb320c8e4158a4f7ae6bf261277f9a4284ffe', '2026-01-27 03:13:00', NULL, '', '', '2026-01-20 03:13:00'),
(82, 7, 'dede9156ce8b831a98528aa5b4554d894f3160caded91ea5e8f5c28081ee00d6', '2026-01-27 03:13:08', NULL, '', '', '2026-01-20 03:13:08'),
(83, 7, 'ef95e3e43bfb9cf8c0d64cf9e232b7d13495e7efee2f70d93292e348bfb3103b', '2026-01-27 03:13:52', NULL, '', '', '2026-01-20 03:13:52'),
(84, 7, '509b8d1a4b6ae911fd78899220bf183a78b5893739ffda1c2a73971d78fabbd9', '2026-01-27 03:14:23', NULL, '', '', '2026-01-20 03:14:23'),
(85, 7, '1f280b56a341a2e87d2184feb7839c5a248ceb2aa4c268ca154b31606d626b40', '2026-01-27 03:24:30', NULL, '', '', '2026-01-20 03:24:30'),
(86, 7, 'e49a2f7582e40c4faee703e1949dd9a7749500f685ba142c7a50a051ab345ead', '2026-01-27 03:35:55', NULL, '', '', '2026-01-20 03:35:55'),
(87, 7, '35b89f2b1540359403d0c11c6594a3829dab3ad6427536ec65af58e28e8f8d8d', '2026-01-27 03:38:47', NULL, '', '', '2026-01-20 03:38:47'),
(88, 1, '757f8097ca62f94e1fdede4c0daa1964c9409d075f690d8b8f2101698a1f6113', '2026-01-27 03:39:10', NULL, '', '', '2026-01-20 03:39:10'),
(89, 1, '5f3d2a21161326c932efeda1b5791ea183ef9cc11b7fa955c7f5c6e073501ab9', '2026-01-27 03:39:29', NULL, '', '', '2026-01-20 03:39:29'),
(90, 1, '4588c3e22ca04ff044c40dda7c5b962c7404575cf9c4e7228a6de603d1fc901e', '2026-01-27 03:42:26', NULL, '', '', '2026-01-20 03:42:26'),
(91, 1, '82bc1046e59836de187fc60cc041b3d91c34bfaf947a2ddc93c8b66a83e7c76d', '2026-01-27 03:43:00', NULL, '', '', '2026-01-20 03:43:00'),
(92, 7, '912696beb110543df2a2576f33042c7bec87f04203442d0ba2acec9e3056798f', '2026-01-27 03:43:14', NULL, '', '', '2026-01-20 03:43:14'),
(93, 7, '52a7e5200ea85a5f23889d5f274b16653e3dc3d445dbc7a851130dbf54a61dd0', '2026-01-27 03:51:22', NULL, '', '', '2026-01-20 03:51:22'),
(94, 12, '671b5b9f4eab05ddabd87bdbe0c996b84e45a237194325674f6c87e3e084ddad', '2026-01-27 04:39:36', NULL, '', '', '2026-01-20 04:39:36'),
(95, 7, '41b91ae76fe38c780adfa4fec34f3fcfb682bc5006f88449353ae02e49a89340', '2026-01-27 04:45:56', NULL, '', '', '2026-01-20 04:45:56'),
(96, 7, '078394421ada52af0f3ce1fc95284550834e730d136d50436cc7b9b1b83810aa', '2026-01-27 05:03:43', NULL, '', '', '2026-01-20 05:03:43'),
(97, 7, 'a8ea2ad138f34e49ec75e6374d8ae61ef9af45e3eb9c6ce13b1d0bb81a91ae79', '2026-01-27 05:03:52', NULL, '', '', '2026-01-20 05:03:52'),
(98, 7, 'f4876476ab5ed43b3d48afc3df420aa9edc00932271611fd17c78af43e89ed9f', '2026-01-27 05:09:04', NULL, '', '', '2026-01-20 05:09:04'),
(99, 7, '38b15225b835bccb57e632f23898c7ce1a642ccf26933b08613dac149b3409d4', '2026-01-27 05:11:57', NULL, '', '', '2026-01-20 05:11:57'),
(100, 7, 'c00c1dbd409e48265a6b1fcc816f5b60069ce028641fe2fc7ca1f681556eed4b', '2026-01-27 05:18:09', NULL, '', '', '2026-01-20 05:18:09'),
(101, 7, '44e80ca96e11eca5e24c3f566718a2f4645cc6d7420b3d68db58f7c74ae19883', '2026-01-27 05:18:15', NULL, '', '', '2026-01-20 05:18:15'),
(102, 7, '3e360f4c576b9d4e422ea73e34be89b628da2f3aebe609e4e9bdaaa15e043c49', '2026-01-27 05:27:36', NULL, '', '', '2026-01-20 05:27:36'),
(103, 7, 'a483bfa3ffb8c4953ce18df9020b74647859df40a75b02dbebc7164a39c3f9cc', '2026-01-27 05:34:03', NULL, '', '', '2026-01-20 05:34:03'),
(104, 7, '05299d6089289c5829f03cd3a16612d8f0d07f148642e564a95ffcb79194ba2f', '2026-01-27 05:54:36', NULL, '', '', '2026-01-20 05:54:36'),
(105, 7, 'e22b19117cf558ecf5027ccfdb2a38ac5d8350b299d1d127296a8b6a06a819b9', '2026-01-27 06:09:40', NULL, '', '', '2026-01-20 06:09:40'),
(106, 7, 'f78e739c2a6b94e91e075398b0a2e9709a9fe1d6142bfb8c8a6ba7e4b4fc3fd9', '2026-01-27 06:26:53', NULL, '', '', '2026-01-20 06:26:53'),
(107, 7, '3c45cb8142020c1660fe4f35de367bafe6fff73c512cf25c71ea6373ecc0cee6', '2026-01-27 06:27:02', NULL, '', '', '2026-01-20 06:27:02'),
(108, 7, '567134ad71cfab1b6396537faad83802f86ebaea5f5c10d274f03fc36c410500', '2026-01-27 06:27:09', NULL, '', '', '2026-01-20 06:27:09'),
(109, 7, '52c742fe1d8d1d136a768fa947c7f90cb9c045e97a03ba4ef961fc8b0d585490', '2026-01-27 06:27:18', NULL, '', '', '2026-01-20 06:27:18'),
(110, 7, '4ebb80ffa28db3ad3aca43c203ba5ddd15a297241f2eb2683f3dd2a6b978af87', '2026-01-27 06:31:58', NULL, '', '', '2026-01-20 06:31:58'),
(111, 7, '2dac3522608bc84455a48efe8e7d794d58cf54c07a373e01d979b65365830673', '2026-01-27 06:34:17', NULL, '', '', '2026-01-20 06:34:17'),
(112, 7, '04792a89c031d49cf3090a1f1370631150523fa1137a7d16ac2a1534e83d6f3d', '2026-01-27 06:35:17', NULL, '', '', '2026-01-20 06:35:17'),
(113, 7, 'ebf7cfc30224915ae0e325470f1a0633056d9f4db1229224475d90213333a994', '2026-01-27 06:37:57', NULL, '', '', '2026-01-20 06:37:57'),
(114, 7, '7092720f84b43092afe6eb57f1615b00e5e705863d0a2089d690d929e97b4a67', '2026-01-27 06:43:09', NULL, '', '', '2026-01-20 06:43:09'),
(115, 7, '39f955f9ac80739ca001edfbfd5a961af0102e60c9db8e7d72b9fa0eb1618095', '2026-01-27 06:43:16', NULL, '', '', '2026-01-20 06:43:16'),
(116, 7, '7830a709c10076cf406c4ffa9f3b92853f08883bc80109167e42ebc4d407689c', '2026-01-27 06:44:14', NULL, '', '', '2026-01-20 06:44:14'),
(117, 7, '92cd1a132e95f06f8b505997a4cc4da83aee2c9647f74984bead31d0f9b52a81', '2026-01-27 07:01:11', NULL, '', '', '2026-01-20 07:01:11'),
(118, 7, '5e830f1ec835c53822a30451f352791af5baed81dd91637735befeade7ba212b', '2026-01-27 07:14:54', NULL, '', '', '2026-01-20 07:14:54'),
(119, 1, 'f793b289ec992152ee29fe6a67ee09755930d8f4ff140e2259b5e21c8d6da028', '2026-01-27 07:23:29', NULL, '', '', '2026-01-20 07:23:29'),
(120, 7, '8edbedc91fd5afe59d2bf6a1713d84b7ed87d28fb1dadffc77564a11a8741e2e', '2026-01-27 07:23:49', NULL, '', '', '2026-01-20 07:23:49'),
(121, 7, 'ddae9b2f1490b3139245b1ee15ae0b6573c28a3645caab214a823933574e92c7', '2026-01-27 07:31:40', NULL, '', '', '2026-01-20 07:31:40'),
(122, 1, '9278c0609b57d2bdb7b54f29dc36743c2acac55a7dc08745bb01fb7ac6128069', '2026-01-27 07:31:40', NULL, '', '', '2026-01-20 07:31:40'),
(123, 7, '6297cfd594eea06f6ccad1bc355742e217f445906b288b70c8ef1c54fb6357b5', '2026-01-27 07:31:47', NULL, '', '', '2026-01-20 07:31:47'),
(124, 1, '33f173c2e4512119fa09a9364bac2f041d992111095c89688e98790d515c7a04', '2026-01-27 07:31:55', NULL, '', '', '2026-01-20 07:31:55'),
(125, 7, '88d204fda219c6f34d9a99ce421e68ad9f8f66c1a9cfbf5a3015e1279842f310', '2026-01-27 07:38:23', NULL, '', '', '2026-01-20 07:38:23'),
(126, 7, '1e1fa7e7e3855209fa2162ff1e7765e8f74aa410362ab29d72ed4f361fcf2fef', '2026-01-27 07:38:33', NULL, '', '', '2026-01-20 07:38:33'),
(127, 1, '80c99e42686621fa56568b13ce1996d267ec936bec3c910c479be6e79da142ec', '2026-01-27 07:39:26', NULL, '', '', '2026-01-20 07:39:26'),
(128, 7, '61fd0cd20ef183f1acec515956896bb394b1284b17fc5d5f43109ac9007483cd', '2026-01-27 07:41:36', NULL, '', '', '2026-01-20 07:41:36'),
(129, 7, '059b253aa2a70bd13d6e23759a4cc5cbb5000b009201b4b75da94ff041997c54', '2026-01-27 07:46:42', NULL, '', '', '2026-01-20 07:46:42'),
(130, 7, 'be538b0bc7b7cf3e00d2c60ba19e0bb7e58814c7a6c880a6ab42502643df0d7e', '2026-01-27 08:28:06', NULL, '', '', '2026-01-20 08:28:06');

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
(7, 'info@karyadeveloperindonesia.com', '$2a$10$yGx97vEQaIJOEeuDHlIGQONjqgqfdFhWUbBrWWxsJNpmeEXf5reWW', 'company', 'Admin', '0881036480285', NULL, 1, 1, NULL, '2026-01-18 09:13:55', '2026-01-20 04:45:49', NULL),
(9, 'company.testing1768808682@karirnusantara.com', '$2a$10$1RpgFQaG58KBk.Fhy8laye2yhBj.C.ogGGaqxkktzmew7TMrlLu.W', 'company', 'Budi Santoso', '081234567890', NULL, 1, 0, NULL, '2026-01-19 07:44:42', '2026-01-19 07:44:42', NULL),
(10, 'company.testing1768808883@karirnusantara.com', '$2a$10$GrXGCrL6oNq.NJAJwmhAg.aQMID5txWxyP7Df/XJuVQ39zSKhioDy', 'company', 'Budi Santoso', '081234567890', NULL, 1, 0, NULL, '2026-01-19 07:48:03', '2026-01-19 07:48:03', NULL),
(12, 'test.company@example.com', '$2a$10$Q7OZlZckA0Tw5RarsLHdHOctnn3n9Sh2J.SM1OwVX58DxtA8Q74le', 'company', 'Test Company Manager', '081234567890', NULL, 1, 0, NULL, '2026-01-20 04:39:36', '2026-01-20 04:39:36', NULL);

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
  MODIFY `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=2;

--
-- AUTO_INCREMENT for table `company_quotas`
--
ALTER TABLE `company_quotas`
  MODIFY `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=4;

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
  MODIFY `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=17;

--
-- AUTO_INCREMENT for table `job_skills`
--
ALTER TABLE `job_skills`
  MODIFY `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=30;

--
-- AUTO_INCREMENT for table `notifications`
--
ALTER TABLE `notifications`
  MODIFY `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `payments`
--
ALTER TABLE `payments`
  MODIFY `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=4;

--
-- AUTO_INCREMENT for table `refresh_tokens`
--
ALTER TABLE `refresh_tokens`
  MODIFY `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=131;

--
-- AUTO_INCREMENT for table `saved_jobs`
--
ALTER TABLE `saved_jobs`
  MODIFY `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `users`
--
ALTER TABLE `users`
  MODIFY `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=13;

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
