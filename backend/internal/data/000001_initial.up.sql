-- Users
CREATE TABLE IF NOT EXISTS
    users (
        id SERIAL PRIMARY KEY,
        username VARCHAR(50) UNIQUE NOT NULL,
        "description" TEXT,
        photo VARCHAR NOT NULL,
        password_hash VARCHAR(256) NOT NULL,
        created_at TIMESTAMP NOT NULL
    );

-- Roles
CREATE TABLE IF NOT EXISTS
    seekers (
        id INTEGER PRIMARY KEY,
        user_id INTEGER,
        "description" VARCHAR,
        "location" VARCHAR,
        rating FLOAT
    );

CREATE TABLE IF NOT EXISTS
    keepers (
        id INTEGER PRIMARY KEY,
        user_id INTEGER,
        "description" VARCHAR,
        "location" VARCHAR,
        rating FLOAT
    );

CREATE TABLE IF NOT EXISTS
    vets (
        id INTEGER PRIMARY KEY,
        user_id INTEGER,
        "description" VARCHAR,
        "location" VARCHAR,
        rating FLOAT
    );

CREATE TABLE IF NOT EXISTS
    roles_users (
        id INTEGER PRIMARY KEY,
        "role" VARCHAR(50),
        user_id INTEGER
    );

ALTER TABLE seekers
ADD FOREIGN KEY (user_id) REFERENCES users (id);

ALTER TABLE keepers
ADD FOREIGN KEY (user_id) REFERENCES users (id);

ALTER TABLE vets
ADD FOREIGN KEY (user_id) REFERENCES users (id);

ALTER TABLE roles_users
ADD FOREIGN KEY (user_id) REFERENCES users (id);

-- Reviews
CREATE TABLE IF NOT EXISTS
    vet_reviews (
        id INTEGER PRIMARY KEY,
        author_id INTEGER,
        "text" VARCHAR,
        grade INTEGER,
        vet_id INTEGER
    );

CREATE TABLE IF NOT EXISTS
    keeper_reviews (
        id INTEGER PRIMARY KEY,
        author_id INTEGER,
        "text" VARCHAR,
        grade INTEGER,
        keeper_id INTEGER
    );

ALTER TABLE vet_reviews
ADD FOREIGN KEY (author_id) REFERENCES users (id);

ALTER TABLE vet_reviews
ADD FOREIGN KEY (vet_id) REFERENCES vets (id);

ALTER TABLE keeper_reviews
ADD FOREIGN KEY (author_id) REFERENCES users (id);

ALTER TABLE keeper_reviews
ADD FOREIGN KEY (keeper_id) REFERENCES keepers (id);

-- Animals
CREATE TYPE animal_type AS ENUM('cat', 'dog');

CREATE TYPE animal_status AS ENUM('found', 'lost');

CREATE TYPE animal_gender AS ENUM('male', 'female');

CREATE TABLE IF NOT EXISTS
    animals (
        id INTEGER PRIMARY KEY,
        "type" animal_type,
        age INTEGER,
        color VARCHAR(30),
        gender animal_gender,
        "description" TEXT,
        "status" animal_status,
        keeper_id INTEGER,
        created_at TIMESTAMP,
        updated_at TIMESTAMP
    );

ALTER TABLE animals
ADD FOREIGN KEY (keeper_id) REFERENCES keepers (id);

-- Messeges
CREATE TABLE IF NOT EXISTS
    messages (
        id INTEGER PRIMARY KEY,
        user_id INTEGER,
        created_at TIMESTAMP,
        updated_at TIMESTAMP,
        conversation_id INTEGER
    );

CREATE TABLE IF NOT EXISTS
    conversations (
        id INTEGER PRIMARY KEY,
        user1_id INTEGER,
        user2_id INTEGER,
        "type" VARCHAR,
    );

ALTER TABLE messages
ADD FOREIGN KEY (user_id) REFERENCES users (id);

ALTER TABLE messages
ADD FOREIGN KEY (conversation_id) REFERENCES conversations (id);

ALTER TABLE conversations
ADD FOREIGN KEY (user1_id) REFERENCES users (id);

ALTER TABLE conversations
ADD FOREIGN KEY (user2_id) REFERENCES users (id);

-- Posts
CREATE TABLE IF NOT EXISTS
    posts (
        id INTEGER PRIMARY KEY,
        body TEXT,
        user_id INTEGER,
        created_at TIMESTAMP,
        updated_at TIMESTAMP,
        animal_id INTEGER
    );

CREATE TABLE IF NOT EXISTS
    post_response (
        id INTEGER PRIMARY KEY,
        post_id INTEGER,
        responser_id INTEGER,
        "text" VARCHAR
    );

CREATE TABLE IF NOT EXISTS
    post_likes (
        id INTEGER PRIMARY KEY,
        post_id INTEGER,
        user_id INTEGER,
        created_at TIMESTAMP
    );

ALTER TABLE posts
ADD FOREIGN KEY (user_id) REFERENCES users (id);

ALTER TABLE posts
ADD FOREIGN KEY (animal_id) REFERENCES animals (id);

ALTER TABLE post_response
ADD FOREIGN KEY (post_id) REFERENCES posts (id);

ALTER TABLE post_response
ADD FOREIGN KEY (responser_id) REFERENCES users (id);

ALTER TABLE post_likes
ADD FOREIGN KEY (post_id) REFERENCES posts (id);

ALTER TABLE post_likes
ADD FOREIGN KEY (user_id) REFERENCES users (id);

-- Comments
CREATE TABLE IF NOT EXISTS
    "comments" (
        id INTEGER PRIMARY KEY,
        "text" TEXT,
        user_id INTEGER,
        created_at TIMESTAMP,
        updated_at TIMESTAMP,
        posts_id INTEGER
    );

CREATE TABLE IF NOT EXISTS
    comment_likes (
        id INTEGER PRIMARY KEY,
        comment_id INTEGER,
        user_id INTEGER,
        created_at TIMESTAMP
    );

ALTER TABLE "comments"
ADD FOREIGN KEY (user_id) REFERENCES users (id);

ALTER TABLE "comments"
ADD FOREIGN KEY (posts_id) REFERENCES posts (id);

ALTER TABLE comment_likes
ADD FOREIGN KEY (comment_id) REFERENCES "comments" (id);

ALTER TABLE comment_likes
ADD FOREIGN KEY (user_id) REFERENCES users (id);

-- Favourites
CREATE TABLE IF NOT EXISTS
    favourite_persons (
        person_id INTEGER,
        user_id INTEGER,
        created_at TIMESTAMP
    );

CREATE TABLE IF NOT EXISTS
    favourite_posts (
        post_id INTEGER,
        user_id INTEGER,
        created_at TIMESTAMP
    );

CREATE TABLE IF NOT EXISTS
    favourite_comments (
        comment_id INTEGER,
        user_id INTEGER,
        created_at TIMESTAMP
    );

ALTER TABLE favourite_posts
ADD FOREIGN KEY (post_id) REFERENCES posts (id);

ALTER TABLE favourite_posts
ADD FOREIGN KEY (user_id) REFERENCES users (id);

ALTER TABLE favourite_persons
ADD FOREIGN KEY (person_id) REFERENCES users (id);

ALTER TABLE favourite_persons
ADD FOREIGN KEY (user_id) REFERENCES users (id);

ALTER TABLE favourite_comments
ADD FOREIGN KEY (comment_id) REFERENCES "comments" (id);

ALTER TABLE favourite_comments
ADD FOREIGN KEY (user_id) REFERENCES users (id);