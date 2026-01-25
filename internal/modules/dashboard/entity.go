package dashboard

// DashboardStats represents company dashboard statistics
type DashboardStats struct {
	ActiveJobs         int               `json:"active_jobs"`
	TotalApplicants    int               `json:"total_applicants"`
	UnderReview        int               `json:"under_review"`
	AcceptedCandidates int               `json:"accepted_candidates"`
	RemainingFreeQuota int               `json:"remaining_free_quota"`
	PendingPayments    int               `json:"pending_payments"`
	RecentApplicants   []RecentApplicant `json:"recent_applicants"`
	ActiveJobsList     []ActiveJob       `json:"active_jobs_list"`
}

// RecentApplicant represents a recent application summary
type RecentApplicant struct {
	ID             uint64 `json:"id"`
	HashID         string `json:"hash_id"`
	ApplicantName  string `json:"applicant_name"`
	ApplicantPhoto string `json:"applicant_photo,omitempty"`
	JobID          uint64 `json:"job_id"`
	JobTitle       string `json:"job_title"`
	Status         string `json:"status"`
	StatusLabel    string `json:"status_label"`
	AppliedAt      string `json:"applied_at"`
}

// ActiveJob represents an active job summary for dashboard
type ActiveJob struct {
	ID              uint64 `json:"id"`
	HashID          string `json:"hash_id"`
	Title           string `json:"title"`
	Status          string `json:"status"`
	ApplicantsCount int    `json:"applicants_count"`
	ViewsCount      int    `json:"views_count"`
	CreatedAt       string `json:"created_at"`
}
