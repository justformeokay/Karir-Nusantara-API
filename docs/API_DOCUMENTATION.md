# Karir Nusantara API Documentation

## Base URL
```
http://localhost:8081/api/v1
```

## Authentication
All protected endpoints require a Bearer token in the Authorization header:
```
Authorization: Bearer <access_token>
```

---

## Table of Contents
1. [Authentication](#authentication-endpoints)
2. [Jobs](#jobs-endpoints)
3. [CV](#cv-endpoints)
4. [Applications](#applications-endpoints)

---

## Authentication Endpoints

### Register User
```
POST /auth/register
```

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "Password123!",
  "full_name": "John Doe",
  "role": "job_seeker"
}
```

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `email` | string | ✅ | Valid email address |
| `password` | string | ✅ | Min 8 chars, must contain uppercase, lowercase, number |
| `full_name` | string | ✅ | User's full name |
| `role` | string | ✅ | Must be `job_seeker` or `company` |
| `phone` | string | ❌ | Phone number (required for company) |
| `company_name` | string | ❌ | Company name (required for company role) |

**⚠️ Important Notes:**
- Use `full_name`, NOT `name`
- Use `job_seeker` (with underscore), NOT `jobseeker`

**Response (201):**
```json
{
  "success": true,
  "message": "Registration successful",
  "data": {
    "user": {
      "id": 1,
      "email": "user@example.com",
      "role": "job_seeker",
      "full_name": "John Doe"
    },
    "access_token": "eyJhbGciOiJIUzI1NiIs...",
    "refresh_token": "uuid-refresh-token",
    "expires_in": 900
  }
}
```

---

### Login
```
POST /auth/login
```

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "Password123!"
}
```

**Response (200):**
```json
{
  "success": true,
  "message": "Login successful",
  "data": {
    "user": {
      "id": 1,
      "email": "user@example.com",
      "role": "job_seeker",
      "full_name": "John Doe"
    },
    "access_token": "eyJhbGciOiJIUzI1NiIs...",
    "refresh_token": "uuid-refresh-token",
    "expires_in": 900
  }
}
```

---

### Refresh Token
```
POST /auth/refresh
```

**Request Body:**
```json
{
  "refresh_token": "uuid-refresh-token"
}
```

---

### Get Current User
```
GET /auth/me
```
**Auth Required:** ✅

---

## Jobs Endpoints

### List Jobs (Public)
```
GET /jobs
```

**Query Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `page` | int | Page number (default: 1) |
| `per_page` | int | Items per page (default: 20, max: 100) |
| `search` | string | Search in title and description |
| `city` | string | Filter by city |
| `province` | string | Filter by province |
| `job_type` | string | Filter by type: `full_time`, `part_time`, `contract`, `internship`, `freelance` |
| `experience_level` | string | Filter by level: `entry`, `junior`, `mid`, `senior`, `lead`, `executive` |
| `is_remote` | boolean | Filter remote jobs |
| `salary_min` | int | Minimum salary |
| `salary_max` | int | Maximum salary |
| `sort_by` | string | Sort field: `created_at`, `salary_min`, `views_count` |
| `sort_order` | string | Sort order: `asc`, `desc` |

---

### Get Job by ID (Public)
```
GET /jobs/{id}
```

---

### Get Job by Slug (Public)
```
GET /jobs/slug/{slug}
```

---

### Create Job
```
POST /jobs
```
**Auth Required:** ✅ (Company only)

**Request Body:**
```json
{
  "title": "Senior Software Engineer",
  "description": "Job description here...",
  "requirements": "Requirements here...",
  "responsibilities": "Responsibilities here...",
  "benefits": "Benefits here...",
  "city": "Jakarta Selatan",
  "province": "DKI Jakarta",
  "is_remote": true,
  "job_type": "full_time",
  "experience_level": "senior",
  "salary_min": 15000000,
  "salary_max": 25000000,
  "salary_currency": "IDR",
  "is_salary_visible": true,
  "skills": ["Go", "Python", "Docker"]
}
```

| Field | Type | Required | Values |
|-------|------|----------|--------|
| `title` | string | ✅ | Job title |
| `description` | string | ✅ | Job description |
| `city` | string | ✅ | City name |
| `province` | string | ✅ | Province name |
| `job_type` | string | ✅ | `full_time`, `part_time`, `contract`, `internship`, `freelance` |
| `experience_level` | string | ✅ | `entry`, `junior`, `mid`, `senior`, `lead`, `executive` |
| `is_remote` | boolean | ❌ | Default: false |
| `requirements` | string | ❌ | Job requirements |
| `responsibilities` | string | ❌ | Job responsibilities |
| `benefits` | string | ❌ | Job benefits |
| `salary_min` | int | ❌ | Minimum salary |
| `salary_max` | int | ❌ | Maximum salary |
| `is_salary_visible` | boolean | ❌ | Show salary to applicants |
| `skills` | string[] | ❌ | Required skills |

**⚠️ Important Notes:**
- Use `job_type`, NOT `type`
- Use `city` and `province`, NOT `location`
- New jobs are created with status `draft`

---

### Update Job
```
PUT /jobs/{id}
```
**Auth Required:** ✅ (Company only - owner)

**Request Body:** Same as Create, all fields optional

---

### Delete Job
```
DELETE /jobs/{id}
```
**Auth Required:** ✅ (Company only - owner)

---

### Publish Job (Draft → Active)
```
PATCH /jobs/{id}/publish
```
**Auth Required:** ✅ (Company only - owner)

**Description:** Changes job status from `draft` to `active`, making it visible to job seekers.

---

### Pause Job (Active → Paused)
```
PATCH /jobs/{id}/pause
```
**Auth Required:** ✅ (Company only - owner)

**Description:** Temporarily hides the job from listings.

---

### Close Job (Active/Paused → Closed)
```
PATCH /jobs/{id}/close
```
**Auth Required:** ✅ (Company only - owner)

**Description:** Closes the job permanently (can be reopened).

---

### Reopen Job (Closed/Paused → Active)
```
PATCH /jobs/{id}/reopen
```
**Auth Required:** ✅ (Company only - owner)

**Description:** Reopens a closed or paused job.

---

### Job Status Flow
```
                   ┌─────────┐
                   │  draft  │
                   └────┬────┘
                        │ publish
                        ▼
     ┌─────────────►┌────────┐◄───────────────┐
     │   reopen     │ active │     reopen     │
     │              └───┬────┘                │
     │         pause/   │   \close            │
     │              ▼   │    ▼                │
     │        ┌──────┐  │  ┌────────┐         │
     │        │paused│──┼──│ closed │─────────┘
     │        └──────┘  │  └────────┘
     │                  │
     └──────────────────┘
```

---

## CV Endpoints

### Create or Update CV
```
POST /cv
```
**Auth Required:** ✅ (Job Seeker only)

**Request Body:**
```json
{
  "personal_info": {
    "full_name": "John Doe",
    "email": "john@example.com",
    "phone": "+6281234567890",
    "city": "Jakarta",
    "province": "DKI Jakarta",
    "summary": "Experienced software engineer...",
    "linkedin": "https://linkedin.com/in/johndoe",
    "portfolio": "https://johndoe.dev"
  },
  "education": [
    {
      "institution": "Universitas Indonesia",
      "degree": "S1 Teknik Informatika",
      "field_of_study": "Computer Science",
      "start_date": "2015-08-01",
      "end_date": "2019-07-01",
      "gpa": "3.8",
      "description": "Relevant coursework..."
    }
  ],
  "experience": [
    {
      "company": "PT Tech Company",
      "position": "Software Engineer",
      "location": "Jakarta",
      "start_date": "2019-08-01",
      "end_date": "",
      "is_current": true,
      "description": "Building web applications...",
      "achievements": ["Increased performance by 50%"]
    }
  ],
  "skills": [
    {
      "name": "Go",
      "level": "advanced",
      "category": "Programming Language"
    }
  ],
  "certifications": [
    {
      "name": "AWS Solutions Architect",
      "issuer": "Amazon Web Services",
      "issue_date": "2023-01-15",
      "expiry_date": "2026-01-15",
      "credential_id": "ABC123",
      "credential_url": "https://aws.amazon.com/verify/ABC123"
    }
  ],
  "languages": [
    {
      "name": "English",
      "proficiency": "fluent"
    }
  ],
  "projects": [
    {
      "name": "E-commerce Platform",
      "description": "Built a full-stack e-commerce...",
      "url": "https://github.com/johndoe/ecommerce",
      "skills": ["Go", "React", "PostgreSQL"]
    }
  ]
}
```

**⚠️ Important Notes:**
- Both `experience` AND `experiences` are supported (backward compatibility)
- Skill levels: `beginner`, `intermediate`, `advanced`, `expert`
- Language proficiency: `basic`, `conversational`, `proficient`, `fluent`, `native`

---

### Get My CV
```
GET /cv
```
**Auth Required:** ✅ (Job Seeker only)

---

### Delete CV
```
DELETE /cv
```
**Auth Required:** ✅ (Job Seeker only)

---

## Applications Endpoints

### Apply for Job
```
POST /applications
```
**Auth Required:** ✅ (Job Seeker only)

**Request Body:**
```json
{
  "job_id": 1,
  "cover_letter": "I am excited to apply for this position..."
}
```

**Note:** CV must be created before applying.

---

### Get My Applications (Job Seeker)
```
GET /applications/me
```
**Auth Required:** ✅ (Job Seeker only)

**Query Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `page` | int | Page number |
| `per_page` | int | Items per page |
| `status` | string | Filter by status |

---

### Get Company Applications (Company)
```
GET /applications/company
```
**Auth Required:** ✅ (Company only)

---

### Get Applications by Job (Company)
```
GET /jobs/{jobId}/applications
```
**Auth Required:** ✅ (Company only - job owner)

---

### Get Application by ID
```
GET /applications/{id}
```
**Auth Required:** ✅ (Owner or Company)

---

### Get Application Timeline
```
GET /applications/{id}/timeline
```
**Auth Required:** ✅ (Owner or Company)

---

### Update Application Status (Company)
```
PATCH /applications/{id}/status
```
**Auth Required:** ✅ (Company only - job owner)

**Request Body:**
```json
{
  "status": "shortlisted",
  "note": "Candidate meets initial requirements"
}
```

For interview scheduling:
```json
{
  "status": "interview_scheduled",
  "note": "Interview with HR",
  "scheduled_at": "2026-01-20T10:00:00+07:00",
  "scheduled_location": "Office Jakarta",
  "scheduled_notes": "Please bring your portfolio"
}
```

---

### Application Status Values
| Status | Label (ID) | Description |
|--------|------------|-------------|
| `submitted` | Lamaran Terkirim | Application received |
| `viewed` | Sedang Ditinjau | Company viewed application |
| `shortlisted` | Masuk Shortlist | Added to shortlist |
| `interview_scheduled` | Interview Dijadwalkan | Interview scheduled |
| `interview_completed` | Interview Selesai | Interview done |
| `assessment` | Tahap Assessment | Assessment phase |
| `offer_sent` | Penawaran Dikirim | Offer letter sent |
| `offer_accepted` | Penawaran Diterima | Candidate accepted |
| `hired` | Diterima | Candidate hired ✅ |
| `rejected` | Tidak Lolos | Application rejected |
| `withdrawn` | Dibatalkan | Candidate withdrew |

---

### Application Status Flow
```
submitted → viewed → shortlisted → interview_scheduled → interview_completed
                                                              ↓
                                         assessment ←─────────┘
                                              ↓
                                         offer_sent → offer_accepted → hired
                                              
(rejected can be applied from any non-terminal status)
```

---

### Withdraw Application (Job Seeker)
```
POST /applications/{id}/withdraw
```
**Auth Required:** ✅ (Job Seeker - application owner)

**Request Body:**
```json
{
  "reason": "Found another opportunity"
}
```

---

## Error Responses

### Validation Error (422)
```json
{
  "success": false,
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Validation failed",
    "details": {
      "email": "This field is required",
      "password": "Must be at least 8 characters"
    }
  }
}
```

### Authentication Error (401)
```json
{
  "success": false,
  "error": {
    "code": "UNAUTHORIZED",
    "message": "Invalid or expired token"
  }
}
```

### Not Found (404)
```json
{
  "success": false,
  "error": {
    "code": "NOT_FOUND",
    "message": "Job not found"
  }
}
```

### Forbidden (403)
```json
{
  "success": false,
  "error": {
    "code": "FORBIDDEN",
    "message": "You don't have permission to access this resource"
  }
}
```

---

## Rate Limiting
- 100 requests per minute per IP
- 1000 requests per hour per user

---

## Changelog

### v1.0.0 (2026-01-17)
- Initial release
- Added job status management endpoints (`/publish`, `/close`, `/pause`, `/reopen`)
- Fixed CV experience field to support both `experience` and `experiences`
