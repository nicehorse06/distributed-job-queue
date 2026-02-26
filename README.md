# distributed-job-queue

Production-oriented Golang job queue API scaffold using Gin.
Includes a Rust gRPC compute-engine example wired through the Go API.

## Run locally

### With Go

```bash
make run
```

### Manual local development (without Docker)

Prerequisites:
- Go 1.22+
- Rust (stable toolchain) + Cargo

1. Start Rust gRPC compute-engine (Terminal A):

```bash
cd compute-engine
cargo run
```

2. Start Go API and point to local Rust service (Terminal B):

```bash
PORT=8080 COMPUTE_ADDR=localhost:50051 go run ./cmd/api
```

3. Test endpoints (Terminal C):

```bash
curl -sS http://localhost:8080/healthz
curl -sS http://localhost:8080/
curl -sS -X POST http://localhost:8080/compute/square \
  -H "Content-Type: application/json" \
  -d '{"value":12}'
```

If port `8080` is occupied, use another port like `PORT=8081`.

### With Docker Compose

```bash
docker compose up --build
```

This starts:
- `api` (Go + Gin) on `:8081` (container internal port is `8080`)
- `compute-engine` (Rust + gRPC) on `:50051`
- `postgres` on `:5432`

## Development

### Local development workflow

```bash
make tidy
make test
make run
```

### Common commands

- `make run`: start the API (default port `8080`)
- `make test`: run all tests
- `make build`: build binary to `bin/distributed-job-queue`

## Endpoints

- GET /healthz -> {"status":"ok"}
- GET / -> {"message":"distributed-job-queue api is running"}
- POST /compute/square -> forwards to Rust gRPC compute service

Example:

```bash
curl -sS -X POST http://localhost:8081/compute/square \
  -H "Content-Type: application/json" \
  -d '{"value":12}'
```

Response:

```json
{"engine":"rust-grpc","square":144,"value":12}
```

## Configuration

- PORT (default: 8080)
- COMPUTE_ADDR (default: localhost:50051)
