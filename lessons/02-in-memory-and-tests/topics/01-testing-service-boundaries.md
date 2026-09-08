# Testing service boundaries

## Why boundaries matter

A focused unit test exercises one part of the service and replaces the dependencies beyond that boundary. This makes failures easier to understand and lets business behavior be tested without starting databases or other services.

Lesson 2 uses three test boundaries:

| Test target         | What is real                                                 | What is replaced               |
| ------------------- | ------------------------------------------------------------ | ------------------------------ |
| HTTP handler        | request decoding, transport validation, and response mapping | application-service dependency |
| Application service | business logic                                               | stores and external providers  |
| In-memory store     | store implementation                                         | nothing                        |

These boundaries do not prescribe package names or individual test cases. Choose the scenarios needed to demonstrate the required behavior from the project specification.

## Test isolation

Each test should create its own store, mocks, handler, and mutable fixtures. Tests should not depend on execution order, arbitrary sleeps, shared mutable state, or services left running by another test.

## Essential reading

- [Using Subtests and Sub-benchmarks](https://go.dev/blog/subtests) introduces table-driven tests, named subtests, and selective test execution. A table is useful when several cases exercise the same behavior; a single case can remain a regular test.
- [`net/http/httptest`](https://pkg.go.dev/net/http/httptest) provides requests and response recorders for invoking an HTTP handler without opening a network port.
- [Go maps in action: concurrency](https://go.dev/blog/maps#hdr-Concurrency) explains why a map needs synchronization when it is accessed concurrently. The [`sync` package](https://pkg.go.dev/sync) documents the available synchronization primitives.
- [The Go data race detector](https://go.dev/doc/articles/race_detector) explains `go test -race`. It finds races only in code paths executed by the tests, so a passing result is useful evidence rather than proof that every path is safe.
- [Minimock](https://github.com/gojuno/minimock) documents the recommended mock generator and its `go generate` workflow.

## Optional reading

- [Go Wiki: TableDrivenTests](https://go.dev/wiki/TableDrivenTests) provides a shorter table-driven test example.
- [Testify](https://github.com/stretchr/testify) provides optional assertion helpers.
- [The cover story](https://go.dev/blog/cover) explains Go statement coverage and its reports in more depth.
- [`go generate`](https://pkg.go.dev/cmd/go#hdr-Generate_Go_files_by_processing_source) documents generation directives. Generation is not run automatically by `go build` or `go test`.

Coverage measures which statements executed; it does not show whether assertions were meaningful or whether important scenarios were selected. Treat the lesson's coverage threshold as a minimum feedback signal, not the definition of a sufficient test suite.
