package interviewtests

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
)

// MiddlewareFunc is a type alias for middleware functions
type MiddlewareFunc func(http.Handler) http.Handler

// Module represents the interview tests module (admin use)
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

// CompanyModule represents the interview tests module for company portal
type CompanyModule struct {
	handler *CompanyHandler
}

// NewCompanyModule creates a new company interview tests module
func NewCompanyModule(db *sqlx.DB, companyService CompanyServiceInterface) *CompanyModule {
	repo := NewRepository(db)
	service := NewService(repo)
	handler := NewCompanyHandler(service, companyService)

	return &CompanyModule{
		handler: handler,
	}
}

// RegisterRoutes registers company interview test routes
func (m *CompanyModule) RegisterRoutes(r chi.Router, authenticate func(http.Handler) http.Handler) {
	r.Route("/company/interview-tests", func(r chi.Router) {
		r.Use(authenticate)

		// Library: browse public admin tests
		r.Get("/library", m.handler.GetLibrary)

		// Own tests
		r.Get("/", m.handler.GetMyTests)
		r.Post("/", m.handler.CreateMyTest)
		r.Put("/{id}", m.handler.UpdateMyTest)
		r.Delete("/{id}", m.handler.DeleteMyTest)
		r.Post("/{id}/publish", m.handler.PublishMyTest)

		// Copy admin test
		r.Post("/{id}/copy", m.handler.CopyFromAdmin)
	})

	// Application-level assignment endpoints
	r.Route("/company/applications/{applicationId}", func(r chi.Router) {
		r.Use(authenticate)
		r.Post("/assign-test", m.handler.AssignTest)
		r.Get("/interview-tests", m.handler.GetApplicationSubmissions)
	})
}

// JobSeekerModule represents the interview tests module for job seeker portal
type JobSeekerModule struct {
	handler *JobSeekerHandler
}

// NewJobSeekerModule creates a new job seeker interview tests module
func NewJobSeekerModule(db *sqlx.DB) *JobSeekerModule {
	repo := NewRepository(db)
	service := NewService(repo)
	handler := NewJobSeekerHandler(service)

	return &JobSeekerModule{
		handler: handler,
	}
}

// RegisterJobSeekerRoutes registers job seeker interview test routes
func (m *JobSeekerModule) RegisterJobSeekerRoutes(r chi.Router, authenticate func(http.Handler) http.Handler) {
	r.Route("/jobseeker/interview-tests", func(r chi.Router) {
		r.Use(authenticate)
		r.Get("/", m.handler.GetMySubmissions)
		r.Get("/{submissionId}", m.handler.GetTestForSubmission)
		r.Post("/{submissionId}/submit", m.handler.SubmitAnswers)
	})
}
