-- ========================================
-- Interview Test Management Tables
-- Migration created: 2026-02-19
-- Tables untuk fitur Interview Test & Psychometric Test
-- ========================================

-- -------------------------------------------------------
-- Table: interview_tests
-- Menyimpan data Tes Psikotes/Interview
-- -------------------------------------------------------
CREATE TABLE IF NOT EXISTS `interview_tests` (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  `title` varchar(255) NOT NULL COMMENT 'Nama Tes',
  `description` longtext NOT NULL COMMENT 'Deskripsi Tes',
  `duration_minutes` int(11) NOT NULL COMMENT 'Durasi Tes dalam menit',
  `total_points` int(11) NOT NULL DEFAULT 0 COMMENT 'Total Poin dari semua pertanyaan',
  `passing_score` int(11) NOT NULL DEFAULT 70 COMMENT 'Nilai Kelulusan dalam persen (0-100)',
  `shuffle_questions` tinyint(1) NOT NULL DEFAULT 0 COMMENT 'Acak urutan soal',
  `show_results_immediately` tinyint(1) NOT NULL DEFAULT 0 COMMENT 'Tampilkan hasil langsung setelah submit',
  `status` enum('draft','active','archived') NOT NULL DEFAULT 'draft' COMMENT 'Status Tes',
  `created_by` bigint(20) UNSIGNED NOT NULL COMMENT 'Admin user_id yang membuat',
  `updated_by` bigint(20) UNSIGNED DEFAULT NULL COMMENT 'Admin user_id yang terakhir update',
  `created_at` timestamp NOT NULL DEFAULT current_timestamp(),
  `updated_at` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  `deleted_at` timestamp NULL DEFAULT NULL COMMENT 'Soft delete timestamp',
  FOREIGN KEY (`created_by`) REFERENCES `users`(`id`) ON DELETE RESTRICT,
  FOREIGN KEY (`updated_by`) REFERENCES `users`(`id`) ON DELETE SET NULL,
  KEY `idx_status` (`status`),
  KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Tabel untuk menyimpan data interview test/psikotes';

-- -------------------------------------------------------
-- Table: interview_questions
-- Menyimpan Pertanyaan untuk Interview Test
-- -------------------------------------------------------
CREATE TABLE IF NOT EXISTS `interview_questions` (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  `interview_test_id` bigint(20) UNSIGNED NOT NULL COMMENT 'Foreign key ke interview_tests',
  `question_text` longtext NOT NULL COMMENT 'Teks Pertanyaan',
  `question_type` enum('multiple_choice','essay') NOT NULL COMMENT 'Tipe Pertanyaan',
  `points` int(11) NOT NULL DEFAULT 1 COMMENT 'Poin untuk Pertanyaan ini',
  `difficulty` enum('easy','medium','hard') NOT NULL DEFAULT 'medium' COMMENT 'Tingkat Kesulitan',
  `order` int(11) NOT NULL DEFAULT 0 COMMENT 'Urutan Pertanyaan dalam Tes',
  `explanation` longtext DEFAULT NULL COMMENT 'Penjelasan jawaban (ditampilkan setelah tes)',
  `created_at` timestamp NOT NULL DEFAULT current_timestamp(),
  `updated_at` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  FOREIGN KEY (`interview_test_id`) REFERENCES `interview_tests` (`id`) ON DELETE CASCADE,
  KEY `idx_test_id` (`interview_test_id`),
  KEY `idx_difficulty` (`difficulty`),
  KEY `idx_order` (`interview_test_id`, `order`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Tabel untuk menyimpan pertanyaan di interview test';

-- -------------------------------------------------------
-- Table: interview_question_options
-- Menyimpan Opsi Jawaban untuk Pertanyaan Multiple Choice
-- -------------------------------------------------------
CREATE TABLE IF NOT EXISTS `interview_question_options` (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  `interview_question_id` bigint(20) UNSIGNED NOT NULL COMMENT 'Foreign key ke interview_questions',
  `option_text` longtext NOT NULL COMMENT 'Teks Opsi Jawaban',
  `is_correct` tinyint(1) NOT NULL DEFAULT 0 COMMENT 'Apakah ini jawaban yang benar',
  `order` int(11) NOT NULL DEFAULT 0 COMMENT 'Urutan Opsi (A, B, C, D, etc)',
  `created_at` timestamp NOT NULL DEFAULT current_timestamp(),
  `updated_at` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  FOREIGN KEY (`interview_question_id`) REFERENCES `interview_questions` (`id`) ON DELETE CASCADE,
  KEY `idx_question_id` (`interview_question_id`),
  KEY `idx_correct` (`is_correct`),
  KEY `idx_order` (`interview_question_id`, `order`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Tabel untuk menyimpan pilihan jawaban multiple choice questions';

-- -------------------------------------------------------
-- Table: interview_test_submissions
-- Menyimpan Submisi/Hasil Test dari Candidate
-- -------------------------------------------------------
CREATE TABLE IF NOT EXISTS `interview_test_submissions` (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  `interview_test_id` bigint(20) UNSIGNED NOT NULL COMMENT 'Foreign key ke interview_tests',
  `user_id` bigint(20) UNSIGNED NOT NULL COMMENT 'Candidate user_id yang mengerjakan tes',
  `application_id` bigint(20) UNSIGNED DEFAULT NULL COMMENT 'Link ke applications (optional)',
  `status` enum('in_progress','submitted','grading','completed') NOT NULL DEFAULT 'in_progress' COMMENT 'Status pengerjaan tes',
  `score` int(11) DEFAULT NULL COMMENT 'Total skor yang didapat',
  `percentage` decimal(5,2) DEFAULT NULL COMMENT 'Persentase skor (0-100)',
  `is_passed` tinyint(1) DEFAULT NULL COMMENT 'Apakah lulus sesuai passing_score',
  `started_at` timestamp NOT NULL DEFAULT current_timestamp() COMMENT 'Waktu mulai mengerjakan tes',
  `submitted_at` timestamp NULL DEFAULT NULL COMMENT 'Waktu submit tes',
  `graded_at` timestamp NULL DEFAULT NULL COMMENT 'Waktu penilaian selesai',
  `graded_by` bigint(20) UNSIGNED DEFAULT NULL COMMENT 'Admin user_id yang melakukan penilaian (essay)',
  `notes` longtext DEFAULT NULL COMMENT 'Catatan dari grader atau internal',
  `created_at` timestamp NOT NULL DEFAULT current_timestamp(),
  `updated_at` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  FOREIGN KEY (`interview_test_id`) REFERENCES `interview_tests` (`id`) ON DELETE CASCADE,
  FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE,
  FOREIGN KEY (`application_id`) REFERENCES `applications` (`id`) ON DELETE SET NULL,
  FOREIGN KEY (`graded_by`) REFERENCES `users` (`id`) ON DELETE SET NULL,
  KEY `idx_test_id` (`interview_test_id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_status` (`status`),
  KEY `idx_is_passed` (`is_passed`),
  KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Tabel untuk menyimpan submisi/hasil test dari candidate';

-- -------------------------------------------------------
-- Table: interview_test_answers
-- Menyimpan Jawaban Candidate untuk Setiap Pertanyaan
-- -------------------------------------------------------
CREATE TABLE IF NOT EXISTS `interview_test_answers` (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  `submission_id` bigint(20) UNSIGNED NOT NULL COMMENT 'Foreign key ke interview_test_submissions',
  `interview_question_id` bigint(20) UNSIGNED NOT NULL COMMENT 'Foreign key ke interview_questions',
  `question_type` enum('multiple_choice','essay') NOT NULL COMMENT 'Tipe pertanyaan',
  `answer_text` longtext DEFAULT NULL COMMENT 'Jawaban teks (untuk essay)',
  `selected_option_id` bigint(20) UNSIGNED DEFAULT NULL COMMENT 'ID opsi jawaban yang dipilih',
  `is_correct` tinyint(1) DEFAULT NULL COMMENT 'Apakah jawaban benar (NULL untuk essay pending)',
  `points_earned` int(11) DEFAULT NULL COMMENT 'Poin yang didapat untuk jawaban ini',
  `grader_feedback` longtext DEFAULT NULL COMMENT 'Feedback dari grader untuk essay answer',
  `answered_at` timestamp NOT NULL DEFAULT current_timestamp() COMMENT 'Waktu menjawab soal',
  `created_at` timestamp NOT NULL DEFAULT current_timestamp(),
  `updated_at` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  FOREIGN KEY (`submission_id`) REFERENCES `interview_test_submissions` (`id`) ON DELETE CASCADE,
  FOREIGN KEY (`interview_question_id`) REFERENCES `interview_questions` (`id`) ON DELETE CASCADE,
  FOREIGN KEY (`selected_option_id`) REFERENCES `interview_question_options` (`id`) ON DELETE SET NULL,
  KEY `idx_submission_id` (`submission_id`),
  KEY `idx_question_id` (`interview_question_id`),
  KEY `idx_option_id` (`selected_option_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Tabel untuk menyimpan jawaban candidate untuk setiap pertanyaan';

-- ========================================
-- Triggers untuk Auto Update total_points pada interview_tests
-- ========================================
DELIMITER $$
CREATE TRIGGER `after_interview_question_insert` AFTER INSERT ON `interview_questions` FOR EACH ROW
BEGIN
  UPDATE `interview_tests`
  SET `total_points` = (
    SELECT COALESCE(SUM(`points`), 0) FROM `interview_questions`
    WHERE `interview_test_id` = NEW.interview_test_id AND `deleted_at` IS NULL
  )
  WHERE `id` = NEW.interview_test_id;
END$$
DELIMITER ;

DELIMITER $$
CREATE TRIGGER `after_interview_question_update` AFTER UPDATE ON `interview_questions` FOR EACH ROW
BEGIN
  UPDATE `interview_tests`
  SET `total_points` = (
    SELECT COALESCE(SUM(`points`), 0) FROM `interview_questions`
    WHERE `interview_test_id` = NEW.interview_test_id AND `deleted_at` IS NULL
  )
  WHERE `id` = NEW.interview_test_id;
END$$
DELIMITER ;

DELIMITER $$
CREATE TRIGGER `after_interview_question_delete` AFTER DELETE ON `interview_questions` FOR EACH ROW
BEGIN
  UPDATE `interview_tests`
  SET `total_points` = (
    SELECT COALESCE(SUM(`points`), 0) FROM `interview_questions`
    WHERE `interview_test_id` = OLD.interview_test_id AND `deleted_at` IS NULL
  )
  WHERE `id` = OLD.interview_test_id;
END$$
DELIMITER ;

-- ========================================
-- End of interview test tables
-- ========================================