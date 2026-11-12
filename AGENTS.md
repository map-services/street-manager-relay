# AGENTS.md

Guidelines for AI agents working on the Street Manager Relay project.

## Project Overview

Go-based relay service for the GOV.UK Street Manager API. It receives SNS notifications, stores events in SQLite, and provides a REST search API with spatial query support via R-tree indexing.

- **Go version**: 1.27
- **Web framework**: Gin
- **Database**: SQLite3 with R-tree spatial index
- **CLI**: Cobra
- **Module**: `github.com/map-services/street-manager-relay`

## Code Style

- Follow the existing codebase style: indentation, naming, import ordering, and comments.
- Use `github.com/cockroachdb/errors` for error wrapping. Alias the standard `errors` package as `stderrors` when both are needed in the same file to avoid conflicts.
- One responsibility per function. Keep handlers thin and delegate to repository methods.
- Embed SQL files using `//go:embed` directives in `internal/db.go`.

## Project Structure

- `main.go` — Entry point; registers Cobra commands.
- `cmd/` — CLI commands (`api-server`, `bulk-loader`, `regen`).
- `internal/` — Core logic: database layer, routes, SNS handling.
- `internal/routes/` — HTTP handlers for search, reference data, and SNS.
- `models/` — Data models (Event, BoundingBox, Facets, etc.).
- `internal/sql/` — Embedded SQL scripts for schema, search, and reference data.

## Testing

- Run tests with `go test ./...`.
- Build the project with `go build ./...`.
- Skip tests for non-functional updates (documentation, comments only).

## Conventions

- **Error handling**: Always wrap errors with context using `errors.Wrap()`. Check `sql.ErrTxDone` before logging rollback errors after a successful commit.
- **Logging**: Use `log/slog` structured logging throughout.
- **Database transactions**: Use `defer tx.Rollback()` with `!stderrors.Is(err, sql.ErrTxDone)` guards to avoid spurious error logs.
- **Pinned dependencies**: Do not upgrade dependencies unless explicitly requested.

## Where to Find Things

- API endpoints: `internal/routes/search.go`, `internal/routes/refdata.go`, `internal/routes/sns.go`
- Database schema: `internal/sql/create_db.sql`
- Search query logic: `internal/sql/search.sql`
- Reference data query: `internal/sql/ref_data.sql`
- Server setup and middleware: `cmd/api_server.go`
