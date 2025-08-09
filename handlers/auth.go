package handlers

import (
	"net/http"
	"portobloglio/models"
	"portobloglio/utils"
	"time"
)

func (h *Handler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		errorMessage := r.URL.Query().Get("error")

		data := newPageData("login", errorMessage)
		RenderTemplate(w, "login.html", data)
	}

	if r.Method == http.MethodPost {

		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Failed to parse form", http.StatusBadRequest)
		}
		username := r.FormValue("username")
		pw := []byte(r.FormValue("password"))

		user, err := models.GetUser(username)
		if err != nil {
			http.Redirect(w, r, "/secret-admin-page-pls-dont-hack-me?error=invalid_credentials", http.StatusSeeOther)
			return
		}

		valid := utils.CheckPassword([]byte(user.Password), pw)

		if !valid {
			http.Redirect(w, r, "/secret-admin-page-pls-dont-hack-me?error=invalid_credentials", http.StatusSeeOther)
			return
		}
		var token string
		token, err = models.GenerateSessionToken(username)
		if err != nil {
			http.Error(w, "error generating token session", http.StatusInternalServerError)
			return
		}
		http.SetCookie(w, &http.Cookie{
			Name:     "Session",
			Value:    token,
			HttpOnly: true,
			Path:     "/",
			SameSite: http.SameSiteStrictMode,
			Expires:  time.Now().Add(8 * time.Hour),
		})
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}
