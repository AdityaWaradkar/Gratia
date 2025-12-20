-- Bootstrap initial ADMIN user
-- This migration is safe to run multiple times

BEGIN;

-- Insert admin only if not already present
INSERT INTO users (
    email,
    password_hash,
    role,
    email_verified,
    created_at,
    updated_at
)
SELECT
    'admin@gratia.com',
    '$2a$10$p7SfP73Dr.6Kk.AAAbCSdOfzf6hwul3AGDJ3vv2iAncUorOsWS5aa',
    'ADMIN',
    true,
    now(),
    now()
WHERE NOT EXISTS (
    SELECT 1 FROM users WHERE email = 'admin@gratia.com'
);

COMMIT;
