# Order Service

Order Service manages carts, orders, restaurant decisions, and delivery requests.

## Prerequisites

- Go 1.27.1 or higher
- [Task](https://taskfile.dev/) CLI for running build commands

## Configuration

Configure the service using environment variables. See `.env.example` for available configuration options.

```bash
cp .env.example .env
# Edit .env with your settings
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

Or run the service binary directly:

```bash
task run-order
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
