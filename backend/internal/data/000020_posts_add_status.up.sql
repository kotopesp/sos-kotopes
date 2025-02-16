ALTER TABLE IF EXISTS posts
DROP COLUMN IF EXISTS deleted_at,
DROP COLUMN IF EXISTS is_deleted;

CREATE TYPE post_status AS ENUM ('deleted', 'on_moderation', 'published');

ALTER TABLE IF EXISTS posts
    ADD COLUMN status post_status NOT NULL DEFAULT 'published';
