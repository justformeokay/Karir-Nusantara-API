package jobs

import (
	"database/sql"
	"time"

	"github.com/karirnusantara/api/internal/shared/hashid"
)

// Job types
const (
	JobTypeFullTime   = "full_time"
	JobTypePartTime   = "part_time"
	JobTypeContract   = "contract"
	JobTypeInternship = "internship"
	JobTypeFreelance  = "freelance"
)

// Experience levels
const (
	ExpLevelEntry     = "entry"
	ExpLevelJunior    = "junior"
	ExpLevelMid       = "mid"
	ExpLevelSenior    = "senior"
	ExpLevelLead      = "lead"
	ExpLevelExecutive = "executive"
)

// Job statuses
const (
	JobStatusDraft  = "draft"
	JobStatusActive = "active"
	JobStatusPaused = "paused"
	JobStatusClosed = "closed"
	JobStatusFilled = "filled"
)

// Job represents a job posting
type Job struct {
	ID                  uint64         `db:"id" json:"id"`
	CompanyID           uint64         `db:"company_id" json:"company_id"`
	Title               string         `db:"title" json:"title"`
	Category            string         `db:"category" json:"category"`
	Slug                string         `db:"slug" json:"slug"`
	Description         string         `db:"description" json:"description"`
	Requirements        sql.NullString `db:"requirements" json:"requirements,omitempty"`
	Responsibilities    sql.NullString `db:"responsibilities" json:"responsibilities,omitempty"`
	Benefits            sql.NullString `db:"benefits" json:"benefits,omitempty"`
	City                string         `db:"city" json:"city"`
	Province            string         `db:"province" json:"province"`
	IsRemote            bool           `db:"is_remote" json:"is_remote"`
	JobType             string         `db:"job_type" json:"job_type"`
	ExperienceLevel     string         `db:"experience_level" json:"experience_level"`
	SalaryMin           sql.NullInt64  `db:"salary_min" json:"salary_min,omitempty"`
	SalaryMax           sql.NullInt64  `db:"salary_max" json:"salary_max,omitempty"`
	SalaryCurrency      string         `db:"salary_currency" json:"salary_currency"`
	IsSalaryVisible     bool           `db:"is_salary_visible" json:"is_salary_visible"`
	IsSalaryFixed       bool           `db:"is_salary_fixed" json:"is_salary_fixed"`
	ApplicationDeadline sql.NullTime   `db:"application_deadline" json:"application_deadline,omitempty"`
	MaxApplications     sql.NullInt64  `db:"max_applications" json:"max_applications,omitempty"`
	Status              string         `db:"status" json:"status"`
	ViewsCount          uint64         `db:"views_count" json:"views_count"`
	ApplicationsCount   uint64         `db:"applications_count" json:"applications_count"`
	SharesCount         uint64         `db:"shares_count" json:"shares_count"`
	PublishedAt         sql.NullTime   `db:"published_at" json:"published_at,omitempty"`
	CreatedAt           time.Time      `db:"created_at" json:"created_at"`
	UpdatedAt           time.Time      `db:"updated_at" json:"updated_at"`
	DeletedAt           sql.NullTime   `db:"deleted_at" json:"-"`

	// Relationships (loaded separately)
	Company *CompanyInfo `db:"-" json:"company,omitempty"`
	Skills  []JobSkill   `db:"-" json:"skills,omitempty"`
}

// CompanyInfo represents minimal company information for job listing
type CompanyInfo struct {
	ID       uint64 `json:"id"`
	HashID   string `json:"hash_id,omitempty"`
	Name     string `json:"name"`
	LogoURL  string `json:"logo_url,omitempty"`
	Website  string `json:"website,omitempty"`
	City     string `json:"city,omitempty"`
	Province string `json:"province,omitempty"`
}

// WithHashID adds hash_id to CompanyInfo
func (c *CompanyInfo) WithHashID() *CompanyInfo {
	if c != nil {
		c.HashID = hashid.Encode(c.ID)
	}
	return c
}

// JobSkill represents a required skill for a job
type JobSkill struct {
	ID         uint64 `db:"id" json:"id"`
	JobID      uint64 `db:"job_id" json:"job_id"`
	SkillName  string `db:"skill_name" json:"skill_name"`
	IsRequired bool   `db:"is_required" json:"is_required"`
}

// JobResponse represents the job response for API
type JobResponse struct {
	ID               uint64       `json:"id"`
	HashID           string       `json:"hash_id"`
	Title            string       `json:"title"`
	Category         string       `json:"category"`
	Slug             string       `json:"slug"`
	Description      string       `json:"description"`
	Requirements     string       `json:"requirements,omitempty"`
	Responsibilities string       `json:"responsibilities,omitempty"`
	Benefits         string       `json:"benefits,omitempty"`
	Location         LocationInfo `json:"location"`
	JobType          string       `json:"job_type"`
	ExperienceLevel  string       `json:"experience_level"`
	Salary           *SalaryInfo  `json:"salary,omitempty"`
	// Raw salary fields for company portal editing
	SalaryMin           *int64       `json:"salary_min,omitempty"`
	SalaryMax           *int64       `json:"salary_max,omitempty"`
	SalaryCurrency      string       `json:"salary_currency,omitempty"`
	IsSalaryVisible     bool         `json:"is_salary_visible"`
	IsSalaryFixed       bool         `json:"is_salary_fixed"`
	ApplicationDeadline string       `json:"application_deadline,omitempty"`
	Status              string       `json:"status"`
	ViewsCount          uint64       `json:"views_count"`
	ApplicationsCount   uint64       `json:"applications_count"`
	SharesCount         uint64       `json:"shares_count"`
	PublishedAt         string       `json:"published_at,omitempty"`
	CreatedAt           string       `json:"created_at"`
	Company             *CompanyInfo `json:"company,omitempty"`
	Skills              []string     `json:"skills,omitempty"`
}

// LocationInfo represents job location
type LocationInfo struct {
	City     string `json:"city"`
	Province string `json:"province"`
	IsRemote bool   `json:"is_remote"`
}

// SalaryInfo represents salary information
type SalaryInfo struct {
	Min      int64  `json:"min,omitempty"`
	Max      int64  `json:"max,omitempty"`
	Currency string `json:"currency"`
}

// ToResponse converts Job to JobResponse
func (j *Job) ToResponse() *JobResponse {
	resp := &JobResponse{
		ID:          j.ID,
		HashID:      hashid.Encode(j.ID),
		Title:       j.Title,
		Category:    j.Category,
		Slug:        j.Slug,
		Description: j.Description,
		Location: LocationInfo{
			City:     j.City,
			Province: j.Province,
			IsRemote: j.IsRemote,
		},
		JobType:           j.JobType,
		ExperienceLevel:   j.ExperienceLevel,
		Status:            j.Status,
		ViewsCount:        j.ViewsCount,
		ApplicationsCount: j.ApplicationsCount,
		SharesCount:       j.SharesCount,
		CreatedAt:         j.CreatedAt.Format(time.RFC3339),
		Company:           j.Company.WithHashID(),
		// Always include raw salary fields for company portal
		SalaryCurrency:  j.SalaryCurrency,
		IsSalaryVisible: j.IsSalaryVisible,
		IsSalaryFixed:   j.IsSalaryFixed,
	}

	// Always populate raw salary fields if they exist
	if j.SalaryMin.Valid {
		resp.SalaryMin = &j.SalaryMin.Int64
	}
	if j.SalaryMax.Valid {
		resp.SalaryMax = &j.SalaryMax.Int64
	}

	if j.Requirements.Valid {
		resp.Requirements = j.Requirements.String
	}
	if j.Responsibilities.Valid {
		resp.Responsibilities = j.Responsibilities.String
	}
	if j.Benefits.Valid {
		resp.Benefits = j.Benefits.String
	}
	if j.ApplicationDeadline.Valid {
		resp.ApplicationDeadline = j.ApplicationDeadline.Time.Format("2006-01-02")
	}
	if j.PublishedAt.Valid {
		resp.PublishedAt = j.PublishedAt.Time.Format(time.RFC3339)
	}

	// Include salary object if visible (for public display)
	if j.IsSalaryVisible && j.SalaryMin.Valid {
		resp.Salary = &SalaryInfo{
			Currency: j.SalaryCurrency,
		}
		if j.IsSalaryFixed {
			// For fixed salary, only Min is used as the fixed amount
			resp.Salary.Min = j.SalaryMin.Int64
			resp.Salary.Max = 0 // Max is 0 for fixed salary
		} else {
			// For range salary
			resp.Salary.Min = j.SalaryMin.Int64
			if j.SalaryMax.Valid {
				resp.Salary.Max = j.SalaryMax.Int64
			}
		}
	}

	// Convert skills
	if len(j.Skills) > 0 {
		resp.Skills = make([]string, len(j.Skills))
		for i, skill := range j.Skills {
			resp.Skills[i] = skill.SkillName
		}
	}

	return resp
}

// Request DTOs

// CreateJobRequest represents a job creation request
type CreateJobRequest struct {
	Title               string   `json:"title" validate:"required,min=5,max=255"`
	Category            string   `json:"category" validate:"required,max=50"`
	Description         string   `json:"description" validate:"required,min=50"`
	Requirements        string   `json:"requirements,omitempty"`
	Responsibilities    string   `json:"responsibilities,omitempty"`
	Benefits            string   `json:"benefits,omitempty"`
	City                string   `json:"city" validate:"required"`
	Province            string   `json:"province" validate:"required"`
	IsRemote            bool     `json:"is_remote"`
	JobType             string   `json:"job_type" validate:"required,oneof=full_time part_time contract internship freelance"`
	ExperienceLevel     string   `json:"experience_level" validate:"required,oneof=entry junior mid senior lead executive"`
	SalaryMin           *int64   `json:"salary_min,omitempty" validate:"omitempty,gte=0"`
	SalaryMax           *int64   `json:"salary_max,omitempty" validate:"omitempty,gtefield=SalaryMin"`
	SalaryCurrency      string   `json:"salary_currency,omitempty"`
	IsSalaryVisible     bool     `json:"is_salary_visible"`
	IsSalaryFixed       bool     `json:"is_salary_fixed"` // If true, only salary_min is used
	ApplicationDeadline string   `json:"application_deadline,omitempty"`
	Skills              []string `json:"skills,omitempty"`
	Status              string   `json:"status,omitempty" validate:"omitempty,oneof=draft active"`
}

// UpdateJobRequest represents a job update request
type UpdateJobRequest struct {
	Title               *string  `json:"title,omitempty" validate:"omitempty,min=5,max=255"`
	Category            *string  `json:"category,omitempty" validate:"omitempty,max=50"`
	Description         *string  `json:"description,omitempty" validate:"omitempty,min=50"`
	Requirements        *string  `json:"requirements,omitempty"`
	Responsibilities    *string  `json:"responsibilities,omitempty"`
	Benefits            *string  `json:"benefits,omitempty"`
	City                *string  `json:"city,omitempty"`
	Province            *string  `json:"province,omitempty"`
	IsRemote            *bool    `json:"is_remote,omitempty"`
	JobType             *string  `json:"job_type,omitempty" validate:"omitempty,oneof=full_time part_time contract internship freelance"`
	ExperienceLevel     *string  `json:"experience_level,omitempty" validate:"omitempty,oneof=entry junior mid senior lead executive"`
	SalaryMin           *int64   `json:"salary_min,omitempty"`
	SalaryMax           *int64   `json:"salary_max,omitempty"`
	IsSalaryVisible     *bool    `json:"is_salary_visible,omitempty"`
	IsSalaryFixed       *bool    `json:"is_salary_fixed,omitempty"` // If true, only salary_min is used
	ApplicationDeadline *string  `json:"application_deadline,omitempty"`
	Skills              []string `json:"skills,omitempty"`
	Status              *string  `json:"status,omitempty" validate:"omitempty,oneof=draft active paused closed filled"`
}

// JobListParams represents job list query parameters
type JobListParams struct {
	Page            int      `json:"page"`
	PerPage         int      `json:"per_page"`
	Search          string   `json:"search"`
	City            string   `json:"city"`
	Province        string   `json:"province"`
	JobType         string   `json:"job_type"`
	ExperienceLevel string   `json:"experience_level"`
	IsRemote        *bool    `json:"is_remote"`
	SalaryMin       *int64   `json:"salary_min"`
	SalaryMax       *int64   `json:"salary_max"`
	Skills          []string `json:"skills"`
	CompanyID       *uint64  `json:"company_id"`
	Status          string   `json:"status"`
	SortBy          string   `json:"sort_by"`
	SortOrder       string   `json:"sort_order"`
}

// DefaultJobListParams returns default list parameters
func DefaultJobListParams() JobListParams {
	return JobListParams{
		Page:      1,
		PerPage:   20,
		Status:    JobStatusActive,
		SortBy:    "published_at",
		SortOrder: "desc",
	}
}

// JobView represents a unique job view record
type JobView struct {
	ID       uint64    `db:"id" json:"id"`
	JobID    uint64    `db:"job_id" json:"job_id"`
	UserID   uint64    `db:"user_id" json:"user_id"`
	ViewedAt time.Time `db:"viewed_at" json:"viewed_at"`
}

// JobShare represents a job share record
type JobShare struct {
	ID       uint64    `db:"id" json:"id"`
	JobID    uint64    `db:"job_id" json:"job_id"`
	UserID   uint64    `db:"user_id" json:"user_id,omitempty"`
	Platform string    `db:"platform" json:"platform,omitempty"`
	SharedAt time.Time `db:"shared_at" json:"shared_at"`
}

// TrackShareRequest represents a share tracking request
type TrackShareRequest struct {
	Platform string `json:"platform,omitempty"` // whatsapp, telegram, facebook, twitter, copy_link, etc
}

// JobStatsResponse represents job statistics
type JobStatsResponse struct {
	JobID             uint64 `json:"job_id"`
	HashID            string `json:"hash_id"`
	Title             string `json:"title"`
	ViewsCount        uint64 `json:"views_count"`
	ApplicationsCount uint64 `json:"applications_count"`
	SharesCount       uint64 `json:"shares_count"`
}
