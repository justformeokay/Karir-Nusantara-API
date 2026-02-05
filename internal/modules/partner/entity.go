package partner

import (
	"database/sql"
	"time"
)

// Partner status constants
const (
	StatusActive    = "active"
	StatusInactive  = "inactive"
	StatusPending   = "pending"
	StatusSuspended = "suspended"
)

// Commission status constants
const (
	CommissionPending   = "pending"
	CommissionApproved  = "approved"
	CommissionPaid      = "paid"
	CommissionCancelled = "cancelled"
)

// Payout status constants
const (
	PayoutPending    = "pending"
	PayoutProcessing = "processing"
	PayoutCompleted  = "completed"
	PayoutFailed     = "failed"
	PayoutCancelled  = "cancelled"
)

// ReferralPartner represents the referral_partners table
type ReferralPartner struct {
	ID                uint64         `db:"id" json:"id"`
	UserID            uint64         `db:"user_id" json:"user_id"`
	ReferralCode      string         `db:"referral_code" json:"referral_code"`
	CommissionRate    float64        `db:"commission_rate" json:"commission_rate"`
	Status            string         `db:"status" json:"status"`
	BankName          sql.NullString `db:"bank_name" json:"bank_name,omitempty"`
	BankAccountNumber sql.NullString `db:"bank_account_number" json:"bank_account_number,omitempty"`
	BankAccountHolder sql.NullString `db:"bank_account_holder" json:"bank_account_holder,omitempty"`
	IsBankVerified    bool           `db:"is_bank_verified" json:"is_bank_verified"`
	TotalReferrals    int            `db:"total_referrals" json:"total_referrals"`
	TotalCommission   int64          `db:"total_commission" json:"total_commission"`
	AvailableBalance  int64          `db:"available_balance" json:"available_balance"`
	PendingBalance    int64          `db:"pending_balance" json:"pending_balance"`
	PaidAmount        int64          `db:"paid_amount" json:"paid_amount"`
	ApprovedBy        sql.NullInt64  `db:"approved_by" json:"approved_by,omitempty"`
	ApprovedAt        sql.NullTime   `db:"approved_at" json:"approved_at,omitempty"`
	Notes             sql.NullString `db:"notes" json:"notes,omitempty"`
	CreatedAt         time.Time      `db:"created_at" json:"created_at"`
	UpdatedAt         time.Time      `db:"updated_at" json:"updated_at"`
}

// PartnerReferral represents the partner_referrals table
type PartnerReferral struct {
	ID               uint64         `db:"id" json:"id"`
	PartnerID        uint64         `db:"partner_id" json:"partner_id"`
	CompanyID        uint64         `db:"company_id" json:"company_id"`
	ReferralCodeUsed string         `db:"referral_code_used" json:"referral_code_used"`
	RegisteredAt     time.Time      `db:"registered_at" json:"registered_at"`
	IsVerified       bool           `db:"is_verified" json:"is_verified"`
	FirstPaymentAt   sql.NullTime   `db:"first_payment_at" json:"first_payment_at,omitempty"`
	Notes            sql.NullString `db:"notes" json:"notes,omitempty"`
}

// PartnerCommission represents the partner_commissions table
type PartnerCommission struct {
	ID                uint64        `db:"id" json:"id"`
	PartnerID         uint64        `db:"partner_id" json:"partner_id"`
	ReferralID        uint64        `db:"referral_id" json:"referral_id"`
	PaymentID         uint64        `db:"payment_id" json:"payment_id"`
	CompanyID         uint64        `db:"company_id" json:"company_id"`
	TransactionAmount int64         `db:"transaction_amount" json:"transaction_amount"`
	CommissionRate    float64       `db:"commission_rate" json:"commission_rate"`
	CommissionAmount  int64         `db:"commission_amount" json:"commission_amount"`
	JobQuota          int           `db:"job_quota" json:"job_quota"`
	Status            string        `db:"status" json:"status"`
	ApprovedBy        sql.NullInt64 `db:"approved_by" json:"approved_by,omitempty"`
	ApprovedAt        sql.NullTime  `db:"approved_at" json:"approved_at,omitempty"`
	PaidAt            sql.NullTime  `db:"paid_at" json:"paid_at,omitempty"`
	PayoutID          sql.NullInt64 `db:"payout_id" json:"payout_id,omitempty"`
	CreatedAt         time.Time     `db:"created_at" json:"created_at"`
	UpdatedAt         time.Time     `db:"updated_at" json:"updated_at"`
}

// PartnerPayout represents the partner_payouts table
type PartnerPayout struct {
	ID                uint64         `db:"id" json:"id"`
	PartnerID         uint64         `db:"partner_id" json:"partner_id"`
	Amount            int64          `db:"amount" json:"amount"`
	BankName          string         `db:"bank_name" json:"bank_name"`
	BankAccountNumber string         `db:"bank_account_number" json:"bank_account_number"`
	BankAccountHolder string         `db:"bank_account_holder" json:"bank_account_holder"`
	Status            string         `db:"status" json:"status"`
	TransferRef       sql.NullString `db:"transfer_ref" json:"transfer_ref,omitempty"`
	FailureReason     sql.NullString `db:"failure_reason" json:"failure_reason,omitempty"`
	RequestedAt       time.Time      `db:"requested_at" json:"requested_at"`
	ProcessedBy       sql.NullInt64  `db:"processed_by" json:"processed_by,omitempty"`
	ProcessedAt       sql.NullTime   `db:"processed_at" json:"processed_at,omitempty"`
	CompletedAt       sql.NullTime   `db:"completed_at" json:"completed_at,omitempty"`
	Notes             sql.NullString `db:"notes" json:"notes,omitempty"`
	CreatedAt         time.Time      `db:"created_at" json:"created_at"`
	UpdatedAt         time.Time      `db:"updated_at" json:"updated_at"`
}

// PartnerUser represents user with partner data joined
type PartnerUser struct {
	// User fields
	ID         uint64    `db:"id" json:"id"`
	Email      string    `db:"email" json:"email"`
	FullName   string    `db:"full_name" json:"full_name"`
	Role       string    `db:"role" json:"role"`
	IsActive   bool      `db:"is_active" json:"is_active"`
	IsVerified bool      `db:"is_verified" json:"is_verified"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`

	// Partner fields
	PartnerID         uint64         `db:"partner_id" json:"partner_id"`
	ReferralCode      string         `db:"referral_code" json:"referral_code"`
	CommissionRate    float64        `db:"commission_rate" json:"commission_rate"`
	PartnerStatus     string         `db:"partner_status" json:"partner_status"`
	BankName          sql.NullString `db:"bank_name" json:"bank_name,omitempty"`
	BankAccountNumber sql.NullString `db:"bank_account_number" json:"bank_account_number,omitempty"`
	BankAccountHolder sql.NullString `db:"bank_account_holder" json:"bank_account_holder,omitempty"`
	IsBankVerified    bool           `db:"is_bank_verified" json:"is_bank_verified"`
	TotalReferrals    int            `db:"total_referrals" json:"total_referrals"`
	TotalCommission   int64          `db:"total_commission" json:"total_commission"`
	AvailableBalance  int64          `db:"available_balance" json:"available_balance"`
	PendingBalance    int64          `db:"pending_balance" json:"pending_balance"`
	PaidAmount        int64          `db:"paid_amount" json:"paid_amount"`
}

// ========================
// DTOs (Request/Response)
// ========================

// RegisterRequest for partner registration
type RegisterRequest struct {
	FullName string `json:"full_name" validate:"required,min=2,max=100"`
	Email    string `json:"email" validate:"required,email"`
	Phone    string `json:"phone" validate:"required,min=10,max=15"`
	Password string `json:"password" validate:"required,min=8"`
}

// RegisterResponse for partner registration
type RegisterResponse struct {
	User    *PartnerUserResponse `json:"user"`
	Message string               `json:"message"`
}

// ForgotPasswordRequest for forgot password
type ForgotPasswordRequest struct {
	Email string `json:"email" validate:"required,email"`
}

// ResetPasswordRequest for reset password
type ResetPasswordRequest struct {
	Token       string `json:"token" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=8"`
}

// LoginRequest for partner login
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

// LoginResponse for partner login
type LoginResponse struct {
	User         *PartnerUserResponse `json:"user"`
	AccessToken  string               `json:"access_token"`
	RefreshToken string               `json:"refresh_token"`
	ExpiresIn    int64                `json:"expires_in"`
}

// PartnerUserResponse is the safe user response for partners
type PartnerUserResponse struct {
	ID            uint64       `json:"id"`
	Name          string       `json:"name"`
	Email         string       `json:"email"`
	ReferralCode  string       `json:"referral_code"`
	AccountStatus string       `json:"account_status"`
	BankAccount   *BankAccount `json:"bank_account,omitempty"`
	CreatedAt     string       `json:"created_at"`
}

// BankAccount represents bank account info
type BankAccount struct {
	BankName      string `json:"bank_name"`
	AccountNumber string `json:"account_number"` // Masked
	AccountHolder string `json:"account_holder"`
	IsVerified    bool   `json:"is_verified"`
}

// DashboardStatsResponse for dashboard statistics
type DashboardStatsResponse struct {
	TotalCompanies    int   `json:"total_companies"`
	TotalTransactions int   `json:"total_transactions"`
	TotalCommission   int64 `json:"total_commission"`
	AvailableBalance  int64 `json:"available_balance"`
	PaidCommission    int64 `json:"paid_commission"`
}

// MonthlyDataResponse for chart data
type MonthlyDataResponse struct {
	Month      string `json:"month"`
	Commission int64  `json:"commission"`
	Companies  int    `json:"companies"`
}

// ReferralInfoResponse for referral info
type ReferralInfoResponse struct {
	ReferralCode   string  `json:"referral_code"`
	ReferralLink   string  `json:"referral_link"`
	CommissionRate float64 `json:"commission_rate"`
}

// CompanyResponse for company list
type CompanyResponse struct {
	ID               string `json:"id"`
	Name             string `json:"name"`
	RegistrationDate string `json:"registration_date"`
	Status           string `json:"status"`
	TotalJobPosts    int    `json:"total_job_posts"`
	TotalRevenue     int64  `json:"total_revenue"`
	CommissionEarned int64  `json:"commission_earned"`
}

// CompanySummaryResponse for companies summary
type CompanySummaryResponse struct {
	TotalCompanies  int   `json:"total_companies"`
	ActiveCompanies int   `json:"active_companies"`
	TotalRevenue    int64 `json:"total_revenue"`
	TotalCommission int64 `json:"total_commission"`
}

// TransactionResponse for transaction list
type TransactionResponse struct {
	ID          string `json:"id"`
	Date        string `json:"date"`
	CompanyName string `json:"company_name"`
	CompanyID   string `json:"company_id"`
	JobQuota    int    `json:"job_quota"`
	Amount      int64  `json:"amount"`
	Commission  int64  `json:"commission"`
	Status      string `json:"status"`
}

// TransactionSummaryResponse for transaction summary
type TransactionSummaryResponse struct {
	TotalCommission   int64 `json:"total_commission"`
	PendingCommission int64 `json:"pending_commission"`
	PaidCommission    int64 `json:"paid_commission"`
	TotalTransactions int   `json:"total_transactions"`
}

// PayoutInfoResponse for payout information
type PayoutInfoResponse struct {
	Balance     *BalanceInfo    `json:"balance"`
	BankAccount *BankAccount    `json:"bank_account"`
	Schedule    *PayoutSchedule `json:"schedule"`
}

// BalanceInfo for balance details
type BalanceInfo struct {
	AvailableBalance int64  `json:"available_balance"`
	PendingPayout    int64  `json:"pending_payout"`
	PaidAmount       int64  `json:"paid_amount"`
	LastPayoutDate   string `json:"last_payout_date,omitempty"`
}

// PayoutSchedule for payout schedule info
type PayoutSchedule struct {
	ProcessingDay int    `json:"processing_day"`
	MinimumPayout int64  `json:"minimum_payout"`
	TransferTime  string `json:"transfer_time"`
}

// PayoutHistoryResponse for payout history
type PayoutHistoryResponse struct {
	ID          string `json:"id"`
	Amount      int64  `json:"amount"`
	Status      string `json:"status"`
	TransferRef string `json:"transfer_ref,omitempty"`
	ProcessedAt string `json:"processed_at,omitempty"`
	CreatedAt   string `json:"created_at"`
}

// UpdateProfileRequest for profile update
type UpdateProfileRequest struct {
	Name string `json:"name" validate:"required,min=2,max=100"`
}

// ChangePasswordRequest for password change
type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" validate:"required,min=6"`
	NewPassword     string `json:"new_password" validate:"required,min=6"`
}

// PaginationResponse for paginated responses
type PaginationResponse struct {
	CurrentPage  int `json:"current_page"`
	TotalPages   int `json:"total_pages"`
	TotalItems   int `json:"total_items"`
	ItemsPerPage int `json:"items_per_page"`
}

// Helper methods

// MaskAccountNumber masks bank account number
func MaskAccountNumber(accountNumber string) string {
	if len(accountNumber) <= 4 {
		return accountNumber
	}
	return "****" + accountNumber[len(accountNumber)-4:]
}

// ToPartnerUserResponse converts PartnerUser to PartnerUserResponse
func (p *PartnerUser) ToPartnerUserResponse() *PartnerUserResponse {
	resp := &PartnerUserResponse{
		ID:            p.ID,
		Name:          p.FullName,
		Email:         p.Email,
		ReferralCode:  p.ReferralCode,
		AccountStatus: p.PartnerStatus,
		CreatedAt:     p.CreatedAt.Format(time.RFC3339),
	}

	if p.BankName.Valid && p.BankName.String != "" {
		resp.BankAccount = &BankAccount{
			BankName:      p.BankName.String,
			AccountNumber: MaskAccountNumber(p.BankAccountNumber.String),
			AccountHolder: p.BankAccountHolder.String,
			IsVerified:    p.IsBankVerified,
		}
	}

	return resp
}
