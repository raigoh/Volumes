package database

import (
	"fmt"
	"strconv"
)

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
