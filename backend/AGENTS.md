# Backend agent harness

## Authority and scope

- This file governs every agent working anywhere under `backend/`.
- Work only inside the `backend/` directory. Read, create, edit, move, and delete files only beneath this directory.
- The only exceptions are read-only access to `../challenge.md` and `../endpoints.md`. Read both at the start of every task. `challenge.md` is the source of truth for project requirements, and `endpoints.md` is the source of truth for the API contract and how endpoints are handled.
- Do not inspect or modify the repository root, `.git`, `frontend/`, or any other sibling directory or file. In particular, do not read root documentation or infer frontend implementation details.
- Never modify `../challenge.md` or `../endpoints.md`. If a request conflicts with either source of truth or requires a change outside `backend/`, stop and report the conflict or required handoff instead of crossing the boundary.
- Run all commands with `backend/` as the working directory. Do not use recursive searches rooted at `..`, absolute paths outside this directory, or commands that traverse parent directories. The only permitted parent-path commands are those that read `../challenge.md` or `../endpoints.md`.
- Treat network services and external systems as out of scope unless the user explicitly authorizes them.

## Mission

Build and maintain the Go REST microservice required by `../challenge.md`, following the API contract in `../endpoints.md`. It must expose calculator operations, validate requests, handle arithmetic edge cases, return consistent JSON, and remain simple to test and understand.

Required operations:

- Addition
- Subtraction
- Multiplication
- Division, including division-by-zero handling

Optional operations may be implemented only after the required behavior is complete and tested:

- Exponentiation
- Square root
- Percentage

## Existing stack and conventions

- Language: Go, using the version declared in `go.mod`.
- Router: `github.com/go-chi/chi/v5`.
- Keep transport concerns in HTTP handlers/router code and arithmetic rules in the service layer.
- Prefer the Go standard library unless a dependency clearly reduces complexity.
- Keep exported APIs minimal, use idiomatic error handling, and format changed Go files with `gofmt`.
- Implement the API contract exactly as defined in `../endpoints.md`. Do not silently change endpoint paths, methods, request/response shapes, or status-code behavior. Keep them explicit in backend-local code and tests.

## API quality bar

- Accept only documented operations and finite numeric operands.
- Reject malformed JSON, missing operands, unsupported operations, invalid numeric values, and unsafe edge cases with an appropriate `4xx` status.
- Return successful results as JSON with a stable shape and correct `Content-Type`.
- Return errors as JSON with a stable, useful, non-sensitive message shape.
- Do not panic or emit `NaN`/`Infinity` as JSON.
- Make browser consumption possible with a narrowly configured CORS policy when the architecture requires cross-origin local development.
- Keep the server port/configuration easy to run locally; prefer environment configuration with a documented default.

## Required test coverage

Use table-driven tests where they improve clarity. Cover at least:

- Every required arithmetic operation, with representative positive, negative, zero, and decimal inputs.
- Division by zero.
- Unsupported operations and malformed or missing input.
- HTTP method, route, status code, response headers, and JSON body behavior.
- Any optional operation and its domain errors if implemented.

Avoid tests that only restate implementation details. Test observable service and HTTP behavior. Do not call a live external service from tests.

## Work loop

1. Read `../challenge.md` and `../endpoints.md`, then inspect only files beneath `backend/` that are relevant to the task.
2. State assumptions in the work product when the requirements leave an API choice open.
3. Make the smallest coherent change that satisfies the requirement.
4. Add or update tests alongside behavior changes.
5. Format and verify locally from `backend/`.
6. Review the diff for accidental scope expansion, generated artifacts, secrets, and unclear API behavior.
7. Report what changed, the verification performed, and any frontend contract or root-documentation handoff needed. Do not perform that handoff yourself.

## Verification commands

Run the strongest applicable checks before declaring success:

```sh
gofmt -w <changed-go-files>
go vet ./...
go test ./...
go test -coverprofile=coverage.out ./...
```

`coverage.out` is a disposable local artifact; do not commit it unless explicitly requested. If a check cannot run, report the exact command and reason. Never claim a check passed unless it was executed successfully.

## Completion criteria

Backend work is complete only when the requested behavior follows `../challenge.md` and the API follows `../endpoints.md`, required paths and edge cases have tests, `go vet ./...` and `go test ./...` pass, formatting is clean, and no file outside `backend/` was changed.
