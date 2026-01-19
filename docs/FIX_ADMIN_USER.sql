-- Fix Admin User with Valid Bcrypt Hash
-- Password: admin123
-- Valid Hash (generated and tested): $2a$10$llmqSju0EZuHEklAzhE7jeubjam0w2DoG252O2cgb1lG73gGL0AEG

-- Step 1: Check current admin user
SELECT id, email, password_hash, role, is_active, created_at FROM users WHERE email = 'admin@karirnusantara.com';

-- Step 2: Replace admin user with valid bcrypt hash and is_active = 1
REPLACE INTO users (id, email, password_hash, full_name, role, is_active, created_at, updated_at) 
VALUES (1, 'admin@karirnusantara.com', '$2a$10$llmqSju0EZuHEklAzhE7jeubjam0w2DoG252O2cgb1lG73gGL0AEG', 'Admin User', 'admin', 1, NOW(), NOW());

-- Step 3: Verify the update
SELECT id, email, password_hash, role, is_active FROM users WHERE email = 'admin@karirnusantara.com';
