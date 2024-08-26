package post

import (
	"literary-lions-forum/internal/models"
	"literary-lions-forum/internal/utils"
	"literary-lions-forum/pkg/database"
	"net/http"
	"strconv"
	"strings"
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

func GetPostByID(postID int) (*models.Post, error) {
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

func GetFilteredPosts(categoryID, userID int, likedOnly bool, limit int) ([]models.Post, error) {
	query := `
			SELECT DISTINCT p.id, p.title, p.content, p.created_at, 
						 u.id AS user_id, u.username,
						 c.id AS category_id, c.name AS category_name
			FROM posts p
			JOIN users u ON p.user_id = u.id
			LEFT JOIN post_categories pc ON p.id = pc.post_id
			LEFT JOIN categories c ON pc.category_id = c.id
	`

	var args []interface{}
	var conditions []string

	if categoryID > 0 {
		conditions = append(conditions, "c.id = ?")
		args = append(args, categoryID)
	}

	if userID > 0 {
		conditions = append(conditions, "p.user_id = ?")
		args = append(args, userID)
	}

	if likedOnly && userID > 0 {
		query += " INNER JOIN likes l ON p.id = l.post_id"
		conditions = append(conditions, "l.user_id = ? AND l.is_like = TRUE")
		args = append(args, userID)
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	query += " ORDER BY p.created_at DESC LIMIT ?"
	args = append(args, limit)

	rows, err := database.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var p models.Post
		err := rows.Scan(
			&p.ID, &p.Title, &p.Content, &p.CreatedAt,
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

func PostListHandler(w http.ResponseWriter, r *http.Request) {
	posts, err := database.GetPosts() // Implement this function to fetch posts with like/dislike counts
	if err != nil {
		http.Error(w, "Error fetching posts", http.StatusInternalServerError)
		return
	}

	data := models.PageData{
		Title: "Literary Lions Forum",
		Page:  "home",
		Data: map[string]interface{}{
			"Posts": posts,
		},
	}

	utils.RenderTemplate(w, "home.html", data)
}
