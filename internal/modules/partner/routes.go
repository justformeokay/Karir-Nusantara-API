package partner

import (
	"context"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/karirnusantara/api/internal/shared/response"
)

// MiddlewareFunc defines the middleware function type
type MiddlewareFunc func(http.Handler) http.Handler

// PartnerMiddleware handles partner-specific authentication
type PartnerMiddleware struct {
	service Service
}

// NewPartnerMiddleware creates a new partner middleware
func NewPartnerMiddleware(service Service) *PartnerMiddleware {
	return &PartnerMiddleware{service: service}
}

// Authenticate validates the JWT token and sets partner info in context
func (m *PartnerMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract token from Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			response.Unauthorized(w, "Missing authorization header")
			return
		}

		// Check Bearer prefix
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			response.Unauthorized(w, "Invalid authorization header format")
			return
		}

		tokenString := parts[1]

		// Validate token
		claims, err := m.service.ValidateAccessToken(tokenString)
		if err != nil {
			response.Unauthorized(w, "Invalid or expired token")
			return
		}

		// Verify this is a partner token
		if claims.Role != "partner" {
			response.Forbidden(w, "Access denied. Partner account required.")
			return
		}

		// Set partner info in context
		ctx := context.WithValue(r.Context(), "user_id", claims.UserID)
		ctx = context.WithValue(ctx, "user_email", claims.Email)
		ctx = context.WithValue(ctx, "user_role", claims.Role)
		ctx = context.WithValue(ctx, "partner_id", claims.PartnerID)
		ctx = context.WithValue(ctx, "referral_code", claims.ReferralCode)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RegisterRoutes registers partner routes
func RegisterRoutes(r chi.Router, h *Handler, middleware *PartnerMiddleware) {
	r.Route("/partner", func(r chi.Router) {
		// Public routes (no authentication required)
		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", h.Register)
			r.Post("/login", h.Login)
			r.Post("/forgot-password", h.ForgotPassword)
			r.Post("/reset-password", h.ResetPassword)
		})

		// Protected routes (authentication required)
		r.Group(func(r chi.Router) {
			r.Use(middleware.Authenticate)

			// Auth routes
			r.Post("/auth/logout", h.Logout)

			// Profile routes
			r.Get("/profile", h.GetProfile)
			r.Patch("/profile", h.UpdateProfile)
			r.Post("/password/change", h.ChangePassword)

			// Dashboard routes
			r.Get("/dashboard/stats", h.GetDashboardStats)
			r.Get("/dashboard/monthly", h.GetMonthlyData)

			// Referral routes
			r.Get("/referral", h.GetReferralInfo)

			// Companies routes
			r.Get("/companies", h.GetCompanies)
			r.Get("/companies/summary", h.GetCompaniesSummary)

			// Transactions routes
			r.Get("/transactions", h.GetTransactions)
			r.Get("/transactions/summary", h.GetTransactionsSummary)

			// Payouts routes
			r.Get("/payouts", h.GetPayoutInfo)
			r.Get("/payouts/history", h.GetPayoutHistory)
		})
	})
}
