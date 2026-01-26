package passwordreset

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

// Repository handles database operations for password reset
type Repository struct {
	db *sqlx.DB
}

// NewRepository creates a new password reset repository
func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

// GenerateToken generates a secure random token
func (r *Repository) GenerateToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// CreateResetToken creates a new password reset token
func (r *Repository) CreateResetToken(email, token string) error {
	expiresAt := time.Now().Add(TokenExpiration)

	query := `
		INSERT INTO password_reset_tokens (email, token, expires_at)
		VALUES (?, ?, ?)
	`

	_, err := r.db.Exec(query, email, token, expiresAt)
	return err
}

// GetTokenByValue retrieves a reset token by its value
func (r *Repository) GetTokenByValue(token string) (*PasswordResetToken, error) {
	query := `
		SELECT id, email, token, expires_at, used_at, created_at
		FROM password_reset_tokens
		WHERE token = ?
		LIMIT 1
	`

	var resetToken PasswordResetToken
	err := r.db.QueryRow(query, token).Scan(
		&resetToken.ID,
		&resetToken.Email,
		&resetToken.Token,
		&resetToken.ExpiresAt,
		&resetToken.UsedAt,
		&resetToken.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &resetToken, nil
}

// MarkTokenAsUsed marks a token as used
func (r *Repository) MarkTokenAsUsed(token string) error {
	query := `
		UPDATE password_reset_tokens
		SET used_at = NOW()
		WHERE token = ?
	`

	_, err := r.db.Exec(query, token)
	return err
}

// GetUserByEmail retrieves a user by email
func (r *Repository) GetUserByEmail(email string) (int, string, error) {
	query := `
		SELECT id, full_name
		FROM users
		WHERE email = ? AND role = 'job_seeker'
		LIMIT 1
	`

	var userID int
	var fullName string
	err := r.db.QueryRow(query, email).Scan(&userID, &fullName)

	if err != nil {
		if err == sql.ErrNoRows {
			return 0, "", fmt.Errorf("user not found")
		}
		return 0, "", err
	}

	return userID, fullName, nil
}

// UpdatePassword updates user's password
func (r *Repository) UpdatePassword(email, hashedPassword string) error {
	query := `
		UPDATE users
		SET password = ?, updated_at = NOW()
		WHERE email = ? AND role = 'job_seeker'
	`

	_, err := r.db.Exec(query, hashedPassword, email)
	return err
}

// DeleteExpiredTokens deletes expired tokens (cleanup)
func (r *Repository) DeleteExpiredTokens() error {
	query := `
		DELETE FROM password_reset_tokens
		WHERE expires_at < NOW()
		OR (used_at IS NOT NULL AND created_at < DATE_SUB(NOW(), INTERVAL 1 DAY))
	`

	_, err := r.db.Exec(query)
	return err
}

// InvalidateUserTokens invalidates all existing tokens for a user
func (r *Repository) InvalidateUserTokens(email string) error {
	query := `
		UPDATE password_reset_tokens
		SET used_at = NOW()
		WHERE email = ? AND used_at IS NULL
	`

	_, err := r.db.Exec(query, email)
	return err
}
