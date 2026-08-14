## Table of Contents
1. [Overview](#overview)
2. [Architecture](#architecture)
3. [Configuration](#configuration)
4. [Database Schema](#database-schema)
5. [Core Components](#core-components)
6. [API Endpoints](#api-endpoints)
7. [Security Implementation](#security-implementation)
8. [Error Handling](#error-handling)
9. [Logging & Monitoring](#logging--monitoring)
10. [Deployment](#deployment)
11. [Development Guidelines](#development-guidelines)

---

## Overview

The Auth Service is a production-ready, microservices-based authentication system built in Go. It provides comprehensive user management, secure authentication, session handling, and role-based access control (RBAC). This service is designed to be the foundational authentication layer for the Gratia platform, with strict security practices and scalability considerations.

### Key Features
- **User Registration & Authentication** with email/password
- **JWT-based Authentication** with access and refresh tokens
- **Role-Based Access Control** (Donor, NGO, Admin, User)
- **Session Management** with device tracking
- **Password Reset Flow** with secure tokens
- **Graceful Shutdown** and health checking
- **Structured JSON Logging** for observability
- **Docker Ready** with multi-stage builds

---

## Architecture

### High-Level Architecture
```
┌─────────────────────────────────────────────────────────────┐
│                     HTTP Client                             │
└─────────────────────┬───────────────────────────────────────┘
                      │
┌─────────────────────▼───────────────────────────────────────┐
│                    HTTP Server                              │
│  ┌─────────────────────────────────────────────────────┐   │
│  │              Router (gorilla/mux)                   │   │
│  │  ┌─────────────────────────────────────────────┐   │   │
│  │  │        Authentication Middleware             │   │   │
│  │  │    (JWT Validation & Context Injection)      │   │   │
│  │  └─────────────────────────────────────────────┘   │   │
│  └─────────────────────────────────────────────────────┘   │
└─────────────────────┬───────────────────────────────────────┘
                      │
┌─────────────────────▼───────────────────────────────────────┐
│              HTTP Handlers (auth/handlers.go)               │
│  ┌─────────────────────────────────────────────────────┐   │
│  │  RegisterUser | LoginUser | RefreshTokens          │   │
│  │  Logout | ForgotPassword | ResetPassword           │   │
│  │  GetCurrentUser | HealthCheck                      │   │
│  └─────────────────────────────────────────────────────┘   │
└─────────────────────┬───────────────────────────────────────┘
                      │
┌─────────────────────▼───────────────────────────────────────┐
│                Service Layer (auth/service.go)              │
│  ┌─────────────────────────────────────────────────────┐   │
│  │  Business Logic & Validation                        │   │
│  │  Password Hashing | Token Generation                │   │
│  │  Role Enforcement | Session Management              │   │
│  └─────────────────────────────────────────────────────┘   │
└─────────────────────┬───────────────────────────────────────┘
                      │
┌─────────────────────▼───────────────────────────────────────┐
│              Repository Layer (auth/repository.go)          │
│  ┌─────────────────────────────────────────────────────┐   │
│  │  Database Operations (CRUD)                         │   │
│  │  User Management | Token Storage                    │   │
│  │  Session Tracking | Password Reset                  │   │
│  └─────────────────────────────────────────────────────┘   │
└─────────────────────┬───────────────────────────────────────┘
                      │
┌─────────────────────▼───────────────────────────────────────┐
│          PostgreSQL Database (pgxpool Connection)           │
│  ┌─────────────────────────────────────────────────────┐   │
│  │  users | refresh_tokens | sessions                  │   │
│  │  password_resets                                    │   │
│  └─────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────┘
```

### Directory Structure
```
auth_service/
├── cmd/
│   └── server/
│       └── main.go               # Application entry point
├── internal/
│   ├── auth/
│   │   ├── handlers.go           # HTTP handlers
│   │   ├── service.go            # Business logic
│   │   ├── repository.go         # Data access layer
│   │   ├── models.go             # Domain models
│   │   ├── role.go               # Role definitions
│   │   └── validator.go          # Input validation
│   ├── config/
│   │   └── config.go             # Configuration management
│   ├── db/
│   │   └── db.go                 # Database connection
│   ├── logger/
│   │   └── logger.go             # Structured logging
│   ├── middleware/
│   │   └── auth_middleware.go    # JWT authentication
│   ├── server/
│   │   └── http.go               # Route registration
│   └── utils/
│       ├── hash.go               # Password hashing
│       └── token.go              # Token generation
├── Dockerfile                     # Multi-stage build
├── go.mod
├── go.sum
└── .env.example                   # Environment variables template
```

---

## Configuration

### Environment Variables

| Variable | Type | Required | Default | Description |
|----------|------|----------|---------|-------------|
| `PORT` | string | No | `8080` | HTTP server port |
| `DATABASE_URL` | string | Yes | - | PostgreSQL connection string |
| `JWT_SECRET` | string | Yes | - | Secret key for JWT signing (minimum 32 characters) |
| `ACCESS_TOKEN_MINUTES` | int | No | `15` | JWT access token validity in minutes |
| `REFRESH_TOKEN_DAYS` | int | No | `30` | Refresh token validity in days |
| `ENV` | string | No | `development` | Environment (development/staging/production) |
| `LOG_LEVEL` | string | No | `info` | Log level (debug/info/warn/error) |

### Configuration Loading Process

```go
// config/config.go
func Load() *Config {
    _ = godotenv.Load() // Optional .env file loading
    
    return &Config{
        Port:            getEnv("PORT", "8080"),
        DatabaseURL:     mustEnv("DATABASE_URL"), // Required
        JWTSecret:       mustEnv("JWT_SECRET"),   // Required
        AccessTokenTTL:  time.Duration(getIntEnv("ACCESS_TOKEN_MINUTES", 15)) * time.Minute,
        RefreshTokenTTL: time.Duration(getIntEnv("REFRESH_TOKEN_DAYS", 30)) * 24 * time.Hour,
        Env:             getEnv("ENV", "development"),
        LogLevel:        getEnv("LOG_LEVEL", "info"),
    }
}
```

### Environment-Specific Configuration

**Development:**
```env
ENV=development
LOG_LEVEL=debug
ACCESS_TOKEN_MINUTES=15
REFRESH_TOKEN_DAYS=30
```

**Production:**
```env
ENV=production
LOG_LEVEL=info
ACCESS_TOKEN_MINUTES=5
REFRESH_TOKEN_DAYS=7
```

---

## Database Schema

### SQL Migration (PostgreSQL)

```sql
-- Users table
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL DEFAULT 'USER',
    email_verified BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Refresh tokens table
CREATE TABLE refresh_tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token VARCHAR(255) UNIQUE NOT NULL,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    revoked BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    INDEX idx_refresh_tokens_user_id (user_id),
    INDEX idx_refresh_tokens_token (token)
);

-- Sessions table
CREATE TABLE sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    refresh_token_id UUID NOT NULL REFERENCES refresh_tokens(id) ON DELETE CASCADE,
    user_agent TEXT,
    ip_address VARCHAR(45),
    is_current BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    INDEX idx_sessions_user_id (user_id),
    INDEX idx_sessions_refresh_token_id (refresh_token_id)
);

-- Password reset tokens table
CREATE TABLE password_resets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token VARCHAR(255) UNIQUE NOT NULL,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    used BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    INDEX idx_password_resets_user_id (user_id),
    INDEX idx_password_resets_token (token)
);
```

### Entity Relationship Diagram

```
┌─────────────────┐
│     users       │
│─────────────────│
│ id (PK)         │───┐
│ email           │   │
│ password_hash   │   │
│ role            │   │
│ email_verified  │   │
│ created_at      │   │
│ updated_at      │   │
└─────────────────┘   │
        │             │
        │             │
        │             │
┌───────▼─────────────┼──────┐
│ refresh_tokens      │      │
│─────────────────────│      │
│ id (PK)            │      │
│ user_id (FK)───────┘      │
│ token (UK)                │
│ expires_at                │
│ revoked                   │
│ created_at                │
└───────────────────────────┘
        │
        │
        │
┌───────▼──────────┐
│    sessions      │
│──────────────────│
│ id (PK)          │
│ user_id (FK)─────┘
│ refresh_token_id │
│ user_agent       │
│ ip_address       │
│ is_current       │
│ created_at       │
└──────────────────┘

┌─────────────────┐
│ password_resets │
│─────────────────│
│ id (PK)         │
│ user_id (FK)────┘
│ token (UK)      │
│ expires_at      │
│ used            │
│ created_at      │
│ updated_at      │
└─────────────────┘
```

---

## Core Components

### 1. Domain Models (`internal/auth/models.go`)

#### User
```go
type User struct {
    ID            string    `json:"id"`
    Email         string    `json:"email"`
    PasswordHash  string    `json:"-"`          // Never serialized to JSON
    Role          Role      `json:"role"`
    EmailVerified bool      `json:"emailVerified"`
    CreatedAt     time.Time `json:"createdAt"`
    UpdatedAt     time.Time `json:"updatedAt"`
}
```

#### RefreshToken
```go
type RefreshToken struct {
    ID        string    `json:"id"`
    UserID    string    `json:"userId"`
    Token     string    `json:"token"`
    ExpiresAt time.Time `json:"expiresAt"`
    Revoked   bool      `json:"revoked"`
    CreatedAt time.Time `json:"createdAt"`
}
```

#### Session
```go
type Session struct {
    ID             string    `json:"id"`
    UserID         string    `json:"userId"`
    RefreshTokenID string    `json:"refreshTokenId"`
    UserAgent      string    `json:"userAgent"`
    IPAddress      string    `json:"ipAddress"`
    IsCurrent      bool      `json:"isCurrent"`
    CreatedAt      time.Time `json:"createdAt"`
}
```

#### PasswordReset
```go
type PasswordReset struct {
    ID        string    `json:"id"`
    UserID    string    `json:"userId"`
    Token     string    `json:"-"`            // Never serialized to JSON
    ExpiresAt time.Time `json:"expiresAt"`
    Used      bool      `json:"used"`
    CreatedAt time.Time `json:"createdAt"`
    UpdatedAt time.Time `json:"updatedAt"`
}
```

### 2. Role Management (`internal/auth/role.go`)

```go
type Role string

const (
    RoleDonor Role = "DONOR"
    RoleNGO   Role = "NGO"
    RoleAdmin Role = "ADMIN"
    RoleUser  Role = "USER"
)

var ValidRoles = map[Role]bool{
    RoleDonor: true,
    RoleNGO:   true,
    RoleAdmin: true,
    RoleUser:  true,
}
```

**Role Hierarchy:**
- **USER**: Default role, basic platform access
- **DONOR**: Can make donations and view campaigns
- **NGO**: Can create campaigns and manage donations
- **ADMIN**: Full system access, can manage users and roles

### 3. Service Layer (`internal/auth/service.go`)

The service layer implements all business logic:

#### RegisterUser
```go
func (s *Service) RegisterUser(ctx context.Context, input RegisterInput, callerRole Role) (*User, error)
```
- Validates email format and password strength
- Checks for existing users
- Enforces role assignment permissions
- Hashes password using bcrypt
- Creates user with `email_verified = false`

#### LoginUser
```go
func (s *Service) LoginUser(ctx context.Context, input LoginInput, userAgent, ip string) (*TokenPair, error)
```
- Validates credentials
- Generates JWT access token (short-lived)
- Creates refresh token (long-lived)
- Creates session with device metadata
- Returns token pair

#### RefreshTokens
```go
func (s *Service) RefreshTokens(ctx context.Context, refreshToken, userAgent, ip string) (*TokenPair, error)
```
- Validates refresh token (not revoked, not expired)
- Revokes old refresh token
- Creates new refresh token and session
- Returns new token pair

#### Logout
```go
func (s *Service) Logout(ctx context.Context, refreshToken string) error
```
- Revokes refresh token
- Deletes associated session

#### Password Reset Flow
```go
func (s *Service) GenerateResetToken(ctx context.Context, email string) (string, error)
func (s *Service) ResetPassword(ctx context.Context, resetToken, newPassword string) error
```
- Generates UUID-based reset token
- Stores token with 1-hour expiry
- Validates token and resets password

### 4. Repository Layer (`internal/auth/repository.go`)

Implements data access operations:

```go
type Repository interface {
    // User operations
    CreateUser(ctx context.Context, user *User) error
    GetUserByEmail(ctx context.Context, email string) (*User, error)
    GetUserByID(ctx context.Context, id string) (*User, error)
    UpdateEmailVerified(ctx context.Context, userID string) error
    UpdatePassword(ctx context.Context, userID, newHash string) error
    UpdateUserRole(ctx context.Context, userID, role string) error
    
    // Refresh token operations
    SaveRefreshToken(ctx context.Context, token *RefreshToken) error
    GetRefreshToken(ctx context.Context, token string) (*RefreshToken, error)
    RevokeRefreshToken(ctx context.Context, tokenID string) error
    
    // Session operations
    SaveSession(ctx context.Context, session *Session) error
    DeleteSessionByToken(ctx context.Context, refreshTokenID string) error
    
    // Password reset operations
    StoreResetToken(ctx context.Context, userID, token string) error
    FindByResetToken(ctx context.Context, token string) (*User, error)
    ClearResetToken(ctx context.Context, userID string) error
}
```

### 5. HTTP Handlers (`internal/auth/handlers.go`)

Manages HTTP request/response:

```go
type Handler struct {
    service *Service
}

// Handler methods
func (h *Handler) RegisterUser(w http.ResponseWriter, r *http.Request)
func (h *Handler) LoginUser(w http.ResponseWriter, r *http.Request)
func (h *Handler) RefreshTokens(w http.ResponseWriter, r *http.Request)
func (h *Handler) Logout(w http.ResponseWriter, r *http.Request)
func (h *Handler) ForgotPassword(w http.ResponseWriter, r *http.Request)
func (h *Handler) ResetPassword(w http.ResponseWriter, r *http.Request)
func (h *Handler) GetCurrentUser(w http.ResponseWriter, r *http.Request)
func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request)
```

---

## API Endpoints

### Base URL: `http://localhost:8080`

### Public Endpoints (No Authentication Required)

#### 1. Health Check
```http
GET /health
```
**Response:**
```json
{
  "status": "healthy",
  "service": "auth_service"
}
```

#### 2. Register User
```http
POST /register
```
**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "SecurePassword123",
  "role": "DONOR"  // Optional, defaults to "USER"
}
```
**Response (201 Created):**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "email": "user@example.com",
  "role": "DONOR",
  "emailVerified": false,
  "createdAt": "2026-08-08T10:00:00Z",
  "updatedAt": "2026-08-08T10:00:00Z"
}
```

#### 3. Login User
```http
POST /login
```
**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "SecurePassword123"
}
```
**Response (200 OK):**
```json
{
  "accessToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refreshToken": "m7x9k3p5q8r2t6w1y4u7a1n5e9c3g6v0"
}
```

#### 4. Refresh Tokens
```http
POST /refresh
```
**Request Body:**
```json
{
  "refreshToken": "m7x9k3p5q8r2t6w1y4u7a1n5e9c3g6v0"
}
```
**Response (200 OK):**
```json
{
  "accessToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refreshToken": "n8y0l4q6s9u3v7x2z5w8b2o6f0d4h7w1"
}
```

#### 5. Forgot Password
```http
POST /forgot-password
```
**Request Body:**
```json
{
  "email": "user@example.com"
}
```
**Response (200 OK - Always returns success to prevent user enumeration):**
```json
{
  "message": "if the email exists, a reset link was sent"
}
```
**Development Response (Includes token for testing):**
```json
{
  "resetToken": "550e8400-e29b-41d4-a716-446655440000"
}
```

#### 6. Reset Password
```http
POST /reset-password
```
**Request Body:**
```json
{
  "token": "550e8400-e29b-41d4-a716-446655440000",
  "newPassword": "NewSecurePassword123"
}
```
**Response:** 204 No Content

### Protected Endpoints (Requires Authentication)

#### 7. Get Current User
```http
GET /me
```
**Headers:**
```
Authorization: Bearer <accessToken>
```
**Response (200 OK):**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "email": "user@example.com",
  "role": "DONOR",
  "emailVerified": false,
  "createdAt": "2026-08-08T10:00:00Z",
  "updatedAt": "2026-08-08T10:00:00Z"
}
```

#### 8. Logout
```http
POST /logout
```
**Headers:**
```
Authorization: Bearer <accessToken>
```
**Request Body:**
```json
{
  "refreshToken": "m7x9k3p5q8r2t6w1y4u7a1n5e9c3g6v0"
}
```
**Response:** 204 No Content

---

## Security Implementation

### 1. Password Security

- **Hashing Algorithm**: bcrypt with cost factor 10
- **Minimum Requirements**: At least 8 characters
- **Storage**: Password hashes only, never plain text
- **Implementation**:
```go
// internal/utils/hash.go
func HashPassword(password string) (string, error) {
    hash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
    if err != nil {
        return "", err
    }
    return string(hash), nil
}
```

### 2. JWT Token Security

#### Access Token
- **Algorithm**: HS256 (HMAC-SHA256)
- **Claims**: user_id, role, exp, iat
- **Lifetime**: Configurable (default: 15 minutes)
- **Stateless**: No server-side storage required
- **Implementation**:
```go
// internal/utils/token.go
func GenerateAccessToken(userID string, role string, secret string, ttl time.Duration) (string, error) {
    claims := jwt.MapClaims{
        "sub":  userID,
        "role": role,
        "exp":  time.Now().Add(ttl).Unix(),
        "iat":  time.Now().Unix(),
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(secret))
}
```

#### Refresh Token
- **Format**: Cryptographically random 32-byte tokens (base64 URL encoded)
- **Storage**: Server-side with database record
- **Lifetime**: Configurable (default: 30 days)
- **Features**: Revocable, trackable
- **Implementation**:
```go
// internal/utils/token.go
func GenerateSecureToken(length int) (string, error) {
    b := make([]byte, length)
    if _, err := rand.Read(b); err != nil {
        return "", err
    }
    return base64.RawURLEncoding.EncodeToString(b), nil
}
```

### 3. Authentication Middleware

The middleware authenticates requests by validating JWT tokens:

```go
// internal/middleware/auth_middleware.go
func Auth(jwtSecret string) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            // Extract and validate token
            // Extract user information from claims
            // Inject into context
            // Continue to handler
        })
    }
}
```

**Context Values:**
- `user_id`: User's UUID
- `user_email`: User's email
- `user_role`: User's role

### 4. Session Management

- **Device Tracking**: Captures User-Agent and IP address
- **Active Sessions**: Each refresh token has an associated session
- **Session Revocation**: Logout revokes both token and session

### 5. Role-Based Access Control (RBAC)

**Role Assignment Rules:**
- Users with `USER` role can only create `USER` and `DONOR` accounts
- Users with `DONOR` role can only create `DONOR` and `USER` accounts
- Users with `NGO` role can only create `NGO` and `USER` accounts
- Users with `ADMIN` role can create any role

### 6. Input Validation

```go
// internal/auth/validator.go
func ValidateEmail(email string) error {
    // RFC 5322 compliant email validation
    emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
    // ...
}

func ValidatePassword(password string) error {
    if len(password) < 8 {
        return errors.New("password must be at least 8 characters")
    }
    return nil
}
```

---

## Error Handling

### Error Response Format

All errors follow a consistent format:
```json
{
  "error": "descriptive error message"
}
```

### HTTP Status Codes

| Status Code | Description |
|-------------|-------------|
| 200 | Success (GET/POST) |
| 201 | Created (POST) |
| 204 | No Content (Successful DELETE) |
| 400 | Bad Request (Invalid JSON/Missing fields) |
| 401 | Unauthorized (Invalid/Expired token) |
| 403 | Forbidden (Insufficient permissions) |
| 404 | Not Found |
| 409 | Conflict (Duplicate email) |
| 500 | Internal Server Error |

### Common Error Messages

| Error | Scenario |
|-------|----------|
| "invalid json payload" | Malformed JSON in request body |
| "email required" | Missing email field |
| "invalid email format" | Email doesn't match regex pattern |
| "password must be at least 8 characters" | Password too short |
| "user already exists" | Duplicate email registration |
| "invalid credentials" | Wrong email/password combination |
| "invalid refresh token" | Expired or revoked refresh token |
| "invalid reset token" | Expired or used reset token |
| "unauthorized role assignment" | Insufficient permissions for role |
| "incomplete token claims" | JWT missing required claims |

---

## Logging & Monitoring

### Structured Logging

The service uses `slog` with JSON formatting for structured logging:

```go
// internal/logger/logger.go
func New(serviceName string, logLevel string) *slog.Logger {
    opts := &slog.HandlerOptions{
        Level: level,
    }
    handler := slog.NewJSONHandler(os.Stdout, opts)
    logger := slog.New(handler).With(slog.String("service", serviceName))
    slog.SetDefault(logger)
    return logger
}
```

### Log Levels

| Level | Usage |
|-------|-------|
| `debug` | Detailed debugging information, development only |
| `info` | Standard operational messages (default) |
| `warn` | Potential issues that don't affect functionality |
| `error` | Operational errors that require attention |

### Log Examples

**Startup:**
```json
{
  "time": "2026-08-08T10:00:00Z",
  "level": "INFO",
  "msg": "Starting server",
  "service": "auth_service",
  "port": "8080",
  "env": "development"
}
```

**Error:**
```json
{
  "time": "2026-08-08T10:01:00Z",
  "level": "ERROR",
  "msg": "Failed to initialize database pool",
  "service": "auth_service",
  "error": "connection refused"
}
```

### Health Check Endpoint

```
GET /health
```

Used by container orchestrators for readiness and liveness probes.

---

## Deployment

### Docker Configuration

**Multi-stage Dockerfile:**

```dockerfile
# Build stage
FROM golang:1.25-alpine AS builder
WORKDIR /app
RUN apk add --no-cache ca-certificates tzdata
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-w -s" -o auth_service ./cmd/server

# Runtime stage
FROM alpine:latest
WORKDIR /app
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
COPY --from=builder /app/auth_service .
USER appuser
EXPOSE 8080
HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
  CMD wget --spider -q http://localhost:8080/health || exit 1
CMD ["./auth_service"]
```

### Environment Configuration

**`.env` file:**
```env
PORT=8080
DATABASE_URL=postgresql://gratia:gratia123@localhost:5432/gratia?sslmode=disable
JWT_SECRET=your_super_secure_jwt_secret
ACCESS_TOKEN_MINUTES=15
REFRESH_TOKEN_DAYS=30
ENV=development
LOG_LEVEL=info
```

### Docker Compose Example

```yaml
version: '3.8'

services:
  auth_service:
    build: .
    ports:
      - "8080:8080"
    environment:
      - DATABASE_URL=postgresql://gratia:gratia123@postgres:5432/gratia?sslmode=disable
      - JWT_SECRET=${JWT_SECRET}
      - ENV=production
      - LOG_LEVEL=info
    depends_on:
      - postgres
    networks:
      - gratia_network

  postgres:
    image: postgres:16-alpine
    environment:
      - POSTGRES_USER=gratia
      - POSTGRES_PASSWORD=gratia123
      - POSTGRES_DB=gratia
    volumes:
      - postgres_data:/var/lib/postgresql/data
    networks:
      - gratia_network

networks:
  gratia_network:
    driver: bridge

volumes:
  postgres_data:
```

### Graceful Shutdown

The service implements graceful shutdown with a 10-second timeout:

```go
// main.go
quit := make(chan os.Signal, 1)
signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
<-quit

shutdownCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
defer cancel()

if err := srv.Shutdown(shutdownCtx); err != nil {
    log.Error("Server forced to shutdown abruptly", "error", err)
    os.Exit(1)
}
```

---

## Development Guidelines

### Adding New Features

#### 1. New Role
1. Add role to `internal/auth/role.go`
2. Update role validation logic
3. Update role assignment permissions in `service.go`

#### 2. New Authentication Method
1. Add new handler in `handlers.go`
2. Add business logic in `service.go`
3. Add database operations in `repository.go`
4. Register route in `server/http.go`

#### 3. New Protected Endpoint
1. Add handler in `handlers.go`
2. Register route under `protected` subrouter in `server/http.go`

### Testing

**Unit Tests Example:**
```go
// service_test.go
func TestRegisterUser(t *testing.T) {
    mockRepo := &MockRepository{}
    service := NewService(mockRepo, "secret", 15*time.Minute, 30*24*time.Hour)
    
    input := RegisterInput{
        Email: "test@example.com",
        Password: "SecurePassword123",
        Role: RoleUser,
    }
    
    user, err := service.RegisterUser(context.Background(), input, RoleAdmin)
    assert.NoError(t, err)
    assert.Equal(t, "test@example.com", user.Email)
}
```

### Common Patterns

#### Context Usage
Always pass context through all layers:
```go
// Handler
user, err := h.service.RegisterUser(r.Context(), input, callerRole)

// Service
func (s *Service) RegisterUser(ctx context.Context, input RegisterInput, callerRole Role) (*User, error)

// Repository
func (r *repository) CreateUser(ctx context.Context, user *User) error
```

#### Error Wrapping
```go
if err := s.repo.CreateUser(ctx, user); err != nil {
    return nil, fmt.Errorf("failed to create user: %w", err)
}
```

### Performance Considerations

1. **Connection Pooling**: pgxpool manages PostgreSQL connections efficiently
2. **Dependency Injection**: Promotes testability and loose coupling
3. **Zero-Allocation**: Use `[]byte` where possible for performance
4. **Context Timeout**: All database operations respect context deadlines

### Security Best Practices

1. **Never log passwords** or sensitive tokens
2. **Validate all input** before processing
3. **Use constant-time comparisons** for password verification
4. **Rate limit** authentication endpoints
5. **Implement audit logging** for security events
6. **Rotate JWT secrets** periodically
7. **Use HTTPS** in production
8. **Implement CORS** policies correctly

---

## Monitoring & Observability

### Metrics to Track

| Metric | Description |
|--------|-------------|
| `auth_requests_total` | Total authentication requests |
| `auth_failures_total` | Failed authentication attempts |
| `auth_errors_total` | Authentication errors |
| `auth_success_duration` | Authentication duration |
| `active_sessions` | Current active sessions |
| `user_registrations` | New user registrations |

### Log Analysis

**Authentication Failures:**
```json
{
  "time": "...",
  "level": "WARN",
  "msg": "Login failed",
  "service": "auth_service",
  "email": "user@example.com",
  "reason": "invalid credentials"
}
```

**Security Events:**
```json
{
  "time": "...",
  "level": "INFO",
  "msg": "User logged in",
  "service": "auth_service",
  "user_id": "550e8400-...",
  "ip": "192.168.1.1",
  "user_agent": "Mozilla/5.0..."
}
```

---

## Troubleshooting

### Common Issues

#### Database Connection Failed
```
Failed to initialize database pool: failed to connect to database: connection refused
```
**Solution:** Verify PostgreSQL is running and the connection string is correct.

#### JWT Token Validation Failed
```
invalid token claims
```
**Solution:** Ensure JWT_SECRET matches across services.

#### Token Expired
```
invalid or expired token
```
**Solution:** Use the refresh token endpoint to obtain new tokens.

#### Role Assignment Denied
```
unauthorized role assignment
```
**Solution:** Check the caller's role permissions.

---

## Appendix

### A. Dependencies

| Package | Version | Purpose |
|---------|---------|---------|
| `github.com/gorilla/mux` | v1.8.1 | HTTP routing |
| `github.com/jackc/pgx/v5` | v5.5.0 | PostgreSQL driver and connection pool |
| `github.com/golang-jwt/jwt/v5` | v5.0.0 | JWT token generation/validation |
| `golang.org/x/crypto` | latest | bcrypt password hashing |
| `github.com/joho/godotenv` | v1.5.1 | Environment variable loading |
| `github.com/google/uuid` | latest | UUID generation |

### B. Environment Variables Quick Reference

```bash
# Required
DATABASE_URL=postgresql://user:pass@host:5432/db?sslmode=disable
JWT_SECRET=32+_character_secret_key

# Optional
PORT=8080
ACCESS_TOKEN_MINUTES=15
REFRESH_TOKEN_DAYS=30
ENV=development
LOG_LEVEL=info
```

### C. Useful Commands

```bash
# Build
go build -o auth_service ./cmd/server

# Run
./auth_service

# Docker build
docker build -t auth_service .

# Docker run
docker run -p 8080:8080 --env-file .env auth_service

# Test API
curl http://localhost:8080/health

# Test protected endpoint
curl -H "Authorization: Bearer <token>" http://localhost:8080/me
```

### D. Migration Script

```sql
-- Create extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Run migration scripts in order
\i migrations/001_create_users_table.sql
\i migrations/002_create_refresh_tokens_table.sql
\i migrations/003_create_sessions_table.sql
\i migrations/004_create_password_resets_table.sql
```

---

## Future Enhancements

### Proposed Features
1. **Email Verification**: Send verification emails on registration
2. **Two-Factor Authentication**: Support TOTP or SMS
3. **Social Login**: OAuth2 integration (Google, Facebook)
4. **Rate Limiting**: Prevent brute force attacks
5. **Password History**: Prevent password reuse
6. **Account Lockout**: Lock accounts after failed attempts
7. **Activity Logging**: Track all user activities
8. **Admin APIs**: User management endpoints
9. **Webhook Support**: Notify other services of events

### Scaling Considerations
1. **Horizontal Scaling**: Stateless JWT allows multiple instances
2. **Session Storage**: Consider Redis for session data
3. **Caching**: Cache user data to reduce database load
4. **Load Balancing**: Distribute traffic across instances
5. **Database Replication**: Read replicas for scaling queries

---

## Contact & Support

For questions or issues regarding the Auth Service:
- **Service Owner**: Platform Team
- **Repository**: https://github.com/adityawaradkar/gratia/auth_service
- **Documentation**: https://docs.gratia.com/auth_service

---
