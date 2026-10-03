package usecase

import (
	"context"
	"errors"
	"fmt"
	"net/mail"
	"strings"

	"github.com/akr24/petter-help/backend/internal/domain"
	"github.com/akr24/petter-help/backend/internal/password"
)

// MinPasswordLen follows NIST 800-63B: length is the main defence, so no
// composition rules, but at least 8 characters.
const MinPasswordLen = 8

var ErrInvalidInput = errors.New("invalid input")

// Auth holds the account use cases.
type Auth struct {
	users domain.UserRepository
}

func NewAuth(users domain.UserRepository) *Auth {
	return &Auth{users: users}
}

// Register creates an account. The email is normalised to lower case so the
// unique index catches case-variant duplicates.
func (a *Auth) Register(ctx context.Context, email, plain string, role domain.Role) (*domain.User, error) {
	email, err := normaliseEmail(email)
	if err != nil {
		return nil, err
	}
	if len(plain) < MinPasswordLen {
		return nil, fmt.Errorf("%w: password must be at least %d characters", ErrInvalidInput, MinPasswordLen)
	}
	if !role.Valid() {
		return nil, fmt.Errorf("%w: role must be seeker or purveyor", ErrInvalidInput)
	}

	hash, err := password.Hash(plain)
	if err != nil {
		return nil, err
	}

	u := &domain.User{Email: email, Role: role, PasswordHash: hash}
	if err := a.users.Create(ctx, u); err != nil {
		return nil, err
	}
	return u, nil
}

// Login verifies the credentials and returns the user. Unknown emails and
// wrong passwords both yield ErrInvalidCredentials so callers can't tell
// which accounts exist.
func (a *Auth) Login(ctx context.Context, email, plain string) (*domain.User, error) {
	email, err := normaliseEmail(email)
	if err != nil {
		return nil, domain.ErrInvalidCredentials
	}

	u, err := a.users.FindByEmail(ctx, email)
	if errors.Is(err, domain.ErrNotFound) {
		// Still run a hash so timing doesn't reveal whether the email exists.
		_, _ = password.Verify(plain, dummyHash)
		return nil, domain.ErrInvalidCredentials
	}
	if err != nil {
		return nil, err
	}

	ok, err := password.Verify(plain, u.PasswordHash)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, domain.ErrInvalidCredentials
	}
	return u, nil
}

func normaliseEmail(email string) (string, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	addr, err := mail.ParseAddress(email)
	if err != nil || addr.Address != email {
		return "", fmt.Errorf("%w: invalid email", ErrInvalidInput)
	}
	return email, nil
}

// dummyHash is a valid argon2id hash of a throwaway string, used only to
// equalise login timing for unknown emails.
var dummyHash = func() string {
	h, err := password.Hash("not-a-real-password")
	if err != nil {
		panic(err)
	}
	return h
}()
