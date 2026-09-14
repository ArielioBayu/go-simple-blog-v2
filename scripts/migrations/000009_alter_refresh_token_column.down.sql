ALTER TABLE refresh_tokens DROP INDEX unique_refresh_token;
ALTER TABLE refresh_tokens MODIFY COLUMN refresh_token TEXT NOT NULL;
