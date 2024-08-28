package user

import (
	"literary-lions-forum/internal/models"
	"literary-lions-forum/internal/utils"
	"literary-lions-forum/pkg/database"
	"net/http"
	"strings"
)

func UserProfileHandler(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from URL path
	path := strings.TrimPrefix(r.URL.Path, "/user/")
	userID := strings.TrimSuffix(path, "/")

	if userID == "" {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	user, err := database.GetUserByID(userID)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	posts, err := database.GetUserPosts(userID)
	if err != nil {
		http.Error(w, "Error fetching user posts", http.StatusInternalServerError)
		return
	}

	data := struct {
		User  models.User
		Posts []models.Post
	}{
		User:  user,
		Posts: posts,
	}

	utils.RenderTemplate(w, "user-profile.html", data)
}
