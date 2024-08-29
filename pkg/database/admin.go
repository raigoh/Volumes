package database

import "fmt"

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
