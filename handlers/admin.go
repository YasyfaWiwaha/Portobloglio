package handlers

import (
	"net/http"
)

func (h *Handler) AdminDashboardHandler(w http.ResponseWriter, r *http.Request) {
	data := newPageData("admin", nil)
	RenderTemplate(w, "admin_dashboard", data)
}
