package models

import (
	"database/sql"
	"time"
)

type Comment struct {
	ID        int
	Content   string
	UserID    int
	PostID    int
	Timestamp time.Time
	Author    *User
}

// Modify CreateComment to accept *sql.DB
func CreateComment(db *sql.DB, comment *Comment) error {
	query := `INSERT INTO comments (content, user_id, post_id) VALUES (?, ?, ?)`
	result, err := db.Exec(query, comment.Content, comment.UserID, comment.PostID)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	comment.ID = int(id)
	return nil
}

// Modify GetCommentsByPostID to accept *sql.DB
func GetCommentsByPostID(db *sql.DB, postID int) ([]*Comment, error) {
	query := `SELECT comment_id, content, user_id, post_id, timestamp FROM comments WHERE post_id = ? ORDER BY timestamp ASC`
	rows, err := db.Query(query, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []*Comment
	for rows.Next() {
		comment := &Comment{}
		if err := rows.Scan(&comment.ID, &comment.Content, &comment.UserID, &comment.PostID, &comment.Timestamp); err != nil {
			return nil, err
		}
		comment.Author, err = GetUserByID(db, comment.UserID) // Pass db to GetUserByID
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}
	return comments, nil
}
