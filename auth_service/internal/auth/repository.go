package auth

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository handles database operations for the auth service
type Repository struct {
	db *pgxpool.Pool
}

// NewRepository creates a new Repository instance
func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

// CreateUser inserts a new user into the database
func (r *Repository) CreateUser(ctx context.Context, user *User) error {
	query := `
		INSERT INTO users (email, password_hash, role, email_verified)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at
	`
	return r.db.QueryRow(ctx, query,
		user.Email,
		user.PasswordHash,
		user.Role,
		user.EmailVerified,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
}

// GetUserByEmail finds a user by email
func (r *Repository) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	query := `
		SELECT id, email, password_hash, role, email_verified, created_at, updated_at
		FROM users
		WHERE email = $1
	`
	user := &User{}
	err := r.db.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
		&user.EmailVerified,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// GetUserByID finds a user by ID
func (r *Repository) GetUserByID(ctx context.Context, id string) (*User, error) {
	query := `
		SELECT id, email, password_hash, role, email_verified, created_at, updated_at
		FROM users
		WHERE id = $1
	`
	user := &User{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
		&user.EmailVerified,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// UpdateEmailVerified sets email_verified = true
func (r *Repository) UpdateEmailVerified(ctx context.Context, userID string) error {
	query := `
		UPDATE users
		SET email_verified = true, updated_at = $2
		WHERE id = $1
	`
	cmdTag, err := r.db.Exec(ctx, query, userID, time.Now())
	if err != nil {
		return err
	}
	if cmdTag.RowsAffected() == 0 {
		return errors.New("user not found")
	}
	return nil
}

// UpdatePassword updates a user’s password
func (r *Repository) UpdatePassword(ctx context.Context, userID string, newHash string) error {
	query := `
		UPDATE users
		SET password_hash = $2, updated_at = $3
		WHERE id = $1
	`
	cmdTag, err := r.db.Exec(ctx, query, userID, newHash, time.Now())
	if err != nil {
		return err
	}
	if cmdTag.RowsAffected() == 0 {
		return errors.New("user not found")
	}
	return nil
}
