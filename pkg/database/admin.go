package database

import "fmt"

// addIsAdminColumn adds an 'is_admin' column to the users table if it doesn't already exist
func addIsAdminColumn() error {
	// Check if the column already exists in the users table
	var count int
	err := DB.QueryRow("SELECT COUNT(*) FROM pragma_table_info('users') WHERE name='is_admin'").Scan(&count)
	if err != nil {
		// If there's an error checking for the column, wrap and return the error
		return fmt.Errorf("error checking for is_admin column: %v", err)
	}

	// If the column doesn't exist (count is 0), add it to the table
	if count == 0 {
		// Execute the ALTER TABLE command to add the new column
		_, err := DB.Exec("ALTER TABLE users ADD COLUMN is_admin BOOLEAN DEFAULT FALSE")
		if err != nil {
			// If there's an error adding the column, wrap and return the error
			return fmt.Errorf("error adding is_admin column: %v", err)
		}
		// Log a message indicating the column was successfully added
		fmt.Println("Added is_admin column to users table")
	}

	// If the column already existed or was successfully added, return nil (no error)
	return nil
}
