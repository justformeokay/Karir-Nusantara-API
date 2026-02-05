package partner

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/karirnusantara/api/internal/middleware"
	"github.com/karirnusantara/api/internal/shared/response"
	"github.com/karirnusantara/api/internal/shared/validator"
)

// Handler handles partner HTTP requests
type Handler struct {
	service   Service
	validator *validator.Validator
}

// NewHandler creates a new partner handler
func NewHandler(service Service, validator *validator.Validator) *Handler {
	return &Handler{
		service:   service,
		validator: validator,
	}
}

// Register handles partner registration
// POST /api/v1/partner/auth/register
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid request body")
		return
	}

	// Validate request
	if errors := h.validator.Validate(&req); errors != nil {
		response.UnprocessableEntity(w, "Validation failed", errors)
		return
	}

	// Register
	authResp, err := h.service.Register(r.Context(), &req)
	if err != nil {
		handleError(w, err)
		return
	}

	response.Created(w, "Registration successful", authResp)
}

// Login handles partner login
// POST /api/v1/partner/auth/login
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid request body")
		return
	}

	// Validate request
	if errors := h.validator.Validate(&req); errors != nil {
		response.UnprocessableEntity(w, "Validation failed", errors)
		return
	}

	// Login
	authResp, err := h.service.Login(r.Context(), &req)
	if err != nil {
		handleError(w, err)
		return
	}

	response.OK(w, "Login successful", authResp)
}

// Logout handles partner logout
// POST /api/v1/partner/auth/logout
func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	// In a full implementation, we would revoke the refresh token
	// For now, just return success (frontend will clear tokens)
	response.OK(w, "Logged out successfully", nil)
}

// ForgotPassword handles forgot password request
// POST /api/v1/partner/auth/forgot-password
func (h *Handler) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	var req ForgotPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid request body")
		return
	}

	// Validate request
	if errors := h.validator.Validate(&req); errors != nil {
		response.UnprocessableEntity(w, "Validation failed", errors)
		return
	}

	// Send reset email
	if err := h.service.ForgotPassword(r.Context(), &req); err != nil {
		handleError(w, err)
		return
	}

	response.OK(w, "If an account with that email exists, a password reset link has been sent", nil)
}

// ResetPassword handles password reset
// POST /api/v1/partner/auth/reset-password
func (h *Handler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var req ResetPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid request body")
		return
	}

	// Validate request
	if errors := h.validator.Validate(&req); errors != nil {
		response.UnprocessableEntity(w, "Validation failed", errors)
		return
	}

	// Reset password
	if err := h.service.ResetPassword(r.Context(), &req); err != nil {
		handleError(w, err)
		return
	}

	response.OK(w, "Password reset successfully", nil)
}

// GetProfile handles get profile request
// GET /api/v1/partner/profile
func (h *Handler) GetProfile(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r)
	if userID == 0 {
		response.Unauthorized(w, "Unauthorized")
		return
	}

	partnerID := getPartnerIDFromContext(r)

	// Get partner info
	partner, err := h.service.GetReferralInfo(r.Context(), partnerID)
	if err != nil {
		handleError(w, err)
		return
	}

	// Get payout info for bank account details
	payoutInfo, err := h.service.GetPayoutInfo(r.Context(), partnerID)
	if err != nil {
		handleError(w, err)
		return
	}

	// Construct profile response
	profileResp := map[string]interface{}{
		"id":             userID,
		"name":           getNameFromContext(r),
		"email":          getEmailFromContext(r),
		"referral_code":  partner.ReferralCode,
		"account_status": "active", // TODO: Get from partner status
		"created_at":     "",       // TODO: Get from user
	}

	if payoutInfo.BankAccount != nil {
		profileResp["bank_account"] = payoutInfo.BankAccount
	}

	response.OK(w, "Profile retrieved successfully", profileResp)
}

// UpdateProfile handles profile update
// PATCH /api/v1/partner/profile
func (h *Handler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r)
	if userID == 0 {
		response.Unauthorized(w, "Unauthorized")
		return
	}

	var req UpdateProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid request body")
		return
	}

	if errors := h.validator.Validate(&req); errors != nil {
		response.UnprocessableEntity(w, "Validation failed", errors)
		return
	}

	if err := h.service.UpdateProfile(r.Context(), userID, &req); err != nil {
		handleError(w, err)
		return
	}

	response.OK(w, "Profile updated successfully", nil)
}

// ChangePassword handles password change
// POST /api/v1/partner/password/change
func (h *Handler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r)
	if userID == 0 {
		response.Unauthorized(w, "Unauthorized")
		return
	}

	var req ChangePasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid request body")
		return
	}

	if errors := h.validator.Validate(&req); errors != nil {
		response.UnprocessableEntity(w, "Validation failed", errors)
		return
	}

	if err := h.service.ChangePassword(r.Context(), userID, &req); err != nil {
		handleError(w, err)
		return
	}

	response.OK(w, "Password changed successfully", nil)
}

// GetDashboardStats handles dashboard stats request
// GET /api/v1/partner/dashboard/stats
func (h *Handler) GetDashboardStats(w http.ResponseWriter, r *http.Request) {
	partnerID := getPartnerIDFromContext(r)
	if partnerID == 0 {
		response.Unauthorized(w, "Unauthorized")
		return
	}

	stats, err := h.service.GetDashboardStats(r.Context(), partnerID)
	if err != nil {
		handleError(w, err)
		return
	}

	response.OK(w, "Dashboard stats retrieved successfully", stats)
}

// GetMonthlyData handles monthly chart data request
// GET /api/v1/partner/dashboard/monthly
func (h *Handler) GetMonthlyData(w http.ResponseWriter, r *http.Request) {
	partnerID := getPartnerIDFromContext(r)
	if partnerID == 0 {
		response.Unauthorized(w, "Unauthorized")
		return
	}

	data, err := h.service.GetMonthlyData(r.Context(), partnerID)
	if err != nil {
		handleError(w, err)
		return
	}

	response.OK(w, "Monthly data retrieved successfully", data)
}

// GetReferralInfo handles referral info request
// GET /api/v1/partner/referral
func (h *Handler) GetReferralInfo(w http.ResponseWriter, r *http.Request) {
	partnerID := getPartnerIDFromContext(r)
	if partnerID == 0 {
		response.Unauthorized(w, "Unauthorized")
		return
	}

	info, err := h.service.GetReferralInfo(r.Context(), partnerID)
	if err != nil {
		handleError(w, err)
		return
	}

	response.OK(w, "Referral info retrieved successfully", info)
}

// GetCompanies handles companies list request
// GET /api/v1/partner/companies
func (h *Handler) GetCompanies(w http.ResponseWriter, r *http.Request) {
	partnerID := getPartnerIDFromContext(r)
	if partnerID == 0 {
		response.Unauthorized(w, "Unauthorized")
		return
	}

	page := getIntQuery(r, "page", 1)
	limit := getIntQuery(r, "limit", 15)
	search := r.URL.Query().Get("search")

	companies, pagination, err := h.service.GetCompanies(r.Context(), partnerID, page, limit, search)
	if err != nil {
		handleError(w, err)
		return
	}

	response.OK(w, "Companies retrieved successfully", map[string]interface{}{
		"companies":  companies,
		"pagination": pagination,
	})
}

// GetCompaniesSummary handles companies summary request
// GET /api/v1/partner/companies/summary
func (h *Handler) GetCompaniesSummary(w http.ResponseWriter, r *http.Request) {
	partnerID := getPartnerIDFromContext(r)
	if partnerID == 0 {
		response.Unauthorized(w, "Unauthorized")
		return
	}

	summary, err := h.service.GetCompaniesSummary(r.Context(), partnerID)
	if err != nil {
		handleError(w, err)
		return
	}

	response.OK(w, "Companies summary retrieved successfully", summary)
}

// GetTransactions handles transactions list request
// GET /api/v1/partner/transactions
func (h *Handler) GetTransactions(w http.ResponseWriter, r *http.Request) {
	partnerID := getPartnerIDFromContext(r)
	if partnerID == 0 {
		response.Unauthorized(w, "Unauthorized")
		return
	}

	page := getIntQuery(r, "page", 1)
	limit := getIntQuery(r, "limit", 15)
	search := r.URL.Query().Get("search")
	status := r.URL.Query().Get("status")

	transactions, pagination, err := h.service.GetTransactions(r.Context(), partnerID, page, limit, search, status)
	if err != nil {
		handleError(w, err)
		return
	}

	response.OK(w, "Transactions retrieved successfully", map[string]interface{}{
		"transactions": transactions,
		"pagination":   pagination,
	})
}

// GetTransactionsSummary handles transactions summary request
// GET /api/v1/partner/transactions/summary
func (h *Handler) GetTransactionsSummary(w http.ResponseWriter, r *http.Request) {
	partnerID := getPartnerIDFromContext(r)
	if partnerID == 0 {
		response.Unauthorized(w, "Unauthorized")
		return
	}

	summary, err := h.service.GetTransactionsSummary(r.Context(), partnerID)
	if err != nil {
		handleError(w, err)
		return
	}

	response.OK(w, "Transactions summary retrieved successfully", summary)
}

// GetPayoutInfo handles payout info request
// GET /api/v1/partner/payouts
func (h *Handler) GetPayoutInfo(w http.ResponseWriter, r *http.Request) {
	partnerID := getPartnerIDFromContext(r)
	if partnerID == 0 {
		response.Unauthorized(w, "Unauthorized")
		return
	}

	info, err := h.service.GetPayoutInfo(r.Context(), partnerID)
	if err != nil {
		handleError(w, err)
		return
	}

	response.OK(w, "Payout info retrieved successfully", info)
}

// GetPayoutHistory handles payout history request
// GET /api/v1/partner/payouts/history
func (h *Handler) GetPayoutHistory(w http.ResponseWriter, r *http.Request) {
	partnerID := getPartnerIDFromContext(r)
	if partnerID == 0 {
		response.Unauthorized(w, "Unauthorized")
		return
	}

	page := getIntQuery(r, "page", 1)
	limit := getIntQuery(r, "limit", 15)

	payouts, pagination, err := h.service.GetPayoutHistory(r.Context(), partnerID, page, limit)
	if err != nil {
		handleError(w, err)
		return
	}

	response.OK(w, "Payout history retrieved successfully", map[string]interface{}{
		"payouts":    payouts,
		"pagination": pagination,
	})
}

// Helper functions

func getUserIDFromContext(r *http.Request) uint64 {
	if userID, ok := r.Context().Value(middleware.UserIDKey).(uint64); ok {
		return userID
	}
	if userID, ok := r.Context().Value("user_id").(uint64); ok {
		return userID
	}
	return 0
}

func getPartnerIDFromContext(r *http.Request) uint64 {
	if partnerID, ok := r.Context().Value("partner_id").(uint64); ok {
		return partnerID
	}
	return 0
}

func getNameFromContext(r *http.Request) string {
	if name, ok := r.Context().Value("user_name").(string); ok {
		return name
	}
	return ""
}

func getEmailFromContext(r *http.Request) string {
	if email, ok := r.Context().Value(middleware.UserEmailKey).(string); ok {
		return email
	}
	if email, ok := r.Context().Value("user_email").(string); ok {
		return email
	}
	return ""
}

func getIntQuery(r *http.Request, key string, defaultVal int) int {
	val := r.URL.Query().Get(key)
	if val == "" {
		return defaultVal
	}
	intVal, err := strconv.Atoi(val)
	if err != nil {
		return defaultVal
	}
	return intVal
}

func handleError(w http.ResponseWriter, err error) {
	// Simple error handling - in production, use proper error types
	response.InternalServerError(w, err.Error())
}
