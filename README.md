# PetterHelp

Matches dog seekers in Cambridge, MA with adoptable dogs from shelters,
rescues, breeders and individuals, based on lifestyle and personality rather
than breed or size alone.

## Stack

- **Backend:** Go 1.25, `net/http`, pgx
- **Frontend:** React + TypeScript via Vite
- **Database:** Postgres 18
- **Local infra:** Docker Compose, Makefile

## Quick start

```sh
make up        # builds images, starts Postgres, backend and frontend
open http://localhost:5173
```

For faster iteration, run the database in Docker and the apps natively:

```sh
make install   # once
make dev       # Postgres in Docker, Go + Vite with hot reload
```

`make help` lists every target.

## Layout

```
backend/
  cmd/server/            entrypoint; wires the layers together
  internal/domain/       entities and repository interfaces (no dependencies)
  internal/usecase/      business rules (depends on domain only)
  internal/adapter/
    http/                Echo handlers and routing
    postgres/            pgx implementations of the repositories
  internal/password/     argon2id hashing
frontend/                Vite React app
db/migrations/           SQL, applied on first database start
docker-compose.db.yml    Postgres only
docker-compose.yml       backend + frontend (includes the db file)
```

Dependencies point inward: adapters import use cases and domain, use cases
import domain, domain imports nothing from the project.

Migrations run only when the Postgres volume is first created. After adding
one, `make db-reset` drops the volume so it is applied on the next `make up`.

## API

| Method | Path                 | Body                                       |
|--------|----------------------|--------------------------------------------|
| GET    | `/healthz`           |                                            |
| GET    | `/api/dogs`          |                                            |
| POST   | `/api/auth/register` | `{"email", "password", "role"}` role is `seeker` or `purveyor` |
| POST   | `/api/auth/login`    | `{"email", "password"}`                    |

Passwords are hashed with argon2id (random per-user salt, PHC-encoded).
Login does not yet issue a session or token.

## Ports

| Service  | Port |
|----------|------|
| Frontend | 5173 |
| Backend  | 8080 |
| Postgres | 5433 |

Override any of these in a `.env` file (see `.env.example`).
