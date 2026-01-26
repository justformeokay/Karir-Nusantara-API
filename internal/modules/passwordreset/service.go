package passwordreset

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/karirnusantara/api/internal/shared/email"
	"golang.org/x/crypto/bcrypt"
)

// Service handles business logic for password reset
type Service struct {
	repo         *Repository
	emailService *email.Service
}

// NewService creates a new password reset service
func NewService(repo *Repository, emailService *email.Service) *Service {
	return &Service{
		repo:         repo,
		emailService: emailService,
	}
}

// RequestPasswordReset initiates password reset process
func (s *Service) RequestPasswordReset(emailAddr string) error {
	// Check if user exists
	_, fullName, err := s.repo.GetUserByEmail(emailAddr)
	if err != nil {
		// Don't reveal if user exists or not for security
		return nil
	}

	// Invalidate any existing tokens for this user
	if err := s.repo.InvalidateUserTokens(emailAddr); err != nil {
		return fmt.Errorf("failed to invalidate existing tokens: %w", err)
	}

	// Generate new token
	token, err := s.repo.GenerateToken()
	if err != nil {
		return fmt.Errorf("failed to generate token: %w", err)
	}

	// Save token to database
	if err := s.repo.CreateResetToken(emailAddr, token); err != nil {
		return fmt.Errorf("failed to create reset token: %w", err)
	}

	// Send email
	if err := s.emailService.SendPasswordResetEmail(emailAddr, token, fullName); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

// VerifyToken checks if a token is valid
func (s *Service) VerifyToken(token string) (bool, string, error) {
	resetToken, err := s.repo.GetTokenByValue(token)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, "", nil
		}
		return false, "", err
	}

	// Check if token is expired
	if time.Now().After(resetToken.ExpiresAt) {
		return false, "", nil
	}

	// Check if token has been used
	if resetToken.UsedAt != nil {
		return false, "", nil
	}

	return true, resetToken.Email, nil
}

// ResetPassword resets user's password using token
func (s *Service) ResetPassword(token, newPassword string) error {
	// Verify token first
	resetToken, err := s.repo.GetTokenByValue(token)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("token tidak valid")
		}
		return err
	}

	// Check if token is expired
	if time.Now().After(resetToken.ExpiresAt) {
		return fmt.Errorf("token telah kadaluarsa")
	}

	// Check if token has been used
	if resetToken.UsedAt != nil {
		return fmt.Errorf("token telah digunakan")
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// Update password
	if err := s.repo.UpdatePassword(resetToken.Email, string(hashedPassword)); err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	// Mark token as used
	if err := s.repo.MarkTokenAsUsed(token); err != nil {
		return fmt.Errorf("failed to mark token as used: %w", err)
	}

	// Get user info for email
	_, fullName, err := s.repo.GetUserByEmail(resetToken.Email)
	if err == nil {
		// Send confirmation email (don't fail if this errors)
		_ = s.emailService.SendPasswordChangeConfirmationEmail(resetToken.Email, fullName)
	}

	return nil
}

// CleanupExpiredTokens removes expired tokens from database
func (s *Service) CleanupExpiredTokens() error {
	return s.repo.DeleteExpiredTokens()
}
