-- Users
CREATE TABLE IF NOT EXISTS "users" (
                                       "id" SERIAL PRIMARY KEY,
                                       "username" varchar UNIQUE NOT NULL,
                                       "password" varchar NOT NULL,
                                       "created_at" timestamp NOT NULL
);
CREATE TABLE IF NOT EXISTS "favourite_persons" (
                                                   "person_id" integer,
                                                   "user_id" integer,
                                                   "created_at" timestamp
);
CREATE TABLE IF NOT EXISTS "favourite_posts" (
                                                 "post_id" integer,
                                                 "user_id" integer,
                                                 "created_at" timestamp
);
CREATE TABLE IF NOT EXISTS "favourite_comments" (
                                                    "comment_id" integer,
                                                    "user_id" integer,
                                                    "created_at" timestamp
);
-- Roles
CREATE TABLE IF NOT EXISTS "roles_users" (
                                             "id" integer PRIMARY KEY,
                                             "role" varchar(50),
    "user_id" integer
    );
CREATE TABLE IF NOT EXISTS "seekers" (
                                         "id" integer PRIMARY KEY,
                                         "user_id" integer,
                                         "description" varchar,
                                         "location" varchar
);
CREATE TABLE IF NOT EXISTS "keepers" (
                                         "id" integer PRIMARY KEY,
                                         "user_id" integer,
                                         "description" varchar,
                                         "location" varchar
);
CREATE TABLE IF NOT EXISTS "vets" (
                                      "id" integer PRIMARY KEY,
                                      "user_id" integer,
                                      "description" varchar,
                                      "location" varchar
);
-- Reviews
CREATE TABLE IF NOT EXISTS "vet_reviews" (
                                             "id" integer PRIMARY KEY,
                                             "author_id" integer,
                                             "text" varchar,
                                             "grade" integer,
                                             "vet_id" integer
);
CREATE TABLE IF NOT EXISTS "keeper_reviews" (
                                                "id" integer PRIMARY KEY,
                                                "author_id" integer,
                                                "text" varchar,
                                                "grade" integer,
                                                "keeper_id" integer
);
-- Animals
CREATE TYPE animal_type AS ENUM ('cat', 'dog');
CREATE TYPE animal_status AS ENUM ('found', 'lost');
CREATE TYPE animal_gender AS ENUM ('male', 'female');
CREATE TABLE IF NOT EXISTS "animals" (
                                         "id" integer PRIMARY KEY,
                                         "type" animal_type,
                                         "age" integer,
                                         "color" varchar (30),
    "gender" animal_gender,
    "description" text,
    "status" animal_status,
    "keeper_id" integer,
    "created_at" timestamp,
    "updated_at" timestamp
    );
-- Messeges
CREATE TABLE IF NOT EXISTS "messages" (
                                          "id" integer PRIMARY KEY,
                                          "user_id" integer,
                                          "created_at" timestamp,
                                          "updated_at" timestamp,
                                          "conversation_id" integer
);
CREATE TABLE IF NOT EXISTS "conversations" (
                                               "id" integer PRIMARY KEY,
                                               "user1_id" integer,
                                               "user2_id" integer
);
-- Posts
CREATE TABLE IF NOT EXISTS "posts" (
                                       "id" integer PRIMARY KEY,
                                       "body" text,
                                       "user_id" integer,
                                       "created_at" timestamp,
                                       "updated_at" timestamp,
                                       "animal_id" integer
);
CREATE TABLE IF NOT EXISTS "post_response" (
                                               "id" integer PRIMARY KEY,
                                               "post_id" integer,
                                               "responser_id" integer,
                                               "text" varchar
);
CREATE TABLE IF NOT EXISTS "post_likes" (
                                            "id" integer PRIMARY KEY,
                                            "post_id" integer,
                                            "user_id" integer,
                                            "created_at" timestamp
);
-- Comments
CREATE TABLE IF NOT EXISTS "comments" (
                                          "id" integer PRIMARY KEY,
                                          "text" text,
                                          "user_id" integer,
                                          "created_at" timestamp,
                                          "updated_at" timestamp,
                                          "posts_id" integer
);
CREATE TABLE IF NOT EXISTS "comment_likes" (
                                               "id" integer PRIMARY KEY,
                                               "comment_id" integer,
                                               "user_id" integer,
                                               "created_at" timestamp
);
ALTER TABLE
    "vet_reviews"
    ADD
        FOREIGN KEY ("author_id") REFERENCES "users" ("id");
ALTER TABLE
    "vet_reviews"
    ADD
        FOREIGN KEY ("vet_id") REFERENCES "vets" ("id");
ALTER TABLE
    "keeper_reviews"
    ADD
        FOREIGN KEY ("author_id") REFERENCES "users" ("id");
ALTER TABLE
    "keeper_reviews"
    ADD
        FOREIGN KEY ("keeper_id") REFERENCES "keepers" ("id");
ALTER TABLE
    "seekers"
    ADD
        FOREIGN KEY ("user_id") REFERENCES "users" ("id");
ALTER TABLE
    "keepers"
    ADD
        FOREIGN KEY ("user_id") REFERENCES "users" ("id");
ALTER TABLE
    "vets"
    ADD
        FOREIGN KEY ("user_id") REFERENCES "users" ("id");
ALTER TABLE
    "animals"
    ADD
        FOREIGN KEY ("keeper_id") REFERENCES "keepers" ("id");
ALTER TABLE
    "messages"
    ADD
        FOREIGN KEY ("user_id") REFERENCES "users" ("id");
ALTER TABLE
    "messages"
    ADD
        FOREIGN KEY ("conversation_id") REFERENCES "conversations" ("id");
ALTER TABLE
    "conversations"
    ADD
        FOREIGN KEY ("user1_id") REFERENCES "users" ("id");
ALTER TABLE
    "conversations"
    ADD
        FOREIGN KEY ("user2_id") REFERENCES "users" ("id");
ALTER TABLE
    "favourite_posts"
    ADD
        FOREIGN KEY ("post_id") REFERENCES "posts" ("id");
ALTER TABLE
    "favourite_posts"
    ADD
        FOREIGN KEY ("user_id") REFERENCES "users" ("id");
ALTER TABLE
    "favourite_persons"
    ADD
        FOREIGN KEY ("person_id") REFERENCES "users" ("id");
ALTER TABLE
    "favourite_persons"
    ADD
        FOREIGN KEY ("user_id") REFERENCES "users" ("id");
ALTER TABLE
    "favourite_comments"
    ADD
        FOREIGN KEY ("comment_id") REFERENCES "comments" ("id");
ALTER TABLE
    "favourite_comments"
    ADD
        FOREIGN KEY ("user_id") REFERENCES "users" ("id");
ALTER TABLE
    "roles_users"
    ADD
        FOREIGN KEY ("user_id") REFERENCES "users" ("id");
ALTER TABLE
    "posts"
    ADD
        FOREIGN KEY ("user_id") REFERENCES "users" ("id");
ALTER TABLE
    "posts"
    ADD
        FOREIGN KEY ("animal_id") REFERENCES "animals" ("id");
ALTER TABLE
    "post_response"
    ADD
        FOREIGN KEY ("post_id") REFERENCES "posts" ("id");
ALTER TABLE
    "post_response"
    ADD
        FOREIGN KEY ("responser_id") REFERENCES "users" ("id");
ALTER TABLE
    "comment_likes"
    ADD
        FOREIGN KEY ("comment_id") REFERENCES "comments" ("id");
ALTER TABLE
    "comment_likes"
    ADD
        FOREIGN KEY ("user_id") REFERENCES "users" ("id");
ALTER TABLE
    "post_likes"
    ADD
        FOREIGN KEY ("post_id") REFERENCES "posts" ("id");
ALTER TABLE
    "post_likes"
    ADD
        FOREIGN KEY ("user_id") REFERENCES "users" ("id");
ALTER TABLE
    "comments"
    ADD
        FOREIGN KEY ("user_id") REFERENCES "users" ("id");
ALTER TABLE
    "comments"
    ADD
        FOREIGN KEY ("posts_id") REFERENCES "posts" ("id");