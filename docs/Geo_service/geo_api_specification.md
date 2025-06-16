***


# Geo Search Service API Specification

***


## Overview

The Geo Search Service enables efficient location-based querying of donations, users, and services within the Gratia platform. It supports proximity-based search, reverse geocoding, and geo-tagging using GPS coordinates. This service enhances location-awareness for NGOs, restaurants, and admins.

***


## Base URL

    /api/v1/geo

***


## Endpoints

***


### 1. Find Nearby Donations

**GET** `/donations/nearby`

Returns a list of available donations sorted by proximity to the provided location.


#### Headers

    Authorization: Bearer <access_token>


#### Query Parameters

| Parameter   | Type    | Description                               |
| ----------- | ------- | ----------------------------------------- |
| `lat`       | float   | Latitude of the search center (required)  |
| `lon`       | float   | Longitude of the search center (required) |
| `radius_km` | float   | Search radius in kilometers (default: 10) |
| `limit`     | integer | Number of results to return (default: 10) |


#### Response (200 OK)

```json
{
  "donations": [
    {
      "donation_id": "uuid-generated",
      "food_description": "50 sandwiches",
      "location": {
        "address": "123 Street, Mumbai",
        "latitude": 19.076,
        "longitude": 72.8777
      },
      "distance_km": 2.4
    }
  ]
}
```

***


### 2. Reverse Geocode

**GET** `/reverse-geocode`

Converts GPS coordinates to a human-readable address.


#### Query Parameters

| Parameter | Type  | Description          |
| --------- | ----- | -------------------- |
| `lat`     | float | Latitude (required)  |
| `lon`     | float | Longitude (required) |


#### Response (200 OK)

```json
{
  "address": "123 Charity Street, Bandra West, Mumbai, MH, India"
}
```

***


### 3. Geotag Donation Location

**POST** `/donations/{donation_id}/location`

Assigns or updates geolocation data for a donation.


#### Headers

    Authorization: Bearer <access_token>


#### Path Parameters

| Parameter     | Type   | Description                       |
| ------------- | ------ | --------------------------------- |
| `donation_id` | string | Unique identifier of the donation |


#### Request Body

```json
{
  "latitude": 19.076,
  "longitude": 72.8777,
  "address": "123 Street, Mumbai"
}
```


#### Response (200 OK)

```json
{
  "message": "Donation location updated successfully."
}
```

***


### 4. Get Service Coverage Area (Admin Use)

**GET** `/coverage`

Used by admins to visualize the service coverage of donations or users.


#### Query Parameters

| Parameter | Type   | Description                             |
| --------- | ------ | --------------------------------------- |
| `type`    | string | One of: `ngo`, `restaurant`, `donation` |
| `city`    | string | Filter by city (optional)               |


#### Response (200 OK)

```json
{
  "locations": [
    {
      "id": "uuid",
      "type": "ngo",
      "name": "Helping Hands",
      "latitude": 19.11,
      "longitude": 72.9
    }
  ]
}
```

***


## Geo Distance Calculation

All distance calculations are based on the **Haversine formula**, assuming Earth’s radius as 6371 km.

***


## Authentication & Authorization

- Most endpoints require a valid JWT token.

- Admin-only access applies to `/coverage`.

- Internal services may use service tokens or mTLS for geotagging.

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


## Security & Performance

- Use of caching (e.g., Redis) for repeated reverse geocode lookups.

- Rate limiting for public-facing geolocation endpoints.

- Validates coordinate ranges (`-90 ≤ lat ≤ 90`, `-180 ≤ lon ≤ 180`).

***


## Future Enhancements

- **Clustering** for heatmaps in admin dashboards.

- **Offline map support** with caching.

- **Geo-fencing** for pickup and delivery validation.

- **ElasticSearch or PostGIS** integration for scalable geo queries.

***
