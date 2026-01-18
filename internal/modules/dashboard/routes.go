package dashboard

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

// MiddlewareFunc defines the middleware function type
type MiddlewareFunc func(http.Handler) http.Handler

// RegisterRoutes registers all dashboard routes
func RegisterRoutes(r chi.Router, h *Handler, authenticate, requireCompany MiddlewareFunc) {
	r.Route("/company/dashboard", func(r chi.Router) {
		r.Use(authenticate)
		r.Use(requireCompany)

		r.Get("/stats", h.GetStats)
		r.Get("/recent-applicants", h.GetRecentApplicants)
		r.Get("/active-jobs", h.GetActiveJobs)
	})
}
