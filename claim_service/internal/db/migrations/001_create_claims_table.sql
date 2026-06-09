BEGIN;

CREATE TABLE IF NOT EXISTS claims (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    food_listing_id UUID NOT NULL,
    donor_user_id TEXT NOT NULL,
    ngo_user_id TEXT NOT NULL,

    status TEXT NOT NULL CHECK (
        status IN (
            'CREATED',
            'ACCEPTED',
            'REJECTED',
            'PICKED_UP',
            'DELIVERED',
            'CANCELLED'
        )
    ),

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    accepted_at TIMESTAMPTZ,
    rejected_at TIMESTAMPTZ,
    picked_up_at TIMESTAMPTZ,
    delivered_at TIMESTAMPTZ,
    cancelled_at TIMESTAMPTZ
);

-- Only one active claim can exist for a food listing.
CREATE UNIQUE INDEX IF NOT EXISTS uniq_active_claim_per_food
ON claims (food_listing_id)
WHERE status IN (
    'CREATED',
    'ACCEPTED',
    'PICKED_UP'
);

-- Query optimisation

CREATE INDEX IF NOT EXISTS idx_claims_food_listing_id
ON claims (food_listing_id);

CREATE INDEX IF NOT EXISTS idx_claims_ngo_user_id
ON claims (ngo_user_id);

CREATE INDEX IF NOT EXISTS idx_claims_donor_user_id
ON claims (donor_user_id);

CREATE INDEX IF NOT EXISTS idx_claims_status
ON claims (status);

CREATE INDEX IF NOT EXISTS idx_claims_created_at
ON claims (created_at DESC);

COMMIT;