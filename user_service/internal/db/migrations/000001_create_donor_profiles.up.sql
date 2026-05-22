CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE donor_profiles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    user_id UUID NOT NULL UNIQUE,

    name TEXT NOT NULL,

    phone TEXT,

    address TEXT,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);