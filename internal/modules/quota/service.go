package quota

import (
	"fmt"

	"github.com/karirnusantara/api/internal/modules/systemsettings"
)

// Service handles business logic for quota
type Service struct {
	repo                  *Repository
	systemSettingsService *systemsettings.Service
}

// NewService creates a new quota service
func NewService(repo *Repository, ssService *systemsettings.Service) *Service {
	return &Service{repo: repo, systemSettingsService: ssService}
}

// getFreeQuotaLimit returns the current free quota limit from settings
func (s *Service) getFreeQuotaLimit() int {
	if s.systemSettingsService != nil {
		if limit, err := s.systemSettingsService.GetFreeQuotaLimit(); err == nil {
			return limit
		}
	}
	return FreeQuotaLimit // fallback to constant
}

// getPricePerJob returns the current price per job from settings
func (s *Service) getPricePerJob() int64 {
	if s.systemSettingsService != nil {
		if price, err := s.systemSettingsService.GetPricePerJob(); err == nil {
			return price
		}
	}
	return PricePerJob // fallback to constant
}

// GetQuota gets the quota information for a company
func (s *Service) GetQuota(companyID uint64) (*QuotaResponse, error) {
	quota, err := s.repo.GetOrCreateCompanyQuota(companyID)
	if err != nil {
		return nil, err
	}

	freeQuotaLimit := s.getFreeQuotaLimit()
	pricePerJob := s.getPricePerJob()

	remainingFree := freeQuotaLimit - quota.FreeQuotaUsed
	if remainingFree < 0 {
		remainingFree = 0
	}

	return &QuotaResponse{
		FreeQuota:          freeQuotaLimit,
		UsedFreeQuota:      quota.FreeQuotaUsed,
		RemainingFreeQuota: remainingFree,
		PaidQuota:          quota.PaidQuota,
		PricePerJob:        pricePerJob,
	}, nil
}

// CanPublishJob checks if a company can publish a job
func (s *Service) CanPublishJob(companyID uint64) (bool, string, error) {
	quota, err := s.repo.GetOrCreateCompanyQuota(companyID)
	if err != nil {
		return false, "", err
	}

	freeQuotaLimit := s.getFreeQuotaLimit()

	// Check if has free quota
	if quota.FreeQuotaUsed < freeQuotaLimit {
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

	freeQuotaLimit := s.getFreeQuotaLimit()

	// Use free quota first
	if quota.FreeQuotaUsed < freeQuotaLimit {
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
	// Default: single posting payment uses current price per job
	var amount int64 = s.getPricePerJob()
	var quotaAmount int = 1

	// If package specified, get package details from database
	if packageID != nil && *packageID != "" {
		// Try to get from database first
		found := false
		if s.systemSettingsService != nil {
			pricing, err := s.systemSettingsService.GetPublicPricing()
			if err == nil && pricing != nil {
				for _, pkg := range pricing.QuotaPackages {
					if pkg.PackageID == *packageID {
						amount = pkg.Price
						quotaAmount = pkg.Quota + pkg.BonusQuota
						found = true
						break
					}
				}
			}
		}

		// Fallback to hardcoded packages only if not found in DB
		if !found {
			pkg := GetPackageByID(*packageID)
			if pkg != nil {
				amount = pkg.Price
				quotaAmount = pkg.TotalQuota
			}
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

// GetPaymentByID gets a single payment by ID
func (s *Service) GetPaymentByID(paymentID uint64) (*Payment, error) {
	return s.repo.GetPaymentByID(paymentID)
}
