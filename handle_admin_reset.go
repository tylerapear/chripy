package main

import (
    "net/http"
    "fmt"
)

func (cfg *apiConfig) handle_admin_reset (w http.ResponseWriter, r *http.Request){

    // Reset Hitcount
    cfg.fileserverHits.Store(0)

    // Reset Users
    err := cfg.dbQueries.ResetUsers(r.Context())
    if err != nil {
        respondWithError(w, 500, fmt.Sprintf("Something went wrong: %s\n", err))
        return
    }

    w.Header().Set("Content-Type", "text/plain; charset=utf-8")
    w.WriteHeader(200)
    data := []byte("Hit count and users reset.\n")
    w.Write(data)
}
