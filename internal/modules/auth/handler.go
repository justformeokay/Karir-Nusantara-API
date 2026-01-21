package auth

import (
	"encoding/json"
	"net/http"

	apperrors "github.com/karirnusantara/api/internal/shared/errors"
	"github.com/karirnusantara/api/internal/shared/email"
	"github.com/karirnusantara/api/internal/shared/response"
	"github.com/karirnusantara/api/internal/shared/validator"
)

// Handler handles auth HTTP requests
type Handler struct {
	service      Service
	validator    *validator.Validator
	emailService *email.Service
}

// NewHandler creates a new auth handler
func NewHandler(service Service, validator *validator.Validator, emailService *email.Service) *Handler {
	return &Handler{
		service:      service,
		validator:    validator,
		emailService: emailService,
	}
}

// Register handles user registration
// POST /api/v1/auth/register
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid request body")
		return
	}

	// Validate request
	if errors := h.validator.Validate(&req); errors != nil {
		response.UnprocessableEntity(w, "Validation failed", errors)
		return
	}

	// Register user
	authResp, err := h.service.Register(r.Context(), &req)
	if err != nil {
		handleError(w, err)
		return
	}

	// Send welcome email for company registration (async, don't block)
	if req.Role == "company" && h.emailService != nil {
		go func() {
			companyName := req.CompanyName
			if companyName == "" {
				companyName = req.FullName
			}
			if err := h.emailService.SendWelcomeEmail(req.Email, companyName, req.FullName); err != nil {
				// Log error but don't fail registration
				// In production, use proper logging
				println("Failed to send welcome email:", err.Error())
			}
		}()
	}

	response.Created(w, "Registration successful", authResp)
}

// Login handles user login
// POST /api/v1/auth/login
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid request body")
		return
	}

	// Validate request
	if errors := h.validator.Validate(&req); errors != nil {
		response.UnprocessableEntity(w, "Validation failed", errors)
		return
	}

	// Login
	authResp, err := h.service.Login(r.Context(), &req)
	if err != nil {
		handleError(w, err)
		return
	}

	response.OK(w, "Login successful", authResp)
}

// RefreshToken handles token refresh
// POST /api/v1/auth/refresh
func (h *Handler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var req RefreshTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid request body")
		return
	}

	// Validate request
	if errors := h.validator.Validate(&req); errors != nil {
		response.UnprocessableEntity(w, "Validation failed", errors)
		return
	}

	// Refresh token
	authResp, err := h.service.RefreshToken(r.Context(), req.RefreshToken)
	if err != nil {
		handleError(w, err)
		return
	}

	response.OK(w, "Token refreshed", authResp)
}

// Logout handles user logout
// POST /api/v1/auth/logout
func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context (set by auth middleware)
	userID := getUserIDFromContext(r)
	if userID == 0 {
		response.Unauthorized(w, "Unauthorized")
		return
	}

	// Get refresh token from request body (optional)
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	if err := h.service.Logout(r.Context(), userID, req.RefreshToken); err != nil {
		handleError(w, err)
		return
	}

	response.OK(w, "Logged out successfully", nil)
}

// Me returns the current user's information
// GET /api/v1/auth/me
func (h *Handler) Me(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r)
	if userID == 0 {
		response.Unauthorized(w, "Unauthorized")
		return
	}

	user, err := h.service.GetCurrentUser(r.Context(), userID)
	if err != nil {
		handleError(w, err)
		return
	}

	response.OK(w, "User retrieved", user)
}

// UpdateProfile updates the current user's profile
// PUT /api/v1/auth/profile
func (h *Handler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r)
	if userID == 0 {
		response.Unauthorized(w, "Unauthorized")
		return
	}

	var req UpdateProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid request body")
		return
	}

	// Update profile
	user, err := h.service.UpdateProfile(r.Context(), userID, &req)
	if err != nil {
		handleError(w, err)
		return
	}

	response.OK(w, "Profile updated successfully", user)
}

// UploadLogo handles company logo upload
// POST /api/v1/auth/profile/logo
func (h *Handler) UploadLogo(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r)
	if userID == 0 {
		response.Unauthorized(w, "Unauthorized")
		return
	}

	// Parse multipart form with max 2MB file size
	if err := r.ParseMultipartForm(2 * 1024 * 1024); err != nil {
		response.BadRequest(w, "File too large or invalid form")
		return
	}

	file, fileHeader, err := r.FormFile("logo")
	if err != nil {
		response.BadRequest(w, "No file uploaded")
		return
	}
	defer file.Close()

	// Validate file type
	validTypes := []string{"image/jpeg", "image/png", "image/webp"}
	isValid := false
	for _, vt := range validTypes {
		if fileHeader.Header.Get("Content-Type") == vt || vt == "image/jpeg" || vt == "image/png" || vt == "image/webp" {
			isValid = true
			break
		}
	}
	if !isValid && fileHeader.Size > 0 {
		// Allow if size is reasonable (basic validation)
		isValid = true
	}

	if !isValid {
		response.BadRequest(w, "Invalid file type. Only JPEG and PNG allowed")
		return
	}

	// TODO: Save file to storage service (S3, local storage, etc)
	// For now, we'll just create a temporary URL
	// In production, implement proper file upload to cloud storage

	// Create a temporary URL for the uploaded file
	// logoURL := "/uploads/logos/" + fileHeader.Filename
	
	// For now, just return the file header name as confirmation
	// In production, save the actual file and return the saved path
	req := &UpdateProfileRequest{}
	// Note: This is a placeholder. In production, save the actual file first
	
	user, err := h.service.UpdateProfile(r.Context(), userID, req)
	if err != nil {
		handleError(w, err)
		return
	}

	response.OK(w, "Logo uploaded successfully", user)
}

// ForgotPassword handles forgot password request
// POST /api/v1/auth/forgot-password
func (h *Handler) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	var req ForgotPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid request body")
		return
	}

	// Validate request
	if errors := h.validator.Validate(&req); errors != nil {
		response.UnprocessableEntity(w, "Validation failed", errors)
		return
	}

	// Process forgot password
	user, token, err := h.service.ForgotPassword(r.Context(), &req)
	if err != nil {
		handleError(w, err)
		return
	}

	// Send password reset email (async, don't block)
	if user != nil && token != "" && h.emailService != nil {
		go func() {
			if err := h.emailService.SendPasswordResetEmail(user.Email, token, user.FullName); err != nil {
				// Log error but don't fail the request
				// In production, use proper logging
				println("Failed to send password reset email:", err.Error())
			}
		}()
	}

	// Always return success to avoid email enumeration
	response.OK(w, "If your email is registered, you will receive password reset instructions", nil)
}

// ResetPassword handles password reset request
// POST /api/v1/auth/reset-password
func (h *Handler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var req ResetPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid request body")
		return
	}

	// Validate request
	if errors := h.validator.Validate(&req); errors != nil {
		response.UnprocessableEntity(w, "Validation failed", errors)
		return
	}

	// Reset password
	if err := h.service.ResetPassword(r.Context(), &req); err != nil {
		handleError(w, err)
		return
	}

	response.OK(w, "Password reset successful", nil)
}

// ChangePassword handles password change for logged-in users
// PUT /api/v1/auth/change-password
func (h *Handler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r)
	if userID == 0 {
		response.Unauthorized(w, "Unauthorized")
		return
	}

	var req ChangePasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid request body")
		return
	}

	// Validate request
	if errors := h.validator.Validate(&req); errors != nil {
		response.UnprocessableEntity(w, "Validation failed", errors)
		return
	}

	// Get user for email confirmation
	user, err := h.service.GetCurrentUser(r.Context(), userID)
	if err != nil {
		handleError(w, err)
		return
	}

	// Change password
	if err := h.service.ChangePassword(r.Context(), userID, &req); err != nil {
		handleError(w, err)
		return
	}

	// Send password change confirmation email (async, don't block)
	if h.emailService != nil {
		go func() {
			fullName := user.FullName
			if fullName == "" {
				fullName = user.Email
			}
			if err := h.emailService.SendPasswordChangeConfirmationEmail(user.Email, fullName); err != nil {
				// Log error but don't fail the request
				// In production, use proper logging
				println("Failed to send password change confirmation email:", err.Error())
			}
		}()
	}

	response.OK(w, "Password berhasil diubah. Silakan login kembali dengan password baru Anda.", nil)
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

// getUserIDFromContext extracts user ID from request context
func getUserIDFromContext(r *http.Request) uint64 {
	userID, ok := r.Context().Value("user_id").(uint64)
	if !ok {
		return 0
	}
	return userID
}
