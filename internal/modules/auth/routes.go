package auth

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

// MiddlewareFunc defines the middleware function type
type MiddlewareFunc func(http.Handler) http.Handler

// RegisterRoutes registers auth routes
// authMiddleware is passed as a function type to avoid circular dependency
func RegisterRoutes(r chi.Router, h *Handler, authenticate MiddlewareFunc) {
	r.Route("/auth", func(r chi.Router) {
		// Public routes
		r.Post("/register", h.Register)
		r.Post("/login", h.Login)
		r.Post("/refresh", h.RefreshToken)
		r.Post("/forgot-password", h.ForgotPassword)
		r.Post("/reset-password", h.ResetPassword)

		// Protected routes
		r.Group(func(r chi.Router) {
			r.Use(authenticate)
			r.Post("/logout", h.Logout)
			r.Get("/me", h.Me)
			r.Put("/profile", h.UpdateProfile)
			r.Post("/profile/logo", h.UploadLogo)
			r.Put("/change-password", h.ChangePassword)
		})
	})
}
