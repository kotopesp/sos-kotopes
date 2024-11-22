ALTER TYPE chat_types 
    RENAME TO chat_types_old;
CREATE TYPE chat_types 
    AS ENUM('keeper', 'seeker', 'vet', 'response', '');
ALTER TABLE chats ALTER COLUMN chat_type
    TYPE chat_types USING chat_type::text::chat_types;
DROP TYPE chat_types_old;