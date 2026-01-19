package admin

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/karirnusantara/api/internal/middleware"
	"github.com/karirnusantara/api/internal/shared/response"
)

// Handler handles admin HTTP requests
type Handler struct {
	service Service
}

// NewHandler creates a new admin handler
func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// ============================================
// AUTHENTICATION
// ============================================

// Login handles admin login
// POST /api/v1/admin/auth/login
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req AdminLoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_REQUEST", "Format request tidak valid")
		return
	}

	// Validate required fields
	if req.Email == "" || req.Password == "" {
		response.Error(w, http.StatusBadRequest, "VALIDATION_ERROR", "Email dan password wajib diisi")
		return
	}

	result, err := h.service.Login(r.Context(), &req)
	if err != nil {
		if errors.Is(err, ErrInvalidCredentials) {
			response.Error(w, http.StatusUnauthorized, "INVALID_CREDENTIALS", err.Error())
			return
		}
		if errors.Is(err, ErrAccountInactive) {
			response.Error(w, http.StatusForbidden, "ACCOUNT_INACTIVE", err.Error())
			return
		}
		response.Error(w, http.StatusInternalServerError, "LOGIN_FAILED", "Gagal melakukan login")
		return
	}

	response.Success(w, http.StatusOK, "Login berhasil", result)
}

// GetCurrentAdmin handles getting current admin info
// GET /api/v1/admin/auth/me
func (h *Handler) GetCurrentAdmin(w http.ResponseWriter, r *http.Request) {
	adminID := middleware.GetUserID(r.Context())
	if adminID == 0 {
		response.Error(w, http.StatusUnauthorized, "UNAUTHORIZED", "Tidak terautentikasi")
		return
	}

	admin, err := h.service.GetCurrentAdmin(r.Context(), adminID)
	if err != nil {
		if errors.Is(err, ErrAdminNotFound) {
			response.Error(w, http.StatusNotFound, "NOT_FOUND", err.Error())
			return
		}
		response.Error(w, http.StatusInternalServerError, "FETCH_FAILED", "Gagal mengambil data admin")
		return
	}

	response.Success(w, http.StatusOK, "Data admin berhasil diambil", admin)
}

// ============================================
// DASHBOARD
// ============================================

// GetDashboardStats handles getting dashboard statistics
// GET /api/v1/admin/dashboard/stats
func (h *Handler) GetDashboardStats(w http.ResponseWriter, r *http.Request) {
	stats, err := h.service.GetDashboardStats(r.Context())
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "FETCH_FAILED", "Gagal mengambil statistik")
		return
	}

	response.Success(w, http.StatusOK, "Statistik dashboard berhasil diambil", stats)
}

// ============================================
// COMPANY MANAGEMENT
// ============================================

// GetCompanies handles listing companies
// GET /api/v1/admin/companies
func (h *Handler) GetCompanies(w http.ResponseWriter, r *http.Request) {
	filter := CompanyFilter{
		Status:   r.URL.Query().Get("status"),
		Search:   r.URL.Query().Get("search"),
		Page:     parseIntOrDefault(r.URL.Query().Get("page"), 1),
		PageSize: parseIntOrDefault(r.URL.Query().Get("page_size"), 10),
	}

	result, err := h.service.GetCompanies(r.Context(), filter)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "FETCH_FAILED", "Gagal mengambil daftar perusahaan")
		return
	}

	response.SuccessWithMeta(w, http.StatusOK, "Daftar perusahaan berhasil diambil", result.Data, &response.Meta{
		Page:       result.Page,
		PerPage:    result.PageSize,
		TotalItems: int64(result.Total),
		TotalPages: result.TotalPages,
	})
}

// GetCompanyByID handles getting a single company
// GET /api/v1/admin/companies/{id}
func (h *Handler) GetCompanyByID(w http.ResponseWriter, r *http.Request) {
	id := parseIDFromRequest(r)
	if id == 0 {
		response.Error(w, http.StatusBadRequest, "INVALID_ID", "ID tidak valid")
		return
	}

	company, err := h.service.GetCompanyByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, ErrCompanyNotFound) {
			response.Error(w, http.StatusNotFound, "NOT_FOUND", err.Error())
			return
		}
		response.Error(w, http.StatusInternalServerError, "FETCH_FAILED", "Gagal mengambil data perusahaan")
		return
	}

	response.Success(w, http.StatusOK, "Data perusahaan berhasil diambil", company)
}

// VerifyCompany handles company verification
// POST /api/v1/admin/companies/{id}/verify
func (h *Handler) VerifyCompany(w http.ResponseWriter, r *http.Request) {
	id := parseIDFromRequest(r)
	if id == 0 {
		response.Error(w, http.StatusBadRequest, "INVALID_ID", "ID tidak valid")
		return
	}

	var req CompanyVerificationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_REQUEST", "Format request tidak valid")
		return
	}

	if req.Action == "" {
		response.Error(w, http.StatusBadRequest, "VALIDATION_ERROR", "Action wajib diisi (approve/reject)")
		return
	}

	adminID := middleware.GetUserID(r.Context())
	if err := h.service.VerifyCompany(r.Context(), id, &req, adminID); err != nil {
		if errors.Is(err, ErrCompanyNotFound) {
			response.Error(w, http.StatusNotFound, "NOT_FOUND", err.Error())
			return
		}
		if errors.Is(err, ErrInvalidAction) {
			response.Error(w, http.StatusBadRequest, "INVALID_ACTION", err.Error())
			return
		}
		response.Error(w, http.StatusInternalServerError, "UPDATE_FAILED", "Gagal memverifikasi perusahaan")
		return
	}

	response.Success(w, http.StatusOK, "Perusahaan berhasil diverifikasi", nil)
}

// UpdateCompanyStatus handles company status change (suspend/reactivate)
// PATCH /api/v1/admin/companies/{id}/status
func (h *Handler) UpdateCompanyStatus(w http.ResponseWriter, r *http.Request) {
	id := parseIDFromRequest(r)
	if id == 0 {
		response.Error(w, http.StatusBadRequest, "INVALID_ID", "ID tidak valid")
		return
	}

	var req CompanyStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_REQUEST", "Format request tidak valid")
		return
	}

	if req.Action == "" {
		response.Error(w, http.StatusBadRequest, "VALIDATION_ERROR", "Action wajib diisi (suspend/reactivate)")
		return
	}

	adminID := middleware.GetUserID(r.Context())
	if err := h.service.UpdateCompanyStatus(r.Context(), id, &req, adminID); err != nil {
		if errors.Is(err, ErrCompanyNotFound) {
			response.Error(w, http.StatusNotFound, "NOT_FOUND", err.Error())
			return
		}
		if errors.Is(err, ErrInvalidAction) {
			response.Error(w, http.StatusBadRequest, "INVALID_ACTION", err.Error())
			return
		}
		response.Error(w, http.StatusInternalServerError, "UPDATE_FAILED", "Gagal mengubah status perusahaan")
		return
	}

	response.Success(w, http.StatusOK, "Status perusahaan berhasil diubah", nil)
}

// ============================================
// JOB MANAGEMENT
// ============================================

// GetJobs handles listing jobs
// GET /api/v1/admin/jobs
func (h *Handler) GetJobs(w http.ResponseWriter, r *http.Request) {
	filter := JobFilter{
		CompanyID: parseUint64OrDefault(r.URL.Query().Get("company_id"), 0),
		Status:    r.URL.Query().Get("status"),
		DateFrom:  r.URL.Query().Get("date_from"),
		DateTo:    r.URL.Query().Get("date_to"),
		Search:    r.URL.Query().Get("search"),
		Page:      parseIntOrDefault(r.URL.Query().Get("page"), 1),
		PageSize:  parseIntOrDefault(r.URL.Query().Get("page_size"), 10),
	}

	result, err := h.service.GetJobs(r.Context(), filter)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "FETCH_FAILED", "Gagal mengambil daftar lowongan")
		return
	}

	response.SuccessWithMeta(w, http.StatusOK, "Daftar lowongan berhasil diambil", result.Data, &response.Meta{
		Page:       result.Page,
		PerPage:    result.PageSize,
		TotalItems: int64(result.Total),
		TotalPages: result.TotalPages,
	})
}

// GetJobByID handles getting a single job
// GET /api/v1/admin/jobs/{id}
func (h *Handler) GetJobByID(w http.ResponseWriter, r *http.Request) {
	id := parseIDFromRequest(r)
	if id == 0 {
		response.Error(w, http.StatusBadRequest, "INVALID_ID", "ID tidak valid")
		return
	}

	job, err := h.service.GetJobByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, ErrJobNotFound) {
			response.Error(w, http.StatusNotFound, "NOT_FOUND", err.Error())
			return
		}
		response.Error(w, http.StatusInternalServerError, "FETCH_FAILED", "Gagal mengambil data lowongan")
		return
	}

	response.Success(w, http.StatusOK, "Data lowongan berhasil diambil", job)
}

// ModerateJob handles job moderation
// POST /api/v1/admin/jobs/{id}/moderate
func (h *Handler) ModerateJob(w http.ResponseWriter, r *http.Request) {
	id := parseIDFromRequest(r)
	if id == 0 {
		response.Error(w, http.StatusBadRequest, "INVALID_ID", "ID tidak valid")
		return
	}

	var req JobActionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_REQUEST", "Format request tidak valid")
		return
	}

	if req.Action == "" {
		response.Error(w, http.StatusBadRequest, "VALIDATION_ERROR", "Action wajib diisi")
		return
	}

	adminID := middleware.GetUserID(r.Context())
	if err := h.service.ModerateJob(r.Context(), id, &req, adminID); err != nil {
		if errors.Is(err, ErrJobNotFound) {
			response.Error(w, http.StatusNotFound, "NOT_FOUND", err.Error())
			return
		}
		if errors.Is(err, ErrInvalidAction) {
			response.Error(w, http.StatusBadRequest, "INVALID_ACTION", err.Error())
			return
		}
		response.Error(w, http.StatusInternalServerError, "UPDATE_FAILED", "Gagal memoderasi lowongan")
		return
	}

	response.Success(w, http.StatusOK, "Lowongan berhasil dimoderasi", nil)
}

// ============================================
// PAYMENT MANAGEMENT
// ============================================

// GetPayments handles listing payments
// GET /api/v1/admin/payments
func (h *Handler) GetPayments(w http.ResponseWriter, r *http.Request) {
	filter := PaymentFilter{
		CompanyID: parseUint64OrDefault(r.URL.Query().Get("company_id"), 0),
		Status:    r.URL.Query().Get("status"),
		DateFrom:  r.URL.Query().Get("date_from"),
		DateTo:    r.URL.Query().Get("date_to"),
		Page:      parseIntOrDefault(r.URL.Query().Get("page"), 1),
		PageSize:  parseIntOrDefault(r.URL.Query().Get("page_size"), 10),
	}

	result, err := h.service.GetPayments(r.Context(), filter)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "FETCH_FAILED", "Gagal mengambil daftar pembayaran")
		return
	}

	response.SuccessWithMeta(w, http.StatusOK, "Daftar pembayaran berhasil diambil", result.Data, &response.Meta{
		Page:       result.Page,
		PerPage:    result.PageSize,
		TotalItems: int64(result.Total),
		TotalPages: result.TotalPages,
	})
}

// GetPaymentByID handles getting a single payment
// GET /api/v1/admin/payments/{id}
func (h *Handler) GetPaymentByID(w http.ResponseWriter, r *http.Request) {
	id := parseIDFromRequest(r)
	if id == 0 {
		response.Error(w, http.StatusBadRequest, "INVALID_ID", "ID tidak valid")
		return
	}

	payment, err := h.service.GetPaymentByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, ErrPaymentNotFound) {
			response.Error(w, http.StatusNotFound, "NOT_FOUND", err.Error())
			return
		}
		response.Error(w, http.StatusInternalServerError, "FETCH_FAILED", "Gagal mengambil data pembayaran")
		return
	}

	response.Success(w, http.StatusOK, "Data pembayaran berhasil diambil", payment)
}

// ProcessPayment handles payment approval/rejection
// POST /api/v1/admin/payments/{id}/process
func (h *Handler) ProcessPayment(w http.ResponseWriter, r *http.Request) {
	id := parseIDFromRequest(r)
	if id == 0 {
		response.Error(w, http.StatusBadRequest, "INVALID_ID", "ID tidak valid")
		return
	}

	var req PaymentActionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_REQUEST", "Format request tidak valid")
		return
	}

	if req.Action == "" {
		response.Error(w, http.StatusBadRequest, "VALIDATION_ERROR", "Action wajib diisi (approve/reject)")
		return
	}

	adminID := middleware.GetUserID(r.Context())
	if err := h.service.ProcessPayment(r.Context(), id, &req, adminID); err != nil {
		if errors.Is(err, ErrPaymentNotFound) {
			response.Error(w, http.StatusNotFound, "NOT_FOUND", err.Error())
			return
		}
		if errors.Is(err, ErrInvalidAction) {
			response.Error(w, http.StatusBadRequest, "INVALID_ACTION", err.Error())
			return
		}
		response.Error(w, http.StatusInternalServerError, "UPDATE_FAILED", "Gagal memproses pembayaran")
		return
	}

	response.Success(w, http.StatusOK, "Pembayaran berhasil diproses", nil)
}

// ============================================
// JOB SEEKER MANAGEMENT
// ============================================

// GetJobSeekers handles listing job seekers
// GET /api/v1/admin/job-seekers
func (h *Handler) GetJobSeekers(w http.ResponseWriter, r *http.Request) {
	filter := JobSeekerFilter{
		Status:   r.URL.Query().Get("status"),
		Search:   r.URL.Query().Get("search"),
		Page:     parseIntOrDefault(r.URL.Query().Get("page"), 1),
		PageSize: parseIntOrDefault(r.URL.Query().Get("page_size"), 10),
	}

	result, err := h.service.GetJobSeekers(r.Context(), filter)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "FETCH_FAILED", "Gagal mengambil daftar pencari kerja")
		return
	}

	response.SuccessWithMeta(w, http.StatusOK, "Daftar pencari kerja berhasil diambil", result.Data, &response.Meta{
		Page:       result.Page,
		PerPage:    result.PageSize,
		TotalItems: int64(result.Total),
		TotalPages: result.TotalPages,
	})
}

// GetJobSeekerByID handles getting a single job seeker
// GET /api/v1/admin/job-seekers/{id}
func (h *Handler) GetJobSeekerByID(w http.ResponseWriter, r *http.Request) {
	id := parseIDFromRequest(r)
	if id == 0 {
		response.Error(w, http.StatusBadRequest, "INVALID_ID", "ID tidak valid")
		return
	}

	jobSeeker, err := h.service.GetJobSeekerByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, ErrJobSeekerNotFound) {
			response.Error(w, http.StatusNotFound, "NOT_FOUND", err.Error())
			return
		}
		response.Error(w, http.StatusInternalServerError, "FETCH_FAILED", "Gagal mengambil data pencari kerja")
		return
	}

	response.Success(w, http.StatusOK, "Data pencari kerja berhasil diambil", jobSeeker)
}

// UpdateJobSeekerStatus handles job seeker status change
// PATCH /api/v1/admin/job-seekers/{id}/status
func (h *Handler) UpdateJobSeekerStatus(w http.ResponseWriter, r *http.Request) {
	id := parseIDFromRequest(r)
	if id == 0 {
		response.Error(w, http.StatusBadRequest, "INVALID_ID", "ID tidak valid")
		return
	}

	var req JobSeekerActionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_REQUEST", "Format request tidak valid")
		return
	}

	if req.Action == "" {
		response.Error(w, http.StatusBadRequest, "VALIDATION_ERROR", "Action wajib diisi")
		return
	}

	adminID := middleware.GetUserID(r.Context())
	if err := h.service.UpdateJobSeekerStatus(r.Context(), id, &req, adminID); err != nil {
		if errors.Is(err, ErrJobSeekerNotFound) {
			response.Error(w, http.StatusNotFound, "NOT_FOUND", err.Error())
			return
		}
		if errors.Is(err, ErrInvalidAction) {
			response.Error(w, http.StatusBadRequest, "INVALID_ACTION", err.Error())
			return
		}
		response.Error(w, http.StatusInternalServerError, "UPDATE_FAILED", "Gagal mengubah status pencari kerja")
		return
	}

	response.Success(w, http.StatusOK, "Status pencari kerja berhasil diubah", nil)
}

// ============================================
// HELPERS
// ============================================

func parseIntOrDefault(s string, defaultVal int) int {
	if s == "" {
		return defaultVal
	}
	val, err := strconv.Atoi(s)
	if err != nil {
		return defaultVal
	}
	return val
}

func parseUint64OrDefault(s string, defaultVal uint64) uint64 {
	if s == "" {
		return defaultVal
	}
	val, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return defaultVal
	}
	return val
}

func parseIDFromRequest(r *http.Request) uint64 {
	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		return 0
	}
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return 0
	}
	return id
}
