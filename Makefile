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


# Development with hot reload
# 1. Start dependencies (PostgreSQL, Redis, RabbitMQ) in Docker
deps-up:
	docker compose -f $(DOCKER_COMPOSE) up -d

# 2. Run API service
run-api:
	air -c $(CMD_DIR)/api/.air.toml

# 3. Run worker service (If working on worker)
run-worker:
	air -c $(CMD_DIR)/worker/.air.toml

# Stop dependencies (PostgreSQL, Redis, RabbitMQ)
deps-up:
	docker compose -f $(DOCKER_COMPOSE) down -d

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
	@echo "  deps-up        - Start all dependencies"
	@echo "  run-api        - Run the API service"
	@echo "  run-worker     - Run the worker service"
	@echo "  fmt            - Format code"
	@echo "  deps           - Install dependencies"
	@echo "  clean          - Clean build artifacts"

.PHONY: build deps-up run-api run-worker fmt deps clean help