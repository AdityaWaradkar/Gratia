# Restaurant Service API Specification

---

## Overview

The Restaurant Service manages all operations related to registered restaurants on the platform. It allows authenticated restaurants to manage their profile, list excess food for donation, and track donation history. The service enforces role-based access, ensuring only users with the `restaurant` role can access or modify restaurant-specific resources.

---

## Base URL

```
/api/v1/restaurant
```

---

## Endpoints

---

### 1. Get Restaurant Profile

**GET** `/me`

Returns the authenticated restaurant's profile information.

#### Headers

```
Authorization: Bearer <access_token>
```

#### Response (200 OK)

```json
{
  "restaurant_id": "e7a35f52-7891-478c-9f6e-c67b8d4eec3c",
  "name": "Shree Krupa Hotel",
  "email": "contact@shreekrupa.com",
  "phone_number": "+919812345678",
  "address": "MG Road, Pune, Maharashtra",
  "is_active": true,
  "created_at": "2025-06-01T12:00:00Z"
}
```

#### Errors

| Status | Message                  |
| ------ | ------------------------ |
| 401    | Invalid or expired token |
| 500    | Internal server error    |

---

### 2. Update Restaurant Profile

**PUT** `/me`

Allows a restaurant to update its profile information.

#### Headers

```
Authorization: Bearer <access_token>
```

#### Request Body

```json
{
  "name": "Shree Krupa Hotel",
  "phone_number": "+919812345678",
  "address": "Updated MG Road, Pune"
}
```

#### Response (200 OK)

```json
{
  "message": "Profile updated successfully."
}
```

#### Errors

| Status | Message               |
| ------ | --------------------- |
| 400    | Invalid input data    |
| 401    | Unauthorized          |
| 500    | Internal server error |

---

### 3. Create Food Donation Listing

**POST** `/donations`

Creates a new food donation listing.

#### Headers

```
Authorization: Bearer <access_token>
```

#### Request Body

```json
{
  "food_name": "Veg Pulao",
  "quantity": "10 kg",
  "pickup_time": "2025-06-01T18:00:00Z",
  "address": "MG Road, Pune",
  "notes": "Please bring containers."
}
```

#### Response (201 Created)

```json
{
  "message": "Donation listing created successfully.",
  "donation_id": "d9bb85a9-4cb0-4379-a64f-cf152e3b37e3"
}
```

#### Errors

| Status | Message               |
| ------ | --------------------- |
| 400    | Invalid input data    |
| 401    | Unauthorized          |
| 500    | Internal server error |

---

### 4. Get All Donation Listings by Restaurant

**GET** `/donations`

Returns all donation listings created by the authenticated restaurant.

#### Headers

```
Authorization: Bearer <access_token>
```

#### Response (200 OK)

```json
[
  {
    "donation_id": "d9bb85a9-4cb0-4379-a64f-cf152e3b37e3",
    "food_name": "Veg Pulao",
    "quantity": "10 kg",
    "pickup_time": "2025-06-01T18:00:00Z",
    "status": "available",
    "created_at": "2025-06-01T12:30:00Z"
  }
]
```

#### Errors

| Status | Message               |
| ------ | --------------------- |
| 401    | Unauthorized          |
| 500    | Internal server error |

---

### 5. Delete a Donation Listing

**DELETE** `/donations/{donation_id}`

Deletes a specific donation listing created by the restaurant.

#### Headers

```
Authorization: Bearer <access_token>
```

#### Response (200 OK)

```json
{
  "message": "Donation listing deleted successfully."
}
```

#### Errors

| Status | Message                    |
| ------ | -------------------------- |
| 401    | Unauthorized               |
| 404    | Donation listing not found |
| 500    | Internal server error      |

---

## Security Considerations

* JWT tokens are required for all endpoints.
* Only users with the `restaurant` role may access this service.
* Role-based access control (RBAC) is enforced by middleware.
* Input validation and sanitization are performed server-side.

---

## Status Codes Summary

| Code | Meaning               |
| ---- | --------------------- |
| 200  | OK                    |
| 201  | Created               |
| 400  | Bad Request           |
| 401  | Unauthorized          |
| 404  | Not Found             |
| 500  | Internal Server Error |

---

## Future Enhancements

* Restaurant donation analytics dashboard
* Food pickup confirmation via OTP or digital signature
* Ability to mark donations as "claimed" or "fulfilled"
* Notification system for donation status changes
