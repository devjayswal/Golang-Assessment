# Reasoning: User API (DOB + Age)

Hey! This doc walks through how I built this thing, why I made certain choices, and how it all ties back to the assessment requirements. I'll keep it real—just practical decisions and where to look in the code.

---

## So, What Did I Build?

A simple REST API where you can create users with a `name` and `dob` (date of birth), and when you fetch them back, it calculates their `age` on the fly. Nothing fancy—just clean code with GoFiber, PostgreSQL, and a layered architecture (handler → service → repository).

---

## How It's Organized

Here's how the pieces fit together:

| Layer | What it does | Where to look |
|-------|--------------|---------------|
| Routes | Wires up endpoints | [internal/routes/routes.go](internal/routes/routes.go) |
| Handlers | Parses requests, returns responses | [internal/handler/user_handler.go](internal/handler/user_handler.go) |
| Service | Business logic, validation, age math | [internal/service](internal/service) |
| Repository | Talks to Postgres | [db/sqlc](db/sqlc) |
| Middleware | Request IDs, logging | [internal/middleware](internal/middleware) |
| Logger | Zap setup | [internal/logger/logger.go](internal/logger/logger.go) |
| Entry point | Boots everything | [cmd/server/main.go](cmd/server/main.go) |

**Request flow in plain terms:**
1. Request comes in → middleware slaps on a request ID and starts a timer.
2. Handler checks if the JSON is valid and the date makes sense.
3. Service does the real work—validates the model, hits the repo, calculates age if needed.
4. Repo runs the SQL against Postgres.
5. Response goes back out, middleware logs how long it took.

---

## The Data Model

Dead simple. One table:

```sql
users(id SERIAL PRIMARY KEY, name TEXT NOT NULL, dob DATE NOT NULL)
```

Migration script is at [db/migrations/001_init.sql](db/migrations/001_init.sql). The Go struct lives in [internal/models/user.go](internal/models/user.go).

---

## Why I Picked These Tools

- **Fiber** – Fast, feels like Express, gets out of my way.
- **PostgreSQL** – Reliable, everyone knows it, great for this kind of app.
- **sqlc** – Keeps SQL in `.sql` files where it belongs and generates type-safe Go code.
- **Zap** – Structured logging that's actually fast. Production-ready out of the box.
- **go-playground/validator** – Tag-based validation. Clean and simple.

---

## The Endpoints

| Method | Path | What it does | Status |
|--------|------|--------------|--------|
| POST | `/users` | Create a user | 201 |
| GET | `/users/:id` | Get user + calculated age | 200 |
| PUT | `/users/:id` | Update name/dob | 200 |
| DELETE | `/users/:id` | Remove user | 204 |
| GET | `/users` | List all (with pagination) | 200 |

All handlers are in [internal/handler/user_handler.go](internal/handler/user_handler.go).

---

## How Age Calculation Works

This was the fun part. The logic is in [internal/service/age.go](internal/service/age.go):

1. Take the difference in years between `dob` and today.
2. If the birthday hasn't happened yet this year, subtract one.
3. If someone passes a future date (weird, but possible), just return 0 instead of negative.

I inject a `clockNow()` function into the service so I can freeze time in tests. Speaking of which—tests are in [internal/service/age_test.go](internal/service/age_test.go) and cover:
- Birthday hasn't happened yet this year
- Birthday is today
- Birthday already passed
- Future dates (edge case)

---

## Validation & Error Handling

I split validation into two layers:

1. **Handler level** – Is it valid JSON? Is the date in `YYYY-MM-DD` format?
2. **Service level** – Is the name present? Does it meet length requirements?

Error responses are straightforward:
- Bad input → `400`
- Can't find the user → `404`
- Created something → `201`
- Deleted something → `204`
- Something blew up → `500`

---

## Pagination

`GET /users` supports `limit` and `offset` query params. Defaults to 50 items, max 500. This keeps things predictable and prevents someone from accidentally pulling the entire database.

---

## Logging & Observability

Every request gets:
- A unique `X-Request-ID` (generated if not provided, echoed back in the response)
- A log entry with method, path, status code, request ID, and duration

The middleware is in [internal/middleware](internal/middleware). Fiber's built-in recover middleware catches panics so the server doesn't crash.

---

## SQLC Setup

I've got sqlc configured at [db/sqlc/sqlc.yaml](db/sqlc/sqlc.yaml) with queries in [internal/repository/queries.sql](internal/repository/queries.sql).

Right now, there's a minimal hand-written implementation in `db/sqlc` so the app runs without needing to generate anything. But if you want to switch to real sqlc:

1. Make sure queries have `-- name:` annotations (they already do).
2. Run `sqlc generate`.
3. Swap in the generated code—the service already uses the `Querier` interface, so it's plug-and-play.

---

## Docker Setup

Everything's containerized:

- **docker-compose.yml** – Spins up Postgres and the app
- **Dockerfile** – Multi-stage build for a small final image

The `db` service has a healthcheck, so the app waits for Postgres to be ready. If the table doesn't exist yet, just run the migration SQL and restart the app container.

---

## Tests

Run `go test ./...` and you'll see the age calculation tests pass.

What's tested:
- Age calculation edge cases

What could be added:
- Integration tests hitting real handlers
- Validation error scenarios

---

## Trade-offs & What I'd Do Next

A few things I'd improve with more time:

- **Validation** – Right now it accepts future dates and returns age 0. Could reject them outright with a 400.
- **Migrations** – Using raw SQL files. A tool like `golang-migrate` would be cleaner for versioning.
- **API docs** – An OpenAPI spec would help document the contract.
- **Developer experience** – A Makefile with `make run`, `make test`, `make sqlc` would be nice.
- **Observability** – Metrics, trace IDs, structured error responses.

---

## How This Covers the Requirements

| Requirement | How I handled it |
|-------------|------------------|
| Store `dob` in DB | Postgres table with DATE column |
| Return dynamic `age` | Calculated in service layer using Go's `time` package |
| Use SQLC | Queries annotated and ready; interface-based design |
| Validate inputs | `go-playground/validator` + date format checks |
| Log with Zap | Middleware logs every request |
| Clean HTTP codes | 201/200/204 for success, 400/404/500 for errors |
| **Bonus:** Docker | Full docker-compose setup |
| **Bonus:** Pagination | `limit`/`offset` on list endpoint |
| **Bonus:** Unit tests | Age calculation thoroughly tested |
| **Bonus:** Middleware | Request ID + request duration logging |

---

## Where to Start Reading

If you're diving into the code:

1. Start at [cmd/server/main.go](cmd/server/main.go) – see how everything boots up
2. Check [internal/routes/routes.go](internal/routes/routes.go) – how endpoints are wired
3. Look at [internal/handler/user_handler.go](internal/handler/user_handler.go) – request/response handling
4. Dig into [internal/service/user_service.go](internal/service/user_service.go) – the actual logic
5. See [db/sqlc](db/sqlc) – database queries

The README has copy-paste commands to get it running locally or with Docker. Should be up in under a minute.

