package models

import (
	"database/sql"
)

type LikeDislike struct {
	ID        int
	UserID    int
	PostID    *int
	CommentID *int
	IsLike    bool
}

func CreateLikeDislike(ld *LikeDislike) error {
	query := `INSERT INTO likes_dislikes (user_id, post_id, comment_id, is_like) VALUES (?, ?, ?, ?)`
	result, err := db.Exec(query, ld.UserID, ld.PostID, ld.CommentID, ld.IsLike)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	ld.ID = int(id)
	return nil
}

func GetLikesDislikes(postID, commentID *int) (likes, dislikes int, err error) {
	var query string
	var args []interface{}

	if postID != nil {
		query = `SELECT COUNT(*) FROM likes_dislikes WHERE post_id = ? AND is_like = ?`
		args = []interface{}{*postID}
	} else if commentID != nil {
		query = `SELECT COUNT(*) FROM likes_dislikes WHERE comment_id = ? AND is_like = ?`
		args = []interface{}{*commentID}
	} else {
		return 0, 0, sql.ErrNoRows
	}

	err = db.QueryRow(query, append(args, true)...).Scan(&likes)
	if err != nil {
		return 0, 0, err
	}

	err = db.QueryRow(query, append(args, false)...).Scan(&dislikes)
	if err != nil {
		return 0, 0, err
	}

	return likes, dislikes, nil
}

// LikePost adds a like to a post
func LikePost(postID int) error {
	likeDislike := &LikeDislike{
		UserID: 1, // Replace with actual user ID from session
		PostID: &postID,
		IsLike: true,
	}
	return CreateLikeDislike(likeDislike)
}

// DislikePost adds a dislike to a post
func DislikePost(postID int) error {
	likeDislike := &LikeDislike{
		UserID: 1, // Replace with actual user ID from session
		PostID: &postID,
		IsLike: false,
	}
	return CreateLikeDislike(likeDislike)
}

// LikeComment adds a like to a comment
func LikeComment(commentID int) error {
	likeDislike := &LikeDislike{
		UserID:    1, // Replace with actual user ID from session
		CommentID: &commentID,
		IsLike:    true,
	}
	return CreateLikeDislike(likeDislike)
}

// DislikeComment adds a dislike to a comment
func DislikeComment(commentID int) error {
	likeDislike := &LikeDislike{
		UserID:    1, // Replace with actual user ID from session
		CommentID: &commentID,
		IsLike:    false,
	}
	return CreateLikeDislike(likeDislike)
}
