package dashboard

import (
	"net/http"
	"strconv"

	"github.com/karirnusantara/api/internal/middleware"
	"github.com/karirnusantara/api/internal/shared/response"
)

// Handler handles HTTP requests for dashboard
type Handler struct {
	service *Service
}

// NewHandler creates a new dashboard handler
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// GetStats godoc
// @Summary Get dashboard statistics
// @Description Get comprehensive dashboard statistics for company
// @Tags Dashboard
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response{data=DashboardStats}
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /company/dashboard/stats [get]
func (h *Handler) GetStats(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context (set by auth middleware)
	userID := middleware.GetUserID(r.Context())
	if userID == 0 {
		response.Error(w, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
		return
	}

	// Lookup company_id from user_id
	companyID, err := h.service.GetCompanyIDByUserID(userID)
	if err != nil {
		response.Error(w, http.StatusNotFound, "COMPANY_NOT_FOUND", "Company not found for this user")
		return
	}

	stats, err := h.service.GetDashboardStats(companyID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "DASHBOARD_ERROR", "Failed to get dashboard stats")
		return
	}

	response.Success(w, http.StatusOK, "Dashboard stats retrieved successfully", stats)
}

// GetRecentApplicants godoc
// @Summary Get recent applicants
// @Description Get list of recent applicants for company jobs
// @Tags Dashboard
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param limit query int false "Limit results" default(10)
// @Success 200 {object} response.Response{data=[]RecentApplicant}
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /company/dashboard/recent-applicants [get]
func (h *Handler) GetRecentApplicants(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == 0 {
		response.Error(w, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
		return
	}

	// Lookup company_id from user_id
	companyID, err := h.service.GetCompanyIDByUserID(userID)
	if err != nil {
		response.Error(w, http.StatusNotFound, "COMPANY_NOT_FOUND", "Company not found for this user")
		return
	}

	limit := 10
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	applicants, err := h.service.GetRecentApplicants(companyID, limit)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "DASHBOARD_ERROR", "Failed to get recent applicants")
		return
	}

	response.Success(w, http.StatusOK, "Recent applicants retrieved successfully", applicants)
}

// GetActiveJobs godoc
// @Summary Get active jobs
// @Description Get list of active jobs for company
// @Tags Dashboard
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param limit query int false "Limit results" default(10)
// @Success 200 {object} response.Response{data=[]ActiveJob}
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /company/dashboard/active-jobs [get]
func (h *Handler) GetActiveJobs(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == 0 {
		response.Error(w, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
		return
	}

	// Lookup company_id from user_id
	companyID, err := h.service.GetCompanyIDByUserID(userID)
	if err != nil {
		response.Error(w, http.StatusNotFound, "COMPANY_NOT_FOUND", "Company not found for this user")
		return
	}

	limit := 10
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	jobs, err := h.service.GetActiveJobsList(companyID, limit)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "DASHBOARD_ERROR", "Failed to get active jobs")
		return
	}

	response.Success(w, http.StatusOK, "Active jobs retrieved successfully", jobs)
}
