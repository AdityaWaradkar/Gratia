# 1. Overview

The database design of the `auth_service` is focused on securely storing authentication-related identity data required for user authentication and authorization.

The service follows a database-per-service architecture, where the `auth_service` owns and manages its own isolated database schema independently from other services in the platform.

The database is intentionally designed to remain:

- Lightweight
    
- Security-focused
    
- Authentication-oriented
    
- Independent from business-domain data
    

The database does not store application-specific profile information or business workflow data.

---

# 2. Purpose of Database Design

The primary purpose of the database layer is to:

- Store user authentication data
    
- Maintain identity records
    
- Support secure login workflows
    
- Persist authorization roles
    
- Enable JWT-based authentication
    

The database acts as the persistent identity store for the platform.

---

# 3. Database Ownership

The `auth_service` exclusively owns:

- User authentication records
    
- Password hashes
    
- Authorization roles
    
- Authentication-related metadata
    

Other services are not allowed to directly access this database.

All inter-service communication must occur through APIs.

This preserves:

- Service isolation
    
- Clear ownership boundaries
    
- Loose coupling
    
- Better maintainability
    

---

# 4. Database Architecture

The `auth_service` uses PostgreSQL as the primary relational database.

The schema is intentionally minimal to reduce:

- Complexity
    
- Security exposure
    
- Unnecessary coupling
    

The database focuses strictly on authentication concerns.

---

# 5. Primary Database Tables

The initial database design contains the following primary table:

|Table|Purpose|
|---|---|
|users|Stores authentication identity records|

Future extensions may introduce additional tables such as:

- refresh_tokens
    
- audit_logs
    

However, the initial implementation remains intentionally simple.

---

# 6. users Table

The `users` table acts as the core authentication table of the service.

It stores:

- User identity
    
- Login credentials
    
- Authorization role
    

---

## Table Structure

|Column|Type|Description|
|---|---|---|
|user_id|UUID|Primary user identifier|
|email|VARCHAR|Unique login identifier|
|password_hash|VARCHAR|bcrypt hashed password|
|role|ENUM|User authorization role|
|created_at|TIMESTAMP|Record creation timestamp|
|updated_at|TIMESTAMP|Record update timestamp|

---

# 7. Column Responsibilities

---

## 7.1 user_id

The `user_id` acts as the globally unique identity of the user.

### Characteristics

- UUID-based identifier
    
- Primary key
    
- Immutable
    
- Used across services for identity reference
    

Using UUIDs improves:

- Distributed system compatibility
    
- Security
    
- Scalability
    

---

## 7.2 email

The `email` column acts as the unique login identifier.

### Characteristics

- Unique per user
    
- Required during authentication
    
- Used for credential lookup
    

### Constraints

- Must be unique
    
- Cannot be null
    

---

## 7.3 password_hash

The `password_hash` stores the bcrypt-generated password hash.

### Characteristics

- Never stores plain-text passwords
    
- Used during login verification
    
- Security-critical field
    

### Security Rules

- Plain-text passwords are never persisted
    
- Hashes must never be exposed externally
    

---

## 7.4 role

The `role` column stores authorization roles used for RBAC.

### Supported Values

- DONOR
    
- NGO
    
- ADMIN
    

### Purpose

- Permission enforcement
    
- Route authorization
    
- Access control validation
    

---

## 7.5 created_at

Stores the timestamp when the user record was created.

Used for:

- Auditing
    
- Record tracking
    
- Administrative visibility
    

---

## 7.6 updated_at

Stores the timestamp of the latest record update.

Used for:

- Change tracking
    
- Administrative monitoring
    

---

# 8. Database Constraints

The schema includes several constraints to maintain consistency and integrity.

|Constraint|Purpose|
|---|---|
|PRIMARY KEY(user_id)|Ensure entity uniqueness|
|UNIQUE(email)|Prevent duplicate accounts|
|NOT NULL(password_hash)|Maintain authentication integrity|
|NOT NULL(role)|Maintain authorization integrity|

These constraints help prevent invalid authentication records.

---

# 9. ENUM Design for Roles

The `role` field is implemented using an ENUM type.

### Allowed Values

```text
DONOR
NGO
ADMIN
```

Using ENUMs ensures:

- Controlled role values
    
- Better data consistency
    
- Reduced invalid authorization states
    

---

# 10. Database Security Considerations

The authentication database contains highly sensitive information.

Security measures include:

- bcrypt password hashing
    
- No plain-text password storage
    
- Restricted database access
    
- Environment-based credentials
    
- Isolated service ownership
    

The database should never expose:

- Password hashes
    
- JWT secrets
    
- Internal authentication metadata
    

---

# 11. Database Isolation Strategy

The platform follows a strict database-per-service architecture.

This means:

- `auth_service` owns its database exclusively
    
- Other services cannot directly query authentication tables
    
- Inter-service access occurs only through APIs
    

Advantages include:

- Loose coupling
    
- Better scalability
    
- Independent schema evolution
    
- Improved maintainability
    

---

# 12. Migration Strategy

Database schema changes are managed using version-controlled SQL migrations.

Migration responsibilities include:

- Table creation
    
- Constraint definition
    
- ENUM creation
    
- Initial admin bootstrap
    
- Schema evolution
    

Example migration operations:

- Create `users` table
    
- Insert default ADMIN account
    
- Add indexes and constraints
    

---

# 13. Initial Admin Bootstrapping

The database initialization process includes creation of a default administrative account.

This ensures:

- Platform operability after deployment
    
- Initial administrative access
    
- Elimination of manual database modification
    

The default admin password must be securely hashed before insertion.

---

# 14. Scalability Considerations

The database design intentionally remains lightweight.

Advantages:

- Faster authentication lookups
    
- Reduced query complexity
    
- Better scaling behavior
    
- Lower operational overhead
    

Authentication databases should prioritize:

- Fast reads
    
- High reliability
    
- Minimal complexity
    

---

# 15. Future Database Extensions

The schema can later be extended with:

- refresh_tokens table
    
- audit logging
    
- login attempt tracking
    
- token revocation tracking
    
- MFA-related metadata
    

The current design focuses on a clean foundational authentication schema.

---

# 16. Recommended Database Structure

```text
auth_service_db
    └── users
```

Potential future structure:

```text
auth_service_db
    ├── users
    ├── refresh_tokens
    └── auth_audit_logs
```

---
