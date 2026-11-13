-- 000038_add_nama_lengkap_and_foto_to_pengguna.up.sql
ALTER TABLE pengguna 
ADD COLUMN IF NOT EXISTS nama_lengkap VARCHAR(255) NOT NULL DEFAULT '',
ADD COLUMN IF NOT EXISTS foto_pengguna VARCHAR(255) DEFAULT '';
