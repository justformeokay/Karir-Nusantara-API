package cvs

import (
	"encoding/json"
	"net/http"

	"github.com/karirnusantara/api/internal/middleware"
	apperrors "github.com/karirnusantara/api/internal/shared/errors"
	"github.com/karirnusantara/api/internal/shared/response"
	"github.com/karirnusantara/api/internal/shared/validator"
)

// Handler handles CV HTTP requests
type Handler struct {
	service   Service
	validator *validator.Validator
}

// NewHandler creates a new CV handler
func NewHandler(service Service, validator *validator.Validator) *Handler {
	return &Handler{
		service:   service,
		validator: validator,
	}
}

// CreateOrUpdate handles CV creation/update
// POST /api/v1/cv
func (h *Handler) CreateOrUpdate(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == 0 {
		response.Unauthorized(w, "Unauthorized")
		return
	}

	var req CreateCVRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid request body")
		return
	}

	if errors := h.validator.Validate(&req); errors != nil {
		response.UnprocessableEntity(w, "Validation failed", errors)
		return
	}

	cv, err := h.service.CreateOrUpdate(r.Context(), userID, &req)
	if err != nil {
		handleError(w, err)
		return
	}

	response.OK(w, "CV saved successfully", cv)
}

// Get handles getting the current user's CV
// GET /api/v1/cv
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == 0 {
		response.Unauthorized(w, "Unauthorized")
		return
	}

	cv, err := h.service.GetByUserID(r.Context(), userID)
	if err != nil {
		handleError(w, err)
		return
	}

	response.OK(w, "CV retrieved", cv)
}

// Delete handles CV deletion
// DELETE /api/v1/cv
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == 0 {
		response.Unauthorized(w, "Unauthorized")
		return
	}

	if err := h.service.Delete(r.Context(), userID); err != nil {
		handleError(w, err)
		return
	}

	response.NoContent(w)
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
