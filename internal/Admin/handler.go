package admin

import (
	"literary-lions-forum/internal/models"
	"literary-lions-forum/internal/utils"
	"literary-lions-forum/pkg/database"
	"log"
	"net/http"
)

// AdminDashboardHandler handles requests to the admin dashboard.
// It fetches various statistics and recent activity data to display on the dashboard.
func AdminDashboardHandler(w http.ResponseWriter, r *http.Request) {
	// Fetch total number of users
	totalUsers, err := database.GetTotalUsers()
	if err != nil {
		log.Printf("Error getting total users: %v", err)
		totalUsers = 0 // Set to 0 if there's an error
	}

	// Fetch total number of posts
	totalPosts, err := database.GetTotalPosts()
	if err != nil {
		log.Printf("Error getting total posts: %v", err)
		totalPosts = 0
	}

	// Fetch total number of comments
	totalComments, err := database.GetTotalComments()
	if err != nil {
		log.Printf("Error getting total comments: %v", err)
		totalComments = 0
	}

	// Fetch number of active users in the last 30 days
	activeUsers, err := database.GetActiveUsers(30)
	if err != nil {
		log.Printf("Error getting active users: %v", err)
		activeUsers = 0
	}

	// Fetch the last 10 activities
	recentActivity, err := database.GetRecentActivity(10)
	if err != nil {
		log.Printf("Error getting recent activity: %v", err)
		recentActivity = []models.Activity{} // Empty slice if there's an error
	}

	// Prepare data for the template
	data := models.PageData{
		Title:          "Admin Dashboard - Literary Lions Forum",
		Page:           "admin-dashboard",
		TotalUsers:     totalUsers,
		TotalPosts:     totalPosts,
		TotalComments:  totalComments,
		ActiveUsers:    activeUsers,
		RecentActivity: recentActivity,
	}

	// Render the admin dashboard template
	utils.RenderTemplate(w, "admin-dashboard.html", data)
}

// UserManagementHandler handles requests to the user management page.
// It supports both GET (display users) and POST (delete user) methods.
func UserManagementHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		// Handle GET request: display user management page

		// Fetch all users from the database
		users, err := database.GetAllUsers()
		if err != nil {
			log.Printf("Error getting users: %v", err)
			//http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			utils.RenderErrorTemplate(w, err, http.StatusInternalServerError, "Internal server error")
			return
		}

		// Prepare data for the template
		data := models.PageData{
			Title: "User Management - Literary Lions Forum",
			Page:  "user_management",
			Users: users,
		}

		// Render the user management template
		utils.RenderTemplate(w, "user-management.html", data)

	} else if r.Method == "POST" {
		// Handle POST request: process user deletion

		// Parse form data
		err := r.ParseForm()
		if err != nil {
			log.Printf("Error parsing form: %v", err)
			//http.Error(w, "Bad Request", http.StatusBadRequest)
			utils.RenderErrorTemplate(w, err, http.StatusBadRequest, "Error parsing form.. ")
			//return
		}

		// Extract user ID and action from form data
		userID := r.FormValue("user_id")
		action := r.FormValue("action")

		// Process the delete action
		if action == "delete" {
			err := database.DeleteUser(userID)
			if err != nil {
				log.Printf("Error deleting user: %v", err)
				//http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				utils.RenderErrorTemplate(w, err, http.StatusInternalServerError, "Internal server error")
				//return
			}
		}

		// Redirect back to the user management page after processing
		http.Redirect(w, r, "/admin/users", http.StatusSeeOther)

	} else {
		// Handle unsupported HTTP methods
		//http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		utils.RenderErrorTemplate(w, nil, http.StatusMethodNotAllowed, "Method not allowed")
	}
}
