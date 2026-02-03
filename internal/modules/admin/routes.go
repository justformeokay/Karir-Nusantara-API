package admin

import (
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"

	"github.com/karirnusantara/api/internal/config"
	"github.com/karirnusantara/api/internal/middleware"
	"github.com/karirnusantara/api/internal/modules/announcements"
	"github.com/karirnusantara/api/internal/modules/quota"
	"github.com/karirnusantara/api/internal/shared/email"
	"github.com/karirnusantara/api/internal/shared/invoice"
)

// Module represents the admin module
type Module struct {
	handler             *Handler
	authMiddleware      *middleware.AuthMiddleware
	announcementsModule *announcements.Module
}

// NewModule creates a new admin module
func NewModule(db *sqlx.DB, cfg *config.Config, authMiddleware *middleware.AuthMiddleware) *Module {
	repo := NewRepository(db)
	service := NewService(repo, cfg)
	handler := NewHandler(service)

	return &Module{
		handler:        handler,
		authMiddleware: authMiddleware,
	}
}

// NewModuleWithQuota creates a new admin module with quota service
func NewModuleWithQuota(db *sqlx.DB, cfg *config.Config, authMiddleware *middleware.AuthMiddleware, quotaSvc *quota.Service, emailSvc *email.Service, invoiceSvc *invoice.Service) *Module {
	repo := NewRepository(db)
	service := NewServiceComplete(repo, cfg, quotaSvc, emailSvc, invoiceSvc)
	handler := NewHandler(service)

	// Initialize announcements module
	announcementsModule := announcements.NewModule(db, authMiddleware)

	return &Module{
		handler:             handler,
		authMiddleware:      authMiddleware,
		announcementsModule: announcementsModule,
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
			r.Get("/dashboard/pending-companies", m.handler.GetPendingCompanies)
			r.Get("/dashboard/pending-payments", m.handler.GetPendingPayments)
			r.Get("/dashboard/open-tickets", m.handler.GetOpenSupportTickets)

			// Company management
			r.Route("/companies", func(r chi.Router) {
				r.Get("/", m.handler.GetCompanies)
				r.Get("/{id}", m.handler.GetCompanyByID)
				r.Get("/{id}/detail", m.handler.GetCompanyDetail)
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

			// Announcements management (notifications, banners, information)
			if m.announcementsModule != nil {
				m.announcementsModule.RegisterAdminRoutes(r)
			}
		})
	})
}

// GetAnnouncementsModule returns the announcements module for public routes registration
func (m *Module) GetAnnouncementsModule() *announcements.Module {
	return m.announcementsModule
}
