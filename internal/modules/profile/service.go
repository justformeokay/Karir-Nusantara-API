package profile

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	apperrors "github.com/karirnusantara/api/internal/shared/errors"
)

// Service defines the profile service interface
type Service interface {
	// Profile operations
	GetProfile(ctx context.Context, userID uint64) (*ProfileResponse, error)
	CreateOrUpdateProfile(ctx context.Context, userID uint64, req *UpdateProfileRequest) (*ProfileResponse, error)
	DeleteProfile(ctx context.Context, userID uint64) error

	// Avatar operations
	UpdateAvatar(ctx context.Context, userID uint64, avatarURL string) error

	// Document operations
	GetDocuments(ctx context.Context, userID uint64) ([]*DocumentResponse, error)
	GetDocumentByID(ctx context.Context, docID uint64, userID uint64) (*DocumentResponse, error)
	CreateDocument(ctx context.Context, userID uint64, doc *ApplicantDocument) (*DocumentResponse, error)
	UpdateDocument(ctx context.Context, docID uint64, userID uint64, name string, description string) (*DocumentResponse, error)
	DeleteDocument(ctx context.Context, docID uint64, userID uint64) error
	SetPrimaryDocument(ctx context.Context, docID uint64, userID uint64) error
}

type service struct {
	repo Repository
}

// NewService creates a new profile service
func NewService(repo Repository) Service {
	return &service{repo: repo}
}

// ========================================
// Profile Operations
// ========================================

// GetProfile retrieves the user's profile
func (s *service) GetProfile(ctx context.Context, userID uint64) (*ProfileResponse, error) {
	profile, err := s.repo.GetProfileByUserID(ctx, userID)
	if err != nil {
		return nil, apperrors.NewInternalError("Failed to get profile", err)
	}

	// If no profile exists, return an empty profile response
	if profile == nil {
		return &ProfileResponse{
			UserID:              userID,
			ProfileCompleteness: 0,
		}, nil
	}

	return profile.ToResponse(), nil
}

// CreateOrUpdateProfile creates or updates a user's profile
func (s *service) CreateOrUpdateProfile(ctx context.Context, userID uint64, req *UpdateProfileRequest) (*ProfileResponse, error) {
	// Check if profile exists
	existingProfile, err := s.repo.GetProfileByUserID(ctx, userID)
	if err != nil {
		return nil, apperrors.NewInternalError("Failed to check existing profile", err)
	}

	profile := &ApplicantProfile{
		UserID: userID,
	}

	if existingProfile != nil {
		profile.ID = existingProfile.ID
	}

	// Apply request data to profile
	s.applyRequestToProfile(profile, req)

	// Calculate profile completeness
	profile.ProfileCompleteness = s.calculateCompleteness(profile)

	if existingProfile != nil {
		// Update existing profile
		if err := s.repo.UpdateProfile(ctx, profile); err != nil {
			return nil, apperrors.NewInternalError("Failed to update profile", err)
		}
	} else {
		// Create new profile
		if err := s.repo.CreateProfile(ctx, profile); err != nil {
			return nil, apperrors.NewInternalError("Failed to create profile", err)
		}
	}

	// Fetch updated profile
	updatedProfile, err := s.repo.GetProfileByUserID(ctx, userID)
	if err != nil {
		return nil, apperrors.NewInternalError("Failed to get updated profile", err)
	}

	return updatedProfile.ToResponse(), nil
}

// DeleteProfile deletes the user's profile
func (s *service) DeleteProfile(ctx context.Context, userID uint64) error {
	profile, err := s.repo.GetProfileByUserID(ctx, userID)
	if err != nil {
		return apperrors.NewInternalError("Failed to get profile", err)
	}
	if profile == nil {
		return apperrors.NewNotFoundError("Profile")
	}

	if err := s.repo.DeleteProfile(ctx, userID); err != nil {
		return apperrors.NewInternalError("Failed to delete profile", err)
	}

	return nil
}

// UpdateAvatar updates the user's avatar URL
func (s *service) UpdateAvatar(ctx context.Context, userID uint64, avatarURL string) error {
	if err := s.repo.UpdateUserAvatar(ctx, userID, avatarURL); err != nil {
		return apperrors.NewInternalError("Failed to update avatar", err)
	}
	return nil
}

// applyRequestToProfile applies request data to profile entity
func (s *service) applyRequestToProfile(profile *ApplicantProfile, req *UpdateProfileRequest) {
	// Personal Information
	if req.DateOfBirth != nil {
		if t, err := time.Parse("2006-01-02", *req.DateOfBirth); err == nil {
			profile.DateOfBirth = sql.NullTime{Time: t, Valid: true}
		}
	}
	if req.Gender != nil {
		profile.Gender = sql.NullString{String: *req.Gender, Valid: true}
	}
	if req.Nationality != nil {
		profile.Nationality = sql.NullString{String: *req.Nationality, Valid: true}
	}
	if req.MaritalStatus != nil {
		profile.MaritalStatus = sql.NullString{String: *req.MaritalStatus, Valid: true}
	}

	// Identity
	if req.NIK != nil {
		profile.NIK = sql.NullString{String: *req.NIK, Valid: true}
	}

	// Address
	if req.Address != nil {
		profile.Address = sql.NullString{String: *req.Address, Valid: true}
	}
	if req.City != nil {
		profile.City = sql.NullString{String: *req.City, Valid: true}
	}
	if req.Province != nil {
		profile.Province = sql.NullString{String: *req.Province, Valid: true}
	}
	if req.PostalCode != nil {
		profile.PostalCode = sql.NullString{String: *req.PostalCode, Valid: true}
	}
	if req.Country != nil {
		profile.Country = sql.NullString{String: *req.Country, Valid: true}
	}

	// Professional Links
	if req.LinkedInURL != nil {
		profile.LinkedInURL = sql.NullString{String: *req.LinkedInURL, Valid: true}
	}
	if req.GithubURL != nil {
		profile.GithubURL = sql.NullString{String: *req.GithubURL, Valid: true}
	}
	if req.PortfolioURL != nil {
		profile.PortfolioURL = sql.NullString{String: *req.PortfolioURL, Valid: true}
	}
	if req.PersonalWebsite != nil {
		profile.PersonalWebsite = sql.NullString{String: *req.PersonalWebsite, Valid: true}
	}

	// Bio/Summary
	if req.ProfessionalSummary != nil {
		profile.ProfessionalSummary = sql.NullString{String: *req.ProfessionalSummary, Valid: true}
	}
	if req.Headline != nil {
		profile.Headline = sql.NullString{String: *req.Headline, Valid: true}
	}

	// Job Preferences
	if req.ExpectedSalaryMin != nil {
		profile.ExpectedSalaryMin = sql.NullInt64{Int64: *req.ExpectedSalaryMin, Valid: true}
	}
	if req.ExpectedSalaryMax != nil {
		profile.ExpectedSalaryMax = sql.NullInt64{Int64: *req.ExpectedSalaryMax, Valid: true}
	}
	if req.PreferredJobTypes != nil {
		if jsonBytes, err := json.Marshal(req.PreferredJobTypes); err == nil {
			profile.PreferredJobTypes = sql.NullString{String: string(jsonBytes), Valid: true}
		}
	}
	if req.PreferredLocations != nil {
		if jsonBytes, err := json.Marshal(req.PreferredLocations); err == nil {
			profile.PreferredLocations = sql.NullString{String: string(jsonBytes), Valid: true}
		}
	}
	if req.AvailableFrom != nil {
		if t, err := time.Parse("2006-01-02", *req.AvailableFrom); err == nil {
			profile.AvailableFrom = sql.NullTime{Time: t, Valid: true}
		}
	}
	if req.WillingToRelocate != nil {
		profile.WillingToRelocate = *req.WillingToRelocate
	}
}

// calculateCompleteness calculates the profile completeness percentage
func (s *service) calculateCompleteness(profile *ApplicantProfile) int {
	totalFields := 20
	filledFields := 0

	// Personal Information (4 fields)
	if profile.DateOfBirth.Valid {
		filledFields++
	}
	if profile.Gender.Valid {
		filledFields++
	}
	if profile.Nationality.Valid {
		filledFields++
	}
	if profile.MaritalStatus.Valid {
		filledFields++
	}

	// Address (5 fields)
	if profile.Address.Valid {
		filledFields++
	}
	if profile.City.Valid {
		filledFields++
	}
	if profile.Province.Valid {
		filledFields++
	}
	if profile.PostalCode.Valid {
		filledFields++
	}
	if profile.Country.Valid {
		filledFields++
	}

	// Professional Links (4 fields)
	if profile.LinkedInURL.Valid {
		filledFields++
	}
	if profile.GithubURL.Valid {
		filledFields++
	}
	if profile.PortfolioURL.Valid {
		filledFields++
	}
	if profile.PersonalWebsite.Valid {
		filledFields++
	}

	// Bio/Summary (2 fields)
	if profile.ProfessionalSummary.Valid {
		filledFields++
	}
	if profile.Headline.Valid {
		filledFields++
	}

	// Job Preferences (5 fields)
	if profile.ExpectedSalaryMin.Valid {
		filledFields++
	}
	if profile.ExpectedSalaryMax.Valid {
		filledFields++
	}
	if profile.PreferredJobTypes.Valid {
		filledFields++
	}
	if profile.PreferredLocations.Valid {
		filledFields++
	}
	if profile.AvailableFrom.Valid {
		filledFields++
	}

	return (filledFields * 100) / totalFields
}

// ========================================
// Document Operations
// ========================================

// GetDocuments retrieves all documents for a user
func (s *service) GetDocuments(ctx context.Context, userID uint64) ([]*DocumentResponse, error) {
	docs, err := s.repo.GetDocumentsByUserID(ctx, userID)
	if err != nil {
		return nil, apperrors.NewInternalError("Failed to get documents", err)
	}

	responses := make([]*DocumentResponse, len(docs))
	for i, doc := range docs {
		responses[i] = doc.ToResponse()
	}

	return responses, nil
}

// GetDocumentByID retrieves a specific document
func (s *service) GetDocumentByID(ctx context.Context, docID uint64, userID uint64) (*DocumentResponse, error) {
	doc, err := s.repo.GetDocumentByID(ctx, docID)
	if err != nil {
		return nil, apperrors.NewInternalError("Failed to get document", err)
	}
	if doc == nil {
		return nil, apperrors.NewNotFoundError("Document")
	}

	// Verify ownership
	if doc.UserID != userID {
		return nil, apperrors.NewForbiddenError("You don't have permission to access this document")
	}

	return doc.ToResponse(), nil
}

// CreateDocument creates a new document
func (s *service) CreateDocument(ctx context.Context, userID uint64, doc *ApplicantDocument) (*DocumentResponse, error) {
	doc.UserID = userID

	if err := s.repo.CreateDocument(ctx, doc); err != nil {
		return nil, apperrors.NewInternalError("Failed to create document", err)
	}

	// If marked as primary, set it as primary
	if doc.IsPrimary {
		if err := s.repo.SetPrimaryDocument(ctx, userID, doc.ID, doc.DocType); err != nil {
			// Log error but don't fail
			println("Warning: failed to set primary document:", err.Error())
		}
	}

	return doc.ToResponse(), nil
}

// UpdateDocument updates a document's metadata
func (s *service) UpdateDocument(ctx context.Context, docID uint64, userID uint64, name string, description string) (*DocumentResponse, error) {
	doc, err := s.repo.GetDocumentByID(ctx, docID)
	if err != nil {
		return nil, apperrors.NewInternalError("Failed to get document", err)
	}
	if doc == nil {
		return nil, apperrors.NewNotFoundError("Document")
	}

	// Verify ownership
	if doc.UserID != userID {
		return nil, apperrors.NewForbiddenError("You don't have permission to update this document")
	}

	if name != "" {
		doc.DocName = name
	}
	if description != "" {
		doc.Description = sql.NullString{String: description, Valid: true}
	}

	if err := s.repo.UpdateDocument(ctx, doc); err != nil {
		return nil, apperrors.NewInternalError("Failed to update document", err)
	}

	return doc.ToResponse(), nil
}

// DeleteDocument deletes a document
func (s *service) DeleteDocument(ctx context.Context, docID uint64, userID uint64) error {
	doc, err := s.repo.GetDocumentByID(ctx, docID)
	if err != nil {
		return apperrors.NewInternalError("Failed to get document", err)
	}
	if doc == nil {
		return apperrors.NewNotFoundError("Document")
	}

	// Verify ownership
	if doc.UserID != userID {
		return apperrors.NewForbiddenError("You don't have permission to delete this document")
	}

	if err := s.repo.DeleteDocument(ctx, docID); err != nil {
		return apperrors.NewInternalError("Failed to delete document", err)
	}

	return nil
}

// SetPrimaryDocument sets a document as the primary CV
func (s *service) SetPrimaryDocument(ctx context.Context, docID uint64, userID uint64) error {
	doc, err := s.repo.GetDocumentByID(ctx, docID)
	if err != nil {
		return apperrors.NewInternalError("Failed to get document", err)
	}
	if doc == nil {
		return apperrors.NewNotFoundError("Document")
	}

	// Verify ownership
	if doc.UserID != userID {
		return apperrors.NewForbiddenError("You don't have permission to modify this document")
	}

	if err := s.repo.SetPrimaryDocument(ctx, userID, docID, doc.DocType); err != nil {
		return apperrors.NewInternalError("Failed to set primary document", err)
	}

	return nil
}
