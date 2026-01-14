package main

import _ "github.com/lib/pq"

import (
    "fmt"
    "net/http"
    "log"
    "sync/atomic"
    "os"
    "database/sql"

    "github.com/joho/godotenv"
    
    "github.com/tylerapear/chirpy/internal/database"
)

type apiConfig struct {
    fileserverHits atomic.Int32
    dbQueries *database.Queries
    jwtSecret string
    polkaAPIKey string
}

func main() {

    // DEFINITIONS
    godotenv.Load()
    dbURL := os.Getenv("DB_URL")

    const port = ":8080"
    const filepathRoot = "."

    mux := http.NewServeMux()
    server := http.Server{
        Addr: port,
        Handler: mux,
    }

    db, err := sql.Open("postgres", dbURL)
    if err != nil {
        fmt.Printf("Error opening database connection: %s\n", err)
        os.Exit(1)
    }

    apiCfg := apiConfig{
        fileserverHits: atomic.Int32{},
        dbQueries: database.New(db),
        jwtSecret: os.Getenv("JWT_SECRET"),
        polkaAPIKey: os.Getenv("POLKA_KEY"),
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
    mux.HandleFunc("POST /api/chirps", apiCfg.handle_api_chirps_post)
    mux.HandleFunc("GET /api/chirps", apiCfg.handle_api_chirps_get)
    mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.handle_api_chirps_get_by_id)
    mux.HandleFunc("DELETE /api/chirps/{chirpID}", apiCfg.handle_api_chirps_delete)
    mux.HandleFunc("POST /api/users", apiCfg.handle_api_users)
    mux.HandleFunc("PUT /api/users", apiCfg.handle_api_users_put)
    mux.HandleFunc("POST /api/login", apiCfg.handle_api_login)
    mux.HandleFunc("POST /api/refresh", apiCfg.handle_api_refresh)
    mux.HandleFunc("POST /api/revoke", apiCfg.handle_api_revoke)
    mux.HandleFunc("POST /api/polka/webhooks", apiCfg.handle_api_polka_webhooks_post)

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

