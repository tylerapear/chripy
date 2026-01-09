package main

import (
    "net/http"
)

func handle_api_healthz (w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/plain; charset=utf-8")
    w.WriteHeader(200)
    data := []byte("OK\n")
    w.Write(data)
}
