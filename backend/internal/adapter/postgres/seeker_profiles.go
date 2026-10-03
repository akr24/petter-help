package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/akr24/petter-help/backend/internal/domain"
)

type SeekerProfileRepository struct {
	pool *pgxpool.Pool
}

func NewSeekerProfileRepository(pool *pgxpool.Pool) *SeekerProfileRepository {
	return &SeekerProfileRepository{pool: pool}
}

func (r *SeekerProfileRepository) Upsert(ctx context.Context, p *domain.SeekerProfile) error {
	return r.pool.QueryRow(ctx, `
		INSERT INTO seeker_profiles (
			user_id, home_type, has_yard, children, has_dogs, has_cats,
			activity_level, hours_alone, experience, training_commitment,
			grooming_commitment, needs_hypoallergenic,
			size_preferences, age_preferences, notes
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
		ON CONFLICT (user_id) DO UPDATE SET
			home_type = EXCLUDED.home_type,
			has_yard = EXCLUDED.has_yard,
			children = EXCLUDED.children,
			has_dogs = EXCLUDED.has_dogs,
			has_cats = EXCLUDED.has_cats,
			activity_level = EXCLUDED.activity_level,
			hours_alone = EXCLUDED.hours_alone,
			experience = EXCLUDED.experience,
			training_commitment = EXCLUDED.training_commitment,
			grooming_commitment = EXCLUDED.grooming_commitment,
			needs_hypoallergenic = EXCLUDED.needs_hypoallergenic,
			size_preferences = EXCLUDED.size_preferences,
			age_preferences = EXCLUDED.age_preferences,
			notes = EXCLUDED.notes,
			updated_at = now()
		RETURNING created_at, updated_at`,
		p.UserID, p.HomeType, p.HasYard, p.Children, p.HasDogs, p.HasCats,
		p.ActivityLevel, p.HoursAlone, p.Experience, p.TrainingCommitment,
		p.GroomingCommitment, p.NeedsHypoallergenic,
		toStrings(p.SizePreferences), toStrings(p.AgePreferences), p.Notes,
	).Scan(&p.CreatedAt, &p.UpdatedAt)
}

func (r *SeekerProfileRepository) FindByUserID(ctx context.Context, userID int64) (*domain.SeekerProfile, error) {
	var (
		p     domain.SeekerProfile
		sizes []string
		ages  []string
	)
	err := r.pool.QueryRow(ctx, `
		SELECT user_id, home_type, has_yard, children, has_dogs, has_cats,
		       activity_level, hours_alone, experience, training_commitment,
		       grooming_commitment, needs_hypoallergenic,
		       size_preferences, age_preferences, notes, created_at, updated_at
		FROM seeker_profiles
		WHERE user_id = $1`, userID,
	).Scan(
		&p.UserID, &p.HomeType, &p.HasYard, &p.Children, &p.HasDogs, &p.HasCats,
		&p.ActivityLevel, &p.HoursAlone, &p.Experience, &p.TrainingCommitment,
		&p.GroomingCommitment, &p.NeedsHypoallergenic,
		&sizes, &ages, &p.Notes, &p.CreatedAt, &p.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	p.SizePreferences = fromStrings[domain.Size](sizes)
	p.AgePreferences = fromStrings[domain.DogAge](ages)
	return &p, nil
}

// pgx maps TEXT[] to []string; these convert to and from the typed slices.
func toStrings[T ~string](xs []T) []string {
	out := make([]string, len(xs))
	for i, x := range xs {
		out[i] = string(x)
	}
	return out
}

func fromStrings[T ~string](xs []string) []T {
	out := make([]T, len(xs))
	for i, x := range xs {
		out[i] = T(x)
	}
	return out
}
