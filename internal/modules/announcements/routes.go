package announcements

import (
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"

	"github.com/karirnusantara/api/internal/middleware"
)

// Module represents the announcements module
type Module struct {
	handler        *Handler
	authMiddleware *middleware.AuthMiddleware
}

// NewModule creates a new announcements module
func NewModule(db *sqlx.DB, authMiddleware *middleware.AuthMiddleware) *Module {
	repo := NewRepository(db)
	service := NewService(repo)
	handler := NewHandler(service)

	return &Module{
		handler:        handler,
		authMiddleware: authMiddleware,
	}
}

// RegisterRoutes registers announcement routes
func (m *Module) RegisterRoutes(r chi.Router) {
	// Public routes - for fetching active announcements
	r.Route("/announcements", func(r chi.Router) {
		r.Get("/", m.handler.GetAllPublic)
		r.Get("/notifications", m.handler.GetPublicNotifications)
		r.Get("/banners", m.handler.GetPublicBanners)
		r.Get("/information", m.handler.GetPublicInformation)
	})
}

// RegisterAdminRoutes registers admin-only announcement routes
func (m *Module) RegisterAdminRoutes(r chi.Router) {
	// These routes are already wrapped with admin authentication in admin module
	r.Route("/announcements", func(r chi.Router) {
		r.Get("/", m.handler.List)
		r.Post("/", m.handler.Create)
		r.Get("/{id}", m.handler.GetByID)
		r.Put("/{id}", m.handler.Update)
		r.Delete("/{id}", m.handler.Delete)
		r.Patch("/{id}/toggle", m.handler.ToggleStatus)
	})
}
