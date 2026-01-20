package dashboard

// Service handles business logic for dashboard
type Service struct {
	repo *Repository
}

// NewService creates a new dashboard service
func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// GetDashboardStats returns all dashboard statistics for a company
func (s *Service) GetDashboardStats(companyID uint64) (*DashboardStats, error) {
	stats := &DashboardStats{}
	var err error

	// Get active jobs count
	stats.ActiveJobs, err = s.repo.GetActiveJobsCount(companyID)
	if err != nil {
		return nil, err
	}

	// Get total applicants
	stats.TotalApplicants, err = s.repo.GetTotalApplicantsCount(companyID)
	if err != nil {
		return nil, err
	}

	// Get under review count
	stats.UnderReview, err = s.repo.GetUnderReviewCount(companyID)
	if err != nil {
		return nil, err
	}

	// Get accepted candidates
	stats.AcceptedCandidates, err = s.repo.GetAcceptedCandidatesCount(companyID)
	if err != nil {
		return nil, err
	}

	// Get recent applicants (limit 5)
	stats.RecentApplicants, err = s.repo.GetRecentApplicants(companyID, 5)
	if err != nil {
		return nil, err
	}

	// Get active jobs list (limit 5)
	stats.ActiveJobsList, err = s.repo.GetActiveJobs(companyID, 5)
	if err != nil {
		return nil, err
	}

	return stats, nil
}

// GetCompanyIDByUserID retrieves company ID for a given user ID
func (s *Service) GetCompanyIDByUserID(userID uint64) (uint64, error) {
	return s.repo.GetCompanyIDByUserID(userID)
}

// GetRecentApplicants returns paginated recent applicants
func (s *Service) GetRecentApplicants(companyID uint64, limit int) ([]RecentApplicant, error) {
	if limit <= 0 {
		limit = 10
	}
	return s.repo.GetRecentApplicants(companyID, limit)
}

// GetActiveJobsList returns paginated active jobs
func (s *Service) GetActiveJobsList(companyID uint64, limit int) ([]ActiveJob, error) {
	if limit <= 0 {
		limit = 10
	}
	return s.repo.GetActiveJobs(companyID, limit)
}
