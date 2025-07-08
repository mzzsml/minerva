package http

import (
    "net/http"
    "log"
)

func StartHttpServer() {
    // TODO: in futuro, passare parametri in questa funzione,
    // al posso di hardcodarli (come ad esempio porta, ecc).

    mux := http.NewServeMux()

    mux.HandleFunc("POST /api/scan/{$}", handleNewScan) // faccio una post, con il body che e' il json, fa il parsing e mette i risultati in un db
    mux.HandleFunc("GET /api/host/{hostAddr}", handleGetHostDetails)

    server := &http.Server{
        Addr:    ":8080",
        Handler: mux,
    }

    log.Printf("minerva is running.")
    err := server.ListenAndServe()
    if err != nil {
        log.Fatal(err)
    }
}
