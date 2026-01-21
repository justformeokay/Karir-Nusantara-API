# Change Password API Documentation

## Overview
API endpoint untuk mengubah password user yang sedang login. Memerlukan autentikasi dengan access token dan verifikasi password lama. Setelah berhasil mengubah password, sistem akan:
- Mengirim email konfirmasi otomatis
- Membatalkan semua refresh token aktif (user harus login ulang)
- Memaksa re-authentication untuk keamanan

## Endpoint

### Change Password
**Endpoint:** `PUT /api/v1/auth/change-password`  
**Authentication:** Required (Bearer Token)  
**Content-Type:** `application/json`

---

## Request

### Headers
```
Authorization: Bearer {access_token}
Content-Type: application/json
```

### Body Parameters
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `old_password` | string | Yes | Password lama user saat ini |
| `new_password` | string | Yes | Password baru (min 8 karakter, harus mengandung huruf besar, huruf kecil, dan angka) |

### Example Request
```bash
curl -X PUT http://localhost:8081/api/v1/auth/change-password \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  -d '{
    "old_password": "OldPassword123!",
    "new_password": "NewPassword456!"
  }'
```

---

## Response

### Success Response (200 OK)
```json
{
  "success": true,
  "message": "Password berhasil diubah. Silakan login kembali dengan password baru Anda."
}
```

**Side Effects:**
1. ✅ Password di database berhasil diupdate
2. ✅ Email konfirmasi otomatis dikirim ke user
3. ✅ Semua refresh token aktif dibatalkan
4. ✅ User harus login ulang dengan password baru

### Error Responses

#### 400 Bad Request - Password Lama Salah
```json
{
  "success": false,
  "error": {
    "code": "BAD_REQUEST",
    "message": "Password lama tidak sesuai"
  }
}
```

#### 400 Bad Request - Password Sama
```json
{
  "success": false,
  "error": {
    "code": "BAD_REQUEST",
    "message": "Password baru tidak boleh sama dengan password lama"
  }
}
```

#### 401 Unauthorized - Token Tidak Ada
```json
{
  "success": false,
  "error": {
    "code": "UNAUTHORIZED",
    "message": "Missing authorization header"
  }
}
```

#### 401 Unauthorized - Token Invalid/Expired
```json
{
  "success": false,
  "error": {
    "code": "UNAUTHORIZED",
    "message": "Invalid or expired token"
  }
}
```

#### 422 Unprocessable Entity - Validation Error
```json
{
  "success": false,
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Validation failed",
    "details": {
      "new_password": [
        "Password must be at least 8 characters",
        "Password must contain uppercase, lowercase, and number"
      ]
    }
  }
}
```

---

## Email Notification

Setelah password berhasil diubah, user akan menerima email konfirmasi dengan format:

**Subject:** `Password Berhasil Diubah - Karir Nusantara`

**Content:**
- Konfirmasi perubahan password
- Timestamp perubahan
- Informasi bahwa semua sesi login aktif telah diakhiri
- Link untuk login kembali
- Peringatan keamanan (jika user tidak melakukan perubahan)

**Email Template Preview:**
```html
✓ Password Berhasil Diubah

Halo [Nama User],

Password akun Anda telah berhasil diubah.

Untuk keamanan akun Anda:
- Semua sesi login aktif telah diakhiri
- Anda perlu login kembali dengan password baru
- Pastikan password Anda tersimpan dengan aman

⚠️ PERHATIAN:
Jika Anda tidak melakukan perubahan password ini, segera hubungi 
tim dukungan kami dan reset password Anda.

[Login Sekarang]
```

---

## Security Features

### 1. Password Validation
- Minimum 8 karakter
- Harus mengandung huruf besar
- Harus mengandung huruf kecil
- Harus mengandung angka

### 2. Old Password Verification
Sistem memverifikasi password lama sebelum mengizinkan perubahan untuk mencegah unauthorized password change.

### 3. Same Password Prevention
Mencegah user menggunakan password yang sama dengan password lama.

### 4. Token Revocation
Setelah password diubah, semua refresh token aktif dibatalkan sehingga user harus login ulang di semua device.

### 5. Email Confirmation
User menerima notifikasi email segera setelah password diubah untuk mendeteksi unauthorized changes.

### 6. Authentication Required
Endpoint ini hanya bisa diakses oleh user yang sudah login dengan valid access token.

---

## Use Cases

### Use Case 1: Regular Password Change
User ingin mengubah password secara berkala untuk keamanan:
1. User login dengan password lama
2. User memanggil API change password dengan old_password dan new_password
3. Sistem validasi old_password
4. Password diupdate di database
5. Email konfirmasi dikirim
6. User logout dan login kembali dengan password baru

### Use Case 2: Security Breach Response
User mencurigai akun mereka telah diakses orang lain:
1. User segera login (jika masih bisa)
2. User langsung change password
3. Semua sesi aktif dibatalkan
4. User menerima email konfirmasi
5. Unauthorized user kehilangan akses

---

## Testing

### Test Scenario 1: Successful Password Change
```bash
# 1. Login first
LOGIN_RESPONSE=$(curl -s -X POST http://localhost:8081/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "OldPassword123!"
  }')

# Extract access token
TOKEN=$(echo $LOGIN_RESPONSE | jq -r '.data.access_token')

# 2. Change password
curl -X PUT http://localhost:8081/api/v1/auth/change-password \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "old_password": "OldPassword123!",
    "new_password": "NewPassword456!"
  }'

# Expected: {"success": true, "message": "Password berhasil diubah..."}

# 3. Verify old password doesn't work
curl -X POST http://localhost:8081/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "OldPassword123!"
  }'

# Expected: {"success": false, "error": {"code": "INVALID_CREDENTIALS"}}

# 4. Verify new password works
curl -X POST http://localhost:8081/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "NewPassword456!"
  }'

# Expected: {"success": true, "data": {...}}
```

### Test Scenario 2: Wrong Old Password
```bash
curl -X PUT http://localhost:8081/api/v1/auth/change-password \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "old_password": "WrongPassword!",
    "new_password": "NewPassword456!"
  }'

# Expected: {"success": false, "error": {"message": "Password lama tidak sesuai"}}
```

### Test Scenario 3: Same Password
```bash
curl -X PUT http://localhost:8081/api/v1/auth/change-password \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "old_password": "CurrentPassword123!",
    "new_password": "CurrentPassword123!"
  }'

# Expected: {"success": false, "error": {"message": "Password baru tidak boleh sama..."}}
```

### Test Scenario 4: Unauthorized Access
```bash
curl -X PUT http://localhost:8081/api/v1/auth/change-password \
  -H "Content-Type: application/json" \
  -d '{
    "old_password": "OldPassword123!",
    "new_password": "NewPassword456!"
  }'

# Expected: {"success": false, "error": {"code": "UNAUTHORIZED"}}
```

---

## Integration with Frontend

### React Example
```typescript
import { authApi } from '@/api/auth';

const ChangePasswordForm = () => {
  const [formData, setFormData] = useState({
    oldPassword: '',
    newPassword: '',
    confirmPassword: ''
  });

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    
    if (formData.newPassword !== formData.confirmPassword) {
      toast.error('Password baru dan konfirmasi tidak cocok');
      return;
    }

    try {
      await authApi.changePassword({
        old_password: formData.oldPassword,
        new_password: formData.newPassword
      });
      
      toast.success('Password berhasil diubah. Silakan login kembali.');
      
      // Logout and redirect to login
      localStorage.removeItem('access_token');
      window.location.href = '/login';
      
    } catch (error: any) {
      const errorMsg = error.response?.data?.error?.message 
        || 'Gagal mengubah password';
      toast.error(errorMsg);
    }
  };

  // ... form JSX
};
```

---

## Notes

1. **Email Delivery**: Email konfirmasi dikirim secara asynchronous untuk tidak memblokir response API. Jika email gagal terkirim, request tetap success (password tetap berhasil diubah).

2. **Token Expiry**: Access token yang digunakan untuk change password akan tetap valid sampai expiry time-nya, tapi refresh token sudah dibatalkan. Best practice: logout setelah change password.

3. **Password Complexity**: Validasi dilakukan di backend. Frontend sebaiknya juga melakukan validasi yang sama untuk UX yang lebih baik.

4. **Rate Limiting**: Pertimbangkan menambahkan rate limiting untuk endpoint ini untuk mencegah brute force attacks.

5. **Password History**: Untuk keamanan lebih lanjut, pertimbangkan menyimpan hash dari N password terakhir dan mencegah reuse.

---

## Related Endpoints

- [POST /api/v1/auth/forgot-password](PASSWORD_RESET_EMAIL_API.md) - For password reset via email
- [POST /api/v1/auth/reset-password](PASSWORD_RESET_EMAIL_API.md) - Complete password reset flow
- [POST /api/v1/auth/login](../README.md) - Login endpoint
- [POST /api/v1/auth/logout](../README.md) - Logout endpoint
