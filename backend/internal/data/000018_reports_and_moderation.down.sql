DROP TABLE IF EXISTS banned_users;

ALTER TABLE IF EXISTS users 
    ADD COLUMN is_deleted BOOLEAN DEFAULT FALSE,
    ADD COLUMN deleted_at TIMESTAMP;

UPDATE users 
SET 
    is_deleted = CASE WHEN status IN ('deleted', 'banned') THEN TRUE ELSE FALSE END,
    deleted_at = CASE WHEN status IN ('deleted', 'banned') THEN NOW() ELSE NULL END;

ALTER TABLE IF EXISTS users 
    DROP COLUMN IF EXISTS status;

DROP TYPE IF EXISTS user_status;

DROP TABLE IF EXISTS moderators;

DROP INDEX IF EXISTS idx_posts_status;
DROP INDEX IF EXISTS idx_comments_status;

ALTER TABLE IF EXISTS posts 
    ADD COLUMN is_deleted BOOLEAN DEFAULT FALSE;

UPDATE posts 
SET is_deleted = CASE WHEN status = 'deleted' THEN TRUE ELSE FALSE END;

ALTER TABLE IF EXISTS posts 
    DROP COLUMN IF EXISTS status;
	
ALTER TABLE IF EXISTS comments 
    ADD COLUMN is_deleted BOOLEAN DEFAULT FALSE;

UPDATE comments 
SET is_deleted = CASE WHEN status = 'deleted' THEN TRUE ELSE FALSE END;

ALTER TABLE IF EXISTS comments 
    DROP COLUMN IF EXISTS status;

DROP TYPE IF EXISTS status;

DROP INDEX IF EXISTS idx_reports_reportable;

DROP TABLE IF EXISTS reports;

DROP TYPE IF EXISTS report_reason;
