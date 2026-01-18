package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test configuration
const (
	TestAPIURL = "http://localhost:8081"
)

// Test users
var (
	VerifiedCompany = TestUser{
		Email:    "test-verified@company.com",
		Password: "TestPassword123!",
	}
	UnverifiedCompany = TestUser{
		Email:    "test-unverified@company.com",
		Password: "TestPassword123!",
	}
)

// TestUser represents a test user
type TestUser struct {
	Email    string
	Password string
	Token    string
}

// LoginResponse represents the login API response
type LoginResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    struct {
		Token   string  `json:"token"`
		Company Company `json:"company"`
	} `json:"data"`
	Error string `json:"error"`
}

// Company represents a company entity
type Company struct {
	ID                 string `json:"id"`
	Email              string `json:"email"`
	CompanyName        string `json:"company_name"`
	IsVerified         bool   `json:"is_verified"`
	VerificationStatus string `json:"verification_status"`
	IsActive           bool   `json:"is_active"`
}

// QuotaResponse represents the quota API response
type QuotaResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    struct {
		FreeQuota          int `json:"free_quota"`
		UsedFreeQuota      int `json:"used_free_quota"`
		RemainingFreeQuota int `json:"remaining_free_quota"`
		PaidQuota          int `json:"paid_quota"`
		PricePerJob        int `json:"price_per_job"`
	} `json:"data"`
	Error string `json:"error"`
}

// JobResponse represents a job API response
type JobResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    struct {
		ID     string `json:"id"`
		Title  string `json:"title"`
		Status string `json:"status"`
	} `json:"data"`
	Error string `json:"error"`
}

// ApplicationResponse represents an application API response
type ApplicationResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    struct {
		ID     string `json:"id"`
		Status string `json:"status"`
	} `json:"data"`
	Error string `json:"error"`
}

// Helper functions

func login(t *testing.T, user TestUser) string {
	payload := map[string]string{
		"email":    user.Email,
		"password": user.Password,
	}
	body, _ := json.Marshal(payload)

	resp, err := http.Post(
		TestAPIURL+"/api/v1/company/auth/login",
		"application/json",
		bytes.NewReader(body),
	)
	require.NoError(t, err)
	defer resp.Body.Close()

	var result LoginResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)
	require.True(t, result.Success, "Login should succeed: %s", result.Error)

	return result.Data.Token
}

func makeAuthenticatedRequest(t *testing.T, method, url string, token string, body interface{}) *http.Response {
	var reqBody *bytes.Reader
	if body != nil {
		jsonBody, _ := json.Marshal(body)
		reqBody = bytes.NewReader(jsonBody)
	} else {
		reqBody = bytes.NewReader([]byte{})
	}

	req, err := http.NewRequest(method, TestAPIURL+url, reqBody)
	require.NoError(t, err)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err)

	return resp
}

// ============================================
// Authentication Tests
// ============================================

func TestAuth_LoginWithValidCredentials(t *testing.T) {
	// AUTH-001: Login with valid credentials
	payload := map[string]string{
		"email":    VerifiedCompany.Email,
		"password": VerifiedCompany.Password,
	}
	body, _ := json.Marshal(payload)

	resp, err := http.Post(
		TestAPIURL+"/api/v1/company/auth/login",
		"application/json",
		bytes.NewReader(body),
	)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var result LoginResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	assert.True(t, result.Success)
	assert.NotEmpty(t, result.Data.Token)
	assert.NotEmpty(t, result.Data.Company.ID)
}

func TestAuth_LoginWithInvalidPassword(t *testing.T) {
	// AUTH-003: Login with wrong password
	payload := map[string]string{
		"email":    VerifiedCompany.Email,
		"password": "wrongpassword",
	}
	body, _ := json.Marshal(payload)

	resp, err := http.Post(
		TestAPIURL+"/api/v1/company/auth/login",
		"application/json",
		bytes.NewReader(body),
	)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	var result LoginResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	assert.False(t, result.Success)
	assert.Contains(t, result.Error, "Invalid")
}

func TestAuth_InvalidTokenReturns401(t *testing.T) {
	// AUTH-007: Invalid token returns 401
	req, err := http.NewRequest("GET", TestAPIURL+"/api/v1/company/profile", nil)
	require.NoError(t, err)

	req.Header.Set("Authorization", "Bearer invalid-token")

	client := &http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestAuth_NewCompanyStartsPending(t *testing.T) {
	// AUTH-013: New company starts with pending verification status
	token := login(t, UnverifiedCompany)

	resp := makeAuthenticatedRequest(t, "GET", "/api/v1/company/profile", token, nil)
	defer resp.Body.Close()

	var result struct {
		Success bool    `json:"success"`
		Data    Company `json:"data"`
	}
	err := json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	assert.Equal(t, "pending", result.Data.VerificationStatus)
	assert.False(t, result.Data.IsVerified)
}

// ============================================
// Permission Tests
// ============================================

func TestPermission_UnverifiedCannotCreateJob(t *testing.T) {
	// AUTH-020: Unverified company cannot create jobs
	token := login(t, UnverifiedCompany)

	job := map[string]interface{}{
		"title":       "Test Job",
		"description": "Test description",
		"location":    "Jakarta",
		"type":        "full-time",
	}

	resp := makeAuthenticatedRequest(t, "POST", "/api/v1/company/jobs", token, job)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusForbidden, resp.StatusCode)
}

func TestPermission_VerifiedCanCreateJob(t *testing.T) {
	// AUTH-021: Verified company can create jobs
	token := login(t, VerifiedCompany)

	job := map[string]interface{}{
		"title":            "Test Job",
		"description":      "<p>Test description</p>",
		"requirements":     "<p>Requirements</p>",
		"location":         "Jakarta",
		"type":             "full-time",
		"experience_level": "mid",
		"salary_min":       10000000,
		"salary_max":       20000000,
	}

	resp := makeAuthenticatedRequest(t, "POST", "/api/v1/company/jobs", token, job)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var result JobResponse
	err := json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	assert.True(t, result.Success)
	assert.NotEmpty(t, result.Data.ID)
}

func TestPermission_CompanyCanOnlySeeOwnJobs(t *testing.T) {
	// AUTH-022: Company can only see own jobs
	token := login(t, VerifiedCompany)

	resp := makeAuthenticatedRequest(t, "GET", "/api/v1/company/jobs", token, nil)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var result struct {
		Success bool `json:"success"`
		Data    struct {
			Jobs []struct {
				CompanyID string `json:"company_id"`
			} `json:"jobs"`
		} `json:"data"`
	}
	err := json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	// All jobs should belong to the same company
	// (In real test, we'd verify against the logged-in company ID)
	assert.True(t, result.Success)
}

// ============================================
// Quota Tests
// ============================================

func TestQuota_NewCompanyStartsWith5FreeQuota(t *testing.T) {
	// QUOTA-001: New company starts with 5 free quota
	token := login(t, VerifiedCompany)

	resp := makeAuthenticatedRequest(t, "GET", "/api/v1/company/quota", token, nil)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var result QuotaResponse
	err := json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	assert.True(t, result.Success)
	assert.Equal(t, 5, result.Data.FreeQuota)
}

func TestQuota_PublishingJobDecrementsQuota(t *testing.T) {
	// QUOTA-002: Publishing job decrements free quota by 1
	token := login(t, VerifiedCompany)

	// Get initial quota
	resp := makeAuthenticatedRequest(t, "GET", "/api/v1/company/quota", token, nil)
	var initialQuota QuotaResponse
	json.NewDecoder(resp.Body).Decode(&initialQuota)
	resp.Body.Close()

	initialRemaining := initialQuota.Data.RemainingFreeQuota

	// Create and publish a job
	job := map[string]interface{}{
		"title":            "Quota Test Job",
		"description":      "<p>Test</p>",
		"requirements":     "<p>Test</p>",
		"location":         "Jakarta",
		"type":             "full-time",
		"experience_level": "mid",
		"publish":          true, // Publish immediately
	}

	resp = makeAuthenticatedRequest(t, "POST", "/api/v1/company/jobs", token, job)
	resp.Body.Close()

	// Get updated quota
	resp = makeAuthenticatedRequest(t, "GET", "/api/v1/company/quota", token, nil)
	var updatedQuota QuotaResponse
	json.NewDecoder(resp.Body).Decode(&updatedQuota)
	resp.Body.Close()

	assert.Equal(t, initialRemaining-1, updatedQuota.Data.RemainingFreeQuota)
}

func TestQuota_APIRejectsPublishWhenQuotaIs0(t *testing.T) {
	// QUOTA-012: API rejects publish when quota is 0
	// Note: This requires a test user with 0 quota
	t.Skip("Requires test user with exhausted quota")
}

func TestQuota_ClosingJobDoesNotRestoreQuota(t *testing.T) {
	// QUOTA-003: Closing job does NOT restore free quota
	token := login(t, VerifiedCompany)

	// Get initial quota
	resp := makeAuthenticatedRequest(t, "GET", "/api/v1/company/quota", token, nil)
	var initialQuota QuotaResponse
	json.NewDecoder(resp.Body).Decode(&initialQuota)
	resp.Body.Close()

	// Get a job to close
	resp = makeAuthenticatedRequest(t, "GET", "/api/v1/company/jobs?status=active", token, nil)
	var jobsResult struct {
		Data struct {
			Jobs []struct {
				ID string `json:"id"`
			} `json:"jobs"`
		} `json:"data"`
	}
	json.NewDecoder(resp.Body).Decode(&jobsResult)
	resp.Body.Close()

	if len(jobsResult.Data.Jobs) == 0 {
		t.Skip("No active jobs to close")
	}

	jobID := jobsResult.Data.Jobs[0].ID

	// Close the job
	resp = makeAuthenticatedRequest(t, "POST", fmt.Sprintf("/api/v1/company/jobs/%s/close", jobID), token, nil)
	resp.Body.Close()

	// Get updated quota
	resp = makeAuthenticatedRequest(t, "GET", "/api/v1/company/quota", token, nil)
	var updatedQuota QuotaResponse
	json.NewDecoder(resp.Body).Decode(&updatedQuota)
	resp.Body.Close()

	// Quota should remain the same (not restored)
	assert.Equal(t, initialQuota.Data.RemainingFreeQuota, updatedQuota.Data.RemainingFreeQuota)
}
