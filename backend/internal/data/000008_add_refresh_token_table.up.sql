CREATE TABLE IF NOT EXISTS refresh_sessions (
    id SERIAL PRIMARY KEY,
    user_id SERIAL REFERENCES users(id) ON DELETE CASCADE,
    refresh_token UUID NOT NULL,
    expires_at TIMESTAMP NOT NULL
);