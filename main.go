package main

import (
	"html/template"
	"net/http"

	"github.com/heimdallr/module/crtsh"
)

func testTemplateHandler(w http.ResponseWriter, r *http.Request) {
	//var templates *template.Template

	templ_files := []string{
		"ui/template/index.html",
	}

	templates := template.Must(template.ParseFiles(templ_files...))
	templates.Execute(w, "index")
}

func crtshHandler(w http.ResponseWriter, h *http.Request) {
	domain := h.FormValue("domain")
	crtsh.QueryCrtsh(w, h, domain)
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", testTemplateHandler)
	mux.HandleFunc("/crtsh", crtshHandler)

	server := &http.Server{
		Addr:    "0.0.0.0:9999",
		Handler: mux,
	}
	server.ListenAndServe()
}
