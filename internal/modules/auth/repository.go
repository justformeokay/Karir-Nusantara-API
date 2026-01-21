package auth

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
)

// Repository defines the auth repository interface
type Repository interface {
	// User operations
	CreateUser(ctx context.Context, user *User) error
	GetUserByID(ctx context.Context, id uint64) (*User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	UpdateUser(ctx context.Context, user *User) error
	DeleteUser(ctx context.Context, id uint64) error
	EmailExists(ctx context.Context, email string) (bool, error)
	GetCompanyByUserID(ctx context.Context, userID uint64) (*CompanyData, error)

	// Refresh token operations
	CreateRefreshToken(ctx context.Context, token *RefreshToken) error
	GetRefreshTokenByHash(ctx context.Context, hash string) (*RefreshToken, error)
	RevokeRefreshToken(ctx context.Context, hash string) error
	RevokeAllUserTokens(ctx context.Context, userID uint64) error
	CleanupExpiredTokens(ctx context.Context) error

	// Password reset operations
	CreatePasswordResetToken(ctx context.Context, token *PasswordResetToken) error
	GetPasswordResetToken(ctx context.Context, tokenStr string) (*PasswordResetToken, error)
	MarkPasswordResetTokenAsUsed(ctx context.Context, id uint64) error
	DeleteExpiredResetTokens(ctx context.Context) error
}

// mysqlRepository implements Repository for MySQL
type mysqlRepository struct {
	db *sqlx.DB
}

// NewRepository creates a new auth repository
func NewRepository(db *sqlx.DB) Repository {
	return &mysqlRepository{db: db}
}

// CreateUser creates a new user
func (r *mysqlRepository) CreateUser(ctx context.Context, user *User) error {
	query := `
		INSERT INTO users (
			email, password_hash, role, full_name, phone, avatar_url,
			is_active, is_verified, created_at, updated_at
		) VALUES (
			?, ?, ?, ?, ?, ?,
			?, ?, NOW(), NOW()
		)
	`

	result, err := r.db.ExecContext(ctx, query,
		user.Email, user.PasswordHash, user.Role, user.FullName, user.Phone, user.AvatarURL,
		user.IsActive, user.IsVerified,
	)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	user.ID = uint64(id)
	return nil
}

// GetUserByID retrieves a user by ID
func (r *mysqlRepository) GetUserByID(ctx context.Context, id uint64) (*User, error) {
	query := `
		SELECT id, email, password_hash, role, full_name, phone, avatar_url,
			   is_active, is_verified, email_verified_at, created_at, updated_at, deleted_at
		FROM users
		WHERE id = ? AND deleted_at IS NULL
	`

	var user User
	if err := r.db.GetContext(ctx, &user, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}

	return &user, nil
}

// GetCompanyByUserID retrieves company data for a user
func (r *mysqlRepository) GetCompanyByUserID(ctx context.Context, userID uint64) (*CompanyData, error) {
	query := `
		SELECT company_name, company_logo_url, company_description, company_website,
		       company_industry, company_size, company_location, company_phone, company_email,
		       company_address, company_city, company_province, company_postal_code,
		       established_year, employee_count, company_status,
		       ktp_founder_url, akta_pendirian_url, npwp_url, nib_url
		FROM companies
		WHERE user_id = ? AND deleted_at IS NULL
	`

	var company CompanyData
	if err := r.db.GetContext(ctx, &company, query, userID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get company by user id: %w", err)
	}

	return &company, nil
}

// GetUserByEmail retrieves a user by email
func (r *mysqlRepository) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	query := `
		SELECT id, email, password_hash, role, full_name, phone, avatar_url,
			   is_active, is_verified, email_verified_at, created_at, updated_at, deleted_at
		FROM users
		WHERE email = ? AND deleted_at IS NULL
	`

	var user User
	if err := r.db.GetContext(ctx, &user, query, email); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	return &user, nil
}

// UpdateUser updates a user
func (r *mysqlRepository) UpdateUser(ctx context.Context, user *User) error {
	query := `
		UPDATE users SET
			password_hash = ?,
			full_name = ?,
			phone = ?,
			avatar_url = ?,
			is_active = ?,
			is_verified = ?,
			updated_at = NOW()
		WHERE id = ? AND deleted_at IS NULL
	`

	_, err := r.db.ExecContext(ctx, query,
		user.PasswordHash, user.FullName, user.Phone, user.AvatarURL,
		user.IsActive, user.IsVerified, user.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

// DeleteUser soft deletes a user
func (r *mysqlRepository) DeleteUser(ctx context.Context, id uint64) error {
	query := `UPDATE users SET deleted_at = NOW() WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}

// EmailExists checks if an email already exists
func (r *mysqlRepository) EmailExists(ctx context.Context, email string) (bool, error) {
	query := `SELECT COUNT(*) FROM users WHERE email = ? AND deleted_at IS NULL`
	var count int
	if err := r.db.GetContext(ctx, &count, query, email); err != nil {
		return false, fmt.Errorf("failed to check email existence: %w", err)
	}
	return count > 0, nil
}

// CreateRefreshToken creates a new refresh token
func (r *mysqlRepository) CreateRefreshToken(ctx context.Context, token *RefreshToken) error {
	query := `
		INSERT INTO refresh_tokens (user_id, token_hash, expires_at, device_info, ip_address, created_at)
		VALUES (?, ?, ?, ?, ?, NOW())
	`

	result, err := r.db.ExecContext(ctx, query,
		token.UserID, token.TokenHash, token.ExpiresAt, token.DeviceInfo, token.IPAddress,
	)
	if err != nil {
		return fmt.Errorf("failed to create refresh token: %w", err)
	}

	id, _ := result.LastInsertId()
	token.ID = uint64(id)
	return nil
}

// GetRefreshTokenByHash retrieves a refresh token by hash
func (r *mysqlRepository) GetRefreshTokenByHash(ctx context.Context, hash string) (*RefreshToken, error) {
	query := `
		SELECT id, user_id, token_hash, expires_at, revoked_at, device_info, ip_address, created_at
		FROM refresh_tokens
		WHERE token_hash = ?
	`

	var token RefreshToken
	if err := r.db.GetContext(ctx, &token, query, hash); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get refresh token: %w", err)
	}

	return &token, nil
}

// RevokeRefreshToken revokes a refresh token
func (r *mysqlRepository) RevokeRefreshToken(ctx context.Context, hash string) error {
	query := `UPDATE refresh_tokens SET revoked_at = NOW() WHERE token_hash = ?`
	_, err := r.db.ExecContext(ctx, query, hash)
	if err != nil {
		return fmt.Errorf("failed to revoke refresh token: %w", err)
	}
	return nil
}

// RevokeAllUserTokens revokes all refresh tokens for a user
func (r *mysqlRepository) RevokeAllUserTokens(ctx context.Context, userID uint64) error {
	query := `UPDATE refresh_tokens SET revoked_at = NOW() WHERE user_id = ? AND revoked_at IS NULL`
	_, err := r.db.ExecContext(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("failed to revoke all user tokens: %w", err)
	}
	return nil
}

// CleanupExpiredTokens removes expired tokens
func (r *mysqlRepository) CleanupExpiredTokens(ctx context.Context) error {
	query := `DELETE FROM refresh_tokens WHERE expires_at < NOW() OR revoked_at IS NOT NULL`
	_, err := r.db.ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to cleanup expired tokens: %w", err)
	}
	return nil
}

// CreatePasswordResetToken creates a new password reset token
func (r *mysqlRepository) CreatePasswordResetToken(ctx context.Context, token *PasswordResetToken) error {
	query := `
		INSERT INTO password_reset_tokens (email, token, expires_at, created_at)
		VALUES (?, ?, ?, NOW())
	`

	result, err := r.db.ExecContext(ctx, query, token.Email, token.Token, token.ExpiresAt)
	if err != nil {
		return fmt.Errorf("failed to create password reset token: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	token.ID = uint64(id)
	return nil
}

// GetPasswordResetToken retrieves a password reset token
func (r *mysqlRepository) GetPasswordResetToken(ctx context.Context, tokenStr string) (*PasswordResetToken, error) {
	query := `
		SELECT id, email, token, expires_at, used_at, created_at
		FROM password_reset_tokens
		WHERE token = ? AND used_at IS NULL AND expires_at > NOW()
	`

	var token PasswordResetToken
	if err := r.db.GetContext(ctx, &token, query, tokenStr); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get password reset token: %w", err)
	}

	return &token, nil
}

// MarkPasswordResetTokenAsUsed marks a password reset token as used
func (r *mysqlRepository) MarkPasswordResetTokenAsUsed(ctx context.Context, id uint64) error {
	query := `UPDATE password_reset_tokens SET used_at = NOW() WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to mark password reset token as used: %w", err)
	}
	return nil
}

// DeleteExpiredResetTokens deletes expired password reset tokens
func (r *mysqlRepository) DeleteExpiredResetTokens(ctx context.Context) error {
	query := `DELETE FROM password_reset_tokens WHERE expires_at < NOW() OR used_at IS NOT NULL`
	_, err := r.db.ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to delete expired reset tokens: %w", err)
	}
	return nil
}
