// Command server is the PetterHelp HTTP API.
//
// It wires the layers together: Postgres adapters implement the domain
// repositories, use cases hold the business rules, and the Echo adapter
// exposes them over HTTP. Nothing else in the tree imports cmd/server.
package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	httpadapter "github.com/akr24/petter-help/backend/internal/adapter/http"
	"github.com/akr24/petter-help/backend/internal/adapter/postgres"
	"github.com/akr24/petter-help/backend/internal/usecase"
)

func main() {
	addr := ":" + envOr("PORT", "8080")

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	pool, err := postgres.Connect(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("database: %v", err)
	}
	defer pool.Close()

	e := httpadapter.New(httpadapter.Deps{
		DB:         pool,
		Dogs:       usecase.NewDogs(postgres.NewDogRepository(pool)),
		Auth:       usecase.NewAuth(postgres.NewUserRepository(pool)),
		CORSOrigin: envOr("CORS_ORIGIN", "*"),
	})
	e.Server.ReadHeaderTimeout = 5 * time.Second

	go func() {
		log.Printf("listening on %s", addr)
		if err := e.Start(addr); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("server: %v", err)
		}
	}()

	<-ctx.Done()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(shutdownCtx); err != nil {
		log.Printf("shutdown: %v", err)
	}
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
