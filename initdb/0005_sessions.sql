-- Table to store user sessionss for logins
CREATE TABLE IF NOT EXISTS sessions (
    -- unique id
    session_id  public.xid PRIMARY KEY DEFAULT xid(),
    -- Time the session was created
    created_at  TIMESTAMP   DEFAULT NOW(),
    -- is the session currently valid
    valid       BOOLEAN     DEFAULT TRUE,
    -- The ID of the task the image is related too
    user_id     public.xid NOT NULL,
    CONSTRAINT fk_task
        FOREIGN KEY (user_id)
            REFERENCES users(id)
);
