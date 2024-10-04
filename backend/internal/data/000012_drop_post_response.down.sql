CREATE TABLE IF NOT EXISTS
    post_response
(
    id         SERIAL PRIMARY KEY,
    post_id    INTEGER   NOT NULL REFERENCES posts (id),
    author_id  INTEGER   NOT NULL REFERENCES users (id),
    content    VARCHAR   NOT NULL,
    is_deleted BOOLEAN   NOT NULL DEFAULT false,
    deleted_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);
