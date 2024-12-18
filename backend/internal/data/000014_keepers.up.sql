CREATE TYPE animal_acceptance_type AS ENUM ('home', 'homeless', 'homeless-hadhome', 'depends');
CREATE TYPE animal_category_type AS ENUM ('cat', 'dog', 'both');
CREATE TYPE boarding_compensation_type AS ENUM ('paid', 'free', 'depends');
CREATE TYPE boarding_duration_type AS ENUM ('hours', 'days', 'weeks', 'months', 'depends');

ALTER TABLE IF EXISTS keepers
    ADD has_cage              BOOLEAN NOT NULL,
    ADD boarding_duration     boarding_duration_type NOT NULL,
    ADD boarding_compensation boarding_compensation_type NOT NULL,
    ADD animal_acceptance     animal_acceptance_type NOT NULL,
    ADD animal_category       animal_category_type NOT NULL,
    ADD is_deleted            BOOLEAN NOT NULL DEFAULT false,
    ADD deleted_at            TIMESTAMP;
