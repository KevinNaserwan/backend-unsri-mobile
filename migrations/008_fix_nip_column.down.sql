-- Rollback: Re-add n_ip columns (not recommended, but for migration rollback support)
-- Note: This is only for migration rollback. The n_ip column should not be used.

-- Rollback dosen table
-- ALTER TABLE dosen ADD COLUMN n_ip TEXT;
-- CREATE UNIQUE INDEX IF NOT EXISTS dosen_n_ip_key ON dosen(n_ip);
-- CREATE INDEX IF NOT EXISTS idx_dosen_n_ip ON dosen(n_ip);

-- Rollback staff table
-- ALTER TABLE staff ADD COLUMN n_ip TEXT;
-- CREATE UNIQUE INDEX IF NOT EXISTS staff_n_ip_key ON staff(n_ip);
-- CREATE INDEX IF NOT EXISTS idx_staff_n_ip ON staff(n_ip);

-- Rollback mahasiswa table
-- ALTER TABLE mahasiswa ADD COLUMN n_i_m TEXT;
-- CREATE UNIQUE INDEX IF NOT EXISTS mahasiswa_n_i_m_key ON mahasiswa(n_i_m);
-- CREATE INDEX IF NOT EXISTS idx_mahasiswa_n_i_m ON mahasiswa(n_i_m);

-- Note: This rollback is commented out because n_ip/n_i_m columns are incorrect
-- and should not be recreated. If rollback is needed, uncomment above lines.
