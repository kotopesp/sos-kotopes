CREATE TABLE IF NOT EXISTS refresh_sessions (
    id SERIAL PRIMARY KEY,
    user_id SERIAL REFERENCES users(id) ON DELETE CASCADE,
    refresh_token UUID NOT NULL,
    fingerprint CHARACTER VARYING(200) NOT NULL UNIQUE,
    expires_at TIMESTAMP NOT NULL
);