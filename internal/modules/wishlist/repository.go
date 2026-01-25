package wishlist

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

// Repository defines the wishlist repository interface
type Repository interface {
	Save(ctx context.Context, userID, jobID uint64) (*SavedJob, error)
	Remove(ctx context.Context, userID, jobID uint64) error
	GetByUserAndJob(ctx context.Context, userID, jobID uint64) (*SavedJob, error)
	ListByUser(ctx context.Context, userID uint64, params ListParams) ([]*SavedJob, int64, error)
	IsSaved(ctx context.Context, userID, jobID uint64) (bool, error)
	CountByUser(ctx context.Context, userID uint64) (int64, error)
}

type repository struct {
	db *sqlx.DB
}

// NewRepository creates a new wishlist repository
func NewRepository(db *sqlx.DB) Repository {
	return &repository{db: db}
}

// Save adds a job to user's wishlist
func (r *repository) Save(ctx context.Context, userID, jobID uint64) (*SavedJob, error) {
	query := `
		INSERT INTO saved_jobs (user_id, job_id, created_at)
		VALUES (?, ?, NOW())
		ON DUPLICATE KEY UPDATE created_at = created_at
	`

	result, err := r.db.ExecContext(ctx, query, userID, jobID)
	if err != nil {
		return nil, err
	}

	id, _ := result.LastInsertId()
	if id == 0 {
		// Was duplicate, get existing
		return r.GetByUserAndJob(ctx, userID, jobID)
	}

	return &SavedJob{
		ID:     uint64(id),
		UserID: userID,
		JobID:  jobID,
	}, nil
}

// Remove removes a job from user's wishlist
func (r *repository) Remove(ctx context.Context, userID, jobID uint64) error {
	query := `DELETE FROM saved_jobs WHERE user_id = ? AND job_id = ?`
	_, err := r.db.ExecContext(ctx, query, userID, jobID)
	return err
}

// GetByUserAndJob gets a saved job by user and job ID
func (r *repository) GetByUserAndJob(ctx context.Context, userID, jobID uint64) (*SavedJob, error) {
	query := `
		SELECT id, user_id, job_id, created_at
		FROM saved_jobs
		WHERE user_id = ? AND job_id = ?
	`

	var saved SavedJob
	err := r.db.GetContext(ctx, &saved, query, userID, jobID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &saved, nil
}

// ListByUser lists all saved jobs for a user with job details
func (r *repository) ListByUser(ctx context.Context, userID uint64, params ListParams) ([]*SavedJob, int64, error) {
	// Count total
	countQuery := `SELECT COUNT(*) FROM saved_jobs WHERE user_id = ?`
	var total int64
	if err := r.db.GetContext(ctx, &total, countQuery, userID); err != nil {
		return nil, 0, err
	}

	if total == 0 {
		return []*SavedJob{}, 0, nil
	}

	// List with job details
	query := `
		SELECT 
			sj.id, sj.user_id, sj.job_id, sj.created_at,
			j.id as "job.id", j.title as "job.title", j.slug as "job.slug",
			j.city as "job.city", j.province as "job.province", j.is_remote as "job.is_remote",
			j.job_type as "job.job_type", j.experience_level as "job.experience_level",
			j.salary_min as "job.salary_min", j.salary_max as "job.salary_max",
			j.status as "job.status", j.created_at as "job.created_at",
			c.id as "job.company.id", c.company_name as "job.company.name", 
			c.company_logo_url as "job.company.logo_url"
		FROM saved_jobs sj
		JOIN jobs j ON sj.job_id = j.id
		JOIN companies c ON j.company_id = c.id
		WHERE sj.user_id = ?
		ORDER BY sj.created_at DESC
		LIMIT ? OFFSET ?
	`

	offset := (params.Page - 1) * params.PerPage

	rows, err := r.db.QueryxContext(ctx, query, userID, params.PerPage, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var savedJobs []*SavedJob
	for rows.Next() {
		var sj SavedJob
		var job JobInfo
		var company CompanyInfo
		var salaryMin, salaryMax sql.NullInt64
		var logoURL sql.NullString

		err := rows.Scan(
			&sj.ID, &sj.UserID, &sj.JobID, &sj.CreatedAt,
			&job.ID, &job.Title, &job.Slug,
			&job.City, &job.Province, &job.IsRemote,
			&job.JobType, &job.ExperienceLevel,
			&salaryMin, &salaryMax,
			&job.Status, &job.CreatedAt,
			&company.ID, &company.Name, &logoURL,
		)
		if err != nil {
			return nil, 0, err
		}

		if salaryMin.Valid {
			job.SalaryMin = &salaryMin.Int64
		}
		if salaryMax.Valid {
			job.SalaryMax = &salaryMax.Int64
		}
		if logoURL.Valid {
			company.LogoURL = &logoURL.String
		}

		job.Company = &company
		sj.Job = &job
		savedJobs = append(savedJobs, &sj)
	}

	return savedJobs, total, nil
}

// IsSaved checks if a job is saved by the user
func (r *repository) IsSaved(ctx context.Context, userID, jobID uint64) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM saved_jobs WHERE user_id = ? AND job_id = ?)`
	var exists bool
	err := r.db.GetContext(ctx, &exists, query, userID, jobID)
	return exists, err
}

// CountByUser returns the count of saved jobs for a user
func (r *repository) CountByUser(ctx context.Context, userID uint64) (int64, error) {
	query := `SELECT COUNT(*) FROM saved_jobs WHERE user_id = ?`
	var count int64
	err := r.db.GetContext(ctx, &count, query, userID)
	return count, err
}
