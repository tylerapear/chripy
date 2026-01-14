package main

import (
	"fmt"
	"net/http"
	"encoding/json"
	"database/sql"

	"github.com/tylerapear/chirpy/internal/database"
	"github.com/tylerapear/chirpy/internal/auth"

	"github.com/google/uuid"
)

func (cfg apiConfig) handle_api_polka_webhooks_post(w http.ResponseWriter, r *http.Request) {

	type polkaRequest struct {
		Event string `json:"event"`
		Data struct{
			UserID string `json:"user_id"`
		} `json:"data"`
	}

	type response struct{
		Message string `json:"message"`
	}

	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil || apiKey != cfg.polkaAPIKey {
		respondWithError(w, 401, fmt.Sprintf("Invalid API Key: %s\n", err))
        return
	}


	decoder := json.NewDecoder(r.Body)
    params := polkaRequest{}
    err = decoder.Decode(&params)
    if err != nil {
        respondWithError(w, 500, fmt.Sprintf("Something went wrong: %s\n", err))
        return
    }

	if params.Event != "user.upgraded" {
		respondWithJSON(w, 204, response{})
		return
	}

	user_id, err := uuid.Parse(params.Data.UserID)
	if err != nil {
		respondWithError(w, 500, fmt.Sprintf("Error parsing user id: %s\n", err))
		return
	}

	_, err = cfg.dbQueries.UpdateUserIsChirpyRed(r.Context(), database.UpdateUserIsChirpyRedParams{
		ID: user_id,
		IsChirpyRed: sql.NullBool{Bool: true, Valid: true},
	})
	if err != nil {
        respondWithError(w, 404, fmt.Sprintf("Error updating user: %s\n", err))
        return
    }

	respondWithJSON(w, 204, response{})
	return

}