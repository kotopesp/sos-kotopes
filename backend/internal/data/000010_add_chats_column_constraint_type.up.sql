  
ALTER TABLE IF EXISTS chat_member 
    DROP CONSTRAINT IF EXISTS chat_member_pkey;

ALTER TABLE IF EXISTS chat_member
    DROP COLUMN IF EXISTS id;

ALTER TABLE IF EXISTS chat_member
    ADD CONSTRAINT chat_member_pkey PRIMARY KEY (user_id, chat_id);

ALTER TYPE chat_types
    ADD VALUE '';