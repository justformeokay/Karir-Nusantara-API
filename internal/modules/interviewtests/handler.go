package interviewtests

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/karirnusantara/api/internal/shared/response"
)

// Handler handles HTTP requests for interview tests
type Handler struct {
	service *Service
}

// NewHandler creates a new interview tests handler
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// Create handles POST /admin/interview-tests
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateInterviewTestRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}

	// Get admin ID from context
	adminID := getAdminIDFromContext(r)

	result, err := h.service.Create(r.Context(), req, adminID)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	response.Success(w, http.StatusCreated, "Interview test created successfully", result)
}

// GetByID handles GET /admin/interview-tests/{id}
func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_ID", "Invalid test ID")
		return
	}

	result, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	response.Success(w, http.StatusOK, "Interview test retrieved successfully", result)
}

// GetAll handles GET /admin/interview-tests
func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")

	results, err := h.service.GetAll(r.Context(), status)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	response.Success(w, http.StatusOK, "Interview tests retrieved successfully", results)
}

// Update handles PUT /admin/interview-tests/{id}
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_ID", "Invalid test ID")
		return
	}

	var req UpdateInterviewTestRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}

	// Get admin ID from context
	adminID := getAdminIDFromContext(r)

	result, err := h.service.Update(r.Context(), id, req, adminID)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	response.Success(w, http.StatusOK, "Interview test updated successfully", result)
}

// Delete handles DELETE /admin/interview-tests/{id}
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_ID", "Invalid test ID")
		return
	}

	err = h.service.Delete(r.Context(), id)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	response.Success(w, http.StatusOK, "Interview test deleted successfully", nil)
}

// Publish handles POST /admin/interview-tests/{id}/publish
func (h *Handler) Publish(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_ID", "Invalid test ID")
		return
	}

	// Get admin ID from context
	adminID := getAdminIDFromContext(r)

	result, err := h.service.Publish(r.Context(), id, adminID)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	response.Success(w, http.StatusOK, "Interview test published successfully", result)
}

// Duplicate handles POST /admin/interview-tests/{id}/duplicate
func (h *Handler) Duplicate(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_ID", "Invalid test ID")
		return
	}

	// Get admin ID from context
	adminID := getAdminIDFromContext(r)

	result, err := h.service.Duplicate(r.Context(), id, adminID)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	response.Success(w, http.StatusCreated, "Interview test duplicated successfully", result)
}

// =============== Helper functions ===============

// getAdminIDFromContext retrieves admin ID from request context
func getAdminIDFromContext(r *http.Request) uint64 {
	userID := r.Context().Value("user_id")
	if userID == nil {
		return 0
	}

	// Try to convert to uint64
	switch v := userID.(type) {
	case uint64:
		return v
	case int64:
		return uint64(v)
	case float64:
		return uint64(v)
	case string:
		id, _ := strconv.ParseUint(v, 10, 64)
		return id
	default:
		return 0
	}
}

// handleServiceError handles service layer errors
func handleServiceError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrTestNotFound):
		response.Error(w, http.StatusNotFound, "NOT_FOUND", err.Error())
	case errors.Is(err, ErrInvalidStatus):
		response.Error(w, http.StatusBadRequest, "INVALID_STATUS", err.Error())
	case errors.Is(err, ErrInvalidType):
		response.Error(w, http.StatusBadRequest, "INVALID_TYPE", err.Error())
	case errors.Is(err, ErrInvalidDifficulty):
		response.Error(w, http.StatusBadRequest, "INVALID_DIFFICULTY", err.Error())
	case errors.Is(err, ErrNoQuestions):
		response.Error(w, http.StatusBadRequest, "NO_QUESTIONS", err.Error())
	case errors.Is(err, ErrInvalidOptions):
		response.Error(w, http.StatusBadRequest, "INVALID_OPTIONS", err.Error())
	case errors.Is(err, context.DeadlineExceeded):
		response.Error(w, http.StatusRequestTimeout, "TIMEOUT", "Request timeout")
	case errors.Is(err, context.Canceled):
		response.Error(w, http.StatusRequestTimeout, "CANCELED", "Request canceled")
	default:
		response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
	}
}

// =============== Company Handler ===============

// CompanyServiceInterface defines the interface for company service used by company interview test handler
type CompanyServiceInterface interface {
	GetCompanyIDByUserID(ctx context.Context, userID uint64) (uint64, error)
}

// CompanyHandler handles HTTP requests for company interview tests
type CompanyHandler struct {
	service        *Service
	companyService CompanyServiceInterface
}

// NewCompanyHandler creates a new company interview tests handler
func NewCompanyHandler(service *Service, companyService CompanyServiceInterface) *CompanyHandler {
	return &CompanyHandler{service: service, companyService: companyService}
}

// getCompanyIDFromContext resolves company ID for the authenticated user
func (h *CompanyHandler) getCompanyIDFromContext(r *http.Request) (uint64, uint64, error) {
	userID := getAdminIDFromContext(r)
	if userID == 0 {
		return 0, 0, errors.New("unauthorized")
	}
	companyID, err := h.companyService.GetCompanyIDByUserID(r.Context(), userID)
	if err != nil {
		return 0, 0, err
	}
	if companyID == 0 {
		return 0, 0, errors.New("company not found")
	}
	return userID, companyID, nil
}

// GetLibrary handles GET /company/interview-tests/library
func (h *CompanyHandler) GetLibrary(w http.ResponseWriter, r *http.Request) {
	results, err := h.service.GetPublicAdminTests(r.Context())
	if err != nil {
		handleServiceError(w, err)
		return
	}
	response.Success(w, http.StatusOK, "Library retrieved successfully", results)
}

// GetMyTests handles GET /company/interview-tests
func (h *CompanyHandler) GetMyTests(w http.ResponseWriter, r *http.Request) {
	_, companyID, err := h.getCompanyIDFromContext(r)
	if err != nil {
		response.Error(w, http.StatusUnauthorized, "UNAUTHORIZED", err.Error())
		return
	}

	status := r.URL.Query().Get("status")
	results, err := h.service.GetByCompanyID(r.Context(), companyID, status)
	if err != nil {
		handleServiceError(w, err)
		return
	}
	response.Success(w, http.StatusOK, "Interview tests retrieved successfully", results)
}

// CreateMyTest handles POST /company/interview-tests
func (h *CompanyHandler) CreateMyTest(w http.ResponseWriter, r *http.Request) {
	userID, companyID, err := h.getCompanyIDFromContext(r)
	if err != nil {
		response.Error(w, http.StatusUnauthorized, "UNAUTHORIZED", err.Error())
		return
	}

	var req CreateInterviewTestRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}

	result, err := h.service.CreateForCompany(r.Context(), req, companyID, userID)
	if err != nil {
		handleServiceError(w, err)
		return
	}
	response.Success(w, http.StatusCreated, "Interview test created successfully", result)
}

// UpdateMyTest handles PUT /company/interview-tests/{id}
func (h *CompanyHandler) UpdateMyTest(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_ID", "Invalid test ID")
		return
	}

	userID, companyID, err := h.getCompanyIDFromContext(r)
	if err != nil {
		response.Error(w, http.StatusUnauthorized, "UNAUTHORIZED", err.Error())
		return
	}

	var req UpdateInterviewTestRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}

	result, err := h.service.UpdateForCompany(r.Context(), id, req, companyID, userID)
	if err != nil {
		handleServiceError(w, err)
		return
	}
	response.Success(w, http.StatusOK, "Interview test updated successfully", result)
}

// DeleteMyTest handles DELETE /company/interview-tests/{id}
func (h *CompanyHandler) DeleteMyTest(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_ID", "Invalid test ID")
		return
	}

	_, companyID, err := h.getCompanyIDFromContext(r)
	if err != nil {
		response.Error(w, http.StatusUnauthorized, "UNAUTHORIZED", err.Error())
		return
	}

	if err := h.service.DeleteForCompany(r.Context(), id, companyID); err != nil {
		handleServiceError(w, err)
		return
	}
	response.Success(w, http.StatusOK, "Interview test deleted successfully", nil)
}

// PublishMyTest handles POST /company/interview-tests/{id}/publish
func (h *CompanyHandler) PublishMyTest(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_ID", "Invalid test ID")
		return
	}

	userID, companyID, err := h.getCompanyIDFromContext(r)
	if err != nil {
		response.Error(w, http.StatusUnauthorized, "UNAUTHORIZED", err.Error())
		return
	}

	result, err := h.service.PublishForCompany(r.Context(), id, companyID, userID)
	if err != nil {
		handleServiceError(w, err)
		return
	}
	response.Success(w, http.StatusOK, "Interview test published successfully", result)
}

// CopyFromAdmin handles POST /company/interview-tests/{id}/copy
func (h *CompanyHandler) CopyFromAdmin(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_ID", "Invalid test ID")
		return
	}

	userID, companyID, err := h.getCompanyIDFromContext(r)
	if err != nil {
		response.Error(w, http.StatusUnauthorized, "UNAUTHORIZED", err.Error())
		return
	}

	result, err := h.service.CopyFromAdmin(r.Context(), id, companyID, userID)
	if err != nil {
		handleServiceError(w, err)
		return
	}
	response.Success(w, http.StatusCreated, "Test copied successfully", result)
}

// AssignTest handles POST /company/applications/{applicationId}/assign-test
func (h *CompanyHandler) AssignTest(w http.ResponseWriter, r *http.Request) {
	appIDStr := chi.URLParam(r, "applicationId")
	applicationID, err := strconv.ParseUint(appIDStr, 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_ID", "Invalid application ID")
		return
	}

	_, companyID, err := h.getCompanyIDFromContext(r)
	if err != nil {
		response.Error(w, http.StatusUnauthorized, "UNAUTHORIZED", err.Error())
		return
	}

	var req struct {
		InterviewTestID uint64 `json:"interview_test_id"`
		CandidateUserID uint64 `json:"candidate_user_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}
	if req.InterviewTestID == 0 || req.CandidateUserID == 0 {
		response.Error(w, http.StatusBadRequest, "MISSING_FIELDS", "interview_test_id and candidate_user_id are required")
		return
	}

	result, err := h.service.AssignTestToCandidate(r.Context(), req.InterviewTestID, req.CandidateUserID, applicationID, companyID)
	if err != nil {
		handleServiceError(w, err)
		return
	}
	response.Success(w, http.StatusCreated, "Test assigned successfully", result)
}

// GetApplicationSubmissions handles GET /company/applications/{applicationId}/interview-tests
func (h *CompanyHandler) GetApplicationSubmissions(w http.ResponseWriter, r *http.Request) {
	appIDStr := chi.URLParam(r, "applicationId")
	applicationID, err := strconv.ParseUint(appIDStr, 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_ID", "Invalid application ID")
		return
	}

	_, _, err = h.getCompanyIDFromContext(r)
	if err != nil {
		response.Error(w, http.StatusUnauthorized, "UNAUTHORIZED", err.Error())
		return
	}

	results, err := h.service.GetSubmissionsForApplication(r.Context(), applicationID)
	if err != nil {
		handleServiceError(w, err)
		return
	}
	response.Success(w, http.StatusOK, "Submissions retrieved successfully", results)
}

// =============== Job Seeker Handler ===============

// JobSeekerHandler handles HTTP requests for job seeker test activities
type JobSeekerHandler struct {
	service *Service
}

// NewJobSeekerHandler creates a new job seeker interview test handler
func NewJobSeekerHandler(service *Service) *JobSeekerHandler {
	return &JobSeekerHandler{service: service}
}

// GetMySubmissions handles GET /jobseeker/interview-tests
func (h *JobSeekerHandler) GetMySubmissions(w http.ResponseWriter, r *http.Request) {
	userID := getAdminIDFromContext(r)
	if userID == 0 {
		response.Error(w, http.StatusUnauthorized, "UNAUTHORIZED", "Unauthorized")
		return
	}

	results, err := h.service.GetSubmissionsForUser(r.Context(), userID)
	if err != nil {
		handleServiceError(w, err)
		return
	}
	response.Success(w, http.StatusOK, "Submissions retrieved successfully", results)
}

// GetTestForSubmission handles GET /jobseeker/interview-tests/{submissionId}
func (h *JobSeekerHandler) GetTestForSubmission(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "submissionId")
	submissionID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_ID", "Invalid submission ID")
		return
	}

	userID := getAdminIDFromContext(r)
	if userID == 0 {
		response.Error(w, http.StatusUnauthorized, "UNAUTHORIZED", "Unauthorized")
		return
	}

	result, err := h.service.GetTestForSubmission(r.Context(), submissionID, userID)
	if err != nil {
		response.Error(w, http.StatusNotFound, "NOT_FOUND", err.Error())
		return
	}
	response.Success(w, http.StatusOK, "Test retrieved successfully", result)
}

// SubmitAnswers handles POST /jobseeker/interview-tests/{submissionId}/submit
func (h *JobSeekerHandler) SubmitAnswers(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "submissionId")
	submissionID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_ID", "Invalid submission ID")
		return
	}

	userID := getAdminIDFromContext(r)
	if userID == 0 {
		response.Error(w, http.StatusUnauthorized, "UNAUTHORIZED", "Unauthorized")
		return
	}

	var req SubmitAnswersRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}

	result, err := h.service.SubmitTestAnswers(r.Context(), submissionID, userID, req.Answers)
	if err != nil {
		handleServiceError(w, err)
		return
	}
	response.Success(w, http.StatusOK, "Test submitted successfully", result)
}
