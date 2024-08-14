package user

import (
	"literary-lions-forum/internal/utils"
	"literary-lions-forum/pkg/database"
	"net/http"
	"strconv"
)

func UserManagementHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		users, err := database.ListUsers()
		if err != nil {
			http.Error(w, "Error fetching users: "+err.Error(), http.StatusInternalServerError)
			return
		}
		utils.RenderTemplate(w, "user-management.html", map[string]interface{}{
			"Users": users,
		})
	} else if r.Method == http.MethodPost {
		action := r.FormValue("action")
		userID, err := strconv.Atoi(r.FormValue("user_id"))
		if err != nil {
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}

		if action == "delete" {
			err = database.DeleteUser(userID)
			if err != nil {
				http.Error(w, "Error deleting user: "+err.Error(), http.StatusInternalServerError)
				return
			}
			http.Redirect(w, r, "/admin/users", http.StatusSeeOther)
		} else {
			http.Error(w, "Invalid action", http.StatusBadRequest)
		}
	}
}
