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
	ID            uint64         `db:"id" json:"id"`
	Email         string         `db:"email" json:"email"`
	PasswordHash  string         `db:"password_hash" json:"-"`
	Role          string         `db:"role" json:"role"`
	FullName      string         `db:"full_name" json:"full_name"`
	Phone         sql.NullString `db:"phone" json:"phone,omitempty"`
	AvatarURL     sql.NullString `db:"avatar_url" json:"avatar_url,omitempty"`
	IsActive      bool           `db:"is_active" json:"is_active"`
	IsVerified    bool           `db:"is_verified" json:"is_verified"`
	EmailVerifiedAt sql.NullTime   `db:"email_verified_at" json:"email_verified_at,omitempty"`
	CreatedAt     time.Time      `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time      `db:"updated_at" json:"updated_at"`
	DeletedAt     sql.NullTime   `db:"deleted_at" json:"-"`
}

// UserResponse represents the user response (safe for public exposure)
type UserResponse struct {
	ID         uint64 `json:"id"`
	Email      string `json:"email"`
	Role       string `json:"role"`
	FullName   string `json:"full_name"`
	Phone      string `json:"phone,omitempty"`
	AvatarURL  string `json:"avatar_url,omitempty"`
	IsActive   bool   `json:"is_active"`
	IsVerified bool   `json:"is_verified"`
	CreatedAt  string `json:"created_at"`
}

// UserWithCompanyResponse includes user data and company information
type UserWithCompanyResponse struct {
	ID                  uint64 `json:"id"`
	Email               string `json:"email"`
	Role                string `json:"role"`
	FullName            string `json:"full_name"`
	Phone               string `json:"phone,omitempty"`
	AvatarURL           string `json:"avatar_url,omitempty"`
	IsActive            bool   `json:"is_active"`
	IsVerified          bool   `json:"is_verified"`
	CreatedAt           string `json:"created_at"`
	CompanyName         string `json:"company_name,omitempty"`
	CompanyLogoURL      string `json:"company_logo_url,omitempty"`
	CompanyDescription  string `json:"company_description,omitempty"`
	CompanyWebsite      string `json:"company_website,omitempty"`
	CompanyIndustry     string `json:"company_industry,omitempty"`
	CompanySize         string `json:"company_size,omitempty"`
	CompanyLocation     string `json:"company_location,omitempty"`
	CompanyPhone        string `json:"company_phone,omitempty"`
	CompanyEmail        string `json:"company_email,omitempty"`
	CompanyAddress      string `json:"company_address,omitempty"`
	CompanyCity         string `json:"company_city,omitempty"`
	CompanyProvince     string `json:"company_province,omitempty"`
	CompanyPostalCode   string `json:"company_postal_code,omitempty"`
	EstablishedYear     int64  `json:"established_year,omitempty"`
	EmployeeCount       int64  `json:"employee_count,omitempty"`
	CompanyStatus       string `json:"company_status,omitempty"`
	KTPFounderURL       string `json:"ktp_founder_url,omitempty"`
	AktaPendirianURL    string `json:"akta_pendirian_url,omitempty"`
	NPWPURL             string `json:"npwp_url,omitempty"`
	NIBURL              string `json:"nib_url,omitempty"`
}

// CompanyData represents company information from companies table
type CompanyData struct {
	CompanyName        sql.NullString `db:"company_name"`
	CompanyLogoURL     sql.NullString `db:"company_logo_url"`
	CompanyDescription sql.NullString `db:"company_description"`
	CompanyWebsite     sql.NullString `db:"company_website"`
	CompanyIndustry    sql.NullString `db:"company_industry"`
	CompanySize        sql.NullString `db:"company_size"`
	CompanyLocation    sql.NullString `db:"company_location"`
	CompanyPhone       sql.NullString `db:"company_phone"`
	CompanyEmail       sql.NullString `db:"company_email"`
	CompanyAddress     sql.NullString `db:"company_address"`
	CompanyCity        sql.NullString `db:"company_city"`
	CompanyProvince    sql.NullString `db:"company_province"`
	CompanyPostalCode  sql.NullString `db:"company_postal_code"`
	EstablishedYear    sql.NullInt64  `db:"established_year"`
	EmployeeCount      sql.NullInt64  `db:"employee_count"`
	CompanyStatus      sql.NullString `db:"company_status"`
	KTPFounderURL      sql.NullString `db:"ktp_founder_url"`
	AktaPendirianURL   sql.NullString `db:"akta_pendirian_url"`
	NPWPURL            sql.NullString `db:"npwp_url"`
	NIBURL             sql.NullString `db:"nib_url"`
}

// ToResponseWithCompany converts User + CompanyData to UserWithCompanyResponse
func (u *User) ToResponseWithCompany(company *CompanyData) *UserWithCompanyResponse {
	resp := &UserWithCompanyResponse{
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

	// Populate company fields if company data exists
	if company != nil {
		if company.CompanyName.Valid {
			resp.CompanyName = company.CompanyName.String
		}
		if company.CompanyLogoURL.Valid {
			resp.CompanyLogoURL = company.CompanyLogoURL.String
		}
		if company.CompanyDescription.Valid {
			resp.CompanyDescription = company.CompanyDescription.String
		}
		if company.CompanyWebsite.Valid {
			resp.CompanyWebsite = company.CompanyWebsite.String
		}
		if company.CompanyIndustry.Valid {
			resp.CompanyIndustry = company.CompanyIndustry.String
		}
		if company.CompanySize.Valid {
			resp.CompanySize = company.CompanySize.String
		}
		if company.CompanyLocation.Valid {
			resp.CompanyLocation = company.CompanyLocation.String
		}
		if company.CompanyPhone.Valid {
			resp.CompanyPhone = company.CompanyPhone.String
		}
		if company.CompanyEmail.Valid {
			resp.CompanyEmail = company.CompanyEmail.String
		}
		if company.CompanyAddress.Valid {
			resp.CompanyAddress = company.CompanyAddress.String
		}
		if company.CompanyCity.Valid {
			resp.CompanyCity = company.CompanyCity.String
		}
		if company.CompanyProvince.Valid {
			resp.CompanyProvince = company.CompanyProvince.String
		}
		if company.CompanyPostalCode.Valid {
			resp.CompanyPostalCode = company.CompanyPostalCode.String
		}
		if company.EstablishedYear.Valid {
			resp.EstablishedYear = company.EstablishedYear.Int64
		}
		if company.EmployeeCount.Valid {
			resp.EmployeeCount = company.EmployeeCount.Int64
		}
		if company.CompanyStatus.Valid {
			resp.CompanyStatus = company.CompanyStatus.String
		}
		if company.KTPFounderURL.Valid {
			resp.KTPFounderURL = company.KTPFounderURL.String
		}
		if company.AktaPendirianURL.Valid {
			resp.AktaPendirianURL = company.AktaPendirianURL.String
		}
		if company.NPWPURL.Valid {
			resp.NPWPURL = company.NPWPURL.String
		}
		if company.NIBURL.Valid {
			resp.NIBURL = company.NIBURL.String
		}
	}

	return resp
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

// PasswordResetToken represents a password reset token entity
type PasswordResetToken struct {
	ID        uint64       `db:"id"`
	Email     string       `db:"email"`
	Token     string       `db:"token"`
	ExpiresAt time.Time    `db:"expires_at"`
	UsedAt    sql.NullTime `db:"used_at"`
	CreatedAt time.Time    `db:"created_at"`
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

// ForgotPasswordRequest represents a forgot password request
type ForgotPasswordRequest struct {
	Email string `json:"email" validate:"required,email"`
}

// ResetPasswordRequest represents a reset password request
type ResetPasswordRequest struct {
	Token       string `json:"token" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=8,password"`
}

// ChangePasswordRequest represents a change password request (for logged-in users)
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=8,password"`
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
	User         interface{} `json:"user"` // Can be UserResponse or UserWithCompanyResponse
	AccessToken  string      `json:"access_token"`
	RefreshToken string      `json:"refresh_token"`
	ExpiresIn    int64       `json:"expires_in"`
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
