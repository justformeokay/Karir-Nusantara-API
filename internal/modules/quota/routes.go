package quota

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

// MiddlewareFunc defines the middleware function type
type MiddlewareFunc func(http.Handler) http.Handler

// RegisterRoutes registers quota routes
func RegisterRoutes(r chi.Router, h *Handler, authenticate, requireCompany MiddlewareFunc) {
	// Company quota routes (authenticated + company role required)
	r.Route("/company", func(r chi.Router) {
		r.Use(authenticate)
		r.Use(requireCompany)

		r.Get("/quota", h.GetQuota)
		r.Get("/payments", h.GetPayments)
		r.Get("/payments/info", h.GetPaymentInfo)
		r.Post("/payments/proof", h.SubmitPaymentProof)
	})
}
