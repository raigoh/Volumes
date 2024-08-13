package post

import (
	"literary-lions-forum/internal/models"
	"literary-lions-forum/pkg/database"
)

func GetLatestPosts(limit int) ([]models.Post, error) {
	posts := []models.Post{}
	query := `SELECT id, user_id, title, content, created_at, updated_at 
              FROM posts 
              ORDER BY created_at DESC 
              LIMIT ?`

	rows, err := database.DB.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var post models.Post
		err := rows.Scan(&post.ID, &post.UserID, &post.Title, &post.Content, &post.CreatedAt, &post.UpdatedAt)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	return posts, nil
}

func CreatePost(userID, categoryID int, title, content string) (int, error) {
	result, err := database.DB.Exec("INSERT INTO posts (user_id, title, content) VALUES (?, ?, ?)", userID, title, content)
	if err != nil {
		return 0, err
	}

	postID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	_, err = database.DB.Exec("INSERT INTO post_categories (post_id, category_id) VALUES (?, ?)", postID, categoryID)
	if err != nil {
		return 0, err
	}

	return int(postID), nil
}

func GetPostByID(postID int) (models.Post, error) {
	var post models.Post
	err := database.DB.QueryRow(`
		SELECT p.id, p.user_id, p.title, p.content, p.created_at, u.username, c.id, c.name
		FROM posts p
		JOIN users u ON p.user_id = u.id
		JOIN post_categories pc ON p.id = pc.post_id
		JOIN categories c ON pc.category_id = c.id
		WHERE p.id = ?`, postID).Scan(
		&post.ID, &post.UserID, &post.Title, &post.Content, &post.CreatedAt,
		&post.User.Username, &post.Category.ID, &post.Category.Name)
	return post, err
}
