-- 000037_add_activation_and_admin_fields.up.sql

-- 1. Penambahan kolom keamanan dan aktivasi akun pada tabel pengguna
ALTER TABLE pengguna 
ADD COLUMN IF NOT EXISTS token_version INT NOT NULL DEFAULT 1,
ADD COLUMN IF NOT EXISTS kode_aktivasi VARCHAR(50) DEFAULT NULL,
ADD COLUMN IF NOT EXISTS kode_aktivasi_expired TIMESTAMP DEFAULT NULL,
ADD COLUMN IF NOT EXISTS is_first_login BOOLEAN NOT NULL DEFAULT FALSE;

-- 2. Penambahan kolom periode dan berkas CV pada tabel anggota
ALTER TABLE anggota
ADD COLUMN IF NOT EXISTS periode VARCHAR(20) DEFAULT NULL,
ADD COLUMN IF NOT EXISTS cv_url VARCHAR(255) DEFAULT NULL;

-- 3. Penambahan kolom status_aktif untuk soft-delete / arsip departemen dan divisi
ALTER TABLE departemen 
ADD COLUMN IF NOT EXISTS status_aktif BOOLEAN NOT NULL DEFAULT TRUE;

ALTER TABLE divisi_departemen 
ADD COLUMN IF NOT EXISTS status_aktif BOOLEAN NOT NULL DEFAULT TRUE;
