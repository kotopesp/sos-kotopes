CREATE TYPE chat_member_roles AS ENUM ('admin', 'user');

ALTER TABLE IF EXISTS messages
    ADD content VARCHAR NOT NULL DEFAULT '';

ALTER TABLE IF EXISTS messages
    ALTER COLUMN content DROP DEFAULT;

ALTER TABLE IF EXISTS chat_members
    ADD role       chat_member_roles NOT NULL DEFAULT 'user',
    ADD updated_at TIMESTAMP         NOT NULL DEFAULT now();
