# Password Management APIs - Quick Reference

Referensi cepat untuk semua API yang berkaitan dengan password management di Karir Nusantara.

---

## 📚 API Endpoints Overview

| Endpoint | Method | Auth | Purpose |
|----------|--------|------|---------|
| `/auth/forgot-password` | POST | ❌ No | Request password reset via email |
| `/auth/reset-password` | POST | ❌ No | Reset password using token from email |
| `/auth/change-password` | PUT | ✅ Yes | Change password (logged-in users) |

---

## 🔐 1. Forgot Password

**Untuk user yang lupa password mereka.**

### Request
```bash
curl -X POST http://localhost:8081/api/v1/auth/forgot-password \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com"
  }'
```

### Response
```json
{
  "success": true,
  "message": "If your email is registered, you will receive password reset instructions"
}
```

### What Happens
1. ✅ Generate secure reset token (64 chars, 1 hour expiry)
2. ✅ Save token to database
3. ✅ Send password reset email with link
4. ✅ Always return success (prevent email enumeration)

### Email Sent
- **Subject:** Reset Password - Karir Nusantara
- **Contains:** Reset link with token
- **Valid for:** 1 hour
- **Single use:** Token can only be used once

---

## 🔄 2. Reset Password

**Complete password reset using token from email.**

### Request
```bash
curl -X POST http://localhost:8081/api/v1/auth/reset-password \
  -H "Content-Type: application/json" \
  -d '{
    "token": "abc123...xyz789",
    "new_password": "NewPassword123!"
  }'
```

### Response
```json
{
  "success": true,
  "message": "Password reset successful"
}
```

### What Happens
1. ✅ Validate token (exists, not expired, not used)
2. ✅ Hash new password
3. ✅ Update password in database
4. ✅ Mark token as used
5. ✅ Revoke all refresh tokens (force re-login all devices)

### Errors
- Token invalid/expired: `"Invalid or expired reset token"`
- Token already used: `"Invalid or expired reset token"`

---

## 🔑 3. Change Password

**For logged-in users to change their password.**

### Request
```bash
curl -X PUT http://localhost:8081/api/v1/auth/change-password \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer {access_token}" \
  -d '{
    "old_password": "CurrentPassword123!",
    "new_password": "NewPassword456!"
  }'
```

### Response
```json
{
  "success": true,
  "message": "Password berhasil diubah. Silakan login kembali dengan password baru Anda."
}
```

### What Happens
1. ✅ Verify old password
2. ✅ Check new password != old password
3. ✅ Hash new password
4. ✅ Update password in database
5. ✅ Revoke all refresh tokens (force re-login)
6. ✅ Send confirmation email

### Email Sent
- **Subject:** Password Berhasil Diubah - Karir Nusantara
- **Contains:** Confirmation with security warnings
- **Purpose:** Notify user of password change (security)

### Errors
- Wrong old password: `"Password lama tidak sesuai"`
- Same password: `"Password baru tidak boleh sama dengan password lama"`
- Unauthorized: `"Missing authorization header"`

---

## 📋 Comparison Table

| Feature | Forgot Password | Reset Password | Change Password |
|---------|----------------|----------------|-----------------|
| **Auth Required** | ❌ No | ❌ No | ✅ Yes (Bearer) |
| **Old Password Required** | ❌ No | ❌ No | ✅ Yes |
| **Token Required** | ❌ No | ✅ Yes (from email) | ❌ No |
| **Email Sent** | ✅ Reset link | ❌ No | ✅ Confirmation |
| **Token Revocation** | ❌ No | ✅ Yes | ✅ Yes |
| **Use Case** | Forgot password | Complete reset flow | Regular password change |

---

## 🔄 Complete Workflows

### Workflow 1: Forgot Password Flow
```
User clicks "Forgot Password"
    ↓
POST /auth/forgot-password
    ↓
Email sent with reset link
    ↓
User clicks link (opens /reset-password?token=xxx)
    ↓
User enters new password
    ↓
POST /auth/reset-password
    ↓
Password updated
    ↓
User redirects to login
    ↓
User logs in with new password
```

### Workflow 2: Change Password Flow (Logged In)
```
User navigates to Settings/Profile
    ↓
User clicks "Change Password"
    ↓
User enters old + new password
    ↓
PUT /auth/change-password (with Bearer token)
    ↓
Password updated + Email sent
    ↓
All sessions logged out
    ↓
User redirects to login
    ↓
User logs in with new password
```

---

## 🛡️ Security Features

### All Endpoints Include:
- ✅ Password complexity validation (min 8, uppercase, lowercase, number)
- ✅ Bcrypt password hashing
- ✅ Token/session revocation after password change
- ✅ Rate limiting recommended (to prevent brute force)

### Forgot/Reset Specific:
- ✅ Secure token generation (crypto/rand, 32 bytes)
- ✅ Token expiry (1 hour)
- ✅ Single-use tokens
- ✅ Email enumeration prevention

### Change Password Specific:
- ✅ Old password verification
- ✅ Same password prevention
- ✅ Email notification for unauthorized change detection

---

## 📧 Email Templates

### 1. Password Reset Email
```
Subject: Reset Password - Karir Nusantara

Halo [Nama],

Kami menerima permintaan untuk mereset password akun Anda.

[Reset Password Button]

Link ini hanya berlaku selama 1 jam.
Jika Anda tidak meminta reset password, abaikan email ini.
```

### 2. Password Change Confirmation Email
```
Subject: Password Berhasil Diubah - Karir Nusantara

Halo [Nama],

✓ Password akun Anda telah berhasil diubah.

Perubahan ini dilakukan pada: Baru saja

Untuk keamanan akun:
• Semua sesi login aktif telah diakhiri
• Anda perlu login kembali dengan password baru

⚠️ Jika Anda tidak melakukan perubahan ini, segera 
hubungi tim dukungan kami.
```

---

## 🧪 Quick Testing

### Test All Password APIs
```bash
# 1. Forgot Password
curl -X POST http://localhost:8081/api/v1/auth/forgot-password \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com"}'

# 2. Reset Password (use token from email)
curl -X POST http://localhost:8081/api/v1/auth/reset-password \
  -H "Content-Type: application/json" \
  -d '{"token":"TOKEN_FROM_EMAIL","new_password":"NewPass123!"}'

# 3. Change Password (need to login first)
TOKEN=$(curl -s -X POST http://localhost:8081/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"NewPass123!"}' \
  | jq -r '.data.access_token')

curl -X PUT http://localhost:8081/api/v1/auth/change-password \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"old_password":"NewPass123!","new_password":"AnotherPass456!"}'
```

### Run Test Scripts
```bash
# Forgot/Reset Password test
./tests/password_reset_test.sh

# Change Password test
./tests/change_password_test.sh
```

---

## 🎯 When to Use Which?

### Use **Forgot Password** when:
- ❓ User doesn't remember their password
- 🔒 User is locked out of their account
- 📧 User has access to their email

### Use **Reset Password** when:
- ✉️ User received reset email with token
- 🔗 User clicked reset link from email
- 🆕 User wants to set new password (without old password)

### Use **Change Password** when:
- ✅ User is logged in
- 🔐 User knows their current password
- 🔄 User wants to update password for security
- ⏰ Regular password rotation policy

---

## 📖 Full Documentation

- **Change Password:** `/docs/CHANGE_PASSWORD_API.md`
- **Change Password Implementation:** `/docs/CHANGE_PASSWORD_IMPLEMENTATION.md`
- **Forgot/Reset Password:** `/docs/PASSWORD_RESET_EMAIL_API.md`
- **Main API Docs:** `/README.md`

---

## 🚨 Important Notes

### Email Delivery
- Emails are sent **asynchronously** (non-blocking)
- Check SMTP configuration in `.env`
- Email failures don't affect API success response
- Monitor email delivery in production

### Security Best Practices
1. Always use HTTPS in production
2. Implement rate limiting (e.g., max 5 requests per hour per email)
3. Consider adding CAPTCHA for forgot password
4. Log all password change attempts
5. Consider 2FA for additional security

### Database
- Password reset tokens stored in `password_reset_tokens` table
- Tokens auto-expire after 1 hour
- Consider adding cleanup job for expired tokens

### Frontend Integration
- Show loading states during API calls
- Handle all error cases gracefully
- Redirect to login after password change
- Clear auth tokens after password change
- Show email sent confirmation (for UX)

---

## ✅ Checklist for Production

- [ ] Configure SMTP server properly
- [ ] Set up email templates with company branding
- [ ] Add rate limiting middleware
- [ ] Implement password history (prevent reusing last N passwords)
- [ ] Add audit logging for password changes
- [ ] Set up monitoring for failed login attempts
- [ ] Configure proper CORS settings
- [ ] Add CAPTCHA for forgot password endpoint
- [ ] Test email delivery in production environment
- [ ] Set up cleanup job for expired tokens
- [ ] Review and update security policies
- [ ] Document password requirements for users

---

**Last Updated:** January 21, 2026  
**API Version:** v1  
**Status:** ✅ Production Ready
