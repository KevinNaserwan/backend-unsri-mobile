-- Fix NIP column issue: Remove duplicate n_ip column from dosen and staff tables
-- This migration removes the incorrectly named n_ip column that was created by GORM
-- The correct column name is 'nip' (already exists)
-- 
-- Issue: GORM was converting NIP to n_i_p (snake_case), creating duplicate columns
-- Solution: Explicitly specify column name in GORM tags: gorm:"column:nip"

-- Fix dosen table: Drop constraints and indexes first, then column
ALTER TABLE dosen DROP CONSTRAINT IF EXISTS dosen_n_ip_key;
DROP INDEX IF EXISTS idx_dosen_n_ip;
ALTER TABLE dosen DROP COLUMN IF EXISTS n_ip;

-- Fix staff table: Drop constraints and indexes first, then column
ALTER TABLE staff DROP CONSTRAINT IF EXISTS staff_n_ip_key;
DROP INDEX IF EXISTS idx_staff_n_ip;
ALTER TABLE staff DROP COLUMN IF EXISTS n_ip;

-- Fix mahasiswa table (for consistency, check n_i_m column if exists)
ALTER TABLE mahasiswa DROP CONSTRAINT IF EXISTS mahasiswa_n_i_m_key;
DROP INDEX IF EXISTS idx_mahasiswa_n_i_m;
ALTER TABLE mahasiswa DROP COLUMN IF EXISTS n_i_m;
