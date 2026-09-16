ALTER TABLE comments 
DROP FOREIGN KEY fk_post_id_posts, 
ADD CONSTRAINT fk_comments_post_id FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE;
