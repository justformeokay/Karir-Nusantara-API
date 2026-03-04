package interviewtests

import (
	"database/sql"
	"time"
)

// InterviewTestStatus represents the status of an interview test
type InterviewTestStatus string

const (
	StatusDraft    InterviewTestStatus = "draft"
	StatusActive   InterviewTestStatus = "active"
	StatusArchived InterviewTestStatus = "archived"
)

// QuestionType represents the type of question
type QuestionType string

const (
	TypeMultipleChoice QuestionType = "multiple_choice"
	TypeEssay          QuestionType = "essay"
)

// QuestionDifficulty represents the difficulty level
type QuestionDifficulty string

const (
	DifficultyEasy   QuestionDifficulty = "easy"
	DifficultyMedium QuestionDifficulty = "medium"
	DifficultyHard   QuestionDifficulty = "hard"
)

// SubmissionStatus represents the submission status
type SubmissionStatus string

const (
	SubmissionInProgress SubmissionStatus = "in_progress"
	SubmissionSubmitted  SubmissionStatus = "submitted"
	SubmissionGrading    SubmissionStatus = "grading"
	SubmissionCompleted  SubmissionStatus = "completed"
)

// OwnerType represents the type of owner for an interview test
type OwnerType string

const (
	OwnerSuperAdmin OwnerType = "super_admin"
	OwnerCompany    OwnerType = "company"
)

// InterviewTest represents an interview test entity
type InterviewTest struct {
	ID                     uint64              `db:"id" json:"id"`
	Title                  string              `db:"title" json:"title"`
	Description            string              `db:"description" json:"description"`
	DurationMinutes        int                 `db:"duration_minutes" json:"duration_minutes"`
	TotalPoints            int                 `db:"total_points" json:"total_points"`
	PassingScore           int                 `db:"passing_score" json:"passing_score"`
	ShuffleQuestions       bool                `db:"shuffle_questions" json:"shuffle_questions"`
	ShowResultsImmediately bool                `db:"show_results_immediately" json:"show_results_immediately"`
	Status                 InterviewTestStatus `db:"status" json:"status"`
	OwnerType              OwnerType           `db:"owner_type" json:"owner_type"`
	OwnerCompanyID         sql.NullInt64       `db:"owner_company_id" json:"owner_company_id,omitempty"`
	IsPublic               bool                `db:"is_public" json:"is_public"`
	CreatedBy              uint64              `db:"created_by" json:"created_by"`
	UpdatedBy              sql.NullInt64       `db:"updated_by" json:"updated_by,omitempty"`
	CreatedAt              time.Time           `db:"created_at" json:"created_at"`
	UpdatedAt              time.Time           `db:"updated_at" json:"updated_at"`
	DeletedAt              sql.NullTime        `db:"deleted_at" json:"deleted_at,omitempty"`
}

// InterviewQuestion represents a question entity
type InterviewQuestion struct {
	ID              uint64             `db:"id" json:"id"`
	InterviewTestID uint64             `db:"interview_test_id" json:"interview_test_id"`
	QuestionText    string             `db:"question_text" json:"question_text"`
	QuestionType    QuestionType       `db:"question_type" json:"question_type"`
	Points          int                `db:"points" json:"points"`
	Difficulty      QuestionDifficulty `db:"difficulty" json:"difficulty"`
	Order           int                `db:"order" json:"order"`
	Explanation     sql.NullString     `db:"explanation" json:"explanation,omitempty"`
	CreatedAt       time.Time          `db:"created_at" json:"created_at"`
	UpdatedAt       time.Time          `db:"updated_at" json:"updated_at"`
}

// QuestionOption represents a question option for multiple choice
type QuestionOption struct {
	ID                  uint64    `db:"id" json:"id"`
	InterviewQuestionID uint64    `db:"interview_question_id" json:"interview_question_id"`
	OptionText          string    `db:"option_text" json:"option_text"`
	IsCorrect           bool      `db:"is_correct" json:"is_correct"`
	Order               int       `db:"order" json:"order"`
	CreatedAt           time.Time `db:"created_at" json:"created_at"`
	UpdatedAt           time.Time `db:"updated_at" json:"updated_at"`
}

// InterviewTestSubmission represents a test submission
type InterviewTestSubmission struct {
	ID              uint64           `db:"id" json:"id"`
	InterviewTestID uint64           `db:"interview_test_id" json:"interview_test_id"`
	UserID          uint64           `db:"user_id" json:"user_id"`
	ApplicationID   sql.NullInt64    `db:"application_id" json:"application_id,omitempty"`
	Status          SubmissionStatus `db:"status" json:"status"`
	Score           sql.NullInt64    `db:"score" json:"score,omitempty"`
	Percentage      sql.NullFloat64  `db:"percentage" json:"percentage,omitempty"`
	IsPassed        sql.NullBool     `db:"is_passed" json:"is_passed,omitempty"`
	StartedAt       time.Time        `db:"started_at" json:"started_at"`
	SubmittedAt     sql.NullTime     `db:"submitted_at" json:"submitted_at,omitempty"`
	GradedAt        sql.NullTime     `db:"graded_at" json:"graded_at,omitempty"`
	GradedBy        sql.NullInt64    `db:"graded_by" json:"graded_by,omitempty"`
	Notes           sql.NullString   `db:"notes" json:"notes,omitempty"`
	CreatedAt       time.Time        `db:"created_at" json:"created_at"`
	UpdatedAt       time.Time        `db:"updated_at" json:"updated_at"`
}

// InterviewTestAnswer represents an answer
type InterviewTestAnswer struct {
	ID                  uint64         `db:"id" json:"id"`
	SubmissionID        uint64         `db:"submission_id" json:"submission_id"`
	InterviewQuestionID uint64         `db:"interview_question_id" json:"interview_question_id"`
	QuestionType        QuestionType   `db:"question_type" json:"question_type"`
	AnswerText          sql.NullString `db:"answer_text" json:"answer_text,omitempty"`
	SelectedOptionID    sql.NullInt64  `db:"selected_option_id" json:"selected_option_id,omitempty"`
	IsCorrect           sql.NullBool   `db:"is_correct" json:"is_correct,omitempty"`
	PointsEarned        sql.NullInt64  `db:"points_earned" json:"points_earned,omitempty"`
	GraderFeedback      sql.NullString `db:"grader_feedback" json:"grader_feedback,omitempty"`
	AnsweredAt          time.Time      `db:"answered_at" json:"answered_at"`
	CreatedAt           time.Time      `db:"created_at" json:"created_at"`
	UpdatedAt           time.Time      `db:"updated_at" json:"updated_at"`
}

// =============== Response DTOs ===============

// InterviewTestResponse is the API response for interview test
type InterviewTestResponse struct {
	ID                     uint64                      `json:"id"`
	Title                  string                      `json:"title"`
	Description            string                      `json:"description"`
	DurationMinutes        int                         `json:"duration_minutes"`
	TotalPoints            int                         `json:"total_points"`
	PassingScore           int                         `json:"passing_score"`
	ShuffleQuestions       bool                        `json:"shuffle_questions"`
	ShowResultsImmediately bool                        `json:"show_results_immediately"`
	Status                 string                      `json:"status"`
	OwnerType              string                      `json:"owner_type"`
	OwnerCompanyID         *int64                      `json:"owner_company_id,omitempty"`
	IsPublic               bool                        `json:"is_public"`
	Questions              []InterviewQuestionResponse `json:"questions,omitempty"`
	CreatedBy              uint64                      `json:"created_by"`
	CreatedAt              string                      `json:"created_at"`
	UpdatedAt              string                      `json:"updated_at"`
}

// InterviewQuestionResponse is the API response for questions
type InterviewQuestionResponse struct {
	ID           uint64                   `json:"id"`
	QuestionText string                   `json:"question_text"`
	QuestionType string                   `json:"question_type"`
	Points       int                      `json:"points"`
	Difficulty   string                   `json:"difficulty"`
	Order        int                      `json:"order"`
	Explanation  *string                  `json:"explanation,omitempty"`
	Options      []QuestionOptionResponse `json:"options,omitempty"`
}

// QuestionOptionResponse is the API response for options
type QuestionOptionResponse struct {
	ID         uint64 `json:"id"`
	OptionText string `json:"option_text"`
	IsCorrect  bool   `json:"is_correct"`
	Order      int    `json:"order"`
}

// ToResponse converts InterviewTest to InterviewTestResponse
func (t *InterviewTest) ToResponse() InterviewTestResponse {
	resp := InterviewTestResponse{
		ID:                     t.ID,
		Title:                  t.Title,
		Description:            t.Description,
		DurationMinutes:        t.DurationMinutes,
		TotalPoints:            t.TotalPoints,
		PassingScore:           t.PassingScore,
		ShuffleQuestions:       t.ShuffleQuestions,
		ShowResultsImmediately: t.ShowResultsImmediately,
		Status:                 string(t.Status),
		OwnerType:              string(t.OwnerType),
		IsPublic:               t.IsPublic,
		CreatedBy:              t.CreatedBy,
		CreatedAt:              t.CreatedAt.Format(time.RFC3339),
		UpdatedAt:              t.UpdatedAt.Format(time.RFC3339),
	}
	if t.OwnerCompanyID.Valid {
		resp.OwnerCompanyID = &t.OwnerCompanyID.Int64
	}
	return resp
}

// ToResponse converts InterviewQuestion to InterviewQuestionResponse
func (q *InterviewQuestion) ToResponse() InterviewQuestionResponse {
	resp := InterviewQuestionResponse{
		ID:           q.ID,
		QuestionText: q.QuestionText,
		QuestionType: string(q.QuestionType),
		Points:       q.Points,
		Difficulty:   string(q.Difficulty),
		Order:        q.Order,
	}

	if q.Explanation.Valid {
		resp.Explanation = &q.Explanation.String
	}

	return resp
}

// ToResponse converts QuestionOption to QuestionOptionResponse
func (o *QuestionOption) ToResponse() QuestionOptionResponse {
	return QuestionOptionResponse{
		ID:         o.ID,
		OptionText: o.OptionText,
		IsCorrect:  o.IsCorrect,
		Order:      o.Order,
	}
}

// =============== Request DTOs ===============

// CreateInterviewTestRequest represents the request to create a test
type CreateInterviewTestRequest struct {
	Title                  string                  `json:"title"`
	Description            string                  `json:"description"`
	DurationMinutes        int                     `json:"duration_minutes"`
	PassingScore           int                     `json:"passing_score"`
	ShuffleQuestions       *bool                   `json:"shuffle_questions,omitempty"`
	ShowResultsImmediately *bool                   `json:"show_results_immediately,omitempty"`
	Questions              []CreateQuestionRequest `json:"questions"`
}

// CreateQuestionRequest represents the request to create a question
type CreateQuestionRequest struct {
	QuestionText string                `json:"question_text"`
	QuestionType string                `json:"question_type"`
	Points       int                   `json:"points"`
	Difficulty   string                `json:"difficulty"`
	Explanation  *string               `json:"explanation,omitempty"`
	Options      []CreateOptionRequest `json:"options,omitempty"`
}

// CreateOptionRequest represents the request to create an option
type CreateOptionRequest struct {
	OptionText string `json:"option_text"`
	IsCorrect  bool   `json:"is_correct"`
}

// UpdateInterviewTestRequest represents the request to update a test
type UpdateInterviewTestRequest struct {
	Title                  string                  `json:"title"`
	Description            string                  `json:"description"`
	DurationMinutes        int                     `json:"duration_minutes"`
	PassingScore           int                     `json:"passing_score"`
	ShuffleQuestions       *bool                   `json:"shuffle_questions,omitempty"`
	ShowResultsImmediately *bool                   `json:"show_results_immediately,omitempty"`
	Questions              []CreateQuestionRequest `json:"questions"`
}

// AssignTestRequest represents the request to assign a test to a candidate
type AssignTestRequest struct {
	InterviewTestID uint64 `json:"interview_test_id"`
	CandidateUserID uint64 `json:"candidate_user_id"`
	ApplicationID   uint64 `json:"application_id"`
}

// SubmitAnswerRequest represents a single answer submitted by a job seeker
type SubmitAnswerRequest struct {
	QuestionID       uint64 `json:"question_id"`
	QuestionType     string `json:"question_type"`
	SelectedOptionID uint64 `json:"selected_option_id,omitempty"`
	AnswerText       string `json:"answer_text,omitempty"`
}

// SubmitAnswersRequest wraps multiple answers in one submission
type SubmitAnswersRequest struct {
	Answers []SubmitAnswerRequest `json:"answers"`
}

// SubmissionResponse is the API response for a test submission
type SubmissionResponse struct {
	ID            uint64                `json:"id"`
	Status        string                `json:"status"`
	Score         *int64                `json:"score,omitempty"`
	Percentage    *float64              `json:"percentage,omitempty"`
	IsPassed      *bool                 `json:"is_passed,omitempty"`
	StartedAt     *string               `json:"started_at,omitempty"`
	SubmittedAt   *string               `json:"submitted_at,omitempty"`
	ApplicationID *uint64               `json:"application_id,omitempty"`
	Test          InterviewTestResponse `json:"test"`
}

// TestForSubmissionResponse wraps the test + submission context for a job seeker to take
type TestForSubmissionResponse struct {
	SubmissionID uint64                `json:"submission_id"`
	Status       string                `json:"status"`
	Test         InterviewTestResponse `json:"test"`
}
