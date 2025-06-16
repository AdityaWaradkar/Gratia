# Functional Requirements Document

## Project: Gratia - Cloud-Native Food Donation Platform

---

### Overview

Gratia is a full-stack, microservices-based cloud-native platform built to connect restaurants and catering services with NGOs for the purpose of managing and redistributing excess food. This document outlines and finalizes all functional requirements for each microservice, elaborates on role-based access control (RBAC), security standards, and the notification workflows involved.

---

## Microservices Functional Requirements

### 1. Auth Service

**Purpose:** Central identity and access management for all user roles.

#### Functional Requirements:

- **User Registration**
  - Role-based registration: Restaurant, NGO, Admin.
  - Input validation for email, password, phone number, and location.
  - Secure password hashing (bcrypt or Argon2).
  - Email verification link dispatch post-registration.

- **Login & Token Management**
  - Credential validation with JWT issuance (Access Token & Refresh Token).
  - Refresh endpoint for rotating access tokens.
  - Support for OAuth2 (optional, for social login integration).

- **Password Management**
  - Forgot password flow with time-bound secure token sent via email.
  - Authenticated password update endpoint.

- **RBAC Enforcement**
  - Role embedded in JWT.
  - Middleware to validate JWT and extract role for access decisions.

- **Session Handling**
  - Token revocation on logout.
  - Blacklist support for compromised token detection.

---

### 2. Restaurant Service

**Purpose:** Allow restaurants to list, update, and track food donations.

#### Functional Requirements:

- **Create Donation**
  - Input: Title, Description, Quantity, Expiry Time, Pickup Location (Geo).
  - Triggers NGO notification based on location.

- **Manage Listings**
  - Update: status (open, claimed, picked up, expired).
  - Archive completed donations.

- **History & Analytics**
  - Past donation log with filters.
  - Dashboard metrics: total donations, most active NGOs, impact reports.

- **RBAC Enforcement**
  - Only authenticated users with role `restaurant` can create or update listings.

---

### 3. NGO Service

**Purpose:** Enable NGOs to discover, claim, and track food listings.

#### Functional Requirements:

- **Discover Donations**
  - Listings filtered by geolocation, quantity, expiry.
  - Distance-based sorting using Geo Search service.

- **Claim Donation**
  - Lock mechanism to ensure one NGO can claim a listing.
  - Triggers restaurant notification on claim.

- **Manage Claimed Donations**
  - Update donation status: en-route, received, or cancel.

- **Feedback Mechanism**
  - NGO can leave a review or feedback per donation.

- **RBAC Enforcement**
  - Access restricted to users with `ngo` role.

---

### 4. Admin Dashboard Service

**Purpose:** Platform-wide oversight and management.

#### Functional Requirements:

- **User Management**
  - View all registered users.
  - Suspend/reactivate accounts.

- **Content Moderation**
  - Review flagged listings.
  - Monitor suspicious behaviors.

- **Platform Analytics**
  - Donation volumes, NGO activity, heatmaps of food flow.

- **Audit Log Access**
  - Secure access to platform-wide audit records.

- **RBAC Enforcement**
  - Admin-only access; enforced via JWT and gateway.

---

### 5. Notification Service

**Purpose:** Orchestrates user alerts via email, SMS, or in-app delivery.

#### Functional Requirements:

- **Trigger Points**
  - New nearby donation listing (for NGOs).
  - Donation claimed (for Restaurants).
  - Password recovery token.
  - Admin alerts (suspension, warnings).

- **Multi-Channel Delivery**
  - Email via third-party provider (e.g., SendGrid).
  - SMS via provider (e.g., Twilio).
  - Internal push messages (future integration).

- **Retry & Failure Handling**
  - Retries failed deliveries.
  - Logs all notification attempts for auditing.

- **User Preferences**
  - Configurable channels and opt-in/opt-out.

---

### 6. Geo Search Service

**Purpose:** Provides location-based search functionality for donation discovery.

#### Functional Requirements:

- **Proximity Search**
  - Radius-based geo queries using coordinates.

- **Distance Computation**
  - Precise real-time calculation of distance between NGO and restaurant.

- **Address Resolution**
  - Reverse geocoding using Google Maps API.

- **Caching**
  - Frequently queried locations cached in Redis for speed.

---

### 7. Audit Log Service

**Purpose:** Provides immutable system-wide logging for traceability and compliance.

#### Functional Requirements:

- **Event Recording**
  - Logs authentication, donation events, status updates, and admin actions.

- **Structured Logs**
  - Categorized logs (INFO, WARN, ERROR) in JSON format.

- **Query Interface**
  - Search logs by actor, time range, or action type.

- **Integration**
  - Compatible with Loki and ELK stack.

---

### 8. Cleanup Service

**Purpose:** Manages periodic system hygiene and data integrity.

#### Functional Requirements:

- **Donation Auto-Expire**
  - Marks listings as expired once expiry time is passed.

- **Data Pruning**
  - Deletes orphaned or stale records (e.g., unclaimed listings).

- **Queue Cleanup**
  - Purges notification attempts after retry threshold.

- **Health Endpoint**
  - Exposes readiness and liveness for monitoring.

---

## Role-Based Access Control (RBAC)

| Role       | Permissions                                                     |
|------------|-----------------------------------------------------------------|
| Restaurant | Create/manage listings, view claim status, update pickup status |
| NGO        | View/claim listings, update delivery status, give feedback      |
| Admin      | Manage users, moderate content, access audit logs               |

- All access is mediated through JWT.
- Gateway intercepts and verifies JWT before passing to internal services.

---

## Security Requirements

- **Authentication**
  - JWTs signed with a secret key, expiring every 15 minutes.
  - Refresh tokens valid for 7 days, stored securely.

- **Password Handling**
  - All passwords hashed using bcrypt with salt.

- **Secrets Management**
  - All API keys and credentials managed via HashiCorp Vault.

- **Transport Security**
  - All traffic enforced over HTTPS with TLS 1.3.

- **Service-to-Service Communication**
  - Mutual TLS for inter-microservice communication.

---

## Notification Workflow

### Flow: New Donation Created

1. Restaurant creates a listing.
2. Geo Search identifies nearby NGOs.
3. Notification Service sends email/SMS to those NGOs.

### Flow: Donation Claimed

1. NGO claims donation.
2. Restaurant receives notification (email/SMS).

### Flow: Password Reset

1. User requests reset.
2. Token emailed via Notification Service.

### Flow: Account Suspended

1. Admin suspends user.
2. Notification dispatched with reason.

---

## Conclusion

This document serves as the finalized blueprint for Gratia's backend functional design. The microservices are scoped with clear responsibilities, strong security standards, and operational independence, aligned with DevOps and modern cloud-native principles.
