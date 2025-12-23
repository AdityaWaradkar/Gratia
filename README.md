# Gratia

Gratia is a cloud-native backend platform designed to connect donors and NGOs for efficient and accountable redistribution of surplus resources.\
The system is built using a **microservices architecture**, with a strong emphasis on **clean boundaries, security, and long-term maintainability**.

This repository contains the **backend foundation** of the Gratia platform.

***


## Overview

Gratia follows a strict separation of responsibilities across services:

- **Authentication and identity** are handled by a dedicated authentication service.

- **Domain-specific user data** (donors and NGOs) are managed by a separate user service.

- **Business services** validate users via service-to-service communication, not shared databases.

This architecture avoids tight coupling, enables independent scaling, and keeps the system easy to reason about as it grows.

***


## Services

### Authentication Service (`auth_service`)

Responsible for identity and access management.

**Responsibilities**

- User registration and login

- JWT access and refresh token issuance

- Token refresh and logout

- Password reset workflows

- Authentication middleware

**Out of scope**

- Donor or NGO profile data

- Business or domain logic

***


### User Service (`user_service`)

Responsible for domain-specific user information.

**Responsibilities**

- Donor profile creation and updates

- NGO profile creation

- Admin verification of NGO profiles

- Exposing internal APIs for other services to validate donor/NGO status

**Out of scope**

- Authentication

- Password management

- Token issuance

***


### Food Service (`food_service`)

Responsible for managing food listings created by donors.

**Responsibilities**

- Create food listings (donor-only)

- List available (open) food listings

- Retrieve individual listings

- Update listings by the owning donor

- Enforce domain rules such as expiry and status transitions

**Design principles**

- Never reads user data directly from another service’s database

- Validates donor identity via internal APIs exposed by `user_service`

***


## Architecture

    Client (Web / Mobile)
            |
            v
    Authentication Service
    (JWT Issuance)
            |
            v
    Downstream Services
    (User Service, Food Service, ...)

- JWT tokens issued by the authentication service are trusted by downstream services.

- Each service owns its **own database schema**.

- Cross-service linkage is done using `user_id` from JWT claims.

- Role-specific validation (donor / NGO) is handled by the user service via internal APIs.

***


## Database Ownership

- `users` table → owned by **auth\_service**

- `donor_profiles`, `ngo_profiles` → owned by **user\_service**

- `food_listings` → owned by **food\_service**

No service directly modifies another service’s data.\
This prevents accidental coupling and data corruption.

***


## Technology Stack

- Language: Go

- Database: PostgreSQL

- Authentication: JWT (Access & Refresh Tokens)

- Architecture: Microservices

- API Style: REST

- Containerization: Docker

***


## Project Structure

    auth_service/
      ├── cmd/
      ├── internal/
      │   ├── auth
      │   ├── middleware
      │   ├── db
      │   └── config
      └── api/

    user_service/
      ├── cmd/
      ├── internal/
      │   ├── user
      │   ├── middleware
      │   ├── db
      │   └── config
      └── api/

    food_service/
      ├── cmd/
      ├── internal/
      │   ├── food
      │   ├── middleware
      │   ├── db
      │   └── config
      └── api/

Each service follows a consistent layered architecture:

- Handler (HTTP layer)

- Service (business logic)

- Repository (data access)

***


## Configuration

Each service is configured using environment variables.

Example:

    PORT=8080
    DATABASE_URL=postgresql://...
    JWT_SECRET=your_secret
    USER_SERVICE_URL=http://user_service:8081

Secrets are never committed to version control.

***


## API Testing

All APIs are tested using Postman with a shared environment.

Recommended variables:

    AUTH_BASE_URL
    USER_BASE_URL
    FOOD_BASE_URL
    ACCESS_TOKEN
    REFRESH_TOKEN
    ADMIN_ACCESS_TOKEN

Tokens issued by the authentication service are reused across services.

***


## Development Status

- Authentication Service: Complete

- User Service: Complete

- Food Service: Complete

Planned services:

- Claim Service

- Notification Service

- Admin Service

- API Gateway

- Analytics and reporting

***


## Author

Aditya Waradkar\
Backend Engineer (Go, PostgreSQL, Microservices)

***


## License

This project is licensed under the MIT License.

***
