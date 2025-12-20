BEGIN;

-- Drop old role constraint
ALTER TABLE user_profiles
DROP CONSTRAINT IF EXISTS user_profiles_role_check;

-- Normalize existing data (VERY IMPORTANT)
UPDATE user_profiles
SET role = 'USER'
WHERE role = 'Donor';

UPDATE user_profiles
SET role = 'ADMIN'
WHERE role = 'Admin';

-- Add new strict constraint
ALTER TABLE user_profiles
ADD CONSTRAINT user_profiles_role_check
CHECK (role IN ('USER', 'NGO', 'ADMIN'));

COMMIT;
