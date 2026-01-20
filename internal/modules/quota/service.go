package quota

import (
	"fmt"
)

// Service handles business logic for quota
type Service struct {
	repo *Repository
}

// NewService creates a new quota service
func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// GetQuota gets the quota information for a company
func (s *Service) GetQuota(companyID uint64) (*QuotaResponse, error) {
	quota, err := s.repo.GetOrCreateCompanyQuota(companyID)
	if err != nil {
		return nil, err
	}
	
	remainingFree := FreeQuotaLimit - quota.FreeQuotaUsed
	if remainingFree < 0 {
		remainingFree = 0
	}
	
	return &QuotaResponse{
		FreeQuota:          FreeQuotaLimit,
		UsedFreeQuota:      quota.FreeQuotaUsed,
		RemainingFreeQuota: remainingFree,
		PaidQuota:          quota.PaidQuota,
		PricePerJob:        PricePerJob,
	}, nil
}

// CanPublishJob checks if a company can publish a job
func (s *Service) CanPublishJob(companyID uint64) (bool, string, error) {
	quota, err := s.repo.GetOrCreateCompanyQuota(companyID)
	if err != nil {
		return false, "", err
	}
	
	// Check if has free quota
	if quota.FreeQuotaUsed < FreeQuotaLimit {
		return true, "free", nil
	}
	
	// Check if has paid quota
	if quota.PaidQuota > 0 {
		return true, "paid", nil
	}
	
	return false, "", nil
}

// ConsumeQuota consumes quota for publishing a job
func (s *Service) ConsumeQuota(companyID uint64) error {
	quota, err := s.repo.GetOrCreateCompanyQuota(companyID)
	if err != nil {
		return err
	}
	
	// Use free quota first
	if quota.FreeQuotaUsed < FreeQuotaLimit {
		return s.repo.IncrementFreeQuotaUsed(companyID)
	}
	
	// Use paid quota
	if quota.PaidQuota > 0 {
		return s.repo.DecrementPaidQuota(companyID)
	}
	
	return fmt.Errorf("no quota available")
}

// SubmitPaymentProof submits a payment proof
func (s *Service) SubmitPaymentProof(companyID uint64, jobID *uint64, packageID *string, proofImageURL string) (*Payment, error) {
	// Determine amount and quota based on package or single payment
	var amount int64 = PricePerJob
	var quotaAmount int = 1
	
	if packageID != nil {
		pkg := GetPackageByID(*packageID)
		if pkg != nil {
			amount = pkg.Price
			quotaAmount = pkg.TotalQuota // Include bonus quota
		}
	}
	
	payment := &Payment{
		CompanyID:   companyID,
		Amount:      amount,
		QuotaAmount: quotaAmount,
		Status:      PaymentStatusPending,
	}
	
	if jobID != nil {
		payment.JobID.Valid = true
		payment.JobID.Int64 = int64(*jobID)
	}
	
	if packageID != nil {
		payment.PackageID.Valid = true
		payment.PackageID.String = *packageID
	}
	
	if proofImageURL != "" {
		payment.ProofImageURL.Valid = true
		payment.ProofImageURL.String = proofImageURL
	}
	
	err := s.repo.CreatePayment(payment)
	if err != nil {
		return nil, err
	}
	
	return payment, nil
}

// ConfirmPayment confirms a payment (admin only)
func (s *Service) ConfirmPayment(paymentID uint64, adminID uint64, note string) error {
	payment, err := s.repo.GetPaymentByID(paymentID)
	if err != nil {
		return err
	}
	
	if payment.Status != PaymentStatusPending {
		return fmt.Errorf("payment already processed")
	}
	
	// Update payment status
	err = s.repo.UpdatePaymentStatus(paymentID, PaymentStatusConfirmed, &adminID, note)
	if err != nil {
		return err
	}
	
	// Add paid quota to company (use QuotaAmount from payment which includes bonus)
	quotaToAdd := payment.QuotaAmount
	if quotaToAdd == 0 {
		quotaToAdd = 1 // Fallback for old payments
	}
	return s.repo.AddPaidQuota(payment.CompanyID, quotaToAdd)
}

// RejectPayment rejects a payment (admin only)
func (s *Service) RejectPayment(paymentID uint64, adminID uint64, note string) error {
	payment, err := s.repo.GetPaymentByID(paymentID)
	if err != nil {
		return err
	}
	
	if payment.Status != PaymentStatusPending {
		return fmt.Errorf("payment already processed")
	}
	
	return s.repo.UpdatePaymentStatus(paymentID, PaymentStatusRejected, &adminID, note)
}

// GetPayments gets the payment history for a company
func (s *Service) GetPayments(params PaymentListParams) ([]PaymentResponse, int, error) {
	payments, total, err := s.repo.ListPayments(params)
	if err != nil {
		return nil, 0, err
	}
	
	responses := make([]PaymentResponse, len(payments))
	for i, p := range payments {
		responses[i] = *p.ToResponse()
	}
	
	return responses, total, nil
}

// GetPendingPaymentsCount gets the count of pending payments
func (s *Service) GetPendingPaymentsCount(companyID uint64) (int, error) {
	return s.repo.GetPendingPaymentsCount(companyID)
}
