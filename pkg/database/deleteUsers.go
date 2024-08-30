package database

import (
	"fmt"
	"strconv"
)

// DeleteUser deletes a user from the database by their ID
func DeleteUser(userID string) error {
	// Convert the userID string to an integer
	id, err := strconv.Atoi(userID)
	if err != nil {
		// If conversion fails, return an error with context
		return fmt.Errorf("invalid user ID: %v", err)
	}

	// Execute a DELETE SQL command to remove the user with the given ID
	_, err = DB.Exec("DELETE FROM users WHERE id = ?", id)
	if err != nil {
		// If the deletion fails, return an error with context
		return fmt.Errorf("error deleting user: %v", err)
	}

	// If deletion is successful, return nil (no error)
	return nil
}

// DeleteUserByEmail deletes a user from the database by their email
func DeleteUserByEmail(email string) error {
	// Execute a DELETE SQL command to remove the user with the given email
	_, err := DB.Exec("DELETE FROM users WHERE email = ?", email)
	if err != nil {
		// If the deletion fails, return an error with context
		return fmt.Errorf("error deleting user by email: %v", err)
	}

	// If deletion is successful, return nil (no error)
	return nil
}
