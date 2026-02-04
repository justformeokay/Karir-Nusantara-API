package partner

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/karirnusantara/api/internal/config"
	apperrors "github.com/karirnusantara/api/internal/shared/errors"
	"golang.org/x/crypto/bcrypt"
)

// Service defines the partner service interface
type Service interface {
	// Authentication
	Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error)
	ValidateAccessToken(tokenString string) (*TokenClaims, error)

	// Profile
	GetProfile(ctx context.Context, userID uint64) (*PartnerUserResponse, error)
	UpdateProfile(ctx context.Context, userID uint64, req *UpdateProfileRequest) error
	ChangePassword(ctx context.Context, userID uint64, req *ChangePasswordRequest) error

	// Dashboard
	GetDashboardStats(ctx context.Context, partnerID uint64) (*DashboardStatsResponse, error)
	GetMonthlyData(ctx context.Context, partnerID uint64) ([]MonthlyDataResponse, error)

	// Referral
	GetReferralInfo(ctx context.Context, partnerID uint64) (*ReferralInfoResponse, error)

	// Companies
	GetCompanies(ctx context.Context, partnerID uint64, page, limit int, search string) ([]CompanyResponse, *PaginationResponse, error)
	GetCompaniesSummary(ctx context.Context, partnerID uint64) (*CompanySummaryResponse, error)

	// Transactions
	GetTransactions(ctx context.Context, partnerID uint64, page, limit int, search, status string) ([]TransactionResponse, *PaginationResponse, error)
	GetTransactionsSummary(ctx context.Context, partnerID uint64) (*TransactionSummaryResponse, error)

	// Payouts
	GetPayoutInfo(ctx context.Context, partnerID uint64) (*PayoutInfoResponse, error)
	GetPayoutHistory(ctx context.Context, partnerID uint64, page, limit int) ([]PayoutHistoryResponse, *PaginationResponse, error)
}

// TokenClaims represents JWT claims for partner tokens
type TokenClaims struct {
	UserID       uint64 `json:"user_id"`
	PartnerID    uint64 `json:"partner_id"`
	Email        string `json:"email"`
	Role         string `json:"role"`
	ReferralCode string `json:"referral_code"`
	jwt.RegisteredClaims
}

// service implements Service
type service struct {
	repo    Repository
	config  *config.JWTConfig
	baseURL string
}

// NewService creates a new partner service
func NewService(repo Repository, cfg *config.JWTConfig, baseURL string) Service {
	return &service{
		repo:    repo,
		config:  cfg,
		baseURL: baseURL,
	}
}

// Login authenticates a partner
func (s *service) Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error) {
	// Get partner user by email
	partnerUser, err := s.repo.GetPartnerUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, apperrors.NewInternalError("Failed to get user", err)
	}
	if partnerUser == nil {
		return nil, apperrors.NewInvalidCredentialsError()
	}

	// Check if user is active
	if !partnerUser.IsActive {
		return nil, apperrors.NewForbiddenError("Account is deactivated")
	}

	// Check partner status
	if partnerUser.PartnerStatus == StatusSuspended {
		return nil, apperrors.NewForbiddenError("Partner account is suspended")
	}

	if partnerUser.PartnerStatus == StatusPending {
		return nil, apperrors.NewForbiddenError("Partner account is pending approval")
	}

	// Get password hash and verify
	passwordHash, err := s.repo.GetUserPasswordHash(ctx, partnerUser.ID)
	if err != nil {
		return nil, apperrors.NewInternalError("Failed to verify credentials", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(req.Password)); err != nil {
		return nil, apperrors.NewInvalidCredentialsError()
	}

	// Generate tokens
	accessToken, err := s.generateAccessToken(partnerUser)
	if err != nil {
		return nil, apperrors.NewInternalError("Failed to generate access token", err)
	}

	refreshToken := s.generateRefreshToken()

	return &LoginResponse{
		User:         partnerUser.ToPartnerUserResponse(),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(s.config.AccessExpiry.Seconds()),
	}, nil
}

// ValidateAccessToken validates JWT token and returns claims
func (s *service) ValidateAccessToken(tokenString string) (*TokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.config.Secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*TokenClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

// GetProfile retrieves partner profile
func (s *service) GetProfile(ctx context.Context, userID uint64) (*PartnerUserResponse, error) {
	partnerUser, err := s.repo.GetPartnerUserByEmail(ctx, "")
	if err != nil {
		return nil, apperrors.NewInternalError("Failed to get profile", err)
	}

	// Need to get by user ID instead
	partner, err := s.repo.GetPartnerByUserID(ctx, userID)
	if err != nil {
		return nil, apperrors.NewInternalError("Failed to get partner", err)
	}
	if partner == nil {
		return nil, apperrors.NewNotFoundError("Partner")
	}

	// Construct partner user response
	resp := &PartnerUserResponse{
		ID:            userID,
		ReferralCode:  partner.ReferralCode,
		AccountStatus: partner.Status,
		CreatedAt:     partner.CreatedAt.Format(time.RFC3339),
	}

	if partner.BankName.Valid && partner.BankName.String != "" {
		resp.BankAccount = &BankAccount{
			BankName:      partner.BankName.String,
			AccountNumber: MaskAccountNumber(partner.BankAccountNumber.String),
			AccountHolder: partner.BankAccountHolder.String,
			IsVerified:    partner.IsBankVerified,
		}
	}

	// Get user info separately
	partnerUser, err = s.repo.GetPartnerUserByEmail(ctx, "")
	if partnerUser != nil {
		resp.Name = partnerUser.FullName
		resp.Email = partnerUser.Email
	}

	return resp, nil
}

// UpdateProfile updates partner display name
func (s *service) UpdateProfile(ctx context.Context, userID uint64, req *UpdateProfileRequest) error {
	if err := s.repo.UpdateUserName(ctx, userID, req.Name); err != nil {
		return apperrors.NewInternalError("Failed to update profile", err)
	}
	return nil
}

// ChangePassword changes partner password
func (s *service) ChangePassword(ctx context.Context, userID uint64, req *ChangePasswordRequest) error {
	// Verify current password
	currentHash, err := s.repo.GetUserPasswordHash(ctx, userID)
	if err != nil {
		return apperrors.NewInternalError("Failed to verify password", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(currentHash), []byte(req.CurrentPassword)); err != nil {
		return apperrors.NewValidationError("Current password is incorrect", nil)
	}

	// Hash new password
	newHash, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return apperrors.NewInternalError("Failed to hash password", err)
	}

	// Update password
	if err := s.repo.UpdateUserPassword(ctx, userID, string(newHash)); err != nil {
		return apperrors.NewInternalError("Failed to update password", err)
	}

	return nil
}

// GetDashboardStats retrieves dashboard statistics
func (s *service) GetDashboardStats(ctx context.Context, partnerID uint64) (*DashboardStatsResponse, error) {
	stats, err := s.repo.GetDashboardStats(ctx, partnerID)
	if err != nil {
		return nil, apperrors.NewInternalError("Failed to get dashboard stats", err)
	}
	return stats, nil
}

// GetMonthlyData retrieves monthly chart data
func (s *service) GetMonthlyData(ctx context.Context, partnerID uint64) ([]MonthlyDataResponse, error) {
	data, err := s.repo.GetMonthlyData(ctx, partnerID, 7)
	if err != nil {
		return nil, apperrors.NewInternalError("Failed to get monthly data", err)
	}
	return data, nil
}

// GetReferralInfo retrieves referral information
func (s *service) GetReferralInfo(ctx context.Context, partnerID uint64) (*ReferralInfoResponse, error) {
	partner, err := s.repo.GetPartnerByID(ctx, partnerID)
	if err != nil {
		return nil, apperrors.NewInternalError("Failed to get referral info", err)
	}
	if partner == nil {
		return nil, apperrors.NewNotFoundError("Partner")
	}

	return &ReferralInfoResponse{
		ReferralCode:   partner.ReferralCode,
		ReferralLink:   fmt.Sprintf("%s/register?ref=%s", s.baseURL, partner.ReferralCode),
		CommissionRate: partner.CommissionRate,
	}, nil
}

// GetCompanies retrieves referred companies
func (s *service) GetCompanies(ctx context.Context, partnerID uint64, page, limit int, search string) ([]CompanyResponse, *PaginationResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 15
	}

	companies, total, err := s.repo.GetReferredCompanies(ctx, partnerID, page, limit, search)
	if err != nil {
		return nil, nil, apperrors.NewInternalError("Failed to get companies", err)
	}

	totalPages := (total + limit - 1) / limit

	pagination := &PaginationResponse{
		CurrentPage:  page,
		TotalPages:   totalPages,
		TotalItems:   total,
		ItemsPerPage: limit,
	}

	return companies, pagination, nil
}

// GetCompaniesSummary retrieves companies summary
func (s *service) GetCompaniesSummary(ctx context.Context, partnerID uint64) (*CompanySummaryResponse, error) {
	summary, err := s.repo.GetCompaniesSummary(ctx, partnerID)
	if err != nil {
		return nil, apperrors.NewInternalError("Failed to get companies summary", err)
	}
	return summary, nil
}

// GetTransactions retrieves transactions
func (s *service) GetTransactions(ctx context.Context, partnerID uint64, page, limit int, search, status string) ([]TransactionResponse, *PaginationResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 15
	}

	transactions, total, err := s.repo.GetTransactions(ctx, partnerID, page, limit, search, status)
	if err != nil {
		return nil, nil, apperrors.NewInternalError("Failed to get transactions", err)
	}

	totalPages := (total + limit - 1) / limit

	pagination := &PaginationResponse{
		CurrentPage:  page,
		TotalPages:   totalPages,
		TotalItems:   total,
		ItemsPerPage: limit,
	}

	return transactions, pagination, nil
}

// GetTransactionsSummary retrieves transactions summary
func (s *service) GetTransactionsSummary(ctx context.Context, partnerID uint64) (*TransactionSummaryResponse, error) {
	summary, err := s.repo.GetTransactionsSummary(ctx, partnerID)
	if err != nil {
		return nil, apperrors.NewInternalError("Failed to get transactions summary", err)
	}
	return summary, nil
}

// GetPayoutInfo retrieves payout information
func (s *service) GetPayoutInfo(ctx context.Context, partnerID uint64) (*PayoutInfoResponse, error) {
	info, err := s.repo.GetPayoutInfo(ctx, partnerID)
	if err != nil {
		return nil, apperrors.NewInternalError("Failed to get payout info", err)
	}
	return info, nil
}

// GetPayoutHistory retrieves payout history
func (s *service) GetPayoutHistory(ctx context.Context, partnerID uint64, page, limit int) ([]PayoutHistoryResponse, *PaginationResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 15
	}

	payouts, total, err := s.repo.GetPayoutHistory(ctx, partnerID, page, limit)
	if err != nil {
		return nil, nil, apperrors.NewInternalError("Failed to get payout history", err)
	}

	totalPages := (total + limit - 1) / limit

	pagination := &PaginationResponse{
		CurrentPage:  page,
		TotalPages:   totalPages,
		TotalItems:   total,
		ItemsPerPage: limit,
	}

	return payouts, pagination, nil
}

// Private helper methods

func (s *service) generateAccessToken(user *PartnerUser) (string, error) {
	claims := TokenClaims{
		UserID:       user.ID,
		PartnerID:    user.PartnerID,
		Email:        user.Email,
		Role:         "partner",
		ReferralCode: user.ReferralCode,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.config.AccessExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   fmt.Sprintf("partner_%d", user.ID),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.config.Secret))
}

func (s *service) generateRefreshToken() string {
	// Generate a random token
	data := fmt.Sprintf("%d-%d", time.Now().UnixNano(), time.Now().Unix())
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}
