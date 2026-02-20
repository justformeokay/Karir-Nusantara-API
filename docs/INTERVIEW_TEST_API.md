# Interview Test API Documentation

API untuk mengelola Interview Test/Psychometric Test di Karir Nusantara Admin Dashboard.

## Base URL
```
http://localhost:8080/api/v1/admin/interview-tests
```

## Authentication
Semua endpoint memerlukan authentication sebagai admin. Sertakan token JWT di header:
```
Authorization: Bearer <admin_token>
```

## Endpoints

### 1. Get All Interview Tests
Mengambil daftar semua tes interview.

**Endpoint:** `GET /admin/interview-tests`

**Query Parameters:**
- `status` (optional): Filter berdasarkan status (`draft`, `active`, `archived`)

**Response:**
```json
{
  "status": "success",
  "message": "Interview tests retrieved successfully",
  "data": [
    {
      "id": 1,
      "title": "Tes Psikologi Dasar",
      "description": "Tes untuk mengukur kemampuan dasar psikologi kandidat",
      "duration_minutes": 60,
      "total_points": 100,
      "passing_score": 70,
      "shuffle_questions": false,
      "show_results_immediately": false,
      "status": "active",
      "created_by": 1,
      "created_at": "2026-02-19T10:00:00Z",
      "updated_at": "2026-02-19T10:00:00Z"
    }
  ]
}
```

### 2. Get Interview Test by ID
Mengambil detail tes interview termasuk pertanyaan dan opsinya.

**Endpoint:** `GET /admin/interview-tests/:id`

**Response:**
```json
{
  "status": "success",
  "message": "Interview test retrieved successfully",
  "data": {
    "id": 1,
    "title": "Tes Psikologi Dasar",
    "description": "Tes untuk mengukur kemampuan dasar psikologi kandidat",
    "duration_minutes": 60,
    "total_points": 100,
    "passing_score": 70,
    "shuffle_questions": false,
    "show_results_immediately": false,
    "status": "draft",
    "created_by": 1,
    "created_at": "2026-02-19T10:00:00Z",
    "updated_at": "2026-02-19T10:00:00Z",
    "questions": [
      {
        "id": 1,
        "question_text": "Apa itu psikologi?",
        "question_type": "multiple_choice",
        "points": 10,
        "difficulty": "easy",
        "order": 1,
        "explanation": "Psikologi adalah ilmu yang mempelajari perilaku dan proses mental",
        "options": [
          {
            "id": 1,
            "option_text": "Ilmu tentang perilaku",
            "is_correct": true,
            "order": 1
          },
          {
            "id": 2,
            "option_text": "Ilmu tentang matematika",
            "is_correct": false,
            "order": 2
          }
        ]
      },
      {
        "id": 2,
        "question_text": "Jelaskan pengalaman kerja Anda",
        "question_type": "essay",
        "points": 20,
        "difficulty": "medium",
        "order": 2,
        "explanation": null,
        "options": []
      }
    ]
  }
}
```

### 3. Create Interview Test
Membuat tes interview baru dengan pertanyaan.

**Endpoint:** `POST /admin/interview-tests`

**Request Body:**
```json
{
  "title": "Tes Psikologi Dasar",
  "description": "Tes untuk mengukur kemampuan dasar psikologi kandidat",
  "duration_minutes": 60,
  "passing_score": 70,
  "shuffle_questions": false,
  "show_results_immediately": false,
  "questions": [
    {
      "question_text": "Apa itu psikologi?",
      "question_type": "multiple_choice",
      "points": 10,
      "difficulty": "easy",
      "explanation": "Penjelasan jawaban",
      "options": [
        {
          "option_text": "Ilmu tentang perilaku",
          "is_correct": true
        },
        {
          "option_text": "Ilmu tentang matematika",
          "is_correct": false
        },
        {
          "option_text": "Ilmu tentang fisika",
          "is_correct": false
        }
      ]
    },
    {
      "question_text": "Jelaskan motivasi Anda melamar pekerjaan ini",
      "question_type": "essay",
      "points": 20,
      "difficulty": "medium",
      "explanation": null
    }
  ]
}
```

**Response:**
```json
{
  "status": "success",
  "message": "Interview test created successfully",
  "data": {
    "id": 1,
    "title": "Tes Psikologi Dasar",
    "...": "..."
  }
}
```

**Validation Rules:**
- `title`: Required
- `description`: Required
- `duration_minutes`: Required, must be > 0
- `passing_score`: Required, must be 0-100
- `questions`: Required, must have at least 1 question
- For `multiple_choice` questions:
  - Must have at least 2 options
  - Must have at least 1 correct answer
- For `essay` questions:
  - No options needed

### 4. Update Interview Test
Mengupdate tes interview yang sudah ada.

**Endpoint:** `PUT /admin/interview-tests/:id`

**Request Body:** (sama seperti Create)

**Response:**
```json
{
  "status": "success",
  "message": "Interview test updated successfully",
  "data": {
    "id": 1,
    "...": "..."
  }
}
```

### 5. Delete Interview Test
Menghapus tes interview (soft delete).

**Endpoint:** `DELETE /admin/interview-tests/:id`

**Response:**
```json
{
  "status": "success",
  "message": "Interview test deleted successfully",
  "data": null
}
```

### 6. Publish Interview Test
Mengubah status tes dari draft menjadi active.

**Endpoint:** `POST /admin/interview-tests/:id/publish`

**Response:**
```json
{
  "status": "success",
  "message": "Interview test published successfully",
  "data": {
    "id": 1,
    "status": "active",
    "...": "..."
  }
}
```

**Notes:**
- Hanya tes dengan status `draft` yang bisa dipublish
- Tes harus memiliki minimal 1 pertanyaan

### 7. Duplicate Interview Test
Membuat salinan dari tes yang sudah ada.

**Endpoint:** `POST /admin/interview-tests/:id/duplicate`

**Response:**
```json
{
  "status": "success",
  "message": "Interview test duplicated successfully",
  "data": {
    "id": 2,
    "title": "Tes Psikologi Dasar (Copy)",
    "status": "draft",
    "...": "..."
  }
}
```

**Notes:**
- Semua pertanyaan dan opsi akan disalin
- Status tes baru akan menjadi `draft`
- Title akan ditambahkan suffix " (Copy)"

## Error Responses

### Validation Error
```json
{
  "status": "error",
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Title is required"
  }
}
```

### Not Found
```json
{
  "status": "error",
  "error": {
    "code": "NOT_FOUND",
    "message": "Interview test not found"
  }
}
```

### Internal Error
```json
{
  "status": "error",
  "error": {
    "code": "INTERNAL_ERROR",
    "message": "Error message here"
  }
}
```

## Data Types

### Status
- `draft`: Tes masih dalam tahap pembuatan
- `active`: Tes sudah dipublish dan bisa digunakan
- `archived`: Tes sudah tidak aktif

### Question Type
- `multiple_choice`: Pertanyaan pilihan ganda
- `essay`: Pertanyaan essay

### Difficulty
- `easy`: Mudah
- `medium`: Sedang
- `hard`: Sulit

## Database Schema

### Tables
1. `interview_tests` - Tabel utama untuk tes
2. `interview_questions` - Tabel pertanyaan
3. `interview_question_options` - Tabel opsi jawaban (untuk multiple choice)
4. `interview_test_submissions` - Tabel submission kandidat (untuk future)
5. `interview_test_answers` - Tabel jawaban kandidat (untuk future)

## Frontend Integration

API ini sudah terintegrasi dengan frontend yang ada di:
- `karir-nusantara-hub/src/pages/InterviewTest.tsx`
- `karir-nusantara-hub/src/api/interview-test.ts`
- `karir-nusantara-hub/src/components/interview-test/`

Endpoint frontend akan memanggil API dengan base URL:
```typescript
const API_BASE_URL = 'http://localhost:8080/api/v1';
```

## Testing

Jalankan backend:
```bash
cd karir-nusantara-api
make run
# atau
go run cmd/api/main.go
```

Server akan berjalan di `http://localhost:8080`

## Next Steps (Future Implementation)

1. **Candidate Test Taking**
   - Endpoint untuk kandidat mengambil tes
   - Timer management
   - Auto-submit saat waktu habis

2. **Grading System**
   - Auto-grading untuk multiple choice
   - Manual grading untuk essay questions
   - Grade review interface

3. **Test Analytics**
   - Statistics per test
   - Question difficulty analysis
   - Candidate performance tracking

4. **Test Assignment**
   - Assign test ke job posting
   - Assign test ke specific candidate
   - Required vs optional tests
