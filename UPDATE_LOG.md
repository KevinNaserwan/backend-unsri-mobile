# Update Log

Log perubahan dan update pada project UNSRI Backend Mobile.

---

## Update Log - 2025-11-27

**Author:** anop  
**Date:** 2025-11-27  
**Version:** Development

### 📦 Modules Development

#### 1. Master Data Service (Port: 8096)
**Status:** ✅ Completed

**Features:**
- Study Program management (CRUD)
- Academic Period management dengan validasi active period
- Room management dengan capacity tracking
- Full REST API dengan authentication & authorization
- Integration dengan API Gateway
- Docker configuration

**Files Created:**
- `internal/master-data/config/config.go`
- `internal/master-data/repository/repository.go`
- `internal/master-data/service/service.go`
- `internal/master-data/handler/handler.go`
- `internal/master-data/handler/routes.go`
- `internal/master-data/middleware/auth.go`
- `internal/shared/models/master_data.go`
- `cmd/master-data-service/main.go`
- `deployments/docker/Dockerfile.master-data-service`
- `internal/master-data/service/service_test.go`

**API Endpoints:**
- `GET /api/v1/study-programs` - List study programs
- `GET /api/v1/study-programs/:id` - Get by ID
- `POST /api/v1/study-programs` - Create
- `PUT /api/v1/study-programs/:id` - Update
- `DELETE /api/v1/study-programs/:id` - Delete
- Similar endpoints for academic-periods and rooms

---

#### 2. Enrollment Module
**Status:** ✅ Completed

**Features:**
- Student enrollment workflow
- Enrollment status management (PENDING → APPROVED/REJECTED)
- Grade and score tracking
- Enrollment by class filtering
- Full CRUD operations

**Files Modified:**
- `internal/shared/models/course.go` - Updated Enrollment model
- `internal/course/service/service.go` - Added enrollment methods
- `internal/course/repository/repository.go` - Added enrollment repository methods
- `internal/course/handler/handler.go` - Added enrollment handlers
- `internal/course/handler/routes.go` - Added enrollment routes
- `internal/course/service/service_test.go` - Added enrollment tests

**API Endpoints:**
- `GET /api/v1/enrollments` - List enrollments
- `GET /api/v1/enrollments/:id` - Get by ID
- `GET /api/v1/classes/:id/enrollments` - Get by class
- `POST /api/v1/enrollments` - Create enrollment
- `PUT /api/v1/enrollments/:id/status` - Update status
- `PUT /api/v1/enrollments/:id/grade` - Update grade
- `DELETE /api/v1/enrollments/:id` - Delete enrollment

---

#### 3. Work Attendance Module (HRIS)
**Status:** ✅ Completed

**Features:**
- Shift Pattern management
- User Shift assignment
- Work Schedule management
- Check-in/Check-out dengan validasi
- Auto-detect late in & early out
- Location tracking
- WiFi validation
- Work Attendance Records

**Files Created:**
- `internal/shared/models/work_attendance.go`
- `internal/attendance/service/service.go` - Added work attendance methods
- `internal/attendance/repository/repository.go` - Added work attendance repository methods
- `internal/attendance/handler/handler.go` - Added work attendance handlers
- `internal/attendance/handler/routes.go` - Added work attendance routes
- `internal/attendance/service/service_test.go` - Added work attendance tests

**API Endpoints:**
- `POST /api/v1/work-attendance/check-in` - Check-in
- `POST /api/v1/work-attendance/check-out` - Check-out
- `GET /api/v1/work-attendance/records` - Get records
- `GET /api/v1/work-attendance/shifts` - List shift patterns
- `POST /api/v1/work-attendance/shifts` - Create shift pattern
- `GET /api/v1/work-attendance/user-shifts/:userId` - Get user shifts
- `POST /api/v1/work-attendance/user-shifts` - Assign shift to user
- `GET /api/v1/work-attendance/schedules` - List work schedules
- `POST /api/v1/work-attendance/schedules` - Create work schedule

---

#### 4. Leave Management Module (HRIS) - Port: 8097
**Status:** ✅ Completed

**Features:**
- Leave Request management
- Leave Quota management
- Leave types: Annual, Sick, Personal, Emergency, Unpaid, Other
- Approval workflow
- Automatic quota validation
- Quota calculation

**Files Created:**
- `internal/leave/config/config.go`
- `internal/leave/repository/repository.go`
- `internal/leave/service/service.go`
- `internal/leave/handler/handler.go`
- `internal/leave/handler/routes.go`
- `internal/leave/middleware/auth.go`
- `internal/shared/models/leave.go`
- `cmd/leave-service/main.go`
- `deployments/docker/Dockerfile.leave-service`
- `internal/leave/service/service_test.go`

**API Endpoints:**
- `POST /api/v1/leave-requests` - Create leave request
- `GET /api/v1/leave-requests` - List leave requests
- `GET /api/v1/leave-requests/:id` - Get by ID
- `PUT /api/v1/leave-requests/:id/approve` - Approve
- `PUT /api/v1/leave-requests/:id/reject` - Reject
- `PUT /api/v1/leave-requests/:id/cancel` - Cancel
- `GET /api/v1/leave-quotas` - List quotas
- `POST /api/v1/leave-quotas` - Create quota
- `GET /api/v1/leave-quotas/by-user/:userId` - Get by user

---

### 🧪 Testing

**Unit Tests Created:**
- ✅ Master Data Service tests
- ✅ Leave Management Service tests
- ✅ Auth Service tests
- ✅ Course Service tests
- ✅ Attendance Service tests
- ✅ User Service tests
- ✅ Schedule Service tests
- ✅ QR Service tests
- ✅ Broadcast Service tests
- ✅ Notification Service tests
- ✅ Calendar Service tests
- ✅ Location Service tests
- ✅ Access Service tests
- ✅ Quick Actions Service tests
- ✅ File Storage Service tests
- ✅ Search Service tests
- ✅ Report Service tests

**Test Results:**
- Total Services Tested: 17
- All Tests Passing: ✅ 17/17
- Test Coverage: Configured
- Linter Errors: 0 (internal code)

---

### 🔄 CI/CD

**GitHub Actions Workflows:**
- ✅ CI Workflow - Automated testing, linting, building
- ✅ CD Workflow - Automated deployment to staging/production
- ✅ Test Workflow - Parallel testing with matrix strategy

**Configuration:**
- ✅ PostgreSQL service for integration tests
- ✅ Go module caching
- ✅ Test coverage reporting
- ✅ Docker image building & pushing
- ✅ Artifact uploads

---

### 📚 Documentation

**Files Created:**
- ✅ `WEB_ADMIN_DESIGN_GUIDE.md` - Design guide for web admin
- ✅ `README_TESTING.md` - Testing guide
- ✅ `.github/workflows/README.md` - CI/CD documentation
- ✅ `CHANGELOG.md` - This file
- ✅ `UPDATE_LOG.md` - Detailed update log

---

### 🔧 Configuration

**VS Code:**
- ✅ `.vscode/settings.json` - YAML validation config
- ✅ `.vscode/extensions.json` - Recommended extensions
- ✅ `.github/workflows/.vscode/settings.json` - Suppress false positives

**Other:**
- ✅ `.yaml-lint.yml` - YAML linting config
- ✅ Updated `.gitignore` - Coverage files

---

### 🐳 Docker & Deployment

**Docker:**
- ✅ `Dockerfile.leave-service` - Leave service container
- ✅ Updated `docker-compose.yml` - Added leave-service
- ✅ Updated API Gateway environment variables

**Integration:**
- ✅ API Gateway routes for all new services
- ✅ Proxy handlers updated
- ✅ Service dependencies configured

---

### 📊 Statistics

**Code Metrics:**
- New Services: 2 (master-data, leave)
- Enhanced Services: 2 (course, attendance)
- New Models: 8
- New API Endpoints: 50+
- Test Files: 17
- Lines of Test Code: ~2000+

**Module Completion:**
- Master Data: 100% ✅
- Enrollment: 100% ✅
- Work Attendance: 100% ✅
- Leave Management: 100% ✅
- Unit Testing: 100% ✅
- CI/CD: 100% ✅

---

### 🔐 Security

- ✅ JWT authentication implemented
- ✅ Role-based access control (RBAC)
- ✅ Middleware for authorization
- ✅ Input validation
- ✅ Error handling

---

### ⚠️ Known Issues

1. **VS Code Linter Warnings** (False Positive)
   - Error "Unable to resolve action" di VS Code adalah normal
   - Linter lokal tidak bisa akses GitHub Actions marketplace
   - Workflow akan bekerja dengan baik di GitHub
   - **Status:** Tidak perlu diperbaiki, normal behavior

2. **Secrets Warning** (False Positive)
   - Warning "Context access might be invalid: DOCKER_USERNAME" adalah normal
   - Secrets akan di-set di GitHub repository settings
   - **Status:** Tidak perlu diperbaiki, normal behavior

---

### 📝 Notes

- Semua module development mengikuti best practices
- Code structure konsisten dengan existing services
- Error handling menggunakan custom error types
- Validation menggunakan GORM tags dan business logic
- Soft delete menggunakan GORM DeletedAt
- Pagination dan filtering diimplementasikan

---

### 👤 Author

**anop** - 2025-11-27
- Module development
- Unit testing
- CI/CD setup
- Documentation

---

## Next Steps

1. Integration testing di staging environment
2. Performance testing
3. Security audit
4. API documentation (Swagger)
5. Load testing
6. Monitoring setup

---

*Last Updated: 2025-11-27 by anop*

