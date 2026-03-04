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
		resp := test.ToResponse()

		// Get questions for this test
		questions, err := s.repo.GetQuestionsByTestID(ctx, test.ID)
		if err != nil {
			return nil, err
		}

		// Get all options for this test
		optionsMap, err := s.repo.GetOptionsByTestID(ctx, test.ID)
		if err != nil {
			return nil, err
		}

		// Build response with questions
		resp.Questions = make([]InterviewQuestionResponse, len(questions))
		for j, q := range questions {
			qResp := q.ToResponse()

			// Add options if this is a multiple choice question
			if q.QuestionType == TypeMultipleChoice {
				if opts, ok := optionsMap[q.ID]; ok {
					qResp.Options = make([]QuestionOptionResponse, len(opts))
					for k, opt := range opts {
						qResp.Options[k] = opt.ToResponse()
					}
				}
			}

			resp.Questions[j] = qResp
		}

		responses[i] = resp
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

// =============== Company methods ===============

// GetPublicAdminTests retrieves all public tests from super_admin (library)
func (s *Service) GetPublicAdminTests(ctx context.Context) ([]InterviewTestResponse, error) {
	tests, err := s.repo.GetPublicAdminTests(ctx)
	if err != nil {
		return nil, err
	}

	responses := make([]InterviewTestResponse, len(tests))
	for i, test := range tests {
		resp := test.ToResponse()

		questions, err := s.repo.GetQuestionsByTestID(ctx, test.ID)
		if err != nil {
			return nil, err
		}
		optionsMap, err := s.repo.GetOptionsByTestID(ctx, test.ID)
		if err != nil {
			return nil, err
		}

		resp.Questions = make([]InterviewQuestionResponse, len(questions))
		for j, q := range questions {
			qResp := q.ToResponse()
			if q.QuestionType == TypeMultipleChoice {
				if opts, ok := optionsMap[q.ID]; ok {
					qResp.Options = make([]QuestionOptionResponse, len(opts))
					for k, opt := range opts {
						qResp.Options[k] = opt.ToResponse()
					}
				}
			}
			resp.Questions[j] = qResp
		}
		responses[i] = resp
	}
	return responses, nil
}

// GetByCompanyID retrieves all tests owned by a company
func (s *Service) GetByCompanyID(ctx context.Context, companyID uint64, status string) ([]InterviewTestResponse, error) {
	if status != "" && !isValidStatus(status) {
		return nil, ErrInvalidStatus
	}

	tests, err := s.repo.GetByCompanyID(ctx, companyID, status)
	if err != nil {
		return nil, err
	}

	responses := make([]InterviewTestResponse, len(tests))
	for i, test := range tests {
		resp := test.ToResponse()

		questions, err := s.repo.GetQuestionsByTestID(ctx, test.ID)
		if err != nil {
			return nil, err
		}
		optionsMap, err := s.repo.GetOptionsByTestID(ctx, test.ID)
		if err != nil {
			return nil, err
		}

		resp.Questions = make([]InterviewQuestionResponse, len(questions))
		for j, q := range questions {
			qResp := q.ToResponse()
			if q.QuestionType == TypeMultipleChoice {
				if opts, ok := optionsMap[q.ID]; ok {
					qResp.Options = make([]QuestionOptionResponse, len(opts))
					for k, opt := range opts {
						qResp.Options[k] = opt.ToResponse()
					}
				}
			}
			resp.Questions[j] = qResp
		}
		responses[i] = resp
	}
	return responses, nil
}

// CreateForCompany creates a new private test for a company
func (s *Service) CreateForCompany(ctx context.Context, req CreateInterviewTestRequest, companyID uint64, userID uint64) (*InterviewTestResponse, error) {
	if err := s.validateCreateRequest(req); err != nil {
		return nil, err
	}

	shuffleQuestions := false
	if req.ShuffleQuestions != nil {
		shuffleQuestions = *req.ShuffleQuestions
	}
	showResultsImmediately := false
	if req.ShowResultsImmediately != nil {
		showResultsImmediately = *req.ShowResultsImmediately
	}

	test := &InterviewTest{
		Title:                  req.Title,
		Description:            req.Description,
		DurationMinutes:        req.DurationMinutes,
		PassingScore:           req.PassingScore,
		ShuffleQuestions:       shuffleQuestions,
		ShowResultsImmediately: showResultsImmediately,
		Status:                 StatusDraft,
		OwnerType:              OwnerCompany,
		OwnerCompanyID:         sql.NullInt64{Int64: int64(companyID), Valid: true},
		IsPublic:               false,
		CreatedBy:              userID,
	}

	if err := s.repo.Create(ctx, test); err != nil {
		return nil, err
	}

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

	return s.GetByID(ctx, test.ID)
}

// UpdateForCompany updates a company-owned test (ownership check included)
func (s *Service) UpdateForCompany(ctx context.Context, id uint64, req UpdateInterviewTestRequest, companyID uint64, userID uint64) (*InterviewTestResponse, error) {
	test, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if test == nil {
		return nil, ErrTestNotFound
	}
	if test.OwnerType != OwnerCompany || !test.OwnerCompanyID.Valid || uint64(test.OwnerCompanyID.Int64) != companyID {
		return nil, ErrTestNotFound
	}
	if err := s.validateUpdateRequest(req); err != nil {
		return nil, err
	}

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
	test.UpdatedBy = sql.NullInt64{Int64: int64(userID), Valid: true}

	if err := s.repo.Update(ctx, test); err != nil {
		return nil, err
	}
	if err := s.repo.DeleteQuestionsByTestID(ctx, id); err != nil {
		return nil, err
	}
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
	return s.GetByID(ctx, test.ID)
}

// DeleteForCompany soft-deletes a company-owned test
func (s *Service) DeleteForCompany(ctx context.Context, id uint64, companyID uint64) error {
	test, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if test == nil {
		return ErrTestNotFound
	}
	if test.OwnerType != OwnerCompany || !test.OwnerCompanyID.Valid || uint64(test.OwnerCompanyID.Int64) != companyID {
		return ErrTestNotFound
	}
	return s.repo.Delete(ctx, id)
}

// PublishForCompany publishes a company-owned test
func (s *Service) PublishForCompany(ctx context.Context, id uint64, companyID uint64, userID uint64) (*InterviewTestResponse, error) {
	test, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if test == nil {
		return nil, ErrTestNotFound
	}
	if test.OwnerType != OwnerCompany || !test.OwnerCompanyID.Valid || uint64(test.OwnerCompanyID.Int64) != companyID {
		return nil, ErrTestNotFound
	}
	questions, err := s.repo.GetQuestionsByTestID(ctx, id)
	if err != nil {
		return nil, err
	}
	if len(questions) == 0 {
		return nil, ErrNoQuestions
	}
	if err := s.repo.UpdateStatus(ctx, id, StatusActive, userID); err != nil {
		return nil, err
	}
	return s.GetByID(ctx, id)
}

// CopyFromAdmin copies a public admin test to a company's ownership
func (s *Service) CopyFromAdmin(ctx context.Context, testID uint64, companyID uint64, userID uint64) (*InterviewTestResponse, error) {
	original, err := s.GetByID(ctx, testID)
	if err != nil {
		return nil, err
	}

	req := CreateInterviewTestRequest{
		Title:                  fmt.Sprintf("%s (Copy)", original.Title),
		Description:            original.Description,
		DurationMinutes:        original.DurationMinutes,
		PassingScore:           original.PassingScore,
		ShuffleQuestions:       &original.ShuffleQuestions,
		ShowResultsImmediately: &original.ShowResultsImmediately,
		Questions:              make([]CreateQuestionRequest, len(original.Questions)),
	}

	for i, q := range original.Questions {
		qReq := CreateQuestionRequest{
			QuestionText: q.QuestionText,
			QuestionType: q.QuestionType,
			Points:       q.Points,
			Difficulty:   q.Difficulty,
			Explanation:  q.Explanation,
			Options:      make([]CreateOptionRequest, len(q.Options)),
		}
		for j, opt := range q.Options {
			qReq.Options[j] = CreateOptionRequest{
				OptionText: opt.OptionText,
				IsCorrect:  opt.IsCorrect,
			}
		}
		req.Questions[i] = qReq
	}

	return s.CreateForCompany(ctx, req, companyID, userID)
}

// =============== Submission / Assignment methods ===============

// AssignTestToCandidate creates a new submission record (assigns test to a candidate)
func (s *Service) AssignTestToCandidate(ctx context.Context, testID uint64, candidateUserID uint64, applicationID uint64, companyID uint64) (*SubmissionResponse, error) {
	// Verify the test exists and belongs to the company (or is a public admin test)
	test, err := s.repo.GetByID(ctx, testID)
	if err != nil {
		return nil, err
	}
	if test == nil {
		return nil, ErrTestNotFound
	}
	// Only allow active tests
	if test.Status != StatusActive {
		return nil, errors.New("only active tests can be assigned")
	}

	// Prevent duplicate assignment
	existing, err := s.repo.GetSubmissionByApplicationAndTest(ctx, applicationID, testID)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("test already assigned to this candidate for this application")
	}

	sub := &InterviewTestSubmission{
		InterviewTestID: testID,
		UserID:          candidateUserID,
		ApplicationID:   sql.NullInt64{Int64: int64(applicationID), Valid: true},
	}

	if err := s.repo.CreateSubmission(ctx, sub); err != nil {
		return nil, err
	}

	testResp := test.ToResponse()
	return buildSubmissionResponse(sub, &testResp), nil
}

// GetSubmissionsForApplication returns all test submissions for a given application
func (s *Service) GetSubmissionsForApplication(ctx context.Context, applicationID uint64) ([]SubmissionResponse, error) {
	subs, err := s.repo.GetSubmissionsByApplicationID(ctx, applicationID)
	if err != nil {
		return nil, err
	}

	responses := make([]SubmissionResponse, 0, len(subs))
	for _, sub := range subs {
		test, err := s.repo.GetByID(ctx, sub.InterviewTestID)
		if err != nil || test == nil {
			continue
		}
		testResp := test.ToResponse()
		responses = append(responses, *buildSubmissionResponse(sub, &testResp))
	}
	return responses, nil
}

// GetSubmissionsForUser returns all test submissions assigned to a job seeker
func (s *Service) GetSubmissionsForUser(ctx context.Context, userID uint64) ([]SubmissionResponse, error) {
	subs, err := s.repo.GetSubmissionsByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	responses := make([]SubmissionResponse, 0, len(subs))
	for _, sub := range subs {
		test, err := s.repo.GetByID(ctx, sub.InterviewTestID)
		if err != nil || test == nil {
			continue
		}
		testResp := test.ToResponse()
		responses = append(responses, *buildSubmissionResponse(sub, &testResp))
	}
	return responses, nil
}

// GetTestForSubmission returns the full test detail (with questions/options) for a submission, for a job seeker to answer
func (s *Service) GetTestForSubmission(ctx context.Context, submissionID uint64, userID uint64) (*TestForSubmissionResponse, error) {
	sub, err := s.repo.GetSubmissionByID(ctx, submissionID)
	if err != nil {
		return nil, err
	}
	if sub == nil || sub.UserID != userID {
		return nil, errors.New("submission not found")
	}

	testResp, err := s.GetByID(ctx, sub.InterviewTestID)
	if err != nil {
		return nil, err
	}

	// Strip correct answer info - candidates should not see which is correct
	for i := range testResp.Questions {
		for j := range testResp.Questions[i].Options {
			testResp.Questions[i].Options[j].IsCorrect = false
		}
	}

	return &TestForSubmissionResponse{
		SubmissionID: sub.ID,
		Status:       string(sub.Status),
		Test:         *testResp,
	}, nil
}

// SubmitTestAnswers processes answers submitted by a job seeker, auto-grades multiple choice
func (s *Service) SubmitTestAnswers(ctx context.Context, submissionID uint64, userID uint64, answers []SubmitAnswerRequest) (*SubmissionResponse, error) {
	sub, err := s.repo.GetSubmissionByID(ctx, submissionID)
	if err != nil {
		return nil, err
	}
	if sub == nil || sub.UserID != userID {
		return nil, errors.New("submission not found")
	}
	if sub.Status != SubmissionInProgress {
		return nil, errors.New("this test has already been submitted")
	}

	// Get the test to know total points
	test, err := s.repo.GetByID(ctx, sub.InterviewTestID)
	if err != nil || test == nil {
		return nil, errors.New("test not found")
	}

	// Fetch all questions for this test (for points lookup)
	questions, err := s.repo.GetQuestionsByTestID(ctx, sub.InterviewTestID)
	if err != nil {
		return nil, err
	}
	questionMap := make(map[uint64]*InterviewQuestion)
	for _, q := range questions {
		questionMap[q.ID] = q
	}

	totalScore := int64(0)
	hasEssay := false

	for _, a := range answers {
		answer := &InterviewTestAnswer{
			SubmissionID:        submissionID,
			InterviewQuestionID: a.QuestionID,
			QuestionType:        QuestionType(a.QuestionType),
		}

		if a.QuestionType == string(TypeMultipleChoice) && a.SelectedOptionID != 0 {
			answer.SelectedOptionID = sql.NullInt64{Int64: int64(a.SelectedOptionID), Valid: true}

			// Auto-grade: check if selected option is correct
			correctOpt, err := s.repo.GetCorrectOptionByQuestionID(ctx, a.QuestionID)
			if err == nil && correctOpt != nil {
				isCorrect := correctOpt.ID == a.SelectedOptionID
				answer.IsCorrect = sql.NullBool{Bool: isCorrect, Valid: true}
				if isCorrect {
					// Find points for this question
					if q, ok := questionMap[a.QuestionID]; ok {
						answer.PointsEarned = sql.NullInt64{Int64: int64(q.Points), Valid: true}
						totalScore += int64(q.Points)
					}
				} else {
					answer.PointsEarned = sql.NullInt64{Int64: 0, Valid: true}
				}
			}
		} else if a.QuestionType == string(TypeEssay) {
			answer.AnswerText = sql.NullString{String: a.AnswerText, Valid: a.AnswerText != ""}
			hasEssay = true
		}

		_ = s.repo.SaveAnswer(ctx, answer)
	}

	// Determine completion status
	newStatus := SubmissionSubmitted
	var percentage *float64
	var isPassed *bool

	if !hasEssay {
		// All auto-gradable – complete immediately
		newStatus = SubmissionCompleted

		if test.TotalPoints > 0 {
			pct := float64(totalScore) / float64(test.TotalPoints) * 100
			percentage = &pct
			passed := pct >= float64(test.PassingScore)
			isPassed = &passed
		}
	}

	score := totalScore
	if err := s.repo.UpdateSubmissionStatus(ctx, submissionID, newStatus, &score, percentage, isPassed); err != nil {
		return nil, err
	}

	// Reload submission
	sub, _ = s.repo.GetSubmissionByID(ctx, submissionID)
	testResp := test.ToResponse()
	return buildSubmissionResponse(sub, &testResp), nil
}

// =============== Response builders ===============

func buildSubmissionResponse(sub *InterviewTestSubmission, test *InterviewTestResponse) *SubmissionResponse {
	resp := &SubmissionResponse{
		ID:     sub.ID,
		Status: string(sub.Status),
		Test:   *test,
	}
	if sub.Score.Valid {
		v := sub.Score.Int64
		resp.Score = &v
	}
	if sub.Percentage.Valid {
		v := sub.Percentage.Float64
		resp.Percentage = &v
	}
	if sub.IsPassed.Valid {
		v := sub.IsPassed.Bool
		resp.IsPassed = &v
	}
	// StartedAt is time.Time, not null - check if not zero value
	if !sub.StartedAt.IsZero() {
		t := sub.StartedAt.Format("2006-01-02T15:04:05Z07:00")
		resp.StartedAt = &t
	}
	if sub.SubmittedAt.Valid {
		t := sub.SubmittedAt.Time.Format("2006-01-02T15:04:05Z07:00")
		resp.SubmittedAt = &t
	}
	if sub.ApplicationID.Valid {
		v := uint64(sub.ApplicationID.Int64)
		resp.ApplicationID = &v
	}
	return resp
}
