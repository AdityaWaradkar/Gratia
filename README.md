# Gratia

Gratia is a cloud-native backend platform designed to connect donors and NGOs for efficient surplus resource distribution.\
The project is built using a microservices architecture with a strong focus on clean design, security, and maintainability.

This repository represents the backend foundation of the Gratia platform.

***


## Overview

Gratia follows a clear separation of concerns:

- **Authentication and identity** are handled by a dedicated service.

- **Domain data** such as donor and NGO profiles are handled by a separate service.

- Services communicate implicitly through JWT claims rather than direct coupling.

This design ensures scalability, security, and long-term maintainability.

***


## Services

### Authentication Service (`auth_service`)

Responsible for identity and access management.

**Responsibilities**

- User registration and login

- JWT access and refresh token generation

- Token refresh and logout

- Password reset flow

- Authentication middleware

**Not responsible for**

- Donor or NGO profile data

- Business domain logic

***


### User Service (`user_service`)

Responsible for domain-specific user data.

**Responsibilities**

- Donor profile creation and updates

- NGO profile creation

- Admin verification of NGO profiles

- Domain validation and authorization

**Not responsible for**

- Authentication

- Passwords

- Token issuance

***


## Architecture

    Client (Web / Mobile)
            |
            v
    Authentication Service
    (JWT Issuance & Validation)
            |
            v
    User Service
    (Donor & NGO Profiles)

- JWT tokens issued by the authentication service are trusted by downstream services.

- Each service owns its database tables and logic.

- Cross-service linkage is done using `user_id` from JWT claims.

***


## Database Design

- `users` table is owned by the authentication service.

- `donor_profiles` and `ngo_profiles` tables are owned by the user service.

- No service directly modifies another service’s data.

This avoids tight coupling and accidental data corruption.

***


## Technology Stack

- Language: Go

- Database: PostgreSQL

- Authentication: JWT (Access and Refresh tokens)

- Containerization: Docker

- Architecture: Microservices

- API Style: REST

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

Each service follows a layered structure:

- Handler

- Service

- Repository

***


## Configuration

Each service is configured via environment variables.

Example:

    PORT=8080
    DATABASE_URL=postgresql://...
    JWT_SECRET=your_secret

Secrets are never committed to version control.

***


## API Testing

All APIs are tested using Postman with a shared environment.

Recommended environment variables:

    AUTH_BASE_URL

    USER_BASE_URL

    ACCESS_TOKEN

    REFRESH_TOKEN

    ADMIN_ACCESS_TOKEN

Tokens issued by the authentication service are reused for user service requests.

***


## Development Status

- Authentication Service: Complete

- User Service: Complete

- Future planned services:

  - Donation management

  - Matching and allocation

  - Notifications and events

***


## Author

Aditya Waradkar\
Backend Engineer (Go, PostgreSQL, Microservices)

***


## License

This project is licensed under the MIT License.

***
