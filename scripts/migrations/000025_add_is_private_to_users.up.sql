ALTER TABLE users 
ADD COLUMN is_private BOOLEAN NOT NULL DEFAULT FALSE AFTER is_verified;

CREATE INDEX idx_users_is_private ON users (is_private);
