package jobreports

import (
	"github.com/go-chi/chi/v5"
	"github.com/karirnusantara/api/internal/middleware"
)

// RegisterRoutes registers job report routes
func RegisterRoutes(r chi.Router, h *Handler, authMiddleware *middleware.AuthMiddleware) {
	// Job Seeker routes - nested under /jobs/{jobId}/reports
	r.Route("/jobs/{jobId}/reports", func(r chi.Router) {
		r.Use(authMiddleware.Authenticate)
		r.Use(authMiddleware.RequireJobSeeker)

		r.Post("/", h.CreateReport)
		r.Get("/status", h.CheckReportStatus)
	})

	// Admin routes
	r.Route("/admin/job-reports", func(r chi.Router) {
		r.Use(authMiddleware.Authenticate)
		r.Use(authMiddleware.RequireAdmin)

		r.Get("/", h.AdminGetAllReports)
		r.Get("/pending-count", h.AdminGetPendingCount)
		r.Get("/{id}", h.AdminGetReport)
		r.Patch("/{id}/status", h.AdminUpdateReportStatus)
		r.Post("/ban-company", h.AdminBanCompany)
	})

	// Admin route for job-specific reports
	r.Route("/admin/jobs/{jobId}/reports", func(r chi.Router) {
		r.Use(authMiddleware.Authenticate)
		r.Use(authMiddleware.RequireAdmin)

		r.Get("/", h.AdminGetReportsByJob)
	})
}
