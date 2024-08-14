package admin

import (
	"literary-lions-forum/pkg/database"
)

func IsUserAdmin(userID int) (bool, error) {
	var isAdmin bool
	err := database.DB.QueryRow("SELECT is_admin FROM users WHERE id = ?", userID).Scan(&isAdmin)
	return isAdmin, err
}

// func CreateAdminUser(username, email, password string) error {
// 	// Hash the password
// 	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
// 	if err != nil {
// 		return fmt.Errorf("error hashing password: %v", err)
// 	}

// 	// Insert the admin user
// 	_, err = database.DB.Exec(`
// 		INSERT INTO users (username, email, password, is_admin)
// 		VALUES (?, ?, ?, TRUE)
// 	`, username, email, string(hashedPassword))

// 	if err != nil {
// 		return fmt.Errorf("error creating admin user: %v", err)
// 	}

// 	fmt.Printf("Admin user '%s' created successfully\n", username)
// 	return nil
// }
