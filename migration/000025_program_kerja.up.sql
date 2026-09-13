CREATE TABLE IF NOT EXISTS program_kerja (
    id_proker BIGSERIAL PRIMARY KEY,
    id_departemen BIGINT REFERENCES departemen(id_departemen) ON DELETE CASCADE,
    nama_proker VARCHAR(255) NOT NULL,
    deskripsi TEXT,
    foto VARCHAR(255),
    status VARCHAR(50) DEFAULT 'Belum Terlaksana',
    dibuat_pada TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    diperbarui_pada TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
