package handlers

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"portobloglio/models"
)

var templates = template.Must(template.ParseGlob("templates/*.html"))

func RenderTemplate(w http.ResponseWriter, name string, data interface{}) {
	var buf bytes.Buffer
	err := templates.ExecuteTemplate(&buf, name, data)
	if err != nil {
		log.Printf("Template rendering error for %s: %v", name, err)
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
	}
	buf.WriteTo(w)
}

func newPageData(pageCssPath string, content any) models.PageData {
	pageCSS := ""
	if pageCssPath != "" {
		pageCSS = fmt.Sprintf("/static/css/%s.css", pageCssPath)
	}

	return models.PageData{
		GlobalCSSPath: "/static/css/global.css",
		PageCSSPath:   pageCSS,
		HTMXPath:      "/static/scripts/htmx.min.js",
		Content:       content,
	}
}

func (h *Handler) LandingHandler(w http.ResponseWriter, r *http.Request) {
	data := newPageData("index", nil)
	RenderTemplate(w, "index.html", data)
}

func (h *Handler) AboutHandler(w http.ResponseWriter, r *http.Request) {
	data := newPageData("about", nil)
	RenderTemplate(w, "about.html", data)
}
