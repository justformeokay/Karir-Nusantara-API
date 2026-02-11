-- Migration: Add indexes for job management optimization
-- Created: 2026-02-11

-- Index untuk sorting dan filtering berdasarkan created_at
CREATE INDEX IF NOT EXISTS idx_jobs_created_at ON jobs(created_at DESC);

-- Index untuk filter berdasarkan status
CREATE INDEX IF NOT EXISTS idx_jobs_status ON jobs(status);

-- Index untuk filter berdasarkan company_id
CREATE INDEX IF NOT EXISTS idx_jobs_company_id ON jobs(company_id);

-- Index untuk search berdasarkan title (LIKE search)
CREATE INDEX IF NOT EXISTS idx_jobs_title ON jobs(title);

-- Composite index untuk common filter combinations
-- Filter: status + created_at
CREATE INDEX IF NOT EXISTS idx_jobs_status_created_at ON jobs(status, created_at DESC);

-- Composite index untuk company + status
CREATE INDEX IF NOT EXISTS idx_jobs_company_status ON jobs(company_id, status);

-- Composite index untuk search + status (search untuk job title + filter status)
CREATE INDEX IF NOT EXISTS idx_jobs_title_status ON jobs(title, status);

-- Index untuk published_at (untuk active jobs)
CREATE INDEX IF NOT EXISTS idx_jobs_published_at ON jobs(published_at);

-- Index untuk views_count dan applications_count (untuk sorting statistik)
CREATE INDEX IF NOT EXISTS idx_jobs_views_count ON jobs(views_count DESC);
CREATE INDEX IF NOT EXISTS idx_jobs_applications_count ON jobs(applications_count DESC);
