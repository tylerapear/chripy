package main

import (
    "fmt"
    "net/http"
    "log"
    "sync/atomic"

    "github.com/lib/pq"
)

type apiConfig struct {
    fileserverHits atomic.Int32
}

func main() {

    // DEFINITIONS
    apiCfg := apiConfig{
        fileserverHits: atomic.Int32{},
    }
    
    const port = ":8080"
    const filepathRoot = "."

    mux := http.NewServeMux()
    server := http.Server{
        Addr: port,
        Handler: mux,
    }

    // HANDLERS
    mux.Handle(
        "/app/", 
        http.StripPrefix(
            "/app/", 
            apiCfg.middlewareMetricsInc(http.FileServer(http.Dir(filepathRoot))),
        ),
    )
    mux.HandleFunc("GET /api/healthz", handle_api_healthz)
    mux.HandleFunc("GET /admin/metrics", apiCfg.handle_admin_metrics)
    mux.HandleFunc("POST /admin/reset", apiCfg.handle_admin_reset)
    mux.HandleFunc("POST /api/validate_chirp", handle_validate_chirp)

    // START SERVER
    fmt.Printf("Server listening on %s and serving files from %s\n", port, filepathRoot)
    log.Fatal(server.ListenAndServe())

    return
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
   return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
        cfg.fileserverHits.Add(1)
        next.ServeHTTP(w, r)
   }) 
}

