package post

import (
	"literary-lions-forum/internal/category"
	"literary-lions-forum/internal/like"
	"literary-lions-forum/internal/models"
	"literary-lions-forum/internal/utils"
	"literary-lions-forum/pkg/database"
	"literary-lions-forum/pkg/session"
	"log"
	"net/http"
	"strconv"
)

// AllPostsHandler handles the request to display all posts on the forum.
// It supports filtering posts by category, user, and showing only liked posts.
// Additionally, it includes search functionality to find posts by a query.
func AllPostsHandler(w http.ResponseWriter, r *http.Request) {
	// Simulate an intentional error to test error handling
	// panic("intentional error for testing")

	// Retrieve the current session and user ID (if the user is logged in).
	sess, _ := session.GetSession(w, r)
	userID := session.GetUserID(sess)

	// Initialize a variable to hold the user data (if logged in).
	var user *models.User
	if userID != 0 {
		// Fetch the user details by their ID if the user is logged in.
		var err error
		user, err = session.GetUserByID(userID)
		if err != nil {
			// Render an error template if fetching user data fails.
			utils.RenderErrorTemplate(w, err, http.StatusInternalServerError, "Error fetching user data")
			return
		}
	}

	// Variables to hold filter and search criteria.
	var categoryID int     // To filter posts by category.
	var filterUserID int   // To filter posts by a specific user.
	var likedOnly bool     // To filter only liked posts.
	var searchQuery string // Search query for finding posts.

	// Handle POST method to apply filters and search criteria from form submission.
	if r.Method == http.MethodPost {
		// Parse the form data.
		r.ParseForm()
		searchQuery = r.FormValue("query")                    // Retrieve the search query from the form.
		categoryID, _ = strconv.Atoi(r.FormValue("category")) // Convert category filter to integer.
		filterUserID, _ = strconv.Atoi(r.FormValue("user"))   // Convert user filter to integer.
		likedOnly = r.FormValue("liked") == "true"            // Determine if "liked only" is checked.
	} else {
		// Handle GET method, retrieving filters from query parameters in the URL.
		categoryID, _ = strconv.Atoi(r.URL.Query().Get("category"))  // Get category from URL.
		filterUserID, _ = strconv.Atoi(r.URL.Query().Get("user"))    // Get user filter from URL.
		likedOnly, _ = strconv.ParseBool(r.URL.Query().Get("liked")) // Get liked-only flag from URL.
		searchQuery = r.URL.Query().Get("query")                     // Get search query from URL.
	}

	// If the user is not logged in, disable the "liked only" filter.
	if userID == 0 {
		likedOnly = false
	}

	// Variable to hold the list of posts and an error if any.
	var posts []models.Post
	var err error

	// If there's a search query, perform a search for posts matching the query.
	if searchQuery != "" {
		posts, err = EnhancedSearch(searchQuery, 10) // Limit results to 10 posts.
	} else {
		// Otherwise, get posts filtered by category, user, and liked status.
		posts, err = GetFilteredPosts(categoryID, filterUserID, likedOnly, 10) // Limit results to 10 posts.
	}

	// Handle any errors that occur while fetching posts.
	if err != nil {
		log.Printf("Error fetching posts: %v", err) // Log the error.
		posts = []models.Post{}                     // If error occurs, set posts to an empty list.
	}

	// For each post, retrieve the count of likes and dislikes.
	for i, p := range posts {
		likes, dislikes, err := like.GetLikesCount(p.ID, "post") // Get like/dislike counts for the post.
		if err != nil {
			log.Printf("Error fetching like counts for post %d: %v", p.ID, err) // Log any errors.
			// Render an error template if there's an issue fetching like counts.
			utils.RenderErrorTemplate(w, err, http.StatusInternalServerError, "Server error, who would like this even?!?")
		} else {
			// Update the post with the retrieved like and dislike counts.
			posts[i].Likes = likes
			posts[i].Dislikes = dislikes
		}
	}

	// Fetch all available categories for filtering options in the UI.
	allCategories, err := category.GetCategories()
	if err != nil {
		log.Printf("Error fetching all categories: %v", err) // Log error if fetching categories fails.
		allCategories = []models.Category{}                  // Set to empty list on failure.
	}

	// Fetch all users for potential filtering or display purposes.
	allUsers, err := database.GetAllUsers()
	if err != nil {
		log.Printf("Error fetching all users: %v", err) // Log error if fetching users fails.
		allUsers = []models.User{}                      // Set to empty list on failure.
	}

	// Prepare data to pass to the template for rendering.
	data := models.PageData{
		Title: "All Posts - Literary Lions Forum", // Page title.
		Page:  "all-posts",                        // Current page identifier.
		User:  user,                               // User information (if logged in).
		Data: map[string]interface{}{ // Data to pass to the template.
			"Posts":            posts,         // List of posts to display.
			"AllCategories":    allCategories, // List of all categories for filtering.
			"SelectedCategory": categoryID,    // Selected category filter.
			"FilterUserID":     filterUserID,  // Filtered user ID (if any).
			"LikedOnly":        likedOnly,     // Whether to show only liked posts.
			"SearchQuery":      searchQuery,   // The search query (if any).
		},
		Users: allUsers, // List of all users.
	}

	// Render the "all-posts.html" template with the prepared data.
	utils.RenderTemplate(w, "all-posts.html", data)
}
