***


# Notification Service API Specification

***


## Overview

The Notification Service is responsible for managing all system-generated and user-triggered notifications within the Gratia platform. It supports real-time and persistent notifications for NGOs, restaurants, and admins. It provides RESTful APIs to fetch, read, and manage notifications securely.

***


## Base URL

    /api/v1/notifications

***


## Endpoints

***


### 1. Get Notifications for User

**GET** `/`

Fetches all notifications for the authenticated user (NGO, restaurant, or admin), with support for pagination and filtering by status.


#### Headers

    Authorization: Bearer <access_token>


#### Query Parameters (optional)

| Parameter | Type    | Description                            |
| --------- | ------- | -------------------------------------- |
| `status`  | string  | Filter by `read` or `unread`           |
| `limit`   | integer | Number of items per page (default: 10) |
| `page`    | integer | Page number (default: 1)               |


#### Response (200 OK)

```json
{
  "notifications": [
    {
      "id": "uuid",
      "title": "Donation Claimed",
      "message": "Your claim for donation ID d1c2... has been confirmed.",
      "type": "donation_update",
      "status": "unread",
      "created_at": "2025-06-15T14:32:00Z"
    }
  ],
  "page": 1,
  "total_pages": 3
}
```

***


### 2. Mark Notification as Read

**PATCH** `/{notification_id}/read`

Marks a specific notification as read.


#### Headers

    Authorization: Bearer <access_token>


#### Path Parameters

| Parameter         | Type   | Description                    |
| ----------------- | ------ | ------------------------------ |
| `notification_id` | string | ID of the notification to mark |


#### Response (200 OK)

```json
{
  "message": "Notification marked as read."
}
```

***


### 3. Mark All as Read

**PATCH** `/read-all`

Marks **all** unread notifications for the user as read.


#### Headers

    Authorization: Bearer <access_token>


#### Response (200 OK)

```json
{
  "message": "All notifications marked as read."
}
```

***


### 4. Delete Notification

**DELETE** `/{notification_id}`

Deletes a specific notification from the user's inbox (soft delete).


#### Headers

    Authorization: Bearer <access_token>


#### Path Parameters

| Parameter         | Type   | Description                   |
| ----------------- | ------ | ----------------------------- |
| `notification_id` | string | Unique ID of the notification |


#### Response (200 OK)

```json
{
  "message": "Notification deleted successfully."
}
```

***


### 5. Send Notification (Internal/Event Triggered)

**POST** `/send`

Used internally by services (like donation service or admin service) to trigger a new notification. Requires service-to-service authorization.


#### Headers

    X-Internal-Token: <service_token>
    Content-Type: application/json


#### Request Body

```json
{
  "recipient_id": "uuid-of-user",
  "recipient_type": "ngo", // or "restaurant", "admin"
  "title": "New Donation Available",
  "message": "A new food donation is now available near your location.",
  "type": "donation_available"
}
```


#### Response (201 Created)

```json
{
  "message": "Notification sent successfully."
}
```

***


## Notification Types

| Type                 | Description                                              |
| -------------------- | -------------------------------------------------------- |
| `donation_update`    | Status updates on donations (claimed, picked, cancelled) |
| `donation_available` | New donation posted, near location                       |
| `admin_broadcast`    | System-wide broadcast by admins                          |
| `reminder`           | Pickup or expiration reminder                            |
| `system`             | System alerts, e.g. verification success                 |

***


## Authentication & Authorization

- **User-facing endpoints** (`GET`, `PATCH`, `DELETE`) require a valid JWT in the `Authorization` header.

- **Internal service-to-service** notification publishing requires an `X-Internal-Token` for trusted service access.

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
| 500  | Internal Server Error |

***


## Security Considerations

- Rate limiting on read/delete endpoints to avoid abuse.

- Notifications are user-specific and cannot be accessed by others.

- Internal token required to trigger cross-service notifications.

- Encryption and secure storage for sensitive notification content.

***


## Future Enhancements

- WebSocket support for **real-time delivery**.

- Push notification integration with Firebase / APNs.

- Grouping notifications by type/date.

- UI badge counter sync for unread notifications.
