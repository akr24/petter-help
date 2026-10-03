package usecase

import (
	"errors"
	"testing"

	"github.com/akr24/petter-help/backend/internal/domain"
)

func validProfile() domain.SeekerProfile {
	return domain.SeekerProfile{
		HomeType:           domain.HomeApartment,
		Children:           domain.ChildrenNone,
		ActivityLevel:      domain.LevelModerate,
		HoursAlone:         6,
		Experience:         domain.ExperienceSome,
		TrainingCommitment: domain.LevelModerate,
		GroomingCommitment: domain.LevelLow,
		SizePreferences:    []domain.Size{domain.SizeSmall, domain.SizeSmall, domain.SizeMedium},
		AgePreferences:     nil,
		Notes:              "  works from home  ",
	}
}

func TestValidateProfileNormalises(t *testing.T) {
	p := validProfile()
	if err := validateProfile(&p); err != nil {
		t.Fatal(err)
	}
	if len(p.SizePreferences) != 2 {
		t.Errorf("want duplicates removed, got %v", p.SizePreferences)
	}
	if p.AgePreferences == nil {
		t.Error("want non-nil empty slice for age_preferences")
	}
	if p.Notes != "works from home" {
		t.Errorf("want trimmed notes, got %q", p.Notes)
	}
}

func TestValidateProfileRejects(t *testing.T) {
	cases := map[string]func(*domain.SeekerProfile){
		"home_type":   func(p *domain.SeekerProfile) { p.HomeType = "boat" },
		"children":    func(p *domain.SeekerProfile) { p.Children = "" },
		"hours_alone": func(p *domain.SeekerProfile) { p.HoursAlone = 25 },
		"experience":  func(p *domain.SeekerProfile) { p.Experience = "pro" },
		"size":        func(p *domain.SeekerProfile) { p.SizePreferences = []domain.Size{"giant"} },
		"age":         func(p *domain.SeekerProfile) { p.AgePreferences = []domain.DogAge{"teen"} },
		"notes":       func(p *domain.SeekerProfile) { p.Notes = string(make([]byte, MaxNotesLen+1)) },
	}
	for name, mutate := range cases {
		p := validProfile()
		mutate(&p)
		if err := validateProfile(&p); !errors.Is(err, ErrInvalidInput) {
			t.Errorf("%s: want ErrInvalidInput, got %v", name, err)
		}
	}
}
