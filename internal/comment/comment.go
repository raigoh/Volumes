package comment

import (
	"literary-lions-forum/internal/models"
	"literary-lions-forum/pkg/database"
)

func GetCommentsByPostID(postID int) ([]models.Comment, error) {
	comments := []models.Comment{}
	rows, err := database.DB.Query(`
		SELECT c.id, c.user_id, c.content, c.created_at, u.username
		FROM comments c
		JOIN users u ON c.user_id = u.id
		WHERE c.post_id = ?
		ORDER BY c.created_at DESC`, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var comment models.Comment
		err := rows.Scan(&comment.ID, &comment.UserID, &comment.Content, &comment.CreatedAt, &comment.User.Username)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}

	return comments, nil
}

func CreateComment(userID, postID int, content string) error {
	_, err := database.DB.Exec("INSERT INTO comments (user_id, post_id, content) VALUES (?, ?, ?)", userID, postID, content)
	return err
}
