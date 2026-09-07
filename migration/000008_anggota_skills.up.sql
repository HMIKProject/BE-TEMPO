CREATE TABLE IF NOT EXISTS keahlian_anggota (
    id_anggota BIGINT REFERENCES anggota(id_anggota) ON DELETE CASCADE,
    id_keahlian BIGINT REFERENCES keahlian(id_keahlian) ON DELETE CASCADE,
    tingkat_penguasaan INT NOT NULL DEFAULT 1,
    PRIMARY KEY (id_anggota, id_keahlian)
);