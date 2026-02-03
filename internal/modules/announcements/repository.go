package announcements

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
)

// Repository handles database operations for announcements
type Repository struct {
	db *sqlx.DB
}

// NewRepository creates a new announcements repository
func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

// Create inserts a new announcement
func (r *Repository) Create(ctx context.Context, announcement *Announcement) error {
	query := `
		INSERT INTO announcements (
			title, content, type, target_audience, is_active, priority,
			start_date, end_date, created_by, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW())
	`

	result, err := r.db.ExecContext(ctx, query,
		announcement.Title,
		announcement.Content,
		announcement.Type,
		announcement.TargetAudience,
		announcement.IsActive,
		announcement.Priority,
		announcement.StartDate,
		announcement.EndDate,
		announcement.CreatedBy,
	)
	if err != nil {
		return fmt.Errorf("failed to create announcement: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}
	announcement.ID = uint64(id)

	return nil
}

// GetByID retrieves an announcement by ID
func (r *Repository) GetByID(ctx context.Context, id uint64) (*Announcement, error) {
	var announcement Announcement
	query := `SELECT * FROM announcements WHERE id = ?`

	err := r.db.GetContext(ctx, &announcement, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get announcement: %w", err)
	}

	return &announcement, nil
}

// Update updates an existing announcement
func (r *Repository) Update(ctx context.Context, announcement *Announcement) error {
	query := `
		UPDATE announcements SET
			title = ?,
			content = ?,
			type = ?,
			target_audience = ?,
			is_active = ?,
			priority = ?,
			start_date = ?,
			end_date = ?,
			updated_by = ?,
			updated_at = NOW()
		WHERE id = ?
	`

	_, err := r.db.ExecContext(ctx, query,
		announcement.Title,
		announcement.Content,
		announcement.Type,
		announcement.TargetAudience,
		announcement.IsActive,
		announcement.Priority,
		announcement.StartDate,
		announcement.EndDate,
		announcement.UpdatedBy,
		announcement.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update announcement: %w", err)
	}

	return nil
}

// Delete removes an announcement
func (r *Repository) Delete(ctx context.Context, id uint64) error {
	query := `DELETE FROM announcements WHERE id = ?`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete announcement: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("announcement not found")
	}

	return nil
}

// ToggleStatus toggles the active status of an announcement
func (r *Repository) ToggleStatus(ctx context.Context, id uint64, isActive bool, updatedBy uint64) error {
	query := `UPDATE announcements SET is_active = ?, updated_by = ?, updated_at = NOW() WHERE id = ?`

	result, err := r.db.ExecContext(ctx, query, isActive, updatedBy, id)
	if err != nil {
		return fmt.Errorf("failed to toggle status: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("announcement not found")
	}

	return nil
}

// List retrieves announcements with filtering and pagination
func (r *Repository) List(ctx context.Context, filter AnnouncementFilter) ([]Announcement, int64, error) {
	var announcements []Announcement
	var total int64

	// Build WHERE clause
	var conditions []string
	var args []interface{}

	if filter.Type != "" {
		conditions = append(conditions, "type = ?")
		args = append(args, filter.Type)
	}

	if filter.TargetAudience != "" {
		conditions = append(conditions, "target_audience = ?")
		args = append(args, filter.TargetAudience)
	}

	if filter.IsActive != nil {
		conditions = append(conditions, "is_active = ?")
		args = append(args, *filter.IsActive)
	}

	if filter.Search != "" {
		conditions = append(conditions, "(title LIKE ? OR content LIKE ?)")
		searchTerm := "%" + filter.Search + "%"
		args = append(args, searchTerm, searchTerm)
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	// Count total
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM announcements %s", whereClause)
	err := r.db.GetContext(ctx, &total, countQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count announcements: %w", err)
	}

	// Get paginated results
	if filter.Limit <= 0 {
		filter.Limit = 10
	}
	if filter.Page <= 0 {
		filter.Page = 1
	}
	offset := (filter.Page - 1) * filter.Limit

	listQuery := fmt.Sprintf(`
		SELECT * FROM announcements %s
		ORDER BY priority DESC, created_at DESC
		LIMIT ? OFFSET ?
	`, whereClause)

	args = append(args, filter.Limit, offset)

	err = r.db.SelectContext(ctx, &announcements, listQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list announcements: %w", err)
	}

	return announcements, total, nil
}

// GetActiveByType retrieves active announcements by type for public endpoints
func (r *Repository) GetActiveByType(ctx context.Context, announcementType string, targetAudience string) ([]Announcement, error) {
	var announcements []Announcement

	query := `
		SELECT * FROM announcements 
		WHERE type = ? 
		AND is_active = true
		AND (target_audience = 'all' OR target_audience = ?)
		AND (start_date IS NULL OR start_date <= NOW())
		AND (end_date IS NULL OR end_date >= NOW())
		ORDER BY priority DESC, created_at DESC
	`

	err := r.db.SelectContext(ctx, &announcements, query, announcementType, targetAudience)
	if err != nil {
		return nil, fmt.Errorf("failed to get active announcements: %w", err)
	}

	return announcements, nil
}

// GetAllActive retrieves all active announcements for a specific audience
func (r *Repository) GetAllActive(ctx context.Context, targetAudience string) ([]Announcement, error) {
	var announcements []Announcement

	query := `
		SELECT * FROM announcements 
		WHERE is_active = true
		AND (target_audience = 'all' OR target_audience = ?)
		AND (start_date IS NULL OR start_date <= NOW())
		AND (end_date IS NULL OR end_date >= NOW())
		ORDER BY type, priority DESC, created_at DESC
	`

	err := r.db.SelectContext(ctx, &announcements, query, targetAudience)
	if err != nil {
		return nil, fmt.Errorf("failed to get active announcements: %w", err)
	}

	return announcements, nil
}
