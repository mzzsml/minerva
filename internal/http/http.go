package http

import (
    "encoding/json"
    "fmt"
    "html/template"
    "io"
    "log"
    "mime"
    "mime/multipart"
    "net/http"
    "strings"

    "github.com/mzzsml/nparse"

    "github.com/mzzsml/minerva/internal/storage"
)

type Server struct {
    Addr  string
    Store *storage.DB
}

func (s *Server) Start() {
    mux := http.NewServeMux()

    mux.HandleFunc("POST /upload", s.handleUpload)
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

func (s Server) handleUpload(w http.ResponseWriter, r *http.Request) {
    mediaType, params, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
    if err != nil {
        log.Printf("error: %s\n", err)
        return
    }
    if strings.HasPrefix(mediaType, "multipart/") {
        mr := multipart.NewReader(r.Body, params["boundary"])
        for {
            p, err := mr.NextPart()
            if err == io.EOF {
                return
            }
            if err != nil {
                w.WriteHeader(http.StatusInternalServerError)
                log.Printf("error: %s\n", err)
                return
            }
            slurp, err := io.ReadAll(p)
            if err != nil {
                w.WriteHeader(http.StatusInternalServerError)
                log.Printf("error: %s\n", err)
                return
            }
            nmapScan, _ := nparse.NewNmapScan(slurp)
            err = s.Store.InsertHosts(nmapScan)
            if err != nil {
                w.WriteHeader(http.StatusInternalServerError)
                fmt.Fprintf(w, "Error in inserting new host in db: %s\n", err)
                return
            } else {
                // AS OF TODAY 03/12/25 IT RETURNS NO ERROR IF HOST ALREADY
                // EXISTS, SO IT WILL PRINT UPLOADED SUCC... EVEN IF IT DIDNT
                // ACTUALLY UPLOAD IT.
                w.WriteHeader(http.StatusOK)
                fmt.Fprintf(w, "Uploaded successfully.\n")
            }
        }
    }
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
    host, err := s.Store.GetHostInfoByIPv4(r.PathValue("addr"))
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
