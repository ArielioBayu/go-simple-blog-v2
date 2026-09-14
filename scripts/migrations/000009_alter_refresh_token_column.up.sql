ALTER TABLE refresh_tokens MODIFY COLUMN refresh_token VARCHAR(255) NOT NULL;
ALTER TABLE refresh_tokens ADD CONSTRAINT unique_refresh_token UNIQUE (refresh_token);
