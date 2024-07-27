CREATE TABLE IF NOT EXISTS
    favourite_comments
(
    id         SERIAL PRIMARY KEY,
    comment_id INTEGER   NOT NULL REFERENCES comments (id),
    user_id    INTEGER   NOT NULL REFERENCES users (id),
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);