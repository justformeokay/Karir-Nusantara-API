package passwordreset

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/karirnusantara/api/internal/shared/response"
)

// Handler handles HTTP requests for password reset
type Handler struct {
	service *Service
}

// NewHandler creates a new password reset handler
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// ForgotPassword handles forgot password requests
func (h *Handler) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	var req ForgotPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid request body")
		return
	}

	// Validate email
	if req.Email == "" {
		response.UnprocessableEntity(w, "Validasi gagal", map[string]string{
			"email": "Email harus diisi",
		})
		return
	}

	// Request password reset
	if err := h.service.RequestPasswordReset(req.Email); err != nil {
		log.Printf("[PASSWORD RESET ERROR] Failed to send reset email for %s: %v", req.Email, err)
		response.InternalServerError(w, "Gagal mengirim email reset password")
		return
	}

	// Always return success to prevent email enumeration
	response.OK(w, "Jika email terdaftar, link reset password telah dikirim ke email Anda", ForgotPasswordResponse{
		Message: "Jika email terdaftar, link reset password telah dikirim ke email Anda",
	})
}

// VerifyToken verifies if a reset token is valid
func (h *Handler) VerifyToken(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		response.BadRequest(w, "Token diperlukan")
		return
	}

	valid, email, err := h.service.VerifyToken(token)
	if err != nil {
		response.InternalServerError(w, "Gagal memverifikasi token")
		return
	}

	if !valid {
		response.OK(w, "Token verification result", VerifyTokenResponse{
			Valid: false,
		})
		return
	}

	response.OK(w, "Token is valid", VerifyTokenResponse{
		Valid: valid,
		Email: email,
	})
}

// ResetPassword handles password reset with token
func (h *Handler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var req ResetPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid request body")
		return
	}

	// Validate
	if req.Token == "" {
		response.UnprocessableEntity(w, "Validasi gagal", map[string]string{
			"token": "Token diperlukan",
		})
		return
	}

	if req.NewPassword == "" {
		response.UnprocessableEntity(w, "Validasi gagal", map[string]string{
			"new_password": "Password baru harus diisi",
		})
		return
	}

	if len(req.NewPassword) < 8 {
		response.UnprocessableEntity(w, "Validasi gagal", map[string]string{
			"new_password": "Password minimal 8 karakter",
		})
		return
	}

	// Reset password
	if err := h.service.ResetPassword(req.Token, req.NewPassword); err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	response.OK(w, "Password berhasil diubah", ResetPasswordResponse{
		Message: "Password berhasil diubah",
	})
}
