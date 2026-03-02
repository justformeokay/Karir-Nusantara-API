package systemsettings

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/karirnusantara/api/internal/middleware"
	"github.com/karirnusantara/api/internal/shared/validator"
)

// MiddlewareFunc is a standard HTTP middleware function type
type MiddlewareFunc func(http.Handler) http.Handler

// RegisterAdminRoutes registers admin-protected system settings routes
// All routes under /admin/system-settings require authentication + admin role
func RegisterAdminRoutes(r chi.Router, h *Handler, authenticate, requireAdmin MiddlewareFunc) {
	r.Route("/admin/system-settings", func(r chi.Router) {
		r.Use(authenticate)
		r.Use(requireAdmin)

		// GET    /api/v1/admin/system-settings         - Get all settings
		r.Get("/", h.GetSettings)

		// PUT    /api/v1/admin/system-settings         - Update general settings
		r.Put("/", h.UpdateSettings)

		// POST   /api/v1/admin/system-settings/packages       - Create / update a package
		r.Post("/packages", h.UpsertPackage)

		// DELETE /api/v1/admin/system-settings/packages/{id} - Delete a package
		r.Delete("/packages/{id}", h.DeletePackage)
	})
}

// RegisterPublicRoutes registers the public pricing endpoint
// GET /api/v1/pricing  — no authentication required
func RegisterPublicRoutes(r chi.Router, h *Handler) {
	r.Get("/pricing", h.GetPublicPricing)
}

// Module wraps the full system settings module for easy wiring in main.go
type Module struct {
	handler *Handler
	service *Service
}

// NewModuleFromDB creates a ready-to-use system settings module
func NewModuleFromDB(repo *Repository, v *validator.Validator) *Module {
	svc := NewService(repo)
	return &Module{
		handler: NewHandler(svc, v),
		service: svc,
	}
}

// Service returns the underlying service (used by other modules to read live settings)
func (m *Module) Service() *Service {
	return m.service
}

// RegisterRoutes is a convenience helper that registers both admin and public routes
func (m *Module) RegisterRoutes(r chi.Router, authMiddleware *middleware.AuthMiddleware) {
	RegisterAdminRoutes(r, m.handler, authMiddleware.Authenticate, authMiddleware.RequireAdmin)
	RegisterPublicRoutes(r, m.handler)
}
