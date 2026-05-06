# 1. Overview

The middleware architecture in the `auth_service` is responsible for enforcing authentication and authorization rules across protected APIs.

Middleware acts as an intermediate processing layer between incoming HTTP requests and backend route handlers. It is used to validate requests before they reach application business logic.

In the Gratia platform, middleware is primarily responsible for:

- JWT authentication
    
- Authorization enforcement
    
- Request validation
    
- User identity extraction
    
- Protected route handling
    

The middleware layer is designed to provide centralized, reusable, and consistent security enforcement across all backend services.

---

# 2. Purpose of Middleware

The primary purpose of middleware is to separate security and request-processing concerns from application business logic.

Instead of implementing authentication and authorization checks inside every route handler, middleware performs these operations centrally before requests are processed.

This approach improves:

- Code reusability
    
- Maintainability
    
- Security consistency
    
- Service modularity
    

---

# 3. Middleware Flow

The request lifecycle follows this sequence:

```text
Incoming Request
        │
        ▼
Authentication Middleware
        │
        ▼
Authorization Middleware
        │
        ▼
Route Handler
        │
        ▼
Response
```

If middleware validation fails at any stage:

- Request processing stops
    
- Error response is returned
    
- Route handler is never executed
    

---

# 4. Authentication Middleware

Authentication middleware is responsible for validating JWT tokens and establishing authenticated user identity.

---

## Responsibilities

The authentication middleware performs the following operations:

- Extract JWT token from request headers
    
- Validate token signature
    
- Verify token expiration
    
- Parse JWT claims
    
- Reconstruct authenticated user identity
    
- Reject unauthorized requests
    

---

## Authentication Workflow

The authentication middleware follows this sequence:

```text
Request Received
        │
        ▼
Extract Authorization Header
        │
        ▼
Extract JWT Token
        │
        ▼
Validate Signature
        │
        ▼
Validate Expiration
        │
        ▼
Parse Claims
        │
        ▼
Inject User Context
```

If any validation step fails:

- Authentication fails
    
- Request is rejected
    
- `401 Unauthorized` response is returned
    

---

## Authorization Header Format

JWT tokens are transmitted using the `Authorization` header.

Example:

```http
Authorization: Bearer <jwt_token>
```

Requests without valid tokens are denied access to protected endpoints.

---

# 5. Request Context Injection

After successful authentication, middleware injects authenticated user information into the request context.

Typical context values include:

- user_id
    
- role
    

This allows downstream handlers and services to access authenticated user identity without re-validating the token.

---

# 6. Authorization Middleware

Authorization middleware is responsible for enforcing Role-Based Access Control (RBAC).

This middleware operates only after successful authentication.

---

## Responsibilities

The authorization middleware performs:

- Role extraction
    
- Permission validation
    
- Route access enforcement
    
- Forbidden request rejection
    

---

## Authorization Workflow

```text
Authenticated Request
        │
        ▼
Extract User Role
        │
        ▼
Compare Against Allowed Roles
        │
        ├── Allowed → Continue Request
        │
        └── Denied → Reject Request
```

If role validation fails:

- Request is rejected
    
- `403 Forbidden` response is returned
    

---

# 7. Protected Route Handling

Middleware is applied to protected endpoints that require authenticated access.

Example protected routes:

|Route|Required Role|
|---|---|
|/food/create|DONOR|
|/claim/create|NGO|
|/admin/*|ADMIN|

Public routes such as login and signup bypass authentication middleware.

---

# 8. Middleware Layer Responsibilities

The middleware layer is responsible only for:

- Security enforcement
    
- Request validation
    
- User identity extraction
    
- Access restriction
    

It does not contain:

- Business logic
    
- Database business operations
    
- Domain workflows
    
- Service orchestration
    

This separation maintains clean architectural boundaries.

---

# 9. Advantages of Middleware Architecture

The middleware-based design provides several architectural advantages.

---

## Centralized Security Enforcement

Authentication and authorization logic remains centralized instead of duplicated across handlers.

---

## Reusability

Middleware can be reused across:

- Multiple routes
    
- Multiple services
    
- Different deployment environments
    

---

## Improved Maintainability

Security changes can be implemented in one location instead of modifying every endpoint.

---

## Consistent Request Validation

All protected requests follow identical authentication and authorization validation rules.

---

# 10. Error Handling

Middleware is responsible for handling invalid or unauthorized requests consistently.

---

## Authentication Failure

Authentication failures return:

```json
{
  "error": "Unauthorized"
}
```

HTTP Status:

```text
401 Unauthorized
```

---

## Authorization Failure

Authorization failures return:

```json
{
  "error": "Forbidden"
}
```

HTTP Status:

```text
403 Forbidden
```

---

# 11. Security Considerations

The middleware layer is a critical security boundary within the platform.

The middleware system helps prevent:

- Unauthorized API access
    
- Invalid token usage
    
- Expired token access
    
- Privilege escalation
    
- Cross-role access violations
    

Since all protected requests pass through middleware, it becomes one of the most security-sensitive layers in the architecture.

---

# 12. Scalability Characteristics

The middleware architecture is designed to support stateless distributed systems.

Because JWT validation is performed locally:

- No centralized session storage is required
    
- Services remain horizontally scalable
    
- Authentication latency remains low
    

This aligns well with:

- Kubernetes deployments
    
- Container orchestration
    
- Microservices architectures
    

---

# 13. Reusability Across Services

The middleware architecture is intentionally designed to be reusable across all backend services.

Services such as:

- `user_service`
    
- `food_service`
    
- `claim_service`
    

can share the same authentication and authorization middleware logic.

This ensures:

- Consistent security enforcement
    
- Standardized authentication behavior
    
- Reduced implementation duplication
    

---

# 14. Recommended Middleware Structure

```text
/middleware
    auth_middleware.go
    rbac_middleware.go
```

---

# 15. Future Enhancements

The middleware layer can later support:

- Rate limiting
    
- Request tracing
    
- API key validation
    
- OAuth integration
    
- Permission-based authorization
    
- Distributed tracing
    

The current design focuses on secure foundational authentication and authorization enforcement.

---
