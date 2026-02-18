package jobreports

import (
	"context"
	"fmt"
	"time"
)

// Service interface for job reports
type Service interface {
	CreateReport(ctx context.Context, userID, jobID uint64, req *CreateReportRequest) (*JobReportWithDetails, error)
	GetReportByID(ctx context.Context, id uint64) (*JobReportWithDetails, error)
	GetAllReports(ctx context.Context, status string, page, limit int) ([]JobReportWithDetails, int, error)
	GetReportsByJobID(ctx context.Context, jobID uint64) ([]JobReportWithDetails, error)
	UpdateReportStatus(ctx context.Context, id uint64, req *UpdateReportStatusRequest, adminID uint64) error
	HasUserReportedJob(ctx context.Context, userID, jobID uint64) (bool, error)
	GetPendingReportsCount(ctx context.Context) (int, error)
	GetReportCountByCompanyID(ctx context.Context, companyID uint64) (int, error)
	BanCompany(ctx context.Context, companyID uint64, adminID uint64, reason string) error
}

type service struct {
	repo *Repository
}

// NewService creates a new job reports service
func NewService(repo *Repository) Service {
	return &service{repo: repo}
}

// CreateReport creates a new job report
func (s *service) CreateReport(ctx context.Context, userID, jobID uint64, req *CreateReportRequest) (*JobReportWithDetails, error) {
	// Check if user has already reported this job
	hasReported, err := s.repo.HasUserReportedJob(ctx, userID, jobID)
	if err != nil {
		return nil, err
	}
	if hasReported {
		return nil, fmt.Errorf("Anda sudah pernah melaporkan lowongan ini")
	}

	// Check cooldown - only 1 report per 24 hours regardless of job
	lastReportTime, err := s.repo.GetLastReportTime(ctx, userID)
	if err != nil {
		return nil, err
	}
	if lastReportTime != nil {
		timeSinceLastReport := time.Since(*lastReportTime)
		if timeSinceLastReport < CooldownDuration {
			remainingTime := CooldownDuration - timeSinceLastReport
			hours := int(remainingTime.Hours())
			minutes := int(remainingTime.Minutes()) % 60
			return nil, fmt.Errorf("Anda dapat membuat laporan baru dalam %d jam %d menit lagi", hours, minutes)
		}
	}

	report := &JobReport{
		JobID:       jobID,
		UserID:      userID,
		Reason:      req.Reason,
		Description: req.Description,
		Status:      ReportStatusPending,
	}

	if err := s.repo.Create(ctx, report); err != nil {
		return nil, err
	}

	// Fetch the complete report with details
	return s.repo.GetByID(ctx, report.ID)
}

// GetReportByID gets a job report by ID
func (s *service) GetReportByID(ctx context.Context, id uint64) (*JobReportWithDetails, error) {
	return s.repo.GetByID(ctx, id)
}

// GetAllReports gets all reports with optional status filter
func (s *service) GetAllReports(ctx context.Context, status string, page, limit int) ([]JobReportWithDetails, int, error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	return s.repo.GetAll(ctx, status, page, limit)
}

// GetReportsByJobID gets all reports for a specific job
func (s *service) GetReportsByJobID(ctx context.Context, jobID uint64) ([]JobReportWithDetails, error) {
	return s.repo.GetByJobID(ctx, jobID)
}

// UpdateReportStatus updates the status of a report
func (s *service) UpdateReportStatus(ctx context.Context, id uint64, req *UpdateReportStatusRequest, adminID uint64) error {
	// Verify report exists
	report, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if report == nil {
		return fmt.Errorf("laporan tidak ditemukan")
	}

	return s.repo.UpdateStatus(ctx, id, req.Status, req.AdminNote, adminID)
}

// HasUserReportedJob checks if user has already reported a job
func (s *service) HasUserReportedJob(ctx context.Context, userID, jobID uint64) (bool, error) {
	return s.repo.HasUserReportedJob(ctx, userID, jobID)
}

// GetPendingReportsCount gets the count of pending reports
func (s *service) GetPendingReportsCount(ctx context.Context) (int, error) {
	return s.repo.GetPendingReportsCount(ctx)
}

// GetReportCountByCompanyID gets the number of reports for all jobs from a company
func (s *service) GetReportCountByCompanyID(ctx context.Context, companyID uint64) (int, error) {
	return s.repo.GetReportCountByCompanyID(ctx, companyID)
}

// BanCompany bans a company with reason
func (s *service) BanCompany(ctx context.Context, companyID uint64, adminID uint64, reason string) error {
	return s.repo.BanCompany(ctx, companyID, adminID, reason)
}
