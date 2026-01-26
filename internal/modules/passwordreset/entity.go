package passwordreset

import (
	"time"
)

// PasswordResetToken represents a password reset token
type PasswordResetToken struct {
	ID        int        `json:"id" db:"id"`
	Email     string     `json:"email" db:"email"`
	Token     string     `json:"token" db:"token"`
	ExpiresAt time.Time  `json:"expires_at" db:"expires_at"`
	UsedAt    *time.Time `json:"used_at" db:"used_at"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
}

// Request types
type ForgotPasswordRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type ResetPasswordRequest struct {
	Token       string `json:"token" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=8"`
}

type VerifyTokenRequest struct {
	Token string `json:"token" validate:"required"`
}

// Response types
type ForgotPasswordResponse struct {
	Message string `json:"message"`
}

type ResetPasswordResponse struct {
	Message string `json:"message"`
}

type VerifyTokenResponse struct {
	Valid bool   `json:"valid"`
	Email string `json:"email,omitempty"`
}

// Constants
const (
	TokenExpiration = 5 * time.Minute // 5 minutes
)
