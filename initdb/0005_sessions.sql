-- Table to store user sessionss for logins
CREATE TABLE IF NOT EXISTS sessions (
    -- unique id
    session_id  CHAR(25) PRIMARY KEY,
    -- Time the session was created
    created_at  TIMESTAMP   DEFAULT NOW(),
    -- The ID of the task the image is related too
    user_id     public.xid NOT NULL,
    CONSTRAINT fk_task
        FOREIGN KEY (user_id)
            REFERENCES users(id)
);
