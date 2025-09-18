# Chorus Controller

A control plane for Chorus Worker that manages S3 replication jobs.

## Project Structure

```
chorus-controller/
├── cmd/                    # Application entry points
│   └── main.go           # Main application entry point
├── internal/              # Private application code
│   ├── config/           # Configuration management
│   │   └── config.go     # Configuration structs and loading
│   ├── domain/           # Domain models and interfaces
│   │   ├── models.go     # Domain models and DTOs
│   │   └── interfaces.go # Service interfaces and contracts
│   ├── service/          # Business logic layer
│   │   ├── replication.go # Replication business logic
│   │   └── storage.go    # Storage business logic
│   ├── repository/       # Data access layer
│   │   └── worker.go     # Worker gRPC client implementation
│   ├── handler/          # HTTP handlers
│   │   ├── health.go     # Health check endpoints
│   │   ├── replication.go # Replication management endpoints
│   │   ├── storage.go    # Storage management endpoints
│   │   └── middleware.go # HTTP middleware and error handling
│   ├── server/           # HTTP server configuration
│   │   └── server.go     # Server setup and routing
│   └── errors/           # Error handling
│       └── errors.go     # Error types and utilities
├── docs/                 # API documentation
│   ├── docs.go          # Swagger documentation setup
│   ├── swagger.json     # OpenAPI specification
│   └── swagger.yaml     # OpenAPI specification (YAML)
├── build/                # Build artifacts
├── go.mod               # Go module dependencies
├── go.sum               # Go module checksums
├── Makefile             # Build and development commands
└── README.md            
```

## Key Components

- **`cmd/main.go`**: Application entry point that initializes all layers using dependency injection
- **`internal/config/`**: Configuration management for worker gRPC address and HTTP port
- **`internal/domain/`**: Domain models, DTOs, and service interfaces defining business contracts
- **`internal/service/`**: Business logic layer implementing domain services for replication and storage operations
- **`internal/repository/`**: Data access layer handling gRPC communication with Chorus Worker service
- **`internal/handler/`**: HTTP handlers organized by feature (health, storage, replication) with middleware
- **`internal/server/`**: HTTP server configuration and routing setup
- **`internal/errors/`**: Centralized error handling with custom error types and HTTP status mapping
- **`docs/`**: Auto-generated Swagger/OpenAPI documentation

## Architecture

The project follows **Clean Architecture** principles with clear separation of concerns:

### **Layers**
1. **Domain Layer** (`internal/domain/`): Contains business models and interfaces
2. **Service Layer** (`internal/service/`): Implements business logic
3. **Repository Layer** (`internal/repository/`): Handles data access
4. **Handler Layer** (`internal/handler/`): Manages HTTP requests/responses
5. **Server Layer** (`internal/server/`): Configures and runs the HTTP server

### **Benefits**
- **Testability**: Each layer can be tested independently
- **Maintainability**: Clear separation makes code easier to understand and modify
- **Scalability**: New features can be added without affecting existing code
- **Reusability**: Services can be reused across different handlers
- **Error Consistency**: Centralized error handling ensures consistent API responses

## API Endpoints

- `GET /health` - Health check
- `GET /storages` - List all configured storages
- `GET /buckets` - List buckets available for replication
- `POST /replications` - Create new replication job
- `GET /replications` - List all replication jobs
- `POST /replications/pause` - Pause replication job
- `POST /replications/resume` - Resume replication job
- `DELETE /replications` - Delete replication job
- `POST /replications/switch/zero-downtime` - Switch buckets without downtime

## Development

### Prerequisites
- Go 1.25+
- Access to Chorus Worker gRPC service

### Building
```bash
make build
# or
go build ./cmd/main.go
```

### Running
```bash
make run
# or
go run ./cmd/main.go
```

### API Documentation
Swagger UI available at `/swagger/*` when running the server.
