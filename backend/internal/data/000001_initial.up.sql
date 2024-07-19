-- Users
CREATE TABLE IF NOT EXISTS
    users
(
    id            SERIAL PRIMARY KEY,
    username      VARCHAR UNIQUE NOT NULL,
    firstname     VARCHAR,
    lastname      VARCHAR,
    description   VARCHAR,
    photo         VARCHAR,
    password_hash VARCHAR        NOT NULL,
    is_deleted    BOOLEAN        NOT NULL,
    deleted_at    TIMESTAMP,
    created_at    TIMESTAMP      NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMP      NOT NULL
);

-- Roles
CREATE TABLE IF NOT EXISTS
    seekers
(
    id          SERIAL PRIMARY KEY,
    user_id     INTEGER REFERENCES users (id),
    description VARCHAR,
    created_at  TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS
    keepers
(
    id          SERIAL PRIMARY KEY,
    user_id     INTEGER   NOT NULL REFERENCES users (id),
    description VARCHAR,
    created_at  TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS
    vets
(
    id          SERIAL PRIMARY KEY,
    user_id     INTEGER   NOT NULL REFERENCES users (id),
    description VARCHAR,
    created_at  TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMP NOT NULL
);

CREATE VIEW user_roles_view AS
SELECT user_id, 'seeker' AS role_type
FROM seekers
UNION ALL
SELECT user_id, 'keeper' AS role_type
FROM keepers
UNION ALL
SELECT user_id, 'vet' AS role_type
FROM vets;

-- Reviews
CREATE TABLE IF NOT EXISTS
    vet_reviews
(
    id         SERIAL PRIMARY KEY,
    author_id  INTEGER   NOT NULL REFERENCES users (id),
    vet_id     INTEGER   NOT NULL REFERENCES vets (id),
    content    VARCHAR,
    grade      INTEGER   NOT NULL,
    is_deleted BOOLEAN,
    deleted_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS
    keeper_reviews
(
    id         SERIAL PRIMARY KEY,
    author_id  INTEGER   NOT NULL REFERENCES users (id),
    keeper_id  INTEGER   NOT NULL REFERENCES keepers (id),
    content    VARCHAR,
    grade      INTEGER   NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL
);

-- Animals
CREATE TYPE animal_types AS ENUM ('cat', 'dog');

CREATE TYPE animal_statuses AS ENUM ('found', 'lost', 'need_home');

CREATE TYPE animal_genders AS ENUM ('male', 'female');

CREATE TABLE IF NOT EXISTS
    animals
(
    id          SERIAL PRIMARY KEY,
    keeper_id   INTEGER REFERENCES keepers (id),
    animal_type animal_types,
    age         INTEGER,
    color       VARCHAR,
    gender      animal_genders,
    description VARCHAR,
    status      animal_statuses,
    created_at  TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMP NOT NULL
);

-- Messages
CREATE TYPE chat_types AS ENUM ('keeper', 'seeker');

CREATE TABLE IF NOT EXISTS
    chats
(
    id         SERIAL PRIMARY KEY,
    chat_type  chat_types,
    is_deleted BOOLEAN   NOT NULL,
    deleted_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS
    messages
(
    id         SERIAL PRIMARY KEY,
    user_id    INTEGER   NOT NULL REFERENCES users (id),
    chat_id    INTEGER   NOT NULL REFERENCES chats (id),
    is_deleted BOOLEAN   NOT NULL,
    deleted_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS
    chat_members
(
    id         SERIAL PRIMARY KEY,
    user_id    INTEGER   NOT NULL REFERENCES users (id),
    chat_id    INTEGER REFERENCES chats (id),
    is_deleted BOOLEAN   NOT NULL,
    deleted_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Posts
CREATE TABLE IF NOT EXISTS
    posts
(
    id         SERIAL PRIMARY KEY,
    author_id  INTEGER   NOT NULL REFERENCES users (id),
    animal_id  INTEGER   NOT NULL REFERENCES animals (id),
    title      VARCHAR   NOT NULL,
    content    VARCHAR,
    is_deleted BOOLEAN   NOT NULL,
    deleted_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS
    post_response
(
    id         SERIAL PRIMARY KEY,
    post_id    INTEGER   NOT NULL REFERENCES posts (id),
    author_id  INTEGER   NOT NULL REFERENCES users (id),
    content    VARCHAR   NOT NULL,
    is_deleted BOOLEAN   NOT NULL,
    deleted_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS
    post_likes
(
    id         SERIAL PRIMARY KEY,
    post_id    INTEGER   NOT NULL REFERENCES posts (id),
    user_id    INTEGER   NOT NULL REFERENCES users (id),
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Comments
CREATE TABLE IF NOT EXISTS
    comments
(
    id         SERIAL PRIMARY KEY,
    content    VARCHAR   NOT NULL,
    author_id  INTEGER   NOT NULL REFERENCES users (id),
    posts_id   INTEGER   NOT NULL REFERENCES posts (id),
    is_deleted BOOLEAN   NOT NULL,
    deleted_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS
    comment_likes
(
    id         SERIAL PRIMARY KEY,
    comment_id INTEGER   NOT NULL REFERENCES comments (id),
    user_id    INTEGER   NOT NULL REFERENCES users (id),
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Favourites
CREATE TABLE IF NOT EXISTS
    favourite_persons
(
    id         SERIAL PRIMARY KEY,
    person_id  INTEGER   NOT NULL REFERENCES users (id),
    user_id    INTEGER   NOT NULL REFERENCES users (id),
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS
    favourite_posts
(
    id         SERIAL PRIMARY KEY,
    post_id    INTEGER   NOT NULL REFERENCES posts (id),
    user_id    INTEGER   NOT NULL REFERENCES users (id),
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS
    favourite_comments
(
    id         SERIAL PRIMARY KEY,
    comment_id INTEGER   NOT NULL REFERENCES comments (id),
    user_id    INTEGER   NOT NULL REFERENCES users (id),
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
