package tickets

import (
	"github.com/go-chi/chi/v5"
	"github.com/karirnusantara/api/internal/middleware"
)

// RegisterRoutes registers ticket routes
func RegisterRoutes(r chi.Router, h *Handler, authMiddleware *middleware.AuthMiddleware) {
	// Job Seeker routes
	r.Route("/tickets", func(r chi.Router) {
		r.Use(authMiddleware.Authenticate)
		r.Use(authMiddleware.RequireJobSeeker)

		r.Post("/", h.CreateTicket)
		r.Get("/", h.GetMyTickets)
		r.Get("/cooldown", h.CheckCooldown)
		r.Get("/{id}", h.GetTicket)
		r.Post("/{id}/responses", h.AddResponse)
	})

	// Admin routes
	r.Route("/admin/tickets", func(r chi.Router) {
		r.Use(authMiddleware.Authenticate)
		r.Use(authMiddleware.RequireAdmin)

		r.Get("/", h.AdminGetAllTickets)
		r.Get("/{id}", h.AdminGetTicket)
		r.Post("/{id}/responses", h.AdminAddResponse)
		r.Patch("/{id}/status", h.AdminUpdateTicketStatus)
		r.Post("/{id}/close", h.AdminCloseTicket)
		r.Post("/{id}/resolve", h.AdminResolveTicket)
	})
}
