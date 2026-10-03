package usecase

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"github.com/akr24/petter-help/backend/internal/domain"
)

// MaxNotesLen bounds the free-text field.
const MaxNotesLen = 2000

// SeekerProfiles holds the lifestyle-profile use cases.
type SeekerProfiles struct {
	repo domain.SeekerProfileRepository
}

func NewSeekerProfiles(repo domain.SeekerProfileRepository) *SeekerProfiles {
	return &SeekerProfiles{repo: repo}
}

// Get returns the caller's profile, or domain.ErrNotFound if they haven't
// filled one in yet.
func (s *SeekerProfiles) Get(ctx context.Context, userID int64) (*domain.SeekerProfile, error) {
	return s.repo.FindByUserID(ctx, userID)
}

// Save validates and stores the whole profile for userID. It is a full
// replace, not a patch, so the client always sends every answer.
func (s *SeekerProfiles) Save(ctx context.Context, userID int64, p domain.SeekerProfile) (*domain.SeekerProfile, error) {
	p.UserID = userID
	if err := validateProfile(&p); err != nil {
		return nil, err
	}
	if err := s.repo.Upsert(ctx, &p); err != nil {
		return nil, err
	}
	return &p, nil
}

func validateProfile(p *domain.SeekerProfile) error {
	bad := func(field, allowed string) error {
		return fmt.Errorf("%w: %s must be one of %s", ErrInvalidInput, field, allowed)
	}
	switch {
	case !p.HomeType.Valid():
		return bad("home_type", "apartment, house, other")
	case !p.Children.Valid():
		return bad("children", "none, young, older")
	case !p.ActivityLevel.Valid():
		return bad("activity_level", "low, moderate, high")
	case p.HoursAlone < 0 || p.HoursAlone > 24:
		return fmt.Errorf("%w: hours_alone must be between 0 and 24", ErrInvalidInput)
	case !p.Experience.Valid():
		return bad("experience", "first_time, some, experienced")
	case !p.TrainingCommitment.Valid():
		return bad("training_commitment", "low, moderate, high")
	case !p.GroomingCommitment.Valid():
		return bad("grooming_commitment", "low, moderate, high")
	}

	p.SizePreferences = dedupe(p.SizePreferences)
	for _, sz := range p.SizePreferences {
		if !sz.Valid() {
			return bad("size_preferences", "small, medium, large")
		}
	}
	p.AgePreferences = dedupe(p.AgePreferences)
	for _, a := range p.AgePreferences {
		if !a.Valid() {
			return bad("age_preferences", "puppy, adult, senior")
		}
	}

	p.Notes = strings.TrimSpace(p.Notes)
	if len(p.Notes) > MaxNotesLen {
		return fmt.Errorf("%w: notes must be at most %d characters", ErrInvalidInput, MaxNotesLen)
	}
	return nil
}

// dedupe returns a non-nil copy of xs with duplicates removed, so JSON
// output is [] rather than null and the DB array check sees unique values.
func dedupe[T comparable](xs []T) []T {
	out := make([]T, 0, len(xs))
	for _, x := range xs {
		if !slices.Contains(out, x) {
			out = append(out, x)
		}
	}
	return out
}
