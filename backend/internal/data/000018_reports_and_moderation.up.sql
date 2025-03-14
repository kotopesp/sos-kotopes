CREATE TABLE IF NOT EXISTS reports
(
    id SERIAL PRIMARY KEY,
    user_id INTEGER   NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    post_id INTEGER NOT NULL REFERENCES posts (id) ON DELETE CASCADE,
    reason 	     VARCHAR(255) NOT NULL,
    created_at    TIMESTAMP      NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS
    moderators
(
    user_id INTEGER PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    created_at  TIMESTAMP      NOT NULL DEFAULT NOW()
);


CREATE TYPE post_status AS ENUM ('deleted', 'on_moderation', 'published');

ALTER TABLE IF EXISTS posts
    ADD COLUMN IF NOT EXISTS status post_status;

UPDATE posts
    SET status = CASE
        WHEN is_deleted = TRUE THEN 'deleted'::post_status
        ELSE 'published'::post_status
END;

ALTER TABLE IF EXISTS posts
ALTER COLUMN status SET NOT NULL;

ALTER TABLE IF EXISTS posts
ALTER COLUMN status SET DEFAULT 'published';

ALTER TABLE IF EXISTS posts
DROP COLUMN IF EXISTS is_deleted;
