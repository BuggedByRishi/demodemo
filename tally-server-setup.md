# Tally Server — Backend Setup Documentation

Complete step-by-step record of setting up the Go backend server and connecting it to a PostgreSQL database.

---

## Table of Contents

1. [Project Goal](#1-project-goal)
2. [Tech Stack](#2-tech-stack)
3. [Project Structure](#3-project-structure)
4. [Setting Up the Go Module](#4-setting-up-the-go-module)
5. [Writing the Core Files](#5-writing-the-core-files)
6. [Database Migrations](#6-database-migrations)
7. [Environment Variables](#7-environment-variables)
8. [Installing Dependencies](#8-installing-dependencies)
9. [Fixing Common Errors](#9-fixing-common-errors)
10. [Running the Server](#10-running-the-server)
11. [Verifying the Connection](#11-verifying-the-connection)

---

## 1. Project Goal

Build a production-style Go HTTP backend server that:
- Connects to a PostgreSQL database
- Follows a clean layered architecture (handler → service → repository → DB)
- Uses existing migration files to set up the database schema
- Reads configuration from environment variables

---

## 2. Tech Stack

| Tool | Purpose |
|---|---|
| Go 1.22+ | Backend language |
| `net/http` | HTTP server and router (standard library, no framework) |
| `pgx/v5` | PostgreSQL driver (`pgxpool` for connection pooling) |
| `goose` | Database migration runner |
| PostgreSQL | Database |

---

## 3. Project Structure

```
Tally-server/
├── cmd/
│   └── server/
│       └── main.go              ← entry point, wires all layers
├── internal/
│   ├── config/
│   │   └── config.go            ← reads DB_* env vars, builds DSN
│   ├── db/
│   │   └── db.go                ← opens pgxpool connection pool
│   ├── handler/
│   │   ├── health_handler.go    ← GET /health
│   │   ├── organization_handler.go
│   │   ├── agent_handler.go
│   │   ├── reference_handler.go
│   │   └── helpers.go           ← shared writeJSON, lastPathSegment
│   ├── model/
│   │   ├── auth.go              ← auth.users, sessions, accounts structs
│   │   ├── core.go              ← organizations, agents, partners structs
│   │   ├── proxmox.go           ← clusters, nodes, instances structs
│   │   └── errors.go            ← ErrNotFound, ErrConflict sentinels
│   ├── repository/
│   │   ├── organization_repo.go ← SQL for core.organizations
│   │   ├── agent_repo.go        ← SQL for core.agents
│   │   └── reference_repo.go    ← SQL for countries, states, themes
│   └── service/
│       ├── organization_service.go
│       └── agent_service.go
├── migrations/
│   ├── 00001_schemas.sql
│   ├── 00002_roles_and_permissions.sql
│   ├── 00003_extensions.sql
│   ├── 00004_types.sql
│   ├── 00005_auth_sessions.sql
│   ├── 00006_core_references.sql
│   ├── 00007_core_partners.sql
│   ├── 00008_core_organizations.sql
│   └── 00009_core_agents.sql
├── .env                         ← never commit this
├── .env.example
└── go.mod
```

### Why this structure?

- `cmd/` — entry point only. No business logic here.
- `internal/` — Go enforces that nothing outside this project can import these packages.
- Each layer only talks to the layer directly below it. Handlers never touch SQL. Repositories never know about HTTP.

---

## 4. Setting Up the Go Module

### go.mod

```
module Tally-server

go 1.22

require (
    github.com/jackc/pgx/v5 v5.10.0
)
```

> **Important:** The module name (`Tally-server`) is just a namespace — it has nothing to do with GitHub. Every import inside the project uses this prefix: `"Tally-server/internal/config"`, `"Tally-server/internal/db"` etc.

---

## 5. Writing the Core Files

### `internal/config/config.go`

Reads individual `DB_*` environment variables and assembles a PostgreSQL DSN string.

```go
package config

import (
    "fmt"
    "log"
    "os"
)

type Config struct {
    DatabaseURL string
    Port        string
    Env         string
}

func Load() Config {
    host     := getenvOr("DB_HOST", "localhost")
    port     := getenvOr("DB_PORT", "5432")
    user     := mustGetenv("DB_USER")
    password := mustGetenv("DB_PASSWORD")
    dbname   := mustGetenv("DB_NAME")
    sslmode  := getenvOr("DB_SSLMODE", "disable")

    dsn := fmt.Sprintf(
        "postgres://%s:%s@%s:%s/%s?sslmode=%s",
        user, password, host, port, dbname, sslmode,
    )

    return Config{
        DatabaseURL: dsn,
        Port:        getenvOr("PORT", "8080"),
        Env:         getenvOr("APP_ENV", "development"),
    }
}

func mustGetenv(key string) string {
    v := os.Getenv(key)
    if v == "" {
        log.Fatalf("required environment variable %q is not set", key)
    }
    return v
}

func getenvOr(key, fallback string) string {
    if v := os.Getenv(key); v != "" {
        return v
    }
    return fallback
}
```

---

### `internal/db/db.go`

Opens a `pgxpool` connection pool. Uses `pgxpool` directly instead of `database/sql` wrapper to avoid driver registration issues.

```go
package db

import (
    "context"
    "fmt"
    "time"

    "github.com/jackc/pgx/v5/pgxpool"
)

func New(dsn string) (*pgxpool.Pool, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    config, err := pgxpool.ParseConfig(dsn)
    if err != nil {
        return nil, fmt.Errorf("parse db config: %w", err)
    }

    config.MaxConns = 25
    config.MinConns = 5
    config.MaxConnLifetime = 5 * time.Minute
    config.MaxConnIdleTime = 1 * time.Minute

    pool, err := pgxpool.NewWithConfig(ctx, config)
    if err != nil {
        return nil, fmt.Errorf("open db pool: %w", err)
    }

    if err := pool.Ping(ctx); err != nil {
        return nil, fmt.Errorf("ping db: %w", err)
    }

    return pool, nil
}
```

> **Why pgxpool instead of database/sql?**
> Using `sql.Open("pgx", dsn)` with `pgx/v5` caused a `unknown driver "pgx"` error because the driver name changed between versions. Using `pgxpool` directly bypasses this entirely and is the recommended approach for pgx/v5.

---

### `cmd/server/main.go`

Wires all layers together and registers routes.

```go
package main

import (
    "log"
    "net/http"

    "Tally-server/internal/config"
    "Tally-server/internal/db"
    "Tally-server/internal/handler"
    "Tally-server/internal/repository"
    "Tally-server/internal/service"
)

func main() {
    // 1. Load config from DB_* env vars
    cfg := config.Load()

    // 2. Connect to PostgreSQL
    database, err := db.New(cfg.DatabaseURL)
    if err != nil {
        log.Fatalf("cannot connect to database: %v", err)
    }
    defer database.Close()
    log.Println("database connected successfully")

    // 3. Wire layers: repository → service → handler
    orgRepo   := repository.NewOrganizationRepository(database)
    agentRepo := repository.NewAgentRepository(database)
    refRepo   := repository.NewReferenceRepository(database)

    orgSvc   := service.NewOrganizationService(orgRepo)
    agentSvc := service.NewAgentService(agentRepo)

    healthH := handler.NewHealthHandler(database)
    orgH    := handler.NewOrganizationHandler(orgSvc)
    agentH  := handler.NewAgentHandler(agentSvc)
    refH    := handler.NewReferenceHandler(refRepo)

    // 4. Register routes
    mux := http.NewServeMux()

    mux.HandleFunc("GET /health", healthH.Health)

    mux.HandleFunc("GET /countries", refH.ListCountries)
    mux.HandleFunc("GET /countries/{id}/states", refH.ListStates)
    mux.HandleFunc("GET /themes", refH.ListThemes)

    mux.HandleFunc("POST /organizations", orgH.Create)
    mux.HandleFunc("GET /organizations", orgH.List)
    mux.HandleFunc("GET /organizations/{id}", orgH.GetByID)

    mux.HandleFunc("POST /agents", agentH.Create)
    mux.HandleFunc("GET /agents/{id}", agentH.GetByID)
    mux.HandleFunc("GET /organizations/{id}/agents", agentH.ListByOrganization)

    // 5. Start server
    addr := ":" + cfg.Port
    log.Printf("server starting on %s (env: %s)", addr, cfg.Env)
    if err := http.ListenAndServe(addr, mux); err != nil {
        log.Fatalf("server error: %v", err)
    }
}
```

---

## 6. Database Migrations

Migrations are written in SQL using [goose](https://github.com/pressly/goose) format. Each file has `-- +goose Up` and `-- +goose Down` sections.

### Migration order

| File | What it creates |
|---|---|
| `00001_schemas.sql` | All schemas: `auth`, `core`, `proxmox`, `guacamole`, `tally`, `extensions` |
| `00002_roles_and_permissions.sql` | DB roles and their privileges |
| `00003_extensions.sql` | `citext` extension for case-insensitive emails |
| `00004_types.sql` | Custom ENUM type `core.apps` |
| `00005_auth_sessions.sql` | `auth.users`, `auth.sessions`, `auth.accounts`, `auth.verifications` |
| `00006_core_references.sql` | `core.countries`, `core.states`, `core.frontend_themes` (with seed data) |
| `00007_core_partners.sql` | `core.partners` |
| `00008_core_organizations.sql` | `core.organizations`, `core.organization_teams` |
| `00009_core_agents.sql` | `core.agents`, `core.agent_teams`, `core.agent_oauth_associations` |

### Running migrations

```bash
# Install goose
go install github.com/pressly/goose/v3/cmd/goose@latest

# Run all migrations
goose -dir ./migrations postgres \
  "postgres://rishi:rishi@localhost:5432/tally_on_cloud?sslmode=disable" \
  up
```

---

## 7. Environment Variables

### `.env` file

```dotenv
DB_HOST=localhost
DB_PORT=5432
DB_USER=rishi
DB_PASSWORD=rishi
DB_NAME=tally_on_cloud
DB_SSLMODE=disable
```

> **Never commit `.env` to git.** Only commit `.env.example` with blank values.

### Loading env vars

`source .env` had issues with line endings. The reliable way:

```bash
export $(grep -v '^#' .env | xargs)
```

---

## 8. Installing Dependencies

```bash
# Install pgx connection pool
go get github.com/jackc/pgx/v5/pgxpool

# Tidy up go.mod and go.sum
go mod tidy
```

---

## 9. Fixing Common Errors

### `import is a program, not an importable package`

**Cause:** A file inside `internal/` had `package main` instead of its correct package name.

**Fix:** Check all package declarations and correct them:
```bash
grep -r "^package" internal/
# then fix any wrong ones, e.g.:
sed -i 's/package main/package db/' internal/db/db.go
```

### `found packages handler (x.go) and main (y.go)`

**Cause:** Two files in the same folder had different package names. One still said `package main`.

**Fix:** Run `sed` on the offending file to fix the package name.

### `required environment variable "DATABASE_URL" is not set`

**Cause:** The old `config.go` expected a single `DATABASE_URL` but the `.env` file had individual `DB_*` variables.

**Fix:** Rewrote `config.go` to read `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`, `DB_SSLMODE` separately and assemble the DSN.

### `required environment variable "DB_USER" is not set` (after sourcing .env)

**Cause:** `source .env` failed silently due to Windows-style line endings (CRLF) in the `.env` file.

**Fix:**
```bash
export $(grep -v '^#' .env | xargs)
```

### `sql: unknown driver "pgx"`

**Cause:** `pgx/v5/stdlib` changed its driver registration. Using `sql.Open("pgx", dsn)` no longer works reliably.

**Fix:** Switched from `database/sql` + stdlib wrapper to `pgxpool` directly:
```go
// Before (broken)
import _ "github.com/jackc/pgx/v5/stdlib"
db, err := sql.Open("pgx", dsn)

// After (correct)
import "github.com/jackc/pgx/v5/pgxpool"
pool, err := pgxpool.NewWithConfig(ctx, config)
```

### `cannot use *pgxpool.Pool as *sql.DB`

**Cause:** After switching to `pgxpool`, all files that accepted `*sql.DB` needed to be updated.

**Fix:** Updated all handler, repository, and health check files to accept `*pgxpool.Pool`. Also updated method calls:

| `database/sql` | `pgxpool` |
|---|---|
| `db.QueryRowContext(ctx, ...)` | `pool.QueryRow(ctx, ...)` |
| `db.QueryContext(ctx, ...)` | `pool.Query(ctx, ...)` |
| `db.ExecContext(ctx, ...)` | `pool.Exec(ctx, ...)` |
| `db.PingContext(ctx)` | `pool.Ping(ctx)` |
| `sql.ErrNoRows` | `pgx.ErrNoRows` |

---

## 10. Running the Server

```bash
# Load environment variables
export $(grep -v '^#' .env | xargs)

# Run the server
go run ./cmd/server
```

Expected output:
```
2026/06/04 11:54:30 database connected successfully
2026/06/04 11:54:30 server starting on :8080 (env: development)
```

---

## 11. Verifying the Connection

### Check the database in psql

```bash
psql -U rishi -d tally_on_cloud

# Verify schemas
\dn

# Verify tables
\dt core.*
\dt auth.*

# Verify seeded data
SELECT * FROM core.countries;
SELECT theme_name FROM core.frontend_themes;
\q
```

### Test API endpoints

```bash
# Health check — confirms server and DB are alive
curl http://localhost:8080/health
# → {"status":"ok"}

# Countries — confirms migrations ran and seeded data exists
curl http://localhost:8080/countries
# → [{"country_id":1,"country_name":"India"}, ...]

# States for India
curl http://localhost:8080/countries/1/states

# Themes
curl http://localhost:8080/themes
# → list of 6 colour themes
```

If all three return data, the full stack is working:
```
curl → HTTP server → handler → repository → PostgreSQL → response
```

---

*Documentation generated: June 2026*
