package home

import (
	"literary-lions-forum/internal/category"
	"literary-lions-forum/internal/models"
	"literary-lions-forum/internal/post"
	"literary-lions-forum/internal/utils"
	"literary-lions-forum/pkg/database"
	"literary-lions-forum/pkg/session"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func SearchHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")
	log.Printf("Searching for: %s", query)

	sess, _ := session.GetSession(w, r)
	userID := session.GetUserID(sess)

	var user *models.User
	if userID != 0 {
		var err error
		user, err = session.GetUserByID(userID)
		if err != nil {
			log.Printf("Error fetching user: %v", err)
			utils.RenderErrorTemplate(w, err, http.StatusInternalServerError, "Server error, the database is being picky today.. We are still training him")
		}
	}

	posts, err := EnhancedSearch(query, 10)
	if err != nil {
		log.Printf("Error searching: %v", err)
		posts = []models.Post{}
	}
	log.Printf("Found %d posts", len(posts))

	pageData := models.PageData{
		Title: "Search Results - Literary Lions Forum",
		Page:  "search",
		User:  user,
		Data: map[string]interface{}{
			"Posts":       posts,
			"SearchQuery": query,
		},
	}

	utils.RenderTemplate(w, "search.html", pageData)
}

func EnhancedSearch(query string, limit int) ([]models.Post, error) {
	query = strings.TrimSpace(strings.ToLower(query))

	// First, check if the query matches a username
	matchedUser, err := GetUserByUsername(query)
	if err == nil && matchedUser.ID != "" {
		// Convert the user ID from string to int
		userID, err := strconv.Atoi(matchedUser.ID)
		if err == nil {
			// If a user is found and ID conversion is successful, return their posts
			return post.GetFilteredPosts(0, userID, false, limit)
		}
	}

	// Next, check if the query matches a category name
	categories, err := category.GetCategories()
	if err == nil {
		for _, cat := range categories {
			if strings.ToLower(cat.Name) == query {
				// If a category is found, return posts from that category
				return post.GetFilteredPosts(cat.ID, 0, false, limit)
			}
		}
	}

	// If no user or category match, perform a general search
	return post.SearchPosts(query, limit)
}

// GetUserByUsername retrieves a user from the database by their username
func GetUserByUsername(username string) (models.User, error) {
	var user models.User
	err := database.DB.QueryRow("SELECT id, username, email FROM users WHERE LOWER(username) = LOWER(?)", username).
		Scan(&user.ID, &user.Username, &user.Email)
	if err != nil {
		return user, err
	}
	return user, nil
}
