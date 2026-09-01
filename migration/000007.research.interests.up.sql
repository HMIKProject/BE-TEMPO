CREATE TABLE IF NOT EXISTS research_interests (
    id_interest BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    nama_interest VARCHAR(100) NOT NULL,
    deskripsi VARCHAR(255)
);