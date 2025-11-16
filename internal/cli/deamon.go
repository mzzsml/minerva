package cli

import (
	//"net/http"

	"github.com/mzzsml/minerva/internal/http"
	"github.com/mzzsml/minerva/internal/storage"
)

func getHostDetailsHandler(storage storage.Db) {
}

func startDeamon(s http.Server) {
	// qui richiamo StartHttpServer
	// attacco un metodo a db, che poi posso usare nell'handler come funzione
	s.Start()
}
