package http

import (
    "net/http"
    "log"

    "fmt"

    "github.com/mzzsml/minerva/internal/storage"
)

type Server struct {
    Addr string
    Store *storage.Db
}

func (s *Server) Start() {
    // TODO: in futuro, passare parametri in questa funzione,
    // al posso di hardcodarli (come ad esempio porta, ecc).

    mux := http.NewServeMux()

    mux.HandleFunc("POST /api/scan/{$}", handleNewScan) // faccio una post, con il body che e' il json, fa il parsing e mette i risultati in un db
    //mux.HandleFunc("GET /api/host/{hostAddr}", handleGetHostDetails)
    mux.HandleFunc("GET /api/host/{hostAddr}", s.handleGetHostDetails)

    server := &http.Server{
        Addr:    s.Addr,
        Handler: mux,
    }

    log.Printf("minerva is running.")
    err := server.ListenAndServe()
    if err != nil {
        log.Fatal(err)
    }
}

func (s *Server) handleGetHostDetails(w http.ResponseWriter, r *http.Request) {
    addr := r.PathValue("hostAddr")
 
    res, err := s.Store.GetHostDetails(addr)
    if err != nil {
        log.Fatalf("error: handleGetHostDetails (new version): %s", err)
    }

    fmt.Fprintf(w, "%v\n", res)
}
