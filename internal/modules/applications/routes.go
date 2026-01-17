package applications

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

// MiddlewareFunc defines the middleware function type
type MiddlewareFunc func(http.Handler) http.Handler

// RegisterRoutes registers the applications routes
func RegisterRoutes(r chi.Router, h *Handler, authenticate, requireJobSeeker, requireCompany MiddlewareFunc) {
	r.Route("/applications", func(r chi.Router) {
		// Protected routes - require authentication
		r.Group(func(r chi.Router) {
			r.Use(authenticate)

			// Get single application (both job seeker and company)
			r.Get("/{id}", h.GetByID)

			// Get application timeline
			r.Get("/{id}/timeline", h.GetTimeline)
		})

		// Job seeker only routes
		r.Group(func(r chi.Router) {
			r.Use(authenticate)
			r.Use(requireJobSeeker)

			// Apply for a job
			r.Post("/", h.Apply)

			// List my applications
			r.Get("/me", h.ListMyApplications)

			// Withdraw application
			r.Post("/{id}/withdraw", h.Withdraw)
		})

		// Company only routes
		r.Group(func(r chi.Router) {
			r.Use(authenticate)
			r.Use(requireCompany)

			// List all applications for company
			r.Get("/company", h.ListCompanyApplications)

			// Update application status
			r.Patch("/{id}/status", h.UpdateStatus)
		})
	})

	// Company route for listing applications by job
	r.Route("/jobs/{jobId}/applications", func(r chi.Router) {
		r.Use(authenticate)
		r.Use(requireCompany)

		r.Get("/", h.ListJobApplications)
	})
}
