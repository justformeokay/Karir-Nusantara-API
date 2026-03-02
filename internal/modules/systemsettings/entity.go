package systemsettings

import "time"

// Default values — fallback when DB is empty
const (
	DefaultFreeQuotaLimit = 3
	DefaultPricePerJob    = 20000
	DefaultCurrency       = "IDR"
)

// Setting keys
const (
	KeyFreeQuotaLimit = "free_quota_limit"
	KeyPricePerJob    = "price_per_job"
	KeyCurrency       = "currency"
)

// DataType defines the value type stored in system_settings
type DataType string

const (
	DataTypeString  DataType = "string"
	DataTypeInteger DataType = "integer"
	DataTypeDecimal DataType = "decimal"
	DataTypeBoolean DataType = "boolean"
	DataTypeJSON    DataType = "json"
)

// SystemSetting represents a single key-value configuration
type SystemSetting struct {
	ID          uint64    `db:"id"          json:"id"`
	Key         string    `db:"key"         json:"key"`
	Value       string    `db:"value"       json:"value"`
	DataType    DataType  `db:"data_type"   json:"data_type"`
	Description string    `db:"description" json:"description"`
	UpdatedBy   *uint64   `db:"updated_by"  json:"updated_by,omitempty"`
	CreatedAt   time.Time `db:"created_at"  json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"  json:"updated_at"`
}

// QuotaPackage represents a purchasable quota bundle
type QuotaPackage struct {
	ID           uint64    `db:"id"            json:"id"`
	PackageID    string    `db:"package_id"    json:"package_id"`
	Name         string    `db:"name"          json:"name"`
	Quota        int       `db:"quota"         json:"quota"`
	BonusQuota   int       `db:"bonus_quota"   json:"bonus_quota"`
	Price        int64     `db:"price"         json:"price"`
	Description  string    `db:"description"   json:"description"`
	IsBestValue  bool      `db:"is_best_value" json:"is_best_value"`
	IsActive     bool      `db:"is_active"     json:"is_active"`
	DisplayOrder int       `db:"display_order" json:"display_order"`
	CreatedBy    *uint64   `db:"created_by"    json:"created_by,omitempty"`
	UpdatedBy    *uint64   `db:"updated_by"    json:"updated_by,omitempty"`
	CreatedAt    time.Time `db:"created_at"    json:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"    json:"updated_at"`
}

// ---- Request / Response DTOs ----

// GetSettingsResponse is the payload returned to the admin
type GetSettingsResponse struct {
	FreeQuotaLimit int            `json:"free_quota_limit"`
	PricePerJob    int64          `json:"price_per_job"`
	Currency       string         `json:"currency"`
	QuotaPackages  []QuotaPackage `json:"quota_packages"`
}

// UpdateSettingsRequest is the payload to update general settings
type UpdateSettingsRequest struct {
	FreeQuotaLimit *int    `json:"free_quota_limit"`
	PricePerJob    *int64  `json:"price_per_job"`
	Currency       *string `json:"currency"`
}

// UpsertQuotaPackageRequest is the payload for creating/updating a package
type UpsertQuotaPackageRequest struct {
	PackageID    string `json:"package_id"    validate:"required"`
	Name         string `json:"name"          validate:"required"`
	Quota        int    `json:"quota"         validate:"required,min=1"`
	BonusQuota   int    `json:"bonus_quota"`
	Price        int64  `json:"price"         validate:"required,min=0"`
	Description  string `json:"description"`
	IsBestValue  bool   `json:"is_best_value"`
	IsActive     bool   `json:"is_active"`
	DisplayOrder int    `json:"display_order"`
}

// PublicPricingResponse is used by company / job-seeker frontends (no admin-only data)
type PublicPricingResponse struct {
	FreeQuotaLimit int            `json:"free_quota_limit"`
	PricePerJob    int64          `json:"price_per_job"`
	Currency       string         `json:"currency"`
	QuotaPackages  []QuotaPackage `json:"quota_packages"`
}
