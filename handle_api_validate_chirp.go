package main

import (
    "encoding/json"
    "net/http"
    "strings"
    "fmt"
    "time"

    "github.com/tylerapear/chirpy/internal/database"

    "github.com/google/uuid"
)

func (cfg apiConfig) handle_api_chirps (w http.ResponseWriter, r *http.Request){

    type chirpResponseVals struct {
        ID uuid.UUID `json:"id"`
        CreatedAt time.Time `json:"created_at"`
        UpdatedAt time.Time `json:"updated_at"`
        Body string `json:"body"`
        UserID uuid.UUID `json:"user_id"`
    }

    type chirpParams struct {
        Body string `json:"body"`
        UserID uuid.UUID `json:"user_id"`
    }

    decoder := json.NewDecoder(r.Body)
    params := chirpParams{}
    err := decoder.Decode(&params)
    if err != nil {
        respondWithError(w, 500, fmt.Sprintf("Something went wrong: %s\n", err))
        return
    }
    fmt.Printf("userid: %s\n", params.UserID)

    cleaned_chirp, valid := validate_chirp(params.Body)
    if !valid {
        respondWithError(w, 400, "Chirp message too long")
        return
    }
    
    createChirpParams := database.CreateChirpParams{
        Body: cleaned_chirp,
        UserID: params.UserID,
    }
   
    created_chirp, err := cfg.dbQueries.CreateChirp(r.Context(), createChirpParams)
    if err != nil {
        respondWithError(w, 500, fmt.Sprintf("Error creating chrip: %s", err))
        return
    }

    respBody := chirpResponseVals{
        ID: created_chirp.ID,
        CreatedAt: created_chirp.CreatedAt,
        UpdatedAt: created_chirp.UpdatedAt,
        Body: created_chirp.Body,
        UserID: created_chirp.UserID,
    }
    respondWithJSON(w, 201, respBody)
    return
}
 
func validate_chirp(body string) (string, bool){

    const maxPostLen = 140

    if len(body) <= maxPostLen {
        return cleanProfanity(body), true
    }
    
    return body, false
        
}

func cleanProfanity(post string) string {
    fmt.Println("checking words")
    forbidden_words := []string{"kerfuffle", "sharbert", "fornax"}
    fmt.Printf("words: %v\n", forbidden_words)
    post_words := strings.Split(post, " ")

    for i, word := range post_words {
        fmt.Printf("checking if %s is in list\n", word)
        if contains(forbidden_words, strings.ToLower(word)) {
            post_words[i] = "****"
        } 
    }

    return strings.Join(post_words, " ")

}

func contains (list []string, target string) bool {
    for _, s := range list {
        if s == target {
            return true
        }
    }
    return false
}
