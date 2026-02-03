package announcements

import (
	"context"
	"database/sql"
	"errors"
	"math"
	"time"
)

var (
	ErrAnnouncementNotFound = errors.New("announcement not found")
	ErrInvalidType          = errors.New("invalid announcement type")
	ErrInvalidAudience      = errors.New("invalid target audience")
)

// Service handles business logic for announcements
type Service struct {
	repo *Repository
}

// NewService creates a new announcements service
func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// Create creates a new announcement
func (s *Service) Create(ctx context.Context, req CreateAnnouncementRequest, adminID uint64) (*AnnouncementResponse, error) {
	// Validate type
	if !isValidType(req.Type) {
		return nil, ErrInvalidType
	}

	// Validate target audience
	if !isValidAudience(req.TargetAudience) {
		return nil, ErrInvalidAudience
	}

	// Set defaults
	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	priority := 0
	if req.Priority != nil {
		priority = *req.Priority
	}

	announcement := &Announcement{
		Title:          req.Title,
		Content:        req.Content,
		Type:           AnnouncementType(req.Type),
		TargetAudience: TargetAudience(req.TargetAudience),
		IsActive:       isActive,
		Priority:       priority,
		CreatedBy:      sql.NullInt64{Int64: int64(adminID), Valid: true},
	}

	// Parse dates if provided
	if req.StartDate != "" {
		startDate, err := time.Parse(time.RFC3339, req.StartDate)
		if err != nil {
			startDate, err = time.Parse("2006-01-02", req.StartDate)
			if err != nil {
				return nil, errors.New("invalid start_date format")
			}
		}
		announcement.StartDate = sql.NullTime{Time: startDate, Valid: true}
	}

	if req.EndDate != "" {
		endDate, err := time.Parse(time.RFC3339, req.EndDate)
		if err != nil {
			endDate, err = time.Parse("2006-01-02", req.EndDate)
			if err != nil {
				return nil, errors.New("invalid end_date format")
			}
		}
		announcement.EndDate = sql.NullTime{Time: endDate, Valid: true}
	}

	err := s.repo.Create(ctx, announcement)
	if err != nil {
		return nil, err
	}

	// Fetch the created announcement
	created, err := s.repo.GetByID(ctx, announcement.ID)
	if err != nil {
		return nil, err
	}

	resp := created.ToResponse()
	return &resp, nil
}

// GetByID retrieves an announcement by ID
func (s *Service) GetByID(ctx context.Context, id uint64) (*AnnouncementResponse, error) {
	announcement, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if announcement == nil {
		return nil, ErrAnnouncementNotFound
	}

	resp := announcement.ToResponse()
	return &resp, nil
}

// Update updates an existing announcement
func (s *Service) Update(ctx context.Context, id uint64, req UpdateAnnouncementRequest, adminID uint64) (*AnnouncementResponse, error) {
	announcement, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if announcement == nil {
		return nil, ErrAnnouncementNotFound
	}

	// Update fields if provided
	if req.Title != nil {
		announcement.Title = *req.Title
	}
	if req.Content != nil {
		announcement.Content = *req.Content
	}
	if req.Type != nil {
		if !isValidType(*req.Type) {
			return nil, ErrInvalidType
		}
		announcement.Type = AnnouncementType(*req.Type)
	}
	if req.TargetAudience != nil {
		if !isValidAudience(*req.TargetAudience) {
			return nil, ErrInvalidAudience
		}
		announcement.TargetAudience = TargetAudience(*req.TargetAudience)
	}
	if req.IsActive != nil {
		announcement.IsActive = *req.IsActive
	}
	if req.Priority != nil {
		announcement.Priority = *req.Priority
	}

	// Parse dates if provided
	if req.StartDate != nil {
		if *req.StartDate == "" {
			announcement.StartDate = sql.NullTime{Valid: false}
		} else {
			startDate, err := time.Parse(time.RFC3339, *req.StartDate)
			if err != nil {
				startDate, err = time.Parse("2006-01-02", *req.StartDate)
				if err != nil {
					return nil, errors.New("invalid start_date format")
				}
			}
			announcement.StartDate = sql.NullTime{Time: startDate, Valid: true}
		}
	}

	if req.EndDate != nil {
		if *req.EndDate == "" {
			announcement.EndDate = sql.NullTime{Valid: false}
		} else {
			endDate, err := time.Parse(time.RFC3339, *req.EndDate)
			if err != nil {
				endDate, err = time.Parse("2006-01-02", *req.EndDate)
				if err != nil {
					return nil, errors.New("invalid end_date format")
				}
			}
			announcement.EndDate = sql.NullTime{Time: endDate, Valid: true}
		}
	}

	announcement.UpdatedBy = sql.NullInt64{Int64: int64(adminID), Valid: true}

	err = s.repo.Update(ctx, announcement)
	if err != nil {
		return nil, err
	}

	// Fetch updated announcement
	updated, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	resp := updated.ToResponse()
	return &resp, nil
}

// Delete removes an announcement
func (s *Service) Delete(ctx context.Context, id uint64) error {
	return s.repo.Delete(ctx, id)
}

// ToggleStatus toggles the active status of an announcement
func (s *Service) ToggleStatus(ctx context.Context, id uint64, isActive bool, adminID uint64) (*AnnouncementResponse, error) {
	err := s.repo.ToggleStatus(ctx, id, isActive, adminID)
	if err != nil {
		return nil, err
	}

	announcement, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if announcement == nil {
		return nil, ErrAnnouncementNotFound
	}

	resp := announcement.ToResponse()
	return &resp, nil
}

// List retrieves announcements with filtering and pagination
func (s *Service) List(ctx context.Context, filter AnnouncementFilter) (*AnnouncementListResponse, error) {
	announcements, total, err := s.repo.List(ctx, filter)
	if err != nil {
		return nil, err
	}

	// Convert to response
	data := make([]AnnouncementResponse, len(announcements))
	for i, a := range announcements {
		data[i] = a.ToResponse()
	}

	// Calculate pagination
	if filter.Limit <= 0 {
		filter.Limit = 10
	}
	if filter.Page <= 0 {
		filter.Page = 1
	}
	totalPages := int(math.Ceil(float64(total) / float64(filter.Limit)))

	return &AnnouncementListResponse{
		Data: data,
		Pagination: PaginationMeta{
			CurrentPage: filter.Page,
			PerPage:     filter.Limit,
			TotalItems:  total,
			TotalPages:  totalPages,
		},
	}, nil
}

// GetActiveByType retrieves active announcements by type for public endpoints
func (s *Service) GetActiveByType(ctx context.Context, announcementType string, targetAudience string) ([]AnnouncementResponse, error) {
	announcements, err := s.repo.GetActiveByType(ctx, announcementType, targetAudience)
	if err != nil {
		return nil, err
	}

	data := make([]AnnouncementResponse, len(announcements))
	for i, a := range announcements {
		data[i] = a.ToResponse()
	}

	return data, nil
}

// GetAllActive retrieves all active announcements for a specific audience
func (s *Service) GetAllActive(ctx context.Context, targetAudience string) ([]AnnouncementResponse, error) {
	announcements, err := s.repo.GetAllActive(ctx, targetAudience)
	if err != nil {
		return nil, err
	}

	data := make([]AnnouncementResponse, len(announcements))
	for i, a := range announcements {
		data[i] = a.ToResponse()
	}

	return data, nil
}

// Helper functions
func isValidType(t string) bool {
	switch t {
	case "notification", "banner", "information":
		return true
	}
	return false
}

func isValidAudience(a string) bool {
	switch a {
	case "all", "company", "candidate", "partner":
		return true
	}
	return false
}
