package wishlist

import (
	"context"
	"fmt"

	apperrors "github.com/karirnusantara/api/internal/shared/errors"
)

// Service defines the wishlist service interface
type Service interface {
	SaveJob(ctx context.Context, userID uint64, req *SaveJobRequest) (*WishlistResponse, error)
	RemoveJob(ctx context.Context, userID, jobID uint64) error
	ListSavedJobs(ctx context.Context, userID uint64, params ListParams) ([]*WishlistResponse, int64, error)
	IsSaved(ctx context.Context, userID, jobID uint64) (bool, error)
	GetStats(ctx context.Context, userID uint64) (*WishlistStats, error)
}

// WishlistStats represents wishlist statistics
type WishlistStats struct {
	TotalSaved int64 `json:"total_saved"`
}

type service struct {
	repo Repository
}

// NewService creates a new wishlist service
func NewService(repo Repository) Service {
	return &service{repo: repo}
}

// SaveJob adds a job to user's wishlist
func (s *service) SaveJob(ctx context.Context, userID uint64, req *SaveJobRequest) (*WishlistResponse, error) {
	// Check if already saved
	existing, err := s.repo.GetByUserAndJob(ctx, userID, req.JobID)
	if err != nil {
		return nil, fmt.Errorf("failed to check saved job: %w", err)
	}

	if existing != nil {
		return nil, apperrors.NewConflictError("Job already saved to wishlist")
	}

	// Save the job
	saved, err := s.repo.Save(ctx, userID, req.JobID)
	if err != nil {
		return nil, fmt.Errorf("failed to save job: %w", err)
	}

	// Get with details
	items, _, err := s.repo.ListByUser(ctx, userID, ListParams{Page: 1, PerPage: 100})
	if err != nil {
		return nil, fmt.Errorf("failed to get saved job details: %w", err)
	}

	// Find the just-saved item
	for _, item := range items {
		if item.ID == saved.ID {
			return &WishlistResponse{
				ID:      item.ID,
				JobID:   item.JobID,
				SavedAt: item.CreatedAt,
				Job:     item.Job,
			}, nil
		}
	}

	return &WishlistResponse{
		ID:      saved.ID,
		JobID:   saved.JobID,
		SavedAt: saved.CreatedAt,
	}, nil
}

// RemoveJob removes a job from user's wishlist
func (s *service) RemoveJob(ctx context.Context, userID, jobID uint64) error {
	// Check if saved
	existing, err := s.repo.GetByUserAndJob(ctx, userID, jobID)
	if err != nil {
		return fmt.Errorf("failed to check saved job: %w", err)
	}

	if existing == nil {
		return apperrors.NewNotFoundError("Saved job not found")
	}

	if err := s.repo.Remove(ctx, userID, jobID); err != nil {
		return fmt.Errorf("failed to remove saved job: %w", err)
	}

	return nil
}

// ListSavedJobs lists all saved jobs for a user
func (s *service) ListSavedJobs(ctx context.Context, userID uint64, params ListParams) ([]*WishlistResponse, int64, error) {
	if params.Page < 1 {
		params.Page = 1
	}
	if params.PerPage < 1 {
		params.PerPage = 20
	}
	if params.PerPage > 100 {
		params.PerPage = 100
	}

	items, total, err := s.repo.ListByUser(ctx, userID, params)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list saved jobs: %w", err)
	}

	responses := make([]*WishlistResponse, len(items))
	for i, item := range items {
		responses[i] = &WishlistResponse{
			ID:      item.ID,
			JobID:   item.JobID,
			SavedAt: item.CreatedAt,
			Job:     item.Job,
		}
	}

	return responses, total, nil
}

// IsSaved checks if a job is saved by the user
func (s *service) IsSaved(ctx context.Context, userID, jobID uint64) (bool, error) {
	return s.repo.IsSaved(ctx, userID, jobID)
}

// GetStats returns wishlist statistics for a user
func (s *service) GetStats(ctx context.Context, userID uint64) (*WishlistStats, error) {
	count, err := s.repo.CountByUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get wishlist stats: %w", err)
	}

	return &WishlistStats{
		TotalSaved: count,
	}, nil
}
