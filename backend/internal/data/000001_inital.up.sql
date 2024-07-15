CREATE TABLE users
(
    id            serial       not null unique,
    username      varchar(255) not null unique,
    password_hash varchar(255) not null
    created_at time             not null
);

CREATE TABLE users_profiles
(
    id          serial       not null unique,
    title       varchar(255) not null,
    description varchar(255)
);