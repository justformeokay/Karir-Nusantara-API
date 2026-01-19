# Admin API Documentation

## Overview

The Super Admin API provides endpoints for managing all aspects of the Karir Nusantara job portal platform. All protected endpoints require a valid admin JWT token.

## Base URL

```
/api/v1/admin
```

## Authentication

### Admin Login

Authenticate an admin user and receive a JWT token.

**Endpoint:** `POST /api/v1/admin/auth/login`

**Request Body:**
```json
{
  "email": "admin@karirnusantara.com",
  "password": "admin123"
}
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Login berhasil",
  "data": {
    "admin": {
      "id": 1,
      "email": "admin@karirnusantara.com",
      "full_name": "Super Admin",
      "role": "admin",
      "is_active": true,
      "created_at": "2024-01-01T00:00:00Z"
    },
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "expires_in": 86400
  }
}
```

### Get Current Admin

Get the currently authenticated admin's information.

**Endpoint:** `GET /api/v1/admin/auth/me`

**Headers:**
```
Authorization: Bearer <access_token>
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Data admin berhasil diambil",
  "data": {
    "id": 1,
    "email": "admin@karirnusantara.com",
    "full_name": "Super Admin",
    "role": "admin",
    "is_active": true,
    "created_at": "2024-01-01T00:00:00Z"
  }
}
```

---

## Dashboard

### Get Dashboard Statistics

Retrieve overall platform statistics.

**Endpoint:** `GET /api/v1/admin/dashboard/stats`

**Headers:**
```
Authorization: Bearer <access_token>
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Statistik dashboard berhasil diambil",
  "data": {
    "total_companies": 150,
    "pending_verifications": 12,
    "verified_companies": 130,
    "suspended_companies": 8,
    "total_jobs": 500,
    "active_jobs": 320,
    "pending_jobs": 45,
    "flagged_jobs": 5,
    "total_job_seekers": 5000,
    "active_job_seekers": 4500,
    "total_payments": 200,
    "pending_payments": 15,
    "total_revenue": 50000000
  }
}
```

---

## Company Management

### List Companies

Get a paginated list of companies with optional filtering.

**Endpoint:** `GET /api/v1/admin/companies`

**Headers:**
```
Authorization: Bearer <access_token>
```

**Query Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| status | string | Filter by status: `pending`, `verified`, `rejected`, `suspended` |
| search | string | Search by company name, email, or full name |
| page | integer | Page number (default: 1) |
| page_size | integer | Items per page (default: 10) |

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Daftar perusahaan berhasil diambil",
  "data": [
    {
      "id": 1,
      "email": "company@example.com",
      "full_name": "John Doe",
      "phone": "08123456789",
      "company_name": "PT Example",
      "company_description": "A great company",
      "company_website": "https://example.com",
      "company_logo_url": "https://example.com/logo.png",
      "company_status": "verified",
      "is_active": true,
      "is_verified": true,
      "email_verified_at": "2024-01-01T00:00:00Z",
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z",
      "jobs_count": 10,
      "active_jobs_count": 5,
      "total_applications": 100
    }
  ],
  "meta": {
    "page": 1,
    "per_page": 10,
    "total_items": 150,
    "total_pages": 15
  }
}
```

### Get Company Details

Get detailed information about a specific company.

**Endpoint:** `GET /api/v1/admin/companies/{id}`

**Headers:**
```
Authorization: Bearer <access_token>
```

### Verify Company

Approve or reject a pending company registration.

**Endpoint:** `POST /api/v1/admin/companies/{id}/verify`

**Headers:**
```
Authorization: Bearer <access_token>
```

**Request Body:**
```json
{
  "action": "approve",  // or "reject"
  "reason": "Company documents verified successfully"  // required for reject
}
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Perusahaan berhasil diverifikasi",
  "data": null
}
```

### Update Company Status

Suspend or reactivate a company account.

**Endpoint:** `PATCH /api/v1/admin/companies/{id}/status`

**Headers:**
```
Authorization: Bearer <access_token>
```

**Request Body:**
```json
{
  "action": "suspend",  // or "reactivate"
  "reason": "Violation of terms of service"  // required for suspend
}
```

---

## Job Management

### List Jobs

Get a paginated list of job postings with optional filtering.

**Endpoint:** `GET /api/v1/admin/jobs`

**Headers:**
```
Authorization: Bearer <access_token>
```

**Query Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| company_id | integer | Filter by company ID |
| status | string | Filter by status: `draft`, `pending`, `active`, `paused`, `closed`, `filled`, `rejected` |
| date_from | string | Filter jobs created from date (YYYY-MM-DD) |
| date_to | string | Filter jobs created to date (YYYY-MM-DD) |
| search | string | Search by job title or company name |
| page | integer | Page number (default: 1) |
| page_size | integer | Items per page (default: 10) |

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Daftar lowongan berhasil diambil",
  "data": [
    {
      "id": 1,
      "company_id": 1,
      "company_name": "PT Example",
      "title": "Software Engineer",
      "slug": "software-engineer-pt-example",
      "description": "Job description here...",
      "requirements": "Job requirements here...",
      "city": "Jakarta",
      "province": "DKI Jakarta",
      "is_remote": false,
      "job_type": "full_time",
      "experience_level": "mid",
      "salary_min": 10000000,
      "salary_max": 15000000,
      "status": "active",
      "admin_status": "approved",
      "admin_note": "",
      "views_count": 100,
      "applications_count": 25,
      "published_at": "2024-01-01T00:00:00Z",
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    }
  ],
  "meta": {
    "page": 1,
    "per_page": 10,
    "total_items": 500,
    "total_pages": 50
  }
}
```

### Get Job Details

Get detailed information about a specific job posting.

**Endpoint:** `GET /api/v1/admin/jobs/{id}`

### Moderate Job

Approve, reject, close, flag, or unflag a job posting.

**Endpoint:** `POST /api/v1/admin/jobs/{id}/moderate`

**Headers:**
```
Authorization: Bearer <access_token>
```

**Request Body:**
```json
{
  "action": "approve",  // "approve", "reject", "close", "flag", "unflag"
  "reason": "Content is appropriate"  // optional, required for reject/flag
}
```

---

## Payment Management

### List Payments

Get a paginated list of payments with optional filtering.

**Endpoint:** `GET /api/v1/admin/payments`

**Headers:**
```
Authorization: Bearer <access_token>
```

**Query Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| company_id | integer | Filter by company ID |
| status | string | Filter by status: `pending`, `confirmed`, `rejected` |
| date_from | string | Filter payments from date (YYYY-MM-DD) |
| date_to | string | Filter payments to date (YYYY-MM-DD) |
| page | integer | Page number (default: 1) |
| page_size | integer | Items per page (default: 10) |

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Daftar pembayaran berhasil diambil",
  "data": [
    {
      "id": 1,
      "company_id": 1,
      "company_name": "PT Example",
      "job_id": 5,
      "job_title": "Software Engineer",
      "amount": 500000,
      "proof_image_url": "https://storage.example.com/proof.jpg",
      "status": "pending",
      "status_label": "Menunggu Konfirmasi",
      "note": "",
      "submitted_at": "2024-01-15T10:30:00Z"
    }
  ],
  "meta": {
    "page": 1,
    "per_page": 10,
    "total_items": 200,
    "total_pages": 20
  }
}
```

### Get Payment Details

Get detailed information about a specific payment.

**Endpoint:** `GET /api/v1/admin/payments/{id}`

### Process Payment

Approve or reject a pending payment.

**Endpoint:** `POST /api/v1/admin/payments/{id}/process`

**Headers:**
```
Authorization: Bearer <access_token>
```

**Request Body:**
```json
{
  "action": "approve",  // or "reject"
  "note": "Payment verified successfully"  // required for reject
}
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Pembayaran berhasil diproses",
  "data": null
}
```

---

## Job Seeker Management

### List Job Seekers

Get a paginated list of job seekers with optional filtering.

**Endpoint:** `GET /api/v1/admin/job-seekers`

**Headers:**
```
Authorization: Bearer <access_token>
```

**Query Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| status | string | Filter by status: `active`, `inactive` |
| search | string | Search by name or email |
| page | integer | Page number (default: 1) |
| page_size | integer | Items per page (default: 10) |

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Daftar pencari kerja berhasil diambil",
  "data": [
    {
      "id": 1,
      "email": "jobseeker@example.com",
      "full_name": "Jane Doe",
      "phone": "08123456789",
      "avatar_url": "https://example.com/avatar.jpg",
      "is_active": true,
      "is_verified": true,
      "email_verified_at": "2024-01-01T00:00:00Z",
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z",
      "applications_count": 15,
      "has_cv": true
    }
  ],
  "meta": {
    "page": 1,
    "per_page": 10,
    "total_items": 5000,
    "total_pages": 500
  }
}
```

### Get Job Seeker Details

Get detailed information about a specific job seeker.

**Endpoint:** `GET /api/v1/admin/job-seekers/{id}`

### Update Job Seeker Status

Suspend, deactivate, or reactivate a job seeker account.

**Endpoint:** `PATCH /api/v1/admin/job-seekers/{id}/status`

**Headers:**
```
Authorization: Bearer <access_token>
```

**Request Body:**
```json
{
  "action": "suspend",  // "suspend", "deactivate", "reactivate"
  "reason": "Account flagged for suspicious activity"  // required for suspend
}
```

---

## Error Responses

All API endpoints return consistent error responses:

**400 Bad Request:**
```json
{
  "success": false,
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Action wajib diisi"
  }
}
```

**401 Unauthorized:**
```json
{
  "success": false,
  "error": {
    "code": "UNAUTHORIZED",
    "message": "Tidak terautentikasi"
  }
}
```

**403 Forbidden:**
```json
{
  "success": false,
  "error": {
    "code": "FORBIDDEN",
    "message": "Akses ditolak"
  }
}
```

**404 Not Found:**
```json
{
  "success": false,
  "error": {
    "code": "NOT_FOUND",
    "message": "Perusahaan tidak ditemukan"
  }
}
```

**500 Internal Server Error:**
```json
{
  "success": false,
  "error": {
    "code": "INTERNAL_ERROR",
    "message": "Terjadi kesalahan internal"
  }
}
```

---

## Database Migration

Before using the admin APIs, run the migration to add required database fields:

```sql
-- migrations/003_add_admin_fields.sql
```

This migration adds:
- `admin_status`, `admin_note`, `flag_reason` columns to `jobs` table
- `confirmed_by_id` column to `payments` table
- Creates necessary indexes for performance
- Inserts a default admin user (admin@karirnusantara.com / admin123)

---

## Default Admin Credentials

For development purposes, a default admin account is created:

- **Email:** admin@karirnusantara.com
- **Password:** admin123

⚠️ **Important:** Change these credentials in production!
