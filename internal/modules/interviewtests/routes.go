package interviewtests

import (
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
)

// Module represents the interview tests module
type Module struct {
	handler *Handler
}

// NewModule creates a new interview tests module
func NewModule(db *sqlx.DB) *Module {
	repo := NewRepository(db)
	service := NewService(repo)
	handler := NewHandler(service)

	return &Module{
		handler: handler,
	}
}

// RegisterAdminRoutes registers admin-only interview test routes
func (m *Module) RegisterAdminRoutes(r chi.Router) {
	// These routes are already wrapped with admin authentication
	r.Route("/interview-tests", func(r chi.Router) {
		r.Get("/", m.handler.GetAll)
		r.Post("/", m.handler.Create)
		r.Get("/{id}", m.handler.GetByID)
		r.Put("/{id}", m.handler.Update)
		r.Delete("/{id}", m.handler.Delete)
		r.Post("/{id}/publish", m.handler.Publish)
		r.Post("/{id}/duplicate", m.handler.Duplicate)
	})
}
