package wishlist

import (
	"time"
)

// SavedJob represents a job saved by a user (wishlist item)
type SavedJob struct {
	ID        uint64    `json:"id" db:"id"`
	UserID    uint64    `json:"user_id" db:"user_id"`
	JobID     uint64    `json:"job_id" db:"job_id"`
	CreatedAt time.Time `json:"saved_at" db:"created_at"`

	// Populated from joins
	Job *JobInfo `json:"job,omitempty"`
}

// JobInfo represents job details for wishlist
type JobInfo struct {
	ID              uint64    `json:"id" db:"id"`
	Title           string    `json:"title" db:"title"`
	Slug            string    `json:"slug" db:"slug"`
	City            string    `json:"city" db:"city"`
	Province        string    `json:"province" db:"province"`
	IsRemote        bool      `json:"is_remote" db:"is_remote"`
	JobType         string    `json:"job_type" db:"job_type"`
	ExperienceLevel string    `json:"experience_level" db:"experience_level"`
	SalaryMin       *int64    `json:"salary_min,omitempty" db:"salary_min"`
	SalaryMax       *int64    `json:"salary_max,omitempty" db:"salary_max"`
	Status          string    `json:"status" db:"status"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`

	// Company info
	Company *CompanyInfo `json:"company,omitempty"`
}

// CompanyInfo represents company details
type CompanyInfo struct {
	ID      uint64  `json:"id" db:"id"`
	Name    string  `json:"name" db:"company_name"`
	LogoURL *string `json:"logo_url,omitempty" db:"company_logo_url"`
}

// SaveJobRequest represents the request to save a job
type SaveJobRequest struct {
	JobID uint64 `json:"job_id" validate:"required"`
}

// ListParams represents pagination parameters
type ListParams struct {
	Page    int `json:"page"`
	PerPage int `json:"per_page"`
}

// WishlistResponse represents a saved job with full details
type WishlistResponse struct {
	ID      uint64    `json:"id"`
	JobID   uint64    `json:"job_id"`
	SavedAt time.Time `json:"saved_at"`
	Job     *JobInfo  `json:"job"`
}
