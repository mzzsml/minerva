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

    mux.HandleFunc("/{$}", s.handleIndex)
    mux.HandleFunc("POST /upload", s.handleUpload)
    mux.HandleFunc("GET /hosts/{$}", s.handleGetHosts)
    mux.HandleFunc("GET /hosts/{addr}", s.handleGetHostDetails)

    server := &http.Server{
        Addr:    s.Addr,
        Handler: mux,
    }

    log.Printf("minerva is running on %v", server.Addr)
    err := server.ListenAndServe()
    if err != nil {
        log.Fatal(err)
    }
}

func (s Server) handleIndex(w http.ResponseWriter, r *http.Request) {
    t := template.Must(template.New("upload").ParseFiles("./internal/template/upload.html"))
    err := t.Execute(w, nil)
    if err != nil {
        log.Printf("error: %s\n", err)
        return
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
