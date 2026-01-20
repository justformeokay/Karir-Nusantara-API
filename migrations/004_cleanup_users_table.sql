-- Drop company-related columns from users table since they are now in companies table
ALTER TABLE users DROP COLUMN IF EXISTS company_name;
ALTER TABLE users DROP COLUMN IF EXISTS company_description;
ALTER TABLE users DROP COLUMN IF EXISTS company_website;
ALTER TABLE users DROP COLUMN IF EXISTS company_logo_url;
ALTER TABLE users DROP COLUMN IF EXISTS company_status;
ALTER TABLE users DROP COLUMN IF EXISTS company_industry;
ALTER TABLE users DROP COLUMN IF EXISTS company_size;
ALTER TABLE users DROP COLUMN IF EXISTS company_location;
