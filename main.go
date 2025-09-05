package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"portobloglio/handlers"
	"portobloglio/middlewares"
	"portobloglio/models"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("sqlite3", "db/portobloglio.db")
	if err != nil {
		panic("failed to open database: " + err.Error())
	}
	models.SetDB(db)
}

func sessionCleanupLoop(db *sql.DB, interval time.Duration) {
	ticker := time.NewTicker(interval)

	cleanExpiredSessions(db)
	defer ticker.Stop()
	for range ticker.C {
		fmt.Println("loop")
		cleanExpiredSessions(db)
	}
}

func cleanExpiredSessions(db *sql.DB) {
	_, err := db.Exec(`DELETE FROM sessions WHERE expired_at < CURRENT_TIMESTAMP`)
	if err != nil {
		log.Printf("Session cleanup error: %v", err)
	}
	fmt.Println("cleaned expiread sessions")
}

func main() {
	defer db.Close()
	go sessionCleanupLoop(db, time.Hour)

	h := &handlers.Handler{
		DB:     db,
		Logger: log.Default(),
	}

	mux := http.NewServeMux()
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	mux.HandleFunc("/", h.LandingHandler)
	mux.HandleFunc("/projects", h.ProjectsHandler)
	mux.HandleFunc("/blogs", middlewares.OptionalAuthMiddleware(h.BlogsHandler))
	mux.HandleFunc("/blogs/", h.BlogDetailsHandler)
	mux.HandleFunc("/about", h.AboutHandler)
	mux.HandleFunc("/secret-admin-page-pls-dont-hack-me", h.LoginHandler)
	mux.HandleFunc("/dashboard", middlewares.AuthMiddleware(h.AdminDashboardHandler))

	fmt.Println("Server listening on port 5050")
	err := http.ListenAndServe(":5050", mux)
	if err != nil {
		fmt.Println("Error starting server", err)
	}
}
