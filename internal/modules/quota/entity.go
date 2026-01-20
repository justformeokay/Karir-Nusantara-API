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
	FreeQuotaLimit = 10         // Free job postings per company
	PricePerJob    = 10000      // IDR 10,000 per additional job
)

// TopUpPackage represents a top-up package option
type TopUpPackage struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Quota       int    `json:"quota"`       // Number of job posts purchased
	BonusQuota  int    `json:"bonus_quota"` // Bonus job posts
	TotalQuota  int    `json:"total_quota"` // Total quota received
	Price       int64  `json:"price"`       // Price in IDR
	PricePerJob int64  `json:"price_per_job"` // Effective price per job
	IsBestValue bool   `json:"is_best_value"` // Highlight as best value
	Description string `json:"description"`
}

// GetTopUpPackages returns all available top-up packages
func GetTopUpPackages() []TopUpPackage {
	return []TopUpPackage{
		{
			ID:          "single",
			Name:        "1 Posting",
			Quota:       1,
			BonusQuota:  0,
			TotalQuota:  1,
			Price:       10000,
			PricePerJob: 10000,
			IsBestValue: false,
			Description: "Bayar untuk 1 lowongan",
		},
		{
			ID:          "pack5",
			Name:        "5 Posting",
			Quota:       5,
			BonusQuota:  0,
			TotalQuota:  5,
			Price:       50000,
			PricePerJob: 10000,
			IsBestValue: false,
			Description: "Hemat waktu, beli 5 sekaligus",
		},
		{
			ID:          "pack10",
			Name:        "10 Posting + 2 GRATIS",
			Quota:       10,
			BonusQuota:  2,
			TotalQuota:  12,
			Price:       100000,
			PricePerJob: 8333, // 100000 / 12
			IsBestValue: true,
			Description: "Beli 10 dapat 12! Hemat Rp 20.000",
		},
		{
			ID:          "pack20",
			Name:        "20 Posting + 5 GRATIS",
			Quota:       20,
			BonusQuota:  5,
			TotalQuota:  25,
			Price:       200000,
			PricePerJob: 8000, // 200000 / 25
			IsBestValue: false,
			Description: "Beli 20 dapat 25! Hemat Rp 50.000",
		},
	}
}

// GetPackageByID returns a package by its ID
func GetPackageByID(id string) *TopUpPackage {
	for _, pkg := range GetTopUpPackages() {
		if pkg.ID == id {
			return &pkg
		}
	}
	return nil
}

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
	PackageID     sql.NullString `db:"package_id" json:"package_id,omitempty"`
	QuotaAmount   int            `db:"quota_amount" json:"quota_amount"`
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
	PackageID     string  `json:"package_id,omitempty"`
	PackageName   string  `json:"package_name,omitempty"`
	QuotaAmount   int     `json:"quota_amount"`
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
		QuotaAmount: p.QuotaAmount,
		Amount:      p.Amount,
		Status:      p.Status,
		StatusLabel: GetStatusLabel(p.Status),
		SubmittedAt: p.SubmittedAt.Format(time.RFC3339),
	}

	if p.JobID.Valid {
		jobID := uint64(p.JobID.Int64)
		resp.JobID = &jobID
	}
	if p.PackageID.Valid {
		resp.PackageID = p.PackageID.String
		// Get package name
		if pkg := GetPackageByID(p.PackageID.String); pkg != nil {
			resp.PackageName = pkg.Name
		}
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
	JobID     uint64 `json:"job_id" validate:"omitempty"`
	PackageID string `json:"package_id" validate:"omitempty"` // Package ID for top-up (single, pack5, pack10, pack20)
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
