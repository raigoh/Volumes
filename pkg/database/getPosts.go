package database

import (
	"fmt"
	"literary-lions-forum/internal/models"
	"log"
	"strconv"
)

// GetPosts retrieves all posts from the database with associated user, category, and like information
func GetPosts() ([]models.Post, error) {
	// SQL query to fetch posts with related information
	query := `
			SELECT p.id, p.user_id, p.title, p.content, p.created_at, u.username,
						 c.id AS category_id, c.name AS category_name,
						 (SELECT COUNT(*) FROM likes WHERE post_id = p.id AND is_like = 1) AS likes,
						 (SELECT COUNT(*) FROM likes WHERE post_id = p.id AND is_like = 0) AS dislikes
			FROM posts p
			JOIN users u ON p.user_id = u.id
			LEFT JOIN post_categories pc ON p.id = pc.post_id
			LEFT JOIN categories c ON pc.category_id = c.id
			ORDER BY p.created_at DESC
	`

	// Execute the query
	rows, err := DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []models.Post
	// Iterate through the result rows
	for rows.Next() {
		var p models.Post
		// Scan the row into a Post struct
		err := rows.Scan(
			&p.ID, &p.User.ID, &p.Title, &p.Content, &p.CreatedAt, &p.User.Username,
			&p.Category.ID, &p.Category.Name, &p.Likes, &p.Dislikes,
		)
		if err != nil {
			return nil, err
		}
		log.Printf("Fetched post ID: %d, Likes: %d, Dislikes: %d", p.ID, p.Likes, p.Dislikes)
		posts = append(posts, p)
	}

	return posts, nil
}

// GetUserPosts retrieves all posts from a specific user
func GetUserPosts(userID string) ([]models.Post, error) {
	// Convert string userID to integer
	intID, err := strconv.Atoi(userID)
	if err != nil {
		return nil, err
	}

	// SQL query to fetch posts for a specific user
	query := `
		SELECT p.id, p.title, p.content, p.created_at, 
			   c.id AS category_id, c.name AS category_name,
			   (SELECT COUNT(*) FROM likes WHERE post_id = p.id AND is_like = 1) AS likes,
			   (SELECT COUNT(*) FROM likes WHERE post_id = p.id AND is_like = 0) AS dislikes
		FROM posts p
		LEFT JOIN post_categories pc ON p.id = pc.post_id
		LEFT JOIN categories c ON pc.category_id = c.id
		WHERE p.user_id = ?
		ORDER BY p.created_at DESC
	`

	// Execute the query with the user ID as a parameter
	rows, err := DB.Query(query, intID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []models.Post
	// Iterate through the result rows
	for rows.Next() {
		var p models.Post
		// Scan the row into a Post struct
		err := rows.Scan(
			&p.ID, &p.Title, &p.Content, &p.CreatedAt,
			&p.Category.ID, &p.Category.Name,
			&p.Likes, &p.Dislikes,
		)
		if err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}

	return posts, nil
}

// GetTotalPosts returns the total number of posts in the database
func GetTotalPosts() (int, error) {
	var count int
	// Execute a simple COUNT query
	err := DB.QueryRow("SELECT COUNT(*) FROM posts").Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("error counting posts: %v", err)
	}
	return count, nil
}
