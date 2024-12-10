ALTER TABLE IF EXISTS seekers
    ADD location      VARCHAR   NOT NULL,
    ADD equipment     VARCHAR   NOT NULL,
    ADD price         INTEGER   NOT NULL,
    ADD car           BOOLEAN   NOT NULL,
    ADD is_deleted    BOOLEAN   NOT NULL DEFAULT false,
    ADD deleted_at    TIMESTAMP;