package post

import (
	"literary-lions-forum/internal/models"
	"literary-lions-forum/pkg/database"
)

func GetLatestPosts(limit int) ([]models.Post, error) {
	query := `
			SELECT p.id, p.title, p.created_at, 
						 u.id AS user_id, u.username,
						 c.id AS category_id, c.name AS category_name
			FROM posts p
			JOIN users u ON p.user_id = u.id
			LEFT JOIN post_categories pc ON p.id = pc.post_id
			LEFT JOIN categories c ON pc.category_id = c.id
			ORDER BY p.created_at DESC
			LIMIT ?
	`
	rows, err := database.DB.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var p models.Post
		err := rows.Scan(
			&p.ID, &p.Title, &p.CreatedAt,
			&p.User.ID, &p.User.Username,
			&p.Category.ID, &p.Category.Name,
		)
		if err != nil {
			return nil, err
		}
		posts = append(posts, p)
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

func GetPostByID(id int) (models.Post, error) {
	var post models.Post
	query := `
			SELECT p.id, p.title, p.content, p.created_at, p.updated_at, 
						 u.id AS user_id, u.username,
						 c.id AS category_id, c.name AS category_name
			FROM posts p
			JOIN users u ON p.user_id = u.id
			LEFT JOIN post_categories pc ON p.id = pc.post_id
			LEFT JOIN categories c ON pc.category_id = c.id
			WHERE p.id = ?
	`
	err := database.DB.QueryRow(query, id).Scan(
		&post.ID, &post.Title, &post.Content, &post.CreatedAt, &post.UpdatedAt,
		&post.User.ID, &post.User.Username,
		&post.Category.ID, &post.Category.Name,
	)
	if err != nil {
		return post, err
	}
	return post, nil
}
