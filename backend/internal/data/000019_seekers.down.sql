ALTER TABLE IF EXISTS seekers
    DROP COLUMN IF EXISTS animal_type,
    DROP COLUMN IF EXISTS location_id,
    DROP COLUMN IF EXISTS equipment_rental
    DROP COLUMN IF EXISTS equipment
    DROP COLUMN IF EXISTS price,
    DROP COLUMN IF EXISTS have_car,
    DROP COLUMN IF EXISTS willingness_carry,
    DROP COLUMN IF EXISTS is_deleted,
    DROP COLUMN IF EXISTS deleted_at;

DROP TYPE IF EXISTS willingness_carry_type CASCADE;
