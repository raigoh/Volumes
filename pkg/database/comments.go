package database

import "fmt"

// GetTotalComments returns the total number of comments in the database
func GetTotalComments() (int, error) {
	// Variable to store the count of comments
	var count int

	// Execute a SQL query to count all rows in the comments table
	// The result of COUNT(*) is then scanned into the count variable
	err := DB.QueryRow("SELECT COUNT(*) FROM comments").Scan(&count)

	// Check if there was an error during the query execution or scanning
	if err != nil {
		// If an error occurred, return 0 for the count and wrap the error with additional context
		return 0, fmt.Errorf("error counting comments: %v", err)
	}

	// If successful, return the count and nil for the error
	return count, nil
}
