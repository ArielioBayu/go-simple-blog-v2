ALTER TABLE activities ADD CONSTRAINT unique_user_post UNIQUE (post_id, user_id);
