# Password Reset & Email Notification API

## Overview
API ini menyediakan fitur reset password dan notifikasi email otomatis untuk:
1. **Welcome Email** - Dikirim saat company mendaftar
2. **Password Reset Email** - Dikirim saat forgot password

## Email Configuration

### SMTP Settings
```
Host: mail.karyadeveloperindonesia.com
Port: 587 (TLS)
User: no-reply@karyadeveloperindonesia.com
Password: Justformeokay23
From Name: Karir Nusantara
```

### Environment Variables
Tambahkan di file `.env`:
```env
SMTP_HOST=mail.karyadeveloperindonesia.com
SMTP_PORT=587
SMTP_USER=no-reply@karyadeveloperindonesia.com
SMTP_PASSWORD=Justformeokay23
MAIL_FROM_NAME=Karir Nusantara
MAIL_FROM_EMAIL=no-reply@karyadeveloperindonesia.com
```

## API Endpoints

### 1. Forgot Password
Request reset password token.

**Endpoint:** `POST /api/v1/auth/forgot-password`

**Request Body:**
```json
{
  "email": "info@karyadeveloperindonesia.com"
}
```

**Response (Success):**
```json
{
  "success": true,
  "message": "If your email is registered, you will receive password reset instructions",
  "data": null
}
```

**Notes:**
- Selalu return success untuk mencegah email enumeration
- Email akan dikirim jika email terdaftar
- Token berlaku selama 1 jam
- Link reset: `https://company.karirnusantara.com/reset-password?token=<TOKEN>`

### 2. Reset Password
Reset password menggunakan token.

**Endpoint:** `POST /api/v1/auth/reset-password`

**Request Body:**
```json
{
  "token": "abc123def456...",
  "new_password": "NewPassword123!"
}
```

**Response (Success):**
```json
{
  "success": true,
  "message": "Password reset successful",
  "data": null
}
```

**Error Responses:**
```json
{
  "success": false,
  "error": "bad_request",
  "message": "Invalid or expired reset token"
}
```

**Notes:**
- Token hanya bisa digunakan sekali
- Token expired setelah 1 jam
- Semua refresh token user akan di-revoke setelah reset password

### 3. Register (Updated with Email)
Registrasi company sekarang otomatis mengirim welcome email.

**Endpoint:** `POST /api/v1/auth/register`

**Request Body:**
```json
{
  "email": "company@example.com",
  "password": "SecurePassword123!",
  "full_name": "John Doe",
  "role": "company",
  "company_name": "PT Example Indonesia"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Registration successful",
  "data": {
    "user": { ... },
    "access_token": "...",
    "refresh_token": "...",
    "expires_in": 3600
  }
}
```

**Notes:**
- Welcome email dikirim secara async (tidak memblokir response)
- Email hanya dikirim untuk role "company"

## Database Schema

### Table: password_reset_tokens
```sql
CREATE TABLE `password_reset_tokens` (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT,
  `email` varchar(255) NOT NULL,
  `token` varchar(255) NOT NULL,
  `expires_at` timestamp NOT NULL,
  `used_at` timestamp NULL DEFAULT NULL,
  `created_at` timestamp NOT NULL DEFAULT current_timestamp(),
  PRIMARY KEY (`id`),
  KEY `idx_email` (`email`),
  KEY `idx_token` (`token`),
  KEY `idx_expires_at` (`expires_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
```

## Email Templates

### Welcome Email
- **Subject:** "Selamat Datang di Karir Nusantara"
- **Content:** Welcome message dengan link ke dashboard
- **CTA:** Login ke Dashboard

### Password Reset Email
- **Subject:** "Reset Password - Karir Nusantara"
- **Content:** Instruksi reset password dengan link
- **CTA:** Reset Password button
- **Warning:** Link valid 1 jam, jangan bagikan ke siapapun

## Testing

### Test Forgot Password
```bash
curl -X POST http://localhost:8081/api/v1/auth/forgot-password \
  -H "Content-Type: application/json" \
  -d '{
    "email": "info@karyadeveloperindonesia.com"
  }'
```

### Test Reset Password
```bash
curl -X POST http://localhost:8081/api/v1/auth/reset-password \
  -H "Content-Type: application/json" \
  -d '{
    "token": "YOUR_TOKEN_HERE",
    "new_password": "NewPassword123!"
  }'
```

### Test Registration with Email
```bash
curl -X POST http://localhost:8081/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "SecurePass123!",
    "full_name": "Test User",
    "role": "company",
    "company_name": "PT Test Company"
  }'
```

## Security Features

1. **Email Enumeration Prevention**
   - Forgot password selalu return success
   - Tidak memberi tahu apakah email terdaftar atau tidak

2. **Token Security**
   - Token generated secara cryptographically secure
   - Token length: 64 characters (32 bytes hex)
   - Token hanya valid 1 jam
   - Token hanya bisa digunakan sekali

3. **Password Security**
   - Password di-hash dengan bcrypt
   - Minimum 8 karakter
   - Harus memenuhi kompleksitas password

4. **Session Management**
   - Semua refresh token di-revoke setelah reset password
   - User harus login ulang setelah reset

## Monitoring & Logging

- Email sending dilakukan async (goroutine)
- Error email tidak memblokir request
- Error di-log ke console (production: use proper logging)
- Monitor failed email deliveries

## Troubleshooting

### Email tidak terkirim
1. Check SMTP credentials
2. Check SMTP host & port accessible
3. Check firewall rules
4. Verify sender email domain

### Token invalid/expired
1. Token hanya valid 1 jam
2. Token hanya bisa digunakan sekali
3. Request token baru jika sudah expired

### Password validation failed
1. Minimum 8 karakter
2. Harus memenuhi kompleksitas (huruf besar, kecil, angka)

## Next Steps

1. **Frontend Implementation**
   - Buat halaman forgot password
   - Buat halaman reset password dengan token validation
   - Handle error messages

2. **Email Template Customization**
   - Sesuaikan template dengan brand
   - Tambahkan logo perusahaan
   - Customize colors & styling

3. **Monitoring**
   - Setup email delivery monitoring
   - Alert untuk failed email deliveries
   - Track reset password usage

## Migration

Jalankan migration untuk membuat tabel:
```bash
cd /Users/putramac/Desktop/Loker/karir-nusantara-api
/Applications/XAMPP/xamppfiles/bin/mysql -u root karir_nusantara < migrations/006_password_reset_tokens.sql
```

Atau gunakan script:
```bash
./migrations/run_006_migration.sh
```
