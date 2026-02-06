package admin

import (
	"context"
	"errors"
	"fmt"
	"time"
)

// PartnerService handles partner-related business logic for admin
type PartnerService interface {
	// Partner management
	GetPartners(ctx context.Context, status string, search string, page, limit int) (*AdminPartnerListResponse, error)
	GetPartnerByID(ctx context.Context, id uint64) (*AdminPartnerDetailResponse, error)
	UpdatePartnerStatus(ctx context.Context, id uint64, req UpdatePartnerStatusRequest) error
	ApprovePartner(ctx context.Context, id uint64, approvedBy uint64, req ApprovePartnerRequest) error

	// Referred companies
	GetReferredCompanies(ctx context.Context, search string, page, limit int) (*AdminReferredCompanyListResponse, error)

	// Payouts
	GetPayouts(ctx context.Context, status string, search string, page, limit int) (*AdminPayoutListResponse, error)
	CreatePayout(ctx context.Context, req CreatePayoutRequest) (uint64, error)
	ProcessPayout(ctx context.Context, id uint64, req ProcessPayoutRequest) error

	// Partner balances
	GetPartnerBalances(ctx context.Context, page, limit int) (*AdminPartnerBalanceListResponse, error)

	// Stats
	GetReferralStats(ctx context.Context) (*AdminReferralStatsResponse, error)
	GetPayoutStats(ctx context.Context) (*AdminPayoutStatsResponse, error)
}

type partnerService struct {
	repo PartnerRepository
}

// NewPartnerService creates a new partner service for admin
func NewPartnerService(repo PartnerRepository) PartnerService {
	return &partnerService{repo: repo}
}

// GetPartners returns paginated list of partners
func (s *partnerService) GetPartners(ctx context.Context, status string, search string, page, limit int) (*AdminPartnerListResponse, error) {
	partners, total, err := s.repo.GetPartners(ctx, status, search, page, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get partners: %w", err)
	}

	items := make([]AdminPartnerItem, 0, len(partners))
	for _, p := range partners {
		items = append(items, AdminPartnerItem{
			ID:               p.ID,
			FullName:         p.FullName,
			Email:            p.Email,
			Phone:            stringPtrValue(p.Phone),
			ReferralCode:     p.ReferralCode,
			CommissionRate:   p.CommissionRate,
			Status:           p.Status,
			TotalReferrals:   p.TotalReferrals,
			TotalCommission:  p.TotalCommission,
			AvailableBalance: p.AvailableBalance,
			PaidAmount:       p.PaidAmount,
			BankInfo: &BankInfo{
				BankName:          stringPtrValue(p.BankName),
				BankAccountNumber: stringPtrValue(p.BankAccountNumber),
				BankAccountHolder: stringPtrValue(p.BankAccountHolder),
				IsVerified:        p.IsBankVerified,
			},
			CreatedAt: p.CreatedAt,
		})
	}

	totalPages := (total + limit - 1) / limit

	return &AdminPartnerListResponse{
		Partners: items,
		Pagination: Pagination{
			Page:       page,
			Limit:      limit,
			Total:      total,
			TotalPages: totalPages,
		},
	}, nil
}

// GetPartnerByID returns a partner by ID
func (s *partnerService) GetPartnerByID(ctx context.Context, id uint64) (*AdminPartnerDetailResponse, error) {
	p, err := s.repo.GetPartnerByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get partner: %w", err)
	}
	if p == nil {
		return nil, errors.New("partner not found")
	}

	// Get last payout date
	lastPayoutDate, _ := s.repo.GetLastPayoutDate(ctx, id)

	return &AdminPartnerDetailResponse{
		ID:             p.ID,
		UserID:         p.UserID,
		FullName:       p.FullName,
		Email:          p.Email,
		Phone:          stringPtrValue(p.Phone),
		ReferralCode:   p.ReferralCode,
		CommissionRate: p.CommissionRate,
		Status:         p.Status,
		BankInfo: &BankInfo{
			BankName:          stringPtrValue(p.BankName),
			BankAccountNumber: stringPtrValue(p.BankAccountNumber),
			BankAccountHolder: stringPtrValue(p.BankAccountHolder),
			IsVerified:        p.IsBankVerified,
		},
		TotalReferrals:   p.TotalReferrals,
		TotalCommission:  p.TotalCommission,
		AvailableBalance: p.AvailableBalance,
		PendingBalance:   p.PendingBalance,
		PaidAmount:       p.PaidAmount,
		ApprovedAt:       p.ApprovedAt,
		LastPayoutDate:   lastPayoutDate,
		Notes:            stringPtrValue(p.Notes),
		CreatedAt:        p.CreatedAt,
		UpdatedAt:        p.UpdatedAt,
	}, nil
}

// UpdatePartnerStatus updates partner status (suspend/activate)
func (s *partnerService) UpdatePartnerStatus(ctx context.Context, id uint64, req UpdatePartnerStatusRequest) error {
	// Validate status
	validStatuses := map[string]bool{
		"active":    true,
		"suspended": true,
		"pending":   true,
	}
	if !validStatuses[req.Status] {
		return errors.New("invalid status, must be one of: active, suspended, pending")
	}

	// Check partner exists
	p, err := s.repo.GetPartnerByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get partner: %w", err)
	}
	if p == nil {
		return errors.New("partner not found")
	}

	return s.repo.UpdatePartnerStatus(ctx, id, req.Status, req.Notes)
}

// ApprovePartner approves a pending partner
func (s *partnerService) ApprovePartner(ctx context.Context, id uint64, approvedBy uint64, req ApprovePartnerRequest) error {
	// Check partner exists
	p, err := s.repo.GetPartnerByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get partner: %w", err)
	}
	if p == nil {
		return errors.New("partner not found")
	}

	// Check partner is pending
	if p.Status != "pending" {
		return errors.New("partner is not in pending status")
	}

	// Default commission rate if not provided
	commissionRate := 10.0 // 10%
	if req.CommissionRate > 0 {
		commissionRate = req.CommissionRate
	}

	return s.repo.ApprovePartner(ctx, id, approvedBy, commissionRate, req.Notes)
}

// GetReferredCompanies returns paginated list of referred companies
func (s *partnerService) GetReferredCompanies(ctx context.Context, search string, page, limit int) (*AdminReferredCompanyListResponse, error) {
	companies, total, err := s.repo.GetReferredCompanies(ctx, search, page, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get referred companies: %w", err)
	}

	items := make([]AdminReferredCompanyItem, 0, len(companies))
	for _, c := range companies {
		// Get transaction stats for this company
		transactions, revenue, commission, _ := s.repo.GetReferredCompanyStats(ctx, c.CompanyID)

		items = append(items, AdminReferredCompanyItem{
			ID:          c.ID,
			CompanyID:   c.CompanyID,
			CompanyName: c.CompanyName,
			PartnerInfo: PartnerReferrerInfo{
				ID:           c.PartnerID,
				Name:         c.PartnerName,
				ReferralCode: c.ReferralCode,
			},
			TotalTransactions:     transactions,
			TotalRevenueGenerated: revenue,
			TotalCommission:       commission,
			RegistrationDate:      c.RegistrationDate,
			Status:                c.CompanyStatus,
		})
	}

	totalPages := (total + limit - 1) / limit

	return &AdminReferredCompanyListResponse{
		Companies: items,
		Pagination: Pagination{
			Page:       page,
			Limit:      limit,
			Total:      total,
			TotalPages: totalPages,
		},
	}, nil
}

// GetPayouts returns paginated list of payouts
func (s *partnerService) GetPayouts(ctx context.Context, status string, search string, page, limit int) (*AdminPayoutListResponse, error) {
	payouts, total, err := s.repo.GetPayouts(ctx, status, search, page, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get payouts: %w", err)
	}

	items := make([]AdminPayoutItem, 0, len(payouts))
	for _, p := range payouts {
		items = append(items, AdminPayoutItem{
			ID:          p.ID,
			PartnerID:   p.PartnerID,
			PartnerName: p.PartnerName,
			Amount:      p.Amount,
			Status:      p.Status,
			BankInfo: &BankInfo{
				BankName:          stringPtrValue(p.BankName),
				BankAccountNumber: stringPtrValue(p.BankAccountNumber),
				BankAccountHolder: stringPtrValue(p.BankAccountHolder),
			},
			PayoutProofURL: stringPtrValue(p.PayoutProofURL),
			RequestedAt:    p.RequestedAt,
			PaidAt:         p.PaidAt,
			Notes:          stringPtrValue(p.Notes),
		})
	}

	totalPages := (total + limit - 1) / limit

	return &AdminPayoutListResponse{
		Payouts: items,
		Pagination: Pagination{
			Page:       page,
			Limit:      limit,
			Total:      total,
			TotalPages: totalPages,
		},
	}, nil
}

// CreatePayout creates a new payout request (admin initiated)
func (s *partnerService) CreatePayout(ctx context.Context, req CreatePayoutRequest) (uint64, error) {
	// Check partner exists
	p, err := s.repo.GetPartnerByID(ctx, req.PartnerID)
	if err != nil {
		return 0, fmt.Errorf("failed to get partner: %w", err)
	}
	if p == nil {
		return 0, errors.New("partner not found")
	}

	// Check partner has enough balance
	if p.AvailableBalance < req.Amount {
		return 0, errors.New("insufficient balance")
	}

	return s.repo.CreatePayout(ctx, req.PartnerID, req.Amount, req.Notes)
}

// ProcessPayout marks a payout as paid
func (s *partnerService) ProcessPayout(ctx context.Context, id uint64, req ProcessPayoutRequest) error {
	// Check payout exists
	p, err := s.repo.GetPayoutByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get payout: %w", err)
	}
	if p == nil {
		return errors.New("payout not found")
	}

	// Check payout is pending or processing
	if p.Status != "pending" && p.Status != "processing" {
		return errors.New("payout is not in pending or processing status")
	}

	return s.repo.ProcessPayout(ctx, id, req.PayoutProofURL, req.Notes)
}

// GetPartnerBalances returns partners with available balance
func (s *partnerService) GetPartnerBalances(ctx context.Context, page, limit int) (*AdminPartnerBalanceListResponse, error) {
	partners, total, err := s.repo.GetPartnersWithBalance(ctx, page, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get partner balances: %w", err)
	}

	items := make([]AdminPartnerBalanceItem, 0, len(partners))
	for _, p := range partners {
		// Get last payout date
		lastPayoutDate, _ := s.repo.GetLastPayoutDate(ctx, p.ID)

		items = append(items, AdminPartnerBalanceItem{
			ID:               p.ID,
			PartnerName:      p.FullName,
			Email:            p.Email,
			AvailableBalance: p.AvailableBalance,
			PendingBalance:   p.PendingBalance,
			TotalPaid:        p.PaidAmount,
			BankInfo: &BankInfo{
				BankName:          stringPtrValue(p.BankName),
				BankAccountNumber: stringPtrValue(p.BankAccountNumber),
				BankAccountHolder: stringPtrValue(p.BankAccountHolder),
				IsVerified:        p.IsBankVerified,
			},
			LastPayoutDate: lastPayoutDate,
		})
	}

	totalPages := (total + limit - 1) / limit

	return &AdminPartnerBalanceListResponse{
		Partners: items,
		Pagination: Pagination{
			Page:       page,
			Limit:      limit,
			Total:      total,
			TotalPages: totalPages,
		},
	}, nil
}

// GetReferralStats returns overall referral program statistics
func (s *partnerService) GetReferralStats(ctx context.Context) (*AdminReferralStatsResponse, error) {
	return s.repo.GetReferralStats(ctx)
}

// GetPayoutStats returns payout statistics
func (s *partnerService) GetPayoutStats(ctx context.Context) (*AdminPayoutStatsResponse, error) {
	return s.repo.GetPayoutStats(ctx)
}

// Helper function to convert *string to string
func stringPtrValue(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// Helper type for pagination
type Pagination struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

// Response types with pagination
type AdminPartnerListResponse struct {
	Partners   []AdminPartnerItem `json:"partners"`
	Pagination Pagination         `json:"pagination"`
}

type AdminReferredCompanyListResponse struct {
	Companies  []AdminReferredCompanyItem `json:"companies"`
	Pagination Pagination                 `json:"pagination"`
}

type AdminPayoutListResponse struct {
	Payouts    []AdminPayoutItem `json:"payouts"`
	Pagination Pagination        `json:"pagination"`
}

type AdminPartnerBalanceListResponse struct {
	Partners   []AdminPartnerBalanceItem `json:"partners"`
	Pagination Pagination                `json:"pagination"`
}

// Partner list item
type AdminPartnerItem struct {
	ID               uint64    `json:"id"`
	FullName         string    `json:"full_name"`
	Email            string    `json:"email"`
	Phone            string    `json:"phone,omitempty"`
	ReferralCode     string    `json:"referral_code"`
	CommissionRate   float64   `json:"commission_rate"`
	Status           string    `json:"status"`
	TotalReferrals   int       `json:"total_referrals"`
	TotalCommission  int64     `json:"total_commission"`
	AvailableBalance int64     `json:"available_balance"`
	PaidAmount       int64     `json:"paid_amount"`
	BankInfo         *BankInfo `json:"bank_info,omitempty"`
	CreatedAt        time.Time `json:"created_at"`
}

// Partner detail response
type AdminPartnerDetailResponse struct {
	ID               uint64     `json:"id"`
	UserID           uint64     `json:"user_id"`
	FullName         string     `json:"full_name"`
	Email            string     `json:"email"`
	Phone            string     `json:"phone,omitempty"`
	ReferralCode     string     `json:"referral_code"`
	CommissionRate   float64    `json:"commission_rate"`
	Status           string     `json:"status"`
	BankInfo         *BankInfo  `json:"bank_info,omitempty"`
	TotalReferrals   int        `json:"total_referrals"`
	TotalCommission  int64      `json:"total_commission"`
	AvailableBalance int64      `json:"available_balance"`
	PendingBalance   int64      `json:"pending_balance"`
	PaidAmount       int64      `json:"paid_amount"`
	ApprovedAt       *time.Time `json:"approved_at,omitempty"`
	LastPayoutDate   *time.Time `json:"last_payout_date,omitempty"`
	Notes            string     `json:"notes,omitempty"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
}

// Bank info
type BankInfo struct {
	BankName          string `json:"bank_name,omitempty"`
	BankAccountNumber string `json:"bank_account_number,omitempty"`
	BankAccountHolder string `json:"bank_account_holder,omitempty"`
	IsVerified        bool   `json:"is_verified"`
}

// Referred company item
type AdminReferredCompanyItem struct {
	ID                    uint64              `json:"id"`
	CompanyID             uint64              `json:"company_id"`
	CompanyName           string              `json:"company_name"`
	PartnerInfo           PartnerReferrerInfo `json:"partner_info"`
	TotalTransactions     int                 `json:"total_transactions"`
	TotalRevenueGenerated int64               `json:"total_revenue_generated"`
	TotalCommission       int64               `json:"total_commission"`
	RegistrationDate      time.Time           `json:"registration_date"`
	Status                string              `json:"status"`
}

// Partner referrer info
type PartnerReferrerInfo struct {
	ID           uint64 `json:"id"`
	Name         string `json:"name"`
	ReferralCode string `json:"referral_code"`
}

// Payout item
type AdminPayoutItem struct {
	ID             uint64     `json:"id"`
	PartnerID      uint64     `json:"partner_id"`
	PartnerName    string     `json:"partner_name"`
	Amount         int64      `json:"amount"`
	Status         string     `json:"status"`
	BankInfo       *BankInfo  `json:"bank_info,omitempty"`
	PayoutProofURL string     `json:"payout_proof_url,omitempty"`
	RequestedAt    time.Time  `json:"requested_at"`
	PaidAt         *time.Time `json:"paid_at,omitempty"`
	Notes          string     `json:"notes,omitempty"`
}

// Partner balance item
type AdminPartnerBalanceItem struct {
	ID               uint64     `json:"id"`
	PartnerName      string     `json:"partner_name"`
	Email            string     `json:"email"`
	AvailableBalance int64      `json:"available_balance"`
	PendingBalance   int64      `json:"pending_balance"`
	TotalPaid        int64      `json:"total_paid"`
	BankInfo         *BankInfo  `json:"bank_info,omitempty"`
	LastPayoutDate   *time.Time `json:"last_payout_date,omitempty"`
}
