package chat

import (
	"github.com/go-chi/chi/v5"
	"github.com/karirnusantara/api/internal/middleware"
)

// RegisterRoutes registers chat routes
func RegisterRoutes(r chi.Router, h *Handler, authMiddleware *middleware.AuthMiddleware) {
	// Company routes - require company authentication
	r.Route("/company/chat", func(r chi.Router) {
		r.Use(authMiddleware.Authenticate)
		r.Use(authMiddleware.RequireRole("company"))
		
		r.Post("/conversations", h.CreateConversation)
		r.Get("/conversations", h.GetMyConversations)
		r.Get("/conversations/{id}", h.GetConversation)
		r.Post("/conversations/{id}/messages", h.SendMessage)
		r.Patch("/conversations/{id}/close", h.CloseConversation)
		r.Get("/conversations/{id}/pdf", h.DownloadConversationPDF)
		r.Post("/upload", h.UploadAttachment)
	})
	
	// Admin routes - require admin authentication
	r.Route("/admin/chat", func(r chi.Router) {
		r.Use(authMiddleware.Authenticate)
		r.Use(authMiddleware.RequireRole("admin"))
		
		r.Get("/conversations", h.GetAllConversations)
		r.Get("/conversations/{id}", h.GetConversation)
		r.Post("/conversations/{id}/messages", h.SendMessage)
		r.Patch("/conversations/{id}/status", h.UpdateConversationStatus)
		r.Get("/conversations/{id}/pdf", h.DownloadConversationPDF)
		r.Post("/upload", h.UploadAttachment)
	})
}
