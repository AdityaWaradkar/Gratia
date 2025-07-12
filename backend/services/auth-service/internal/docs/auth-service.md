
# Gratia Auth Service Documentation

## Overview

The **Auth Service** is a microservice responsible for user authentication, registration, session management, and secure token issuance using JWT. It supports role-based access control for the following user roles:

- `restaurant`

- `ngo`

- `admin`

It exposes REST API endpoints for registration, login, token refresh, user info retrieval, and logout.

***


## Folder Structure

    auth-service/
    │
    ├── cmd/
    │   └── server/
    │       └── main.go              # Entry point to start the auth server
    │
    ├── internal/
    │   ├── config/
    │   │   └── config.go            # Configuration loading from environment
    │   ├── handlers/
    │   │   └── auth_handler.go      # HTTP handlers for auth endpoints
    │   ├── middleware/
    │   │   └── jwt_middleware.go    # JWT validation middleware
    │   ├── models/
    │   │   └── user.go              # User model struct
    │   ├── repository/
    │   │   └── user_repository.go   # User database CRUD operations
    │   ├── service/
    │   │   └── auth.go              # Business logic: registration, login, JWT
    │   └── utils/
    │       └── jwt.go               # JWT token creation and validation utils
    │
    ├── docs/
    │   └── openapi.yaml             # OpenAPI 3.0 specification for API documentation
    │
    ├── .env                        # Environment variables (DB URL, JWT secret, port)
    ├── Dockerfile                  # Docker multi-stage build for production image
    ├── go.mod
    └── go.sum

***


## Key Components

### 1. `cmd/server/main.go`

- Loads environment variables.

- Connects to PostgreSQL database using `sqlx`.

- Initializes repository, service, and HTTP handlers.

- Starts HTTP server on configured port (`8081` default).

- Adds CORS middleware for cross-origin requests.

- Handles graceful shutdown on interrupt signal.

***


### 2. Configuration (`internal/config/config.go`)

- Loads configuration from `.env` or system environment.

- Required variables: `DATABASE_URL`, `JWT_SECRET`, `SERVER_PORT`.

- Validates presence of required configs, else exits.

***


### 3. Models (`internal/models/user.go`)

- Defines the `User` struct with fields like:

  - `UserID` (UUID)

  - `Email`, `FullName`, `PhoneNumber`

  - `PasswordHash` (not JSON exposed)

  - `Role` (restaurant/ngo/admin)

  - `IsActive` (bool pointer)

  - Timestamps (`CreatedAt`, `UpdatedAt`)

- Supports JSON serialization and DB mapping.

***


### 4. Repository (`internal/repository/user_repository.go`)

- Implements `UserRepository` interface with methods:

  - `CreateUser(ctx, user)`

  - `GetUserByEmail(ctx, email)`

  - `GetUserByID(ctx, userID)`

- Uses `sqlx` to query Postgres `users` table.

- Inserts new users and retrieves users by email or ID.

***


### 5. Service Layer (`internal/service/auth.go`)

- Contains business logic for auth workflows.

- Methods:

  - `Register(ctx, user, password)`:

    - Checks if email already exists.

    - Hashes password with bcrypt.

    - Inserts user with `IsActive=true` by default.

  - `Login(ctx, email, password)`:

    - Validates credentials.

    - Returns access and refresh JWT tokens.

  - `GenerateJWT(user, duration)`:

    - Creates signed JWT with claims `sub`, `email`, `role`, and expiry.

  - `GetUserByID(ctx, id)`:

    - Retrieves user details by UUID.

- Exposes errors for common cases (`ErrUserExists`, `ErrInvalidCredentials`).

***


### 6. HTTP Handlers (`internal/handlers/auth_handler.go`)

- Defines `AuthHandler` struct that wraps the service.

- Implements endpoints:

  - `POST /register` — Registers a new user.

  - `POST /login` — Authenticates user, returns tokens.

  - `POST /refresh` — Refreshes access token using refresh token.

  - `GET /me` — Returns current user profile info (protected).

  - `POST /logout` — Logs out the user (protected).

  - `POST /forgot-password` — Stub for password reset initiation.

  - `POST /reset-password` — Stub for password reset completion.

- Uses JSON for request/response.

- Handles input validation and errors.

- Uses JWT middleware for protected routes.

***


### 7. Middleware (`internal/middleware/jwt_middleware.go`)

- JWT Authentication middleware extracts Bearer token.

- Validates token using secret and `utils.ValidateJWT`.

- Stores user ID from JWT claims in request context.

- Rejects requests with invalid or missing tokens.

***


### 8. Utility Functions (`internal/utils/jwt.go`)

- Functions to create and validate JWT tokens.

- Uses `github.com/golang-jwt/jwt/v4`.

- Encapsulates JWT secret management.

***


### 9. OpenAPI Documentation (`docs/openapi.yaml`)

- OpenAPI 3 specification of all auth endpoints.

- Includes request/response schemas, error codes, and security schemes.

- Supports tools like Swagger UI for interactive docs.

***


### 10. Dockerfile

- Multi-stage build:

  - Builds Go binary in golang:1.24.4 image.

  - Produces minimal alpine image with compiled binary.

- Exposes port 8081.

- Runs `./auth-service` on container start.

***


## API Endpoints Summary

| Method | Path                           | Description                      | Auth Required |
| ------ | ------------------------------ | -------------------------------- | ------------- |
| POST   | `/api/v1/auth/register`        | Register new user                | No            |
| POST   | `/api/v1/auth/login`           | Authenticate and get tokens      | No            |
| POST   | `/api/v1/auth/refresh`         | Refresh access token             | No            |
| GET    | `/api/v1/auth/me`              | Get current user profile         | Yes (JWT)     |
| POST   | `/api/v1/auth/logout`          | Logout user (invalidate session) | Yes (JWT)     |
| POST   | `/api/v1/auth/forgot-password` | Initiate password reset (stub)   | No            |
| POST   | `/api/v1/auth/reset-password`  | Reset password (stub)            | No            |

***


## Security Considerations

- Passwords hashed with bcrypt with salt.

- JWT tokens signed securely; access tokens expire after 15 minutes.

- Refresh tokens expire after 7 days.

- Middleware validates JWT tokens on protected routes.

- CORS configured to allow cross-origin requests.

- Rate limiting can be added on sensitive endpoints (future enhancement).

- Service designed to be deployed behind HTTPS.

***


## How to Run Locally

1. Create `.env` file with the following:

       DATABASE_URL=<your_postgres_connection_string>
       JWT_SECRET=<your_jwt_secret_key>
       SERVER_PORT=8081

2. Start PostgreSQL with a `users` table (see schema below).

3. Build and run service:

   ```bash
   go build -o auth-service ./cmd/server/main.go
   ./auth-service
   ```

4. Alternatively, build and run using Docker:

   ```bash
   docker build -t gratia-auth-service .
   docker run -p 8081:8081 gratia-auth-service
   ```

***


## Database Schema (PostgreSQL)

```sql
CREATE TABLE public.users (
    user_id UUID PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    full_name VARCHAR(255),
    phone_number VARCHAR(20),
    role VARCHAR(50) NOT NULL,
    location_lat FLOAT8,
    location_lng FLOAT8,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE
);
```

***


## Future Enhancements

- Password reset and recovery workflows.

- Email verification on registration.

- OAuth support (Google, GitHub).

- Two-factor authentication (2FA).

- Rate limiting middleware.

- Admin impersonation for troubleshooting.

- Detailed logging and monitoring integration.

***


## Contact & Support

For issues or questions, reach out to:

**Aditya Waradkar**\
Email: [aditya](mailto:aditya@example.com)waradkar1801\@gmail.com\
GitHub: <https://github.com/AdityaWaradkar>

***
