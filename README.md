***


# Gratia — Food Donation Management Platform

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)]()

***


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

***


## Project Overview

**Gratia** is a full-stack, cloud-native platform designed to efficiently connect restaurants and catering services with NGOs for the purpose of donating surplus food. It strives to minimize food waste and ensure excess food reaches those who need it most.

Restaurants can list their available surplus food, while NGOs can search, filter by location, and claim donations on a first-come, first-served basis. The platform offers real-time location filtering, interactive maps for easy navigation, and an integrated messaging system to streamline communication between donors and recipients.

Built with modern DevOps practices and a microservices architecture, Gratia ensures reliability, scalability, and maintainability by leveraging containerization, Kubernetes orchestration, infrastructure as code, and CI/CD automation.

***


## Features

- **Restaurant Portal:** Easily list and manage surplus food donations.

- **NGO Portal:** Browse, filter by location, and claim available food donations.

- **Real-time Location Mapping:** Interactive maps to visualize donation and pickup points.

- **Messaging System:** Secure communication channel between restaurants and NGOs.

- **Microservices Architecture:** Modular and scalable services for core functionalities.

- **Cloud-Native Deployment:** Docker containers orchestrated via Kubernetes (EKS).

- **Infrastructure as Code:** Provision cloud resources with Terraform scripts.

- **CI/CD Automation:** Continuous integration and deployment using GitHub Actions.

- **Monitoring & Observability:** Integrated Prometheus and Grafana dashboards.

- **Robust Database:** PostgreSQL for reliable data storage and querying.

***


## Tech Stack

| Component        | Technology / Tool       |
| ---------------- | ----------------------- |
| Frontend         | React.js                |
| Backend          | Go (Golang)             |
| Database         | PostgreSQL              |
| Containerization | Docker                  |
| Orchestration    | Kubernetes (Amazon EKS) |
| Infrastructure   | Terraform               |
| CI/CD            | GitHub Actions          |
| Monitoring       | Prometheus, Grafana     |
| Version Control  | Git, GitHub             |

***


## Architecture

Gratia employs a **microservices architecture**, where discrete services manage user authentication, donation handling, messaging, and notifications. Each service runs inside Docker containers, deployed and managed with Kubernetes on AWS EKS for high availability and scalability.

Infrastructure provisioning is automated via Terraform, enabling seamless environment setup and maintenance. Automated CI/CD pipelines ensure robust testing and smooth deployments. Observability is achieved through Prometheus and Grafana integration, providing monitoring and alerting capabilities.

***


## Installation

### Prerequisites

- [Go](https://golang.org/dl/) (for backend services)

- [Node.js](https://nodejs.org/en/) v14+ (for frontend )

- [Docker](https://www.docker.com/get-started)

- [kubectl](https://kubernetes.io/docs/tasks/tools/)

- [Terraform](https://www.terraform.io/downloads)

- AWS CLI configured with credentials (for EKS cluster)

- PostgreSQL instance (local or managed)


### Setup Locally

1. **Clone the repository**

   ```bash
   git clone https://github.com/AdityaWaradkar/Gratia.git
   cd Gratia
   ```

2. **Install backend dependencies**

   ```bash
   cd backend
   go mod download
   ```

3. **Install frontend dependencies**

   ```bash
   cd ../frontend
   npm install
   ```

4. **Configure environment variables**

   Create `.env` files in both backend and frontend directories with necessary settings (database URLs, JWT secrets, API keys).

5. **Start services**

   - Start PostgreSQL locally or via Docker.

   - Run backend service(s):

     ```bash
     cd backend/services/{service_you_want_to_run}
     go run ./cmd/server/main.go
     ```

   - Run frontend development server:

     ```bash
     cd ../frontend
     npm start
     ```

***


## Usage

- Access the frontend app at `http://localhost:3000` (default).

- Restaurants can register and log in to list surplus food donations.

- NGOs can browse and claim available donations, filtering by location.

- Use the messaging system to coordinate donation pickups efficiently.

***


## API Documentation

Comprehensive RESTful API specifications for all microservices are available inside the `/docs` folder. Swagger/OpenAPI documentation is generated and can be hosted or viewed locally for easy integration and testing.

***


## Contributing

Contributions, bug reports, and feature requests are welcome!

- Open issues to discuss bugs or enhancements.

- Fork the repository and submit pull requests.

- Help improve documentation, tests, or code quality.

Please adhere to the coding style and commit message conventions.

***


## License

This project is licensed under the MIT License — see the [LICENSE]() file for details.

***


## Contact

**Aditya Waradkar**

- LinkedIn: [linkedin.com/in/aditya-waradkar-9a03b92a5](https://www.linkedin.com/in/aditya-waradkar-9a03b92a5/)

- Email: [adityawaradkar1801@gmail.co](mailto:adityawaradkar1801@gmail.com)m

***
