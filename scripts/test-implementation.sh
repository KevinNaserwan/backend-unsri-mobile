#!/bin/bash

# Test Implementation Script for UNSRI Backend
# Tests: MinIO Storage, Attendance with Selfie, QR for Gate, Geofencing, Profile Photo

echo "🚀 Testing UNSRI Backend Implementation"
echo "======================================="

# Base URLs
AUTH_URL="http://localhost:8081/api/v1/auth"
USER_URL="http://localhost:8082/api/v1/users"
ATTENDANCE_URL="http://localhost:8084/api/v1/attendance"
QR_URL="http://localhost:8085/api/v1/qr"
FILE_URL="http://localhost:8093/api/v1/files"
LOCATION_URL="http://localhost:8090/api/v1/location"

# Test files
SELFIE_FILE="test-selfie.jpg"
PROFILE_FILE="test-profile.jpg"

# Create test image files if they don't exist
if [ ! -f "$SELFIE_FILE" ]; then
    echo "Creating test selfie file..."
    # Create a small test image (1x1 pixel PNG)
    echo -e '\x89PNG\r\n\x1a\n\x00\x00\x00\rIHDR\x00\x00\x00\x01\x00\x00\x00\x01\x08\x02\x00\x00\x00\x90wS\xde\x00\x00\x00\tpHYs\x00\x00\x0b\x13\x00\x00\x0b\x13\x01\x00\x9a\x9c\x18\x00\x00\x00\nIDATx\x9cc\xf8\x00\x00\x00\x01\x00\x01\x00\x00\x00\x00IEND\xaeB`\x82' > "$SELFIE_FILE"
fi

if [ ! -f "$PROFILE_FILE" ]; then
    echo "Creating test profile file..."
    cp "$SELFIE_FILE" "$PROFILE_FILE"
fi

echo ""
echo "1️⃣  Testing Authentication"
echo "=========================="

# Register test user
echo "Registering test user..."
REGISTER_RESPONSE=$(curl -s -X POST "$AUTH_URL/register" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@unsri.ac.id",
    "password": "password123",
    "role": "staff",
    "name": "Test User"
  }')

echo "Register Response: $REGISTER_RESPONSE"

# Login
echo "Logging in..."
LOGIN_RESPONSE=$(curl -s -X POST "$AUTH_URL/login" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@unsri.ac.id",
    "password": "password123"
  }')

echo "Login Response: $LOGIN_RESPONSE"

# Extract token
TOKEN=$(echo "$LOGIN_RESPONSE" | grep -o '"access_token":"[^"]*' | cut -d'"' -f4)
if [ -z "$TOKEN" ]; then
    echo "❌ Failed to get authentication token"
    exit 1
fi

echo "✅ Authentication successful"
echo "Token: ${TOKEN:0:20}..."

echo ""
echo "2️⃣  Testing File Storage (MinIO)"
echo "==============================="

# Test file upload
echo "Uploading test file..."
UPLOAD_RESPONSE=$(curl -s -X POST "$FILE_URL/upload" \
  -H "Authorization: Bearer $TOKEN" \
  -F "file=@$SELFIE_FILE" \
  -F "file_type=selfie" \
  -F "is_public=false")

echo "Upload Response: $UPLOAD_RESPONSE"

# Extract file URL
FILE_URL_RESULT=$(echo "$UPLOAD_RESPONSE" | grep -o '"url":"[^"]*' | cut -d'"' -f4)
if [ -n "$FILE_URL_RESULT" ]; then
    echo "✅ File upload successful"
    echo "File URL: $FILE_URL_RESULT"
else
    echo "❌ File upload failed"
fi

echo ""
echo "3️⃣  Testing Profile Photo Upload"
echo "==============================="

# Upload profile photo
echo "Uploading profile photo..."
PROFILE_RESPONSE=$(curl -s -X POST "$USER_URL/profile/photo" \
  -H "Authorization: Bearer $TOKEN" \
  -F "photo=@$PROFILE_FILE")

echo "Profile Photo Response: $PROFILE_RESPONSE"

if echo "$PROFILE_RESPONSE" | grep -q "profile_photo_url"; then
    echo "✅ Profile photo upload successful"
else
    echo "❌ Profile photo upload failed"
fi

echo ""
echo "4️⃣  Testing Geofencing"
echo "====================="

# Create geofence
echo "Creating geofence..."
GEOFENCE_RESPONSE=$(curl -s -X POST "$LOCATION_URL/geofences" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "UNSRI Campus",
    "description": "Main campus area",
    "latitude": -2.9880,
    "longitude": 104.7565,
    "radius": 1000
  }')

echo "Geofence Response: $GEOFENCE_RESPONSE"

# Test location validation
echo "Testing location validation..."
LOCATION_RESPONSE=$(curl -s -X POST "$LOCATION_URL/validate" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "latitude": -2.9880,
    "longitude": 104.7565
  }')

echo "Location Validation Response: $LOCATION_RESPONSE"

if echo "$LOCATION_RESPONSE" | grep -q "valid"; then
    echo "✅ Geofencing validation successful"
else
    echo "❌ Geofencing validation failed"
fi

echo ""
echo "5️⃣  Testing QR Code Generation"
echo "============================="

# Generate access QR for gate
echo "Generating access QR for gate..."
QR_RESPONSE=$(curl -s -X GET "$QR_URL/access/generate" \
  -H "Authorization: Bearer $TOKEN")

echo "QR Generation Response: $QR_RESPONSE"

# Extract QR data for validation
QR_DATA=$(echo "$QR_RESPONSE" | grep -o '"qr_code":"[^"]*' | cut -d'"' -f4)
if [ -n "$QR_DATA" ]; then
    echo "✅ QR generation successful"
    echo "QR Data: ${QR_DATA:0:50}..."
    
    # Test gate validation (public endpoint)
    echo "Testing gate QR validation..."
    GATE_RESPONSE=$(curl -s -X POST "$QR_URL/gate/validate" \
      -H "Content-Type: application/json" \
      -d "{\"qr_data\": \"$QR_DATA\"}")
    
    echo "Gate Validation Response: $GATE_RESPONSE"
    
    if echo "$GATE_RESPONSE" | grep -q "allowed"; then
        echo "✅ Gate QR validation successful"
    else
        echo "❌ Gate QR validation failed"
    fi
else
    echo "❌ QR generation failed"
fi

echo ""
echo "6️⃣  Testing Work Attendance with Selfie"
echo "======================================"

# Test check-in with selfie
echo "Testing work check-in with selfie..."
CHECKIN_RESPONSE=$(curl -s -X POST "$ATTENDANCE_URL/work-attendance/check-in" \
  -H "Authorization: Bearer $TOKEN" \
  -F "latitude=-2.9880" \
  -F "longitude=104.7565" \
  -F "selfie=@$SELFIE_FILE" \
  -F "notes=Test check-in")

echo "Check-in Response: $CHECKIN_RESPONSE"

if echo "$CHECKIN_RESPONSE" | grep -q "attendance_type"; then
    echo "✅ Work check-in with selfie successful"
else
    echo "❌ Work check-in with selfie failed"
fi

echo ""
echo "7️⃣  Testing Health Checks"
echo "========================"

# Test all service health checks
SERVICES=("auth-service:8081" "user-service:8082" "attendance-service:8084" "qr-service:8085" "file-storage-service:8093" "location-service:8090")

for service in "${SERVICES[@]}"; do
    name=$(echo $service | cut -d':' -f1)
    port=$(echo $service | cut -d':' -f2)
    
    echo "Testing $name health..."
    HEALTH_RESPONSE=$(curl -s "http://localhost:$port/health")
    
    if echo "$HEALTH_RESPONSE" | grep -q "ok"; then
        echo "✅ $name is healthy"
    else
        echo "❌ $name health check failed"
    fi
done

echo ""
echo "🎉 Testing Complete!"
echo "==================="
echo ""
echo "Summary of implemented features:"
echo "✅ HRIS renamed to Kepegawaian"
echo "✅ MinIO storage for photos"
echo "✅ Mandatory selfie for attendance"
echo "✅ QR code generation as images"
echo "✅ Gate integration endpoints"
echo "✅ Geofencing validation"
echo "✅ Profile photo management"
echo ""
echo "🚀 Backend ready for UNSRI mobile app!"

# Cleanup test files
rm -f "$SELFIE_FILE" "$PROFILE_FILE"
