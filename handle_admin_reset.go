package main

import (
    "net/http"
)

func (cfg *apiConfig) handle_admin_reset (w http.ResponseWriter, r *http.Request){ 
    cfg.fileserverHits.Store(0)
    w.Header().Set("Content-Type", "text/plain; charset=utf-8")
    w.WriteHeader(200)
    data := []byte("Hit count reset.\n")
    w.Write(data)
}
