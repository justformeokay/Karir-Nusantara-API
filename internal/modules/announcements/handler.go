package announcements

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/karirnusantara/api/internal/shared/response"
)

// Handler handles HTTP requests for announcements
type Handler struct {
	service *Service
}

// NewHandler creates a new announcements handler
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// Create handles POST /admin/announcements
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateAnnouncementRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}

	// Validate required fields
	if req.Title == "" {
		response.Error(w, http.StatusBadRequest, "VALIDATION_ERROR", "Title is required")
		return
	}
	if req.Content == "" {
		response.Error(w, http.StatusBadRequest, "VALIDATION_ERROR", "Content is required")
		return
	}
	if req.Type == "" {
		response.Error(w, http.StatusBadRequest, "VALIDATION_ERROR", "Type is required")
		return
	}
	if req.TargetAudience == "" {
		response.Error(w, http.StatusBadRequest, "VALIDATION_ERROR", "Target audience is required")
		return
	}

	// Get admin ID from context
	adminID := getAdminIDFromContext(r)

	result, err := h.service.Create(r.Context(), req, adminID)
	if err != nil {
		if errors.Is(err, ErrInvalidType) {
			response.Error(w, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid announcement type")
			return
		}
		if errors.Is(err, ErrInvalidAudience) {
			response.Error(w, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid target audience")
			return
		}
		response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	response.Success(w, http.StatusCreated, "Announcement created successfully", result)
}

// GetByID handles GET /admin/announcements/{id}
func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_ID", "Invalid announcement ID")
		return
	}

	result, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, ErrAnnouncementNotFound) {
			response.Error(w, http.StatusNotFound, "NOT_FOUND", "Announcement not found")
			return
		}
		response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	response.Success(w, http.StatusOK, "Announcement retrieved successfully", result)
}

// Update handles PUT /admin/announcements/{id}
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_ID", "Invalid announcement ID")
		return
	}

	var req UpdateAnnouncementRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}

	adminID := getAdminIDFromContext(r)

	result, err := h.service.Update(r.Context(), id, req, adminID)
	if err != nil {
		if errors.Is(err, ErrAnnouncementNotFound) {
			response.Error(w, http.StatusNotFound, "NOT_FOUND", "Announcement not found")
			return
		}
		if errors.Is(err, ErrInvalidType) {
			response.Error(w, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid announcement type")
			return
		}
		if errors.Is(err, ErrInvalidAudience) {
			response.Error(w, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid target audience")
			return
		}
		response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	response.Success(w, http.StatusOK, "Announcement updated successfully", result)
}

// Delete handles DELETE /admin/announcements/{id}
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_ID", "Invalid announcement ID")
		return
	}

	err = h.service.Delete(r.Context(), id)
	if err != nil {
		if err.Error() == "announcement not found" {
			response.Error(w, http.StatusNotFound, "NOT_FOUND", "Announcement not found")
			return
		}
		response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	response.Success(w, http.StatusOK, "Announcement deleted successfully", nil)
}

// ToggleStatus handles PATCH /admin/announcements/{id}/toggle
func (h *Handler) ToggleStatus(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_ID", "Invalid announcement ID")
		return
	}

	var req ToggleStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}

	adminID := getAdminIDFromContext(r)

	result, err := h.service.ToggleStatus(r.Context(), id, req.IsActive, adminID)
	if err != nil {
		if err.Error() == "announcement not found" {
			response.Error(w, http.StatusNotFound, "NOT_FOUND", "Announcement not found")
			return
		}
		response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	response.Success(w, http.StatusOK, "Announcement status updated successfully", result)
}

// List handles GET /admin/announcements
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	filter := AnnouncementFilter{
		Type:           r.URL.Query().Get("type"),
		TargetAudience: r.URL.Query().Get("target_audience"),
		Search:         r.URL.Query().Get("search"),
	}

	// Parse is_active
	if isActiveStr := r.URL.Query().Get("is_active"); isActiveStr != "" {
		isActive := isActiveStr == "true"
		filter.IsActive = &isActive
	}

	// Parse pagination
	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		if page, err := strconv.Atoi(pageStr); err == nil {
			filter.Page = page
		}
	}
	if filter.Page <= 0 {
		filter.Page = 1
	}

	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil {
			filter.Limit = limit
		}
	}
	if filter.Limit <= 0 {
		filter.Limit = 10
	}

	result, err := h.service.List(r.Context(), filter)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	response.Success(w, http.StatusOK, "Announcements retrieved successfully", result)
}

// GetPublicNotifications handles GET /announcements/notifications
func (h *Handler) GetPublicNotifications(w http.ResponseWriter, r *http.Request) {
	audience := r.URL.Query().Get("audience")
	if audience == "" {
		audience = "all"
	}

	result, err := h.service.GetActiveByType(r.Context(), "notification", audience)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	response.Success(w, http.StatusOK, "Notifications retrieved successfully", result)
}

// GetPublicBanners handles GET /announcements/banners
func (h *Handler) GetPublicBanners(w http.ResponseWriter, r *http.Request) {
	audience := r.URL.Query().Get("audience")
	if audience == "" {
		audience = "all"
	}

	result, err := h.service.GetActiveByType(r.Context(), "banner", audience)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	response.Success(w, http.StatusOK, "Banners retrieved successfully", result)
}

// GetPublicInformation handles GET /announcements/information
func (h *Handler) GetPublicInformation(w http.ResponseWriter, r *http.Request) {
	audience := r.URL.Query().Get("audience")
	if audience == "" {
		audience = "all"
	}

	result, err := h.service.GetActiveByType(r.Context(), "information", audience)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	response.Success(w, http.StatusOK, "Information retrieved successfully", result)
}

// GetAllPublic handles GET /announcements
func (h *Handler) GetAllPublic(w http.ResponseWriter, r *http.Request) {
	audience := r.URL.Query().Get("audience")
	if audience == "" {
		audience = "all"
	}

	result, err := h.service.GetAllActive(r.Context(), audience)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	response.Success(w, http.StatusOK, "Announcements retrieved successfully", result)
}

// Helper function to get admin ID from context
func getAdminIDFromContext(r *http.Request) uint64 {
	// Try to get from context (set by auth middleware)
	if userID, ok := r.Context().Value("user_id").(uint64); ok {
		return userID
	}
	// Fallback: try string format
	if userIDStr, ok := r.Context().Value("user_id").(string); ok {
		if id, err := strconv.ParseUint(userIDStr, 10, 64); err == nil {
			return id
		}
	}
	return 0
}
