CREATE TABLE comments (
    comment_id BIGSERIAL PRIMARY KEY,
    user_name TEXT NOT NULL, 
    user_avatar TEXT NOT NULL,
    post_id BIGINT NOT NULL,
    parent_comment_id BIGINT,
    content TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    FOREIGN KEY (post_id) REFERENCES posts(post_id),
    FOREIGN KEY (parent_comment_id) REFERENCES comments(comment_id)
);
