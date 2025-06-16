***


# Cleanup Service API Specification

***


## Overview

The Cleanup Service is responsible for maintaining data hygiene across the Gratia platform. It performs automated and manual cleanup of expired or orphaned records including donations, logs, sessions, notifications, and media. This service improves performance, reduces storage costs, and ensures system reliability.

***


## Base URL

    /api/v1/cleanup

***


## Endpoints

***


### 1. Trigger Manual Cleanup (Admin Only)

**POST** `/run`

Initiates a cleanup operation manually. Only accessible to admins or internal services.


#### Headers

    Authorization: Bearer <admin_access_token>


#### Request Body (optional)

```json
{
  "types": ["expired_donations", "logs", "media"], // optional
  "dry_run": false // if true, performs analysis only
}
```


#### Response (200 OK)

```json
{
  "message": "Cleanup executed successfully.",
  "summary": {
    "expired_donations": 12,
    "logs_deleted": 230,
    "media_removed": 4
  }
}
```

***


### 2. Get Cleanup Report

**GET** `/reports/latest`

Returns a summary of the most recent cleanup task including stats, errors (if any), and duration.


#### Headers

    Authorization: Bearer <admin_access_token>


#### Response (200 OK)

```json
{
  "timestamp": "2025-06-16T03:00:00Z",
  "executed_by": "system",
  "duration_ms": 1435,
  "stats": {
    "expired_donations": 20,
    "logs_deleted": 500,
    "media_removed": 10,
    "stale_sessions": 14
  },
  "status": "success",
  "errors": []
}
```

***


### 3. Schedule Cleanup Job (Internal)

**POST** `/schedule`

Used by the orchestration system (e.g., cron/Kubernetes CronJob) to register or modify scheduled jobs.


#### Headers

    X-Internal-Token: <service_token>


#### Request Body

```json
{
  "job_name": "daily_cleanup",
  "cron_expression": "0 3 * * *", // every day at 3 AM
  "enabled": true
}
```


#### Response (201 Created)

```json
{
  "message": "Cleanup job scheduled successfully."
}
```

***


## Cleanup Types Handled

| Cleanup Type        | Description                                           |
| ------------------- | ----------------------------------------------------- |
| `expired_donations` | Deletes food donations past their pickup time         |
| `logs`              | Deletes audit/system logs older than retention window |
| `media`             | Removes unused images/files (e.g., unlinked photos)   |
| `stale_sessions`    | Clears old authentication tokens or sessions          |
| `orphan_records`    | Removes references with missing relations             |
| `notifications`     | Deletes unread notifications older than X days        |

***


## Authentication & Authorization

- All **manual and scheduled triggers** are admin- or service-only.

- Cleanup jobs must be run by services with valid `X-Internal-Token`.

- Dry runs allow testing without deletion.

***


## Status Codes Summary

| Code | Meaning               |
| ---- | --------------------- |
| 200  | OK                    |
| 201  | Created               |
| 400  | Invalid request       |
| 401  | Unauthorized          |
| 403  | Forbidden             |
| 500  | Internal Server Error |

***


## Security & Safety Features

- **Dry-run mode** for safe testing before applying deletions.

- **Audit trail** of cleanup jobs stored in Audit Log Service.

- **Rate limiting** and **batch deletion** to avoid overload.

- **Soft delete** optional for donations/logs before hard delete.

***


## Future Enhancements

- Dashboard UI for cleanup scheduling and monitoring.

- Slack/email alerts for failed or successful cleanup jobs.

- Integration with S3/GCS to delete orphaned cloud files.

- Support for TTL policies in PostgreSQL, Redis, and object storage.

***
