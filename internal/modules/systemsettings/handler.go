package systemsettings

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/karirnusantara/api/internal/middleware"
	"github.com/karirnusantara/api/internal/shared/response"
	"github.com/karirnusantara/api/internal/shared/validator"
)

// Handler handles HTTP requests for system settings
type Handler struct {
	service   *Service
	validator *validator.Validator
}

// NewHandler creates a new Handler
func NewHandler(service *Service, v *validator.Validator) *Handler {
	return &Handler{service: service, validator: v}
}

// -----------------------------------------------------------
// GET /api/v1/admin/system-settings
// Returns full settings for the admin dashboard
// -----------------------------------------------------------
func (h *Handler) GetSettings(w http.ResponseWriter, r *http.Request) {
	result, err := h.service.GetSettings()
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "FETCH_FAILED", "Gagal mengambil pengaturan sistem")
		return
	}
	response.Success(w, http.StatusOK, "Pengaturan sistem berhasil diambil", result)
}

// -----------------------------------------------------------
// PUT /api/v1/admin/system-settings
// Updates general settings (free quota limit, price per job, currency)
// -----------------------------------------------------------
func (h *Handler) UpdateSettings(w http.ResponseWriter, r *http.Request) {
	adminID := middleware.GetUserID(r.Context())
	if adminID == 0 {
		response.Error(w, http.StatusUnauthorized, "UNAUTHORIZED", "Tidak terautentikasi")
		return
	}

	var req UpdateSettingsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_REQUEST", "Format request tidak valid")
		return
	}

	if err := h.service.UpdateSettings(&req, &adminID); err != nil {
		response.Error(w, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}

	// Return the updated settings
	result, err := h.service.GetSettings()
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "FETCH_FAILED", "Gagal mengambil pengaturan setelah update")
		return
	}

	response.Success(w, http.StatusOK, "Pengaturan sistem berhasil diperbarui", result)
}

// -----------------------------------------------------------
// POST /api/v1/admin/system-settings/packages
// Create or update a quota package
// -----------------------------------------------------------
func (h *Handler) UpsertPackage(w http.ResponseWriter, r *http.Request) {
	adminID := middleware.GetUserID(r.Context())
	if adminID == 0 {
		response.Error(w, http.StatusUnauthorized, "UNAUTHORIZED", "Tidak terautentikasi")
		return
	}

	var req UpsertQuotaPackageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_REQUEST", "Format request tidak valid")
		return
	}

	if h.validator != nil {
		if errs := h.validator.Validate(req); errs != nil {
			response.UnprocessableEntity(w, "Validasi gagal", errs)
			return
		}
	}

	pkg, err := h.service.UpsertPackage(&req, &adminID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "UPSERT_FAILED", "Gagal menyimpan paket kuota")
		return
	}

	response.Success(w, http.StatusOK, "Paket kuota berhasil disimpan", pkg)
}

// -----------------------------------------------------------
// DELETE /api/v1/admin/system-settings/packages/{id}
// Delete a quota package
// -----------------------------------------------------------
func (h *Handler) DeletePackage(w http.ResponseWriter, r *http.Request) {
	adminID := middleware.GetUserID(r.Context())
	if adminID == 0 {
		response.Error(w, http.StatusUnauthorized, "UNAUTHORIZED", "Tidak terautentikasi")
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_ID", "ID tidak valid")
		return
	}

	if err := h.service.DeletePackage(id); err != nil {
		if err.Error() == "package not found" {
			response.Error(w, http.StatusNotFound, "NOT_FOUND", "Paket kuota tidak ditemukan")
			return
		}
		response.Error(w, http.StatusInternalServerError, "DELETE_FAILED", "Gagal menghapus paket kuota")
		return
	}

	response.Success(w, http.StatusOK, "Paket kuota berhasil dihapus", nil)
}

// -----------------------------------------------------------
// GET /api/v1/pricing   (public — no auth required)
// Returns public pricing info consumed by all frontends
// -----------------------------------------------------------
func (h *Handler) GetPublicPricing(w http.ResponseWriter, r *http.Request) {
	result, err := h.service.GetPublicPricing()
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "FETCH_FAILED", "Gagal mengambil informasi harga")
		return
	}
	response.Success(w, http.StatusOK, "Informasi harga berhasil diambil", result)
}
