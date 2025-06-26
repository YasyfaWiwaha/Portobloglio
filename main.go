package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var templates = template.Must(template.ParseGlob("templates/*.html"))

var db *sql.DB

type BaseData struct {
	GlobalCSSPath string
	PageCSSPath   string
	HTMXPath      string
}

type Project struct {
	URL         string
	Title       string
	Description string
}

type Blog struct {
	Id        string
	Title     string
	Content   string
	Category  Category
	Tags      []Tag
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Category struct {
	ID   string
	Name string
}

type Tag struct {
	ID   string
	Name string
}

type ProjectsData struct {
	BaseData
	Projects []Project
}

type BlogsData struct {
	BaseData
	Blogs []Blog
}

func init() {
	var err error
	db, err = sql.Open("sqlite3", "db/portobloglio.db")
	if err != nil {
		panic("failed to open database: " + err.Error())
	}
}

func getAllProjects() ([]Project, error) {
	rows, err := db.Query("SELECT * FROM projects")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []Project
	for rows.Next() {
		var p Project
		err := rows.Scan(&p.URL, &p.Title, &p.Description)
		if err != nil {
			return nil, err
		}
		projects = append(projects, p)
	}
	return projects, nil
}

func getAllBlogs() ([]Blog, error) {
	rows, err := db.Query("SELECT id, title, content, created_at, updated_at FROM blogs")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var blogs []Blog
	for rows.Next() {
		var b Blog
		err := rows.Scan(&b.Id, &b.Title, &b.Content, &b.CreatedAt, &b.UpdatedAt)
		if err != nil {
			return nil, err
		}
		b.Id = fmt.Sprintf("%x", b.Id)
		blogs = append(blogs, b)
	}
	return blogs, nil
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
	blogs, err := getAllBlogs()
	if err != nil {
		http.Error(w, "Failed to load blogs", http.StatusInternalServerError)
	}
	data := BlogsData{
		BaseData: BaseData{
			GlobalCSSPath: "/static/css/global.css",
			PageCSSPath:   "/static/css/projects.css",
			HTMXPath:      "/static/scripts/htmx.min.js",
		},
		Blogs: blogs,
	}
	renderTemplate(w, "blogs", data)
}

func projectsHandler(w http.ResponseWriter, r *http.Request) {
	projects, err := getAllProjects()
	if err != nil {
		http.Error(w, "Failed to load projects", http.StatusInternalServerError)
		return
	}
	data := ProjectsData{
		BaseData: BaseData{
			GlobalCSSPath: "/static/css/global.css",
			PageCSSPath:   "/static/css/projects.css",
			HTMXPath:      "/static/scripts/htmx.min.js",
		},
		Projects: projects,
	}
	renderTemplate(w, "projects", data)
}

func aboutHandler(w http.ResponseWriter, r *http.Request) {
	data := BaseData{
		GlobalCSSPath: "static/css/global.css",
		PageCSSPath:   "static/css/about.css",
		HTMXPath:      "static/scripts/htmx.min.js",
	}
	renderTemplate(w, "about.html", data)
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
