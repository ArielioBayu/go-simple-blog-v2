CREATE TABLE IF NOT EXISTS comments(
    id INT AUTO_INCREMENT PRIMARY KEY,
    post_id INT NOT NULL,
    user_id BIGINT NOT NULL,
    comment_content LONGTEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(50) NOT NULL,
    updated_by VARCHAR(50) NOT NULL,
    CONSTRAINT fk_post_id_posts FOREIGN KEY(post_id) REFERENCES posts(id),
    CONSTRAINT fk_user_id_users FOREIGN KEY(user_id) REFERENCES users(id)
);