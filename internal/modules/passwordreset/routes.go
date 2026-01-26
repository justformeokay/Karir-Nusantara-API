package passwordreset

import (
	"github.com/go-chi/chi/v5"
)

// RegisterRoutes registers password reset routes
func RegisterRoutes(r chi.Router, handler *Handler) {
	r.Route("/password-reset", func(r chi.Router) {
		r.Post("/forgot", handler.ForgotPassword)
		r.Get("/verify", handler.VerifyToken)
		r.Post("/reset", handler.ResetPassword)
	})
}

// SetupRoutes is an alias for RegisterRoutes for consistency
func SetupRoutes(r chi.Router, handler *Handler) {
	RegisterRoutes(r, handler)
}
