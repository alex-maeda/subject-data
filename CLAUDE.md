# subject-data

Go API server for Sovra's subject data service. Manages Subjects, Properties, and the Records used to derive those properties. Server-to-server only — never called directly from the web app.

## Tech Stack

- Go 1.26
- chi v5 for HTTP routing
- golangci-lint for linting
- Target deploy: AWS ECS Fargate (region `us-east-2`)

## Commands

- `make build` — compile binary to `bin/run_service`
- `make test` — run all tests
- `make run` — start dev server
- `make lint` — run golangci-lint
- `make fmt` — format all Go files

## Project Structure

- `cmd/run_service/` — entry point, signal handling
- `internal/` — all application code
  - `api/` — HTTP server, router, handlers, request/response types
  - `middleware/` — Bearer token auth middleware

## Code Style

- Use `internal/` for all application code
- Use `log/slog` for structured logging — never `fmt.Println` or `log.Printf`
- Errors: wrap with `%w`, return early on error (no else after error check)
- Naming: `userID`, `httpClient` (Go conventions)
- Reuse shared test wrappers/mocks; don't duplicate per-file

## API Direction

- **REST API** — CRUD operations for Subjects, Properties, and Records
- **Server-to-server auth only** — the web app never calls this service directly
- Bearer token auth via `AUTH_TOKENS` env var

## Pulumi / Infrastructure

- **ALWAYS run Pulumi in Docker** — never run `pulumi` commands directly on host
- **Inside Docker, ALWAYS check IP first** before running any other commands (e.g., `curl -s https://checkip.amazonaws.com` inside the container). Only proceed with Pulumi/AWS commands after confirming the IP is correct.
- After IP check, run `aws sts get-caller-identity` to verify the correct IAM identity
- Use `./scripts/pulumi-docker.sh` for all Pulumi operations (preview, up, destroy, etc.)
- Use `docker-compose.yml` services: `pulumi` (AWS-backed state) or `pulumi-local` (file-backed state)
- `go build ./...` on host is acceptable for compilation checks only
- All IP checks, AWS commands, and Pulumi operations must happen inside Docker, never on host

## Related Repos

- `sovraai/subject-data` — structured subject search API. Project conventions (CI/CD, infra, code style) are mirrored from here.
