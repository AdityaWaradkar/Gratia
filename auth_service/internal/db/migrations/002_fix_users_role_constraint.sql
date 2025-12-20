BEGIN;

-- 1. Drop old role constraint
ALTER TABLE users
DROP CONSTRAINT IF EXISTS users_role_check;

-- 2. Add new constraint aligned with auth_service
ALTER TABLE users
ADD CONSTRAINT users_role_check
CHECK (role IN ('USER', 'ADMIN'));

COMMIT;
