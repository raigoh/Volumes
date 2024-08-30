package comment

import (
	"literary-lions-forum/internal/models"
	"literary-lions-forum/pkg/database"
	"strconv"
)

// GetCommentsByPostID retrieves all comments for a specific post, including user information and like/dislike counts.
// It returns a slice of Comment structs and an error if any occurs during the database operation.
func GetCommentsByPostID(postID int) ([]models.Comment, error) {
	// SQL query to fetch comments with user info and like/dislike counts
	query := `
			SELECT c.id, c.user_id, c.content, c.created_at, u.username,
						 (SELECT COUNT(*) FROM likes WHERE comment_id = c.id AND is_like = 1) as likes,
						 (SELECT COUNT(*) FROM likes WHERE comment_id = c.id AND is_like = 0) as dislikes
			FROM comments c
			JOIN users u ON c.user_id = u.id
			WHERE c.post_id = ?
			ORDER BY c.created_at DESC
	`
	// Execute the query with the provided post ID
	rows, err := database.DB.Query(query, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []models.Comment
	// Iterate through the result set
	for rows.Next() {
		var comment models.Comment
		var userID int
		// Scan the row into the comment struct and separate variables
		err := rows.Scan(
			&comment.ID, &userID, &comment.Content, &comment.CreatedAt,
			&comment.User.Username, &comment.Likes, &comment.Dislikes,
		)
		if err != nil {
			return nil, err
		}
		// Convert userID to string for the User struct
		comment.User.ID = strconv.Itoa(userID)
		comments = append(comments, comment)
	}

	return comments, nil
}

// CreateComment adds a new comment to the database.
// It takes the user ID, post ID, and comment content as parameters.
// Returns an error if the database operation fails.
func CreateComment(userID, postID int, content string) error {
	// Uncomment the following line for debugging purposes
	//fmt.Println("New comment to ", postID, " from ", userID, " with content: ", content)

	// Execute an INSERT query to add the new comment
	_, err := database.DB.Exec("INSERT INTO comments (user_id, post_id, content) VALUES (?, ?, ?)", userID, postID, content)
	return err
}
