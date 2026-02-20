package interviewtests

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/karirnusantara/api/internal/shared/response"
)

// Handler handles HTTP requests for interview tests
type Handler struct {
	service *Service
}

// NewHandler creates a new interview tests handler
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// Create handles POST /admin/interview-tests
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateInterviewTestRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}

	// Get admin ID from context
	adminID := getAdminIDFromContext(r)

	result, err := h.service.Create(r.Context(), req, adminID)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	response.Success(w, http.StatusCreated, "Interview test created successfully", result)
}

// GetByID handles GET /admin/interview-tests/{id}
func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_ID", "Invalid test ID")
		return
	}

	result, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	response.Success(w, http.StatusOK, "Interview test retrieved successfully", result)
}

// GetAll handles GET /admin/interview-tests
func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")

	results, err := h.service.GetAll(r.Context(), status)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	response.Success(w, http.StatusOK, "Interview tests retrieved successfully", results)
}

// Update handles PUT /admin/interview-tests/{id}
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_ID", "Invalid test ID")
		return
	}

	var req UpdateInterviewTestRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}

	// Get admin ID from context
	adminID := getAdminIDFromContext(r)

	result, err := h.service.Update(r.Context(), id, req, adminID)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	response.Success(w, http.StatusOK, "Interview test updated successfully", result)
}

// Delete handles DELETE /admin/interview-tests/{id}
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_ID", "Invalid test ID")
		return
	}

	err = h.service.Delete(r.Context(), id)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	response.Success(w, http.StatusOK, "Interview test deleted successfully", nil)
}

// Publish handles POST /admin/interview-tests/{id}/publish
func (h *Handler) Publish(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_ID", "Invalid test ID")
		return
	}

	// Get admin ID from context
	adminID := getAdminIDFromContext(r)

	result, err := h.service.Publish(r.Context(), id, adminID)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	response.Success(w, http.StatusOK, "Interview test published successfully", result)
}

// Duplicate handles POST /admin/interview-tests/{id}/duplicate
func (h *Handler) Duplicate(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_ID", "Invalid test ID")
		return
	}

	// Get admin ID from context
	adminID := getAdminIDFromContext(r)

	result, err := h.service.Duplicate(r.Context(), id, adminID)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	response.Success(w, http.StatusCreated, "Interview test duplicated successfully", result)
}

// =============== Helper functions ===============

// getAdminIDFromContext retrieves admin ID from request context
func getAdminIDFromContext(r *http.Request) uint64 {
	userID := r.Context().Value("user_id")
	if userID == nil {
		return 0
	}

	// Try to convert to uint64
	switch v := userID.(type) {
	case uint64:
		return v
	case int64:
		return uint64(v)
	case float64:
		return uint64(v)
	case string:
		id, _ := strconv.ParseUint(v, 10, 64)
		return id
	default:
		return 0
	}
}

// handleServiceError handles service layer errors
func handleServiceError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrTestNotFound):
		response.Error(w, http.StatusNotFound, "NOT_FOUND", err.Error())
	case errors.Is(err, ErrInvalidStatus):
		response.Error(w, http.StatusBadRequest, "INVALID_STATUS", err.Error())
	case errors.Is(err, ErrInvalidType):
		response.Error(w, http.StatusBadRequest, "INVALID_TYPE", err.Error())
	case errors.Is(err, ErrInvalidDifficulty):
		response.Error(w, http.StatusBadRequest, "INVALID_DIFFICULTY", err.Error())
	case errors.Is(err, ErrNoQuestions):
		response.Error(w, http.StatusBadRequest, "NO_QUESTIONS", err.Error())
	case errors.Is(err, ErrInvalidOptions):
		response.Error(w, http.StatusBadRequest, "INVALID_OPTIONS", err.Error())
	case errors.Is(err, context.DeadlineExceeded):
		response.Error(w, http.StatusRequestTimeout, "TIMEOUT", "Request timeout")
	case errors.Is(err, context.Canceled):
		response.Error(w, http.StatusRequestTimeout, "CANCELED", "Request canceled")
	default:
		response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
	}
}
