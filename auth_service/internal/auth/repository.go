package auth

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// RepositoryInterface defines all methods the service expects
type RepositoryInterface interface {
	CreateUser(ctx context.Context, user *User) error
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	GetUserByID(ctx context.Context, id string) (*User, error)
	UpdateEmailVerified(ctx context.Context, userID string) error
	UpdatePassword(ctx context.Context, userID, newHash string) error
	SaveRefreshToken(ctx context.Context, rt *RefreshToken) error
	GetRefreshToken(ctx context.Context, token string) (*RefreshToken, error)
	RevokeRefreshToken(ctx context.Context, id string) error
	SaveSession(ctx context.Context, s *Session) error
	DeleteSessionByToken(ctx context.Context, refreshTokenID string) error
	SavePasswordReset(ctx context.Context, pr *PasswordReset) error
	StoreResetToken(ctx context.Context, userID, token string) error
	FindByResetToken(ctx context.Context, token string) (*User, error)
	ClearResetToken(ctx context.Context, userID string) error
}

// Repository handles database operations
type Repository struct {
	db *pgxpool.Pool
}

// NewRepository creates a new Repository instance
func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

// ---------------- Users ----------------

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

func (r *Repository) UpdateEmailVerified(ctx context.Context, userID string) error {
	query := `UPDATE users SET email_verified = true, updated_at = $2 WHERE id = $1`
	cmdTag, err := r.db.Exec(ctx, query, userID, time.Now())
	if err != nil {
		return err
	}
	if cmdTag.RowsAffected() == 0 {
		return errors.New("user not found")
	}
	return nil
}

func (r *Repository) UpdatePassword(ctx context.Context, userID, newHash string) error {
	query := `UPDATE users SET password_hash = $2, updated_at = $3 WHERE id = $1`
	cmdTag, err := r.db.Exec(ctx, query, userID, newHash, time.Now())
	if err != nil {
		return err
	}
	if cmdTag.RowsAffected() == 0 {
		return errors.New("user not found")
	}
	return nil
}

// ---------------- Refresh Tokens ----------------

func (r *Repository) SaveRefreshToken(ctx context.Context, rt *RefreshToken) error {
	query := `
		INSERT INTO refresh_tokens (user_id, token, expires_at, revoked, created_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`
	return r.db.QueryRow(ctx, query,
		rt.UserID,
		rt.TokenHash,
		rt.ExpiresAt,
		rt.Revoked,
		rt.CreatedAt,
	).Scan(&rt.ID)
}

func (r *Repository) GetRefreshToken(ctx context.Context, token string) (*RefreshToken, error) {
	query := `
		SELECT id, user_id, token, expires_at, revoked, created_at
		FROM refresh_tokens
		WHERE token = $1
	`
	rt := &RefreshToken{}
	err := r.db.QueryRow(ctx, query, token).Scan(
		&rt.ID,
		&rt.UserID,
		&rt.TokenHash,
		&rt.ExpiresAt,
		&rt.Revoked,
		&rt.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return rt, nil
}

func (r *Repository) RevokeRefreshToken(ctx context.Context, id string) error {
	query := `UPDATE refresh_tokens SET revoked = true WHERE id = $1`
	cmdTag, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	if cmdTag.RowsAffected() == 0 {
		return errors.New("token not found")
	}
	return nil
}

// ---------------- Sessions ----------------

func (r *Repository) SaveSession(ctx context.Context, s *Session) error {
	query := `
		INSERT INTO sessions (user_id, refresh_token_id, user_agent, ip_address, is_current, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`
	return r.db.QueryRow(ctx, query,
		s.UserID,
		s.RefreshTokenID,
		s.UserAgent,
		s.IPAddress,
		s.IsCurrent,
		s.CreatedAt,
	).Scan(&s.ID)
}

func (r *Repository) DeleteSessionByToken(ctx context.Context, refreshTokenID string) error {
	query := `DELETE FROM sessions WHERE refresh_token_id = $1`
	cmdTag, err := r.db.Exec(ctx, query, refreshTokenID)
	if err != nil {
		return err
	}
	if cmdTag.RowsAffected() == 0 {
		return errors.New("session not found")
	}
	return nil
}

// ---------------- Password Reset ----------------

func (r *Repository) SavePasswordReset(ctx context.Context, pr *PasswordReset) error {
	query := `
		INSERT INTO password_resets (user_id, token, expires_at, used, created_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`
	return r.db.QueryRow(ctx, query,
		pr.UserID,
		pr.Token,
		pr.ExpiresAt,
		pr.Used,
		pr.CreatedAt,
	).Scan(&pr.ID)
}

// StoreResetToken inserts a new reset token, invalidates old tokens, and logs expiry
func (r *Repository) StoreResetToken(ctx context.Context, userID, token string) error {
	now := time.Now()
	expiresAt := now.Add(1 * time.Hour) // token valid for 1 hour

	// Invalidate previous unused tokens
	cmdTag, err := r.db.Exec(ctx, `
		UPDATE password_resets
		SET used = true, updated_at = $2
		WHERE user_id = $1 AND used = false
	`, userID, now)
	if err != nil {
		return err
	}
	if cmdTag.RowsAffected() > 0 {
		log.Printf("ℹ️ Invalidated %d previous reset token(s) for user %s", cmdTag.RowsAffected(), userID)
	}

	// Insert the new token
	query := `
		INSERT INTO password_resets (user_id, token, expires_at, used, created_at, updated_at)
		VALUES ($1, $2, $3, false, $4, $5)
		RETURNING id, expires_at, used
	`
	var id string
	var storedExpires time.Time
	var used bool
	err = r.db.QueryRow(ctx, query, userID, token, expiresAt, now, now).Scan(&id, &storedExpires, &used)
	if err != nil {
		return err
	}

	log.Printf("✅ Created new reset token for user %s: %s (expires at %s, used=%v)", userID, token, storedExpires.Format(time.RFC3339), used)
	return nil
}

// FindByResetToken retrieves a user by a valid, unused reset token and logs token status
func (r *Repository) FindByResetToken(ctx context.Context, token string) (*User, error) {
	query := `
		SELECT u.id, u.email, u.password_hash, u.role, u.email_verified, u.created_at, u.updated_at
		FROM users u
		INNER JOIN password_resets pr ON pr.user_id = u.id
		WHERE pr.token = $1 AND pr.used = false AND pr.expires_at > NOW()
	`
	user := &User{}
	err := r.db.QueryRow(ctx, query, token).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
		&user.EmailVerified,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		log.Printf("❌ Failed to find user by reset token '%s': %v", token, err)
		return nil, err
	}
	log.Printf("✅ Found user %s by reset token '%s'", user.ID, token)
	return user, nil
}

func (r *Repository) ClearResetToken(ctx context.Context, userID string) error {
	query := `
		UPDATE password_resets
		SET used = true, updated_at = $2
		WHERE user_id = $1 AND used = false
	`
	_, err := r.db.Exec(ctx, query, userID, time.Now())
	return err
}
