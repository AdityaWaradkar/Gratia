CREATE TABLE IF NOT EXISTS user_profiles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    user_id TEXT NOT NULL UNIQUE,
    role TEXT NOT NULL CHECK (role IN ('Donor', 'NGO', 'Admin')),

    name TEXT NOT NULL,
    phone TEXT,
    address TEXT,

    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_user_profiles_user_id
ON user_profiles(user_id);
