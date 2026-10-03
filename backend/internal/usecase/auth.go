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

// MaxNameLen keeps display names to something a UI can render.
const MaxNameLen = 100

var ErrInvalidInput = errors.New("invalid input")

// TokenIssuer turns an authenticated user into a bearer token. The usecase
// layer depends on this interface, not on a JWT library.
type TokenIssuer interface {
	Issue(u *domain.User) (string, error)
}

// Session is what a successful register or login returns.
type Session struct {
	User  *domain.User `json:"user"`
	Token string       `json:"token"`
}

// RegisterInput is everything needed to create an account.
type RegisterInput struct {
	Email    string
	Password string
	Name     string
	Role     domain.Role
}

// Auth holds the account use cases.
type Auth struct {
	users  domain.UserRepository
	tokens TokenIssuer
}

func NewAuth(users domain.UserRepository, tokens TokenIssuer) *Auth {
	return &Auth{users: users, tokens: tokens}
}

// Register creates an account and signs the new user in. The email is
// normalised to lower case so the unique index catches case-variant duplicates.
func (a *Auth) Register(ctx context.Context, in RegisterInput) (*Session, error) {
	email, err := normaliseEmail(in.Email)
	if err != nil {
		return nil, err
	}
	if len(in.Password) < MinPasswordLen {
		return nil, fmt.Errorf("%w: password must be at least %d characters", ErrInvalidInput, MinPasswordLen)
	}
	name := strings.TrimSpace(in.Name)
	if name == "" || len(name) > MaxNameLen {
		return nil, fmt.Errorf("%w: name is required and at most %d characters", ErrInvalidInput, MaxNameLen)
	}
	if !in.Role.Valid() {
		return nil, fmt.Errorf("%w: role must be seeker or purveyor", ErrInvalidInput)
	}

	hash, err := password.Hash(in.Password)
	if err != nil {
		return nil, err
	}

	u := &domain.User{Email: email, Name: name, Role: in.Role, PasswordHash: hash}
	if err := a.users.Create(ctx, u); err != nil {
		return nil, err
	}
	return a.session(u)
}

// Login verifies the credentials and signs the user in. Unknown emails and
// wrong passwords both yield ErrInvalidCredentials so callers can't tell
// which accounts exist.
func (a *Auth) Login(ctx context.Context, email, plain string) (*Session, error) {
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
	return a.session(u)
}

// Current loads the account behind a verified token's user id. A deleted
// account yields ErrUnauthenticated so stale tokens stop working.
func (a *Auth) Current(ctx context.Context, userID int64) (*domain.User, error) {
	u, err := a.users.FindByID(ctx, userID)
	if errors.Is(err, domain.ErrNotFound) {
		return nil, domain.ErrUnauthenticated
	}
	return u, err
}

func (a *Auth) session(u *domain.User) (*Session, error) {
	tok, err := a.tokens.Issue(u)
	if err != nil {
		return nil, err
	}
	return &Session{User: u, Token: tok}, nil
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
