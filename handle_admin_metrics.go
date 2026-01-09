package main

import (
    "net/http"
    "fmt"
)

func (cfg *apiConfig) handle_admin_metrics (w http.ResponseWriter, r *http.Request) {
 
    templates := loadTemplates()

    w.Header().Set("Content-Type", "text/html; charset=utf-8")
    w.WriteHeader(200)
    data := []byte(
        fmt.Sprintf(
            templates["admin_metrics_template"],
            cfg.fileserverHits.Load(),
        ),
    )
    w.Write(data)
}
