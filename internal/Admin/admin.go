package admin

import (
	"literary-lions-forum/pkg/database"
)

func IsUserAdmin(userID int) (bool, error) {
	var isAdmin bool
	err := database.DB.QueryRow("SELECT is_admin FROM users WHERE id = ?", userID).Scan(&isAdmin)
	return isAdmin, err
}
