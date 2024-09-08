package home

import (
	"literary-lions-forum/internal/category"
	"literary-lions-forum/internal/like"
	"literary-lions-forum/internal/models"
	"literary-lions-forum/internal/post"
	"literary-lions-forum/internal/utils"
	"literary-lions-forum/pkg/database"
	"literary-lions-forum/pkg/session"
	"log"
	"net/http"
	"strconv"
)

// HomeHandler processes requests for the main page of the forum.
// It fetches and prepares all necessary data for rendering the home page.
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	// Set headers to prevent caching of the home page
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")

	// Get the user's session and user ID
	sess, _ := session.GetSession(w, r)
	userID := session.GetUserID(sess)

	// Fetch user data if logged in
	var user *models.User
	if userID != 0 {
		var err error
		user, err = session.GetUserByID(userID)
		if err != nil {
			log.Printf("Error fetching user: %v", err)
			utils.RenderErrorTemplate(w, err, http.StatusInternalServerError, "Server error, the database is being picky today.. We are still training him")
		}
	}

	// Initialize variables for filtering and search
	var categoryID int
	var filterUserID int
	var likedOnly bool
	var searchQuery string

	if r.Method == http.MethodPost {
		r.ParseForm()

		// Handle both search and filter options
		searchQuery = r.FormValue("query")
		categoryID, _ = strconv.Atoi(r.FormValue("category"))
		filterUserID, _ = strconv.Atoi(r.FormValue("user"))
		likedOnly = r.FormValue("liked") == "true"
	} else {
		// Handle GET parameters (you might want to remove this eventually)
		categoryID, _ = strconv.Atoi(r.URL.Query().Get("category"))
		filterUserID, _ = strconv.Atoi(r.URL.Query().Get("user"))
		likedOnly, _ = strconv.ParseBool(r.URL.Query().Get("liked"))
		searchQuery = r.URL.Query().Get("query")
	}

	// Only apply the likedOnly filter if the user is logged in
	if userID == 0 {
		likedOnly = false
	}

	var posts []models.Post
	var err error

	// Fetch posts based on search query or filters
	if searchQuery != "" {
		// If there's a search query, use it to search posts
		posts, err = EnhancedSearch(searchQuery, 10) // Limit to 10 results
	} else {
		// If no search query, use the existing filter logic
		posts, err = post.GetFilteredPosts(categoryID, filterUserID, likedOnly, 10)
	}

	if err != nil {
		log.Printf("Error fetching posts: %v", err)
		posts = []models.Post{}
	}

	// Fetch fresh like/dislike counts for each post
	for i, p := range posts {
		likes, dislikes, err := like.GetLikesCount(p.ID, "post")
		if err != nil {
			log.Printf("Error fetching like counts for post %d: %v", p.ID, err)
			utils.RenderErrorTemplate(w, err, http.StatusInternalServerError, "Server error, who would like this even?!?")
		} else {
			posts[i].Likes = likes
			posts[i].Dislikes = dislikes
		}
	}

	// Fetch popular categories
	popularCategories, err := category.GetPopularCategories(5)
	if err != nil {
		log.Printf("Error fetching popular categories: %v", err)
		popularCategories = []models.Category{}
	}

	// Fetch all categories
	allCategories, err := category.GetCategories()
	if err != nil {
		log.Printf("Error fetching all categories: %v", err)
		allCategories = []models.Category{}
	}

	// Fetch all users
	allUsers, err := database.GetAllUsers()
	if err != nil {
		log.Printf("Error fetching all users: %v", err)
		allUsers = []models.User{}
	}

	// Prepare page data for rendering
	pageData := models.PageData{
		Title: "Home - Literary Lions Forum",
		Page:  "home",
		User:  user,
		Data: map[string]interface{}{
			"Posts":             posts,
			"PopularCategories": popularCategories,
			"AllCategories":     allCategories,
			"SelectedCategory":  categoryID,
			"FilterUserID":      filterUserID,
			"LikedOnly":         likedOnly,
			"CurrentUserID":     userID,
			"SearchQuery":       searchQuery,
		},
		Users: allUsers,
	}

	// Render the home page template
	utils.RenderTemplate(w, "home.html", pageData)
}
