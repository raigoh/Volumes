package database

import (
	"database/sql"
	"fmt"
	"os"

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
