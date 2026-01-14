package auth

import (
	"fmt"
	"time"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/golang-jwt/jwt/v5"
)

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {

	mySigningKey := []byte(tokenSecret)
	
	claims := &jwt.RegisteredClaims{
		Issuer: "chirpy",
		IssuedAt: &jwt.NumericDate{time.Now()},
		ExpiresAt: &jwt.NumericDate{time.Now().Add(expiresIn)},
		Subject: userID.String(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(mySigningKey)
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func (token *jwt.Token) (any, error) {
		return []byte(tokenSecret), nil
	})

	if err != nil {
		return uuid.Nil, fmt.Errorf("Error validating JWT: %s", err)

	} else if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok {
		userIDString, err := claims.GetSubject()
		if err != nil {
			return uuid.Nil, fmt.Errorf("Error getting token uuid: %s", err)
		}

		id, err := uuid.Parse(userIDString)
		if err != nil {
			return uuid.Nil, fmt.Errorf("Invalid user ID: %w", err)
		}

		return id, nil

	} else {
		return uuid.Nil, fmt.Errorf("Invalid claims")
	}
}

func GetBearerToken(headers http.Header) (string, error) { 

	if len(headers["Authorization"]) < 1 {
		return "", fmt.Errorf("Authorization header is nil")
	}

	authHeader := headers["Authorization"][0]
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return "", fmt.Errorf("Authorization header has incorrect formatting")
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	return tokenString, nil
}

func GetAPIKey(headers http.Header) (string, error) {

	if len(headers["Authorization"]) < 1 {
		return "", fmt.Errorf("Authorization header is nil")
	}

	authHeader := headers["Authorization"][0]
	if !strings.HasPrefix(authHeader, "ApiKey ") {
		return "", fmt.Errorf("Authorization header has incorrect formatting")
	}

	keyString := strings.TrimPrefix(authHeader, "ApiKey ")
	return keyString, nil
}