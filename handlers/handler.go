package handlers

import (
	"database/sql"
	"log"
)

type Handler struct {
	DB     *sql.DB
	Logger *log.Logger
}
