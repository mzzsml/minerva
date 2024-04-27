package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"

	"github.com/minerva/api"
	"github.com/minerva/modules/crtsh"
	"github.com/minerva/types"
)

type User struct {
	Name     string
	Projects []string
}

func newUser() *User {
	u := User{Name: "test"}
	return &u
}

func testTemplateHandler(w http.ResponseWriter, r *http.Request) {
	//var templates *template.Template

	templ_files := []string{
		"ui/template/index.html",
	}

	templates := template.Must(template.ParseFiles(templ_files...))
	templates.Execute(w, "index")
}

func selectProjectHandler(w http.ResponseWriter, r *http.Request) {

	projectsJson := []byte(api.GetProjects())

	var projects []types.Project

	err := json.Unmarshal(projectsJson, &projects)
	if err != nil {
		fmt.Println("error:", err)
	}

	templateFiles := []string{
		"ui/templates/project.html",
	}
	template := template.Must(template.ParseFiles(templateFiles...))

	err = template.Execute(w, projects)
	if err != nil {
		panic(err)
	}
}

func crtshHandler(w http.ResponseWriter, h *http.Request) {
	domain := h.FormValue("domain")
	domains := crtsh.QueryCrtsh(w, h, domain)

	templ_files := []string{
		"ui/templates/crtsh.html",
	}

	template := template.Must(template.ParseFiles(templ_files...))
	if err := template.Execute(w, domains); err != nil {
		panic(err)
	}
	//for _, d := range domains {
	//	fmt.Fprintln(w, d)
	//}
}

func projectHandler(w http.ResponseWriter, h *http.Request) {
	projectName := h.PathValue("projectName")
	fmt.Println(projectName)
}

func main() {
	mux := http.NewServeMux()
	files := http.FileServer(http.Dir("assets"))
	mux.Handle("/static/", http.StripPrefix("/static/", files))

	//mux.HandleFunc("/", testTemplateHandler)
	mux.HandleFunc("/{$}", selectProjectHandler)
	mux.HandleFunc("/crtsh", crtshHandler)
	mux.HandleFunc("/project/{projectName}/", projectHandler)

	server := &http.Server{
		Addr:    "0.0.0.0:9999",
		Handler: mux,
	}
	err := server.ListenAndServe()

	// err = godotenv.Load()
	// if err != nil {
	// 	log.Fatal("Error loading .env file")
	// }

	if err != nil {
		fmt.Println("error:", err)
	}
}
