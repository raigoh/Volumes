package admin

import (
	"literary-lions-forum/internal/models"
	"literary-lions-forum/internal/utils"
	"literary-lions-forum/pkg/database"
	"log"
	"net/http"
)

func adminDashboardHandler(w http.ResponseWriter, r *http.Request) {
	// Fetch the required data from your database
	totalUsers, err := database.GetTotalUsers()
	if err != nil {
		log.Printf("Error getting total users: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	totalPosts, err := database.GetTotalPosts()
	if err != nil {
		log.Printf("Error getting total posts: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	totalComments, err := database.GetTotalComments()
	if err != nil {
		log.Printf("Error getting total comments: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	activeUsers, err := database.GetActiveUsers(30) // Active in last 30 days
	if err != nil {
		log.Printf("Error getting active users: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	recentActivity, err := database.GetRecentActivity(10) // Get last 10 activities
	if err != nil {
		log.Printf("Error getting recent activity: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data := models.PageData{
		Title:          "Admin Dashboard - Literary Lions Forum",
		Page:           "admin_dashboard",
		TotalUsers:     totalUsers,
		TotalPosts:     totalPosts,
		TotalComments:  totalComments,
		ActiveUsers:    activeUsers,
		RecentActivity: recentActivity,
	}

	utils.RenderTemplate(w, "admin-dashboard.html", data)
}

func UserManagementHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		// Fetch users from the database
		users, err := database.GetAllUsers()
		if err != nil {
			log.Printf("Error getting users: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		data := models.PageData{
			Title: "User Management - Literary Lions Forum",
			Page:  "user_management",
			Users: users,
		}

		utils.RenderTemplate(w, "user-management.html", data)
	} else if r.Method == "POST" {
		// Handle user deletion
		err := r.ParseForm()
		if err != nil {
			log.Printf("Error parsing form: %v", err)
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		userID := r.FormValue("user_id")
		action := r.FormValue("action")

		if action == "delete" {
			err := database.DeleteUser(userID)
			if err != nil {
				log.Printf("Error deleting user: %v", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
		}

		// Redirect back to the user management page
		http.Redirect(w, r, "/admin/users", http.StatusSeeOther)
	} else {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}
