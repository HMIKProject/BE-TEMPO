CREATE TABLE IF NOT EXISTS anggota_research_interests (
    id_anggota BIGINT REFERENCES anggota(id_anggota) ON DELETE CASCADE,
    id_interest BIGINT REFERENCES research_interests(id_interest) ON DELETE CASCADE,
    PRIMARY KEY (id_anggota, id_interest)
);