package token

import (
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/akr24/petter-help/backend/internal/domain"
)

const secret = "0123456789abcdef0123456789abcdef"

func TestIssueAndParse(t *testing.T) {
	iss, err := New(secret, time.Hour)
	if err != nil {
		t.Fatal(err)
	}
	raw, err := iss.Issue(&domain.User{ID: 42, Role: domain.RolePurveyor})
	if err != nil {
		t.Fatal(err)
	}
	c, err := iss.Parse(raw)
	if err != nil {
		t.Fatal(err)
	}
	if c.UserID != 42 || c.Role != domain.RolePurveyor {
		t.Fatalf("got %+v", c)
	}
}

func TestRejectsShortSecret(t *testing.T) {
	if _, err := New("short", time.Hour); err == nil {
		t.Fatal("want error for short secret")
	}
}

func TestRejectsExpired(t *testing.T) {
	iss, _ := New(secret, time.Minute)
	start := time.Now()
	iss.now = func() time.Time { return start }
	raw, _ := iss.Issue(&domain.User{ID: 1, Role: domain.RoleSeeker})

	iss.now = func() time.Time { return start.Add(2 * time.Minute) }
	if _, err := iss.Parse(raw); !errors.Is(err, ErrInvalid) {
		t.Fatalf("want ErrInvalid, got %v", err)
	}
}

func TestRejectsTamperedAndForeign(t *testing.T) {
	iss, _ := New(secret, time.Hour)
	other, _ := New(strings.Repeat("x", 32), time.Hour)
	raw, _ := other.Issue(&domain.User{ID: 1, Role: domain.RoleSeeker})

	for _, bad := range []string{"", "not.a.jwt", raw, raw[:len(raw)-2] + "zz"} {
		if _, err := iss.Parse(bad); err == nil {
			t.Errorf("Parse(%q) should fail", bad)
		}
	}
}
