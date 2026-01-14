package auth

import (
	"testing"
	"regexp"
)

func TestHashPassword(t *testing.T) {
	password := "1234"
	want := regexp.MustCompile("^\\$argon2id\\$")
    hash, err := HashPassword(password)
    if !want.MatchString(hash) || err != nil {
        t.Errorf(`HashPassword("1234") = %s, %v, want match for %#q, nil`, hash, err, want)
    }
}

func TestCheckPasswordHash(t *testing.T) {
	password := "1234"
	hash, err := HashPassword(password)
	if err != nil {
		t.Errorf("Error hashing password: %s", err)
	}
	matches, err := CheckPasswordHash(password, hash)
    if !matches || err != nil {
        t.Errorf(`CheckPasswordHash("%s, %s") = %v, %v, want true, nil`, password, hash, matches, err)
    }
}