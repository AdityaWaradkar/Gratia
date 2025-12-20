BEGIN;

-- 1. Add role check constraint only if it does not exist
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM pg_constraint
        WHERE conname = 'user_profiles_role_check'
    ) THEN
        ALTER TABLE user_profiles
        ADD CONSTRAINT user_profiles_role_check
        CHECK (role IN ('USER', 'ADMIN'));
    END IF;
END $$;

-- 2. Ensure verified_at is timestamptz (safe conversion)
ALTER TABLE user_profiles
ALTER COLUMN verified_at
TYPE timestamptz
USING verified_at AT TIME ZONE 'UTC';

COMMIT;
