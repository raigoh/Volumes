package database

import "fmt"

// GetTotalComments returns the total number of comments in the database
func GetTotalComments() (int, error) {
	var count int
	err := DB.QueryRow("SELECT COUNT(*) FROM comments").Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("error counting comments: %v", err)
	}
	return count, nil
}
