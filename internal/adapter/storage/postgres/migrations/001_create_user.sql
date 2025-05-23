CREATE TABLE users (
    user_id TEXT PRIMARY KEY,
    name VARCHAR(70) NOT NULL,
    avatar VARCHAR(70) NOT NULL,
    expires_at TIMESTAMPTZ
);