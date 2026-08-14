package claim

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Domain-level sentinel errors
var (
	ErrClaimNotFound = errors.New("claim not found")
)

// Repository defines the interface for claim data access operations
type Repository interface {
	Create(ctx context.Context, claim *Claim) error
	GetByID(ctx context.Context, id string) (*Claim, error)
	GetActiveByFoodID(ctx context.Context, foodListingID string) (*Claim, error)
	GetByFoodID(ctx context.Context, foodListingID string) ([]Claim, error)
	GetByNGOUserID(ctx context.Context, ngoUserID string) ([]Claim, error)
	GetByDonorUserID(ctx context.Context, donorUserID string) ([]Claim, error)
	UpdateStatus(ctx context.Context, id string, status ClaimStatus) error
}

type repository struct {
	db *pgxpool.Pool
}

// NewRepository creates a new claim repository instance
func NewRepository(db *pgxpool.Pool) Repository {
	return &repository{
		db: db,
	}
}

// Create inserts a new claim record into the database
func (r *repository) Create(ctx context.Context, claim *Claim) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
		INSERT INTO claims (
			id,
			food_listing_id,
			ngo_user_id,
			donor_user_id,
			status,
			created_at,
			updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING created_at, updated_at
	`

	err := r.db.QueryRow(
		ctx,
		query,
		claim.ID,
		claim.FoodListingID,
		claim.NGOUserID,
		claim.DonorUserID,
		claim.Status,
		claim.CreatedAt,
		claim.UpdatedAt,
	).Scan(&claim.CreatedAt, &claim.UpdatedAt)

	return err
}

// GetByID retrieves a claim by its unique identifier
func (r *repository) GetByID(ctx context.Context, id string) (*Claim, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var claim Claim

	query := `
		SELECT
			id,
			food_listing_id,
			ngo_user_id,
			donor_user_id,
			status,
			created_at,
			updated_at,
			accepted_at,
			rejected_at,
			picked_up_at,
			delivered_at,
			cancelled_at
		FROM claims
		WHERE id = $1
	`

	err := r.db.QueryRow(ctx, query, id).Scan(
		&claim.ID,
		&claim.FoodListingID,
		&claim.NGOUserID,
		&claim.DonorUserID,
		&claim.Status,
		&claim.CreatedAt,
		&claim.UpdatedAt,
		&claim.AcceptedAt,
		&claim.RejectedAt,
		&claim.PickedUpAt,
		&claim.DeliveredAt,
		&claim.CancelledAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrClaimNotFound
		}
		return nil, err
	}

	return &claim, nil
}

// GetActiveByFoodID retrieves an active claim for a specific food listing
func (r *repository) GetActiveByFoodID(ctx context.Context, foodListingID string) (*Claim, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var claim Claim

	query := `
		SELECT
			id,
			food_listing_id,
			ngo_user_id,
			donor_user_id,
			status,
			created_at,
			updated_at,
			accepted_at,
			rejected_at,
			picked_up_at,
			delivered_at,
			cancelled_at
		FROM claims
		WHERE food_listing_id = $1
		  AND status IN ('CREATED', 'ACCEPTED', 'PICKED_UP')
		LIMIT 1
	`

	err := r.db.QueryRow(ctx, query, foodListingID).Scan(
		&claim.ID,
		&claim.FoodListingID,
		&claim.NGOUserID,
		&claim.DonorUserID,
		&claim.Status,
		&claim.CreatedAt,
		&claim.UpdatedAt,
		&claim.AcceptedAt,
		&claim.RejectedAt,
		&claim.PickedUpAt,
		&claim.DeliveredAt,
		&claim.CancelledAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrClaimNotFound
		}
		return nil, err
	}

	return &claim, nil
}

// GetByFoodID retrieves all claims for a specific food listing
func (r *repository) GetByFoodID(ctx context.Context, foodListingID string) ([]Claim, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
		SELECT
			id,
			food_listing_id,
			ngo_user_id,
			donor_user_id,
			status,
			created_at,
			updated_at,
			accepted_at,
			rejected_at,
			picked_up_at,
			delivered_at,
			cancelled_at
		FROM claims
		WHERE food_listing_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(ctx, query, foodListingID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var claims []Claim

	for rows.Next() {
		var claim Claim

		if err := rows.Scan(
			&claim.ID,
			&claim.FoodListingID,
			&claim.NGOUserID,
			&claim.DonorUserID,
			&claim.Status,
			&claim.CreatedAt,
			&claim.UpdatedAt,
			&claim.AcceptedAt,
			&claim.RejectedAt,
			&claim.PickedUpAt,
			&claim.DeliveredAt,
			&claim.CancelledAt,
		); err != nil {
			return nil, err
		}

		claims = append(claims, claim)
	}

	return claims, nil
}

// GetByNGOUserID retrieves all claims made by a specific NGO user
func (r *repository) GetByNGOUserID(ctx context.Context, ngoUserID string) ([]Claim, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
		SELECT
			id,
			food_listing_id,
			ngo_user_id,
			donor_user_id,
			status,
			created_at,
			updated_at,
			accepted_at,
			rejected_at,
			picked_up_at,
			delivered_at,
			cancelled_at
		FROM claims
		WHERE ngo_user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(ctx, query, ngoUserID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var claims []Claim

	for rows.Next() {
		var claim Claim

		if err := rows.Scan(
			&claim.ID,
			&claim.FoodListingID,
			&claim.NGOUserID,
			&claim.DonorUserID,
			&claim.Status,
			&claim.CreatedAt,
			&claim.UpdatedAt,
			&claim.AcceptedAt,
			&claim.RejectedAt,
			&claim.PickedUpAt,
			&claim.DeliveredAt,
			&claim.CancelledAt,
		); err != nil {
			return nil, err
		}

		claims = append(claims, claim)
	}

	return claims, nil
}

// GetByDonorUserID retrieves all claims for a specific donor user
func (r *repository) GetByDonorUserID(ctx context.Context, donorUserID string) ([]Claim, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
		SELECT
			id,
			food_listing_id,
			ngo_user_id,
			donor_user_id,
			status,
			created_at,
			updated_at,
			accepted_at,
			rejected_at,
			picked_up_at,
			delivered_at,
			cancelled_at
		FROM claims
		WHERE donor_user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(ctx, query, donorUserID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var claims []Claim

	for rows.Next() {
		var claim Claim

		if err := rows.Scan(
			&claim.ID,
			&claim.FoodListingID,
			&claim.NGOUserID,
			&claim.DonorUserID,
			&claim.Status,
			&claim.CreatedAt,
			&claim.UpdatedAt,
			&claim.AcceptedAt,
			&claim.RejectedAt,
			&claim.PickedUpAt,
			&claim.DeliveredAt,
			&claim.CancelledAt,
		); err != nil {
			return nil, err
		}

		claims = append(claims, claim)
	}

	return claims, nil
}

// UpdateStatus updates the status of a claim with timestamp tracking
func (r *repository) UpdateStatus(ctx context.Context, id string, status ClaimStatus) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	now := time.Now()

	query := `
		UPDATE claims
		SET
			status = $1,
			updated_at = $2,
			accepted_at = CASE WHEN $1 = 'ACCEPTED' THEN $2 ELSE accepted_at END,
			rejected_at = CASE WHEN $1 = 'REJECTED' THEN $2 ELSE rejected_at END,
			picked_up_at = CASE WHEN $1 = 'PICKED_UP' THEN $2 ELSE picked_up_at END,
			delivered_at = CASE WHEN $1 = 'DELIVERED' THEN $2 ELSE delivered_at END,
			cancelled_at = CASE WHEN $1 = 'CANCELLED' THEN $2 ELSE cancelled_at END
		WHERE id = $3
	`

	cmd, err := r.db.Exec(ctx, query, status, now, id)
	if err != nil {
		return err
	}

	if cmd.RowsAffected() == 0 {
		return ErrClaimNotFound
	}

	return nil
}