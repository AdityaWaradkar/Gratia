package auth

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repo interface {
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
	StoreEmailVerification(ctx context.Context, userID, token string, expiresAt time.Time) error
	VerifyEmailToken(ctx context.Context, token string) (*User, error)
	MarkEmailVerified(ctx context.Context, userID, token string) error
}

type repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repo {
	return &repository{db: db}
}

func (r *repository) CreateUser(ctx context.Context, user *User) error {
	query := `
		INSERT INTO users (email, password_hash, role, email_verified)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at
	`
	return r.db.QueryRow(ctx, query,
		user.Email, user.PasswordHash, user.Role, user.EmailVerified,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
}

func (r *repository) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	query := `
		SELECT id, email, password_hash, role, email_verified, created_at, updated_at
		FROM users
		WHERE email = $1
	`
	user := &User{}
	err := r.db.QueryRow(ctx, query, email).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.Role,
		&user.EmailVerified, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *repository) GetUserByID(ctx context.Context, id string) (*User, error) {
	query := `
		SELECT id, email, password_hash, role, email_verified, created_at, updated_at
		FROM users
		WHERE id = $1
	`
	user := &User{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.Role,
		&user.EmailVerified, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *repository) UpdateEmailVerified(ctx context.Context, userID string) error {
	query := `
		UPDATE users
		SET email_verified = true, updated_at = $2
		WHERE id = $1
	`
	cmd, err := r.db.Exec(ctx, query, userID, time.Now())
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return errors.New("user not found")
	}
	return nil
}

func (r *repository) UpdatePassword(ctx context.Context, userID, newHash string) error {
	query := `
		UPDATE users
		SET password_hash = $2, updated_at = $3
		WHERE id = $1
	`
	cmd, err := r.db.Exec(ctx, query, userID, newHash, time.Now())
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return errors.New("user not found")
	}
	return nil
}

func (r *repository) SaveRefreshToken(ctx context.Context, rt *RefreshToken) error {
	query := `
		INSERT INTO refresh_tokens (user_id, token, expires_at, revoked, created_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`
	return r.db.QueryRow(ctx, query,
		rt.UserID, rt.TokenHash, rt.ExpiresAt, rt.Revoked, rt.CreatedAt,
	).Scan(&rt.ID)
}

func (r *repository) GetRefreshToken(ctx context.Context, token string) (*RefreshToken, error) {
	query := `
		SELECT id, user_id, token, expires_at, revoked, created_at
		FROM refresh_tokens
		WHERE token = $1
	`
	rt := &RefreshToken{}
	err := r.db.QueryRow(ctx, query, token).Scan(
		&rt.ID, &rt.UserID, &rt.TokenHash, &rt.ExpiresAt, &rt.Revoked, &rt.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return rt, nil
}

func (r *repository) RevokeRefreshToken(ctx context.Context, id string) error {
	cmd, err := r.db.Exec(ctx,
		`UPDATE refresh_tokens SET revoked = true WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return errors.New("token not found")
	}
	return nil
}

func (r *repository) SaveSession(ctx context.Context, s *Session) error {
	query := `
		INSERT INTO sessions (user_id, refresh_token_id, user_agent, ip_address, is_current, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`
	return r.db.QueryRow(ctx, query,
		s.UserID, s.RefreshTokenID, s.UserAgent, s.IPAddress, s.IsCurrent, s.CreatedAt,
	).Scan(&s.ID)
}

func (r *repository) DeleteSessionByToken(ctx context.Context, refreshTokenID string) error {
	cmd, err := r.db.Exec(ctx,
		`DELETE FROM sessions WHERE refresh_token_id = $1`, refreshTokenID)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return errors.New("session not found")
	}
	return nil
}

func (r *repository) SavePasswordReset(ctx context.Context, pr *PasswordReset) error {
	query := `
		INSERT INTO password_resets (user_id, token, expires_at, used, created_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`
	return r.db.QueryRow(ctx, query,
		pr.UserID, pr.Token, pr.ExpiresAt, pr.Used, pr.CreatedAt,
	).Scan(&pr.ID)
}

func (r *repository) StoreResetToken(ctx context.Context, userID, token string) error {
	now := time.Now()
	exp := now.Add(time.Hour)

	_, err := r.db.Exec(ctx,
		`UPDATE password_resets SET used = true, updated_at = $2 WHERE user_id = $1 AND used = false`,
		userID, now)
	if err != nil {
		return err
	}

	_, err = r.db.Exec(ctx,
		`INSERT INTO password_resets (user_id, token, expires_at, used, created_at, updated_at)
		 VALUES ($1, $2, $3, false, $4, $5)`,
		userID, token, exp, now, now)

	return err
}

func (r *repository) FindByResetToken(ctx context.Context, token string) (*User, error) {
	query := `
		SELECT u.id, u.email, u.password_hash, u.role, u.email_verified, u.created_at, u.updated_at
		FROM users u
		INNER JOIN password_resets pr ON pr.user_id = u.id
		WHERE pr.token = $1 AND pr.used = false AND pr.expires_at > NOW()
	`
	user := &User{}
	err := r.db.QueryRow(ctx, query, token).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.Role,
		&user.EmailVerified, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *repository) ClearResetToken(ctx context.Context, userID string) error {
	_, err := r.db.Exec(ctx,
		`UPDATE password_resets SET used = true, updated_at = $2 WHERE user_id = $1 AND used = false`,
		userID, time.Now())
	return err
}

func (r *repository) StoreEmailVerification(ctx context.Context, userID, token string, expiresAt time.Time) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO email_verifications (user_id, token, expires_at, created_at, used)
		 VALUES ($1, $2, $3, $4, false)`,
		userID, token, expiresAt, time.Now())
	return err
}

func (r *repository) VerifyEmailToken(ctx context.Context, token string) (*User, error) {
	query := `
		SELECT u.id, u.email, u.password_hash, u.role, u.email_verified, u.created_at, u.updated_at
		FROM users u
		INNER JOIN email_verifications ev ON ev.user_id = u.id
		WHERE ev.token = $1 AND ev.used = false AND ev.expires_at > NOW()
	`
	user := &User{}
	err := r.db.QueryRow(ctx, query, token).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.Role,
		&user.EmailVerified, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *repository) MarkEmailVerified(ctx context.Context, userID, token string) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx,
		`UPDATE users SET email_verified = true, updated_at = $2 WHERE id = $1`,
		userID, time.Now())
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx,
		`UPDATE email_verifications SET used = true, updated_at = $3
		 WHERE user_id = $1 AND token = $2`,
		userID, token, time.Now())
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}
