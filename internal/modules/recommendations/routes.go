package recommendations

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

// MiddlewareFunc defines the middleware function type
type MiddlewareFunc func(http.Handler) http.Handler

// RegisterRoutes registers recommendation routes
func RegisterRoutes(r chi.Router, h *Handler, authenticate MiddlewareFunc) {
	r.Route("/recommendations", func(r chi.Router) {
		// Authenticated routes only
		r.Use(authenticate)
		r.Get("/", h.GetRecommendations)
	})
}
