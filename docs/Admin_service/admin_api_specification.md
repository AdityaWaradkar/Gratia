# Admin Service API Specification



## Overview

The Admin Service provides role-restricted access to system administrators for managing NGO and restaurant users, monitoring donation activities, verifying accounts, and overseeing platform analytics and system health. It supports secure JWT-based access and enforces admin-only permissions.

***


## Base URL

    /api/v1/admin

***


## Endpoints

***


### 1. Get Admin Profile

**GET** `/profile`

Retrieves the authenticated admin's profile information.


#### Headers

    Authorization: Bearer <access_token>


#### Response (200 OK)

```json
{
  "admin_id": "uuid-generated",
  "name": "Admin Name",
  "email": "admin@gratia.org",
  "role": "super_admin",
  "created_at": "2025-06-01T10:00:00Z"
}
```

***


### 2. List All NGOs

**GET** `/ngos`

Fetches all registered NGOs, with optional filters for verification and pagination.


#### Headers

    Authorization: Bearer <access_token>


#### Query Parameters (optional)

| Parameter  | Type    | Description                   |
| ---------- | ------- | ----------------------------- |
| `verified` | boolean | Filter by verification status |
| `page`     | integer | Page number                   |
| `limit`    | integer | Items per page                |


#### Response (200 OK)

```json
{
  "ngos": [
    {
      "ngo_id": "uuid",
      "name": "Helping Hands",
      "email": "ngo@example.com",
      "is_verified": true,
      "created_at": "2025-06-01T12:00:00Z"
    }
  ],
  "page": 1,
  "total_items": 100
}
```

***


### 3. Verify NGO Account

**PATCH** `/ngos/{ngo_id}/verify`

Marks an NGO account as verified.


#### Path Parameters

| Parameter | Type   | Description             |
| --------- | ------ | ----------------------- |
| `ngo_id`  | string | NGO’s unique identifier |


#### Headers

    Authorization: Bearer <access_token>


#### Response (200 OK)

```json
{
  "message": "NGO account verified successfully."
}
```

***


### 4. List All Restaurants

**GET** `/restaurants`

Fetches all registered restaurants with optional filters.


#### Headers

    Authorization: Bearer <access_token>


#### Query Parameters (optional)

| Parameter  | Type    | Description              |
| ---------- | ------- | ------------------------ |
| `verified` | boolean | Filter by verification   |
| `page`     | integer | Pagination - page number |
| `limit`    | integer | Pagination - page size   |

***


### 5. Verify Restaurant Account

**PATCH** `/restaurants/{restaurant_id}/verify`

Approves a restaurant's verification.


#### Headers

    Authorization: Bearer <access_token>


#### Path Parameters

| Parameter       | Type   | Description                 |
| --------------- | ------ | --------------------------- |
| `restaurant_id` | string | Unique ID of the restaurant |

***


### 6. View All Donations

**GET** `/donations`

Fetches all donation entries across the platform.


#### Headers

    Authorization: Bearer <access_token>


#### Query Parameters (optional)

| Parameter | Type    | Description                                       |
| --------- | ------- | ------------------------------------------------- |
| `status`  | string  | Filter by status (available, claimed, picked\_up) |
| `page`    | integer | Page number                                       |
| `limit`   | integer | Items per page                                    |

***


### 7. View Donation Stats

**GET** `/stats/donations`

Provides high-level donation statistics.


#### Headers

    Authorization: Bearer <access_token>


#### Response (200 OK)

```json
{
  "total_donations": 200,
  "claimed_donations": 150,
  "pending_donations": 30,
  "picked_up_donations": 120,
  "total_meals_donated": 8200
}
```

***


### 8. System Health Check

**GET** `/health`

Returns current system and service health metrics.


#### Headers

    Authorization: Bearer <access_token>


#### Response (200 OK)

```json
{
  "services": {
    "auth_service": "healthy",
    "ngo_service": "healthy",
    "restaurant_service": "healthy",
    "donation_service": "healthy"
  },
  "uptime": "48 days, 4 hours",
  "database_connection": "connected"
}
```

***


### 9. Admin Broadcast Message (Future)

**POST** `/broadcast`

Sends a platform-wide message to NGOs or restaurants.


#### Headers

    Authorization: Bearer <access_token>


#### Request Body

```json
{
  "target_audience": "ngo", // or "restaurant" or "all"
  "message": "Please ensure pickups are completed by 8 PM today."
}
```


#### Response (200 OK)

```json
{
  "message": "Broadcast sent successfully."
}
```

***


## Authentication & Authorization

- All endpoints require a valid admin JWT in the Authorization header.

- Role-based restrictions apply (`admin`, `super_admin`).

- Token middleware ensures admin-only access.

***


## Status Codes Summary

| Code | Meaning               |
| ---- | --------------------- |
| 200  | OK                    |
| 201  | Created               |
| 400  | Bad Request           |
| 401  | Unauthorized          |
| 403  | Forbidden             |
| 404  | Not Found             |
| 409  | Conflict              |
| 500  | Internal Server Error |

***


## Security Considerations

- Only authenticated admins may access these endpoints.

- Admin actions are logged for audit trails.

- All endpoints are protected by rate limiting and input sanitization.

***


## Future Enhancements

- Role-based admin actions (e.g., read-only, super admin).

- Graph-based analytics for meals saved, donor impact.

- Export data reports (CSV, PDF).

- Real-time dashboard with live metrics (WebSockets).

***
