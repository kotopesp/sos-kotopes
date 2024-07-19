CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    ext_id INTEGER UNIQUE,
    username VARCHAR UNIQUE NOT NULL,
    firstname VARCHAR,
    lastname VARCHAR,
    description TEXT,
    photo BYTEA,
    password_hash VARCHAR NOT NULL,
    is_deleted BOOLEAN,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);