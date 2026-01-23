package company

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/karirnusantara/api/internal/shared/hashid"
	"github.com/karirnusantara/api/internal/shared/response"
)

// Handler handles company-related HTTP requests
type Handler struct {
	service      Service
	fileService  *FileService
}

// NewHandler creates a new company handler
func NewHandler(service Service, fileService *FileService) *Handler {
	return &Handler{
		service:     service,
		fileService: fileService,
	}
}

// GetCompanyProfile retrieves the current user's company profile
// GET /api/v1/company/profile
func (h *Handler) GetCompanyProfile(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(uint64)
	if !ok || userID == 0 {
		response.Unauthorized(w, "Unauthorized")
		return
	}

	company, err := h.service.GetCompanyByUserID(r.Context(), userID)
	if err != nil {
		response.InternalServerError(w, "Failed to get company profile")
		return
	}

	if company == nil {
		response.NotFound(w, "Company profile not found")
		return
	}

	response.OK(w, "Company profile retrieved", company)
}

// UpdateCompanyProfile updates the current user's company profile
// PUT /api/v1/company/profile
func (h *Handler) UpdateCompanyProfile(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(uint64)
	if !ok || userID == 0 {
		response.Unauthorized(w, "Unauthorized")
		return
	}

	var req UpdateCompanyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid request body")
		return
	}

	company, err := h.service.CreateOrUpdateCompany(r.Context(), userID, &req)
	if err != nil {
		response.InternalServerError(w, "Failed to update company profile")
		return
	}

	response.OK(w, "Company profile updated successfully", company)
}

// UploadCompanyLogo uploads a company logo
// POST /api/v1/companies/logo
func (h *Handler) UploadCompanyLogo(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(uint64)
	if !ok || userID == 0 {
		response.Unauthorized(w, "Unauthorized")
		return
	}

	// Parse multipart form
	if err := r.ParseMultipartForm(5 * 1024 * 1024); err != nil {
		response.BadRequest(w, "File too large or invalid form")
		return
	}

	file, fileHeader, err := r.FormFile("logo")
	if err != nil {
		response.BadRequest(w, "No file uploaded")
		return
	}
	defer file.Close()

	// Validate file
	if !ValidateImageFile(fileHeader.Filename) {
		response.BadRequest(w, "Invalid file type. Only JPG, PNG, and PDF allowed")
		return
	}
	if !ValidateImageFileSize(fileHeader.Size, 5*1024*1024) {
		response.BadRequest(w, "File too large. Maximum 5MB")
		return
	}

	// Get actual company entity from database
	company, err := h.service.GetCompanyEntityByUserID(r.Context(), userID)
	if err != nil {
		response.InternalServerError(w, "Failed to get company")
		return
	}
	if company == nil {
		response.NotFound(w, "Company not found. Please complete your company profile first")
		return
	}

	// Save file
	filePath, err := h.fileService.SaveCompanyDocument(company.ID, "logo", file, fileHeader.Filename)
	if err != nil {
		response.InternalServerError(w, "Failed to save file")
		return
	}

	// Update company logo URL
	if err := h.service.UpdateCompanyLogoURL(r.Context(), company.ID, "/docs/companies/"+filePath); err != nil {
		response.InternalServerError(w, "Failed to update company")
		return
	}

	response.OK(w, "Logo uploaded successfully", map[string]interface{}{
		"logo_url": "/docs/companies/" + filePath,
	})
}

// UploadCompanyDocument uploads a company legal document
// POST /api/v1/companies/documents
func (h *Handler) UploadCompanyDocument(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(uint64)
	if !ok || userID == 0 {
		response.Unauthorized(w, "Unauthorized")
		return
	}

	// Parse multipart form
	if err := r.ParseMultipartForm(10 * 1024 * 1024); err != nil {
		response.BadRequest(w, "File too large or invalid form")
		return
	}

	docType := r.FormValue("doc_type")
	if docType == "" {
		response.BadRequest(w, "doc_type is required")
		return
	}

	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		response.BadRequest(w, "No file uploaded")
		return
	}
	defer file.Close()

	// Validate file
	if !ValidateImageFile(fileHeader.Filename) {
		response.BadRequest(w, "Invalid file type. Only JPG, PNG, and PDF allowed")
		return
	}
	if !ValidateImageFileSize(fileHeader.Size, 10*1024*1024) {
		response.BadRequest(w, "File too large. Maximum 10MB")
		return
	}

	// Get actual company entity from database
	company, err := h.service.GetCompanyEntityByUserID(r.Context(), userID)
	if err != nil {
		response.InternalServerError(w, "Failed to get company")
		return
	}
	if company == nil {
		response.NotFound(w, "Company not found. Please complete your company profile first")
		return
	}

	// Save file
	filePath, err := h.fileService.SaveCompanyDocument(company.ID, docType, file, fileHeader.Filename)
	if err != nil {
		response.InternalServerError(w, "Failed to save file")
		return
	}

	// Update company document URL based on type
	fullPath := "/docs/companies/" + filePath
	if err := h.service.UpdateCompanyDocument(r.Context(), company.ID, docType, fullPath); err != nil {
		response.InternalServerError(w, "Failed to update company")
		return
	}

	response.OK(w, fmt.Sprintf("%s uploaded successfully", docType), map[string]interface{}{
		"doc_type": docType,
		"file_url": fullPath,
	})
}

// GetPublicCompanyProfile retrieves a company profile for public viewing (by job seekers)
// GET /api/v1/companies/{id}
func (h *Handler) GetPublicCompanyProfile(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		response.BadRequest(w, "Company ID is required")
		return
	}

	// Try to parse as hash_id first, then as numeric ID
	var companyID uint64
	var err error

	// First try hash_id decode
	companyID, err = hashid.Decode(idStr)
	if err != nil {
		// If hash_id decode fails, this could be a direct numeric ID (for backwards compatibility)
		// but for security reasons, we reject numeric IDs
		response.BadRequest(w, "Invalid company ID format")
		return
	}

	company, err := h.service.GetPublicCompanyByID(r.Context(), companyID)
	if err != nil {
		response.InternalServerError(w, "Failed to get company profile")
		return
	}

	if company == nil {
		response.NotFound(w, "Company not found")
		return
	}

	response.OK(w, "Company profile retrieved", company)
}
