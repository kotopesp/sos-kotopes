ALTER TABLE IF EXISTS keepers
    ADD location_id INTEGER NOT NULL REFERENCES locations (id),
    DROP COLUMN IF EXISTS location;
