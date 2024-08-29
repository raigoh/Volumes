package database

import (
	"fmt"
	"literary-lions-forum/internal/models"
	"strconv"
)

// GetAllUsers returns all users in the database
func GetAllUsers() ([]models.User, error) {
	rows, err := DB.Query("SELECT id, username, email FROM users")
	if err != nil {
		return nil, fmt.Errorf("error querying users: %v", err)
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		var id int
		err := rows.Scan(&id, &user.Username, &user.Email)
		if err != nil {
			return nil, fmt.Errorf("error scanning user row: %v", err)
		}
		user.ID = strconv.Itoa(id) // Convert int to string
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating user rows: %v", err)
	}

	return users, nil
}

// ListUsers returns all users in the database
func ListUsers() ([]models.User, error) {
	rows, err := DB.Query("SELECT id, username, email FROM users")
	if err != nil {
		return nil, fmt.Errorf("error querying users: %v", err)
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.ID, &user.Username, &user.Email)
		if err != nil {
			return nil, fmt.Errorf("error scanning user row: %v", err)
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating user rows: %v", err)
	}

	return users, nil
}

// GetTotalUsers returns the total number of users in the database
func GetTotalUsers() (int, error) {
	var count int
	err := DB.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("error counting users: %v", err)
	}
	return count, nil
}

func GetUserByID(id string) (models.User, error) {
	var user models.User
	intID, err := strconv.Atoi(id)
	if err != nil {
		return user, err
	}

	err = DB.QueryRow("SELECT id, username, email, created_at FROM users WHERE id = ?", intID).Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt)
	if err != nil {
		return user, err
	}

	return user, nil
}
