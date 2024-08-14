ALTER TABLE chats
    ADD post_id INTEGER REFERENCES posts (id);

ALTER TYPE
    chat_types ADD VALUE 'response';