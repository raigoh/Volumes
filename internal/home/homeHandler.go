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

// HomeHandler is the main handler for rendering the home page of the Literary Lions Forum.
// It fetches data such as user session information, popular categories, all categories, latest posts, and user details.
// This data is then passed to a template for rendering the homepage.
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	// Simulate a panic for testing
	panic(nil)

	// Set HTTP headers to prevent caching of the home page to ensure users see the most updated content.
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")

	// Retrieve the current session of the user.
	sess, _ := session.GetSession(w, r)

	// Fetch the user ID from the session. If the user is not logged in, userID will be 0.
	userID := session.GetUserID(sess)

	var user *models.User
	// Check if the user is logged in by verifying if userID is non-zero.
	if userID != 0 {
		var err error
		// Fetch the user details by userID. This would be used to personalize the homepage.
		user, err = session.GetUserByID(userID)
		if err != nil {
			// Log the error and render an error page if there is an issue fetching the user data.
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
		posts, err = post.EnhancedSearch(searchQuery, 10) // Limit to 10 results
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
			return
		} else {
			posts[i].Likes = likes
			posts[i].Dislikes = dislikes
		}
	}

	// Fetch popular categories
	popularCategories, err := category.GetPopularCategories(5)
	if err != nil {
		// Log the error if there is an issue fetching the popular categories, and set the data to an empty list to avoid breaking the page.
		log.Printf("Error fetching popular categories: %v", err)
		popularCategories = []models.Category{}
	}

	// Fetch all available categories in the forum. This will display the full list of categories.
	allCategories, err := category.GetCategories()
	if err != nil {
		// Log the error and fall back to an empty list of categories in case of a failure.
		log.Printf("Error fetching all categories: %v", err)
		allCategories = []models.Category{}
	}

	// Fetch all users from the database. This could be used to display user-related data on the homepage.
	allUsers, err := database.GetAllUsers()
	if err != nil {
		// Log the error and return an empty list if there’s an issue fetching user data.
		log.Printf("Error fetching all users: %v", err)
		allUsers = []models.User{}
	}

	// Fetch the latest posts, limited to 5. This will show the most recent discussions or posts on the homepage.
	latestPosts, err := post.GetLatestPosts(5)
	if err != nil {
		// Log the error and set the posts list to empty if there is a problem fetching the latest posts.
		log.Printf("Error fetching latest posts: %v", err)
		latestPosts = []models.Post{}
	}

	// For each post in the latestPosts list, fetch the like and dislike counts.
	for i, p := range latestPosts {
		// Retrieve the number of likes and dislikes for the post using the post's ID.
		likes, dislikes, err := like.GetLikesCount(p.ID, "post")
		if err != nil {
			// Log the error if there is a problem retrieving the like and dislike counts for the post.
			log.Printf("Error fetching like counts for post %d: %v", p.ID, err)
		} else {
			// If no error occurs, update the Likes and Dislikes fields for the current post in the slice.
			latestPosts[i].Likes = likes
			latestPosts[i].Dislikes = dislikes
		}
	}

	// Prepare the data to be passed to the HTML template for rendering the homepage.
	pageData := models.PageData{
		Title: "Home - Literary Lions Forum", // Title of the webpage, displayed in the browser tab.
		Page:  "home",                        // Identifier for the page, can be used to load page-specific scripts/styles.
		User:  user,                          // The user object, which is `nil` if the user is not logged in.
		Data: map[string]interface{}{
			"PopularCategories": popularCategories, // Popular categories data for display.
			"AllCategories":     allCategories,     // All available categories for display.
			"CurrentUserID":     userID,            // The current user's ID (0 if not logged in).
			"Posts":             latestPosts,       // Latest posts data for display on the homepage.
		},
		Users: allUsers, // Data for all users, potentially used to display user-related information.
	}

	// Render the "home.html" template, passing the `pageData` that contains all the necessary information.
	utils.RenderTemplate(w, "home.html", pageData)
}
