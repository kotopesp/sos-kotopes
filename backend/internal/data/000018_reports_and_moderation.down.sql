ALTER TABLE IF EXISTS posts
    ADD COLUMN IF NOT EXISTS is_deleted BOOLEAN DEFAULT FALSE;

UPDATE posts
SET is_deleted = CASE
                     WHEN status = 'deleted' THEN TRUE
                     ELSE FALSE
    END;

ALTER TABLE IF EXISTS posts
DROP COLUMN IF EXISTS status;

DROP TYPE IF EXISTS post_status;

DROP INDEX IF EXISTS idx_unique_post_user;

DROP TABLE IF EXISTS reports;
DROP TABLE IF EXISTS moderators;
