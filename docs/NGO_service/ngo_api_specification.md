
# NGO Service API Specification

---

## Overview

The NGO Service is responsible for managing operations related to NGO users, including profile management, viewing and claiming food donations, donation history, and communication with restaurants. The service ensures secure, role-based access limited to authenticated NGO users.

---

## Base URL

```
/api/v1/ngo
```

---

## Endpoints

---

### 1. Get NGO Profile

**GET** `/profile`

Retrieves the profile information of the authenticated NGO user.

#### Headers

```
Authorization: Bearer <access_token>
```

#### Response (200 OK)

```json
{
  "ngo_id": "uuid-generated",
  "name": "Helping Hands NGO",
  "email": "contact@helpinghands.org",
  "phone_number": "+919812345678",
  "address": "123 Charity Street, City, State, Zip",
  "is_verified": true,
  "created_at": "2025-06-01T10:00:00Z"
}
```

#### Errors

| Status | Message               |
| ------ | --------------------- |
| 401    | Unauthorized          |
| 404    | NGO profile not found |
| 500    | Internal server error |

---

### 2. Update NGO Profile

**PUT** `/profile`

Updates NGO profile details.

#### Headers

```
Authorization: Bearer <access_token>
```

#### Request Body

```json
{
  "name": "Helping Hands NGO",
  "phone_number": "+919812345678",
  "address": "123 Charity Street, City, State, Zip"
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
| 404    | NGO profile not found |
| 500    | Internal server error |

---

### 3. List Available Donations

**GET** `/donations`

Lists all available food donations for claim by the NGO, with optional filtering and pagination.

#### Headers

```
Authorization: Bearer <access_token>
```

#### Query Parameters (optional)

| Parameter  | Type    | Description                            |
| ---------- | ------- | ------------------------------------ |
| `page`     | integer | Page number (default: 1)              |
| `limit`    | integer | Number of donations per page (default: 10) |
| `location` | string  | Filter donations by location proximity (city or GPS coordinates) |

#### Response (200 OK)

```json
{
  "donations": [
    {
      "donation_id": "uuid-generated",
      "restaurant_id": "uuid-generated",
      "food_description": "50 meals of cooked rice and vegetables",
      "quantity": 50,
      "pickup_time": "2025-06-05T18:00:00Z",
      "location": "Restaurant Address or GPS coordinates",
      "status": "available"
    }
  ],
  "page": 1,
  "total_pages": 5,
  "total_items": 50
}
```

#### Errors

| Status | Message       |
| ------ | ------------- |
| 401    | Unauthorized  |
| 500    | Internal error|

---

### 4. Claim Donation

**POST** `/donations/{donation_id}/claim`

Allows the NGO to claim a specific donation. Claims are processed on a first-come, first-served basis.

#### Headers

```
Authorization: Bearer <access_token>
```

#### Path Parameters

| Parameter     | Type   | Description                   |
| ------------- | ------ | -----------------------------|
| `donation_id` | string | Unique identifier of donation |

#### Response (200 OK)

```json
{
  "message": "Donation claimed successfully.",
  "donation_id": "uuid-generated"
}
```

#### Errors

| Status | Message                  |
| ------ | ------------------------ |
| 400    | Donation already claimed |
| 401    | Unauthorized             |
| 404    | Donation not found       |
| 500    | Internal server error    |

---

### 5. View Donation History

**GET** `/donations/history`

Returns a list of donations previously claimed by the NGO.

#### Headers

```
Authorization: Bearer <access_token>
```

#### Response (200 OK)

```json
{
  "claimed_donations": [
    {
      "donation_id": "uuid-generated",
      "restaurant_name": "Foodies Delight",
      "food_description": "30 meals of pasta",
      "quantity": 30,
      "claimed_at": "2025-05-25T15:00:00Z",
      "status": "picked_up"
    }
  ]
}
```

#### Errors

| Status | Message           |
| ------ | ----------------- |
| 401    | Unauthorized      |
| 500    | Internal server error |

---

### 6. Messaging with Restaurants (Future Feature)

**POST** `/messages`

Allows NGOs to send messages to restaurants regarding donations.

#### Headers

```
Authorization: Bearer <access_token>
```

#### Request Body

```json
{
  "restaurant_id": "uuid-generated",
  "donation_id": "uuid-generated",
  "message": "Can we schedule pickup for tomorrow morning?"
}
```

#### Response (200 OK)

```json
{
  "message": "Message sent successfully."
}
```

#### Errors

| Status | Message                |
| ------ | ---------------------- |
| 400    | Invalid input data     |
| 401    | Unauthorized           |
| 404    | Restaurant or donation not found |
| 500    | Internal server error  |

---

## Authentication & Authorization

- All endpoints require a valid JWT bearer token in the Authorization header.
- Access is restricted to users with the `ngo` role.
- Middleware validates token integrity and user role before allowing access.

---

## Status Codes Summary

| Code | Meaning                 |
| ---- | ----------------------- |
| 200  | OK                      |
| 201  | Created                 |
| 400  | Bad Request             |
| 401  | Unauthorized            |
| 403  | Forbidden               |
| 404  | Not Found               |
| 409  | Conflict                |
| 500  | Internal Server Error   |

---

## Security Considerations

- All communications occur over HTTPS (TLS 1.3).
- Rate limiting protects critical endpoints from abuse.
- Input validation prevents injection and other attacks.
- Role-based access control ensures only authorized NGO users can perform actions.

---

## Future Enhancements

- Real-time notifications for new donations.
- Map integration for donation locations and navigation.
- Enhanced messaging and chat capabilities.
- NGO analytics dashboard to monitor donation history and impact.

---
