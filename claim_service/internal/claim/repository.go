package claim

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
)

var ErrClaimNotFound = errors.New("claim not found")

type ClaimRepository interface {
	Create(ctx context.Context, tx *sqlx.Tx, claim *Claim) error
	GetByID(ctx context.Context, id string) (*Claim, error)
	GetActiveByFoodID(ctx context.Context, foodListingID string) (*Claim, error)
	GetByFoodID(ctx context.Context, foodListingID string) ([]Claim, error)
	GetByNGOUserID(ctx context.Context, ngoUserID string) ([]Claim, error)
	UpdateStatus(ctx context.Context, tx *sqlx.Tx, id string, status ClaimStatus) error
}

type claimRepository struct {
	db *sqlx.DB
}

func NewClaimRepository(db *sqlx.DB) ClaimRepository {
	return &claimRepository{db: db}
}

func (r *claimRepository) Create(
	ctx context.Context,
	tx *sqlx.Tx,
	claim *Claim,
) error {
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
		VALUES (
			:id,
			:food_listing_id,
			:ngo_user_id,
			:donor_user_id,
			:status,
			:created_at,
			:updated_at
		)
	`
	_, err := tx.NamedExecContext(ctx, query, claim)
	return err
}

func (r *claimRepository) GetByID(
	ctx context.Context,
	id string,
) (*Claim, error) {
	var claim Claim

	query := `
		SELECT *
		FROM claims
		WHERE id = $1
		  AND deleted_at IS NULL
	`

	err := r.db.GetContext(ctx, &claim, query, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrClaimNotFound
	}
	return &claim, err
}

func (r *claimRepository) GetActiveByFoodID(
	ctx context.Context,
	foodListingID string,
) (*Claim, error) {
	var claim Claim

	query := `
		SELECT *
		FROM claims
		WHERE food_listing_id = $1
		  AND status IN ('REQUESTED', 'APPROVED', 'PICKED_UP')
		  AND deleted_at IS NULL
		LIMIT 1
	`

	err := r.db.GetContext(ctx, &claim, query, foodListingID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrClaimNotFound
	}
	return &claim, err
}

func (r *claimRepository) GetByFoodID(
	ctx context.Context,
	foodListingID string,
) ([]Claim, error) {
	var claims []Claim

	query := `
		SELECT *
		FROM claims
		WHERE food_listing_id = $1
		  AND deleted_at IS NULL
		ORDER BY created_at DESC
	`

	err := r.db.SelectContext(ctx, &claims, query, foodListingID)
	return claims, err
}

func (r *claimRepository) GetByNGOUserID(
	ctx context.Context,
	ngoUserID string,
) ([]Claim, error) {
	var claims []Claim

	query := `
		SELECT *
		FROM claims
		WHERE ngo_user_id = $1
		  AND deleted_at IS NULL
		ORDER BY created_at DESC
	`

	err := r.db.SelectContext(ctx, &claims, query, ngoUserID)
	return claims, err
}

func (r *claimRepository) UpdateStatus(
	ctx context.Context,
	tx *sqlx.Tx,
	id string,
	status ClaimStatus,
) error {
	query := `
		UPDATE claims
		SET status = $1,
		    updated_at = NOW()
		WHERE id = $2
		  AND deleted_at IS NULL
	`

	res, err := tx.ExecContext(ctx, query, status, id)
	if err != nil {
		return err
	}

	affected, _ := res.RowsAffected()
	if affected == 0 {
		return ErrClaimNotFound
	}
	return nil
}
