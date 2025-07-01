package main

import (
	"bytes"
	"database/sql"
	"encoding/hex"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
	"time"

	// markdown rendering
	_ "github.com/mattn/go-sqlite3"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/renderer/html"
)

var templates = template.Must(template.ParseGlob("templates/*.html"))

var db *sql.DB

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

type PageData struct {
	GlobalCSSPath string
	PageCSSPath   string
	HTMXPath      string
	Content       any
}

func newPageData(pageCssPath string, content any) PageData {
	pageCSS := ""
	if pageCssPath != "" {
		pageCSS = fmt.Sprintf("/static/css/%s.css", pageCssPath)
	}

	return PageData{
		GlobalCSSPath: "/static/css/global.css",
		PageCSSPath:   pageCSS,
		HTMXPath:      "/static/scripts/htmx.min.js",
		Content:       content,
	}
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

func formatUUID(b []byte) (string, error) {
	if len(b) != 16 {
		return "", fmt.Errorf("invalid UUID length: %d", len(b))
	}
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		uint32(b[0])<<24|uint32(b[1])<<16|uint32(b[2])<<8|uint32(b[3]),
		uint16(b[4])<<8|uint16(b[5]),
		uint16(b[6])<<8|uint16(b[7]),
		uint16(b[8])<<8|uint16(b[9]),
		b[10:]), nil
}

func parseUUID(uuidStr string) ([]byte, error) {
	clean := strings.ReplaceAll(uuidStr, "-", "")
	if len(clean) != 32 {
		return nil, fmt.Errorf("invalid UUID length: %d", len(clean))
	}

	bytes, err := hex.DecodeString(clean)
	if err != nil {
		return nil, fmt.Errorf("failed to decode UUID hex: %w", err)
	}

	return bytes, nil
}

func parseMarkDown(source string) (template.HTML, error) {
	md := goldmark.New(
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
		),
	)
	var buf bytes.Buffer

	err := md.Convert([]byte(source), &buf)
	if err != nil {
		return "", err
	}
	return template.HTML(buf.String()), nil
}

func getAllBlogs() ([]Blog, error) {
	rows, err := db.Query("SELECT id, title, content, created_at, updated_at FROM blogs")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var blogs []Blog
	for rows.Next() {
		var rawID []byte
		var b Blog

		err := rows.Scan(&rawID, &b.Title, &b.Content, &b.CreatedAt, &b.UpdatedAt)
		if err != nil {
			return nil, err
		}

		b.Id, err = formatUUID(rawID)
		if err != nil {
			return nil, err
		}

		blogs = append(blogs, b)
	}

	return blogs, nil
}

func getBlogByID(id string) (Blog, error) {
	var b Blog
	idBytes, err := parseUUID(id)
	if err != nil {
		return Blog{}, fmt.Errorf("invalid uuid: %w", err)
	}

	err = db.QueryRow(
		"SELECT id, title, content, created_at, updated_at FROM blogs WHERE id=?", idBytes).
		Scan(&b.Id, &b.Title, &b.Content, &b.CreatedAt, &b.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return Blog{}, fmt.Errorf("Blog not found! %w", err)
		}
		return Blog{}, fmt.Errorf("query failed! %w", err)
	}
	return b, nil
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
	data := newPageData("index", nil)
	renderTemplate(w, "index.html", data)
}

func blogsHandler(w http.ResponseWriter, r *http.Request) {
	blogs, err := getAllBlogs()
	if err != nil {
		http.Error(w, "Failed to load blogs", http.StatusInternalServerError)
	}
	data := newPageData("blogs", blogs)
	renderTemplate(w, "blogs", data)
}

func blogDetailsHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/blogs/")
	if id == "" {
		http.Redirect(w, r, "/blogs", http.StatusFound)
		return
	}
	blog, err := getBlogByID(id)
	if err != nil {
		http.Error(w, "Failed to load blog", http.StatusInternalServerError)
	}

	htmlContent, err := parseMarkDown(blog.Content)
	if err != nil {
		http.Error(w, "Failed to render blog markdown", http.StatusInternalServerError)
		return
	}

	view := struct {
		Blog
		HTML template.HTML
	}{
		Blog: blog,
		HTML: htmlContent,
	}

	data := newPageData("blog_details", view)
	renderTemplate(w, "blog-details", data)
}

func projectsHandler(w http.ResponseWriter, r *http.Request) {
	projects, err := getAllProjects()
	if err != nil {
		http.Error(w, "Failed to load projects", http.StatusInternalServerError)
		return
	}
	data := newPageData("projects", projects)
	renderTemplate(w, "projects", data)
}

func aboutHandler(w http.ResponseWriter, r *http.Request) {
	data := newPageData("", nil)
	renderTemplate(w, "about.html", data)
}

func main() {
	defer db.Close()

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", landingHandler)
	http.HandleFunc("/blogs", blogsHandler)
	http.HandleFunc("/blogs/", blogDetailsHandler)
	http.HandleFunc("/projects", projectsHandler)
	http.HandleFunc("/about", aboutHandler)

	fmt.Println("Server listening on port 5050")
	err := http.ListenAndServe(":5050", nil)
	if err != nil {
		fmt.Println("Error starting server", err)
	}
}
