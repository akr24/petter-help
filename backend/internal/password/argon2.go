// Package password hashes and verifies passwords with argon2id.
//
// argon2id is the OWASP-recommended algorithm for password storage. Each hash
// uses a fresh random salt and is encoded in the PHC string format, so the
// parameters travel with the hash and can be raised later without breaking
// existing accounts.
package password

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

// Params follow the OWASP minimum for argon2id: 19 MiB memory, 2 iterations,
// 1 lane. Raise memory or time if the host can afford it.
type Params struct {
	Memory  uint32 // KiB
	Time    uint32
	Threads uint8
	SaltLen uint32
	KeyLen  uint32
}

// Default is the parameter set used for new hashes.
var Default = Params{
	Memory:  19 * 1024,
	Time:    2,
	Threads: 1,
	SaltLen: 16,
	KeyLen:  32,
}

var ErrMalformedHash = errors.New("password: malformed hash")

// Hash returns the PHC-encoded argon2id hash of plain using Default params.
func Hash(plain string) (string, error) {
	return HashWith(plain, Default)
}

// HashWith is Hash with explicit parameters.
func HashWith(plain string, p Params) (string, error) {
	salt := make([]byte, p.SaltLen)
	if _, err := rand.Read(salt); err != nil {
		return "", fmt.Errorf("password: generate salt: %w", err)
	}
	key := argon2.IDKey([]byte(plain), salt, p.Time, p.Memory, p.Threads, p.KeyLen)

	b64 := base64.RawStdEncoding
	return fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version, p.Memory, p.Time, p.Threads,
		b64.EncodeToString(salt), b64.EncodeToString(key)), nil
}

// Verify reports whether plain matches encoded. The comparison is constant
// time, and the parameters are read from the hash itself.
func Verify(plain, encoded string) (bool, error) {
	parts := strings.Split(encoded, "$")
	if len(parts) != 6 || parts[1] != "argon2id" {
		return false, ErrMalformedHash
	}

	var version int
	if _, err := fmt.Sscanf(parts[2], "v=%d", &version); err != nil || version != argon2.Version {
		return false, ErrMalformedHash
	}

	var p Params
	if _, err := fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &p.Memory, &p.Time, &p.Threads); err != nil {
		return false, ErrMalformedHash
	}

	b64 := base64.RawStdEncoding
	salt, err := b64.DecodeString(parts[4])
	if err != nil {
		return false, ErrMalformedHash
	}
	want, err := b64.DecodeString(parts[5])
	if err != nil {
		return false, ErrMalformedHash
	}

	got := argon2.IDKey([]byte(plain), salt, p.Time, p.Memory, p.Threads, uint32(len(want)))
	return subtle.ConstantTimeCompare(got, want) == 1, nil
}
