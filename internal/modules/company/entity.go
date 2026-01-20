package company

import (
	"database/sql"
	"time"
)

// Company represents a company entity
type Company struct {
	ID                    uint64         `db:"id" json:"id"`
	UserID                uint64         `db:"user_id" json:"user_id"`
	CompanyName           sql.NullString `db:"company_name" json:"company_name,omitempty"`
	CompanyDescription    sql.NullString `db:"company_description" json:"company_description,omitempty"`
	CompanyWebsite        sql.NullString `db:"company_website" json:"company_website,omitempty"`
	CompanyLogoURL        sql.NullString `db:"company_logo_url" json:"company_logo_url,omitempty"`
	CompanyIndustry       sql.NullString `db:"company_industry" json:"company_industry,omitempty"`
	CompanySize           sql.NullString `db:"company_size" json:"company_size,omitempty"`
	CompanyLocation       sql.NullString `db:"company_location" json:"company_location,omitempty"`
	CompanyPhone          sql.NullString `db:"company_phone" json:"company_phone,omitempty"`
	CompanyEmail          sql.NullString `db:"company_email" json:"company_email,omitempty"`
	CompanyAddress        sql.NullString `db:"company_address" json:"company_address,omitempty"`
	CompanyCity           sql.NullString `db:"company_city" json:"company_city,omitempty"`
	CompanyProvince       sql.NullString `db:"company_province" json:"company_province,omitempty"`
	CompanyPostalCode     sql.NullString `db:"company_postal_code" json:"company_postal_code,omitempty"`
	EstablishedYear       sql.NullInt64  `db:"established_year" json:"established_year,omitempty"`
	EmployeeCount         sql.NullInt64  `db:"employee_count" json:"employee_count,omitempty"`
	CompanyStatus         string         `db:"company_status" json:"company_status"`
	KTPFounderURL         sql.NullString `db:"ktp_founder_url" json:"ktp_founder_url,omitempty"`
	AktaPendirianURL      sql.NullString `db:"akta_pendirian_url" json:"akta_pendirian_url,omitempty"`
	NPWPURL               sql.NullString `db:"npwp_url" json:"npwp_url,omitempty"`
	NIBURL                sql.NullString `db:"nib_url" json:"nib_url,omitempty"`
	DocumentsVerifiedAt   sql.NullTime   `db:"documents_verified_at" json:"documents_verified_at,omitempty"`
	DocumentsVerifiedBy   sql.NullInt64  `db:"documents_verified_by" json:"documents_verified_by,omitempty"`
	VerificationNotes     sql.NullString `db:"verification_notes" json:"verification_notes,omitempty"`
	CreatedAt             time.Time      `db:"created_at" json:"created_at"`
	UpdatedAt             time.Time      `db:"updated_at" json:"updated_at"`
	DeletedAt             sql.NullTime   `db:"deleted_at" json:"-"`
}

// CompanyResponse represents company data for API responses
type CompanyResponse struct {
	ID                  uint64 `json:"id"`
	UserID              uint64 `json:"user_id"`
	CompanyName         string `json:"company_name,omitempty"`
	CompanyDescription  string `json:"company_description,omitempty"`
	CompanyWebsite      string `json:"company_website,omitempty"`
	CompanyLogoURL      string `json:"company_logo_url,omitempty"`
	CompanyIndustry     string `json:"company_industry,omitempty"`
	CompanySize         string `json:"company_size,omitempty"`
	CompanyLocation     string `json:"company_location,omitempty"`
	CompanyPhone        string `json:"company_phone,omitempty"`
	CompanyEmail        string `json:"company_email,omitempty"`
	CompanyAddress      string `json:"company_address,omitempty"`
	CompanyCity         string `json:"company_city,omitempty"`
	CompanyProvince     string `json:"company_province,omitempty"`
	CompanyPostalCode   string `json:"company_postal_code,omitempty"`
	EstablishedYear     int    `json:"established_year,omitempty"`
	EmployeeCount       int    `json:"employee_count,omitempty"`
	CompanyStatus       string `json:"company_status"`
	CreatedAt           string `json:"created_at"`
}

// ToResponse converts Company to CompanyResponse
func (c *Company) ToResponse() *CompanyResponse {
	resp := &CompanyResponse{
		ID:            c.ID,
		UserID:        c.UserID,
		CompanyStatus: c.CompanyStatus,
		CreatedAt:     c.CreatedAt.Format(time.RFC3339),
	}

	if c.CompanyName.Valid {
		resp.CompanyName = c.CompanyName.String
	}
	if c.CompanyDescription.Valid {
		resp.CompanyDescription = c.CompanyDescription.String
	}
	if c.CompanyWebsite.Valid {
		resp.CompanyWebsite = c.CompanyWebsite.String
	}
	if c.CompanyLogoURL.Valid {
		resp.CompanyLogoURL = c.CompanyLogoURL.String
	}
	if c.CompanyIndustry.Valid {
		resp.CompanyIndustry = c.CompanyIndustry.String
	}
	if c.CompanySize.Valid {
		resp.CompanySize = c.CompanySize.String
	}
	if c.CompanyLocation.Valid {
		resp.CompanyLocation = c.CompanyLocation.String
	}
	if c.CompanyPhone.Valid {
		resp.CompanyPhone = c.CompanyPhone.String
	}
	if c.CompanyEmail.Valid {
		resp.CompanyEmail = c.CompanyEmail.String
	}
	if c.CompanyAddress.Valid {
		resp.CompanyAddress = c.CompanyAddress.String
	}
	if c.CompanyCity.Valid {
		resp.CompanyCity = c.CompanyCity.String
	}
	if c.CompanyProvince.Valid {
		resp.CompanyProvince = c.CompanyProvince.String
	}
	if c.CompanyPostalCode.Valid {
		resp.CompanyPostalCode = c.CompanyPostalCode.String
	}
	if c.EstablishedYear.Valid {
		resp.EstablishedYear = int(c.EstablishedYear.Int64)
	}
	if c.EmployeeCount.Valid {
		resp.EmployeeCount = int(c.EmployeeCount.Int64)
	}

	return resp
}
