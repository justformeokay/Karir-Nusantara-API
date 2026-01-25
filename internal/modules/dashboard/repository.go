package dashboard

import (
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/karirnusantara/api/internal/shared/hashid"
)

// Repository handles database operations for dashboard
type Repository struct {
	db *sqlx.DB
}

// NewRepository creates a new dashboard repository
func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

// GetActiveJobsCount returns the count of active jobs for a company
func (r *Repository) GetActiveJobsCount(companyID uint64) (int, error) {
	var count int
	err := r.db.Get(&count, `
		SELECT COUNT(*) FROM jobs 
		WHERE company_id = ? AND status = 'active' AND deleted_at IS NULL
	`, companyID)
	return count, err
}

// GetTotalApplicantsCount returns total applicants for all company jobs
func (r *Repository) GetTotalApplicantsCount(companyID uint64) (int, error) {
	var count int
	err := r.db.Get(&count, `
		SELECT COUNT(*) FROM applications a
		JOIN jobs j ON a.job_id = j.id
		WHERE j.company_id = ? AND j.deleted_at IS NULL
	`, companyID)
	return count, err
}

// GetUnderReviewCount returns applicants under review
func (r *Repository) GetUnderReviewCount(companyID uint64) (int, error) {
	var count int
	err := r.db.Get(&count, `
		SELECT COUNT(*) FROM applications a
		JOIN jobs j ON a.job_id = j.id
		WHERE j.company_id = ? AND j.deleted_at IS NULL
		AND a.current_status IN ('submitted', 'viewed', 'shortlisted')
	`, companyID)
	return count, err
}

// GetAcceptedCandidatesCount returns hired/accepted candidates
func (r *Repository) GetAcceptedCandidatesCount(companyID uint64) (int, error) {
	var count int
	err := r.db.Get(&count, `
		SELECT COUNT(*) FROM applications a
		JOIN jobs j ON a.job_id = j.id
		WHERE j.company_id = ? AND j.deleted_at IS NULL
		AND a.current_status IN ('offer_accepted', 'hired')
	`, companyID)
	return count, err
}

// recentApplicantRow is used for scanning database results
type recentApplicantRow struct {
	ID             uint64    `db:"id"`
	ApplicantName  string    `db:"applicant_name"`
	ApplicantPhoto string    `db:"applicant_photo"`
	JobID          uint64    `db:"job_id"`
	JobTitle       string    `db:"job_title"`
	Status         string    `db:"status"`
	AppliedAt      time.Time `db:"applied_at"`
}

// GetRecentApplicants returns recent applicants for company jobs
func (r *Repository) GetRecentApplicants(companyID uint64, limit int) ([]RecentApplicant, error) {
	var rows []recentApplicantRow
	err := r.db.Select(&rows, `
		SELECT 
			a.id,
			u.full_name as applicant_name,
			COALESCE(u.avatar_url, '') as applicant_photo,
			j.id as job_id,
			j.title as job_title,
			a.current_status as status,
			a.applied_at
		FROM applications a
		JOIN jobs j ON a.job_id = j.id
		JOIN users u ON a.user_id = u.id
		WHERE j.company_id = ? AND j.deleted_at IS NULL
		ORDER BY a.applied_at DESC
		LIMIT ?
	`, companyID, limit)
	if err != nil {
		return nil, err
	}

	applicants := make([]RecentApplicant, len(rows))
	for i, row := range rows {
		applicants[i] = RecentApplicant{
			ID:             row.ID,
			HashID:         hashid.Encode(row.ID),
			ApplicantName:  row.ApplicantName,
			ApplicantPhoto: row.ApplicantPhoto,
			JobID:          row.JobID,
			JobTitle:       row.JobTitle,
			Status:         row.Status,
			StatusLabel:    getStatusLabel(row.Status),
			AppliedAt:      row.AppliedAt.Format(time.RFC3339),
		}
	}

	return applicants, nil
}

// activeJobRow is used for scanning database results
type activeJobRow struct {
	ID              uint64    `db:"id"`
	Title           string    `db:"title"`
	Status          string    `db:"status"`
	ApplicantsCount int       `db:"applicants_count"`
	ViewsCount      int       `db:"views_count"`
	CreatedAt       time.Time `db:"created_at"`
}

// GetActiveJobs returns active jobs for dashboard
func (r *Repository) GetActiveJobs(companyID uint64, limit int) ([]ActiveJob, error) {
	var rows []activeJobRow
	err := r.db.Select(&rows, `
		SELECT 
			id,
			title,
			status,
			applications_count as applicants_count,
			views_count,
			created_at
		FROM jobs 
		WHERE company_id = ? AND status = 'active' AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT ?
	`, companyID, limit)
	if err != nil {
		return nil, err
	}

	jobs := make([]ActiveJob, len(rows))
	for i, row := range rows {
		jobs[i] = ActiveJob{
			ID:              row.ID,
			HashID:          hashid.Encode(row.ID),
			Title:           row.Title,
			Status:          row.Status,
			ApplicantsCount: row.ApplicantsCount,
			ViewsCount:      row.ViewsCount,
			CreatedAt:       row.CreatedAt.Format(time.RFC3339),
		}
	}

	return jobs, nil
}

func getStatusLabel(status string) string {
	labels := map[string]string{
		"submitted":           "Baru Melamar",
		"viewed":              "Sedang Ditinjau",
		"shortlisted":         "Masuk Shortlist",
		"interview_scheduled": "Interview Dijadwalkan",
		"interview_completed": "Interview Selesai",
		"assessment":          "Tahap Assessment",
		"offer_sent":          "Penawaran Dikirim",
		"offer_accepted":      "Penawaran Diterima",
		"hired":               "Diterima",
		"rejected":            "Ditolak",
		"withdrawn":           "Dibatalkan",
	}
	if label, ok := labels[status]; ok {
		return label
	}
	return status
}

// GetCompanyIDByUserID retrieves company ID from user ID
func (r *Repository) GetCompanyIDByUserID(userID uint64) (uint64, error) {
	var companyID uint64
	err := r.db.Get(&companyID, `
		SELECT id FROM companies 
		WHERE user_id = ? AND deleted_at IS NULL
	`, userID)
	return companyID, err
}
