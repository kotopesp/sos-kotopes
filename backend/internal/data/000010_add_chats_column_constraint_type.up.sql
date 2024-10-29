ALTER TABLE IF EXISTS chat_members 
    DROP CONSTRAINT IF EXISTS chat_members_pkey;

ALTER TABLE IF EXISTS chat_members
    DROP COLUMN IF EXISTS id;

ALTER TABLE IF EXISTS chat_members
    ADD CONSTRAINT chat_members_pkey PRIMARY KEY (user_id, chat_id);

ALTER TYPE chat_types
    ADD VALUE '';
