package repository

import (
	"context"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/adityawaradkar/gratia-auth/internal/models"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *models.User) error
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	GetUserByID(ctx context.Context, userID uuid.UUID) (*models.User, error)
}

type userRepo struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) UserRepository {
	return &userRepo{db: db}
}

func (r *userRepo) CreateUser(ctx context.Context, user *models.User) error {
	query := `
		INSERT INTO public.users 
		(user_id, email, password_hash, full_name, phone_number, role, location_lat, location_lng, is_active, created_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,NOW())
	`
	if user.UserID == uuid.Nil {
		user.UserID = uuid.New()
	}

	_, err := r.db.ExecContext(ctx, query,
		user.UserID,
		user.Email,
		user.PasswordHash,
		user.FullName,
		user.PhoneNumber,
		user.Role,
		user.LocationLat,
		user.LocationLng,
		user.IsActive,
	)
	return err
}

func (r *userRepo) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	query := `SELECT * FROM public.users WHERE email = $1`
	err := r.db.GetContext(ctx, &user, query, email)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepo) GetUserByID(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	var user models.User
	query := `SELECT * FROM public.users WHERE user_id = $1`
	err := r.db.GetContext(ctx, &user, query, userID)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
