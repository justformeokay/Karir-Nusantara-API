package recommendations

import (
	"net/http"
	"strconv"

	"github.com/karirnusantara/api/internal/middleware"
	"github.com/karirnusantara/api/internal/modules/cvs"
	"github.com/karirnusantara/api/internal/modules/jobs"
	"github.com/karirnusantara/api/internal/modules/profile"
	"github.com/karirnusantara/api/internal/shared/response"
)

// Handler handles recommendation HTTP requests
type Handler struct {
	service        *Service
	jobsService    jobs.Service
	cvsService     cvs.Service
	profileService profile.Service
}

// NewHandler creates a new recommendations handler
func NewHandler(
	service *Service,
	jobsService jobs.Service,
	cvsService cvs.Service,
	profileService profile.Service,
) *Handler {
	return &Handler{
		service:        service,
		jobsService:    jobsService,
		cvsService:     cvsService,
		profileService: profileService,
	}
}

// GetRecommendations returns personalized job recommendations
// GET /api/v1/recommendations
func (h *Handler) GetRecommendations(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == 0 {
		response.Unauthorized(w, "Unauthorized")
		return
	}

	// Parse limit from query
	limit := 20
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 50 {
			limit = l
		}
	}

	// Get user's CV
	cv, err := h.cvsService.GetByUserID(r.Context(), userID)
	if err != nil {
		// CV not found is okay, we'll use default profile
		cv = nil
	}

	// Get user's profile
	userProfile, err := h.profileService.GetProfile(r.Context(), userID)
	if err != nil {
		userProfile = nil
	}

	// Build user profile for matching
	// cvsService.GetByUserID already returns *CVResponse
	cvResponse := cv // cv is already *CVResponse from service

	// Extract profile preferences
	preferredLocations := []string{}
	preferredJobTypes := []string{}
	var expectedSalaryMin, expectedSalaryMax int64 = 0, 0
	province := ""

	if userProfile != nil {
		if len(userProfile.PreferredLocations) > 0 {
			preferredLocations = userProfile.PreferredLocations
		}
		if len(userProfile.PreferredJobTypes) > 0 {
			preferredJobTypes = userProfile.PreferredJobTypes
		}
		if userProfile.ExpectedSalaryMin != nil {
			expectedSalaryMin = *userProfile.ExpectedSalaryMin
		}
		if userProfile.ExpectedSalaryMax != nil {
			expectedSalaryMax = *userProfile.ExpectedSalaryMax
		}
		if userProfile.Province != nil {
			province = *userProfile.Province
		}
	}

	// Get user name from context or CV
	userName := ""
	if cvResponse != nil && cvResponse.PersonalInfo.FullName != "" {
		userName = cvResponse.PersonalInfo.FullName
	}

	matchProfile := BuildUserProfile(
		userID,
		userName,
		cvResponse,
		province,
		preferredLocations,
		preferredJobTypes,
		expectedSalaryMin,
		expectedSalaryMax,
	)

	// Get all active jobs
	params := jobs.DefaultJobListParams()
	params.PerPage = 100 // Get more jobs for recommendations
	params.Status = jobs.JobStatusActive

	jobResponses, _, err := h.jobsService.List(r.Context(), params)
	if err != nil {
		response.InternalServerError(w, "Failed to get jobs")
		return
	}

	// Service.List already returns []*JobResponse, no conversion needed

	// Get recommendations
	recommendations := h.service.GetRecommendations(matchProfile, jobResponses, limit)

	response.Success(w, http.StatusOK, "Recommendations retrieved successfully", recommendations)
}
