package applications

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
)

// Repository defines the applications repository interface
type Repository interface {
	Create(ctx context.Context, app *Application) error
	GetByID(ctx context.Context, id uint64) (*Application, error)
	GetByUserAndJob(ctx context.Context, userID, jobID uint64) (*Application, error)
	Update(ctx context.Context, app *Application) error
	List(ctx context.Context, params ApplicationListParams) ([]*Application, int64, error)

	// Timeline
	AddTimelineEvent(ctx context.Context, event *TimelineEvent) error
	GetTimeline(ctx context.Context, applicationID uint64) ([]TimelineEvent, error)
	GetTimelineForApplicant(ctx context.Context, applicationID uint64) ([]TimelineEvent, error)

	// Related data
	GetJobInfo(ctx context.Context, jobID uint64) (*JobInfo, error)
	GetApplicantInfo(ctx context.Context, userID uint64) (*ApplicantInfo, error)
	GetCVSnapshotInfo(ctx context.Context, snapshotID uint64) (*CVSnapshotInfo, error)
	GetCompanyIDByUserID(ctx context.Context, userID uint64) (uint64, error)
}

type mysqlRepository struct {
	db *sqlx.DB
}

// NewRepository creates a new applications repository
func NewRepository(db *sqlx.DB) Repository {
	return &mysqlRepository{db: db}
}

// Create creates a new application
func (r *mysqlRepository) Create(ctx context.Context, app *Application) error {
	query := `
		INSERT INTO applications (
			user_id, job_id, cv_snapshot_id, cover_letter, current_status,
			applied_at, last_status_update, created_at, updated_at
		) VALUES (
			?, ?, ?, ?, ?,
			NOW(), NOW(), NOW(), NOW()
		)
	`

	result, err := r.db.ExecContext(ctx, query,
		app.UserID, app.JobID, app.CVSnapshotID, app.CoverLetter, app.CurrentStatus,
	)
	if err != nil {
		return fmt.Errorf("failed to create application: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	app.ID = uint64(id)
	return nil
}

// GetByID retrieves an application by ID
func (r *mysqlRepository) GetByID(ctx context.Context, id uint64) (*Application, error) {
	query := `
		SELECT id, user_id, job_id, cv_snapshot_id, cover_letter, current_status,
			   applied_at, last_status_update, created_at, updated_at
		FROM applications
		WHERE id = ?
	`

	var app Application
	if err := r.db.GetContext(ctx, &app, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get application: %w", err)
	}

	return &app, nil
}

// GetByUserAndJob retrieves an application by user and job ID
func (r *mysqlRepository) GetByUserAndJob(ctx context.Context, userID, jobID uint64) (*Application, error) {
	query := `
		SELECT id, user_id, job_id, cv_snapshot_id, cover_letter, current_status,
			   applied_at, last_status_update, created_at, updated_at
		FROM applications
		WHERE user_id = ? AND job_id = ?
	`

	var app Application
	if err := r.db.GetContext(ctx, &app, query, userID, jobID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get application: %w", err)
	}

	return &app, nil
}

// Update updates an application
func (r *mysqlRepository) Update(ctx context.Context, app *Application) error {
	query := `
		UPDATE applications SET
			current_status = ?,
			last_status_update = NOW(),
			updated_at = NOW()
		WHERE id = ?
	`

	_, err := r.db.ExecContext(ctx, query, app.CurrentStatus, app.ID)
	if err != nil {
		return fmt.Errorf("failed to update application: %w", err)
	}

	return nil
}

// List retrieves applications with filtering
func (r *mysqlRepository) List(ctx context.Context, params ApplicationListParams) ([]*Application, int64, error) {
	var conditions []string
	var args []interface{}

	if params.UserID != nil {
		conditions = append(conditions, "a.user_id = ?")
		args = append(args, *params.UserID)
	}

	if params.JobID != nil {
		conditions = append(conditions, "a.job_id = ?")
		args = append(args, *params.JobID)
	}

	if params.CompanyID != nil {
		conditions = append(conditions, "j.company_id = ?")
		args = append(args, *params.CompanyID)
	}

	if params.Status != "" {
		conditions = append(conditions, "a.current_status = ?")
		args = append(args, params.Status)
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	// Count total
	countQuery := fmt.Sprintf(`
		SELECT COUNT(*)
		FROM applications a
		JOIN jobs j ON a.job_id = j.id
		%s
	`, whereClause)

	var total int64
	if err := r.db.GetContext(ctx, &total, countQuery, args...); err != nil {
		return nil, 0, fmt.Errorf("failed to count applications: %w", err)
	}

	// Build query
	orderBy := "a.applied_at DESC"
	if params.SortBy != "" {
		order := "ASC"
		if params.SortOrder == "desc" {
			order = "DESC"
		}
		orderBy = fmt.Sprintf("a.%s %s", params.SortBy, order)
	}

	offset := (params.Page - 1) * params.PerPage
	query := fmt.Sprintf(`
		SELECT a.id, a.user_id, a.job_id, a.cv_snapshot_id, a.cover_letter, a.current_status,
			   a.applied_at, a.last_status_update, a.created_at, a.updated_at
		FROM applications a
		JOIN jobs j ON a.job_id = j.id
		%s
		ORDER BY %s
		LIMIT ? OFFSET ?
	`, whereClause, orderBy)

	args = append(args, params.PerPage, offset)

	var apps []*Application
	if err := r.db.SelectContext(ctx, &apps, query, args...); err != nil {
		return nil, 0, fmt.Errorf("failed to list applications: %w", err)
	}

	return apps, total, nil
}

// AddTimelineEvent adds a timeline event
func (r *mysqlRepository) AddTimelineEvent(ctx context.Context, event *TimelineEvent) error {
	query := `
		INSERT INTO application_timelines (
			application_id, status, note, is_visible_to_applicant,
			updated_by_type, updated_by_id, scheduled_at, scheduled_location, scheduled_notes,
			interview_type, meeting_link, meeting_platform, interview_address, contact_person, contact_phone,
			created_at
		) VALUES (
			?, ?, ?, ?,
			?, ?, ?, ?, ?,
			?, ?, ?, ?, ?, ?,
			NOW()
		)
	`

	result, err := r.db.ExecContext(ctx, query,
		event.ApplicationID, event.Status, event.Note, event.IsVisibleToApplicant,
		event.UpdatedByType, event.UpdatedByID, event.ScheduledAt, event.ScheduledLocation, event.ScheduledNotes,
		event.InterviewType, event.MeetingLink, event.MeetingPlatform, event.InterviewAddress, event.ContactPerson, event.ContactPhone,
	)
	if err != nil {
		return fmt.Errorf("failed to add timeline event: %w", err)
	}

	id, _ := result.LastInsertId()
	event.ID = uint64(id)
	return nil
}

// GetTimeline retrieves all timeline events for an application
func (r *mysqlRepository) GetTimeline(ctx context.Context, applicationID uint64) ([]TimelineEvent, error) {
	query := `
		SELECT id, application_id, status, note, is_visible_to_applicant,
			   updated_by_type, updated_by_id, scheduled_at, scheduled_location, scheduled_notes,
			   interview_type, meeting_link, meeting_platform, interview_address, contact_person, contact_phone,
			   created_at
		FROM application_timelines
		WHERE application_id = ?
		ORDER BY created_at ASC
	`

	var events []TimelineEvent
	if err := r.db.SelectContext(ctx, &events, query, applicationID); err != nil {
		return nil, fmt.Errorf("failed to get timeline: %w", err)
	}

	return events, nil
}

// GetTimelineForApplicant retrieves visible timeline events for an applicant
func (r *mysqlRepository) GetTimelineForApplicant(ctx context.Context, applicationID uint64) ([]TimelineEvent, error) {
	query := `
		SELECT id, application_id, status, note, is_visible_to_applicant,
			   updated_by_type, updated_by_id, scheduled_at, scheduled_location, scheduled_notes,
			   interview_type, meeting_link, meeting_platform, interview_address, contact_person, contact_phone,
			   created_at
		FROM application_timelines
		WHERE application_id = ? AND is_visible_to_applicant = TRUE
		ORDER BY created_at ASC
	`

	var events []TimelineEvent
	if err := r.db.SelectContext(ctx, &events, query, applicationID); err != nil {
		return nil, fmt.Errorf("failed to get timeline: %w", err)
	}

	return events, nil
}

// GetJobInfo retrieves job info for an application
func (r *mysqlRepository) GetJobInfo(ctx context.Context, jobID uint64) (*JobInfo, error) {
	query := `
		SELECT j.id, j.title, j.city, j.province, j.status,
			   c.id as company_id, c.company_name, c.company_logo_url
		FROM jobs j
		JOIN companies c ON j.company_id = c.id
		WHERE j.id = ? AND j.deleted_at IS NULL
	`

	var result struct {
		ID          uint64         `db:"id"`
		Title       string         `db:"title"`
		City        string         `db:"city"`
		Province    string         `db:"province"`
		Status      string         `db:"status"`
		CompanyID   uint64         `db:"company_id"`
		CompanyName sql.NullString `db:"company_name"`
		CompanyLogo sql.NullString `db:"company_logo_url"`
	}

	if err := r.db.GetContext(ctx, &result, query, jobID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get job info: %w", err)
	}

	return &JobInfo{
		ID:       result.ID,
		Title:    result.Title,
		City:     result.City,
		Province: result.Province,
		Status:   result.Status,
		Company: CompanyInfo{
			ID:      result.CompanyID,
			Name:    result.CompanyName.String,
			LogoURL: result.CompanyLogo.String,
		},
	}, nil
}

// GetApplicantInfo retrieves applicant info
func (r *mysqlRepository) GetApplicantInfo(ctx context.Context, userID uint64) (*ApplicantInfo, error) {
	query := `SELECT id, full_name, email, phone, avatar_url FROM users WHERE id = ?`

	var result struct {
		ID        uint64         `db:"id"`
		Name      string         `db:"full_name"`
		Email     string         `db:"email"`
		Phone     sql.NullString `db:"phone"`
		AvatarURL sql.NullString `db:"avatar_url"`
	}

	if err := r.db.GetContext(ctx, &result, query, userID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get applicant info: %w", err)
	}

	return &ApplicantInfo{
		ID:        result.ID,
		Name:      result.Name,
		Email:     result.Email,
		Phone:     result.Phone.String,
		AvatarURL: result.AvatarURL.String,
	}, nil
}

// GetCVSnapshotInfo retrieves CV snapshot info with full details
func (r *mysqlRepository) GetCVSnapshotInfo(ctx context.Context, snapshotID uint64) (*CVSnapshotInfo, error) {
	query := `SELECT id, completeness_score, personal_info, education, experience, skills, certifications, languages, projects, created_at FROM cv_snapshots WHERE id = ?`

	var result struct {
		ID                uint64    `db:"id"`
		CompletenessScore int       `db:"completeness_score"`
		PersonalInfo      string    `db:"personal_info"`
		Education         string    `db:"education"`
		Experience        string    `db:"experience"`
		Skills            string    `db:"skills"`
		Certifications    string    `db:"certifications"`
		Languages         string    `db:"languages"`
		Projects          string    `db:"projects"`
		CreatedAt         time.Time `db:"created_at"`
	}

	if err := r.db.GetContext(ctx, &result, query, snapshotID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get CV snapshot info: %w", err)
	}

	snapshot := &CVSnapshotInfo{
		ID:                result.ID,
		CompletenessScore: result.CompletenessScore,
		CreatedAt:         result.CreatedAt.Format(time.RFC3339),
	}

	// Parse JSON fields
	if result.PersonalInfo != "" && result.PersonalInfo != "null" {
		json.Unmarshal([]byte(result.PersonalInfo), &snapshot.PersonalInfo)
	}
	if result.Education != "" && result.Education != "null" {
		json.Unmarshal([]byte(result.Education), &snapshot.Education)
	}
	if result.Experience != "" && result.Experience != "null" {
		json.Unmarshal([]byte(result.Experience), &snapshot.Experience)
	}
	if result.Skills != "" && result.Skills != "null" {
		json.Unmarshal([]byte(result.Skills), &snapshot.Skills)
	}
	if result.Certifications != "" && result.Certifications != "null" {
		json.Unmarshal([]byte(result.Certifications), &snapshot.Certifications)
	}
	if result.Languages != "" && result.Languages != "null" {
		json.Unmarshal([]byte(result.Languages), &snapshot.Languages)
	}
	if result.Projects != "" && result.Projects != "null" {
		json.Unmarshal([]byte(result.Projects), &snapshot.Projects)
	}

	return snapshot, nil
}

// GetCompanyIDByUserID retrieves company ID for a given user ID
func (r *mysqlRepository) GetCompanyIDByUserID(ctx context.Context, userID uint64) (uint64, error) {
	var companyID uint64
	query := `SELECT id FROM companies WHERE user_id = ? AND deleted_at IS NULL`
	if err := r.db.GetContext(ctx, &companyID, query, userID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, fmt.Errorf("company not found for user ID %d", userID)
		}
		return 0, fmt.Errorf("failed to get company ID: %w", err)
	}
	return companyID, nil
}
