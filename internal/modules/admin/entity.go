package admin

import (
	"database/sql"
	"time"

	"github.com/karirnusantara/api/internal/shared/hashid"
)

// ============================================
// ADMIN-SPECIFIC ENTITIES
// ============================================

// AdminUser represents an admin user (extends base User with admin-specific fields)
type AdminUser struct {
	ID           uint64    `db:"id" json:"id"`
	Email        string    `db:"email" json:"email"`
	PasswordHash string    `db:"password_hash" json:"-"`
	FullName     string    `db:"full_name" json:"full_name"`
	Role         string    `db:"role" json:"role"`
	IsActive     bool      `db:"is_active" json:"is_active"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time `db:"updated_at" json:"updated_at"`
}

// SimpleUser represents a simple user for basic queries
type SimpleUser struct {
	ID       uint64 `db:"id"`
	Email    string `db:"email"`
	FullName string `db:"full_name"`
	Role     string `db:"role"`
}

// AdminUserResponse represents the admin user response
type AdminUserResponse struct {
	ID        uint64 `json:"id"`
	HashID    string `json:"hash_id"`
	Email     string `json:"email"`
	FullName  string `json:"full_name"`
	Role      string `json:"role"`
	IsActive  bool   `json:"is_active"`
	CreatedAt string `json:"created_at"`
}

func (u *AdminUser) ToResponse() *AdminUserResponse {
	return &AdminUserResponse{
		ID:        u.ID,
		HashID:    hashid.Encode(u.ID),
		Email:     u.Email,
		FullName:  u.FullName,
		Role:      u.Role,
		IsActive:  u.IsActive,
		CreatedAt: u.CreatedAt.Format(time.RFC3339),
	}
}

// ============================================
// COMPANY ENTITIES (ADMIN VIEW)
// ============================================

// Company statuses for admin verification workflow
const (
	CompanyStatusPending   = "pending"
	CompanyStatusVerified  = "verified"
	CompanyStatusRejected  = "rejected"
	CompanyStatusSuspended = "suspended"
)

// CompanyAdmin represents company data for admin view
type CompanyAdmin struct {
	ID                 uint64         `db:"id" json:"id"`
	Email              string         `db:"email" json:"email"`
	FullName           string         `db:"full_name" json:"full_name"`
	Phone              sql.NullString `db:"phone" json:"phone,omitempty"`
	CompanyName        sql.NullString `db:"company_name" json:"company_name"`
	CompanyDescription sql.NullString `db:"company_description" json:"company_description,omitempty"`
	CompanyWebsite     sql.NullString `db:"company_website" json:"company_website,omitempty"`
	CompanyLogoURL     sql.NullString `db:"company_logo_url" json:"company_logo_url,omitempty"`
	CompanyStatus      sql.NullString `db:"company_status" json:"company_status"`
	IsActive           bool           `db:"is_active" json:"is_active"`
	IsVerified         bool           `db:"is_verified" json:"is_verified"`
	EmailVerifiedAt    sql.NullTime   `db:"email_verified_at" json:"email_verified_at,omitempty"`
	CreatedAt          time.Time      `db:"created_at" json:"created_at"`
	UpdatedAt          time.Time      `db:"updated_at" json:"updated_at"`

	// Computed fields
	JobsCount         int `db:"-" json:"jobs_count"`
	ActiveJobsCount   int `db:"-" json:"active_jobs_count"`
	TotalApplications int `db:"-" json:"total_applications"`
}

// CompanyAdminResponse represents the company response for admin API
type CompanyAdminResponse struct {
	ID                 uint64 `json:"id"`
	HashID             string `json:"hash_id"`
	Email              string `json:"email"`
	FullName           string `json:"full_name"`
	Phone              string `json:"phone,omitempty"`
	CompanyName        string `json:"company_name"`
	CompanyDescription string `json:"company_description,omitempty"`
	CompanyWebsite     string `json:"company_website,omitempty"`
	CompanyLogoURL     string `json:"company_logo_url,omitempty"`
	CompanyStatus      string `json:"company_status"`
	IsActive           bool   `json:"is_active"`
	IsVerified         bool   `json:"is_verified"`
	EmailVerifiedAt    string `json:"email_verified_at,omitempty"`
	CreatedAt          string `json:"created_at"`
	UpdatedAt          string `json:"updated_at"`
	JobsCount          int    `json:"jobs_count"`
	ActiveJobsCount    int    `json:"active_jobs_count"`
	TotalApplications  int    `json:"total_applications"`
}

func (c *CompanyAdmin) ToResponse() *CompanyAdminResponse {
	resp := &CompanyAdminResponse{
		ID:                c.ID,
		HashID:            hashid.Encode(c.ID),
		Email:             c.Email,
		FullName:          c.FullName,
		IsActive:          c.IsActive,
		IsVerified:        c.IsVerified,
		CreatedAt:         c.CreatedAt.Format(time.RFC3339),
		UpdatedAt:         c.UpdatedAt.Format(time.RFC3339),
		JobsCount:         c.JobsCount,
		ActiveJobsCount:   c.ActiveJobsCount,
		TotalApplications: c.TotalApplications,
	}

	if c.Phone.Valid {
		resp.Phone = c.Phone.String
	}
	if c.CompanyName.Valid {
		resp.CompanyName = c.CompanyName.String
	}
	if c.CompanyDescription.Valid {
		resp.CompanyDescription = c.CompanyDescription.String
	}
	if c.CompanyWebsite.Valid {
		resp.CompanyWebsite = c.CompanyWebsite.String
	}
	if c.CompanyLogoURL.Valid {
		resp.CompanyLogoURL = c.CompanyLogoURL.String
	}
	if c.CompanyStatus.Valid {
		resp.CompanyStatus = c.CompanyStatus.String
	}
	if c.EmailVerifiedAt.Valid {
		resp.EmailVerifiedAt = c.EmailVerifiedAt.Time.Format(time.RFC3339)
	}

	return resp
}

// CompanyDetailResponse represents the detailed company information for admin
type CompanyDetailResponse struct {
	// Basic Info
	ID       uint64 `json:"id"`
	HashID   string `json:"hash_id"`
	Email    string `json:"email"`
	FullName string `json:"full_name"`
	Phone    string `json:"phone,omitempty"`

	// Company Info
	CompanyName        string `json:"company_name"`
	CompanyDescription string `json:"company_description,omitempty"`
	CompanyWebsite     string `json:"company_website,omitempty"`
	CompanyLogoURL     string `json:"company_logo_url,omitempty"`
	CompanyIndustry    string `json:"company_industry,omitempty"`
	CompanySize        string `json:"company_size,omitempty"`
	CompanyLocation    string `json:"company_location,omitempty"`
	CompanyAddress     string `json:"company_address,omitempty"`
	CompanyCity        string `json:"company_city,omitempty"`
	CompanyProvince    string `json:"company_province,omitempty"`
	PostalCode         string `json:"postal_code,omitempty"`
	EstablishedYear    int    `json:"established_year,omitempty"`
	EmployeeCount      int    `json:"employee_count,omitempty"`

	// Verification Status
	CompanyStatus       string `json:"company_status"`
	IsActive            bool   `json:"is_active"`
	IsVerified          bool   `json:"is_verified"`
	EmailVerifiedAt     string `json:"email_verified_at,omitempty"`
	DocumentsVerifiedAt string `json:"documents_verified_at,omitempty"`
	VerificationNotes   string `json:"verification_notes,omitempty"`

	// Legal Documents
	LegalDocuments struct {
		KtpFounderURL    string `json:"ktp_founder_url,omitempty"`
		AktaPendirianURL string `json:"akta_pendirian_url,omitempty"`
		NpwpURL          string `json:"npwp_url,omitempty"`
		NibURL           string `json:"nib_url,omitempty"`
	} `json:"legal_documents"`

	// Job & Application Stats
	JobsCount         int `json:"jobs_count"`
	ActiveJobsCount   int `json:"active_jobs_count"`
	TotalApplications int `json:"total_applications"`

	// Quota Info
	QuotaInfo struct {
		FreeQuotaUsed  int `json:"free_quota_used"`
		FreeQuotaTotal int `json:"free_quota_total"` // Usually 5
		PaidQuota      int `json:"paid_quota"`
		TotalQuota     int `json:"total_quota"` // free + paid

		// Job posting details
		FreeJobsActive  int `json:"free_jobs_active"`  // Active jobs using free quota
		PaidJobsActive  int `json:"paid_jobs_active"`  // Active jobs using paid quota
		TotalJobsActive int `json:"total_jobs_active"` // Total active jobs
		DraftJobsCount  int `json:"draft_jobs_count"`  // Draft jobs
	} `json:"quota_info"`

	// Timestamps
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// ============================================
// JOB ENTITIES (ADMIN VIEW)
// ============================================

// Job statuses for admin oversight
const (
	JobStatusDraft    = "draft"
	JobStatusPending  = "pending"
	JobStatusActive   = "active"
	JobStatusPaused   = "paused"
	JobStatusClosed   = "closed"
	JobStatusFilled   = "filled"
	JobStatusRejected = "rejected"
	JobStatusFlagged  = "flagged"
)

// JobAdmin represents job data for admin view
type JobAdmin struct {
	ID                uint64         `db:"id" json:"id"`
	CompanyID         uint64         `db:"company_id" json:"company_id"`
	CompanyName       sql.NullString `db:"company_name" json:"company_name"`
	Title             string         `db:"title" json:"title"`
	Slug              string         `db:"slug" json:"slug"`
	Description       string         `db:"description" json:"description"`
	Requirements      sql.NullString `db:"requirements" json:"requirements,omitempty"`
	City              string         `db:"city" json:"city"`
	Province          string         `db:"province" json:"province"`
	IsRemote          bool           `db:"is_remote" json:"is_remote"`
	JobType           string         `db:"job_type" json:"job_type"`
	ExperienceLevel   string         `db:"experience_level" json:"experience_level"`
	SalaryMin         sql.NullInt64  `db:"salary_min" json:"salary_min,omitempty"`
	SalaryMax         sql.NullInt64  `db:"salary_max" json:"salary_max,omitempty"`
	Status            string         `db:"status" json:"status"`
	AdminStatus       sql.NullString `db:"admin_status" json:"admin_status,omitempty"`
	AdminNote         sql.NullString `db:"admin_note" json:"admin_note,omitempty"`
	FlagReason        sql.NullString `db:"flag_reason" json:"flag_reason,omitempty"`
	ViewsCount        uint64         `db:"views_count" json:"views_count"`
	ApplicationsCount uint64         `db:"applications_count" json:"applications_count"`
	PublishedAt       sql.NullTime   `db:"published_at" json:"published_at,omitempty"`
	CreatedAt         time.Time      `db:"created_at" json:"created_at"`
	UpdatedAt         time.Time      `db:"updated_at" json:"updated_at"`
}

// JobAdminResponse represents the job response for admin API
type JobAdminResponse struct {
	ID                uint64 `json:"id"`
	CompanyID         uint64 `json:"company_id"`
	CompanyName       string `json:"company_name"`
	Title             string `json:"title"`
	Slug              string `json:"slug"`
	Description       string `json:"description"`
	Requirements      string `json:"requirements,omitempty"`
	City              string `json:"city"`
	Province          string `json:"province"`
	IsRemote          bool   `json:"is_remote"`
	JobType           string `json:"job_type"`
	ExperienceLevel   string `json:"experience_level"`
	SalaryMin         *int64 `json:"salary_min,omitempty"`
	SalaryMax         *int64 `json:"salary_max,omitempty"`
	Status            string `json:"status"`
	AdminStatus       string `json:"admin_status,omitempty"`
	AdminNote         string `json:"admin_note,omitempty"`
	FlagReason        string `json:"flag_reason,omitempty"`
	ViewsCount        uint64 `json:"views_count"`
	ApplicationsCount uint64 `json:"applications_count"`
	PublishedAt       string `json:"published_at,omitempty"`
	CreatedAt         string `json:"created_at"`
	UpdatedAt         string `json:"updated_at"`
}

func (j *JobAdmin) ToResponse() *JobAdminResponse {
	resp := &JobAdminResponse{
		ID:                j.ID,
		CompanyID:         j.CompanyID,
		Title:             j.Title,
		Slug:              j.Slug,
		Description:       j.Description,
		City:              j.City,
		Province:          j.Province,
		IsRemote:          j.IsRemote,
		JobType:           j.JobType,
		ExperienceLevel:   j.ExperienceLevel,
		Status:            j.Status,
		ViewsCount:        j.ViewsCount,
		ApplicationsCount: j.ApplicationsCount,
		CreatedAt:         j.CreatedAt.Format(time.RFC3339),
		UpdatedAt:         j.UpdatedAt.Format(time.RFC3339),
	}

	if j.CompanyName.Valid {
		resp.CompanyName = j.CompanyName.String
	}
	if j.Requirements.Valid {
		resp.Requirements = j.Requirements.String
	}
	if j.SalaryMin.Valid {
		resp.SalaryMin = &j.SalaryMin.Int64
	}
	if j.SalaryMax.Valid {
		resp.SalaryMax = &j.SalaryMax.Int64
	}
	if j.AdminStatus.Valid {
		resp.AdminStatus = j.AdminStatus.String
	}
	if j.AdminNote.Valid {
		resp.AdminNote = j.AdminNote.String
	}
	if j.FlagReason.Valid {
		resp.FlagReason = j.FlagReason.String
	}
	if j.PublishedAt.Valid {
		resp.PublishedAt = j.PublishedAt.Time.Format(time.RFC3339)
	}

	return resp
}

// ============================================
// PAYMENT ENTITIES (ADMIN VIEW)
// ============================================

// Payment statuses
const (
	PaymentStatusPending   = "pending"
	PaymentStatusConfirmed = "confirmed"
	PaymentStatusRejected  = "rejected"
)

// PaymentAdmin represents payment data for admin view
type PaymentAdmin struct {
	ID            uint64         `db:"id" json:"id"`
	CompanyID     uint64         `db:"company_id" json:"company_id"`
	CompanyName   sql.NullString `db:"company_name" json:"company_name"`
	JobID         sql.NullInt64  `db:"job_id" json:"job_id,omitempty"`
	JobTitle      sql.NullString `db:"job_title" json:"job_title,omitempty"`
	Amount        int64          `db:"amount" json:"amount"`
	ProofImageURL sql.NullString `db:"proof_image_url" json:"proof_image_url,omitempty"`
	Status        string         `db:"status" json:"status"`
	Note          sql.NullString `db:"note" json:"note,omitempty"`
	ConfirmedByID sql.NullInt64  `db:"confirmed_by_id" json:"confirmed_by_id,omitempty"`
	SubmittedAt   time.Time      `db:"submitted_at" json:"submitted_at"`
	ConfirmedAt   sql.NullTime   `db:"confirmed_at" json:"confirmed_at,omitempty"`
	CreatedAt     time.Time      `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time      `db:"updated_at" json:"updated_at"`
}

// PaymentAdminResponse represents the payment response for admin API
type PaymentAdminResponse struct {
	ID            uint64  `json:"id"`
	CompanyID     uint64  `json:"company_id"`
	CompanyName   string  `json:"company_name"`
	JobID         *uint64 `json:"job_id,omitempty"`
	JobTitle      string  `json:"job_title,omitempty"`
	Amount        int64   `json:"amount"`
	ProofImageURL string  `json:"proof_image_url,omitempty"`
	Status        string  `json:"status"`
	StatusLabel   string  `json:"status_label"`
	Note          string  `json:"note,omitempty"`
	ConfirmedByID *uint64 `json:"confirmed_by_id,omitempty"`
	SubmittedAt   string  `json:"submitted_at"`
	ConfirmedAt   string  `json:"confirmed_at,omitempty"`
}

func (p *PaymentAdmin) ToResponse() *PaymentAdminResponse {
	resp := &PaymentAdminResponse{
		ID:          p.ID,
		CompanyID:   p.CompanyID,
		Amount:      p.Amount,
		Status:      p.Status,
		StatusLabel: getPaymentStatusLabel(p.Status),
		SubmittedAt: p.SubmittedAt.Format(time.RFC3339),
	}

	if p.CompanyName.Valid {
		resp.CompanyName = p.CompanyName.String
	}
	if p.JobID.Valid {
		jobID := uint64(p.JobID.Int64)
		resp.JobID = &jobID
	}
	if p.JobTitle.Valid {
		resp.JobTitle = p.JobTitle.String
	}
	if p.ProofImageURL.Valid {
		resp.ProofImageURL = p.ProofImageURL.String
	}
	if p.Note.Valid {
		resp.Note = p.Note.String
	}
	if p.ConfirmedByID.Valid {
		confirmedBy := uint64(p.ConfirmedByID.Int64)
		resp.ConfirmedByID = &confirmedBy
	}
	if p.ConfirmedAt.Valid {
		resp.ConfirmedAt = p.ConfirmedAt.Time.Format(time.RFC3339)
	}

	return resp
}

func getPaymentStatusLabel(status string) string {
	switch status {
	case PaymentStatusPending:
		return "Menunggu Konfirmasi"
	case PaymentStatusConfirmed:
		return "Dikonfirmasi"
	case PaymentStatusRejected:
		return "Ditolak"
	default:
		return status
	}
}

// ============================================
// JOB SEEKER ENTITIES (ADMIN VIEW)
// ============================================

// JobSeekerAdmin represents job seeker data for admin view
type JobSeekerAdmin struct {
	ID              uint64         `db:"id" json:"id"`
	Email           string         `db:"email" json:"email"`
	FullName        string         `db:"full_name" json:"full_name"`
	Phone           sql.NullString `db:"phone" json:"phone,omitempty"`
	AvatarURL       sql.NullString `db:"avatar_url" json:"avatar_url,omitempty"`
	IsActive        bool           `db:"is_active" json:"is_active"`
	IsVerified      bool           `db:"is_verified" json:"is_verified"`
	EmailVerifiedAt sql.NullTime   `db:"email_verified_at" json:"email_verified_at,omitempty"`
	CreatedAt       time.Time      `db:"created_at" json:"created_at"`
	UpdatedAt       time.Time      `db:"updated_at" json:"updated_at"`

	// Computed fields
	ApplicationsCount int  `db:"-" json:"applications_count"`
	HasCV             bool `db:"-" json:"has_cv"`
}

// JobSeekerAdminResponse represents the job seeker response for admin API
type JobSeekerAdminResponse struct {
	ID                uint64 `json:"id"`
	Email             string `json:"email"`
	FullName          string `json:"full_name"`
	Phone             string `json:"phone,omitempty"`
	AvatarURL         string `json:"avatar_url,omitempty"`
	IsActive          bool   `json:"is_active"`
	IsVerified        bool   `json:"is_verified"`
	EmailVerifiedAt   string `json:"email_verified_at,omitempty"`
	CreatedAt         string `json:"created_at"`
	UpdatedAt         string `json:"updated_at"`
	ApplicationsCount int    `json:"applications_count"`
	HasCV             bool   `json:"has_cv"`
}

func (js *JobSeekerAdmin) ToResponse() *JobSeekerAdminResponse {
	resp := &JobSeekerAdminResponse{
		ID:                js.ID,
		Email:             js.Email,
		FullName:          js.FullName,
		IsActive:          js.IsActive,
		IsVerified:        js.IsVerified,
		CreatedAt:         js.CreatedAt.Format(time.RFC3339),
		UpdatedAt:         js.UpdatedAt.Format(time.RFC3339),
		ApplicationsCount: js.ApplicationsCount,
		HasCV:             js.HasCV,
	}

	if js.Phone.Valid {
		resp.Phone = js.Phone.String
	}
	if js.AvatarURL.Valid {
		resp.AvatarURL = js.AvatarURL.String
	}
	if js.EmailVerifiedAt.Valid {
		resp.EmailVerifiedAt = js.EmailVerifiedAt.Time.Format(time.RFC3339)
	}

	return resp
}

// ============================================
// DASHBOARD STATS
// ============================================

// DashboardStats represents admin dashboard statistics
type DashboardStats struct {
	TotalCompanies       int   `json:"total_companies"`
	PendingVerifications int   `json:"pending_verifications"`
	VerifiedCompanies    int   `json:"verified_companies"`
	SuspendedCompanies   int   `json:"suspended_companies"`
	TotalJobs            int   `json:"total_jobs"`
	ActiveJobs           int   `json:"active_jobs"`
	PendingJobs          int   `json:"pending_jobs"`
	FlaggedJobs          int   `json:"flagged_jobs"`
	TotalJobSeekers      int   `json:"total_job_seekers"`
	ActiveJobSeekers     int   `json:"active_job_seekers"`
	TotalPayments        int   `json:"total_payments"`
	PendingPayments      int   `json:"pending_payments"`
	TotalRevenue         int64 `json:"total_revenue"`
}

// ============================================
// ADMIN ACTION LOG
// ============================================

// AdminActionLog represents an admin action log entry
type AdminActionLog struct {
	ID         uint64         `db:"id" json:"id"`
	AdminID    uint64         `db:"admin_id" json:"admin_id"`
	AdminName  string         `db:"admin_name" json:"admin_name"`
	Action     string         `db:"action" json:"action"`
	EntityType string         `db:"entity_type" json:"entity_type"`
	EntityID   uint64         `db:"entity_id" json:"entity_id"`
	Details    sql.NullString `db:"details" json:"details,omitempty"`
	IPAddress  sql.NullString `db:"ip_address" json:"ip_address,omitempty"`
	CreatedAt  time.Time      `db:"created_at" json:"created_at"`
}

// ============================================
// REQUEST DTOs
// ============================================

// AdminLoginRequest represents admin login request
type AdminLoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// AdminAuthResponse represents admin authentication response
type AdminAuthResponse struct {
	Admin       *AdminUserResponse `json:"admin"`
	AccessToken string             `json:"access_token"`
	ExpiresIn   int64              `json:"expires_in"`
}

// CompanyVerificationRequest represents company verification action
type CompanyVerificationRequest struct {
	Action string `json:"action" validate:"required,oneof=approve reject"`
	Reason string `json:"reason,omitempty" validate:"required_if=Action reject"`
}

// CompanyStatusRequest represents company status change action
type CompanyStatusRequest struct {
	Action string `json:"action" validate:"required,oneof=suspend reactivate"`
	Reason string `json:"reason,omitempty" validate:"required_if=Action suspend"`
}

// JobActionRequest represents job moderation action
type JobActionRequest struct {
	Action string `json:"action" validate:"required,oneof=approve reject close flag unflag"`
	Reason string `json:"reason,omitempty"`
}

// PaymentActionRequest represents payment approval/rejection action
type PaymentActionRequest struct {
	Action string `json:"action" validate:"required,oneof=approve reject"`
	Note   string `json:"note,omitempty" validate:"required_if=Action reject"`
}

// JobSeekerActionRequest represents job seeker moderation action
type JobSeekerActionRequest struct {
	Action string `json:"action" validate:"required,oneof=suspend reactivate deactivate"`
	Reason string `json:"reason,omitempty" validate:"required_if=Action suspend"`
}

// ============================================
// FILTER PARAMETERS
// ============================================

// CompanyFilter represents filter parameters for companies
type CompanyFilter struct {
	Status   string `json:"status"`
	Search   string `json:"search"`
	Page     int    `json:"page"`
	PageSize int    `json:"page_size"`
}

// JobFilter represents filter parameters for jobs
type JobFilter struct {
	CompanyID uint64 `json:"company_id"`
	Status    string `json:"status"`
	DateFrom  string `json:"date_from"`
	DateTo    string `json:"date_to"`
	Search    string `json:"search"`
	Page      int    `json:"page"`
	PageSize  int    `json:"page_size"`
}

// PaymentFilter represents filter parameters for payments
type PaymentFilter struct {
	CompanyID uint64 `json:"company_id"`
	Status    string `json:"status"`
	DateFrom  string `json:"date_from"`
	DateTo    string `json:"date_to"`
	Page      int    `json:"page"`
	PageSize  int    `json:"page_size"`
}

// JobSeekerFilter represents filter parameters for job seekers
type JobSeekerFilter struct {
	Status   string `json:"status"`
	Search   string `json:"search"`
	Page     int    `json:"page"`
	PageSize int    `json:"page_size"`
}

// ============================================
// PAGINATED RESPONSE
// ============================================

// PaginatedResponse represents a paginated response
type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Total      int         `json:"total"`
	Page       int         `json:"page"`
	PageSize   int         `json:"page_size"`
	TotalPages int         `json:"total_pages"`
}

func NewPaginatedResponse(data interface{}, total, page, pageSize int) *PaginatedResponse {
	totalPages := total / pageSize
	if total%pageSize > 0 {
		totalPages++
	}
	return &PaginatedResponse{
		Data:       data,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}
}
