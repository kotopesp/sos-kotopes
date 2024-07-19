-- Users
CREATE TABLE IF NOT EXISTS
    users (
        id SERIAL PRIMARY KEY,
        username VARCHAR UNIQUE NOT NULL,
        firstname VARCHAR,
        lastname VARCHAR,
        description TEXT,
        photo VARCHAR,
        password_hash VARCHAR NOT NULL,
        is_deleted BOOLEAN NOT NULL,
        created_at TIMESTAMP NOT NULL DEFAULT NOW()
    );

-- Roles
CREATE TABLE IF NOT EXISTS
    seekers (
        id SERIAL PRIMARY KEY,
        user_id INTEGER REFERENCES users (id),
        description TEXT,
        location VARCHAR,
        created_at TIMESTAMP NOT NULL DEFAULT NOW()
    );

CREATE TABLE IF NOT EXISTS
    keepers (
        id SERIAL PRIMARY KEY,
        user_id INTEGER NOT NULL REFERENCES users (id),
        description TEXT,
        location VARCHAR,
        created_at TIMESTAMP NOT NULL DEFAULT NOW()
    );

CREATE TABLE IF NOT EXISTS
    vets (
        id SERIAL PRIMARY KEY,
        user_id INTEGER NOT NULL REFERENCES users (id),
        description VARCHAR,
        location VARCHAR,
        created_at TIMESTAMP NOT NULL DEFAULT NOW()
    );

CREATE TABLE IF NOT EXISTS
    roles_users (
        id SERIAL PRIMARY KEY,
        role VARCHAR NOT NULL,
        user_id INTEGER NOT NULL REFERENCES users (id),
        created_at TIMESTAMP NOT NULL DEFAULT NOW()
    );

-- Reviews
CREATE TABLE IF NOT EXISTS
    vet_reviews (
        id SERIAL PRIMARY KEY,
        author_id INTEGER NOT NULL REFERENCES users (id),
        content VARCHAR,
        grade INTEGER NOT NULL,
        vet_id INTEGER NOT NULL REFERENCES vets (id),
        created_at TIMESTAMP NOT NULL DEFAULT NOW()
    );

CREATE TABLE IF NOT EXISTS
    keeper_reviews (
        id SERIAL PRIMARY KEY,
        author_id INTEGER NOT NULL REFERENCES users (id),
        content VARCHAR,
        grade INTEGER NOT NULL,
        keeper_id INTEGER NOT NULL REFERENCES keepers (id),
        created_at TIMESTAMP NOT NULL DEFAULT NOW()
    );

-- Animals
CREATE TYPE animal_types AS ENUM('cat', 'dog');

CREATE TYPE animal_statuses AS ENUM('found', 'lost');

CREATE TYPE animal_genders AS ENUM('male', 'female');

CREATE TABLE IF NOT EXISTS
    animals (
        id SERIAL PRIMARY KEY,
        animal_type animal_type,
        age INTEGER,
        color VARCHAR,
        gender animal_gender,
        description TEXT,
        status animal_status,
        keeper_id INTEGER REFERENCES keepers (id),
        created_at TIMESTAMP NOT NULL DEFAULT NOW(),
        updated_at TIMESTAMP NOT NULL
    );

-- Messages
CREATE TABLE IF NOT EXISTS
    conversations (
        id SERIAL PRIMARY KEY,
        conversation_type VARCHAR,
        is_deleted BOOLEAN NOT NULL,
        created_at TIMESTAMP NOT NULL DEFAULT NOW()
    );

CREATE TABLE IF NOT EXISTS
    messages (
        id SERIAL PRIMARY KEY,
        user_id INTEGER NOT NULL REFERENCES users (id),
        created_at TIMESTAMP NOT NULL DEFAULT NOW(),
        updated_at TIMESTAMP NOT NULL,
        is_deleted BOOLEAN NOT NULL,
        conversation_id INTEGER NOT NULL REFERENCES conversations (id)
    );

CREATE TABLE IF NOT EXISTS
    conversation_participants (
        id SERIAL PRIMARY KEY,
        user_id INTEGER NOT NULL REFERENCES users (id),
        created_at TIMESTAMP NOT NULL DEFAULT NOW(),
        updated_at TIMESTAMP NOT NULL,
        is_deleted BOOLEAN NOT NULL
    )

-- Posts
CREATE TABLE IF NOT EXISTS
    posts (
        id SERIAL PRIMARY KEY,
        title VARCHAR NOT NULL,
        body TEXT,
        author_id INTEGER NOT NULL REFERENCES users (id),
        created_at TIMESTAMP NOT NULL DEFAULT NOW(),
        updated_at TIMESTAMP NOT NULL,
        is_deleted BOOLEAN NOT NULL,
        animal_id INTEGER REFERENCES animals (id)
    );

CREATE TABLE IF NOT EXISTS
    post_response (
        id SERIAL PRIMARY KEY,
        post_id INTEGER NOT NULL REFERENCES posts (id),
        author_id INTEGER NOT NULL REFERENCES users (id),
        content VARCHAR NOT NULL,
        is_deleted BOOLEAN NOT NULL,
        created_at TIMESTAMP NOT NULL DEFAULT NOW()
    );

CREATE TABLE IF NOT EXISTS
    post_likes (
        id SERIAL PRIMARY KEY,
        post_id INTEGER NOT NULL REFERENCES posts (id),
        user_id INTEGER NOT NULL REFERENCES users (id),
        created_at TIMESTAMP NOT NULL DEFAULT NOW()
    );

-- Comments
CREATE TABLE IF NOT EXISTS
    comments (
        id SERIAL PRIMARY KEY,
        content TEXT NOT NULL,
        author_id INTEGER NOT NULL REFERENCES users (id),
        created_at TIMESTAMP NOT NULL DEFAULT NOW(),
        updated_at TIMESTAMP NOT NULL,
        is_deleted BOOLEAN NOT NULL,
        posts_id INTEGER NOT NULL REFERENCES posts (id)
    );

CREATE TABLE IF NOT EXISTS
    comment_likes (
        id SERIAL PRIMARY KEY,
        comment_id INTEGER NOT NULL REFERENCES comments (id),
        user_id INTEGER NOT NULL REFERENCES users (id),
        created_at TIMESTAMP NOT NULL DEFAULT NOW()
    );

-- Favourites
CREATE TABLE IF NOT EXISTS
    favourite_persons (
        id SERIAL PRIMARY KEY,
        person_id INTEGER NOT NULL REFERENCES users (id),
        user_id INTEGER NOT NULL REFERENCES users (id),
        created_at TIMESTAMP NOT NULL DEFAULT NOW()
    );

CREATE TABLE IF NOT EXISTS
    favourite_posts (
        id SERIAL PRIMARY KEY,
        post_id INTEGER NOT NULL REFERENCES posts (id),
        user_id INTEGER NOT NULL REFERENCES users (id),
        created_at TIMESTAMP NOT NULL DEFAULT NOW()
    );

CREATE TABLE IF NOT EXISTS
    favourite_comments (
        id SERIAL PRIMARY KEY,
        comment_id INTEGER NOT NULL REFERENCES comments (id),
        user_id INTEGER NOT NULL REFERENCES users (id),
        created_at TIMESTAMP NOT NULL DEFAULT NOW()
    );