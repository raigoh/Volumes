package user

import (
	"literary-lions-forum/internal/models"
	"literary-lions-forum/internal/utils"
	"literary-lions-forum/pkg/database"
	"net/http"
	"strings"
)

// UserProfileHandler handles HTTP requests for displaying a user's profile
func UserProfileHandler(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from URL path
	// This assumes the URL structure is /user/{userID}/
	path := strings.TrimPrefix(r.URL.Path, "/user/")
	userID := strings.TrimSuffix(path, "/")

	// Validate that a user ID was provided
	if userID == "" {
		//http.Error(w, "Invalid user ID", http.StatusBadRequest)
		utils.RenderErrorTemplate(w, nil, http.StatusBadRequest, "Invalid user ID")
		return
	}

	// Fetch the user data from the database
	user, err := database.GetUserByID(userID)
	if err != nil {
		// If the user is not found, return a 400 error
		//http.Error(w, "User not found", http.StatusNotFound)
		utils.RenderErrorTemplate(w, nil, http.StatusBadRequest, "Invalid user ID")
		return
	}

	// Fetch the user's posts from the database
	posts, err := database.GetUserPosts(userID)
	if err != nil {
		// If there's an error fetching posts, return a 500 error
		//http.Error(w, "Error fetching user posts", http.StatusInternalServerError)
		utils.RenderErrorTemplate(w, err, http.StatusInternalServerError, "Error fetching user posts")
		return
	}

	// Prepare the data for the template
	data := struct {
		User  models.User
		Posts []models.Post
	}{
		User:  user,
		Posts: posts,
	}

	// Render the user profile template with the prepared data
	utils.RenderTemplate(w, "user-profile.html", data)
}
