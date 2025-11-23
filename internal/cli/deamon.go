package cli

import (
    "github.com/mzzsml/minerva/internal/http"
)

func startDeamon(s http.Server) {
    s.Start()
}
