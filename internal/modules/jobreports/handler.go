package jobreports

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/karirnusantara/api/internal/middleware"
	"github.com/karirnusantara/api/internal/shared/response"
	"github.com/karirnusantara/api/internal/shared/validator"
)

// Handler handles job report HTTP requests
type Handler struct {
	service   Service
	validator *validator.Validator
}

// NewHandler creates a new job reports handler
func NewHandler(service Service, validator *validator.Validator) *Handler {
	return &Handler{
		service:   service,
		validator: validator,
	}
}

// CreateReport handles creating a new job report
// @Summary Report a job listing
// @Tags JobReports
// @Accept json
// @Produce json
// @Param jobId path string true "Job ID or Hash ID"
// @Param request body CreateReportRequest true "Report details"
// @Success 201 {object} JobReportWithDetails
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 429 {object} response.ErrorResponse "Cooldown active or already reported"
// @Router /jobs/{jobId}/reports [post]
func (h *Handler) CreateReport(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == 0 {
		response.Error(w, http.StatusUnauthorized, "unauthorized", "Anda harus login terlebih dahulu")
		return
	}

	jobIDStr := chi.URLParam(r, "jobId")
	if jobIDStr == "" {
		response.Error(w, http.StatusBadRequest, "invalid_request", "Job ID diperlukan")
		return
	}

	// Try to parse as numeric ID first
	jobID, err := strconv.ParseUint(jobIDStr, 10, 64)
	if err != nil {
		// If not numeric, try to resolve hash_id - this would need jobs service
		// For now, we'll just use numeric IDs
		response.Error(w, http.StatusBadRequest, "invalid_request", "Job ID tidak valid")
		return
	}

	var req CreateReportRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid_request", "Format request tidak valid")
		return
	}

	if errors := h.validator.Validate(&req); errors != nil {
		response.UnprocessableEntity(w, "Validation failed", errors)
		return
	}

	report, err := h.service.CreateReport(r.Context(), userID, jobID, &req)
	if err != nil {
		// Check if it's a validation error (cooldown, already reported)
		errMsg := err.Error()
		if errMsg[:4] == "Anda" {
			response.Error(w, http.StatusTooManyRequests, "report_limit", errMsg)
			return
		}
		response.Error(w, http.StatusInternalServerError, "create_failed", errMsg)
		return
	}

	response.JSON(w, http.StatusCreated, map[string]interface{}{
		"message": "Terima kasih atas laporan Anda. Tim kami akan meninjau laporan ini.",
		"report":  report,
	})
}

// CheckReportStatus checks if user has reported a specific job
// @Summary Check if user has reported a job
// @Tags JobReports
// @Produce json
// @Param jobId path string true "Job ID"
// @Success 200 {object} map[string]bool
// @Router /jobs/{jobId}/reports/status [get]
func (h *Handler) CheckReportStatus(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == 0 {
		response.JSON(w, http.StatusOK, map[string]interface{}{
			"has_reported": false,
		})
		return
	}

	jobIDStr := chi.URLParam(r, "jobId")
	jobID, err := strconv.ParseUint(jobIDStr, 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid_request", "Job ID tidak valid")
		return
	}

	hasReported, err := h.service.HasUserReportedJob(r.Context(), userID, jobID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "check_failed", err.Error())
		return
	}

	response.JSON(w, http.StatusOK, map[string]interface{}{
		"has_reported": hasReported,
	})
}

// AdminGetAllReports handles getting all reports (admin only)
// @Summary Get all job reports (Admin)
// @Tags JobReports
// @Produce json
// @Param status query string false "Filter by status"
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Success 200 {object} map[string]interface{}
// @Router /admin/job-reports [get]
func (h *Handler) AdminGetAllReports(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")

	page := 1
	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil {
			page = p
		}
	}

	limit := 20
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			limit = l
		}
	}

	reports, total, err := h.service.GetAllReports(r.Context(), status, page, limit)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "fetch_failed", err.Error())
		return
	}

	response.JSON(w, http.StatusOK, map[string]interface{}{
		"reports": reports,
		"total":   total,
		"page":    page,
		"limit":   limit,
	})
}

// AdminGetReport handles getting a single report (admin only)
// @Summary Get a job report (Admin)
// @Tags JobReports
// @Produce json
// @Param id path int true "Report ID"
// @Success 200 {object} JobReportWithDetails
// @Router /admin/job-reports/{id} [get]
func (h *Handler) AdminGetReport(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid_request", "Report ID tidak valid")
		return
	}

	report, err := h.service.GetReportByID(r.Context(), id)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "fetch_failed", err.Error())
		return
	}
	if report == nil {
		response.Error(w, http.StatusNotFound, "not_found", "Laporan tidak ditemukan")
		return
	}

	response.JSON(w, http.StatusOK, map[string]interface{}{
		"report": report,
	})
}

// AdminUpdateReportStatus handles updating report status (admin only)
// @Summary Update job report status (Admin)
// @Tags JobReports
// @Accept json
// @Produce json
// @Param id path int true "Report ID"
// @Param request body UpdateReportStatusRequest true "Status update"
// @Success 200 {object} map[string]string
// @Router /admin/job-reports/{id}/status [patch]
func (h *Handler) AdminUpdateReportStatus(w http.ResponseWriter, r *http.Request) {
	adminID := middleware.GetUserID(r.Context())
	if adminID == 0 {
		response.Error(w, http.StatusUnauthorized, "unauthorized", "Unauthorized")
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid_request", "Report ID tidak valid")
		return
	}

	var req UpdateReportStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid_request", "Format request tidak valid")
		return
	}

	if errors := h.validator.Validate(&req); errors != nil {
		response.UnprocessableEntity(w, "Validation failed", errors)
		return
	}

	if err := h.service.UpdateReportStatus(r.Context(), id, &req, adminID); err != nil {
		response.Error(w, http.StatusInternalServerError, "update_failed", err.Error())
		return
	}

	response.JSON(w, http.StatusOK, map[string]interface{}{
		"message": "Status laporan berhasil diperbarui",
	})
}

// AdminGetReportsByJob handles getting all reports for a specific job (admin only)
// @Summary Get reports for a job (Admin)
// @Tags JobReports
// @Produce json
// @Param jobId path int true "Job ID"
// @Success 200 {object} map[string]interface{}
// @Router /admin/jobs/{jobId}/reports [get]
func (h *Handler) AdminGetReportsByJob(w http.ResponseWriter, r *http.Request) {
	jobIDStr := chi.URLParam(r, "jobId")
	jobID, err := strconv.ParseUint(jobIDStr, 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid_request", "Job ID tidak valid")
		return
	}

	reports, err := h.service.GetReportsByJobID(r.Context(), jobID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "fetch_failed", err.Error())
		return
	}

	response.JSON(w, http.StatusOK, map[string]interface{}{
		"reports": reports,
		"total":   len(reports),
	})
}

// AdminGetPendingCount handles getting count of pending reports (admin only)
// @Summary Get pending reports count (Admin)
// @Tags JobReports
// @Produce json
// @Success 200 {object} map[string]int
// @Router /admin/job-reports/pending-count [get]
func (h *Handler) AdminGetPendingCount(w http.ResponseWriter, r *http.Request) {
	count, err := h.service.GetPendingReportsCount(r.Context())
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "fetch_failed", err.Error())
		return
	}

	response.JSON(w, http.StatusOK, map[string]interface{}{
		"pending_count": count,
	})
}

// AdminBanCompany handles banning a company (admin only)
// @Summary Ban a company (Admin)
// @Tags JobReports
// @Accept json
// @Produce json
// @Param request body BanCompanyRequest true "Ban details"
// @Success 200 {object} map[string]string
// @Router /admin/job-reports/ban-company [post]
func (h *Handler) AdminBanCompany(w http.ResponseWriter, r *http.Request) {
	adminID := middleware.GetUserID(r.Context())
	if adminID == 0 {
		response.Error(w, http.StatusUnauthorized, "unauthorized", "Anda harus login sebagai admin")
		return
	}

	var req BanCompanyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid_request", "Format request tidak valid")
		return
	}

	if errors := h.validator.Validate(&req); errors != nil {
		response.UnprocessableEntity(w, "Validation failed", errors)
		return
	}

	err := h.service.BanCompany(r.Context(), req.CompanyID, adminID, req.Reason)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "ban_failed", err.Error())
		return
	}

	response.JSON(w, http.StatusOK, map[string]interface{}{
		"message": "Perusahaan berhasil dibanned",
	})
}
