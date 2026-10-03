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
| POST   | `/api/auth/register` | `{"email", "password", "name", "role"}` role is `seeker` or `purveyor` |
| POST   | `/api/auth/login`    | `{"email", "password"}`                    |
| GET    | `/api/auth/me`       | requires `Authorization: Bearer <token>`   |
| GET    | `/api/seeker/profile`| seeker token; 404 until a profile is saved |
| PUT    | `/api/seeker/profile`| seeker token; full replace, see below      |

Seeker profile body (every field required except the booleans, arrays and
notes, which default to false, `[]` and `""`):

```json
{
  "home_type": "apartment | house | other",
  "has_yard": false,
  "children": "none | young | older",
  "has_dogs": false,
  "has_cats": false,
  "activity_level": "low | moderate | high",
  "hours_alone": 0,
  "experience": "first_time | some | experienced",
  "training_commitment": "low | moderate | high",
  "grooming_commitment": "low | moderate | high",
  "needs_hypoallergenic": false,
  "size_preferences": ["small", "medium", "large"],
  "age_preferences": ["puppy", "adult", "senior"],
  "notes": ""
}
```

The rationale for each field is in `db/migrations/004_seeker_profiles.sql`.

Register and login both return `{"user": {...}, "token": "..."}`. The token
is an HS256 JWT signed with `JWT_SECRET`, valid for 24 hours, carrying only
the user id and role. Passwords are hashed with argon2id (random per-user
salt, PHC-encoded).

The `users` table holds identity and login only; lifestyle profiles,
purveyor details and listings go in their own tables keyed by `users.id`.

## Ports

| Service  | Port |
|----------|------|
| Frontend | 5173 |
| Backend  | 8080 |
| Postgres | 5433 |

Override any of these in a `.env` file (see `.env.example`).
