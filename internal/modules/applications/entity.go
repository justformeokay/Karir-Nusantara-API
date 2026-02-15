package applications

import (
	"database/sql"
	"time"

	"github.com/karirnusantara/api/internal/shared/hashid"
)

// Application statuses
const (
	StatusSubmitted          = "submitted"
	StatusViewed             = "viewed"
	StatusShortlisted        = "shortlisted"
	StatusInterviewScheduled = "interview_scheduled"
	StatusInterviewCompleted = "interview_completed"
	StatusAssessment         = "assessment"
	StatusOfferSent          = "offer_sent"
	StatusOfferAccepted      = "offer_accepted"
	StatusHired              = "hired"
	StatusRejected           = "rejected"
	StatusWithdrawn          = "withdrawn"
)

// Application represents a job application
type Application struct {
	ID                 uint64         `db:"id" json:"id"`
	UserID             uint64         `db:"user_id" json:"user_id"`
	JobID              uint64         `db:"job_id" json:"job_id"`
	CVSnapshotID       uint64         `db:"cv_snapshot_id" json:"cv_snapshot_id"`
	CVSource           string         `db:"cv_source" json:"cv_source"`                                 // 'built' or 'uploaded'
	UploadedDocumentID sql.NullInt64  `db:"uploaded_document_id" json:"uploaded_document_id,omitempty"` // Reference to document if cv_source='uploaded'
	CoverLetter        sql.NullString `db:"cover_letter" json:"cover_letter,omitempty"`
	CurrentStatus      string         `db:"current_status" json:"current_status"`
	AppliedAt          time.Time      `db:"applied_at" json:"applied_at"`
	LastStatusUpdate   time.Time      `db:"last_status_update" json:"last_status_update"`
	CreatedAt          time.Time      `db:"created_at" json:"created_at"`
	UpdatedAt          time.Time      `db:"updated_at" json:"updated_at"`

	// Relationships (loaded separately)
	Job        *JobInfo        `db:"-" json:"job,omitempty"`
	Applicant  *ApplicantInfo  `db:"-" json:"applicant,omitempty"`
	CVSnapshot *CVSnapshotInfo `db:"-" json:"cv_snapshot,omitempty"`
	Timeline   []TimelineEvent `db:"-" json:"timeline,omitempty"`
}

// JobInfo represents minimal job information for application
type JobInfo struct {
	ID       uint64      `json:"id"`
	HashID   string      `json:"hash_id"`
	Title    string      `json:"title"`
	Company  CompanyInfo `json:"company"`
	City     string      `json:"city"`
	Province string      `json:"province"`
	Status   string      `json:"status"`
}

// CompanyInfo represents minimal company information
type CompanyInfo struct {
	ID      uint64 `json:"id"`
	HashID  string `json:"hash_id"`
	Name    string `json:"name"`
	LogoURL string `json:"logo_url,omitempty"`
}

// ApplicantInfo represents minimal applicant information
type ApplicantInfo struct {
	ID        uint64 `json:"id"`
	HashID    string `json:"hash_id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Phone     string `json:"phone,omitempty"`
	AvatarURL string `json:"avatar_url,omitempty"`
}

// CVSnapshotInfo represents CV snapshot information
type CVSnapshotInfo struct {
	ID                uint64                   `json:"id"`
	CompletenessScore int                      `json:"completeness_score"`
	PersonalInfo      map[string]interface{}   `json:"personal_info,omitempty"`
	Education         []map[string]interface{} `json:"education,omitempty"`
	Experience        []map[string]interface{} `json:"experience,omitempty"`
	Skills            []map[string]interface{} `json:"skills,omitempty"`
	Certifications    []map[string]interface{} `json:"certifications,omitempty"`
	Languages         []map[string]interface{} `json:"languages,omitempty"`
	Projects          []map[string]interface{} `json:"projects,omitempty"`
	CreatedAt         string                   `json:"created_at"`
}

// TimelineEvent represents an application timeline event
type TimelineEvent struct {
	ID                   uint64         `db:"id" json:"id"`
	ApplicationID        uint64         `db:"application_id" json:"application_id"`
	Status               string         `db:"status" json:"status"`
	Note                 sql.NullString `db:"note" json:"note,omitempty"`
	IsVisibleToApplicant bool           `db:"is_visible_to_applicant" json:"is_visible_to_applicant"`
	UpdatedByType        string         `db:"updated_by_type" json:"updated_by_type"`
	UpdatedByID          sql.NullInt64  `db:"updated_by_id" json:"updated_by_id,omitempty"`
	ScheduledAt          sql.NullTime   `db:"scheduled_at" json:"scheduled_at,omitempty"`
	ScheduledLocation    sql.NullString `db:"scheduled_location" json:"scheduled_location,omitempty"`
	ScheduledNotes       sql.NullString `db:"scheduled_notes" json:"scheduled_notes,omitempty"`
	InterviewType        sql.NullString `db:"interview_type" json:"interview_type,omitempty"`
	MeetingLink          sql.NullString `db:"meeting_link" json:"meeting_link,omitempty"`
	MeetingPlatform      sql.NullString `db:"meeting_platform" json:"meeting_platform,omitempty"`
	InterviewAddress     sql.NullString `db:"interview_address" json:"interview_address,omitempty"`
	ContactPerson        sql.NullString `db:"contact_person" json:"contact_person,omitempty"`
	ContactPhone         sql.NullString `db:"contact_phone" json:"contact_phone,omitempty"`
	CreatedAt            time.Time      `db:"created_at" json:"created_at"`
}

// Request DTOs

// ApplyJobRequest represents a job application request
type ApplyJobRequest struct {
	JobID              uint64 `json:"job_id" validate:"required"`
	CoverLetter        string `json:"cover_letter,omitempty"`
	CVSource           string `json:"cv_source,omitempty" validate:"omitempty,oneof=built uploaded"`                     // Optional: 'built' or 'uploaded', defaults to 'built'
	UploadedDocumentID uint64 `json:"uploaded_document_id,omitempty" validate:"omitempty,required_if=CVSource uploaded"` // Required if cv_source='uploaded'
}

// UpdateStatusRequest represents a status update request (by company)
type UpdateStatusRequest struct {
	Status            string `json:"status" validate:"required,oneof=viewed shortlisted interview_scheduled interview_completed assessment offer_sent offer_accepted hired rejected submitted"`
	Note              string `json:"note,omitempty"`
	ScheduledAt       string `json:"scheduled_at,omitempty"`
	ScheduledLocation string `json:"scheduled_location,omitempty"`
	ScheduledNotes    string `json:"scheduled_notes,omitempty"`
	InterviewType     string `json:"interview_type,omitempty"`    // online, offline, whatsapp_notification
	MeetingLink       string `json:"meeting_link,omitempty"`      // Zoom/Google Meet link for online
	MeetingPlatform   string `json:"meeting_platform,omitempty"`  // zoom, google_meet, microsoft_teams, etc.
	InterviewAddress  string `json:"interview_address,omitempty"` // Full address for offline interview
	ContactPerson     string `json:"contact_person,omitempty"`    // Contact person name
	ContactPhone      string `json:"contact_phone,omitempty"`     // Contact phone number
}

// WithdrawRequest represents a withdrawal request
type WithdrawRequest struct {
	Reason string `json:"reason,omitempty"`
}

// ApplicationListParams represents application list parameters
type ApplicationListParams struct {
	Page      int
	PerPage   int
	UserID    *uint64
	JobID     *uint64
	CompanyID *uint64
	Status    string
	SortBy    string
	SortOrder string
}

// DefaultApplicationListParams returns default list parameters
func DefaultApplicationListParams() ApplicationListParams {
	return ApplicationListParams{
		Page:      1,
		PerPage:   20,
		SortBy:    "applied_at",
		SortOrder: "desc",
	}
}

// Response DTOs

// ApplicationResponse represents the application response
type ApplicationResponse struct {
	ID               uint64                  `json:"id"`
	HashID           string                  `json:"hash_id"`
	Job              *JobInfo                `json:"job"`
	Applicant        *ApplicantInfo          `json:"applicant,omitempty"`
	CVSnapshot       *CVSnapshotInfo         `json:"cv_snapshot,omitempty"`
	CoverLetter      string                  `json:"cover_letter,omitempty"`
	CurrentStatus    string                  `json:"current_status"`
	StatusLabel      string                  `json:"status_label"`
	AppliedAt        string                  `json:"applied_at"`
	LastStatusUpdate string                  `json:"last_status_update"`
	Timeline         []TimelineEventResponse `json:"timeline,omitempty"`
}

// TimelineEventResponse represents a timeline event response
// TimelineEventResponse represents a timeline event response
type TimelineEventResponse struct {
	ID                uint64 `json:"id"`
	Status            string `json:"status"`
	StatusLabel       string `json:"status_label"`
	Note              string `json:"note,omitempty"`
	ScheduledAt       string `json:"scheduled_at,omitempty"`
	ScheduledLocation string `json:"scheduled_location,omitempty"`
	ScheduledNotes    string `json:"scheduled_notes,omitempty"`
	InterviewType     string `json:"interview_type,omitempty"`
	MeetingLink       string `json:"meeting_link,omitempty"`
	MeetingPlatform   string `json:"meeting_platform,omitempty"`
	InterviewAddress  string `json:"interview_address,omitempty"`
	ContactPerson     string `json:"contact_person,omitempty"`
	ContactPhone      string `json:"contact_phone,omitempty"`
	CreatedAt         string `json:"created_at"`
}

// ToResponse converts Application to ApplicationResponse
func (a *Application) ToResponse() *ApplicationResponse {
	resp := &ApplicationResponse{
		ID:               a.ID,
		HashID:           hashid.Encode(a.ID),
		Job:              a.Job,
		Applicant:        a.Applicant,
		CVSnapshot:       a.CVSnapshot,
		CurrentStatus:    a.CurrentStatus,
		StatusLabel:      GetStatusLabel(a.CurrentStatus),
		AppliedAt:        a.AppliedAt.Format(time.RFC3339),
		LastStatusUpdate: a.LastStatusUpdate.Format(time.RFC3339),
	}

	// Add hash_id to Job
	if resp.Job != nil {
		resp.Job.HashID = hashid.Encode(resp.Job.ID)
		resp.Job.Company.HashID = hashid.Encode(resp.Job.Company.ID)
	}

	// Add hash_id to Applicant
	if resp.Applicant != nil {
		resp.Applicant.HashID = hashid.Encode(resp.Applicant.ID)
	}

	if a.CoverLetter.Valid {
		resp.CoverLetter = a.CoverLetter.String
	}

	// Convert timeline
	if len(a.Timeline) > 0 {
		resp.Timeline = make([]TimelineEventResponse, len(a.Timeline))
		for i, event := range a.Timeline {
			resp.Timeline[i] = TimelineEventResponse{
				ID:          event.ID,
				Status:      event.Status,
				StatusLabel: GetStatusLabel(event.Status),
				CreatedAt:   event.CreatedAt.Format(time.RFC3339),
			}
			if event.Note.Valid {
				resp.Timeline[i].Note = event.Note.String
			}
			if event.ScheduledAt.Valid {
				resp.Timeline[i].ScheduledAt = event.ScheduledAt.Time.Format(time.RFC3339)
			}
			if event.ScheduledLocation.Valid {
				resp.Timeline[i].ScheduledLocation = event.ScheduledLocation.String
			}
			if event.ScheduledNotes.Valid {
				resp.Timeline[i].ScheduledNotes = event.ScheduledNotes.String
			}
			// Add new interview scheduling fields
			if event.InterviewType.Valid {
				resp.Timeline[i].InterviewType = event.InterviewType.String
			}
			if event.MeetingLink.Valid {
				resp.Timeline[i].MeetingLink = event.MeetingLink.String
			}
			if event.MeetingPlatform.Valid {
				resp.Timeline[i].MeetingPlatform = event.MeetingPlatform.String
			}
			if event.InterviewAddress.Valid {
				resp.Timeline[i].InterviewAddress = event.InterviewAddress.String
			}
			if event.ContactPerson.Valid {
				resp.Timeline[i].ContactPerson = event.ContactPerson.String
			}
			if event.ContactPhone.Valid {
				resp.Timeline[i].ContactPhone = event.ContactPhone.String
			}
		}
	}

	return resp
}

// GetStatusLabel returns human-readable status label
func GetStatusLabel(status string) string {
	labels := map[string]string{
		StatusSubmitted:          "Lamaran Terkirim",
		StatusViewed:             "Sedang Ditinjau",
		StatusShortlisted:        "Masuk Shortlist",
		StatusInterviewScheduled: "Interview Dijadwalkan",
		StatusInterviewCompleted: "Interview Selesai",
		StatusAssessment:         "Tahap Assessment",
		StatusOfferSent:          "Penawaran Dikirim",
		StatusOfferAccepted:      "Penawaran Diterima",
		StatusHired:              "Diterima",
		StatusRejected:           "Tidak Lolos",
		StatusWithdrawn:          "Dibatalkan",
	}

	if label, ok := labels[status]; ok {
		return label
	}
	return status
}

// IsTerminalStatus checks if status is terminal (no further updates)
func IsTerminalStatus(status string) bool {
	return status == StatusHired || status == StatusRejected || status == StatusWithdrawn
}

// IsValidStatusTransition checks if status transition is valid
func IsValidStatusTransition(from, to string) bool {
	// Define valid transitions - made more flexible for business needs
	// Companies should be able to skip stages if needed
	validTransitions := map[string][]string{
		StatusSubmitted: {
			StatusViewed, StatusShortlisted, StatusInterviewScheduled, StatusRejected,
		},
		StatusViewed: {
			StatusShortlisted, StatusInterviewScheduled, StatusRejected,
		},
		StatusShortlisted: {
			StatusInterviewScheduled, StatusRejected,
		},
		StatusInterviewScheduled: {
			StatusInterviewCompleted, StatusRejected,
		},
		StatusInterviewCompleted: {
			StatusAssessment, StatusOfferSent, StatusHired, StatusRejected,
		},
		StatusAssessment: {
			StatusOfferSent, StatusHired, StatusRejected,
		},
		StatusOfferSent: {
			StatusOfferAccepted, StatusRejected,
		},
		StatusOfferAccepted: {
			StatusHired,
		},
	}

	allowed, ok := validTransitions[from]
	if !ok {
		return false
	}

	for _, s := range allowed {
		if s == to {
			return true
		}
	}
	return false
}
