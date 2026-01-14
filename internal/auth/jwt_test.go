package auth

import (
	"testing"
	"time"
	"net/http"

	"github.com/google/uuid"
)

func TestMakeJWTValidateJWT(t *testing.T) {
	userID := uuid.New()
	tokenSecret := "tokenSecret"
	expiresIn := time.Minute

	token, err := MakeJWT(userID, tokenSecret, expiresIn)
	want := (token != "") && (err == nil)
	if !want {
		t.Errorf("MakeJWT(userID, tokenSecret, expiresIn) returns %s, %s, want string, nil", token, err)
	}

	id, err := ValidateJWT(token, tokenSecret)
	want = (id != uuid.Nil) && (err == nil)
	if !want {
		t.Errorf("MakeValidate(tokenString, tokenSecret) returns %s, %s, want uuid.UUID, nil", id, err)
	}
	
}

func TestGetBearerToken(t *testing.T) {
	validReq, err := http.NewRequest(http.MethodGet, "http://example.com", nil)
	if err != nil {
		t.Errorf("Error creating request: %s", err)
	}
	validReq.Header.Add("Authorization", "Bearer 1234")

	invalidReq, err := http.NewRequest(http.MethodGet, "http://example.com", nil)
	if err != nil {
		t.Errorf("Error creating request: %s", err)
	}
	invalidReq.Header.Add("Authorization", "foobar")

	want := "1234"
	validToken, err := GetBearerToken(validReq.Header)
	if !(validToken == want) || err != nil {
		t.Errorf("GetBearerToken(validRequest) returns %s, %s, wanted %s, %s", validToken, err, "1234", "nil")
	}

	want = ""
	invalidToken, err :=  GetBearerToken(invalidReq.Header)
	if !(invalidToken == want) || err == nil {
		t.Errorf("GetBearerToken(validRequest) returns %s, %s, wanted %s, %s", invalidToken, err, "", "error")
	}

}