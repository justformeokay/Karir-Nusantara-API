package jobreports

import (
	"database/sql"
	"time"
)

// JobReport represents a job listing report from job seekers
type JobReport struct {
	ID          uint64        `json:"id" db:"id"`
	JobID       uint64        `json:"job_id" db:"job_id"`
	UserID      uint64        `json:"user_id" db:"user_id"`
	Reason      string        `json:"reason" db:"reason"`
	Description string        `json:"description" db:"description"`
	Status      string        `json:"status" db:"status"`
	AdminNote   string        `json:"admin_note" db:"admin_note"`
	ReviewedBy  sql.NullInt64 `json:"reviewed_by" db:"reviewed_by"`
	ReviewedAt  sql.NullTime  `json:"reviewed_at" db:"reviewed_at"`
	CreatedAt   time.Time     `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at" db:"updated_at"`
}

// JobReportWithDetails includes job and reporter info
type JobReportWithDetails struct {
	JobReport
	JobTitle      string `json:"job_title" db:"job_title"`
	JobStatus     string `json:"job_status" db:"job_status"`
	CompanyID     uint64 `json:"company_id" db:"company_id"`
	CompanyHashID string `json:"company_hash_id" db:"company_hash_id"`
	CompanyName   string `json:"company_name" db:"company_name"`
	CompanyStatus string `json:"company_status" db:"company_status"`
	ReporterName  string `json:"reporter_name" db:"reporter_name"`
	ReporterEmail string `json:"reporter_email" db:"reporter_email"`
	ReviewerName  string `json:"reviewer_name" db:"reviewer_name"`
	TotalReports  int    `json:"total_reports" db:"total_reports"`
}

// CreateReportRequest represents request to create a new job report
type CreateReportRequest struct {
	Reason      string `json:"reason" validate:"required,oneof=scam misleading inappropriate discriminatory other"`
	Description string `json:"description" validate:"required,min=10,max=1000"`
}

// UpdateReportStatusRequest represents request to update report status
type UpdateReportStatusRequest struct {
	Status    string `json:"status" validate:"required,oneof=pending reviewed dismissed action_taken"`
	AdminNote string `json:"admin_notes" validate:"omitempty,max=500"`
}

// BanCompanyRequest represents request to ban a company
type BanCompanyRequest struct {
	CompanyID uint64 `json:"company_id" validate:"required"`
	Reason    string `json:"reason" validate:"required,min=10,max=500"`
}

// Report status constants
const (
	ReportStatusPending     = "pending"
	ReportStatusReviewed    = "reviewed"
	ReportStatusDismissed   = "dismissed"
	ReportStatusActionTaken = "action_taken"
)

// Report reason constants
const (
	ReportReasonScam           = "scam"
	ReportReasonMisleading     = "misleading"
	ReportReasonInappropriate  = "inappropriate"
	ReportReasonDiscriminatory = "discriminatory"
	ReportReasonOther          = "other"
)

// CooldownDuration is the minimum time between report creations per job (24 hours)
const CooldownDuration = 24 * time.Hour

// GetReasonLabel returns human-readable reason label
func GetReasonLabel(reason string) string {
	labels := map[string]string{
		"scam":           "Penipuan / Scam",
		"misleading":     "Informasi Menyesatkan",
		"inappropriate":  "Konten Tidak Pantas",
		"discriminatory": "Diskriminatif",
		"other":          "Lainnya",
	}
	if label, ok := labels[reason]; ok {
		return label
	}
	return reason
}

// GetStatusLabel returns human-readable status label
func GetStatusLabel(status string) string {
	labels := map[string]string{
		"pending":      "Menunggu Review",
		"reviewed":     "Sudah Direview",
		"dismissed":    "Diabaikan",
		"action_taken": "Tindakan Diambil",
	}
	if label, ok := labels[status]; ok {
		return label
	}
	return status
}
