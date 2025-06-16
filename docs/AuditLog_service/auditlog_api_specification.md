***


# Audit Log Service API Specification

***


## Overview

The Audit Log Service captures and stores detailed records of critical events and user/system actions across all Gratia microservices. It supports secure retrieval, filtering, and analysis of logs by administrators for compliance, monitoring, and debugging.

***


## Base URL

    /api/v1/audit

***


## Endpoints

***


### 1. Record Audit Event (Internal)

**POST** `/record`

Used internally by microservices to record an audit event. Not accessible by end users.


#### Headers

    X-Internal-Token: <service_token>
    Content-Type: application/json


#### Request Body

```json
{
  "actor_id": "uuid-of-user-or-system",
  "actor_type": "ngo", // or "restaurant", "admin", "system"
  "action": "donation_claimed",
  "target_id": "donation-uuid",
  "target_type": "donation",
  "metadata": {
    "ip": "192.168.1.12",
    "user_agent": "Mozilla/5.0",
    "notes": "Claimed donation ID successfully"
  },
  "timestamp": "2025-06-16T10:00:00Z"
}
```


#### Response (201 Created)

```json
{
  "message": "Audit event recorded successfully."
}
```

***


### 2. Get Logs (Admin Only)

**GET** `/logs`

Allows admin users to fetch and filter audit logs based on actor, action, or time range.


#### Headers

    Authorization: Bearer <admin_access_token>


#### Query Parameters (optional)

| Parameter    | Type     | Description                                     |
| ------------ | -------- | ----------------------------------------------- |
| `actor_id`   | string   | Filter logs by user/system ID                   |
| `actor_type` | string   | Filter by actor type (ngo, restaurant, admin)   |
| `action`     | string   | Filter by action name (e.g. `donation_claimed`) |
| `target_id`  | string   | Filter by target object ID                      |
| `from`       | datetime | Start time (ISO8601)                            |
| `to`         | datetime | End time (ISO8601)                              |
| `limit`      | integer  | Number of logs per page (default: 50)           |
| `page`       | integer  | Page number (default: 1)                        |


#### Response (200 OK)

```json
{
  "logs": [
    {
      "id": "uuid-log-entry",
      "actor_id": "uuid",
      "actor_type": "ngo",
      "action": "donation_claimed",
      "target_type": "donation",
      "target_id": "donation-uuid",
      "metadata": {
        "ip": "192.168.1.12",
        "user_agent": "Mozilla/5.0"
      },
      "timestamp": "2025-06-16T10:00:00Z"
    }
  ],
  "page": 1,
  "total_pages": 20
}
```

***


### 3. Get Actor's Log (Self)

**GET** `/my-logs`

Allows authenticated users (NGO/Restaurant) to view their own audit trail (limited actions).


#### Headers

    Authorization: Bearer <access_token>


#### Response (200 OK)

```json
{
  "logs": [
    {
      "action": "profile_updated",
      "target_type": "ngo_profile",
      "target_id": "uuid",
      "timestamp": "2025-06-15T08:30:00Z"
    }
  ]
}
```

***


## Common Audit Actions

| Action                | Description                        |
| --------------------- | ---------------------------------- |
| `donation_claimed`    | NGO claimed a donation             |
| `donation_posted`     | Restaurant created a donation      |
| `profile_updated`     | User profile was updated           |
| `login_success`       | User logged in                     |
| `login_failed`        | Failed login attempt               |
| `message_sent`        | NGO sent a message to a restaurant |
| `verification_passed` | NGO or restaurant got verified     |
| `admin_broadcasted`   | Admin sent out a system message    |

***


## Status Codes Summary

| Code | Meaning               |
| ---- | --------------------- |
| 200  | OK                    |
| 201  | Created               |
| 400  | Bad Request           |
| 401  | Unauthorized          |
| 403  | Forbidden (non-admin) |
| 500  | Internal Server Error |

***


## Authentication & Security

- **Recording (`POST /record`)** requires a valid internal service token.

- **Reading (`/logs`)** requires JWT with `admin` role.

- **Self logs (`/my-logs`)** are available to the logged-in user only.

- Timestamps are UTC in ISO 8601 format.

- PII in `metadata` should be minimal and encrypted if sensitive.

***


## Storage & Performance Notes

- Data stored in a time-series database (or append-only log store).

- Logs are immutable and append-only for security compliance.

- Indexed by actor\_id, action, and timestamp for fast filtering.

- Optionally archived or exported to S3/BigQuery for analytics.

***


## Future Enhancements

- Real-time stream to SIEM tools like Elastic Stack or Datadog.

- Anomaly detection on suspicious actions.

- Graph-based actor-object interaction history.

- Visual log timeline in admin dashboard.

***
