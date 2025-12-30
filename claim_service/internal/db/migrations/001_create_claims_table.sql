BEGIN;

CREATE TABLE IF NOT EXISTS claims (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    food_listing_id UUID NOT NULL,
    donor_user_id TEXT NOT NULL,
    ngo_user_id TEXT NOT NULL,

    status TEXT NOT NULL CHECK (
        status IN (
            'REQUESTED',
            'APPROVED',
            'REJECTED',
            'PICKED_UP',
            'DELIVERED',
            'CANCELLED'
        )
    ),

    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now()
);

-- Ensure only one active claim per food listing
CREATE UNIQUE INDEX IF NOT EXISTS uniq_active_claim_per_food
ON claims (food_listing_id)
WHERE status IN ('REQUESTED', 'APPROVED', 'PICKED_UP');

-- Indexes for query performance
CREATE INDEX IF NOT EXISTS idx_claims_ngo_user
ON claims (ngo_user_id);

CREATE INDEX IF NOT EXISTS idx_claims_donor_user
ON claims (donor_user_id);

CREATE INDEX IF NOT EXISTS idx_claims_food
ON claims (food_listing_id);

COMMIT;
