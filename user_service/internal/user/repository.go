package user

import (
	"context"
	"database/sql"
	"errors"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetByUserID(ctx context.Context, userID string) (*UserProfile, error) {
	query := `
		SELECT id, user_id, role, name, phone, address, created_at, updated_at
		FROM user_profiles
		WHERE user_id = $1
	`

	var profile UserProfile
	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&profile.ID,
		&profile.UserID,
		&profile.Role,
		&profile.Name,
		&profile.Phone,
		&profile.Address,
		&profile.CreatedAt,
		&profile.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("profile not found")
	}

	if err != nil {
		return nil, err
	}

	return &profile, nil
}

func (r *Repository) UpdateByUserID(
	ctx context.Context,
	userID string,
	name string,
	phone string,
	address string,
) (*UserProfile, error) {

	query := `
		UPDATE user_profiles
		SET name = $1, phone = $2, address = $3, updated_at = NOW()
		WHERE user_id = $4
		RETURNING id, user_id, role, name, phone, address, created_at, updated_at
	`

	var profile UserProfile
	err := r.db.QueryRowContext(ctx, query, name, phone, address, userID).Scan(
		&profile.ID,
		&profile.UserID,
		&profile.Role,
		&profile.Name,
		&profile.Phone,
		&profile.Address,
		&profile.CreatedAt,
		&profile.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("profile not found")
	}

	if err != nil {
		return nil, err
	}

	return &profile, nil
}

func (r *Repository) CreateProfile(
	ctx context.Context,
	userID string,
	role string,
) (*UserProfile, error) {

	query := `
		INSERT INTO user_profiles (user_id, role, name)
		VALUES ($1, $2, '')
		RETURNING id, user_id, role, name, phone, address, created_at, updated_at
	`

	var profile UserProfile
	err := r.db.QueryRowContext(ctx, query, userID, role).Scan(
		&profile.ID,
		&profile.UserID,
		&profile.Role,
		&profile.Name,
		&profile.Phone,
		&profile.Address,
		&profile.CreatedAt,
		&profile.UpdatedAt,
	)

	return &profile, err
}
