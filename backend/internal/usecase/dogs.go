package usecase

import (
	"context"

	"github.com/akr24/petter-help/backend/internal/domain"
)

// Dogs holds the dog-listing use cases.
type Dogs struct {
	repo domain.DogRepository
}

func NewDogs(repo domain.DogRepository) *Dogs {
	return &Dogs{repo: repo}
}

// List returns every listing. Ranking by seeker profile comes later.
func (d *Dogs) List(ctx context.Context) ([]domain.Dog, error) {
	return d.repo.List(ctx)
}
