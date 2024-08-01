CREATE TABLE IF NOT EXISTS
    post_photos
(
    id          SERIAL PRIMARY KEY,
    post_id     INTEGER   NOT NULL REFERENCES posts (id),
    photo       BYTEA,
    is_deleted  BOOLEAN            DEFAULT false,
    deleted_at  TIMESTAMP,
    created_at  TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMP NOT NULL DEFAULT NOW()
);

ALTER TABLE IF EXISTS posts
    DROP COLUMN IF EXISTS photo;