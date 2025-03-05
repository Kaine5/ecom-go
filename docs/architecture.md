# Architecture Documentation

## Overview

The E-commerce Go project implements a simplified version of clean architecture principles and domain-driven design concepts. This structure is designed as a learning scaffold rather than a production-ready implementation. The document explains the basic architecture, its layers.

This project serves as an introduction to these architectural concepts. For more comprehensive implementations or production applications, we recommend exploring the references provided below in greater depth.

## References and Further Reading
- [The Clean Architecture blog post by Robert C. Martin](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Domain-Driven Design Quickly](https://www.infoq.com/minibooks/domain-driven-design-quickly/)
- [Standard Go Project Layout](https://github.com/golang-standards/project-layout)

## Project Structure

```
ecommerce-api/
├── cmd/                    # Application entry points
│   ├── api/                # API entry point
│   │   ├── .air.toml       # Configuration for hot-reload with air server
│   │   └── main.go
│   ├── worker/             # Worker entry point
│   │   ├── .air.toml       # Configuration for hot-reload with air server
│   │   └── main.go
├── internal/               # Private application code
│   ├── models/             # Domain models/entities
│   │   ├── user.go         # User model
│   │   ├── product.go      # Product model
│   │   └── order.go        # Order model
│   ├── repository/         # Data access layer
│   │   ├── user_interface.go  # User repository interface
│   │   ├── user_repo.go     # User repository implementation
│   │   ├── db.go           # Database connection setup
│   │   └── factory.go      # Repository factory
│   ├── service/            # Business logic layer
│   │   ├── user.go         # User service
│   │   ├── product.go      # Product service
│   │   └── order.go        # Order service
│   ├── handler/            # HTTP handlers
│   │   ├── user.go         # User handler
│   │   ├── product.go      # Product handler
│   │   └── order.go        # Order handler
│   ├── middleware/         # HTTP middlewares
│   │   ├── auth.go         # Authentication middleware
│   │   └── logging.go      # Logging middleware
│   ├── cache/              # Cache layer (optional)
│   │   └── redis.go        # Redis implementations
│   └── config/             # Configuration
│       └── config.go       # Configuration management
├── pkg/                    # Public libraries that could be used by other projects
│   ├── logger/             # Standardized logging
│   │   └── logger.go
│   ├── errors/             # Error handling
│   │   ├── base.go         # Error interface
│   │   └── errors.go       # Error implementations
│   └── http/               # HTTP utilities
│       └── response.go     # Response formatting
├── docker/                # Docker configuration
│   ├── api/
│   │   └── Dockerfile     # API service Dockerfile
│   ├── worker/
│   │   └── Dockerfile     # Worker service Dockerfile
│   └── docker-compose.yml # Docker Compose for local development dependencies (PostgreSQL, RabitMQ, Redis)
├── infra/                 # Infrastructure and deployment
│   ├── k8s/               # Common Kubernetes resources
│   │   ├── db.yaml
│   │   ├── redis.yaml
│   │   ├── rabbitmq.yaml
│   │   └── configmap.yaml
│   ├── minikube/          # Minikube-specific resources
│   │   ├── api.yaml
│   │   └── worker.yaml
│   └── eks/               # EKS-specific resources
│       ├── api.yaml
│       ├── worker.yaml
│       └── autoscaling.yaml
├── scripts/               # Utility scripts
│   ├── deploy-minikube.sh # Deploy to Minikube
│   └── deploy-eks.sh      # Deploy to AWS EKS
├── config.yaml            # Application configuration
├── Makefile               # Build and development commands
├── docs/                  # Documentation
│   ├── api.md             # API documentation
│   └── architecture.md    # Architecture documentation
├── go.mod                 # Go modules file
└── README.md              # Project documentation
```

## Layer Descriptions

### Models Layer (`internal/models`)

Contains domain models and business logic that is intrinsic to the entities.

### Repository Layer (`internal/repository`)

Handles data persistence and retrieval, abstracting the database interactions.

### Service Layer (`internal/service`)

Implements the business logic of the application, orchestrating calls to repositories.

Example: The `UserService` handles operations like user registration, ensuring business rules are followed (e.g., checking for duplicate emails).

### Handler Layer (`internal/handler`)

Manages HTTP requests and responses, converting between HTTP and domain objects.

### Middleware Layer (`internal/middleware`)

Provides cross-cutting concerns that can be applied to HTTP requests.

### Config Layer (`internal/config`)

Manages application configuration from various sources.

### Utility Packages (`pkg`)

Shared utilities used across the application.


### Infrastructure Layer (`infra/`)

Contains Kubernetes and deployment configurations for different environments.

**Contains:**
- **Common Kubernetes Resources**: Shared configurations
  - Database, Redis, RabbitMQ configurations
  - ConfigMaps for application settings

- **Environment-Specific Resources**: Customized for different platforms
  - Minikube configurations for local development
  - EKS configurations for AWS cloud deployment

### Docker Configuration (`docker/`)

Contains Dockerfiles for building container images.

**Contains:**
- **API Dockerfile**: Multi-stage build for the API service
- **Worker Dockerfile**: Multi-stage build for the worker service
- **Docker-compose**: For local development

### Scripts (`scripts/`)

Utility scripts for deploying to EKS and Minikube.

## Database Management

The application uses both automatic migrations through GORM and explicit SQL migrations:

1. **GORM AutoMigrate**: For development and rapid prototyping, the application can automatically create and update database tables based on model definitions.

2. **SQL Migrations**: For more controlled schema management, explicit SQL migration files can be used with the migration tool.

## Request Flow

1. HTTP request is received by the Gin router
2. Request passes through middleware (logging, auth, etc.)
3. Handler parses and validates the request
4. Handler calls the appropriate service method
5. Service implements business logic, calling repositories as needed
6. Repository performs database operations
7. Results flow back up through the layers
8. Handler formats the response and sends it back

## Error Handling

The application uses a standardized error handling approach with typed errors:

```go
// Base error interface
type BaseError interface {
    ToResponseError() *ResponseError
    Error() string
    Type() string
}

// Standard response error structure
type ResponseError struct {
    Type       string      `json:"type"`
    Errors     []ErrorItem `json:"errors"`
    StatusCode int         `json:"status_code"`
}

// Error implementations include:
// - ServerError
// - NotFoundError
// - BadRequestError
// - UnauthorizedError
// - ForbiddenError
```

This approach allows for:
- Consistent error responses across the API
- Type-based error handling in the application code
- Proper HTTP status codes based on error types

## Deployment Architecture

### Minikube Development

For local Kubernetes development:
- Single replicas of API and worker services
- LoadBalancer service (via Minikube tunnel)
- Development-friendly configurations
- Deploy using `./scripts/deploy-to-minikube.sh`

### Cloud Deployment (AWS EKS)

For cloud deployment, the application can be deployed to AWS EKS:

- Multiple replicas of API and worker services
- Load balancing with AWS Network Load Balancer
- Auto-scaling based on CPU utilization
- Containerized PostgreSQL (in a real production environment, consider AWS RDS)
- Containerized Redis (in a real production environment, consider ElastiCache)
- Containerized RabbitMQ (in a real production environment, consider Amazon MQ)
- Deploy using `./scripts/deploy-to-eks.sh`
