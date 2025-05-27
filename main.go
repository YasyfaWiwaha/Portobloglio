package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	_ "github.com/mattn/go-sqlite3"
)

var tmpl = template.Must(template.ParseFiles("templates/index.html"))
var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("sqlite3", "db/portobloglio.db")
	if err != nil {
		panic("failed to open database: " + err.Error())
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	data := struct{ 
		Greeting string 
		CSSPath string
		HTMXPath string
		}{
			Greeting: "Hey Boss!!!",
			CSSPath: "static/css/styles.css",
			HTMXPath: "static/scripts/htmx.min.js",
		}
	err := tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}

func boomHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, `<div style="padding: 1rem; background: #ffc; border: 2px solid #fa0; margin-top: 1rem;">💥 BOOM!</div>`)
}


func main() {
	defer db.Close()
	
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", handler)
	http.HandleFunc("/boom", boomHandler)

	fmt.Println("Server listening on port 5050")
	err := http.ListenAndServe(":5050", nil)
	if err != nil {
		fmt.Println("Error starting server", err)
	}
}
