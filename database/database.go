package database

import (
	"database/sql"
)

func InitializeTables(db *sql.DB) error {
	statements := []string{
		`CREATE TABLE IF NOT EXISTS users (
			user_id INTEGER PRIMARY KEY AUTOINCREMENT,
			email TEXT UNIQUE NOT NULL,
			username TEXT UNIQUE NOT NULL,
			password TEXT NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS categories (
			category_id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT UNIQUE NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS posts (
			post_id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL,
			content TEXT NOT NULL,
			user_id INTEGER NOT NULL,
			category_id INTEGER NOT NULL,
			timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(user_id),
			FOREIGN KEY (category_id) REFERENCES categories(category_id)
		)`,
		`CREATE TABLE IF NOT EXISTS comments (
			comment_id INTEGER PRIMARY KEY AUTOINCREMENT,
			content TEXT NOT NULL,
			user_id INTEGER NOT NULL,
			post_id INTEGER NOT NULL,
			timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(user_id),
			FOREIGN KEY (post_id) REFERENCES posts(post_id)
		)`,
		`CREATE TABLE IF NOT EXISTS likes_dislikes (
			like_id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			post_id INTEGER,
			comment_id INTEGER,
			is_like BOOLEAN NOT NULL,
			FOREIGN KEY (user_id) REFERENCES users(user_id),
			FOREIGN KEY (post_id) REFERENCES posts(post_id),
			FOREIGN KEY (comment_id) REFERENCES comments(comment_id),
			CHECK (post_id IS NOT NULL OR comment_id IS NOT NULL)
		)`,
	}

	for _, stmt := range statements {
		_, err := db.Exec(stmt)
		if err != nil {
			return err
		}
	}

	return nil
}
