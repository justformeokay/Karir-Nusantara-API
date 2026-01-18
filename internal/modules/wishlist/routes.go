package wishlist

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

// MiddlewareFunc defines the middleware function type
type MiddlewareFunc func(http.Handler) http.Handler

// RegisterRoutes registers the wishlist routes
func RegisterRoutes(r chi.Router, h *Handler, authenticate, requireJobSeeker MiddlewareFunc) {
	r.Route("/wishlist", func(r chi.Router) {
		// All routes require authentication as job seeker
		r.Group(func(r chi.Router) {
			r.Use(authenticate)
			r.Use(requireJobSeeker)

			// List saved jobs
			r.Get("/", h.ListSavedJobs)

			// Save a job
			r.Post("/", h.SaveJob)

			// Get wishlist stats
			r.Get("/stats", h.GetStats)

			// Check if job is saved
			r.Get("/check/{jobId}", h.CheckSaved)

			// Remove a saved job
			r.Delete("/{jobId}", h.RemoveJob)
		})
	})
}
