package policies

import (
	"github.com/go-chi/chi/v5"
)

// RegisterRoutes registers policies routes
func RegisterRoutes(r chi.Router) {
	h := &Handler{}

	r.Route("/policies", func(r chi.Router) {
		r.Get("/privacy-policy/pdf", h.GeneratePrivacyPolicyPDF)
		r.Get("/terms-of-service/pdf", h.GenerateTermsOfServicePDF)
	})
}
