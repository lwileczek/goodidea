-- Table to store users so they can login
CREATE TABLE users (
    -- unique id
    id          CHAR(21)    NOT NULL PRIMARY KEY DEFAULT 'u' || xid(),
    -- some name for the user
    name        VARCHAR(48) NOT NULL,
    -- salt for this user's password
    salt        CHAR(8)     NOT NULL,
    -- hashed password
    passwd      CHAR(64)    NOT NULL,
    -- Time the user was created
    created_at  DATE        DEFAULT NOW()::date,
    -- Admin: if the user is an admin user or not
    admin       BOOLEAN     DEFAULT FALSE
);

INSERT INTO users(name, salt, passwd, admin) VALUES('admin', 'CWWo5vK7', 'efc48adf01fc6f785244f28677bee8a884541861c4d2b3a4ff4db691cab6259e', true);
