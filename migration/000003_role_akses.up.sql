CREATE TABLE IF NOT EXISTS role_akses (
    id_role BIGINT REFERENCES role(id_role) ON DELETE CASCADE,
    id_akses BIGINT REFERENCES akses(id_akses) ON DELETE CASCADE
)