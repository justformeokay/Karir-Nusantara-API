package tests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ============================================
// Application Status Transition Tests
// ============================================

// Valid transitions as defined in the system
var validTransitions = map[string][]string{
	"submitted":           {"viewed", "rejected", "withdrawn"},
	"viewed":              {"shortlisted", "rejected", "withdrawn"},
	"shortlisted":         {"interview_scheduled", "rejected", "withdrawn"},
	"interview_scheduled": {"interview_completed", "rejected", "withdrawn"},
	"interview_completed": {"assessment", "offer_sent", "rejected", "withdrawn"},
	"assessment":          {"offer_sent", "rejected", "withdrawn"},
	"offer_sent":          {"offer_accepted", "rejected", "withdrawn"},
	"offer_accepted":      {"hired", "rejected", "withdrawn"},
	"hired":               {},
	"rejected":            {},
	"withdrawn":           {},
}

func TestStatusTransition_SubmittedToViewed(t *testing.T) {
	// CAND-001: submitted → viewed transition allowed
	token := login(t, VerifiedCompany)

	// Get an application with submitted status
	applicationID := getApplicationWithStatus(t, token, "submitted")
	if applicationID == "" {
		t.Skip("No application with submitted status")
	}

	// Update status to viewed
	resp := updateApplicationStatus(t, token, applicationID, "viewed")
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var result ApplicationResponse
	err := json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	assert.True(t, result.Success)
	assert.Equal(t, "viewed", result.Data.Status)
}

func TestStatusTransition_ViewedToShortlisted(t *testing.T) {
	// CAND-002: viewed → shortlisted transition allowed
	token := login(t, VerifiedCompany)

	applicationID := getApplicationWithStatus(t, token, "viewed")
	if applicationID == "" {
		t.Skip("No application with viewed status")
	}

	resp := updateApplicationStatus(t, token, applicationID, "shortlisted")
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var result ApplicationResponse
	json.NewDecoder(resp.Body).Decode(&result)

	assert.True(t, result.Success)
	assert.Equal(t, "shortlisted", result.Data.Status)
}

func TestStatusTransition_ShortlistedToInterviewScheduled(t *testing.T) {
	// CAND-003: shortlisted → interview_scheduled allowed
	token := login(t, VerifiedCompany)

	applicationID := getApplicationWithStatus(t, token, "shortlisted")
	if applicationID == "" {
		t.Skip("No application with shortlisted status")
	}

	resp := updateApplicationStatus(t, token, applicationID, "interview_scheduled")
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var result ApplicationResponse
	json.NewDecoder(resp.Body).Decode(&result)

	assert.True(t, result.Success)
	assert.Equal(t, "interview_scheduled", result.Data.Status)
}

func TestStatusTransition_InterviewScheduledToCompleted(t *testing.T) {
	// CAND-004: interview_scheduled → interview_completed allowed
	token := login(t, VerifiedCompany)

	applicationID := getApplicationWithStatus(t, token, "interview_scheduled")
	if applicationID == "" {
		t.Skip("No application with interview_scheduled status")
	}

	resp := updateApplicationStatus(t, token, applicationID, "interview_completed")
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var result ApplicationResponse
	json.NewDecoder(resp.Body).Decode(&result)

	assert.True(t, result.Success)
	assert.Equal(t, "interview_completed", result.Data.Status)
}

func TestStatusTransition_InterviewCompletedToOfferSent(t *testing.T) {
	// CAND-005: interview_completed → offer_sent allowed
	token := login(t, VerifiedCompany)

	applicationID := getApplicationWithStatus(t, token, "interview_completed")
	if applicationID == "" {
		t.Skip("No application with interview_completed status")
	}

	resp := updateApplicationStatus(t, token, applicationID, "offer_sent")
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var result ApplicationResponse
	json.NewDecoder(resp.Body).Decode(&result)

	assert.True(t, result.Success)
	assert.Equal(t, "offer_sent", result.Data.Status)
}

func TestStatusTransition_OfferSentToAccepted(t *testing.T) {
	// CAND-006: offer_sent → offer_accepted allowed
	token := login(t, VerifiedCompany)

	applicationID := getApplicationWithStatus(t, token, "offer_sent")
	if applicationID == "" {
		t.Skip("No application with offer_sent status")
	}

	resp := updateApplicationStatus(t, token, applicationID, "offer_accepted")
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var result ApplicationResponse
	json.NewDecoder(resp.Body).Decode(&result)

	assert.True(t, result.Success)
	assert.Equal(t, "offer_accepted", result.Data.Status)
}

func TestStatusTransition_OfferAcceptedToHired(t *testing.T) {
	// CAND-007: offer_accepted → hired allowed
	token := login(t, VerifiedCompany)

	applicationID := getApplicationWithStatus(t, token, "offer_accepted")
	if applicationID == "" {
		t.Skip("No application with offer_accepted status")
	}

	resp := updateApplicationStatus(t, token, applicationID, "hired")
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var result ApplicationResponse
	json.NewDecoder(resp.Body).Decode(&result)

	assert.True(t, result.Success)
	assert.Equal(t, "hired", result.Data.Status)
}

func TestStatusTransition_AnyToRejected(t *testing.T) {
	// CAND-008: Any status → rejected allowed
	token := login(t, VerifiedCompany)

	// Test rejection from various statuses
	testCases := []string{"submitted", "viewed", "shortlisted", "interview_scheduled"}

	for _, status := range testCases {
		t.Run(fmt.Sprintf("From_%s", status), func(t *testing.T) {
			applicationID := getApplicationWithStatus(t, token, status)
			if applicationID == "" {
				t.Skipf("No application with %s status", status)
			}

			resp := updateApplicationStatus(t, token, applicationID, "rejected")
			defer resp.Body.Close()

			assert.Equal(t, http.StatusOK, resp.StatusCode)

			var result ApplicationResponse
			json.NewDecoder(resp.Body).Decode(&result)

			assert.True(t, result.Success)
			assert.Equal(t, "rejected", result.Data.Status)
		})
	}
}

// ============================================
// Invalid Status Transition Tests
// ============================================

func TestStatusTransition_SubmittedToHiredNotAllowed(t *testing.T) {
	// CAND-020: submitted → hired directly NOT allowed
	token := login(t, VerifiedCompany)

	applicationID := getApplicationWithStatus(t, token, "submitted")
	if applicationID == "" {
		t.Skip("No application with submitted status")
	}

	resp := updateApplicationStatus(t, token, applicationID, "hired")
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	var result ApplicationResponse
	json.NewDecoder(resp.Body).Decode(&result)

	assert.False(t, result.Success)
	assert.Contains(t, result.Error, "Invalid")
}

func TestStatusTransition_RejectedToAnyNotAllowed(t *testing.T) {
	// CAND-021: rejected → any status NOT allowed
	token := login(t, VerifiedCompany)

	applicationID := getApplicationWithStatus(t, token, "rejected")
	if applicationID == "" {
		t.Skip("No application with rejected status")
	}

	// Try various transitions
	invalidStatuses := []string{"viewed", "shortlisted", "hired"}

	for _, status := range invalidStatuses {
		t.Run(fmt.Sprintf("To_%s", status), func(t *testing.T) {
			resp := updateApplicationStatus(t, token, applicationID, status)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

			var result ApplicationResponse
			json.NewDecoder(resp.Body).Decode(&result)

			assert.False(t, result.Success)
		})
	}
}

func TestStatusTransition_HiredToRejectedNotAllowed(t *testing.T) {
	// CAND-022: hired → rejected NOT allowed
	token := login(t, VerifiedCompany)

	applicationID := getApplicationWithStatus(t, token, "hired")
	if applicationID == "" {
		t.Skip("No application with hired status")
	}

	resp := updateApplicationStatus(t, token, applicationID, "rejected")
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	var result ApplicationResponse
	json.NewDecoder(resp.Body).Decode(&result)

	assert.False(t, result.Success)
}

func TestStatusTransition_WithdrawnToAnyNotAllowed(t *testing.T) {
	// CAND-023: withdrawn → any status NOT allowed
	token := login(t, VerifiedCompany)

	applicationID := getApplicationWithStatus(t, token, "withdrawn")
	if applicationID == "" {
		t.Skip("No application with withdrawn status")
	}

	invalidStatuses := []string{"viewed", "shortlisted", "hired"}

	for _, status := range invalidStatuses {
		t.Run(fmt.Sprintf("To_%s", status), func(t *testing.T) {
			resp := updateApplicationStatus(t, token, applicationID, status)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

			var result ApplicationResponse
			json.NewDecoder(resp.Body).Decode(&result)

			assert.False(t, result.Success)
		})
	}
}

// ============================================
// Job Status Transition Tests
// ============================================

func TestJobStatus_DraftToActive(t *testing.T) {
	// JOB-010: draft → active (publish) allowed
	token := login(t, VerifiedCompany)

	jobID := getJobWithStatus(t, token, "draft")
	if jobID == "" {
		t.Skip("No job with draft status")
	}

	resp := makeAuthenticatedRequest(t, "POST", fmt.Sprintf("/api/v1/company/jobs/%s/publish", jobID), token, nil)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var result JobResponse
	json.NewDecoder(resp.Body).Decode(&result)

	assert.True(t, result.Success)
	assert.Equal(t, "active", result.Data.Status)
}

func TestJobStatus_ActiveToClosed(t *testing.T) {
	// JOB-012: active → closed allowed
	token := login(t, VerifiedCompany)

	jobID := getJobWithStatus(t, token, "active")
	if jobID == "" {
		t.Skip("No job with active status")
	}

	resp := makeAuthenticatedRequest(t, "POST", fmt.Sprintf("/api/v1/company/jobs/%s/close", jobID), token, nil)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var result JobResponse
	json.NewDecoder(resp.Body).Decode(&result)

	assert.True(t, result.Success)
	assert.Equal(t, "closed", result.Data.Status)
}

func TestJobStatus_ClosedToActive(t *testing.T) {
	// JOB-014: closed → active (reopen) allowed
	token := login(t, VerifiedCompany)

	jobID := getJobWithStatus(t, token, "closed")
	if jobID == "" {
		t.Skip("No job with closed status")
	}

	resp := makeAuthenticatedRequest(t, "POST", fmt.Sprintf("/api/v1/company/jobs/%s/reopen", jobID), token, nil)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var result JobResponse
	json.NewDecoder(resp.Body).Decode(&result)

	assert.True(t, result.Success)
	assert.Equal(t, "active", result.Data.Status)
}

// ============================================
// Helper Functions
// ============================================

func getApplicationWithStatus(t *testing.T, token, status string) string {
	resp := makeAuthenticatedRequest(t, "GET", fmt.Sprintf("/api/v1/company/applications?status=%s&limit=1", status), token, nil)
	defer resp.Body.Close()

	var result struct {
		Success bool `json:"success"`
		Data    struct {
			Applications []struct {
				ID string `json:"id"`
			} `json:"applications"`
		} `json:"data"`
	}

	json.NewDecoder(resp.Body).Decode(&result)

	if len(result.Data.Applications) > 0 {
		return result.Data.Applications[0].ID
	}
	return ""
}

func getJobWithStatus(t *testing.T, token, status string) string {
	resp := makeAuthenticatedRequest(t, "GET", fmt.Sprintf("/api/v1/company/jobs?status=%s&limit=1", status), token, nil)
	defer resp.Body.Close()

	var result struct {
		Success bool `json:"success"`
		Data    struct {
			Jobs []struct {
				ID string `json:"id"`
			} `json:"jobs"`
		} `json:"data"`
	}

	json.NewDecoder(resp.Body).Decode(&result)

	if len(result.Data.Jobs) > 0 {
		return result.Data.Jobs[0].ID
	}
	return ""
}

func updateApplicationStatus(t *testing.T, token, applicationID, newStatus string) *http.Response {
	payload := map[string]string{
		"status": newStatus,
	}
	return makeAuthenticatedRequest(t, "PATCH", fmt.Sprintf("/api/v1/company/applications/%s/status", applicationID), token, payload)
}
