BEGIN;

-- 1. Normalize existing roles
UPDATE users
SET role = 'DONOR'
WHERE role NOT IN ('DONOR', 'NGO', 'ADMIN');

-- 2. Drop old role constraint if it exists (safe)
ALTER TABLE users
DROP CONSTRAINT IF EXISTS users_role_check;

-- 3. Add strict role constraint
ALTER TABLE users
ADD CONSTRAINT users_role_check
CHECK (role IN ('DONOR', 'NGO', 'ADMIN'));

-- 4. Set correct default role
ALTER TABLE users
ALTER COLUMN role SET DEFAULT 'DONOR';

COMMIT;
