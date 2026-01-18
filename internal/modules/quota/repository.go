package quota

import (
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
)

// Repository handles database operations for quota
type Repository struct {
	db *sqlx.DB
}

// NewRepository creates a new quota repository
func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

// GetOrCreateCompanyQuota gets or creates a company's quota record
func (r *Repository) GetOrCreateCompanyQuota(companyID uint64) (*CompanyQuota, error) {
	quota := &CompanyQuota{}
	
	err := r.db.Get(quota, `
		SELECT * FROM company_quotas WHERE company_id = ?
	`, companyID)
	
	if err == sql.ErrNoRows {
		// Create new quota record
		result, err := r.db.Exec(`
			INSERT INTO company_quotas (company_id, free_quota_used, paid_quota, created_at, updated_at)
			VALUES (?, 0, 0, NOW(), NOW())
		`, companyID)
		if err != nil {
			return nil, fmt.Errorf("failed to create quota: %w", err)
		}
		
		id, _ := result.LastInsertId()
		return r.GetByID(uint64(id))
	}
	
	if err != nil {
		return nil, fmt.Errorf("failed to get quota: %w", err)
	}
	
	return quota, nil
}

// GetByID gets a quota record by ID
func (r *Repository) GetByID(id uint64) (*CompanyQuota, error) {
	quota := &CompanyQuota{}
	err := r.db.Get(quota, `SELECT * FROM company_quotas WHERE id = ?`, id)
	if err != nil {
		return nil, err
	}
	return quota, nil
}

// IncrementFreeQuotaUsed increments the free quota used count
func (r *Repository) IncrementFreeQuotaUsed(companyID uint64) error {
	_, err := r.db.Exec(`
		UPDATE company_quotas 
		SET free_quota_used = free_quota_used + 1, updated_at = NOW() 
		WHERE company_id = ?
	`, companyID)
	return err
}

// DecrementPaidQuota decrements the paid quota count
func (r *Repository) DecrementPaidQuota(companyID uint64) error {
	_, err := r.db.Exec(`
		UPDATE company_quotas 
		SET paid_quota = paid_quota - 1, updated_at = NOW() 
		WHERE company_id = ? AND paid_quota > 0
	`, companyID)
	return err
}

// AddPaidQuota adds to paid quota
func (r *Repository) AddPaidQuota(companyID uint64, count int) error {
	_, err := r.db.Exec(`
		UPDATE company_quotas 
		SET paid_quota = paid_quota + ?, updated_at = NOW() 
		WHERE company_id = ?
	`, count, companyID)
	return err
}

// CreatePayment creates a new payment record
func (r *Repository) CreatePayment(payment *Payment) error {
	result, err := r.db.Exec(`
		INSERT INTO payments (company_id, job_id, amount, proof_image_url, status, submitted_at, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, NOW(), NOW(), NOW())
	`, payment.CompanyID, payment.JobID, payment.Amount, payment.ProofImageURL, payment.Status)
	
	if err != nil {
		return fmt.Errorf("failed to create payment: %w", err)
	}
	
	id, _ := result.LastInsertId()
	payment.ID = uint64(id)
	return nil
}

// GetPaymentByID gets a payment by ID
func (r *Repository) GetPaymentByID(id uint64) (*Payment, error) {
	payment := &Payment{}
	err := r.db.Get(payment, `SELECT * FROM payments WHERE id = ?`, id)
	if err != nil {
		return nil, err
	}
	return payment, nil
}

// UpdatePaymentStatus updates a payment's status
func (r *Repository) UpdatePaymentStatus(id uint64, status string, confirmedByID *uint64, note string) error {
	query := `
		UPDATE payments 
		SET status = ?, note = ?, updated_at = NOW()
	`
	args := []interface{}{status, note}
	
	if status == PaymentStatusConfirmed || status == PaymentStatusRejected {
		query += `, confirmed_by_id = ?, confirmed_at = NOW()`
		args = append(args, confirmedByID)
	}
	
	query += ` WHERE id = ?`
	args = append(args, id)
	
	_, err := r.db.Exec(query, args...)
	return err
}

// ListPayments lists payments for a company
func (r *Repository) ListPayments(params PaymentListParams) ([]Payment, int, error) {
	var payments []Payment
	var total int
	
	baseQuery := `FROM payments WHERE company_id = ?`
	args := []interface{}{params.CompanyID}
	
	if params.Status != "" {
		baseQuery += ` AND status = ?`
		args = append(args, params.Status)
	}
	
	// Get total count
	countQuery := `SELECT COUNT(*) ` + baseQuery
	err := r.db.Get(&total, countQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	
	// Get paginated results
	offset := (params.Page - 1) * params.PerPage
	dataQuery := `SELECT * ` + baseQuery + ` ORDER BY submitted_at DESC LIMIT ? OFFSET ?`
	args = append(args, params.PerPage, offset)
	
	err = r.db.Select(&payments, dataQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	
	return payments, total, nil
}

// GetPendingPaymentsCount returns count of pending payments for a company
func (r *Repository) GetPendingPaymentsCount(companyID uint64) (int, error) {
	var count int
	err := r.db.Get(&count, `
		SELECT COUNT(*) FROM payments WHERE company_id = ? AND status = ?
	`, companyID, PaymentStatusPending)
	return count, err
}
