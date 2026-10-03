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
backend/     Go API (cmd/server is the entrypoint)
frontend/    Vite React app
db/          SQL migrations, applied on first database start
docker-compose.db.yml   Postgres only
docker-compose.yml      backend + frontend (includes the db file)
```

## Ports

| Service  | Port |
|----------|------|
| Frontend | 5173 |
| Backend  | 8080 |
| Postgres | 5433 |

Override any of these in a `.env` file (see `.env.example`).
