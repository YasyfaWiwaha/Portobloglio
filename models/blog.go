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

func CreateBlog(db *sql.DB, title string, content string, cat Category, tags []Tag) (*Blog, error) {
	blobUUID, strUUID, err := utils.GenerateUUIDv7()
	if err != nil {
		return nil, fmt.Errorf("failed to generate UUID: %w", err)
	}

	currentTime := time.Now()
	newBlog := &Blog{
		Id:        strUUID,
		Title:     title,
		Content:   content,
		Category:  cat,
		Tags:      tags,
		CreatedAt: currentTime,
		UpdatedAt: currentTime,
	}

	query := "INSERT INTO blogs (id, title, content, category_id, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)"
	_, err = db.Exec(query, blobUUID, title, content, cat.ID, currentTime, currentTime)

	if err != nil {
		return nil, fmt.Errorf("failed to insert to database")
	}
	return newBlog, err
}

func GetAllCategories(db *sql.DB) ([]Category, error) {
	rows, err := db.Query("SELECT id, name FROM categories")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []Category
	for rows.Next() {
		var c Category
		err := rows.Scan(&c.ID, &c.Name)
		if err != nil {
			return nil, err
		}

		categories = append(categories, c)
	}

	return categories, nil
}

func GetCategoryByID(db *sql.DB, id string) (Category, error) {
	var c Category
	err := db.QueryRow(
		"SELECT id, name FROM categories WHERE id=?", id).
		Scan(c)

	if err != nil {
		if err == sql.ErrNoRows {
			return Category{}, fmt.Errorf("category not found! %w", err)
		}
		return Category{}, fmt.Errorf("query failed! %w", err)
	}
	return c, nil
}
