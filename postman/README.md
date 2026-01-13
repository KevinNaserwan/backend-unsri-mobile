# Postman Collection

Postman collection untuk testing UNSRI Backend API.

## 📦 Files

- `UNSRI_Backend_API.postman_collection.json` - API Collection
- `UNSRI_Backend_Environment.postman_environment.json` - Environment variables

## 🚀 Quick Start

### 1. Import Collection

1. Buka Postman
2. Klik **Import** button
3. Pilih file `UNSRI_Backend_API.postman_collection.json`
4. Klik **Import**

### 2. Import Environment

1. Klik **Environments** di sidebar kiri
2. Klik **Import**
3. Pilih file `UNSRI_Backend_Environment.postman_environment.json`
4. Klik **Import**

### 3. Set Environment Variables

1. Pilih environment **UNSRI Backend - Local**
2. Set `base_url` sesuai environment Anda:
   - Local: `http://localhost:8080`
   - Production: `https://api.unsri.ac.id`

### 4. Get Access Token

1. Buka request **Authentication > Login**
2. Update email dan password di body
3. Klik **Send**
4. Copy `access_token` dari response
5. Set variable `access_token` di environment

### 5. Test API

Sekarang Anda bisa test semua API endpoints. Token akan otomatis digunakan untuk authenticated requests.

## 📋 Collection Structure

Collection dirapikan menjadi tiga kategori utama:

- Core
  - Authentication (Register, Login)
  - Users (Get Profile, Update Profile, Search Users)
  - Search (Search, Global Search)
  - File Storage (Upload, List, Get by ID, Download, Upload Document)
  - Access (Validate QR, History, Log, Permissions)
  - Notifications, Broadcasts, Calendar, Master Data

- Academic
  - Attendance (Scan QR, History)
  - QR Code (Generate Class QR)
  - Courses & Schedules (List Courses, List Schedules)
  - Academic Reports

- HRIS
  - Work Attendance (Check In, Check Out, Records)
  - Location (Tap In, Tap Out)
  - HRIS Reports

## 🧭 Alur Penggunaan (Guided Flows)

### Absensi Mahasiswa (Academic)

1. Siapkan QR kelas
   - Endpoint: POST /api/v1/qr/class/generate
   - Prasyarat:
     - Role dosen atau staff
     - Kelas/schedule tersedia (course & schedule terdaftar)
     - Body berisi course_id dan schedule_id
2. Mahasiswa scan QR untuk absen
   - Endpoint: POST /api/v1/attendance/qr/scan
   - Prasyarat:
     - Mahasiswa login (punya access_token)
     - QR valid dan aktif untuk sesi kelas
3. Cek riwayat absensi
   - Endpoint: GET /api/v1/attendance/history
   - Opsional: page, per_page

### Kehadiran Kerja (HRIS)

1. Pastikan shift/sesi kerja (opsional)
   - Endpoint terkait:
     - GET /api/v1/work-attendance/shifts
     - GET /api/v1/work-attendance/sessions
   - Tujuan: ada shift/sesi aktif untuk user sesuai kebijakan HR
2. Check-in
   - Endpoint: POST /api/v1/work-attendance/check-in
   - Body (multipart/form-data):
     - file `selfie` (wajib)
     - `latitude`, `longitude` (opsional)
     - `schedule_id` (opsional)
     - `is_via_unsri_wifi` (boolean, opsional/tergantung kebijakan)
   - Prasyarat:
     - User login
     - Berada dalam geofence/wifi yang diset (jika diaktifkan)
3. Check-out
   - Endpoint: POST /api/v1/work-attendance/check-out
   - Body (multipart/form-data):
     - file `selfie` (wajib)
     - `latitude`, `longitude` (opsional)
     - `schedule_id` (opsional)
     - `is_via_unsri_wifi` (boolean, opsional/tergantung kebijakan)
4. Lihat riwayat kehadiran kerja
   - Endpoint: GET /api/v1/work-attendance/records
   - Opsional: page, per_page, start_date, end_date, user_id

### Akses Gate Kampus (Core)

1. Generate QR akses personal
   - Endpoint: GET /api/v1/qr/access/generate
   - Prasyarat: user login
2. Validasi QR di gate
   - Endpoint: POST /api/v1/qr/gate/validate
   - Catatan: endpoint publik untuk sistem gate (tanpa token)

## 🔧 Environment Variables

- `base_url` - Base URL untuk API (default: http://localhost:8080)
- `access_token` - JWT access token (akan di-set setelah login)
- `refresh_token` - JWT refresh token
- `user_id` - Current user ID

## 💡 Tips

1. **Auto-save Token**: Buat test script di Login request untuk auto-save token:
```javascript
if (pm.response.code === 200) {
    var jsonData = pm.response.json();
    pm.environment.set("access_token", jsonData.data.access_token);
}
```

2. **Pre-request Script**: Untuk auto-include token di semua requests, tambahkan di collection level:
```javascript
pm.request.headers.add({
    key: "Authorization",
    value: "Bearer " + pm.environment.get("access_token")
});
```

## 📚 Documentation

Lihat [API Documentation](../docs/API.md) untuk dokumentasi lengkap semua endpoints.
