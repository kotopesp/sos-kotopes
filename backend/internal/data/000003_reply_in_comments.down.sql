ALTER TABLE IF EXISTS comments
    DROP COLUMN IF EXISTS parent_id,
    DROP COLUMN IF EXISTS reply_id;
