### Overview
`chorus-controller` is the control plane that manages S3 replication jobs run by `chorus-worker`. This guide explains how to set up a development environment, run and build the application, and manage database migrations with Atlas.

### Prerequisites
- Go 1.25+
- Docker/Docker Compose (recommended for Postgres)
- Atlas CLI (for migrations)
- A running `chorus-worker` gRPC service (default `localhost:9670`) or an accessible endpoint

### Configuration (Environment Variables)
The application reads configuration from environment variables:
- `WORKER_GRPC_ADDR` (default: `localhost:9670`) — gRPC address of `chorus-worker`.
- `POSTGRES_DSN` (default: `postgres://postgres:postgres@localhost:5432/chorus?sslmode=disable`) — Postgres connection DSN.
- `HTTP_PORT` (default: `8081`) — HTTP server port.

Quick set up:
```bash
export WORKER_GRPC_ADDR=localhost:9670
export POSTGRES_DSN="postgres://postgres:postgres@localhost:5432/chorus?sslmode=disable"
export HTTP_PORT=8081
```

### Start Postgres for Development
Use the repo’s Docker Compose or run a standalone Postgres container.

- Option 1: Repo Docker Compose
```bash
# From /Users/hant/Github/chorus-controller
docker compose up -d postgres
```
- Option 2: Standalone container
```bash
docker run -d --name chorus-pg \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=chorus \
  -p 5432:5432 \
  postgres:16
```

### Migrations (Atlas)
This project uses Atlas instead of GORM AutoMigrate. SQL files are in `migrations/`. Atlas configuration is in `atlas.hcl`.

Install Atlas CLI (if not installed):
```bash
curl -sSf https://atlasgo.sh | sh
```

Common commands (via Makefile):
- Apply all migrations: `make migrate-up`
- Roll back last migration: `make migrate-down`
- Check status: `make migrate-status`
- Create new migration: `make migrate-new` (prompts for name)
- Baseline an existing DB: `make migrate-baseline`
- Generate checksums: `make migrate-hash`

First-time setup:
```bash
# 1) Ensure Postgres is running
# 2) Apply migrations
make migrate-up
```

When changing schema:
```bash
make migrate-new      # enter a migration name
# edit the generated SQL in migrations/
make migrate-up
```

### Run the Application
Three common ways:

- Run with Go:
```bash
make run
# or
go run ./cmd
```

- Docker Compose (controller + postgres):
```bash
docker compose up --build -d
# Controller listens on :8081 by default, Swagger at /swagger/*
```

- Standalone Docker:
```bash
make docker-build
make migrate-up  # apply migrations to local DB before running the container
make docker-run  # uses default env mapping; adjust if needed
```

### Build Binary
```bash
make build
# produces bin/chorus-controller

# or
go build -o bin/chorus-controller ./cmd
```

### Testing
```bash
make test
# or
go test ./...
```

### API and Swagger
- When running, Swagger UI is available at: `http://localhost:8081/swagger/index.html`.
- Key endpoints:
  - `GET /health` — health check
  - `GET /storages` — list storages
  - `GET /buckets` — list buckets
  - `POST /storages` — create storage (persisted in DB)
  - `GET /storages/db` — list DB storages
  - `POST /replications` — create replication job
  - `GET /replications` — list replication jobs
  - `POST /replications/pause` — pause job
  - `POST /replications/resume` — resume job
  - `DELETE /replications` — delete job
  - `POST /replications/switch/zero-downtime` — zero-downtime bucket switch

### Integration with chorus-worker
`chorus-controller` calls the `chorus-worker` gRPC API via `WORKER_GRPC_ADDR` (default `localhost:9670`). Ensure the worker is running and reachable.

- If developing `chorus-worker` locally, run it per its repo guide and configure:
```bash
export WORKER_GRPC_ADDR=localhost:9670
```

### Quick Dev Workflow
```bash
# 1) Start Postgres
docker compose up -d postgres

# 2) Apply migrations
make migrate-up

# 3) Start worker (from the chorus-worker repo, port 9670)
#    ... see worker repo instructions ...

# 4) Run controller
make run
# Open http://localhost:8081/swagger/index.html to try the API
```

### Troubleshooting
- Cannot connect to Postgres:
  - Check container: `docker ps | grep chorus-pg`
  - Verify DSN: `echo $POSTGRES_DSN`
  - Check health: `docker compose ps`
- Migration issues:
  - Status: `make migrate-status`
  - Retry: `make migrate-up`
  - Baseline if DB already has schema: `make migrate-baseline`
- Cannot reach worker:
  - Verify worker is running and `WORKER_GRPC_ADDR` is correct
  - In Docker, use `host.docker.internal` to reach host services
- Swagger not showing:
  - Ensure route `/swagger/*any` exists and open `http://localhost:8081/swagger/index.html`

### Makefile Quick Reference
```bash
make build            # build binary
make run              # run server
make test             # run tests
make clean            # remove bin/
make migrate-up       # apply migrations
make migrate-down     # rollback most recent
make migrate-status   # show migration status
make migrate-new      # create new migration file
make migrate-baseline # baseline existing DB
make migrate-hash     # generate checksums
make docker-build     # build docker image
make docker-run       # run container (env mapping)
```
