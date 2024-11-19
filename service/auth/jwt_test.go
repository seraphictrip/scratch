package auth_test

import (
	"scratch/service/auth"
	"strconv"
	"testing"
)

var CreteJWTTests = []struct {
	secret []byte
	userId int
}{
	{[]byte("secret"), 1},
}

func TestCreteJWT(t *testing.T) {
	for i, e := range CreteJWTTests {
		e := e
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			// code
			actual, err := auth.CreateJWT(e.secret, e.userId)
			if err != nil {
				t.Fatalf("unexpected error, got %v, want %v", err, nil)
			}
			if actual == "" {
				t.Fatalf("expected CreateJWT(%v, %v) to not be empty", e.secret, e.userId)
			}

		})
	}
}
