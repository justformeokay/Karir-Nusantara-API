package jobs

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

// MiddlewareFunc defines the middleware function type
type MiddlewareFunc func(http.Handler) http.Handler

// RegisterRoutes registers job routes
func RegisterRoutes(r chi.Router, h *Handler, authenticate, requireCompany MiddlewareFunc) {
	r.Route("/jobs", func(r chi.Router) {
		// Public routes
		r.Get("/", h.List)
		r.Get("/{id}", h.GetByID)
		r.Get("/slug/{slug}", h.GetBySlug)

		// Company-only routes
		r.Group(func(r chi.Router) {
			r.Use(authenticate)
			r.Use(requireCompany)
			r.Post("/", h.Create)
			r.Put("/{id}", h.Update)
			r.Delete("/{id}", h.Delete)

			// Job status management endpoints
			r.Patch("/{id}/publish", h.Publish)
			r.Patch("/{id}/close", h.Close)
			r.Patch("/{id}/pause", h.Pause)
			r.Patch("/{id}/reopen", h.Reopen)
		})

		// Company-specific list route
		r.Group(func(r chi.Router) {
			r.Use(authenticate)
			r.Use(requireCompany)
			r.Get("/company/list", h.ListByCompany)
		})
	})
}
