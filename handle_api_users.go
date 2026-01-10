package main

import (
    "fmt"
    "encoding/json"
    "net/http"
    "time"

    "github.com/google/uuid"
)

func (cfg apiConfig) handle_api_users(w http.ResponseWriter, r *http.Request) {

    type userParameters struct {
        Email string `json:"email"`
    }

    type userResponseVals struct {
        ID uuid.UUID `json:"id"`
        Email string `json:"email"`
        CreatedAt time.Time `json:"created_at"`
        UpdatedAt time.Time `json:"updated_at"`
    }

    decoder := json.NewDecoder(r.Body)
    params := userParameters{}
    err := decoder.Decode(&params)
    if err != nil {
        respondWithError(w, 500, fmt.Sprintf("Something went wrong: %s\n", err))
        return
    }

    created_user, err := cfg.dbQueries.CreateUser(r.Context(), params.Email)
    if err != nil {
        respondWithError(w, 500, fmt.Sprintf("Error creating user: %s\n", err))
        return
    }

    respBody := userResponseVals{
        ID: created_user.ID,
        Email: created_user.Email,
        CreatedAt: created_user.CreatedAt,
        UpdatedAt: created_user.UpdatedAt,
    }
    respondWithJSON(w, 201, respBody)
    return

}
