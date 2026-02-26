package models

import (
	"database/sql"
	"fmt"
	"portobloglio/utils"
	"time"
)

// ==================== BLOGS ====================

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
	var rawID []byte

	idBytes, err := utils.ParseUUID(id)
	if err != nil {
		return Blog{}, fmt.Errorf("invalid uuid: %w", err)
	}

	err = db.QueryRow(
		"SELECT id, title, content, created_at, updated_at FROM blogs WHERE id=?", idBytes).
		Scan(&rawID, &b.Title, &b.Content, &b.CreatedAt, &b.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return Blog{}, fmt.Errorf("blog not found: %w", err)
		}
		return Blog{}, fmt.Errorf("query failed: %w", err)
	}

	b.Id, err = utils.FormatUUID(rawID)
	if err != nil {
		return Blog{}, fmt.Errorf("failed to format uuid: %w", err)
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
		return nil, fmt.Errorf("failed to insert blog: %w", err)
	}

	return newBlog, nil
}


func UpdateBlog(db *sql.DB, id string, title string, content string, cat Category, tags []Tag) (*Blog, error) {
	idBytes, err := utils.ParseUUID(id)
	if err != nil {
		return nil, fmt.Errorf("invalid uuid: %w", err)
	}

	updatedAt := time.Now()

	query := "UPDATE blogs SET title=?, content=?, category_id=?, updated_at=? WHERE id=?"
	result, err := db.Exec(query, title, content, cat.ID, updatedAt, idBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to update blog: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("failed to check rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return nil, fmt.Errorf("blog not found")
	}

	// Replace all existing tags for this blog
	err = ReplaceBlogTags(db, idBytes, tags)
	if err != nil {
		return nil, fmt.Errorf("failed to update blog tags: %w", err)
	}

	updatedBlog := &Blog{
		Id:       id,
		Title:    title,
		Content:  content,
		Category: cat,
		Tags:     tags,
		UpdatedAt: updatedAt,
	}

	return updatedBlog, nil
}

func DeleteBlogByID(db *sql.DB, id string) error {
	idBytes, err := utils.ParseUUID(id)
	if err != nil {
		return fmt.Errorf("invalid uuid: %w", err)
	}

	result, err := db.Exec("DELETE FROM blogs WHERE id=?", idBytes)
	if err != nil {
		return fmt.Errorf("query failed: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("blog not found")
	}

	return nil
}

// ==================== CATEGORIES ====================

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
	err := db.QueryRow("SELECT id, name FROM categories WHERE id=?", id).
		Scan(&c.ID, &c.Name) // Fixed: was Scan(c)
	if err != nil {
		if err == sql.ErrNoRows {
			return Category{}, fmt.Errorf("category not found: %w", err)
		}
		return Category{}, fmt.Errorf("query failed: %w", err)
	}
	return c, nil
}

func CreateCategory(db *sql.DB, id string, name string) (*Category, error) {
	newCategory := &Category{
		ID:   id,
		Name: name,
	}

	_, err := db.Exec("INSERT INTO categories (id, name) VALUES (?, ?)", id, name)
	if err != nil {
		return nil, fmt.Errorf("failed to insert category: %w", err)
	}

	return newCategory, nil
}

func UpdateCategory(db *sql.DB, id string, name string) (*Category, error) {
	result, err := db.Exec("UPDATE categories SET name=? WHERE id=?", name, id)
	if err != nil {
		return nil, fmt.Errorf("failed to update category: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("failed to check rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return nil, fmt.Errorf("category not found")
	}

	return &Category{ID: id, Name: name}, nil
}

func DeleteCategoryByID(db *sql.DB, id string) error {
	result, err := db.Exec("DELETE FROM categories WHERE id=?", id)
	if err != nil {
		return fmt.Errorf("query failed: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("category not found")
	}

	return nil
}

// ==================== TAGS ====================

func GetAllTags(db *sql.DB) ([]Tag, error) {
	rows, err := db.Query("SELECT id, name FROM tags")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []Tag
	for rows.Next() {
		var t Tag
		err := rows.Scan(&t.ID, &t.Name)
		if err != nil {
			return nil, err
		}
		tags = append(tags, t)
	}

	return tags, nil
}

func GetTagByID(db *sql.DB, id string) (Tag, error) {
	var t Tag
	err := db.QueryRow("SELECT id, name FROM tags WHERE id=?", id).
		Scan(&t.ID, &t.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			return Tag{}, fmt.Errorf("tag not found: %w", err)
		}
		return Tag{}, fmt.Errorf("query failed: %w", err)
	}
	return t, nil
}

func CreateTag(db *sql.DB, id string, name string) (*Tag, error) {
	newTag := &Tag{
		ID:   id,
		Name: name,
	}

	_, err := db.Exec("INSERT INTO tags (id, name) VALUES (?, ?)", id, name) // Fixed: was title, was undeclared err
	if err != nil {
		return nil, fmt.Errorf("failed to insert tag: %w", err)
	}

	return newTag, nil
}

func UpdateTag(db *sql.DB, id string, name string) (*Tag, error) {
	result, err := db.Exec("UPDATE tags SET name=? WHERE id=?", name, id)
	if err != nil {
		return nil, fmt.Errorf("failed to update tag: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("failed to check rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return nil, fmt.Errorf("tag not found")
	}

	return &Tag{ID: id, Name: name}, nil
}

func DeleteTagByID(db *sql.DB, id string) error {
	result, err := db.Exec("DELETE FROM tags WHERE id=?", id)
	if err != nil {
		return fmt.Errorf("query failed: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("tag not found")
	}

	return nil
}

// ==================== BLOG TAGS (join table) ====================

// ReplaceBlogTags deletes all existing tag associations for a blog and inserts the new ones.
// Assumes a join table: blog_tags (blog_id, tag_id)
func ReplaceBlogTags(db *sql.DB, blogIDBytes []byte, tags []Tag) error {
	_, err := db.Exec("DELETE FROM blog_tags WHERE blog_id=?", blogIDBytes)
	if err != nil {
		return fmt.Errorf("failed to delete old blog tags: %w", err)
	}

	for _, tag := range tags {
		tagIDBytes, err := utils.ParseUUID(tag.ID)
		if err != nil {
			return fmt.Errorf("invalid tag uuid %s: %w", tag.ID, err)
		}
		_, err = db.Exec("INSERT INTO blog_tags (blog_id, tag_id) VALUES (?, ?)", blogIDBytes, tagIDBytes)
		if err != nil {
			return fmt.Errorf("failed to insert blog tag %s: %w", tag.ID, err)
		}
	}

	return nil
}