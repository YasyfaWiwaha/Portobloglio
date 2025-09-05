package models

import "database/sql"

type PageData struct {
	GlobalCSSPath string
	PageCSSPath   string
	HTMXPath      string
	Content       any
}

var db *sql.DB

func SetDB(database *sql.DB) {
	db = database
}
