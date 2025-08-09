package models

import (
	"database/sql"
)

type Project struct {
	URL         string
	Title       string
	Description string
}

func GetAllProjects(db *sql.DB) ([]Project, error) {
	rows, err := db.Query("SELECT * FROM projects")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []Project
	for rows.Next() {
		var p Project
		err := rows.Scan(&p.URL, &p.Title, &p.Description)
		if err != nil {
			return nil, err
		}
		projects = append(projects, p)
	}
	return projects, nil
}
