package models

import (
	"database/sql"
	"fmt"
	"portobloglio/utils"
	"time"
)

type Blog struct {
	Id        string
	Title     string
	Content   string
	Category  Category
	Tags      []Tag
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Category struct {
	ID   string
	Name string
}

type Tag struct {
	ID   string
	Name string
}

func GetAllBlogs(db *sql.DB) ([]Blog, error) {
	rows, err := db.Query("SELECT id, title, content, created_at, updated_at FROM blogs")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var blogs []Blog
	for rows.Next() {
		var rawID []byte
		var b Blog

		err := rows.Scan(&rawID, &b.Title, &b.Content, &b.CreatedAt, &b.UpdatedAt)
		if err != nil {
			return nil, err
		}

		b.Id, err = utils.FormatUUID(rawID)
		if err != nil {
			return nil, err
		}

		blogs = append(blogs, b)
	}

	return blogs, nil
}

func GetBlogByID(db *sql.DB, id string) (Blog, error) {
	var b Blog
	idBytes, err := utils.ParseUUID(id)
	if err != nil {
		return Blog{}, fmt.Errorf("invalid uuid: %w", err)
	}

	err = db.QueryRow(
		"SELECT id, title, content, created_at, updated_at FROM blogs WHERE id=?", idBytes).
		Scan(&b.Id, &b.Title, &b.Content, &b.CreatedAt, &b.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return Blog{}, fmt.Errorf("Blog not found! %w", err)
		}
		return Blog{}, fmt.Errorf("query failed! %w", err)
	}
	return b, nil
}
