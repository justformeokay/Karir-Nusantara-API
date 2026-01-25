package wishlist

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/karirnusantara/api/internal/middleware"
	apperrors "github.com/karirnusantara/api/internal/shared/errors"
	"github.com/karirnusantara/api/internal/shared/hashid"
	"github.com/karirnusantara/api/internal/shared/response"
	"github.com/karirnusantara/api/internal/shared/validator"
)

// Handler handles HTTP requests for wishlist
type Handler struct {
	service   Service
	validator *validator.Validator
}

// NewHandler creates a new wishlist handler
func NewHandler(service Service, validator *validator.Validator) *Handler {
	return &Handler{
		service:   service,
		validator: validator,
	}
}

// parseJobID parses a job ID which can be either a numeric ID or a hash_id
func parseJobID(idStr string) (uint64, error) {
	// First try to parse as numeric ID
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err == nil {
		return id, nil
	}

	// If not numeric, try to decode as hash_id (starts with "kn_")
	if strings.HasPrefix(idStr, "kn_") {
		id, err = hashid.Decode(idStr)
		if err != nil {
			return 0, err
		}
		return id, nil
	}

	return 0, err
}

// SaveJob handles saving a job to wishlist
func (h *Handler) SaveJob(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())

	var req SaveJobRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid request body")
		return
	}

	// Parse hash_id if provided
	if req.HashID != "" && req.JobID == 0 {
		decoded, err := hashid.Decode(req.HashID)
		if err != nil {
			response.BadRequest(w, "Invalid hash_id")
			return
		}
		req.JobID = decoded
	}

	// Validate that we have a job_id
	if req.JobID == 0 {
		response.BadRequest(w, "job_id or hash_id is required")
		return
	}

	saved, err := h.service.SaveJob(r.Context(), userID, &req)
	if err != nil {
		handleError(w, err)
		return
	}

	response.Created(w, "Job saved to wishlist", saved)
}

// RemoveJob handles removing a job from wishlist
func (h *Handler) RemoveJob(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())

	jobID, err := parseJobID(chi.URLParam(r, "jobId"))
	if err != nil {
		response.BadRequest(w, "Invalid job ID")
		return
	}

	if err := h.service.RemoveJob(r.Context(), userID, jobID); err != nil {
		handleError(w, err)
		return
	}

	response.OK(w, "Job removed from wishlist", nil)
}

// ListSavedJobs handles listing saved jobs
func (h *Handler) ListSavedJobs(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())

	params := ListParams{
		Page:    1,
		PerPage: 20,
	}

	if p := r.URL.Query().Get("page"); p != "" {
		if page, err := strconv.Atoi(p); err == nil && page > 0 {
			params.Page = page
		}
	}
	if pp := r.URL.Query().Get("per_page"); pp != "" {
		if perPage, err := strconv.Atoi(pp); err == nil && perPage > 0 {
			params.PerPage = perPage
		}
	}

	items, total, err := h.service.ListSavedJobs(r.Context(), userID, params)
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

	response.SuccessWithMeta(w, http.StatusOK, "Saved jobs retrieved", items, meta)
}

// CheckSaved handles checking if a job is saved
func (h *Handler) CheckSaved(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())

	jobID, err := parseJobID(chi.URLParam(r, "jobId"))
	if err != nil {
		response.BadRequest(w, "Invalid job ID")
		return
	}

	isSaved, err := h.service.IsSaved(r.Context(), userID, jobID)
	if err != nil {
		handleError(w, err)
		return
	}

	response.OK(w, "Check result", map[string]bool{"is_saved": isSaved})
}

// GetStats handles getting wishlist statistics
func (h *Handler) GetStats(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())

	stats, err := h.service.GetStats(r.Context(), userID)
	if err != nil {
		handleError(w, err)
		return
	}

	response.OK(w, "Wishlist stats", stats)
}

// handleError handles service errors and returns appropriate HTTP response
func handleError(w http.ResponseWriter, err error) {
	if appErr := apperrors.GetAppError(err); appErr != nil {
		switch appErr.Code {
		case apperrors.ErrCodeNotFound:
			response.NotFound(w, appErr.Message)
		case apperrors.ErrCodeConflict:
			response.Conflict(w, appErr.Message)
		case apperrors.ErrCodeValidation:
			response.UnprocessableEntity(w, appErr.Message, appErr.Details)
		default:
			response.InternalServerError(w, "An error occurred")
		}
		return
	}
	response.InternalServerError(w, "An error occurred")
}
