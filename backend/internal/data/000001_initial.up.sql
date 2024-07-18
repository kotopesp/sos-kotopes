-- Users
CREATE TABLE IF NOT EXISTS
    users (
        id SERIAL PRIMARY KEY,
        username VARCHAR(50) UNIQUE NOT NULL,
        "description" TEXT,
        photo VARCHAR(100) NOT NULL,
        password_hash VARCHAR(256) NOT NULL,
        created_at TIMESTAMP NOT NULL
    );
-- Roles
CREATE TABLE IF NOT EXISTS
    seekers (
        id SERIAL PRIMARY KEY,
        user_id INTEGER,
        "description" TEXT,
        "location" VARCHAR(100),
        rating FLOAT
    );
CREATE TABLE IF NOT EXISTS
    keepers (
        id SERIAL PRIMARY KEY,
        user_id INTEGER NOT NULL,
        "description" TEXT,
        "location" VARCHAR(100),
        rating FLOAT NOT NULL
    );
CREATE TABLE IF NOT EXISTS
    vets (
        id SERIAL PRIMARY KEY,
        user_id INTEGER NOT NULL,
        "description" VARCHAR,
        "location" VARCHAR(100),
        rating FLOAT NOT NULL
    );
CREATE TABLE IF NOT EXISTS
    roles_users (
        id SERIAL PRIMARY KEY,
        "role" VARCHAR(50) NOT NULL,
        user_id INTEGER NOT NULL
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
        id SERIAL PRIMARY KEY,
        author_id INTEGER NOT NULL,
        "text" VARCHAR,
        grade INTEGER NOT NULL,
        vet_id INTEGER NOT NULL
    );
CREATE TABLE IF NOT EXISTS
    keeper_reviews (
        id SERIAL PRIMARY KEY,
        author_id INTEGER NOT NULL,
        "text" VARCHAR,
        grade INTEGER NOT NULL,
        keeper_id INTEGER NOT NULL
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
        id SERIAL PRIMARY KEY,
        "type" animal_type,
        age INTEGER,
        color VARCHAR(30),
        gender animal_gender,
        "description" TEXT,
        "status" animal_status,
        keeper_id INTEGER,
        created_at TIMESTAMP NOT NULL,
        updated_at TIMESTAMP NOT NULL
    );
ALTER TABLE animals
ADD FOREIGN KEY (keeper_id) REFERENCES keepers (id);
-- Messeges
CREATE TABLE IF NOT EXISTS
    messages (
        id SERIAL PRIMARY KEY,
        user_id INTEGER NOT NULL,
        created_at TIMESTAMP NOT NULL,
        updated_at TIMESTAMP NOT NULL,
        conversation_id INTEGER NOT NULL
    );
CREATE TABLE IF NOT EXISTS
    conversations (
        id SERIAL PRIMARY KEY,
        user1_id INTEGER NOT NULL,
        user2_id INTEGER NOT NULL,
        "type" VARCHAR
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
        id SERIAL PRIMARY KEY,
        title VARCHAR NOT NULL,
        body TEXT,
        user_id INTEGER NOT NULL,
        created_at TIMESTAMP NOT NULL,
        updated_at TIMESTAMP NOT NULL,
        animal_id INTEGER
    );
CREATE TABLE IF NOT EXISTS
    post_response (
        id SERIAL PRIMARY KEY,
        post_id INTEGER NOT NULL,
        responser_id INTEGER NOT NULL,
        "text" VARCHAR NOT NULL
    );
CREATE TABLE IF NOT EXISTS
    post_likes (
        id SERIAL PRIMARY KEY,
        post_id INTEGER NOT NULL,
        user_id INTEGER NOT NULL,
        created_at TIMESTAMP NOT NULL
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
        id SERIAL PRIMARY KEY,
        "text" TEXT NOT NULL,
        user_id INTEGER NOT NULL,
        created_at TIMESTAMP NOT NULL,
        updated_at TIMESTAMP NOT NULL,
        posts_id INTEGER NOT NULL
    );
CREATE TABLE IF NOT EXISTS
    comment_likes (
        id SERIAL PRIMARY KEY,
        comment_id INTEGER NOT NULL,
        user_id INTEGER NOT NULL,
        created_at TIMESTAMP NOT NULL
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
        id SERIAL PRIMARY KEY,
        person_id INTEGER NOT NULL,
        user_id INTEGER NOT NULL,
        created_at TIMESTAMP NOT NULL
    );
CREATE TABLE IF NOT EXISTS
    favourite_posts (
        id SERIAL PRIMARY KEY,
        post_id INTEGER NOT NULL,
        user_id INTEGER NOT NULL,
        created_at TIMESTAMP NOT NULL
    );
CREATE TABLE IF NOT EXISTS
    favourite_comments (
        id SERIAL PRIMARY KEY,
        comment_id INTEGER NOT NULL,
        user_id INTEGER NOT NULL,
        created_at TIMESTAMP NOT NULL
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