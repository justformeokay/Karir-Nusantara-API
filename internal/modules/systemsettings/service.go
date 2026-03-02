package systemsettings

import (
	"fmt"
	"strconv"
)

// Service contains the business logic for system settings
type Service struct {
	repo *Repository
}

// NewService creates a new Service
func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// GetSettings returns the full admin settings payload
func (s *Service) GetSettings() (*GetSettingsResponse, error) {
	freeQuota, err := s.repo.GetIntSetting(KeyFreeQuotaLimit, DefaultFreeQuotaLimit)
	if err != nil {
		return nil, fmt.Errorf("get free_quota_limit: %w", err)
	}

	pricePerJob, err := s.repo.GetInt64Setting(KeyPricePerJob, DefaultPricePerJob)
	if err != nil {
		return nil, fmt.Errorf("get price_per_job: %w", err)
	}

	currency, err := s.repo.GetStringSetting(KeyCurrency, DefaultCurrency)
	if err != nil {
		return nil, fmt.Errorf("get currency: %w", err)
	}

	packages, err := s.repo.GetAllPackages(false) // admin sees all (active + inactive)
	if err != nil {
		return nil, fmt.Errorf("get packages: %w", err)
	}

	return &GetSettingsResponse{
		FreeQuotaLimit: freeQuota,
		PricePerJob:    pricePerJob,
		Currency:       currency,
		QuotaPackages:  packages,
	}, nil
}

// GetPublicPricing returns only the publicly-visible pricing info
func (s *Service) GetPublicPricing() (*PublicPricingResponse, error) {
	freeQuota, err := s.repo.GetIntSetting(KeyFreeQuotaLimit, DefaultFreeQuotaLimit)
	if err != nil {
		return nil, err
	}

	pricePerJob, err := s.repo.GetInt64Setting(KeyPricePerJob, DefaultPricePerJob)
	if err != nil {
		return nil, err
	}

	currency, err := s.repo.GetStringSetting(KeyCurrency, DefaultCurrency)
	if err != nil {
		return nil, err
	}

	packages, err := s.repo.GetAllPackages(true) // public only sees active packages
	if err != nil {
		return nil, err
	}

	return &PublicPricingResponse{
		FreeQuotaLimit: freeQuota,
		PricePerJob:    pricePerJob,
		Currency:       currency,
		QuotaPackages:  packages,
	}, nil
}

// UpdateSettings saves the general settings (null fields are skipped)
func (s *Service) UpdateSettings(req *UpdateSettingsRequest, adminID *uint64) error {
	if req.FreeQuotaLimit != nil {
		if *req.FreeQuotaLimit < 0 {
			return fmt.Errorf("free_quota_limit must be >= 0")
		}
		if err := s.repo.Upsert(
			KeyFreeQuotaLimit,
			strconv.Itoa(*req.FreeQuotaLimit),
			DataTypeInteger,
			"Jumlah kuota gratis untuk perusahaan baru yang mendaftar",
			adminID,
		); err != nil {
			return err
		}
	}

	if req.PricePerJob != nil {
		if *req.PricePerJob < 0 {
			return fmt.Errorf("price_per_job must be >= 0")
		}
		if err := s.repo.Upsert(
			KeyPricePerJob,
			strconv.FormatInt(*req.PricePerJob, 10),
			DataTypeInteger,
			"Harga dasar per posting lowongan kerja (IDR)",
			adminID,
		); err != nil {
			return err
		}
	}

	if req.Currency != nil {
		if err := s.repo.Upsert(
			KeyCurrency,
			*req.Currency,
			DataTypeString,
			"Mata uang yang digunakan pada platform",
			adminID,
		); err != nil {
			return err
		}
	}

	return nil
}

// GetFreeQuotaLimit returns the current free quota limit (used by the quota module)
func (s *Service) GetFreeQuotaLimit() (int, error) {
	return s.repo.GetIntSetting(KeyFreeQuotaLimit, DefaultFreeQuotaLimit)
}

// GetPricePerJob returns the current price per job posting (used by billing)
func (s *Service) GetPricePerJob() (int64, error) {
	return s.repo.GetInt64Setting(KeyPricePerJob, DefaultPricePerJob)
}

// ---- Quota Package management ----

// UpsertPackage creates or updates a quota package
func (s *Service) UpsertPackage(req *UpsertQuotaPackageRequest, adminID *uint64) (*QuotaPackage, error) {
	return s.repo.UpsertPackage(req, adminID)
}

// DeletePackage deletes a quota package by DB ID
func (s *Service) DeletePackage(id uint64) error {
	pkg, err := s.repo.GetPackageByID(id)
	if err != nil {
		return err
	}
	if pkg == nil {
		return fmt.Errorf("package not found")
	}
	return s.repo.DeletePackage(id)
}
