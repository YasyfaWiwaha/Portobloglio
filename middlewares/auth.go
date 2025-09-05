package middlewares

import (
	"context"
	"net/http"
	"portobloglio/models"
)

type ContextKey string

const (
	UserContextKey ContextKey = "user"
)

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionCookies, err := r.Cookie("Session")
		if err != nil {
			http.Redirect(w, r, "/secret-admin-page-pls-dont-hack-me", http.StatusSeeOther)
			return
		}

		user, err := models.ValidateSession(sessionCookies.Value)
		if err != nil {
			// Invalid or expired session
			// Clear the invalid cookie
			http.SetCookie(w, &http.Cookie{
				Name:     "Session",
				Value:    "",
				HttpOnly: true,
				Path:     "/",
				MaxAge:   -1, // Delete the cookie
			})
			http.Redirect(w, r, "/secret-admin-page-pls-dont-hack-me", http.StatusSeeOther)
			return
		}

		ctx := context.WithValue(r.Context(), UserContextKey, user)
		r = r.WithContext(ctx)
		next(w, r)
	}
}

func OptionalAuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionCookies, err := r.Cookie("Session")
		// if session cookies exist and valid then add user to context, otherwise continue
		if err == nil {
			user, err := models.ValidateSession(sessionCookies.Value)
			if err == nil {
				ctx := context.WithValue(r.Context(), UserContextKey, user)
				r = r.WithContext(ctx)
			}
		}
		next(w, r)
	}
}

func GetUserFromContext(r *http.Request) *models.User {
	user, ok := r.Context().Value(UserContextKey).(*models.User)
	if ok {
		return user
	}
	return nil
}

func IsAuthenticated(r *http.Request) bool {
	return GetUserFromContext(r) != nil
}
