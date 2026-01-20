-- Create companies table to store complete company information
CREATE TABLE IF NOT EXISTS companies (
    id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    user_id BIGINT UNSIGNED NOT NULL UNIQUE,
    company_name VARCHAR(255) NOT NULL,
    company_description LONGTEXT,
    company_website VARCHAR(255),
    company_logo_url VARCHAR(500),
    company_industry VARCHAR(100),
    company_size VARCHAR(50),
    company_location VARCHAR(255),
    company_phone VARCHAR(20),
    company_email VARCHAR(255),
    company_address LONGTEXT,
    company_city VARCHAR(100),
    company_province VARCHAR(100),
    company_postal_code VARCHAR(20),
    established_year YEAR,
    employee_count INT,
    
    -- Verification and status fields
    company_status ENUM('pending', 'verified', 'rejected', 'suspended') DEFAULT 'pending',
    
    -- Legal documents
    ktp_founder_url VARCHAR(500),
    akta_pendirian_url VARCHAR(500),
    npwp_url VARCHAR(500),
    nib_url VARCHAR(500),
    documents_verified_at TIMESTAMP NULL,
    documents_verified_by BIGINT UNSIGNED NULL,
    verification_notes LONGTEXT,
    
    -- Timestamps
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    
    -- Foreign keys
    CONSTRAINT fk_companies_user_id 
        FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_companies_verified_by 
        FOREIGN KEY (documents_verified_by) REFERENCES users(id) ON DELETE SET NULL,
    
    -- Indexes
    INDEX idx_company_status (company_status),
    INDEX idx_deleted_at (deleted_at),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
