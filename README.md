# Gratia: A Microservices-Based Leftover Food Redistribution Platform

A full-stack food rescue platform connecting donors with NGOs to reduce food waste and fight hunger.

[![Go Version](https://img.shields.io/badge/Go-1.25+-00ADD8?style=flat&logo=go)](https://golang.org)
[![Next.js Version](https://img.shields.io/badge/Next.js-14+-000000?style=flat&logo=next.js)](https://nextjs.org)
[![PostgreSQL Version](https://img.shields.io/badge/PostgreSQL-16-4169E1?style=flat&logo=postgresql)](https://postgresql.org)
[![Docker](https://img.shields.io/badge/Docker-24.0+-2496ED?style=flat&logo=docker)](https://docker.com)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

---

## Table of Contents

- [Overview](#overview)
- [Problem Statement](#problem-statement)
- [System Architecture](#system-architecture)
- [Services](#services)
- [Technology Stack](#technology-stack)
- [Getting Started](#getting-started)
- [API Documentation](#api-documentation)
- [Frontend](#frontend)
- [Database Schema](#database-schema)
- [Security](#security)
- [Deployment](#deployment)
- [Roadmap](#roadmap)
- [Contributing](#contributing)
- [License](#license)

---

## Overview

Gratia is a cloud-native platform that bridges the gap between food donors (restaurants, cafes, bakeries) and NGOs serving communities in need. The system manages the complete lifecycle of food donations - from listing surplus food to tracking its delivery to those who require it.

Built with a microservices architecture and a modern Next.js frontend, Gratia ensures scalability, maintainability, and a seamless user experience across all devices.

### Core Features

- Secure authentication with JWT-based access and refresh tokens
- Role-based access control for Donors, NGOs, and Administrators
- Food listing management with expiry tracking
- Complete claim lifecycle management from request to delivery
- NGO verification system with administrative oversight
- Responsive design compatible with desktop, tablet, and mobile devices
- Containerized deployment with Docker Compose
- Service-to-service communication with validation

---

## Problem Statement

### The Challenge

Food waste represents a significant global issue, with approximately 1.3 billion tons of food discarded annually. Concurrently, food insecurity affects millions of individuals who lack consistent access to nutritious meals. The connection between these two realities remains inefficient - donors lack visibility into where food is needed, NGOs lack information about available food sources, and no trusted mechanism exists to facilitate legitimate exchanges.

### Our Solution

Gratia creates a verified digital marketplace that enables:
- Donors to list surplus food with expiry tracking and location details
- Verified NGOs to claim available food for their communities
- Administrators to oversee platform operations and verify NGO legitimacy
- All participants to track the complete journey from donation to delivery

---

## System Architecture

### Microservices Overview

```
┌─────────────────────────────────────────────────────────────────────────┐
│                         Frontend (Next.js 14)                          │
│                      Single Page Application with SSR                  │
└─────────────────────────────────┬───────────────────────────────────────┘
                                  │
                                  ▼
┌─────────────────────────────────────────────────────────────────────────┐
│                         API Gateway (Planned)                          │
└─────────────────────────────────┬───────────────────────────────────────┘
                                  │
        ┌─────────────────────────┼─────────────────────────┐
        │                         │                         │
        ▼                         ▼                         ▼
┌───────────────┐       ┌───────────────┐       ┌───────────────┐
│   Auth        │       │   User        │       │   Food        │
│   Service     │       │   Service     │       │   Service     │
│   Port: 8080  │       │   Port: 8081  │       │   Port: 8082  │
└───────────────┘       └───────────────┘       └───────────────┘
        │                         │                         │
        └─────────────────────────┼─────────────────────────┘
                                  │
                                  ▼
                        ┌───────────────┐
                        │   Claim       │
                        │   Service     │
                        │   Port: 8083  │
                        └───────────────┘
                                  │
                                  ▼
                        ┌───────────────┐
                        │  PostgreSQL   │
                        │  Database     │
                        │   Version 16  │
                        └───────────────┘
```

### Architecture Principles

1. Service Independence: Each service owns its database schema and business logic
2. Clear Boundaries: No service directly accesses another service's database
3. Token-Based Authentication: JWT tokens issued by Auth Service are trusted by all services
4. Role Validation: Services validate user roles via User Service internal APIs
5. Horizontal Scalability: Each service can scale independently based on load

---

## Services

### Authentication Service (auth_service)

Purpose: Identity and access management for the entire platform.

| Endpoint | Method | Description |
|----------|--------|-------------|
| /register | POST | Register a new user |
| /login | POST | User login with email/password |
| /refresh | POST | Refresh access token |
| /logout | POST | Logout and revoke refresh token |
| /forgot-password | POST | Request password reset |
| /reset-password | POST | Reset password with token |
| /me | GET | Get current user information |
| /health | GET | Health check endpoint |

Technical Details:
- Passwords hashed with bcrypt (cost factor 10)
- JWT access tokens (15 minute expiry)
- Refresh tokens (7 days expiry) stored in database
- Token rotation for enhanced security
- Session tracking with IP address and user agent

---

### User Service (user_service)

Purpose: Manage donor and NGO profile data.

| Endpoint | Method | Description |
|----------|--------|-------------|
| /donors/profile | POST | Create donor profile |
| /donors/profile/me | GET | Get my donor profile |
| /donors/profile/me | PUT | Update my donor profile |
| /ngos/profile | POST | Create NGO profile |
| /ngos/profile/me | GET | Get my NGO profile |
| /admin/ngos/verify | PUT | Verify NGO (Admin only) |
| /internal/users/{id}/donor | GET | Internal donor validation |
| /internal/users/{id}/ngo | GET | Internal NGO validation |
| /health | GET | Health check endpoint |

Technical Details:
- One profile per user (unique constraint on user_id)
- NGO verification requires admin role
- Internal endpoints for service-to-service communication
- Profile metadata stored separately from authentication

---

### Food Service (food_service)

Purpose: Manage food donation listings and availability.

| Endpoint | Method | Description |
|----------|--------|-------------|
| /foods | POST | Create food listing (Donor only) |
| /foods | GET | List available food |
| /foods/{id} | GET | Get food listing details |
| /foods/{id} | PUT | Update food listing (Owner only) |
| /foods/{id} | DELETE | Cancel food listing (Owner only) |
| /internal/foods/{id}/validate | GET | Validate food for claiming |
| /internal/foods/{id}/claim | PATCH | Mark food as claimed |
| /health | GET | Health check endpoint |

Technical Details:
- Donor validation via User Service internal API
- Automatic expiry via background worker (runs every minute)
- Status tracking: AVAILABLE, CLAIMED, EXPIRED, CANCELLED
- Expiry time must be in the future

---

### Claim Service (claim_service)

Purpose: Manage the claim lifecycle from request to delivery.

| Endpoint | Method | Description |
|----------|--------|-------------|
| /claims | POST | Create a claim (NGO only) |
| /claims/{id}/approve | POST | Approve claim (Donor only) |
| /claims/{id}/reject | POST | Reject claim (Donor only) |
| /claims/{id}/cancel | POST | Cancel claim (NGO only) |
| /claims/{id}/pickup | POST | Mark as picked up (NGO only) |
| /claims/{id}/deliver | POST | Mark as delivered (NGO only) |
| /health | GET | Health check endpoint |

Claim Status Flow:

```
                  ┌─────────────────────────────────────────────┐
                  │                                             │
                  ▼                                             │
    ┌─────────────┐    ┌─────────────┐    ┌─────────────┐     │
    │   CREATED   │───▶│   ACCEPTED  │───▶│  PICKED_UP  │─────┘
    └─────────────┘    └─────────────┘    └─────────────┘
          │                  │                   │
          ▼                  ▼                   ▼
    ┌─────────────┐    ┌─────────────┐    ┌─────────────┐
    │   REJECTED  │    │  CANCELLED  │    │  DELIVERED  │
    └─────────────┘    └─────────────┘    └─────────────┘
```

Technical Details:
- Validates NGO status via User Service
- Validates food availability via Food Service
- Uses FOR UPDATE row locking to prevent race conditions
- Unique partial index prevents duplicate active claims
- Timestamps for each status transition

---

## Technology Stack

### Backend

| Technology | Version | Purpose |
|------------|---------|---------|
| Go | 1.25+ | Primary programming language |
| PostgreSQL | 16 | Relational database |
| pgx | v5 | PostgreSQL driver and toolkit |
| gorilla/mux | v1.8 | HTTP router |
| JWT | v5 | Authentication tokens |
| bcrypt | - | Password hashing |
| Docker | 24.0+ | Containerization |

### Frontend

| Technology | Version | Purpose |
|------------|---------|---------|
| Next.js | 14+ | React framework with SSR |
| TypeScript | 5+ | Type-safe JavaScript |
| Tailwind CSS | 4+ | Utility-first styling |
| React Query | 5+ | Data fetching and caching |
| React Hook Form | 7+ | Form management |
| Zod | 3+ | Schema validation |
| Radix UI | - | Headless UI components |
| Framer Motion | 11+ | Animations |

### DevOps and Infrastructure

| Technology | Version | Purpose |
|------------|---------|---------|
| Docker | 24.0+ | Containerization |
| Docker Compose | 2.0+ | Multi-container orchestration |
| Git | - | Version control |
| GitHub | - | Source code hosting |

---

## Getting Started

### Prerequisites

- Docker and Docker Compose
- Node.js 18+ (for frontend development)
- Go 1.25+ (for backend development)
- PostgreSQL 16 (if running without Docker)

### Quick Start with Docker (Recommended)

```bash
# 1. Clone the repository
git clone https://github.com/yourusername/gratia.git
cd gratia

# 2. Start all services
docker compose up -d --build

# 3. Check service health
curl http://localhost:8080/health  # Auth Service
curl http://localhost:8081/health  # User Service
curl http://localhost:8082/health  # Food Service
curl http://localhost:8083/health  # Claim Service

# 4. Start the frontend
cd frontend
npm install
npm run dev

# 5. Open your browser
# Frontend: http://localhost:3000
```

### Manual Setup

#### Backend Services

```bash
# Auth Service
cd auth_service
go mod download
go run cmd/server/main.go

# User Service (in a new terminal)
cd user_service
go mod download
go run cmd/server/main.go

# Food Service (in a new terminal)
cd food_service
go mod download
go run cmd/server/main.go

# Claim Service (in a new terminal)
cd claim_service
go mod download
go run cmd/server/main.go
```

#### Frontend

```bash
cd frontend
npm install
npm run dev
```

### Environment Variables

Each service uses environment variables for configuration:

#### Auth Service
```env
PORT=8080
DATABASE_URL=postgresql://gratia:gratia123@localhost:5432/gratia
JWT_SECRET=your-super-secret-jwt-key
ACCESS_TOKEN_TTL=15m
REFRESH_TOKEN_TTL=168h
ENV=development
LOG_LEVEL=info
```

#### User Service
```env
PORT=8081
DATABASE_URL=postgresql://gratia:gratia123@localhost:5432/gratia
JWT_SECRET=your-super-secret-jwt-key
LOG_LEVEL=info
ENV=development
```

#### Food Service
```env
PORT=8082
DATABASE_URL=postgresql://gratia:gratia123@localhost:5432/gratia
JWT_SECRET=your-super-secret-jwt-key
USER_SERVICE_URL=http://localhost:8081
LOG_LEVEL=info
ENV=development
```

#### Claim Service
```env
PORT=8083
DATABASE_URL=postgresql://gratia:gratia123@localhost:5432/gratia
JWT_SECRET=your-super-secret-jwt-key
USER_SERVICE_URL=http://localhost:8081
FOOD_SERVICE_URL=http://localhost:8082
LOG_LEVEL=info
ENV=development
```

#### Frontend
```env
NEXT_PUBLIC_AUTH_API_URL=http://localhost:8080
NEXT_PUBLIC_USER_API_URL=http://localhost:8081
NEXT_PUBLIC_FOOD_API_URL=http://localhost:8082
NEXT_PUBLIC_CLAIM_API_URL=http://localhost:8083
```

---

## API Documentation

### Authentication Flow

**1. Register a User**
```http
POST /register
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "SecurePass123",
  "role": "user"
}
```

**2. Login**
```http
POST /login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "SecurePass123"
}
```

Response:
```json
{
  "accessToken": "eyJhbGciOiJIUzI1NiIs...",
  "refreshToken": "a1b2c3d4e5f6..."
}
```

**3. Use Access Token**
```http
GET /foods
Authorization: Bearer eyJhbGciOiJIUzI1NiIs...
```

**4. Refresh Token**
```http
POST /refresh
Content-Type: application/json

{
  "refreshToken": "a1b2c3d4e5f6..."
}
```

### Testing Flow

1. Register: POST /auth/register
2. Login: POST /auth/login (save tokens)
3. Create Donor Profile: POST /user/donors/profile
4. Create Food Listing: POST /food/foods
5. Create NGO Profile: POST /user/ngos/profile
6. Verify NGO (Admin): PUT /user/admin/ngos/verify?userId={id}
7. Create Claim: POST /claim/claims
8. Approve Claim: POST /claim/claims/{id}/approve
9. Pickup: POST /claim/claims/{id}/pickup
10. Deliver: POST /claim/claims/{id}/deliver

---

## Frontend

### Pages and Routes

| Route | Description | Access |
|-------|-------------|--------|
| / | Landing page | Public |
| /login | User login | Public |
| /register | User registration | Public |
| /forgot-password | Password reset | Public |
| /donor | Donor dashboard | Donor |
| /donor/listings | My food listings | Donor |
| /donor/create | Create food listing | Donor |
| /ngo | NGO dashboard | NGO |
| /ngo/foods | Available food | NGO |
| /ngo/claims | My claims | NGO |
| /admin | Admin dashboard | Admin |
| /admin/users | User management | Admin |
| /admin/ngos | NGO management | Admin |

### Key Components

- UI Components: Button, Input, Card, Badge, Avatar, Select
- Forms: Login, Register, Forgot Password, Create Listing
- Layout: Sidebar, Navbar, Footer
- Dashboard: Stats cards, listing cards, claim cards

### Styling

- Tailwind CSS v4 for utility-first styling
- Custom theme with green primary color scheme
- Light theme with clean, minimal design
- Mobile-first responsive approach

---

## Database Schema

### Auth Service

```sql
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    role VARCHAR(50) NOT NULL,
    email_verified BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE refresh_tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id),
    token TEXT UNIQUE NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    revoked BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id),
    refresh_token_id UUID REFERENCES refresh_tokens(id),
    user_agent TEXT,
    ip_address TEXT,
    is_current BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### User Service

```sql
CREATE TABLE donor_profiles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) UNIQUE,
    name VARCHAR(255) NOT NULL,
    phone VARCHAR(50),
    address TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE ngo_profiles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) UNIQUE,
    organization VARCHAR(255) NOT NULL,
    registration_no VARCHAR(255) NOT NULL,
    verified BOOLEAN DEFAULT FALSE,
    verified_by UUID REFERENCES users(id),
    verified_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### Food Service

```sql
CREATE TABLE food_listings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    donor_user_id UUID REFERENCES users(id),
    title VARCHAR(255) NOT NULL,
    description TEXT,
    quantity INTEGER NOT NULL,
    unit VARCHAR(50) NOT NULL,
    expiry_time TIMESTAMP NOT NULL,
    location TEXT NOT NULL,
    image_url TEXT,
    status VARCHAR(50) DEFAULT 'AVAILABLE',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_food_listings_status ON food_listings(status);
CREATE INDEX idx_food_listings_donor ON food_listings(donor_user_id);
```

### Claim Service

```sql
CREATE TABLE claims (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    food_listing_id UUID REFERENCES food_listings(id),
    ngo_user_id UUID REFERENCES users(id),
    donor_user_id UUID REFERENCES users(id),
    status VARCHAR(50) DEFAULT 'CREATED',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    accepted_at TIMESTAMP,
    rejected_at TIMESTAMP,
    picked_up_at TIMESTAMP,
    delivered_at TIMESTAMP,
    cancelled_at TIMESTAMP
);

CREATE UNIQUE INDEX idx_unique_active_claim 
ON claims(food_listing_id) 
WHERE status IN ('CREATED', 'ACCEPTED', 'PICKED_UP');

CREATE INDEX idx_claims_food_listing ON claims(food_listing_id);
CREATE INDEX idx_claims_ngo_user ON claims(ngo_user_id);
```

---

## Security

### Authentication and Authorization

- JWT Tokens: Access tokens (15 min) and refresh tokens (7 days)
- Password Security: Bcrypt hashing with cost factor 10
- Role-Based Access: Donor, NGO, Admin with specific permissions
- Token Revocation: Refresh tokens can be revoked on logout

### Service-to-Service Security

- Internal APIs: Not exposed publicly
- Trusted Network: Services communicate within Docker network
- Validation: Each service validates user roles via User Service

### Data Security

- Database Isolation: Each service owns its database schema
- No Direct Access: Services never access another service's database
- Input Validation: All inputs validated with Zod schemas
- SQL Injection Prevention: Parameterized queries using pgx

### Best Practices

- Environment Variables: Secrets never committed to version control
- Non-Root Users: Services run as non-root users in containers
- Health Checks: Services have health check endpoints
- Graceful Shutdown: Services handle SIGTERM gracefully

---

## Deployment

### Docker Compose

```bash
# Build and start all services
docker compose up -d --build

# View logs
docker compose logs -f

# Stop all services
docker compose down

# Stop and remove volumes
docker compose down -v
```

### Production Deployment

1. Update Environment Variables: Use production database URLs and secrets
2. Build Images: docker compose build --no-cache
3. Push to Registry: docker push your-registry/gratia-*
4. Deploy: Pull images and start containers on production server
5. Configure Reverse Proxy: Nginx/Traefik for HTTPS
6. Set Up Monitoring: Prometheus/Grafana for metrics

### CI/CD Pipeline (Planned)

```yaml
name: Deploy
on:
  push:
    branches: [main]
jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: docker/setup-buildx-action@v2
      - name: Build and push Docker images
        run: docker compose build && docker compose push
      - name: Deploy to server
        run: ssh user@server 'cd /app && docker compose pull && docker compose up -d'
```

---

## Roadmap

### Short Term (Q1 2025)
- API Gateway with rate limiting
- Logging aggregation (ELK Stack)
- Monitoring with Prometheus/Grafana
- Real-time notifications (WebSocket)
- Email notifications for claim updates

### Medium Term (Q2 2025)
- Mobile App (React Native)
- Advanced search and filtering
- Analytics dashboard for impact tracking
- Multiple image upload for food listings
- PDF reports for donors and NGOs

### Long Term (Q3 2025+)
- AI-powered food matching algorithm
- Chat system between donors and NGOs
- Payment integration for premium features
- Multi-language support
- Integration with food delivery services
- Blockchain for transparent tracking

---

## Contributing

Contributions are welcome. Please follow these steps:

1. Fork the repository
2. Create a feature branch (git checkout -b feature/amazing-feature)
3. Commit your changes (git commit -m 'Add amazing feature')
4. Push to the branch (git push origin feature/amazing-feature)
5. Open a Pull Request

### Development Guidelines

- Go: Follow standard Go conventions
- Frontend: Use TypeScript with strict mode
- Testing: Write tests for new features
- Documentation: Update README and API docs
- Commits: Use conventional commit messages

### Code Style

- Backend: go fmt and go vet
- Frontend: ESLint and Prettier
- Database: Use migrations for schema changes

---

## License

This project is licensed under the MIT License - see the LICENSE file for details.

---

## Acknowledgments

Built to address food waste and food insecurity. Thanks to all contributors and open-source projects.

---

## Support

- Issues: GitHub Issues
- Discussions: GitHub Discussions
- Email: adityawaradkar2004@gmail.com
