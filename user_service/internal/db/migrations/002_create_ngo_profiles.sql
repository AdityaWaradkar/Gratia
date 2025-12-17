CREATE TABLE IF NOT EXISTS ngo_profiles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    user_id TEXT NOT NULL UNIQUE,
    organization TEXT NOT NULL,
    registration_no TEXT NOT NULL,

    verified BOOLEAN DEFAULT FALSE,
    verified_by TEXT,
    verified_at TIMESTAMP WITH TIME ZONE,

    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),

    CONSTRAINT fk_user
        FOREIGN KEY(user_id)
        REFERENCES user_profiles(user_id)
        ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_ngo_profiles_user_id
ON ngo_profiles(user_id);
