package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3" // Import SQLite driver
)

// DB is a global variable to hold the database connection
var DB *sql.DB

// InitDB initializes the database connection and sets up tables
func InitDB() error {
	var err error

	// Get the database path from environment variable or use default
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = ".data/forum.db" // Default path if not set
	}

	// Open the SQLite database
	DB, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		return fmt.Errorf("failed to open database: %v", err)
	}

	// Create necessary tables in the database
	err = createTables()
	if err != nil {
		return fmt.Errorf("failed to create tables: %v", err)
	}

	// Add 'is_admin' column to users table if it doesn't exist
	err = addIsAdminColumn()
	if err != nil {
		return fmt.Errorf("failed to add is_admin column: %v", err)
	}

	fmt.Println("Database and tables initialized successfully at:", dbPath)
	return nil
}
