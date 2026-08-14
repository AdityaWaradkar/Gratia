# Food Service Documentation

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
12. [Service Integration](#service-integration)

---

## Overview

The Food Service is a production-ready microservice built in Go that manages food donation listings for the Gratia platform. It handles the complete lifecycle of food donations—from creation by donors to claiming by NGOs. This service is designed to work alongside the Auth Service (for authentication) and User Service (for user profile validation), with clear separation of concerns.

### Key Features
- **Food Listing Management** with CRUD operations
- **Donor Validation** via User Service integration
- **Claim Lifecycle Support** for NGOs
- **Automatic Expiry** of time-sensitive listings
- **Internal APIs** for service-to-service communication
- **Graceful Shutdown** and health checking
- **Structured JSON Logging** for observability
- **Docker Ready** with multi-stage builds
- **Background Worker** for automated expiry processing

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
│  │              Router (http.ServeMux)                 │   │
│  │  ┌─────────────────────────────────────────────┐   │   │
│  │  │        Authentication Middleware             │   │   │
│  │  │    (JWT Validation & Context Injection)      │   │   │
│  │  └─────────────────────────────────────────────┘   │   │
│  └─────────────────────────────────────────────────────┘   │
└─────────────────────┬───────────────────────────────────────┘
                      │
┌─────────────────────▼───────────────────────────────────────┐
│              HTTP Handlers (food/handlers.go)               │
│  ┌─────────────────────────────────────────────────────┐   │
│  │  CreateFoodListing | GetFoodListing                │   │
│  │  ListAvailableFoodListings | UpdateFoodListing     │   │
│  │  CancelFoodListing | ValidateFoodClaim            │   │
│  │  MarkFoodClaimed | HealthCheck                    │   │
│  └─────────────────────────────────────────────────────┘   │
└─────────────────────┬───────────────────────────────────────┘
                      │
┌─────────────────────▼───────────────────────────────────────┐
│                Service Layer (food/service.go)              │
│  ┌─────────────────────────────────────────────────────┐   │
│  │  Business Logic & Validation                        │   │
│  │  Donor Verification (via User Client)              │   │
│  │  Listing Lifecycle Management                      │   │
│  │  Claim Validation                                  │   │
│  └─────────────────────────────────────────────────────┘   │
└─────────────────────┬───────────────────────────────────────┘
                      │
┌─────────────────────▼───────────────────────────────────────┐
│              Repository Layer (food/repository.go)          │
│  ┌─────────────────────────────────────────────────────┐   │
│  │  Database Operations (CRUD)                         │   │
│  │  Food Listing Management                            │   │
│  │  Status Updates                                     │   │
│  │  Expiry Sweep                                       │   │
│  └─────────────────────────────────────────────────────┘   │
└─────────────────────┬───────────────────────────────────────┘
                      │
┌─────────────────────▼───────────────────────────────────────┐
│          PostgreSQL Database (pgxpool Connection)           │
│  ┌─────────────────────────────────────────────────────┐   │
│  │  food_listings                                      │   │
│  └─────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────┘
```

### Service Interactions

```
┌──────────────┐         ┌──────────────┐
│   Client     │         │   Client     │
│   (Donor)    │         │   (NGO)      │
└──────┬───────┘         └──────┬───────┘
       │                        │
       ▼                        ▼
┌──────────────────────────────────────────────┐
│            Food Service                      │
│  ┌────────────────────────────────────────┐  │
│  │  Create Listing │ Validate Claim      │  │
│  │  List Listings  │ Mark Claimed        │  │
│  └────────────────────────────────────────┘  │
└──────────┬─────────────────────┬─────────────┘
           │                     │
           ▼                     ▼
┌──────────────────┐   ┌──────────────────┐
│  User Service    │   │  Auth Service    │
│  (Donor/NGO      │   │  (JWT            │
│   Validation)    │   │   Validation)    │
└──────────────────┘   └──────────────────┘
```

### Directory Structure
```
food_service/
├── cmd/
│   └── server/
│       └── main.go               # Application entry point
├── internal/
│   ├── food/
│   │   ├── handlers.go           # HTTP handlers
│   │   ├── service.go            # Business logic
│   │   ├── repository.go         # Data access layer
│   │   ├── models.go             # Domain models
│   │   └── user_client.go        # User service HTTP client
│   ├── config/
│   │   └── config.go             # Configuration management
│   ├── db/
│   │   └── db.go                 # Database connection
│   ├── logger/
│   │   └── logger.go             # Structured logging
│   ├── middleware/
│   │   └── auth_middleware.go    # JWT authentication
│   └── server/
│       └── http.go               # Route registration
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
| `PORT` | string | No | `8082` | HTTP server port |
| `DATABASE_URL` | string | Yes | - | PostgreSQL connection string |
| `JWT_SECRET` | string | Yes | - | Secret key for JWT validation (minimum 32 characters) |
| `USER_SERVICE_URL` | string | Yes | - | URL for user service (donor/NGO validation) |
| `ENV` | string | No | `development` | Environment (development/staging/production) |
| `LOG_LEVEL` | string | No | `info` | Log level (debug/info/warn/error) |

### Configuration Loading Process

```go
// config/config.go
func Load() *Config {
    _ = godotenv.Load()
    
    jwtSecret := mustEnv("JWT_SECRET")
    if len(jwtSecret) < 32 {
        log.Fatal("Fatal: JWT_SECRET must be at least 32 characters")
    }

    return &Config{
        Port:           getEnv("PORT", "8082"),
        DatabaseURL:    mustEnv("DATABASE_URL"),
        JWTSecret:      jwtSecret,
        UserServiceURL: mustEnv("USER_SERVICE_URL"),
        Env:            getEnv("ENV", "development"),
        LogLevel:       getEnv("LOG_LEVEL", "info"),
    }
}
```

### Environment-Specific Configuration

**Development:**
```env
ENV=development
LOG_LEVEL=debug
USER_SERVICE_URL=http://localhost:8081
```

**Production:**
```env
ENV=production
LOG_LEVEL=info
USER_SERVICE_URL=http://user_service:8081
```

---

## Database Schema

### SQL Migration (PostgreSQL)

```sql
-- Food listings table
CREATE TABLE IF NOT EXISTS food_listings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    donor_user_id UUID NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    quantity INTEGER NOT NULL CHECK (quantity > 0),
    unit VARCHAR(50) NOT NULL,
    expiry_time TIMESTAMP WITH TIME ZONE NOT NULL,
    location VARCHAR(255) NOT NULL,
    image_url TEXT,
    status VARCHAR(50) NOT NULL DEFAULT 'AVAILABLE',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    INDEX idx_food_listings_donor_user_id (donor_user_id),
    INDEX idx_food_listings_status (status),
    INDEX idx_food_listings_expiry_time (expiry_time),
    INDEX idx_food_listings_created_at (created_at)
);
```

### Entity Relationship Diagram

```
┌─────────────────────────────┐
│      food_listings          │
│─────────────────────────────│
│ id (PK)                     │
│ donor_user_id (FK) ─────────┼───┐
│ title                       │   │
│ description                 │   │
│ quantity                    │   │
│ unit                        │   │
│ expiry_time                 │   │
│ location                    │   │
│ image_url                   │   │
│ status                      │   │
│ created_at                  │   │
│ updated_at                  │   │
└─────────────────────────────┘   │
                                  │
                                  │
                          ┌───────▼───────────────┐
                          │      users            │
                          │  (in auth_service)    │
                          │───────────────────────│
                          │  id (PK)              │
                          │  email                │
                          │  role                 │
                          └───────────────────────┘
```

---

## Core Components

### 1. Domain Models (`internal/food/models.go`)

#### FoodListing
```go
type FoodListing struct {
    ID          string     `json:"id"`
    DonorUserID string     `json:"donorUserId"`
    Title       string     `json:"title"`
    Description *string    `json:"description,omitempty"`
    Quantity    int        `json:"quantity"`
    Unit        string     `json:"unit"`
    ExpiryTime  time.Time  `json:"expiryTime"`
    Location    string     `json:"location"`
    ImageURL    *string    `json:"imageUrl,omitempty"`
    Status      string     `json:"status"`
    CreatedAt   time.Time  `json:"createdAt"`
    UpdatedAt   time.Time  `json:"updatedAt"`
}
```

#### Food Status Constants
```go
const (
    FoodStatusAvailable = "AVAILABLE"
    FoodStatusClaimed   = "CLAIMED"
    FoodStatusExpired   = "EXPIRED"
    FoodStatusCancelled = "CANCELLED"
)
```

### 2. Service Layer (`internal/food/service.go`)

The service layer implements all business logic and cross-service validation:

#### CreateFoodListing
```go
func (s *Service) CreateFoodListing(
    ctx context.Context,
    userID string,
    title string,
    description *string,
    quantity int,
    unit string,
    expiryTime time.Time,
    location string,
    imageURL *string,
) (*FoodListing, error)
```
- Validates donor status via User Service
- Validates input fields (title, quantity, unit, location, expiry)
- Creates listing with `AVAILABLE` status

#### GetFoodListing
```go
func (s *Service) GetFoodListing(ctx context.Context, id string) (*FoodListing, error)
```
- Retrieves a single listing by ID
- Returns `ErrListingNotFound` if not found

#### ListAvailableFoodListings
```go
func (s *Service) ListAvailableFoodListings(ctx context.Context) ([]*FoodListing, error)
```
- Lists all listings with `AVAILABLE` status
- Ordered by creation date (newest first)

#### UpdateFoodListing
```go
func (s *Service) UpdateFoodListing(ctx context.Context, userID string, listing *FoodListing) error
```
- Validates donor owns the listing
- Ensures listing is still `AVAILABLE`
- Re-validates input fields
- Updates the listing

#### CancelFoodListing
```go
func (s *Service) CancelFoodListing(ctx context.Context, userID string, listingID string) error
```
- Validates donor owns the listing
- Ensures listing is still `AVAILABLE`
- Updates status to `CANCELLED`

#### IsClaimable
```go
func (s *Service) IsClaimable(ctx context.Context, id string) (bool, string, error)
```
- Checks if a listing can be claimed
- Returns reason if not claimable (expired, already claimed, etc.)

#### MarkClaimed
```go
func (s *Service) MarkClaimed(ctx context.Context, id string) error
```
- Updates listing status to `CLAIMED`
- Used by internal claim service

#### ExpireListings
```go
func (s *Service) ExpireListings(ctx context.Context) error
```
- Sweeps and expires outdated listings
- Called by background worker

### 3. Repository Layer (`internal/food/repository.go`)

```go
type Repository interface {
    CreateFoodListing(ctx context.Context, listing *FoodListing) error
    GetFoodListingByID(ctx context.Context, id string) (*FoodListing, error)
    ListFoodListings(ctx context.Context, status string) ([]*FoodListing, error)
    UpdateFoodListing(ctx context.Context, listing *FoodListing) error
    UpdateFoodStatus(ctx context.Context, id string, status string) error
    ExpireFoodListings(ctx context.Context) error
}
```

### 4. User Client (`internal/food/user_client.go`)

Internal HTTP client for communication with User Service:

```go
type UserClient interface {
    IsDonor(ctx context.Context, userID string) (bool, error)
    IsVerifiedNGO(ctx context.Context, userID string) (bool, error)
}
```

**Implementation:**
- 5-second timeout to prevent cascading failures
- Proper error handling for service unavailability
- JSON response parsing for NGO verification

### 5. HTTP Handlers (`internal/food/handlers.go`)

Manages HTTP request/response:

```go
type Handler struct {
    service *Service
}

// Handler methods
func (h *Handler) CreateFoodListing(w http.ResponseWriter, r *http.Request)
func (h *Handler) GetFoodListing(w http.ResponseWriter, r *http.Request)
func (h *Handler) ListAvailableFoodListings(w http.ResponseWriter, r *http.Request)
func (h *Handler) UpdateFoodListing(w http.ResponseWriter, r *http.Request)
func (h *Handler) CancelFoodListing(w http.ResponseWriter, r *http.Request)
func (h *Handler) ValidateFoodClaim(w http.ResponseWriter, r *http.Request)
func (h *Handler) MarkFoodClaimed(w http.ResponseWriter, r *http.Request)
func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request)
```

### 6. Background Worker (`cmd/server/main.go`)

```go
func startExpiryWorker(ctx context.Context, service *food.Service, log *slog.Logger) {
    ticker := time.NewTicker(1 * time.Minute)
    defer ticker.Stop()

    for {
        select {
        case <-ctx.Done():
            log.Info("expiry worker shutting down")
            return
        case <-ticker.C:
            expireCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
            if err := service.ExpireListings(expireCtx); err != nil {
                log.Error("expiry worker error", slog.String("error", err.Error()))
            }
            cancel()
        }
    }
}
```

---

## API Endpoints

### Base URL: `http://localhost:8082`

### Public Endpoints (No Authentication Required)

#### 1. Health Check
```http
GET /health
```
**Response:**
```json
{
  "status": "healthy",
  "service": "food_service"
}
```

### Protected Endpoints (Requires Authentication)

#### 2. Create Food Listing
```http
POST /foods
```
**Headers:**
```
Authorization: Bearer <accessToken>
```
**Request Body:**
```json
{
  "title": "Vegetable Biryani",
  "description": "Freshly prepared lunch for 25 people",
  "quantity": 25,
  "unit": "plates",
  "expiryTime": "2026-08-09T12:00:00Z",
  "location": "Navi Mumbai",
  "imageUrl": "https://example.com/food.jpg"
}
```
**Response (201 Created):**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "donorUserId": "550e8400-e29b-41d4-a716-446655440000",
  "title": "Vegetable Biryani",
  "description": "Freshly prepared lunch for 25 people",
  "quantity": 25,
  "unit": "plates",
  "expiryTime": "2026-08-09T12:00:00Z",
  "location": "Navi Mumbai",
  "imageUrl": "https://example.com/food.jpg",
  "status": "AVAILABLE",
  "createdAt": "2026-08-08T10:00:00Z",
  "updatedAt": "2026-08-08T10:00:00Z"
}
```

#### 3. List Available Food Listings
```http
GET /foods
```
**Headers:**
```
Authorization: Bearer <accessToken>
```
**Response (200 OK):**
```json
[
  {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "donorUserId": "550e8400-e29b-41d4-a716-446655440000",
    "title": "Vegetable Biryani",
    "quantity": 25,
    "unit": "plates",
    "expiryTime": "2026-08-09T12:00:00Z",
    "location": "Navi Mumbai",
    "status": "AVAILABLE",
    "createdAt": "2026-08-08T10:00:00Z",
    "updatedAt": "2026-08-08T10:00:00Z"
  }
]
```

#### 4. Get Food Listing by ID
```http
GET /foods/{id}
```
**Headers:**
```
Authorization: Bearer <accessToken>
```
**Response (200 OK):** FoodListing object

#### 5. Update Food Listing
```http
PUT /foods/{id}
```
**Headers:**
```
Authorization: Bearer <accessToken>
```
**Request Body:** Same as Create request
**Response:** 204 No Content

#### 6. Cancel Food Listing
```http
DELETE /foods/{id}
```
**Headers:**
```
Authorization: Bearer <accessToken>
```
**Response:** 204 No Content

### Internal Endpoints (Service-to-Service, No Auth)

#### 7. Validate Food Claim
```http
GET /internal/foods/{id}/validate
```
**Response (200 OK):**
```json
{
  "foodId": "550e8400-e29b-41d4-a716-446655440000",
  "claimable": true,
  "reason": ""
}
```

#### 8. Mark Food as Claimed
```http
PATCH /internal/foods/{id}/claim
```
**Response:** 204 No Content

---

## Security Implementation

### 1. Authentication Middleware

Same pattern as auth_service and user_service:

```go
// internal/middleware/auth_middleware.go
func Auth(jwtSecret string) func(http.Handler) http.Handler {
    // JWT validation and context injection
}
```

**Context Values:**
- `user_id`: User's UUID
- `user_role`: User's role

### 2. Donor Validation

Before creating a listing, the service validates donor status:

```go
isDonor, err := s.userClient.IsDonor(ctx, userID)
if err != nil || !isDonor {
    return nil, ErrForbidden
}
```

### 3. Ownership Validation

Update and cancel operations verify the user owns the listing:

```go
if existing.DonorUserID != userID {
    return ErrForbidden
}
```

### 4. Status Validation

Only `AVAILABLE` listings can be modified or claimed:

```go
if existing.Status != FoodStatusAvailable {
    return ErrStatusNotAvailable
}
```

### 5. Input Validation

| Field | Validation |
|-------|-----------|
| `title` | Required, non-empty |
| `quantity` | Must be > 0 |
| `unit` | Required, non-empty |
| `location` | Required, non-empty |
| `expiryTime` | Must be in the future |

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
| 401 | Unauthorized |
| 403 | Forbidden |
| 404 | Not Found |
| 409 | Conflict |

### Common Error Messages

| Error | Scenario |
|-------|----------|
| "invalid request body" | Malformed JSON |
| "invalid input: title is required" | Missing required field |
| "unauthorized" | Invalid or missing token |
| "forbidden" | Not a donor or not owner |
| "food listing not found" | Listing ID doesn't exist |
| "only available listings can be modified" | Attempting to modify non-available listing |
| "expiry time must be in the future" | Past expiry date |

---

## Logging & Monitoring

### Structured Logging

Same pattern as other services:

```go
// internal/logger/logger.go
func New(serviceName string, logLevel string) *slog.Logger {
    // JSON logging with service name injection
}
```

### Log Examples

**Listing Creation:**
```json
{
  "time": "2026-08-08T10:00:00Z",
  "level": "INFO",
  "msg": "Food listing created",
  "service": "food_service",
  "listing_id": "550e8400-...",
  "donor_user_id": "550e8400-..."
}
```

**Expiry Worker:**
```json
{
  "time": "2026-08-08T10:01:00Z",
  "level": "INFO",
  "msg": "Expired listings processed",
  "service": "food_service",
  "count": 5
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
    go build -ldflags="-w -s" -o food_service ./cmd/server

# Runtime stage
FROM alpine:latest
WORKDIR /app
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
COPY --from=builder /app/food_service .
USER appuser
EXPOSE 8082
HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
  CMD wget --spider -q http://localhost:8082/health || exit 1
CMD ["./food_service"]
```

### Docker Compose Example

```yaml
food_service:
  build: .
  ports:
    - "8082:8082"
  environment:
    - DATABASE_URL=postgresql://gratia:gratia123@postgres:5432/gratia?sslmode=disable
    - JWT_SECRET=${JWT_SECRET}
    - USER_SERVICE_URL=http://user_service:8081
    - ENV=production
    - LOG_LEVEL=info
  depends_on:
    - postgres
    - user_service
  networks:
    - gratia_network
```

### Graceful Shutdown

The service implements graceful shutdown with:
- 10-second timeout for active requests
- Background worker cancellation
- Proper signal handling

```go
// main.go
stop := make(chan os.Signal, 1)
signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
<-stop

cancelWorker()
srv.Shutdown(shutdownCtx)
```

---

## Service Integration

### Internal APIs

The Food Service exposes internal endpoints for the Claim Service:

| Endpoint | Method | Purpose |
|----------|--------|---------|
| `/internal/foods/{id}/validate` | GET | Validate if a listing can be claimed |
| `/internal/foods/{id}/claim` | PATCH | Mark a listing as claimed |

### External Dependencies

| Service | Purpose | URL Variable |
|---------|---------|--------------|
| Auth Service | JWT Validation | `JWT_SECRET` (shared) |
| User Service | Donor/NGO Validation | `USER_SERVICE_URL` |

### User Client Interface

```go
type UserClient interface {
    IsDonor(ctx context.Context, userID string) (bool, error)
    IsVerifiedNGO(ctx context.Context, userID string) (bool, error)
}
```

---

## Development Guidelines

### Adding New Features

#### 1. New Food Status
1. Add status constant in `models.go`
2. Update status validation logic
3. Update API documentation

#### 2. New Internal Endpoint
1. Add handler method in `handlers.go`
2. Add service method in `service.go`
3. Add repository method in `repository.go`
4. Register route in `server/http.go`

### Testing

**Unit Test Example:**
```go
func TestCreateFoodListing(t *testing.T) {
    mockRepo := &MockRepository{}
    mockUserClient := &MockUserClient{}
    service := NewService(mockRepo, mockUserClient)
    
    listing, err := service.CreateFoodListing(
        context.Background(),
        "user-123",
        "Test Food",
        nil,
        10,
        "plates",
        time.Now().Add(24*time.Hour),
        "Mumbai",
        nil,
    )
    
    assert.NoError(t, err)
    assert.Equal(t, "AVAILABLE", listing.Status)
}
```

### Performance Considerations

1. **Cross-Service Calls**: 5-second timeout prevents cascading failures
2. **Background Worker**: Dedicated goroutine for expiry processing
3. **Connection Pooling**: pgxpool for efficient database connections
4. **Idempotency**: Internal endpoints are idempotent

---

## Background Worker

The expiry worker runs every minute to clean up expired listings:

```go
func startExpiryWorker(ctx context.Context, service *food.Service, log *slog.Logger) {
    ticker := time.NewTicker(1 * time.Minute)
    defer ticker.Stop()

    for {
        select {
        case <-ctx.Done():
            log.Info("expiry worker shutting down")
            return
        case <-ticker.C:
            expireCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
            if err := service.ExpireListings(expireCtx); err != nil {
                log.Error("expiry worker error", slog.String("error", err.Error()))
            }
            cancel()
        }
    }
}
```

---

## Future Enhancements

1. **Image Upload**: S3/Cloud storage for food images
2. **Search & Filter**: Advanced listing search with filters
3. **Geolocation**: Location-based listing discovery
4. **Notifications**: Alert NGOs of new listings
5. **Caching**: Redis for listing data
6. **Webhooks**: Notify claim service on status changes
7. **Metrics**: Prometheus metrics for monitoring
8. **Rate Limiting**: Prevent abuse of listing creation

---

## Dependencies

| Package                        | Version | Purpose               |
| ------------------------------ | ------- | --------------------- |
| `github.com/jackc/pgx/v5`      | v5.5.0  | PostgreSQL driver     |
| `github.com/golang-jwt/jwt/v5` | v5.0.0  | JWT token validation  |
| `github.com/joho/godotenv`     | v1.5.1  | Environment variables |
| `github.com/google/uuid`       | latest  | UUID generation       |