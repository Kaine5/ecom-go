# E-commerce Go Learning Project

This project is designed as a learning resource for Golang development, focusing on clean architecture principles and domain-driven design. The application consists of two main components: an API service and a worker service for handling asynchronous tasks.

## Core Features

* **Product Management**: List, view, create, update, delete products
* **Order Processing**: Create orders, view order details, list orders
* **Advanced Features**: Caching with Redis, asynchronous processing with RabbitMQ, concurrent handling

## Technologies

* **Backend**: Go (Golang) with Gin web framework
* **Database**: PostgreSQL with GORM
* **Cache**: Redis (optional for advanced features)
* **Message Queue**: RabbitMQ (optional for advanced features)
* **Container Orchestration**: Kubernetes (Minikube for local, EKS for AWS)

## Getting Started

### Prerequisites

* Go 1.23+
* Docker and Docker Compose
* Minikube (optional for local Kubernetes deployment)
* AWS CLI and eksctl (optional for EKS deployment)

### Setup Development Environment

1. Clone the repository
   ```bash
   git clone https://github.com/ducthang310/ecom-go.git
   cd ecom-go
   ```

2. Install dependencies
   ```bash
   go mod download
   go install github.com/air-verse/air@latest
   ```

3. Run the application with Docker Compose

   For Windows
   ```bash
   # Start all services using Docker Compose
   docker-compose -f docker/docker-compose.yml up
   
   # 2. Run API service
   air -c cmd/api/.air.toml


   # 3. Run worker service (If working on worker)
   air -c cmd/worker/.air.toml
   ```
   
   For macOS/Linux
   ```bash
   # Start all services using Docker Compose
   docker-compose -f docker/docker-compose.yml up
   
   # 2. Run API service
   air -c cmd/api/.air.toml --build.cmd "go build -o ./tmp/api ./cmd/api" --build.bin="./tmp/api"

   # 3. Run worker service (If working on worker)
   air -c cmd/worker/.air.toml --build.cmd "go build -o ./tmp/worker ./cmd/worker" --build.bin="./tmp/worker"
   ```

   Alternatively, you can use the Makefile:
   ```bash
   # Start all services using Docker Compose
   make deps-up
   
   # Run API service
   make run-api
   
   # Run worker service
   make run-worker
   ```

### Deployment Options

#### Local Kubernetes (Minikube)
```bash
# Deploy to Minikube
./scripts/deploy-to-minikube.sh
```

#### AWS EKS
```bash
# Deploy to EKS
./scripts/deploy-to-eks.sh
```

## Project Overview

```
ecom-go/
├── cmd/                    # Application entry points
│   ├── api/                # API entry point
│   ├── worker/             # Worker entry point
├── internal/               # Private application code
│   ├── models/             # Domain models/entities
│   ├── repository/         # Data access layer
│   ├── service/            # Business logic layer
│   ├── handler/            # HTTP handlers
│   ├── middleware/         # HTTP middlewares
│   ├── cache/              # Cache layer
│   └── config/             # Configuration
├── pkg/                    # Public shared libraries
├── docker/                 # Docker configurations
├── infra/                  # Kubernetes manifests
├── scripts/                # Helper scripts
└── docs/                   # Documentation
```

For detailed information about the architecture, please see the [architecture documentation](./docs/architecture.md).

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the LICENSE file for details.