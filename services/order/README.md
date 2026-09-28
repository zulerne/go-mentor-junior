# Order Service

Order Service manages carts, orders, restaurant decisions, and delivery requests.

## Prerequisites

- Go 1.27.1 or higher
- [Task](https://taskfile.dev/) CLI for running build commands

## Configuration

Configuration is loaded from environment variables with the `` prefix. Defaults are built in.

| Variable | Default | Description |
|---|---|---|
| `ENV` | `local` | Environment: `local` \| `production` |
| `HTTP_ADDRESS` | `:8080` | HTTP listen address |
| `HTTP_TIMEOUT` | `5s` | Read/write timeout |
| `HTTP_IDLETIMEOUT` | `60s` | Idle connection timeout |
| `HTTP_SHUTDOWNTIMEOUT` | `10s` | Graceful shutdown timeout |

Copy `.env.example` to `.env` for local overrides — Taskfile loads it automatically:

```bash
cp .env.example .env
```

## Build and Run

### Build

```bash
task build
```

### Run

```bash
task run
```

## Format and Vet

### Format

```bash
task fmt
```

### Vet

```bash
task vet
```

### Lint

```bash
task lint
```

### Lint Fix

```bash
task lint-fix
```

### Test

```bash
task test
```

### Test with Race Detector

```bash
task test-race
```

### Coverage

```bash
task coverage
```

### Regenerate mocks

```bash
go generate ./...
```

## Project Structure

```
cmd/order/           - Application entry point
internal/
├── api/             - HTTP handlers and middleware
├── config/          - Configuration
├── domain/          - Domain models
├── logger/          - Logging setup
├── provider/        - Data providers (in-memory implementations)
└── services/        - Business logic (Customer and Restaurant services)
```

## API Contract

See the [project specification](../../docs/project.md) for full API contract details.

### Customer Endpoints

- `GET /cart`
- `PUT /cart/items/{menu_item_id}`
- `DELETE /cart/items/{menu_item_id}`
- `POST /orders`
- `GET /orders/{order_id}`
- `GET /orders`
- `POST /orders/{order_id}/cancel`

### Restaurant Endpoints

- `GET /restaurant/orders`
- `POST /restaurant/orders/{order_id}/accept`
- `POST /restaurant/orders/{order_id}/reject`
- `POST /restaurant/orders/{order_id}/start-preparation`
- `POST /restaurant/orders/{order_id}/ready`

## Important Implementation Decisions

See comments in the codebase (marked with `TODO (review)`) for open design questions.
