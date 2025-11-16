# Food Listing Service 

## 1. Overview

The **Food Listing Service** is a core microservice in the Gratia platform, responsible for managing all operations related to donated food items.\
It enables restaurants and catering services (donors) to create, update, and delete listings of leftover food that can be claimed by NGOs.

This service functions independently and interacts with other services such as:

- **Auth Service** – for user authentication and role verification.

- **Notification Service** – for notifying NGOs about new listings.

- **NGO Claim Service** – for managing claims on available listings.

***


## 2. Responsibilities

| Function                | Description                                                                        |
| ----------------------- | ---------------------------------------------------------------------------------- |
| **Create Food Listing** | Donors can create new food donation listings.                                      |
| **Update Listing**      | Donors can modify food details or change the listing status.                       |
| **Delete Listing**      | Donors can delete their own listings that are not claimed or already expired.      |
| **Get Listings**        | Retrieves available listings based on filters such as location or food type.       |
| **Get Single Listing**  | Returns detailed information about a specific listing.                             |
| **Set Availability**    | Updates the status of listings (open, claimed, expired) manually or automatically. |

***


## 3. Service Architecture

- **Service Name:** `food_listing_service`

- **Language:** Go (Golang)

- **Framework:** Fiber / Chi (recommended lightweight REST framework)

- **Database:** PostgreSQL

- **Authentication:** JWT tokens issued by Auth Service

- **Deployment:** Dockerized, managed via Kubernetes (EKS)

- **Monitoring:** Prometheus metrics with Grafana dashboards

- **Logging:** Structured logging using Zap or Logrus

- **Tracing:** OpenTelemetry with Jaeger

***


## 4. Database Schema

### Table: `food_listings`

| Field         | Type                               | Description                                                             |
| ------------- | ---------------------------------- | ----------------------------------------------------------------------- |
| `id`          | UUID                               | Primary key                                                             |
| `donor_id`    | UUID                               | Foreign key referencing users(id)                                       |
| `title`       | VARCHAR(100)                       | Short title describing the donation                                     |
| `description` | TEXT                               | Detailed description of the food                                        |
| `food_type`   | VARCHAR(50)                        | Example: “Veg”, “Non-Veg”, “Packed”, “Cooked”                           |
| `quantity`    | INTEGER                            | Number of portions or servings                                          |
| `expiry_time` | TIMESTAMP                          | Time until which the food is valid for consumption                      |
| `location`    | JSONB                              | Example: `{ "lat": 19.07, "lng": 72.87, "address": "Andheri, Mumbai" }` |
| `status`      | ENUM(`open`, `claimed`, `expired`) | Current state of the listing                                            |
| `created_at`  | TIMESTAMP                          | Record creation timestamp                                               |
| `updated_at`  | TIMESTAMP                          | Record update timestamp                                                 |

***


## 5. Authentication and Authorization

All endpoints require a valid JWT Bearer Token issued by the **Auth Service**.


### Roles and Permissions

| Role      | Permissions                                     |
| --------- | ----------------------------------------------- |
| **DONOR** | Create, update, delete, and view own listings   |
| **NGO**   | View all open listings and filter by parameters |
| **ADMIN** | View and delete any listing                     |

***


## 6. API Endpoints

All routes are prefixed with `/api/v1/food`.


### 6.1 Create Food Listing

**POST** `/api/v1/food`\
**Authorization:** `Bearer <JWT>` (Role: DONOR)

**Request Body:**

```json
{
  "title": "Surplus Lunch Boxes",
  "description": "20 boxes of veg biryani, packed and fresh.",
  "food_type": "Veg",
  "quantity": 20,
  "expiry_time": "2025-11-02T20:00:00Z",
  "location": {
    "lat": 19.07,
    "lng": 72.87,
    "address": "Andheri, Mumbai"
  }
}
```

**Response:**

```json
{
  "id": "8b1c89d0-efb3-4d6a-9f0d-5628e64c4b11",
  "message": "Food listing created successfully"
}
```

***


### 6.2 Get All Listings

**GET** `/api/v1/food`\
**Authorization:** Optional

**Query Parameters (Optional):**

| Parameter     | Example          | Description                                |
| ------------- | ---------------- | ------------------------------------------ |
| `status`      | `open`           | Filter by status                           |
| `type`        | `Veg`            | Filter by food type                        |
| `limit`       | `20`             | Pagination limit                           |
| `offset`      | `0`              | Pagination offset                          |
| `lat` / `lng` | `19.07`, `72.87` | Filter by proximity (within a 10km radius) |

**Response:**

```json
[
  {
    "id": "8b1c89d0-efb3-4d6a-9f0d-5628e64c4b11",
    "title": "Surplus Lunch Boxes",
    "food_type": "Veg",
    "quantity": 20,
    "status": "open",
    "expiry_time": "2025-11-02T20:00:00Z",
    "location": {
      "address": "Andheri, Mumbai"
    },
    "created_at": "2025-11-01T18:00:00Z"
  }
]
```

***


### 6.3 Get Single Listing

**GET** `/api/v1/food/{id}`\
**Authorization:** Optional

**Response:**

```json
{
  "id": "8b1c89d0-efb3-4d6a-9f0d-5628e64c4b11",
  "donor_id": "4f3e9ad1-ff23-4729-8b9a-c3c4d8a298b1",
  "title": "Surplus Lunch Boxes",
  "description": "20 boxes of veg biryani, packed and fresh.",
  "food_type": "Veg",
  "quantity": 20,
  "expiry_time": "2025-11-02T20:00:00Z",
  "location": {
    "lat": 19.07,
    "lng": 72.87,
    "address": "Andheri, Mumbai"
  },
  "status": "open",
  "created_at": "2025-11-01T18:00:00Z",
  "updated_at": "2025-11-01T19:00:00Z"
}
```

***


### 6.4 Update Food Listing

**PUT** `/api/v1/food/{id}`\
**Authorization:** `Bearer <JWT>` (Role: DONOR or ADMIN)

**Request Body:**

```json
{
  "title": "Surplus Dinner Boxes",
  "quantity": 15,
  "expiry_time": "2025-11-02T23:00:00Z"
}
```

**Response:**

```json
{ "message": "Food listing updated successfully" }
```

***


### 6.5 Delete Food Listing

**DELETE** `/api/v1/food/{id}`\
**Authorization:** `Bearer <JWT>` (Role: DONOR or ADMIN)

**Response:**

```json
{ "message": "Listing deleted successfully" }
```

***


### 6.6 Auto-Expire Listings (Internal)

**PATCH** `/api/v1/food/expire`\
**Authorization:** Internal service token only (not public)

- Marks all listings as `expired` where `expiry_time < NOW()`.

- Typically triggered by a scheduled CronJob every hour.

***


## 7. Business Logic and Rules

1. **Validation on Creation**

    - Only authenticated donors can create listings.

    - Expiry time must be in the future.

    - Quantity must be greater than zero.

2. **Auto Expiry**

    - Listings are automatically marked as expired after expiry time.

3. **Claim Integration**

    - Once an NGO claims a listing, its status changes to `claimed`.

    - An event is published to the message broker for synchronization with the Claim Service.

4. **Data Retention**

    - Expired and claimed listings are archived after 7 days for analytical purposes.

***


## 8. Service Layer Structure

| Layer                        | Responsibility                                                       |
| ---------------------------- | -------------------------------------------------------------------- |
| **Handler (HTTP)**           | Routing, validation, and response generation                         |
| **Service (Business Logic)** | Implements creation, update, expiry, and validation rules            |
| **Repository (Database)**    | Performs CRUD operations using `pgxpool`                             |
| **Model (Entities)**         | Defines data structures for listings                                 |
| **Event Layer (Kafka)**      | Publishes listing lifecycle events (`created`, `claimed`, `expired`) |

***


## 9. Example Data Structures

```text
type FoodListing struct {
	ID          string    `json:"id"`
	DonorID     string    `json:"donor_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	FoodType    string    `json:"food_type"`
	Quantity    int       `json:"quantity"`
	ExpiryTime  time.Time `json:"expiry_time"`
	Location    Location  `json:"location"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Location struct {
	Lat     float64 `json:"lat"`
	Lng     float64 `json:"lng"`
	Address string  `json:"address"`
}
```

***


## 10. Event Model

**Kafka Topic:** `listing-events`

**Example Event — `listing.created`:**

```json
{
  "event": "listing.created",
  "data": {
    "listing_id": "8b1c89d0-efb3-4d6a-9f0d-5628e64c4b11",
    "donor_id": "4f3e9ad1-ff23-4729-8b9a-c3c4d8a298b1",
    "status": "open",
    "timestamp": "2025-11-01T19:30:00Z"
  }
}
```

***


## 11. Directory Structure

    food_listing_service/
    │
    ├── cmd/
    │   └── main.go
    ├── internal/
    │   ├── handler/
    │   │   └── food_handler.go
    │   ├── service/
    │   │   └── food_service.go
    │   ├── repository/
    │   │   └── food_repository.go
    │   ├── model/
    │   │   └── food_model.go
    │   └── utils/
    │       └── validator.go
    ├── pkg/
    │   └── middleware/
    │       └── auth_middleware.go
    └── Dockerfile
