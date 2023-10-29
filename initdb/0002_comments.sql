-- comments posted on a task to allow a discussion or clarification
CREATE TABLE comments (
    -- The unique record ID
    id           BIGSERIAL NOT NULL PRIMARY KEY,
    -- All comments are "on" a task so relate back to the specific task
    task_id      INTEGER NOT NULL,
    -- Optionally, users can leave a name if they choose
    username     VARCHAR(48) DEFAULT 'Unknown',
    --The body of the comment
    content      TEXT,
    -- When the task was created
    created_at   TIMESTAMP without time zone DEFAULT NOW(),
    -- Apply Foreign Key
    CONSTRAINT fk_task
        FOREIGN KEY (task_id)
            REFERENCES tasks(id) ON DELETE CASCADE
);

-- Index our comments table on task ID since this will be used to join almost every time it's read
CREATE INDEX task_comments ON comments (task_id);
