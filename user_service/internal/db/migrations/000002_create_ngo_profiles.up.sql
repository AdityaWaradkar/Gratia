CREATE TABLE ngo_profiles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    user_id UUID NOT NULL UNIQUE,

    organization TEXT NOT NULL,

    registration_no TEXT NOT NULL,

    verified BOOLEAN NOT NULL DEFAULT FALSE,

    verified_by UUID,

    verified_at TIMESTAMPTZ,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);