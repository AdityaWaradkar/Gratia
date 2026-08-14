package food

import (
    "context"
    "errors"
    "time"

    "github.com/jackc/pgx/v5"
    "github.com/jackc/pgx/v5/pgxpool"
)

// Repository defines database operations for food service
type Repository interface {
    CreateFoodListing(ctx context.Context, listing *FoodListing) error
    GetFoodListingByID(ctx context.Context, id string) (*FoodListing, error)
    ListFoodListings(ctx context.Context, status string) ([]*FoodListing, error)
    UpdateFoodListing(ctx context.Context, listing *FoodListing) error
    UpdateFoodStatus(ctx context.Context, id string, status string) error
    ExpireFoodListings(ctx context.Context) error
}

type repository struct {
    db *pgxpool.Pool
}

// NewRepository creates a new food repository instance
func NewRepository(db *pgxpool.Pool) Repository {
    return &repository{
        db: db,
    }
}

// CreateFoodListing inserts a new food listing into the database
func (r *repository) CreateFoodListing(ctx context.Context, listing *FoodListing) error {
    ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
    defer cancel()

    query := `
        INSERT INTO food_listings (
            donor_user_id,
            title,
            description,
            quantity,
            unit,
            expiry_time,
            location,
            image_url,
            status
        )
        VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
        RETURNING id, created_at, updated_at
    `

    return r.db.QueryRow(
        ctx,
        query,
        listing.DonorUserID,
        listing.Title,
        listing.Description,
        listing.Quantity,
        listing.Unit,
        listing.ExpiryTime,
        listing.Location,
        listing.ImageURL,
        listing.Status,
    ).Scan(
        &listing.ID,
        &listing.CreatedAt,
        &listing.UpdatedAt,
    )
}

// GetFoodListingByID retrieves a food listing by its unique identifier
func (r *repository) GetFoodListingByID(
    ctx context.Context,
    id string,
) (*FoodListing, error) {
    ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
    defer cancel()

    query := `
        SELECT
            id,
            donor_user_id,
            title,
            description,
            quantity,
            unit,
            expiry_time,
            location,
            image_url,
            status,
            created_at,
            updated_at
        FROM food_listings
        WHERE id = $1
    `

    listing := &FoodListing{}

    err := r.db.QueryRow(ctx, query, id).Scan(
        &listing.ID,
        &listing.DonorUserID,
        &listing.Title,
        &listing.Description,
        &listing.Quantity,
        &listing.Unit,
        &listing.ExpiryTime,
        &listing.Location,
        &listing.ImageURL,
        &listing.Status,
        &listing.CreatedAt,
        &listing.UpdatedAt,
    )

    if err != nil {
        if errors.Is(err, pgx.ErrNoRows) {
            return nil, errors.New("food listing not found")
        }
        return nil, err
    }

    return listing, nil
}

// ListFoodListings retrieves all food listings with a specific status
func (r *repository) ListFoodListings(
    ctx context.Context,
    status string,
) ([]*FoodListing, error) {
    ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
    defer cancel()

    query := `
        SELECT
            id,
            donor_user_id,
            title,
            description,
            quantity,
            unit,
            expiry_time,
            location,
            image_url,
            status,
            created_at,
            updated_at
        FROM food_listings
        WHERE status = $1
        ORDER BY created_at DESC
    `

    rows, err := r.db.Query(ctx, query, status)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var listings []*FoodListing = []*FoodListing{}

    for rows.Next() {
        listing := &FoodListing{}

        if err := rows.Scan(
            &listing.ID,
            &listing.DonorUserID,
            &listing.Title,
            &listing.Description,
            &listing.Quantity,
            &listing.Unit,
            &listing.ExpiryTime,
            &listing.Location,
            &listing.ImageURL,
            &listing.Status,
            &listing.CreatedAt,
            &listing.UpdatedAt,
        ); err != nil {
            return nil, err
        }

        listings = append(listings, listing)
    }

    return listings, nil
}

// UpdateFoodListing updates an existing food listing
func (r *repository) UpdateFoodListing(
    ctx context.Context,
    listing *FoodListing,
) error {
    ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
    defer cancel()

    query := `
        UPDATE food_listings
        SET
            title = $2,
            description = $3,
            quantity = $4,
            unit = $5,
            expiry_time = $6,
            location = $7,
            image_url = $8,
            updated_at = $9
        WHERE id = $1
    `

    cmd, err := r.db.Exec(
        ctx,
        query,
        listing.ID,
        listing.Title,
        listing.Description,
        listing.Quantity,
        listing.Unit,
        listing.ExpiryTime,
        listing.Location,
        listing.ImageURL,
        time.Now().UTC(),
    )

    if err != nil {
        return err
    }

    if cmd.RowsAffected() == 0 {
        return errors.New("food listing not found")
    }

    return nil
}

// UpdateFoodStatus updates the status of a food listing
func (r *repository) UpdateFoodStatus(
    ctx context.Context,
    id string,
    status string,
) error {
    ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
    defer cancel()

    cmd, err := r.db.Exec(
        ctx,
        `
        UPDATE food_listings
        SET
            status = $2,
            updated_at = $3
        WHERE id = $1
        `,
        id,
        status,
        time.Now().UTC(),
    )

    if err != nil {
        return err
    }

    if cmd.RowsAffected() == 0 {
        return errors.New("food listing not found")
    }

    return nil
}

// ExpireFoodListings updates expired food listings to expired status
func (r *repository) ExpireFoodListings(
    ctx context.Context,
) error {
    ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
    defer cancel()

    _, err := r.db.Exec(
        ctx,
        `
        UPDATE food_listings
        SET
            status = $1,
            updated_at = $2
        WHERE
            status = $3
            AND expiry_time < NOW()
        `,
        FoodStatusExpired,
        time.Now().UTC(),
        FoodStatusAvailable,
    )

    return err
}