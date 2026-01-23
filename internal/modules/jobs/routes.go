package jobs

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

// MiddlewareFunc defines the middleware function type
type MiddlewareFunc func(http.Handler) http.Handler

// RegisterRoutes registers job routes
func RegisterRoutes(r chi.Router, h *Handler, authenticate, requireCompany, requireJobSeeker MiddlewareFunc) {
	r.Route("/jobs", func(r chi.Router) {
		// Public routes (must list specific paths first to avoid conflicts)
		r.Get("/", h.List)
		r.Get("/slug/{slug}", h.GetBySlug)

		// Company-specific list route (specific path before {id})
		r.Group(func(r chi.Router) {
			r.Use(authenticate)
			r.Use(requireCompany)
			r.Get("/company/list", h.ListByCompany)
		})

		// Routes with {id} parameter
		r.Route("/{id}", func(r chi.Router) {
			// Public route
			r.Get("/", h.GetByID)

			// Job seeker tracking routes
			r.Group(func(r chi.Router) {
				r.Use(authenticate)
				r.Use(requireJobSeeker)
				r.Post("/view", h.TrackView)
			})

			// Share tracking - authenticated users
			r.Group(func(r chi.Router) {
				r.Use(authenticate)
				r.Post("/share", h.TrackShare)
			})

			// Company-only routes
			r.Group(func(r chi.Router) {
				r.Use(authenticate)
				r.Use(requireCompany)
				r.Put("/", h.Update)
				r.Delete("/", h.Delete)
				r.Patch("/publish", h.Publish)
				r.Patch("/close", h.Close)
				r.Patch("/pause", h.Pause)
				r.Patch("/reopen", h.Reopen)
				r.Get("/stats", h.GetJobStats)
			})
		})

		// Company create route
		r.Group(func(r chi.Router) {
			r.Use(authenticate)
			r.Use(requireCompany)
			r.Post("/", h.Create)
		})
	})
}
