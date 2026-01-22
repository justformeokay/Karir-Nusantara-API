package chat

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jung-kurt/gofpdf"
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

// CloseConversation allows company to close their own conversation
// PATCH /company/chat/conversations/{id}/close
func (h *Handler) CloseConversation(w http.ResponseWriter, r *http.Request) {
	companyID := middleware.GetUserID(r.Context())
	if companyID == 0 {
		response.Error(w, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
		return
	}
	
	conversationIDStr := chi.URLParam(r, "id")
	conversationID, err := strconv.ParseUint(conversationIDStr, 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_ID", "Invalid conversation ID")
		return
	}
	
	// Verify conversation belongs to company
	conv, _, err := h.service.GetConversation(r.Context(), conversationID, companyID)
	if err != nil {
		response.Error(w, http.StatusNotFound, "NOT_FOUND", "Conversation not found")
		return
	}
	
	// Verify ownership
	if conv.CompanyID != companyID {
		response.Error(w, http.StatusForbidden, "FORBIDDEN", "You don't have permission to close this conversation")
		return
	}
	
	// Close conversation
	req := UpdateConversationStatusRequest{Status: "closed"}
	err = h.service.UpdateConversationStatus(r.Context(), conversationID, &req)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "UPDATE_ERROR", "Failed to close conversation")
		return
	}
	
	response.Success(w, http.StatusOK, "Conversation closed successfully", nil)
}

// DownloadConversationPDF generates and downloads conversation as PDF
// GET /company/chat/conversations/{id}/pdf or /admin/chat/conversations/{id}/pdf
func (h *Handler) DownloadConversationPDF(w http.ResponseWriter, r *http.Request) {
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
	
	// Get conversation with messages
	conv, messages, err := h.service.GetConversation(r.Context(), conversationID, userID)
	if err != nil {
		response.Error(w, http.StatusNotFound, "NOT_FOUND", "Conversation not found")
		return
	}
	
	// Generate PDF content
	pdfContent := h.generateConversationPDF(conv, messages)
	
	// Set headers for PDF download
	filename := fmt.Sprintf("conversation_%d_%s.pdf", conversationID, time.Now().Format("20060102"))
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(pdfContent)))
	
	// Write PDF to response
	w.Write(pdfContent)
}

// generateConversationPDF generates a professional PDF from conversation using gofpdf
func (h *Handler) generateConversationPDF(conv *ConversationWithDetails, messages []*ChatMessageWithSender) []byte {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetAutoPageBreak(true, 15)
	pdf.AddPage()
	
	// Header - Company Logo/Branding Area
	pdf.SetFillColor(41, 128, 185) // Blue header
	pdf.Rect(0, 0, 210, 40, "F")
	
	// Company Name
	pdf.SetTextColor(255, 255, 255)
	pdf.SetFont("Arial", "B", 24)
	pdf.SetY(12)
	pdf.CellFormat(0, 10, "KARIR NUSANTARA", "", 1, "C", false, 0, "")
	
	pdf.SetFont("Arial", "", 10)
	pdf.SetY(24)
	pdf.CellFormat(0, 6, "Platform Pencarian Kerja Terpercaya Indonesia", "", 1, "C", false, 0, "")
	pdf.SetY(30)
	pdf.CellFormat(0, 6, "Laporan Percakapan Support", "", 1, "C", false, 0, "")
	
	pdf.Ln(15)
	
	// Document Info Section
	pdf.SetTextColor(0, 0, 0)
	pdf.SetFont("Arial", "B", 12)
	pdf.CellFormat(0, 7, "INFORMASI PERCAKAPAN", "", 1, "L", false, 0, "")
	pdf.SetLineWidth(0.5)
	pdf.SetDrawColor(41, 128, 185)
	pdf.Line(10, pdf.GetY(), 200, pdf.GetY())
	pdf.Ln(5)
	
	// Conversation Details in Table Format
	pdf.SetFont("Arial", "", 10)
	pdf.SetFillColor(245, 245, 245)
	
	// Row 1: ID & Status
	pdf.SetFont("Arial", "B", 10)
	pdf.CellFormat(40, 7, "ID Percakapan", "1", 0, "L", true, 0, "")
	pdf.SetFont("Arial", "", 10)
	pdf.CellFormat(60, 7, fmt.Sprintf("#%d", conv.ID), "1", 0, "L", false, 0, "")
	pdf.SetFont("Arial", "B", 10)
	pdf.CellFormat(30, 7, "Status", "1", 0, "L", true, 0, "")
	pdf.SetFont("Arial", "", 10)
	statusText := strings.ToUpper(conv.Status)
	pdf.CellFormat(60, 7, statusText, "1", 1, "L", false, 0, "")
	
	// Row 2: Company
	pdf.SetFont("Arial", "B", 10)
	pdf.CellFormat(40, 7, "Perusahaan", "1", 0, "L", true, 0, "")
	pdf.SetFont("Arial", "", 10)
	pdf.CellFormat(150, 7, conv.CompanyName, "1", 1, "L", false, 0, "")
	
	// Row 3: Title
	pdf.SetFont("Arial", "B", 10)
	pdf.CellFormat(40, 7, "Judul", "1", 0, "L", true, 0, "")
	pdf.SetFont("Arial", "", 10)
	pdf.MultiCell(150, 7, conv.Title, "1", "L", false)
	
	// Row 4: Subject
	pdf.SetFont("Arial", "B", 10)
	pdf.CellFormat(40, 7, "Subjek", "1", 0, "L", true, 0, "")
	pdf.SetFont("Arial", "", 10)
	pdf.MultiCell(150, 7, conv.Subject, "1", "L", false)
	
	// Row 5: Category
	pdf.SetFont("Arial", "B", 10)
	pdf.CellFormat(40, 7, "Kategori", "1", 0, "L", true, 0, "")
	pdf.SetFont("Arial", "", 10)
	categoryMap := map[string]string{
		"complaint": "Komplain",
		"helpdesk":  "Help Desk",
		"general":   "Umum",
		"urgent":    "Urgent",
	}
	categoryText := categoryMap[conv.Category]
	if categoryText == "" {
		categoryText = conv.Category
	}
	pdf.CellFormat(150, 7, categoryText, "1", 1, "L", false, 0, "")
	
	// Row 6: Created At
	pdf.SetFont("Arial", "B", 10)
	pdf.CellFormat(40, 7, "Dibuat Tanggal", "1", 0, "L", true, 0, "")
	pdf.SetFont("Arial", "", 10)
	pdf.CellFormat(60, 7, conv.CreatedAt.Format("02 Jan 2006 15:04"), "1", 0, "L", false, 0, "")
	
	// Row 6 continued: Total Messages
	pdf.SetFont("Arial", "B", 10)
	pdf.CellFormat(30, 7, "Total Pesan", "1", 0, "L", true, 0, "")
	pdf.SetFont("Arial", "", 10)
	pdf.CellFormat(60, 7, fmt.Sprintf("%d pesan", len(messages)), "1", 1, "L", false, 0, "")
	
	pdf.Ln(8)
	
	// Messages Section
	pdf.SetFont("Arial", "B", 12)
	pdf.CellFormat(0, 7, "RIWAYAT PERCAKAPAN", "", 1, "L", false, 0, "")
	pdf.SetLineWidth(0.5)
	pdf.Line(10, pdf.GetY(), 200, pdf.GetY())
	pdf.Ln(5)
	
	// Messages
	for i, msg := range messages {
		// Check if we need a new page
		if pdf.GetY() > 250 {
			pdf.AddPage()
		}
		
		// Message Header
		pdf.SetFillColor(240, 248, 255) // Light blue for company, light gray for admin
		if msg.SenderType == "admin" {
			pdf.SetFillColor(248, 248, 248)
		}
		
		pdf.SetFont("Arial", "B", 10)
		senderLabel := "Perusahaan"
		if msg.SenderType == "admin" {
			senderLabel = "Admin Support"
		}
		
		headerText := fmt.Sprintf("#%d - %s: %s", i+1, senderLabel, msg.SenderName)
		pdf.CellFormat(140, 6, headerText, "LTR", 0, "L", true, 0, "")
		
		pdf.SetFont("Arial", "", 8)
		pdf.SetTextColor(100, 100, 100)
		timeText := msg.CreatedAt.Format("02 Jan 2006 15:04:05")
		pdf.CellFormat(50, 6, timeText, "RTL", 1, "R", true, 0, "")
		pdf.SetTextColor(0, 0, 0)
		
		// Message Content
		pdf.SetFont("Arial", "", 10)
		
		// Handle attachment info
		messageText := msg.Message
		if msg.AttachmentType.Valid && msg.AttachmentType.String != "" {
			attachmentInfo := ""
			if msg.AttachmentType.String == "image" {
				attachmentInfo = "[Lampiran: Gambar]"
			} else if msg.AttachmentType.String == "audio" {
				attachmentInfo = "[Lampiran: Pesan Suara]"
			}
			if msg.AttachmentFilename.Valid && msg.AttachmentFilename.String != "" {
				attachmentInfo += " " + msg.AttachmentFilename.String
			}
			messageText = attachmentInfo + "\n" + messageText
		}
		
		// Message body with word wrap
		pdf.MultiCell(190, 6, messageText, "LRB", "L", true)
		
		pdf.Ln(3)
	}
	
	// Footer with generation info
	pdf.SetY(-25)
	pdf.SetFont("Arial", "I", 8)
	pdf.SetTextColor(128, 128, 128)
	pdf.CellFormat(0, 5, fmt.Sprintf("Dokumen ini digenerate otomatis pada %s", time.Now().Format("02 January 2006 15:04:05")), "", 1, "C", false, 0, "")
	pdf.CellFormat(0, 5, "Karir Nusantara - Platform Pencarian Kerja Terpercaya", "", 1, "C", false, 0, "")
	pdf.CellFormat(0, 5, "www.karirnusantara.com", "", 1, "C", false, 0, "")
	
	// Generate PDF to bytes
	var buf bytes.Buffer
	err := pdf.Output(&buf)
	if err != nil {
		log.Printf("Error generating PDF: %v", err)
		return []byte{}
	}
	
	return buf.Bytes()
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
