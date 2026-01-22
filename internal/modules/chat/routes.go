package chat

import (
	"github.com/go-chi/chi/v5"
	"github.com/karirnusantara/api/internal/middleware"
)

// RegisterRoutes registers chat routes
func RegisterRoutes(r chi.Router, h *Handler, authMiddleware *middleware.AuthMiddleware) {
	// Company routes - require company authentication
	r.Group(func(r chi.Router) {
		r.Use(authMiddleware.Authenticate)
		r.Use(authMiddleware.RequireRole("company"))
		
		r.Post("/company/chat/conversations", h.CreateConversation)
		r.Get("/company/chat/conversations", h.GetMyConversations)
		r.Get("/company/chat/conversations/{id}", h.GetConversation)
		r.Post("/company/chat/conversations/{id}/messages", h.SendMessage)
		r.Post("/company/chat/upload", h.UploadAttachment)
	})
	
	// Admin routes - require admin authentication
	r.Group(func(r chi.Router) {
		r.Use(authMiddleware.Authenticate)
		r.Use(authMiddleware.RequireRole("admin"))
		
		r.Get("/admin/chat/conversations", h.GetAllConversations)
		r.Get("/admin/chat/conversations/{id}", h.GetConversation)
		r.Post("/admin/chat/conversations/{id}/messages", h.SendMessage)
		r.Patch("/admin/chat/conversations/{id}/status", h.UpdateConversationStatus)
		r.Post("/admin/chat/upload", h.UploadAttachment)
	})
}
