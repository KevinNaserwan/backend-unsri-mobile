# API Endpoints & Payloads

Dokumentasi komprehensif semua endpoint dan payload berdasarkan implementasi kode saat ini.

## Base URL
- Local: http://localhost:8080
- Production: https://api.unsri.ac.id

## Authentication
- Header Authorization diperlukan untuk endpoint bertanda protected: Authorization: Bearer <access_token>

### Auth
- POST /api/v1/auth/register
  - Body JSON:
    - email (required)
    - password (required, min 8)
    - role (required, oneof: mahasiswa dosen staff)
    - nim (mahasiswa)
    - nip (dosen/staff)
    - nama (required)
    - prodi (optional)
    - angkatan (mahasiswa, optional)
    - jabatan (staff, optional)
    - unit (staff, optional)
  - Referensi: [service.go](file:///Users/ahmadnaufalmuzakki/Documents/KERJAAN/Meetsin.Id/2025/UNSRI/unsri%20app/backend-unsri-mobile/internal/auth/service/service.go#L103-L115)
  - Contoh Response:
    ```json
    {
      "success": true,
      "data": {
        "id": "user-uuid",
        "email": "student@unsri.ac.id",
        "role": "mahasiswa",
        "is_active": true,
        "mahasiswa": { "nim": "123", "nama": "Student" }
      }
    }
    ```
  - Sumber Data: users, mahasiswa/dosen/staff (AuthRepository, GORM)
- POST /api/v1/auth/login
  - Body JSON:
    - email (required)
    - password (required)
  - Referensi: [service.go](file:///Users/ahmadnaufalmuzakki/Documents/KERJAAN/Meetsin.Id/2025/UNSRI/unsri%20app/backend-unsri-mobile/internal/auth/service/service.go#L31-L35)
  - Contoh Response:
    ```json
    {
      "success": true,
      "data": {
        "access_token": "jwt-access",
        "refresh_token": "jwt-refresh",
        "user": { "id": "user-uuid", "email": "student@unsri.ac.id", "role": "mahasiswa" }
      }
    }
    ```
  - Sumber Data: users (AuthRepository) dan token dari pkg/jwt
- POST /api/v1/auth/refresh
  - Body JSON:
    - refresh_token (required)
  - Referensi: [service.go](file:///Users/ahmadnaufalmuzakki/Documents/KERJAAN/Meetsin.Id/2025/UNSRI/unsri%20app/backend-unsri-mobile/internal/auth/service/service.go#L288-L297)
  - Contoh Response:
    ```json
    {
      "success": true,
      "data": {
        "access_token": "new-jwt-access",
        "refresh_token": "new-jwt-refresh"
      }
    }
    ```
  - Sumber Data: validasi token & generasi baru via pkg/jwt; user via AuthRepository
- GET /api/v1/auth/verify
  - Header:
    - Authorization: Bearer <access_token>
  - Referensi: [handler.go](file:///Users/ahmadnaufalmuzakki/Documents/KERJAAN/Meetsin.Id/2025/UNSRI/unsri%20app/backend-unsri-mobile/internal/auth/handler/handler.go#L78-L98)
  - Contoh Response:
    ```json
    {
      "success": true,
      "data": { "id": "user-uuid", "email": "student@unsri.ac.id", "role": "mahasiswa", "is_active": true }
    }
    ```
  - Sumber Data: validasi via pkg/jwt, data user via AuthRepository (users + preload role)

## Users
- GET /api/v1/users/profile
  - Protected
  - Contoh Response:
    ```json
    { "success": true, "data": { "id": "user-uuid", "email": "x@unsri.ac.id", "role": "mahasiswa", "is_active": true } }
    ```
  - Sumber Data: users (UserRepository, GORM)
- PUT /api/v1/users/profile
  - Body JSON:
    - email (optional)
    - mahasiswa: { nama?, prodi?, angkatan? } (optional, untuk role mahasiswa)
    - dosen: { nama?, prodi? } (optional)
    - staff: { nama?, jabatan?, unit? } (optional)
  - Referensi: [service.go](file:///Users/ahmadnaufalmuzakki/Documents/KERJAAN/Meetsin.Id/2025/UNSRI/unsri%20app/backend-unsri-mobile/internal/user/service/service.go#L26-L52)
  - Contoh Response:
    ```json
    { "success": true, "data": { "id": "user-uuid", "email": "x@unsri.ac.id", "mahasiswa": { "nama": "Baru" } } }
    ```
  - Sumber Data: users + tabel role terkait (mahasiswa/dosen/staff)
- POST /api/v1/users/avatar
  - Body multipart/form-data:
    - avatar (file) required
  - Referensi: [handler.go](file:///Users/ahmadnaufalmuzakki/Documents/KERJAAN/Meetsin.Id/2025/UNSRI/unsri%20app/backend-unsri-mobile/internal/user/handler/handler.go#L136-L173)
  - Contoh Response:
    ```json
    { "success": true, "data": { "id": "file-uuid", "url": "https://...", "file_type": "avatar" } }
    ```
  - Sumber Data: files (FileStorageService, models.File)
- GET /api/v1/users/search
  - Query:
    - q, role, page, per_page
  - Contoh Response:
    ```json
    { "success": true, "data": [ { "id": "user-uuid", "email": "x@unsri.ac.id" } ], "meta": { "page": 1, "per_page": 20, "total": 1, "total_pages": 1 } }
    ```
  - Sumber Data: users (filter berdasarkan query)
- GET /api/v1/users/:id
  - Contoh Response:
    ```json
    { "success": true, "data": { "id": "user-uuid", "email": "x@unsri.ac.id" } }
    ```
  - Sumber Data: users
- GET /api/v1/users/mahasiswa/:nim
- GET /api/v1/users/dosen/:nip
- GET /api/v1/users/staff/:nip
  - Contoh Response:
    ```json
    { "success": true, "data": { "id": "user-uuid", "mahasiswa": { "nim": "123" } } }
    ```
  - Sumber Data: mahasiswa/dosen/staff + relasi ke users

## Attendance (Academic & Campus)
- POST /api/v1/attendance/qr/generate
  - Body JSON:
    - schedule_id (optional)
    - type (required, oneof: kelas kampus)
    - duration (optional, menit; default 15)
  - Referensi: [service.go](file:///Users/ahmadnaufalmuzakki/Documents/KERJAAN/Meetsin.Id/2025/UNSRI/unsri%20app/backend-unsri-mobile/internal/attendance/service/service.go#L29-L35)
  - Contoh Response:
    ```json
    { "success": true, "data": { "qr": "base64...", "expires_at": "2026-01-13T10:00:00Z" } }
    ```
  - Sumber Data: generate via pkg/qrcode; metadata dari AttendanceService
- POST /api/v1/attendance/qr/scan
  - Body JSON:
    - qr_data (required)
    - latitude, longitude (optional)
  - Referensi: [service.go](file:///Users/ahmadnaufalmuzakki/Documents/KERJAAN/Meetsin.Id/2025/UNSRI/unsri%20app/backend-unsri-mobile/internal/attendance/service/service.go#L98-L104)
  - Contoh Response:
    ```json
    { "success": true, "data": { "status": "hadir", "schedule_id": "sched-uuid", "record_id": "att-uuid" } }
    ```
  - Sumber Data: attendances (models.Attendance) melalui AttendanceRepository
- GET /api/v1/attendance
  - Query:
    - user_id (opsional, admin/staff/dosen saja), start_date, end_date, page, per_page
  - Contoh Response:
    ```json
    { "success": true, "data": [ { "id": "att-uuid", "user_id": "user-uuid", "date": "2026-01-12", "status": "hadir" } ], "meta": { "page": 1, "per_page": 20, "total": 50, "total_pages": 3 } }
    ```
  - Sumber Data: attendances (AttendanceRepository.GetAttendancesByUserID)
- GET /api/v1/attendance/statistics
  - Query:
    - user_id (opsional untuk dosen/staff), start_date, end_date
  - Contoh Response:
    ```json
    { "success": true, "data": { "total": 30, "by_status": { "hadir": 24, "izin": 2, "sakit": 1, "alpa": 3 }, "by_type": { "kelas": 28, "kampus": 2 }, "attendance_rate": 80.0 } }
    ```
  - Sumber Data: agregasi dari table attendances (AttendanceRepository.GetAttendanceStatistics). Presentase = hadir/total*100.
- GET /api/v1/attendance/overview
  - Contoh Response:
    ```json
    {
      "success": true,
      "data": {
        "today_schedules": [ { "id": "sched-uuid", "course_name": "Algoritma" } ],
        "upcoming_schedules": [ { "id": "sched-uuid2", "course_name": "Basis Data" } ],
        "monthly_statistics": { "total": 30, "attendance_rate": 80.0 },
        "current_tap_in": { "status": "IN", "time": "2026-01-13T07:05:00Z" }
      }
    }
    ```
  - Sumber Data: schedules dari repository schedules (via AttendanceRepository helper), statistik dari attendances; status tap in dari repositori tap-in status (AttendanceRepository.GetCurrentTapInStatus)
- GET /api/v1/attendance/history
  - Contoh Response: sama dengan GET /api/v1/attendance (alias)
  - Sumber Data: attendances
- GET /api/v1/attendance/by-course/:courseId
  - Query: start_date, end_date
  - Contoh Response:
    ```json
    { "success": true, "data": [ { "id": "att-uuid", "schedule_id": "sched-uuid", "status": "hadir" } ] }
    ```
  - Sumber Data: join attendances dengan schedules (AttendanceRepository.GetAttendancesByCourseID)
- GET /api/v1/attendance/by-student/:studentId
  - Query: start_date, end_date
  - Contoh Response:
    ```json
    { "success": true, "data": [ { "id": "att-uuid", "user_id": "studentId", "status": "hadir" } ] }
    ```
  - Sumber Data: attendances (AttendanceRepository.GetAttendancesByStudentID)
- POST /api/v1/attendance/manual
  - Body JSON:
    - user_id (required)
    - schedule_id (optional)
    - type (required, oneof: kelas kampus)
    - status (required, oneof: hadir izin sakit alpa terlambat)
    - date (required, YYYY-MM-DD)
    - notes (optional)
  - Referensi: [service.go](file:///Users/ahmadnaufalmuzakki/Documents/KERJAAN/Meetsin.Id/2025/UNSRI/unsri%20app/backend-unsri-mobile/internal/attendance/service/service.go#L224-L232)
  - Contoh Response:
    ```json
    { "success": true, "data": { "id": "att-uuid", "user_id": "user-uuid", "status": "izin" } }
    ```
  - Sumber Data: attendances (insert via AttendanceRepository.CreateAttendance)
- PUT /api/v1/attendance/:id
  - Body JSON:
    - status (required, oneof: hadir izin sakit alpa terlambat)
    - notes (optional)
  - Referensi: [service.go](file:///Users/ahmadnaufalmuzakki/Documents/KERJAAN/Meetsin.Id/2025/UNSRI/unsri%20app/backend-unsri-mobile/internal/attendance/service/service.go#L268-L272)
  - Catatan: endpoint tap-in/tap-out di attendance dihapus

## Work Attendance (HRIS)
- POST /api/v1/work-attendance/check-in
  - Body multipart/form-data:
    - selfie (file) required
    - schedule_id (optional)
    - latitude (optional, number)
    - longitude (optional, number)
    - is_via_unsri_wifi (optional, boolean; wajib true untuk dosen/staff)
    - notes (optional)
  - Referensi: [handler.go](file:///Users/ahmadnaufalmuzakki/Documents/KERJAAN/Meetsin.Id/2025/UNSRI/unsri%20app/backend-unsri-mobile/internal/attendance/handler/handler.go#L554-L615)
  - Contoh Response:
    ```json
    { "success": true, "data": { "id": "record-uuid", "attendance_type": "CHECK_IN", "recorded_at": "2026-01-13T07:05:00Z", "selfie_url": "https://..." } }
    ```
  - Sumber Data: work_attendance_records (AttendanceRepository); file selfie tersimpan di files (FileStorageService)
- POST /api/v1/work-attendance/check-out
  - Body multipart/form-data:
    - selfie (file) required
    - schedule_id, latitude, longitude, is_via_unsri_wifi, notes (optional)
  - Referensi: [handler.go](file:///Users/ahmadnaufalmuzakki/Documents/KERJAAN/Meetsin.Id/2025/UNSRI/unsri%20app/backend-unsri-mobile/internal/attendance/handler/handler.go#L617-L678)
  - Contoh Response:
    ```json
    { "success": true, "data": { "id": "record-uuid", "attendance_type": "CHECK_OUT", "recorded_at": "2026-01-13T17:12:00Z", "selfie_url": "https://..." } }
    ```
  - Sumber Data: work_attendance_records; file selfie di files
- GET /api/v1/work-attendance/records
  - Query:
    - user_id, start_date, end_date, page, per_page
  - Contoh Response:
    ```json
    { "success": true, "data": [ { "id": "record-uuid", "status": "CHECK_IN" } ], "meta": { "page": 1, "per_page": 20, "total": 10, "total_pages": 1 } }
    ```
  - Sumber Data: work_attendance_records (AttendanceRepository.GetWorkAttendanceRecordsByUserID)
- GET /api/v1/work-attendance/shifts
  - Query: is_active?, page, per_page
  - Contoh Response:
    ```json
    { "success": true, "data": [ { "id": "shift-uuid", "shift_name": "Pagi", "start_time": "08:00:00" } ] }
    ```
  - Sumber Data: shift_patterns
- GET /api/v1/work-attendance/shifts/:id
- POST /api/v1/work-attendance/shifts
  - Body JSON:
    - shift_name (required)
    - shift_code (required)
    - start_time (required HH:MM)
    - end_time (required HH:MM)
    - break_duration_minutes? (optional)
    - is_night_shift? (optional)
  - Referensi: [service.go](file:///Users/ahmadnaufalmuzakki/Documents/KERJAAN/Meetsin.Id/2025/UNSRI/unsri%20app/backend-unsri-mobile/internal/attendance/service/service.go#L622-L631)
  - Contoh Response:
    ```json
    { "success": true, "data": { "id": "shift-uuid", "shift_name": "Pagi" } }
    ```
  - Sumber Data: shift_patterns
- PUT /api/v1/work-attendance/shifts/:id
  - Body JSON:
    - shift_name?, start_time?, end_time?, break_duration_minutes?, is_night_shift?, is_active?
- DELETE /api/v1/work-attendance/shifts/:id
- GET /api/v1/work-attendance/user-shifts/:userId
  - Query: date? (YYYY-MM-DD)
  - Contoh Response:
    ```json
    { "success": true, "data": [ { "user_id": "user-uuid", "shift_id": "shift-uuid", "effective_from": "2026-01-01" } ] }
    ```
  - Sumber Data: user_shifts
- POST /api/v1/work-attendance/user-shifts
  - Body JSON:
    - user_id (required)
    - shift_id (required)
    - effective_from (required YYYY-MM-DD)
    - effective_until (optional YYYY-MM-DD)
  - Referensi: [service.go](file:///Users/ahmadnaufalmuzakki/Documents/KERJAAN/Meetsin.Id/2025/UNSRI/unsri%20app/backend-unsri-mobile/internal/attendance/service/service.go#L758-L764)
- GET /api/v1/work-attendance/schedules
  - Query: user_id?, start_date?, end_date?, page, per_page
  - Contoh Response:
    ```json
    { "success": true, "data": [ { "user_id": "user-uuid", "schedule_date": "2026-01-13", "start_time": "08:00:00" } ] }
    ```
  - Sumber Data: work_schedules
- POST /api/v1/work-attendance/schedules
  - Body JSON:
    - user_id (required)
    - schedule_date (required YYYY-MM-DD)
    - shift_id (optional)
    - start_time (required HH:MM)
    - end_time (required HH:MM)
    - work_type? (optional)
    - location? (optional)
    - is_holiday? (optional)

## Courses, Classes, Enrollments
- GET /api/v1/courses
  - Query: prodi?, is_active?, page, per_page
  - Contoh Response:
    ```json
    { "success": true, "data": [ { "id": "course-uuid", "code": "IF101", "name": "Algoritma" } ], "meta": { "page": 1, "per_page": 20, "total": 5, "total_pages": 1 } }
    ```
  - Sumber Data: courses
- GET /api/v1/courses/:id
  - Contoh Response:
    ```json
    { "success": true, "data": { "id": "course-uuid", "code": "IF101", "name": "Algoritma" } }
    ```
  - Sumber Data: courses
- POST /api/v1/courses
  - Body JSON:
    - code, name, credits (required)
    - name_en?, semester?, prodi?, description? (optional)
  - Referensi: [service.go](file:///Users/ahmadnaufalmuzakki/Documents/KERJAAN/Meetsin.Id/2025/UNSRI/unsri%20app/backend-unsri-mobile/internal/course/service/service.go#L22-L31)
  - Contoh Response:
    ```json
    { "success": true, "data": { "id": "course-uuid", "code": "IF101", "name": "Algoritma" } }
    ```
  - Sumber Data: courses
- PUT /api/v1/courses/:id
  - Body JSON:
    - name?, name_en?, credits?, semester?, prodi?, description?, is_active?
  - Contoh Response:
    ```json
    { "success": true, "data": { "id": "course-uuid", "name": "Algoritma Lanjut" } }
    ```
  - Sumber Data: courses
- DELETE /api/v1/courses/:id
  - Contoh Response:
    ```json
    { "success": true, "data": { "message": "deleted" } }
    ```
  - Sumber Data: courses
- GET /api/v1/classes
  - Query: course_id?, dosen_id?, semester?, page, per_page
  - Contoh Response:
    ```json
    { "success": true, "data": [ { "id": "class-uuid", "class_code": "A", "course_id": "course-uuid" } ] }
    ```
  - Sumber Data: classes
- GET /api/v1/classes/:id
  - Contoh Response:
    ```json
    { "success": true, "data": { "id": "class-uuid", "class_code": "A", "course_id": "course-uuid" } }
    ```
  - Sumber Data: classes
- POST /api/v1/classes
  - Body JSON:
    - course_id, class_code, semester, dosen_id, day_of_week, start_time, end_time (required)
    - class_name?, academic_year?, capacity?, assistant_dosen_id?, room? (optional)
  - Referensi: [service.go](file:///Users/ahmadnaufalmuzakki/Documents/KERJAAN/Meetsin.Id/2025/UNSRI/unsri%20app/backend-unsri-mobile/internal/course/service/service.go#L141-L155)
  - Contoh Response:
    ```json
    { "success": true, "data": { "id": "class-uuid", "class_code": "A" } }
    ```
  - Sumber Data: classes
- GET /api/v1/courses/by-student/:studentId
- GET /api/v1/courses/by-lecturer/:lecturerId
  - Contoh Response:
    ```json
    { "success": true, "data": [ { "id": "course-uuid", "code": "IF101" } ] }
    ```
  - Sumber Data: relasi courses/classes terhadap user
- GET /api/v1/enrollments
  - Query: student_id?, class_id?, status?, page, per_page
  - Contoh Response:
    ```json
    { "success": true, "data": [ { "id": "enroll-uuid", "student_id": "student-uuid", "class_id": "class-uuid", "status": "APPROVED" } ] }
    ```
  - Sumber Data: enrollments
- GET /api/v1/enrollments/:id
  - Contoh Response:
    ```json
    { "success": true, "data": { "id": "enroll-uuid", "student_id": "student-uuid", "class_id": "class-uuid" } }
    ```
  - Sumber Data: enrollments
- POST /api/v1/enrollments
  - Body JSON:
    - student_id, class_id, enrollment_date (required)
    - notes? (optional)
  - Referensi: [service.go](file:///Users/ahmadnaufalmuzakki/Documents/KERJAAN/Meetsin.Id/2025/UNSRI/unsri%20app/backend-unsri-mobile/internal/course/service/service.go#L250-L256)
  - Contoh Response:
    ```json
    { "success": true, "data": { "id": "enroll-uuid", "status": "PENDING" } }
    ```
  - Sumber Data: enrollments
- PUT /api/v1/enrollments/:id/status
  - Body JSON:
    - status (required, oneof: PENDING APPROVED REJECTED COMPLETED DROPPED FAILED)
    - notes? (optional)
- PUT /api/v1/enrollments/:id/grade
  - Body JSON:
    - grade? (oneof: A B C D E)
    - score? (number)
    - notes? (optional)
- DELETE /api/v1/enrollments/:id
  - Contoh Response:
    ```json
    { "success": true, "data": { "message": "deleted" } }
    ```
  - Sumber Data: enrollments

## Schedules
- GET /api/v1/schedules
  - Query: dosen_id?, start_date?, end_date?, page, per_page
  - Contoh Response:
    ```json
    { "success": true, "data": [ { "id": "sched-uuid", "course_name": "Algoritma", "date": "2026-01-13", "start_time": "09:00:00" } ] }
    ```
  - Sumber Data: schedules
- GET /api/v1/schedules/today
  - Contoh Response:
    ```json
    { "success": true, "data": [ { "id": "sched-uuid", "course_name": "Algoritma" } ] }
    ```
  - Sumber Data: schedules (service.GetTodaySchedules)
- GET /api/v1/schedules/upcoming
  - Query: limit?
  - Contoh Response:
    ```json
    { "success": true, "data": [ { "id": "sched-uuid", "course_name": "Basis Data" } ] }
    ```
  - Sumber Data: schedules (service.GetUpcomingSchedules)
- GET /api/v1/schedules/calendar/:year/:month
  - Contoh Response:
    ```json
    { "success": true, "data": [ { "date": "2026-01-15", "course_name": "Jaringan Komputer" } ] }
    ```
  - Sumber Data: schedules (service.GetCalendarView)
- GET /api/v1/schedules/:id
  - Contoh Response:
    ```json
    { "success": true, "data": { "id": "sched-uuid", "course_code": "IF101", "course_name": "Algoritma" } }
    ```
  - Sumber Data: schedules
- POST /api/v1/schedules
  - Body JSON:
    - course_id? (optional)
    - course_code, course_name, dosen_id, day_of_week, start_time, end_time, date (required)
  - Referensi: [service.go](file:///Users/ahmadnaufalmuzakki/Documents/KERJAAN/Meetsin.Id/2025/UNSRI/unsri%20app/backend-unsri-mobile/internal/schedule/service/service.go#L22-L33)
  - Contoh Response:
    ```json
    { "success": true, "data": { "id": "sched-uuid", "course_code": "IF101", "date": "2026-01-20" } }
    ```
  - Sumber Data: schedules
- PUT /api/v1/schedules/:id
  - Body JSON:
    - course_code?, course_name?, room?, start_time?, end_time?, is_active?
  - Contoh Response:
    ```json
    { "success": true, "data": { "id": "sched-uuid", "room": "A101" } }
    ```
  - Sumber Data: schedules
- DELETE /api/v1/schedules/:id
  - Contoh Response:
    ```json
    { "success": true, "data": { "message": "deleted" } }
    ```
  - Sumber Data: schedules

## QR Service
- POST /api/v1/qr/generate
  - Body JSON:
    - data (object, required)
    - type? (optional)
    - duration? (optional, menit)
  - Referensi: [service.go](file:///Users/ahmadnaufalmuzakki/Documents/KERJAAN/Meetsin.Id/2025/UNSRI/unsri%20app/backend-unsri-mobile/internal/qr/service/service.go#L23-L31)
  - Contoh Response:
    ```json
    { "success": true, "data": { "id": "qr-uuid", "qr": "base64...", "expires_at": "2026-01-13T10:00:00Z" } }
    ```
  - Sumber Data: generate via pkg/qrcode; metadata disusun di QRService
- POST /api/v1/qr/validate
  - Body JSON:
    - qr_data (required)
  - Contoh Response:
    ```json
    { "success": true, "data": { "valid": true, "payload": { "schedule_id": "sched-uuid" } } }
    ```
  - Sumber Data: validasi payload terenkode QR (QRService)
- GET /api/v1/qr/:id
  - Contoh Response:
    ```json
    { "success": true, "data": { "id": "qr-uuid", "qr": "base64..." } }
    ```
  - Sumber Data: penyimpanan metadata QR (jika persist), atau cache; mengikuti QRService
- POST /api/v1/qr/class/generate
  - Body JSON:
    - course_id, schedule_id (umum), duration? (opsional) tergantung implementasi kelas
  - Referensi: [handler.go](file:///Users/ahmadnaufalmuzakki/Documents/KERJAAN/Meetsin.Id/2025/UNSRI/unsri%20app/backend-unsri-mobile/internal/qr/handler/handler.go#L57-L80)
  - Contoh Response:
    ```json
    { "success": true, "data": { "schedule_id": "sched-uuid", "qr": "base64..." } }
    ```
  - Sumber Data: pkg/qrcode + jadwal kelas (schedules)
- GET /api/v1/qr/access/generate
  - No body; protected
  - Referensi: [routes.go](file:///Users/ahmadnaufalmuzakki/Documents/KERJAAN/Meetsin.Id/2025/UNSRI/unsri%20app/backend-unsri-mobile/internal/qr/handler/routes.go#L26-L33)
  - Contoh Response:
    ```json
    { "success": true, "data": { "session_id": "sess-uuid", "qr": "base64..." } }
    ```
  - Sumber Data: session akses sementara di QRService
- GET /api/v1/qr/access/validate/:session_id
  - No body; protected
  - Contoh Response:
    ```json
    { "success": true, "data": { "valid": true, "session_id": "sess-uuid" } }
    ```
  - Sumber Data: validasi session terhadap store QRService
- POST /api/v1/qr/gate/validate
  - Public
  - Body JSON:
    - qr_data (required)
  - Referensi: [service.go](file:///Users/ahmadnaufalmuzakki/Documents/KERJAAN/Meetsin.Id/2025/UNSRI/unsri%20app/backend-unsri-mobile/internal/qr/service/service.go#L466-L473)
  - Contoh Response:
    ```json
    { "success": true, "data": { "allowed": true, "reason": null } }
    ```
  - Sumber Data: validasi QR payload; logging dilakukan di Access service saat /access/log

## Location
- POST /api/v1/location/tap-in
  - Body JSON:
    - latitude (required)
    - longitude (required)
  - Referensi: [service.go](file:///Users/ahmadnaufalmuzakki/Documents/KERJAAN/Meetsin.Id/2025/UNSRI/unsri%20app/backend-unsri-mobile/internal/location/service/service.go#L21-L25)
  - Contoh Response:
    ```json
    { "success": true, "data": { "status": "IN", "time": "2026-01-13T07:05:00Z" } }
    ```
  - Sumber Data: log lokasi/kehadiran kampus (repository location/attendance)
- POST /api/v1/location/tap-out
  - Body JSON:
    - latitude (required)
    - longitude (required)
  - Contoh Response:
    ```json
    { "success": true, "data": { "status": "OUT", "time": "2026-01-13T17:00:00Z" } }
    ```
  - Sumber Data: log lokasi/kehadiran kampus
- GET /api/v1/location/check-in-status
  - Contoh Response:
    ```json
    { "success": true, "data": { "current": "IN", "since": "2026-01-13T07:05:00Z" } }
    ```
  - Sumber Data: status tap-in aktif dari repositori lokasi
- GET /api/v1/location/history
  - Query: page, per_page
  - Contoh Response:
    ```json
    { "success": true, "data": [ { "status": "IN", "time": "2026-01-13T07:05:00Z" } ], "meta": { "page": 1, "per_page": 20, "total": 10, "total_pages": 1 } }
    ```
  - Sumber Data: log lokasi
- GET /api/v1/location/geofences
  - Contoh Response:
    ```json
    { "success": true, "data": [ { "id": "geo-uuid", "name": "Kampus Indralaya", "radius": 100 } ] }
    ```
  - Sumber Data: geofences
- POST /api/v1/location/validate
  - Body JSON:
    - latitude, longitude (required)
  - Contoh Response:
    ```json
    { "success": true, "data": { "inside": true, "geofence_id": "geo-uuid" } }
    ```
  - Sumber Data: perhitungan lokasi terhadap geofence tersimpan
- POST /api/v1/location/geofences
  - Body JSON:
    - name (required)
    - description? (optional)
    - latitude, longitude, radius (required)
  - Contoh Response:
    ```json
    { "success": true, "data": { "id": "geo-uuid", "name": "Kampus Indralaya" } }
    ```
  - Sumber Data: geofences

## Broadcasts
- GET /api/v1/broadcasts
  - Query: type?, is_published?, page, per_page
  - Contoh Response:
    ```json
    { "success": true, "data": [ { "id": "bc-uuid", "title": "Pengumuman", "type": "general" } ] }
    ```
  - Sumber Data: broadcasts
- GET /api/v1/broadcasts/general
- GET /api/v1/broadcasts/class
- GET /api/v1/broadcasts/:id
  - Contoh Response:
    ```json
    { "success": true, "data": { "id": "bc-uuid", "title": "Pengumuman" } }
    ```
  - Sumber Data: broadcasts
- POST /api/v1/broadcasts
  - Body JSON:
    - title, content, type (required; oneof: general class campus)
    - priority?, class_id?, scheduled_at?, expires_at? (optional, RFC3339)
    - audiences? [{ user_id?, role?, prodi? }]
  - Referensi: [service.go](file:///Users/ahmadnaufalmuzakki/Documents/KERJAAN/Meetsin.Id/2025/UNSRI/unsri%20app/backend-unsri-mobile/internal/broadcast/service/service.go#L22-L39)
  - Contoh Response:
    ```json
    { "success": true, "data": { "id": "bc-uuid", "title": "Pengumuman" } }
    ```
  - Sumber Data: broadcasts
- PUT /api/v1/broadcasts/:id
  - Body JSON:
    - title?, content?, priority?, is_published?
  - Contoh Response:
    ```json
    { "success": true, "data": { "id": "bc-uuid", "is_published": true } }
    ```
  - Sumber Data: broadcasts
- DELETE /api/v1/broadcasts/:id
  - Contoh Response:
    ```json
    { "success": true, "data": { "message": "deleted" } }
    ```
  - Sumber Data: broadcasts
- POST /api/v1/broadcasts/:id/schedule
  - Body JSON:
    - scheduled_at (required, RFC3339)
  - Contoh Response:
    ```json
    { "success": true, "data": { "id": "bc-uuid", "scheduled_at": "2026-01-20T09:00:00Z" } }
    ```
  - Sumber Data: broadcasts
- POST /api/v1/broadcasts/search
  - Body JSON:
    - q (search query), page?, per_page? (mengikuti implementasi handler)
  - Contoh Response:
    ```json
    { "success": true, "data": [ { "id": "bc-uuid", "title": "Pengumuman" } ], "meta": { "page": 1, "per_page": 20, "total": 1, "total_pages": 1 } }
    ```
  - Sumber Data: broadcasts

## Search
- GET /api/v1/search
  - Query:
    - q (required), type (users|courses|schedules|...), role?, filters?, page, per_page
  - Referensi: [service.go](file:///Users/ahmadnaufalmuzakki/Documents/KERJAAN/Meetsin.Id/2025/UNSRI/unsri%20app/backend-unsri-mobile/internal/search/service/service.go#L20-L29)
  - Contoh Response:
    ```json
    { "success": true, "data": { "users": [ { "id": "user-uuid" } ], "courses": [ { "id": "course-uuid" } ], "schedules": [ { "id": "sched-uuid" } ] }, "meta": { "page": 1, "per_page": 20, "total": 3, "total_pages": 1 } }
    ```
  - Sumber Data: tabel terkait (users, courses, schedules) via SearchRepository (contoh: join schedules-courses di repository)
- GET /api/v1/search/global
  - Query:
    - q (required), types? (comma), limit?
  - Contoh Response:
    ```json
    { "success": true, "data": { "users": [ { "id": "user-uuid" } ] } }
    ```
  - Sumber Data: seperti di atas, dengan limit

## Reports
- GET /api/v1/reports/attendance
  - Query:
    - student_id?, course_id?, start_date (required), end_date (required), summary?
  - Referensi: [service.go](file:///Users/ahmadnaufalmuzakki/Documents/KERJAAN/Meetsin.Id/2025/UNSRI/unsri%20app/backend-unsri-mobile/internal/report/service/service.go#L21-L28)
  - Contoh Response:
    ```json
    { "success": true, "data": { "student_id": "student-uuid", "period": { "start": "2026-01-01", "end": "2026-01-31" }, "summary": { "hadir": 24, "izin": 2, "sakit": 1, "alpa": 3 }, "attendance_rate": 80.0 } }
    ```
  - Sumber Data: agregasi dari attendances; perhitungan rate seperti statistik
- GET /api/v1/reports/academic
  - Query:
    - student_id (required), semester?
  - Contoh Response:
    ```json
    { "success": true, "data": { "student_id": "student-uuid", "semester": "2025/2026-Ganjil", "ipk": 3.75, "courses": [ { "code": "IF101", "grade": "A" } ] } }
    ```
  - Sumber Data: transcript/KRS/grades (quick_actions related tables)
- GET /api/v1/reports/course
  - Query:
    - course_id (required), start_date (required), end_date (required)
  - Contoh Response:
    ```json
    { "success": true, "data": { "course_id": "course-uuid", "attendance": { "hadir": 80, "alpa": 20 } } }
    ```
  - Sumber Data: attendances per course (join schedules)
- GET /api/v1/reports/daily
  - Query:
    - date? (YYYY-MM-DD; default: today)
  - Contoh Response:
    ```json
    { "success": true, "data": { "date": "2026-01-13", "attendance_summary": { "hadir": 100, "izin": 5 }, "work_summary": { "CHECK_IN": 80 } } }
    ```
  - Sumber Data: agregasi harian dari attendances dan work_attendance_records

## Master Data
- Study Programs (/api/v1/study-programs)
  - GET: faculty?, is_active?, page, per_page
  - Contoh Response:
    ```json
    { "success": true, "data": [ { "id": "sp-uuid", "code": "IF", "name": "Informatika" } ] }
    ```
  - Sumber Data: study_programs
  - GET /:id
  - Contoh Response:
    ```json
    { "success": true, "data": { "id": "sp-uuid", "code": "IF" } }
    ```
  - Sumber Data: study_programs
  - POST: code, name (required); name_en?, faculty?, degree_level?, accreditation? (optional)
  - Contoh Response:
    ```json
    { "success": true, "data": { "id": "sp-uuid", "code": "IF" } }
    ```
  - Sumber Data: study_programs
  - PUT /:id: name?, name_en?, faculty?, degree_level?, accreditation?, is_active?
  - DELETE /:id
  - Contoh Response (DELETE):
    ```json
    { "success": true, "data": { "message": "deleted" } }
    ```
  - Sumber Data: study_programs
- Academic Periods (/api/v1/academic-periods)
  - GET: academic_year?, semester_type?, is_active?, page, per_page
  - Contoh Response:
    ```json
    { "success": true, "data": [ { "id": "ap-uuid", "academic_year": "2025/2026", "semester_type": "Ganjil" } ] }
    ```
  - Sumber Data: academic_periods
  - GET /active
  - Contoh Response:
    ```json
    { "success": true, "data": { "id": "ap-uuid", "academic_year": "2025/2026", "semester_type": "Ganjil" } }
    ```
  - Sumber Data: academic_periods
  - GET /:id
  - POST: code, name, academic_year, semester_type, start_date, end_date (required); registration_start?, registration_end?, is_active?
  - PUT /:id: name?, academic_year?, semester_type?, start_date?, end_date?, registration_start?, registration_end?, is_active?
  - DELETE /:id
- Rooms (/api/v1/rooms)
  - GET: building?, room_type?, is_active?, page, per_page
  - Contoh Response:
    ```json
    { "success": true, "data": [ { "id": "room-uuid", "code": "A101", "name": "Ruang Kuliah A101" } ] }
    ```
  - Sumber Data: rooms
  - GET /:id
  - POST: code, name (required); building?, floor?, capacity?, room_type?, facilities?
  - PUT /:id: name?, building?, floor?, capacity?, room_type?, facilities?, is_active?
  - DELETE /:id

## Access (Gate)
- POST /api/v1/access/qr/validate
  - Body JSON:
    - qr_token (required)
    - gate_id (required)
  - Referensi: [service.go](file:///Users/ahmadnaufalmuzakki/Documents/KERJAAN/Meetsin.Id/2025/UNSRI/unsri%20app/backend-unsri-mobile/internal/access/service/service.go#L22-L26)
  - Contoh Response:
    ```json
    { "success": true, "data": { "allowed": true, "reason": null } }
    ```
  - Sumber Data: validasi QR dan kebijakan akses; log dicatat di access_logs
- GET /api/v1/access/history
  - Query:
    - user_id?, gate_id?, page, per_page
  - Contoh Response:
    ```json
    { "success": true, "data": [ { "user_id": "user-uuid", "gate_id": "GATE-1", "access_type": "entry", "is_allowed": true } ] }
    ```
  - Sumber Data: access_logs
- POST /api/v1/access/log
  - Body JSON:
    - user_id (required)
    - gate_id (required)
    - access_type (required, oneof: entry exit)
    - is_allowed (boolean)
    - reason? (optional)
    - qr_code_id? (optional)
  - Contoh Response:
    ```json
    { "success": true, "data": { "id": "access-uuid", "is_allowed": true } }
    ```
  - Sumber Data: access_logs
- GET /api/v1/access/permissions/:userId
- POST /api/v1/access/permissions
  - Body JSON:
    - user_id (required)
    - gate_id (required)
    - is_allowed (boolean)
    - valid_from?, valid_until? (RFC3339)
  - Contoh Response:
    ```json
    { "success": true, "data": { "user_id": "user-uuid", "gate_id": "GATE-1", "is_allowed": true } }
    ```
  - Sumber Data: tabel permissions (impl. Access service)

## Error Responses
- Format umum:
```json
{
  "success": false,
  "error": {
    "code": "ERROR_CODE",
    "message": "Error message"
  }
}
```

Catatan: Rujukan langsung ke implementasi handler/service disertakan agar tim dapat menelusuri sumber kode dengan cepat saat perlu konfirmasi detail.

- `UNAUTHORIZED` - Token tidak valid atau expired
- `FORBIDDEN` - Tidak memiliki permission
- `NOT_FOUND` - Resource tidak ditemukan
- `BAD_REQUEST` - Request tidak valid
- `VALIDATION_FAILED` - Validasi gagal
- `CONFLICT` - Resource conflict
- `INTERNAL_ERROR` - Server error

## Rate Limiting

API memiliki rate limiting:
- **Default**: 100 requests per minute per IP
- **Authenticated**: 1000 requests per minute per user

## Pagination

Endpoints yang support pagination menggunakan query parameters:
- `page` - Page number (default: 1)
- `per_page` - Items per page (default: 20, max: 100)

Response format:
```json
{
  "success": true,
  "data": [...],
  "meta": {
    "page": 1,
    "per_page": 20,
    "total": 100,
    "total_pages": 5
  }
}
```

## Status Codes
 
 - `200 OK` - Request berhasil
 - `201 Created` - Resource berhasil dibuat
 - `400 Bad Request` - Request tidak valid
 - `401 Unauthorized` - Tidak terautentikasi
 - `403 Forbidden` - Tidak memiliki permission
 - `404 Not Found` - Resource tidak ditemukan
 - `409 Conflict` - Resource conflict
 - `500 Internal Server Error` - Server error

## Notifications
- POST /api/v1/notifications/send
  - Body JSON:
    - user_id (required)
    - title (required)
    - message (required)
    - type (required, oneof: info warning error success)
    - data? (optional)
  - Referensi: [service.go](file:///Users/ahmadnaufalmuzakki/Documents/KERJAAN/Meetsin.Id/2025/UNSRI/unsri%20app/backend-unsri-mobile/internal/notification/service/service.go#L22-L29)
  - Contoh Response:
    ```json
    { "success": true, "data": { "id": "notif-uuid", "user_id": "user-uuid", "title": "Info", "is_read": false } }
    ```
  - Sumber Data: notifications; push FCM (TODO)
- GET /api/v1/notifications
  - Query:
    - is_read? (boolean), page?, per_page?
  - Contoh Response:
    ```json
    { "success": true, "data": [ { "id": "notif-uuid", "title": "Info", "is_read": false } ], "meta": { "page": 1, "per_page": 20, "total": 5, "total_pages": 1 } }
    ```
  - Sumber Data: notifications
- PUT /api/v1/notifications/:id/read
  - Contoh Response:
    ```json
    { "success": true, "data": { "message": "Notification marked as read" } }
    ```
  - Sumber Data: notifications (flag is_read, read_at)
- PUT /api/v1/notifications/read-all
  - Contoh Response:
    ```json
    { "success": true, "data": { "message": "All notifications marked as read" } }
    ```
  - Sumber Data: notifications
- POST /api/v1/notifications/register-device
  - Body JSON:
    - token (required)
    - platform (required, oneof: ios android web)
  - Referensi: [service.go](file:///Users/ahmadnaufalmuzakki/Documents/KERJAAN/Meetsin.Id/2025/UNSRI/unsri%20app/backend-unsri-mobile/internal/notification/service/service.go#L100-L104)
  - Contoh Response:
    ```json
    { "success": true, "data": { "id": "device-uuid", "token": "fcm-token", "platform": "android" } }
    ```
  - Sumber Data: device_tokens
- DELETE /api/v1/notifications/device/:token
  - Contoh Response:
    ```json
    { "success": true, "data": { "message": "Device token unregistered" } }
    ```
  - Sumber Data: device_tokens

## File Storage
- POST /api/v1/files/upload
  - Body multipart/form-data:
    - file (file) required
    - file_type? (document|avatar|selfie; default: document)
    - is_public? (boolean; default: false)
  - Referensi: [handler.go](file:///Users/ahmadnaufalmuzakki/Documents/KERJAAN/Meetsin.Id/2025/UNSRI/unsri%20app/backend-unsri-mobile/internal/file-storage/handler/handler.go#L26-L53)
  - Contoh Response:
    ```json
    { "success": true, "data": { "id": "file-uuid", "url": "https://...", "file_type": "document" } }
    ```
  - Sumber Data: files
- GET /api/v1/files
  - Query: page?, per_page?
  - Contoh Response:
    ```json
    { "success": true, "data": [ { "id": "file-uuid", "file_name": "doc.pdf" } ], "meta": { "page": 1, "per_page": 20, "total": 2, "total_pages": 1 } }
    ```
  - Sumber Data: files
- GET /api/v1/files/:id
  - Contoh Response:
    ```json
    { "success": true, "data": { "id": "file-uuid", "url": "https://..." } }
    ```
  - Sumber Data: files
- GET /api/v1/files/:id/download
  - Contoh Response:
    Binary content dengan MIME type sesuai file
  - Sumber Data: files (content dari storage)
- DELETE /api/v1/files/:id
  - Contoh Response:
    ```json
    { "success": true, "data": { "message": "File deleted successfully" } }
    ```
  - Sumber Data: files
- POST /api/v1/files/avatar
  - Body multipart/form-data:
    - file (file) required
  - Contoh Response:
    ```json
    { "success": true, "data": { "id": "file-uuid", "file_type": "avatar" } }
    ```
  - Sumber Data: files
- POST /api/v1/files/document
  - Body multipart/form-data:
    - file (file) required
    - is_public? (boolean; default: false)
  - Contoh Response:
    ```json
    { "success": true, "data": { "id": "file-uuid", "file_type": "document" } }
    ```
  - Sumber Data: files

## Quick Actions
- GET /api/v1/quick-actions
  - Contoh Response:
    ```json
    { "success": true, "data": [ { "key": "transcript", "label": "Lihat Transkrip" } ] }
    ```
  - Sumber Data: daftar statis/dinamis dari QuickActionsService
- GET /api/v1/quick-actions/transcript/:studentId
- GET /api/v1/quick-actions/krs/:studentId
  - Query: semester?
- GET /api/v1/quick-actions/bimbingan
  - Query: limit? (default: 10)
  - Referensi: [handler.go](file:///Users/ahmadnaufalmuzakki/Documents/KERJAAN/Meetsin.Id/2025/UNSRI/unsri%20app/backend-unsri-mobile/internal/quick-actions/handler/handler.go#L72-L91)
  - Contoh Response (bimbingan):
    ```json
    { "success": true, "data": [ { "id": "bimb-uuid", "dosen_id": "dosen-uuid", "topic": "TA", "date": "2026-01-10" } ] }
    ```
  - Sumber Data: bimbingans

## Calendar
- GET /api/v1/calendar/events
- GET /api/v1/calendar/events/upcoming
- GET /api/v1/calendar/events/month/:year/:month
- GET /api/v1/calendar/events/:id
- POST /api/v1/calendar/events
  - Body JSON:
    - title, description?, start_at, end_at (required; RFC3339)
    - location?, audience? (optional)
- PUT /api/v1/calendar/events/:id
- DELETE /api/v1/calendar/events/:id
  - Referensi: [routes.go](file:///Users/ahmadnaufalmuzakki/Documents/KERJAAN/Meetsin.Id/2025/UNSRI/unsri%20app/backend-unsri-mobile/internal/calendar/handler/routes.go#L12-L24)
  - Contoh Response:
    ```json
    { "success": true, "data": [ { "id": "event-uuid", "title": "UAS", "start_date": "2026-01-25", "end_date": "2026-01-28" } ] }
    ```
  - Sumber Data: academic_events
