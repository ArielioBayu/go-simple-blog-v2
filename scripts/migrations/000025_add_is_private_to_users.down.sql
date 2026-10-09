DROP INDEX IF EXISTS idx_users_is_private ON users;
ALTER TABLE users DROP COLUMN is_private;
