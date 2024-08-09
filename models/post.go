package models

import (
	"database/sql"
	"time"
)

type Post struct {
	ID         int
	Title      string
	Content    string
	UserID     int
	CategoryID int
	Timestamp  time.Time
	Author     *User
	Category   *Category
}

// Modify CreatePost to accept *sql.DB
func CreatePost(db *sql.DB, post *Post) error {
	query := `INSERT INTO posts (title, content, user_id, category_id) VALUES (?, ?, ?, ?)`
	result, err := db.Exec(query, post.Title, post.Content, post.UserID, post.CategoryID)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	post.ID = int(id)
	return nil
}

func GetPostByID(db *sql.DB, id int) (*Post, error) {
	post := &Post{}
	query := `SELECT post_id, title, content, user_id, category_id, timestamp FROM posts WHERE post_id = ?`
	err := db.QueryRow(query, id).Scan(&post.ID, &post.Title, &post.Content, &post.UserID, &post.CategoryID, &post.Timestamp)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	post.Author, err = GetUserByID(db, post.UserID) // Pass db to GetUserByID
	if err != nil {
		return nil, err
	}

	post.Category, err = GetCategoryByID(post.CategoryID) // No db passed here
	if err != nil {
		return nil, err
	}

	return post, nil
}

// Modify GetPostsByCategory to accept *sql.DB
func GetPostsByCategory(db *sql.DB, categoryID int) ([]*Post, error) {
	query := `SELECT post_id, title, content, user_id, category_id, timestamp FROM posts WHERE category_id = ? ORDER BY timestamp DESC`
	rows, err := db.Query(query, categoryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*Post
	for rows.Next() {
		post := &Post{}
		if err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.UserID, &post.CategoryID, &post.Timestamp); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}

// Modify GetPostsByUser to accept *sql.DB
func GetPostsByUser(db *sql.DB, userID int) ([]*Post, error) {
	query := `SELECT post_id, title, content, user_id, category_id, timestamp FROM posts WHERE user_id = ? ORDER BY timestamp DESC`
	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*Post
	for rows.Next() {
		post := &Post{}
		if err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.UserID, &post.CategoryID, &post.Timestamp); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}
