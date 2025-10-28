CREATE TYPE willingness_carry_type AS ENUM ('yes', 'no', 'situational');

ALTER TABLE IF EXISTS seekers
    ADD animal_type       animal_category_type   NOT NULL,
    ADD location_id       INTEGER NOT NULL REFERENCES locations (id),
    ADD equipment_rental  INTEGER   NOT NULL,
    ADD have_metal_cage   BOOLEAN   NOT NULL,
    ADD have_plastic_cage BOOLEAN   NOT NULL,
    ADD have_net          BOOLEAN   NOT NULL,
    ADD have_ladder       BOOLEAN   NOT NULL,
    ADD have_other        VARCHAR,
    ADD price             INTEGER   NOT NULL,
    ADD have_car          BOOLEAN   NOT NULL,
    ADD willingness_carry willingness_carry_type   NOT NULL,
    ADD is_deleted        BOOLEAN   NOT NULL DEFAULT false,
    ADD deleted_at        TIMESTAMP;
