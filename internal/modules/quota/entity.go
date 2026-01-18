package quota

import (
	"database/sql"
	"time"
)

// Payment statuses
const (
	PaymentStatusPending   = "pending"
	PaymentStatusConfirmed = "confirmed"
	PaymentStatusRejected  = "rejected"
)

// Constants
const (
	FreeQuotaLimit = 5          // Free job postings per company
	PricePerJob    = 30000      // IDR 30,000 per additional job
)

// CompanyQuota represents a company's job posting quota
type CompanyQuota struct {
	ID              uint64    `db:"id" json:"id"`
	CompanyID       uint64    `db:"company_id" json:"company_id"`
	FreeQuotaUsed   int       `db:"free_quota_used" json:"free_quota_used"`
	PaidQuota       int       `db:"paid_quota" json:"paid_quota"`
	CreatedAt       time.Time `db:"created_at" json:"created_at"`
	UpdatedAt       time.Time `db:"updated_at" json:"updated_at"`
}

// Payment represents a payment for job posting quota
type Payment struct {
	ID            uint64         `db:"id" json:"id"`
	CompanyID     uint64         `db:"company_id" json:"company_id"`
	JobID         sql.NullInt64  `db:"job_id" json:"job_id,omitempty"`
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

// QuotaResponse represents the quota response for API
type QuotaResponse struct {
	FreeQuota          int   `json:"free_quota"`
	UsedFreeQuota      int   `json:"used_free_quota"`
	RemainingFreeQuota int   `json:"remaining_free_quota"`
	PaidQuota          int   `json:"paid_quota"`
	PricePerJob        int64 `json:"price_per_job"`
}

// PaymentResponse represents a payment response
type PaymentResponse struct {
	ID            uint64  `json:"id"`
	JobID         *uint64 `json:"job_id,omitempty"`
	JobTitle      string  `json:"job_title,omitempty"`
	Amount        int64   `json:"amount"`
	ProofImageURL string  `json:"proof_image_url,omitempty"`
	Status        string  `json:"status"`
	StatusLabel   string  `json:"status_label"`
	Note          string  `json:"note,omitempty"`
	SubmittedAt   string  `json:"submitted_at"`
	ConfirmedAt   string  `json:"confirmed_at,omitempty"`
}

// GetStatusLabel returns the Indonesian label for payment status
func GetStatusLabel(status string) string {
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

// ToResponse converts Payment to PaymentResponse
func (p *Payment) ToResponse() *PaymentResponse {
	resp := &PaymentResponse{
		ID:          p.ID,
		Amount:      p.Amount,
		Status:      p.Status,
		StatusLabel: GetStatusLabel(p.Status),
		SubmittedAt: p.SubmittedAt.Format(time.RFC3339),
	}

	if p.JobID.Valid {
		jobID := uint64(p.JobID.Int64)
		resp.JobID = &jobID
	}
	if p.ProofImageURL.Valid {
		resp.ProofImageURL = p.ProofImageURL.String
	}
	if p.Note.Valid {
		resp.Note = p.Note.String
	}
	if p.ConfirmedAt.Valid {
		resp.ConfirmedAt = p.ConfirmedAt.Time.Format(time.RFC3339)
	}

	return resp
}

// Request DTOs

// SubmitPaymentProofRequest represents a payment proof submission request
type SubmitPaymentProofRequest struct {
	JobID uint64 `json:"job_id" validate:"omitempty"`
}

// PaymentListParams represents payment list parameters
type PaymentListParams struct {
	Page      int
	PerPage   int
	CompanyID uint64
	Status    string
}

// DefaultPaymentListParams returns default list parameters
func DefaultPaymentListParams() PaymentListParams {
	return PaymentListParams{
		Page:    1,
		PerPage: 20,
	}
}
