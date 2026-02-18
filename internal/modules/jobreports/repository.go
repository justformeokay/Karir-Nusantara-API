package jobreports

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/karirnusantara/api/internal/shared/hashid"
)

// Repository handles database operations for job reports
type Repository struct {
	db *sqlx.DB
}

// NewRepository creates a new job reports repository
func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

// Create creates a new job report
func (r *Repository) Create(ctx context.Context, report *JobReport) error {
	query := `
		INSERT INTO job_reports (job_id, user_id, reason, description, status, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, NOW(), NOW())
	`
	result, err := r.db.ExecContext(ctx, query, report.JobID, report.UserID, report.Reason, report.Description, ReportStatusPending)
	if err != nil {
		return fmt.Errorf("failed to create job report: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}
	report.ID = uint64(id)
	return nil
}

// GetByID gets a job report by ID
func (r *Repository) GetByID(ctx context.Context, id uint64) (*JobReportWithDetails, error) {
	query := `
		SELECT 
			jr.id, jr.job_id, jr.user_id, jr.reason, jr.description, jr.status, 
			COALESCE(jr.admin_note, '') as admin_note,
			jr.reviewed_by, jr.reviewed_at, jr.created_at, jr.updated_at,
			j.title as job_title, j.status as job_status,
			c.id as company_id, c.company_name, c.company_status,
			u.full_name as reporter_name, u.email as reporter_email,
			COALESCE(reviewer.full_name, '') as reviewer_name,
			(SELECT COUNT(*) FROM job_reports WHERE job_id = jr.job_id) as total_reports
		FROM job_reports jr
		JOIN jobs j ON jr.job_id = j.id
		JOIN companies c ON j.company_id = c.id
		JOIN users u ON jr.user_id = u.id
		LEFT JOIN users reviewer ON jr.reviewed_by = reviewer.id
		WHERE jr.id = ?
	`

	var report JobReportWithDetails
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&report.ID, &report.JobID, &report.UserID, &report.Reason, &report.Description, &report.Status,
		&report.AdminNote, &report.ReviewedBy, &report.ReviewedAt, &report.CreatedAt, &report.UpdatedAt,
		&report.JobTitle, &report.JobStatus,
		&report.CompanyID, &report.CompanyName, &report.CompanyStatus,
		&report.ReporterName, &report.ReporterEmail, &report.ReviewerName, &report.TotalReports,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get job report: %w", err)
	}
	report.CompanyHashID = hashid.Encode(report.CompanyID)
	return &report, nil
}

// GetAll gets all job reports with optional filtering
func (r *Repository) GetAll(ctx context.Context, status string, page, limit int) ([]JobReportWithDetails, int, error) {
	offset := (page - 1) * limit

	// Count query
	countQuery := `
		SELECT COUNT(*)
		FROM job_reports jr
		JOIN jobs j ON jr.job_id = j.id
		JOIN companies c ON j.company_id = c.id
		WHERE (? = '' OR jr.status = ?)
	`
	var total int
	if err := r.db.QueryRowContext(ctx, countQuery, status, status).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("failed to count reports: %w", err)
	}

	// Data query
	query := `
		SELECT 
			jr.id, jr.job_id, jr.user_id, jr.reason, jr.description, jr.status, 
			COALESCE(jr.admin_note, '') as admin_note,
			jr.reviewed_by, jr.reviewed_at, jr.created_at, jr.updated_at,
			j.title as job_title, j.status as job_status,
			c.id as company_id, c.company_name, c.company_status,
			u.full_name as reporter_name, u.email as reporter_email,
			COALESCE(reviewer.full_name, '') as reviewer_name,
			(SELECT COUNT(*) FROM job_reports WHERE job_id = jr.job_id) as total_reports
		FROM job_reports jr
		JOIN jobs j ON jr.job_id = j.id
		JOIN companies c ON j.company_id = c.id
		JOIN users u ON jr.user_id = u.id
		LEFT JOIN users reviewer ON jr.reviewed_by = reviewer.id
		WHERE (? = '' OR jr.status = ?)
		ORDER BY jr.created_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.QueryContext(ctx, query, status, status, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get job reports: %w", err)
	}
	defer rows.Close()

	var reports []JobReportWithDetails
	for rows.Next() {
		var report JobReportWithDetails
		err := rows.Scan(
			&report.ID, &report.JobID, &report.UserID, &report.Reason, &report.Description, &report.Status,
			&report.AdminNote, &report.ReviewedBy, &report.ReviewedAt, &report.CreatedAt, &report.UpdatedAt,
			&report.JobTitle, &report.JobStatus,
			&report.CompanyID, &report.CompanyName, &report.CompanyStatus,
			&report.ReporterName, &report.ReporterEmail, &report.ReviewerName, &report.TotalReports,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan job report: %w", err)
		}
		report.CompanyHashID = hashid.Encode(report.CompanyID)
		reports = append(reports, report)
	}

	return reports, total, nil
}

// GetByJobID gets all reports for a specific job
func (r *Repository) GetByJobID(ctx context.Context, jobID uint64) ([]JobReportWithDetails, error) {
	query := `
		SELECT 
			jr.id, jr.job_id, jr.user_id, jr.reason, jr.description, jr.status, 
			COALESCE(jr.admin_note, '') as admin_note,
			jr.reviewed_by, jr.reviewed_at, jr.created_at, jr.updated_at,
			j.title as job_title, j.status as job_status,
			c.id as company_id, c.company_name, c.company_status,
			u.full_name as reporter_name, u.email as reporter_email,
			COALESCE(reviewer.full_name, '') as reviewer_name,
			(SELECT COUNT(*) FROM job_reports WHERE job_id = jr.job_id) as total_reports
		FROM job_reports jr
		JOIN jobs j ON jr.job_id = j.id
		JOIN companies c ON j.company_id = c.id
		JOIN users u ON jr.user_id = u.id
		LEFT JOIN users reviewer ON jr.reviewed_by = reviewer.id
		WHERE jr.job_id = ?
		ORDER BY jr.created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, jobID)
	if err != nil {
		return nil, fmt.Errorf("failed to get job reports: %w", err)
	}
	defer rows.Close()

	var reports []JobReportWithDetails
	for rows.Next() {
		var report JobReportWithDetails
		err := rows.Scan(
			&report.ID, &report.JobID, &report.UserID, &report.Reason, &report.Description, &report.Status,
			&report.AdminNote, &report.ReviewedBy, &report.ReviewedAt, &report.CreatedAt, &report.UpdatedAt,
			&report.JobTitle, &report.JobStatus,
			&report.CompanyID, &report.CompanyName, &report.CompanyStatus,
			&report.ReporterName, &report.ReporterEmail, &report.ReviewerName, &report.TotalReports,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan job report: %w", err)
		}
		report.CompanyHashID = hashid.Encode(report.CompanyID)
		reports = append(reports, report)
	}

	return reports, nil
}

// UpdateStatus updates the status of a job report
func (r *Repository) UpdateStatus(ctx context.Context, id uint64, status, adminNote string, reviewedBy uint64) error {
	query := `
		UPDATE job_reports 
		SET status = ?, admin_note = ?, reviewed_by = ?, reviewed_at = NOW(), updated_at = NOW()
		WHERE id = ?
	`
	_, err := r.db.ExecContext(ctx, query, status, adminNote, reviewedBy, id)
	if err != nil {
		return fmt.Errorf("failed to update report status: %w", err)
	}
	return nil
}

// HasUserReportedJob checks if user has already reported a job
func (r *Repository) HasUserReportedJob(ctx context.Context, userID, jobID uint64) (bool, error) {
	query := `SELECT COUNT(*) FROM job_reports WHERE user_id = ? AND job_id = ?`
	var count int
	err := r.db.QueryRowContext(ctx, query, userID, jobID).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check if user has reported job: %w", err)
	}
	return count > 0, nil
}

// GetLastReportTime gets the last time user reported any job
func (r *Repository) GetLastReportTime(ctx context.Context, userID uint64) (*time.Time, error) {
	query := `SELECT MAX(created_at) FROM job_reports WHERE user_id = ?`
	var lastTime sql.NullTime
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&lastTime)
	if err != nil {
		return nil, fmt.Errorf("failed to get last report time: %w", err)
	}
	if !lastTime.Valid {
		return nil, nil
	}
	return &lastTime.Time, nil
}

// GetReportCountByJobID gets the number of reports for a job
func (r *Repository) GetReportCountByJobID(ctx context.Context, jobID uint64) (int, error) {
	query := `SELECT COUNT(*) FROM job_reports WHERE job_id = ?`
	var count int
	err := r.db.QueryRowContext(ctx, query, jobID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get report count: %w", err)
	}
	return count, nil
}

// GetReportCountByCompanyID gets the number of reports for all jobs from a company
func (r *Repository) GetReportCountByCompanyID(ctx context.Context, companyID uint64) (int, error) {
	query := `
		SELECT COUNT(*) 
		FROM job_reports jr
		JOIN jobs j ON jr.job_id = j.id
		WHERE j.company_id = ?
	`
	var count int
	err := r.db.QueryRowContext(ctx, query, companyID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get report count by company: %w", err)
	}
	return count, nil
}

// GetPendingReportsCount gets the count of pending reports
func (r *Repository) GetPendingReportsCount(ctx context.Context) (int, error) {
	query := `SELECT COUNT(*) FROM job_reports WHERE status = ?`
	var count int
	err := r.db.QueryRowContext(ctx, query, ReportStatusPending).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get pending reports count: %w", err)
	}
	return count, nil
}

// BanCompany bans a company with reason
func (r *Repository) BanCompany(ctx context.Context, companyID uint64, adminID uint64, reason string) error {
	// Start transaction
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback()

	// Update company status to banned
	companyQuery := `
		UPDATE companies 
		SET is_banned = TRUE, 
		    company_status = 'suspended',
		    banned_at = NOW(), 
		    banned_reason = ?, 
		    banned_by = ?,
		    updated_at = NOW()
		WHERE id = ?
	`
	_, err = tx.ExecContext(ctx, companyQuery, reason, adminID, companyID)
	if err != nil {
		return fmt.Errorf("failed to ban company: %w", err)
	}

	// Update all job reports related to this company's jobs to "action_taken"
	reportsQuery := `
		UPDATE job_reports jr
		JOIN jobs j ON jr.job_id = j.id
		SET jr.status = 'action_taken',
		    jr.admin_note = CONCAT(COALESCE(jr.admin_note, ''), ' [Perusahaan telah dibanned]'),
		    jr.reviewed_by = ?,
		    jr.reviewed_at = NOW(),
		    jr.updated_at = NOW()
		WHERE j.company_id = ? 
		  AND jr.status IN ('pending', 'reviewed')
	`
	_, err = tx.ExecContext(ctx, reportsQuery, adminID, companyID)
	if err != nil {
		return fmt.Errorf("failed to update job reports: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
