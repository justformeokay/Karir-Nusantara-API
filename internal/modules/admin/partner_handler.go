package admin

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/karirnusantara/api/internal/middleware"
)

// PartnerHandler handles partner-related HTTP requests for admin
type PartnerHandler struct {
	service PartnerService
}

// NewPartnerHandler creates a new partner handler for admin
func NewPartnerHandler(service PartnerService) *PartnerHandler {
	return &PartnerHandler{service: service}
}

// GetPartners godoc
// @Summary Get list of partners
// @Description Get paginated list of all partners with optional filters
// @Tags Admin Partners
// @Accept json
// @Produce json
// @Param status query string false "Filter by status (active, pending, suspended, all)"
// @Param search query string false "Search by name, email, or referral code"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} AdminPartnerListResponse
// @Security BearerAuth
// @Router /admin/partners [get]
func (h *PartnerHandler) GetPartners(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")
	search := r.URL.Query().Get("search")
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	result, err := h.service.GetPartners(r.Context(), status, search, page, limit)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    result,
	})
}

// GetPartnerByID godoc
// @Summary Get partner by ID
// @Description Get detailed information about a specific partner
// @Tags Admin Partners
// @Accept json
// @Produce json
// @Param id path int true "Partner ID"
// @Success 200 {object} AdminPartnerDetailResponse
// @Security BearerAuth
// @Router /admin/partners/{id} [get]
func (h *PartnerHandler) GetPartnerByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid partner ID")
		return
	}

	result, err := h.service.GetPartnerByID(r.Context(), id)
	if err != nil {
		if err.Error() == "partner not found" {
			respondWithError(w, http.StatusNotFound, err.Error())
			return
		}
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    result,
	})
}

// UpdatePartnerStatus godoc
// @Summary Update partner status
// @Description Update partner status (activate, suspend, etc.)
// @Tags Admin Partners
// @Accept json
// @Produce json
// @Param id path int true "Partner ID"
// @Param body body UpdatePartnerStatusRequest true "Status update request"
// @Success 200 {object} map[string]interface{}
// @Security BearerAuth
// @Router /admin/partners/{id}/status [patch]
func (h *PartnerHandler) UpdatePartnerStatus(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid partner ID")
		return
	}

	var req UpdatePartnerStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Status == "" {
		respondWithError(w, http.StatusBadRequest, "status is required")
		return
	}

	if err := h.service.UpdatePartnerStatus(r.Context(), id, req); err != nil {
		if err.Error() == "partner not found" {
			respondWithError(w, http.StatusNotFound, err.Error())
			return
		}
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Partner status updated successfully",
	})
}

// ApprovePartner godoc
// @Summary Approve partner
// @Description Approve a pending partner application
// @Tags Admin Partners
// @Accept json
// @Produce json
// @Param id path int true "Partner ID"
// @Param body body ApprovePartnerRequest true "Approval request"
// @Success 200 {object} map[string]interface{}
// @Security BearerAuth
// @Router /admin/partners/{id}/approve [post]
func (h *PartnerHandler) ApprovePartner(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid partner ID")
		return
	}

	var req ApprovePartnerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		// Allow empty body - use defaults
		req = ApprovePartnerRequest{}
	}

	// Get admin ID from context using middleware helper
	adminID := middleware.GetUserID(r.Context())

	if err := h.service.ApprovePartner(r.Context(), id, adminID, req); err != nil {
		if err.Error() == "partner not found" {
			respondWithError(w, http.StatusNotFound, err.Error())
			return
		}
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Partner approved successfully",
	})
}

// EditPartner godoc
// @Summary Edit partner
// @Description Update partner details (bank info, commission rate, notes)
// @Tags Admin Partners
// @Accept json
// @Produce json
// @Param id path int true "Partner ID"
// @Param body body EditPartnerRequest true "Edit request"
// @Success 200 {object} map[string]interface{}
// @Security BearerAuth
// @Router /admin/partners/{id} [put]
func (h *PartnerHandler) EditPartner(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid partner ID")
		return
	}

	var req EditPartnerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.service.EditPartner(r.Context(), id, req); err != nil {
		if err.Error() == "partner not found" {
			respondWithError(w, http.StatusNotFound, err.Error())
			return
		}
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Partner updated successfully",
	})
}

// RejectPartner godoc
// @Summary Reject partner
// @Description Reject a pending partner application
// @Tags Admin Partners
// @Accept json
// @Produce json
// @Param id path int true "Partner ID"
// @Param body body RejectPartnerRequest true "Rejection request"
// @Success 200 {object} map[string]interface{}
// @Security BearerAuth
// @Router /admin/partners/{id}/reject [post]
func (h *PartnerHandler) RejectPartner(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid partner ID")
		return
	}

	var req RejectPartnerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Reason == "" {
		respondWithError(w, http.StatusBadRequest, "rejection reason is required")
		return
	}

	// Get admin ID from context
	adminID := middleware.GetUserID(r.Context())

	if err := h.service.RejectPartner(r.Context(), id, adminID, req); err != nil {
		if err.Error() == "partner not found" {
			respondWithError(w, http.StatusNotFound, err.Error())
			return
		}
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Partner rejected successfully",
	})
}

// DeletePartner godoc
// @Summary Delete partner
// @Description Soft-delete a partner (deactivate account)
// @Tags Admin Partners
// @Accept json
// @Produce json
// @Param id path int true "Partner ID"
// @Param body body DeletePartnerRequest true "Delete request"
// @Success 200 {object} map[string]interface{}
// @Security BearerAuth
// @Router /admin/partners/{id} [delete]
func (h *PartnerHandler) DeletePartner(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid partner ID")
		return
	}

	var req DeletePartnerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		// Allow empty body
		req = DeletePartnerRequest{}
	}

	// Get admin ID from context
	adminID := middleware.GetUserID(r.Context())

	if err := h.service.DeletePartner(r.Context(), id, adminID, req); err != nil {
		if err.Error() == "partner not found" {
			respondWithError(w, http.StatusNotFound, err.Error())
			return
		}
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Partner deleted successfully",
	})
}

// GetReferredCompanies godoc
// @Summary Get referred companies
// @Description Get paginated list of companies referred by partners
// @Tags Admin Partners
// @Accept json
// @Produce json
// @Param search query string false "Search by company name or partner name"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} AdminReferredCompanyListResponse
// @Security BearerAuth
// @Router /admin/referrals/companies [get]
func (h *PartnerHandler) GetReferredCompanies(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	result, err := h.service.GetReferredCompanies(r.Context(), search, page, limit)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    result,
	})
}

// GetReferralStats godoc
// @Summary Get referral program statistics
// @Description Get overall statistics for the referral program
// @Tags Admin Partners
// @Accept json
// @Produce json
// @Success 200 {object} AdminReferralStatsResponse
// @Security BearerAuth
// @Router /admin/referrals/stats [get]
func (h *PartnerHandler) GetReferralStats(w http.ResponseWriter, r *http.Request) {
	result, err := h.service.GetReferralStats(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    result,
	})
}

// GetPayouts godoc
// @Summary Get payouts
// @Description Get paginated list of commission payouts
// @Tags Admin Payouts
// @Accept json
// @Produce json
// @Param status query string false "Filter by status (pending, paid, rejected, all)"
// @Param search query string false "Search by partner name"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} AdminPayoutListResponse
// @Security BearerAuth
// @Router /admin/payouts [get]
func (h *PartnerHandler) GetPayouts(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")
	search := r.URL.Query().Get("search")
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	result, err := h.service.GetPayouts(r.Context(), status, search, page, limit)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    result,
	})
}

// GetPayoutStats godoc
// @Summary Get payout statistics
// @Description Get overall payout statistics
// @Tags Admin Payouts
// @Accept json
// @Produce json
// @Success 200 {object} AdminPayoutStatsResponse
// @Security BearerAuth
// @Router /admin/payouts/stats [get]
func (h *PartnerHandler) GetPayoutStats(w http.ResponseWriter, r *http.Request) {
	result, err := h.service.GetPayoutStats(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    result,
	})
}

// GetPartnerBalances godoc
// @Summary Get partner balances
// @Description Get partners with available balance for payout
// @Tags Admin Payouts
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} AdminPartnerBalanceListResponse
// @Security BearerAuth
// @Router /admin/payouts/balances [get]
func (h *PartnerHandler) GetPartnerBalances(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	result, err := h.service.GetPartnerBalances(r.Context(), page, limit)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    result,
	})
}

// CreatePayout godoc
// @Summary Create payout
// @Description Create a new payout request for a partner
// @Tags Admin Payouts
// @Accept json
// @Produce json
// @Param body body CreatePayoutRequest true "Payout request"
// @Success 201 {object} map[string]interface{}
// @Security BearerAuth
// @Router /admin/payouts [post]
func (h *PartnerHandler) CreatePayout(w http.ResponseWriter, r *http.Request) {
	var req CreatePayoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.PartnerID == 0 {
		respondWithError(w, http.StatusBadRequest, "partner_id is required")
		return
	}

	if req.Amount <= 0 {
		respondWithError(w, http.StatusBadRequest, "amount must be greater than 0")
		return
	}

	payoutID, err := h.service.CreatePayout(r.Context(), req)
	if err != nil {
		if err.Error() == "partner not found" {
			respondWithError(w, http.StatusNotFound, err.Error())
			return
		}
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, map[string]interface{}{
		"success":   true,
		"message":   "Payout created successfully",
		"payout_id": payoutID,
	})
}

// ProcessPayout godoc
// @Summary Process payout
// @Description Mark a payout as paid
// @Tags Admin Payouts
// @Accept json
// @Produce json
// @Param id path int true "Payout ID"
// @Param body body ProcessPayoutRequest true "Process payout request"
// @Success 200 {object} map[string]interface{}
// @Security BearerAuth
// @Router /admin/payouts/{id}/process [post]
func (h *PartnerHandler) ProcessPayout(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid payout ID")
		return
	}

	var req ProcessPayoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.service.ProcessPayout(r.Context(), id, req); err != nil {
		if err.Error() == "payout not found" {
			respondWithError(w, http.StatusNotFound, err.Error())
			return
		}
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Payout processed successfully",
	})
}

// Helper functions
func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]interface{}{
		"success": false,
		"error":   message,
	})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
