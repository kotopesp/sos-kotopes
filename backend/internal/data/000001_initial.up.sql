- Убрать ковычки с некоторых слов +
- Названия поля
- Фронтенд хранит дефолтное +
- Убрать ограничения на VARCHAR +
- DEFAULT now() +
- имя и фамилия в юзерах +
- deleted в юзерах +
- Продумать локацию пользователей (таблица местоположений согласно дизайну)
- Rating не должно быть
- Roles лишняя таблица, либо можно оставить
- References изменить 
- content вместо text
- Продумать удаление объектов (плохо полностью удалять из базы)
- Найти другое слово для animal_type (чтобы не было type)
- Опечатка сообщения Messages
- Пока что можно удалить сообщения (спросить у лизы)
- использовать author_id и content везде
- Животное обязательно, оно всегда есть
- responser_id исправить так не пишут и text на content
- comments : posts_id


-- Users
CREATE TABLE IF NOT EXISTS
    users (
        id SERIAL PRIMARY KEY,
        username VARCHAR UNIQUE NOT,
        firstname VARCHAR,
        lastname VARCHAR,
        description TEXT,
        photo VARCHAR,
        password_hash VARCHAR NOT NULL,
        is_deleted BOOLEAN,
        created_at TIMESTAMP NOT NULL DEFAULT NOW()
    );

-- Roles
CREATE TABLE IF NOT EXISTS
    seekers (
        id SERIAL PRIMARY KEY,
        user_id INTEGER,
        description TEXT,
        location VARCHAR,
        created_at TIMESTAMP NOT NULL DEFAULT NOW()
    );

CREATE TABLE IF NOT EXISTS
    keepers (
        id SERIAL PRIMARY KEY,
        user_id INTEGER NOT NULL,
        description TEXT,
        location VARCHAR,
        created_at TIMESTAMP NOT NULL DEFAULT NOW()
    );

CREATE TABLE IF NOT EXISTS
    vets (
        id SERIAL PRIMARY KEY,
        user_id INTEGER NOT NULL,
        description VARCHAR,
        location VARCHAR,
        created_at TIMESTAMP NOT NULL DEFAULT NOW()
    );

CREATE TABLE IF NOT EXISTS
    roles_users (
        id SERIAL PRIMARY KEY,
        role VARCHAR NOT NULL,
        user_id INTEGER NOT NULL,
        created_at TIMESTAMP NOT NULL DEFAULT NOW()
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
        text VARCHAR,
        grade INTEGER NOT NULL,
        vet_id INTEGER NOT NULL,
        created_at TIMESTAMP NOT NULL DEFAULT NOW()
    );

CREATE TABLE IF NOT EXISTS
    keeper_reviews (
        id SERIAL PRIMARY KEY,
        author_id INTEGER NOT NULL,
        text VARCHAR,
        grade INTEGER NOT NULL,
        keeper_id INTEGER NOT NULL,
        created_at TIMESTAMP NOT NULL DEFAULT NOW()
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
        type animal_type,
        age INTEGER,
        color VARCHAR,
        gender animal_gender,
        description TEXT,
        status animal_status,
        keeper_id INTEGER,
        created_at TIMESTAMP NOT NULL DEFAULT NOW(),
        updated_at TIMESTAMP NOT NULL
    );

ALTER TABLE animals
ADD FOREIGN KEY (keeper_id) REFERENCES keepers (id);

-- Messeges
CREATE TABLE IF NOT EXISTS
    messages (
        id SERIAL PRIMARY KEY,
        user_id INTEGER NOT NULL,
        created_at TIMESTAMP NOT NULL DEFAULT NOW(),
        updated_at TIMESTAMP NOT NULL,
        is_deleted BOOLEAN,
        conversation_id INTEGER NOT NULL
    );

CREATE TABLE IF NOT EXISTS
    conversations (
        id SERIAL PRIMARY KEY,
        user1_id INTEGER NOT NULL,
        user2_id INTEGER NOT NULL,
        type VARCHAR,
        is_deleted BOOLEAN
        created_at TIMESTAMP NOT NULL DEFAULT NOW()
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
        created_at TIMESTAMP NOT NULL DEFAULT NOW(),
        updated_at TIMESTAMP NOT NULL,
        is_deleted BOOLEAN
        animal_id INTEGER
    );

CREATE TABLE IF NOT EXISTS
    post_response (
        id SERIAL PRIMARY KEY,
        post_id INTEGER NOT NULL,
        responser_id INTEGER NOT NULL,
        text VARCHAR NOT NULL,
        is_deleted BOOLEAN
        created_at TIMESTAMP NOT NULL DEFAULT NOW()
    );

CREATE TABLE IF NOT EXISTS
    post_likes (
        id SERIAL PRIMARY KEY,
        post_id INTEGER NOT NULL,
        user_id INTEGER NOT NULL,
        created_at TIMESTAMP NOT NULL DEFAULT NOW()
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
    comments (
        id SERIAL PRIMARY KEY,
        text TEXT NOT NULL,
        user_id INTEGER NOT NULL,
        created_at TIMESTAMP NOT NULL DEFAULT NOW(),
        updated_at TIMESTAMP NOT NULL,
        is_deleted BOOLEAN
        posts_id INTEGER NOT NULL
    );

CREATE TABLE IF NOT EXISTS
    comment_likes (
        id SERIAL PRIMARY KEY,
        comment_id INTEGER NOT NULL,
        user_id INTEGER NOT NULL,
        created_at TIMESTAMP NOT NULL DEFAULT NOW()
    );

ALTER TABLE comments
ADD FOREIGN KEY (user_id) REFERENCES users (id);

ALTER TABLE comments
ADD FOREIGN KEY (posts_id) REFERENCES posts (id);

ALTER TABLE comment_likes
ADD FOREIGN KEY (comment_id) REFERENCES comments (id);

ALTER TABLE comment_likes
ADD FOREIGN KEY (user_id) REFERENCES users (id);

-- Favourites
CREATE TABLE IF NOT EXISTS
    favourite_persons (
        id SERIAL PRIMARY KEY,
        person_id INTEGER NOT NULL,
        user_id INTEGER NOT NULL,
        created_at TIMESTAMP NOT NULL DEFAULT NOW()
    );

CREATE TABLE IF NOT EXISTS
    favourite_posts (
        id SERIAL PRIMARY KEY,
        post_id INTEGER NOT NULL,
        user_id INTEGER NOT NULL,
        created_at TIMESTAMP NOT NULL DEFAULT NOW()
    );

CREATE TABLE IF NOT EXISTS
    favourite_comments (
        id SERIAL PRIMARY KEY,
        comment_id INTEGER NOT NULL,
        user_id INTEGER NOT NULL,
        created_at TIMESTAMP NOT NULL DEFAULT NOW()
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
ADD FOREIGN KEY (comment_id) REFERENCES comments (id);

ALTER TABLE favourite_comments
ADD FOREIGN KEY (user_id) REFERENCES users (id);