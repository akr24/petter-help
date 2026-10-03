package domain

import (
	"context"
	"errors"
	"time"
)

// Role distinguishes the two account types the sprint calls for.
type Role string

const (
	RoleSeeker   Role = "seeker"
	RolePurveyor Role = "purveyor"
)

// Valid reports whether r is one of the known roles.
func (r Role) Valid() bool {
	return r == RoleSeeker || r == RolePurveyor
}

// User is an account holder. Anything beyond identity and login (lifestyle
// profile, purveyor details, listings) lives in its own table keyed by ID.
// PasswordHash is the encoded argon2id string and is never serialised.
type User struct {
	ID           int64     `json:"id"`
	Email        string    `json:"email"`
	Name         string    `json:"name"`
	Role         Role      `json:"role"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
}

var (
	ErrEmailTaken         = errors.New("email already registered")
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrUnauthenticated    = errors.New("unauthenticated")
	ErrNotFound           = errors.New("not found")
)

// UserRepository is the storage boundary for accounts.
type UserRepository interface {
	Create(ctx context.Context, u *User) error
	FindByID(ctx context.Context, id int64) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
}
