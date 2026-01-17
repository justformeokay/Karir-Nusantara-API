-- ============================================
-- KARIR NUSANTARA DATABASE SCHEMA
-- ============================================
-- Database: MySQL 8.0+
-- Charset: utf8mb4
-- Collation: utf8mb4_unicode_ci
-- ============================================

-- Create database (run separately if needed)
-- CREATE DATABASE IF NOT EXISTS karir_nusantara CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
-- USE karir_nusantara;

-- ============================================
-- USERS TABLE
-- ============================================
-- Supports multiple roles: job_seeker, company, admin
-- Extensible for future features (social login, 2FA, etc.)

CREATE TABLE IF NOT EXISTS users (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    role ENUM('job_seeker', 'company', 'admin') NOT NULL DEFAULT 'job_seeker',
    
    -- Profile fields
    full_name VARCHAR(255) NOT NULL,
    phone VARCHAR(20),
    avatar_url VARCHAR(500),
    
    -- Company-specific fields (nullable for job seekers)
    company_name VARCHAR(255),
    company_description TEXT,
    company_website VARCHAR(500),
    company_logo_url VARCHAR(500),
    
    -- Account status
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    is_verified BOOLEAN NOT NULL DEFAULT FALSE,
    email_verified_at TIMESTAMP NULL,
    
    -- Timestamps
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    
    -- Indexes
    INDEX idx_users_email (email),
    INDEX idx_users_role (role),
    INDEX idx_users_is_active (is_active),
    INDEX idx_users_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================
-- CVS TABLE
-- ============================================
-- Each user can have one active CV draft
-- CV data stored as JSON for flexibility
-- Supports versioning through cv_snapshots

CREATE TABLE IF NOT EXISTS cvs (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT UNSIGNED NOT NULL UNIQUE,
    
    -- Personal Information (JSON object)
    -- Structure: { "full_name", "email", "phone", "address", "city", "province", "summary", "linkedin", "portfolio" }
    personal_info JSON NOT NULL,
    
    -- Education (JSON array)
    -- Structure: [{ "institution", "degree", "field_of_study", "start_date", "end_date", "gpa", "description" }]
    education JSON NOT NULL DEFAULT '[]',
    
    -- Work Experience (JSON array)
    -- Structure: [{ "company", "position", "location", "start_date", "end_date", "is_current", "description", "achievements" }]
    experience JSON NOT NULL DEFAULT '[]',
    
    -- Skills (JSON array)
    -- Structure: [{ "name", "level", "category" }]
    skills JSON NOT NULL DEFAULT '[]',
    
    -- Certifications (JSON array)
    -- Structure: [{ "name", "issuer", "issue_date", "expiry_date", "credential_id", "credential_url" }]
    certifications JSON NOT NULL DEFAULT '[]',
    
    -- Additional sections (extensible)
    languages JSON NOT NULL DEFAULT '[]',
    projects JSON NOT NULL DEFAULT '[]',
    
    -- Metadata
    last_updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    completeness_score INT UNSIGNED NOT NULL DEFAULT 0,
    
    -- Timestamps
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    -- Foreign Keys
    CONSTRAINT fk_cvs_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    
    -- Indexes
    INDEX idx_cvs_user_id (user_id),
    INDEX idx_cvs_completeness (completeness_score)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================
-- CV SNAPSHOTS TABLE
-- ============================================
-- Immutable snapshots of CV at application time
-- Ensures application history integrity

CREATE TABLE IF NOT EXISTS cv_snapshots (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    cv_id BIGINT UNSIGNED NOT NULL,
    user_id BIGINT UNSIGNED NOT NULL,
    
    -- Complete CV data snapshot (denormalized for immutability)
    personal_info JSON NOT NULL,
    education JSON NOT NULL,
    experience JSON NOT NULL,
    skills JSON NOT NULL,
    certifications JSON NOT NULL,
    languages JSON NOT NULL DEFAULT '[]',
    projects JSON NOT NULL DEFAULT '[]',
    
    -- Snapshot metadata
    snapshot_hash VARCHAR(64) NOT NULL,
    completeness_score INT UNSIGNED NOT NULL DEFAULT 0,
    
    -- Timestamps
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    -- Foreign Keys
    CONSTRAINT fk_cv_snapshots_cv FOREIGN KEY (cv_id) REFERENCES cvs(id) ON DELETE CASCADE,
    CONSTRAINT fk_cv_snapshots_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    
    -- Indexes
    INDEX idx_cv_snapshots_cv_id (cv_id),
    INDEX idx_cv_snapshots_user_id (user_id),
    INDEX idx_cv_snapshots_hash (snapshot_hash)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================
-- JOBS TABLE
-- ============================================
-- Job postings by companies
-- Supports filtering, search, and status management

CREATE TABLE IF NOT EXISTS jobs (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    company_id BIGINT UNSIGNED NOT NULL,
    
    -- Job Details
    title VARCHAR(255) NOT NULL,
    slug VARCHAR(300) NOT NULL UNIQUE,
    description TEXT NOT NULL,
    requirements TEXT,
    responsibilities TEXT,
    benefits TEXT,
    
    -- Location
    city VARCHAR(100) NOT NULL,
    province VARCHAR(100) NOT NULL,
    is_remote BOOLEAN NOT NULL DEFAULT FALSE,
    
    -- Job Type & Level
    job_type ENUM('full_time', 'part_time', 'contract', 'internship', 'freelance') NOT NULL DEFAULT 'full_time',
    experience_level ENUM('entry', 'junior', 'mid', 'senior', 'lead', 'executive') NOT NULL DEFAULT 'entry',
    
    -- Salary (optional, for transparency)
    salary_min BIGINT UNSIGNED,
    salary_max BIGINT UNSIGNED,
    salary_currency VARCHAR(3) DEFAULT 'IDR',
    is_salary_visible BOOLEAN NOT NULL DEFAULT TRUE,
    
    -- Application settings
    application_deadline DATE,
    max_applications INT UNSIGNED,
    
    -- Status
    status ENUM('draft', 'active', 'paused', 'closed', 'filled') NOT NULL DEFAULT 'draft',
    
    -- Metadata
    views_count INT UNSIGNED NOT NULL DEFAULT 0,
    applications_count INT UNSIGNED NOT NULL DEFAULT 0,
    
    -- Timestamps
    published_at TIMESTAMP NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    
    -- Foreign Keys
    CONSTRAINT fk_jobs_company FOREIGN KEY (company_id) REFERENCES users(id) ON DELETE CASCADE,
    
    -- Indexes
    INDEX idx_jobs_company_id (company_id),
    INDEX idx_jobs_status (status),
    INDEX idx_jobs_job_type (job_type),
    INDEX idx_jobs_city (city),
    INDEX idx_jobs_province (province),
    INDEX idx_jobs_salary (salary_min, salary_max),
    INDEX idx_jobs_published_at (published_at),
    INDEX idx_jobs_deleted_at (deleted_at),
    
    -- Full-text search index
    FULLTEXT INDEX idx_jobs_search (title, description, requirements)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================
-- JOB SKILLS TABLE
-- ============================================
-- Many-to-many relationship for job required skills

CREATE TABLE IF NOT EXISTS job_skills (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    job_id BIGINT UNSIGNED NOT NULL,
    skill_name VARCHAR(100) NOT NULL,
    is_required BOOLEAN NOT NULL DEFAULT TRUE,
    
    -- Foreign Keys
    CONSTRAINT fk_job_skills_job FOREIGN KEY (job_id) REFERENCES jobs(id) ON DELETE CASCADE,
    
    -- Unique constraint
    UNIQUE KEY uk_job_skill (job_id, skill_name),
    
    -- Indexes
    INDEX idx_job_skills_skill_name (skill_name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================
-- APPLICATIONS TABLE
-- ============================================
-- Job applications by job seekers
-- Links user, job, and CV snapshot

CREATE TABLE IF NOT EXISTS applications (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT UNSIGNED NOT NULL,
    job_id BIGINT UNSIGNED NOT NULL,
    cv_snapshot_id BIGINT UNSIGNED NOT NULL,
    
    -- Application details
    cover_letter TEXT,
    
    -- Current status (denormalized for quick access)
    current_status ENUM(
        'submitted',
        'viewed',
        'shortlisted',
        'interview_scheduled',
        'interview_completed',
        'assessment',
        'offer_sent',
        'offer_accepted',
        'hired',
        'rejected',
        'withdrawn'
    ) NOT NULL DEFAULT 'submitted',
    
    -- Timestamps
    applied_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    last_status_update TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    -- Foreign Keys
    CONSTRAINT fk_applications_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_applications_job FOREIGN KEY (job_id) REFERENCES jobs(id) ON DELETE CASCADE,
    CONSTRAINT fk_applications_cv_snapshot FOREIGN KEY (cv_snapshot_id) REFERENCES cv_snapshots(id) ON DELETE RESTRICT,
    
    -- Prevent duplicate applications
    UNIQUE KEY uk_user_job (user_id, job_id),
    
    -- Indexes
    INDEX idx_applications_user_id (user_id),
    INDEX idx_applications_job_id (job_id),
    INDEX idx_applications_status (current_status),
    INDEX idx_applications_applied_at (applied_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================
-- APPLICATION TIMELINES TABLE
-- ============================================
-- Immutable, append-only timeline events
-- Core transparency feature

CREATE TABLE IF NOT EXISTS application_timelines (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    application_id BIGINT UNSIGNED NOT NULL,
    
    -- Event details
    status ENUM(
        'submitted',
        'viewed',
        'shortlisted',
        'interview_scheduled',
        'interview_completed',
        'assessment',
        'offer_sent',
        'offer_accepted',
        'hired',
        'rejected',
        'withdrawn'
    ) NOT NULL,
    
    -- Event metadata
    note TEXT,
    is_visible_to_applicant BOOLEAN NOT NULL DEFAULT TRUE,
    
    -- Actor tracking
    updated_by_type ENUM('system', 'company', 'applicant') NOT NULL DEFAULT 'system',
    updated_by_id BIGINT UNSIGNED,
    
    -- Scheduled event support (for interviews)
    scheduled_at TIMESTAMP NULL,
    scheduled_location VARCHAR(500),
    scheduled_notes TEXT,
    
    -- Immutable timestamp
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    -- Foreign Keys
    CONSTRAINT fk_timelines_application FOREIGN KEY (application_id) REFERENCES applications(id) ON DELETE CASCADE,
    
    -- Indexes
    INDEX idx_timelines_application_id (application_id),
    INDEX idx_timelines_status (status),
    INDEX idx_timelines_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================
-- SAVED JOBS TABLE
-- ============================================
-- Users can save jobs for later

CREATE TABLE IF NOT EXISTS saved_jobs (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT UNSIGNED NOT NULL,
    job_id BIGINT UNSIGNED NOT NULL,
    
    -- Timestamps
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    -- Foreign Keys
    CONSTRAINT fk_saved_jobs_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_saved_jobs_job FOREIGN KEY (job_id) REFERENCES jobs(id) ON DELETE CASCADE,
    
    -- Prevent duplicates
    UNIQUE KEY uk_saved_job (user_id, job_id),
    
    -- Indexes
    INDEX idx_saved_jobs_user_id (user_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================
-- REFRESH TOKENS TABLE
-- ============================================
-- For JWT refresh token management

CREATE TABLE IF NOT EXISTS refresh_tokens (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT UNSIGNED NOT NULL,
    token_hash VARCHAR(255) NOT NULL UNIQUE,
    expires_at TIMESTAMP NOT NULL,
    revoked_at TIMESTAMP NULL,
    
    -- Device/session tracking
    device_info VARCHAR(500),
    ip_address VARCHAR(45),
    
    -- Timestamps
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    -- Foreign Keys
    CONSTRAINT fk_refresh_tokens_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    
    -- Indexes
    INDEX idx_refresh_tokens_user_id (user_id),
    INDEX idx_refresh_tokens_expires_at (expires_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================
-- NOTIFICATIONS TABLE (Future-ready)
-- ============================================

CREATE TABLE IF NOT EXISTS notifications (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT UNSIGNED NOT NULL,
    
    -- Notification content
    type VARCHAR(50) NOT NULL,
    title VARCHAR(255) NOT NULL,
    message TEXT NOT NULL,
    data JSON,
    
    -- Status
    is_read BOOLEAN NOT NULL DEFAULT FALSE,
    read_at TIMESTAMP NULL,
    
    -- Timestamps
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    -- Foreign Keys
    CONSTRAINT fk_notifications_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    
    -- Indexes
    INDEX idx_notifications_user_id (user_id),
    INDEX idx_notifications_is_read (is_read),
    INDEX idx_notifications_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================
-- AUDIT LOGS TABLE (Future-ready)
-- ============================================

CREATE TABLE IF NOT EXISTS audit_logs (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT UNSIGNED,
    
    -- Action details
    action VARCHAR(100) NOT NULL,
    entity_type VARCHAR(50) NOT NULL,
    entity_id BIGINT UNSIGNED,
    
    -- Change tracking
    old_values JSON,
    new_values JSON,
    
    -- Request context
    ip_address VARCHAR(45),
    user_agent VARCHAR(500),
    
    -- Timestamps
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    -- Indexes
    INDEX idx_audit_logs_user_id (user_id),
    INDEX idx_audit_logs_entity (entity_type, entity_id),
    INDEX idx_audit_logs_action (action),
    INDEX idx_audit_logs_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
