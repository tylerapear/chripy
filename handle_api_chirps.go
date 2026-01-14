package main

import (
    "encoding/json"
    "net/http"
    "strings"
    "fmt"
    "time"

    "github.com/tylerapear/chirpy/internal/database"
    "github.com/tylerapear/chirpy/internal/auth"

    "github.com/google/uuid"
)


type Chirp struct {
    ID uuid.UUID `json:"id"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
    Body string `json:"body"`
    UserID uuid.UUID `json:"user_id"`
}

func (cfg apiConfig) handle_api_chirps_get_by_id (w http.ResponseWriter, r *http.Request) {

    chirp_id_param := r.PathValue("chirpID")
    chirp_id, err := uuid.Parse(chirp_id_param)
    if err != nil {
        respondWithError(w, 400, fmt.Sprintf("Error parsing chirpID: %s\n", err))
        return
    }
    
    resp, err := cfg.dbQueries.GetChirpById(r.Context(), chirp_id)
    if err != nil {
        respondWithError(w, 404, fmt.Sprintf("Cannot find chrip with ID %s", chirp_id))
        return
    }

    chirp := Chirp{
        ID: resp.ID,
        CreatedAt: resp.CreatedAt,
        UpdatedAt: resp.UpdatedAt,
        Body: resp.Body,
        UserID: resp.UserID,
    }

    respondWithJSON(w, 200, chirp)
    return
    
}

func (cfg apiConfig) handle_api_chirps_get (w http.ResponseWriter, r *http.Request) {
    
    resp, err := cfg.dbQueries.GetChirps(r.Context())
    if err != nil {
        respondWithError(w, 500, fmt.Sprintf("Error retrieving chirps: %s\n", err))
        return
    }

    chirps := make([]Chirp, len(resp))
    for i, chrip := range resp{
        chirps[i] = Chirp{
            ID: chrip.ID,
            CreatedAt: chrip.CreatedAt,
            UpdatedAt: chrip.UpdatedAt,
            Body: chrip.Body,
            UserID: chrip.UserID,
        }
    }
    respondWithJSON(w, 200, chirps)
    return
}

func (cfg apiConfig) handle_api_chirps_post (w http.ResponseWriter, r *http.Request){

    type chirpParams struct {
        Body string `json:"body"`
    }

    decoder := json.NewDecoder(r.Body)
    params := chirpParams{}
    err := decoder.Decode(&params)
    if err != nil {
        respondWithError(w, 500, fmt.Sprintf("Something went wrong: %s\n", err))
        return
    }
    
    token, err := auth.GetBearerToken(r.Header)
    if err != nil {
        respondWithError(w, 401, fmt.Sprintf("Error paring Bearer Token: %s", err))
        return
    }

    userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
    if err != nil {
        respondWithError(w, 401, fmt.Sprintf("Error validating token: %s", err))
        return
    }

    cleaned_chirp, valid := validate_chirp(params.Body)
    if !valid {
        respondWithError(w, 400, "Chirp message too long")
        return
    }
    
    createChirpParams := database.CreateChirpParams{
        Body: cleaned_chirp,
        UserID: userID,
    }
   
    created_chirp, err := cfg.dbQueries.CreateChirp(r.Context(), createChirpParams)
    if err != nil {
        respondWithError(w, 500, fmt.Sprintf("Error creating chrip: %s", err))
        return
    }

    respBody := Chirp{
        ID: created_chirp.ID,
        CreatedAt: created_chirp.CreatedAt,
        UpdatedAt: created_chirp.UpdatedAt,
        Body: created_chirp.Body,
        UserID: created_chirp.UserID,
    }
    respondWithJSON(w, 201, respBody)
    return
}

func (cfg apiConfig) handle_api_chirps_delete (w http.ResponseWriter, r *http.Request){
    type response struct {
        Message string `json:"message"`
    }

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

    chirp_id_param := r.PathValue("chirpID")
    chirp_id, err := uuid.Parse(chirp_id_param)
    if err != nil {
        respondWithError(w, 400, fmt.Sprintf("Error parsing chirpID: %s\n", err))
        return
    }

    chirp, err := cfg.dbQueries.GetChirpById(r.Context(), chirp_id)
    if err != nil {
        respondWithError(w, 404, fmt.Sprintf("Cannot find chrip with ID %s", chirp_id))
        return
    }
    
    if chirp.UserID != user_id {
        respondWithError(w, 403, "You are not the owner of this chrip")
        return
    }

    err = cfg.dbQueries.DeleteChirp(r.Context(), chirp.ID)
    if err != nil {
        respondWithError(w, 500, fmt.Sprintf("Error deleting chrip: %s", err))
        return
    }

    respondWithJSON(w, 204, response{
        Message: "Successfully deleted chrip.",
    })
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
    forbidden_words := []string{"kerfuffle", "sharbert", "fornax"}
    post_words := strings.Split(post, " ")

    for i, word := range post_words {
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

