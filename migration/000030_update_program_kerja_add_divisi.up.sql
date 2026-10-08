ALTER TABLE program_kerja
ADD COLUMN id_divisi BIGINT NULL,
ADD CONSTRAINT fk_program_kerja_divisi FOREIGN KEY (id_divisi) REFERENCES divisi_departemen (id_divisi) ON DELETE SET NULL;
