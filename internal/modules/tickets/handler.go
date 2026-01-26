package tickets

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/karirnusantara/api/internal/middleware"
	"github.com/karirnusantara/api/internal/shared/response"
	"github.com/karirnusantara/api/internal/shared/validator"
)

// Handler handles ticket HTTP requests
type Handler struct {
	service   Service
	validator *validator.Validator
}

// NewHandler creates a new ticket handler
func NewHandler(service Service, validator *validator.Validator) *Handler {
	return &Handler{
		service:   service,
		validator: validator,
	}
}

// CreateTicket handles creating a new support ticket
// @Summary Create a new support ticket
// @Tags Tickets
// @Accept json
// @Produce json
// @Param request body CreateTicketRequest true "Ticket details"
// @Success 201 {object} TicketWithDetails
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 429 {object} response.ErrorResponse "Cooldown active"
// @Router /tickets [post]
func (h *Handler) CreateTicket(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == 0 {
		response.Error(w, http.StatusUnauthorized, "unauthorized", "Anda harus login terlebih dahulu")
		return
	}

	var req CreateTicketRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid_request", "Format request tidak valid")
		return
	}

	if errors := h.validator.Validate(&req); errors != nil {
		response.UnprocessableEntity(w, "Validation failed", errors)
		return
	}

	ticket, err := h.service.CreateTicket(r.Context(), userID, &req)
	if err != nil {
		// Check if it's a cooldown error
		if err.Error()[:4] == "Anda" {
			response.Error(w, http.StatusTooManyRequests, "cooldown_active", err.Error())
			return
		}
		response.Error(w, http.StatusInternalServerError, "create_failed", err.Error())
		return
	}

	response.JSON(w, http.StatusCreated, map[string]interface{}{
		"message": "Ticket berhasil dibuat! Tim support akan merespons dalam waktu 24 jam.",
		"ticket":  ticket,
	})
}

// GetMyTickets handles listing user's tickets
// @Summary Get user's tickets
// @Tags Tickets
// @Produce json
// @Success 200 {array} TicketWithDetails
// @Router /tickets [get]
func (h *Handler) GetMyTickets(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == 0 {
		response.Error(w, http.StatusUnauthorized, "unauthorized", "Anda harus login terlebih dahulu")
		return
	}

	tickets, err := h.service.GetMyTickets(r.Context(), userID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "fetch_failed", err.Error())
		return
	}

	response.JSON(w, http.StatusOK, map[string]interface{}{
		"tickets": tickets,
		"total":   len(tickets),
	})
}

// GetTicket handles getting a specific ticket with responses
// @Summary Get ticket details with responses
// @Tags Tickets
// @Produce json
// @Param id path int true "Ticket ID"
// @Success 200 {object} map[string]interface{}
// @Router /tickets/{id} [get]
func (h *Handler) GetTicket(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == 0 {
		response.Error(w, http.StatusUnauthorized, "unauthorized", "Anda harus login terlebih dahulu")
		return
	}

	ticketID, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid_id", "ID ticket tidak valid")
		return
	}

	ticket, responses, err := h.service.GetTicket(r.Context(), ticketID, userID)
	if err != nil {
		if err.Error() == "ticket not found" {
			response.Error(w, http.StatusNotFound, "not_found", "Ticket tidak ditemukan")
			return
		}
		if err.Error() == "unauthorized to view this ticket" {
			response.Error(w, http.StatusForbidden, "forbidden", "Anda tidak memiliki akses ke ticket ini")
			return
		}
		response.Error(w, http.StatusInternalServerError, "fetch_failed", err.Error())
		return
	}

	response.JSON(w, http.StatusOK, map[string]interface{}{
		"ticket":    ticket,
		"responses": responses,
	})
}

// AddResponse handles adding a response to a ticket
// @Summary Add response to ticket
// @Tags Tickets
// @Accept json
// @Produce json
// @Param id path int true "Ticket ID"
// @Param request body AddResponseRequest true "Response message"
// @Success 201 {object} TicketResponseWithSender
// @Router /tickets/{id}/responses [post]
func (h *Handler) AddResponse(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == 0 {
		response.Error(w, http.StatusUnauthorized, "unauthorized", "Anda harus login terlebih dahulu")
		return
	}

	ticketID, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid_id", "ID ticket tidak valid")
		return
	}

	var req AddResponseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid_request", "Format request tidak valid")
		return
	}

	if errors := h.validator.Validate(&req); errors != nil {
		response.UnprocessableEntity(w, "Validation failed", errors)
		return
	}

	resp, err := h.service.AddResponse(r.Context(), ticketID, userID, "user", &req)
	if err != nil {
		if err.Error() == "ticket not found" {
			response.Error(w, http.StatusNotFound, "not_found", "Ticket tidak ditemukan")
			return
		}
		if err.Error() == "unauthorized to respond to this ticket" {
			response.Error(w, http.StatusForbidden, "forbidden", "Anda tidak memiliki akses ke ticket ini")
			return
		}
		response.Error(w, http.StatusInternalServerError, "create_failed", err.Error())
		return
	}

	response.JSON(w, http.StatusCreated, map[string]interface{}{
		"message":  "Respons berhasil dikirim!",
		"response": resp,
	})
}

// CheckCooldown handles checking user's cooldown status
// @Summary Check if user can create a new ticket
// @Tags Tickets
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /tickets/cooldown [get]
func (h *Handler) CheckCooldown(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == 0 {
		response.Error(w, http.StatusUnauthorized, "unauthorized", "Anda harus login terlebih dahulu")
		return
	}

	canCreate, remaining, err := h.service.CheckCooldown(r.Context(), userID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "check_failed", err.Error())
		return
	}

	response.JSON(w, http.StatusOK, map[string]interface{}{
		"can_create":        canCreate,
		"remaining_seconds": int(remaining.Seconds()),
		"remaining_minutes": int(remaining.Minutes()),
		"cooldown_duration": int(CooldownDuration.Minutes()),
	})
}

// ============== ADMIN ENDPOINTS ==============

// AdminGetAllTickets handles listing all tickets for admin
// @Summary Get all tickets (Admin only)
// @Tags Tickets Admin
// @Produce json
// @Param status query string false "Filter by status"
// @Param priority query string false "Filter by priority"
// @Success 200 {array} TicketWithDetails
// @Router /admin/tickets [get]
func (h *Handler) AdminGetAllTickets(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")
	priority := r.URL.Query().Get("priority")

	tickets, err := h.service.GetAllTickets(r.Context(), status, priority)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "fetch_failed", err.Error())
		return
	}

	// Calculate stats
	stats := map[string]int{
		"total":            len(tickets),
		"open":             0,
		"in_progress":      0,
		"pending_response": 0,
		"resolved":         0,
		"closed":           0,
	}

	for _, t := range tickets {
		stats[t.Status]++
	}

	response.JSON(w, http.StatusOK, map[string]interface{}{
		"tickets": tickets,
		"stats":   stats,
	})
}

// AdminGetTicket handles getting a specific ticket for admin
// @Summary Get ticket details (Admin only)
// @Tags Tickets Admin
// @Produce json
// @Param id path int true "Ticket ID"
// @Success 200 {object} map[string]interface{}
// @Router /admin/tickets/{id} [get]
func (h *Handler) AdminGetTicket(w http.ResponseWriter, r *http.Request) {
	ticketID, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid_id", "ID ticket tidak valid")
		return
	}

	// For admin, we use a different method that doesn't check ownership
	ticket, err := h.service.(*service).repo.GetTicketByID(r.Context(), ticketID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "fetch_failed", err.Error())
		return
	}

	if ticket == nil {
		response.Error(w, http.StatusNotFound, "not_found", "Ticket tidak ditemukan")
		return
	}

	responses, err := h.service.(*service).repo.ListResponsesByTicket(r.Context(), ticketID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "fetch_failed", err.Error())
		return
	}

	response.JSON(w, http.StatusOK, map[string]interface{}{
		"ticket":    ticket,
		"responses": responses,
	})
}

// AdminAddResponse handles adding admin response to a ticket
// @Summary Add admin response to ticket
// @Tags Tickets Admin
// @Accept json
// @Produce json
// @Param id path int true "Ticket ID"
// @Param request body AddResponseRequest true "Response message"
// @Success 201 {object} TicketResponseWithSender
// @Router /admin/tickets/{id}/responses [post]
func (h *Handler) AdminAddResponse(w http.ResponseWriter, r *http.Request) {
	adminID := middleware.GetUserID(r.Context())
	if adminID == 0 {
		response.Error(w, http.StatusUnauthorized, "unauthorized", "Anda harus login terlebih dahulu")
		return
	}

	ticketID, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid_id", "ID ticket tidak valid")
		return
	}

	var req AddResponseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid_request", "Format request tidak valid")
		return
	}

	if errors := h.validator.Validate(&req); errors != nil {
		response.UnprocessableEntity(w, "Validation failed", errors)
		return
	}

	// Admin response - skip ownership check by directly using service with "admin" sender type
	ticket, err := h.service.(*service).repo.GetTicketByID(r.Context(), ticketID)
	if err != nil || ticket == nil {
		response.Error(w, http.StatusNotFound, "not_found", "Ticket tidak ditemukan")
		return
	}

	if ticket.Status == TicketStatusClosed {
		response.Error(w, http.StatusBadRequest, "ticket_closed", "Ticket sudah ditutup")
		return
	}

	respData := &TicketResponse{
		TicketID:   ticketID,
		SenderID:   adminID,
		SenderType: "admin",
		Message:    req.Message,
		CreatedAt:  time.Now(),
	}

	err = h.service.(*service).repo.CreateResponse(r.Context(), respData)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "create_failed", err.Error())
		return
	}

	response.JSON(w, http.StatusCreated, map[string]interface{}{
		"message": "Respons admin berhasil dikirim!",
		"response": &TicketResponseWithSender{
			TicketResponse: *respData,
			SenderName:     "Admin Support",
			SenderEmail:    "support@karirnusantara.id",
		},
	})
}

// AdminUpdateTicketStatus handles updating ticket status
// @Summary Update ticket status (Admin only)
// @Tags Tickets Admin
// @Accept json
// @Produce json
// @Param id path int true "Ticket ID"
// @Param request body UpdateTicketStatusRequest true "New status"
// @Success 200 {object} map[string]interface{}
// @Router /admin/tickets/{id}/status [patch]
func (h *Handler) AdminUpdateTicketStatus(w http.ResponseWriter, r *http.Request) {
	ticketID, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid_id", "ID ticket tidak valid")
		return
	}

	var req UpdateTicketStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid_request", "Format request tidak valid")
		return
	}

	if errors := h.validator.Validate(&req); errors != nil {
		response.UnprocessableEntity(w, "Validation failed", errors)
		return
	}

	err = h.service.UpdateTicketStatus(r.Context(), ticketID, &req)
	if err != nil {
		if err.Error() == "ticket not found" {
			response.Error(w, http.StatusNotFound, "not_found", "Ticket tidak ditemukan")
			return
		}
		response.Error(w, http.StatusInternalServerError, "update_failed", err.Error())
		return
	}

	response.JSON(w, http.StatusOK, map[string]interface{}{
		"message": "Status ticket berhasil diperbarui ke " + GetStatusLabel(req.Status),
	})
}

// AdminCloseTicket handles closing a ticket
// @Summary Close ticket (Admin only)
// @Tags Tickets Admin
// @Produce json
// @Param id path int true "Ticket ID"
// @Success 200 {object} map[string]interface{}
// @Router /admin/tickets/{id}/close [post]
func (h *Handler) AdminCloseTicket(w http.ResponseWriter, r *http.Request) {
	ticketID, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid_id", "ID ticket tidak valid")
		return
	}

	err = h.service.CloseTicket(r.Context(), ticketID)
	if err != nil {
		if err.Error() == "ticket not found" {
			response.Error(w, http.StatusNotFound, "not_found", "Ticket tidak ditemukan")
			return
		}
		response.Error(w, http.StatusInternalServerError, "close_failed", err.Error())
		return
	}

	response.JSON(w, http.StatusOK, map[string]interface{}{
		"message": "Ticket berhasil ditutup",
	})
}

// AdminResolveTicket handles resolving a ticket
// @Summary Resolve ticket (Admin only)
// @Tags Tickets Admin
// @Produce json
// @Param id path int true "Ticket ID"
// @Success 200 {object} map[string]interface{}
// @Router /admin/tickets/{id}/resolve [post]
func (h *Handler) AdminResolveTicket(w http.ResponseWriter, r *http.Request) {
	ticketID, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid_id", "ID ticket tidak valid")
		return
	}

	err = h.service.ResolveTicket(r.Context(), ticketID)
	if err != nil {
		if err.Error() == "ticket not found" {
			response.Error(w, http.StatusNotFound, "not_found", "Ticket tidak ditemukan")
			return
		}
		response.Error(w, http.StatusInternalServerError, "resolve_failed", err.Error())
		return
	}

	response.JSON(w, http.StatusOK, map[string]interface{}{
		"message": "Ticket berhasil ditandai sebagai terselesaikan",
	})
}
