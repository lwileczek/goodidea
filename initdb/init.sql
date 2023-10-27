--- Init SQL run when starting up a new DB
CREATE TABLE tasks (
    -- The unique record ID
    id           SERIAL not null PRIMARY KEY,
    -- Status, is this complete or not
    status       BOOLEAN DEFAULT FALSE,
    -- The title of the request
    title        VARCHAR(64),
    -- A detailed text body explaining the request
    body         TEXT,
    -- The score for this request based on user's votes cast
    score        INTEGER DEFAULT 0,
    -- The time the task was marked complete
    completed_at TIMESTAMP without time zone,
    -- When the task was created
    created_at   TIMESTAMP without time zone DEFAULT NOW(),
    -- The date and time this was marked deleted
    deleted_at   TIMESTAMP without time zone
);
-- Index tasks table to quickly get non-deleted tasks
CREATE INDEX deleted_task_key ON tasks(deleted_at);

-- Always start with a welcome task
INSERT INTO tasks(title, body)
    VALUES('Welcome!', 'Create a task, leave a detailed description, and upvote it for importance!');


CREATE TABLE comments (
    -- The unique record ID
    id           BIGSERIAL not null PRIMARY KEY,
    -- All comments are "on" a task so relate back to the specific task
    task_id      INTEGER NOT NULL,
    FOREIGN KEY (task_id) REFERENCES tasks(id) ON DELETE CASCADE,
    -- Optionally, users can leave a name if they choose
    user         string(48),
    --The body of the comment
    content      TEXT,
    -- When the task was created
    created_at   TIMESTAMP without time zone DEFAULT NOW(),
);

-- Index our comments table on task ID since this will be used to join almost every time it's read
CREATE INDEX task_comments ON comments (task_id);
