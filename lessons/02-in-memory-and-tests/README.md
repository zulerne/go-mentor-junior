# Lesson 2: in-memory business logic and tests

## Goal

Replace the Lesson 1 stubs with a complete in-memory implementation of the assigned service and verify its behavior with automated tests.

The required behavior is defined by the [project specification](../../docs/project.md). This lesson does not change the API contract.

Each student completes the common tasks and one assigned service track.

## Support material

[Testing service boundaries](topics/01-testing-service-boundaries.md) provides optional guidance and references for the testing techniques used in this lesson.

## What you will practice

- turning business rules into observable test cases;
- choosing between HTTP handler, application-service, and store tests;
- writing table-driven tests and subtests;
- using mocks and stubs instead of running external services;
- testing HTTP handlers without opening a network port;
- detecting unsafe access to mutable in-memory state with the race detector.

## Common tasks

1. Implement the complete business logic of the assigned service using concurrency-safe in-memory storage.
2. Keep dependencies outside the assigned service as configurable in-process stubs.
3. Add unit tests for HTTP handlers, the application-service layer, and in-memory stores.
4. Use mocks for provider dependencies in application-service tests.
5. Reach at least 60% total statement coverage for the service.
6. Make `go test ./...` and `go test -race ./...` pass.
7. Preserve the Lesson 1 requirements: build, formatting, static analysis, health probes, and graceful shutdown.
8. Document commands for running tests and checking coverage in the service README. If mocks are generated, document how to regenerate them.

## Testing expectations

All tests in this lesson are unit tests:

- use `httptest` to test HTTP validation and request/response handling;
- use mocks to test application-service business logic independently from providers;
- test the real in-memory store implementations directly, without mocks.

The lesson defines which layers must be tested, but students choose the individual test scenarios needed to demonstrate that their implementation is correct. Tests must be isolated from one another and must not require sleeps, real network services, or a particular execution order.

The Go standard library is sufficient for writing the tests. [Testify](https://github.com/stretchr/testify) is recommended for assertions and [Minimock](https://github.com/gojuno/minimock) is the recommended mock generator. [GoMock](https://github.com/uber-go/mock) and [Mockery](https://github.com/vektra/mockery) are acceptable alternatives. When generated mocks are used, pin the generator version and provide a reproducible generation command, preferably through `go generate`.

Measure coverage with:

```bash
go test -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
```

The reported total statement coverage must be at least 60%.

## Track A: Order Service

1. Implement all Order Service operations and business rules from the project specification.
2. Implement in-memory Cart and Order stores.
3. Provide configurable in-process Identity, Restaurant, and Delivery stubs.

## Track B: Restaurant Service

1. Implement all Restaurant Service operations and business rules from the project specification.
2. Implement in-memory Restaurant and Menu Item stores.

## Acceptance criteria

1. Every endpoint owned by the assigned service implements the behavior defined by the project specification.
2. All application data is held in memory, and mutable in-memory stores are safe for concurrent use.
3. No database or running external service is required.
4. HTTP handler tests verify transport validation and request/response handling.
5. Application-service tests verify business logic using mocks for provider dependencies.
6. The real in-memory store implementations have their own tests.
7. Total statement coverage is at least 60%.
8. `go test ./...` and `go test -race ./...` succeed.
9. `go fmt ./...`, `go vet ./...`, and `go build ./...` continue to succeed.
10. The service still starts, shuts down gracefully, and exposes `GET /livez` and `GET /readyz`.
11. Test and coverage commands are documented in the service README, together with the mock-generation command when one is used.

## Not required

- PostgreSQL, migrations, or another persistent database;
- integration tests requiring containers;
- real communication between the Order and Restaurant Services;
- OpenAPI, Protobuf, or generated API code;
- a specific assertion or mocking library;
- fuzz tests or benchmarks;
- API versioning or compatibility policies;
- tests for private implementation details.

In-memory data may be lost when the service restarts. Persistent storage and database integration tests are introduced in a later lesson.
