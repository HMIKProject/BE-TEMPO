CREATE TABLE IF NOT EXISTS divisi_departemen (
    id_divisi BIGSERIAL PRIMARY KEY,
    id_departemen BIGINT NOT NULL,
    nama_divisi VARCHAR(255) NOT NULL,
    deskripsi TEXT,
    dibuat_pada TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    diperbarui_pada TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_divisi_departemen FOREIGN KEY (id_departemen) REFERENCES departemen (id_departemen) ON DELETE CASCADE
);
