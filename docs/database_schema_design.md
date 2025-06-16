

# Database Schema Design for Gratia

## Overview

This document details the planned database schema for Gratia — a cloud-native food donation platform connecting restaurants and NGOs. It includes entity definitions, tables, key fields, and relationships to support the microservices architecture and functional requirements.

***


## Database Design Principles

- **Normalization:** Designed to at least 3NF to minimize redundancy.

- **Scalability:** Schema supports horizontal scaling and partitioning.

- **Security:** Sensitive data (e.g., passwords) stored securely.

- **Extensibility:** Fields and tables designed for future feature expansion.

- **Multi-Tenancy:** Support multiple organizations (restaurants/NGOs) with proper access control.

***


## Entities and Tables

### 1. Users

| Field          | Type         | Description                              | Constraints      |
| -------------- | ------------ | ---------------------------------------- | ---------------- |
| user\_id       | UUID (PK)    | Unique user identifier                   | Primary Key      |
| email          | VARCHAR(255) | User email (login)                       | Unique, Not Null |
| password\_hash | VARCHAR(255) | Hashed password                          | Not Null         |
| full\_name     | VARCHAR(255) | User full name                           |                  |
| phone\_number  | VARCHAR(20)  | Contact phone                            |                  |
| role           | ENUM         | User role (`restaurant`, `ngo`, `admin`) | Not Null         |
| location\_lat  | DECIMAL(9,6) | Latitude for user location               | Nullable         |
| location\_lng  | DECIMAL(9,6) | Longitude for user location              | Nullable         |
| is\_active     | BOOLEAN      | Account active status                    | Default TRUE     |
| created\_at    | TIMESTAMP    | Account creation timestamp               | Default NOW()    |
| updated\_at    | TIMESTAMP    | Last update timestamp                    |                  |

***


### 2. Restaurants (Extension of Users)

| Field            | Type          | Description                       | Constraints              |
| ---------------- | ------------- | --------------------------------- | ------------------------ |
| user\_id         | UUID (PK, FK) | References `users(user_id)`       | Primary Key, Foreign Key |
| restaurant\_name | VARCHAR(255)  | Official restaurant/catering name |                          |
| address          | TEXT          | Full address                      |                          |
| verified         | BOOLEAN       | Whether restaurant is verified    | Default FALSE            |

_Note: Restaurants are a subset of users with `role = 'restaurant'`._

***


### 3. NGOs (Extension of Users)

| Field     | Type          | Description                 | Constraints              |
| --------- | ------------- | --------------------------- | ------------------------ |
| user\_id  | UUID (PK, FK) | References `users(user_id)` | Primary Key, Foreign Key |
| ngo\_name | VARCHAR(255)  | NGO organization name       |                          |
| address   | TEXT          | NGO office address          |                          |
| verified  | BOOLEAN       | NGO verification status     | Default FALSE            |

_Note: NGOs are a subset of users with `role = 'ngo'`._

***


### 4. Donations

| Field                 | Type         | Description                                                             | Constraints    |
| --------------------- | ------------ | ----------------------------------------------------------------------- | -------------- |
| donation\_id          | UUID (PK)    | Unique donation identifier                                              | Primary Key    |
| restaurant\_id        | UUID (FK)    | References `restaurants(user_id)`                                       | Foreign Key    |
| title                 | VARCHAR(255) | Donation title                                                          | Not Null       |
| description           | TEXT         | Detailed description                                                    |                |
| quantity              | INT          | Quantity of food items donated                                          | Not Null       |
| expiry\_time          | TIMESTAMP    | Expiration datetime for donation                                        | Not Null       |
| pickup\_location\_lat | DECIMAL(9,6) | Pickup location latitude                                                | Not Null       |
| pickup\_location\_lng | DECIMAL(9,6) | Pickup location longitude                                               | Not Null       |
| status                | ENUM         | Donation status (`open`, `claimed`, `picked_up`, `expired`, `archived`) | Default `open` |
| created\_at           | TIMESTAMP    | Donation creation timestamp                                             | Default NOW()  |
| updated\_at           | TIMESTAMP    | Last update timestamp                                                   |                |

***


### 5. Claims

| Field        | Type      | Description                                        | Constraints        |
| ------------ | --------- | -------------------------------------------------- | ------------------ |
| claim\_id    | UUID (PK) | Unique claim identifier                            | Primary Key        |
| donation\_id | UUID (FK) | References `donations(donation_id)`                | Foreign Key        |
| ngo\_id      | UUID (FK) | References `ngos(user_id)`                         | Foreign Key        |
| status       | ENUM      | Claim status (`en_route`, `received`, `cancelled`) | Default `en_route` |
| claimed\_at  | TIMESTAMP | Time claim was made                                | Default NOW()      |
| updated\_at  | TIMESTAMP | Last status update timestamp                       |                    |

_Note: Enforces one-to-one lock on claims to prevent multiple NGOs claiming same donation._

***


### 6. Feedback

| Field        | Type      | Description                         | Constraints   |
| ------------ | --------- | ----------------------------------- | ------------- |
| feedback\_id | UUID (PK) | Unique feedback identifier          | Primary Key   |
| donation\_id | UUID (FK) | References `donations(donation_id)` | Foreign Key   |
| ngo\_id      | UUID (FK) | References `ngos(user_id)`          | Foreign Key   |
| rating       | INT       | Rating score (e.g., 1-5)            | Not Null      |
| comments     | TEXT      | Feedback comments                   | Nullable      |
| created\_at  | TIMESTAMP | Feedback submission time            | Default NOW() |

***


### 7. Notifications

| Field            | Type        | Description                                   | Constraints       |
| ---------------- | ----------- | --------------------------------------------- | ----------------- |
| notification\_id | UUID (PK)   | Unique notification identifier                | Primary Key       |
| user\_id         | UUID (FK)   | Recipient user                                | Foreign Key       |
| type             | VARCHAR(50) | Notification type (email, SMS, in-app)        | Not Null          |
| channel          | VARCHAR(50) | Delivery channel                              | Not Null          |
| message          | TEXT        | Notification content                          | Not Null          |
| status           | ENUM        | Delivery status (`pending`, `sent`, `failed`) | Default `pending` |
| created\_at      | TIMESTAMP   | Timestamp of notification creation            | Default NOW()     |
| updated\_at      | TIMESTAMP   | Timestamp of last status update               |                   |

***


### 8. Audit Logs

| Field           | Type         | Description                                                  | Constraints   |
| --------------- | ------------ | ------------------------------------------------------------ | ------------- |
| audit\_id       | UUID (PK)    | Unique audit log entry                                       | Primary Key   |
| actor\_user\_id | UUID (FK)    | User who triggered the event                                 | Nullable      |
| event\_type     | VARCHAR(100) | Type of event (login, donation\_create, claim\_update, etc.) | Not Null      |
| event\_data     | JSONB        | Structured event details                                     |               |
| created\_at     | TIMESTAMP    | Timestamp of event occurrence                                | Default NOW() |

***


## Relationships Diagram (Summary)

- **Users** have a role that determines if they are **Restaurants** or **NGOs** (one-to-one extension).

- **Restaurants** create multiple **Donations** (one-to-many).

- **Donations** can be **Claimed** by exactly one **NGO** (one-to-one lock via Claims).

- **NGOs** provide **Feedback** on Donations they claimed (one-to-many).

- **Users** receive **Notifications**.

- All critical actions generate entries in **Audit Logs**.

***


## Notes on Implementation

- Use UUIDs as primary keys for uniqueness across distributed systems.

- Passwords stored as salted bcrypt hashes in `password_hash` field.

- ENUM types enforce valid statuses.

- Timestamps managed with timezone awareness.

- Geospatial data stored as decimal lat/lng; optionally use PostGIS extension for advanced geo queries.

- JSONB in audit logs allows flexible, schema-less event data.

***


## Conclusion

This schema provides a robust foundation for the Gratia platform supporting microservices interactions, security, and data integrity. It balances normalization with query performance and prepares for scaling and extensibility.

***


# Supabase Integration for Gratia

Supabase is an excellent backend-as-a-service platform built on top of PostgreSQL, offering a fully managed, scalable, and developer-friendly solution.


### Why Supabase?

- **Managed PostgreSQL:** No need to worry about database setup, backups, or maintenance.

- **Free Tier:** Offers up to 500 MB of database storage, perfect for early-stage projects like Gratia.

- **Authentication:** Built-in user management and secure authentication.

- **Realtime:** Supports real-time subscriptions to database changes for instant UI updates.

- **APIs:** Automatically generates RESTful and GraphQL APIs for your database.

- **Storage & Edge Functions:** Additional services for file storage and serverless functions.


### How Gratia Benefits from Supabase

- Use your existing PostgreSQL schema as-is on Supabase.

- Leverage Supabase Auth to manage user sign-ups, sign-ins, and roles (`restaurant`, `ngo`, `admin`).

- Realtime functionality can power live donation and claim updates in the app.

- Built-in security policies (Row Level Security) can enforce multi-tenancy and data privacy.

- Easy integration with frontend frameworks and backend microservices.

- Low cost or free during development with simple upgrade paths.
