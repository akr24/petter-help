// Package token issues and verifies the JWTs that identify a signed-in user.
//
// Tokens are HS256-signed with a server secret and carry only the user id,
// role and expiry; anything else is looked up fresh so a role change or
// deleted account takes effect without waiting for the token to expire.
package token

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/akr24/petter-help/backend/internal/domain"
)

const issuer = "petter-help"

var ErrInvalid = errors.New("token: invalid")

// Issuer signs and parses tokens with one secret.
type Issuer struct {
	secret []byte
	ttl    time.Duration
	now    func() time.Time
}

// New requires a secret of at least 32 bytes so HS256 isn't brute-forceable.
func New(secret string, ttl time.Duration) (*Issuer, error) {
	if len(secret) < 32 {
		return nil, fmt.Errorf("token: JWT_SECRET must be at least 32 bytes, got %d", len(secret))
	}
	return &Issuer{secret: []byte(secret), ttl: ttl, now: time.Now}, nil
}

// Claims are what a token carries.
type Claims struct {
	UserID int64
	Role   domain.Role
}

type jwtClaims struct {
	Role domain.Role `json:"role"`
	jwt.RegisteredClaims
}

// Issue returns a signed token for u.
func (i *Issuer) Issue(u *domain.User) (string, error) {
	now := i.now()
	claims := jwtClaims{
		Role: u.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    issuer,
			Subject:   strconv.FormatInt(u.ID, 10),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(i.ttl)),
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(i.secret)
}

// Parse verifies the signature, issuer and expiry and returns the claims.
func (i *Issuer) Parse(raw string) (Claims, error) {
	var c jwtClaims
	_, err := jwt.ParseWithClaims(raw, &c, func(t *jwt.Token) (any, error) {
		return i.secret, nil
	},
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
		jwt.WithIssuer(issuer),
		jwt.WithExpirationRequired(),
		jwt.WithTimeFunc(i.now),
	)
	if err != nil {
		return Claims{}, fmt.Errorf("%w: %v", ErrInvalid, err)
	}
	id, err := strconv.ParseInt(c.Subject, 10, 64)
	if err != nil || !c.Role.Valid() {
		return Claims{}, ErrInvalid
	}
	return Claims{UserID: id, Role: c.Role}, nil
}
