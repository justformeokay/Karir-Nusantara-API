package cvs

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
)

// Repository defines the CV repository interface
type Repository interface {
	Create(ctx context.Context, cv *CV) error
	GetByUserID(ctx context.Context, userID uint64) (*CV, error)
	GetByID(ctx context.Context, id uint64) (*CV, error)
	Update(ctx context.Context, cv *CV) error
	Delete(ctx context.Context, id uint64) error

	// Snapshots
	CreateSnapshot(ctx context.Context, snapshot *CVSnapshot) error
	GetSnapshotByID(ctx context.Context, id uint64) (*CVSnapshot, error)
	GetSnapshotsByUserID(ctx context.Context, userID uint64) ([]*CVSnapshot, error)
}

type mysqlRepository struct {
	db *sqlx.DB
}

// NewRepository creates a new CV repository
func NewRepository(db *sqlx.DB) Repository {
	return &mysqlRepository{db: db}
}

// Create creates a new CV
func (r *mysqlRepository) Create(ctx context.Context, cv *CV) error {
	query := `
		INSERT INTO cvs (
			user_id, personal_info, education, experience, skills,
			certifications, languages, projects, completeness_score,
			last_updated_at, created_at, updated_at
		) VALUES (
			?, ?, ?, ?, ?,
			?, ?, ?, ?,
			NOW(), NOW(), NOW()
		)
	`

	result, err := r.db.ExecContext(ctx, query,
		cv.UserID, cv.PersonalInfo, cv.Education, cv.Experience, cv.Skills,
		cv.Certifications, cv.Languages, cv.Projects, cv.CompletenessScore,
	)
	if err != nil {
		return fmt.Errorf("failed to create CV: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	cv.ID = uint64(id)
	return nil
}

// GetByUserID retrieves a CV by user ID
func (r *mysqlRepository) GetByUserID(ctx context.Context, userID uint64) (*CV, error) {
	query := `
		SELECT id, user_id, personal_info, education, experience, skills,
			   certifications, languages, projects, completeness_score,
			   last_updated_at, created_at, updated_at
		FROM cvs
		WHERE user_id = ?
	`

	var cv CV
	if err := r.db.GetContext(ctx, &cv, query, userID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get CV by user id: %w", err)
	}

	return &cv, nil
}

// GetByID retrieves a CV by ID
func (r *mysqlRepository) GetByID(ctx context.Context, id uint64) (*CV, error) {
	query := `
		SELECT id, user_id, personal_info, education, experience, skills,
			   certifications, languages, projects, completeness_score,
			   last_updated_at, created_at, updated_at
		FROM cvs
		WHERE id = ?
	`

	var cv CV
	if err := r.db.GetContext(ctx, &cv, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get CV by id: %w", err)
	}

	return &cv, nil
}

// Update updates a CV
func (r *mysqlRepository) Update(ctx context.Context, cv *CV) error {
	query := `
		UPDATE cvs SET
			personal_info = ?,
			education = ?,
			experience = ?,
			skills = ?,
			certifications = ?,
			languages = ?,
			projects = ?,
			completeness_score = ?,
			last_updated_at = NOW(),
			updated_at = NOW()
		WHERE id = ?
	`

	_, err := r.db.ExecContext(ctx, query,
		cv.PersonalInfo, cv.Education, cv.Experience, cv.Skills,
		cv.Certifications, cv.Languages, cv.Projects, cv.CompletenessScore,
		cv.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update CV: %w", err)
	}

	return nil
}

// Delete deletes a CV
func (r *mysqlRepository) Delete(ctx context.Context, id uint64) error {
	query := `DELETE FROM cvs WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete CV: %w", err)
	}
	return nil
}

// CreateSnapshot creates a CV snapshot
func (r *mysqlRepository) CreateSnapshot(ctx context.Context, snapshot *CVSnapshot) error {
	query := `
		INSERT INTO cv_snapshots (
			cv_id, user_id, personal_info, education, experience, skills,
			certifications, languages, projects, snapshot_hash, completeness_score,
			created_at
		) VALUES (
			?, ?, ?, ?, ?, ?,
			?, ?, ?, ?, ?,
			NOW()
		)
	`

	result, err := r.db.ExecContext(ctx, query,
		snapshot.CVID, snapshot.UserID, snapshot.PersonalInfo, snapshot.Education,
		snapshot.Experience, snapshot.Skills, snapshot.Certifications, snapshot.Languages,
		snapshot.Projects, snapshot.SnapshotHash, snapshot.CompletenessScore,
	)
	if err != nil {
		return fmt.Errorf("failed to create CV snapshot: %w", err)
	}

	id, _ := result.LastInsertId()
	snapshot.ID = uint64(id)
	return nil
}

// GetSnapshotByID retrieves a CV snapshot by ID
func (r *mysqlRepository) GetSnapshotByID(ctx context.Context, id uint64) (*CVSnapshot, error) {
	query := `
		SELECT id, cv_id, user_id, personal_info, education, experience, skills,
			   certifications, languages, projects, snapshot_hash, completeness_score, created_at
		FROM cv_snapshots
		WHERE id = ?
	`

	var snapshot CVSnapshot
	if err := r.db.GetContext(ctx, &snapshot, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get CV snapshot: %w", err)
	}

	return &snapshot, nil
}

// GetSnapshotsByUserID retrieves all CV snapshots for a user
func (r *mysqlRepository) GetSnapshotsByUserID(ctx context.Context, userID uint64) ([]*CVSnapshot, error) {
	query := `
		SELECT id, cv_id, user_id, personal_info, education, experience, skills,
			   certifications, languages, projects, snapshot_hash, completeness_score, created_at
		FROM cv_snapshots
		WHERE user_id = ?
		ORDER BY created_at DESC
	`

	var snapshots []*CVSnapshot
	if err := r.db.SelectContext(ctx, &snapshots, query, userID); err != nil {
		return nil, fmt.Errorf("failed to get CV snapshots: %w", err)
	}

	return snapshots, nil
}
