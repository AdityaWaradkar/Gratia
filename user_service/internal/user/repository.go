package user

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository defines database operations for user service
type Repository interface {
	// Donor profile operations
	CreateDonorProfile(ctx context.Context, donor *DonorProfile) error
	GetDonorProfileByUserID(ctx context.Context, userID string) (*DonorProfile, error)
	UpdateDonorProfile(ctx context.Context, donor *DonorProfile) error

	// NGO profile operations
	CreateNGOProfile(ctx context.Context, ngo *NGOProfile) error
	GetNGOProfileByUserID(ctx context.Context, userID string) (*NGOProfile, error)
	UpdateNGOVerification(
		ctx context.Context,
		userID string,
		verified bool,
		verifiedBy *string,
	) error
}

type repository struct {
	db *pgxpool.Pool
}

// NewRepository creates a new user repository instance
func NewRepository(db *pgxpool.Pool) Repository {
	return &repository{db: db}
}

// CreateDonorProfile inserts a new donor profile into the database
func (r *repository) CreateDonorProfile(
	ctx context.Context,
	donor *DonorProfile,
) error {

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
		INSERT INTO donor_profiles (
			user_id, name, phone, address
		)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRow(
		ctx,
		query,
		donor.UserID,
		donor.Name,
		donor.Phone,
		donor.Address,
	).Scan(
		&donor.ID,
		&donor.CreatedAt,
		&donor.UpdatedAt,
	)

	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return ErrDonorProfileExists
			}
		}

		return err
	}

	return nil
}

// GetDonorProfileByUserID retrieves a donor profile by user ID
func (r *repository) GetDonorProfileByUserID(
	ctx context.Context,
	userID string,
) (*DonorProfile, error) {

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
		SELECT
			id,
			user_id,
			name,
			phone,
			address,
			created_at,
			updated_at
		FROM donor_profiles
		WHERE user_id = $1
	`

	d := &DonorProfile{}

	err := r.db.QueryRow(ctx, query, userID).Scan(
		&d.ID,
		&d.UserID,
		&d.Name,
		&d.Phone,
		&d.Address,
		&d.CreatedAt,
		&d.UpdatedAt,
	)

	if err != nil {
		return nil, ErrDonorProfileNotFound
	}

	return d, nil
}

// UpdateDonorProfile updates an existing donor profile
func (r *repository) UpdateDonorProfile(
	ctx context.Context,
	donor *DonorProfile,
) error {

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	cmd, err := r.db.Exec(
		ctx,
		`
		UPDATE donor_profiles
		SET
			name = $2,
			phone = $3,
			address = $4,
			updated_at = $5
		WHERE user_id = $1
		`,
		donor.UserID,
		donor.Name,
		donor.Phone,
		donor.Address,
		time.Now(),
	)

	if err != nil {
		return err
	}

	if cmd.RowsAffected() == 0 {
		return ErrDonorProfileNotFound
	}

	return nil
}

// CreateNGOProfile inserts a new NGO profile into the database
func (r *repository) CreateNGOProfile(
	ctx context.Context,
	ngo *NGOProfile,
) error {

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
		INSERT INTO ngo_profiles (
			user_id,
			organization,
			registration_no,
			verified
		)
		VALUES ($1, $2, $3, false)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRow(
		ctx,
		query,
		ngo.UserID,
		ngo.Organization,
		ngo.RegistrationNo,
	).Scan(
		&ngo.ID,
		&ngo.CreatedAt,
		&ngo.UpdatedAt,
	)

	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return ErrNGOProfileExists
			}
		}

		return err
	}

	return nil
}

// GetNGOProfileByUserID retrieves an NGO profile by user ID
func (r *repository) GetNGOProfileByUserID(
	ctx context.Context,
	userID string,
) (*NGOProfile, error) {

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
		SELECT
			id,
			user_id,
			organization,
			registration_no,
			verified,
			verified_by,
			verified_at,
			created_at,
			updated_at
		FROM ngo_profiles
		WHERE user_id = $1
	`

	n := &NGOProfile{}

	err := r.db.QueryRow(ctx, query, userID).Scan(
		&n.ID,
		&n.UserID,
		&n.Organization,
		&n.RegistrationNo,
		&n.Verified,
		&n.VerifiedBy,
		&n.VerifiedAt,
		&n.CreatedAt,
		&n.UpdatedAt,
	)

	if err != nil {
		return nil, ErrNGOProfileNotFound
	}

	return n, nil
}

// UpdateNGOVerification updates the verification status of an NGO profile
func (r *repository) UpdateNGOVerification(
	ctx context.Context,
	userID string,
	verified bool,
	verifiedBy *string,
) error {

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	now := time.Now()

	cmd, err := r.db.Exec(
		ctx,
		`
		UPDATE ngo_profiles
		SET
			verified = $2,
			verified_by = $3,
			verified_at = $4,
			updated_at = $4
		WHERE user_id = $1
		`,
		userID,
		verified,
		verifiedBy,
		now,
	)

	if err != nil {
		return err
	}

	if cmd.RowsAffected() == 0 {
		return ErrNGOProfileNotFound
	}

	return nil
}