package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
)

var tmpl = template.Must(template.ParseFiles("templates/index.html"))
var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("sqlite3", "./mywebsite.db")
	if err != nil {
		// Handle error - perhaps log it and exit or panic
		panic("failed to open database: " + err.Error())
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	data := struct{ Greeting string }{Greeting: "Hey Boss!!!"}
	err := tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}

func main() {
	defer db.Close()

	http.HandleFunc("/", handler)
	fmt.Println("Server listening on port 8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error starting server", err)
	}
}

// please ffs
