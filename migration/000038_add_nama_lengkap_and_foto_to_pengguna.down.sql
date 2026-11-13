-- 000038_add_nama_lengkap_and_foto_to_pengguna.down.sql
ALTER TABLE pengguna 
DROP COLUMN IF EXISTS nama_lengkap,
DROP COLUMN IF EXISTS foto_pengguna;
