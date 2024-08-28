package database

import (
	"database/sql"
	"fmt"
	"literary-lions-forum/internal/models"
	"log"
	"os"
	"strconv"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB() error {
	var err error
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "./forum.db" // Default path if not set
	}
	DB, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		return fmt.Errorf("failed to open database: %v", err)
	}

	err = createTables()
	if err != nil {
		return fmt.Errorf("failed to create tables: %v", err)
	}

	// Check if is_admin column exists, if not, add it
	err = addIsAdminColumn()
	if err != nil {
		return fmt.Errorf("failed to add is_admin column: %v", err)
	}

	fmt.Println("Database and tables initialized successfully at:", dbPath)
	return nil
}

func addIsAdminColumn() error {
	// Check if the column exists
	var count int
	err := DB.QueryRow("SELECT COUNT(*) FROM pragma_table_info('users') WHERE name='is_admin'").Scan(&count)
	if err != nil {
		return fmt.Errorf("error checking for is_admin column: %v", err)
	}

	// If the column doesn't exist, add it
	if count == 0 {
		_, err := DB.Exec("ALTER TABLE users ADD COLUMN is_admin BOOLEAN DEFAULT FALSE")
		if err != nil {
			return fmt.Errorf("error adding is_admin column: %v", err)
		}
		fmt.Println("Added is_admin column to users table")
	}

	return nil
}

func createTables() error {
	createUsersTable := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT UNIQUE NOT NULL,
		email TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		is_admin BOOLEAN DEFAULT FALSE
	);`

	createPostsTable := `
	CREATE TABLE IF NOT EXISTS posts (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER,
		title TEXT NOT NULL,
		content TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(id)
	);`

	createCommentsTable := `
	CREATE TABLE IF NOT EXISTS comments (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		post_id INTEGER,
		user_id INTEGER,
		content TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (post_id) REFERENCES posts(id),
		FOREIGN KEY (user_id) REFERENCES users(id)
	);`

	createLikesTable := `
	CREATE TABLE IF NOT EXISTS likes (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER,
    post_id INTEGER,
    comment_id INTEGER,
    is_like BOOLEAN,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (post_id) REFERENCES posts(id),
    FOREIGN KEY (comment_id) REFERENCES comments(id),
    UNIQUE(user_id, post_id, comment_id)
	);`

	createCategoriesTable := `
	CREATE TABLE IF NOT EXISTS categories (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT UNIQUE NOT NULL
	);`

	createPostCategoriesTable := `
	CREATE TABLE IF NOT EXISTS post_categories (
		post_id INTEGER,
		category_id INTEGER,
		PRIMARY KEY (post_id, category_id),
		FOREIGN KEY (post_id) REFERENCES posts(id),
		FOREIGN KEY (category_id) REFERENCES categories(id)
	);`

	tables := []string{
		createUsersTable,
		createPostsTable,
		createCommentsTable,
		createLikesTable,
		createCategoriesTable,
		createPostCategoriesTable,
	}

	for _, table := range tables {
		_, err := DB.Exec(table)
		if err != nil {
			return err
		}
	}

	return nil
}

// GetAllUsers returns all users in the database
func GetAllUsers() ([]models.User, error) {
	rows, err := DB.Query("SELECT id, username, email FROM users")
	if err != nil {
		return nil, fmt.Errorf("error querying users: %v", err)
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		var id int
		err := rows.Scan(&id, &user.Username, &user.Email)
		if err != nil {
			return nil, fmt.Errorf("error scanning user row: %v", err)
		}
		user.ID = strconv.Itoa(id) // Convert int to string
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating user rows: %v", err)
	}

	return users, nil
}

// DeleteUser deletes a user from the database by their ID
func DeleteUser(userID string) error {
	id, err := strconv.Atoi(userID)
	if err != nil {
		return fmt.Errorf("invalid user ID: %v", err)
	}

	_, err = DB.Exec("DELETE FROM users WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("error deleting user: %v", err)
	}
	return nil
}

// DeleteUserByEmail deletes a user from the database by their email
func DeleteUserByEmail(email string) error {
	_, err := DB.Exec("DELETE FROM users WHERE email = ?", email)
	if err != nil {
		return fmt.Errorf("error deleting user by email: %v", err)
	}
	return nil
}

// ListUsers returns all users in the database
func ListUsers() ([]models.User, error) {
	rows, err := DB.Query("SELECT id, username, email FROM users")
	if err != nil {
		return nil, fmt.Errorf("error querying users: %v", err)
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.ID, &user.Username, &user.Email)
		if err != nil {
			return nil, fmt.Errorf("error scanning user row: %v", err)
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating user rows: %v", err)
	}

	return users, nil
}

// GetTotalUsers returns the total number of users in the database
func GetTotalUsers() (int, error) {
	var count int
	err := DB.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("error counting users: %v", err)
	}
	return count, nil
}

// GetTotalPosts returns the total number of posts in the database
func GetTotalPosts() (int, error) {
	var count int
	err := DB.QueryRow("SELECT COUNT(*) FROM posts").Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("error counting posts: %v", err)
	}
	return count, nil
}

// GetTotalComments returns the total number of comments in the database
func GetTotalComments() (int, error) {
	var count int
	err := DB.QueryRow("SELECT COUNT(*) FROM comments").Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("error counting comments: %v", err)
	}
	return count, nil
}

// GetActiveUsers returns the number of users who have been active in the last n days
func GetActiveUsers(days int) (int, error) {
	var count int
	cutoffDate := time.Now().AddDate(0, 0, -days)
	query := `
		SELECT COUNT(DISTINCT user_id) 
		FROM (
			SELECT user_id FROM posts WHERE created_at > ?
			UNION
			SELECT user_id FROM comments WHERE created_at > ?
		)
	`
	err := DB.QueryRow(query, cutoffDate, cutoffDate).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("error counting active users: %v", err)
	}
	return count, nil
}

// GetRecentActivity returns the n most recent activities (posts or comments)
func GetRecentActivity(n int) ([]models.Activity, error) {
	query := `
		SELECT 'post' as type, u.username, p.title as content, p.created_at as created_at
		FROM posts p
		JOIN users u ON p.user_id = u.id
		UNION ALL
		SELECT 'comment' as type, u.username, c.content, c.created_at as created_at
		FROM comments c
		JOIN users u ON c.user_id = u.id
		ORDER BY created_at DESC
		LIMIT ?
	`
	rows, err := DB.Query(query, n)
	if err != nil {
		return nil, fmt.Errorf("error querying recent activity: %v", err)
	}
	defer rows.Close()

	var activities []models.Activity
	for rows.Next() {
		var a models.Activity
		var createdAt time.Time
		err := rows.Scan(&a.Type, &a.Username, &a.Content, &createdAt)
		if err != nil {
			return nil, fmt.Errorf("error scanning activity row: %v", err)
		}
		a.Date = createdAt.Format("02.01.2006 15:04")
		activities = append(activities, a)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating activity rows: %v", err)
	}

	return activities, nil
}

func GetPosts() ([]models.Post, error) {
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

	rows, err := DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var p models.Post
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

func GetUserByID(id string) (models.User, error) {
	var user models.User
	intID, err := strconv.Atoi(id)
	if err != nil {
		return user, err
	}

	err = DB.QueryRow("SELECT id, username, email, created_at FROM users WHERE id = ?", intID).Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt)
	if err != nil {
		return user, err
	}

	return user, nil
}

func GetUserPosts(userID string) ([]models.Post, error) {
	intID, err := strconv.Atoi(userID)
	if err != nil {
		return nil, err
	}

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

	rows, err := DB.Query(query, intID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var p models.Post
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
