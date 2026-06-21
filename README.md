# Restaurant Menu API

A REST API for managing restaurant menus — categories, items, prices, tags, and a read-optimized menu card view. This repository was built as a **learning project for idiomatic Go**: standard library HTTP, explicit layering, minimal dependencies, and patterns you can reuse in future API services.

---

## Table of Contents

- [Overview](#overview)
- [Tech Stack](#tech-stack)
- [Quick Start](#quick-start)
- [Project Structure](#project-structure)
- [Architecture](#architecture)
- [Design Decisions](#design-decisions)
- [Database Schema](#database-schema)
- [API Reference](#api-reference)
- [Sample Data](#sample-data)
- [Testing](#testing)
- [Development Workflow](#development-workflow)
- [Docker](#docker)
- [CI / Releases](#ci--releases)
- [Patterns for Future Projects](#patterns-for-future-projects)

---

## Overview

The API models a restaurant menu with:

- **Categories** — group menu items (e.g. Cupcakes, Brownies)
- **Menu items** — individual dishes with vegetarian/availability flags
- **Menu prices** — multiple portion sizes per item (e.g. HALF / FULL)
- **Tags** — labels like "Spicy" (linked via a join table)
- **Menu card** — a denormalized read endpoint that joins items, prices, and tags for display

Data flows through four layers: **handlers → services → repositories → database**. Wiring happens in `internal/app/router.go`; the `main` package only boots infrastructure and starts the server.

---

## Tech Stack

| Layer | Choice | Why |
|-------|--------|-----|
| Language | Go 1.26.4 | Current toolchain; `go tool` for dev deps |
| HTTP | `net/http` + `http.ServeMux` (Go 1.22+ method routing) | No framework; idiomatic stdlib |
| Database | MySQL 8.0 | Relational model with foreign keys |
| Driver | `github.com/go-sql-driver/mysql` | De facto MySQL driver for Go |
| Migrations | `CREATE TABLE IF NOT EXISTS` at startup | Simple for learning; no migration tool |
| Hot reload | [Air](https://github.com/air-verse/air) via `go tool` | Dev convenience without a runtime dep |
| API testing (manual) | [Bruno](https://www.usebruno.com/) collection | Lightweight alternative to Postman |
| Container | Multi-stage Docker + distroless image | Small, secure production image |

**Intentionally not used:** Gin/Echo/Fiber, ORM (GORM/sqlc), migration frameworks, structured logging libraries, DI containers.

---

## Quick Start

### Prerequisites

- Go 1.26.4+
- Docker & Docker Compose (for MySQL and full stack)
- Python 3 + `requests` (only for the sample data script)

### 1. Environment

```bash
cp .env.sample .env
```

Default values match `docker-compose.yaml`:

```env
DB_HOST=localhost        # use "db" when running inside Docker Compose
DB_PORT=3306
DB_USER=root
DB_PASSWORD=password
DB_NAME=restaurant_menu
PORT=8080
HOST=0.0.0.0
```

### 2. Start MySQL

```bash
docker compose up db -d
```

Tables are created automatically on first app start via `internal/db/migrate.go`.

### 3. Run the API

```bash
make run          # build + run
make live         # hot reload with Air
```

Health check: `GET http://localhost:8080/`

### 4. Seed sample data (optional)

With the API running:

```bash
pip install requests   # if not already installed
python create_sample_data.py
```

### 5. Run tests

Unit tests always run. Integration tests need a test database:

```bash
# Add to .env:
# TEST_DB_DSN=root:password@tcp(localhost:3306)/restaurant_menu_test

make test
```

---

## Project Structure

```
restaurant-menu-api/
├── main.go                    # Entry point: router + ListenAndServe
├── init.go                    # init() hooks: logger, config, DB (runs before main)
├── go.mod / go.sum            # Module definition; Air as `tool` directive
├── Makefile                   # build, run, test, audit, live
├── Dockerfile                 # Multi-stage: golang builder → distroless runtime
├── docker-compose.yaml        # MySQL + app services
├── air.toml                   # Air hot-reload config
├── .env.sample                # Environment variable template
│
├── create_sample_data.py      # Seeds DB via HTTP POST calls
├── sample_data.json           # Source dataset (bakery menu items)
│
├── bruno_api_collection/      # Bruno HTTP requests for manual API testing
│   ├── opencollection.yml
│   ├── environments/Local.yml
│   └── *.yml                  # One file per endpoint
│
├── .github/workflows/ci.yaml  # Tag-triggered release pipeline
│
└── internal/                  # All application code (not importable externally)
    ├── app/
    │   ├── router.go          # Dependency wiring: DB → Service → Handler → routes
    │   └── router_test.go     # Integration tests (app_test package)
    │
    ├── config/
    │   └── config.go          # Env-based config with sync.Once singleton
    │
    ├── db/
    │   ├── db.go              # sql.DB connection pool + sync.Once
    │   └── migrate.go         # CREATE TABLE IF NOT EXISTS for all tables
    │
    ├── logger/
    │   └── logger.go          # Thin wrapper around log.Logger
    │
    ├── middleware/
    │   └── middelware.go      # Panic recovery middleware
    │
    ├── models/                # Domain types + request DTOs (JSON tags)
    │   ├── category.go
    │   ├── menu.go
    │   ├── menu_price.go
    │   ├── menu_price_list.go
    │   ├── menu_tag_list.go
    │   ├── menu_card_item.go
    │   └── tag.go
    │
    ├── handlers/              # HTTP layer: parse request, call service, write response
    │   ├── handler.go         # Handler struct + constructor
    │   ├── health.go          # DB connection pool stats
    │   └── menu.go            # All resource handlers
    │
    ├── services/              # Business logic: validation, error mapping
    │   ├── service.go         # Service struct; wires all repositories
    │   ├── category.go
    │   ├── menu.go
    │   ├── menu_price.go
    │   ├── tag.go
    │   └── menu_card_service.go
    │
    ├── repository/            # Data access: prepared statements, SQL
    │   ├── errors.go          # MySQL duplicate → ErrDuplicate
    │   ├── category_repository.go
    │   ├── menu_repository.go
    │   ├── menu_price_repository.go
    │   ├── menu_card_repository.go  # Complex JOIN + GROUP_CONCAT query
    │   └── tag_repository.go
    │
    ├── utils/
    │   ├── utils.go           # JSON request/response helpers
    │   └── utils_test.go
    │
    └── testutil/              # Test-only helpers (not shipped in binary)
        ├── db.go              # Test DB setup, truncate, tag linking
        └── http.go            # httptest helpers + generic response decoding
```

### Package responsibilities

| Package | Responsibility | Imports from |
|---------|----------------|--------------|
| `main` | Bootstrap only | `app`, `config`, `db`, `logger` |
| `internal/app` | Route table + dependency graph | `handlers`, `services`, `middleware` |
| `internal/handlers` | HTTP in/out | `services`, `models`, `utils` |
| `internal/services` | Business rules | `repository`, `models` |
| `internal/repository` | SQL / persistence | `models`, `database/sql` |
| `internal/models` | Structs only | nothing internal |
| `internal/db` | Connection + schema | `config`, `logger` |
| `internal/config` | Environment config | `os`, `logger` |
| `internal/middleware` | Cross-cutting HTTP concerns | `utils`, `logger` |
| `internal/utils` | Shared HTTP/JSON utilities | stdlib only |
| `internal/testutil` | Integration test infrastructure | `db` |

### File naming conventions

- **Repositories:** `{entity}_repository.go` with `New{Entity}Repository(db)` constructor
- **Services:** one file per domain area (`category.go`, `menu.go`, …) with methods on `*Service`
- **Handlers:** grouped in `menu.go` (could be split per resource as the API grows)
- **Models:** separate `Foo` (response) and `FooRequest` (create payload) where shapes differ
- **Tests:** `router_test.go` in `app_test` package (black-box testing of the HTTP layer)

---

## Architecture

```
  Client (curl / Bruno / create_sample_data.py)
                    │
                    ▼
            ┌───────────────┐
            │  middleware   │  Recover (panic → 500)
            └───────┬───────┘
                    ▼
            ┌───────────────┐
            │   handlers    │  Parse JSON, status codes, response envelope
            └───────┬───────┘
                    ▼
            ┌───────────────┐
            │   services    │  Validation, domain errors, logging
            └───────┬───────┘
                    ▼
            ┌───────────────┐
            │  repository   │  Prepared statements, row scanning
            └───────┬───────┘
                    ▼
            ┌───────────────┐
            │  MySQL 8.0    │
            └───────────────┘
```

### Request lifecycle (example: `POST /menu`)

1. `router.go` matches `POST /menu` and wraps handler with `middleware.Recover`
2. `Handler.PostMenu` decodes body into `models.MenuItemRequest` via `utils.ParseRequestBody`
3. `Service.CreateMenuItem` checks category exists, then inserts via `MenuRepository`
4. Duplicate key → `repository.mapError` → `ErrDuplicate` → user-facing `"menu item already exists"`
5. `utils.CreateResponse` writes `{"data": ...}` with appropriate status code

### Response envelope

All JSON responses use a consistent shape:

```json
{
  "data": <payload>
}
```

Errors on write operations return the error message string in `data` with `4xx` status codes.

---

## Design Decisions

### 1. Standard library HTTP only

`http.ServeMux` with Go 1.22+ pattern `"METHOD /path"` keeps routing declarative without a third-party router. For a learning project and small APIs, this is the idiomatic baseline before reaching for Chi or similar.

### 2. `internal/` package boundary

Everything under `internal/` cannot be imported by external modules. This enforces that `main` is the only public entry point and documents intent: this is an application, not a library.

### 3. Layered architecture with explicit wiring

Dependencies are constructed manually in `app.NewRouter` and `services.New` — no reflection-based DI. This makes the dependency graph easy to read and is the recommended starting point before introducing interfaces everywhere.

```
database *sql.DB
    → repositories (constructed in services.New)
        → *services.Service
            → *handlers.Handler
                → http.ServeMux
```

### 4. `init.go` for infrastructure bootstrap

Logger, config, and DB initialize in `init()` functions so `main()` stays focused on serving HTTP. Trade-off: init order is implicit and panics fail fast at startup (acceptable for this app size).

### 5. Singletons with `sync.Once`

`config`, `db`, and `logger` each use `sync.Once` for lazy, thread-safe single initialization. This avoids passing config through every function while keeping startup idempotent.

### 6. Prepared statements in repositories

Each repository prepares SQL at construction time (`db.Prepare`). Benefits:

- SQL lives next to the entity it serves
- Statements are reused across requests
- Constructor returns `( *Repo, error)` so wiring fails early if SQL is invalid

Cleanup note: for long-lived repos this is fine; if you add `Close()` later, close statements in reverse construction order.

### 7. Repository-level error mapping

MySQL duplicate key errors are detected by message substring and mapped to `repository.ErrDuplicate`. Services use `errors.Is` to return stable, client-safe messages without leaking SQL details.

### 8. Startup migrations (`CREATE TABLE IF NOT EXISTS`)

No goose/golang-migrate — schema is applied when the app starts. Good for learning and zero extra tooling; replace with versioned migrations before production schema evolution.

### 9. Natural keys as primary keys

`categories.name`, `menu_items.name`, composite keys on price/tag join tables. Simplifies the sample data script (refer by name, not UUID). Trade-off: renaming entities is harder; UUIDs are better for production APIs with mutable names.

### 10. Menu card as a dedicated read model

`GET /menu-card` uses a repository with a single JOIN query and aggregates rows in Go (prices nested per item). This separates the **write model** (normalized tables) from a **read model** optimized for clients — a stepping stone toward CQRS thinking.

### 11. Integration tests in `app_test` package

Tests live beside `router.go` but import the package as `app_test` to exercise only the public HTTP surface. `testutil.SetupTestDB` skips tests when `TEST_DB_DSN` is unset, so `go test ./...` works in CI without MySQL while local dev can opt in.

### 12. Minimal dependencies

Runtime: MySQL driver only. Dev tools (Air, staticcheck, govulncheck) are invoked via `go run` in the Makefile or declared as `tool` in `go.mod` — not linked into the production binary.

### 13. Distroless Docker image

Multi-stage build compiles a static binary (`CGO_ENABLED=0`), copies it into `gcr.io/distroless/static-debian12:nonroot`. No shell, minimal attack surface.

---

## Database Schema

```
categories
├── name (PK)
└── description

tags
├── name (PK)
└── description

menu_items
├── name (PK)
├── description
├── is_vegetarian
├── available (default TRUE)
└── category → categories(name)

menu_price_lists
├── menu_item_name → menu_items(name)
├── portion_size
├── price
├── currency
└── PRIMARY KEY (menu_item_name, portion_size)

menu_tags_list
├── menu_item_name → menu_items(name)
├── tag
└── PRIMARY KEY (menu_item_name, tag)
```

**Insert order for related data:** categories → tags (optional) → menu items → menu prices → menu_tags_list (no HTTP endpoint yet; tests insert directly).

---

## API Reference

Base URL: `http://localhost:8080`

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/` | Health check (DB pool stats) |
| `GET` | `/categories` | List categories |
| `POST` | `/categories` | Create category |
| `GET` | `/menu` | List menu items |
| `POST` | `/menu` | Create menu item |
| `GET` | `/tags` | List tags |
| `POST` | `/tags` | Create tag |
| `GET` | `/menu-price` | List all prices |
| `POST` | `/menu-price` | Add price for a menu item |
| `GET` | `/menu-card` | Available items with nested prices and tags |

### Example requests

**Create category**

```bash
curl -X POST http://localhost:8080/categories \
  -H "Content-Type: application/json" \
  -d '{"name": "Cupcakes", "description": "Cupcakes"}'
```

**Create menu item** (category must exist)

```bash
curl -X POST http://localhost:8080/menu \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Chocolate Cupcake",
    "description": "Moist chocolate cupcake",
    "is_vegetarian": true,
    "available": true,
    "category": "Cupcakes"
  }'
```

**Add price**

```bash
curl -X POST http://localhost:8080/menu-price \
  -H "Content-Type: application/json" \
  -d '{
    "menu_item_name": "Chocolate Cupcake",
    "portion_size": "1 Piece",
    "price": 60,
    "currency": "INR"
  }'
```

**Menu card** (read-optimized)

```bash
curl http://localhost:8080/menu-card
```

Bruno collections under `bruno_api_collection/` mirror these endpoints for interactive testing.

---

## Sample Data

### `sample_data.json`

Flat array of menu entries under a top-level `"data"` key. Each entry includes category, item fields, and a `"prices"` array:

```json
{
  "data": [
    {
      "menu_item_name": "Chocolate Cupcake",
      "description": "Moist chocolate cupcake with creamy chocolate buttercream frosting",
      "available": true,
      "is_vegetarian": true,
      "category": "Cupcakes",
      "prices": [
        { "price": 60, "currency": "INR", "portion_size": "1 Piece" }
      ]
    }
  ]
}
```

The file contains a full bakery-style menu (cupcakes, brownies, cookies, etc.) used to populate a demo database.

### `create_sample_data.py`

A small Python script that **seeds the database through the public API** (not direct SQL). This validates endpoints end-to-end and mirrors how a real client would onboard data.

**What it does:**

1. Loads `sample_data.json`
2. Deduplicates categories, menu items, and prices into dictionaries (order-preserving uniqueness)
3. POSTs in dependency order:
   - `/categories`
   - `/menu`
   - `/menu-price`
4. Prints elapsed time in milliseconds

**Requirements:** API running on `http://localhost:8080`, Python 3, `requests`.

```bash
pip install requests
python create_sample_data.py
```

**Why Python and not Go?** Keeps the seed script independent of the module graph — useful as an external integration smoke test. A `cmd/seed` Go command would be the next idiomatic step for a production repo.

**Note:** Re-running the script will fail on duplicates (by design — the API rejects duplicate categories/items/prices). Truncate tables or use a fresh database to re-seed.

---

## Testing

```bash
make test          # all tests (-race -buildvcs)
make test/cover    # HTML coverage report
make audit         # test + tidy + fmt + vet + staticcheck + govulncheck
```

### Integration tests (`internal/app/router_test.go`)

| Test | Covers |
|------|--------|
| `TestHealthCheck` | `GET /` returns pool stats |
| `TestCreateCategoryAndGetCategories` | Happy path CRUD |
| `TestCreateCategoryDuplicate` | Duplicate → 400 + message |
| `TestCreateMenuItemRequiresCategory` | FK validation in service layer |
| `TestCreateMenuItemAndGetMenu` | Menu item lifecycle |
| `TestGetMenuCardWithPricesAndTags` | JOIN read model + nested prices |
| `TestGetMenuCardExcludesUnavailableItems` | `available = false` filtered out |

Enable integration tests by setting in `.env`:

```env
TEST_DB_DSN=root:password@tcp(localhost:3306)/restaurant_menu_test
```

`testutil` creates the test database if missing, runs migrations, and truncates tables after each test.

### Unit tests

`internal/utils/utils_test.go` covers JSON helper behavior.

---

## Development Workflow

```bash
make help          # list all targets
make tidy          # go mod tidy + go fix + go fmt
make build         # output: tmp/bin/main
make run           # build and run with .env loaded
make live          # Air hot reload
make audit         # full quality gate before release
```

### Recommended local loop

1. `docker compose up db -d`
2. `make live`
3. Hit endpoints via Bruno or curl
4. `make audit` before pushing

---

## Docker

**Database only** (local Go development):

```bash
docker compose up db -d
make run
```

**Full stack** (app + MySQL):

```bash
docker compose up --build
```

The `app` service overrides `DB_HOST=db` and waits for MySQL healthcheck before starting.

---

## CI / Releases

`.github/workflows/ci.yaml` triggers on version tags (`v*`):

1. `make test`
2. `make build`
3. Publishes GitHub release with `tmp/bin/main` attached

Tag a release locally:

```bash
git tag v0.1.0
git push origin v0.1.0
```

---

## Patterns for Future Projects

Use this repo as a template checklist when starting a new Go API:

### Project layout

```
cmd/api/main.go          # or root main.go for single-binary apps
internal/
  app/       → router + wiring
  config/    → env config (sync.Once)
  db/        → pool + migrations
  handlers/  → HTTP adapters
  services/  → business logic
  repository/→ SQL / external stores
  models/    → types + request DTOs
  middleware/
  testutil/  → integration test helpers
```

### Bootstrap checklist

- [ ] `.env.sample` with every required variable documented
- [ ] `init.go` or explicit `main` bootstrap for logger, config, DB
- [ ] Fail fast on missing config / DB ping failure
- [ ] Connection pool limits (`SetMaxOpenConns`, `SetMaxIdleConns`, `SetConnMaxLifetime`)

### HTTP checklist

- [ ] Consistent JSON response envelope
- [ ] Panic recovery middleware
- [ ] Method + path routing in one file (`router.go`)
- [ ] Handlers stay thin — no SQL in handlers

### Data layer checklist

- [ ] One repository per aggregate/table group
- [ ] Prepared statements in repository constructors
- [ ] Map driver errors to domain errors in repository
- [ ] Services translate domain errors to API messages

### Testing checklist

- [ ] `app_test` black-box HTTP tests
- [ ] `TEST_DB_DSN` opt-in integration tests (skip in CI without DB)
- [ ] `testutil` for httptest + DB lifecycle
- [ ] `make audit` script: test, fmt, vet, staticcheck, vulncheck

### Operations checklist

- [ ] Multi-stage Dockerfile, non-root runtime
- [ ] `docker-compose` for local dependencies
- [ ] Makefile as the single interface for common commands
- [ ] Seed script or `cmd/seed` for demo data

### When to evolve beyond this repo

| Growth signal | Consider |
|---------------|----------|
| Many routes / middleware chains | [chi](https://github.com/go-chi/chi) or std middleware composition |
| Schema changes in production | goose, golang-migrate, or Atlas |
| Complex queries | sqlc (codegen) or careful repository growth |
| Structured logs / tracing | `log/slog`, OpenTelemetry |
| Configuration files | Viper or env-only with validation (e.g. envconfig) |
| Authentication | middleware + context values; later JWT/session packages |

---

## License

Learning project — add a license if you open-source or reuse substantially.
