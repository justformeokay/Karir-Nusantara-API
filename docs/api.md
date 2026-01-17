# Karir Nusantara API Documentation

## Base URL

```
Development: http://localhost:8080/api/v1
Production: https://api.karirnusantara.com/api/v1
```

## Authentication

The API uses JWT (JSON Web Token) for authentication. Include the access token in the `Authorization` header:

```
Authorization: Bearer <access_token>
```

### Token Lifecycle
- **Access Token**: Valid for 15 minutes
- **Refresh Token**: Valid for 7 days (stored in HTTP-only cookie)

---

## Endpoints

### Health Check

#### GET /health
Check API health status.

**Response:**
```json
{
  "status": "healthy",
  "service": "karir-nusantara-api",
  "version": "1.0.0"
}
```

---

## Auth Module

### POST /api/v1/auth/register
Register a new user (job seeker or company).

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "SecurePass123!",
  "name": "John Doe",
  "role": "job_seeker",
  // For company registration only:
  "company_name": "PT Example",
  "company_description": "A great company",
  "company_industry": "Technology",
  "company_size": "50-100",
  "company_location": "Jakarta"
}
```

**Validation Rules:**
- `email`: Required, valid email format
- `password`: Required, min 8 chars, must contain uppercase, lowercase, number
- `name`: Required, min 2 chars
- `role`: Required, one of: `job_seeker`, `company`

**Response (201 Created):**
```json
{
  "success": true,
  "data": {
    "user": {
      "id": 1,
      "email": "user@example.com",
      "name": "John Doe",
      "role": "job_seeker",
      "created_at": "2024-01-15T10:00:00Z"
    },
    "access_token": "eyJhbGciOiJIUzI1NiIs...",
    "expires_in": 900
  }
}
```

### POST /api/v1/auth/login
Login with email and password.

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "SecurePass123!"
}
```

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "user": {
      "id": 1,
      "email": "user@example.com",
      "name": "John Doe",
      "role": "job_seeker"
    },
    "access_token": "eyJhbGciOiJIUzI1NiIs...",
    "expires_in": 900
  }
}
```

### POST /api/v1/auth/refresh
Refresh access token using refresh token (from cookie).

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIs...",
    "expires_in": 900
  }
}
```

### POST /api/v1/auth/logout
Logout and invalidate refresh token.

**Headers:** `Authorization: Bearer <access_token>`

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Logged out successfully"
}
```

### GET /api/v1/auth/me
Get current authenticated user.

**Headers:** `Authorization: Bearer <access_token>`

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "id": 1,
    "email": "user@example.com",
    "name": "John Doe",
    "role": "job_seeker",
    "avatar_url": "https://...",
    "created_at": "2024-01-15T10:00:00Z"
  }
}
```

---

## Jobs Module

### GET /api/v1/jobs
List all active jobs (public).

**Query Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| page | int | Page number (default: 1) |
| limit | int | Items per page (default: 10, max: 100) |
| search | string | Search in title, description |
| location | string | Filter by location |
| type | string | Filter by type: full_time, part_time, contract, internship |
| experience | string | Filter by level: entry, mid, senior, lead |
| salary_min | int | Minimum salary filter |
| salary_max | int | Maximum salary filter |
| sort_by | string | Sort field: created_at, salary_min, title |
| sort_order | string | asc or desc |

**Response (200 OK):**
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "title": "Software Engineer",
      "slug": "software-engineer-1",
      "description": "We are looking for...",
      "requirements": "- 3+ years experience...",
      "type": "full_time",
      "experience_level": "mid",
      "location": "Jakarta",
      "is_remote": false,
      "salary_min": 15000000,
      "salary_max": 25000000,
      "salary_currency": "IDR",
      "skills": ["Go", "React", "PostgreSQL"],
      "company": {
        "id": 2,
        "name": "PT Tech Corp",
        "logo_url": "https://..."
      },
      "status": "active",
      "created_at": "2024-01-15T10:00:00Z"
    }
  ],
  "meta": {
    "total": 50,
    "page": 1,
    "limit": 10,
    "total_pages": 5
  }
}
```

### GET /api/v1/jobs/{slug}
Get job details by slug (public).

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "id": 1,
    "title": "Software Engineer",
    "slug": "software-engineer-1",
    "description": "Full job description...",
    "requirements": "Detailed requirements...",
    "benefits": "What we offer...",
    "type": "full_time",
    "experience_level": "mid",
    "location": "Jakarta",
    "is_remote": false,
    "salary_min": 15000000,
    "salary_max": 25000000,
    "salary_currency": "IDR",
    "skills": ["Go", "React", "PostgreSQL"],
    "company": {
      "id": 2,
      "name": "PT Tech Corp",
      "description": "Leading tech company...",
      "industry": "Technology",
      "size": "100-500",
      "location": "Jakarta",
      "logo_url": "https://..."
    },
    "status": "active",
    "application_count": 25,
    "created_at": "2024-01-15T10:00:00Z",
    "expires_at": "2024-02-15T10:00:00Z"
  }
}
```

### POST /api/v1/jobs
Create a new job posting (company only).

**Headers:** `Authorization: Bearer <access_token>`

**Request Body:**
```json
{
  "title": "Software Engineer",
  "description": "We are looking for a talented...",
  "requirements": "- 3+ years experience in Go\n- Experience with React",
  "benefits": "- Competitive salary\n- Remote work options",
  "type": "full_time",
  "experience_level": "mid",
  "location": "Jakarta",
  "is_remote": false,
  "salary_min": 15000000,
  "salary_max": 25000000,
  "salary_currency": "IDR",
  "skills": ["Go", "React", "PostgreSQL"],
  "expires_at": "2024-02-15T10:00:00Z"
}
```

**Response (201 Created):**
```json
{
  "success": true,
  "data": {
    "id": 1,
    "title": "Software Engineer",
    "slug": "software-engineer-1",
    ...
  }
}
```

### PUT /api/v1/jobs/{id}
Update a job posting (company only, own jobs).

**Headers:** `Authorization: Bearer <access_token>`

**Request Body:** Same as POST

**Response (200 OK):** Same as GET by slug

### DELETE /api/v1/jobs/{id}
Delete a job posting (company only, own jobs).

**Headers:** `Authorization: Bearer <access_token>`

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Job deleted successfully"
}
```

### GET /api/v1/jobs/company
List jobs for the authenticated company.

**Headers:** `Authorization: Bearer <access_token>`

**Query Parameters:** Same as GET /jobs

**Response:** Same as GET /jobs

---

## CVs Module

### GET /api/v1/cvs/me
Get the current user's CV (job seeker only).

**Headers:** `Authorization: Bearer <access_token>`

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "id": 1,
    "user_id": 1,
    "personal_info": {
      "full_name": "John Doe",
      "email": "john@example.com",
      "phone": "+62812345678",
      "address": "Jakarta, Indonesia",
      "date_of_birth": "1995-01-15",
      "gender": "male",
      "summary": "Experienced software developer..."
    },
    "education": [
      {
        "institution": "Universitas Indonesia",
        "degree": "S1",
        "field_of_study": "Computer Science",
        "start_date": "2013-09-01",
        "end_date": "2017-07-01",
        "gpa": 3.75,
        "description": "Focus on software engineering"
      }
    ],
    "experience": [
      {
        "company": "PT Tech Corp",
        "position": "Software Engineer",
        "location": "Jakarta",
        "start_date": "2020-01-01",
        "end_date": null,
        "is_current": true,
        "description": "Developing backend services..."
      }
    ],
    "skills": [
      {
        "name": "Go",
        "level": "advanced",
        "years": 3
      },
      {
        "name": "React",
        "level": "intermediate",
        "years": 2
      }
    ],
    "certifications": [
      {
        "name": "AWS Certified Developer",
        "issuer": "Amazon Web Services",
        "issue_date": "2023-06-01",
        "expiry_date": "2026-06-01",
        "credential_id": "ABC123",
        "credential_url": "https://..."
      }
    ],
    "languages": [
      {
        "language": "Indonesian",
        "proficiency": "native"
      },
      {
        "language": "English",
        "proficiency": "professional"
      }
    ],
    "completeness": 85,
    "created_at": "2024-01-15T10:00:00Z",
    "updated_at": "2024-01-20T15:30:00Z"
  }
}
```

### POST /api/v1/cvs
Create or update CV (job seeker only).

**Headers:** `Authorization: Bearer <access_token>`

**Request Body:**
```json
{
  "personal_info": {
    "full_name": "John Doe",
    "email": "john@example.com",
    "phone": "+62812345678",
    ...
  },
  "education": [...],
  "experience": [...],
  "skills": [...],
  "certifications": [...],
  "languages": [...]
}
```

**Response (200 OK or 201 Created):** Same as GET /cvs/me

---

## Applications Module

### POST /api/v1/applications
Apply for a job (job seeker only).

**Headers:** `Authorization: Bearer <access_token>`

**Request Body:**
```json
{
  "job_id": 1,
  "cover_letter": "I am excited to apply for this position..."
}
```

**Response (201 Created):**
```json
{
  "success": true,
  "data": {
    "id": 1,
    "job": {
      "id": 1,
      "title": "Software Engineer",
      "company": {
        "id": 2,
        "name": "PT Tech Corp"
      }
    },
    "cover_letter": "I am excited to apply...",
    "current_status": "submitted",
    "timeline": [
      {
        "id": 1,
        "status": "submitted",
        "note": "Lamaran berhasil dikirim",
        "created_at": "2024-01-15T10:00:00Z"
      }
    ],
    "created_at": "2024-01-15T10:00:00Z"
  }
}
```

### GET /api/v1/applications/me
List my applications (job seeker only).

**Headers:** `Authorization: Bearer <access_token>`

**Query Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| page | int | Page number (default: 1) |
| limit | int | Items per page (default: 10) |
| status | string | Filter by status |

**Response (200 OK):**
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "job": {
        "id": 1,
        "title": "Software Engineer",
        "company": {...}
      },
      "current_status": "under_review",
      "timeline": [...],
      "created_at": "2024-01-15T10:00:00Z"
    }
  ],
  "meta": {
    "total": 5,
    "page": 1,
    "limit": 10,
    "total_pages": 1
  }
}
```

### GET /api/v1/applications/company
List applications for company (company only).

**Headers:** `Authorization: Bearer <access_token>`

**Query Parameters:** Same as /applications/me

**Response:** Similar structure with applicant info included

### GET /api/v1/jobs/{jobId}/applications
List applications for a specific job (company only).

**Headers:** `Authorization: Bearer <access_token>`

**Query Parameters:** Same as /applications/me

### GET /api/v1/applications/{id}
Get application details.

**Headers:** `Authorization: Bearer <access_token>`

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "id": 1,
    "job": {...},
    "applicant": {...}, // Only visible to company
    "cv_snapshot": {...},
    "cover_letter": "...",
    "current_status": "interview_scheduled",
    "timeline": [
      {
        "id": 1,
        "status": "submitted",
        "note": "Lamaran berhasil dikirim",
        "created_at": "2024-01-15T10:00:00Z"
      },
      {
        "id": 2,
        "status": "under_review",
        "note": "CV sedang ditinjau oleh tim HR",
        "created_at": "2024-01-16T09:00:00Z"
      },
      {
        "id": 3,
        "status": "interview_scheduled",
        "note": "Anda diundang untuk interview",
        "scheduled_at": "2024-01-20T10:00:00Z",
        "scheduled_location": "Zoom Meeting",
        "scheduled_notes": "Link akan dikirim via email",
        "created_at": "2024-01-17T14:00:00Z"
      }
    ],
    "created_at": "2024-01-15T10:00:00Z",
    "updated_at": "2024-01-17T14:00:00Z"
  }
}
```

### GET /api/v1/applications/{id}/timeline
Get application timeline only.

**Headers:** `Authorization: Bearer <access_token>`

**Response (200 OK):**
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "status": "submitted",
      "note": "Lamaran berhasil dikirim",
      "created_at": "2024-01-15T10:00:00Z"
    },
    ...
  ]
}
```

### PATCH /api/v1/applications/{id}/status
Update application status (company only).

**Headers:** `Authorization: Bearer <access_token>`

**Request Body:**
```json
{
  "status": "interview_scheduled",
  "note": "Anda diundang untuk interview tahap 1",
  "scheduled_at": "2024-01-20T10:00:00+07:00",
  "scheduled_location": "Kantor Jakarta atau Zoom",
  "scheduled_notes": "Silakan konfirmasi kehadiran"
}
```

**Application Status Values:**
| Status | Description |
|--------|-------------|
| submitted | Lamaran baru dikirim |
| under_review | CV sedang ditinjau |
| shortlisted | Masuk daftar pendek |
| interview_scheduled | Interview dijadwalkan |
| interview_completed | Interview selesai |
| assessment | Tahap assessment |
| offered | Penawaran dikirim |
| hired | Diterima bekerja |
| rejected | Ditolak |
| withdrawn | Dibatalkan pelamar |

**Valid Status Transitions:**
- `submitted` → `under_review`, `rejected`, `withdrawn`
- `under_review` → `shortlisted`, `rejected`, `withdrawn`
- `shortlisted` → `interview_scheduled`, `rejected`, `withdrawn`
- `interview_scheduled` → `interview_completed`, `rejected`, `withdrawn`
- `interview_completed` → `assessment`, `offered`, `rejected`, `withdrawn`
- `assessment` → `offered`, `rejected`, `withdrawn`
- `offered` → `hired`, `rejected`, `withdrawn`

**Response (200 OK):** Same as GET /applications/{id}

### POST /api/v1/applications/{id}/withdraw
Withdraw application (job seeker only).

**Headers:** `Authorization: Bearer <access_token>`

**Request Body (optional):**
```json
{
  "reason": "Mendapat penawaran lain"
}
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Application withdrawn successfully"
}
```

---

## Error Responses

All errors follow a consistent format:

```json
{
  "success": false,
  "error": {
    "code": "error_code",
    "message": "Human readable message"
  }
}
```

### HTTP Status Codes

| Status | Description |
|--------|-------------|
| 200 | Success |
| 201 | Created |
| 400 | Bad Request - Invalid input |
| 401 | Unauthorized - Missing/invalid token |
| 403 | Forbidden - Insufficient permissions |
| 404 | Not Found |
| 409 | Conflict - Duplicate resource |
| 422 | Validation Error |
| 500 | Internal Server Error |

### Error Codes

| Code | Description |
|------|-------------|
| bad_request | Invalid request format or parameters |
| unauthorized | Authentication required or token invalid |
| forbidden | User doesn't have permission |
| not_found | Requested resource not found |
| conflict | Resource already exists |
| validation_error | Input validation failed |
| internal_error | Server error |

---

## Rate Limiting

- **Anonymous requests**: 60 requests per minute
- **Authenticated requests**: 300 requests per minute

Rate limit headers:
```
X-RateLimit-Limit: 300
X-RateLimit-Remaining: 299
X-RateLimit-Reset: 1705312800
```

---

## Pagination

All list endpoints support pagination with these query parameters:

| Parameter | Default | Max |
|-----------|---------|-----|
| page | 1 | - |
| limit | 10 | 100 |

Response includes `meta` object:
```json
{
  "meta": {
    "total": 100,
    "page": 1,
    "limit": 10,
    "total_pages": 10
  }
}
```
