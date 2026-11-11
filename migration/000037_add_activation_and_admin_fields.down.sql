-- 000037_add_activation_and_admin_fields.down.sql

ALTER TABLE divisi_departemen DROP COLUMN IF EXISTS status_aktif;
ALTER TABLE departemen DROP COLUMN IF EXISTS status_aktif;
ALTER TABLE anggota DROP COLUMN IF EXISTS cv_url, DROP COLUMN IF EXISTS periode;
ALTER TABLE pengguna 
DROP COLUMN IF EXISTS is_first_login,
DROP COLUMN IF EXISTS kode_aktivasi_expired,
DROP COLUMN IF EXISTS kode_aktivasi,
DROP COLUMN IF EXISTS token_version;
