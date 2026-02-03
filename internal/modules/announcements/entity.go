package announcements

import (
	"database/sql"
	"time"
)

// AnnouncementType represents the type of announcement
type AnnouncementType string

const (
	TypeNotification AnnouncementType = "notification"
	TypeBanner       AnnouncementType = "banner"
	TypeInformation  AnnouncementType = "information"
)

// TargetAudience represents who should see the announcement
type TargetAudience string

const (
	AudienceAll       TargetAudience = "all"
	AudienceCompany   TargetAudience = "company"
	AudienceCandidate TargetAudience = "candidate"
	AudiencePartner   TargetAudience = "partner"
)

// Announcement represents an announcement entity
type Announcement struct {
	ID             uint64           `db:"id" json:"id"`
	Title          string           `db:"title" json:"title"`
	Content        string           `db:"content" json:"content"`
	Type           AnnouncementType `db:"type" json:"type"`
	TargetAudience TargetAudience   `db:"target_audience" json:"target_audience"`
	IsActive       bool             `db:"is_active" json:"is_active"`
	Priority       int              `db:"priority" json:"priority"`
	StartDate      sql.NullTime     `db:"start_date" json:"start_date,omitempty"`
	EndDate        sql.NullTime     `db:"end_date" json:"end_date,omitempty"`
	CreatedBy      sql.NullInt64    `db:"created_by" json:"created_by,omitempty"`
	UpdatedBy      sql.NullInt64    `db:"updated_by" json:"updated_by,omitempty"`
	CreatedAt      time.Time        `db:"created_at" json:"created_at"`
	UpdatedAt      time.Time        `db:"updated_at" json:"updated_at"`
}

// AnnouncementResponse is the API response format
type AnnouncementResponse struct {
	ID             uint64  `json:"id"`
	Title          string  `json:"title"`
	Content        string  `json:"content"`
	Type           string  `json:"type"`
	TargetAudience string  `json:"target_audience"`
	IsActive       bool    `json:"is_active"`
	Priority       int     `json:"priority"`
	StartDate      *string `json:"start_date,omitempty"`
	EndDate        *string `json:"end_date,omitempty"`
	CreatedBy      *uint64 `json:"created_by,omitempty"`
	UpdatedBy      *uint64 `json:"updated_by,omitempty"`
	CreatedAt      string  `json:"created_at"`
	UpdatedAt      string  `json:"updated_at"`
}

// ToResponse converts Announcement to AnnouncementResponse
func (a *Announcement) ToResponse() AnnouncementResponse {
	resp := AnnouncementResponse{
		ID:             a.ID,
		Title:          a.Title,
		Content:        a.Content,
		Type:           string(a.Type),
		TargetAudience: string(a.TargetAudience),
		IsActive:       a.IsActive,
		Priority:       a.Priority,
		CreatedAt:      a.CreatedAt.Format(time.RFC3339),
		UpdatedAt:      a.UpdatedAt.Format(time.RFC3339),
	}

	if a.StartDate.Valid {
		startDate := a.StartDate.Time.Format(time.RFC3339)
		resp.StartDate = &startDate
	}

	if a.EndDate.Valid {
		endDate := a.EndDate.Time.Format(time.RFC3339)
		resp.EndDate = &endDate
	}

	if a.CreatedBy.Valid {
		createdBy := uint64(a.CreatedBy.Int64)
		resp.CreatedBy = &createdBy
	}

	if a.UpdatedBy.Valid {
		updatedBy := uint64(a.UpdatedBy.Int64)
		resp.UpdatedBy = &updatedBy
	}

	return resp
}

// CreateAnnouncementRequest is the request body for creating an announcement
type CreateAnnouncementRequest struct {
	Title          string `json:"title" validate:"required,min=3,max=255"`
	Content        string `json:"content" validate:"required,min=10"`
	Type           string `json:"type" validate:"required,oneof=notification banner information"`
	TargetAudience string `json:"target_audience" validate:"required,oneof=all company candidate partner"`
	IsActive       *bool  `json:"is_active"`
	Priority       *int   `json:"priority"`
	StartDate      string `json:"start_date,omitempty"`
	EndDate        string `json:"end_date,omitempty"`
}

// UpdateAnnouncementRequest is the request body for updating an announcement
type UpdateAnnouncementRequest struct {
	Title          *string `json:"title,omitempty" validate:"omitempty,min=3,max=255"`
	Content        *string `json:"content,omitempty" validate:"omitempty,min=10"`
	Type           *string `json:"type,omitempty" validate:"omitempty,oneof=notification banner information"`
	TargetAudience *string `json:"target_audience,omitempty" validate:"omitempty,oneof=all company candidate partner"`
	IsActive       *bool   `json:"is_active,omitempty"`
	Priority       *int    `json:"priority,omitempty"`
	StartDate      *string `json:"start_date,omitempty"`
	EndDate        *string `json:"end_date,omitempty"`
}

// ToggleStatusRequest is the request body for toggling announcement status
type ToggleStatusRequest struct {
	IsActive bool `json:"is_active"`
}

// AnnouncementFilter for filtering announcements
type AnnouncementFilter struct {
	Type           string
	TargetAudience string
	IsActive       *bool
	Search         string
	Page           int
	Limit          int
}

// AnnouncementListResponse for paginated list
type AnnouncementListResponse struct {
	Data       []AnnouncementResponse `json:"data"`
	Pagination PaginationMeta         `json:"pagination"`
}

// PaginationMeta contains pagination information
type PaginationMeta struct {
	CurrentPage int   `json:"current_page"`
	PerPage     int   `json:"per_page"`
	TotalItems  int64 `json:"total_items"`
	TotalPages  int   `json:"total_pages"`
}
