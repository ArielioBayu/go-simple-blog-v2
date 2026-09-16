ALTER TABLE activities 
DROP FOREIGN KEY fk_activities_post_id, 
ADD CONSTRAINT fk_post_id_activities FOREIGN KEY (post_id) REFERENCES posts(id);
