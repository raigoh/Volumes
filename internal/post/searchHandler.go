package post

import (
	"literary-lions-forum/internal/category"
	"literary-lions-forum/internal/models"
	"literary-lions-forum/internal/utils"
	"literary-lions-forum/pkg/database"
	"literary-lions-forum/pkg/session"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// SearchHandler handles the search requests from the user.
// It processes the search query and displays the relevant posts based on the query.
func SearchHandler(w http.ResponseWriter, r *http.Request) {
	// Simulate an intentional error to test error handling
	// panic("intentional error for testing")

	// Extract the search query from the URL parameters.
	query := r.URL.Query().Get("query")
	log.Printf("Searching for: %s", query) // Log the search query for debugging.

	// Retrieve the session data for the current user (if any).
	sess, _ := session.GetSession(w, r)
	userID := session.GetUserID(sess) // Get the user ID from the session.

	var user *models.User // Initialize a user variable to store user details if the user is logged in.
	if userID != 0 {      // Check if the user is logged in (userID is not 0).
		var err error
		// Fetch the user's details from the database using the session user ID.
		user, err = session.GetUserByID(userID)
		if err != nil {
			log.Printf("Error fetching user: %v", err) // Log any error in fetching the user.
			// Render an error template with a custom error message if there's an issue fetching the user.
			utils.RenderErrorTemplate(w, err, http.StatusInternalServerError, "Server error, the database is being picky today.. We are still training him")
		}
	}

	// Perform an enhanced search for posts based on the query.
	posts, err := EnhancedSearch(query)
	if err != nil {
		log.Printf("Error searching: %v", err) // Log any search-related errors.
		posts = []models.Post{}                // If there’s an error, set posts to an empty slice to avoid rendering issues.
	}
	log.Printf("Found %d posts", len(posts)) // Log the number of posts found.

	// Create a PageData object to store the data required for rendering the search results page.
	pageData := models.PageData{
		Title: "Search Results - Literary Lions Forum", // Set the page title.
		Page:  "search",                                // Identify the current page as "search".
		User:  user,                                    // Pass the current user's details (if any).
		Data: map[string]interface{}{ // Pass the posts and search query to the template.
			"Posts":       posts,
			"SearchQuery": query,
		},
	}

	// Render the search results template with the provided data.
	utils.RenderTemplate(w, "all-posts.html", pageData)
}

// EnhancedSearch is an improved search function that performs multiple types of searches:
// - Search by username.
// - Search by category name.
// - General post search based on the query.
func EnhancedSearch(query string) ([]models.Post, error) {
	// Trim spaces and convert the query to lowercase for case-insensitive search.
	query = strings.TrimSpace(strings.ToLower(query))

	// First, check if the query matches a username.
	matchedUser, err := GetUserByUsername(query) // Attempt to find a user with the username matching the query.
	if err == nil && matchedUser.ID != "" {      // If a user is found and their ID is valid:
		// Convert the user ID from string to integer.
		userID, err := strconv.Atoi(matchedUser.ID)
		if err == nil {
			// Return posts created by the found user.
			return GetFilteredPosts(0, userID, false)
		}
	}

	// Next, check if the query matches a category name.
	categories, err := category.GetCategories() // Fetch all categories from the database.
	if err == nil {
		for _, cat := range categories { // Loop through all categories.
			if strings.ToLower(cat.Name) == query { // If a category name matches the query:
				// Return posts from the matching category.
				return GetFilteredPosts(cat.ID, 0, false)
			}
		}
	}

	// If no user or category matches, perform a general search on post titles and content.
	return SearchPosts(query)
}

// GetUserByUsername retrieves a user from the database based on their username.
// It performs a case-insensitive search using the `LOWER()` SQL function.
func GetUserByUsername(username string) (models.User, error) {
	var user models.User // Initialize a User object to hold the query result.
	// Query the database to find a user by their username (case-insensitive).
	err := database.DB.QueryRow("SELECT id, username, email FROM users WHERE LOWER(username) = LOWER(?)", username).
		Scan(&user.ID, &user.Username, &user.Email) // Scan the resulting row into the User object.
	if err != nil {
		return user, err // Return an error if the user is not found.
	}
	return user, nil // Return the found user.
}
