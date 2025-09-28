CREATE TYPE report_reason AS ENUM ('spam', 'violent_content', 'violent_speech');

CREATE TABLE reports (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    reportable_id INTEGER NOT NULL,
    reportable_type VARCHAR(20) NOT NULL,
    reason report_reason NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    
    CONSTRAINT unique_user_reportable UNIQUE (user_id, reportable_id, reportable_type)
);


CREATE INDEX idx_reports_reportable ON reports (reportable_type, reportable_id);

CREATE TYPE status AS ENUM ('deleted', 'on_moderation', 'published');

ALTER TABLE IF EXISTS posts
    ADD COLUMN IF NOT EXISTS status status DEFAULT 'published'::status;

UPDATE posts
    SET status = CASE
        WHEN is_deleted = TRUE THEN 'deleted'::status
        ELSE 'published'::status
END;

ALTER TABLE posts
    ALTER COLUMN status SET NOT NULL;

ALTER TABLE posts DROP COLUMN IF EXISTS is_deleted;

ALTER TABLE IF EXISTS comments
    ADD COLUMN IF NOT EXISTS status status DEFAULT 'published'::status;

UPDATE comments
    SET status = CASE
        WHEN is_deleted = TRUE THEN 'deleted'::status
        ELSE 'published'::status
END;

ALTER TABLE comments
    ALTER COLUMN status SET NOT NULL;


ALTER TABLE comments DROP COLUMN is_deleted;

CREATE INDEX IF NOT EXISTS idx_posts_status ON posts (status);
CREATE INDEX IF NOT EXISTS idx_comments_status ON comments (status);

CREATE TABLE IF NOT EXISTS moderators (
    user_id INTEGER PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TYPE user_status AS ENUM ('deleted', 'active', 'banned');

ALTER TABLE IF EXISTS users
    ADD COLUMN IF NOT EXISTS status user_status DEFAULT 'active'::user_status;

UPDATE users 
    SET status = CASE
        WHEN is_deleted = TRUE THEN 'deleted'::user_status
        ELSE 'active'::user_status
END;

ALTER TABLE users
    ALTER COLUMN status SET NOT NULL;

ALTER TABLE users DROP COLUMN IF EXISTS is_deleted;
ALTER TABLE users DROP COLUMN IF EXISTS deleted_at;

CREATE TABLE banned_users (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    moderator_id INTEGER NOT NULL REFERENCES moderators(user_id),
    report_id INTEGER REFERENCES reports(id) ON DELETE SET NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
