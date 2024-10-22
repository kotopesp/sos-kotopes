CREATE TABLE IF NOT EXISTS
    photo_posts
(
    id      SERIAL PRIMARY KEY,
    post_id INTEGER   NOT NULL REFERENCES posts (id),
    photo   BYTEA
);
