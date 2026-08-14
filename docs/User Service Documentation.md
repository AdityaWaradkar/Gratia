# User Service Documentation

---

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

The User Service is a production-ready microservice built in Go that manages user profiles, preferences, and account management for the platform. It handles user data operations, profile updates, and provides user-related information to other services. This service is designed to work alongside the Auth Service, with clear separation of concerns—authentication handled by Auth Service, user data management handled by User Service.

### Key Features
- **User Profile Management** with CRUD operations
- **User Preferences** storage and retrieval
- **User Search** with filtering and pagination
- **Role-Based Access Control** for user data
- **Graceful Shutdown** and health checking
- **Structured JSON Logging** for observability
- **Docker Ready** with multi-stage builds
- **Service-to-Service Authentication** via JWT

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
│              HTTP Handlers (user/handlers.go)               │
│  ┌─────────────────────────────────────────────────────┐   │
│  │  GetProfile | UpdateProfile | DeleteProfile        │   │
│  │  GetPreferences | UpdatePreferences                 │   │
│  │  SearchUsers | GetUserByID | HealthCheck           │   │
│  └─────────────────────────────────────────────────────┘   │
└─────────────────────┬───────────────────────────────────────┘
                      │
┌─────────────────────▼───────────────────────────────────────┐
│                Service Layer (user/service.go)              │
│  ┌─────────────────────────────────────────────────────┐   │
│  │  Business Logic & Validation                        │   │
│  │  Profile Management | Preferences                  │   │
│  │  Search & Filtering | Authorization Checks         │   │
│  └─────────────────────────────────────────────────────┘   │
└─────────────────────┬───────────────────────────────────────┘
                      │
┌─────────────────────▼───────────────────────────────────────┐
│              Repository Layer (user/repository.go)          │
│  ┌─────────────────────────────────────────────────────┐   │
│  │  Database Operations (CRUD)                         │   │
│  │  User Profiles | Preferences | Search              │   │
│  └─────────────────────────────────────────────────────┘   │
└─────────────────────┬───────────────────────────────────────┘
                      │
┌─────────────────────▼───────────────────────────────────────┐
│          PostgreSQL Database (pgxpool Connection)           │
│  ┌─────────────────────────────────────────────────────┐   │
│  │  users (shared with auth) | user_profiles          │   │
│  │  user_preferences | user_activity                  │   │
│  └─────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────┘
```

### Directory Structure
```
user_service/
├── cmd/
│   └── server/
│       └── main.go               # Application entry point
├── internal/
│   ├── user/
│   │   ├── handlers.go           # HTTP handlers
│   │   ├── service.go            # Business logic
│   │   ├── repository.go         # Data access layer
│   │   ├── models.go             # Domain models
│   │   └── validator.go          # Input validation
│   ├── config/
│   │   └── config.go             # Configuration management
│   ├── db/
│   │   └── db.go                 # Database connection
│   ├── logger/
│   │   └── logger.go             # Structured logging
│   ├── middleware/
│   │   ├── auth_middleware.go    # JWT authentication
│   │   └── admin_middleware.go   # Admin authorization
│   ├── server/
│   │   └── http.go               # Route registration
│   └── utils/
│       ├── search.go             # Search helpers
│       └── pagination.go         # Pagination helpers
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
| `PORT` | string | No | `8081` | HTTP server port |
| `DATABASE_URL` | string | Yes | - | PostgreSQL connection string |
| `JWT_SECRET` | string | Yes | - | Secret key for JWT validation (same as auth_service) |
| `ENV` | string | No | `development` | Environment (development/staging/production) |
| `LOG_LEVEL` | string | No | `info` | Log level (debug/info/warn/error) |
| `AUTH_SERVICE_URL` | string | Yes | - | URL for auth service (for user validation) |

### Configuration Loading Process

```go
// config/config.go
func Load() *Config {
    _ = godotenv.Load()
    
    return &Config{
        Port:             getEnv("PORT", "8081"),
        DatabaseURL:      mustEnv("DATABASE_URL"),
        JWTSecret:        mustEnv("JWT_SECRET"),
        Env:              getEnv("ENV", "development"),
        LogLevel:         getEnv("LOG_LEVEL", "info"),
        AuthServiceURL:   mustEnv("AUTH_SERVICE_URL"),
    }
}
```

---

## Database Schema

### SQL Migration (PostgreSQL)

```sql
-- Users table (shared with auth_service)
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL DEFAULT 'USER',
    email_verified BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- User profiles table
CREATE TABLE user_profiles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    display_name VARCHAR(100),
    avatar_url TEXT,
    bio TEXT,
    location VARCHAR(255),
    website VARCHAR(255),
    phone VARCHAR(20),
    date_of_birth DATE,
    gender VARCHAR(20),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(user_id),
    INDEX idx_user_profiles_user_id (user_id),
    INDEX idx_user_profiles_display_name (display_name)
);

-- User preferences table
CREATE TABLE user_preferences (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    theme VARCHAR(50) DEFAULT 'light',
    language VARCHAR(10) DEFAULT 'en',
    timezone VARCHAR(50) DEFAULT 'UTC',
    email_notifications BOOLEAN DEFAULT TRUE,
    push_notifications BOOLEAN DEFAULT TRUE,
    marketing_emails BOOLEAN DEFAULT FALSE,
    privacy_settings JSONB DEFAULT '{}',
    notification_settings JSONB DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(user_id),
    INDEX idx_user_preferences_user_id (user_id)
);

-- User activity log (for audit)
CREATE TABLE user_activity (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    action VARCHAR(100) NOT NULL,
    metadata JSONB,
    ip_address VARCHAR(45),
    user_agent TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    INDEX idx_user_activity_user_id (user_id),
    INDEX idx_user_activity_created_at (created_at)
);
```

### Entity Relationship Diagram

```
┌─────────────────┐
│     users       │
│─────────────────│
│ id (PK)         │───┐
│ email           │   │
│ role            │   │
│ created_at      │   │
│ updated_at      │   │
└─────────────────┘   │
        │             │
        │             │
        │             │
┌───────▼─────────────┼──────┐
│  user_profiles      │      │
│─────────────────────│      │
│ id (PK)            │      │
│ user_id (FK)───────┘      │
│ first_name                │
│ last_name                 │
│ display_name              │
│ avatar_url                │
│ bio                       │
│ location                  │
│ created_at                │
│ updated_at                │
└───────────────────────────┘

┌─────────────────┐
│user_preferences │
│─────────────────│
│ id (PK)         │
│ user_id (FK)────┘
│ theme           │
│ language        │
│ timezone        │
│ email_notif     │
│ push_notif      │
│ privacy_setting │
│ created_at      │
│ updated_at      │
└─────────────────┘

┌─────────────────┐
│ user_activity   │
│─────────────────│
│ id (PK)         │
│ user_id (FK)────┘
│ action          │
│ metadata        │
│ ip_address      │
│ user_agent      │
│ created_at      │
└─────────────────┘
```

---

## Core Components

### 1. Domain Models (`internal/user/models.go`)

#### UserProfile
```go
type UserProfile struct {
    ID          string     `json:"id"`
    UserID      string     `json:"userId"`
    FirstName   string     `json:"firstName"`
    LastName    string     `json:"lastName"`
    DisplayName string     `json:"displayName"`
    AvatarURL   string     `json:"avatarUrl"`
    Bio         string     `json:"bio"`
    Location    string     `json:"location"`
    Website     string     `json:"website"`
    Phone       string     `json:"phone"`
    DateOfBirth *time.Time `json:"dateOfBirth,omitempty"`
    Gender      string     `json:"gender"`
    CreatedAt   time.Time  `json:"createdAt"`
    UpdatedAt   time.Time  `json:"updatedAt"`
}
```

#### UserPreferences
```go
type UserPreferences struct {
    ID                  string          `json:"id"`
    UserID              string          `json:"userId"`
    Theme               string          `json:"theme"`
    Language            string          `json:"language"`
    Timezone            string          `json:"timezone"`
    EmailNotifications  bool            `json:"emailNotifications"`
    PushNotifications   bool            `json:"pushNotifications"`
    MarketingEmails     bool            `json:"marketingEmails"`
    PrivacySettings     json.RawMessage `json:"privacySettings"`
    NotificationSettings json.RawMessage `json:"notificationSettings"`
    CreatedAt           time.Time       `json:"createdAt"`
    UpdatedAt           time.Time       `json:"updatedAt"`
}
```

#### UserActivity
```go
type UserActivity struct {
    ID        string          `json:"id"`
    UserID    string          `json:"userId"`
    Action    string          `json:"action"`
    Metadata  json.RawMessage `json:"metadata"`
    IPAddress string          `json:"ipAddress"`
    UserAgent string          `json:"userAgent"`
    CreatedAt time.Time       `json:"createdAt"`
}
```

### 2. Service Layer (`internal/user/service.go`)

The service layer implements all business logic:

#### GetProfile
```go
func (s *Service) GetProfile(ctx context.Context, userID string) (*UserProfile, error)
```
- Retrieves user profile from database
- Returns error if profile not found

#### UpdateProfile
```go
func (s *Service) UpdateProfile(ctx context.Context, userID string, input UpdateProfileInput) (*UserProfile, error)
```
- Validates input fields
- Updates only provided fields (partial update)
- Returns updated profile

#### SearchUsers
```go
func (s *Service) SearchUsers(ctx context.Context, query string, filters SearchFilters, pagination Pagination) (*SearchResult, error)
```
- Searches users by name, email, or other fields
- Applies filters (role, location, etc.)
- Returns paginated results

#### UpdatePreferences
```go
func (s *Service) UpdatePreferences(ctx context.Context, userID string, input UpdatePreferencesInput) (*UserPreferences, error)
```
- Updates user preferences
- Merges with existing preferences
- Returns updated preferences

### 3. Repository Layer (`internal/user/repository.go`)

```go
type Repository interface {
    // Profile operations
    GetProfile(ctx context.Context, userID string) (*UserProfile, error)
    CreateProfile(ctx context.Context, profile *UserProfile) error
    UpdateProfile(ctx context.Context, profile *UserProfile) error
    DeleteProfile(ctx context.Context, userID string) error
    
    // Preferences operations
    GetPreferences(ctx context.Context, userID string) (*UserPreferences, error)
    CreatePreferences(ctx context.Context, prefs *UserPreferences) error
    UpdatePreferences(ctx context.Context, prefs *UserPreferences) error
    
    // User operations (read-only from shared table)
    GetUserByID(ctx context.Context, userID string) (*User, error)
    GetUserByEmail(ctx context.Context, email string) (*User, error)
    SearchUsers(ctx context.Context, query string, filters SearchFilters, pagination Pagination) ([]*UserProfile, int64, error)
    
    // Activity logging
    LogActivity(ctx context.Context, activity *UserActivity) error
}
```

---

## API Endpoints

### Base URL: `http://localhost:8081`

### Public Endpoints (No Authentication Required)

#### 1. Health Check
```http
GET /health
```
**Response:**
```json
{
  "status": "healthy",
  "service": "user_service"
}
```

### Protected Endpoints (Requires Authentication)

#### 2. Get Current User Profile
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
  "userId": "550e8400-e29b-41d4-a716-446655440000",
  "firstName": "John",
  "lastName": "Doe",
  "displayName": "JohnDoe",
  "avatarUrl": "https://example.com/avatar.jpg",
  "bio": "Software developer",
  "location": "New York, USA",
  "website": "https://johndoe.dev",
  "phone": "+1234567890",
  "dateOfBirth": "1990-01-01T00:00:00Z",
  "gender": "male",
  "createdAt": "2026-08-08T10:00:00Z",
  "updatedAt": "2026-08-08T10:00:00Z"
}
```

#### 3. Update Current User Profile
```http
PUT /me
```
**Headers:**
```
Authorization: Bearer <accessToken>
```
**Request Body:**
```json
{
  "firstName": "John",
  "lastName": "Doe",
  "displayName": "JohnDoe",
  "bio": "Senior software developer",
  "location": "San Francisco, USA",
  "website": "https://johndoe.dev",
  "phone": "+1234567890"
}
```
**Response (200 OK):** Updated profile object

#### 4. Get User Preferences
```http
GET /me/preferences
```
**Headers:**
```
Authorization: Bearer <accessToken>
```
**Response (200 OK):**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "userId": "550e8400-e29b-41d4-a716-446655440000",
  "theme": "dark",
  "language": "en",
  "timezone": "America/New_York",
  "emailNotifications": true,
  "pushNotifications": true,
  "marketingEmails": false,
  "privacySettings": {"profileVisibility": "public"},
  "notificationSettings": {"emailDigest": "daily"},
  "createdAt": "2026-08-08T10:00:00Z",
  "updatedAt": "2026-08-08T10:00:00Z"
}
```

#### 5. Update User Preferences
```http
PUT /me/preferences
```
**Headers:**
```
Authorization: Bearer <accessToken>
```
**Request Body:**
```json
{
  "theme": "dark",
  "language": "en",
  "timezone": "America/New_York",
  "emailNotifications": true,
  "pushNotifications": false,
  "marketingEmails": false
}
```
**Response (200 OK):** Updated preferences object

#### 6. Get User Profile by ID (Admin Only)
```http
GET /users/{id}
```
**Headers:**
```
Authorization: Bearer <accessToken>
```
**Response (200 OK):** User profile object

#### 7. Search Users (Admin Only)
```http
GET /users?query=john&role=USER&page=1&limit=20
```
**Headers:**
```
Authorization: Bearer <accessToken>
```
**Response (200 OK):**
```json
{
  "users": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "email": "john@example.com",
      "firstName": "John",
      "lastName": "Doe",
      "displayName": "JohnDoe",
      "role": "USER"
    }
  ],
  "total": 1,
  "page": 1,
  "limit": 20
}
```

#### 8. Delete User Account
```http
DELETE /me
```
**Headers:**
```
Authorization: Bearer <accessToken>
```
**Response:** 204 No Content

---

## Security Implementation

### 1. Authentication Middleware

```go
// internal/middleware/auth_middleware.go
func Auth(jwtSecret string) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            // Extract and validate JWT token
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

### 2. Admin Authorization Middleware

```go
// internal/middleware/admin_middleware.go
func AdminOnly(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        role := middleware.UserRole(r.Context())
        if role != "ADMIN" {
            writeError(w, http.StatusForbidden, "admin access required")
            return
        }
        next(w, r)
    }
}
```

### 3. Service-to-Service Authentication

For internal service communication:
```go
// Requests from other services include a service token
func ValidateServiceToken(ctx context.Context, token string) bool {
    // Validate using shared JWT secret
    // Check for service-specific claims
    return true
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
| 200 | Success |
| 201 | Created |
| 204 | No Content |
| 400 | Bad Request |
| 401 | Unauthorized (Invalid token) |
| 403 | Forbidden (Insufficient permissions) |
| 404 | Not Found |
| 409 | Conflict |
| 500 | Internal Server Error |

### Common Error Messages

| Error | Scenario |
|-------|----------|
| "invalid json payload" | Malformed JSON in request body |
| "user profile not found" | User doesn't have a profile |
| "user not found" | User ID doesn't exist |
| "unauthorized access" | Attempting to access another user's data |
| "admin access required" | Non-admin accessing admin endpoint |
| "invalid search query" | Malformed search parameters |
| "preferences not found" | User doesn't have preferences |

---

## Logging & Monitoring

### Structured Logging

Same pattern as auth_service:

```go
// internal/logger/logger.go
func New(serviceName string, logLevel string) *slog.Logger {
    // Same implementation as auth_service
}
```

### Log Examples

**Profile Update:**
```json
{
  "time": "2026-08-08T10:00:00Z",
  "level": "INFO",
  "msg": "Profile updated",
  "service": "user_service",
  "user_id": "550e8400-...",
  "fields": ["firstName", "lastName"]
}
```

---

## Deployment

### Docker Configuration

```dockerfile
# Build stage
FROM golang:1.25-alpine AS builder
WORKDIR /app
RUN apk add --no-cache ca-certificates tzdata
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-w -s" -o user_service ./cmd/server

# Runtime stage
FROM alpine:latest
WORKDIR /app
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
COPY --from=builder /app/user_service .
USER appuser
EXPOSE 8081
HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
  CMD wget --spider -q http://localhost:8081/health || exit 1
CMD ["./user_service"]
```

### Docker Compose Example

```yaml
user_service:
  build: .
  ports:
    - "8081:8081"
  environment:
    - DATABASE_URL=postgresql://gratia:gratia123@postgres:5432/gratia?sslmode=disable
    - JWT_SECRET=${JWT_SECRET}
    - AUTH_SERVICE_URL=http://auth_service:8080
    - ENV=production
    - LOG_LEVEL=info
  depends_on:
    - postgres
  networks:
    - gratia_network
```

---

## Development Guidelines

### Adding New Features

1. **New Profile Field**
   - Add field to `UserProfile` model
   - Update database schema
   - Update service validation
   - Update API response

2. **New User Activity Type**
   - Add activity type constant
   - Update activity logging
   - Add to search/filter

### Testing

```go
func TestUpdateProfile(t *testing.T) {
    mockRepo := &MockRepository{}
    service := NewService(mockRepo)
    
    input := UpdateProfileInput{
        FirstName: "John",
        LastName: "Doe",
    }
    
    profile, err := service.UpdateProfile(context.Background(), "user-123", input)
    assert.NoError(t, err)
    assert.Equal(t, "John", profile.FirstName)
}
```

---

## Monitoring & Observability

### Metrics to Track

| Metric | Description |
|--------|-------------|
| `profile_updates_total` | Profile update operations |
| `profile_retrievals_total` | Profile retrieval operations |
| `user_searches_total` | User search operations |
| `profile_update_duration` | Profile update latency |

---

## Future Enhancements

1. **User Avatar Upload**: S3/Cloud storage integration
2. **User Following/Followers**: Social features
3. **User Recommendations**: Based on activity and preferences
4. **User Reports**: Admin dashboard data
5. **Webhook Notifications**: On profile updates
6. **User Import/Export**: Bulk operations
7. **User Verification**: Document upload for NGO/Donor verification

---

## Dependencies

| Package | Version | Purpose |
|---------|---------|---------|
| `github.com/gorilla/mux` | v1.8.1 | HTTP routing |
| `github.com/jackc/pgx/v5` | v5.5.0 | PostgreSQL driver |
| `github.com/golang-jwt/jwt/v5` | v5.0.0 | JWT token validation |
| `github.com/joho/godotenv` | v1.5.1 | Environment variables |
| `github.com/google/uuid` | latest | UUID generation |