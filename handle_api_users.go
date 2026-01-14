package main

import (
    "fmt"
    "encoding/json"
    "net/http"
    "time"
    "database/sql"

    "github.com/tylerapear/chirpy/internal/database"
    "github.com/tylerapear/chirpy/internal/auth"

    "github.com/google/uuid"
)

 type User struct {
    ID uuid.UUID `json:"id"`
    Email string `json:"email"`
    Password string `json:"-"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
    IsChirpyRed bool `json:"is_chirpy_red"`
}

type userParameters struct {
        Password string `json:password`
        Email string `json:"email"`
}

func (cfg apiConfig) handle_api_users(w http.ResponseWriter, r *http.Request) {

    decoder := json.NewDecoder(r.Body)
    params := userParameters{}
    err := decoder.Decode(&params)
    if err != nil {
        respondWithError(w, 500, fmt.Sprintf("Something went wrong: %s\n", err))
        return
    }

    passwordHash, err := auth.HashPassword(params.Password)
    if err != nil {
        respondWithError(w, 500, fmt.Sprintf("Error hashing password: %s", err ))
        return
    }

    createUserParams := database.CreateUserParams{
        HashedPassword: sql.NullString{String: passwordHash, Valid: true},
        Email: params.Email,
    }

    created_user, err := cfg.dbQueries.CreateUser(r.Context(), createUserParams)
    if err != nil {
        respondWithError(w, 500, fmt.Sprintf("Error creating user: %s\n", err))
        return
    }

    respBody := User{
        ID: created_user.ID,
        Email: created_user.Email,
        CreatedAt: created_user.CreatedAt,
        UpdatedAt: created_user.UpdatedAt,
        IsChirpyRed: created_user.IsChirpyRed.Bool,
    }
    respondWithJSON(w, 201, respBody)
    return

}

func (cfg apiConfig) handle_api_users_put(w http.ResponseWriter, r *http.Request) {

    token, err := auth.GetBearerToken(r.Header)
    if err != nil {
        respondWithError(w, 401, fmt.Sprintf("Error paring Bearer Token: %s", err))
        return
    }

    user_id, err := auth.ValidateJWT(token, cfg.jwtSecret)
    if err != nil {
        respondWithError(w, 401, fmt.Sprintf("Invalid Token"))
        return
    }

    decoder := json.NewDecoder(r.Body)
    params := userParameters{}
    err = decoder.Decode(&params)
    if err != nil {
        respondWithError(w, 500, fmt.Sprintf("Something went wrong: %s\n", err))
        return
    }

    passwordHash, err := auth.HashPassword(params.Password)
    if err != nil {
        respondWithError(w, 500, fmt.Sprintf("Error hashing password: %s", err ))
        return
    }

    updateUserParams := database.UpdateUserParams{
        ID: user_id,
        HashedPassword: sql.NullString{String: passwordHash, Valid: true},
        Email: params.Email,
    }

    updated_user, err := cfg.dbQueries.UpdateUser(r.Context(), updateUserParams)
    if err != nil {
        respondWithError(w, 500, fmt.Sprintf("Error updating user: %s\n", err))
        return
    }

    respBody := User{
        ID: updated_user.ID,
        Email: updated_user.Email,
        CreatedAt: updated_user.CreatedAt,
        UpdatedAt: updated_user.UpdatedAt,
        IsChirpyRed: updated_user.IsChirpyRed.Bool,
    }
    respondWithJSON(w, 200, respBody)
    return

}
