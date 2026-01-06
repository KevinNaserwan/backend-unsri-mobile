-- Add profile_photo_url to users table
ALTER TABLE users ADD COLUMN IF NOT EXISTS profile_photo_url VARCHAR(500);

-- Add selfie_url to attendances table
ALTER TABLE attendances ADD COLUMN IF NOT EXISTS selfie_url VARCHAR(500);

-- Add selfie_url to work_attendance_records table
ALTER TABLE work_attendance_records ADD COLUMN IF NOT EXISTS selfie_url VARCHAR(500);

