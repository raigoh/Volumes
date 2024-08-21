package comment

import (
	"literary-lions-forum/internal/models"
	"literary-lions-forum/pkg/database"
	"strconv"
)

func GetCommentsByPostID(postID int) ([]models.Comment, error) {
	query := `
			SELECT c.id, c.user_id, c.content, c.created_at, u.username,
						 (SELECT COUNT(*) FROM likes WHERE comment_id = c.id AND is_like = 1) as likes,
						 (SELECT COUNT(*) FROM likes WHERE comment_id = c.id AND is_like = 0) as dislikes
			FROM comments c
			JOIN users u ON c.user_id = u.id
			WHERE c.post_id = ?
			ORDER BY c.created_at DESC
	`
	rows, err := database.DB.Query(query, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []models.Comment
	for rows.Next() {
		var comment models.Comment
		var userID int
		err := rows.Scan(
			&comment.ID, &userID, &comment.Content, &comment.CreatedAt,
			&comment.User.Username, &comment.Likes, &comment.Dislikes,
		)
		if err != nil {
			return nil, err
		}
		comment.User.ID = strconv.Itoa(userID)
		comments = append(comments, comment)
	}

	return comments, nil
}

func CreateComment(userID, postID int, content string) error {
	//fmt.Println("New comment to ", postID, " from ", userID, " with content: ", content)
	_, err := database.DB.Exec("INSERT INTO comments (user_id, post_id, content) VALUES (?, ?, ?)", userID, postID, content)
	return err
}
