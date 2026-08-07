package auth

import (
    "context"
    "errors"
    "time"

    "github.com/jackc/pgx/v5"
    "github.com/jackc/pgx/v5/pgxpool"
)

// Repository defines all required database operations for the authentication domain
type Repository interface {
    CreateUser(ctx context.Context, user *User) error
    GetUserByEmail(ctx context.Context, email string) (*User, error)
    GetUserByID(ctx context.Context, id string) (*User, error)
    UpdateEmailVerified(ctx context.Context, userID string) error
    UpdatePassword(ctx context.Context, userID, newHash string) error
    UpdateUserRole(ctx context.Context, userID, role string) error
    SaveRefreshToken(ctx context.Context, token *RefreshToken) error
    GetRefreshToken(ctx context.Context, token string) (*RefreshToken, error)
    RevokeRefreshToken(ctx context.Context, tokenID string) error
    SaveSession(ctx context.Context, session *Session) error
    DeleteSessionByToken(ctx context.Context, refreshTokenID string) error
    StoreResetToken(ctx context.Context, userID, token string) error
    FindByResetToken(ctx context.Context, token string) (*User, error)
    ClearResetToken(ctx context.Context, userID string) error
}

// repository implements the Repository interface using a pgxpool connection
type repository struct {
    db *pgxpool.Pool
}

// NewRepository initializes the concrete repository with an active connection pool
func NewRepository(db *pgxpool.Pool) Repository {
    return &repository{db: db}
}

// CreateUser inserts a new user record and retrieves generated metadata
func (r *repository) CreateUser(ctx context.Context, user *User) error {
    query := `
        INSERT INTO users (email, password_hash, role, email_verified)
        VALUES ($1, $2, $3, $4)
        RETURNING id, created_at, updated_at
    `
    return r.db.QueryRow(
        ctx, query, user.Email, user.PasswordHash, user.Role, user.EmailVerified,
    ).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
}

// GetUserByEmail fetches a user by email, safely handling not-found scenarios
func (r *repository) GetUserByEmail(ctx context.Context, email string) (*User, error) {
    query := `
        SELECT id, email, password_hash, role, email_verified, created_at, updated_at
        FROM users WHERE email = $1
    `
    user := &User{}
    err := r.db.QueryRow(ctx, query, email).Scan(
        &user.ID, &user.Email, &user.PasswordHash, &user.Role,
        &user.EmailVerified, &user.CreatedAt, &user.UpdatedAt,
    )
    if err != nil {
        if errors.Is(err, pgx.ErrNoRows) {
            return nil, nil
        }
        return nil, err
    }
    return user, nil
}

// GetUserByID retrieves a specific user by their unique UUID identifier
func (r *repository) GetUserByID(ctx context.Context, id string) (*User, error) {
    query := `
        SELECT id, email, password_hash, role, email_verified, created_at, updated_at
        FROM users WHERE id = $1
    `
    user := &User{}
    err := r.db.QueryRow(ctx, query, id).Scan(
        &user.ID, &user.Email, &user.PasswordHash, &user.Role,
        &user.EmailVerified, &user.CreatedAt, &user.UpdatedAt,
    )
    if err != nil {
        if errors.Is(err, pgx.ErrNoRows) {
            return nil, nil
        }
        return nil, err
    }
    return user, nil
}

// UpdateEmailVerified flags a user's email address as successfully verified
func (r *repository) UpdateEmailVerified(ctx context.Context, userID string) error {
    cmd, err := r.db.Exec(
        ctx,
        `UPDATE users SET email_verified = true, updated_at = $2 WHERE id = $1`,
        userID, time.Now(),
    )
    if err != nil {
        return err
    }
    if cmd.RowsAffected() == 0 {
        return errors.New("user not found")
    }
    return nil
}

// UpdatePassword securely replaces a user's cryptographic password hash
func (r *repository) UpdatePassword(ctx context.Context, userID, newHash string) error {
    cmd, err := r.db.Exec(
        ctx,
        `UPDATE users SET password_hash = $2, updated_at = $3 WHERE id = $1`,
        userID, newHash, time.Now(),
    )
    if err != nil {
        return err
    }
    if cmd.RowsAffected() == 0 {
        return errors.New("user not found")
    }
    return nil
}

// UpdateUserRole elevates or demotes a user's platform access level
func (r *repository) UpdateUserRole(ctx context.Context, userID, role string) error {
    cmd, err := r.db.Exec(
        ctx,
        `UPDATE users SET role = $2, updated_at = $3 WHERE id = $1`,
        userID, role, time.Now(),
    )
    if err != nil {
        return err
    }
    if cmd.RowsAffected() == 0 {
        return errors.New("user not found")
    }
    return nil
}

// SaveRefreshToken persists long-lived session credentials to the database
func (r *repository) SaveRefreshToken(ctx context.Context, token *RefreshToken) error {
    query := `
        INSERT INTO refresh_tokens (user_id, token, expires_at, revoked, created_at)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id
    `
    return r.db.QueryRow(
        ctx, query, token.UserID, token.Token, token.ExpiresAt, token.Revoked, token.CreatedAt,
    ).Scan(&token.ID)
}

// GetRefreshToken retrieves token metadata for validation during rotation
func (r *repository) GetRefreshToken(ctx context.Context, token string) (*RefreshToken, error) {
    query := `
        SELECT id, user_id, token, expires_at, revoked, created_at
        FROM refresh_tokens WHERE token = $1
    `
    rt := &RefreshToken{}
    err := r.db.QueryRow(ctx, query, token).Scan(
        &rt.ID, &rt.UserID, &rt.Token, &rt.ExpiresAt, &rt.Revoked, &rt.CreatedAt,
    )
    if err != nil {
        if errors.Is(err, pgx.ErrNoRows) {
            return nil, nil
        }
        return nil, err
    }
    return rt, nil
}

// RevokeRefreshToken permanently invalidates a specific rotation token
func (r *repository) RevokeRefreshToken(ctx context.Context, tokenID string) error {
    cmd, err := r.db.Exec(
        ctx,
        `UPDATE refresh_tokens SET revoked = true WHERE id = $1`,
        tokenID,
    )
    if err != nil {
        return err
    }
    if cmd.RowsAffected() == 0 {
        return errors.New("token not found")
    }
    return nil
}

// SaveSession logs active device and network details for security tracking
func (r *repository) SaveSession(ctx context.Context, session *Session) error {
    query := `
        INSERT INTO sessions (user_id, refresh_token_id, user_agent, ip_address, is_current, created_at)
        VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING id
    `
    return r.db.QueryRow(
        ctx, query, session.UserID, session.RefreshTokenID, session.UserAgent, session.IPAddress, session.IsCurrent, session.CreatedAt,
    ).Scan(&session.ID)
}

// DeleteSessionByToken completely purges device session metadata upon logout
func (r *repository) DeleteSessionByToken(ctx context.Context, refreshTokenID string) error {
    cmd, err := r.db.Exec(
        ctx,
        `DELETE FROM sessions WHERE refresh_token_id = $1`,
        refreshTokenID,
    )
    if err != nil {
        return err
    }
    if cmd.RowsAffected() == 0 {
        return errors.New("session not found")
    }
    return nil
}

// StoreResetToken safely registers a new password recovery attempt
func (r *repository) StoreResetToken(ctx context.Context, userID, token string) error {
    now := time.Now()
    _, err := r.db.Exec(
        ctx,
        `UPDATE password_resets SET used = true, updated_at = $2 WHERE user_id = $1 AND used = false`,
        userID, now,
    )
    if err != nil {
        return err
    }
    _, err = r.db.Exec(
        ctx,
        `INSERT INTO password_resets (user_id, token, expires_at, used, created_at, updated_at)
         VALUES ($1, $2, $3, false, $4, $5)`,
        userID, token, now.Add(time.Hour), now, now,
    )
    return err
}

// FindByResetToken strictly verifies a password recovery token's validity
func (r *repository) FindByResetToken(ctx context.Context, token string) (*User, error) {
    query := `
        SELECT u.id, u.email, u.password_hash, u.role, u.email_verified, u.created_at, u.updated_at
        FROM users u
        JOIN password_resets pr ON pr.user_id = u.id
        WHERE pr.token = $1 AND pr.used = false AND pr.expires_at > NOW()
    `
    user := &User{}
    err := r.db.QueryRow(ctx, query, token).Scan(
        &user.ID, &user.Email, &user.PasswordHash, &user.Role,
        &user.EmailVerified, &user.CreatedAt, &user.UpdatedAt,
    )
    if err != nil {
        if errors.Is(err, pgx.ErrNoRows) {
            return nil, nil
        }
        return nil, err
    }
    return user, nil
}

// ClearResetToken consumes a recovery token immediately after successful use
func (r *repository) ClearResetToken(ctx context.Context, userID string) error {
    _, err := r.db.Exec(
        ctx,
        `UPDATE password_resets SET used = true, updated_at = $2 WHERE user_id = $1 AND used = false`,
        userID, time.Now(),
    )
    return err
}