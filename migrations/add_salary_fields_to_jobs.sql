-- Add fixed salary support to jobs table
ALTER TABLE `jobs` ADD COLUMN `is_salary_fixed` tinyint(1) NOT NULL DEFAULT 0 AFTER `is_salary_visible`;

-- Existing salary_min/salary_max will be:
-- - is_salary_fixed = 1: salary_min = fixed amount, salary_max = NULL
-- - is_salary_fixed = 0: salary_min/salary_max = range (existing behavior)
