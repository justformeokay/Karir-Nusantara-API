# Announcements Feature Implementation

## Overview
Sistem announcements (notifikasi, banner, dan informasi) telah diimplementasikan untuk platform Karir Nusantara. Fitur ini memungkinkan Super Admin untuk mengelola pengumuman yang akan ditampilkan di berbagai frontend.

## Database
**Table:** `announcements`
**Migration File:** `/karir-nusantara-api/migrations/015_create_announcements.sql`

### Struktur Kolom:
- `id` - Primary key
- `title` - Judul pengumuman
- `content` - Isi pengumuman
- `type` - Tipe: 'notification', 'banner', 'information'
- `target_audience` - Target: 'all', 'company', 'candidate', 'partner'
- `is_active` - Status aktif/nonaktif
- `priority` - Prioritas 1-10 (higher = more important)
- `start_date` - Tanggal mulai (optional)
- `end_date` - Tanggal selesai (optional)
- `created_by`, `updated_by` - Admin yang membuat/update
- `created_at`, `updated_at` - Timestamp

## Backend API

### Public Endpoints (No Authentication Required)
Endpoint ini dapat diakses oleh semua frontend tanpa autentikasi:

1. **GET /api/v1/announcements**
   - Mendapatkan semua announcements aktif
   - Auto-filter: `is_active = true`
   - Sorted by: priority DESC, created_at DESC

2. **GET /api/v1/announcements/notifications**
   - Mendapatkan notifications aktif
   - Filter: `type = 'notification'` AND `is_active = true`

3. **GET /api/v1/announcements/banners**
   - Mendapatkan banners aktif
   - Filter: `type = 'banner'` AND `is_active = true`

4. **GET /api/v1/announcements/information**
   - Mendapatkan information aktif
   - Filter: `type = 'information'` AND `is_active = true`

### Admin Endpoints (Require Authentication & Admin Role)
Endpoint ini hanya untuk Super Admin di karir-nusantara-hub:

1. **GET /api/v1/admin/announcements**
   - List semua announcements dengan pagination
   - Query params: `page`, `per_page`, `type`, `target_audience`, `is_active`, `search`

2. **GET /api/v1/admin/announcements/:id**
   - Mendapatkan detail announcement

3. **POST /api/v1/admin/announcements**
   - Membuat announcement baru

4. **PUT /api/v1/admin/announcements/:id**
   - Update announcement

5. **PATCH /api/v1/admin/announcements/:id/toggle**
   - Toggle status aktif/nonaktif

6. **DELETE /api/v1/admin/announcements/:id**
   - Hapus announcement

## Frontend Implementation

### 1. karir-nusantara-hub (Admin Panel - Port 5175)
**Location:** `/src/pages/Notifications.tsx`

**Features:**
- Full CRUD management untuk announcements
- Tabs: Notifications, Banners, Information
- Filter by type
- Toggle active/inactive status
- Preview modal
- Pagination support
- Real-time updates

**API Client:** `/src/api/announcements.ts`

### 2. karir-nusantara (Job Seekers - Port 8080)
**Location:** `/src/components/announcements/AnnouncementComponents.tsx`

**Components:**
1. **AnnouncementBanner**
   - Ditampilkan di top of page
   - Auto-rotating carousel untuk multiple banners
   - Dismissible
   - Filter: `target_audience = 'all'` OR `'candidate'`

2. **NotificationsList**
   - Menampilkan max 3 notifications
   - Expandable untuk melihat semua
   - Blue themed
   - Filter: relevant untuk job seekers

3. **InformationSection**
   - Menampilkan max 2 information cards
   - Green themed
   - Grid layout (2 columns on desktop)

**API Client:** `/src/api/announcements.ts`

**Integration:** Ditambahkan di `HomePage.tsx`

### 3. karir-nusantara-company (Companies - Port 5174)
**Status:** Ready to implement
**Target Audience Filter:** `'all'` OR `'company'`

### 4. karir-nusantara-partners (Partners - Port 5176)
**Status:** Ready to implement
**Target Audience Filter:** `'all'` OR `'partner'`

## Module Files

### Backend (`/karir-nusantara-api/internal/modules/announcements/`)
- `entity.go` - Data structures & types
- `repository.go` - Database operations
- `service.go` - Business logic
- `handler.go` - HTTP handlers
- `routes.go` - Route definitions

### Integration Points
- `/karir-nusantara-api/cmd/api/main.go` - Public routes registration
- `/karir-nusantara-api/internal/modules/admin/routes.go` - Admin routes registration

## Usage Examples

### Create Announcement (Admin)
```bash
curl -X POST http://localhost:8081/api/v1/admin/announcements \
  -H "Authorization: Bearer <admin_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Pemeliharaan Sistem",
    "content": "Sistem akan maintenance tanggal 15 Feb pukul 00:00-04:00 WIB",
    "type": "notification",
    "target_audience": "all",
    "is_active": true,
    "priority": 10
  }'
```

### Get Public Notifications (Job Seeker)
```bash
curl http://localhost:8081/api/v1/announcements/notifications
```

### Filter by Target Audience (Frontend)
```typescript
import { getNotifications, filterForJobSeekers } from '@/api/announcements';

const response = await getNotifications();
const relevantNotifications = filterForJobSeekers(response.data);
// Returns only announcements with target_audience = 'all' or 'candidate'
```

## Testing

### Test Public API
```bash
# All announcements
curl http://localhost:8081/api/v1/announcements | jq .

# Notifications only
curl http://localhost:8081/api/v1/announcements/notifications | jq .

# Banners only
curl http://localhost:8081/api/v1/announcements/banners | jq .

# Information only
curl http://localhost:8081/api/v1/announcements/information | jq .
```

### Test Admin API
```bash
# Login as admin first
TOKEN=$(curl -s -X POST http://localhost:8081/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@karirnusantara.com","password":"admin123"}' \
  | jq -r '.data.access_token')

# List all announcements
curl -H "Authorization: Bearer $TOKEN" \
  http://localhost:8081/api/v1/admin/announcements | jq .

# Filter by type
curl -H "Authorization: Bearer $TOKEN" \
  "http://localhost:8081/api/v1/admin/announcements?type=notification" | jq .
```

## Important Notes

1. **No Breaking Changes**: Semua endpoint public menggunakan path baru (`/api/v1/announcements/*`), sehingga tidak mengganggu API existing untuk frontend lain.

2. **Authentication**: 
   - Public endpoints: NO authentication required
   - Admin endpoints: Require JWT token + admin role

3. **Filtering**:
   - Backend: Auto-filter `is_active = true` untuk public endpoints
   - Frontend: Additional filter by `target_audience` di client side

4. **Performance**: 
   - Public endpoints sorted by priority & date
   - Pagination available untuk admin panel
   - Frontend components limit display (3 notifications, 2 information)

5. **Future Enhancement**:
   - Email notification untuk announcements penting
   - Push notifications
   - Schedule announcements (start_date, end_date)
   - Analytics tracking

## Maintenance

### Add New Announcement Type
1. Update `type` enum di migration file
2. Update `AnnouncementType` di backend entity.go
3. Update frontend types
4. Add new component if needed

### Change Target Audience
1. Update `target_audience` enum di migration
2. Update `TargetAudience` type di backend & frontend
3. Update filter functions

## Support Contacts
- Backend: karir-nusantara-api
- Admin Frontend: karir-nusantara-hub
- Job Seeker Frontend: karir-nusantara
- Company Frontend: karir-nusantara-company
- Partner Frontend: karir-nusantara-partners
