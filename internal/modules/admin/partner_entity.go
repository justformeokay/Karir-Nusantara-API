package admin

import "time"

// Request types for admin partner management

// UpdatePartnerStatusRequest for updating partner status
type UpdatePartnerStatusRequest struct {
	Status string  `json:"status" validate:"required,oneof=active suspended pending"`
	Notes  *string `json:"notes,omitempty"`
}

// ApprovePartnerRequest for approving a pending partner
type ApprovePartnerRequest struct {
	CommissionRate float64 `json:"commission_rate" validate:"min=0,max=100"`
	Notes          *string `json:"notes,omitempty"`
}

// ProcessPayoutRequest for processing a payout
type ProcessPayoutRequest struct {
	PayoutProofURL string  `json:"payout_proof_url" validate:"url"`
	Notes          *string `json:"notes,omitempty"`
}

// CreatePayoutRequest for creating a manual payout
type CreatePayoutRequest struct {
	PartnerID uint64  `json:"partner_id" validate:"required"`
	Amount    int64   `json:"amount" validate:"required,gt=0"`
	Notes     *string `json:"notes,omitempty"`
}

// Response types for admin APIs

// AdminReferralStatsResponse represents referral program statistics
type AdminReferralStatsResponse struct {
	TotalPartners            int   `json:"total_partners"`
	ActivePartners           int   `json:"active_partners"`
	PendingPartners          int   `json:"pending_partners"`
	TotalReferredCompanies   int   `json:"total_referred_companies"`
	TotalCommissionGenerated int64 `json:"total_commission_generated"`
	PendingPayouts           int64 `json:"pending_payouts"`
	TotalPaidOut             int64 `json:"total_paid_out"`
	PartnersWithBalance      int   `json:"partners_with_balance"`
}

// AdminPayoutStatsResponse represents payout statistics
type AdminPayoutStatsResponse struct {
	TotalCommissionGenerated int64 `json:"total_commission_generated"`
	PendingPayouts           int64 `json:"pending_payouts"`
	TotalPaidOut             int64 `json:"total_paid_out"`
	PartnersWithBalance      int   `json:"partners_with_balance"`
}

// Database row types

// PartnerDBRow represents a row from database join query
type PartnerDBRow struct {
	ID                uint64     `db:"id"`
	UserID            uint64     `db:"user_id"`
	FullName          string     `db:"full_name"`
	Email             string     `db:"email"`
	Phone             *string    `db:"phone"`
	ReferralCode      string     `db:"referral_code"`
	CommissionRate    float64    `db:"commission_rate"`
	Status            string     `db:"status"`
	BankName          *string    `db:"bank_name"`
	BankAccountNumber *string    `db:"bank_account_number"`
	BankAccountHolder *string    `db:"bank_account_holder"`
	IsBankVerified    bool       `db:"is_bank_verified"`
	TotalReferrals    int        `db:"total_referrals"`
	TotalCommission   int64      `db:"total_commission"`
	AvailableBalance  int64      `db:"available_balance"`
	PendingBalance    int64      `db:"pending_balance"`
	PaidAmount        int64      `db:"paid_amount"`
	ApprovedBy        *uint64    `db:"approved_by"`
	ApprovedAt        *time.Time `db:"approved_at"`
	Notes             *string    `db:"notes"`
	CreatedAt         time.Time  `db:"created_at"`
	UpdatedAt         time.Time  `db:"updated_at"`
}

// ReferredCompanyDBRow represents a row from referred companies query
type ReferredCompanyDBRow struct {
	ID               uint64    `db:"id"`
	CompanyID        uint64    `db:"company_id"`
	CompanyName      string    `db:"company_name"`
	PartnerID        uint64    `db:"partner_id"`
	PartnerName      string    `db:"partner_name"`
	ReferralCode     string    `db:"referral_code"`
	RegistrationDate time.Time `db:"registration_date"`
	CompanyStatus    string    `db:"company_status"`
}

// PayoutDBRow represents a row from payouts query
type PayoutDBRow struct {
	ID                uint64     `db:"id"`
	PartnerID         uint64     `db:"partner_id"`
	PartnerName       string     `db:"partner_name"`
	Amount            int64      `db:"amount"`
	Status            string     `db:"status"`
	PayoutProofURL    *string    `db:"payout_proof_url"`
	RequestedAt       time.Time  `db:"requested_at"`
	PaidAt            *time.Time `db:"paid_at"`
	Notes             *string    `db:"notes"`
	BankName          *string    `db:"bank_name"`
	BankAccountNumber *string    `db:"bank_account_number"`
	BankAccountHolder *string    `db:"bank_account_holder"`
}
