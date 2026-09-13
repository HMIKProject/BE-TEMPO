ALTER TABLE anggota
ADD COLUMN IF NOT EXISTS id_departemen BIGINT REFERENCES departemen(id_departemen) ON DELETE SET NULL;
