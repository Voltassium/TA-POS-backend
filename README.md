# go-template

A maintainable Go backend template with layered architecture (controllers, services, repositories), Bun-based database access, MinIO-compatible object storage support, JWT authentication, and structured configuration.

## Quick Start

- **Prerequisites:** Go (1.18+), PostgreSQL (or the DB you configure), MinIO (optional), and `git`.
- **Get dependencies:**

```bash
go mod download
```

- **Configuration:** copy and update the environment file:

- File: [config/files/env.example.yaml](config/files/env.example.yaml)

Copy it to `config/files/env.yaml` and set your database, object storage, JWT and other credentials.

- **Run locally:**

```bash
go run main.go
# or build then run
go build -o bin/app .
./bin/app
```

## What this project contains

- **Layered architecture:** separation between `controllers`, `services`, `repository`, and `domain` for clear responsibilities.
- **Database:** [database/bun.go](database/bun.go) shows Bun usage for DB access.
- **Object storage:** MinIO-compatible implementation in [storage/minio.go](storage/minio.go).
- **Authentication:** JWT utilities in `pkg/authentication` and related domain/service layers.
- **Config management:** configuration types and initialization in `config/` with default example in [config/files/env.example.yaml](config/files/env.example.yaml).
- **HTTP server:** server initialization and middlewares under `http/server/` including error handling and request logging.
- **Logging:** Zap-based logger in `http/logger`.

## Project Layout (high-level)

- `main.go` — application entrypoint.
- `config/` — application config and environment files.
- `internal/` — application code (controllers, services, etc.).
- `domain/` — domain models and business logic types.
- `dto/` — request/response DTOs and pagination helpers.
- `repository/` — data access layer implementations.
- `pkg/` — reusable packages (e.g., `authentication`, encryption, JWT helpers).
- `database/` — Bun DB setup and helpers.
- `http/` — server helpers, middlewares, client wrappers.
- `storage/` — object storage abstraction and MinIO implementation.
- `utils/`, `validation/`, `errors/` — helper utilities and domain-specific errors.

## Configuration details

- Environment example: [config/files/env.example.yaml](config/files/env.example.yaml).
- Important settings:
	- Database DSN (used by `database/bun.go`).
	- Object storage endpoint and credentials (used by `storage/minio.go`).
	- JWT secret and token TTL (used by `pkg/authentication/jwt.go`).

Keep secrets out of VCS. Use your preferred secret manager in production.

## Running & Development

- Start the server locally after configuring `config/files/env.yaml`.
- Run all tests (if present):

```bash
go test ./... 
```

- Useful commands:

```bash
go vet ./...
go fmt ./...
golangci-lint run # if you use golangci-lint
```

## Adding a new API endpoint

1. Add request/response DTOs under `dto/requests` and `dto/response`.
2. Implement business logic in `internal/services`.
3. Add repository functions in `repository/` if DB access is required.
4. Add a controller in `internal/controllers` and register route in `routes/`.

## Notes on components

- Authentication: JWT helpers are in `pkg/authentication` and domain/service wiring lives under `internal` and `domain`.
- Storage: MinIO implementation included; adapt `storage/implementation.go` and `storage/minio.go` for other providers.
- Database: Bun is used — see `database/bun.go` and `repository/` implementations.

## Contributing

Contributions are welcome. Please open issues or PRs with clear descriptions and tests where applicable.

## License

Add your project license here.

---
If you'd like, I can also run `go test ./...`, create a small development Makefile, or commit this change for you.