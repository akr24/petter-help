// Package db wraps the Postgres connection pool and the queries the API uses.
package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Connect opens a pool against url and verifies it with a ping.
func Connect(ctx context.Context, url string) (*pgxpool.Pool, error) {
	if url == "" {
		return nil, fmt.Errorf("DATABASE_URL is not set")
	}
	pool, err := pgxpool.New(ctx, url)
	if err != nil {
		return nil, err
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, err
	}
	return pool, nil
}

// Dog is a single adoptable dog listing.
//
// The columns here are placeholders until the team settles the dog categories
// and questionnaire traits the matching algorithm will score on.
type Dog struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Breed       string `json:"breed"`
	Size        string `json:"size"`
	Personality string `json:"personality"`
	Needs       string `json:"needs"`
	Purveyor    string `json:"purveyor"`
}

// ListDogs returns every dog listing, newest first.
func ListDogs(ctx context.Context, pool *pgxpool.Pool) ([]Dog, error) {
	rows, err := pool.Query(ctx, `
		SELECT id, name, breed, size, personality, needs, purveyor
		FROM dogs
		ORDER BY id DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	dogs := []Dog{}
	for rows.Next() {
		var d Dog
		if err := rows.Scan(&d.ID, &d.Name, &d.Breed, &d.Size, &d.Personality, &d.Needs, &d.Purveyor); err != nil {
			return nil, err
		}
		dogs = append(dogs, d)
	}
	return dogs, rows.Err()
}
