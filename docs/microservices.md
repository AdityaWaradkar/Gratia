# Gratia – Microservices Architecture

## Overview

Gratia is a distributed backend system designed to reduce food waste by connecting donors with NGOs.\
The system follows a **domain-driven microservices architecture**, where each service owns a clearly defined responsibility and data model.

The architecture prioritizes:

- Clear service boundaries

- Minimal coupling

- Strong authorization guarantees

- Incremental scalability (no premature complexity)

***


## Core Design Principles

- **Authentication is centralized**

- **Business domains are isolated**

- **Authorization is enforced at service boundaries**

- **Each service owns its database**

- **Services communicate via JWT claims, not shared state**

***


## Implemented Services

### 1. Auth Service

**Purpose**\
Identity and access control for the entire platform.

**Responsibilities**

- User registration

- User login

- Password reset

- JWT issuance and validation

- Global role management

**Roles Managed**

    USER

    ADMIN

**Key Notes**

- JWT contains `sub` (user\_id), `role`, and email

- Acts as the single source of truth for authentication

- Other services trust JWTs issued by Auth Service

**Out of Scope**

- User profiles

- Donor / NGO business logic

- Domain-specific data

***


### 2. User Service

**Purpose**\
Manage domain-level user data and profiles.

**Responsibilities**

- Donor profile management

- NGO profile management

- NGO verification workflow

- Profile updates and retrieval

**Domain Roles**

- Donor

- NGO

**Authorization Model**

- Auth role (`USER`, `ADMIN`) is read from JWT

- Domain role (Donor / NGO) is derived from profile existence and state

**Admin Capabilities**

- Verify NGO profiles

- Approve or reject NGO legitimacy

**Key Notes**

- User Service does not issue or validate tokens

- Relies entirely on Auth Service JWTs

- Maintains separation between identity and domain data

***


## Planned Services

### 3. Food Service

**Purpose**\
Manage food availability posted by donors.

**Responsibilities**

- Create food listings

- Update or delete listings

- Track food metadata:

  - Type

  - Quantity

  - Expiry

  - Location

- Maintain food status:

  - AVAILABLE

  - CLAIMED

  - EXPIRED

**Authorization**

- Only donors can create food listings

- Admins may moderate listings (future)

**Out of Scope**

- Claim lifecycle

- NGO interactions

- Notifications

***


### 4. Claim Service

**Purpose**\
Handle the lifecycle of food claims by NGOs.

**Responsibilities**

- NGOs request food claims

- Manage claim state transitions:

  - REQUESTED

  - APPROVED

  - PICKED\_UP

  - DELIVERED

  - CANCELLED

- Enforce business rules:

  - One active claim per food item

  - Only verified NGOs can claim food

**Authorization**

- Only NGOs can create claims

- Donors can approve or reject claims

- Admins can intervene if required

***


## Deferred / Future Services

### Notification Handling (Deferred)

- Notifications are currently handled inline within services

- Designed to be abstracted later if scale requires

- Will be extracted only when:

  - Multiple notification channels are introduced

  - Asynchronous processing becomes necessary

***


### Analytics (Deferred)

- No dedicated analytics service initially

- Reporting can be generated via database queries or logs

- A separate analytics service may be introduced later if needed

***


### API Gateway (Deferred)

- Not required in the current phase

- Services are accessed directly

- Gateway may be introduced later for:

  - Rate limiting

  - Centralized request routing

  - External client abstraction

***


## Inter-Service Communication

- JWT is the primary trust mechanism

- Services do not call Auth Service on every request

- Authorization decisions are made locally using JWT claims

- No shared databases or shared schemas

***


## Final Architecture Summary

    Auth Service
      └── Identity, JWT, Global Roles

    User Service
      ├── Donor Profiles
      ├── NGO Profiles
      └── NGO Verification (Admin)

    Food Service
      └── Food Inventory and Availability

    Claim Service
      └── Claim Lifecycle and Rules

***


## Architectural Status

- Auth Service: Implemented

- User Service: Implemented

- Food Service: Next planned

- Claim Service: Planned

- Notifications, Analytics, Gateway: Deferred by design

***
