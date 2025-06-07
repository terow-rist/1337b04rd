CREATE TABLE posts (
    post_id BIGSERIAL PRIMARY KEY,
    user_id TEXT NOT NULL, 
    user_avatar TEXT NOT NULL,
    title VARCHAR(20) NOT NULL,
    content TEXT NOT NULL,
    image TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    archived_at TIMESTAMPTZ,
    FOREIGN KEY (user_id) REFERENCES users(user_id)
);
