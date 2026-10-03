package password

import (
	"strings"
	"testing"
)

func TestHashAndVerify(t *testing.T) {
	h, err := Hash("correct horse battery staple")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(h, "$argon2id$v=19$m=19456,t=2,p=1$") {
		t.Fatalf("unexpected encoding: %s", h)
	}

	ok, err := Verify("correct horse battery staple", h)
	if err != nil || !ok {
		t.Fatalf("want match, got ok=%v err=%v", ok, err)
	}

	ok, err = Verify("wrong", h)
	if err != nil || ok {
		t.Fatalf("want mismatch, got ok=%v err=%v", ok, err)
	}
}

func TestHashUsesFreshSalt(t *testing.T) {
	a, _ := Hash("same")
	b, _ := Hash("same")
	if a == b {
		t.Fatal("two hashes of the same password should differ")
	}
}

func TestVerifyRejectsMalformed(t *testing.T) {
	for _, bad := range []string{"", "plain", "$bcrypt$x$y$z$w", "$argon2id$v=19$m=1,t=1,p=1$!!$!!"} {
		if _, err := Verify("x", bad); err == nil {
			t.Errorf("Verify(%q) should error", bad)
		}
	}
}
