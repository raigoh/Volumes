package post

import (
	"literary-lions-forum/internal/models"
	"literary-lions-forum/pkg/database"
	"strconv"
	"strings"
)

// GetLatestPosts retrieves the most recent posts with a specified limit
func GetLatestPosts(limit int) ([]models.Post, error) {
	// SQL query to fetch latest posts with user and category information
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
	// Execute the query
	rows, err := database.DB.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Iterate through the results and construct Post objects
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

// CreatePost inserts a new post into the database
func CreatePost(userID, categoryID int, title, content string) (int, error) {
	// Insert the post
	result, err := database.DB.Exec("INSERT INTO posts (user_id, title, content) VALUES (?, ?, ?)", userID, title, content)
	if err != nil {
		return 0, err
	}

	// Get the ID of the newly inserted post
	postID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	// Associate the post with the category
	_, err = database.DB.Exec("INSERT INTO post_categories (post_id, category_id) VALUES (?, ?)", postID, categoryID)
	if err != nil {
		return 0, err
	}

	return int(postID), nil
}

// GetPostByID retrieves a single post by its ID, including like and dislike counts
func GetPostByID(postID int) (*models.Post, error) {
	// SQL query to fetch post details with user information and like/dislike counts
	query := `
			SELECT p.id, p.user_id, p.title, p.content, p.created_at, u.username,
						 (SELECT COUNT(*) FROM likes WHERE post_id = p.id AND is_like = 1) as likes,
						 (SELECT COUNT(*) FROM likes WHERE post_id = p.id AND is_like = 0) as dislikes
			FROM posts p
			JOIN users u ON p.user_id = u.id
			WHERE p.id = ?
	`
	var post models.Post
	var userID int
	// Execute the query and scan the result into the post struct
	err := database.DB.QueryRow(query, postID).Scan(
		&post.ID, &userID, &post.Title, &post.Content, &post.CreatedAt,
		&post.User.Username, &post.Likes, &post.Dislikes,
	)
	if err != nil {
		return nil, err
	}
	post.User.ID = strconv.Itoa(userID)
	return &post, nil
}

// GetFilteredPosts retrieves posts based on various filter criteria
func GetFilteredPosts(categoryID, userID int, likedOnly bool, limit int) ([]models.Post, error) {
	// Base SQL query
	query := `
			SELECT DISTINCT p.id, p.title, p.content, p.created_at, 
									 u.id AS user_id, u.username,
									 c.id AS category_id, c.name AS category_name,
									 (SELECT COUNT(*) FROM likes WHERE post_id = p.id AND is_like = 1) as likes,
									 (SELECT COUNT(*) FROM likes WHERE post_id = p.id AND is_like = 0) as dislikes
			FROM posts p
			JOIN users u ON p.user_id = u.id
			LEFT JOIN post_categories pc ON p.id = pc.post_id
			LEFT JOIN categories c ON pc.category_id = c.id
	`

	var args []interface{}
	var conditions []string

	// Add filter conditions based on input parameters
	if categoryID > 0 {
		conditions = append(conditions, "c.id = ?")
		args = append(args, categoryID)
	}

	if userID > 0 {
		conditions = append(conditions, "p.user_id = ?")
		args = append(args, userID)
	}

	if likedOnly {
		query += " INNER JOIN likes l ON p.id = l.post_id"
		conditions = append(conditions, "l.user_id = ? AND l.is_like = TRUE")
		args = append(args, userID)
	}

	// Combine all conditions
	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	// Add ordering and limit
	query += " ORDER BY p.created_at DESC LIMIT ?"
	args = append(args, limit)

	// Execute the query
	rows, err := database.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Iterate through the results and construct Post objects
	var posts []models.Post
	for rows.Next() {
		var p models.Post
		err := rows.Scan(
			&p.ID, &p.Title, &p.Content, &p.CreatedAt,
			&p.User.ID, &p.User.Username,
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

func SearchPosts(query string, limit int) ([]models.Post, error) {
	sqlQuery := `
			SELECT DISTINCT p.id, p.title, p.content, p.created_at, 
									 u.id AS user_id, u.username,
									 c.id AS category_id, c.name AS category_name,
									 (SELECT COUNT(*) FROM likes WHERE post_id = p.id AND is_like = 1) as likes,
									 (SELECT COUNT(*) FROM likes WHERE post_id = p.id AND is_like = 0) as dislikes
			FROM posts p
			JOIN users u ON p.user_id = u.id
			LEFT JOIN post_categories pc ON p.id = pc.post_id
			LEFT JOIN categories c ON pc.category_id = c.id
			WHERE p.title LIKE ? OR p.content LIKE ?
			ORDER BY p.created_at DESC
			LIMIT ?
	`

	rows, err := database.DB.Query(sqlQuery, "%"+query+"%", "%"+query+"%", limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var post models.Post
		err := rows.Scan(
			&post.ID, &post.Title, &post.Content, &post.CreatedAt,
			&post.User.ID, &post.User.Username,
			&post.Category.ID, &post.Category.Name,
			&post.Likes, &post.Dislikes,
		)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	return posts, nil
}
