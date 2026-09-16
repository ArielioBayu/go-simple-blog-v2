ALTER TABLE comments 
DROP FOREIGN KEY fk_comments_post_id, 
ADD CONSTRAINT fk_post_id_posts FOREIGN KEY (post_id) REFERENCES posts(id);
