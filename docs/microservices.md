## **Microservices**

---

### 1. **Auth Service**

- **Function**: User registration, login, password reset

- **Responsibility**:

  - Issue & validate **JWT tokens**

  - Manage roles (**Donor / NGO / Admin**)

  - Secure inter-service communication

***


### 2. **User Service**

- **Function**: Manage donor, NGO, and admin profiles

- **Responsibility**:

  - Store profile details

  - Update user info

  - Handle NGO verification by Admin

***


### 3. **Food Listing Service**

- **Function**: Donors post leftover food

- **Responsibility**:

  - Create, update, delete listings

  - Add food details (type, quantity, expiry, location)

  - Set listing availability (open, claimed, expired)

***


### 4. **Claim Service**

- **Function**: NGOs claim food

- **Responsibility**:

  - Handle claim requests

  - Assign claims to NGOs

  - Update claim status (pending, picked, delivered)

***


### 5. **Notification Service**

- **Function**: Alerts & communication

- **Responsibility**:

  - Notify NGOs when new food is listed nearby

  - Notify donors when food is claimed or delivered

  - Send emails/SMS/push messages

***


### 6. **Admin Service**

- **Function**: Platform moderation

- **Responsibility**:

  - Manage users (approve/reject NGOs, ban users)

  - Handle abuse reports

  - View system-wide stats

***


### 7. **Analytics Service (Optional, can come later)**

- **Function**: Track platform impact

- **Responsibility**:

  - Generate reports (e.g., total meals donated, NGO performance)

  - Provide dashboards/insights

***


### 8. **Gateway / API Gateway**

- **Function**: Single entry point for clients

- **Responsibility**:

  - Route requests to correct microservice

  - Handle rate-limiting, load balancing, authentication middleware


***