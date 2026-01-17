package applications

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/karirnusantara/api/internal/middleware"
	apperrors "github.com/karirnusantara/api/internal/shared/errors"
	"github.com/karirnusantara/api/internal/shared/response"
	"github.com/karirnusantara/api/internal/shared/validator"
)

// Handler handles HTTP requests for applications
type Handler struct {
	service   Service
	validator *validator.Validator
}

// NewHandler creates a new applications handler
func NewHandler(service Service, validator *validator.Validator) *Handler {
	return &Handler{
		service:   service,
		validator: validator,
	}
}

// Apply handles job application submission
func (h *Handler) Apply(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())

	var req ApplyJobRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid request body")
		return
	}

	if errors := h.validator.Validate(&req); errors != nil {
		response.UnprocessableEntity(w, "Validation failed", errors)
		return
	}

	app, err := h.service.Apply(r.Context(), userID, &req)
	if err != nil {
		handleError(w, err)
		return
	}

	response.Created(w, "Application submitted successfully", app)
}

// GetByID handles getting a single application
func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.BadRequest(w, "Invalid application ID")
		return
	}

	userID := middleware.GetUserID(r.Context())
	userRole := middleware.GetUserRole(r.Context())
	isCompany := userRole == "company"

	app, err := h.service.GetByID(r.Context(), id, userID, isCompany)
	if err != nil {
		handleError(w, err)
		return
	}

	response.OK(w, "Application retrieved", app)
}

// ListMyApplications handles listing applications for the authenticated user
func (h *Handler) ListMyApplications(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	params := h.parseListParams(r)

	apps, total, err := h.service.ListByUser(r.Context(), userID, params)
	if err != nil {
		handleError(w, err)
		return
	}

	meta := &response.Meta{
		Page:       params.Page,
		PerPage:    params.PerPage,
		TotalItems: total,
		TotalPages: int((total + int64(params.PerPage) - 1) / int64(params.PerPage)),
	}

	response.SuccessWithMeta(w, http.StatusOK, "Applications retrieved", apps, meta)
}

// ListCompanyApplications handles listing applications for a company
func (h *Handler) ListCompanyApplications(w http.ResponseWriter, r *http.Request) {
	companyID := middleware.GetUserID(r.Context())
	params := h.parseListParams(r)

	apps, total, err := h.service.ListByCompany(r.Context(), companyID, params)
	if err != nil {
		handleError(w, err)
		return
	}

	meta := &response.Meta{
		Page:       params.Page,
		PerPage:    params.PerPage,
		TotalItems: total,
		TotalPages: int((total + int64(params.PerPage) - 1) / int64(params.PerPage)),
	}

	response.SuccessWithMeta(w, http.StatusOK, "Applications retrieved", apps, meta)
}

// ListJobApplications handles listing applications for a specific job
func (h *Handler) ListJobApplications(w http.ResponseWriter, r *http.Request) {
	jobID, err := strconv.ParseUint(chi.URLParam(r, "jobId"), 10, 64)
	if err != nil {
		response.BadRequest(w, "Invalid job ID")
		return
	}

	companyID := middleware.GetUserID(r.Context())
	params := h.parseListParams(r)

	apps, total, err := h.service.ListByJob(r.Context(), jobID, companyID, params)
	if err != nil {
		handleError(w, err)
		return
	}

	meta := &response.Meta{
		Page:       params.Page,
		PerPage:    params.PerPage,
		TotalItems: total,
		TotalPages: int((total + int64(params.PerPage) - 1) / int64(params.PerPage)),
	}

	response.SuccessWithMeta(w, http.StatusOK, "Applications retrieved", apps, meta)
}

// UpdateStatus handles updating application status by company
func (h *Handler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.BadRequest(w, "Invalid application ID")
		return
	}

	companyID := middleware.GetUserID(r.Context())

	var req UpdateStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid request body")
		return
	}

	if errors := h.validator.Validate(&req); errors != nil {
		response.UnprocessableEntity(w, "Validation failed", errors)
		return
	}

	app, err := h.service.UpdateStatus(r.Context(), id, companyID, &req)
	if err != nil {
		handleError(w, err)
		return
	}

	response.OK(w, "Application status updated", app)
}

// Withdraw handles application withdrawal by applicant
func (h *Handler) Withdraw(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.BadRequest(w, "Invalid application ID")
		return
	}

	userID := middleware.GetUserID(r.Context())

	var req WithdrawRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		// Empty body is acceptable
		req = WithdrawRequest{}
	}

	if err := h.service.Withdraw(r.Context(), id, userID, req.Reason); err != nil {
		handleError(w, err)
		return
	}

	response.OK(w, "Application withdrawn successfully", nil)
}

// GetTimeline handles getting application timeline
func (h *Handler) GetTimeline(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.BadRequest(w, "Invalid application ID")
		return
	}

	userID := middleware.GetUserID(r.Context())
	userRole := middleware.GetUserRole(r.Context())
	isCompany := userRole == "company"

	app, err := h.service.GetByID(r.Context(), id, userID, isCompany)
	if err != nil {
		handleError(w, err)
		return
	}

	response.OK(w, "Timeline retrieved", app.Timeline)
}

// parseListParams parses query parameters for listing
func (h *Handler) parseListParams(r *http.Request) ApplicationListParams {
	params := ApplicationListParams{
		Page:    1,
		PerPage: 20,
	}

	if page := r.URL.Query().Get("page"); page != "" {
		if p, err := strconv.Atoi(page); err == nil && p > 0 {
			params.Page = p
		}
	}

	if limit := r.URL.Query().Get("limit"); limit != "" {
		if l, err := strconv.Atoi(limit); err == nil && l > 0 && l <= 100 {
			params.PerPage = l
		}
	}

	if status := r.URL.Query().Get("status"); status != "" {
		params.Status = status
	}

	if sortBy := r.URL.Query().Get("sort_by"); sortBy != "" {
		params.SortBy = sortBy
	}

	if sortOrder := r.URL.Query().Get("sort_order"); sortOrder != "" {
		params.SortOrder = sortOrder
	}

	return params
}

// handleError handles errors and sends appropriate response
func handleError(w http.ResponseWriter, err error) {
	appErr := apperrors.GetAppError(err)
	if appErr != nil {
		if appErr.Details != nil {
			response.ErrorWithDetails(w, appErr.HTTPStatus, appErr.Code, appErr.Message, appErr.Details)
		} else {
			response.Error(w, appErr.HTTPStatus, appErr.Code, appErr.Message)
		}
		return
	}
	response.InternalServerError(w, "An unexpected error occurred")
}
