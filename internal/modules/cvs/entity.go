package cvs

import (
	"encoding/json"
	"time"
)

// PersonalInfo represents CV personal information
type PersonalInfo struct {
	FullName   string `json:"full_name"`
	Email      string `json:"email"`
	Phone      string `json:"phone,omitempty"`
	Address    string `json:"address,omitempty"`
	City       string `json:"city,omitempty"`
	Province   string `json:"province,omitempty"`
	Summary    string `json:"summary,omitempty"`
	LinkedIn   string `json:"linkedin,omitempty"`
	Portfolio  string `json:"portfolio,omitempty"`
	PhotoURL   string `json:"photo_url,omitempty"`
}

// Education represents an education entry
type Education struct {
	Institution  string `json:"institution"`
	Degree       string `json:"degree"`
	FieldOfStudy string `json:"field_of_study"`
	StartDate    string `json:"start_date"`
	EndDate      string `json:"end_date,omitempty"`
	GPA          string `json:"gpa,omitempty"`
	Description  string `json:"description,omitempty"`
}

// Experience represents a work experience entry
type Experience struct {
	Company      string   `json:"company"`
	Position     string   `json:"position"`
	Location     string   `json:"location,omitempty"`
	StartDate    string   `json:"start_date"`
	EndDate      string   `json:"end_date,omitempty"`
	IsCurrent    bool     `json:"is_current"`
	Description  string   `json:"description,omitempty"`
	Achievements []string `json:"achievements,omitempty"`
}

// Skill represents a skill entry
type Skill struct {
	Name     string `json:"name"`
	Level    string `json:"level,omitempty"` // beginner, intermediate, advanced, expert
	Category string `json:"category,omitempty"`
}

// Certification represents a certification entry
type Certification struct {
	Name          string `json:"name"`
	Issuer        string `json:"issuer"`
	IssueDate     string `json:"issue_date"`
	ExpiryDate    string `json:"expiry_date,omitempty"`
	CredentialID  string `json:"credential_id,omitempty"`
	CredentialURL string `json:"credential_url,omitempty"`
}

// Language represents a language skill
type Language struct {
	Name        string `json:"name"`
	Proficiency string `json:"proficiency"` // basic, conversational, proficient, fluent, native
}

// Project represents a project entry
type Project struct {
	Name        string   `json:"name"`
	Description string   `json:"description,omitempty"`
	URL         string   `json:"url,omitempty"`
	StartDate   string   `json:"start_date,omitempty"`
	EndDate     string   `json:"end_date,omitempty"`
	Skills      []string `json:"skills,omitempty"`
}

// CV represents a user's CV
type CV struct {
	ID                uint64          `db:"id" json:"id"`
	UserID            uint64          `db:"user_id" json:"user_id"`
	PersonalInfo      json.RawMessage `db:"personal_info" json:"-"`
	Education         json.RawMessage `db:"education" json:"-"`
	Experience        json.RawMessage `db:"experience" json:"-"`
	Skills            json.RawMessage `db:"skills" json:"-"`
	Certifications    json.RawMessage `db:"certifications" json:"-"`
	Languages         json.RawMessage `db:"languages" json:"-"`
	Projects          json.RawMessage `db:"projects" json:"-"`
	LastUpdatedAt     time.Time       `db:"last_updated_at" json:"last_updated_at"`
	CompletenessScore int             `db:"completeness_score" json:"completeness_score"`
	CreatedAt         time.Time       `db:"created_at" json:"created_at"`
	UpdatedAt         time.Time       `db:"updated_at" json:"updated_at"`

	// Parsed fields (not stored in DB)
	PersonalInfoParsed   *PersonalInfo   `db:"-" json:"personal_info"`
	EducationParsed      []Education     `db:"-" json:"education"`
	ExperienceParsed     []Experience    `db:"-" json:"experience"`
	SkillsParsed         []Skill         `db:"-" json:"skills"`
	CertificationsParsed []Certification `db:"-" json:"certifications"`
	LanguagesParsed      []Language      `db:"-" json:"languages"`
	ProjectsParsed       []Project       `db:"-" json:"projects"`
}

// CVSnapshot represents an immutable snapshot of CV at application time
type CVSnapshot struct {
	ID                uint64          `db:"id" json:"id"`
	CVID              uint64          `db:"cv_id" json:"cv_id"`
	UserID            uint64          `db:"user_id" json:"user_id"`
	PersonalInfo      json.RawMessage `db:"personal_info" json:"-"`
	Education         json.RawMessage `db:"education" json:"-"`
	Experience        json.RawMessage `db:"experience" json:"-"`
	Skills            json.RawMessage `db:"skills" json:"-"`
	Certifications    json.RawMessage `db:"certifications" json:"-"`
	Languages         json.RawMessage `db:"languages" json:"-"`
	Projects          json.RawMessage `db:"projects" json:"-"`
	SnapshotHash      string          `db:"snapshot_hash" json:"snapshot_hash"`
	CompletenessScore int             `db:"completeness_score" json:"completeness_score"`
	CreatedAt         time.Time       `db:"created_at" json:"created_at"`

	// Parsed fields
	PersonalInfoParsed   *PersonalInfo   `db:"-" json:"personal_info"`
	EducationParsed      []Education     `db:"-" json:"education"`
	ExperienceParsed     []Experience    `db:"-" json:"experience"`
	SkillsParsed         []Skill         `db:"-" json:"skills"`
	CertificationsParsed []Certification `db:"-" json:"certifications"`
	LanguagesParsed      []Language      `db:"-" json:"languages"`
	ProjectsParsed       []Project       `db:"-" json:"projects"`
}

// ParseFields parses JSON fields into structured data
func (cv *CV) ParseFields() error {
	if len(cv.PersonalInfo) > 0 {
		cv.PersonalInfoParsed = &PersonalInfo{}
		if err := json.Unmarshal(cv.PersonalInfo, cv.PersonalInfoParsed); err != nil {
			return err
		}
	}
	if len(cv.Education) > 0 {
		if err := json.Unmarshal(cv.Education, &cv.EducationParsed); err != nil {
			return err
		}
	}
	if len(cv.Experience) > 0 {
		if err := json.Unmarshal(cv.Experience, &cv.ExperienceParsed); err != nil {
			return err
		}
	}
	if len(cv.Skills) > 0 {
		if err := json.Unmarshal(cv.Skills, &cv.SkillsParsed); err != nil {
			return err
		}
	}
	if len(cv.Certifications) > 0 {
		if err := json.Unmarshal(cv.Certifications, &cv.CertificationsParsed); err != nil {
			return err
		}
	}
	if len(cv.Languages) > 0 {
		if err := json.Unmarshal(cv.Languages, &cv.LanguagesParsed); err != nil {
			return err
		}
	}
	if len(cv.Projects) > 0 {
		if err := json.Unmarshal(cv.Projects, &cv.ProjectsParsed); err != nil {
			return err
		}
	}
	return nil
}

// ParseFields parses JSON fields for CVSnapshot
func (s *CVSnapshot) ParseFields() error {
	if len(s.PersonalInfo) > 0 {
		s.PersonalInfoParsed = &PersonalInfo{}
		if err := json.Unmarshal(s.PersonalInfo, s.PersonalInfoParsed); err != nil {
			return err
		}
	}
	if len(s.Education) > 0 {
		if err := json.Unmarshal(s.Education, &s.EducationParsed); err != nil {
			return err
		}
	}
	if len(s.Experience) > 0 {
		if err := json.Unmarshal(s.Experience, &s.ExperienceParsed); err != nil {
			return err
		}
	}
	if len(s.Skills) > 0 {
		if err := json.Unmarshal(s.Skills, &s.SkillsParsed); err != nil {
			return err
		}
	}
	if len(s.Certifications) > 0 {
		if err := json.Unmarshal(s.Certifications, &s.CertificationsParsed); err != nil {
			return err
		}
	}
	if len(s.Languages) > 0 {
		if err := json.Unmarshal(s.Languages, &s.LanguagesParsed); err != nil {
			return err
		}
	}
	if len(s.Projects) > 0 {
		if err := json.Unmarshal(s.Projects, &s.ProjectsParsed); err != nil {
			return err
		}
	}
	return nil
}

// Request DTOs

// CreateCVRequest represents a CV creation/update request
type CreateCVRequest struct {
	PersonalInfo   *PersonalInfo   `json:"personal_info" validate:"required"`
	Education      []Education     `json:"education,omitempty"`
	Experience     []Experience    `json:"experience,omitempty"`
	Skills         []Skill         `json:"skills,omitempty"`
	Certifications []Certification `json:"certifications,omitempty"`
	Languages      []Language      `json:"languages,omitempty"`
	Projects       []Project       `json:"projects,omitempty"`
}

// UpdateCVRequest represents a partial CV update
type UpdateCVRequest struct {
	PersonalInfo   *PersonalInfo   `json:"personal_info,omitempty"`
	Education      []Education     `json:"education,omitempty"`
	Experience     []Experience    `json:"experience,omitempty"`
	Skills         []Skill         `json:"skills,omitempty"`
	Certifications []Certification `json:"certifications,omitempty"`
	Languages      []Language      `json:"languages,omitempty"`
	Projects       []Project       `json:"projects,omitempty"`
}

// Response DTOs

// CVResponse represents the CV response
type CVResponse struct {
	ID                uint64          `json:"id"`
	UserID            uint64          `json:"user_id"`
	PersonalInfo      *PersonalInfo   `json:"personal_info"`
	Education         []Education     `json:"education"`
	Experience        []Experience    `json:"experience"`
	Skills            []Skill         `json:"skills"`
	Certifications    []Certification `json:"certifications"`
	Languages         []Language      `json:"languages"`
	Projects          []Project       `json:"projects"`
	CompletenessScore int             `json:"completeness_score"`
	LastUpdatedAt     string          `json:"last_updated_at"`
	CreatedAt         string          `json:"created_at"`
}

// ToResponse converts CV to CVResponse
func (cv *CV) ToResponse() *CVResponse {
	return &CVResponse{
		ID:                cv.ID,
		UserID:            cv.UserID,
		PersonalInfo:      cv.PersonalInfoParsed,
		Education:         cv.EducationParsed,
		Experience:        cv.ExperienceParsed,
		Skills:            cv.SkillsParsed,
		Certifications:    cv.CertificationsParsed,
		Languages:         cv.LanguagesParsed,
		Projects:          cv.ProjectsParsed,
		CompletenessScore: cv.CompletenessScore,
		LastUpdatedAt:     cv.LastUpdatedAt.Format(time.RFC3339),
		CreatedAt:         cv.CreatedAt.Format(time.RFC3339),
	}
}

// CVSnapshotResponse represents the CV snapshot response
type CVSnapshotResponse struct {
	ID                uint64          `json:"id"`
	PersonalInfo      *PersonalInfo   `json:"personal_info"`
	Education         []Education     `json:"education"`
	Experience        []Experience    `json:"experience"`
	Skills            []Skill         `json:"skills"`
	Certifications    []Certification `json:"certifications"`
	Languages         []Language      `json:"languages"`
	Projects          []Project       `json:"projects"`
	CompletenessScore int             `json:"completeness_score"`
	CreatedAt         string          `json:"created_at"`
}

// ToResponse converts CVSnapshot to CVSnapshotResponse
func (s *CVSnapshot) ToResponse() *CVSnapshotResponse {
	return &CVSnapshotResponse{
		ID:                s.ID,
		PersonalInfo:      s.PersonalInfoParsed,
		Education:         s.EducationParsed,
		Experience:        s.ExperienceParsed,
		Skills:            s.SkillsParsed,
		Certifications:    s.CertificationsParsed,
		Languages:         s.LanguagesParsed,
		Projects:          s.ProjectsParsed,
		CompletenessScore: s.CompletenessScore,
		CreatedAt:         s.CreatedAt.Format(time.RFC3339),
	}
}
