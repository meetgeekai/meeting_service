# user_service

This repository provides a production-ready Go microservice template built with [Gin](https://github.com/gin-gonic/gin).  
It is designed with maintainability, observability, and Kubernetes integration in mind.

## Features

- **Dedicated health server** on port `9090`  
  Exposes the following endpoints:  
  - `/health` – generic health check  
  - `/liveness` – liveness probe for Kubernetes  
  - `/readiness` – readiness probe for Kubernetes  

- **Graceful shutdown**  
  Listens for `SIGTERM` signals from Kubernetes and allows up to 30 seconds for ongoing requests to complete before shutting down.

- **Global rate limiting**  
  Provides request rate limiting to protect the service from overload.

- **Repository pattern**  
  Data access is abstracted behind repositories, ensuring a clean separation between business logic and persistence.

- **Layered architecture**  
  - **Handlers**: HTTP endpoints exposed via Gin  
  - **Services**: Business logic, easily testable with mocked dependencies  
  - **Repositories**: Database access and persistence logic  

- **MySQL integration**  
  Includes a MySQL connector with connection pooling for efficient database usage.

## Usage

When creating a new service from this template:

1. **Rename the service**  
   Replace all instances of `go-rest-microservice-template` with your desired service name.

2. **Install development tools**  
   This project relies on the following CLI tools:  
   - [sqlc](https://docs.sqlc.dev/en/latest/overview/install.html)  
   - [golang-migrate](https://github.com/golang-migrate/migrate)  

3. **Obtain environment files**  
   Request the `.env` and `docker-compose.yml` files needed for local development.

## Local Development

Once the environment files are in place:

```bash
# Run the service locally
docker-compose up -d
```
