ALTER TABLE messages
    DROP COLUMN IF EXISTS content;

ALTER TABLE chat_members
    DROP COLUMN IF EXISTS role,
    DROP COLUMN IF EXISTS updated_at;

DROP TYPE IF EXISTS chat_member_roles CASCADE;
