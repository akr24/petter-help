package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/akr24/petter-help/backend/internal/domain"
)

type DogRepository struct {
	pool *pgxpool.Pool
}

func NewDogRepository(pool *pgxpool.Pool) *DogRepository {
	return &DogRepository{pool: pool}
}

func (r *DogRepository) List(ctx context.Context) ([]domain.Dog, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, name, breed, size, personality, needs, purveyor
		FROM dogs
		ORDER BY id DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	dogs := []domain.Dog{}
	for rows.Next() {
		var d domain.Dog
		if err := rows.Scan(&d.ID, &d.Name, &d.Breed, &d.Size, &d.Personality, &d.Needs, &d.Purveyor); err != nil {
			return nil, err
		}
		dogs = append(dogs, d)
	}
	return dogs, rows.Err()
}
