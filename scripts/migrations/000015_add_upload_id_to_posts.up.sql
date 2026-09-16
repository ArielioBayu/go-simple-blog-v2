ALTER TABLE posts 
ADD COLUMN upload_id BIGINT NULL AFTER post_hashtags,
ADD CONSTRAINT fk_posts_upload FOREIGN KEY (upload_id) REFERENCES uploads(id) ON DELETE SET NULL;
