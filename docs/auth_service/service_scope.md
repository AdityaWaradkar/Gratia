# **Auth Service – Service Scope**

---

## **1. Purpose**

The **Auth Service** is responsible for **secure authentication and authorization** across the Gratia platform.\
It ensures that only valid users (Donors, NGOs, Admins) can access the system, issues tokens for authenticated sessions, and enforces role-based access.

This service acts as the **identity provider** for all other microservices.

***


## **2. Responsibilities**

- **User Authentication**

  - Register new users with email & password (hashed).

  - Authenticate users during login.

  - Verify identity using JWT tokens.

- **Authorization**

  - Assign and manage roles (`Donor`, `NGO`, `Admin`).

  - Enforce role-based access to APIs in other services.

- **Token Management**

  - Issue access tokens (short-lived JWTs).

  - Issue refresh tokens (long-lived, stored securely).

  - Support token refresh flow.

  - Invalidate tokens on logout or account deletion.

- **Password Security**

  - Store passwords using strong hashing (bcrypt/argon2).

  - Support password reset via secure token.

- **Inter-Service Security**

  - Provide token verification endpoint for other microservices.

  - Allow internal services to validate JWT claims securely.

***


## **3. Out of Scope**

To keep the service **focused and clean**, the Auth Service will **not** handle:

- User profile details (name, address, NGO license, etc.) → handled by **User Service**.

- Food listing, claims, notifications → handled by their respective services.

- Analytics & reports.

- File uploads (e.g., NGO license proof).

***


## **4. Core Users / Roles**

- **Donor** → can register/login, list food (via Food Listing Service).

- **NGO** → can register/login, claim food (via Claim Service).

- **Admin** → can login, manage approvals & bans (via Admin Service).

***


## **5. APIs (High-Level)**

- `POST /auth/register` – Register a new user.

- `POST /auth/login` – Authenticate user, return access + refresh tokens.

- `POST /auth/refresh` – Generate new access token using refresh token.

- `POST /auth/logout` – Invalidate current tokens.

- `POST /auth/forgot-password` – Request password reset link/token.

- `POST /auth/reset-password` – Reset password using secure token.

- `GET /auth/validate` – Validate access token (used by other services).

***


## **6. Security Considerations**

- Use **JWT (RS256)** for signed tokens → avoids symmetric key leaks.

- Store **refresh tokens** in DB with expiry.

- Implement **rate limiting** on login & register to prevent abuse.

- Enforce **strong password policy**.

- Use **HTTPS only** in production.

- CORS configuration for frontend.

***


## **7. Data Model (Minimal)**

**Table: users**

- `id (UUID)` – primary key

- `email (unique)`

- `password_hash`

- `role (enum: DONOR, NGO, ADMIN)`

- `created_at`

- `updated_at`

**Table: refresh\_tokens**

- `id (UUID)`

- `user_id (FK)`

- `token (hashed)`

- `expires_at`

- `revoked (boolean)`

***


## **8. Dependencies**

- **Database**: PostgreSQL (preferred) or MySQL.

- **Cache (Optional, for scaling)**: Redis (for blacklisting tokens, login throttling).

- **Libraries**:

  - `jwt-go` (JWT handling)

  - `bcrypt` (password hashing)

  - `gin/fiber` (HTTP framework)

***


## **9. Observability**

- Logging: login attempts, failed logins, token refresh.

- Metrics: total registrations, active users, failed login rate.

- Health check endpoint: `/auth/health`.

***


## **10. Success Criteria**

The Auth Service is considered **complete** when:

1. Users can securely register, login, logout, and refresh tokens.

2. Password reset flow works end-to-end.

3. Other services can verify tokens via Auth Service.

4. All security measures (hashing, JWT signing, HTTPS, rate limiting) are in place.

5. It is containerized, tested, documented, and deployed.

---
