package chat

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/karirnusantara/api/internal/middleware"
	"github.com/karirnusantara/api/internal/shared/response"
	"github.com/karirnusantara/api/internal/shared/validator"
)

// Handler handles chat HTTP requests
type Handler struct {
	service      Service
	validator    *validator.Validator
	uploadFolder string
}

// NewHandler creates a new chat handler
func NewHandler(service Service, v *validator.Validator, uploadFolder string) *Handler {
	return &Handler{
		service:      service,
		validator:    v,
		uploadFolder: uploadFolder,
	}
}

// CreateConversation creates a new conversation
// POST /company/chat/conversations
func (h *Handler) CreateConversation(w http.ResponseWriter, r *http.Request) {
	companyID := middleware.GetUserID(r.Context())
	if companyID == 0 {
		response.Error(w, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
		return
	}
	
	var req CreateConversationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}
	
	if errors := h.validator.Validate(&req); errors != nil {
		response.UnprocessableEntity(w, "Validation failed", errors)
		return
	}
	
	conv, err := h.service.CreateConversation(r.Context(), companyID, &req)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "CREATE_ERROR", "Failed to create conversation")
		return
	}
	
	response.Success(w, http.StatusCreated, "Conversation created successfully", conv)
}

// GetMyConversations gets all conversations for logged in company
// GET /company/chat/conversations
func (h *Handler) GetMyConversations(w http.ResponseWriter, r *http.Request) {
	companyID := middleware.GetUserID(r.Context())
	if companyID == 0 {
		response.Error(w, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
		return
	}
	
	convs, err := h.service.GetMyConversations(r.Context(), companyID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "FETCH_ERROR", "Failed to fetch conversations")
		return
	}
	
	response.Success(w, http.StatusOK, "Conversations retrieved successfully", convs)
}

// GetConversation gets a conversation with messages
// GET /company/chat/conversations/{id}
func (h *Handler) GetConversation(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == 0 {
		response.Error(w, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
		return
	}
	
	conversationIDStr := chi.URLParam(r, "id")
	conversationID, err := strconv.ParseUint(conversationIDStr, 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_ID", "Invalid conversation ID")
		return
	}
	
	conv, messages, err := h.service.GetConversation(r.Context(), conversationID, userID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "FETCH_ERROR", err.Error())
		return
	}
	
	result := map[string]interface{}{
		"conversation": conv,
		"messages":     messages,
	}
	
	response.Success(w, http.StatusOK, "Conversation retrieved successfully", result)
}

// SendMessage sends a message in a conversation
// POST /company/chat/conversations/{id}/messages
func (h *Handler) SendMessage(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == 0 {
		response.Error(w, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
		return
	}
	
	// Get sender type from context (set by middleware)
	role, ok := r.Context().Value("user_role").(string)
	if !ok {
		role = "company" // default
	}
	senderType := "company"
	if role == "admin" {
		senderType = "admin"
	}
	
	conversationIDStr := chi.URLParam(r, "id")
	conversationID, err := strconv.ParseUint(conversationIDStr, 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_ID", "Invalid conversation ID")
		return
	}
	
	var req SendMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}
	
	if errors := h.validator.Validate(&req); errors != nil {
		response.UnprocessableEntity(w, "Validation failed", errors)
		return
	}
	
	message, err := h.service.SendMessage(r.Context(), conversationID, userID, senderType, &req)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "SEND_ERROR", err.Error())
		return
	}
	
	response.Success(w, http.StatusCreated, "Message sent successfully", message)
}

// GetAllConversations gets all conversations (admin only)
// GET /admin/chat/conversations
func (h *Handler) GetAllConversations(w http.ResponseWriter, r *http.Request) {
	convs, err := h.service.GetAllConversations(r.Context())
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "FETCH_ERROR", "Failed to fetch conversations")
		return
	}
	
	response.Success(w, http.StatusOK, "Conversations retrieved successfully", convs)
}

// UpdateConversationStatus updates conversation status (admin only)
// PATCH /admin/chat/conversations/{id}/status
func (h *Handler) UpdateConversationStatus(w http.ResponseWriter, r *http.Request) {
	conversationIDStr := chi.URLParam(r, "id")
	conversationID, err := strconv.ParseUint(conversationIDStr, 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_ID", "Invalid conversation ID")
		return
	}
	
	var req UpdateConversationStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}
	
	if errors := h.validator.Validate(&req); errors != nil {
		response.UnprocessableEntity(w, "Validation failed", errors)
		return
	}
	
	err = h.service.UpdateConversationStatus(r.Context(), conversationID, &req)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "UPDATE_ERROR", "Failed to update conversation status")
		return
	}
	
	response.Success(w, http.StatusOK, "Conversation status updated successfully", nil)
}

// UploadAttachment handles file upload for chat attachments
// POST /company/chat/upload or /admin/chat/upload
func (h *Handler) UploadAttachment(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == 0 {
		response.Error(w, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
		return
	}

	// Parse multipart form (max 10MB)
	if err := r.ParseMultipartForm(10 * 1024 * 1024); err != nil {
		response.Error(w, http.StatusBadRequest, "FILE_TOO_LARGE", "File too large. Maximum 10MB")
		return
	}

	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		response.Error(w, http.StatusBadRequest, "NO_FILE", "No file uploaded")
		return
	}
	defer file.Close()

	// Get attachment type from form
	attachmentType := r.FormValue("type") // "image" or "audio"
	if attachmentType != "image" && attachmentType != "audio" {
		response.Error(w, http.StatusBadRequest, "INVALID_TYPE", "Type must be 'image' or 'audio'")
		return
	}

	// Validate file type
	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	if attachmentType == "image" {
		validExts := []string{".jpg", ".jpeg", ".png", ".gif"}
		if !contains(validExts, ext) {
			response.Error(w, http.StatusBadRequest, "INVALID_FILE", "Invalid image file. Only JPG, PNG, GIF allowed")
			return
		}
	} else if attachmentType == "audio" {
		validExts := []string{".mp3", ".wav", ".ogg", ".m4a", ".webm"}
		if !contains(validExts, ext) {
			response.Error(w, http.StatusBadRequest, "INVALID_FILE", "Invalid audio file. Only MP3, WAV, OGG, M4A, WEBM allowed")
			return
		}
	}

	// Validate file size (10MB)
	if fileHeader.Size > 10*1024*1024 {
		response.Error(w, http.StatusBadRequest, "FILE_TOO_LARGE", "File too large. Maximum 10MB")
		return
	}

	// Create chat uploads directory
	uploadDir := filepath.Join(h.uploadFolder, "chat")
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		response.Error(w, http.StatusInternalServerError, "UPLOAD_ERROR", "Failed to create upload directory")
		return
	}

	// Generate unique filename
	timestamp := time.Now().Unix()
	filename := fmt.Sprintf("%d_%d%s", userID, timestamp, ext)
	filePath := filepath.Join(uploadDir, filename)

	// Save file
	dst, err := os.Create(filePath)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "UPLOAD_ERROR", "Failed to save file")
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		response.Error(w, http.StatusInternalServerError, "UPLOAD_ERROR", "Failed to write file")
		return
	}

	// Return file URL (relative path)
	fileURL := fmt.Sprintf("/docs/chat/%s", filename)
	
	response.Success(w, http.StatusOK, "File uploaded successfully", map[string]interface{}{
		"url":      fileURL,
		"type":     attachmentType,
		"filename": fileHeader.Filename,
	})
}

// Helper function to check if slice contains string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
