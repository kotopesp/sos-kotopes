ALTER TABLE IF EXISTS keepers
    ADD has_cage              BOOLEAN NOT NULL,
    ADD boarding_duration     VARCHAR NOT NULL,
    ADD boarding_compensation VARCHAR NOT NULL,
    ADD animal_acceptance     VARCHAR NOT NULL,
    ADD animal_category       VARCHAR NOT NULL;
