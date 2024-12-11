ALTER TABLE IF EXISTS chat_members
    DROP CONSTRAINT IF EXISTS chat_members_pkey;

ALTER TABLE IF EXISTS chat_members
    ADD COLUMN id SERIAL PRIMARY KEY;

ALTER TYPE chat_types
    RENAME TO chat_types_old;
CREATE TYPE chat_types
    AS ENUM('keeper', 'seeker', 'vet', 'response');
ALTER TABLE IF EXISTS chats
    ALTER COLUMN chat_type
        TYPE chat_types USING chat_type::text::chat_types;
DROP TYPE chat_types_old;
