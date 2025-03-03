# Makefile for E-commerce Go project

# Variables
BIN_DIR = bin
CMD_DIR = cmd
DOCKER_COMPOSE = docker/docker-compose.yml

# Build binary executables
build:
	mkdir -p $(BIN_DIR)
	go build -o $(BIN_DIR)/api $(CMD_DIR)/api/main.go
	go build -o $(BIN_DIR)/worker $(CMD_DIR)/worker/main.go

# Run API service
run-api:
	go run $(CMD_DIR)/api/main.go

# Run worker service
run-worker:
	go run $(CMD_DIR)/worker/main.go

# Start Docker containers
docker-up:
	docker compose -f $(DOCKER_COMPOSE) up -d

# Stop Docker containers
docker-down:
	docker compose -f $(DOCKER_COMPOSE) down

# Format code
fmt:
	go fmt ./...

# Install dependencies
deps:
	go mod download

# Clean build artifacts
clean:
	rm -rf $(BIN_DIR)

# Help target
help:
	@echo "Available targets:"
	@echo "  build          - Build the project binaries"
	@echo "  run-api        - Run the API service"
	@echo "  run-worker     - Run the worker service"
	@echo "  docker-up      - Start dependency containers"
	@echo "  docker-down    - Stop dependency containers"
	@echo "  fmt            - Format code"
	@echo "  deps           - Install dependencies"
	@echo "  clean          - Clean build artifacts"

.PHONY: build run-api run-worker docker-up docker-down fmt deps clean help