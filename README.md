# Go Authentication Service (`auth_service`)

This repository contains a modular and secure **Authentication Service** built in **Go (Golang)**. It implements a comprehensive authentication and authorization solution using JWTs, refresh tokens, and a PostgreSQL database. The service adheres to a layered architecture (Handler -> Service -> Repository) for clarity and maintainability.

---

## Technology Stack

* **Language:** Go (Golang)
* **Database:** PostgreSQL (via `pgxpool`)
* **Authentication:** JWT (HS256) and Refresh Tokens
* **Security:** `bcrypt` for password hashing

---

## Features

### 1. Core Authentication and User Management
* **Registration & Login:** Secure creation and authentication of user accounts.
* **JWT Generation:** Issuance of time-bound Access Tokens and persistent Refresh Tokens.
* **Role-Based Access Control (RBAC):** Supports `USER`, `MOD`, and `ADMIN` roles, with authorization checks on privileged actions (e.g., role assignment).
* **User Retrieval:** Endpoint to fetch the profile of the currently authenticated user.

### 2. Session Management
* **Session Tracking:** Records user agent and IP address for active sessions.
* **Token Refresh:** Securely generates new tokens and revokes old ones upon refresh.
* **Logout & Control:** Provides mechanisms for single-session logout and retrieval/deletion of all active sessions.

### 3. Recovery and Verification
* **Password Reset Flow:** Implements secure token generation (`ForgotPassword`) and consumption (`ResetPassword`) for password recovery.
* **Email Verification:** Functionality to generate, send (simulated), and validate email verification tokens.

---

## API Endpoints (14 Total)

The service exposes the following complete set of handler methods:

| Category | Method | Handler Function | Description |
| :--- | :--- | :--- | :--- |
| **Authentication** | `POST` | `RegisterUser` | Creates a new user account with role assignment logic. |
| **Authentication** | `POST` | `LoginUser` | Authenticates user and issues Access/Refresh token pair. |
| **Authentication** | `POST` | `RefreshTokens` | Generates a new Access Token using a valid Refresh Token. |
| **Authentication** | `POST` | `Logout` | Revokes a Refresh Token and deletes the associated session. |
| **User Info** | `GET` | `GetCurrentUser` | Retrieves the profile details of the authenticated user. |
| **Session Control**| `GET` | `GetSessions` | Lists all active sessions for the authenticated user. |
| **Session Control**| `DELETE`| `DeleteSession` | Terminates a specific user session by ID. |
| **Password Reset** | `POST` | `ForgotPassword` | Initiates password recovery by generating a reset token. |
| **Password Reset** | `POST` | `ResetPassword` | Sets a new password using a valid reset token. |
| **Email Verify** | `POST` | `VerifyEmail` | Consumes a token to mark the user's email as verified. |
| **Email Verify** | `POST` | `GenerateEmailVerification` | Generates a verification token for a given email (utility). |
| **Email Verify** | `POST` | `ResendEmailVerification` | Generates and returns a new verification token for an unverified user. |
| **Utility** | `POST` | `ValidateTokenHandler` | Validates an Access Token and returns its claims. |
| **Utility** | `GET` | `HealthCheck` | Checks service operational status and database connectivity. |

---

## Design Principles

The service is built on the following principles:

* **Separation of Concerns:** Business logic resides entirely within the `Service` layer, database operations in the `Repository`, and HTTP concerns in the `Handler`.
* **Context Passing:** `context.Context` is used throughout the layers to manage timeouts and propagate request-scoped values (like `UserID` and `UserRole` from middleware).
* **Database Transactions:** Email verification involves a transaction to ensure both the user's status is updated and the token is marked as used atomically.
