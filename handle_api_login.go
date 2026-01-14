package main

import (
	"fmt"
	"encoding/json"
	"net/http"
	"time"

	"github.com/tylerapear/chirpy/internal/auth"
	"github.com/tylerapear/chirpy/internal/database"
)

type loginParams struct{
	Password string `json:"password"`
	Email string `json:"email"`
}

type response struct{
	User
	Token string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

func (cfg apiConfig) handle_api_login (w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
    params := loginParams{}
    err := decoder.Decode(&params)
    if err != nil {
        respondWithError(w, 500, fmt.Sprintf("Something went wrong: %s\n", err))
        return
    }

	user, err := cfg.dbQueries.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, 404, fmt.Sprintf("User not found"))
        return
	}

	matches, err := auth.CheckPasswordHash(params.Password, user.HashedPassword.String)
	if err != nil {
		respondWithError(w, 500, fmt.Sprintf("Error checking password hash: %s", err))
        return
	}

	if !matches {
		respondWithError(w, 401, fmt.Sprintf("incorrect email or password"))
		return
	}

	jwtExpiresIn := time.Hour
	refreshTokenExpiresAt := time.Now().Add(60 * 24 * time.Hour)

	refresh_token, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(w, 500, fmt.Sprintf("Error creating refresh token: %s", err))
	}
	
	createRefreshTokenParams := database.CreateRefreshTokenParams{
		Token: refresh_token,
		UserID: user.ID,
		ExpiresAt: refreshTokenExpiresAt,
	}
	createdRefreshToken, err := cfg.dbQueries.CreateRefreshToken(r.Context(), createRefreshTokenParams)
	if err != nil {
		respondWithError(w, 500, fmt.Sprintf("Error creating refresh token: %s", err))
	}

	jwtToken, err := auth.MakeJWT(user.ID, cfg.jwtSecret, jwtExpiresIn)
	if err != nil {
		respondWithError(w, 500, fmt.Sprintf("Error creating token: %s", err))
	}

	tokenResp := response{
		User: User{
			ID: user.ID,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Email: user.Email,
			IsChirpyRed: user.IsChirpyRed.Bool,
		},
		Token: jwtToken,
		RefreshToken: createdRefreshToken.Token,
	}

	respondWithJSON(w, 200, tokenResp)
	return

}