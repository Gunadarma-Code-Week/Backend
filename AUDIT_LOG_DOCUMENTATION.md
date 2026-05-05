# Audit Log System Documentation

## Overview
Sistem audit log merekam aktivitas user yang terauthentikasi di backend. Data yang dicatat meliputi:
- **Method**: HTTP method (POST, PUT, DELETE, PATCH, dll) - **GET requests tidak dicatat**
- **Endpoint**: Path/route yang diakses
- **User ID**: ID pengguna yang melakukan aksi
- **Request Body**: Data yang dikirim user
- **Response Code**: HTTP status code dari response
- **IP Address**: Alamat IP pengguna
- **User Agent**: Browser/client info
- **Timestamp**: Waktu aktivitas dicatat

---

## Komponen Sistem

### 1. Entity: AuditLog (`entity/audit_log.go`)
Menyimpan struktur data untuk audit log:
```go
type AuditLog struct {
    ID            uint64
    UserID        uint64            // Foreign key ke User
    Method        string            // HTTP method
    Endpoint      string            // Route/path
    RequestBody   json.RawMessage   // Request body dalam JSON
    ResponseCode  int               // HTTP status code
    ResponseBody  json.RawMessage   // Response body dalam JSON
    IPAddress     string            // Client IP
    UserAgent     string            // Browser/client info
    ErrorMessage  *string           // Pesan error jika ada
    CreatedAt     time.Time         // Waktu pencatatan
    UpdatedAt     time.Time
}
```

### 2. Repository: AuditLogRepository (`repository/audit-log-repository.go`)
Menangani operasi database untuk audit logs:
- `Create()`: Membuat record audit log baru
- `FindByUserID()`: Mencari log berdasarkan user ID
- `FindByEndpoint()`: Mencari log berdasarkan endpoint
- `FindByDateRange()`: Mencari log dalam range tanggal
- `FindByUserIDAndDateRange()`: Kombinasi filter user dan tanggal
- `FindAll()`: Mendapatkan semua log dengan pagination

### 3. Service: AuditLogService (`service/audit-log-service.go`)
Business logic untuk audit logging:
- `RecordActivity()`: Mencatat aktivitas tanpa response
- `RecordActivityWithResponse()`: Mencatat aktivitas dengan response
- `RecordActivityWithError()`: Mencatat aktivitas dengan error
- `GetUserActivityLogs()`: Retrieve log user tertentu
- `GetActivityLogsByDateRange()`: Retrieve log dalam range tanggal

### 4. Middleware: AuditMiddleware (`middleware/audit.go`)
Middleware untuk menangkap dan merekam aktivitas user:
- Dijalankan setelah authentication middleware
- Menangkap method, endpoint, request body, response code
- Mengambil IP address dan user agent
- Secara otomatis merekam setiap request **modifikasi data** (POST, PUT, DELETE, PATCH) dari authenticated user
- **GET dan HEAD requests TIDAK dicatat** untuk mengurangi noise

### 5. Handler: AuditLogHandler (`handler/audit-log.go`)
Handler API untuk mengakses audit logs:
- `GetAllAuditLogs()`: Endpoint untuk mendapatkan semua logs
- `GetUserAuditLogs()`: Endpoint untuk logs user tertentu
- `GetAuditLogsByDateRange()`: Endpoint untuk logs dalam range tanggal

---

## Instalasi & Setup

### 1. Database Migration
Audit log table akan otomatis dibuat melalui AutoMigrate di `config/database.go`.
Tabel `audit_logs` akan memiliki kolom-kolom sesuai dengan entity AuditLog.

### 2. Inisialisasi di Router
Sudah diintegrasikan di `router/setup.go`:
```go
auditLogRepository = repository.NewAuditLogRepository(database)
auditLogService    = service.NewAuditLogService(auditLogRepository)
auditLogHandler    = handler.NewAuditLogHandler(auditLogService)
auditMiddleware    = middleware.NewAuditMiddleware(auditLogService)
```

### 3. Middleware Configuration
Middleware audit ditambahkan ke authenticated routes:
```go
mustAuth := router.Group("")
mustAuth.Use(authMiddleware.JwtAuthMiddleware)
mustAuth.Use(auditMiddleware.AuditLogMiddleware)  // ← Audit logging
```

---

## API Endpoints

### 1. Get All Audit Logs
**Endpoint**: `GET /api/v1/gcw/resources/admin/audit-logs`
**Auth**: Required (Admin only)
**Query Parameters**:
- `page` (optional, default: 1): Halaman untuk pagination
- `limit` (optional, default: 10): Jumlah item per halaman

**Response**:
```json
{
  "code": "success",
  "message": "Successfully retrieved audit logs",
  "data": {
    "logs": [
      {
        "id": 1,
        "user_id": 5,
        "method": "POST",
        "endpoint": "/api/v1/gcw/resources/team/registration/hackathon",
        "request_body": {"team_name": "Team A", "members": [1, 2, 3]},
        "response_code": 201,
        "ip_address": "192.168.1.100",
        "user_agent": "Mozilla/5.0...",
        "created_at": "2026-05-03T10:30:00Z"
      }
    ],
    "total": 150,
    "page": 1,
    "limit": 10,
    "total_pages": 15
  }
}
```

### 2. Get User Audit Logs
**Endpoint**: `GET /api/v1/gcw/resources/admin/audit-logs/user/:user_id`
**Auth**: Required (Admin only)
**Path Parameters**:
- `user_id` (required): ID user yang ingin dilihat lognya

**Query Parameters**:
- `page` (optional, default: 1)
- `limit` (optional, default: 10)

**Response**: Sama seperti endpoint 1, tapi hanya untuk user tertentu

**Contoh**:
```
GET /api/v1/gcw/resources/admin/audit-logs/user/5?page=1&limit=10
```

### 3. Get Audit Logs by Date Range
**Endpoint**: `GET /api/v1/gcw/resources/admin/audit-logs/date-range`
**Auth**: Required (Admin only)
**Query Parameters**:
- `start_date` (required): Format YYYY-MM-DD
- `end_date` (required): Format YYYY-MM-DD
- `user_id` (optional): Filter berdasarkan user ID
- `page` (optional, default: 1)
- `limit` (optional, default: 10)

**Response**: Sama seperti endpoint 1

**Contoh**:
```
GET /api/v1/gcw/resources/admin/audit-logs/date-range?start_date=2026-05-01&end_date=2026-05-03&limit=20
```

Dengan filter user:
```
GET /api/v1/gcw/resources/admin/audit-logs/date-range?start_date=2026-05-01&end_date=2026-05-03&user_id=5&limit=20
```

---

## Usage Examples

### Contoh: Tracking registrasi tim
Ketika user melakukan registrasi team hackathon melalui endpoint:
```
POST /api/v1/gcw/resources/team/registration/hackathon
Content-Type: application/json
Authorization: Bearer <token>

{
  "team_name": "Team Alpha",
  "member_ids": [1, 2, 3]
}
```

Audit log otomatis dicatat:
```
{
  "user_id": 5,
  "method": "POST",
  "endpoint": "/api/v1/gcw/resources/team/registration/hackathon",
  "request_body": {
    "team_name": "Team Alpha",
    "member_ids": [1, 2, 3]
  },
  "response_code": 201,
  "ip_address": "192.168.1.100",
  "user_agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64)...",
  "created_at": "2026-05-03T10:30:00Z"
}
```

### Contoh: Tracking profile update
```
POST /api/v1/gcw/resources/profile/my
Content-Type: application/json
Authorization: Bearer <token>

{
  "name": "John Doe",
  "email": "john@example.com",
  "phone": "081234567890"
}
```

Audit log dicatat dengan detail request dan response code.

---

## Database Schema

Tabel `audit_logs` memiliki struktur:

```sql
CREATE TABLE audit_logs (
  id BIGSERIAL PRIMARY KEY,
  user_id BIGINT NOT NULL,
  method VARCHAR(10) NOT NULL,
  endpoint VARCHAR(255) NOT NULL,
  request_body JSONB,
  response_code INTEGER DEFAULT 0,
  response_body JSONB,
  ip_address VARCHAR(45),
  user_agent TEXT,
  error_message TEXT,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexes untuk performa
CREATE INDEX idx_audit_logs_user_id ON audit_logs(user_id);
CREATE INDEX idx_audit_logs_endpoint ON audit_logs(endpoint);
CREATE INDEX idx_audit_logs_created_at ON audit_logs(created_at);
```

---

## Fitur-Fitur

✅ **Automatic Logging**: Semua authenticated requests (POST, PUT, DELETE, PATCH) otomatis tercatat, GET requests tidak dicatat
✅ **Detailed Tracking**: Method, endpoint, request body, response code, IP, user agent
✅ **Pagination**: Semua endpoint mendukung pagination untuk performa
✅ **Date Range Filter**: Bisa filter logs berdasarkan range tanggal
✅ **User Filter**: Bisa melihat logs user tertentu
✅ **Admin Only**: Hanya admin yang bisa mengakses audit logs
✅ **JSONB Storage**: Request/response body disimpan dalam format JSON untuk fleksibilitas query

---

## Endpoints yang Tercatat

Semua endpoint authenticated dengan method **POST, PUT, DELETE, PATCH** akan otomatis tercatat. **GET requests tidak dicatat** untuk mengurangi noise di audit log.

Endpoints yang tercatat:

- ✅ Profile update: `POST /profile/my`
- ✅ Team registration: `POST /team/registration/*`
- ✅ Submission: `POST /submission/*`
- ✅ Seminar: `POST /seminar/join`
- ✅ Admin routes: `PUT/DELETE /admin/*`
- ✅ Semua method modifikasi data (POST, PUT, DELETE, PATCH) lainnya

Endpoints yang TIDAK dicatat:

- ❌ `GET /profile/my` (read-only)
- ❌ `GET /profile/events` (read-only)
- ❌ `GET /team/registration/find/*` (read-only)
- ❌ `GET /admin/users` (read-only)
- ❌ `GET /seminar/my-ticket` (read-only)
- ❌ Semua GET requests

---

## Customization

### Mengubah Limit Default Logging
Edit di `middleware/audit.go`, fungsi `shouldSkipAuditLog()`:
```go
func shouldSkipAuditLog(path string) bool {
    skipPaths := []string{
        "/health",
        "/metrics",
        // Tambah path yang tidak perlu dilog di sini
    }
    // ...
}
```

### Menambah Info Tambahan
Edit `entity/audit_log.go` untuk menambah field baru:
```go
type AuditLog struct {
    // ... existing fields ...
    CustomField string `gorm:"varchar(255)"`
}
```

Kemudian update migration di `config/database.go`.

---

## Best Practices

1. **Regular Cleanup**: Pertimbangkan untuk archive/delete audit logs lama secara berkala
2. **Performance**: Gunakan index pada kolom `user_id`, `endpoint`, dan `created_at`
3. **Security**: Pastikan hanya admin yang bisa akses audit logs
4. **Privacy**: Pertimbangkan untuk tidak log password atau sensitive data lainnya
5. **Retention Policy**: Tentukan kebijakan berapa lama audit logs disimpan

---

## Troubleshooting

**Q: Audit logs tidak tercatat?**
A: Pastikan:
- Middleware audit sudah ditambahkan ke route group
- User sudah terauthentikasi (token valid)
- Database sudah di-migrate
- Tidak ada error di service layer

**Q: Query terlalu lambat?**
A: Tambahkan index:
```sql
CREATE INDEX idx_audit_logs_user_created ON audit_logs(user_id, created_at DESC);
```

**Q: Storage terlalu besar?**
A: Implementasikan retention policy:
```sql
DELETE FROM audit_logs WHERE created_at < NOW() - INTERVAL '90 days';
```

---

## Kesimpulan

Sistem audit log ini memberikan transparency dan accountability untuk semua aktivitas user terauthentikasi di aplikasi. Admin dapat memonitor aktivitas user, melacak perubahan, dan melakukan investigasi jika diperlukan.
