-- Table to store users so they can login
CREATE TABLE users (
    -- unique id
    id      SERIAL      not null primary key,
    -- some name for the user
    name    VARCHAR(48) not null,
    -- salt for this user's password
    salt    CHAR(8)     not null,
    -- hashed password
    passwd  CHAR(64)    not null
);

INSERT INTO users(name, salt, passwd) VALUES('admin', 'CWWo5vK7', 'efc48adf01fc6f785244f28677bee8a884541861c4d2b3a4ff4db691cab6259e');
