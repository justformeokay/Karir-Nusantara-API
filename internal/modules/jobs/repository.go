package jobs

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
)

// Repository defines the jobs repository interface
type Repository interface {
	Create(ctx context.Context, job *Job) error
	GetByID(ctx context.Context, id uint64) (*Job, error)
	GetBySlug(ctx context.Context, slug string) (*Job, error)
	Update(ctx context.Context, job *Job) error
	Delete(ctx context.Context, id uint64) error
	List(ctx context.Context, params JobListParams) ([]*Job, int64, error)
	IncrementViewCount(ctx context.Context, id uint64) error
	IncrementApplicationCount(ctx context.Context, id uint64) error
	
	// Skills
	AddSkills(ctx context.Context, jobID uint64, skills []string) error
	GetSkills(ctx context.Context, jobID uint64) ([]JobSkill, error)
	DeleteSkills(ctx context.Context, jobID uint64) error

	// Company info
	GetCompanyInfo(ctx context.Context, companyID uint64) (*CompanyInfo, error)
}

type mysqlRepository struct {
	db *sqlx.DB
}

// NewRepository creates a new jobs repository
func NewRepository(db *sqlx.DB) Repository {
	return &mysqlRepository{db: db}
}

// Create creates a new job
func (r *mysqlRepository) Create(ctx context.Context, job *Job) error {
	query := `
		INSERT INTO jobs (
			company_id, title, slug, description, requirements, responsibilities, benefits,
			city, province, is_remote, job_type, experience_level,
			salary_min, salary_max, salary_currency, is_salary_visible,
			application_deadline, max_applications, status, published_at,
			created_at, updated_at
		) VALUES (
			?, ?, ?, ?, ?, ?, ?,
			?, ?, ?, ?, ?,
			?, ?, ?, ?,
			?, ?, ?, ?,
			NOW(), NOW()
		)
	`

	result, err := r.db.ExecContext(ctx, query,
		job.CompanyID, job.Title, job.Slug, job.Description, job.Requirements, job.Responsibilities, job.Benefits,
		job.City, job.Province, job.IsRemote, job.JobType, job.ExperienceLevel,
		job.SalaryMin, job.SalaryMax, job.SalaryCurrency, job.IsSalaryVisible,
		job.ApplicationDeadline, job.MaxApplications, job.Status, job.PublishedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create job: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	job.ID = uint64(id)
	return nil
}

// GetByID retrieves a job by ID
func (r *mysqlRepository) GetByID(ctx context.Context, id uint64) (*Job, error) {
	query := `
		SELECT id, company_id, title, slug, description, requirements, responsibilities, benefits,
			   city, province, is_remote, job_type, experience_level,
			   salary_min, salary_max, salary_currency, is_salary_visible,
			   application_deadline, max_applications, status, views_count, applications_count,
			   published_at, created_at, updated_at, deleted_at
		FROM jobs
		WHERE id = ? AND deleted_at IS NULL
	`

	var job Job
	if err := r.db.GetContext(ctx, &job, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get job by id: %w", err)
	}

	return &job, nil
}

// GetBySlug retrieves a job by slug
func (r *mysqlRepository) GetBySlug(ctx context.Context, slug string) (*Job, error) {
	query := `
		SELECT id, company_id, title, slug, description, requirements, responsibilities, benefits,
			   city, province, is_remote, job_type, experience_level,
			   salary_min, salary_max, salary_currency, is_salary_visible,
			   application_deadline, max_applications, status, views_count, applications_count,
			   published_at, created_at, updated_at, deleted_at
		FROM jobs
		WHERE slug = ? AND deleted_at IS NULL
	`

	var job Job
	if err := r.db.GetContext(ctx, &job, query, slug); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get job by slug: %w", err)
	}

	return &job, nil
}

// Update updates a job
func (r *mysqlRepository) Update(ctx context.Context, job *Job) error {
	query := `
		UPDATE jobs SET
			title = ?, slug = ?, description = ?, requirements = ?, responsibilities = ?, benefits = ?,
			city = ?, province = ?, is_remote = ?, job_type = ?, experience_level = ?,
			salary_min = ?, salary_max = ?, is_salary_visible = ?,
			application_deadline = ?, status = ?, published_at = ?,
			updated_at = NOW()
		WHERE id = ? AND deleted_at IS NULL
	`

	_, err := r.db.ExecContext(ctx, query,
		job.Title, job.Slug, job.Description, job.Requirements, job.Responsibilities, job.Benefits,
		job.City, job.Province, job.IsRemote, job.JobType, job.ExperienceLevel,
		job.SalaryMin, job.SalaryMax, job.IsSalaryVisible,
		job.ApplicationDeadline, job.Status, job.PublishedAt,
		job.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update job: %w", err)
	}

	return nil
}

// Delete soft deletes a job
func (r *mysqlRepository) Delete(ctx context.Context, id uint64) error {
	query := `UPDATE jobs SET deleted_at = NOW() WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete job: %w", err)
	}
	return nil
}

// List retrieves jobs with filtering and pagination
func (r *mysqlRepository) List(ctx context.Context, params JobListParams) ([]*Job, int64, error) {
	// Build WHERE clause
	var conditions []string
	var args []interface{}

	conditions = append(conditions, "deleted_at IS NULL")

	if params.Status != "" {
		conditions = append(conditions, "status = ?")
		args = append(args, params.Status)
	}

	if params.Search != "" {
		conditions = append(conditions, "MATCH(title, description, requirements) AGAINST(? IN NATURAL LANGUAGE MODE)")
		args = append(args, params.Search)
	}

	if params.City != "" {
		conditions = append(conditions, "city = ?")
		args = append(args, params.City)
	}

	if params.Province != "" {
		conditions = append(conditions, "province = ?")
		args = append(args, params.Province)
	}

	if params.JobType != "" {
		conditions = append(conditions, "job_type = ?")
		args = append(args, params.JobType)
	}

	if params.ExperienceLevel != "" {
		conditions = append(conditions, "experience_level = ?")
		args = append(args, params.ExperienceLevel)
	}

	if params.IsRemote != nil {
		conditions = append(conditions, "is_remote = ?")
		args = append(args, *params.IsRemote)
	}

	if params.SalaryMin != nil {
		conditions = append(conditions, "salary_max >= ?")
		args = append(args, *params.SalaryMin)
	}

	if params.SalaryMax != nil {
		conditions = append(conditions, "salary_min <= ?")
		args = append(args, *params.SalaryMax)
	}

	if params.CompanyID != nil {
		conditions = append(conditions, "company_id = ?")
		args = append(args, *params.CompanyID)
	}

	whereClause := strings.Join(conditions, " AND ")

	// Count total
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM jobs WHERE %s", whereClause)
	var total int64
	if err := r.db.GetContext(ctx, &total, countQuery, args...); err != nil {
		return nil, 0, fmt.Errorf("failed to count jobs: %w", err)
	}

	// Build ORDER BY
	orderBy := "published_at DESC"
	if params.SortBy != "" {
		order := "ASC"
		if params.SortOrder == "desc" {
			order = "DESC"
		}
		orderBy = fmt.Sprintf("%s %s", params.SortBy, order)
	}

	// Build query with pagination
	offset := (params.Page - 1) * params.PerPage
	query := fmt.Sprintf(`
		SELECT id, company_id, title, slug, description, requirements, responsibilities, benefits,
			   city, province, is_remote, job_type, experience_level,
			   salary_min, salary_max, salary_currency, is_salary_visible,
			   application_deadline, max_applications, status, views_count, applications_count,
			   published_at, created_at, updated_at, deleted_at
		FROM jobs
		WHERE %s
		ORDER BY %s
		LIMIT ? OFFSET ?
	`, whereClause, orderBy)

	args = append(args, params.PerPage, offset)

	var jobs []*Job
	if err := r.db.SelectContext(ctx, &jobs, query, args...); err != nil {
		return nil, 0, fmt.Errorf("failed to list jobs: %w", err)
	}

	return jobs, total, nil
}

// IncrementViewCount increments the view count
func (r *mysqlRepository) IncrementViewCount(ctx context.Context, id uint64) error {
	query := `UPDATE jobs SET views_count = views_count + 1 WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

// IncrementApplicationCount increments the application count
func (r *mysqlRepository) IncrementApplicationCount(ctx context.Context, id uint64) error {
	query := `UPDATE jobs SET applications_count = applications_count + 1 WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

// AddSkills adds skills to a job
func (r *mysqlRepository) AddSkills(ctx context.Context, jobID uint64, skills []string) error {
	if len(skills) == 0 {
		return nil
	}

	query := `INSERT INTO job_skills (job_id, skill_name, is_required) VALUES (?, ?, TRUE)`
	for _, skill := range skills {
		if _, err := r.db.ExecContext(ctx, query, jobID, skill); err != nil {
			return fmt.Errorf("failed to add skill: %w", err)
		}
	}
	return nil
}

// GetSkills retrieves skills for a job
func (r *mysqlRepository) GetSkills(ctx context.Context, jobID uint64) ([]JobSkill, error) {
	query := `SELECT id, job_id, skill_name, is_required FROM job_skills WHERE job_id = ?`
	var skills []JobSkill
	if err := r.db.SelectContext(ctx, &skills, query, jobID); err != nil {
		return nil, fmt.Errorf("failed to get skills: %w", err)
	}
	return skills, nil
}

// DeleteSkills deletes all skills for a job
func (r *mysqlRepository) DeleteSkills(ctx context.Context, jobID uint64) error {
	query := `DELETE FROM job_skills WHERE job_id = ?`
	_, err := r.db.ExecContext(ctx, query, jobID)
	return err
}

// GetCompanyInfo retrieves company info for a job
func (r *mysqlRepository) GetCompanyInfo(ctx context.Context, companyID uint64) (*CompanyInfo, error) {
	query := `SELECT id, company_name, company_logo_url, company_website FROM users WHERE id = ?`
	
	var result struct {
		ID         uint64         `db:"id"`
		Name       sql.NullString `db:"company_name"`
		LogoURL    sql.NullString `db:"company_logo_url"`
		Website    sql.NullString `db:"company_website"`
	}
	
	if err := r.db.GetContext(ctx, &result, query, companyID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get company info: %w", err)
	}

	return &CompanyInfo{
		ID:      result.ID,
		Name:    result.Name.String,
		LogoURL: result.LogoURL.String,
		Website: result.Website.String,
	}, nil
}
