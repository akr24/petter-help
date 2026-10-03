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

// User is an account holder. PasswordHash is the encoded argon2id string and
// is never serialised to JSON.
type User struct {
	ID           int64     `json:"id"`
	Email        string    `json:"email"`
	Role         Role      `json:"role"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
}

var (
	ErrEmailTaken         = errors.New("email already registered")
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrNotFound           = errors.New("not found")
)

// UserRepository is the storage boundary for accounts.
type UserRepository interface {
	Create(ctx context.Context, u *User) error
	FindByEmail(ctx context.Context, email string) (*User, error)
}
