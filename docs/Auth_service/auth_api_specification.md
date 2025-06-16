# Auth Service API Specification

***


## Overview

The Auth Service handles user authentication, registration, and secure token issuance using JWT (JSON Web Tokens). It supports the following roles:

- `restaurant`

- `ngo`

- `admin`

It provides secure mechanisms for login, logout, role-based access control (RBAC), and session management. All sensitive operations are protected by JWT authentication and role-based authorization.

***


## Base URL

    /api/v1/auth

***


## Endpoints

***


### 1. Register User

**POST** `/register`

Registers a new user with role-specific privileges.


#### Request Body

```json
{
  "full_name": "Aditya Waradkar",
  "email": "aditya@example.com",
  "password": "MySecurePass123!",
  "phone_number": "+919876543210",
  "role": "restaurant"
}
```


#### Response (201 Created)

```json
{
  "message": "User registered successfully.",
  "user_id": "b5dfc0e2-9d1e-45ac-a43c-e687b1c6ecdb"
}
```


#### Errors

| Status | Message               |
| ------ | --------------------- |
| 400    | Invalid input data    |
| 409    | Email already in use  |
| 500    | Internal server error |

***


### 2. Login

**POST** `/login`

Authenticates the user and returns access and refresh tokens.


#### Request Body

```json
{
  "email": "aditya@example.com",
  "password": "MySecurePass123!"
}
```


#### Response (200 OK)

```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6...",
  "refresh_token": "refresh.token.value",
  "expires_in": 900,
  "token_type": "Bearer"
}
```


#### Errors

| Status | Message                   |
| ------ | ------------------------- |
| 401    | Invalid email or password |
| 403    | User account is disabled  |
| 500    | Internal server error     |

***


### 3. Get Current User Info

**GET** `/me`

Returns the authenticated user's profile details.


#### Headers

    Authorization: Bearer <access_token>


#### Response (200 OK)

```json
{
  "user_id": "b5dfc0e2-9d1e-45ac-a43c-e687b1c6ecdb",
  "full_name": "Aditya Waradkar",
  "email": "aditya@example.com",
  "role": "restaurant",
  "phone_number": "+919876543210",
  "is_active": true,
  "created_at": "2025-06-01T14:32:00Z"
}
```


#### Errors

| Status | Message                  |
| ------ | ------------------------ |
| 401    | Invalid or expired token |
| 500    | Internal server error    |

***


### 4. Logout

**POST** `/logout`

Revokes the user's current session token.


#### Headers

    Authorization: Bearer <access_token>


#### Response (200 OK)

```json
{
  "message": "Logged out successfully."
}
```


#### Errors

| Status | Message                  |
| ------ | ------------------------ |
| 401    | Invalid or expired token |
| 500    | Internal server error    |

***


### 5. Refresh Token

**POST** `/refresh`

Issues a new access token using a valid refresh token.


#### Request Body

```json
{
  "refresh_token": "valid.refresh.token"
}
```


#### Response (200 OK)

```json
{
  "access_token": "new.jwt.token",
  "expires_in": 900,
  "token_type": "Bearer"
}
```

***


## Authentication and Authorization

- JWT is used for stateless authentication.

- Tokens include user ID, email, role, and expiry timestamp.

- Middleware is responsible for validating tokens and enforcing role-based access.


### Example JWT Claims

```json
{
  "sub": "b5dfc0e2-9d1e-45ac-a43c-e687b1c6ecdb",
  "email": "aditya@example.com",
  "role": "restaurant",
  "exp": 1717451200
}
```

***


## Password Management (Planned)

### Forgot Password

**POST** `/forgot-password`

Initiates the password reset process by sending a reset link to the user's email.


### Reset Password

**POST** `/reset-password`

Accepts a reset token and a new password to update the user account.

***


## Status Codes Summary

| Code | Meaning                    |
| ---- | -------------------------- |
| 200  | OK                         |
| 201  | Created                    |
| 400  | Bad Request                |
| 401  | Unauthorized               |
| 403  | Forbidden                  |
| 409  | Conflict (e.g., duplicate) |
| 500  | Internal Server Error      |

***


## Security Considerations

- Passwords are hashed using bcrypt with salt.

- All traffic is served over HTTPS (TLS 1.3 enforced).

- JWT tokens are signed using a secure secret key.

- Access token lifespan: 15 minutes.

- Refresh token lifespan: 7 days.

- Rate limiting is applied to authentication endpoints to prevent brute-force attacks.

***


## Future Enhancements

- Password reset and recovery features

- Email verification for new accounts

- OAuth login support (Google, GitHub)

- Two-factor authentication (2FA)

- Admin impersonation for troubleshooting and support
