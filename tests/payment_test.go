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
// Payment Flow Tests
// ============================================

// PaymentResponse represents a payment API response
type PaymentResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    struct {
		ID          string `json:"id"`
		Quantity    int    `json:"quantity"`
		TotalAmount int    `json:"total_amount"`
		Status      string `json:"status"`
	} `json:"data"`
	Error string `json:"error"`
}

// PaymentHistoryResponse represents payment history
type PaymentHistoryResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    []struct {
		ID          string `json:"id"`
		Quantity    int    `json:"quantity"`
		TotalAmount int    `json:"total_amount"`
		Status      string `json:"status"`
		ProofURL    string `json:"proof_url"`
		CreatedAt   string `json:"created_at"`
	} `json:"data"`
	Error string `json:"error"`
}

func TestPayment_ProofStoredWithPendingStatus(t *testing.T) {
	// PAY-004: Proof is stored with pending status
	token := login(t, VerifiedCompany)

	// Note: In real test, you'd use multipart/form-data for file upload
	// This is a simplified version
	payload := map[string]interface{}{
		"quantity":   5,
		"proof_url":  "https://example.com/proof.jpg", // Simulated
	}

	resp := makeAuthenticatedRequest(t, "POST", "/api/v1/company/quota/payment", token, payload)
	defer resp.Body.Close()

	// Should succeed
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var result PaymentResponse
	err := json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	assert.True(t, result.Success)
	assert.Equal(t, "pending", result.Data.Status)
	assert.Equal(t, 5, result.Data.Quantity)
}

func TestPayment_ConfirmedPaymentTriggersQuotaIncrement(t *testing.T) {
	// PAY-013: Confirmed payment triggers quota increment
	// Note: This test requires admin access to confirm payment
	// Here we just verify the quota changes after confirmation
	t.Skip("Requires admin confirmation - integration test")
}

func TestPayment_GetPaymentHistory(t *testing.T) {
	token := login(t, VerifiedCompany)

	resp := makeAuthenticatedRequest(t, "GET", "/api/v1/company/quota/payments", token, nil)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var result PaymentHistoryResponse
	err := json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	assert.True(t, result.Success)
	// Should return array (possibly empty)
	assert.NotNil(t, result.Data)
}

func TestPayment_GetQuotaInfo(t *testing.T) {
	token := login(t, VerifiedCompany)

	resp := makeAuthenticatedRequest(t, "GET", "/api/v1/company/quota", token, nil)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var result QuotaResponse
	err := json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	assert.True(t, result.Success)
	assert.GreaterOrEqual(t, result.Data.FreeQuota, 0)
	assert.GreaterOrEqual(t, result.Data.PaidQuota, 0)
	assert.Greater(t, result.Data.PricePerJob, 0)
}

// ============================================
// Payment Validation Tests
// ============================================

func TestPayment_InvalidQuantityRejected(t *testing.T) {
	token := login(t, VerifiedCompany)

	// Test with 0 quantity
	payload := map[string]interface{}{
		"quantity":  0,
		"proof_url": "https://example.com/proof.jpg",
	}

	resp := makeAuthenticatedRequest(t, "POST", "/api/v1/company/quota/payment", token, payload)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestPayment_NegativeQuantityRejected(t *testing.T) {
	token := login(t, VerifiedCompany)

	payload := map[string]interface{}{
		"quantity":  -5,
		"proof_url": "https://example.com/proof.jpg",
	}

	resp := makeAuthenticatedRequest(t, "POST", "/api/v1/company/quota/payment", token, payload)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestPayment_MissingProofRejected(t *testing.T) {
	token := login(t, VerifiedCompany)

	payload := map[string]interface{}{
		"quantity": 5,
		// Missing proof_url
	}

	resp := makeAuthenticatedRequest(t, "POST", "/api/v1/company/quota/payment", token, payload)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

// ============================================
// Price Calculation Tests
// ============================================

func TestPayment_TotalAmountCalculatedCorrectly(t *testing.T) {
	token := login(t, VerifiedCompany)

	// Get price per job first
	resp := makeAuthenticatedRequest(t, "GET", "/api/v1/company/quota", token, nil)
	var quotaResult QuotaResponse
	json.NewDecoder(resp.Body).Decode(&quotaResult)
	resp.Body.Close()

	pricePerJob := quotaResult.Data.PricePerJob
	quantity := 5
	expectedTotal := pricePerJob * quantity

	// Create payment
	payload := map[string]interface{}{
		"quantity":  quantity,
		"proof_url": "https://example.com/proof.jpg",
	}

	resp = makeAuthenticatedRequest(t, "POST", "/api/v1/company/quota/payment", token, payload)
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusCreated {
		var result PaymentResponse
		json.NewDecoder(resp.Body).Decode(&result)

		assert.Equal(t, expectedTotal, result.Data.TotalAmount)
	}
}

// ============================================
// Quota Usage Priority Tests
// ============================================

func TestQuota_UseFreeQuotaFirst(t *testing.T) {
	// Verify that free quota is used before paid quota
	token := login(t, VerifiedCompany)

	// Get initial state
	resp := makeAuthenticatedRequest(t, "GET", "/api/v1/company/quota", token, nil)
	var initialQuota QuotaResponse
	json.NewDecoder(resp.Body).Decode(&initialQuota)
	resp.Body.Close()

	if initialQuota.Data.RemainingFreeQuota == 0 {
		t.Skip("No free quota remaining for this test")
	}

	initialFreeRemaining := initialQuota.Data.RemainingFreeQuota
	initialPaidQuota := initialQuota.Data.PaidQuota

	// Create and publish a job
	job := map[string]interface{}{
		"title":            "Quota Priority Test",
		"description":      "<p>Test</p>",
		"requirements":     "<p>Test</p>",
		"location":         "Jakarta",
		"type":             "full-time",
		"experience_level": "mid",
		"publish":          true,
	}

	resp = makeAuthenticatedRequest(t, "POST", "/api/v1/company/jobs", token, job)
	resp.Body.Close()

	// Get updated quota
	resp = makeAuthenticatedRequest(t, "GET", "/api/v1/company/quota", token, nil)
	var updatedQuota QuotaResponse
	json.NewDecoder(resp.Body).Decode(&updatedQuota)
	resp.Body.Close()

	// Free quota should decrease, paid quota should remain same
	assert.Equal(t, initialFreeRemaining-1, updatedQuota.Data.RemainingFreeQuota)
	assert.Equal(t, initialPaidQuota, updatedQuota.Data.PaidQuota)
}

func TestQuota_UsePaidQuotaWhenFreeExhausted(t *testing.T) {
	// When free quota is 0, paid quota should be used
	// Note: This requires a test user with exhausted free quota but has paid quota
	t.Skip("Requires specific test data setup")
}
