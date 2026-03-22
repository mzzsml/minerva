package http

import (
    "encoding/json"
    "fmt"
    "io"
    "log"
    // "mime"
    // "mime/multipart"
    "net/http"
    // "strings"

    "github.com/mzzsml/nparse"

    "github.com/mzzsml/minerva/internal/storage"
)

type Server struct {
    Addr  string
    Store *storage.DB
}

func (s *Server) Start() {
    mux := http.NewServeMux()

    mux.HandleFunc("POST /hosts/{$}", s.handleInsertHosts)
    mux.HandleFunc("GET /hosts/{$}", s.handleGetHosts)
    mux.HandleFunc("GET /hosts/{addr}", s.handleGetHostDetails)
    mux.HandleFunc("DELETE /hosts/{addr}", s.handleDeleteHost)

    server := &http.Server{
        Addr:    s.Addr,
        Handler: mux,
    }

    log.Printf("minerva is running on %v", server.Addr)
    if err := server.ListenAndServe(); err != nil {
        log.Fatal(err)
    }
}

func (s Server) handleInsertHosts(w http.ResponseWriter, r *http.Request) {
    if r.Method != "" && r.Method != "POST" {
        //TODO: implementare un errore: la richiesta, per poter caricare cose, deve essere POST.
        return
    }
    body, err := io.ReadAll(r.Body)
    if err != nil {
        // TODO: for now just return. Implementare errore.
        return
    }
    nmapScan, err := nparse.NewNmapScan(body)
    if err != nil {
        // TODO: Implementare errore.
        return
    }
    if err := s.Store.InsertHosts(nmapScan); err != nil {
        // TODO: implementare errore.
        return
    }
    fmt.Fprintf(w, "%s\n", body)
}

func (s Server) handleGetHosts(w http.ResponseWriter, r *http.Request) {
    hosts, err := s.Store.GetHosts()
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        fmt.Fprintf(w, "error: %s\n", err)
        return
    }
    m, _ := json.Marshal(hosts)
    w.WriteHeader(http.StatusOK)
    log.Printf("%s %s\n", r.Method, r.URL)
    fmt.Fprintf(w, "%s\n", m)
}

func (s *Server) handleGetHostDetails(w http.ResponseWriter, r *http.Request) {
    host, err := s.Store.GetHostInfo(r.PathValue("addr"))
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        fmt.Fprintf(w, "error: %s\n", err)
        return
    }

    h, _ := json.Marshal(host)
    w.WriteHeader(http.StatusOK)
    log.Printf("%s %s\n", r.Method, r.URL)
    fmt.Fprintf(w, "%s\n", h)
}

func (s *Server) handleDeleteHost(w http.ResponseWriter, r *http.Request) {
    if ok := s.Store.DeleteHost(r.PathValue("addr")); ok {
        fmt.Fprintf(w, `{"status":"host deleted successully"}`)
        log.Printf("%s %s\n", r.Method, r.URL)
    }
    w.WriteHeader(http.StatusNotModified)
    return
}
