package applications

import (
	"context"
	"database/sql"
	"time"

	"github.com/karirnusantara/api/internal/modules/cvs"
	"github.com/karirnusantara/api/internal/modules/jobs"
	apperrors "github.com/karirnusantara/api/internal/shared/errors"
)

// Service defines the applications service interface
type Service interface {
	Apply(ctx context.Context, userID uint64, req *ApplyJobRequest) (*ApplicationResponse, error)
	GetByID(ctx context.Context, id uint64, viewerID uint64, isCompany bool) (*ApplicationResponse, error)
	ListByUser(ctx context.Context, userID uint64, params ApplicationListParams) ([]*ApplicationResponse, int64, error)
	ListByCompany(ctx context.Context, companyID uint64, params ApplicationListParams) ([]*ApplicationResponse, int64, error)
	ListByJob(ctx context.Context, jobID uint64, companyID uint64, params ApplicationListParams) ([]*ApplicationResponse, int64, error)
	UpdateStatus(ctx context.Context, applicationID uint64, companyID uint64, req *UpdateStatusRequest) (*ApplicationResponse, error)
	Withdraw(ctx context.Context, applicationID uint64, userID uint64, reason string) error
}

type service struct {
	repo       Repository
	cvService  cvs.Service
	jobService jobs.Service
}

// NewService creates a new applications service
func NewService(repo Repository, cvService cvs.Service, jobService jobs.Service) Service {
	return &service{
		repo:       repo,
		cvService:  cvService,
		jobService: jobService,
	}
}

// Apply submits a job application
func (s *service) Apply(ctx context.Context, userID uint64, req *ApplyJobRequest) (*ApplicationResponse, error) {
	// Check if already applied
	existing, err := s.repo.GetByUserAndJob(ctx, userID, req.JobID)
	if err != nil {
		return nil, apperrors.NewInternalError("Failed to check existing application", err)
	}
	if existing != nil {
		return nil, apperrors.NewConflictError("You have already applied to this job")
	}

	// Get job to verify it exists and is active
	job, err := s.jobService.GetByID(ctx, req.JobID)
	if err != nil {
		return nil, err
	}
	if job.Status != jobs.JobStatusActive {
		return nil, apperrors.NewBadRequestError("This job is no longer accepting applications")
	}

	// Create CV snapshot
	snapshot, err := s.cvService.CreateSnapshot(ctx, userID)
	if err != nil {
		// Check if error is "CV not found"
		if appErr := apperrors.GetAppError(err); appErr != nil && appErr.Code == apperrors.ErrCodeNotFound {
			return nil, apperrors.NewBadRequestError("Please create a CV before applying")
		}
		return nil, apperrors.NewInternalError("Failed to create CV snapshot", err)
	}

	// Create application
	app := &Application{
		UserID:        userID,
		JobID:         req.JobID,
		CVSnapshotID:  snapshot.ID,
		CurrentStatus: StatusSubmitted,
	}

	if req.CoverLetter != "" {
		app.CoverLetter = sql.NullString{String: req.CoverLetter, Valid: true}
	}

	if err := s.repo.Create(ctx, app); err != nil {
		return nil, apperrors.NewInternalError("Failed to create application", err)
	}

	// Add initial timeline event
	event := &TimelineEvent{
		ApplicationID:        app.ID,
		Status:               StatusSubmitted,
		Note:                 sql.NullString{String: "Lamaran berhasil dikirim", Valid: true},
		IsVisibleToApplicant: true,
		UpdatedByType:        "system",
	}

	if err := s.repo.AddTimelineEvent(ctx, event); err != nil {
		return nil, apperrors.NewInternalError("Failed to add timeline event", err)
	}

	return s.GetByID(ctx, app.ID, userID, false)
}

// GetByID retrieves an application by ID
func (s *service) GetByID(ctx context.Context, id uint64, viewerID uint64, isCompany bool) (*ApplicationResponse, error) {
	app, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, apperrors.NewInternalError("Failed to get application", err)
	}
	if app == nil {
		return nil, apperrors.NewNotFoundError("Application")
	}

	// Load related data
	if err := s.loadApplicationRelations(ctx, app, isCompany); err != nil {
		return nil, err
	}

	// Authorization check
	if !isCompany {
		// Job seeker can only view their own applications
		if app.UserID != viewerID {
			return nil, apperrors.NewForbiddenError("You don't have permission to view this application")
		}
	}

	return app.ToResponse(), nil
}

// ListByUser retrieves applications for a user
func (s *service) ListByUser(ctx context.Context, userID uint64, params ApplicationListParams) ([]*ApplicationResponse, int64, error) {
	params.UserID = &userID

	apps, total, err := s.repo.List(ctx, params)
	if err != nil {
		return nil, 0, apperrors.NewInternalError("Failed to list applications", err)
	}

	responses := make([]*ApplicationResponse, len(apps))
	for i, app := range apps {
		if err := s.loadApplicationRelations(ctx, app, false); err != nil {
			return nil, 0, err
		}
		responses[i] = app.ToResponse()
	}

	return responses, total, nil
}

// ListByCompany retrieves applications for a company
func (s *service) ListByCompany(ctx context.Context, companyID uint64, params ApplicationListParams) ([]*ApplicationResponse, int64, error) {
	params.CompanyID = &companyID

	apps, total, err := s.repo.List(ctx, params)
	if err != nil {
		return nil, 0, apperrors.NewInternalError("Failed to list applications", err)
	}

	responses := make([]*ApplicationResponse, len(apps))
	for i, app := range apps {
		if err := s.loadApplicationRelations(ctx, app, true); err != nil {
			return nil, 0, err
		}
		responses[i] = app.ToResponse()
	}

	return responses, total, nil
}

// ListByJob retrieves applications for a specific job
func (s *service) ListByJob(ctx context.Context, jobID uint64, companyID uint64, params ApplicationListParams) ([]*ApplicationResponse, int64, error) {
	// Verify job belongs to company
	job, err := s.jobService.GetByID(ctx, jobID)
	if err != nil {
		return nil, 0, err
	}
	if job.Company.ID != companyID {
		return nil, 0, apperrors.NewForbiddenError("You don't have permission to view these applications")
	}

	params.JobID = &jobID

	apps, total, err := s.repo.List(ctx, params)
	if err != nil {
		return nil, 0, apperrors.NewInternalError("Failed to list applications", err)
	}

	responses := make([]*ApplicationResponse, len(apps))
	for i, app := range apps {
		if err := s.loadApplicationRelations(ctx, app, true); err != nil {
			return nil, 0, err
		}
		responses[i] = app.ToResponse()
	}

	return responses, total, nil
}

// UpdateStatus updates the application status (company action)
func (s *service) UpdateStatus(ctx context.Context, applicationID uint64, companyID uint64, req *UpdateStatusRequest) (*ApplicationResponse, error) {
	// Get application
	app, err := s.repo.GetByID(ctx, applicationID)
	if err != nil {
		return nil, apperrors.NewInternalError("Failed to get application", err)
	}
	if app == nil {
		return nil, apperrors.NewNotFoundError("Application")
	}

	// Verify job belongs to company
	job, err := s.repo.GetJobInfo(ctx, app.JobID)
	if err != nil {
		return nil, apperrors.NewInternalError("Failed to get job", err)
	}
	if job.Company.ID != companyID {
		return nil, apperrors.NewForbiddenError("You don't have permission to update this application")
	}

	// Check if status is terminal
	if IsTerminalStatus(app.CurrentStatus) {
		return nil, apperrors.NewBadRequestError("Cannot update a terminal status")
	}

	// Validate status transition
	if !IsValidStatusTransition(app.CurrentStatus, req.Status) {
		return nil, apperrors.NewBadRequestError("Invalid status transition")
	}

	// Update application status
	app.CurrentStatus = req.Status

	if err := s.repo.Update(ctx, app); err != nil {
		return nil, apperrors.NewInternalError("Failed to update application", err)
	}

	// Add timeline event
	event := &TimelineEvent{
		ApplicationID:        app.ID,
		Status:               req.Status,
		IsVisibleToApplicant: true,
		UpdatedByType:        "company",
		UpdatedByID:          sql.NullInt64{Int64: int64(companyID), Valid: true},
	}

	if req.Note != "" {
		event.Note = sql.NullString{String: req.Note, Valid: true}
	}

	// Handle interview scheduling
	if req.Status == StatusInterviewScheduled && req.ScheduledAt != "" {
		scheduledTime, err := time.Parse(time.RFC3339, req.ScheduledAt)
		if err == nil {
			event.ScheduledAt = sql.NullTime{Time: scheduledTime, Valid: true}
		}
		if req.ScheduledLocation != "" {
			event.ScheduledLocation = sql.NullString{String: req.ScheduledLocation, Valid: true}
		}
		if req.ScheduledNotes != "" {
			event.ScheduledNotes = sql.NullString{String: req.ScheduledNotes, Valid: true}
		}
	}

	if err := s.repo.AddTimelineEvent(ctx, event); err != nil {
		return nil, apperrors.NewInternalError("Failed to add timeline event", err)
	}

	return s.GetByID(ctx, app.ID, companyID, true)
}

// Withdraw allows an applicant to withdraw their application
func (s *service) Withdraw(ctx context.Context, applicationID uint64, userID uint64, reason string) error {
	// Get application
	app, err := s.repo.GetByID(ctx, applicationID)
	if err != nil {
		return apperrors.NewInternalError("Failed to get application", err)
	}
	if app == nil {
		return apperrors.NewNotFoundError("Application")
	}

	// Verify ownership
	if app.UserID != userID {
		return apperrors.NewForbiddenError("You don't have permission to withdraw this application")
	}

	// Check if already terminal
	if IsTerminalStatus(app.CurrentStatus) {
		return apperrors.NewBadRequestError("Cannot withdraw an application that is already completed")
	}

	// Update status
	app.CurrentStatus = StatusWithdrawn

	if err := s.repo.Update(ctx, app); err != nil {
		return apperrors.NewInternalError("Failed to update application", err)
	}

	// Add timeline event
	note := "Lamaran dibatalkan oleh pelamar"
	if reason != "" {
		note = reason
	}

	event := &TimelineEvent{
		ApplicationID:        app.ID,
		Status:               StatusWithdrawn,
		Note:                 sql.NullString{String: note, Valid: true},
		IsVisibleToApplicant: true,
		UpdatedByType:        "applicant",
		UpdatedByID:          sql.NullInt64{Int64: int64(userID), Valid: true},
	}

	if err := s.repo.AddTimelineEvent(ctx, event); err != nil {
		return apperrors.NewInternalError("Failed to add timeline event", err)
	}

	return nil
}

// loadApplicationRelations loads related data for an application
func (s *service) loadApplicationRelations(ctx context.Context, app *Application, isCompany bool) error {
	// Load job info
	job, err := s.repo.GetJobInfo(ctx, app.JobID)
	if err != nil {
		return apperrors.NewInternalError("Failed to load job info", err)
	}
	app.Job = job

	// Load applicant info (only for company view)
	if isCompany {
		applicant, err := s.repo.GetApplicantInfo(ctx, app.UserID)
		if err != nil {
			return apperrors.NewInternalError("Failed to load applicant info", err)
		}
		app.Applicant = applicant
	}

	// Load CV snapshot info
	snapshot, err := s.repo.GetCVSnapshotInfo(ctx, app.CVSnapshotID)
	if err != nil {
		return apperrors.NewInternalError("Failed to load CV snapshot info", err)
	}
	app.CVSnapshot = snapshot

	// Load timeline
	var timeline []TimelineEvent
	if isCompany {
		timeline, err = s.repo.GetTimeline(ctx, app.ID)
	} else {
		timeline, err = s.repo.GetTimelineForApplicant(ctx, app.ID)
	}
	if err != nil {
		return apperrors.NewInternalError("Failed to load timeline", err)
	}
	app.Timeline = timeline

	return nil
}
