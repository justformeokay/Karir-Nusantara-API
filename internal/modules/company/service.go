package company

import (
	"context"
	"database/sql"

	apperrors "github.com/karirnusantara/api/internal/shared/errors"
)

// Service defines the company service interface
type Service interface {
	GetCompanyByUserID(ctx context.Context, userID uint64) (*CompanyResponse, error)
	GetCompanyEntityByUserID(ctx context.Context, userID uint64) (*Company, error)
	GetCompanyIDByUserID(ctx context.Context, userID uint64) (uint64, error)
	GetPublicCompanyByID(ctx context.Context, companyID uint64) (*PublicCompanyResponse, error)
	CreateOrUpdateCompany(ctx context.Context, userID uint64, req *UpdateCompanyRequest) (*CompanyResponse, error)
	UpdateCompanyLogoURL(ctx context.Context, companyID uint64, logoURL string) error
	UpdateCompanyDocument(ctx context.Context, companyID uint64, docType, filePath string) error
}

// service implements Service
type service struct {
	repo Repository
}

// NewService creates a new company service
func NewService(repo Repository) Service {
	return &service{
		repo: repo,
	}
}

// GetCompanyByUserID retrieves company information for a user
func (s *service) GetCompanyByUserID(ctx context.Context, userID uint64) (*CompanyResponse, error) {
	company, err := s.repo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, apperrors.NewInternalError("Failed to get company", err)
	}
	if company == nil {
		return nil, nil
	}

	return company.ToResponse(), nil
}

// GetCompanyEntityByUserID retrieves the actual company entity for a user
func (s *service) GetCompanyEntityByUserID(ctx context.Context, userID uint64) (*Company, error) {
	company, err := s.repo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, apperrors.NewInternalError("Failed to get company", err)
	}
	if company == nil {
		return nil, nil
	}

	return company, nil
}

// GetCompanyIDByUserID retrieves only the company ID for a user (optimized for quota lookups)
func (s *service) GetCompanyIDByUserID(ctx context.Context, userID uint64) (uint64, error) {
	company, err := s.repo.GetByUserID(ctx, userID)
	if err != nil {
		return 0, apperrors.NewInternalError("Failed to get company", err)
	}
	if company == nil {
		return 0, nil
	}

	return company.ID, nil
}

// GetPublicCompanyByID retrieves public company information by ID (for job seekers)
func (s *service) GetPublicCompanyByID(ctx context.Context, companyID uint64) (*PublicCompanyResponse, error) {
	company, err := s.repo.GetByID(ctx, companyID)
	if err != nil {
		return nil, apperrors.NewInternalError("Failed to get company", err)
	}
	if company == nil {
		return nil, nil
	}

	// Only return companies with verified status
	if company.CompanyStatus != "verified" && company.CompanyStatus != "pending" {
		return nil, nil
	}

	return company.ToPublicResponse(), nil
}

// CreateOrUpdateCompany creates or updates company information
func (s *service) CreateOrUpdateCompany(ctx context.Context, userID uint64, req *UpdateCompanyRequest) (*CompanyResponse, error) {
	// Check if company exists
	company, err := s.repo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, apperrors.NewInternalError("Failed to get company", err)
	}

	// Track if this was a verified company (for re-verification)
	wasVerified := false
	if company != nil && company.CompanyStatus == "verified" {
		wasVerified = true
	}

	if company == nil {
		// Create new company
		company = &Company{
			UserID:        userID,
			CompanyStatus: "pending",
		}
	}

	// Update fields if provided
	if req.CompanyName != "" {
		company.CompanyName = sql.NullString{String: req.CompanyName, Valid: true}
	}
	if req.CompanyDescription != "" {
		company.CompanyDescription = sql.NullString{String: req.CompanyDescription, Valid: true}
	}
	if req.CompanyWebsite != "" {
		company.CompanyWebsite = sql.NullString{String: req.CompanyWebsite, Valid: true}
	}
	if req.CompanyIndustry != "" {
		company.CompanyIndustry = sql.NullString{String: req.CompanyIndustry, Valid: true}
	}
	if req.CompanySize != "" {
		company.CompanySize = sql.NullString{String: req.CompanySize, Valid: true}
	}
	if req.CompanyLocation != "" {
		company.CompanyLocation = sql.NullString{String: req.CompanyLocation, Valid: true}
	}
	if req.CompanyPhone != "" {
		company.CompanyPhone = sql.NullString{String: req.CompanyPhone, Valid: true}
	}
	if req.CompanyEmail != "" {
		company.CompanyEmail = sql.NullString{String: req.CompanyEmail, Valid: true}
	}
	if req.CompanyAddress != "" {
		company.CompanyAddress = sql.NullString{String: req.CompanyAddress, Valid: true}
	}
	if req.CompanyCity != "" {
		company.CompanyCity = sql.NullString{String: req.CompanyCity, Valid: true}
	}
	if req.CompanyProvince != "" {
		company.CompanyProvince = sql.NullString{String: req.CompanyProvince, Valid: true}
	}
	if req.CompanyPostalCode != "" {
		company.CompanyPostalCode = sql.NullString{String: req.CompanyPostalCode, Valid: true}
	}
	if req.EstablishedYear > 0 {
		company.EstablishedYear = sql.NullInt64{Int64: int64(req.EstablishedYear), Valid: true}
	}
	if req.EmployeeCount > 0 {
		company.EmployeeCount = sql.NullInt64{Int64: int64(req.EmployeeCount), Valid: true}
	}

	// If company was verified, reset status to pending for re-verification
	// This ensures data integrity and prevents manipulation of company information
	if wasVerified {
		company.CompanyStatus = "pending"
		// Clear previous verification data
		company.DocumentsVerifiedAt = sql.NullTime{Valid: false}
		company.DocumentsVerifiedBy = sql.NullInt64{Valid: false}
		company.VerificationNotes = sql.NullString{Valid: false}
	}

	// Create or update
	if company.ID == 0 {
		if err := s.repo.Create(ctx, company); err != nil {
			return nil, apperrors.NewInternalError("Failed to create company", err)
		}
	} else {
		if err := s.repo.Update(ctx, company); err != nil {
			return nil, apperrors.NewInternalError("Failed to update company", err)
		}
	}

	return company.ToResponse(), nil
}

// UpdateCompanyLogoURL updates the company logo URL in database
func (s *service) UpdateCompanyLogoURL(ctx context.Context, companyID uint64, logoURL string) error {
	company, err := s.repo.GetByID(ctx, companyID)
	if err != nil {
		return apperrors.NewInternalError("Failed to get company", err)
	}
	if company == nil {
		return apperrors.NewNotFoundError("Company")
	}

	company.CompanyLogoURL = sql.NullString{String: logoURL, Valid: true}
	if err := s.repo.Update(ctx, company); err != nil {
		return apperrors.NewInternalError("Failed to update company logo", err)
	}

	return nil
}

// UpdateCompanyDocument updates a company document URL in database
func (s *service) UpdateCompanyDocument(ctx context.Context, companyID uint64, docType, filePath string) error {
	company, err := s.repo.GetByID(ctx, companyID)
	if err != nil {
		return apperrors.NewInternalError("Failed to get company", err)
	}
	if company == nil {
		return apperrors.NewNotFoundError("Company")
	}

	switch docType {
	case "ktp":
		company.KTPFounderURL = sql.NullString{String: filePath, Valid: true}
	case "akta":
		company.AktaPendirianURL = sql.NullString{String: filePath, Valid: true}
	case "npwp":
		company.NPWPURL = sql.NullString{String: filePath, Valid: true}
	case "nib":
		company.NIBURL = sql.NullString{String: filePath, Valid: true}
	}

	if err := s.repo.Update(ctx, company); err != nil {
		return apperrors.NewInternalError("Failed to update company document", err)
	}

	return nil
}

// UpdateCompanyRequest represents a request to update company information
type UpdateCompanyRequest struct {
	CompanyName        string `json:"company_name,omitempty"`
	CompanyDescription string `json:"company_description,omitempty"`
	CompanyWebsite     string `json:"company_website,omitempty"`
	CompanyIndustry    string `json:"company_industry,omitempty"`
	CompanySize        string `json:"company_size,omitempty"`
	CompanyLocation    string `json:"company_location,omitempty"`
	CompanyPhone       string `json:"company_phone,omitempty"`
	CompanyEmail       string `json:"company_email,omitempty"`
	CompanyAddress     string `json:"company_address,omitempty"`
	CompanyCity        string `json:"company_city,omitempty"`
	CompanyProvince    string `json:"company_province,omitempty"`
	CompanyPostalCode  string `json:"company_postal_code,omitempty"`
	EstablishedYear    int    `json:"established_year,omitempty"`
	EmployeeCount      int    `json:"employee_count,omitempty"`
}
