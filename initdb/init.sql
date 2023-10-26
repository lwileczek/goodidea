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

-- Always start with a welcome task
INSERT INTO tasks(title, body)
    VALUES('Welcome!', 'Create a task, leave a detailed description, and upvote it for importance!');
