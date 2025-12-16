# GoFiber User API (DOB + Age)

A small REST API in Go using Fiber, PostgreSQL, sqlc, Zap, and go-playground/validator. Users have `name` and `dob`; responses compute `age` on the fly.

## Stack
- Fiber
- PostgreSQL + sqlc
- Zap logging
- go-playground/validator

## Quick Start

### 1) Prepare PostgreSQL
Create a local DB and apply migration:

```sql
CREATE DATABASE golang_assessment;
\c golang_assessment;
CREATE TABLE IF NOT EXISTS users (
  id SERIAL PRIMARY KEY,
  name TEXT NOT NULL,
  dob DATE NOT NULL
);
```

Or run the file at `db/migrations/001_init.sql`.

### 2) Configure Environment
Set `DATABASE_URL` if needed:

```bash
export DATABASE_URL=postgres://postgres:postgres@localhost:5432/golang_assessment?sslmode=disable
```
On Windows PowerShell:
```powershell
$env:DATABASE_URL = "postgres://postgres:postgres@localhost:5432/golang_assessment?sslmode=disable"
```

### 3) Install deps and run
```bash
go mod tidy
go run ./cmd/server
```
Server listens on `:8080`.

### Docker (app + Postgres)
```bash
docker compose up --build
```
App on `:8080`, Postgres on `:5432` with `dev_user/mysecretpassword`.

## Endpoints
- POST /users
- GET /users/:id
- PUT /users/:id
- DELETE /users/:id
- GET /users (supports `limit` and `offset`)

### Examples
Create:
```bash
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{"name":"Alice","dob":"1990-05-10"}'
```

Get:
```bash
curl http://localhost:8080/users/1
```

List:
```bash
curl http://localhost:8080/users
```

## sqlc
Config is at `db/sqlc/sqlc.yaml` and queries at `internal/repository/queries.sql`.
Generate code:
```bash
sqlc generate
```
Then replace `SimpleRepo` wiring with the generated package types.

## Notes
- Age is calculated with the standard birthday logic.
- Request ID is injected via `X-Request-ID` header.
- Request duration + metadata are logged with Zap.

## Optional: Docker
Add a `docker-compose.yml` with Postgres and an app service if desired.
