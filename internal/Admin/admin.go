package admin

import (
	"literary-lions-forum/pkg/database"
)

// IsUserAdmin checks if a user with the given ID has admin privileges.
//
// Parameters:
//   - userID: An integer representing the unique identifier of the user.
//
// Returns:
//   - bool: True if the user is an admin, false otherwise.
//   - error: An error object which will be non-nil if any database error occurs.
//
// This function performs a database query to check the 'is_admin' status of a user.
func IsUserAdmin(userID int) (bool, error) {
	// Initialize a boolean variable to store the admin status
	var isAdmin bool

	// Execute a database query to fetch the 'is_admin' value for the given user ID
	// The query selects the 'is_admin' column from the 'users' table where the 'id' matches the provided userID
	err := database.DB.QueryRow("SELECT is_admin FROM users WHERE id = ?", userID).Scan(&isAdmin)

	// Return the admin status and any error that occurred during the database operation
	// If err is nil, the operation was successful
	// If err is non-nil, it indicates a database error (e.g., no such user, connection issues, etc.)
	return isAdmin, err
}
