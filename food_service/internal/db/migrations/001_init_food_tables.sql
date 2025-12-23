CREATE TABLE food_listings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    donor_user_id TEXT NOT NULL,

    title TEXT NOT NULL,
    description TEXT,

    quantity INTEGER NOT NULL CHECK (quantity > 0),
    unit TEXT NOT NULL, -- plates, kg, packets

    expiry_time TIMESTAMPTZ NOT NULL,

    location TEXT NOT NULL,

    image_url TEXT, -- optional

    status TEXT NOT NULL CHECK (
        status IN ('OPEN', 'CLAIMED', 'EXPIRED')
    ),

    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now()
);
