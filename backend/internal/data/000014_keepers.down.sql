ALTER TABLE IF EXISTS keepers
    DROP COLUMN IF EXISTS has_cage,
    DROP COLUMN IF EXISTS boarding_duration,
    DROP COLUMN IF EXISTS boarding_compensation,
    DROP COLUMN IF EXISTS animal_acceptance,
    DROP COLUMN IF EXISTS animal_category,
    DROP COLUMN IF EXISTS deleted_at,
    DROP COLUMN IF EXISTS is_deleted;
