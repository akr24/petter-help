# PetterHelp local development.
#
#   make db-up      start Postgres only (port 5433)
#   make up         build and start db + backend + frontend in Docker
#   make down       stop everything (data volume is kept)
#   make dev        run backend and frontend natively against the Docker db
#   make help       list every target

SHELL := /bin/bash
COMPOSE    := docker compose
COMPOSE_DB := docker compose -f docker-compose.db.yml

.DEFAULT_GOAL := help

## ---------- Docker ----------

.PHONY: db-up
db-up: ## Start Postgres in Docker
	$(COMPOSE_DB) up -d --wait

.PHONY: db-down
db-down: ## Stop Postgres (keeps data)
	$(COMPOSE_DB) down

.PHONY: db-reset
db-reset: ## Stop Postgres and delete its data volume (re-runs migrations next start)
	$(COMPOSE_DB) down -v

.PHONY: psql
psql: ## Open a psql shell in the Postgres container
	$(COMPOSE_DB) exec db psql -U petter -d petter_help

.PHONY: build
build: ## Build backend and frontend images
	$(COMPOSE) build

.PHONY: up
up: ## Build and start db + backend + frontend
	$(COMPOSE) up -d --build --wait

.PHONY: down
down: ## Stop all containers (keeps data)
	$(COMPOSE) down

.PHONY: logs
logs: ## Tail logs from all containers
	$(COMPOSE) logs -f

.PHONY: ps
ps: ## Show container status
	$(COMPOSE) ps

## ---------- Native (faster iteration) ----------

export DATABASE_URL ?= postgres://petter:petter@localhost:5433/petter_help?sslmode=disable

.PHONY: backend
backend: ## Run the Go API natively on :8080 (needs `make db-up`)
	cd backend && go run ./cmd/server

.PHONY: frontend
frontend: ## Run the Vite dev server natively on :5173
	cd frontend && npm run dev

.PHONY: dev
dev: db-up ## Start db in Docker, then backend and frontend natively
	$(MAKE) -j2 backend frontend

.PHONY: install
install: ## Download Go modules and npm packages
	cd backend && go mod download
	cd frontend && npm install

.PHONY: test
test: ## Run backend tests and frontend lint
	cd backend && go test ./...
	cd frontend && npm run lint

.PHONY: help
help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## ' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-10s\033[0m %s\n", $$1, $$2}'
