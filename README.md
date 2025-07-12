
# Gratia — Seamless Food Donation Management Platform

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

---

## Table of Contents

- [Project Overview](#project-overview)
- [Features](#features)
- [Tech Stack](#tech-stack)
- [Architecture](#architecture)
- [Installation](#installation)
- [Usage](#usage)
- [API Documentation](#api-documentation)
- [Contributing](#contributing)
- [License](#license)
- [Contact](#contact)

---

## Project Overview

**Gratia** is a full-stack cloud-native web platform designed to facilitate food donations by connecting restaurants and catering services with NGOs. This platform aims to reduce food waste and help distribute surplus food efficiently to those in need.

Restaurants can list their excess food donations, and NGOs can browse, filter by location, and claim these donations on a first-come, first-served basis. The platform includes features like real-time location filtering, navigation via maps, and a messaging system to facilitate communication between donors and recipients.

This project is built using modern DevOps practices including containerization, orchestration, and automated CI/CD pipelines to ensure reliability and scalability.

---

## Features

- **Restaurant Portal:** List and manage surplus food donations.
- **NGO Portal:** Browse and claim food donations with location-based filtering.
- **Real-time Location Mapping:** Integrated maps to visualize donation and pickup locations.
- **Messaging System:** Direct communication channel between restaurants and NGOs.
- **Microservices Architecture:** Modular services for scalability and maintainability.
- **Cloud-Native Deployment:** Docker containers orchestrated with Kubernetes (EKS).
- **Infrastructure as Code:** Terraform scripts for infrastructure provisioning.
- **CI/CD Automation:** GitHub Actions for build, test, and deployment pipelines.
- **Monitoring & Alerts:** Prometheus and Grafana integration for observability.
- **Database:** PostgreSQL for robust data storage and querying.

---

## Tech Stack

| Component             | Technology / Tool           |
|-----------------------|----------------------------|
| Frontend              | React.js (MERN stack)       |
| Backend               | Node.js, Express.js         |
| Database              | PostgreSQL                  |
| Containerization      | Docker                     |
| Orchestration         | Kubernetes (EKS)            |
| Infrastructure        | Terraform                  |
| CI/CD                 | GitHub Actions             |
| Monitoring            | Prometheus, Grafana         |
| Version Control       | Git, GitHub                |

---

## Architecture

Gratia follows a microservices architecture where different functionalities such as user management, donation management, messaging, and notifications are handled by separate services communicating via REST APIs. The services run inside Docker containers deployed on Kubernetes managed by Amazon EKS, ensuring high availability and scalability.

Infrastructure provisioning is automated using Terraform, allowing easy environment setup and management. Continuous integration and deployment pipelines ensure automated testing and smooth updates.

---

## Installation

### Prerequisites

- [Node.js](https://nodejs.org/en/) v14 or above
- [Docker](https://www.docker.com/get-started)
- [kubectl](https://kubernetes.io/docs/tasks/tools/)
- [Terraform](https://www.terraform.io/downloads)
- AWS CLI configured (for EKS)
- PostgreSQL installed or access to managed PostgreSQL instance

### Setup Locally

1. **Clone the repository**

```bash
git clone https://github.com/AdityaWaradkar/Gratia.git
cd Gratia
````

2. **Install backend dependencies**

```bash
cd backend
npm install
```

3. **Install frontend dependencies**

```bash
cd ../frontend
npm install
```

4. **Configure environment variables**

Create `.env` files in backend and frontend folders with necessary configuration (database URL, API keys, etc).

5. **Run locally**

- Start PostgreSQL database (locally or via Docker).

- Run backend server:

```bash
cd backend
npm start
```

- Run frontend:

```bash
cd frontend
npm start
```

***


## Usage

- Open the frontend app in your browser (`http://localhost:3000` by default).

- Restaurants can register/login to list food donations.

- NGOs can browse available donations filtered by location and claim them.

- Use the messaging feature to coordinate pickups.

***


## API Documentation

The backend exposes RESTful APIs to manage users, donations, claims, and messaging. You can find detailed API specs in the `/docs` folder or visit the live API documentation hosted via Swagger/OpenAPI (if set up).

***


## Contributing

Contributions are welcome! Feel free to:

- Report bugs or suggest features via GitHub Issues.

- Fork the repository and submit pull requests.

- Improve documentation or add tests.

Please follow the code style and commit message conventions.

***


## License

This project is licensed under the MIT License — see the [LICENSE]() file for details.

***


## Contact

**Aditya Waradkar**

- GitHub: [AdityaWaradkar](https://github.com/AdityaWaradkar)

- LinkedIn: [linkedin.com/in/aditya-waradkar-9a03b92a5](https://www.linkedin.com/in/aditya-waradkar-9a03b92a5/)

- Email: adityawaradkar1801@gmail.com

***

_Thank you for checking out Gratia! Together, we can reduce food waste and help those in need._


