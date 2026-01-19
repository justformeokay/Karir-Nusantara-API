package admin

import (
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"

	"github.com/karirnusantara/api/internal/config"
	"github.com/karirnusantara/api/internal/middleware"
)

// Module represents the admin module
type Module struct {
	handler    *Handler
	authMiddleware *middleware.AuthMiddleware
}

// NewModule creates a new admin module
func NewModule(db *sqlx.DB, cfg *config.Config, authMiddleware *middleware.AuthMiddleware) *Module {
	repo := NewRepository(db)
	service := NewService(repo, cfg)
	handler := NewHandler(service)

	return &Module{
		handler:    handler,
		authMiddleware: authMiddleware,
	}
}

// RegisterRoutes registers admin routes
func (m *Module) RegisterRoutes(r chi.Router) {
	r.Route("/admin", func(r chi.Router) {
		// Public routes (login)
		r.Route("/auth", func(r chi.Router) {
			r.Post("/login", m.handler.Login)
		})

		// Protected admin routes
		r.Group(func(r chi.Router) {
			// Require authentication and admin role
			r.Use(m.authMiddleware.Authenticate)
			r.Use(m.authMiddleware.RequireAdmin)

			// Current admin info
			r.Get("/auth/me", m.handler.GetCurrentAdmin)

			// Dashboard
			r.Get("/dashboard/stats", m.handler.GetDashboardStats)

			// Company management
			r.Route("/companies", func(r chi.Router) {
				r.Get("/", m.handler.GetCompanies)
				r.Get("/{id}", m.handler.GetCompanyByID)
				r.Post("/{id}/verify", m.handler.VerifyCompany)
				r.Patch("/{id}/status", m.handler.UpdateCompanyStatus)
			})

			// Job management
			r.Route("/jobs", func(r chi.Router) {
				r.Get("/", m.handler.GetJobs)
				r.Get("/{id}", m.handler.GetJobByID)
				r.Post("/{id}/moderate", m.handler.ModerateJob)
			})

			// Payment management
			r.Route("/payments", func(r chi.Router) {
				r.Get("/", m.handler.GetPayments)
				r.Get("/{id}", m.handler.GetPaymentByID)
				r.Post("/{id}/process", m.handler.ProcessPayment)
			})

			// Job seeker management
			r.Route("/job-seekers", func(r chi.Router) {
				r.Get("/", m.handler.GetJobSeekers)
				r.Get("/{id}", m.handler.GetJobSeekerByID)
				r.Patch("/{id}/status", m.handler.UpdateJobSeekerStatus)
			})
		})
	})
}
