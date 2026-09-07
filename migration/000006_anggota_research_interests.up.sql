CREATE TABLE IF NOT EXISTS minat_riset_anggota (
    id_anggota BIGINT REFERENCES anggota(id_anggota) ON DELETE CASCADE,
    id_minat BIGINT REFERENCES minat_riset(id_minat) ON DELETE CASCADE,
    PRIMARY KEY (id_anggota, id_minat)
);