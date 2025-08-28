# Chorus Controller

A control plane for Chorus Worker that manages S3 replication jobs.
## Features

- Storage management (add, update, delete)
- Bucket listing for replication
- Replication job management (create, pause, resume, delete)
- Zero-downtime bucket switching
- RESTful API with comprehensive Swagger documentation
- Interactive API documentation via Swagger UI

## Quick Start

```bash
# Build
go build -o build/controller ./cmd

# Run
WORKER_GRPC_ADDR=localhost:9670 HTTP_PORT=8081 ./build/controller
```
