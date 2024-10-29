ALTER TABLE IF EXISTS comments
    ADD parent_id INTEGER REFERENCES comments (id),
    ADD reply_id  INTEGER REFERENCES comments (id);
