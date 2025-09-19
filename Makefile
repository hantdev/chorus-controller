# Variables
BINARY_NAME=controller
BUILD_DIR=build
MAIN_PATH=./cmd/main.go
MODULE_NAME=github.com/hantdev/chorus-controller

# Go build flags
LDFLAGS=-ldflags "-X main.Version=$(shell git describe --tags --always --dirty) -X main.BuildTime=$(shell date -u '+%Y-%m-%d_%H:%M:%S')"

# Default target
.PHONY: all
all: clean build

# Database migration targets
.PHONY: migrate-gen
migrate-gen:
	@if [ -z "$(NAME)" ]; then \
		echo "Usage: make migrate-gen NAME=<migration_name>"; \
		echo "Example: make migrate-gen NAME=add_user_table"; \
		exit 1; \
	fi
	@echo "Generating migration: $(NAME)"
	@echo "Step 1: Generating schema from GORM models..."
	@go run cmd/schema/main.go schema.sql
	@echo "Step 2: Creating migration with Atlas..."
	@GOWORK=off atlas migrate diff "$(NAME)" \
		--dir "file://migrations" \
		--to "file://schema.sql" \
		--dev-url "docker://postgres/16/dev?search_path=public"
	@echo "Step 3: Cleaning up..."
	@rm -f schema.sql
	@echo "✅ Migration generated successfully!"

.PHONY: migrate-apply
migrate-apply:
	@echo "Applying migrations..."
	@GOWORK=off atlas migrate apply --env local

.PHONY: migrate-status
migrate-status:
	@echo "Migration status:"
	@GOWORK=off atlas migrate status --env local

.PHONY: migrate-hash
migrate-hash:
	@echo "Updating migration hash..."
	@GOWORK=off atlas migrate hash --dir "file://migrations"

.PHONY: migrate-reset
migrate-reset:
	@echo "⚠️  WARNING: This will drop all tables and recreate them!"
	@read -p "Are you sure? (y/N): " confirm && [ "$$confirm" = "y" ]
	@GOWORK=off atlas migrate apply --env local --baseline 0

# Build the application
.PHONY: build
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

# Build for specific platforms
.PHONY: build-linux
build-linux:
	@echo "Building for Linux..."
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 $(MAIN_PATH)

.PHONY: build-darwin
build-darwin:
	@echo "Building for macOS..."
	@mkdir -p $(BUILD_DIR)
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 $(MAIN_PATH)

.PHONY: build-windows
build-windows:
	@echo "Building for Windows..."
	@mkdir -p $(BUILD_DIR)
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe $(MAIN_PATH)

# Build for all platforms
.PHONY: build-all
build-all: build-linux build-darwin build-windows

# Run the application
.PHONY: run
run:
	@echo "Running $(BINARY_NAME)..."
	go run $(MAIN_PATH)

# Run with hot reload (requires air)
.PHONY: dev
dev:
	@echo "Running in development mode with hot reload..."
	@if command -v air > /dev/null; then \
		air; \
	else \
		echo "Air not found. Installing air..."; \
		go install github.com/cosmtrek/air@latest; \
		air; \
	fi

# Test the application
.PHONY: test
test:
	@echo "Running tests..."
	go test -v ./...

# Test with coverage
.PHONY: test-coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run tests in watch mode
.PHONY: test-watch
test-watch:
	@echo "Running tests in watch mode..."
	@if command -v gotestsum > /dev/null; then \
		gotestsum --format=short-verbose --watch; \
	else \
		echo "gotestsum not found. Installing gotestsum..."; \
		go install gotest.tools/gotestsum@latest; \
		gotestsum --format=short-verbose --watch; \
	fi

# Benchmark tests
.PHONY: benchmark
benchmark:
	@echo "Running benchmarks..."
	go test -bench=. -benchmem ./...

# Lint the code
.PHONY: lint
lint:
	@echo "Running linter..."
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not found. Installing golangci-lint..."; \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; \
		golangci-lint run; \
	fi

# Format the code
.PHONY: fmt
fmt:
	@echo "Formatting code..."
	go fmt ./...
	@echo "Checking for goimports..."
	@which goimports > /dev/null 2>&1 || (echo "Installing goimports..." && go install golang.org/x/tools/cmd/goimports@latest)
	@echo "Running goimports..."
	@$(shell go env GOPATH)/bin/goimports -w .

# Vet the code
.PHONY: vet
vet:
	@echo "Running go vet..."
	go vet ./...

# Check for security vulnerabilities
.PHONY: security
security:
	@echo "Checking for security vulnerabilities..."
	@if command -v gosec > /dev/null; then \
		echo "gosec found but may not work with Go modules outside GOPATH"; \
		echo "Consider using govulncheck instead for Go module projects"; \
		echo "Installing govulncheck..."; \
		go install golang.org/x/vuln/cmd/govulncheck@latest; \
		echo "Running govulncheck..."; \
		$(shell go env GOPATH)/bin/govulncheck ./...; \
	else \
		echo "Installing govulncheck (modern Go security tool)..."; \
		go install golang.org/x/vuln/cmd/govulncheck@latest; \
		echo "govulncheck installed. Running security check..."; \
		$(shell go env GOPATH)/bin/govulncheck ./...; \
	fi

# Generate mock files
.PHONY: mocks
mocks:
	@echo "Generating mocks..."
	@if command -v mockgen > /dev/null; then \
		mockgen -source=internal/worker/client.go -destination=internal/worker/mocks/client_mock.go; \
	else \
		echo "mockgen not found. Installing mockgen..."; \
		go install github.com/golang/mock/mockgen@latest; \
		mockgen -source=internal/worker/client.go -destination=internal/worker/mocks/client_mock.go; \
	fi

# Clean build artifacts
.PHONY: clean
clean:
	@echo "Cleaning build artifacts..."
	rm -rf $(BUILD_DIR)
	rm -f coverage.out coverage.html
	go clean -cache -testcache

# Install dependencies
.PHONY: deps
deps:
	@echo "Installing dependencies..."
	go mod download
	go mod tidy

# Update dependencies
.PHONY: deps-update
deps-update:
	@echo "Updating dependencies..."
	go get -u ./...
	go mod tidy

# Show dependency tree
.PHONY: deps-tree
deps-tree:
	@echo "Dependency tree:"
	go mod graph

# Docker targets
.PHONY: docker-build
docker-build:
	@echo "Building Docker image..."
	docker build -t $(BINARY_NAME):latest .

.PHONY: docker-run
docker-run:
	@echo "Running Docker container..."
	docker run -p 8081:8081 $(BINARY_NAME):latest

# Migration targets
.PHONY: migrate-up
migrate-up: ## Apply all migrations
	@echo "Applying migrations..."
	@if command -v atlas > /dev/null; then \
		atlas migrate apply --env local; \
	else \
		echo "Atlas CLI not found. Installing Atlas..."; \
		curl -sSf https://atlasgo.sh | sh; \
		atlas migrate apply --env local; \
	fi

.PHONY: migrate-down
migrate-down: ## Rollback last migration
	@echo "Rolling back last migration..."
	@if command -v atlas > /dev/null; then \
		atlas migrate down --env local; \
	else \
		echo "Atlas CLI not found. Please install Atlas first."; \
		exit 1; \
	fi

.PHONY: migrate-status
migrate-status: ## Show migration status
	@echo "Migration status:"
	@if command -v atlas > /dev/null; then \
		atlas migrate status --env local; \
	else \
		echo "Atlas CLI not found. Please install Atlas first."; \
		exit 1; \
	fi

.PHONY: migrate-new
migrate-new: ## Create new migration file
	@echo "Creating new migration..."
	@if command -v atlas > /dev/null; then \
		read -p "Enter migration name: " name; \
		atlas migrate new $$name --env local; \
	else \
		echo "Atlas CLI not found. Please install Atlas first."; \
		exit 1; \
	fi

.PHONY: migrate-baseline
migrate-baseline: ## Baseline existing database
	@echo "Baselining existing database..."
	@if command -v atlas > /dev/null; then \
		atlas migrate baseline --env local; \
	else \
		echo "Atlas CLI not found. Please install Atlas first."; \
		exit 1; \
	fi

.PHONY: migrate-hash
migrate-hash: ## Generate migration checksums
	@echo "Generating migration checksums..."
	@if command -v atlas > /dev/null; then \
		atlas migrate hash --env local; \
	else \
		echo "Atlas CLI not found. Please install Atlas first."; \
		exit 1; \
	fi

# Help target
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  build          - Build the application"
	@echo "  build-linux    - Build for Linux"
	@echo "  build-darwin   - Build for macOS"
	@echo "  build-windows  - Build for Windows"
	@echo "  build-all      - Build for all platforms"
	@echo "  run            - Run the application"
	@echo "  dev            - Run with hot reload (requires air)"
	@echo "  test           - Run tests"
	@echo "  test-coverage  - Run tests with coverage"
	@echo "  test-watch     - Run tests in watch mode"
	@echo "  benchmark      - Run benchmarks"
	@echo "  lint           - Run linter"
	@echo "  fmt            - Format code"
	@echo "  vet            - Run go vet"
	@echo "  security       - Check for security vulnerabilities"
	@echo "  mocks          - Generate mock files"
	@echo "  migrate-up     - Apply all migrations"
	@echo "  migrate-down   - Rollback last migration"
	@echo "  migrate-status - Show migration status"
	@echo "  migrate-new    - Create new migration file"
	@echo "  migrate-baseline - Baseline existing database"
	@echo "  migrate-hash   - Generate migration checksums"
	@echo "  swagger        - Generate Swagger documentation"
	@echo "  swagger-clean  - Clean generated Swagger documentation"
	@echo "  swagger-serve  - Serve Swagger documentation"
	@echo "  clean          - Clean build artifacts"
	@echo "  deps           - Install dependencies"
	@echo "  deps-update    - Update dependencies"
	@echo "  deps-tree      - Show dependency tree"
	@echo "  docker-build   - Build Docker image"
	@echo "  docker-run     - Run Docker container"
	@echo "  help           - Show this help message"

.PHONY: swagger
swagger: ## Generate Swagger documentation
	@echo "Generating Swagger documentation..."
	@if command -v swag > /dev/null; then \
		swag init -g ./cmd/main.go -o ./docs; \
	else \
		echo "swag not found. Installing swag..."; \
		go install github.com/swaggo/swag/cmd/swag@latest; \
		$(shell go env GOPATH)/bin/swag init -g ./cmd/main.go -o ./docs; \
	fi
	@echo "Swagger docs generated in docs/ directory"

.PHONY: swagger-clean
swagger-clean: ## Clean generated Swagger documentation
	@echo "Cleaning Swagger documentation..."
	@rm -rf docs/
	@echo "Swagger docs cleaned"

.PHONY: swagger-serve
swagger-serve: ## Serve Swagger documentation (requires swagger-ui)
	@echo "Serving Swagger documentation..."
	@if command -v swagger > /dev/null; then \
		swagger serve -F swagger docs/swagger.json; \
	else \
		echo "swagger not found. Installing swagger..."; \
		go install github.com/go-swagger/go-swagger/cmd/swagger@latest; \
		swagger serve -F swagger docs/swagger.json; \
	fi