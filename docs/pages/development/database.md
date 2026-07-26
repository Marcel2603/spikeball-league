# Database Behavior

The Spikeball League Tracker uses **SQLite** as its embedded database to store leagues, teams, players, and match history.
This design ensures that the application remains entirely self-contained and easy to deploy without needing external
database services like PostgreSQL.

## Schema & Queries (`sqlc`)

We use [`sqlc`](https://sqlc.dev/) to generate type-safe Go code from SQL queries.

- **Schema:** Defined in `internal/db/schema.sql`.
- **Queries:** Defined in `internal/db/queries.sql`.
- **Generated Code:** When you run `sqlc generate` (or `make generate` if configured), it generates the interface
- `Querier` and the implementation in `internal/db/`.

This approach eliminates the need for reflection-based ORMs and catches SQL syntax errors at compile time.

## Migrations

Database migrations are managed using [`golang-migrate`](https://github.com/golang-migrate/migrate).

1. **Embedded Migrations:** The SQL migration files are located in `internal/db/migrations/` and are embedded directly into
the binary using Go's `//go:embed` directive.
2. **Auto-Migration on Startup:** When the application starts, `db.NewDB(config)` automatically applies any pending migrations.
You do not need to run a separate migration script manually.

## Configuration

The database file location is fully configurable via the `app.yaml` configuration file (or via environment variables).

```yaml
# app.yaml
database:
  path: "spikeball.db"
```

If the database file does not exist, SQLite will automatically create it upon first connection.

## Testing

When writing tests for handlers or database layers, you can use the special `:memory:` path to create a fast, isolated,
in-memory SQLite database instance. This prevents tests from clashing on a physical file and guarantees an empty state for
every test run.

Example usage for testing:

```go
config := config.DatabaseConfig{Path: ":memory:"}
testDB, _ := db.NewDB(config)
```
