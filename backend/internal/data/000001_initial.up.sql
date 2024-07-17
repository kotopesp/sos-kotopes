CREATE TABLE users (
    id SERIAL PRIMARY KEY NOT NULL,
    username VARCHAR(100) UNIQUE NOT NULL,
    password VARCHAR(100) NOT NULL,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE roles (
    id SERIAL PRIMARY KEY NOT NULL,
    name VARCHAR(100) NOT NULL
);

INSERT INTO roles(name) VALUES('seeker');
INSERT INTO roles(name) VALUES('keeper');
INSERT INTO roles(name) VALUES('vet');

CREATE TABLE roles_users (
    id SERIAL PRIMARY KEY NOT NULL,
    role_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    CONSTRAINT fk_user
         FOREIGN KEY(user_id)
            REFERENCES users(id)
            ON DELETE CASCADE,
    CONSTRAINT fk_role
         FOREIGN KEY(role_id)
            REFERENCES roles(id)
            ON DELETE CASCADE
);

CREATE TABLE seekers (
    id SERIAL PRIMARY KEY NOT NULL,
    user_id INTEGER NOT NULL,
    description TEXT,
    location TEXT NOT NULL,
    CONSTRAINT fk_user
         FOREIGN KEY(user_id)
            REFERENCES users(id)
            ON DELETE CASCADE
);

CREATE TABLE keepers (
    id SERIAL PRIMARY KEY NOT NULL,
    user_id INTEGER NOT NULL,
    description TEXT,
    location TEXT NOT NULL,
    CONSTRAINT fk_user
         FOREIGN KEY(user_id)
            REFERENCES users(id)
            ON DELETE CASCADE
);

CREATE TABLE vets (
    id SERIAL PRIMARY KEY NOT NULL,
    user_id INTEGER NOT NULL,
    description TEXT,
    location TEXT NOT NULL,
    CONSTRAINT fk_user
         FOREIGN KEY(user_id)
            REFERENCES users(id)
            ON DELETE CASCADE
);