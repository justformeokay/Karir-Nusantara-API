package admin

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

// PartnerRepository handles partner-related database operations for admin
type PartnerRepository interface {
	// Partner management
	GetPartners(ctx context.Context, status string, search string, page, limit int) ([]PartnerDBRow, int, error)
	GetPartnerByID(ctx context.Context, id uint64) (*PartnerDBRow, error)
	UpdatePartnerStatus(ctx context.Context, id uint64, status string, notes *string) error
	ApprovePartner(ctx context.Context, id uint64, approvedBy uint64, commissionRate float64, notes *string) error

	// Referred companies
	GetReferredCompanies(ctx context.Context, search string, page, limit int) ([]ReferredCompanyDBRow, int, error)
	GetReferredCompanyStats(ctx context.Context, companyID uint64) (int, int64, int64, error) // transactions, revenue, commission

	// Payouts
	GetPayouts(ctx context.Context, status string, search string, page, limit int) ([]PayoutDBRow, int, error)
	GetPayoutByID(ctx context.Context, id uint64) (*PayoutDBRow, error)
	CreatePayout(ctx context.Context, partnerID uint64, amount int64, notes *string) (uint64, error)
	ProcessPayout(ctx context.Context, id uint64, proofURL string, notes *string) error

	// Partner balances
	GetPartnersWithBalance(ctx context.Context, page, limit int) ([]PartnerDBRow, int, error)

	// Stats
	GetReferralStats(ctx context.Context) (*AdminReferralStatsResponse, error)
	GetPayoutStats(ctx context.Context) (*AdminPayoutStatsResponse, error)

	// Get last payout date for partner
	GetLastPayoutDate(ctx context.Context, partnerID uint64) (*time.Time, error)
}

type partnerRepository struct {
	db *sqlx.DB
}

// NewPartnerRepository creates a new partner repository for admin
func NewPartnerRepository(db *sqlx.DB) PartnerRepository {
	return &partnerRepository{db: db}
}

// GetPartners returns paginated list of partners
func (r *partnerRepository) GetPartners(ctx context.Context, status string, search string, page, limit int) ([]PartnerDBRow, int, error) {
	offset := (page - 1) * limit

	// Base query
	baseQuery := `
		FROM referral_partners rp
		JOIN users u ON rp.user_id = u.id
		WHERE u.deleted_at IS NULL
	`
	args := []interface{}{}

	if status != "" && status != "all" {
		baseQuery += " AND rp.status = ?"
		args = append(args, status)
	}

	if search != "" {
		baseQuery += " AND (u.full_name LIKE ? OR u.email LIKE ? OR rp.referral_code LIKE ?)"
		searchTerm := "%" + search + "%"
		args = append(args, searchTerm, searchTerm, searchTerm)
	}

	// Count query
	countQuery := "SELECT COUNT(*) " + baseQuery
	var total int
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("failed to count partners: %w", err)
	}

	// Data query
	dataQuery := `
		SELECT 
			rp.id, rp.user_id, u.full_name, u.email, u.phone,
			rp.referral_code, rp.commission_rate, rp.status,
			rp.bank_name, rp.bank_account_number, rp.bank_account_holder,
			rp.is_bank_verified, rp.total_referrals, rp.total_commission,
			rp.available_balance, rp.pending_balance, rp.paid_amount,
			rp.approved_by, rp.approved_at, rp.notes,
			rp.created_at, rp.updated_at
	` + baseQuery + `
		ORDER BY rp.created_at DESC
		LIMIT ? OFFSET ?
	`
	args = append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, dataQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query partners: %w", err)
	}
	defer rows.Close()

	var partners []PartnerDBRow
	for rows.Next() {
		var p PartnerDBRow
		err := rows.Scan(
			&p.ID, &p.UserID, &p.FullName, &p.Email, &p.Phone,
			&p.ReferralCode, &p.CommissionRate, &p.Status,
			&p.BankName, &p.BankAccountNumber, &p.BankAccountHolder,
			&p.IsBankVerified, &p.TotalReferrals, &p.TotalCommission,
			&p.AvailableBalance, &p.PendingBalance, &p.PaidAmount,
			&p.ApprovedBy, &p.ApprovedAt, &p.Notes,
			&p.CreatedAt, &p.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan partner: %w", err)
		}
		partners = append(partners, p)
	}

	return partners, total, nil
}

// GetPartnerByID returns a partner by ID
func (r *partnerRepository) GetPartnerByID(ctx context.Context, id uint64) (*PartnerDBRow, error) {
	query := `
		SELECT 
			rp.id, rp.user_id, u.full_name, u.email, u.phone,
			rp.referral_code, rp.commission_rate, rp.status,
			rp.bank_name, rp.bank_account_number, rp.bank_account_holder,
			rp.is_bank_verified, rp.total_referrals, rp.total_commission,
			rp.available_balance, rp.pending_balance, rp.paid_amount,
			rp.approved_by, rp.approved_at, rp.notes,
			rp.created_at, rp.updated_at
		FROM referral_partners rp
		JOIN users u ON rp.user_id = u.id
		WHERE rp.id = ? AND u.deleted_at IS NULL
	`

	var p PartnerDBRow
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&p.ID, &p.UserID, &p.FullName, &p.Email, &p.Phone,
		&p.ReferralCode, &p.CommissionRate, &p.Status,
		&p.BankName, &p.BankAccountNumber, &p.BankAccountHolder,
		&p.IsBankVerified, &p.TotalReferrals, &p.TotalCommission,
		&p.AvailableBalance, &p.PendingBalance, &p.PaidAmount,
		&p.ApprovedBy, &p.ApprovedAt, &p.Notes,
		&p.CreatedAt, &p.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get partner: %w", err)
	}

	return &p, nil
}

// UpdatePartnerStatus updates partner status
func (r *partnerRepository) UpdatePartnerStatus(ctx context.Context, id uint64, status string, notes *string) error {
	query := `
		UPDATE referral_partners 
		SET status = ?, notes = COALESCE(?, notes), updated_at = NOW()
		WHERE id = ?
	`
	_, err := r.db.ExecContext(ctx, query, status, notes, id)
	if err != nil {
		return fmt.Errorf("failed to update partner status: %w", err)
	}
	return nil
}

// ApprovePartner approves a pending partner
func (r *partnerRepository) ApprovePartner(ctx context.Context, id uint64, approvedBy uint64, commissionRate float64, notes *string) error {
	query := `
		UPDATE referral_partners 
		SET status = 'active', 
			approved_by = ?, 
			approved_at = NOW(),
			commission_rate = ?,
			notes = COALESCE(?, notes),
			updated_at = NOW()
		WHERE id = ?
	`
	_, err := r.db.ExecContext(ctx, query, approvedBy, commissionRate, notes, id)
	if err != nil {
		return fmt.Errorf("failed to approve partner: %w", err)
	}
	return nil
}

// GetReferredCompanies returns paginated list of referred companies
func (r *partnerRepository) GetReferredCompanies(ctx context.Context, search string, page, limit int) ([]ReferredCompanyDBRow, int, error) {
	offset := (page - 1) * limit

	baseQuery := `
		FROM partner_referrals pr
		JOIN referral_partners rp ON pr.partner_id = rp.id
		JOIN users u ON rp.user_id = u.id
		JOIN companies c ON pr.company_id = c.id
		WHERE 1=1
	`
	args := []interface{}{}

	if search != "" {
		baseQuery += " AND (c.company_name LIKE ? OR u.full_name LIKE ? OR rp.referral_code LIKE ?)"
		searchTerm := "%" + search + "%"
		args = append(args, searchTerm, searchTerm, searchTerm)
	}

	// Count
	countQuery := "SELECT COUNT(*) " + baseQuery
	var total int
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("failed to count referred companies: %w", err)
	}

	// Data
	dataQuery := `
		SELECT 
			pr.id, pr.company_id, c.company_name,
			pr.partner_id, u.full_name as partner_name,
			rp.referral_code, pr.registered_at as registration_date,
			c.verification_status as company_status
	` + baseQuery + `
		ORDER BY pr.registered_at DESC
		LIMIT ? OFFSET ?
	`
	args = append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, dataQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query referred companies: %w", err)
	}
	defer rows.Close()

	var companies []ReferredCompanyDBRow
	for rows.Next() {
		var c ReferredCompanyDBRow
		err := rows.Scan(
			&c.ID, &c.CompanyID, &c.CompanyName,
			&c.PartnerID, &c.PartnerName,
			&c.ReferralCode, &c.RegistrationDate,
			&c.CompanyStatus,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan referred company: %w", err)
		}
		companies = append(companies, c)
	}

	return companies, total, nil
}

// GetReferredCompanyStats returns transaction stats for a company
func (r *partnerRepository) GetReferredCompanyStats(ctx context.Context, companyID uint64) (int, int64, int64, error) {
	query := `
		SELECT 
			COALESCE(COUNT(*), 0),
			COALESCE(SUM(transaction_amount), 0),
			COALESCE(SUM(commission_amount), 0)
		FROM partner_commissions
		WHERE company_id = ?
	`
	var transactions int
	var revenue, commission int64
	err := r.db.QueryRowContext(ctx, query, companyID).Scan(&transactions, &revenue, &commission)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("failed to get company stats: %w", err)
	}
	return transactions, revenue, commission, nil
}

// GetPayouts returns paginated list of payouts
func (r *partnerRepository) GetPayouts(ctx context.Context, status string, search string, page, limit int) ([]PayoutDBRow, int, error) {
	offset := (page - 1) * limit

	baseQuery := `
		FROM partner_payouts pp
		JOIN referral_partners rp ON pp.partner_id = rp.id
		JOIN users u ON rp.user_id = u.id
		WHERE 1=1
	`
	args := []interface{}{}

	if status != "" && status != "all" {
		baseQuery += " AND pp.status = ?"
		args = append(args, status)
	}

	if search != "" {
		baseQuery += " AND u.full_name LIKE ?"
		args = append(args, "%"+search+"%")
	}

	// Count
	countQuery := "SELECT COUNT(*) " + baseQuery
	var total int
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("failed to count payouts: %w", err)
	}

	// Data
	dataQuery := `
		SELECT 
			pp.id, pp.partner_id, u.full_name as partner_name,
			pp.amount, pp.status, pp.transfer_ref,
			pp.created_at as requested_at, pp.completed_at, pp.notes,
			pp.bank_name, pp.bank_account_number, pp.bank_account_holder
	` + baseQuery + `
		ORDER BY pp.created_at DESC
		LIMIT ? OFFSET ?
	`
	args = append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, dataQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query payouts: %w", err)
	}
	defer rows.Close()

	var payouts []PayoutDBRow
	for rows.Next() {
		var p PayoutDBRow
		err := rows.Scan(
			&p.ID, &p.PartnerID, &p.PartnerName,
			&p.Amount, &p.Status, &p.PayoutProofURL,
			&p.RequestedAt, &p.PaidAt, &p.Notes,
			&p.BankName, &p.BankAccountNumber, &p.BankAccountHolder,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan payout: %w", err)
		}
		payouts = append(payouts, p)
	}

	return payouts, total, nil
}

// GetPayoutByID returns a payout by ID
func (r *partnerRepository) GetPayoutByID(ctx context.Context, id uint64) (*PayoutDBRow, error) {
	query := `
		SELECT 
			pp.id, pp.partner_id, u.full_name as partner_name,
			pp.amount, pp.status, pp.transfer_ref,
			pp.created_at as requested_at, pp.completed_at, pp.notes,
			pp.bank_name, pp.bank_account_number, pp.bank_account_holder
		FROM partner_payouts pp
		JOIN referral_partners rp ON pp.partner_id = rp.id
		JOIN users u ON rp.user_id = u.id
		WHERE pp.id = ?
	`

	var p PayoutDBRow
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&p.ID, &p.PartnerID, &p.PartnerName,
		&p.Amount, &p.Status, &p.PayoutProofURL,
		&p.RequestedAt, &p.PaidAt, &p.Notes,
		&p.BankName, &p.BankAccountNumber, &p.BankAccountHolder,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get payout: %w", err)
	}

	return &p, nil
}

// CreatePayout creates a new payout request
func (r *partnerRepository) CreatePayout(ctx context.Context, partnerID uint64, amount int64, notes *string) (uint64, error) {
	// Get partner bank info
	var bankName, bankAccountNumber, bankAccountHolder string
	bankQuery := `SELECT COALESCE(bank_name, ''), COALESCE(bank_account_number, ''), COALESCE(bank_account_holder, '') FROM referral_partners WHERE id = ?`
	if err := r.db.QueryRowContext(ctx, bankQuery, partnerID).Scan(&bankName, &bankAccountNumber, &bankAccountHolder); err != nil {
		return 0, fmt.Errorf("failed to get partner bank info: %w", err)
	}

	query := `
		INSERT INTO partner_payouts (partner_id, amount, bank_name, bank_account_number, bank_account_holder, status, notes, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, 'pending', ?, NOW(), NOW())
	`
	result, err := r.db.ExecContext(ctx, query, partnerID, amount, bankName, bankAccountNumber, bankAccountHolder, notes)
	if err != nil {
		return 0, fmt.Errorf("failed to create payout: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get payout id: %w", err)
	}

	return uint64(id), nil
}

// ProcessPayout marks a payout as paid
func (r *partnerRepository) ProcessPayout(ctx context.Context, id uint64, proofURL string, notes *string) error {
	// Get payout details first
	payout, err := r.GetPayoutByID(ctx, id)
	if err != nil {
		return err
	}
	if payout == nil {
		return fmt.Errorf("payout not found")
	}

	// Update payout status - using transfer_ref to store proof URL
	query := `
		UPDATE partner_payouts 
		SET status = 'completed', transfer_ref = ?, completed_at = NOW(), notes = COALESCE(?, notes), updated_at = NOW()
		WHERE id = ?
	`
	_, err = r.db.ExecContext(ctx, query, proofURL, notes, id)
	if err != nil {
		return fmt.Errorf("failed to update payout: %w", err)
	}

	// Update partner balance
	updateBalanceQuery := `
		UPDATE referral_partners 
		SET available_balance = available_balance - ?, 
			paid_amount = paid_amount + ?,
			updated_at = NOW()
		WHERE id = ?
	`
	_, err = r.db.ExecContext(ctx, updateBalanceQuery, payout.Amount, payout.Amount, payout.PartnerID)
	if err != nil {
		return fmt.Errorf("failed to update partner balance: %w", err)
	}

	return nil
}

// GetPartnersWithBalance returns partners with available balance > 0
func (r *partnerRepository) GetPartnersWithBalance(ctx context.Context, page, limit int) ([]PartnerDBRow, int, error) {
	offset := (page - 1) * limit

	// Count
	countQuery := `
		SELECT COUNT(*) 
		FROM referral_partners rp
		JOIN users u ON rp.user_id = u.id
		WHERE rp.available_balance > 0 AND u.deleted_at IS NULL
	`
	var total int
	if err := r.db.QueryRowContext(ctx, countQuery).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("failed to count partners with balance: %w", err)
	}

	// Data
	query := `
		SELECT 
			rp.id, rp.user_id, u.full_name, u.email, u.phone,
			rp.referral_code, rp.commission_rate, rp.status,
			rp.bank_name, rp.bank_account_number, rp.bank_account_holder,
			rp.is_bank_verified, rp.total_referrals, rp.total_commission,
			rp.available_balance, rp.pending_balance, rp.paid_amount,
			rp.approved_by, rp.approved_at, rp.notes,
			rp.created_at, rp.updated_at
		FROM referral_partners rp
		JOIN users u ON rp.user_id = u.id
		WHERE rp.available_balance > 0 AND u.deleted_at IS NULL
		ORDER BY rp.available_balance DESC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query partners with balance: %w", err)
	}
	defer rows.Close()

	var partners []PartnerDBRow
	for rows.Next() {
		var p PartnerDBRow
		err := rows.Scan(
			&p.ID, &p.UserID, &p.FullName, &p.Email, &p.Phone,
			&p.ReferralCode, &p.CommissionRate, &p.Status,
			&p.BankName, &p.BankAccountNumber, &p.BankAccountHolder,
			&p.IsBankVerified, &p.TotalReferrals, &p.TotalCommission,
			&p.AvailableBalance, &p.PendingBalance, &p.PaidAmount,
			&p.ApprovedBy, &p.ApprovedAt, &p.Notes,
			&p.CreatedAt, &p.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan partner: %w", err)
		}
		partners = append(partners, p)
	}

	return partners, total, nil
}

// GetReferralStats returns overall referral program statistics
func (r *partnerRepository) GetReferralStats(ctx context.Context) (*AdminReferralStatsResponse, error) {
	query := `
		SELECT 
			(SELECT COUNT(*) FROM referral_partners rp JOIN users u ON rp.user_id = u.id WHERE u.deleted_at IS NULL) as total_partners,
			(SELECT COUNT(*) FROM referral_partners rp JOIN users u ON rp.user_id = u.id WHERE rp.status = 'active' AND u.deleted_at IS NULL) as active_partners,
			(SELECT COUNT(*) FROM referral_partners rp JOIN users u ON rp.user_id = u.id WHERE rp.status = 'pending' AND u.deleted_at IS NULL) as pending_partners,
			(SELECT COUNT(*) FROM partner_referrals) as total_referred_companies,
			(SELECT COALESCE(SUM(total_commission), 0) FROM referral_partners) as total_commission,
			(SELECT COALESCE(SUM(amount), 0) FROM partner_payouts WHERE status IN ('pending', 'processing')) as pending_payouts,
			(SELECT COALESCE(SUM(paid_amount), 0) FROM referral_partners) as total_paid_out,
			(SELECT COUNT(*) FROM referral_partners rp JOIN users u ON rp.user_id = u.id WHERE rp.available_balance > 0 AND u.deleted_at IS NULL) as partners_with_balance
	`

	var stats AdminReferralStatsResponse
	err := r.db.QueryRowContext(ctx, query).Scan(
		&stats.TotalPartners,
		&stats.ActivePartners,
		&stats.PendingPartners,
		&stats.TotalReferredCompanies,
		&stats.TotalCommissionGenerated,
		&stats.PendingPayouts,
		&stats.TotalPaidOut,
		&stats.PartnersWithBalance,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get referral stats: %w", err)
	}

	return &stats, nil
}

// GetPayoutStats returns payout statistics
func (r *partnerRepository) GetPayoutStats(ctx context.Context) (*AdminPayoutStatsResponse, error) {
	query := `
		SELECT 
			(SELECT COALESCE(SUM(total_commission), 0) FROM referral_partners) as total_commission,
			(SELECT COALESCE(SUM(amount), 0) FROM partner_payouts WHERE status IN ('pending', 'processing')) as pending_payouts,
			(SELECT COALESCE(SUM(paid_amount), 0) FROM referral_partners) as total_paid_out,
			(SELECT COUNT(*) FROM referral_partners rp JOIN users u ON rp.user_id = u.id WHERE rp.available_balance > 0 AND u.deleted_at IS NULL) as partners_with_balance
	`

	var stats AdminPayoutStatsResponse
	err := r.db.QueryRowContext(ctx, query).Scan(
		&stats.TotalCommissionGenerated,
		&stats.PendingPayouts,
		&stats.TotalPaidOut,
		&stats.PartnersWithBalance,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get payout stats: %w", err)
	}

	return &stats, nil
}

// GetLastPayoutDate returns the last payout date for a partner
func (r *partnerRepository) GetLastPayoutDate(ctx context.Context, partnerID uint64) (*time.Time, error) {
	query := `
		SELECT completed_at FROM partner_payouts 
		WHERE partner_id = ? AND status = 'completed' 
		ORDER BY completed_at DESC LIMIT 1
	`
	var paidAt *time.Time
	err := r.db.QueryRowContext(ctx, query, partnerID).Scan(&paidAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get last payout date: %w", err)
	}
	return paidAt, nil
}
