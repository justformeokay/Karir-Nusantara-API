package admin

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"github.com/karirnusantara/api/internal/config"
)

// Common errors
var (
	ErrInvalidCredentials = errors.New("email atau password salah")
	ErrAccountInactive    = errors.New("akun tidak aktif")
	ErrAdminNotFound      = errors.New("admin tidak ditemukan")
	ErrCompanyNotFound    = errors.New("perusahaan tidak ditemukan")
	ErrJobNotFound        = errors.New("lowongan tidak ditemukan")
	ErrPaymentNotFound    = errors.New("pembayaran tidak ditemukan")
	ErrJobSeekerNotFound  = errors.New("pencari kerja tidak ditemukan")
	ErrInvalidAction      = errors.New("aksi tidak valid")
)

// Service defines the admin service interface
type Service interface {
	// Authentication
	Login(ctx context.Context, req *AdminLoginRequest) (*AdminAuthResponse, error)
	GetCurrentAdmin(ctx context.Context, adminID uint64) (*AdminUserResponse, error)

	// Dashboard
	GetDashboardStats(ctx context.Context) (*DashboardStats, error)

	// Company management
	GetCompanies(ctx context.Context, filter CompanyFilter) (*PaginatedResponse, error)
	GetCompanyByID(ctx context.Context, id uint64) (*CompanyAdminResponse, error)
	VerifyCompany(ctx context.Context, id uint64, req *CompanyVerificationRequest, adminID uint64) error
	UpdateCompanyStatus(ctx context.Context, id uint64, req *CompanyStatusRequest, adminID uint64) error

	// Job management
	GetJobs(ctx context.Context, filter JobFilter) (*PaginatedResponse, error)
	GetJobByID(ctx context.Context, id uint64) (*JobAdminResponse, error)
	ModerateJob(ctx context.Context, id uint64, req *JobActionRequest, adminID uint64) error

	// Payment management
	GetPayments(ctx context.Context, filter PaymentFilter) (*PaginatedResponse, error)
	GetPaymentByID(ctx context.Context, id uint64) (*PaymentAdminResponse, error)
	ProcessPayment(ctx context.Context, id uint64, req *PaymentActionRequest, adminID uint64) error

	// Job seeker management
	GetJobSeekers(ctx context.Context, filter JobSeekerFilter) (*PaginatedResponse, error)
	GetJobSeekerByID(ctx context.Context, id uint64) (*JobSeekerAdminResponse, error)
	UpdateJobSeekerStatus(ctx context.Context, id uint64, req *JobSeekerActionRequest, adminID uint64) error
}

// service implements the Service interface
type service struct {
	repo   Repository
	config *config.Config
}

// NewService creates a new admin service
func NewService(repo Repository, cfg *config.Config) Service {
	return &service{
		repo:   repo,
		config: cfg,
	}
}

// ============================================
// AUTHENTICATION
// ============================================

func (s *service) Login(ctx context.Context, req *AdminLoginRequest) (*AdminAuthResponse, error) {
	// Find admin by email
	admin, err := s.repo.GetAdminByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to find admin: %w", err)
	}
	if admin == nil {
		return nil, ErrInvalidCredentials
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(admin.PasswordHash), []byte(req.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	// Check if active
	if !admin.IsActive {
		return nil, ErrAccountInactive
	}

	// Generate JWT token
	expiresIn := time.Hour * 24 // 24 hours for admin tokens
	token, err := s.generateAccessToken(admin, expiresIn)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &AdminAuthResponse{
		Admin:       admin.ToResponse(),
		AccessToken: token,
		ExpiresIn:   int64(expiresIn.Seconds()),
	}, nil
}

func (s *service) GetCurrentAdmin(ctx context.Context, adminID uint64) (*AdminUserResponse, error) {
	admin, err := s.repo.GetAdminByID(ctx, adminID)
	if err != nil {
		return nil, fmt.Errorf("failed to get admin: %w", err)
	}
	if admin == nil {
		return nil, ErrAdminNotFound
	}

	return admin.ToResponse(), nil
}

// ============================================
// DASHBOARD
// ============================================

func (s *service) GetDashboardStats(ctx context.Context) (*DashboardStats, error) {
	stats, err := s.repo.GetDashboardStats(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get dashboard stats: %w", err)
	}
	return stats, nil
}

// ============================================
// COMPANY MANAGEMENT
// ============================================

func (s *service) GetCompanies(ctx context.Context, filter CompanyFilter) (*PaginatedResponse, error) {
	companies, total, err := s.repo.GetCompanies(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to get companies: %w", err)
	}

	// Convert to responses
	responses := make([]*CompanyAdminResponse, len(companies))
	for i, c := range companies {
		responses[i] = c.ToResponse()
	}

	return NewPaginatedResponse(responses, total, filter.Page, filter.PageSize), nil
}

func (s *service) GetCompanyByID(ctx context.Context, id uint64) (*CompanyAdminResponse, error) {
	company, err := s.repo.GetCompanyByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get company: %w", err)
	}
	if company == nil {
		return nil, ErrCompanyNotFound
	}

	return company.ToResponse(), nil
}

func (s *service) VerifyCompany(ctx context.Context, id uint64, req *CompanyVerificationRequest, adminID uint64) error {
	// Check if company exists
	company, err := s.repo.GetCompanyByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get company: %w", err)
	}
	if company == nil {
		return ErrCompanyNotFound
	}

	var newStatus string
	var action string

	switch req.Action {
	case "approve":
		newStatus = CompanyStatusVerified
		action = "company_verified"
	case "reject":
		newStatus = CompanyStatusRejected
		action = "company_rejected"
	default:
		return ErrInvalidAction
	}

	// Update status
	if err := s.repo.UpdateCompanyStatus(ctx, id, newStatus); err != nil {
		// Log the actual error for debugging
		fmt.Printf("DEBUG: UpdateCompanyStatus failed for id=%d, status=%s, error=%v\n", id, newStatus, err)
		return err
	}

	// Log admin action (optional)
	if req.Reason != "" {
		s.logAction(ctx, adminID, action, "company", id, req.Reason)
	}

	return nil
}

func (s *service) UpdateCompanyStatus(ctx context.Context, id uint64, req *CompanyStatusRequest, adminID uint64) error {
	// Check if company exists
	company, err := s.repo.GetCompanyByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get company: %w", err)
	}
	if company == nil {
		return ErrCompanyNotFound
	}

	var isActive bool
	var status string
	var action string

	switch req.Action {
	case "suspend":
		isActive = false
		status = CompanyStatusSuspended
		action = "company_suspended"
	case "reactivate":
		isActive = true
		status = CompanyStatusVerified
		action = "company_reactivated"
	default:
		return ErrInvalidAction
	}

	// Update active status
	if err := s.repo.UpdateCompanyActive(ctx, id, isActive); err != nil {
		return fmt.Errorf("failed to update company active status: %w", err)
	}

	// Update company_status
	if err := s.repo.UpdateCompanyStatus(ctx, id, status); err != nil {
		return fmt.Errorf("failed to update company status: %w", err)
	}

	// Log admin action
	s.logAction(ctx, adminID, action, "company", id, req.Reason)

	return nil
}

// ============================================
// JOB MANAGEMENT
// ============================================

func (s *service) GetJobs(ctx context.Context, filter JobFilter) (*PaginatedResponse, error) {
	jobs, total, err := s.repo.GetJobs(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to get jobs: %w", err)
	}

	// Convert to responses
	responses := make([]*JobAdminResponse, len(jobs))
	for i, j := range jobs {
		responses[i] = j.ToResponse()
	}

	return NewPaginatedResponse(responses, total, filter.Page, filter.PageSize), nil
}

func (s *service) GetJobByID(ctx context.Context, id uint64) (*JobAdminResponse, error) {
	job, err := s.repo.GetJobByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get job: %w", err)
	}
	if job == nil {
		return nil, ErrJobNotFound
	}

	return job.ToResponse(), nil
}

func (s *service) ModerateJob(ctx context.Context, id uint64, req *JobActionRequest, adminID uint64) error {
	// Check if job exists
	job, err := s.repo.GetJobByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get job: %w", err)
	}
	if job == nil {
		return ErrJobNotFound
	}

	var newStatus string
	var adminStatus string
	var action string

	switch req.Action {
	case "approve":
		newStatus = JobStatusActive
		adminStatus = "approved"
		action = "job_approved"
	case "reject":
		newStatus = JobStatusRejected
		adminStatus = "rejected"
		action = "job_rejected"
	case "close":
		newStatus = JobStatusClosed
		action = "job_closed"
	case "flag":
		adminStatus = JobStatusFlagged
		action = "job_flagged"
	case "unflag":
		adminStatus = ""
		action = "job_unflagged"
	default:
		return ErrInvalidAction
	}

	// Update job status if needed
	if newStatus != "" {
		if err := s.repo.UpdateJobStatus(ctx, id, newStatus); err != nil {
			return fmt.Errorf("failed to update job status: %w", err)
		}
	}

	// Update admin status/note
	if req.Action == "flag" || req.Action == "unflag" || req.Action == "approve" || req.Action == "reject" {
		if err := s.repo.UpdateJobAdminStatus(ctx, id, adminStatus, req.Reason); err != nil {
			return fmt.Errorf("failed to update job admin status: %w", err)
		}
	}

	// Log admin action
	s.logAction(ctx, adminID, action, "job", id, req.Reason)

	return nil
}

// ============================================
// PAYMENT MANAGEMENT
// ============================================

func (s *service) GetPayments(ctx context.Context, filter PaymentFilter) (*PaginatedResponse, error) {
	payments, total, err := s.repo.GetPayments(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to get payments: %w", err)
	}

	// Convert to responses
	responses := make([]*PaymentAdminResponse, len(payments))
	for i, p := range payments {
		responses[i] = p.ToResponse()
	}

	return NewPaginatedResponse(responses, total, filter.Page, filter.PageSize), nil
}

func (s *service) GetPaymentByID(ctx context.Context, id uint64) (*PaymentAdminResponse, error) {
	payment, err := s.repo.GetPaymentByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get payment: %w", err)
	}
	if payment == nil {
		return nil, ErrPaymentNotFound
	}

	return payment.ToResponse(), nil
}

func (s *service) ProcessPayment(ctx context.Context, id uint64, req *PaymentActionRequest, adminID uint64) error {
	// Check if payment exists
	payment, err := s.repo.GetPaymentByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get payment: %w", err)
	}
	if payment == nil {
		return ErrPaymentNotFound
	}

	var newStatus string
	var action string

	switch req.Action {
	case "approve":
		newStatus = PaymentStatusConfirmed
		action = "payment_approved"
	case "reject":
		newStatus = PaymentStatusRejected
		action = "payment_rejected"
	default:
		return ErrInvalidAction
	}

	// Update payment status
	if err := s.repo.UpdatePaymentStatus(ctx, id, newStatus, req.Note, adminID); err != nil {
		return fmt.Errorf("failed to update payment status: %w", err)
	}

	// Log admin action
	s.logAction(ctx, adminID, action, "payment", id, req.Note)

	return nil
}

// ============================================
// JOB SEEKER MANAGEMENT
// ============================================

func (s *service) GetJobSeekers(ctx context.Context, filter JobSeekerFilter) (*PaginatedResponse, error) {
	jobSeekers, total, err := s.repo.GetJobSeekers(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to get job seekers: %w", err)
	}

	// Convert to responses
	responses := make([]*JobSeekerAdminResponse, len(jobSeekers))
	for i, js := range jobSeekers {
		responses[i] = js.ToResponse()
	}

	return NewPaginatedResponse(responses, total, filter.Page, filter.PageSize), nil
}

func (s *service) GetJobSeekerByID(ctx context.Context, id uint64) (*JobSeekerAdminResponse, error) {
	jobSeeker, err := s.repo.GetJobSeekerByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get job seeker: %w", err)
	}
	if jobSeeker == nil {
		return nil, ErrJobSeekerNotFound
	}

	return jobSeeker.ToResponse(), nil
}

func (s *service) UpdateJobSeekerStatus(ctx context.Context, id uint64, req *JobSeekerActionRequest, adminID uint64) error {
	// Check if job seeker exists
	jobSeeker, err := s.repo.GetJobSeekerByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get job seeker: %w", err)
	}
	if jobSeeker == nil {
		return ErrJobSeekerNotFound
	}

	var isActive bool
	var action string

	switch req.Action {
	case "suspend", "deactivate":
		isActive = false
		action = "job_seeker_suspended"
	case "reactivate":
		isActive = true
		action = "job_seeker_reactivated"
	default:
		return ErrInvalidAction
	}

	// Update active status
	if err := s.repo.UpdateJobSeekerActive(ctx, id, isActive); err != nil {
		return fmt.Errorf("failed to update job seeker status: %w", err)
	}

	// Log admin action
	s.logAction(ctx, adminID, action, "job_seeker", id, req.Reason)

	return nil
}

// ============================================
// HELPERS
// ============================================

func (s *service) logAction(ctx context.Context, adminID uint64, action, entityType string, entityID uint64, details string) {
	log := &AdminActionLog{
		AdminID:    adminID,
		Action:     action,
		EntityType: entityType,
		EntityID:   entityID,
		Details:    sql.NullString{String: details, Valid: details != ""},
	}

	// We don't fail if logging fails, just log it
	_ = s.repo.LogAdminAction(ctx, log)
}

// generateAccessToken generates a new access token for admin
func (s *service) generateAccessToken(admin *AdminUser, expiry time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"user_id":    admin.ID,
		"email":      admin.Email,
		"role":       admin.Role,
		"token_type": "access",
		"exp":        time.Now().Add(expiry).Unix(),
		"iat":        time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.config.JWT.Secret))
}
