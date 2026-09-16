ALTER TABLE activities 
DROP FOREIGN KEY fk_post_id_activities, 
ADD CONSTRAINT fk_activities_post_id FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE;
