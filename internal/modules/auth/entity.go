package auth

import (
	"database/sql"
	"time"
)

// User roles
const (
	RoleJobSeeker = "job_seeker"
	RoleCompany   = "company"
	RoleAdmin     = "admin"
)

// User represents a user entity
type User struct {
	ID                uint64         `db:"id" json:"id"`
	Email             string         `db:"email" json:"email"`
	PasswordHash      string         `db:"password_hash" json:"-"`
	Role              string         `db:"role" json:"role"`
	FullName          string         `db:"full_name" json:"full_name"`
	Phone             sql.NullString `db:"phone" json:"phone,omitempty"`
	AvatarURL         sql.NullString `db:"avatar_url" json:"avatar_url,omitempty"`
	CompanyName       sql.NullString `db:"company_name" json:"company_name,omitempty"`
	CompanyDescription sql.NullString `db:"company_description" json:"company_description,omitempty"`
	CompanyWebsite    sql.NullString `db:"company_website" json:"company_website,omitempty"`
	CompanyLogoURL    sql.NullString `db:"company_logo_url" json:"company_logo_url,omitempty"`
	IsActive          bool           `db:"is_active" json:"is_active"`
	IsVerified        bool           `db:"is_verified" json:"is_verified"`
	CompanyStatus     sql.NullString `db:"company_status" json:"company_status,omitempty"`
	EmailVerifiedAt   sql.NullTime   `db:"email_verified_at" json:"email_verified_at,omitempty"`
	CreatedAt         time.Time      `db:"created_at" json:"created_at"`
	UpdatedAt         time.Time      `db:"updated_at" json:"updated_at"`
	DeletedAt         sql.NullTime   `db:"deleted_at" json:"-"`
}

// UserResponse represents the user response (safe for public exposure)
type UserResponse struct {
	ID            uint64 `json:"id"`
	Email         string `json:"email"`
	Role          string `json:"role"`
	FullName      string `json:"full_name"`
	Phone         string `json:"phone,omitempty"`
	AvatarURL     string `json:"avatar_url,omitempty"`
	CompanyName   string `json:"company_name,omitempty"`
	CompanyLogoURL string `json:"company_logo_url,omitempty"`
	IsActive      bool   `json:"is_active"`
	IsVerified    bool   `json:"is_verified"`
	CompanyStatus string `json:"company_status,omitempty"`
	CreatedAt     string `json:"created_at"`
}

// ToResponse converts User to UserResponse
func (u *User) ToResponse() *UserResponse {
	resp := &UserResponse{
		ID:         u.ID,
		Email:      u.Email,
		Role:       u.Role,
		FullName:   u.FullName,
		IsActive:   u.IsActive,
		IsVerified: u.IsVerified,
		CreatedAt:  u.CreatedAt.Format(time.RFC3339),
	}

	if u.Phone.Valid {
		resp.Phone = u.Phone.String
	}
	if u.AvatarURL.Valid {
		resp.AvatarURL = u.AvatarURL.String
	}
	if u.CompanyName.Valid {
		resp.CompanyName = u.CompanyName.String
	}
	if u.CompanyLogoURL.Valid {
		resp.CompanyLogoURL = u.CompanyLogoURL.String
	}
	if u.CompanyStatus.Valid {
		resp.CompanyStatus = u.CompanyStatus.String
	}

	return resp
}

// RefreshToken represents a refresh token entity
type RefreshToken struct {
	ID         uint64       `db:"id"`
	UserID     uint64       `db:"user_id"`
	TokenHash  string       `db:"token_hash"`
	ExpiresAt  time.Time    `db:"expires_at"`
	RevokedAt  sql.NullTime `db:"revoked_at"`
	DeviceInfo string       `db:"device_info"`
	IPAddress  string       `db:"ip_address"`
	CreatedAt  time.Time    `db:"created_at"`
}

// Request DTOs

// RegisterRequest represents a registration request
type RegisterRequest struct {
	Email       string `json:"email" validate:"required,email"`
	Password    string `json:"password" validate:"required,min=8,password"`
	FullName    string `json:"full_name" validate:"required,min=2,max=255"`
	Phone       string `json:"phone,omitempty" validate:"omitempty,phone"`
	Role        string `json:"role" validate:"required,oneof=job_seeker company"`
	CompanyName string `json:"company_name,omitempty" validate:"required_if=Role company"`
}

// LoginRequest represents a login request
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// RefreshTokenRequest represents a refresh token request
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// UpdateProfileRequest represents a profile update request
type UpdateProfileRequest struct {
	FullName           string `json:"full_name,omitempty" validate:"omitempty,min=2,max=255"`
	Phone              string `json:"phone,omitempty" validate:"omitempty,phone"`
	CompanyName        string `json:"company_name,omitempty" validate:"omitempty,min=2,max=255"`
	CompanyDescription string `json:"company_description,omitempty" validate:"omitempty,min=50"`
	CompanyWebsite     string `json:"company_website,omitempty" validate:"omitempty,url"`
	CompanySize        string `json:"company_size,omitempty"`
	CompanyLocation    string `json:"company_location,omitempty"`
	CompanyIndustry    string `json:"company_industry,omitempty"`
}

// Response DTOs

// AuthResponse represents an authentication response
type AuthResponse struct {
	User         *UserResponse `json:"user"`
	AccessToken  string        `json:"access_token"`
	RefreshToken string        `json:"refresh_token"`
	ExpiresIn    int64         `json:"expires_in"`
}

// TokenClaims represents JWT token claims
type TokenClaims struct {
	UserID   uint64 `json:"user_id"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	TokenType string `json:"token_type"`
}

// IsJobSeeker checks if user is a job seeker
func (u *User) IsJobSeeker() bool {
	return u.Role == RoleJobSeeker
}

// IsCompany checks if user is a company
func (u *User) IsCompany() bool {
	return u.Role == RoleCompany
}

// IsAdmin checks if user is an admin
func (u *User) IsAdmin() bool {
	return u.Role == RoleAdmin
}
