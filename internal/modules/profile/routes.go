package profile

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

// MiddlewareFunc defines the middleware function type
type MiddlewareFunc func(http.Handler) http.Handler

// RegisterRoutes registers profile routes
func RegisterRoutes(r chi.Router, h *Handler, authenticate, requireJobSeeker MiddlewareFunc) {
	r.Route("/profile", func(r chi.Router) {
		// All routes require authentication
		r.Use(authenticate)
		// Only job seekers can manage their profile
		r.Use(requireJobSeeker)

		// Profile endpoints
		r.Get("/", h.GetProfile)
		r.Put("/", h.UpdateProfile)
		r.Delete("/", h.DeleteProfile)

		// Document endpoints
		r.Route("/documents", func(r chi.Router) {
			r.Get("/", h.GetDocuments)
			r.Post("/", h.UploadDocument)

			r.Route("/{id}", func(r chi.Router) {
				r.Get("/", h.GetDocument)
				r.Put("/", h.UpdateDocument)
				r.Delete("/", h.DeleteDocument)
				r.Post("/primary", h.SetPrimaryDocument)
			})
		})
	})
}
