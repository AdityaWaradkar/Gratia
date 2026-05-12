# 1. Overview

Role-Based Access Control (RBAC) is the authorization mechanism used in the Gratia platform to regulate access to protected resources and operations.

After a user is successfully authenticated, RBAC determines what actions the user is permitted to perform within the system.

The authorization system is implemented inside the `auth_service` and enforced through middleware-based route protection across backend services.

RBAC enables the platform to maintain controlled access boundaries between different categories of users while ensuring that sensitive operations remain restricted to authorized roles only.

---

# 2. Purpose of RBAC

The primary purpose of RBAC is to enforce authorization policies based on predefined user roles.

RBAC answers the question:

> “What is the authenticated user allowed to access?”

The authorization layer ensures that:

- Users only access permitted resources
    
- Sensitive operations remain protected
    
- Administrative operations remain restricted
    
- Business workflows follow controlled permissions
    

---

# 3. RBAC Architecture

The RBAC system is built around:

- Role assignment
    
- Role validation
    
- Permission enforcement
    
- Middleware-level authorization
    

Each authenticated user is assigned a role during account creation.

That role becomes part of the JWT token claims and is later used during authorization checks.

---

# 4. Supported Roles

The platform currently defines three primary roles.

|Role|Description|
|---|---|
|DONOR|User responsible for donating food|
|NGO|Verified organization claiming food donations|
|ADMIN|Administrative authority managing the platform|

Each role has clearly restricted access boundaries.

---

# 5. Role Responsibilities

## 5.1 DONOR Role

The `DONOR` role represents users who create and manage food donation listings.

### Allowed Operations

- Create food listings
    
- Update food listings
    
- Delete food listings
    
- View owned listings
    
- Manage donation activity
    

### Restricted Operations

- NGO verification
    
- Administrative management
    
- Claim approval workflows belonging to other users
    

---

## 5.2 NGO Role

The `NGO` role represents verified organizations participating in food claim operations.

### Allowed Operations

- Discover food listings
    
- Create food claims
    
- Track claim lifecycle
    
- Manage NGO-related claim activity
    

### Restricted Operations

- Administrative operations
    
- NGO verification approval
    
- Unauthorized food management
    

Only verified NGOs should be allowed to perform claim-related operations.

---

## 5.3 ADMIN Role

The `ADMIN` role represents system administrators.

### Allowed Operations

- Verify NGOs
    
- Manage platform operations
    
- Access administrative endpoints
    
- Monitor system activity
    
- Perform privileged operations
    

The `ADMIN` role has the highest level of authorization within the platform.

---

# 6. RBAC Workflow

The authorization process follows this sequence:

```text
Authenticated Request
          │
          ▼
JWT Token Validation
          │
          ▼
Extract User Role
          │
          ▼
Authorization Middleware
          │
          ├── Role Allowed → Continue Request
          │
          └── Role Denied → Reject Request
```

---

# 7. Role Assignment

Roles are assigned during user registration.

Example:

- Donor accounts receive `DONOR`
    
- NGO accounts receive `NGO`
    
- Administrative accounts receive `ADMIN`
    

The assigned role becomes part of the user's authentication identity.

---

# 8. JWT Role Claims

The user's role is embedded directly inside the JWT payload.

Example:

```json
{
  "user_id": "uuid",
  "role": "DONOR",
  "exp": 1715000000
}
```

This allows services to authorize requests without requiring additional database lookups.

---

# 9. Authorization Middleware

Authorization enforcement is performed using middleware.

The middleware is responsible for:

- Extracting user role from JWT claims
    
- Validating role permissions
    
- Restricting unauthorized access
    
- Rejecting forbidden requests
    

This ensures consistent authorization enforcement across all services.

---

# 10. Route Protection Strategy

Protected routes define which roles are permitted to access specific operations.

---

## Example Access Restrictions

|Endpoint|Allowed Roles|
|---|---|
|/food/create|DONOR|
|/claim/create|NGO|
|/admin/*|ADMIN|

Routes without proper authorization are rejected.

---

# 11. Authorization Rules

The RBAC system follows several core authorization rules.

|Rule|Description|
|---|---|
|Authentication required before authorization|Anonymous requests are denied|
|Every user must have a role|Authorization depends on roles|
|Unauthorized access must be rejected|Prevent privilege escalation|
|Roles define access scope|Permissions remain isolated|
|Administrative routes remain protected|Critical operations stay secure|

---

# 12. Forbidden Access Handling

When a user attempts to access a restricted resource:

- Authorization fails
    
- Request is rejected
    
- `403 Forbidden` response is returned
    

Example:

```json
{
  "error": "Forbidden"
}
```

The system avoids exposing internal authorization details.

---

# 13. Security Considerations

RBAC is a critical part of the platform security architecture.

The authorization system helps prevent:

- Unauthorized resource access
    
- Privilege escalation
    
- Administrative misuse
    
- Cross-role access violations
    

Authorization checks are enforced consistently at middleware level to reduce implementation errors.

---

# 14. Advantages of RBAC

The RBAC approach provides several architectural advantages.

---

## Simplified Permission Management

Permissions are managed through roles instead of individual users.

---

## Improved Maintainability

Authorization logic remains centralized and reusable.

---

## Consistent Access Control

Middleware-based enforcement ensures consistent authorization across services.

---

## Better Security Isolation

Each role operates within clearly defined permission boundaries.

---

# 15. Scalability Characteristics

RBAC integrates efficiently with stateless JWT authentication.

Since user roles are embedded inside JWT claims:

- No repeated database lookups are required
    
- Authorization remains lightweight
    
- Services can validate permissions independently
    

This supports horizontally scalable distributed systems.

---

# 16. Future Extensibility

The RBAC system can later be extended to support:

- Fine-grained permissions
    
- Permission groups
    
- Dynamic policy management
    
- Resource-level authorization
    
- Multi-role users
    
- Hierarchical roles
    

The current implementation is intentionally designed to remain simple and maintainable during the initial platform phase.

---

