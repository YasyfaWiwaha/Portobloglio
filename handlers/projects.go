package handlers

import (
	"net/http"
	"portobloglio/models"
)

func (h *Handler) ProjectsHandler(w http.ResponseWriter, r *http.Request) {
	projects, err := models.GetAllProjects(h.DB)
	if err != nil {
		h.Logger.Println("error fetching projects", err)
		http.Error(w, "Failed to load projects", http.StatusInternalServerError)
		return
	}
	data := newPageData("projects", projects)
	RenderTemplate(w, "projects", data)
}
