package repository

import (
	"context"

	"github.com/adityawaradkar/gratia/auth_service/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	DB *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		DB: db,
	}
}

func (r *UserRepository) CreateUser(
	ctx context.Context,
	user *model.User,
) error {

	query := `
		INSERT INTO users (
			email,
			password_hash,
			role
		)
		VALUES ($1, $2, $3)
		RETURNING
			user_id,
			created_at,
			updated_at
	`

	err := r.DB.QueryRow(
		ctx,
		query,
		user.Email,
		user.PasswordHash,
		user.Role,
	).Scan(
		&user.UserID,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	return err
}

func (r *UserRepository) GetUserByEmail(
	ctx context.Context,
	email string,
) (*model.User, error) {

	query := `
		SELECT
			user_id,
			email,
			password_hash,
			role,
			created_at,
			updated_at
		FROM users
		WHERE email = $1
	`

	var user model.User

	err := r.DB.QueryRow(
		ctx,
		query,
		email,
	).Scan(
		&user.UserID,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetUserByID(
	ctx context.Context,
	userID string,
) (*model.User, error) {

	query := `
		SELECT
			user_id,
			email,
			password_hash,
			role,
			created_at,
			updated_at
		FROM users
		WHERE user_id = $1
	`

	var user model.User

	err := r.DB.QueryRow(
		ctx,
		query,
		userID,
	).Scan(
		&user.UserID,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}