ALTER TABLE IF EXISTS seekers
    ADD animal_type       VARCHAR     NOT NULL,
    ADD location          VARCHAR     NOT NULL,
    ADD equipment_rental  INTEGER     NOT NULL,
    ADD price             INTEGER     NOT NULL,
    ADD have_car          BOOLEAN     NOT NULL,
    ADD willingness_carry VARCHAR     NOT NULL,
    ADD is_deleted        BOOLEAN     NOT NULL DEFAULT false,
    ADD deleted_at        TIMESTAMP;
