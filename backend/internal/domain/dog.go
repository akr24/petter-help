package domain

import "context"

// Dog is a single adoptable dog listing.
//
// The fields here are placeholders until the team settles the dog categories
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

// DogRepository is the storage boundary for dog listings.
type DogRepository interface {
	List(ctx context.Context) ([]Dog, error)
}
