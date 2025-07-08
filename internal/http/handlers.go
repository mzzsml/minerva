package http

import (
    "net/http"
    "log"
    "fmt"

    "github.com/mzzsml/minerva/internal/parser"
    "github.com/mzzsml/minerva/internal/storage"
)

func handleNewScan(w http.ResponseWriter, r *http.Request) {
    // nota implementazione:
    // ottieni il body
    // fai il parsing del json
    //poi verra' inserito in db
    //ritorna uuid
    // cosi' poi faccio GET /scan/<uuid> e vedo i risultati
    body := make([]byte, r.ContentLength)
    r.Body.Read(body)
    //fmt.Fprintf(w, "%s\n", body)

    //faccio il parsing
    var h parser.Host
    err := h.ParseFromJson(body)
    if err != nil {
        log.Fatal(err)
    }
    log.Printf("DEBUG: got %v\n", h)

    //storage.GetHostById(1)
    storage.NewHost(h) 
}

func handleGetHostDetails(w http.ResponseWriter, r *http.Request) {
    // prendo l'id dell'host dall'url (/hosts/1)
    // qui lo simulo hardcodando
    //hostAddr := "192.168.1.48"
    
    //var p parser.Port
    p, err := storage.GetHostDetails(r.PathValue("hostAddr"))
    if err != nil {
        log.Printf("%s", err)
    }
    fmt.Fprintf(w, "%v\n", p)
}
