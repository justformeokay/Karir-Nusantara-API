package quota

import (
	"context"
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

// CompanyService defines the interface for company operations needed by quota
type CompanyService interface {
	GetCompanyIDByUserID(ctx context.Context, userID uint64) (uint64, error)
}

// Handler handles HTTP requests for quota
type Handler struct {
	service        *Service
	validator      *validator.Validator
	companyService CompanyService
}

// NewHandler creates a new quota handler
func NewHandler(service *Service, v *validator.Validator, companyService CompanyService) *Handler {
	return &Handler{
		service:        service,
		validator:      v,
		companyService: companyService,
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
	userID := middleware.GetUserID(r.Context())
	if userID == 0 {
		response.Error(w, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
		return
	}

	// Get actual company ID from user ID
	companyID, err := h.companyService.GetCompanyIDByUserID(r.Context(), userID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "COMPANY_ERROR", "Failed to get company")
		return
	}
	if companyID == 0 {
		response.Error(w, http.StatusNotFound, "COMPANY_NOT_FOUND", "Company not found")
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
	userID := middleware.GetUserID(r.Context())
	if userID == 0 {
		response.Error(w, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
		return
	}

	// Get actual company ID from user ID
	companyID, err := h.companyService.GetCompanyIDByUserID(r.Context(), userID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "COMPANY_ERROR", "Failed to get company")
		return
	}
	if companyID == 0 {
		response.Error(w, http.StatusNotFound, "COMPANY_NOT_FOUND", "Company not found")
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

	// Parse optional package_id (for top-up packages)
	var packageID *string
	if pkgID := r.FormValue("package_id"); pkgID != "" {
		// Validate package exists
		if GetPackageByID(pkgID) == nil {
			response.Error(w, http.StatusBadRequest, "INVALID_PACKAGE", "Invalid package ID")
			return
		}
		packageID = &pkgID
	}

	payment, err := h.service.SubmitPaymentProof(companyID, jobID, packageID, proofImageURL)
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
	userID := middleware.GetUserID(r.Context())
	if userID == 0 {
		response.Error(w, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
		return
	}

	// Get actual company ID from user ID
	companyID, err := h.companyService.GetCompanyIDByUserID(r.Context(), userID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "COMPANY_ERROR", "Failed to get company")
		return
	}
	if companyID == 0 {
		response.Error(w, http.StatusNotFound, "COMPANY_NOT_FOUND", "Company not found")
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
		"packages":       GetTopUpPackages(),
	}

	response.Success(w, http.StatusOK, "Payment info retrieved", info)
}

// GetPackages returns available top-up packages
// @Summary Get top-up packages
// @Description Get all available top-up packages for job posting quota
// @Tags Quota
// @Produce json
// @Success 200 {object} response.Response{data=[]TopUpPackage}
// @Router /company/packages [get]
func (h *Handler) GetPackages(w http.ResponseWriter, r *http.Request) {
	packages := GetTopUpPackages()
	response.Success(w, http.StatusOK, "Packages retrieved successfully", packages)
}

// DownloadInvoice downloads the invoice PDF for a confirmed payment
// @Summary Download payment invoice
// @Description Download invoice PDF for a confirmed payment
// @Tags Quota
// @Security BearerAuth
// @Produce application/pdf
// @Param id path int true "Payment ID"
// @Success 200 {file} binary
// @Router /company/payments/{id}/invoice [get]
func (h *Handler) DownloadInvoice(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == 0 {
		response.Error(w, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
		return
	}

	// Get actual company ID from user ID
	companyID, err := h.companyService.GetCompanyIDByUserID(r.Context(), userID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "COMPANY_ERROR", "Failed to get company")
		return
	}
	if companyID == 0 {
		response.Error(w, http.StatusNotFound, "COMPANY_NOT_FOUND", "Company not found")
		return
	}

	// Get payment ID from URL
	paymentIDStr := r.URL.Query().Get("id")
	if paymentIDStr == "" {
		response.Error(w, http.StatusBadRequest, "INVALID_ID", "Payment ID is required")
		return
	}

	paymentID, err := strconv.ParseUint(paymentIDStr, 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_ID", "Invalid payment ID")
		return
	}

	// Verify payment belongs to company and is confirmed
	payment, err := h.service.GetPaymentByID(paymentID)
	if err != nil {
		response.Error(w, http.StatusNotFound, "PAYMENT_NOT_FOUND", "Payment not found")
		return
	}

	if payment.CompanyID != companyID {
		response.Error(w, http.StatusForbidden, "FORBIDDEN", "You don't have access to this payment")
		return
	}

	if payment.Status != "confirmed" {
		response.Error(w, http.StatusBadRequest, "PAYMENT_NOT_CONFIRMED", "Invoice only available for confirmed payments")
		return
	}

	// Build invoice file path - use confirmed_at date or current date
	var dateStr string
	if payment.ConfirmedAt.Valid {
		dateStr = payment.ConfirmedAt.Time.Format("20060102")
	} else {
		dateStr = time.Now().Format("20060102")
	}

	invoicePath := fmt.Sprintf("./docs/invoices/invoice_%s_%d.pdf", dateStr, paymentID)

	// Check if invoice file exists
	if _, err := os.Stat(invoicePath); os.IsNotExist(err) {
		response.Error(w, http.StatusNotFound, "INVOICE_NOT_FOUND", "Invoice file not found")
		return
	}

	// Read invoice file
	invoiceData, err := os.ReadFile(invoicePath)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "READ_ERROR", "Failed to read invoice file")
		return
	}

	// Set headers for PDF download
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=invoice_%d.pdf", paymentID))
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(invoiceData)))

	// Write PDF data
	w.Write(invoiceData)
}

// parseJSON helper to parse JSON body
func parseJSON(r *http.Request, v interface{}) error {
	return json.NewDecoder(r.Body).Decode(v)
}
