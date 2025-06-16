# Microservices Boundaries and Communication Protocols

---

## Microservices Boundaries and Responsibilities

| Microservice            | Boundaries & Responsibilities                                                                                              | Communication Protocol                 |
|------------------------|----------------------------------------------------------------------------------------------------------------------------|--------------------------------------|
| **Auth Service**         | Manages user identity, authentication, authorization, token management, password flows, and RBAC enforcement.                | REST (HTTP/HTTPS)                    |
| **Restaurant Service**   | Handles creation, update, management, and archival of food donation listings by restaurants.                                | REST                               |
| **NGO Service**          | Provides discovery, claiming, and status updates for donations from NGO users; manages NGO-specific feedback mechanisms.     | REST                               |
| **Admin Dashboard**      | Offers user management, content moderation, platform-wide analytics, and audit log access for admin users.                   | REST                               |
| **Notification Service** | Orchestrates all notifications via email, SMS, and in-app alerts, with retries and user preference handling.                | REST + Message Queue (RabbitMQ/Kafka) |
| **Geo Search Service**   | Handles geospatial queries, distance calculations, and reverse geocoding with caching for donation discovery.                | gRPC (for performance & low latency)|
| **Audit Log Service**    | Records immutable logs of critical system events, supports structured queries, integrates with log aggregation tools.       | REST                               |
| **Cleanup Service**      | Periodically manages auto-expiry, stale data pruning, and queue cleanup with health checks.                                 | REST + Scheduled Jobs              |

---

## Communication Protocols

- **REST**  
  Used as the primary synchronous communication protocol for most user-facing and internal microservice requests due to its simplicity and wide support.

- **gRPC**  
  Employed for performance-critical services like Geo Search where low latency, efficient data serialization, and real-time responses are necessary.

- **Asynchronous Messaging (Message Queues)**  
  Implemented via RabbitMQ or Kafka, especially in Notification Service, to handle event-driven workflows and reliable message delivery with retry and failure handling.

---

