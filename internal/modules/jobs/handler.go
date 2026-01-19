package jobs

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

// Handler handles job HTTP requests
type Handler struct {
	service   Service
	validator *validator.Validator
}

// NewHandler creates a new job handler
func NewHandler(service Service, validator *validator.Validator) *Handler {
	return &Handler{
		service:   service,
		validator: validator,
	}
}

// Create handles job creation
// POST /api/v1/jobs
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	companyID := middleware.GetUserID(r.Context())
	if companyID == 0 {
		response.Unauthorized(w, "Unauthorized")
		return
	}

	var req CreateJobRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid request body")
		return
	}

	if errors := h.validator.Validate(&req); errors != nil {
		response.UnprocessableEntity(w, "Validation failed", errors)
		return
	}

	job, err := h.service.Create(r.Context(), companyID, &req)
	if err != nil {
		handleError(w, err)
		return
	}

	response.Created(w, "Job created successfully", job)
}

// GetByID handles getting a job by ID
// GET /api/v1/jobs/{id}
func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(w, "Invalid job ID")
		return
	}

	job, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		handleError(w, err)
		return
	}

	// Increment view count (async would be better in production)
	go h.service.IncrementViewCount(r.Context(), id)

	response.OK(w, "Job retrieved", job)
}

// GetBySlug handles getting a job by slug
// GET /api/v1/jobs/slug/{slug}
func (h *Handler) GetBySlug(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	if slug == "" {
		response.BadRequest(w, "Slug is required")
		return
	}

	job, err := h.service.GetBySlug(r.Context(), slug)
	if err != nil {
		handleError(w, err)
		return
	}

	response.OK(w, "Job retrieved", job)
}

// Update handles job update
// PUT /api/v1/jobs/{id}
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	companyID := middleware.GetUserID(r.Context())
	if companyID == 0 {
		response.Unauthorized(w, "Unauthorized")
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(w, "Invalid job ID")
		return
	}

	var req UpdateJobRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid request body")
		return
	}

	if errors := h.validator.Validate(&req); errors != nil {
		response.UnprocessableEntity(w, "Validation failed", errors)
		return
	}

	job, err := h.service.Update(r.Context(), id, companyID, &req)
	if err != nil {
		handleError(w, err)
		return
	}

	response.OK(w, "Job updated successfully", job)
}

// Delete handles job deletion
// DELETE /api/v1/jobs/{id}
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	companyID := middleware.GetUserID(r.Context())
	if companyID == 0 {
		response.Unauthorized(w, "Unauthorized")
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(w, "Invalid job ID")
		return
	}

	if err := h.service.Delete(r.Context(), id, companyID); err != nil {
		handleError(w, err)
		return
	}

	response.NoContent(w)
}

// List handles job listing with filters
// GET /api/v1/jobs
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	params := DefaultJobListParams()

	// Parse query parameters
	query := r.URL.Query()

	if page := query.Get("page"); page != "" {
		if p, err := strconv.Atoi(page); err == nil && p > 0 {
			params.Page = p
		}
	}
	if perPage := query.Get("per_page"); perPage != "" {
		if pp, err := strconv.Atoi(perPage); err == nil && pp > 0 && pp <= 100 {
			params.PerPage = pp
		}
	}
	if search := query.Get("search"); search != "" {
		params.Search = search
	}
	if city := query.Get("city"); city != "" {
		params.City = city
	}
	if province := query.Get("province"); province != "" {
		params.Province = province
	}
	if jobType := query.Get("job_type"); jobType != "" {
		params.JobType = jobType
	}
	if expLevel := query.Get("experience_level"); expLevel != "" {
		params.ExperienceLevel = expLevel
	}
	if isRemote := query.Get("is_remote"); isRemote != "" {
		remote := isRemote == "true"
		params.IsRemote = &remote
	}
	if salaryMin := query.Get("salary_min"); salaryMin != "" {
		if min, err := strconv.ParseInt(salaryMin, 10, 64); err == nil {
			params.SalaryMin = &min
		}
	}
	if salaryMax := query.Get("salary_max"); salaryMax != "" {
		if max, err := strconv.ParseInt(salaryMax, 10, 64); err == nil {
			params.SalaryMax = &max
		}
	}
	if sortBy := query.Get("sort_by"); sortBy != "" {
		params.SortBy = sortBy
	}
	if sortOrder := query.Get("sort_order"); sortOrder != "" {
		params.SortOrder = sortOrder
	}

	jobs, total, err := h.service.List(r.Context(), params)
	if err != nil {
		handleError(w, err)
		return
	}

	totalPages := int(total) / params.PerPage
	if int(total)%params.PerPage > 0 {
		totalPages++
	}

	meta := &response.Meta{
		Page:       params.Page,
		PerPage:    params.PerPage,
		TotalItems: total,
		TotalPages: totalPages,
	}

	response.SuccessWithMeta(w, http.StatusOK, "Jobs retrieved", jobs, meta)
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

// Publish handles publishing a draft job
// PATCH /api/v1/jobs/{id}/publish
func (h *Handler) Publish(w http.ResponseWriter, r *http.Request) {
	companyID := middleware.GetUserID(r.Context())
	if companyID == 0 {
		response.Unauthorized(w, "Unauthorized")
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(w, "Invalid job ID")
		return
	}

	job, err := h.service.UpdateStatus(r.Context(), id, companyID, JobStatusActive)
	if err != nil {
		handleError(w, err)
		return
	}

	response.OK(w, "Job published successfully", job)
}

// Close handles closing an active job
// PATCH /api/v1/jobs/{id}/close
func (h *Handler) Close(w http.ResponseWriter, r *http.Request) {
	companyID := middleware.GetUserID(r.Context())
	if companyID == 0 {
		response.Unauthorized(w, "Unauthorized")
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(w, "Invalid job ID")
		return
	}

	job, err := h.service.UpdateStatus(r.Context(), id, companyID, JobStatusClosed)
	if err != nil {
		handleError(w, err)
		return
	}

	response.OK(w, "Job closed successfully", job)
}

// Pause handles pausing an active job
// PATCH /api/v1/jobs/{id}/pause
func (h *Handler) Pause(w http.ResponseWriter, r *http.Request) {
	companyID := middleware.GetUserID(r.Context())
	if companyID == 0 {
		response.Unauthorized(w, "Unauthorized")
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(w, "Invalid job ID")
		return
	}

	job, err := h.service.UpdateStatus(r.Context(), id, companyID, JobStatusPaused)
	if err != nil {
		handleError(w, err)
		return
	}

	response.OK(w, "Job paused successfully", job)
}

// Reopen handles reopening a closed/paused job
// PATCH /api/v1/jobs/{id}/reopen
func (h *Handler) Reopen(w http.ResponseWriter, r *http.Request) {
	companyID := middleware.GetUserID(r.Context())
	if companyID == 0 {
		response.Unauthorized(w, "Unauthorized")
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(w, "Invalid job ID")
		return
	}

	job, err := h.service.UpdateStatus(r.Context(), id, companyID, JobStatusActive)
	if err != nil {
		handleError(w, err)
		return
	}

	response.OK(w, "Job reopened successfully", job)
}

// ListByCompany handles listing jobs for the authenticated company
// GET /api/v1/company/jobs
func (h *Handler) ListByCompany(w http.ResponseWriter, r *http.Request) {
	companyID := middleware.GetUserID(r.Context())
	if companyID == 0 {
		response.Unauthorized(w, "Unauthorized")
		return
	}

	params := DefaultJobListParams()

	// Parse query parameters
	query := r.URL.Query()

	if page := query.Get("page"); page != "" {
		if p, err := strconv.Atoi(page); err == nil && p > 0 {
			params.Page = p
		}
	}
	if perPage := query.Get("per_page"); perPage != "" {
		if pp, err := strconv.Atoi(perPage); err == nil && pp > 0 && pp <= 100 {
			params.PerPage = pp
		}
	}
	if search := query.Get("search"); search != "" {
		params.Search = search
	}
	if sortBy := query.Get("sort_by"); sortBy != "" {
		params.SortBy = sortBy
	}
	if sortOrder := query.Get("sort_order"); sortOrder != "" {
		params.SortOrder = sortOrder
	}

	jobs, total, err := h.service.ListByCompany(r.Context(), companyID, params)
	if err != nil {
		handleError(w, err)
		return
	}

	totalPages := int(total) / params.PerPage
	if int(total)%params.PerPage > 0 {
		totalPages++
	}

	meta := &response.Meta{
		Page:       params.Page,
		PerPage:    params.PerPage,
		TotalItems: total,
		TotalPages: totalPages,
	}

	response.SuccessWithMeta(w, http.StatusOK, "Jobs retrieved", jobs, meta)
}
