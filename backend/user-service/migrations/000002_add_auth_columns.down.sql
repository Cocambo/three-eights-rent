DROP INDEX IF EXISTS idx_users_refresh_token;

ALTER TABLE users
    DROP COLUMN IF EXISTS refresh_token_expires_at,
    DROP COLUMN IF EXISTS refresh_token,
    DROP COLUMN IF EXISTS password_hash;
