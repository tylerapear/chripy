package main

import (
	"fmt"
	"net/http"

	"github.com/tylerapear/chirpy/internal/auth"
)

func (cfg apiConfig) handle_api_revoke (w http.ResponseWriter, r *http.Request) {

	refresh_token, err := auth.GetBearerToken(r.Header)
    if err != nil {
        respondWithError(w, 401, fmt.Sprintf("Error paring Bearer Token: %s", err))
        return
    }

	_, err = cfg.dbQueries.RevokeRefreshToken(r.Context(), refresh_token)
	if err != nil {
        respondWithError(w, 401, fmt.Sprintf("Error revoking refresh token: %s", err))
        return
    }

	respondWithJSON(w, 204, "Successfully revoked refresh token")
	return

}