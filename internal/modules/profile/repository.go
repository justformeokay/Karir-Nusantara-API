package profile

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
)

// Repository defines the profile repository interface
type Repository interface {
	// Profile operations
	CreateProfile(ctx context.Context, profile *ApplicantProfile) error
	GetProfileByUserID(ctx context.Context, userID uint64) (*ApplicantProfile, error)
	GetProfileByID(ctx context.Context, id uint64) (*ApplicantProfile, error)
	UpdateProfile(ctx context.Context, profile *ApplicantProfile) error
	DeleteProfile(ctx context.Context, userID uint64) error

	// Avatar operations
	UpdateUserAvatar(ctx context.Context, userID uint64, avatarURL string) error

	// Document operations
	CreateDocument(ctx context.Context, doc *ApplicantDocument) error
	GetDocumentByID(ctx context.Context, id uint64) (*ApplicantDocument, error)
	GetDocumentsByUserID(ctx context.Context, userID uint64) ([]*ApplicantDocument, error)
	GetPrimaryDocument(ctx context.Context, userID uint64, docType DocumentType) (*ApplicantDocument, error)
	UpdateDocument(ctx context.Context, doc *ApplicantDocument) error
	DeleteDocument(ctx context.Context, id uint64) error
	SetPrimaryDocument(ctx context.Context, userID uint64, docID uint64, docType DocumentType) error
}

type mysqlRepository struct {
	db *sqlx.DB
}

// NewRepository creates a new profile repository
func NewRepository(db *sqlx.DB) Repository {
	return &mysqlRepository{db: db}
}

// ========================================
// Profile Operations
// ========================================

// CreateProfile creates a new applicant profile
func (r *mysqlRepository) CreateProfile(ctx context.Context, profile *ApplicantProfile) error {
	query := `
		INSERT INTO applicant_profiles (
			user_id, date_of_birth, gender, nationality, marital_status,
			nik, address, city, province, postal_code, country,
			linkedin_url, github_url, portfolio_url, personal_website,
			professional_summary, headline,
			expected_salary_min, expected_salary_max, preferred_job_types, preferred_locations,
			available_from, willing_to_relocate, profile_completeness,
			created_at, updated_at
		) VALUES (
			?, ?, ?, ?, ?,
			?, ?, ?, ?, ?, ?,
			?, ?, ?, ?,
			?, ?,
			?, ?, ?, ?,
			?, ?, ?,
			NOW(), NOW()
		)
	`

	result, err := r.db.ExecContext(ctx, query,
		profile.UserID, profile.DateOfBirth, profile.Gender, profile.Nationality, profile.MaritalStatus,
		profile.NIK, profile.Address, profile.City, profile.Province, profile.PostalCode, profile.Country,
		profile.LinkedInURL, profile.GithubURL, profile.PortfolioURL, profile.PersonalWebsite,
		profile.ProfessionalSummary, profile.Headline,
		profile.ExpectedSalaryMin, profile.ExpectedSalaryMax, profile.PreferredJobTypes, profile.PreferredLocations,
		profile.AvailableFrom, profile.WillingToRelocate, profile.ProfileCompleteness,
	)
	if err != nil {
		return fmt.Errorf("failed to create profile: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	profile.ID = uint64(id)
	return nil
}

// GetProfileByUserID retrieves a profile by user ID
func (r *mysqlRepository) GetProfileByUserID(ctx context.Context, userID uint64) (*ApplicantProfile, error) {
	query := `
		SELECT id, user_id, date_of_birth, gender, nationality, marital_status,
			   nik, address, city, province, postal_code, country,
			   linkedin_url, github_url, portfolio_url, personal_website,
			   professional_summary, headline,
			   expected_salary_min, expected_salary_max, preferred_job_types, preferred_locations,
			   available_from, willing_to_relocate, profile_completeness,
			   created_at, updated_at
		FROM applicant_profiles
		WHERE user_id = ?
	`

	var profile ApplicantProfile
	if err := r.db.GetContext(ctx, &profile, query, userID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get profile by user id: %w", err)
	}

	return &profile, nil
}

// GetProfileByID retrieves a profile by ID
func (r *mysqlRepository) GetProfileByID(ctx context.Context, id uint64) (*ApplicantProfile, error) {
	query := `
		SELECT id, user_id, date_of_birth, gender, nationality, marital_status,
			   nik, address, city, province, postal_code, country,
			   linkedin_url, github_url, portfolio_url, personal_website,
			   professional_summary, headline,
			   expected_salary_min, expected_salary_max, preferred_job_types, preferred_locations,
			   available_from, willing_to_relocate, profile_completeness,
			   created_at, updated_at
		FROM applicant_profiles
		WHERE id = ?
	`

	var profile ApplicantProfile
	if err := r.db.GetContext(ctx, &profile, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get profile by id: %w", err)
	}

	return &profile, nil
}

// UpdateProfile updates an existing profile
func (r *mysqlRepository) UpdateProfile(ctx context.Context, profile *ApplicantProfile) error {
	query := `
		UPDATE applicant_profiles SET
			date_of_birth = ?, gender = ?, nationality = ?, marital_status = ?,
			nik = ?, address = ?, city = ?, province = ?, postal_code = ?, country = ?,
			linkedin_url = ?, github_url = ?, portfolio_url = ?, personal_website = ?,
			professional_summary = ?, headline = ?,
			expected_salary_min = ?, expected_salary_max = ?, preferred_job_types = ?, preferred_locations = ?,
			available_from = ?, willing_to_relocate = ?, profile_completeness = ?,
			updated_at = NOW()
		WHERE id = ?
	`

	_, err := r.db.ExecContext(ctx, query,
		profile.DateOfBirth, profile.Gender, profile.Nationality, profile.MaritalStatus,
		profile.NIK, profile.Address, profile.City, profile.Province, profile.PostalCode, profile.Country,
		profile.LinkedInURL, profile.GithubURL, profile.PortfolioURL, profile.PersonalWebsite,
		profile.ProfessionalSummary, profile.Headline,
		profile.ExpectedSalaryMin, profile.ExpectedSalaryMax, profile.PreferredJobTypes, profile.PreferredLocations,
		profile.AvailableFrom, profile.WillingToRelocate, profile.ProfileCompleteness,
		profile.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update profile: %w", err)
	}

	return nil
}

// DeleteProfile deletes a profile by user ID
func (r *mysqlRepository) DeleteProfile(ctx context.Context, userID uint64) error {
	query := `DELETE FROM applicant_profiles WHERE user_id = ?`
	_, err := r.db.ExecContext(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("failed to delete profile: %w", err)
	}
	return nil
}

// UpdateUserAvatar updates user's avatar URL
func (r *mysqlRepository) UpdateUserAvatar(ctx context.Context, userID uint64, avatarURL string) error {
	query := `UPDATE users SET avatar_url = ?, updated_at = NOW() WHERE id = ? AND deleted_at IS NULL`
	_, err := r.db.ExecContext(ctx, query, avatarURL, userID)
	if err != nil {
		return fmt.Errorf("failed to update user avatar: %w", err)
	}
	return nil
}

// ========================================
// Document Operations
// ========================================

// CreateDocument creates a new document
func (r *mysqlRepository) CreateDocument(ctx context.Context, doc *ApplicantDocument) error {
	query := `
		INSERT INTO applicant_documents (
			user_id, document_type, document_name, document_url,
			file_size, mime_type, is_primary, description,
			uploaded_at
		) VALUES (
			?, ?, ?, ?,
			?, ?, ?, ?,
			NOW()
		)
	`

	result, err := r.db.ExecContext(ctx, query,
		doc.UserID, doc.DocType, doc.DocName, doc.DocURL,
		doc.FileSize, doc.MimeType, doc.IsPrimary, doc.Description,
	)
	if err != nil {
		return fmt.Errorf("failed to create document: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	doc.ID = uint64(id)
	return nil
}

// GetDocumentByID retrieves a document by ID
func (r *mysqlRepository) GetDocumentByID(ctx context.Context, id uint64) (*ApplicantDocument, error) {
	query := `
		SELECT id, user_id, document_type, document_name, document_url,
			   file_size, mime_type, is_primary, description,
			   uploaded_at, expires_at
		FROM applicant_documents
		WHERE id = ?
	`

	var doc ApplicantDocument
	if err := r.db.GetContext(ctx, &doc, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get document by id: %w", err)
	}

	return &doc, nil
}

// GetDocumentsByUserID retrieves all documents for a user
func (r *mysqlRepository) GetDocumentsByUserID(ctx context.Context, userID uint64) ([]*ApplicantDocument, error) {
	query := `
		SELECT id, user_id, document_type, document_name, document_url,
			   file_size, mime_type, is_primary, description,
			   uploaded_at, expires_at
		FROM applicant_documents
		WHERE user_id = ?
		ORDER BY is_primary DESC, uploaded_at DESC
	`

	var docs []*ApplicantDocument
	if err := r.db.SelectContext(ctx, &docs, query, userID); err != nil {
		return nil, fmt.Errorf("failed to get documents by user id: %w", err)
	}

	return docs, nil
}

// GetPrimaryDocument retrieves the primary document of a specific type for a user
func (r *mysqlRepository) GetPrimaryDocument(ctx context.Context, userID uint64, docType DocumentType) (*ApplicantDocument, error) {
	query := `
		SELECT id, user_id, document_type, document_name, document_url,
			   file_size, mime_type, is_primary, description,
			   uploaded_at, expires_at
		FROM applicant_documents
		WHERE user_id = ? AND document_type = ? AND is_primary = 1
		LIMIT 1
	`

	var doc ApplicantDocument
	if err := r.db.GetContext(ctx, &doc, query, userID, docType); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get primary document: %w", err)
	}

	return &doc, nil
}

// UpdateDocument updates an existing document
func (r *mysqlRepository) UpdateDocument(ctx context.Context, doc *ApplicantDocument) error {
	query := `
		UPDATE applicant_documents SET
			document_name = ?, is_primary = ?, description = ?
		WHERE id = ?
	`

	_, err := r.db.ExecContext(ctx, query,
		doc.DocName, doc.IsPrimary, doc.Description,
		doc.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update document: %w", err)
	}

	return nil
}

// DeleteDocument deletes a document by ID
func (r *mysqlRepository) DeleteDocument(ctx context.Context, id uint64) error {
	query := `DELETE FROM applicant_documents WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete document: %w", err)
	}
	return nil
}

// SetPrimaryDocument sets a document as primary and unsets others of the same type
func (r *mysqlRepository) SetPrimaryDocument(ctx context.Context, userID uint64, docID uint64, docType DocumentType) error {
	// Start transaction
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback()

	// Unset all primary documents of this type for the user
	unsetQuery := `UPDATE applicant_documents SET is_primary = 0 WHERE user_id = ? AND document_type = ?`
	if _, err := tx.ExecContext(ctx, unsetQuery, userID, docType); err != nil {
		return fmt.Errorf("failed to unset primary documents: %w", err)
	}

	// Set the specified document as primary
	setQuery := `UPDATE applicant_documents SET is_primary = 1 WHERE id = ? AND user_id = ?`
	if _, err := tx.ExecContext(ctx, setQuery, docID, userID); err != nil {
		return fmt.Errorf("failed to set primary document: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
