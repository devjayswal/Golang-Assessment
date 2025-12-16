# User API (Go + Fiber)

A straightforward REST API for managing users with `name` and `dob`. When you fetch a user, it calculates their `age` on the fly—no stored age column, just math.

Built with **GoFiber**, **PostgreSQL**, **sqlc**, **Zap**, and **go-playground/validator**.

---

## What's in the Box

- **Fiber** – Fast HTTP framework
- **PostgreSQL** – Database
- **sqlc** – Type-safe SQL queries
- **Zap** – Structured logging
- **Validator** – Input validation

---

## Getting Started

### Option 1: Run Locally

**1. Set up Postgres**

Create a database and run the migration:

```sql
CREATE DATABASE golang_assessment;
\c golang_assessment;
CREATE TABLE IF NOT EXISTS users (
  id SERIAL PRIMARY KEY,
  name TEXT NOT NULL,
  dob DATE NOT NULL
);
```

Or just run `db/migrations/001_init.sql`.

**2. Set your database URL**

Linux/macOS:
```bash
export DATABASE_URL=postgres://postgres:postgres@localhost:5432/golang_assessment?sslmode=disable
```

Windows PowerShell:
```powershell
$env:DATABASE_URL = "postgres://postgres:postgres@localhost:5432/golang_assessment?sslmode=disable"
```

**3. Run the server**

```bash
go mod tidy
go run ./cmd/server
```

Server starts on `http://localhost:8080`.

---

### Option 2: Use Docker

This is the easy way—everything spins up together.

```bash
docker compose up --build
```

- App runs on `:8080`
- Postgres runs on `:5432` (user: `dev_user`, password: `mysecretpassword`)

**First time?** You'll need to create the table:

```bash
docker compose exec -it db psql -U dev_user -d golang_assessment
```

Then paste the `CREATE TABLE` statement from above (or from `db/migrations/001_init.sql`).

After that:
```bash
docker compose restart app
```

---

## API Endpoints

| Method | Endpoint | What it does |
|--------|----------|--------------|
| POST | `/users` | Create a user |
| GET | `/users/:id` | Get a user (includes calculated age) |
| PUT | `/users/:id` | Update a user |
| DELETE | `/users/:id` | Delete a user |
| GET | `/users` | List all users (supports `limit` & `offset`) |

---

## Quick Examples

**Create a user:**
```bash
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{"name":"Alice","dob":"1990-05-10"}'
```

Response:
```json
{"id":1,"name":"Alice","dob":"1990-05-10"}
```

**Get a user:**
```bash
curl http://localhost:8080/users/1
```

Response:
```json
{"id":1,"name":"Alice","dob":"1990-05-10","age":35}
```

**Update a user:**
```bash
curl -X PUT http://localhost:8080/users/1 \
  -H "Content-Type: application/json" \
  -d '{"name":"Alice Smith","dob":"1990-05-10"}'
```

**Delete a user:**
```bash
curl -X DELETE http://localhost:8080/users/1
```

Returns `204 No Content`.

**List users (with pagination):**
```bash
curl "http://localhost:8080/users?limit=10&offset=0"
```

---

## Running Tests

```bash
go test ./...
```

The age calculation logic has unit tests covering edge cases.

---

## Project Structure

```
cmd/server/main.go      ← Entry point
internal/
  handler/              ← HTTP handlers
  service/              ← Business logic + age calculation
  repository/           ← SQL queries
  models/               ← Data structures
  middleware/           ← Request ID, logging
  routes/               ← Route wiring
  logger/               ← Zap setup
db/
  migrations/           ← SQL migration files
  sqlc/                 ← sqlc config + generated code
```

---

## About sqlc

Queries live in `internal/repository/queries.sql`. The sqlc config is at `db/sqlc/sqlc.yaml`.

Right now there's a hand-written implementation that mirrors what sqlc would generate—so it works out of the box. If you want to regenerate:

```bash
sqlc generate
```

The service layer already uses interfaces, so swapping in the generated code is seamless.

---

## Notes

- **Age calculation** uses standard birthday logic (subtracts a year if birthday hasn't happened yet this year)
- **Request ID** is attached to every request via `X-Request-ID` header
- **Logging** captures method, path, status, duration, and request ID for each request
- **Validation** checks that `name` and `dob` are present and `dob` is in `YYYY-MM-DD` format

---

## What I'd Add Next

- OpenAPI/Swagger docs
- More comprehensive tests (integration, validation errors)
- A proper migration tool like `golang-migrate`
- Makefile for common commands

---

Check out `Reasoning.md` for the full breakdown of design decisions and how this meets the assessment requirements.
