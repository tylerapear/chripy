package main

import (
	"fmt"
	"net/http"
	"time"
	"errors"
	"database/sql"

	"github.com/tylerapear/chirpy/internal/auth"
)

type tokenRespParams struct {
	Token string `json:"token"`
}

func (cfg apiConfig) handle_api_refresh (w http.ResponseWriter, r *http.Request) {

	refresh_token, err := auth.GetBearerToken(r.Header)
    if err != nil {
        respondWithError(w, 401, fmt.Sprintf("Error paring Bearer Token: %s", err))
        return
    }

	dbTokenLookup, err := cfg.dbQueries.GetRefreshToken(r.Context(), refresh_token)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, 401, "Refresh token is invalid")
			return
		}
		respondWithError(w, 500, fmt.Sprintf("Error retrieving refresh token: %s", err))
        return
	}
	if dbTokenLookup.ExpiresAt.Before(time.Now()) || dbTokenLookup.RevokedAt.Valid {
		respondWithError(w, 401, fmt.Sprintf("Refresh token expired/revoked"))
        return
	}

	user, err := cfg.dbQueries.GetUserFromRefreshToken(r.Context(), refresh_token)
	if err != nil {
		respondWithError(w, 500, fmt.Sprintf("Error looking up user: %s", err))
        return
	}

	access_token, err := auth.MakeJWT(user.ID, cfg.jwtSecret, time.Hour)
	if err != nil {
		respondWithError(w, 500, fmt.Sprintf("Error generating JWT: %s", err))
        return
	}

	tokenResp := tokenRespParams{
		Token: access_token,
	}
	respondWithJSON(w, 200, tokenResp)
	return

}