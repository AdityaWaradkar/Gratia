# 1. Overview

The Gratia platform uses JWT-based stateless authentication as the primary mechanism for maintaining authenticated user identity across distributed microservices.

JWT (JSON Web Token) enables the platform to authenticate requests without maintaining server-side session state. Instead of storing authentication sessions in memory or databases, authenticated user identity is embedded directly within signed tokens that are transmitted with every protected request.

This approach is intentionally chosen to support:

- Distributed system architecture
    
- Horizontal scalability
    
- Service independence
    
- Stateless backend services
    

The `auth_service` acts as the centralized authority responsible for generating and validating JWT tokens within the platform.

---

# 2. Purpose of JWT-Based Authentication

The primary purpose of JWT authentication is to provide a portable and secure identity representation that can be verified independently by backend services.

JWT authentication allows:

- Stateless authentication
    
- Secure identity propagation
    
- Reduced dependency on centralized session storage
    
- Efficient inter-service trust establishment
    

This mechanism becomes critical in microservices architectures where multiple services must securely validate user identity.

---

# 3. Stateless Authentication Concept

In traditional session-based authentication systems:

- Session data is stored on the server
    
- Every request requires session lookup
    
- Scaling session management becomes complex
    

In the Gratia platform:

- No server-side authentication sessions are stored
    
- Authentication state is carried inside JWT tokens
    
- Services independently validate tokens
    

This creates a stateless authentication model.

---

# 4. JWT Structure

A JWT token consists of three parts:

```text
HEADER.PAYLOAD.SIGNATURE
```

Each section is Base64URL encoded.

---

## 4.1 Header

The header defines token metadata.

Example:

```json
{
  "alg": "HS256",
  "typ": "JWT"
}
```

### Fields

|Field|Description|
|---|---|
|alg|Signing algorithm|
|typ|Token type|

The platform uses the `HS256` signing algorithm.

---

## 4.2 Payload

The payload contains identity claims.

Example:

```json
{
  "user_id": "uuid",
  "role": "DONOR",
  "exp": 1715000000,
  "iat": 1714996400
}
```

---

### JWT Claims

|Claim|Description|
|---|---|
|user_id|Unique authenticated user identifier|
|role|Authorization role|
|exp|Token expiration timestamp|
|iat|Token issued timestamp|

These claims provide enough identity information for downstream services to authorize requests.

---

## 4.3 Signature

The signature ensures token integrity.

The signature is generated using:

- Header
    
- Payload
    
- Secret signing key
    

If the payload is modified:

- Signature validation fails
    
- Token becomes invalid
    

This prevents token tampering.

---

# 5. JWT Authentication Workflow

The JWT authentication process follows this sequence:

```text
User Login Request
        │
        ▼
auth_service
        │
        ├── Validate credentials
        ├── Generate JWT token
        └── Return token
                │
                ▼
Client Stores JWT
                │
                ▼
Authenticated Requests
                │
                ▼
Backend Services Validate JWT
```

---

# 6. Token Generation Process

After successful authentication:

1. User identity is verified
    
2. JWT claims are created
    
3. Expiration timestamp is assigned
    
4. Token is signed using secret key
    
5. Signed JWT is returned to the client
    

The token then acts as the user's portable authentication identity.

---

# 7. Token Transmission

JWT tokens are transmitted using the `Authorization` header.

Example:

```http
Authorization: Bearer <jwt_token>
```

Every protected request must include a valid JWT token.

---

# 8. Token Validation

When a protected request reaches a service:

1. JWT token is extracted
    
2. Signature is verified
    
3. Expiration is checked
    
4. Claims are parsed
    
5. User identity is reconstructed
    

If validation fails:

- Request is rejected
    
- Access is denied
    

---

# 9. Token Expiration Strategy

JWT tokens are intentionally short-lived.

### Recommended Expiration

- 15 minutes to 60 minutes
    

### Purpose

- Reduce risk of token theft
    
- Limit exposure window
    
- Improve authentication security
    

Expired tokens are considered invalid and require re-authentication.

---

# 10. Secret Key Management

JWT security depends heavily on secure signing key management.

### Requirements

- Secrets must never be hardcoded
    
- Secrets must be stored using environment variables
    
- Secrets must remain private
    
- Secrets should differ across environments
    

Compromised signing keys compromise the entire authentication system.

---

# 11. Advantages of JWT Stateless Authentication

The JWT approach provides several architectural advantages.

---

## 11.1 Stateless Architecture

No session storage is required.

This simplifies:

- Scaling
    
- Load balancing
    
- Distributed deployments
    

---

## 11.2 Horizontal Scalability

Any service instance can validate tokens independently.

This enables:

- Kubernetes scaling
    
- Container orchestration
    
- Multi-instance deployments
    

---

## 11.3 Reduced Database Dependency

Authenticated requests do not require repeated database lookups.

This reduces:

- Database load
    
- Authentication latency
    

---

## 11.4 Inter-Service Trust

Services can independently trust authenticated requests by validating JWT signatures.

This is essential in distributed systems.

---

# 12. Security Considerations

JWT authentication introduces several security requirements.

---

## Token Integrity

- Tokens must be cryptographically signed
    
- Invalid signatures must be rejected
    

---

## Token Expiration

- Expired tokens must be denied
    
- Long-lived tokens should be avoided
    

---

## Secret Protection

- Signing secrets must remain confidential
    
- Secret exposure compromises authentication
    

---

## HTTPS Requirement

JWT tokens should only be transmitted over secure HTTPS connections to prevent interception.

---

# 13. Limitations of JWT Authentication

Although JWT provides strong scalability advantages, it also introduces limitations.

---

## No Immediate Revocation

Since tokens are stateless:

- Revoking active tokens becomes difficult
    
- Tokens remain valid until expiration
    

---

## Token Theft Risk

If a token is stolen:

- It can be reused until expiration
    

This is why short-lived tokens are important.

---

# 14. Why JWT Was Chosen for Gratia

JWT-based authentication aligns strongly with the architectural goals of the platform.

The platform requires:

- Distributed services
    
- Stateless communication
    
- Horizontal scalability
    
- Kubernetes compatibility
    
- Independent service validation
    

JWT authentication supports all of these requirements efficiently.

Traditional session-based authentication would introduce:

- Centralized session storage
    
- Additional infrastructure complexity
    
- Scaling challenges
    

For a microservices architecture, JWT provides a cleaner and more scalable solution.

---

# 15. Future Enhancements

The JWT authentication system can later be extended with:

- Refresh token support
    
- Token revocation mechanisms
    
- Rotating signing keys
    
- Multi-factor authentication
    
- OAuth integration
    

The current architecture is intentionally designed to support these future enhancements without major redesign.

---
