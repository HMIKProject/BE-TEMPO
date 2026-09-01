CREATE TABLE IF NOT EXISTS anggota_skill (
    id_anggota BIGINT REFERENCES anggota(id_anggota) ON DELETE CASCADE,
    id_skill BIGINT REFERENCES skill(id_skill) ON DELETE CASCADE,
    level INT NOT NULL DEFAULT 1,
    PRIMARY KEY (id_anggota, id_skill)
);