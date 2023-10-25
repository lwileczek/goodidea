--- Init SQL run when starting up a new DB
CREATE TABLE tasks (
    id           SERIAL not null PRIMARY KEY,
    status       BOOLEAN DEFAULT FALSE,
    title        VARCHAR(64),
    body         TEXT,
    score        INTEGER DEFAULT 0,
    completed_at TIMESTAMP without time zone
    created_at   TIMESTAMP without time zone DEFAULT NOW(),
    deleted_at   TIMESTAMP without time zone,
);

INSERT INTO tasks(title, body) VALUES('Welcome!', 'Create a task, leave a detailed description, and upvote it for importance!');
