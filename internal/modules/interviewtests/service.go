package interviewtests

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

var (
	ErrTestNotFound      = errors.New("interview test not found")
	ErrInvalidStatus     = errors.New("invalid status")
	ErrInvalidType       = errors.New("invalid question type")
	ErrInvalidDifficulty = errors.New("invalid difficulty")
	ErrNoQuestions       = errors.New("test must have at least one question")
	ErrInvalidOptions    = errors.New("multiple choice question must have at least 2 options with one correct answer")
)

// Service handles business logic for interview tests
type Service struct {
	repo *Repository
}

// NewService creates a new interview tests service
func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// Create creates a new interview test with questions
func (s *Service) Create(ctx context.Context, req CreateInterviewTestRequest, adminID uint64) (*InterviewTestResponse, error) {
	// Validate request
	if err := s.validateCreateRequest(req); err != nil {
		return nil, err
	}

	// Set defaults
	shuffleQuestions := false
	if req.ShuffleQuestions != nil {
		shuffleQuestions = *req.ShuffleQuestions
	}

	showResultsImmediately := false
	if req.ShowResultsImmediately != nil {
		showResultsImmediately = *req.ShowResultsImmediately
	}

	// Create test
	test := &InterviewTest{
		Title:                  req.Title,
		Description:            req.Description,
		DurationMinutes:        req.DurationMinutes,
		PassingScore:           req.PassingScore,
		ShuffleQuestions:       shuffleQuestions,
		ShowResultsImmediately: showResultsImmediately,
		Status:                 StatusDraft,
		CreatedBy:              adminID,
	}

	if err := s.repo.Create(ctx, test); err != nil {
		return nil, err
	}

	// Create questions and options
	for i, qReq := range req.Questions {
		question := &InterviewQuestion{
			InterviewTestID: test.ID,
			QuestionText:    qReq.QuestionText,
			QuestionType:    QuestionType(qReq.QuestionType),
			Points:          qReq.Points,
			Difficulty:      QuestionDifficulty(qReq.Difficulty),
			Order:           i + 1,
		}

		if qReq.Explanation != nil {
			question.Explanation = sql.NullString{String: *qReq.Explanation, Valid: true}
		}

		if err := s.repo.CreateQuestion(ctx, question); err != nil {
			return nil, err
		}

		// Create options for multiple choice questions
		if question.QuestionType == TypeMultipleChoice {
			for j, optReq := range qReq.Options {
				option := &QuestionOption{
					InterviewQuestionID: question.ID,
					OptionText:          optReq.OptionText,
					IsCorrect:           optReq.IsCorrect,
					Order:               j + 1,
				}

				if err := s.repo.CreateOption(ctx, option); err != nil {
					return nil, err
				}
			}
		}
	}

	// Fetch the complete test with questions
	return s.GetByID(ctx, test.ID)
}

// GetByID retrieves an interview test by ID with questions and options
func (s *Service) GetByID(ctx context.Context, id uint64) (*InterviewTestResponse, error) {
	test, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if test == nil {
		return nil, ErrTestNotFound
	}

	// Get questions
	questions, err := s.repo.GetQuestionsByTestID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Get all options for this test
	optionsMap, err := s.repo.GetOptionsByTestID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Build response
	resp := test.ToResponse()
	resp.Questions = make([]InterviewQuestionResponse, len(questions))

	for i, q := range questions {
		qResp := q.ToResponse()

		// Add options if this is a multiple choice question
		if q.QuestionType == TypeMultipleChoice {
			if opts, ok := optionsMap[q.ID]; ok {
				qResp.Options = make([]QuestionOptionResponse, len(opts))
				for j, opt := range opts {
					qResp.Options[j] = opt.ToResponse()
				}
			}
		}

		resp.Questions[i] = qResp
	}

	return &resp, nil
}

// GetAll retrieves all interview tests
func (s *Service) GetAll(ctx context.Context, status string) ([]InterviewTestResponse, error) {
	// Validate status if provided
	if status != "" && !isValidStatus(status) {
		return nil, ErrInvalidStatus
	}

	tests, err := s.repo.GetAll(ctx, status)
	if err != nil {
		return nil, err
	}

	responses := make([]InterviewTestResponse, len(tests))
	for i, test := range tests {
		responses[i] = test.ToResponse()
	}

	return responses, nil
}

// Update updates an existing interview test
func (s *Service) Update(ctx context.Context, id uint64, req UpdateInterviewTestRequest, adminID uint64) (*InterviewTestResponse, error) {
	// Check if test exists
	test, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if test == nil {
		return nil, ErrTestNotFound
	}

	// Validate request
	if err := s.validateUpdateRequest(req); err != nil {
		return nil, err
	}

	// Update test fields
	test.Title = req.Title
	test.Description = req.Description
	test.DurationMinutes = req.DurationMinutes
	test.PassingScore = req.PassingScore

	if req.ShuffleQuestions != nil {
		test.ShuffleQuestions = *req.ShuffleQuestions
	}
	if req.ShowResultsImmediately != nil {
		test.ShowResultsImmediately = *req.ShowResultsImmediately
	}

	test.UpdatedBy = sql.NullInt64{Int64: int64(adminID), Valid: true}

	if err := s.repo.Update(ctx, test); err != nil {
		return nil, err
	}

	// Delete existing questions and options
	if err := s.repo.DeleteQuestionsByTestID(ctx, id); err != nil {
		return nil, err
	}

	// Create new questions and options
	for i, qReq := range req.Questions {
		question := &InterviewQuestion{
			InterviewTestID: test.ID,
			QuestionText:    qReq.QuestionText,
			QuestionType:    QuestionType(qReq.QuestionType),
			Points:          qReq.Points,
			Difficulty:      QuestionDifficulty(qReq.Difficulty),
			Order:           i + 1,
		}

		if qReq.Explanation != nil {
			question.Explanation = sql.NullString{String: *qReq.Explanation, Valid: true}
		}

		if err := s.repo.CreateQuestion(ctx, question); err != nil {
			return nil, err
		}

		// Create options for multiple choice questions
		if question.QuestionType == TypeMultipleChoice {
			for j, optReq := range qReq.Options {
				option := &QuestionOption{
					InterviewQuestionID: question.ID,
					OptionText:          optReq.OptionText,
					IsCorrect:           optReq.IsCorrect,
					Order:               j + 1,
				}

				if err := s.repo.CreateOption(ctx, option); err != nil {
					return nil, err
				}
			}
		}
	}

	// Fetch the updated test
	return s.GetByID(ctx, test.ID)
}

// Delete soft deletes an interview test
func (s *Service) Delete(ctx context.Context, id uint64) error {
	test, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if test == nil {
		return ErrTestNotFound
	}

	return s.repo.Delete(ctx, id)
}

// Publish publishes an interview test (changes status to active)
func (s *Service) Publish(ctx context.Context, id uint64, adminID uint64) (*InterviewTestResponse, error) {
	test, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if test == nil {
		return nil, ErrTestNotFound
	}

	// Check if test has questions
	questions, err := s.repo.GetQuestionsByTestID(ctx, id)
	if err != nil {
		return nil, err
	}
	if len(questions) == 0 {
		return nil, ErrNoQuestions
	}

	// Update status to active
	if err := s.repo.UpdateStatus(ctx, id, StatusActive, adminID); err != nil {
		return nil, err
	}

	return s.GetByID(ctx, id)
}

// Duplicate creates a copy of an existing test
func (s *Service) Duplicate(ctx context.Context, id uint64, adminID uint64) (*InterviewTestResponse, error) {
	// Get the original test
	original, err := s.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Create duplicate request
	req := CreateInterviewTestRequest{
		Title:                  fmt.Sprintf("%s (Copy)", original.Title),
		Description:            original.Description,
		DurationMinutes:        original.DurationMinutes,
		PassingScore:           original.PassingScore,
		ShuffleQuestions:       &original.ShuffleQuestions,
		ShowResultsImmediately: &original.ShowResultsImmediately,
		Questions:              make([]CreateQuestionRequest, len(original.Questions)),
	}

	// Copy questions
	for i, q := range original.Questions {
		qReq := CreateQuestionRequest{
			QuestionText: q.QuestionText,
			QuestionType: q.QuestionType,
			Points:       q.Points,
			Difficulty:   q.Difficulty,
			Explanation:  q.Explanation,
			Options:      make([]CreateOptionRequest, len(q.Options)),
		}

		// Copy options
		for j, opt := range q.Options {
			qReq.Options[j] = CreateOptionRequest{
				OptionText: opt.OptionText,
				IsCorrect:  opt.IsCorrect,
			}
		}

		req.Questions[i] = qReq
	}

	return s.Create(ctx, req, adminID)
}

// =============== Validation helpers ===============

func (s *Service) validateCreateRequest(req CreateInterviewTestRequest) error {
	if req.Title == "" {
		return errors.New("title is required")
	}
	if req.Description == "" {
		return errors.New("description is required")
	}
	if req.DurationMinutes <= 0 {
		return errors.New("duration must be greater than 0")
	}
	if req.PassingScore < 0 || req.PassingScore > 100 {
		return errors.New("passing score must be between 0 and 100")
	}
	if len(req.Questions) == 0 {
		return ErrNoQuestions
	}

	// Validate each question
	for i, q := range req.Questions {
		if err := s.validateQuestion(q); err != nil {
			return fmt.Errorf("question %d: %w", i+1, err)
		}
	}

	return nil
}

func (s *Service) validateUpdateRequest(req UpdateInterviewTestRequest) error {
	if req.Title == "" {
		return errors.New("title is required")
	}
	if req.Description == "" {
		return errors.New("description is required")
	}
	if req.DurationMinutes <= 0 {
		return errors.New("duration must be greater than 0")
	}
	if req.PassingScore < 0 || req.PassingScore > 100 {
		return errors.New("passing score must be between 0 and 100")
	}
	if len(req.Questions) == 0 {
		return ErrNoQuestions
	}

	// Validate each question
	for i, q := range req.Questions {
		if err := s.validateQuestion(q); err != nil {
			return fmt.Errorf("question %d: %w", i+1, err)
		}
	}

	return nil
}

func (s *Service) validateQuestion(q CreateQuestionRequest) error {
	if q.QuestionText == "" {
		return errors.New("question text is required")
	}
	if !isValidType(q.QuestionType) {
		return ErrInvalidType
	}
	if !isValidDifficulty(q.Difficulty) {
		return ErrInvalidDifficulty
	}
	if q.Points <= 0 {
		return errors.New("points must be greater than 0")
	}

	// Validate multiple choice questions
	if q.QuestionType == string(TypeMultipleChoice) {
		if len(q.Options) < 2 {
			return ErrInvalidOptions
		}

		correctCount := 0
		for _, opt := range q.Options {
			if opt.OptionText == "" {
				return errors.New("option text is required")
			}
			if opt.IsCorrect {
				correctCount++
			}
		}

		if correctCount == 0 {
			return errors.New("multiple choice question must have at least one correct answer")
		}
	}

	return nil
}

func isValidStatus(status string) bool {
	return status == string(StatusDraft) ||
		status == string(StatusActive) ||
		status == string(StatusArchived)
}

func isValidType(qType string) bool {
	return qType == string(TypeMultipleChoice) || qType == string(TypeEssay)
}

func isValidDifficulty(difficulty string) bool {
	return difficulty == string(DifficultyEasy) ||
		difficulty == string(DifficultyMedium) ||
		difficulty == string(DifficultyHard)
}
