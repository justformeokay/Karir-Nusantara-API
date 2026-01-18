package jobs

import (
	"context"
	"database/sql"
	"fmt"
	"regexp"
	"strings"
	"time"

	apperrors "github.com/karirnusantara/api/internal/shared/errors"
)

// Service defines the jobs service interface
type Service interface {
	Create(ctx context.Context, companyID uint64, req *CreateJobRequest) (*JobResponse, error)
	GetByID(ctx context.Context, id uint64) (*JobResponse, error)
	GetBySlug(ctx context.Context, slug string) (*JobResponse, error)
	Update(ctx context.Context, id uint64, companyID uint64, req *UpdateJobRequest) (*JobResponse, error)
	UpdateStatus(ctx context.Context, id uint64, companyID uint64, status string) (*JobResponse, error)
	Delete(ctx context.Context, id uint64, companyID uint64) error
	List(ctx context.Context, params JobListParams) ([]*JobResponse, int64, error)
	IncrementViewCount(ctx context.Context, id uint64) error
}

type service struct {
	repo Repository
}

// NewService creates a new jobs service
func NewService(repo Repository) Service {
	return &service{repo: repo}
}

// Create creates a new job posting
func (s *service) Create(ctx context.Context, companyID uint64, req *CreateJobRequest) (*JobResponse, error) {
	// Generate slug from title
	slug := generateSlug(req.Title)

	// Check if slug exists, append timestamp if needed
	existing, err := s.repo.GetBySlug(ctx, slug)
	if err != nil {
		return nil, apperrors.NewInternalError("Failed to check slug", err)
	}
	if existing != nil {
		slug = fmt.Sprintf("%s-%d", slug, time.Now().Unix())
	}

	// Build job entity
	job := &Job{
		CompanyID:       companyID,
		Title:           req.Title,
		Slug:            slug,
		Description:     req.Description,
		City:            req.City,
		Province:        req.Province,
		IsRemote:        req.IsRemote,
		JobType:         req.JobType,
		ExperienceLevel: req.ExperienceLevel,
		SalaryCurrency:  "IDR",
		IsSalaryVisible: req.IsSalaryVisible,
		Status:          JobStatusDraft,
	}

	// Optional fields
	if req.Requirements != "" {
		job.Requirements = sql.NullString{String: req.Requirements, Valid: true}
	}
	if req.Responsibilities != "" {
		job.Responsibilities = sql.NullString{String: req.Responsibilities, Valid: true}
	}
	if req.Benefits != "" {
		job.Benefits = sql.NullString{String: req.Benefits, Valid: true}
	}
	if req.SalaryMin != nil {
		job.SalaryMin = sql.NullInt64{Int64: *req.SalaryMin, Valid: true}
	}
	if req.SalaryMax != nil {
		job.SalaryMax = sql.NullInt64{Int64: *req.SalaryMax, Valid: true}
	}
	if req.SalaryCurrency != "" {
		job.SalaryCurrency = req.SalaryCurrency
	}
	if req.ApplicationDeadline != "" {
		deadline, err := time.Parse("2006-01-02", req.ApplicationDeadline)
		if err == nil {
			job.ApplicationDeadline = sql.NullTime{Time: deadline, Valid: true}
		}
	}

	// Set status and published_at
	if req.Status == JobStatusActive {
		job.Status = JobStatusActive
		job.PublishedAt = sql.NullTime{Time: time.Now(), Valid: true}
	}

	// Create job
	if err := s.repo.Create(ctx, job); err != nil {
		return nil, apperrors.NewInternalError("Failed to create job", err)
	}

	// Add skills
	if len(req.Skills) > 0 {
		if err := s.repo.AddSkills(ctx, job.ID, req.Skills); err != nil {
			return nil, apperrors.NewInternalError("Failed to add skills", err)
		}
	}

	return s.GetByID(ctx, job.ID)
}

// GetByID retrieves a job by ID with all related data
func (s *service) GetByID(ctx context.Context, id uint64) (*JobResponse, error) {
	job, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, apperrors.NewInternalError("Failed to get job", err)
	}
	if job == nil {
		return nil, apperrors.NewNotFoundError("Job")
	}

	// Load related data
	if err := s.loadJobRelations(ctx, job); err != nil {
		return nil, err
	}

	return job.ToResponse(), nil
}

// GetBySlug retrieves a job by slug
func (s *service) GetBySlug(ctx context.Context, slug string) (*JobResponse, error) {
	job, err := s.repo.GetBySlug(ctx, slug)
	if err != nil {
		return nil, apperrors.NewInternalError("Failed to get job", err)
	}
	if job == nil {
		return nil, apperrors.NewNotFoundError("Job")
	}

	// Load related data
	if err := s.loadJobRelations(ctx, job); err != nil {
		return nil, err
	}

	return job.ToResponse(), nil
}

// Update updates a job posting
func (s *service) Update(ctx context.Context, id uint64, companyID uint64, req *UpdateJobRequest) (*JobResponse, error) {
	// Get existing job
	job, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, apperrors.NewInternalError("Failed to get job", err)
	}
	if job == nil {
		return nil, apperrors.NewNotFoundError("Job")
	}

	// Check ownership
	if job.CompanyID != companyID {
		return nil, apperrors.NewForbiddenError("You don't have permission to update this job")
	}

	// Update fields
	if req.Title != nil {
		job.Title = *req.Title
		job.Slug = generateSlug(*req.Title)
	}
	if req.Description != nil {
		job.Description = *req.Description
	}
	if req.Requirements != nil {
		job.Requirements = sql.NullString{String: *req.Requirements, Valid: true}
	}
	if req.Responsibilities != nil {
		job.Responsibilities = sql.NullString{String: *req.Responsibilities, Valid: true}
	}
	if req.Benefits != nil {
		job.Benefits = sql.NullString{String: *req.Benefits, Valid: true}
	}
	if req.City != nil {
		job.City = *req.City
	}
	if req.Province != nil {
		job.Province = *req.Province
	}
	if req.IsRemote != nil {
		job.IsRemote = *req.IsRemote
	}
	if req.JobType != nil {
		job.JobType = *req.JobType
	}
	if req.ExperienceLevel != nil {
		job.ExperienceLevel = *req.ExperienceLevel
	}
	if req.SalaryMin != nil {
		job.SalaryMin = sql.NullInt64{Int64: *req.SalaryMin, Valid: true}
	}
	if req.SalaryMax != nil {
		job.SalaryMax = sql.NullInt64{Int64: *req.SalaryMax, Valid: true}
	}
	if req.IsSalaryVisible != nil {
		job.IsSalaryVisible = *req.IsSalaryVisible
	}
	if req.ApplicationDeadline != nil {
		deadline, err := time.Parse("2006-01-02", *req.ApplicationDeadline)
		if err == nil {
			job.ApplicationDeadline = sql.NullTime{Time: deadline, Valid: true}
		}
	}
	if req.Status != nil {
		// Handle status transition
		if *req.Status == JobStatusActive && job.Status != JobStatusActive {
			job.PublishedAt = sql.NullTime{Time: time.Now(), Valid: true}
		}
		job.Status = *req.Status
	}

	// Update job
	if err := s.repo.Update(ctx, job); err != nil {
		return nil, apperrors.NewInternalError("Failed to update job", err)
	}

	// Update skills if provided
	if req.Skills != nil {
		if err := s.repo.DeleteSkills(ctx, job.ID); err != nil {
			return nil, apperrors.NewInternalError("Failed to update skills", err)
		}
		if len(req.Skills) > 0 {
			if err := s.repo.AddSkills(ctx, job.ID, req.Skills); err != nil {
				return nil, apperrors.NewInternalError("Failed to add skills", err)
			}
		}
	}

	return s.GetByID(ctx, job.ID)
}

// Delete deletes a job posting
func (s *service) Delete(ctx context.Context, id uint64, companyID uint64) error {
	// Get existing job
	job, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return apperrors.NewInternalError("Failed to get job", err)
	}
	if job == nil {
		return apperrors.NewNotFoundError("Job")
	}

	// Check ownership
	if job.CompanyID != companyID {
		return apperrors.NewForbiddenError("You don't have permission to delete this job")
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return apperrors.NewInternalError("Failed to delete job", err)
	}

	return nil
}

// List retrieves jobs with filtering and pagination
func (s *service) List(ctx context.Context, params JobListParams) ([]*JobResponse, int64, error) {
	jobs, total, err := s.repo.List(ctx, params)
	if err != nil {
		return nil, 0, apperrors.NewInternalError("Failed to list jobs", err)
	}

	// Convert to response
	responses := make([]*JobResponse, len(jobs))
	for i, job := range jobs {
		// Load relations for each job
		if err := s.loadJobRelations(ctx, job); err != nil {
			return nil, 0, err
		}
		responses[i] = job.ToResponse()
	}

	return responses, total, nil
}

// IncrementViewCount increments the view count for a job
func (s *service) IncrementViewCount(ctx context.Context, id uint64) error {
	return s.repo.IncrementViewCount(ctx, id)
}

// UpdateStatus updates the job status (publish, close, pause, reopen)
func (s *service) UpdateStatus(ctx context.Context, id uint64, companyID uint64, newStatus string) (*JobResponse, error) {
	// Get existing job
	job, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, apperrors.NewInternalError("Failed to get job", err)
	}
	if job == nil {
		return nil, apperrors.NewNotFoundError("Job")
	}

	// Verify ownership
	if job.CompanyID != companyID {
		return nil, apperrors.NewForbiddenError("You don't have permission to modify this job")
	}

	// Validate status transition
	if !isValidStatusTransition(job.Status, newStatus) {
		return nil, apperrors.NewBadRequestError(fmt.Sprintf("Cannot change status from '%s' to '%s'", job.Status, newStatus))
	}

	// Update status
	job.Status = newStatus

	// Set published_at when publishing
	if newStatus == JobStatusActive && !job.PublishedAt.Valid {
		job.PublishedAt = sql.NullTime{Time: time.Now(), Valid: true}
	}

	if err := s.repo.Update(ctx, job); err != nil {
		return nil, apperrors.NewInternalError("Failed to update job status", err)
	}

	return s.GetByID(ctx, id)
}

// isValidStatusTransition checks if the status transition is allowed
func isValidStatusTransition(from, to string) bool {
	validTransitions := map[string][]string{
		JobStatusDraft:  {JobStatusActive},                          // draft -> active (publish)
		JobStatusActive: {JobStatusPaused, JobStatusClosed},         // active -> paused/closed
		JobStatusPaused: {JobStatusActive, JobStatusClosed},         // paused -> active (reopen) or closed
		JobStatusClosed: {JobStatusActive},                          // closed -> active (reopen)
	}

	allowed, ok := validTransitions[from]
	if !ok {
		return false
	}

	for _, s := range allowed {
		if s == to {
			return true
		}
	}
	return false
}

// loadJobRelations loads related data for a job
func (s *service) loadJobRelations(ctx context.Context, job *Job) error {
	// Load company info
	company, err := s.repo.GetCompanyInfo(ctx, job.CompanyID)
	if err != nil {
		return apperrors.NewInternalError("Failed to load company info", err)
	}
	job.Company = company

	// Load skills
	skills, err := s.repo.GetSkills(ctx, job.ID)
	if err != nil {
		return apperrors.NewInternalError("Failed to load skills", err)
	}
	job.Skills = skills

	return nil
}

// generateSlug generates a URL-friendly slug from a title
func generateSlug(title string) string {
	// Convert to lowercase
	slug := strings.ToLower(title)

	// Replace spaces with hyphens
	slug = strings.ReplaceAll(slug, " ", "-")

	// Remove non-alphanumeric characters (except hyphens)
	reg := regexp.MustCompile("[^a-z0-9-]+")
	slug = reg.ReplaceAllString(slug, "")

	// Remove multiple consecutive hyphens
	reg = regexp.MustCompile("-+")
	slug = reg.ReplaceAllString(slug, "-")

	// Trim hyphens from ends
	slug = strings.Trim(slug, "-")

	return slug
}
