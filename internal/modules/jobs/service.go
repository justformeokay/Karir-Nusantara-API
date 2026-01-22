package jobs

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	apperrors "github.com/karirnusantara/api/internal/shared/errors"
	"github.com/karirnusantara/api/internal/modules/company"
	"github.com/karirnusantara/api/internal/modules/quota"
	"github.com/karirnusantara/api/internal/shared/email"
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
	emailService    *email.Service
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

// NewServiceWithEmail creates a new jobs service with email notification support
func NewServiceWithEmail(repo Repository, companyRepo company.Repository, quotaService *quota.Service, emailService *email.Service) Service {
	return &service{
		repo:         repo,
		companyRepo:  companyRepo,
		quotaService: quotaService,
		emailService: emailService,
	}
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

	// Send email notification if job is published
	if job.Status == JobStatusActive && s.emailService != nil {
		// Use background context for goroutine since request context will be cancelled
		go s.sendJobPostedNotification(context.Background(), job.ID, companyID, userID)
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

// sendJobPostedNotification sends email notification when a job is posted
func (s *service) sendJobPostedNotification(ctx context.Context, jobID uint64, companyID uint64, userID uint64) {
	// Get job details
	job, err := s.GetByID(ctx, jobID)
	if err != nil {
		log.Printf("Failed to get job details for email notification: %v", err)
		return
	}

	// Get company details to get email
	company, err := s.companyRepo.GetByID(ctx, companyID)
	if err != nil {
		log.Printf("Failed to get company details for email notification: %v", err)
		return
	}

	// Get company name and email safely
	companyName := "Perusahaan Anda"
	if company.CompanyName.Valid {
		companyName = company.CompanyName.String
	}

	companyEmail := ""
	if company.CompanyEmail.Valid {
		companyEmail = company.CompanyEmail.String
	}

	if companyEmail == "" {
		log.Printf("Company email is empty for company #%d, skipping notification", companyID)
		return
	}

	// Format location
	location := job.Location.City
	if job.Location.Province != "" {
		location = fmt.Sprintf("%s, %s", job.Location.City, job.Location.Province)
	}
	remoteText := ""
	if job.Location.IsRemote {
		remoteText = "(Remote)"
	}

	// Format job type
	jobTypeMap := map[string]string{
		"full-time":  "Full Time",
		"part-time":  "Part Time",
		"contract":   "Kontrak",
		"internship": "Magang",
		"freelance":  "Freelance",
	}
	jobTypeText := jobTypeMap[job.JobType]
	if jobTypeText == "" {
		jobTypeText = job.JobType
	}

	// Format experience level
	levelMap := map[string]string{
		"entry":    "Entry Level",
		"junior":   "Junior",
		"mid":      "Mid Level",
		"senior":   "Senior",
		"lead":     "Lead",
		"manager":  "Manager",
		"director": "Director",
	}
	levelText := levelMap[job.ExperienceLevel]
	if levelText == "" {
		levelText = job.ExperienceLevel
	}

	// Format published date
	publishedDate := "Baru saja"
	if job.PublishedAt != "" {
		publishedDate = job.PublishedAt
	}

	// Build email content
	subject := fmt.Sprintf("✅ Lowongan '%s' Berhasil Dipublikasikan", job.Title)

	body := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background: linear-gradient(135deg, #2980b9 0%%, #6dd5fa 100%%); color: white; padding: 30px; text-align: center; border-radius: 8px 8px 0 0; }
        .content { background: #f9f9f9; padding: 30px; border: 1px solid #e0e0e0; }
        .job-details { background: white; padding: 20px; margin: 20px 0; border-left: 4px solid #2980b9; border-radius: 4px; }
        .job-details h3 { margin-top: 0; color: #2980b9; }
        .detail-row { margin: 10px 0; }
        .detail-label { font-weight: bold; color: #555; }
        .button { display: inline-block; padding: 12px 30px; background: #2980b9; color: white; text-decoration: none; border-radius: 5px; margin: 20px 0; }
        .footer { text-align: center; padding: 20px; color: #777; font-size: 12px; }
        .alert { background: #d4edda; border: 1px solid #c3e6cb; color: #155724; padding: 15px; border-radius: 4px; margin: 15px 0; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>🎉 Lowongan Berhasil Dipublikasikan!</h1>
        </div>
        
        <div class="content">
            <p>Halo <strong>%s</strong>,</p>
            
            <p>Selamat! Lowongan pekerjaan Anda telah berhasil dipublikasikan di platform <strong>Karir Nusantara</strong> dan sekarang dapat dilihat oleh ribuan pencari kerja di seluruh Indonesia.</p>
            
            <div class="job-details">
                <h3>📋 Detail Lowongan</h3>
                <div class="detail-row">
                    <span class="detail-label">Posisi:</span> %s
                </div>
                <div class="detail-row">
                    <span class="detail-label">Lokasi:</span> %s %s
                </div>
                <div class="detail-row">
                    <span class="detail-label">Tipe Pekerjaan:</span> %s
                </div>
                <div class="detail-row">
                    <span class="detail-label">Level:</span> %s
                </div>
                <div class="detail-row">
                    <span class="detail-label">Status:</span> <span style="color: #28a745; font-weight: bold;">AKTIF</span>
                </div>
                <div class="detail-row">
                    <span class="detail-label">Tanggal Publish:</span> %s
                </div>
            </div>
            
            <div class="alert">
                <strong>🔒 Keamanan & Privasi</strong><br>
                Email ini dikirim untuk mengkonfirmasi aktivitas posting lowongan dari akun Anda. 
                Jika Anda tidak melakukan posting ini, segera hubungi tim support kami melalui chat support di dashboard.
            </div>
            
            <p><strong>Apa yang terjadi selanjutnya?</strong></p>
            <ul>
                <li>Lowongan Anda kini dapat dilihat oleh pencari kerja</li>
                <li>Anda akan menerima notifikasi saat ada lamaran masuk</li>
                <li>Anda dapat mengelola lowongan di menu Dashboard > Lowongan</li>
                <li>Statistik viewing akan diupdate secara real-time</li>
            </ul>
            
            <div style="text-align: center;">
                <a href="http://localhost:5174/dashboard/jobs" class="button">Kelola Lowongan Saya</a>
            </div>
            
            <p style="margin-top: 30px; padding-top: 20px; border-top: 1px solid #ddd;">
                <strong>Butuh bantuan?</strong><br>
                Hubungi kami melalui fitur Chat Support di dashboard atau email ke support@karirnusantara.com
            </p>
        </div>
        
        <div class="footer">
            <p>Karir Nusantara - Platform Pencarian Kerja Terpercaya Indonesia</p>
            <p>Email ini dikirim otomatis oleh sistem. Mohon tidak membalas email ini.</p>
            <p>&copy; 2024 Karir Nusantara. All rights reserved.</p>
        </div>
    </div>
</body>
</html>
	`,
		companyName,
		job.Title,
		location,
		remoteText,
		jobTypeText,
		levelText,
		publishedDate,
	)

	// Send email
	log.Printf("[JOB NOTIFICATION] Attempting to send email to: %s for job #%d (%s)", companyEmail, jobID, job.Title)
	err = s.emailService.SendEmail(companyEmail, subject, body)
	if err != nil {
		log.Printf("[JOB NOTIFICATION ERROR] Failed to send job posted notification email to %s: %v", companyEmail, err)
		// Don't fail the job creation if email fails
	} else {
		log.Printf("[JOB NOTIFICATION SUCCESS] Job posted notification email sent to %s for job #%d", companyEmail, jobID)
	}
}
