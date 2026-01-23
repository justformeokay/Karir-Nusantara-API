package company

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

// MiddlewareFunc defines the middleware function type
type MiddlewareFunc func(http.Handler) http.Handler

// RegisterRoutes registers company routes
func RegisterRoutes(r chi.Router, h *Handler, authenticate MiddlewareFunc) {
	r.Route("/companies", func(r chi.Router) {
		// Public routes - for job seekers to view company details
		r.Get("/{id}", h.GetPublicCompanyProfile)

		// Protected routes - for company owners
		r.Group(func(r chi.Router) {
			r.Use(authenticate)
			r.Get("/profile", h.GetCompanyProfile)
			r.Put("/profile", h.UpdateCompanyProfile)
			r.Post("/logo", h.UploadCompanyLogo)
			r.Post("/documents", h.UploadCompanyDocument)
		})
	})
}
