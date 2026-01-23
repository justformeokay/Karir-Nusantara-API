package profile

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/karirnusantara/api/internal/shared/hashid"
)

// Gender enum
type Gender string

const (
	GenderMale           Gender = "male"
	GenderFemale         Gender = "female"
	GenderOther          Gender = "other"
	GenderPreferNotToSay Gender = "prefer_not_to_say"
)

// MaritalStatus enum
type MaritalStatus string

const (
	MaritalStatusSingle   MaritalStatus = "single"
	MaritalStatusMarried  MaritalStatus = "married"
	MaritalStatusDivorced MaritalStatus = "divorced"
	MaritalStatusWidowed  MaritalStatus = "widowed"
)

// DocumentType enum
type DocumentType string

const (
	DocumentTypeCVUploaded  DocumentType = "cv_uploaded"
	DocumentTypeCVGenerated DocumentType = "cv_generated"
	DocumentTypeCertificate DocumentType = "certificate"
	DocumentTypeTranscript  DocumentType = "transcript"
	DocumentTypePortfolio   DocumentType = "portfolio"
	DocumentTypeKTP         DocumentType = "ktp"
	DocumentTypeOther       DocumentType = "other"
)

// ApplicantProfile represents a job seeker's comprehensive profile
type ApplicantProfile struct {
	ID     uint64 `db:"id" json:"id"`
	UserID uint64 `db:"user_id" json:"user_id"`

	// Personal Information
	DateOfBirth   sql.NullTime   `db:"date_of_birth" json:"-"`
	Gender        sql.NullString `db:"gender" json:"-"`
	Nationality   sql.NullString `db:"nationality" json:"-"`
	MaritalStatus sql.NullString `db:"marital_status" json:"-"`

	// Identity
	NIK sql.NullString `db:"nik" json:"-"`

	// Address
	Address    sql.NullString `db:"address" json:"-"`
	City       sql.NullString `db:"city" json:"-"`
	Province   sql.NullString `db:"province" json:"-"`
	PostalCode sql.NullString `db:"postal_code" json:"-"`
	Country    sql.NullString `db:"country" json:"-"`

	// Professional Links
	LinkedInURL     sql.NullString `db:"linkedin_url" json:"-"`
	GithubURL       sql.NullString `db:"github_url" json:"-"`
	PortfolioURL    sql.NullString `db:"portfolio_url" json:"-"`
	PersonalWebsite sql.NullString `db:"personal_website" json:"-"`

	// Bio/Summary
	ProfessionalSummary sql.NullString `db:"professional_summary" json:"-"`
	Headline            sql.NullString `db:"headline" json:"-"`

	// Job Preferences
	ExpectedSalaryMin  sql.NullInt64  `db:"expected_salary_min" json:"-"`
	ExpectedSalaryMax  sql.NullInt64  `db:"expected_salary_max" json:"-"`
	PreferredJobTypes  sql.NullString `db:"preferred_job_types" json:"-"`
	PreferredLocations sql.NullString `db:"preferred_locations" json:"-"`
	AvailableFrom      sql.NullTime   `db:"available_from" json:"-"`
	WillingToRelocate  bool           `db:"willing_to_relocate" json:"willing_to_relocate"`

	// Profile Completeness
	ProfileCompleteness int `db:"profile_completeness" json:"profile_completeness"`

	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

// ApplicantDocument represents an uploaded document
type ApplicantDocument struct {
	ID          uint64       `db:"id" json:"id"`
	UserID      uint64       `db:"user_id" json:"user_id"`
	DocType     DocumentType `db:"document_type" json:"document_type"`
	DocName     string       `db:"document_name" json:"document_name"`
	DocURL      string       `db:"document_url" json:"document_url"`
	FileSize    sql.NullInt64    `db:"file_size" json:"-"`
	MimeType    sql.NullString   `db:"mime_type" json:"-"`
	IsPrimary   bool         `db:"is_primary" json:"is_primary"`
	Description sql.NullString   `db:"description" json:"-"`
	UploadedAt  time.Time    `db:"uploaded_at" json:"uploaded_at"`
	ExpiresAt   sql.NullTime `db:"expires_at" json:"-"`
}

// ========================================
// Response DTOs
// ========================================

// ProfileResponse is the response DTO for applicant profile
type ProfileResponse struct {
	ID     uint64 `json:"id"`
	HashID string `json:"hash_id"`
	UserID uint64 `json:"user_id"`

	// Personal Information
	DateOfBirth   *string `json:"date_of_birth,omitempty"`
	Gender        *string `json:"gender,omitempty"`
	Nationality   *string `json:"nationality,omitempty"`
	MaritalStatus *string `json:"marital_status,omitempty"`

	// Identity (masked for privacy)
	NIK *string `json:"nik,omitempty"`

	// Address
	Address    *string `json:"address,omitempty"`
	City       *string `json:"city,omitempty"`
	Province   *string `json:"province,omitempty"`
	PostalCode *string `json:"postal_code,omitempty"`
	Country    *string `json:"country,omitempty"`

	// Professional Links
	LinkedInURL     *string `json:"linkedin_url,omitempty"`
	GithubURL       *string `json:"github_url,omitempty"`
	PortfolioURL    *string `json:"portfolio_url,omitempty"`
	PersonalWebsite *string `json:"personal_website,omitempty"`

	// Bio/Summary
	ProfessionalSummary *string `json:"professional_summary,omitempty"`
	Headline            *string `json:"headline,omitempty"`

	// Job Preferences
	ExpectedSalaryMin  *int64   `json:"expected_salary_min,omitempty"`
	ExpectedSalaryMax  *int64   `json:"expected_salary_max,omitempty"`
	PreferredJobTypes  []string `json:"preferred_job_types,omitempty"`
	PreferredLocations []string `json:"preferred_locations,omitempty"`
	AvailableFrom      *string  `json:"available_from,omitempty"`
	WillingToRelocate  bool     `json:"willing_to_relocate"`

	// Profile Completeness
	ProfileCompleteness int `json:"profile_completeness"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// DocumentResponse is the response DTO for applicant document
type DocumentResponse struct {
	ID          uint64 `json:"id"`
	HashID      string `json:"hash_id"`
	DocType     string `json:"document_type"`
	DocName     string `json:"document_name"`
	DocURL      string `json:"document_url"`
	FileSize    *int64 `json:"file_size,omitempty"`
	MimeType    *string `json:"mime_type,omitempty"`
	IsPrimary   bool   `json:"is_primary"`
	Description *string `json:"description,omitempty"`
	UploadedAt  time.Time `json:"uploaded_at"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty"`
}

// ToResponse converts ApplicantProfile to ProfileResponse
func (p *ApplicantProfile) ToResponse() *ProfileResponse {
	resp := &ProfileResponse{
		ID:                  p.ID,
		HashID:              hashid.Encode(p.ID),
		UserID:              p.UserID,
		WillingToRelocate:   p.WillingToRelocate,
		ProfileCompleteness: p.ProfileCompleteness,
		CreatedAt:           p.CreatedAt,
		UpdatedAt:           p.UpdatedAt,
	}

	// Personal Information
	if p.DateOfBirth.Valid {
		dob := p.DateOfBirth.Time.Format("2006-01-02")
		resp.DateOfBirth = &dob
	}
	if p.Gender.Valid {
		resp.Gender = &p.Gender.String
	}
	if p.Nationality.Valid {
		resp.Nationality = &p.Nationality.String
	}
	if p.MaritalStatus.Valid {
		resp.MaritalStatus = &p.MaritalStatus.String
	}

	// Identity
	if p.NIK.Valid {
		resp.NIK = &p.NIK.String
	}

	// Address
	if p.Address.Valid {
		resp.Address = &p.Address.String
	}
	if p.City.Valid {
		resp.City = &p.City.String
	}
	if p.Province.Valid {
		resp.Province = &p.Province.String
	}
	if p.PostalCode.Valid {
		resp.PostalCode = &p.PostalCode.String
	}
	if p.Country.Valid {
		resp.Country = &p.Country.String
	}

	// Professional Links
	if p.LinkedInURL.Valid {
		resp.LinkedInURL = &p.LinkedInURL.String
	}
	if p.GithubURL.Valid {
		resp.GithubURL = &p.GithubURL.String
	}
	if p.PortfolioURL.Valid {
		resp.PortfolioURL = &p.PortfolioURL.String
	}
	if p.PersonalWebsite.Valid {
		resp.PersonalWebsite = &p.PersonalWebsite.String
	}

	// Bio/Summary
	if p.ProfessionalSummary.Valid {
		resp.ProfessionalSummary = &p.ProfessionalSummary.String
	}
	if p.Headline.Valid {
		resp.Headline = &p.Headline.String
	}

	// Job Preferences
	if p.ExpectedSalaryMin.Valid {
		resp.ExpectedSalaryMin = &p.ExpectedSalaryMin.Int64
	}
	if p.ExpectedSalaryMax.Valid {
		resp.ExpectedSalaryMax = &p.ExpectedSalaryMax.Int64
	}
	if p.PreferredJobTypes.Valid {
		var jobTypes []string
		if err := json.Unmarshal([]byte(p.PreferredJobTypes.String), &jobTypes); err == nil {
			resp.PreferredJobTypes = jobTypes
		}
	}
	if p.PreferredLocations.Valid {
		var locations []string
		if err := json.Unmarshal([]byte(p.PreferredLocations.String), &locations); err == nil {
			resp.PreferredLocations = locations
		}
	}
	if p.AvailableFrom.Valid {
		af := p.AvailableFrom.Time.Format("2006-01-02")
		resp.AvailableFrom = &af
	}

	return resp
}

// ToResponse converts ApplicantDocument to DocumentResponse
func (d *ApplicantDocument) ToResponse() *DocumentResponse {
	resp := &DocumentResponse{
		ID:         d.ID,
		HashID:     hashid.Encode(d.ID),
		DocType:    string(d.DocType),
		DocName:    d.DocName,
		DocURL:     d.DocURL,
		IsPrimary:  d.IsPrimary,
		UploadedAt: d.UploadedAt,
	}

	if d.FileSize.Valid {
		resp.FileSize = &d.FileSize.Int64
	}
	if d.MimeType.Valid {
		resp.MimeType = &d.MimeType.String
	}
	if d.Description.Valid {
		resp.Description = &d.Description.String
	}
	if d.ExpiresAt.Valid {
		resp.ExpiresAt = &d.ExpiresAt.Time
	}

	return resp
}

// ========================================
// Request DTOs
// ========================================

// UpdateProfileRequest is the request DTO for updating profile
type UpdateProfileRequest struct {
	// Personal Information
	DateOfBirth   *string `json:"date_of_birth" validate:"omitempty,datetime=2006-01-02"`
	Gender        *string `json:"gender" validate:"omitempty,oneof=male female other prefer_not_to_say"`
	Nationality   *string `json:"nationality" validate:"omitempty,max=100"`
	MaritalStatus *string `json:"marital_status" validate:"omitempty,oneof=single married divorced widowed"`

	// Identity
	NIK *string `json:"nik" validate:"omitempty,max=20"`

	// Address
	Address    *string `json:"address" validate:"omitempty,max=1000"`
	City       *string `json:"city" validate:"omitempty,max=100"`
	Province   *string `json:"province" validate:"omitempty,max=100"`
	PostalCode *string `json:"postal_code" validate:"omitempty,max=10"`
	Country    *string `json:"country" validate:"omitempty,max=100"`

	// Professional Links
	LinkedInURL     *string `json:"linkedin_url" validate:"omitempty,url,max=500"`
	GithubURL       *string `json:"github_url" validate:"omitempty,url,max=500"`
	PortfolioURL    *string `json:"portfolio_url" validate:"omitempty,url,max=500"`
	PersonalWebsite *string `json:"personal_website" validate:"omitempty,url,max=500"`

	// Bio/Summary
	ProfessionalSummary *string `json:"professional_summary" validate:"omitempty,max=5000"`
	Headline            *string `json:"headline" validate:"omitempty,max=255"`

	// Job Preferences
	ExpectedSalaryMin  *int64   `json:"expected_salary_min" validate:"omitempty,min=0"`
	ExpectedSalaryMax  *int64   `json:"expected_salary_max" validate:"omitempty,min=0,gtefield=ExpectedSalaryMin"`
	PreferredJobTypes  []string `json:"preferred_job_types" validate:"omitempty,dive,oneof=full_time part_time contract internship freelance remote"`
	PreferredLocations []string `json:"preferred_locations" validate:"omitempty,max=10,dive,max=100"`
	AvailableFrom      *string  `json:"available_from" validate:"omitempty,datetime=2006-01-02"`
	WillingToRelocate  *bool    `json:"willing_to_relocate"`
}

// UploadDocumentRequest is the request for document upload (multipart form)
type UploadDocumentRequest struct {
	DocType     string `json:"document_type" validate:"required,oneof=cv_uploaded certificate transcript portfolio ktp other"`
	Description string `json:"description" validate:"omitempty,max=500"`
	IsPrimary   bool   `json:"is_primary"`
}
