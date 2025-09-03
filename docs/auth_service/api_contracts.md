# Gratia Auth Service – API Contracts Documentation

---

## Introduction

The Auth Service is responsible for managing authentication and authorization within the Gratia platform. It ensures that only valid users and services can access protected resources. This document defines the set of APIs provided by the Auth Service, their expected inputs and outputs, error handling conventions, and security considerations.

***


## Responsibilities

The Auth Service performs the following functions:

1. User registration

2. User login and token issuance

3. Access token refresh

4. Logout and session invalidation

5. Email verification

6. Password reset and recovery

7. Session management across devices

8. Token validation for internal services

9. Public key exposure for JWT validation

***


## API Endpoints

### 1. Register – `POST /auth/register`

Registers a new user in the system.

**Request**

```json
{
  "email": "ngo@example.org",
  "password": "StrongP@ssw0rd!",
  "role": "NGO"
}
```

**Response (201 Created)**

```json
{
  "user": {
    "id": "uuid",
    "email": "ngo@example.org",
    "role": "NGO",
    "emailVerified": false,
    "createdAt": "2025-09-03T10:00:00Z"
  }
}
```

**Error Codes**

- 400 Bad Request: Invalid input

- 409 Conflict: Email already exists

***


### 2. Login – `POST /auth/login`

Authenticates a user and issues access and refresh tokens.

**Request**

```json
{
  "email": "ngo@example.org",
  "password": "StrongP@ssw0rd!"
}
```

**Response (200 OK)**

```json
{
  "accessToken": "jwt_token_here",
  "expiresIn": 900,
  "refreshToken": "refresh_token_here",
  "refreshExpiresIn": 2592000,
  "user": {
    "id": "uuid",
    "email": "ngo@example.org",
    "role": "NGO",
    "emailVerified": false
  }
}
```

**Error Codes**

- 401 Unauthorized: Invalid credentials

***


### 3. Refresh Token – `POST /auth/refresh`

Issues a new access token and refresh token using a valid refresh token.

**Request**

```json
{
  "refreshToken": "refresh_token_here"
}
```

**Response (200 OK)**

```json
{
  "accessToken": "new_jwt",
  "expiresIn": 900,
  "refreshToken": "new_refresh",
  "refreshExpiresIn": 2592000
}
```

***


### 4. Logout – `POST /auth/logout`

Invalidates the provided refresh token and ends the session.

**Request**

```json
{
  "refreshToken": "refresh_token_here"
}
```

**Response**\
204 No Content

***


### 5. Get Current User – `GET /auth/me`

Returns details of the authenticated user.

**Headers**\
`Authorization: Bearer <access_token>`

**Response (200 OK)**

```json
{
  "id": "uuid",
  "email": "ngo@example.org",
  "role": "NGO",
  "emailVerified": true
}
```

***


### 6. Forgot Password – `POST /auth/forgot-password`

Initiates the password reset process.

**Request**

```json
{
  "email": "ngo@example.org"
}
```

**Response**\
202 Accepted

***


### 7. Reset Password – `POST /auth/reset-password`

Resets the password using a reset token sent via email.

**Request**

```json
{
  "resetToken": "token_from_email",
  "newPassword": "NewP@ssword123"
}
```

**Response**\
204 No Content

***


### 8. Verify Email – `POST /auth/verify-email`

Marks a user’s email as verified.

**Request**

```json
{
  "verificationToken": "token_from_email"
}
```

**Response**\
204 No Content

***


### 9. Validate Token – `GET /auth/validate`

Checks the validity of a given token.

**Request**\
`/auth/validate?token=<jwt_here>`

**Response (200 OK)**

```json
{
  "active": true,
  "sub": "uuid",
  "role": "NGO",
  "exp": 1735982400
}
```

***


### 10. Public Keys – `GET /auth/keys`

Provides public keys for JWT validation.

**Response (200 OK)**

```json
{
  "keys": [
    {
      "kty": "RSA",
      "kid": "key-2025",
      "n": "...",
      "e": "AQAB"
    }
  ]
}
```

***


### 11. List Sessions – `GET /auth/sessions`

Returns all active sessions for the authenticated user.

**Response (200 OK)**

```json
{
  "sessions": [
    {
      "id": "sess123",
      "createdAt": "2025-09-01T10:00:00Z",
      "userAgent": "Chrome on Windows",
      "ip": "203.0.113.10",
      "current": true
    },
    {
      "id": "sess456",
      "createdAt": "2025-08-28T08:00:00Z",
      "userAgent": "Firefox on Mobile",
      "ip": "203.0.113.11",
      "current": false
    }
  ]
}
```

***


### 12. Revoke Session – `DELETE /auth/sessions/{id}`

Terminates a specific session for the authenticated user.

**Response**\
204 No Content

***


## Error Handling

All errors follow a consistent structure:

```json
{
  "error": {
    "code": "INVALID_CREDENTIALS",
    "message": "Email or password is incorrect",
    "details": [
      { "field": "email", "issue": "not_found" }
    ],
    "traceId": "uuid"
  }
}
```

***


## Security Considerations

- Passwords are securely hashed using Argon2id or bcrypt.

- Access tokens are short-lived (e.g., 15 minutes).

- Refresh tokens are long-lived (e.g., 30 days) and rotated on use.

- All cookies, if used, must be `HttpOnly`, `Secure`, and `SameSite=Strict`.

- Login and registration endpoints should be rate-limited to prevent brute force attacks.

- Errors include a `traceId` to assist in monitoring and debugging.

***
