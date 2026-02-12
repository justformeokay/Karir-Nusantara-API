-- Add category column to jobs table
ALTER TABLE `jobs` ADD COLUMN `category` varchar(50) NOT NULL DEFAULT 'Engineering' AFTER `title`;

-- Add index for better query performance
ALTER TABLE `jobs` ADD INDEX `idx_category` (`category`);
