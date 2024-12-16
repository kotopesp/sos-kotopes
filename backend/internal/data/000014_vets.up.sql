ALTER TABLE IF EXISTS vets
    ADD COLUMN is_organization BOOLEAN NOT NULL DEFAULT false,
    ADD COLUMN patronymic TEXT,
    ADD COLUMN education TEXT,
    ADD COLUMN location TEXT,
    ADD COLUMN price NUMERIC(10, 2) NOT NULL DEFAULT 0,
    ADD COLUMN org_email TEXT,
    ADD COLUMN inn_number TEXT,
    ADD COLUMN remote_consulting BOOLEAN NOT NULL DEFAULT false,
    ADD COLUMN is_inpatient BOOLEAN NOT NULL DEFAULT false,
    ADD COLUMN is_deleted BOOLEAN NOT NULL DEFAULT false,
    ADD COLUMN deleted_at TIMESTAMP;