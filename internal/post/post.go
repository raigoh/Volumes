package post

import (
	"literary-lions-forum/internal/models"
	"literary-lions-forum/pkg/database"
	"strconv"
	"strings"
)

// GetLatestPosts retrieves the most recent posts with a specified limit.
// It fetches posts along with their associated user and category information.
func GetLatestPosts(limit int) ([]models.Post, error) {
	// SQL query to fetch the latest posts with user and category data.
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
	// Execute the query with the specified limit on the number of posts.
	rows, err := database.DB.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Initialize a slice to hold the retrieved posts.
	var posts []models.Post

	// Iterate through the rows of the result set.
	for rows.Next() {
		var p models.Post
		// Scan each row's data into the Post object fields.
		err := rows.Scan(
			&p.ID, &p.Title, &p.CreatedAt,
			&p.User.ID, &p.User.Username,
			&p.Category.ID, &p.Category.Name,
		)
		if err != nil {
			return nil, err
		}
		// Append the constructed post to the posts slice.
		posts = append(posts, p)
	}
	return posts, nil
}

// CreatePost inserts a new post into the database with the associated user and category.
// It returns the ID of the newly created post.
func CreatePost(userID, categoryID int, title, content string) (int, error) {
	// Insert the new post into the posts table.
	result, err := database.DB.Exec("INSERT INTO posts (user_id, title, content) VALUES (?, ?, ?)", userID, title, content)
	if err != nil {
		return 0, err
	}

	// Get the ID of the newly inserted post.
	postID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	// Associate the post with the selected category by inserting into the post_categories table.
	_, err = database.DB.Exec("INSERT INTO post_categories (post_id, category_id) VALUES (?, ?)", postID, categoryID)
	if err != nil {
		return 0, err
	}

	// Return the ID of the newly created post.
	return int(postID), nil
}

// GetPostByID retrieves a single post by its ID.
// It includes the post's user information, like, and dislike counts.
func GetPostByID(postID int) (*models.Post, error) {
	// SQL query to fetch post details along with user data and like/dislike counts.
	query := `
        SELECT p.id, p.user_id, p.title, p.content, p.created_at, u.username,
               c.id AS category_id, c.name AS category_name,
               (SELECT COUNT(*) FROM likes WHERE post_id = p.id AND is_like = 1) as likes,
               (SELECT COUNT(*) FROM likes WHERE post_id = p.id AND is_like = 0) as dislikes
        FROM posts p
        JOIN users u ON p.user_id = u.id
        LEFT JOIN post_categories pc ON p.id = pc.post_id
        LEFT JOIN categories c ON pc.category_id = c.id
        WHERE p.id = ?
    `
	var post models.Post
	var userID int

	// Execute the query and scan the results into the Post object.
	err := database.DB.QueryRow(query, postID).Scan(
		&post.ID, &userID, &post.Title, &post.Content, &post.CreatedAt,
		&post.User.Username,
		&post.Category.ID, &post.Category.Name,
		&post.Likes, &post.Dislikes,
	)
	if err != nil {
		return nil, err
	}

	// Convert the userID from integer to string and assign it to the post's User.
	post.User.ID = strconv.Itoa(userID)
	return &post, nil
}

// GetFilteredPosts retrieves posts based on various filter criteria such as category ID, user ID, liked posts only, and a limit.
func GetFilteredPosts(categoryID, userID int, likedOnly bool, limit int) ([]models.Post, error) {
	// Base SQL query for retrieving posts with user and category information, and like/dislike counts.
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

	// Slice to hold filter conditions and arguments for the query.
	var conditions []string
	var args []interface{}

	// Filter posts by category if a valid categoryID is provided.
	if categoryID > 0 {
		conditions = append(conditions, "c.id = ?")
		args = append(args, categoryID)
	}

	// Filter posts by user if a valid userID is provided.
	if userID > 0 {
		if likedOnly {
			// If likedOnly is true, join the likes table and filter by posts liked by the user.
			query += " JOIN likes l ON p.id = l.post_id"
			conditions = append(conditions, "l.user_id = ? AND l.is_like = TRUE")
		} else {
			// Otherwise, filter by posts created by the user.
			conditions = append(conditions, "p.user_id = ?")
		}
		args = append(args, userID)
	}

	// Combine the conditions and add them to the query if they exist.
	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	// Add ordering and limit to the query.
	query += " ORDER BY p.created_at DESC LIMIT ?"
	args = append(args, limit)

	// Execute the query with the constructed arguments.
	rows, err := database.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Initialize a slice to hold the retrieved posts.
	var posts []models.Post

	// Iterate through the results and construct Post objects.
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
		// Append the constructed post to the posts slice.
		posts = append(posts, p)
	}

	return posts, nil
}

// SearchPosts retrieves posts based on a search query (searches both title and content).
// It limits the results to the specified number of posts.
func SearchPosts(query string, limit int) ([]models.Post, error) {
	// SQL query to search for posts where the title or content matches the search query.
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

	// Execute the query with the search query parameter (wildcards added for partial matches).
	rows, err := database.DB.Query(sqlQuery, "%"+query+"%", "%"+query+"%", limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Initialize a slice to hold the retrieved posts.
	var posts []models.Post

	// Iterate through the result set and construct Post objects.
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
		// Append the constructed post to the posts slice.
		posts = append(posts, post)
	}

	return posts, nil
}

// GetLikedPosts retrieves posts that the user has liked.
func GetLikedPosts(userID, limit int) ([]models.Post, error) {
	var posts []models.Post

	query := `
			SELECT p.id, p.title, p.content, p.created_at, 
						 c.id AS category_id, c.name AS category_name,
						 u.username AS user_username
			FROM posts p
			JOIN likes l ON l.post_id = p.id
			JOIN users u ON u.id = p.user_id
			LEFT JOIN post_categories pc ON p.id = pc.post_id
			LEFT JOIN categories c ON pc.category_id = c.id
			WHERE l.user_id = ? AND l.is_like = TRUE
			ORDER BY p.created_at DESC
			LIMIT ?
	`

	rows, err := database.DB.Query(query, userID, limit)
	if err != nil {
		return posts, err
	}
	defer rows.Close()

	for rows.Next() {
		var post models.Post
		err := rows.Scan(
			&post.ID, &post.Title, &post.Content, &post.CreatedAt,
			&post.Category.ID, &post.Category.Name,
			&post.User.Username,
		)
		if err != nil {
			return posts, err
		}
		posts = append(posts, post)
	}

	return posts, nil
}
