package main

import (
	"net/http"
)

type Project struct {
    URL         string
	Title       string
    Description string
}

type ProjectsData struct {
	BaseData
	Projects []Project
}

func projectsHandler(w http.ResponseWriter, r *http.Request) {
	projects, err := getAllProjects()
    if err != nil {
        http.Error(w, "Failed to load projects", http.StatusInternalServerError)
        return
    }
	data := ProjectsData{
		BaseData: BaseData{
			GlobalCSSPath: "/static/css/global.css",
			PageCSSPath:   "/static/css/projects.css",
			HTMXPath:      "/static/scripts/htmx.min.js",
		},
		Projects: projects,
	}
	renderTemplate(w, "base.html", data)
}

func getAllProjects() ([]Project, error) {
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
