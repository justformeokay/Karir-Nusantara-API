# Admin Partner Management API Documentation

API endpoints untuk mengelola referral partners di super admin dashboard (karir-nusantara-hub).

## Base URL
```
http://localhost:8081/api/v1
```

## Authentication
Semua endpoints memerlukan admin authentication:
```
Authorization: Bearer <admin_token>
```

---

## Partner Management

### 1. Get Partners List
Get paginated list of all partners.

**Endpoint:** `GET /admin/partners`

**Query Parameters:**
| Parameter | Type   | Required | Default | Description |
|-----------|--------|----------|---------|-------------|
| status    | string | No       | all     | Filter by status: `active`, `pending`, `suspended`, `all` |
| search    | string | No       | -       | Search by name, email, or referral code |
| page      | int    | No       | 1       | Page number |
| limit     | int    | No       | 10      | Items per page (max: 100) |

**Response:**
```json
{
  "success": true,
  "data": {
    "partners": [
      {
        "id": 1,
        "full_name": "Ahmad Pratama",
        "email": "ahmad.pratama@email.com",
        "phone": "08123456789",
        "referral_code": "AHMAD2024",
        "commission_rate": 40,
        "status": "active",
        "total_referrals": 8,
        "total_commission": 47850000,
        "available_balance": 12500000,
        "paid_amount": 35350000,
        "bank_info": {
          "bank_name": "Bank Central Asia",
          "bank_account_number": "1234567890",
          "bank_account_holder": "Ahmad Pratama",
          "is_verified": true
        },
        "created_at": "2024-01-15T00:00:00+07:00"
      }
    ],
    "pagination": {
      "page": 1,
      "limit": 10,
      "total": 3,
      "total_pages": 1
    }
  }
}
```

---

### 2. Get Partner Detail
Get detailed information about a specific partner.

**Endpoint:** `GET /admin/partners/{id}`

**Response:**
```json
{
  "success": true,
  "data": {
    "id": 1,
    "user_id": 22,
    "full_name": "Ahmad Pratama",
    "email": "ahmad.pratama@email.com",
    "phone": "08123456789",
    "referral_code": "AHMAD2024",
    "commission_rate": 40,
    "status": "active",
    "bank_info": {
      "bank_name": "Bank Central Asia",
      "bank_account_number": "1234567890",
      "bank_account_holder": "Ahmad Pratama",
      "is_verified": true
    },
    "total_referrals": 8,
    "total_commission": 47850000,
    "available_balance": 12500000,
    "pending_balance": 0,
    "paid_amount": 35350000,
    "approved_at": "2026-02-04T20:06:58+07:00",
    "last_payout_date": "2026-02-05T10:00:00+07:00",
    "notes": "Approved by admin",
    "created_at": "2024-01-15T00:00:00+07:00",
    "updated_at": "2026-02-04T20:06:58+07:00"
  }
}
```

---

### 3. Update Partner Status
Update partner status (activate, suspend).

**Endpoint:** `PATCH /admin/partners/{id}/status`

**Request Body:**
```json
{
  "status": "suspended",  // active, suspended, pending
  "notes": "Suspended due to policy violation"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Partner status updated successfully"
}
```

---

### 4. Approve Partner
Approve a pending partner application.

**Endpoint:** `POST /admin/partners/{id}/approve`

**Request Body:**
```json
{
  "commission_rate": 35,  // optional, default: 10
  "notes": "Approved by admin"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Partner approved successfully"
}
```

---

## Referral Management

### 5. Get Referred Companies
Get list of companies referred by partners.

**Endpoint:** `GET /admin/referrals/companies`

**Query Parameters:**
| Parameter | Type   | Required | Default | Description |
|-----------|--------|----------|---------|-------------|
| search    | string | No       | -       | Search by company name or partner name |
| page      | int    | No       | 1       | Page number |
| limit     | int    | No       | 10      | Items per page |

**Response:**
```json
{
  "success": true,
  "data": {
    "companies": [
      {
        "id": 1,
        "company_id": 5,
        "company_name": "PT Karya Developer",
        "partner_info": {
          "id": 1,
          "name": "Ahmad Pratama",
          "referral_code": "AHMAD2024"
        },
        "total_transactions": 5,
        "total_revenue_generated": 25000000,
        "total_commission": 10000000,
        "registration_date": "2024-03-15T00:00:00+07:00",
        "status": "verified"
      }
    ],
    "pagination": {
      "page": 1,
      "limit": 10,
      "total": 1,
      "total_pages": 1
    }
  }
}
```

---

### 6. Get Referral Stats
Get overall referral program statistics.

**Endpoint:** `GET /admin/referrals/stats`

**Response:**
```json
{
  "success": true,
  "data": {
    "total_partners": 3,
    "active_partners": 2,
    "pending_partners": 1,
    "total_referred_companies": 8,
    "total_commission_generated": 47850000,
    "pending_payouts": 5000000,
    "total_paid_out": 35350000,
    "partners_with_balance": 1
  }
}
```

---

## Payout Management

### 7. Get Payouts List
Get list of commission payouts.

**Endpoint:** `GET /admin/payouts`

**Query Parameters:**
| Parameter | Type   | Required | Default | Description |
|-----------|--------|----------|---------|-------------|
| status    | string | No       | all     | Filter: `pending`, `processing`, `completed`, `all` |
| search    | string | No       | -       | Search by partner name |
| page      | int    | No       | 1       | Page number |
| limit     | int    | No       | 10      | Items per page |

**Response:**
```json
{
  "success": true,
  "data": {
    "payouts": [
      {
        "id": 1,
        "partner_id": 1,
        "partner_name": "Ahmad Pratama",
        "amount": 5000000,
        "status": "pending",
        "bank_info": {
          "bank_name": "Bank Central Asia",
          "bank_account_number": "1234567890",
          "bank_account_holder": "Ahmad Pratama"
        },
        "payout_proof_url": "",
        "requested_at": "2026-02-06T02:28:49+07:00",
        "paid_at": null,
        "notes": ""
      }
    ],
    "pagination": {
      "page": 1,
      "limit": 10,
      "total": 1,
      "total_pages": 1
    }
  }
}
```

---

### 8. Get Payout Stats
Get payout statistics.

**Endpoint:** `GET /admin/payouts/stats`

**Response:**
```json
{
  "success": true,
  "data": {
    "total_commission_generated": 47850000,
    "pending_payouts": 5000000,
    "total_paid_out": 35350000,
    "partners_with_balance": 1
  }
}
```

---

### 9. Get Partner Balances
Get partners with available balance for payout.

**Endpoint:** `GET /admin/payouts/balances`

**Query Parameters:**
| Parameter | Type | Required | Default | Description |
|-----------|------|----------|---------|-------------|
| page      | int  | No       | 1       | Page number |
| limit     | int  | No       | 10      | Items per page |

**Response:**
```json
{
  "success": true,
  "data": {
    "partners": [
      {
        "id": 1,
        "partner_name": "Ahmad Pratama",
        "email": "ahmad.pratama@email.com",
        "available_balance": 12500000,
        "pending_balance": 0,
        "total_paid": 35350000,
        "bank_info": {
          "bank_name": "Bank Central Asia",
          "bank_account_number": "1234567890",
          "bank_account_holder": "Ahmad Pratama",
          "is_verified": true
        },
        "last_payout_date": "2026-02-05T10:00:00+07:00"
      }
    ],
    "pagination": {
      "page": 1,
      "limit": 10,
      "total": 1,
      "total_pages": 1
    }
  }
}
```

---

### 10. Create Payout
Create a new payout request for a partner.

**Endpoint:** `POST /admin/payouts`

**Request Body:**
```json
{
  "partner_id": 1,
  "amount": 5000000,
  "notes": "Monthly payout request"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Payout created successfully",
  "payout_id": 1
}
```

---

### 11. Process Payout
Mark a payout as completed (paid).

**Endpoint:** `POST /admin/payouts/{id}/process`

**Request Body:**
```json
{
  "payout_proof_url": "https://example.com/proof/transfer123.jpg",
  "notes": "Transfer completed"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Payout processed successfully"
}
```

---

## Error Responses

All endpoints return error responses in this format:

```json
{
  "success": false,
  "error": "Error message here"
}
```

Common HTTP status codes:
- `400` - Bad Request (invalid input)
- `401` - Unauthorized (missing/invalid token)
- `403` - Forbidden (not admin)
- `404` - Not Found
- `500` - Internal Server Error
