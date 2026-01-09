package main

import (
    "encoding/json"
    "net/http"
    "strings"
    "fmt"
)
 
func handle_validate_chirp(w http.ResponseWriter, r *http.Request){

    type parameters struct {
        Body string `json:"body"`
    }
 
    type respVals struct {
        Valid bool `json:"valid"`
        CleanedBody string `json:"cleaned_body"`
    }

    const maxPostLen = 140

    decoder := json.NewDecoder(r.Body)
    params := parameters{}
    err := decoder.Decode(&params)
    if err != nil {
        respondWithError(w, 500, fmt.Sprintf("Something went wrong: %s", err))
        return
    }


    if len(params.Body) <= maxPostLen {
        cleanedPost := cleanProfanity(params.Body)
        respBody := respVals{
            Valid: true,
            CleanedBody: cleanedPost,
        }
        respondWithJSON(w, 200, respBody)
        return
    }

    respondWithError(w, 400, "Chirp message is too long")
    return
        
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
