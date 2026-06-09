CREATE TABLE food_listings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    donor_user_id TEXT NOT NULL,

    title TEXT NOT NULL,
    description TEXT,

    quantity INTEGER NOT NULL CHECK (quantity > 0),
    unit TEXT NOT NULL,

    expiry_time TIMESTAMPTZ NOT NULL,

    location TEXT NOT NULL,

    image_url TEXT,

    status TEXT NOT NULL CHECK (
        status IN (
            'AVAILABLE',
            'CLAIMED',
            'EXPIRED',
            'CANCELLED'
        )
    ),

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_food_listings_status
ON food_listings(status);

CREATE INDEX idx_food_listings_donor_user_id
ON food_listings(donor_user_id);

CREATE INDEX idx_food_listings_expiry_time
ON food_listings(expiry_time);