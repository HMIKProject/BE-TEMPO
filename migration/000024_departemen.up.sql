CREATE TABLE IF NOT EXISTS departemen (
    id_departemen BIGSERIAL PRIMARY KEY,
    nama_departemen VARCHAR(255) NOT NULL,
    deskripsi TEXT,
    logo VARCHAR(255),
    dibuat_pada TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    diperbarui_pada TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
