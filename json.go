package main

import (
    "encoding/json"
    "log"
    "net/http"
)

func respondWithError(w http.ResponseWriter, code int, msg string) {

    type errorVals struct {
        Error string `json:"error"`
    }

    errorBody := errorVals{
        Error: msg,
    }
    dat, err := json.Marshal(errorBody)
    if err != nil {
        log.Printf("Error marshalling JSON: %s\n", err)
        w.WriteHeader(500)
        return
    }
    w.WriteHeader(code)
    w.Write(dat)
    return
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
    
    dat, err := json.Marshal(payload)
    if err != nil {
        log.Printf("Error marshalling JSON: %s\n", err)
        w.WriteHeader(500)
        return
    }
    w.WriteHeader(code)
    w.Write(dat)
    return
}
