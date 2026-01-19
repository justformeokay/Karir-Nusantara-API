-- ============================================================
-- KARIR NUSANTARA - CREATE SUPER ADMIN USER
-- ============================================================
-- Password: admin123
-- Bcrypt Hash: $2a$10$N9qo8uLOickgx2ZMRZoMyeIjZRGdjGj/n3.2U8xGTv5MBDj2iqpYi
-- ============================================================

-- Step 1: Cek apakah admin sudah ada
-- SELECT * FROM users WHERE email = 'admin@karirnusantara.com';

-- Step 2: Delete admin lama jika ada (opsional)
-- DELETE FROM users WHERE email = 'admin@karirnusantara.com';

-- Step 3: Insert/Update Admin User
-- Gunakan REPLACE untuk insert or update
REPLACE INTO users (
    email, 
    password_hash, 
    full_name, 
    role, 
    is_active, 
    is_verified, 
    email_verified_at, 
    created_at, 
    updated_at
) VALUES (
    'admin@karirnusantara.com',
    '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZRGdjGj/n3.2U8xGTv5MBDj2iqpYi',
    'Super Admin',
    'admin',
    1,
    1,
    NOW(),
    NOW(),
    NOW()
);

-- Step 4: Verify berhasil
SELECT id, email, full_name, role, is_active, is_verified FROM users WHERE email = 'admin@karirnusantara.com';

-- ============================================================
-- CREDENTIALS:
-- Email: admin@karirnusantara.com
-- Password: admin123
-- Role: admin
-- ============================================================
