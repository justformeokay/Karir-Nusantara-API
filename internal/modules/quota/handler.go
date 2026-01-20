package quota

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/karirnusantara/api/internal/middleware"
	"github.com/karirnusantara/api/internal/shared/response"
	"github.com/karirnusantara/api/internal/shared/validator"
)

// Handler handles HTTP requests for quota
type Handler struct {
	service   *Service
	validator *validator.Validator
}

// NewHandler creates a new quota handler
func NewHandler(service *Service, v *validator.Validator) *Handler {
	return &Handler{
		service:   service,
		validator: v,
	}
}

// GetQuota returns the company's quota information
// @Summary Get company quota
// @Description Get job posting quota for the authenticated company
// @Tags Quota
// @Security BearerAuth
// @Produce json
// @Success 200 {object} response.Response{data=QuotaResponse}
// @Router /company/quota [get]
func (h *Handler) GetQuota(w http.ResponseWriter, r *http.Request) {
	companyID := middleware.GetUserID(r.Context())
	if companyID == 0 {
		response.Error(w, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
		return
	}

	quota, err := h.service.GetQuota(companyID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "QUOTA_ERROR", "Failed to get quota")
		return
	}

	response.Success(w, http.StatusOK, "Quota retrieved successfully", quota)
}

// SubmitPaymentProof handles payment proof submission
// @Summary Submit payment proof
// @Description Submit payment proof for additional job quota
// @Tags Quota
// @Security BearerAuth
// @Accept multipart/form-data
// @Produce json
// @Param proof_image formData file true "Payment proof image"
// @Param job_id formData int false "Job ID (optional)"
// @Success 201 {object} response.Response{data=PaymentResponse}
// @Router /company/payments/proof [post]
func (h *Handler) SubmitPaymentProof(w http.ResponseWriter, r *http.Request) {
	companyID := middleware.GetUserID(r.Context())
	if companyID == 0 {
		response.Error(w, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
		return
	}

	// Parse multipart form (max 10MB)
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		response.Error(w, http.StatusBadRequest, "PARSE_ERROR", "Failed to parse form data")
		return
	}

	// Handle file upload
	file, header, err := r.FormFile("proof_image")
	if err != nil {
		response.Error(w, http.StatusBadRequest, "FILE_REQUIRED", "Payment proof image is required")
		return
	}
	defer file.Close()

	// Create upload directory if not exists
	uploadDir := fmt.Sprintf("./docs/payments/%d", companyID)
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		response.Error(w, http.StatusInternalServerError, "UPLOAD_ERROR", "Failed to create upload directory")
		return
	}

	// Generate unique filename
	ext := filepath.Ext(header.Filename)
	filename := fmt.Sprintf("proof_%d%s", time.Now().Unix(), ext)
	filePath := filepath.Join(uploadDir, filename)

	// Create destination file
	dst, err := os.Create(filePath)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "UPLOAD_ERROR", "Failed to save file")
		return
	}
	defer dst.Close()

	// Copy file content
	if _, err := io.Copy(dst, file); err != nil {
		response.Error(w, http.StatusInternalServerError, "UPLOAD_ERROR", "Failed to save file")
		return
	}

	// URL path for the uploaded file
	proofImageURL := fmt.Sprintf("/docs/payments/%d/%s", companyID, filename)

	// Parse optional job_id
	var jobID *uint64
	if jobIDStr := r.FormValue("job_id"); jobIDStr != "" {
		if parsed, err := strconv.ParseUint(jobIDStr, 10, 64); err == nil {
			jobID = &parsed
		}
	}

	payment, err := h.service.SubmitPaymentProof(companyID, jobID, proofImageURL)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "PAYMENT_ERROR", "Failed to submit payment proof")
		return
	}

	response.Success(w, http.StatusCreated, "Payment proof submitted successfully", payment.ToResponse())
}

// GetPayments returns the payment history
// @Summary Get payment history
// @Description Get payment history for the authenticated company
// @Tags Quota
// @Security BearerAuth
// @Produce json
// @Param page query int false "Page number"
// @Param per_page query int false "Items per page"
// @Param status query string false "Filter by status"
// @Success 200 {object} response.Response{data=[]PaymentResponse}
// @Router /company/payments [get]
func (h *Handler) GetPayments(w http.ResponseWriter, r *http.Request) {
	companyID := middleware.GetUserID(r.Context())
	if companyID == 0 {
		response.Error(w, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
		return
	}

	params := DefaultPaymentListParams()
	params.CompanyID = companyID

	if page, err := strconv.Atoi(r.URL.Query().Get("page")); err == nil && page > 0 {
		params.Page = page
	}
	if perPage, err := strconv.Atoi(r.URL.Query().Get("per_page")); err == nil && perPage > 0 {
		params.PerPage = perPage
	}
	if status := r.URL.Query().Get("status"); status != "" {
		params.Status = status
	}

	payments, total, err := h.service.GetPayments(params)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "PAYMENTS_ERROR", "Failed to get payments")
		return
	}

	totalPages := (total + params.PerPage - 1) / params.PerPage

	response.SuccessWithMeta(w, http.StatusOK, "Payments retrieved successfully", payments, &response.Meta{
		Page:       params.Page,
		PerPage:    params.PerPage,
		TotalPages: totalPages,
		TotalItems: int64(total),
	})
}

// GetPaymentInfo returns the payment information (bank account, etc.)
// @Summary Get payment info
// @Description Get bank account information for payment
// @Tags Quota
// @Produce json
// @Success 200 {object} response.Response
// @Router /company/payments/info [get]
func (h *Handler) GetPaymentInfo(w http.ResponseWriter, r *http.Request) {
	info := map[string]interface{}{
		"bank":           "BCA",
		"account_number": "8725164421",
		"account_name":   "Saputra Budianto",
		"price_per_job":  PricePerJob,
	}

	response.Success(w, http.StatusOK, "Payment info retrieved", info)
}

// parseJSON helper to parse JSON body
func parseJSON(r *http.Request, v interface{}) error {
	return json.NewDecoder(r.Body).Decode(v)
}
