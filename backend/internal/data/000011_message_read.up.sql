CREATE TABLE IF NOT EXISTS 
    message_read 
(
    id serial PRIMARY KEY,
    user_id integer NOT NULL REFERENCES users (id),
    message_id integer NOT NULL REFERENCES messages (id),
    read_at timestamp NULL
);

ALTER TABLE IF EXISTS messages
    ADD COLUMN sender_name VARCHAR NOT NULL DEFAULT '';

ALTER TABLE IF EXISTS messages
    ALTER COLUMN sender_name DROP DEFAULT;
