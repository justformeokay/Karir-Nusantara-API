package cvs

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

// MiddlewareFunc defines the middleware function type
type MiddlewareFunc func(http.Handler) http.Handler

// RegisterRoutes registers CV routes
func RegisterRoutes(r chi.Router, h *Handler, authenticate, requireJobSeeker MiddlewareFunc) {
	r.Route("/cv", func(r chi.Router) {
		// All routes require authentication
		r.Use(authenticate)
		// Only job seekers can manage CVs
		r.Use(requireJobSeeker)

		r.Get("/", h.Get)
		r.Post("/", h.CreateOrUpdate)
		r.Put("/", h.CreateOrUpdate)
		r.Delete("/", h.Delete)
	})
}
