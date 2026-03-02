package interviewtests

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
)

// Repository handles database operations for interview tests
type Repository struct {
	db *sqlx.DB
}

// NewRepository creates a new interview tests repository
func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

// =============== Interview Test CRUD ===============

// Create inserts a new interview test
func (r *Repository) Create(ctx context.Context, test *InterviewTest) error {
	query := `
		INSERT INTO interview_tests (
			title, description, duration_minutes, passing_score,
			shuffle_questions, show_results_immediately, status,
			owner_type, owner_company_id, is_public, created_by,
			created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW())
	`

	result, err := r.db.ExecContext(ctx, query,
		test.Title,
		test.Description,
		test.DurationMinutes,
		test.PassingScore,
		test.ShuffleQuestions,
		test.ShowResultsImmediately,
		test.Status,
		test.OwnerType,
		test.OwnerCompanyID,
		test.IsPublic,
		test.CreatedBy,
	)
	if err != nil {
		return fmt.Errorf("failed to create interview test: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}
	test.ID = uint64(id)

	return nil
}

// GetByID retrieves an interview test by ID
func (r *Repository) GetByID(ctx context.Context, id uint64) (*InterviewTest, error) {
	var test InterviewTest
	query := `SELECT * FROM interview_tests WHERE id = ? AND deleted_at IS NULL`

	err := r.db.GetContext(ctx, &test, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get interview test: %w", err)
	}

	return &test, nil
}

// GetAll retrieves all interview tests
func (r *Repository) GetAll(ctx context.Context, status string) ([]*InterviewTest, error) {
	var tests []*InterviewTest

	query := `SELECT * FROM interview_tests WHERE deleted_at IS NULL`
	args := []interface{}{}

	if status != "" {
		query += ` AND status = ?`
		args = append(args, status)
	}

	query += ` ORDER BY created_at DESC`

	err := r.db.SelectContext(ctx, &tests, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get interview tests: %w", err)
	}

	return tests, nil
}

// Update updates an existing interview test
func (r *Repository) Update(ctx context.Context, test *InterviewTest) error {
	query := `
		UPDATE interview_tests SET
			title = ?,
			description = ?,
			duration_minutes = ?,
			passing_score = ?,
			shuffle_questions = ?,
			show_results_immediately = ?,
			status = ?,
			updated_by = ?,
			updated_at = NOW()
		WHERE id = ? AND deleted_at IS NULL
	`

	_, err := r.db.ExecContext(ctx, query,
		test.Title,
		test.Description,
		test.DurationMinutes,
		test.PassingScore,
		test.ShuffleQuestions,
		test.ShowResultsImmediately,
		test.Status,
		test.UpdatedBy,
		test.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update interview test: %w", err)
	}

	return nil
}

// Delete soft deletes an interview test
func (r *Repository) Delete(ctx context.Context, id uint64) error {
	query := `UPDATE interview_tests SET deleted_at = NOW() WHERE id = ?`

	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete interview test: %w", err)
	}

	return nil
}

// UpdateStatus updates the status of an interview test
func (r *Repository) UpdateStatus(ctx context.Context, id uint64, status InterviewTestStatus, adminID uint64) error {
	query := `
		UPDATE interview_tests 
		SET status = ?, updated_by = ?, updated_at = NOW() 
		WHERE id = ? AND deleted_at IS NULL
	`

	_, err := r.db.ExecContext(ctx, query, status, adminID, id)
	if err != nil {
		return fmt.Errorf("failed to update status: %w", err)
	}

	return nil
}

// =============== Question CRUD ===============

// CreateQuestion inserts a new question
func (r *Repository) CreateQuestion(ctx context.Context, question *InterviewQuestion) error {
	query := `
		INSERT INTO interview_questions (
			interview_test_id, question_text, question_type, points,
			difficulty, ` + "`order`" + `, explanation, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, NOW(), NOW())
	`

	result, err := r.db.ExecContext(ctx, query,
		question.InterviewTestID,
		question.QuestionText,
		question.QuestionType,
		question.Points,
		question.Difficulty,
		question.Order,
		question.Explanation,
	)
	if err != nil {
		return fmt.Errorf("failed to create question: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}
	question.ID = uint64(id)

	return nil
}

// GetQuestionsByTestID retrieves all questions for a test
func (r *Repository) GetQuestionsByTestID(ctx context.Context, testID uint64) ([]*InterviewQuestion, error) {
	var questions []*InterviewQuestion
	query := `
		SELECT * FROM interview_questions 
		WHERE interview_test_id = ? 
		ORDER BY ` + "`order`" + ` ASC
	`

	err := r.db.SelectContext(ctx, &questions, query, testID)
	if err != nil {
		return nil, fmt.Errorf("failed to get questions: %w", err)
	}

	return questions, nil
}

// DeleteQuestionsByTestID deletes all questions for a test
func (r *Repository) DeleteQuestionsByTestID(ctx context.Context, testID uint64) error {
	query := `DELETE FROM interview_questions WHERE interview_test_id = ?`

	_, err := r.db.ExecContext(ctx, query, testID)
	if err != nil {
		return fmt.Errorf("failed to delete questions: %w", err)
	}

	return nil
}

// =============== Option CRUD ===============

// CreateOption inserts a new option
func (r *Repository) CreateOption(ctx context.Context, option *QuestionOption) error {
	query := `
		INSERT INTO interview_question_options (
			interview_question_id, option_text, is_correct, ` + "`order`" + `,
			created_at, updated_at
		) VALUES (?, ?, ?, ?, NOW(), NOW())
	`

	result, err := r.db.ExecContext(ctx, query,
		option.InterviewQuestionID,
		option.OptionText,
		option.IsCorrect,
		option.Order,
	)
	if err != nil {
		return fmt.Errorf("failed to create option: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}
	option.ID = uint64(id)

	return nil
}

// GetOptionsByQuestionID retrieves all options for a question
func (r *Repository) GetOptionsByQuestionID(ctx context.Context, questionID uint64) ([]*QuestionOption, error) {
	var options []*QuestionOption
	query := `
		SELECT * FROM interview_question_options 
		WHERE interview_question_id = ? 
		ORDER BY ` + "`order`" + ` ASC
	`

	err := r.db.SelectContext(ctx, &options, query, questionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get options: %w", err)
	}

	return options, nil
}

// GetOptionsByTestID retrieves all options for all questions in a test
func (r *Repository) GetOptionsByTestID(ctx context.Context, testID uint64) (map[uint64][]*QuestionOption, error) {
	var options []*QuestionOption
	query := `
		SELECT o.* FROM interview_question_options o
		INNER JOIN interview_questions q ON o.interview_question_id = q.id
		WHERE q.interview_test_id = ?
		ORDER BY o.interview_question_id, o.` + "`order`" + ` ASC
	`

	err := r.db.SelectContext(ctx, &options, query, testID)
	if err != nil {
		return nil, fmt.Errorf("failed to get options: %w", err)
	}

	// Group options by question ID
	optionsMap := make(map[uint64][]*QuestionOption)
	for _, option := range options {
		optionsMap[option.InterviewQuestionID] = append(optionsMap[option.InterviewQuestionID], option)
	}

	return optionsMap, nil
}

// DeleteOptionsByQuestionID deletes all options for a question
func (r *Repository) DeleteOptionsByQuestionID(ctx context.Context, questionID uint64) error {
	query := `DELETE FROM interview_question_options WHERE interview_question_id = ?`

	_, err := r.db.ExecContext(ctx, query, questionID)
	if err != nil {
		return fmt.Errorf("failed to delete options: %w", err)
	}

	return nil
}

// GetPublicAdminTests retrieves all public tests created by super_admin
func (r *Repository) GetPublicAdminTests(ctx context.Context) ([]*InterviewTest, error) {
	var tests []*InterviewTest
	query := `
		SELECT * FROM interview_tests 
		WHERE deleted_at IS NULL 
		  AND owner_type = 'super_admin' 
		  AND is_public = 1 
		  AND status = 'active'
		ORDER BY created_at DESC
	`

	err := r.db.SelectContext(ctx, &tests, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get public admin tests: %w", err)
	}

	return tests, nil
}

// GetByCompanyID retrieves all tests owned by a specific company
func (r *Repository) GetByCompanyID(ctx context.Context, companyID uint64, status string) ([]*InterviewTest, error) {
	var tests []*InterviewTest

	query := `
		SELECT * FROM interview_tests 
		WHERE deleted_at IS NULL 
		  AND owner_type = 'company' 
		  AND owner_company_id = ?
	`
	args := []interface{}{companyID}

	if status != "" {
		query += ` AND status = ?`
		args = append(args, status)
	}

	query += ` ORDER BY created_at DESC`

	err := r.db.SelectContext(ctx, &tests, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get company tests: %w", err)
	}

	return tests, nil
}
