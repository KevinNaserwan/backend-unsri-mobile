-- Remove selfie_url from work_attendance_records table
ALTER TABLE work_attendance_records DROP COLUMN IF EXISTS selfie_url;

-- Remove selfie_url from attendances table
ALTER TABLE attendances DROP COLUMN IF EXISTS selfie_url;

-- Remove profile_photo_url from users table
ALTER TABLE users DROP COLUMN IF EXISTS profile_photo_url;

