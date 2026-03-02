package systemsettings

import (
	"database/sql"
	"fmt"
	"strconv"

	"github.com/jmoiron/sqlx"
)

// Repository handles all DB operations for system settings & quota packages
type Repository struct {
	db *sqlx.DB
}

// NewRepository creates a new Repository
func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

// ---- System Settings ----

// GetByKey returns a single setting row. Returns nil, nil when not found (will fall back to default).
func (r *Repository) GetByKey(key string) (*SystemSetting, error) {
	row := &SystemSetting{}
	err := r.db.Get(row, `SELECT * FROM system_settings WHERE `+"`key`"+` = ?`, key)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get setting %q: %w", key, err)
	}
	return row, nil
}

// GetAll returns every setting row
func (r *Repository) GetAll() ([]SystemSetting, error) {
	var rows []SystemSetting
	err := r.db.Select(&rows, `SELECT * FROM system_settings ORDER BY `+"`key`")
	if err != nil {
		return nil, fmt.Errorf("get all settings: %w", err)
	}
	return rows, nil
}

// Upsert inserts or updates a setting value
func (r *Repository) Upsert(key, value string, dt DataType, description string, updatedBy *uint64) error {
	_, err := r.db.Exec(`
		INSERT INTO system_settings (`+"`key`"+`, value, data_type, description, updated_by, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, NOW(), NOW())
		ON DUPLICATE KEY UPDATE
			value      = VALUES(value),
			data_type  = VALUES(data_type),
			description= VALUES(description),
			updated_by = VALUES(updated_by),
			updated_at = NOW()
	`, key, value, string(dt), description, updatedBy)
	if err != nil {
		return fmt.Errorf("upsert setting %q: %w", key, err)
	}
	return nil
}

// GetIntSetting reads a setting as int, returning defaultVal when missing or invalid
func (r *Repository) GetIntSetting(key string, defaultVal int) (int, error) {
	row, err := r.GetByKey(key)
	if err != nil {
		return defaultVal, err
	}
	if row == nil {
		return defaultVal, nil
	}
	v, err := strconv.Atoi(row.Value)
	if err != nil {
		return defaultVal, nil
	}
	return v, nil
}

// GetInt64Setting reads a setting as int64, returning defaultVal when missing or invalid
func (r *Repository) GetInt64Setting(key string, defaultVal int64) (int64, error) {
	row, err := r.GetByKey(key)
	if err != nil {
		return defaultVal, err
	}
	if row == nil {
		return defaultVal, nil
	}
	v, err := strconv.ParseInt(row.Value, 10, 64)
	if err != nil {
		return defaultVal, nil
	}
	return v, nil
}

// GetStringSetting reads a setting as string, returning defaultVal when missing
func (r *Repository) GetStringSetting(key string, defaultVal string) (string, error) {
	row, err := r.GetByKey(key)
	if err != nil {
		return defaultVal, err
	}
	if row == nil {
		return defaultVal, nil
	}
	return row.Value, nil
}

// ---- Quota Packages ----

// GetAllPackages returns all packages ordered by display_order
func (r *Repository) GetAllPackages(activeOnly bool) ([]QuotaPackage, error) {
	query := `SELECT * FROM quota_packages`
	if activeOnly {
		query += ` WHERE is_active = 1`
	}
	query += ` ORDER BY display_order ASC, id ASC`

	var pkgs []QuotaPackage
	if err := r.db.Select(&pkgs, query); err != nil {
		return nil, fmt.Errorf("get packages: %w", err)
	}
	return pkgs, nil
}

// GetPackageByID returns a quota package by its database ID
func (r *Repository) GetPackageByID(id uint64) (*QuotaPackage, error) {
	pkg := &QuotaPackage{}
	err := r.db.Get(pkg, `SELECT * FROM quota_packages WHERE id = ?`, id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get package by id: %w", err)
	}
	return pkg, nil
}

// GetPackageByPackageID returns a quota package by its string package_id
func (r *Repository) GetPackageByPackageID(packageID string) (*QuotaPackage, error) {
	pkg := &QuotaPackage{}
	err := r.db.Get(pkg, `SELECT * FROM quota_packages WHERE package_id = ?`, packageID)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get package by package_id: %w", err)
	}
	return pkg, nil
}

// UpsertPackage inserts or updates a quota package using package_id as the unique key
func (r *Repository) UpsertPackage(req *UpsertQuotaPackageRequest, updatedBy *uint64) (*QuotaPackage, error) {
	_, err := r.db.Exec(`
		INSERT INTO quota_packages 
			(package_id, name, quota, bonus_quota, price, description, is_best_value, is_active, display_order, created_by, updated_by, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW())
		ON DUPLICATE KEY UPDATE
			name          = VALUES(name),
			quota         = VALUES(quota),
			bonus_quota   = VALUES(bonus_quota),
			price         = VALUES(price),
			description   = VALUES(description),
			is_best_value = VALUES(is_best_value),
			is_active     = VALUES(is_active),
			display_order = VALUES(display_order),
			updated_by    = VALUES(updated_by),
			updated_at    = NOW()
	`,
		req.PackageID, req.Name, req.Quota, req.BonusQuota, req.Price,
		req.Description, req.IsBestValue, req.IsActive, req.DisplayOrder,
		updatedBy, updatedBy,
	)
	if err != nil {
		return nil, fmt.Errorf("upsert package: %w", err)
	}
	return r.GetPackageByPackageID(req.PackageID)
}

// DeletePackage deletes a package by database ID
func (r *Repository) DeletePackage(id uint64) error {
	_, err := r.db.Exec(`DELETE FROM quota_packages WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("delete package: %w", err)
	}
	return nil
}
