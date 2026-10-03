package domain

import (
	"context"
	"time"
)

// SeekerProfile is a seeker's lifestyle questionnaire: the inputs the
// matching algorithm scores against dog listings. One per seeker.
//
// The option sets below are a first proposal; the team hasn't settled the
// questionnaire yet, so expect them to change.
type SeekerProfile struct {
	UserID int64 `json:"user_id"`

	// Household
	HomeType HomeType `json:"home_type"`
	HasYard  bool     `json:"has_yard"`
	Children Children `json:"children"`
	HasDogs  bool     `json:"has_dogs"`
	HasCats  bool     `json:"has_cats"`

	// Lifestyle and capacity
	ActivityLevel       Level      `json:"activity_level"`
	HoursAlone          int        `json:"hours_alone"`
	Experience          Experience `json:"experience"`
	TrainingCommitment  Level      `json:"training_commitment"`
	GroomingCommitment  Level      `json:"grooming_commitment"`
	NeedsHypoallergenic bool       `json:"needs_hypoallergenic"`

	// Soft preferences; empty means no preference
	SizePreferences []Size   `json:"size_preferences"`
	AgePreferences  []DogAge `json:"age_preferences"`

	Notes string `json:"notes"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type HomeType string

const (
	HomeApartment HomeType = "apartment"
	HomeHouse     HomeType = "house"
	HomeOther     HomeType = "other"
)

func (h HomeType) Valid() bool { return h == HomeApartment || h == HomeHouse || h == HomeOther }

// Children: "young" means under six, where a dog's tolerance matters most.
type Children string

const (
	ChildrenNone  Children = "none"
	ChildrenYoung Children = "young"
	ChildrenOlder Children = "older"
)

func (c Children) Valid() bool { return c == ChildrenNone || c == ChildrenYoung || c == ChildrenOlder }

// Level is a three-step scale reused for activity, training and grooming.
type Level string

const (
	LevelLow      Level = "low"
	LevelModerate Level = "moderate"
	LevelHigh     Level = "high"
)

func (l Level) Valid() bool { return l == LevelLow || l == LevelModerate || l == LevelHigh }

type Experience string

const (
	ExperienceFirstTime   Experience = "first_time"
	ExperienceSome        Experience = "some"
	ExperienceExperienced Experience = "experienced"
)

func (e Experience) Valid() bool {
	return e == ExperienceFirstTime || e == ExperienceSome || e == ExperienceExperienced
}

// Size and DogAge are shared with dog listings so preferences and listings
// compare directly.
type Size string

const (
	SizeSmall  Size = "small"
	SizeMedium Size = "medium"
	SizeLarge  Size = "large"
)

func (s Size) Valid() bool { return s == SizeSmall || s == SizeMedium || s == SizeLarge }

type DogAge string

const (
	DogAgePuppy  DogAge = "puppy"
	DogAgeAdult  DogAge = "adult"
	DogAgeSenior DogAge = "senior"
)

func (a DogAge) Valid() bool { return a == DogAgePuppy || a == DogAgeAdult || a == DogAgeSenior }

// SeekerProfileRepository is the storage boundary for profiles.
type SeekerProfileRepository interface {
	// Upsert inserts or replaces the profile for p.UserID and fills timestamps.
	Upsert(ctx context.Context, p *SeekerProfile) error
	FindByUserID(ctx context.Context, userID int64) (*SeekerProfile, error)
}
