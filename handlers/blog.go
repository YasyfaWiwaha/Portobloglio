package handlers

import (
	"bytes"
	"net/http"

	"html/template"

	"portobloglio/middlewares"
	"portobloglio/models"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/renderer/html"
)

func (h *Handler) BlogsHandler(w http.ResponseWriter, r *http.Request) {
	blogs, err := models.GetAllBlogs(h.DB)
	if err != nil {
		http.Error(w, "Failed to load blogs", http.StatusInternalServerError)
		return
	}
	var data any
	user := middlewares.GetUserFromContext(r)
	data = newPageData("blogs", map[string]any{
		"Blogs": blogs,
		"User":  user,
	})

	RenderTemplate(w, "blogs", data)
}

func (h *Handler) BlogDetailsHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Redirect(w, r, "/blogs", http.StatusFound)
		return
	}
	blog, err := models.GetBlogByID(h.DB, id)
	if err != nil {
		http.Error(w, "Failed to load blog", http.StatusInternalServerError)
	}

	htmlContent, err := parseMarkDown(blog.Content)
	if err != nil {
		http.Error(w, "Failed to render blog markdown", http.StatusInternalServerError)
		return
	}

	view := struct {
		models.Blog
		HTML template.HTML
	}{
		Blog: blog,
		HTML: htmlContent,
	}

	data := newPageData("blog_details", view)
	RenderTemplate(w, "blog_details", data)
}

func (h *Handler) BlogCreateFormHandler(w http.ResponseWriter, r *http.Request) {
	data := newPageData("create_blog", nil)
	RenderTemplate(w, "create_blog", data)
}

func (h *Handler) BlogCreateHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}
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
