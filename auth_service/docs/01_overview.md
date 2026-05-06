# Overview

The `auth_service` is the centralized authentication and authorization service of the Gratia platform. It acts as the primary security layer of the system and is responsible for managing user identity, authentication workflows, JWT token generation, and role-based access control across all microservices.

In the overall architecture, this service establishes trust between users and backend services by ensuring that every protected request is authenticated and properly authorized before access is granted.

The service is intentionally designed as a lightweight, stateless, and infrastructure-oriented component with strictly limited ownership boundaries. Its responsibilities are focused entirely on authentication and authorization concerns and do not include business-domain operations such as user profile management, food listing management, NGO verification, or claim lifecycle handling.

---

# Core Objectives

The primary objectives of the `auth_service` are:

- Authenticate users securely
    
- Generate and validate JWT access tokens
    
- Enforce Role-Based Access Control (RBAC)
    
- Provide reusable authentication middleware
    
- Maintain minimal identity-related data
    
- Establish secure inter-service trust
    
- Support stateless authentication architecture
    

The service forms the security foundation for the entire Gratia platform and is required by all protected backend services.

---

# Responsibilities

The `auth_service` exclusively owns all authentication and authorization-related operations within the platform.

## User Authentication

The service validates user credentials during login operations.

### Responsibilities

- Accept login requests
    
- Validate user credentials
    
- Verify passwords using bcrypt
    
- Reject invalid authentication attempts
    
- Generate JWT tokens after successful authentication
    

### Authentication Flow

1. Client submits email and password
    
2. Service retrieves user credentials from the database
    
3. Stored bcrypt password hash is fetched
    
4. Password verification is performed securely
    
5. JWT token is generated upon successful authentication
    
6. Authentication response is returned to the client
    

### Security Measures

- Passwords are never stored in plain text
    
- bcrypt is used for password hashing
    
- Authentication errors remain generic
    
- Password hashes are never exposed externally
    

---

## JWT Token Management

The `auth_service` acts as the sole authority responsible for JWT token generation and validation.

### Responsibilities

- Generate signed JWT access tokens
    
- Embed user identity claims into tokens
    
- Validate token signatures
    
- Verify token expiration
    
- Extract authenticated user context
    

### JWT Claims

|Claim|Description|
|---|---|
|user_id|Unique user identifier|
|role|Authorization role|
|exp|Token expiration timestamp|
|iat|Token issued timestamp|

### JWT Strategy

The platform uses stateless JWT authentication to support distributed microservices communication without centralized session storage.

### Recommended JWT Configuration

|Property|Configuration|
|---|---|
|Signing Algorithm|HS256|
|Token Expiration|15–60 minutes|
|Secret Storage|Environment variables|
|Transport Method|Authorization header|

---

## Role-Based Access Control (RBAC)

The `auth_service` enforces authorization policies using role-based access control.

### Supported Roles

|Role|Description|
|---|---|
|DONOR|Food donor|
|NGO|Verified NGO user|
|ADMIN|Administrative authority|

### Responsibilities

- Validate user roles
    
- Restrict protected endpoints
    
- Enforce authorization rules
    
- Support middleware-level access control
    

### Example Access Restrictions

|Endpoint|Allowed Roles|
|---|---|
|/admin/*|ADMIN|
|/food/create|DONOR|
|/claim/create|NGO|

Authorization checks are performed only after successful authentication.

---

## Identity Bootstrapping

The service is responsible for initial system bootstrapping during deployment.

### Responsibilities

- Create the default administrative account
    
- Ensure the platform is operable after deployment
    
- Enable initial system administration
    

### Implementation Strategy

A pre-configured admin account is created during database migration execution to eliminate manual database modification after deployment.

---

# Ownership Boundaries

The `auth_service` only manages authentication and authorization concerns.

## Responsibilities Owned by auth_service

- Signup
    
- Login
    
- Password hashing
    
- JWT generation
    
- JWT validation
    
- Role validation
    
- Authentication middleware
    
- Authorization middleware
    

## Responsibilities Not Owned by auth_service

|Responsibility|Owning Service|
|---|---|
|Donor profiles|`user_service`|
|NGO profiles|`user_service`|
|NGO verification|`user_service`|
|Food listings|`food_service`|
|Claim lifecycle|`claim_service`|
|Delivery tracking|`claim_service`|
|Business workflows|Domain services|

This separation ensures low coupling and clear service ownership across the distributed system.

---

# Data Ownership

The `auth_service` maintains only minimal identity-related data required for authentication and authorization.

The database schema is intentionally designed to remain lightweight and security-focused.

---

# Database Design

## users Table

The `users` table acts as the primary authentication table.

|Field|Type|Description|
|---|---|---|
|user_id|UUID|Primary user identifier|
|email|VARCHAR|Unique login identifier|
|password_hash|VARCHAR|bcrypt hashed password|
|role|ENUM|DONOR / NGO / ADMIN|
|created_at|TIMESTAMP|Record creation timestamp|
|updated_at|TIMESTAMP|Record update timestamp|

---

## Database Constraints

|Constraint|Purpose|
|---|---|
|PRIMARY KEY(user_id)|Entity uniqueness|
|UNIQUE(email)|Prevent duplicate accounts|
|NOT NULL(password_hash)|Maintain authentication integrity|
|NOT NULL(role)|Maintain authorization integrity|

---

## Optional refresh_tokens Table

For future extensibility, a `refresh_tokens` table may be introduced.

|Field|Description|
|---|---|
|token_id|Unique token identifier|
|user_id|Associated user|
|expires_at|Token expiration|
|revoked|Revocation status|

This table is optional for the initial implementation phase.

---

# API Responsibilities

The `auth_service` exposes only authentication-related APIs.

## POST /auth/signup

Registers new platform users.

### Request Body

```json
{
  "email": "user@example.com",
  "password": "securePassword",
  "role": "DONOR"
}
```

### Response

```json
{
  "message": "User created successfully"
}
```

---

## POST /auth/login

Authenticates users and issues JWT tokens.

### Request Body

```json
{
  "email": "user@example.com",
  "password": "securePassword"
}
```

### Response

```json
{
  "access_token": "<jwt>",
  "expires_in": 3600
}
```

---

## GET /auth/validate

Validates JWT authenticity.

### Request Headers

```http
Authorization: Bearer <token>
```

### Response

```json
{
  "user_id": "uuid",
  "role": "DONOR"
}
```

---

# Middleware Responsibilities

The `auth_service` provides reusable middleware components for authentication and authorization.

## Authentication Middleware

### Responsibilities

- Extract JWT from request headers
    
- Validate token signature
    
- Validate token expiration
    
- Parse JWT claims
    
- Reject unauthorized requests
    

### Request Context

Authenticated requests expose:

- user_id
    
- role
    

to downstream handlers.

---

## Authorization Middleware

### Responsibilities

- Validate user roles
    
- Restrict protected routes
    
- Enforce RBAC policies
    
- Reject forbidden access attempts
    

---

# Security Requirements

Security is the primary concern of the `auth_service`.

## Password Security

- bcrypt hashing only
    
- Secure hashing cost factor
    
- No plain-text password storage
    
- No password logging
    

## Token Security

- Strong JWT signing secret
    
- Secure environment variable storage
    
- Signature validation
    
- Expiration enforcement
    

## API Security

- Input validation
    
- Request sanitization
    
- Generic authentication errors
    
- Prevention of credential enumeration
    

---

# Inter-Service Authentication Strategy

Other services depend on JWT tokens issued by the `auth_service`.

The platform uses local JWT validation within each service using the shared signing secret.

### Advantages

- Reduced network overhead
    
- Improved scalability
    
- Lower authentication latency
    
- Stateless request processing
    

This approach aligns with production-oriented distributed architectures.

---

# Scalability Characteristics

The `auth_service` is intentionally designed as a stateless service.

### Design Characteristics

- No session storage
    
- No in-memory authentication state
    
- Horizontal scalability support
    
- Kubernetes-friendly deployment
    

JWT tokens act as the portable authentication context across requests.

---

# Reliability Requirements

Since the `auth_service` acts as the central authentication authority, it becomes a critical system dependency.

If the service becomes unavailable:

- Users cannot authenticate
    
- JWT tokens cannot be issued
    
- Protected APIs become inaccessible
    

### Operational Requirements

- Health check endpoints
    
- Structured logging
    
- Error monitoring
    
- Reliable database connectivity
    

---

# Recommended Project Structure

The service follows Clean / Hexagonal Architecture principles.

```text
/auth_service
    /cmd
    /internal
        /handler
        /service
        /repository
        /model
        /middleware
        /utils
        /config
    /migrations
    /docs
```

---

# Future Extensibility

The architecture is designed to support future enhancements without major redesign.

Potential future improvements include:

- Refresh token rotation
    
- OAuth integration
    
- Multi-factor authentication
    
- API gateway integration
    
- Token revocation mechanisms
    
- Service-to-service authentication
    

---
