package admin

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"

	"github.com/karirnusantara/api/internal/shared/hashid"
)

// Repository defines the admin repository interface
type Repository interface {
	// Admin user operations
	GetAdminByEmail(ctx context.Context, email string) (*AdminUser, error)
	GetAdminByID(ctx context.Context, id uint64) (*AdminUser, error)
	GetUserByID(ctx context.Context, id uint64) (*SimpleUser, error)

	// Dashboard
	GetDashboardStats(ctx context.Context) (*DashboardStats, error)

	// Company operations
	GetCompanies(ctx context.Context, filter CompanyFilter) ([]*CompanyAdmin, int, error)
	GetCompanyByID(ctx context.Context, id uint64) (*CompanyAdmin, error)
	GetCompanyDetailByID(ctx context.Context, id uint64) (*CompanyDetailResponse, error)
	UpdateCompanyStatus(ctx context.Context, id uint64, status string) error
	UpdateCompanyActive(ctx context.Context, id uint64, isActive bool) error

	// Job operations
	GetJobs(ctx context.Context, filter JobFilter) ([]*JobAdmin, int, error)
	GetJobByID(ctx context.Context, id uint64) (*JobAdmin, error)
	UpdateJobStatus(ctx context.Context, id uint64, status string) error
	UpdateJobAdminStatus(ctx context.Context, id uint64, adminStatus, note string) error

	// Payment operations
	GetPayments(ctx context.Context, filter PaymentFilter) ([]*PaymentAdmin, int, error)
	GetPaymentByID(ctx context.Context, id uint64) (*PaymentAdmin, error)
	UpdatePaymentStatus(ctx context.Context, id uint64, status, note string, confirmedByID uint64) error

	// Job seeker operations
	GetJobSeekers(ctx context.Context, filter JobSeekerFilter) ([]*JobSeekerAdmin, int, error)
	GetJobSeekerByID(ctx context.Context, id uint64) (*JobSeekerAdmin, error)
	UpdateJobSeekerActive(ctx context.Context, id uint64, isActive bool) error

	// Audit log
	LogAdminAction(ctx context.Context, log *AdminActionLog) error
}

// repository implements the Repository interface
type repository struct {
	db *sqlx.DB
}

// NewRepository creates a new admin repository
func NewRepository(db *sqlx.DB) Repository {
	return &repository{db: db}
}

// ============================================
// ADMIN USER OPERATIONS
// ============================================

func (r *repository) GetAdminByEmail(ctx context.Context, email string) (*AdminUser, error) {
	query := `
		SELECT id, email, password_hash, full_name, role, is_active, created_at, updated_at
		FROM users
		WHERE email = ? AND role = 'admin' AND is_active = true
	`

	var admin AdminUser
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&admin.ID,
		&admin.Email,
		&admin.PasswordHash,
		&admin.FullName,
		&admin.Role,
		&admin.IsActive,
		&admin.CreatedAt,
		&admin.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &admin, nil
}

func (r *repository) GetAdminByID(ctx context.Context, id uint64) (*AdminUser, error) {
	query := `
		SELECT id, email, password_hash, full_name, role, is_active, created_at, updated_at
		FROM users
		WHERE id = ? AND role = 'admin'
	`

	var admin AdminUser
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&admin.ID,
		&admin.Email,
		&admin.PasswordHash,
		&admin.FullName,
		&admin.Role,
		&admin.IsActive,
		&admin.CreatedAt,
		&admin.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &admin, nil
}

func (r *repository) GetUserByID(ctx context.Context, id uint64) (*SimpleUser, error) {
	query := `
		SELECT id, email, full_name, role
		FROM users
		WHERE id = ?
	`

	var user SimpleUser
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Email,
		&user.FullName,
		&user.Role,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

// ============================================
// DASHBOARD OPERATIONS
// ============================================

func (r *repository) GetDashboardStats(ctx context.Context) (*DashboardStats, error) {
	stats := &DashboardStats{}

	// Companies stats
	err := r.db.QueryRowContext(ctx, `
		SELECT 
			COUNT(*) as total,
			SUM(CASE WHEN company_status = 'pending' THEN 1 ELSE 0 END) as pending,
			SUM(CASE WHEN company_status = 'verified' THEN 1 ELSE 0 END) as verified,
			SUM(CASE WHEN company_status = 'suspended' THEN 1 ELSE 0 END) as suspended
		FROM users WHERE role = 'company'
	`).Scan(&stats.TotalCompanies, &stats.PendingVerifications, &stats.VerifiedCompanies, &stats.SuspendedCompanies)
	if err != nil {
		return nil, fmt.Errorf("failed to get company stats: %w", err)
	}

	// Jobs stats
	err = r.db.QueryRowContext(ctx, `
		SELECT 
			COUNT(*) as total,
			SUM(CASE WHEN status = 'active' THEN 1 ELSE 0 END) as active,
			SUM(CASE WHEN status = 'pending' THEN 1 ELSE 0 END) as pending,
			SUM(CASE WHEN admin_status = 'flagged' THEN 1 ELSE 0 END) as flagged
		FROM jobs
	`).Scan(&stats.TotalJobs, &stats.ActiveJobs, &stats.PendingJobs, &stats.FlaggedJobs)
	if err != nil {
		return nil, fmt.Errorf("failed to get job stats: %w", err)
	}

	// Job seekers stats
	err = r.db.QueryRowContext(ctx, `
		SELECT 
			COUNT(*) as total,
			SUM(CASE WHEN is_active = true THEN 1 ELSE 0 END) as active
		FROM users WHERE role = 'job_seeker'
	`).Scan(&stats.TotalJobSeekers, &stats.ActiveJobSeekers)
	if err != nil {
		return nil, fmt.Errorf("failed to get job seeker stats: %w", err)
	}

	// Payments stats
	err = r.db.QueryRowContext(ctx, `
		SELECT 
			COUNT(*) as total,
			SUM(CASE WHEN status = 'pending' THEN 1 ELSE 0 END) as pending,
			SUM(CASE WHEN status = 'confirmed' THEN amount ELSE 0 END) as revenue
		FROM payments
	`).Scan(&stats.TotalPayments, &stats.PendingPayments, &stats.TotalRevenue)
	if err != nil {
		return nil, fmt.Errorf("failed to get payment stats: %w", err)
	}

	return stats, nil
}

// ============================================
// COMPANY OPERATIONS
// ============================================

func (r *repository) GetCompanies(ctx context.Context, filter CompanyFilter) ([]*CompanyAdmin, int, error) {
	// Count query - join users with companies table
	countQuery := `
		SELECT COUNT(*) FROM users u
		JOIN companies c ON u.id = c.user_id
		WHERE u.role = 'company'
	`
	args := []interface{}{}

	if filter.Status != "" {
		countQuery += " AND c.company_status = ?"
		args = append(args, filter.Status)
	}
	if filter.Search != "" {
		countQuery += " AND (c.company_name LIKE ? OR u.email LIKE ? OR u.full_name LIKE ?)"
		searchTerm := "%" + filter.Search + "%"
		args = append(args, searchTerm, searchTerm, searchTerm)
	}

	var total int
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	// Main query - join users with companies
	query := `
		SELECT 
			u.id, u.email, u.full_name, u.phone, c.company_name, c.company_description,
			c.company_website, c.company_logo_url, c.company_status, u.is_active, u.is_verified,
			u.email_verified_at, c.created_at, c.updated_at,
			(SELECT COUNT(*) FROM jobs WHERE company_id = u.id) as jobs_count,
			(SELECT COUNT(*) FROM jobs WHERE company_id = u.id AND status = 'active') as active_jobs,
			(SELECT COUNT(*) FROM applications a JOIN jobs j ON a.job_id = j.id WHERE j.company_id = u.id) as total_applications
		FROM users u
		JOIN companies c ON u.id = c.user_id
		WHERE u.role = 'company'
	`

	queryArgs := []interface{}{}
	if filter.Status != "" {
		query += " AND c.company_status = ?"
		queryArgs = append(queryArgs, filter.Status)
	}
	if filter.Search != "" {
		query += " AND (c.company_name LIKE ? OR u.email LIKE ? OR u.full_name LIKE ?)"
		searchTerm := "%" + filter.Search + "%"
		queryArgs = append(queryArgs, searchTerm, searchTerm, searchTerm)
	}

	query += " ORDER BY c.created_at DESC"

	// Pagination
	if filter.PageSize <= 0 {
		filter.PageSize = 10
	}
	if filter.Page <= 0 {
		filter.Page = 1
	}
	offset := (filter.Page - 1) * filter.PageSize
	query += fmt.Sprintf(" LIMIT %d OFFSET %d", filter.PageSize, offset)

	rows, err := r.db.QueryContext(ctx, query, queryArgs...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var companies []*CompanyAdmin
	for rows.Next() {
		var c CompanyAdmin
		if err := rows.Scan(
			&c.ID, &c.Email, &c.FullName, &c.Phone, &c.CompanyName, &c.CompanyDescription,
			&c.CompanyWebsite, &c.CompanyLogoURL, &c.CompanyStatus, &c.IsActive, &c.IsVerified,
			&c.EmailVerifiedAt, &c.CreatedAt, &c.UpdatedAt,
			&c.JobsCount, &c.ActiveJobsCount, &c.TotalApplications,
		); err != nil {
			return nil, 0, err
		}
		companies = append(companies, &c)
	}

	return companies, total, nil
}

func (r *repository) GetCompanyByID(ctx context.Context, id uint64) (*CompanyAdmin, error) {
	query := `
		SELECT 
			u.id, u.email, u.full_name, u.phone, c.company_name, c.company_description,
			c.company_website, c.company_logo_url, c.company_status, u.is_active, u.is_verified,
			u.email_verified_at, c.created_at, c.updated_at,
			(SELECT COUNT(*) FROM jobs WHERE company_id = u.id) as jobs_count,
			(SELECT COUNT(*) FROM jobs WHERE company_id = u.id AND status = 'active') as active_jobs,
			(SELECT COUNT(*) FROM applications a JOIN jobs j ON a.job_id = j.id WHERE j.company_id = u.id) as total_applications
		FROM users u
		JOIN companies c ON u.id = c.user_id
		WHERE u.id = ? AND u.role = 'company'
	`

	var c CompanyAdmin
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&c.ID, &c.Email, &c.FullName, &c.Phone, &c.CompanyName, &c.CompanyDescription,
		&c.CompanyWebsite, &c.CompanyLogoURL, &c.CompanyStatus, &c.IsActive, &c.IsVerified,
		&c.EmailVerifiedAt, &c.CreatedAt, &c.UpdatedAt,
		&c.JobsCount, &c.ActiveJobsCount, &c.TotalApplications,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &c, nil
}

// GetCompanyDetailByID fetches detailed company information including documents and quota
func (r *repository) GetCompanyDetailByID(ctx context.Context, id uint64) (*CompanyDetailResponse, error) {
	query := `
		SELECT 
			u.id, u.email, u.full_name, COALESCE(u.phone, '') as phone,
			c.company_name, COALESCE(c.company_description, '') as company_description, COALESCE(c.company_website, '') as company_website, COALESCE(c.company_logo_url, '') as company_logo_url,
			COALESCE(c.company_industry, '') as company_industry, COALESCE(c.company_size, '') as company_size, COALESCE(c.company_location, '') as company_location, COALESCE(c.company_address, '') as company_address,
			COALESCE(c.company_city, '') as company_city, COALESCE(c.company_province, '') as company_province, COALESCE(c.company_postal_code, '') as company_postal_code,
			COALESCE(c.established_year, 0) as established_year, COALESCE(c.employee_count, 0) as employee_count,
			c.company_status, u.is_active, u.is_verified, COALESCE(u.email_verified_at, '') as email_verified_at,
			COALESCE(c.documents_verified_at, '') as documents_verified_at, COALESCE(c.verification_notes, '') as verification_notes,
			COALESCE(c.ktp_founder_url, '') as ktp_founder_url, COALESCE(c.akta_pendirian_url, '') as akta_pendirian_url, COALESCE(c.npwp_url, '') as npwp_url, COALESCE(c.nib_url, '') as nib_url,
			c.created_at, c.updated_at,
			COALESCE((SELECT COUNT(*) FROM jobs WHERE company_id = c.id), 0) as jobs_count,
			COALESCE((SELECT COUNT(*) FROM jobs WHERE company_id = c.id AND status = 'active'), 0) as active_jobs_count,
			COALESCE((SELECT COUNT(*) FROM applications a JOIN jobs j ON a.job_id = j.id WHERE j.company_id = c.id), 0) as total_applications,
			COALESCE(cq.free_quota_used, 0) as free_quota_used,
			COALESCE(cq.paid_quota, 0) as paid_quota,
			COALESCE((SELECT COUNT(*) FROM jobs WHERE company_id = c.id AND status = 'draft'), 0) as draft_jobs_count
		FROM users u
		JOIN companies c ON u.id = c.user_id
		LEFT JOIN company_quotas cq ON c.id = cq.company_id
		WHERE u.id = ? AND u.role = 'company'
	`

	var (
		detail CompanyDetailResponse
		companyIndustry, companySize, companyLocation, companyAddress,
		companyCity, companyProvince, postalCode string
		establishedYear, employeeCount           int64
		verificationNotes, phone                 string
		docsVerifiedAt, emailVerifiedAt          string
		ktpURL, aktaURL, npwpURL, nibURL         string
		freeQuotaUsed, paidQuota, draftJobsCount int64
	)

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&detail.ID, &detail.Email, &detail.FullName, &phone,
		&detail.CompanyName, &detail.CompanyDescription, &detail.CompanyWebsite, &detail.CompanyLogoURL,
		&companyIndustry, &companySize, &companyLocation, &companyAddress,
		&companyCity, &companyProvince, &postalCode,
		&establishedYear, &employeeCount,
		&detail.CompanyStatus, &detail.IsActive, &detail.IsVerified, &emailVerifiedAt,
		&docsVerifiedAt, &verificationNotes,
		&ktpURL, &aktaURL, &npwpURL, &nibURL,
		&detail.CreatedAt, &detail.UpdatedAt,
		&detail.JobsCount, &detail.ActiveJobsCount, &detail.TotalApplications,
		&freeQuotaUsed, &paidQuota, &draftJobsCount,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	// Set HashID
	detail.HashID = hashid.Encode(detail.ID)

	// Set fields from string/int variables (they're never nullable now)
	detail.Phone = phone
	detail.CompanyIndustry = companyIndustry
	detail.CompanySize = companySize
	detail.CompanyLocation = companyLocation
	detail.CompanyAddress = companyAddress
	detail.CompanyCity = companyCity
	detail.CompanyProvince = companyProvince
	detail.PostalCode = postalCode
	detail.EstablishedYear = int(establishedYear)
	detail.EmployeeCount = int(employeeCount)
	detail.DocumentsVerifiedAt = docsVerifiedAt
	detail.EmailVerifiedAt = emailVerifiedAt
	detail.VerificationNotes = verificationNotes

	// Set legal documents (no longer nullable)
	detail.LegalDocuments.KtpFounderURL = ktpURL
	detail.LegalDocuments.AktaPendirianURL = aktaURL
	detail.LegalDocuments.NpwpURL = npwpURL
	detail.LegalDocuments.NibURL = nibURL

	// Set quota info
	const FreeQuotaTotal = 10
	detail.QuotaInfo.FreeQuotaUsed = int(freeQuotaUsed)
	detail.QuotaInfo.FreeQuotaTotal = FreeQuotaTotal
	detail.QuotaInfo.PaidQuota = int(paidQuota)
	detail.QuotaInfo.TotalQuota = detail.QuotaInfo.FreeQuotaUsed + detail.QuotaInfo.PaidQuota

	// Calculate job posting details
	// Free jobs: min of (active_jobs_count, free_quota_used)
	// Paid jobs: max of (active_jobs_count - free_quota_used, 0)
	detail.QuotaInfo.FreeJobsActive = detail.ActiveJobsCount
	if detail.QuotaInfo.FreeJobsActive > detail.QuotaInfo.FreeQuotaUsed {
		detail.QuotaInfo.FreeJobsActive = detail.QuotaInfo.FreeQuotaUsed
	}

	detail.QuotaInfo.PaidJobsActive = detail.ActiveJobsCount - detail.QuotaInfo.FreeJobsActive
	if detail.QuotaInfo.PaidJobsActive < 0 {
		detail.QuotaInfo.PaidJobsActive = 0
	}

	detail.QuotaInfo.TotalJobsActive = detail.ActiveJobsCount
	detail.QuotaInfo.DraftJobsCount = int(draftJobsCount)

	return &detail, nil
}

func (r *repository) UpdateCompanyStatus(ctx context.Context, id uint64, status string) error {
	// id is user_id, need to update companies table where user_id = id
	// If status is being set to "verified", also set documents_verified_at
	var query string
	var args []interface{}

	if status == "verified" {
		query = `UPDATE companies SET company_status = ?, documents_verified_at = NOW() WHERE user_id = ?`
		args = []interface{}{status, id}
	} else {
		query = `UPDATE companies SET company_status = ? WHERE user_id = ?`
		args = []interface{}{status, id}
	}

	fmt.Printf("[DEBUG] UpdateCompanyStatus: query=%s, status=%s, user_id=%d\n", query, status, id)

	result, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		fmt.Printf("[DEBUG] ExecContext error: %v\n", err)
		return fmt.Errorf("database error: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		fmt.Printf("[DEBUG] RowsAffected error: %v\n", err)
		return fmt.Errorf("rows affected error: %w", err)
	}

	fmt.Printf("[DEBUG] Rows affected: %d\n", rowsAffected)

	if rowsAffected == 0 {
		return fmt.Errorf("no company found with user_id %d", id)
	}

	return nil
}

func (r *repository) UpdateCompanyActive(ctx context.Context, id uint64, isActive bool) error {
	query := `UPDATE users SET is_active = ?, updated_at = NOW() WHERE id = ? AND role = 'company'`
	_, err := r.db.ExecContext(ctx, query, isActive, id)
	return err
}

// ============================================
// JOB OPERATIONS
// ============================================

func (r *repository) GetJobs(ctx context.Context, filter JobFilter) ([]*JobAdmin, int, error) {
	// Build WHERE clause
	conditions := []string{}
	args := []interface{}{}

	if filter.CompanyID > 0 {
		conditions = append(conditions, "j.company_id = ?")
		args = append(args, filter.CompanyID)
	}
	if filter.Status != "" {
		conditions = append(conditions, "j.status = ?")
		args = append(args, filter.Status)
	}
	if filter.Search != "" {
		conditions = append(conditions, "(j.title LIKE ? OR u.company_name LIKE ?)")
		searchTerm := "%" + filter.Search + "%"
		args = append(args, searchTerm, searchTerm)
	}
	if filter.DateFrom != "" {
		conditions = append(conditions, "j.created_at >= ?")
		args = append(args, filter.DateFrom)
	}
	if filter.DateTo != "" {
		conditions = append(conditions, "j.created_at <= ?")
		args = append(args, filter.DateTo+" 23:59:59")
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = " WHERE " + strings.Join(conditions, " AND ")
	}

	// Count query
	countQuery := `SELECT COUNT(*) FROM jobs j LEFT JOIN users u ON j.company_id = u.id` + whereClause
	var total int
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	// Main query
	query := `
		SELECT 
			j.id, j.company_id, u.company_name, j.title, j.slug, j.description, j.requirements,
			j.city, j.province, j.is_remote, j.job_type, j.experience_level,
			j.salary_min, j.salary_max, j.status, j.admin_status, j.admin_note, j.flag_reason,
			j.views_count, j.applications_count, j.published_at, j.created_at, j.updated_at
		FROM jobs j
		LEFT JOIN users u ON j.company_id = u.id
	` + whereClause + " ORDER BY j.created_at DESC"

	// Pagination
	if filter.PageSize <= 0 {
		filter.PageSize = 10
	}
	if filter.Page <= 0 {
		filter.Page = 1
	}
	offset := (filter.Page - 1) * filter.PageSize
	query += fmt.Sprintf(" LIMIT %d OFFSET %d", filter.PageSize, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var jobs []*JobAdmin
	for rows.Next() {
		var j JobAdmin
		if err := rows.Scan(
			&j.ID, &j.CompanyID, &j.CompanyName, &j.Title, &j.Slug, &j.Description, &j.Requirements,
			&j.City, &j.Province, &j.IsRemote, &j.JobType, &j.ExperienceLevel,
			&j.SalaryMin, &j.SalaryMax, &j.Status, &j.AdminStatus, &j.AdminNote, &j.FlagReason,
			&j.ViewsCount, &j.ApplicationsCount, &j.PublishedAt, &j.CreatedAt, &j.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}
		jobs = append(jobs, &j)
	}

	return jobs, total, nil
}

func (r *repository) GetJobByID(ctx context.Context, id uint64) (*JobAdmin, error) {
	query := `
		SELECT 
			j.id, j.company_id, u.company_name, j.title, j.slug, j.description, j.requirements,
			j.city, j.province, j.is_remote, j.job_type, j.experience_level,
			j.salary_min, j.salary_max, j.status, j.admin_status, j.admin_note, j.flag_reason,
			j.views_count, j.applications_count, j.published_at, j.created_at, j.updated_at
		FROM jobs j
		LEFT JOIN users u ON j.company_id = u.id
		WHERE j.id = ?
	`

	var j JobAdmin
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&j.ID, &j.CompanyID, &j.CompanyName, &j.Title, &j.Slug, &j.Description, &j.Requirements,
		&j.City, &j.Province, &j.IsRemote, &j.JobType, &j.ExperienceLevel,
		&j.SalaryMin, &j.SalaryMax, &j.Status, &j.AdminStatus, &j.AdminNote, &j.FlagReason,
		&j.ViewsCount, &j.ApplicationsCount, &j.PublishedAt, &j.CreatedAt, &j.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &j, nil
}

func (r *repository) UpdateJobStatus(ctx context.Context, id uint64, status string) error {
	query := `UPDATE jobs SET status = ?, updated_at = NOW() WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, status, id)
	return err
}

func (r *repository) UpdateJobAdminStatus(ctx context.Context, id uint64, adminStatus, note string) error {
	query := `UPDATE jobs SET admin_status = ?, admin_note = ?, updated_at = NOW() WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, adminStatus, note, id)
	return err
}

// ============================================
// PAYMENT OPERATIONS
// ============================================

func (r *repository) GetPayments(ctx context.Context, filter PaymentFilter) ([]*PaymentAdmin, int, error) {
	// Build WHERE clause
	conditions := []string{}
	args := []interface{}{}

	if filter.CompanyID > 0 {
		conditions = append(conditions, "p.company_id = ?")
		args = append(args, filter.CompanyID)
	}
	if filter.Status != "" {
		conditions = append(conditions, "p.status = ?")
		args = append(args, filter.Status)
	}
	if filter.DateFrom != "" {
		conditions = append(conditions, "p.submitted_at >= ?")
		args = append(args, filter.DateFrom)
	}
	if filter.DateTo != "" {
		conditions = append(conditions, "p.submitted_at <= ?")
		args = append(args, filter.DateTo+" 23:59:59")
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = " WHERE " + strings.Join(conditions, " AND ")
	}

	// Count query
	countQuery := `SELECT COUNT(*) FROM payments p` + whereClause
	var total int
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	// Main query
	query := `
		SELECT 
			p.id, p.company_id, c.company_name, p.job_id, j.title as job_title,
			p.amount, p.proof_image_url, p.status, p.note, p.confirmed_by_id,
			p.submitted_at, p.confirmed_at, p.created_at, p.updated_at
		FROM payments p
		LEFT JOIN companies c ON p.company_id = c.id
		LEFT JOIN jobs j ON p.job_id = j.id
	` + whereClause + " ORDER BY p.submitted_at DESC"

	// Pagination
	if filter.PageSize <= 0 {
		filter.PageSize = 10
	}
	if filter.Page <= 0 {
		filter.Page = 1
	}
	offset := (filter.Page - 1) * filter.PageSize
	query += fmt.Sprintf(" LIMIT %d OFFSET %d", filter.PageSize, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var payments []*PaymentAdmin
	for rows.Next() {
		var p PaymentAdmin
		if err := rows.Scan(
			&p.ID, &p.CompanyID, &p.CompanyName, &p.JobID, &p.JobTitle,
			&p.Amount, &p.ProofImageURL, &p.Status, &p.Note, &p.ConfirmedByID,
			&p.SubmittedAt, &p.ConfirmedAt, &p.CreatedAt, &p.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}
		payments = append(payments, &p)
	}

	return payments, total, nil
}

func (r *repository) GetPaymentByID(ctx context.Context, id uint64) (*PaymentAdmin, error) {
	query := `
		SELECT 
			p.id, p.company_id, c.company_name, p.job_id, j.title as job_title,
			p.amount, p.proof_image_url, p.status, p.note, p.confirmed_by_id,
			p.submitted_at, p.confirmed_at, p.created_at, p.updated_at
		FROM payments p
		LEFT JOIN companies c ON p.company_id = c.id
		LEFT JOIN jobs j ON p.job_id = j.id
		WHERE p.id = ?
	`

	var p PaymentAdmin
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&p.ID, &p.CompanyID, &p.CompanyName, &p.JobID, &p.JobTitle,
		&p.Amount, &p.ProofImageURL, &p.Status, &p.Note, &p.ConfirmedByID,
		&p.SubmittedAt, &p.ConfirmedAt, &p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &p, nil
}

func (r *repository) UpdatePaymentStatus(ctx context.Context, id uint64, status, note string, confirmedByID uint64) error {
	query := `
		UPDATE payments 
		SET status = ?, note = ?, confirmed_by_id = ?, confirmed_at = NOW(), updated_at = NOW() 
		WHERE id = ?
	`
	_, err := r.db.ExecContext(ctx, query, status, note, confirmedByID, id)
	return err
}

// ============================================
// JOB SEEKER OPERATIONS
// ============================================

func (r *repository) GetJobSeekers(ctx context.Context, filter JobSeekerFilter) ([]*JobSeekerAdmin, int, error) {
	// Count query
	countQuery := `
		SELECT COUNT(*) FROM users 
		WHERE role = 'job_seeker'
	`
	args := []interface{}{}

	if filter.Status == "active" {
		countQuery += " AND is_active = true"
	} else if filter.Status == "inactive" {
		countQuery += " AND is_active = false"
	}
	if filter.Search != "" {
		countQuery += " AND (full_name LIKE ? OR email LIKE ?)"
		searchTerm := "%" + filter.Search + "%"
		args = append(args, searchTerm, searchTerm)
	}

	var total int
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	// Main query
	query := `
		SELECT 
			u.id, u.email, u.full_name, u.phone, u.avatar_url, u.is_active, u.is_verified,
			u.email_verified_at, u.created_at, u.updated_at,
			(SELECT COUNT(*) FROM applications WHERE user_id = u.id) as applications_count,
			(SELECT COUNT(*) > 0 FROM cvs WHERE user_id = u.id) as has_cv
		FROM users u
		WHERE u.role = 'job_seeker'
	`

	queryArgs := []interface{}{}
	if filter.Status == "active" {
		query += " AND u.is_active = true"
	} else if filter.Status == "inactive" {
		query += " AND u.is_active = false"
	}
	if filter.Search != "" {
		query += " AND (u.full_name LIKE ? OR u.email LIKE ?)"
		searchTerm := "%" + filter.Search + "%"
		queryArgs = append(queryArgs, searchTerm, searchTerm)
	}

	query += " ORDER BY u.created_at DESC"

	// Pagination
	if filter.PageSize <= 0 {
		filter.PageSize = 10
	}
	if filter.Page <= 0 {
		filter.Page = 1
	}
	offset := (filter.Page - 1) * filter.PageSize
	query += fmt.Sprintf(" LIMIT %d OFFSET %d", filter.PageSize, offset)

	rows, err := r.db.QueryContext(ctx, query, queryArgs...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var jobSeekers []*JobSeekerAdmin
	for rows.Next() {
		var js JobSeekerAdmin
		if err := rows.Scan(
			&js.ID, &js.Email, &js.FullName, &js.Phone, &js.AvatarURL, &js.IsActive, &js.IsVerified,
			&js.EmailVerifiedAt, &js.CreatedAt, &js.UpdatedAt,
			&js.ApplicationsCount, &js.HasCV,
		); err != nil {
			return nil, 0, err
		}
		jobSeekers = append(jobSeekers, &js)
	}

	return jobSeekers, total, nil
}

func (r *repository) GetJobSeekerByID(ctx context.Context, id uint64) (*JobSeekerAdmin, error) {
	query := `
		SELECT 
			u.id, u.email, u.full_name, u.phone, u.avatar_url, u.is_active, u.is_verified,
			u.email_verified_at, u.created_at, u.updated_at,
			(SELECT COUNT(*) FROM applications WHERE user_id = u.id) as applications_count,
			(SELECT COUNT(*) > 0 FROM cvs WHERE user_id = u.id) as has_cv
		FROM users u
		WHERE u.id = ? AND u.role = 'job_seeker'
	`

	var js JobSeekerAdmin
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&js.ID, &js.Email, &js.FullName, &js.Phone, &js.AvatarURL, &js.IsActive, &js.IsVerified,
		&js.EmailVerifiedAt, &js.CreatedAt, &js.UpdatedAt,
		&js.ApplicationsCount, &js.HasCV,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &js, nil
}

func (r *repository) UpdateJobSeekerActive(ctx context.Context, id uint64, isActive bool) error {
	query := `UPDATE users SET is_active = ?, updated_at = NOW() WHERE id = ? AND role = 'job_seeker'`
	_, err := r.db.ExecContext(ctx, query, isActive, id)
	return err
}

// ============================================
// AUDIT LOG OPERATIONS
// ============================================

func (r *repository) LogAdminAction(ctx context.Context, log *AdminActionLog) error {
	query := `
		INSERT INTO audit_logs (user_id, action, entity_type, entity_id, details, ip_address, created_at)
		VALUES (?, ?, ?, ?, ?, ?, NOW())
	`
	_, err := r.db.ExecContext(ctx, query,
		log.AdminID, log.Action, log.EntityType, log.EntityID, log.Details, log.IPAddress,
	)
	return err
}
