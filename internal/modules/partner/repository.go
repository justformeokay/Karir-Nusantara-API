package partner

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

// Repository defines the partner repository interface
type Repository interface {
	// Authentication
	GetPartnerUserByEmail(ctx context.Context, email string) (*PartnerUser, error)
	GetPartnerByUserID(ctx context.Context, userID uint64) (*ReferralPartner, error)
	GetPartnerByID(ctx context.Context, partnerID uint64) (*ReferralPartner, error)
	GetUserPasswordHash(ctx context.Context, userID uint64) (string, error)
	UpdateUserPassword(ctx context.Context, userID uint64, passwordHash string) error
	UpdateUserName(ctx context.Context, userID uint64, name string) error

	// Registration
	CreatePartnerUser(ctx context.Context, fullName, email, phone, passwordHash, referralCode string) (uint64, uint64, error)
	CheckEmailExists(ctx context.Context, email string) (bool, error)

	// Password Reset
	CreatePasswordResetToken(ctx context.Context, userID uint64, token string, expiresAt time.Time) error
	GetPasswordResetToken(ctx context.Context, token string) (uint64, time.Time, error)
	DeletePasswordResetToken(ctx context.Context, token string) error
	GetUserEmailByID(ctx context.Context, userID uint64) (string, error)

	// Dashboard
	GetDashboardStats(ctx context.Context, partnerID uint64) (*DashboardStatsResponse, error)
	GetMonthlyData(ctx context.Context, partnerID uint64, months int) ([]MonthlyDataResponse, error)

	// Companies
	GetReferredCompanies(ctx context.Context, partnerID uint64, page, limit int, search string) ([]CompanyResponse, int, error)
	GetCompaniesSummary(ctx context.Context, partnerID uint64) (*CompanySummaryResponse, error)

	// Transactions
	GetTransactions(ctx context.Context, partnerID uint64, page, limit int, search, status string) ([]TransactionResponse, int, error)
	GetTransactionsSummary(ctx context.Context, partnerID uint64) (*TransactionSummaryResponse, error)

	// Payouts
	GetPayoutInfo(ctx context.Context, partnerID uint64) (*PayoutInfoResponse, error)
	GetPayoutHistory(ctx context.Context, partnerID uint64, page, limit int) ([]PayoutHistoryResponse, int, error)
	GetLastPayoutDate(ctx context.Context, partnerID uint64) (*time.Time, error)
}

// repository implements Repository
type repository struct {
	db *sqlx.DB
}

// NewRepository creates a new partner repository
func NewRepository(db *sqlx.DB) Repository {
	return &repository{db: db}
}

// GetPartnerUserByEmail retrieves partner user by email
func (r *repository) GetPartnerUserByEmail(ctx context.Context, email string) (*PartnerUser, error) {
	query := `
		SELECT 
			u.id, u.email, u.full_name, u.role, u.is_active, u.is_verified, u.created_at,
			rp.id as partner_id, rp.referral_code, rp.commission_rate, rp.status as partner_status,
			rp.bank_name, rp.bank_account_number, rp.bank_account_holder, rp.is_bank_verified,
			rp.total_referrals, rp.total_commission, rp.available_balance, rp.pending_balance, rp.paid_amount
		FROM users u
		INNER JOIN referral_partners rp ON rp.user_id = u.id
		WHERE u.email = ? AND u.role = 'partner' AND u.deleted_at IS NULL
	`

	var p PartnerUser
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&p.ID, &p.Email, &p.FullName, &p.Role, &p.IsActive, &p.IsVerified, &p.CreatedAt,
		&p.PartnerID, &p.ReferralCode, &p.CommissionRate, &p.PartnerStatus,
		&p.BankName, &p.BankAccountNumber, &p.BankAccountHolder, &p.IsBankVerified,
		&p.TotalReferrals, &p.TotalCommission, &p.AvailableBalance, &p.PendingBalance, &p.PaidAmount,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get partner user by email: %w", err)
	}

	return &p, nil
}

// GetPartnerByUserID retrieves partner by user ID
func (r *repository) GetPartnerByUserID(ctx context.Context, userID uint64) (*ReferralPartner, error) {
	query := `
		SELECT id, user_id, referral_code, commission_rate, status,
			   bank_name, bank_account_number, bank_account_holder, is_bank_verified,
			   total_referrals, total_commission, available_balance, pending_balance, paid_amount,
			   approved_by, approved_at, notes, created_at, updated_at
		FROM referral_partners
		WHERE user_id = ?
	`

	var p ReferralPartner
	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&p.ID, &p.UserID, &p.ReferralCode, &p.CommissionRate, &p.Status,
		&p.BankName, &p.BankAccountNumber, &p.BankAccountHolder, &p.IsBankVerified,
		&p.TotalReferrals, &p.TotalCommission, &p.AvailableBalance, &p.PendingBalance, &p.PaidAmount,
		&p.ApprovedBy, &p.ApprovedAt, &p.Notes, &p.CreatedAt, &p.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get partner by user ID: %w", err)
	}

	return &p, nil
}

// GetPartnerByID retrieves partner by partner ID
func (r *repository) GetPartnerByID(ctx context.Context, partnerID uint64) (*ReferralPartner, error) {
	query := `
		SELECT id, user_id, referral_code, commission_rate, status,
			   bank_name, bank_account_number, bank_account_holder, is_bank_verified,
			   total_referrals, total_commission, available_balance, pending_balance, paid_amount,
			   approved_by, approved_at, notes, created_at, updated_at
		FROM referral_partners
		WHERE id = ?
	`

	var p ReferralPartner
	err := r.db.QueryRowContext(ctx, query, partnerID).Scan(
		&p.ID, &p.UserID, &p.ReferralCode, &p.CommissionRate, &p.Status,
		&p.BankName, &p.BankAccountNumber, &p.BankAccountHolder, &p.IsBankVerified,
		&p.TotalReferrals, &p.TotalCommission, &p.AvailableBalance, &p.PendingBalance, &p.PaidAmount,
		&p.ApprovedBy, &p.ApprovedAt, &p.Notes, &p.CreatedAt, &p.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get partner by ID: %w", err)
	}

	return &p, nil
}

// GetUserPasswordHash retrieves user password hash
func (r *repository) GetUserPasswordHash(ctx context.Context, userID uint64) (string, error) {
	var hash string
	query := `SELECT password_hash FROM users WHERE id = ? AND deleted_at IS NULL`
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&hash)
	if err != nil {
		return "", fmt.Errorf("failed to get user password hash: %w", err)
	}
	return hash, nil
}

// UpdateUserPassword updates user password
func (r *repository) UpdateUserPassword(ctx context.Context, userID uint64, passwordHash string) error {
	query := `UPDATE users SET password_hash = ?, updated_at = NOW() WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, passwordHash, userID)
	if err != nil {
		return fmt.Errorf("failed to update user password: %w", err)
	}
	return nil
}

// UpdateUserName updates user name
func (r *repository) UpdateUserName(ctx context.Context, userID uint64, name string) error {
	query := `UPDATE users SET full_name = ?, updated_at = NOW() WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, name, userID)
	if err != nil {
		return fmt.Errorf("failed to update user name: %w", err)
	}
	return nil
}

// GetDashboardStats retrieves dashboard statistics
func (r *repository) GetDashboardStats(ctx context.Context, partnerID uint64) (*DashboardStatsResponse, error) {
	// Get basic stats from referral_partners
	query := `
		SELECT 
			rp.total_referrals,
			rp.total_commission,
			rp.available_balance,
			rp.paid_amount,
			(SELECT COUNT(*) FROM partner_commissions WHERE partner_id = rp.id) as total_transactions
		FROM referral_partners rp
		WHERE rp.id = ?
	`

	var stats DashboardStatsResponse
	err := r.db.QueryRowContext(ctx, query, partnerID).Scan(
		&stats.TotalCompanies,
		&stats.TotalCommission,
		&stats.AvailableBalance,
		&stats.PaidCommission,
		&stats.TotalTransactions,
	)

	if err == sql.ErrNoRows {
		// Return empty stats if partner not found
		return &DashboardStatsResponse{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get dashboard stats: %w", err)
	}

	return &stats, nil
}

// GetMonthlyData retrieves monthly chart data
func (r *repository) GetMonthlyData(ctx context.Context, partnerID uint64, months int) ([]MonthlyDataResponse, error) {
	query := `
		SELECT 
			DATE_FORMAT(pc.created_at, '%b') as month,
			COALESCE(SUM(pc.commission_amount), 0) as commission,
			COUNT(DISTINCT pr.company_id) as companies
		FROM partner_commissions pc
		JOIN partner_referrals pr ON pr.id = pc.referral_id
		WHERE pc.partner_id = ?
		  AND pc.created_at >= DATE_SUB(NOW(), INTERVAL ? MONTH)
		  AND pc.status IN ('approved', 'paid')
		GROUP BY YEAR(pc.created_at), MONTH(pc.created_at), DATE_FORMAT(pc.created_at, '%b')
		ORDER BY YEAR(pc.created_at), MONTH(pc.created_at)
	`

	rows, err := r.db.QueryContext(ctx, query, partnerID, months)
	if err != nil {
		return nil, fmt.Errorf("failed to get monthly data: %w", err)
	}
	defer rows.Close()

	var result []MonthlyDataResponse
	for rows.Next() {
		var m MonthlyDataResponse
		if err := rows.Scan(&m.Month, &m.Commission, &m.Companies); err != nil {
			return nil, fmt.Errorf("failed to scan monthly data: %w", err)
		}
		result = append(result, m)
	}

	// If no data, return last 7 months with zero values
	if len(result) == 0 {
		now := time.Now()
		for i := months - 1; i >= 0; i-- {
			t := now.AddDate(0, -i, 0)
			result = append(result, MonthlyDataResponse{
				Month:      t.Format("Jan"),
				Commission: 0,
				Companies:  0,
			})
		}
	}

	return result, nil
}

// GetReferredCompanies retrieves referred companies with pagination
func (r *repository) GetReferredCompanies(ctx context.Context, partnerID uint64, page, limit int, search string) ([]CompanyResponse, int, error) {
	offset := (page - 1) * limit

	// Count query
	countQuery := `
		SELECT COUNT(*)
		FROM partner_referrals pr
		JOIN companies c ON c.id = pr.company_id
		WHERE pr.partner_id = ?
	`
	countArgs := []interface{}{partnerID}

	if search != "" {
		countQuery += " AND c.company_name LIKE ?"
		countArgs = append(countArgs, "%"+search+"%")
	}

	var total int
	if err := r.db.QueryRowContext(ctx, countQuery, countArgs...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("failed to count companies: %w", err)
	}

	// Data query
	dataQuery := `
		SELECT 
			c.id,
			c.company_name,
			c.created_at,
			CASE WHEN c.company_status = 'verified' THEN 'active' ELSE 'inactive' END as status,
			COALESCE(j.job_count, 0) as total_job_posts,
			COALESCE(pc.total_revenue, 0) as total_revenue,
			COALESCE(pc.total_commission, 0) as commission_earned
		FROM partner_referrals pr
		JOIN companies c ON c.id = pr.company_id
		LEFT JOIN (
			SELECT company_id, COUNT(*) as job_count
			FROM jobs
			GROUP BY company_id
		) j ON j.company_id = c.id
		LEFT JOIN (
			SELECT 
				company_id,
				SUM(transaction_amount) as total_revenue,
				SUM(commission_amount) as total_commission
			FROM partner_commissions
			WHERE status IN ('approved', 'paid')
			GROUP BY company_id
		) pc ON pc.company_id = c.id
		WHERE pr.partner_id = ?
	`
	dataArgs := []interface{}{partnerID}

	if search != "" {
		dataQuery += " AND c.company_name LIKE ?"
		dataArgs = append(dataArgs, "%"+search+"%")
	}

	dataQuery += " ORDER BY c.created_at DESC LIMIT ? OFFSET ?"
	dataArgs = append(dataArgs, limit, offset)

	rows, err := r.db.QueryContext(ctx, dataQuery, dataArgs...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get companies: %w", err)
	}
	defer rows.Close()

	var companies []CompanyResponse
	for rows.Next() {
		var c CompanyResponse
		var createdAt time.Time
		var id uint64
		if err := rows.Scan(&id, &c.Name, &createdAt, &c.Status, &c.TotalJobPosts, &c.TotalRevenue, &c.CommissionEarned); err != nil {
			return nil, 0, fmt.Errorf("failed to scan company: %w", err)
		}
		c.ID = fmt.Sprintf("%d", id)
		c.RegistrationDate = createdAt.Format(time.RFC3339)
		companies = append(companies, c)
	}

	return companies, total, nil
}

// GetCompaniesSummary retrieves companies summary
func (r *repository) GetCompaniesSummary(ctx context.Context, partnerID uint64) (*CompanySummaryResponse, error) {
	query := `
		SELECT 
			COUNT(*) as total_companies,
			COALESCE(SUM(CASE WHEN c.company_status = 'verified' THEN 1 ELSE 0 END), 0) as active_companies,
			COALESCE(SUM(COALESCE(pc.total_revenue, 0)), 0) as total_revenue,
			COALESCE(SUM(COALESCE(pc.total_commission, 0)), 0) as total_commission
		FROM partner_referrals pr
		JOIN companies c ON c.id = pr.company_id
		LEFT JOIN (
			SELECT 
				company_id,
				SUM(transaction_amount) as total_revenue,
				SUM(commission_amount) as total_commission
			FROM partner_commissions
			WHERE status IN ('approved', 'paid')
			GROUP BY company_id
		) pc ON pc.company_id = c.id
		WHERE pr.partner_id = ?
	`

	var summary CompanySummaryResponse
	err := r.db.QueryRowContext(ctx, query, partnerID).Scan(
		&summary.TotalCompanies,
		&summary.ActiveCompanies,
		&summary.TotalRevenue,
		&summary.TotalCommission,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get companies summary: %w", err)
	}

	return &summary, nil
}

// GetTransactions retrieves transactions with pagination
func (r *repository) GetTransactions(ctx context.Context, partnerID uint64, page, limit int, search, status string) ([]TransactionResponse, int, error) {
	offset := (page - 1) * limit

	// Count query
	countQuery := `
		SELECT COUNT(*)
		FROM partner_commissions pc
		JOIN companies c ON c.id = pc.company_id
		WHERE pc.partner_id = ?
	`
	countArgs := []interface{}{partnerID}

	if search != "" {
		countQuery += " AND c.company_name LIKE ?"
		countArgs = append(countArgs, "%"+search+"%")
	}

	if status != "" && status != "all" {
		countQuery += " AND pc.status = ?"
		countArgs = append(countArgs, status)
	}

	var total int
	if err := r.db.QueryRowContext(ctx, countQuery, countArgs...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("failed to count transactions: %w", err)
	}

	// Data query
	dataQuery := `
		SELECT 
			pc.id,
			pc.created_at,
			c.company_name,
			c.id as company_id,
			pc.job_quota,
			pc.transaction_amount,
			pc.commission_amount,
			pc.status
		FROM partner_commissions pc
		JOIN companies c ON c.id = pc.company_id
		WHERE pc.partner_id = ?
	`
	dataArgs := []interface{}{partnerID}

	if search != "" {
		dataQuery += " AND c.company_name LIKE ?"
		dataArgs = append(dataArgs, "%"+search+"%")
	}

	if status != "" && status != "all" {
		dataQuery += " AND pc.status = ?"
		dataArgs = append(dataArgs, status)
	}

	dataQuery += " ORDER BY pc.created_at DESC LIMIT ? OFFSET ?"
	dataArgs = append(dataArgs, limit, offset)

	rows, err := r.db.QueryContext(ctx, dataQuery, dataArgs...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get transactions: %w", err)
	}
	defer rows.Close()

	var transactions []TransactionResponse
	for rows.Next() {
		var t TransactionResponse
		var createdAt time.Time
		var id, companyID uint64
		if err := rows.Scan(&id, &createdAt, &t.CompanyName, &companyID, &t.JobQuota, &t.Amount, &t.Commission, &t.Status); err != nil {
			return nil, 0, fmt.Errorf("failed to scan transaction: %w", err)
		}
		t.ID = fmt.Sprintf("%d", id)
		t.CompanyID = fmt.Sprintf("%d", companyID)
		t.Date = createdAt.Format(time.RFC3339)
		transactions = append(transactions, t)
	}

	return transactions, total, nil
}

// GetTransactionsSummary retrieves transaction summary
func (r *repository) GetTransactionsSummary(ctx context.Context, partnerID uint64) (*TransactionSummaryResponse, error) {
	query := `
		SELECT 
			COALESCE(SUM(commission_amount), 0) as total_commission,
			COALESCE(SUM(CASE WHEN status = 'pending' THEN commission_amount ELSE 0 END), 0) as pending_commission,
			COALESCE(SUM(CASE WHEN status = 'paid' THEN commission_amount ELSE 0 END), 0) as paid_commission,
			COUNT(*) as total_transactions
		FROM partner_commissions
		WHERE partner_id = ?
	`

	var summary TransactionSummaryResponse
	err := r.db.QueryRowContext(ctx, query, partnerID).Scan(
		&summary.TotalCommission,
		&summary.PendingCommission,
		&summary.PaidCommission,
		&summary.TotalTransactions,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get transactions summary: %w", err)
	}

	return &summary, nil
}

// GetPayoutInfo retrieves payout information
func (r *repository) GetPayoutInfo(ctx context.Context, partnerID uint64) (*PayoutInfoResponse, error) {
	// Get partner info
	query := `
		SELECT 
			available_balance, pending_balance, paid_amount,
			bank_name, bank_account_number, bank_account_holder, is_bank_verified
		FROM referral_partners
		WHERE id = ?
	`

	var availableBalance, pendingBalance, paidAmount int64
	var bankName, accountNumber, accountHolder sql.NullString
	var isBankVerified bool

	err := r.db.QueryRowContext(ctx, query, partnerID).Scan(
		&availableBalance, &pendingBalance, &paidAmount,
		&bankName, &accountNumber, &accountHolder, &isBankVerified,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get payout info: %w", err)
	}

	// Get last payout date
	lastPayoutDate, _ := r.GetLastPayoutDate(ctx, partnerID)

	info := &PayoutInfoResponse{
		Balance: &BalanceInfo{
			AvailableBalance: availableBalance,
			PendingPayout:    pendingBalance,
			PaidAmount:       paidAmount,
		},
		Schedule: &PayoutSchedule{
			ProcessingDay: 20,
			MinimumPayout: 500000,
			TransferTime:  "1-3 Business Days",
		},
	}

	if lastPayoutDate != nil {
		info.Balance.LastPayoutDate = lastPayoutDate.Format(time.RFC3339)
	}

	if bankName.Valid && bankName.String != "" {
		info.BankAccount = &BankAccount{
			BankName:      bankName.String,
			AccountNumber: MaskAccountNumber(accountNumber.String),
			AccountHolder: accountHolder.String,
			IsVerified:    isBankVerified,
		}
	}

	return info, nil
}

// GetPayoutHistory retrieves payout history
func (r *repository) GetPayoutHistory(ctx context.Context, partnerID uint64, page, limit int) ([]PayoutHistoryResponse, int, error) {
	offset := (page - 1) * limit

	// Count query
	var total int
	countQuery := `SELECT COUNT(*) FROM partner_payouts WHERE partner_id = ?`
	if err := r.db.QueryRowContext(ctx, countQuery, partnerID).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("failed to count payouts: %w", err)
	}

	// Data query
	dataQuery := `
		SELECT id, amount, status, transfer_ref, processed_at, created_at
		FROM partner_payouts
		WHERE partner_id = ?
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.QueryContext(ctx, dataQuery, partnerID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get payout history: %w", err)
	}
	defer rows.Close()

	var payouts []PayoutHistoryResponse
	for rows.Next() {
		var p PayoutHistoryResponse
		var id uint64
		var transferRef sql.NullString
		var processedAt sql.NullTime
		var createdAt time.Time

		if err := rows.Scan(&id, &p.Amount, &p.Status, &transferRef, &processedAt, &createdAt); err != nil {
			return nil, 0, fmt.Errorf("failed to scan payout: %w", err)
		}

		p.ID = fmt.Sprintf("%d", id)
		p.CreatedAt = createdAt.Format(time.RFC3339)
		if transferRef.Valid {
			p.TransferRef = transferRef.String
		}
		if processedAt.Valid {
			p.ProcessedAt = processedAt.Time.Format(time.RFC3339)
		}

		payouts = append(payouts, p)
	}

	return payouts, total, nil
}

// GetLastPayoutDate retrieves the last payout date
func (r *repository) GetLastPayoutDate(ctx context.Context, partnerID uint64) (*time.Time, error) {
	query := `
		SELECT completed_at 
		FROM partner_payouts 
		WHERE partner_id = ? AND status = 'completed'
		ORDER BY completed_at DESC
		LIMIT 1
	`

	var completedAt time.Time
	err := r.db.QueryRowContext(ctx, query, partnerID).Scan(&completedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &completedAt, nil
}

// CreatePartnerUser creates a new partner user and referral_partner record
func (r *repository) CreatePartnerUser(ctx context.Context, fullName, email, phone, passwordHash, referralCode string) (uint64, uint64, error) {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback()

	// Insert user
	userQuery := `
		INSERT INTO users (email, password_hash, full_name, phone, role, is_active, is_verified, created_at, updated_at)
		VALUES (?, ?, ?, ?, 'partner', 1, 0, NOW(), NOW())
	`
	result, err := tx.ExecContext(ctx, userQuery, email, passwordHash, fullName, phone)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to insert user: %w", err)
	}

	userID, err := result.LastInsertId()
	if err != nil {
		return 0, 0, fmt.Errorf("failed to get user id: %w", err)
	}

	// Insert referral_partner with pending status
	partnerQuery := `
		INSERT INTO referral_partners (user_id, referral_code, commission_rate, status, created_at, updated_at)
		VALUES (?, ?, 40.00, 'pending', NOW(), NOW())
	`
	result, err = tx.ExecContext(ctx, partnerQuery, userID, referralCode)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to insert partner: %w", err)
	}

	partnerID, err := result.LastInsertId()
	if err != nil {
		return 0, 0, fmt.Errorf("failed to get partner id: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return 0, 0, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return uint64(userID), uint64(partnerID), nil
}

// CheckEmailExists checks if email already exists
func (r *repository) CheckEmailExists(ctx context.Context, email string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE email = ? AND deleted_at IS NULL)`
	var exists bool
	err := r.db.QueryRowContext(ctx, query, email).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check email: %w", err)
	}
	return exists, nil
}

// CreatePasswordResetToken creates a password reset token
func (r *repository) CreatePasswordResetToken(ctx context.Context, userID uint64, token string, expiresAt time.Time) error {
	// Delete any existing tokens for this user
	_, _ = r.db.ExecContext(ctx, `DELETE FROM password_resets WHERE user_id = ?`, userID)

	query := `
		INSERT INTO password_resets (user_id, token, expires_at, created_at)
		VALUES (?, ?, ?, NOW())
	`
	_, err := r.db.ExecContext(ctx, query, userID, token, expiresAt)
	if err != nil {
		return fmt.Errorf("failed to create password reset token: %w", err)
	}
	return nil
}

// GetPasswordResetToken retrieves password reset token info
func (r *repository) GetPasswordResetToken(ctx context.Context, token string) (uint64, time.Time, error) {
	query := `
		SELECT user_id, expires_at 
		FROM password_resets 
		WHERE token = ?
	`
	var userID uint64
	var expiresAt time.Time
	err := r.db.QueryRowContext(ctx, query, token).Scan(&userID, &expiresAt)
	if err == sql.ErrNoRows {
		return 0, time.Time{}, fmt.Errorf("token not found")
	}
	if err != nil {
		return 0, time.Time{}, fmt.Errorf("failed to get token: %w", err)
	}
	return userID, expiresAt, nil
}

// DeletePasswordResetToken deletes a password reset token
func (r *repository) DeletePasswordResetToken(ctx context.Context, token string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM password_resets WHERE token = ?`, token)
	return err
}

// GetUserEmailByID retrieves user email by ID
func (r *repository) GetUserEmailByID(ctx context.Context, userID uint64) (string, error) {
	var email string
	err := r.db.QueryRowContext(ctx, `SELECT email FROM users WHERE id = ?`, userID).Scan(&email)
	if err != nil {
		return "", fmt.Errorf("failed to get user email: %w", err)
	}
	return email, nil
}
