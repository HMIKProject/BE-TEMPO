CREATE TABLE IF NOT EXISTS role_akses (
    id_role BIGINT REFERENCES roles(id_role) ON DELETE CASCADE,
    id_akses BIGINT REFERENCES hak_akses(id_akses) ON DELETE CASCADE,
    PRIMARY KEY (id_role, id_akses)
);