package claim

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/jmoiron/sqlx"
)

var ErrClaimNotFound = errors.New("claim not found")

type ClaimRepository interface {
	Create(ctx context.Context, tx *sqlx.Tx, claim *Claim) error

	GetByID(ctx context.Context, id string) (*Claim, error)

	GetByIDForUpdate(
		ctx context.Context,
		tx *sqlx.Tx,
		id string,
	) (*Claim, error)

	GetActiveByFoodID(
		ctx context.Context,
		foodListingID string,
	) (*Claim, error)

	GetByFoodID(
		ctx context.Context,
		foodListingID string,
	) ([]Claim, error)

	GetByNGOUserID(
		ctx context.Context,
		ngoUserID string,
	) ([]Claim, error)

	GetByDonorUserID(
		ctx context.Context,
		donorUserID string,
	) ([]Claim, error)

	UpdateStatus(
		ctx context.Context,
		tx *sqlx.Tx,
		id string,
		status ClaimStatus,
	) error
}

type claimRepository struct {
	db *sqlx.DB
}

func NewClaimRepository(db *sqlx.DB) ClaimRepository {
	return &claimRepository{
		db: db,
	}
}

/*
Create
*/

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

/*
Reads
*/

func (r *claimRepository) GetByID(
	ctx context.Context,
	id string,
) (*Claim, error) {

	var claim Claim

	query := `
		SELECT *
		FROM claims
		WHERE id = $1
	`

	err := r.db.GetContext(ctx, &claim, query, id)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrClaimNotFound
	}

	if err != nil {
		return nil, err
	}

	return &claim, nil
}

func (r *claimRepository) GetByIDForUpdate(
	ctx context.Context,
	tx *sqlx.Tx,
	id string,
) (*Claim, error) {

	var claim Claim

	query := `
		SELECT *
		FROM claims
		WHERE id = $1
		FOR UPDATE
	`

	err := tx.GetContext(ctx, &claim, query, id)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrClaimNotFound
	}

	if err != nil {
		return nil, err
	}

	return &claim, nil
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
		  AND status IN (
				'CREATED',
				'ACCEPTED',
				'PICKED_UP'
		  )
		LIMIT 1
	`

	err := r.db.GetContext(ctx, &claim, query, foodListingID)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrClaimNotFound
	}

	if err != nil {
		return nil, err
	}

	return &claim, nil
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
		ORDER BY created_at DESC
	`

	if err := r.db.SelectContext(
		ctx,
		&claims,
		query,
		foodListingID,
	); err != nil {
		return nil, err
	}

	return claims, nil
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
		ORDER BY created_at DESC
	`

	if err := r.db.SelectContext(
		ctx,
		&claims,
		query,
		ngoUserID,
	); err != nil {
		return nil, err
	}

	return claims, nil
}

func (r *claimRepository) GetByDonorUserID(
	ctx context.Context,
	donorUserID string,
) ([]Claim, error) {

	var claims []Claim

	query := `
		SELECT *
		FROM claims
		WHERE donor_user_id = $1
		ORDER BY created_at DESC
	`

	if err := r.db.SelectContext(
		ctx,
		&claims,
		query,
		donorUserID,
	); err != nil {
		return nil, err
	}

	return claims, nil
}

/*
Updates
*/

func (r *claimRepository) UpdateStatus(
	ctx context.Context,
	tx *sqlx.Tx,
	id string,
	status ClaimStatus,
) error {

	now := time.Now().UTC()

	query := `
		UPDATE claims
		SET
			status = $1,
			updated_at = $2,

			accepted_at = CASE
				WHEN $1 = 'ACCEPTED' THEN $2
				ELSE accepted_at
			END,

			rejected_at = CASE
				WHEN $1 = 'REJECTED' THEN $2
				ELSE rejected_at
			END,

			picked_up_at = CASE
				WHEN $1 = 'PICKED_UP' THEN $2
				ELSE picked_up_at
			END,

			delivered_at = CASE
				WHEN $1 = 'DELIVERED' THEN $2
				ELSE delivered_at
			END,

			cancelled_at = CASE
				WHEN $1 = 'CANCELLED' THEN $2
				ELSE cancelled_at
			END

		WHERE id = $3
	`

	res, err := tx.ExecContext(
		ctx,
		query,
		status,
		now,
		id,
	)

	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrClaimNotFound
	}

	return nil
}