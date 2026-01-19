-- ============================================================
-- DEBUG - Check Admin User in Database
-- ============================================================

-- 1. Lihat semua user dengan role 'admin'
SELECT 
    id, 
    email, 
    full_name, 
    role, 
    is_active, 
    is_verified,
    password_hash,
    created_at
FROM users 
WHERE role = 'admin'
ORDER BY created_at DESC;

-- 2. Check khusus admin@karirnusantara.com
SELECT 
    id, 
    email, 
    full_name, 
    role, 
    is_active, 
    is_verified,
    LENGTH(password_hash) as hash_length,
    password_hash,
    created_at
FROM users 
WHERE email = 'admin@karirnusantara.com';

-- 3. Jika user tidak ada atau ada yang salah, recreate dengan REPLACE
-- Uncomment untuk jalankan:
/*
REPLACE INTO users (
    id,
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
    1,
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
*/
