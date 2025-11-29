package http

import (
    "log"
    "net/http"
    "fmt"
    "html/template"
    "mime"
    "io"
    "strings"
    "mime/multipart"

    "github.com/mzzsml/minerva/internal/storage"
)

type Server struct {
    Addr  string
    Store *storage.Db
}

func (s *Server) Start() {
    mux := http.NewServeMux()

    mux.HandleFunc("/{$}", s.handleIndex)
    mux.HandleFunc("POST /upload", s.handleUpload)
    mux.HandleFunc("POST /api/scan/{$}", handleNewScan) // faccio una post, con il body che e' il json, fa il parsing e mette i risultati in un db
    mux.HandleFunc("GET /api/host/{hostAddr}", s.handleGetHostDetails)

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
                log.Printf("error: %s\n", err)
                return
            }
            slurp, err := io.ReadAll(p)
            if err != nil {
                log.Printf("error: %s\n", err)
            }
            fmt.Fprintf(w, "%s\n", slurp)
        }
    }

    //r.ParseMultipartForm(10 << 20)

    //file, _, err := r.FormFile("file")
    //if err != nil {
    //    http.Error(w, "Error retrieving the file", http.StatusBadRequest)
    //    log.Printf("%s\n", err)
    //    return
    //}
    //defer file.Close()
    //var f *os.File
    //f.ReadFrom(file)
    //fmt.Fprintf(w, "%v\n", f)
}

func (s *Server) handleGetHostDetails(w http.ResponseWriter, r *http.Request) {
    addr := r.PathValue("hostAddr")

    res, err := s.Store.GetHostDetails(addr)
    if err != nil {
        log.Fatalf("error: handleGetHostDetails (new version): %s", err)
    }

    fmt.Fprintf(w, "%v\n", res)
}
