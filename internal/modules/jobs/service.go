package jobs

import (
	"context"
	"database/sql"
	"fmt"
	"regexp"
	"strings"
	"time"

	apperrors "github.com/karirnusantara/api/internal/shared/errors"
	"github.com/karirnusantara/api/internal/modules/company"
	"github.com/karirnusantara/api/internal/modules/quota"
)

// Service defines the jobs service interface
type Service interface {
	Create(ctx context.Context, companyID uint64, userID uint64, req *CreateJobRequest) (*JobResponse, error)
	GetByID(ctx context.Context, id uint64) (*JobResponse, error)
	GetBySlug(ctx context.Context, slug string) (*JobResponse, error)
	Update(ctx context.Context, id uint64, companyID uint64, req *UpdateJobRequest) (*JobResponse, error)
	UpdateStatus(ctx context.Context, id uint64, companyID uint64, userID uint64, status string) (*JobResponse, error)
	Delete(ctx context.Context, id uint64, companyID uint64) error
	List(ctx context.Context, params JobListParams) ([]*JobResponse, int64, error)
	ListByCompany(ctx context.Context, companyID uint64, params JobListParams) ([]*JobResponse, int64, error)
	IncrementViewCount(ctx context.Context, id uint64) error
	GetCompanyByUserID(ctx context.Context, userID uint64) (*company.Company, error)
}

type service struct {
	repo            Repository
	companyRepo     company.Repository
	quotaService    *quota.Service
}

// NewService creates a new jobs service
func NewService(repo Repository) Service {
	return &service{repo: repo}
}

// NewServiceWithCompanyRepo creates a new jobs service with company repository
func NewServiceWithCompanyRepo(repo Repository, companyRepo company.Repository) Service {
	return &service{repo: repo, companyRepo: companyRepo}
}

// NewServiceWithQuota creates a new jobs service with company and quota repositories
func NewServiceWithQuota(repo Repository, companyRepo company.Repository, quotaService *quota.Service) Service {
	return &service{repo: repo, companyRepo: companyRepo, quotaService: quotaService}
}

// GetCompanyByUserID retrieves company information for a given user ID
func (s *service) GetCompanyByUserID(ctx context.Context, userID uint64) (*company.Company, error) {
	if s.companyRepo == nil {
		return nil, apperrors.NewInternalError("Company repository not available", nil)
	}
	
	company, err := s.companyRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, apperrors.NewInternalError("Failed to get company", err)
	}
	if company == nil {
		return nil, nil
	}
	
	return company, nil
}

// Create creates a new job posting
func (s *service) Create(ctx context.Context, companyID uint64, userID uint64, req *CreateJobRequest) (*JobResponse, error) {
	// Validate company eligibility to create jobs
	if s.companyRepo != nil {
		canCreate, validationErr, err := s.companyRepo.CanCreateJobs(ctx, companyID)
		if err != nil {
			return nil, apperrors.NewInternalError("Failed to validate company", err)
		}
		if validationErr != nil {
			details := map[string]string{
				"code": validationErr.Code,
			}
			if validationErr.Details != "" {
				details["details"] = validationErr.Details
			}
			return nil, apperrors.NewValidationError(validationErr.Message, details)
		}
		if !canCreate {
			return nil, apperrors.NewValidationError("Company not eligible for job posting", nil)
		}
	}

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
		// Check and consume quota if publishing directly
		if s.quotaService != nil {
			canPublish, quotaType, err := s.quotaService.CanPublishJob(userID)
			if err != nil {
				return nil, apperrors.NewInternalError("Failed to check quota", err)
			}
			if !canPublish {
				return nil, apperrors.NewValidationError("Kuota posting habis. Silakan beli kuota tambahan untuk melanjutkan.", map[string]string{
					"code": "QUOTA_EXHAUSTED",
					"details": "Kuota gratis 10 post sudah habis. Harga per posting: Rp 15.000",
				})
			}
			// Consume quota
			if err := s.quotaService.ConsumeQuota(userID); err != nil {
				return nil, apperrors.NewInternalError("Failed to consume quota", err)
			}
			_ = quotaType // Log or track quota type used
		}
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

// ListByCompany lists jobs for a specific company
func (s *service) ListByCompany(ctx context.Context, companyID uint64, params JobListParams) ([]*JobResponse, int64, error) {
	params.CompanyID = &companyID
	jobs, total, err := s.repo.ListByCompany(ctx, companyID, params)
	if err != nil {
		return nil, 0, apperrors.NewInternalError("Failed to list company jobs", err)
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
func (s *service) UpdateStatus(ctx context.Context, id uint64, companyID uint64, userID uint64, newStatus string) (*JobResponse, error) {
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

	// Check and consume quota when publishing (draft -> active or first time active)
	if newStatus == JobStatusActive && !job.PublishedAt.Valid {
		if s.quotaService != nil {
			canPublish, _, err := s.quotaService.CanPublishJob(userID)
			if err != nil {
				return nil, apperrors.NewInternalError("Failed to check quota", err)
			}
			if !canPublish {
				return nil, apperrors.NewValidationError("Kuota posting habis. Silakan beli kuota tambahan untuk melanjutkan.", map[string]string{
					"code": "QUOTA_EXHAUSTED",
					"details": "Kuota gratis 10 post sudah habis. Harga per posting: Rp 15.000",
				})
			}
			// Consume quota
			if err := s.quotaService.ConsumeQuota(userID); err != nil {
				return nil, apperrors.NewInternalError("Failed to consume quota", err)
			}
		}
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
