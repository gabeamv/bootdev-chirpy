package test

import (
	"bootdev-chirpy/internal/auth"
	"fmt"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
)

type testCase struct {
	userID    uuid.UUID
	secret    string
	expiresIn time.Duration
	want      string
}

func TestJWT(t *testing.T) {
	id := uuid.New()
	test_cases := map[string]testCase{
		"make-and-validate": {userID: id, secret: "this-is-secret", expiresIn: time.Hour, want: id.String()},
		"expire":            {userID: id, secret: "this-is-secret", expiresIn: (-1 * time.Second), want: fmt.Sprintf("%v", uuid.UUID{})},
	}
	for name, tc := range test_cases {
		t.Run(name, func(t *testing.T) {
			token, _ := auth.MakeJWT(tc.userID, tc.secret, tc.expiresIn)
			got, _ := auth.ValidateJWT(token, tc.secret)
			if diff := cmp.Diff(got.String(), tc.want); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}

func TestWrongSecret(t *testing.T) {
	want := uuid.New()
	realSecret := "real-secret"
	wrongSecret := "wrong-secret"

	token, _ := auth.MakeJWT(want, wrongSecret, time.Hour)
	got, _ := auth.ValidateJWT(token, realSecret)
	if diff := cmp.Diff(got, want); diff == "" {
		t.Fatal(diff)
	}
}
