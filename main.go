package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

var templates = template.Must(template.ParseGlob("templates/*.html"))

var db *sql.DB

type BaseData struct {
	GlobalCSSPath string
	PageCSSPath   string
	HTMXPath      string
}

func init() {
	var err error
	db, err = sql.Open("sqlite3", "db/portobloglio.db")
	if err != nil {
		panic("failed to open database: " + err.Error())
	}
}

func renderTemplate(w http.ResponseWriter, name string, data interface{}) {
	var buf bytes.Buffer
	err := templates.ExecuteTemplate(&buf, name, data)
	if err != nil {
		log.Printf("Template rendering error for %s: %v", name, err)
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
	}
	buf.WriteTo(w)
}

func landingHandler(w http.ResponseWriter, r *http.Request) {
	data := BaseData{
		GlobalCSSPath: "static/css/global.css",
		PageCSSPath:   "static/css/index.css",
		HTMXPath:      "static/scripts/htmx.min.js",
	}
	renderTemplate(w, "index.html", data)
}

func blogsHandler(w http.ResponseWriter, r *http.Request) {
	data := BaseData{
		GlobalCSSPath: "static/css/global.css",
		PageCSSPath:   "static/css/blogs.css",
		HTMXPath:      "static/scripts/htmx.min.js",
	}
	renderTemplate(w, "blogs.html", data)
}


func aboutHandler(w http.ResponseWriter, r *http.Request) {
	data := BaseData{
		GlobalCSSPath: "static/css/global.css",
		PageCSSPath:   "static/css/about.css",
		HTMXPath:      "static/scripts/htmx.min.js",
	}
	renderTemplate(w, "about.html", data)
}

func projectsHandler(w http.ResponseWriter, r *http.Request) {
	data := BaseData{
		GlobalCSSPath: "/static/css/global.css",
		PageCSSPath:   "/static/css/projects.css",
		HTMXPath:      "/static/scripts/htmx.min.js",
	}
	renderTemplate(w, "projects.html", data)
}

func main() {
	defer db.Close()

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", landingHandler)
	http.HandleFunc("/blogs", blogsHandler)
	http.HandleFunc("/projects", projectsHandler)
	http.HandleFunc("/about", aboutHandler)

	fmt.Println("Server listening on port 5050")
	err := http.ListenAndServe(":5050", nil)
	if err != nil {
		fmt.Println("Error starting server", err)
	}
}
