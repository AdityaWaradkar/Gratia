# Claim Service Documentation

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

The Claim Service is a production-ready microservice built in Go that manages the complete lifecycle of food claims on the Gratia platform. It handles the entire claim process—from creation by NGOs to delivery confirmation—with strict state machine validation and cross-service coordination. This service is designed to work alongside the Auth Service (for authentication), User Service (for NGO verification), and Food Service (for food listing validation), ensuring a seamless claim experience.

### Key Features
- **Claim Lifecycle Management** with strict state transitions
- **NGO Verification** via User Service integration
- **Food Listing Validation** via Food Service integration
- **State Machine** with validation for each transition
- **Actor-Based Authorization** (NGO and Donor roles)
- **Idempotent Operations** preventing duplicate claims
- **Timestamp Tracking** for every state change
- **Graceful Shutdown** and health checking
- **Structured JSON Logging** for observability
- **Docker Ready** with multi-stage builds

### Claim Status Flow
```
                    ┌─────────────┐
                    │   CREATED   │
                    │  (NGO claims)│
                    └──────┬──────┘
                           │
              ┌────────────┼────────────┐
              │            │            │
              ▼            ▼            ▼
       ┌──────────┐ ┌──────────┐ ┌──────────┐
       │ ACCEPTED │ │ REJECTED │ │CANCELLED │
       │(Donor)   │ │(Donor)   │ │(NGO)     │
       └────┬─────┘ └──────────┘ └──────────┘
            │
            ▼
       ┌──────────┐
       │ PICKED_UP│
       │  (NGO)   │
       └────┬─────┘
            │
            ▼
       ┌──────────┐
       │DELIVERED │
       │  (NGO)   │
       └──────────┘
```

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
│              HTTP Handlers (claim/handlers.go)              │
│  ┌─────────────────────────────────────────────────────┐   │
│  │  CreateClaim | ApproveClaim | RejectClaim          │   │
│  │  CancelClaim | MarkPickedUp | MarkDelivered       │   │
│  │  GetClaim | GetClaimsByNGO | GetClaimsByDonor     │   │
│  │  GetClaimsByFood | HealthCheck                    │   │
│  └─────────────────────────────────────────────────────┘   │
└─────────────────────┬───────────────────────────────────────┘
                      │
┌─────────────────────▼───────────────────────────────────────┐
│                Service Layer (claim/service.go)             │
│  ┌─────────────────────────────────────────────────────┐   │
│  │  Business Logic & State Machine                     │   │
│  │  NGO Verification (via User Client)                │   │
│  │  Food Validation (via Food Client)                 │   │
│  │  Claim Lifecycle Management                        │   │
│  │  Authorization Checks                              │   │
│  └─────────────────────────────────────────────────────┘   │
└─────────────────────┬───────────────────────────────────────┘
                      │
┌─────────────────────▼───────────────────────────────────────┐
│              Repository Layer (claim/repository.go)         │
│  ┌─────────────────────────────────────────────────────┐   │
│  │  Database Operations (CRUD)                         │   │
│  │  Claim Management                                   │   │
│  │  Status Updates with Timestamp Tracking             │   │
│  │  Active Claim Checks                                │   │
│  └─────────────────────────────────────────────────────┘   │
└─────────────────────┬───────────────────────────────────────┘
                      │
┌─────────────────────▼───────────────────────────────────────┐
│          PostgreSQL Database (pgxpool Connection)           │
│  ┌─────────────────────────────────────────────────────┐   │
│  │  claims                                             │   │
│  └─────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────┘
```

### Service Interactions

```
┌──────────────┐         ┌──────────────┐
│   Client     │         │   Client     │
│   (NGO)      │         │   (Donor)    │
└──────┬───────┘         └──────┬───────┘
       │                        │
       ▼                        ▼
┌──────────────────────────────────────────────┐
│            Claim Service                     │
│  ┌────────────────────────────────────────┐  │
│  │  Create Claim │ Approve Claim         │  │
│  │  Reject Claim │ Cancel Claim          │  │
│  │  Mark Pickup  │ Mark Delivery         │  │
│  └────────────────────────────────────────┘  │
└──────────┬─────────────────────┬─────────────┘
           │                     │
           ▼                     ▼
┌──────────────────┐   ┌──────────────────┐
│  User Service    │   │  Food Service    │
│  (NGO            │   │  (Food Listing   │
│   Verification)  │   │   Validation)    │
└──────────────────┘   └──────────────────┘
```

### Directory Structure
```
claim_service/
├── cmd/
│   └── server/
│       └── main.go               # Application entry point
├── internal/
│   ├── claim/
│   │   ├── handlers.go           # HTTP handlers
│   │   ├── service.go            # Business logic
│   │   ├── repository.go         # Data access layer
│   │   ├── models.go             # Domain models
│   │   └── client.go             # External service clients
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
| `PORT` | string | No | `8083` | HTTP server port |
| `DATABASE_URL` | string | Yes | - | PostgreSQL connection string |
| `JWT_SECRET` | string | Yes | - | Secret key for JWT validation (minimum 32 characters) |
| `USER_SERVICE_URL` | string | Yes | - | URL for user service (NGO verification) |
| `FOOD_SERVICE_URL` | string | Yes | - | URL for food service (listing validation) |
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
        Port:           getEnv("PORT", "8083"),
        DatabaseURL:    mustEnv("DATABASE_URL"),
        JWTSecret:      jwtSecret,
        UserServiceURL: mustEnv("USER_SERVICE_URL"),
        FoodServiceURL: mustEnv("FOOD_SERVICE_URL"),
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
FOOD_SERVICE_URL=http://localhost:8082
```

**Production:**
```env
ENV=production
LOG_LEVEL=info
USER_SERVICE_URL=http://user_service:8081
FOOD_SERVICE_URL=http://food_service:8082
```

---

## Database Schema

### SQL Migration (PostgreSQL)

```sql
-- Claims table
CREATE TABLE IF NOT EXISTS claims (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    food_listing_id UUID NOT NULL,
    ngo_user_id UUID NOT NULL,
    donor_user_id UUID NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'CREATED',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    accepted_at TIMESTAMP WITH TIME ZONE,
    rejected_at TIMESTAMP WITH TIME ZONE,
    picked_up_at TIMESTAMP WITH TIME ZONE,
    delivered_at TIMESTAMP WITH TIME ZONE,
    cancelled_at TIMESTAMP WITH TIME ZONE,
    CONSTRAINT fk_food_listing FOREIGN KEY (food_listing_id) REFERENCES food_listings(id) ON DELETE CASCADE,
    INDEX idx_claims_food_listing_id (food_listing_id),
    INDEX idx_claims_ngo_user_id (ngo_user_id),
    INDEX idx_claims_donor_user_id (donor_user_id),
    INDEX idx_claims_status (status),
    INDEX idx_claims_created_at (created_at),
    CONSTRAINT unique_active_claim UNIQUE (food_listing_id) WHERE status IN ('CREATED', 'ACCEPTED', 'PICKED_UP')
);
```

### Entity Relationship Diagram

```
┌─────────────────────────────┐
│          claims             │
│─────────────────────────────│
│ id (PK)                     │
│ food_listing_id (FK) ───────┼───┐
│ ngo_user_id (FK) ───────────┼───┼──┐
│ donor_user_id (FK) ─────────┼───┼──┼──┐
│ status                      │   │  │  │
│ created_at                  │   │  │  │
│ updated_at                  │   │  │  │
│ accepted_at                 │   │  │  │
│ rejected_at                 │   │  │  │
│ picked_up_at                │   │  │  │
│ delivered_at                │   │  │  │
│ cancelled_at                │   │  │  │
└─────────────────────────────┘   │  │  │
                                  │  │  │
                          ┌───────▼──┼──┼──┐
                          │ food_    │  │  │
                          │ listings │  │  │
                          │ (food_   │  │  │
                          │ service) │  │  │
                          └──────────┘  │  │
                                        │  │
                                ┌───────▼──┼──┐
                                │   users  │  │
                                │  (auth/  │  │
                                │   user   │  │
                                │ service) │  │
                                └──────────┘  │
                                              │
                                      ┌───────▼──┐
                                      │   users  │
                                      │  (auth/  │
                                      │   user   │
                                      │ service) │
                                      └──────────┘
```

---

## Core Components

### 1. Domain Models (`internal/claim/models.go`)

#### Claim
```go
type Claim struct {
    ID            string     `json:"id"`
    FoodListingID string     `json:"foodListingId"`
    NGOUserID     string     `json:"ngoUserId"`
    DonorUserID   string     `json:"donorUserId"`
    Status        ClaimStatus `json:"status"`
    CreatedAt     time.Time  `json:"createdAt"`
    UpdatedAt     time.Time  `json:"updatedAt"`
    AcceptedAt    *time.Time `json:"acceptedAt,omitempty"`
    RejectedAt    *time.Time `json:"rejectedAt,omitempty"`
    PickedUpAt    *time.Time `json:"pickedUpAt,omitempty"`
    DeliveredAt   *time.Time `json:"deliveredAt,omitempty"`
    CancelledAt   *time.Time `json:"cancelledAt,omitempty"`
}
```

#### Claim Status Constants
```go
type ClaimStatus string

const (
    ClaimStatusCreated   ClaimStatus = "CREATED"
    ClaimStatusAccepted  ClaimStatus = "ACCEPTED"
    ClaimStatusRejected  ClaimStatus = "REJECTED"
    ClaimStatusPickedUp  ClaimStatus = "PICKED_UP"
    ClaimStatusDelivered ClaimStatus = "DELIVERED"
    ClaimStatusCancelled ClaimStatus = "CANCELLED"
)
```

#### Actor Roles
```go
type ActorRole string

const (
    ActorNGO   ActorRole = "NGO"
    ActorDonor ActorRole = "DONOR"
)
```

### 2. State Machine Methods

#### IsActive
```go
func (c Claim) IsActive() bool
```
Returns true if claim is in `CREATED`, `ACCEPTED`, or `PICKED_UP` state.

#### IsTerminal
```go
func (c Claim) IsTerminal() bool
```
Returns true if claim is in `REJECTED`, `CANCELLED`, or `DELIVERED` state.

#### CanTransitionTo
```go
func (c Claim) CanTransitionTo(next ClaimStatus) bool
```
Validates state transitions:

| From | To |
|------|-----|
| CREATED | ACCEPTED, REJECTED, CANCELLED |
| ACCEPTED | PICKED_UP, CANCELLED |
| PICKED_UP | DELIVERED |
| Others | None |

#### CanBeModifiedBy
```go
func (c Claim) CanBeModifiedBy(actor ActorRole, userID string) bool
```
- NGOs can modify claims they created
- Donors can modify claims made on their listings

### 3. Service Layer (`internal/claim/service.go`)

The service layer implements all business logic and cross-service validation:

#### CreateClaim
```go
func (s *Service) CreateClaim(ctx context.Context, foodListingID string, ngoUserID string) (*Claim, error)
```
- Validates NGO is verified via User Service
- Validates food listing is available via Food Service
- Prevents self-claiming (NGO claiming their own food)
- Checks for existing active claims
- Creates claim with `CREATED` status

#### ApproveClaim (Donor Action)
```go
func (s *Service) ApproveClaim(ctx context.Context, claimID string, donorUserID string) error
```
- Validates donor owns the claim's food listing
- Transitions claim to `ACCEPTED`
- Marks food listing as `CLAIMED` via Food Service

#### RejectClaim (Donor Action)
```go
func (s *Service) RejectClaim(ctx context.Context, claimID string, donorUserID string) error
```
- Validates donor owns the claim's food listing
- Transitions claim to `REJECTED`

#### CancelByNGO (NGO Action)
```go
func (s *Service) CancelByNGO(ctx context.Context, claimID string, ngoUserID string) error
```
- Validates NGO owns the claim
- Transitions claim to `CANCELLED`

#### MarkPickedUp (NGO Action)
```go
func (s *Service) MarkPickedUp(ctx context.Context, claimID string, ngoUserID string) error
```
- Validates NGO owns the claim
- Transitions claim to `PICKED_UP`

#### MarkDelivered (NGO Action)
```go
func (s *Service) MarkDelivered(ctx context.Context, claimID string, ngoUserID string) error
```
- Validates NGO owns the claim
- Transitions claim to `DELIVERED`

### 4. Repository Layer (`internal/claim/repository.go`)

```go
type Repository interface {
    Create(ctx context.Context, claim *Claim) error
    GetByID(ctx context.Context, id string) (*Claim, error)
    GetActiveByFoodID(ctx context.Context, foodListingID string) (*Claim, error)
    GetByFoodID(ctx context.Context, foodListingID string) ([]Claim, error)
    GetByNGOUserID(ctx context.Context, ngoUserID string) ([]Claim, error)
    GetByDonorUserID(ctx context.Context, donorUserID string) ([]Claim, error)
    UpdateStatus(ctx context.Context, id string, status ClaimStatus) error
}
```

### 5. Service Clients (`internal/claim/client.go`)

**UserClient:**
```go
type UserClient struct {
    baseURL    string
    httpClient *http.Client
}

func (c *UserClient) IsVerifiedNGO(ctx context.Context, userID string) (bool, error)
```

**FoodClient:**
```go
type FoodClient struct {
    baseURL    string
    httpClient *http.Client
}

func (c *FoodClient) GetFoodForClaim(ctx context.Context, foodListingID string) (donorUserID string, status string, err error)
func (c *FoodClient) MarkFoodClaimed(ctx context.Context, foodListingID string) error
```

### 6. HTTP Handlers (`internal/claim/handlers.go`)

```go
type Handler struct {
    service *Service
}

// Handler methods
func (h *Handler) CreateClaim(w http.ResponseWriter, r *http.Request)
func (h *Handler) ApproveClaim(w http.ResponseWriter, r *http.Request)
func (h *Handler) RejectClaim(w http.ResponseWriter, r *http.Request)
func (h *Handler) CancelClaim(w http.ResponseWriter, r *http.Request)
func (h *Handler) MarkPickedUp(w http.ResponseWriter, r *http.Request)
func (h *Handler) MarkDelivered(w http.ResponseWriter, r *http.Request)
func (h *Handler) GetClaim(w http.ResponseWriter, r *http.Request)
func (h *Handler) GetClaimsByNGO(w http.ResponseWriter, r *http.Request)
func (h *Handler) GetClaimsByDonor(w http.ResponseWriter, r *http.Request)
func (h *Handler) GetClaimsByFood(w http.ResponseWriter, r *http.Request)
func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request)
```

---

## API Endpoints

### Base URL: `http://localhost:8083`

### Public Endpoints (No Authentication Required)

#### 1. Health Check
```http
GET /health
```
**Response:**
```json
{
  "status": "healthy",
  "service": "claim_service"
}
```

### Protected Endpoints (Requires Authentication)

#### 2. Create Claim (NGO Only)
```http
POST /claims
```
**Headers:**
```
Authorization: Bearer <accessToken>
```
**Request Body:**
```json
{
  "foodListingId": "550e8400-e29b-41d4-a716-446655440000"
}
```
**Response (201 Created):**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "foodListingId": "550e8400-e29b-41d4-a716-446655440000",
  "ngoUserId": "550e8400-e29b-41d4-a716-446655440000",
  "donorUserId": "550e8400-e29b-41d4-a716-446655440000",
  "status": "CREATED",
  "createdAt": "2026-08-08T10:00:00Z",
  "updatedAt": "2026-08-08T10:00:00Z",
  "acceptedAt": null,
  "rejectedAt": null,
  "pickedUpAt": null,
  "deliveredAt": null,
  "cancelledAt": null
}
```

#### 3. Get Claim by ID
```http
GET /claims/{id}
```
**Headers:**
```
Authorization: Bearer <accessToken>
```
**Response (200 OK):** Claim object

#### 4. Get Claims by NGO
```http
GET /claims/ngo
```
**Headers:**
```
Authorization: Bearer <accessToken>
```
**Response (200 OK):** Array of claims

#### 5. Get Claims by Donor
```http
GET /claims/donor
```
**Headers:**
```
Authorization: Bearer <accessToken>
```
**Response (200 OK):** Array of claims

#### 6. Get Claims by Food Listing
```http
GET /foods/{foodId}/claims
```
**Headers:**
```
Authorization: Bearer <accessToken>
```
**Response (200 OK):** Array of claims

#### 7. Approve Claim (Donor Only)
```http
POST /claims/{id}/approve
```
**Headers:**
```
Authorization: Bearer <accessToken>
```
**Response:** 204 No Content

#### 8. Reject Claim (Donor Only)
```http
POST /claims/{id}/reject
```
**Headers:**
```
Authorization: Bearer <accessToken>
```
**Response:** 204 No Content

#### 9. Cancel Claim (NGO Only)
```http
POST /claims/{id}/cancel
```
**Headers:**
```
Authorization: Bearer <accessToken>
```
**Response:** 204 No Content

#### 10. Mark as Picked Up (NGO Only)
```http
POST /claims/{id}/pickup
```
**Headers:**
```
Authorization: Bearer <accessToken>
```
**Response:** 204 No Content

#### 11. Mark as Delivered (NGO Only)
```http
POST /claims/{id}/deliver
```
**Headers:**
```
Authorization: Bearer <accessToken>
```
**Response:** 204 No Content

---

## Security Implementation

### 1. Authentication Middleware

Same pattern as other services:

```go
// internal/middleware/auth_middleware.go
func Auth(jwtSecret string) func(http.Handler) http.Handler {
    // JWT validation and context injection
}
```

**Context Values:**
- `user_id`: User's UUID
- `user_role`: User's role

### 2. Actor-Based Authorization

Each endpoint checks the user's role:

```go
// CreateClaim - NGO only
if middleware.UserRole(ctx) != string(ActorNGO) {
    writeError(w, http.StatusForbidden, "insufficient permissions")
    return
}

// ApproveClaim - Donor only
if middleware.UserRole(ctx) != string(ActorDonor) {
    writeError(w, http.StatusForbidden, "insufficient permissions")
    return
}
```

### 3. Ownership Validation

Service layer validates the user owns the resource:

```go
func (c Claim) CanBeModifiedBy(actor ActorRole, userID string) bool {
    switch actor {
    case ActorNGO:
        return c.NGOUserID == userID
    case ActorDonor:
        return c.DonorUserID == userID
    default:
        return false
    }
}
```

### 4. State Machine Validation

Before any status change, the service validates the transition:

```go
func (c Claim) CanTransitionTo(next ClaimStatus) bool {
    switch c.Status {
    case ClaimStatusCreated:
        return next == ClaimStatusAccepted ||
               next == ClaimStatusRejected ||
               next == ClaimStatusCancelled
    case ClaimStatusAccepted:
        return next == ClaimStatusPickedUp ||
               next == ClaimStatusCancelled
    case ClaimStatusPickedUp:
        return next == ClaimStatusDelivered
    default:
        return false
    }
}
```

### 5. Cross-Service Validation

**NGO Verification:**
```go
verified, err := s.userClient.IsVerifiedNGO(ctx, ngoUserID)
if !verified {
    return nil, ErrNGONotVerified
}
```

**Food Listing Validation:**
```go
donorUserID, foodStatus, err := s.foodClient.GetFoodForClaim(ctx, foodListingID)
if foodStatus != "AVAILABLE" {
    return nil, ErrFoodNotOpen
}
```

### 6. Duplicate Claim Prevention

```go
existing, err := s.repo.GetActiveByFoodID(ctx, foodListingID)
if existing != nil {
    return nil, ErrActiveClaimExists
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
| 401 | Unauthorized |
| 403 | Forbidden |
| 404 | Not Found |
| 409 | Conflict |

### Common Error Messages

| Error | Scenario |
|-------|----------|
| "invalid request body" | Malformed JSON |
| "foodListingId is required" | Missing required field |
| "unauthorized action" | User doesn't own the claim |
| "invalid claim state transition" | Attempting invalid status change |
| "active claim already exists" | Food listing already has active claim |
| "ngo is not verified" | NGO not verified by admin |
| "food listing is not open" | Food listing not available |
| "ngo cannot claim its own food listing" | Self-claim attempt |
| "claim not found" | Claim ID doesn't exist |

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

**Claim Creation:**
```json
{
  "time": "2026-08-08T10:00:00Z",
  "level": "INFO",
  "msg": "Claim created",
  "service": "claim_service",
  "claim_id": "550e8400-...",
  "food_listing_id": "550e8400-...",
  "ngo_user_id": "550e8400-..."
}
```

**Status Transition:**
```json
{
  "time": "2026-08-08T10:01:00Z",
  "level": "INFO",
  "msg": "Claim status updated",
  "service": "claim_service",
  "claim_id": "550e8400-...",
  "from_status": "CREATED",
  "to_status": "ACCEPTED",
  "actor": "DONOR",
  "actor_user_id": "550e8400-..."
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
    go build -ldflags="-w -s" -o claim_service ./cmd/server

# Runtime stage
FROM alpine:latest
WORKDIR /app
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
COPY --from=builder /app/claim_service .
USER appuser
EXPOSE 8083
HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
  CMD wget --spider -q http://localhost:8083/health || exit 1
CMD ["./claim_service"]
```

### Docker Compose Example

```yaml
claim_service:
  build: .
  ports:
    - "8083:8083"
  environment:
    - DATABASE_URL=postgresql://gratia:gratia123@postgres:5432/gratia?sslmode=disable
    - JWT_SECRET=${JWT_SECRET}
    - USER_SERVICE_URL=http://user_service:8081
    - FOOD_SERVICE_URL=http://food_service:8082
    - ENV=production
    - LOG_LEVEL=info
  depends_on:
    - postgres
    - user_service
    - food_service
  networks:
    - gratia_network
```

### Graceful Shutdown

The service implements graceful shutdown with a 10-second timeout:

```go
// main.go
quit := make(chan os.Signal, 1)
signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
<-quit

shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
defer shutdownCancel()

if err := httpServer.Shutdown(shutdownCtx); err != nil {
    log.Error("server shutdown failed", slog.String("error", err.Error()))
    os.Exit(1)
}
```

---

## Service Integration

### Internal APIs Called by Claim Service

| Service | Endpoint | Purpose |
|---------|----------|---------|
| User Service | `GET /internal/users/{id}/ngo` | Check if user is verified NGO |
| Food Service | `GET /internal/foods/{id}/validate` | Validate food listing availability |
| Food Service | `PATCH /internal/foods/{id}/claim` | Mark food as claimed |

### External Dependencies

| Service | Purpose | URL Variable |
|---------|---------|--------------|
| Auth Service | JWT Validation | `JWT_SECRET` (shared) |
| User Service | NGO Verification | `USER_SERVICE_URL` |
| Food Service | Food Validation & Claim Marking | `FOOD_SERVICE_URL` |

### Client Interfaces

```go
// UserClient interface
func (c *UserClient) IsVerifiedNGO(ctx context.Context, userID string) (bool, error)

// FoodClient interface
func (c *FoodClient) GetFoodForClaim(ctx context.Context, foodListingID string) (donorUserID string, status string, err error)
func (c *FoodClient) MarkFoodClaimed(ctx context.Context, foodListingID string) error
```

---

## Development Guidelines

### Adding New Features

#### 1. New Claim Status
1. Add status constant in `models.go`
2. Update `CanTransitionTo` method
3. Add timestamp column in database
4. Update `UpdateStatus` repository method
5. Update API documentation

#### 2. New Actor Role
1. Add role constant in `models.go`
2. Update `CanBeModifiedBy` method
3. Add role checks in handlers

### Testing

**Unit Test Example:**
```go
func TestCreateClaim(t *testing.T) {
    mockRepo := &MockRepository{}
    mockUserClient := &MockUserClient{verified: true}
    mockFoodClient := &MockFoodClient{status: "AVAILABLE", donorID: "donor-123"}
    
    service := NewService(mockRepo, mockUserClient, mockFoodClient)
    
    claim, err := service.CreateClaim(
        context.Background(),
        "food-123",
        "ngo-123",
    )
    
    assert.NoError(t, err)
    assert.Equal(t, ClaimStatusCreated, claim.Status)
    assert.Equal(t, "food-123", claim.FoodListingID)
    assert.Equal(t, "ngo-123", claim.NGOUserID)
}
```

### Performance Considerations

1. **Cross-Service Calls**: 5-second timeout prevents cascading failures
2. **Database Indexes**: Proper indexes on `food_listing_id`, `ngo_user_id`, `donor_user_id`
3. **Idempotent Operations**: Unique constraint on active claims
4. **Connection Pooling**: pgxpool for efficient database connections

---

## Future Enhancements

1. **Webhook Notifications**: Notify donors/NGOs on status changes
2. **Claim History**: Audit trail for all status changes
3. **Batch Operations**: Bulk approve/reject claims
4. **Rating System**: Allow donors to rate NGOs after delivery
5. **Claim Analytics**: Metrics dashboard for claims
6. **Push Notifications**: Real-time updates via WebSocket
7. **Claim Expiry**: Auto-cancel claims after timeout
8. **SLA Tracking**: Time tracking for each claim stage

---

## Dependencies

| Package | Version | Purpose |
|---------|---------|---------|
| `github.com/jackc/pgx/v5` | v5.5.0 | PostgreSQL driver |
| `github.com/golang-jwt/jwt/v5` | v5.0.0 | JWT token validation |
| `github.com/joho/godotenv` | v1.5.1 | Environment variables |
| `github.com/google/uuid` | latest | UUID generation |