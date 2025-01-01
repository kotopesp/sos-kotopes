ALTER TABLE IF EXISTS seekers
    ADD animal_type       VARCHAR   NOT NULL,
    ADD location          VARCHAR   NOT NULL,
    ADD equipment_rental  INTEGER   NOT NULL,
    ADD have_metal_cage   BOOLEAN   NOT NULL,
    ADD have_plastic_cage BOOLEAN   NOT NULL,
    ADD have_net          BOOLEAN   NOT NULL,
    ADD have_ladder       BOOLEAN   NOT NULL,
    ADD have_other        VARCHAR   NOT NULL,
    ADD price             INTEGER   NOT NULL,
    ADD have_car          BOOLEAN   NOT NULL,
    ADD willingness_carry VARCHAR   NOT NULL,
    ADD is_deleted        BOOLEAN   NOT NULL DEFAULT false,
    ADD deleted_at        TIMESTAMP;
