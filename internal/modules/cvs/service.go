package cvs

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"

	apperrors "github.com/karirnusantara/api/internal/shared/errors"
)

// Service defines the CV service interface
type Service interface {
	CreateOrUpdate(ctx context.Context, userID uint64, req *CreateCVRequest) (*CVResponse, error)
	GetByUserID(ctx context.Context, userID uint64) (*CVResponse, error)
	GetByID(ctx context.Context, id uint64) (*CVResponse, error)
	Delete(ctx context.Context, userID uint64) error
	CreateSnapshot(ctx context.Context, userID uint64) (*CVSnapshot, error)
	GetSnapshotByID(ctx context.Context, id uint64) (*CVSnapshotResponse, error)
	CalculateCompleteness(cv *CV) int
}

type service struct {
	repo Repository
}

// NewService creates a new CV service
func NewService(repo Repository) Service {
	return &service{repo: repo}
}

// CreateOrUpdate creates or updates a user's CV
func (s *service) CreateOrUpdate(ctx context.Context, userID uint64, req *CreateCVRequest) (*CVResponse, error) {
	// Check if CV exists
	existingCV, err := s.repo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, apperrors.NewInternalError("Failed to check existing CV", err)
	}

	// Convert request to JSON
	personalInfoJSON, _ := json.Marshal(req.PersonalInfo)
	educationJSON, _ := json.Marshal(req.Education)
	experienceJSON, _ := json.Marshal(req.Experience)
	skillsJSON, _ := json.Marshal(req.Skills)
	certificationsJSON, _ := json.Marshal(req.Certifications)
	languagesJSON, _ := json.Marshal(req.Languages)
	projectsJSON, _ := json.Marshal(req.Projects)

	cv := &CV{
		UserID:         userID,
		PersonalInfo:   personalInfoJSON,
		Education:      educationJSON,
		Experience:     experienceJSON,
		Skills:         skillsJSON,
		Certifications: certificationsJSON,
		Languages:      languagesJSON,
		Projects:       projectsJSON,
	}

	// Parse for completeness calculation
	cv.PersonalInfoParsed = req.PersonalInfo
	cv.EducationParsed = req.Education
	cv.ExperienceParsed = req.Experience
	cv.SkillsParsed = req.Skills
	cv.CertificationsParsed = req.Certifications
	cv.LanguagesParsed = req.Languages
	cv.ProjectsParsed = req.Projects

	cv.CompletenessScore = s.CalculateCompleteness(cv)

	if existingCV != nil {
		// Update existing CV
		cv.ID = existingCV.ID
		if err := s.repo.Update(ctx, cv); err != nil {
			return nil, apperrors.NewInternalError("Failed to update CV", err)
		}
	} else {
		// Create new CV
		if err := s.repo.Create(ctx, cv); err != nil {
			return nil, apperrors.NewInternalError("Failed to create CV", err)
		}
	}

	return s.GetByUserID(ctx, userID)
}

// GetByUserID retrieves a CV by user ID
func (s *service) GetByUserID(ctx context.Context, userID uint64) (*CVResponse, error) {
	cv, err := s.repo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, apperrors.NewInternalError("Failed to get CV", err)
	}
	if cv == nil {
		return nil, apperrors.NewNotFoundError("CV")
	}

	if err := cv.ParseFields(); err != nil {
		return nil, apperrors.NewInternalError("Failed to parse CV data", err)
	}

	return cv.ToResponse(), nil
}

// GetByID retrieves a CV by ID
func (s *service) GetByID(ctx context.Context, id uint64) (*CVResponse, error) {
	cv, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, apperrors.NewInternalError("Failed to get CV", err)
	}
	if cv == nil {
		return nil, apperrors.NewNotFoundError("CV")
	}

	if err := cv.ParseFields(); err != nil {
		return nil, apperrors.NewInternalError("Failed to parse CV data", err)
	}

	return cv.ToResponse(), nil
}

// Delete deletes a user's CV
func (s *service) Delete(ctx context.Context, userID uint64) error {
	cv, err := s.repo.GetByUserID(ctx, userID)
	if err != nil {
		return apperrors.NewInternalError("Failed to get CV", err)
	}
	if cv == nil {
		return apperrors.NewNotFoundError("CV")
	}

	if err := s.repo.Delete(ctx, cv.ID); err != nil {
		return apperrors.NewInternalError("Failed to delete CV", err)
	}

	return nil
}

// CreateSnapshot creates an immutable snapshot of the current CV
func (s *service) CreateSnapshot(ctx context.Context, userID uint64) (*CVSnapshot, error) {
	cv, err := s.repo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, apperrors.NewInternalError("Failed to get CV", err)
	}
	if cv == nil {
		return nil, apperrors.NewNotFoundError("CV")
	}

	// Create snapshot hash from all CV data
	hashData := string(cv.PersonalInfo) + string(cv.Education) + string(cv.Experience) +
		string(cv.Skills) + string(cv.Certifications) + string(cv.Languages) + string(cv.Projects)
	hash := sha256.Sum256([]byte(hashData))

	snapshot := &CVSnapshot{
		CVID:              cv.ID,
		UserID:            userID,
		PersonalInfo:      cv.PersonalInfo,
		Education:         cv.Education,
		Experience:        cv.Experience,
		Skills:            cv.Skills,
		Certifications:    cv.Certifications,
		Languages:         cv.Languages,
		Projects:          cv.Projects,
		SnapshotHash:      hex.EncodeToString(hash[:]),
		CompletenessScore: cv.CompletenessScore,
	}

	if err := s.repo.CreateSnapshot(ctx, snapshot); err != nil {
		return nil, apperrors.NewInternalError("Failed to create CV snapshot", err)
	}

	return snapshot, nil
}

// GetSnapshotByID retrieves a CV snapshot by ID
func (s *service) GetSnapshotByID(ctx context.Context, id uint64) (*CVSnapshotResponse, error) {
	snapshot, err := s.repo.GetSnapshotByID(ctx, id)
	if err != nil {
		return nil, apperrors.NewInternalError("Failed to get CV snapshot", err)
	}
	if snapshot == nil {
		return nil, apperrors.NewNotFoundError("CV Snapshot")
	}

	if err := snapshot.ParseFields(); err != nil {
		return nil, apperrors.NewInternalError("Failed to parse CV snapshot data", err)
	}

	return snapshot.ToResponse(), nil
}

// CalculateCompleteness calculates the completeness score of a CV (0-100)
func (s *service) CalculateCompleteness(cv *CV) int {
	score := 0
	maxScore := 100

	// Personal Info (30 points)
	if cv.PersonalInfoParsed != nil {
		if cv.PersonalInfoParsed.FullName != "" {
			score += 10
		}
		if cv.PersonalInfoParsed.Email != "" {
			score += 5
		}
		if cv.PersonalInfoParsed.Phone != "" {
			score += 5
		}
		if cv.PersonalInfoParsed.Summary != "" {
			score += 10
		}
	}

	// Education (20 points)
	if len(cv.EducationParsed) > 0 {
		score += 20
	}

	// Experience (25 points)
	if len(cv.ExperienceParsed) > 0 {
		score += 15
		if len(cv.ExperienceParsed) >= 2 {
			score += 10
		}
	}

	// Skills (15 points)
	if len(cv.SkillsParsed) > 0 {
		score += 10
		if len(cv.SkillsParsed) >= 5 {
			score += 5
		}
	}

	// Certifications (5 points)
	if len(cv.CertificationsParsed) > 0 {
		score += 5
	}

	// Languages (5 points)
	if len(cv.LanguagesParsed) > 0 {
		score += 5
	}

	// Cap at maxScore
	if score > maxScore {
		score = maxScore
	}

	return score
}
