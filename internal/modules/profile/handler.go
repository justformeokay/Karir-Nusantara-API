package profile

import (
	"database/sql"
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
	apperrors "github.com/karirnusantara/api/internal/shared/errors"
	"github.com/karirnusantara/api/internal/shared/hashid"
	"github.com/karirnusantara/api/internal/shared/response"
	"github.com/karirnusantara/api/internal/shared/validator"
)

// Handler handles profile HTTP requests
type Handler struct {
	service   Service
	validator *validator.Validator
	docsPath  string
}

// NewHandler creates a new profile handler
func NewHandler(service Service, validator *validator.Validator, docsPath string) *Handler {
	return &Handler{
		service:   service,
		validator: validator,
		docsPath:  docsPath,
	}
}

// ========================================
// Profile Endpoints
// ========================================

// GetProfile handles GET /api/v1/profile
func (h *Handler) GetProfile(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == 0 {
		response.Unauthorized(w, "Unauthorized")
		return
	}

	profile, err := h.service.GetProfile(r.Context(), userID)
	if err != nil {
		handleError(w, err)
		return
	}

	response.OK(w, "Profile retrieved successfully", profile)
}

// UpdateProfile handles PUT /api/v1/profile
func (h *Handler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == 0 {
		response.Unauthorized(w, "Unauthorized")
		return
	}

	var req UpdateProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid request body")
		return
	}

	if errors := h.validator.Validate(&req); errors != nil {
		response.UnprocessableEntity(w, "Validation failed", errors)
		return
	}

	profile, err := h.service.CreateOrUpdateProfile(r.Context(), userID, &req)
	if err != nil {
		handleError(w, err)
		return
	}

	response.OK(w, "Profile updated successfully", profile)
}

// DeleteProfile handles DELETE /api/v1/profile
func (h *Handler) DeleteProfile(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == 0 {
		response.Unauthorized(w, "Unauthorized")
		return
	}

	if err := h.service.DeleteProfile(r.Context(), userID); err != nil {
		handleError(w, err)
		return
	}

	response.NoContent(w)
}

// ========================================
// Document Endpoints
// ========================================

// GetDocuments handles GET /api/v1/profile/documents
func (h *Handler) GetDocuments(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == 0 {
		response.Unauthorized(w, "Unauthorized")
		return
	}

	docs, err := h.service.GetDocuments(r.Context(), userID)
	if err != nil {
		handleError(w, err)
		return
	}

	response.OK(w, "Documents retrieved successfully", docs)
}

// GetDocument handles GET /api/v1/profile/documents/{id}
func (h *Handler) GetDocument(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == 0 {
		response.Unauthorized(w, "Unauthorized")
		return
	}

	docID, err := h.parseDocumentID(r)
	if err != nil {
		response.BadRequest(w, "Invalid document ID")
		return
	}

	doc, err := h.service.GetDocumentByID(r.Context(), docID, userID)
	if err != nil {
		handleError(w, err)
		return
	}

	response.OK(w, "Document retrieved successfully", doc)
}

// UploadDocument handles POST /api/v1/profile/documents
func (h *Handler) UploadDocument(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == 0 {
		response.Unauthorized(w, "Unauthorized")
		return
	}

	// Parse multipart form (max 10MB)
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		response.BadRequest(w, "Failed to parse form data. Max file size is 10MB")
		return
	}

	// Get document type
	docType := r.FormValue("document_type")
	if docType == "" {
		docType = "cv_uploaded"
	}

	// Validate document type
	validTypes := map[string]bool{
		"cv_uploaded": true,
		"certificate": true,
		"transcript":  true,
		"portfolio":   true,
		"ktp":         true,
		"other":       true,
	}
	if !validTypes[docType] {
		response.BadRequest(w, "Invalid document type")
		return
	}

	// Get file
	file, header, err := r.FormFile("file")
	if err != nil {
		response.BadRequest(w, "File is required")
		return
	}
	defer file.Close()

	// Validate file type
	allowedTypes := map[string]bool{
		"application/pdf":    true,
		"application/msword": true,
		"application/vnd.openxmlformats-officedocument.wordprocessingml.document": true,
		"image/jpeg": true,
		"image/png":  true,
	}

	contentType := header.Header.Get("Content-Type")
	if !allowedTypes[contentType] {
		response.BadRequest(w, "Invalid file type. Allowed: PDF, DOC, DOCX, JPG, PNG")
		return
	}

	// Create directory for user documents
	userDir := filepath.Join(h.docsPath, "applicants", strconv.FormatUint(userID, 10))
	if err := os.MkdirAll(userDir, 0755); err != nil {
		response.InternalServerError(w, "Failed to create directory")
		return
	}

	// Generate unique filename
	ext := filepath.Ext(header.Filename)
	filename := fmt.Sprintf("%s_%d%s", docType, time.Now().Unix(), ext)
	filePath := filepath.Join(userDir, filename)

	// Save file
	dst, err := os.Create(filePath)
	if err != nil {
		response.InternalServerError(w, "Failed to save file")
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		response.InternalServerError(w, "Failed to save file")
		return
	}

	// Create document record
	isPrimary := r.FormValue("is_primary") == "true"
	description := r.FormValue("description")

	doc := &ApplicantDocument{
		UserID:      userID,
		DocType:     DocumentType(docType),
		DocName:     header.Filename,
		DocURL:      fmt.Sprintf("/docs/applicants/%d/%s", userID, filename),
		FileSize:    sql.NullInt64{Int64: header.Size, Valid: true},
		MimeType:    sql.NullString{String: contentType, Valid: true},
		IsPrimary:   isPrimary,
		Description: sql.NullString{String: description, Valid: description != ""},
	}

	docResp, err := h.service.CreateDocument(r.Context(), userID, doc)
	if err != nil {
		handleError(w, err)
		return
	}

	response.Created(w, "Document uploaded successfully", docResp)
}

// UpdateDocument handles PUT /api/v1/profile/documents/{id}
func (h *Handler) UpdateDocument(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == 0 {
		response.Unauthorized(w, "Unauthorized")
		return
	}

	docID, err := h.parseDocumentID(r)
	if err != nil {
		response.BadRequest(w, "Invalid document ID")
		return
	}

	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid request body")
		return
	}

	doc, err := h.service.UpdateDocument(r.Context(), docID, userID, req.Name, req.Description)
	if err != nil {
		handleError(w, err)
		return
	}

	response.OK(w, "Document updated successfully", doc)
}

// DeleteDocument handles DELETE /api/v1/profile/documents/{id}
func (h *Handler) DeleteDocument(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == 0 {
		response.Unauthorized(w, "Unauthorized")
		return
	}

	docID, err := h.parseDocumentID(r)
	if err != nil {
		response.BadRequest(w, "Invalid document ID")
		return
	}

	if err := h.service.DeleteDocument(r.Context(), docID, userID); err != nil {
		handleError(w, err)
		return
	}

	response.NoContent(w)
}

// SetPrimaryDocument handles POST /api/v1/profile/documents/{id}/primary
func (h *Handler) SetPrimaryDocument(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == 0 {
		response.Unauthorized(w, "Unauthorized")
		return
	}

	docID, err := h.parseDocumentID(r)
	if err != nil {
		response.BadRequest(w, "Invalid document ID")
		return
	}

	if err := h.service.SetPrimaryDocument(r.Context(), docID, userID); err != nil {
		handleError(w, err)
		return
	}

	response.OK(w, "Document set as primary successfully", nil)
}

// parseDocumentID parses document ID from URL (supports both hash_id and numeric id)
func (h *Handler) parseDocumentID(r *http.Request) (uint64, error) {
	idParam := chi.URLParam(r, "id")
	if idParam == "" {
		return 0, fmt.Errorf("missing id parameter")
	}

	// Try to decode as hash_id first
	if strings.HasPrefix(idParam, "kn_") {
		return hashid.Decode(idParam)
	}

	// Try to parse as numeric ID
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid id format")
	}

	return id, nil
}

// handleError handles errors and sends appropriate response
func handleError(w http.ResponseWriter, err error) {
	appErr := apperrors.GetAppError(err)
	if appErr != nil {
		if appErr.Details != nil {
			response.ErrorWithDetails(w, appErr.HTTPStatus, appErr.Code, appErr.Message, appErr.Details)
		} else {
			response.Error(w, appErr.HTTPStatus, appErr.Code, appErr.Message)
		}
		return
	}
	response.InternalServerError(w, "An unexpected error occurred")
}
