BEGIN;

/* ============================================================
   1. Create donor_profiles table
   ============================================================ */

CREATE TABLE IF NOT EXISTS donor_profiles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id TEXT NOT NULL UNIQUE,
    name TEXT NOT NULL,
    phone TEXT,
    address TEXT,
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now()
);

/* ============================================================
   2. Migrate data from user_profiles → donor_profiles
   (only non-NGO users)
   ============================================================ */

INSERT INTO donor_profiles (
    user_id,
    name,
    phone,
    address,
    created_at,
    updated_at
)
SELECT
    user_id,
    name,
    phone,
    address,
    created_at,
    updated_at
FROM user_profiles
WHERE role = 'Donor'
ON CONFLICT (user_id) DO NOTHING;

/* ============================================================
   3. Drop foreign key from ngo_profiles → user_profiles
   ============================================================ */

ALTER TABLE ngo_profiles
DROP CONSTRAINT IF EXISTS fk_user;

/* ============================================================
   4. Drop user_profiles table
   ============================================================ */

DROP TABLE IF EXISTS user_profiles;

COMMIT;
