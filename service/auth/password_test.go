package auth_test

import (
	"os"
	"scratch/service/auth"
	"strconv"
	"testing"
)

var HashPasswordTests = []struct {
	password string
}{
	{"pass"},
	{"password"},
	{"secret"},
}

func TestHashPassword(t *testing.T) {
	for i, e := range HashPasswordTests {
		e := e
		os.Environ()
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			// code
			actual, _ := auth.HashPassword(e.password)

			if actual == e.password || actual == "" {
				t.Fatalf("HashPassword(%v) = %v, want hashed", e.password, actual)
			}

		})
	}
}

var ComparePasswordsTests = []struct {
	password string
	expected bool
}{
	{"pass", true},
	{"password", true},
	{"otherpassword", true},
	{"!@#$%^&*()_+", true},
	{"1234567890", true},
	{"abcdefghijklmnopqrstuvwxyz1234567890!@#$%^&*()_+", true},
	// known to not work on passwords longer then 72 bytes
	{"abcdefghijklmnopqrstuvwxyz1234567890!@#$%^&*()_+abcdefghijklmnopqrstuvwxyz1234567890!@#$%^&*()_+abcdefghijklmnopqrstuvwxyz1234567890!@#$%^&*()_+", false},
}

func TestComparePasswords(t *testing.T) {
	for i, e := range ComparePasswordsTests {
		e := e
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			// code
			hashed, _ := auth.HashPassword(e.password)
			match := auth.ComparePasswords(hashed, []byte(e.password))

			if match != e.expected {
				t.Fatalf("ComparePasswords(%v, %v) = %v, want %v", hashed, []byte(e.password), match, e.expected)
			}
		})
	}
}
