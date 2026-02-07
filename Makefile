.PHONY: help build run clean test migrate migrate-down dev docker-build docker-run install

# Variables
BINARY_NAME=vesper
MAIN_PATH=./cmd/server
DATA_DIR=./data
DB_FILE=$(DATA_DIR)/tasks.db

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-15s %s\n", $$1, $$2}'

build: ## Build the application binary
	@echo "Building $(BINARY_NAME)..."
	@go build -o $(BINARY_NAME) $(MAIN_PATH)
	@echo "Build complete: ./$(BINARY_NAME)"

run: build ## Build and run the application
	@echo "Starting $(BINARY_NAME)..."
	@./$(BINARY_NAME)

clean: ## Clean build artifacts
	@echo "Cleaning..."
	@rm -f $(BINARY_NAME)
	@rm -rf tmp/
	@echo "Clean complete"

test: ## Run tests
	@echo "Running tests..."
	@go test -v ./...

migrate: ## Run database migrations
	@echo "Running database migrations..."
	@mkdir -p $(DATA_DIR)
	@go run ./internal/database/migrate/migrate.go up
	@echo "Migrations complete"

migrate-down: ## Rollback database migrations
	@echo "Rolling back database migrations..."
	@go run ./internal/database/migrate/migrate.go down
	@echo "Rollback complete"

dev: ## Run in development mode with hot reload (requires air)
	@echo "Starting development server with air..."
	@air

docker-build: ## Build Docker image
	@echo "Building Docker image..."
	@docker build -t vesper:latest .
	@echo "Docker image built: vesper:latest"

docker-run: docker-build ## Build and run Docker container
	@echo "Running Docker container..."
	@docker run -p 8080:8080 -v $(PWD)/data:/data vesper:latest

install: ## Install dependencies
	@echo "Installing dependencies..."
	@go mod download
	@go mod verify
	@echo "Dependencies installed"

fmt: ## Format Go code
	@echo "Formatting code..."
	@go fmt ./...
	@echo "Formatting complete"

vet: ## Run go vet
	@echo "Running go vet..."
	@go vet ./...
	@echo "Vet complete"

lint: fmt vet ## Run all linters

setup: install migrate ## Complete project setup
	@echo "Project setup complete!"
	@echo "Run 'make run' to start the server"
