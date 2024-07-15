CREATE TABLE users
(
    id            serial       not null unique,
    name          varchar(255) not null,
    username      varchar(255) not null unique,
    password_hash varchar(255) not null
);

CREATE TABLE users_profiles
(
    id          serial       not null unique,
    title       varchar(255) not null,
    description varchar(255)
);