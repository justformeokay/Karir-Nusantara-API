package company

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
)

// Repository defines the company repository interface
type Repository interface {
	GetByUserID(ctx context.Context, userID uint64) (*Company, error)
	GetByID(ctx context.Context, companyID uint64) (*Company, error)
	Create(ctx context.Context, company *Company) error
	Update(ctx context.Context, company *Company) error
}

// mysqlRepository implements Repository
type mysqlRepository struct {
	db *sqlx.DB
}

// NewRepository creates a new company repository
func NewRepository(db *sqlx.DB) Repository {
	return &mysqlRepository{db: db}
}

// GetByUserID retrieves company by user ID
func (r *mysqlRepository) GetByUserID(ctx context.Context, userID uint64) (*Company, error) {
	query := `
		SELECT id, user_id, company_name, company_description, company_website, company_logo_url,
		       company_industry, company_size, company_location, company_phone, company_email,
		       company_address, company_city, company_province, company_postal_code,
		       established_year, employee_count, company_status,
		       ktp_founder_url, akta_pendirian_url, npwp_url, nib_url,
		       documents_verified_at, documents_verified_by, verification_notes,
		       created_at, updated_at, deleted_at
		FROM companies
		WHERE user_id = ? AND deleted_at IS NULL
	`

	var company Company
	if err := r.db.GetContext(ctx, &company, query, userID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get company by user id: %w", err)
	}

	return &company, nil
}

// GetByID retrieves company by company ID
func (r *mysqlRepository) GetByID(ctx context.Context, companyID uint64) (*Company, error) {
	query := `
		SELECT id, user_id, company_name, company_description, company_website, company_logo_url,
		       company_industry, company_size, company_location, company_phone, company_email,
		       company_address, company_city, company_province, company_postal_code,
		       established_year, employee_count, company_status,
		       ktp_founder_url, akta_pendirian_url, npwp_url, nib_url,
		       documents_verified_at, documents_verified_by, verification_notes,
		       created_at, updated_at, deleted_at
		FROM companies
		WHERE id = ? AND deleted_at IS NULL
	`

	var company Company
	if err := r.db.GetContext(ctx, &company, query, companyID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get company by id: %w", err)
	}

	return &company, nil
}

// Create creates a new company record
func (r *mysqlRepository) Create(ctx context.Context, company *Company) error {
	query := `
		INSERT INTO companies (
			user_id, company_name, company_description, company_website, company_logo_url,
			company_industry, company_size, company_location, company_phone, company_email,
			company_address, company_city, company_province, company_postal_code,
			established_year, employee_count, company_status,
			created_at, updated_at
		) VALUES (
			?, ?, ?, ?, ?,
			?, ?, ?, ?, ?,
			?, ?, ?, ?,
			?, ?, ?,
			NOW(), NOW()
		)
	`

	result, err := r.db.ExecContext(ctx, query,
		company.UserID, company.CompanyName, company.CompanyDescription, company.CompanyWebsite, company.CompanyLogoURL,
		company.CompanyIndustry, company.CompanySize, company.CompanyLocation, company.CompanyPhone, company.CompanyEmail,
		company.CompanyAddress, company.CompanyCity, company.CompanyProvince, company.CompanyPostalCode,
		company.EstablishedYear, company.EmployeeCount, company.CompanyStatus,
	)
	if err != nil {
		return fmt.Errorf("failed to create company: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	company.ID = uint64(id)
	return nil
}

// Update updates an existing company record
func (r *mysqlRepository) Update(ctx context.Context, company *Company) error {
	query := `
		UPDATE companies SET
			company_name = ?,
			company_description = ?,
			company_website = ?,
			company_logo_url = ?,
			company_industry = ?,
			company_size = ?,
			company_location = ?,
			company_phone = ?,
			company_email = ?,
			company_address = ?,
			company_city = ?,
			company_province = ?,
			company_postal_code = ?,
			established_year = ?,
			employee_count = ?,
			ktp_founder_url = ?,
			akta_pendirian_url = ?,
			npwp_url = ?,
			nib_url = ?,
			updated_at = NOW()
		WHERE id = ? AND deleted_at IS NULL
	`

	_, err := r.db.ExecContext(ctx, query,
		company.CompanyName, company.CompanyDescription, company.CompanyWebsite, company.CompanyLogoURL,
		company.CompanyIndustry, company.CompanySize, company.CompanyLocation, company.CompanyPhone, company.CompanyEmail,
		company.CompanyAddress, company.CompanyCity, company.CompanyProvince, company.CompanyPostalCode,
		company.EstablishedYear, company.EmployeeCount,
		company.KTPFounderURL, company.AktaPendirianURL, company.NPWPURL, company.NIBURL,
		company.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update company: %w", err)
	}

	return nil
}
