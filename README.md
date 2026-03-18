[![Go Report Card](https://goreportcard.com/badge/github.com/joyyth/go-boilerplate)](https://goreportcard.com/report/github.com/joyyth/go-boilerplate)

A minimal Go + Reactjs boilerplate. Comes with JWT auth, PostgreSQL, and a clean layered architecture — ready to build on.

---

## Tech Stack

### Backend

|             |                                                                       |
| ----------- | --------------------------------------------------------------------- |
| Language    | Go                                                                    |
| Router      | [Chi](https://github.com/go-chi/chi)                                  |
| Database    | PostgreSQL via [pgx](https://github.com/jackc/pgx)                    |
| Migrations  | [Goose](https://github.com/pressly/goose)                             |
| Auth        | JWT (access + refresh tokens) + bcrypt                                |
| Config      | [Koanf](https://github.com/knadh/koanf) + godotenv                    |
| Logging     | [Zerolog](https://github.com/rs/zerolog)                              |
| Validation  | [go-playground/validator](https://github.com/go-playground/validator) |
| Live Reload | [Air](https://github.com/air-verse/air)                               |
| Task Runner | [Task](https://taskfile.dev)                                          |

## Prerequisites

- Go 1.25+
- PostgreSQL
- [Task](https://taskfile.dev/#/installation)

---

## Getting Started

### 1. Clone the repository

```bash
git clone https://github.com/joyyth/go-boilerplate.git my-app
cd my-app
```

### 2. Database Setup

Ensure you have a **PostgreSQL** database running.

```sql
CREATE DATABASE yourapp;
```

### 3. Backend Setup

```bash
cd backend

# copy and configure environment variables
cp .env.sample .env

# start with live reload — migrations run automatically on startup
task dev
```

---

## Using This Boilerplate

When starting a new project:

1. **Rename the module** — find and replace `github.com/joyyth/go-boilerplate` with your module path in `go.mod` and all imports
2. **Rename the env prefix** — find and replace `YOURAPP_` with your app name (e.g. `MYSHOP_`) across all files
3. **Update `ServiceName`** in `cmd/main.go`
4. **Add your domain** — follow the existing `user` pattern: model → repository → service → handler → register in `routes.go` and `container.go`

---

## Backend Tasks

```bash
task dev                  # run with live reload
task run                  # run without live reload
task build                # build binary to ./bin/backend
task test                 # run all tests
task tidy                 # format and tidy dependencies

task migrations:install   # install goose CLI
task migrations:new       # create a new migration  (name=migration_name)
task migrations:up        # apply all pending migrations  (DB_DSN=...)
task migrations:status    # show migration status  (DB_DSN=...)
```
