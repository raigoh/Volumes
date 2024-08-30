package main

import (
	"fmt"
	admin "literary-lions-forum/internal/Admin"
	"literary-lions-forum/internal/auth"
	"literary-lions-forum/internal/category"
	"literary-lions-forum/internal/comment"
	"literary-lions-forum/internal/errors"
	"literary-lions-forum/internal/home"
	"literary-lions-forum/internal/like"
	"literary-lions-forum/internal/post"
	"literary-lions-forum/internal/user"
	"literary-lions-forum/internal/utils"
	"literary-lions-forum/pkg/database"
	"log"
	"net/http"
	"os"
)

func main() {
	// Initialize templates
	// This loads and parses all HTML templates used in the application
	if err := utils.InitTemplates(); err != nil {
		log.Fatal("Failed to initialize templates:", err)
	}

	// Serve static files
	// This sets up a file server to serve static files (CSS, JS, images) from the "web/static" directory
	fs := http.FileServer(http.Dir("web/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Initialize database
	// This establishes a connection to the database and sets it up for use
	err := database.InitDB()
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	// Ensure the database connection is closed when the program exits
	defer database.DB.Close()

	// Insert initial categories
	// This populates the database with predefined categories if they don't already exist
	err = category.InsertInitialCategories()
	if err != nil {
		log.Printf("Failed to insert initial categories: %v", err)
	}

	// Set up routes
	// Each of these lines maps a URL path to a specific handler function

	// Home page
	http.HandleFunc("/", home.HomeHandler)

	// Authentication routes
	http.HandleFunc("/register", auth.RegisterHandler)
	http.HandleFunc("/login", auth.LoginHandler)
	http.HandleFunc("/logout", auth.LogoutHandler)

	// Post-related routes
	http.HandleFunc("/new-post", auth.RequireAuth(post.NewPostHandler)) // Requires authentication
	http.HandleFunc("/post/", post.PostDetailHandler)

	// Comment route
	http.HandleFunc("/comment", auth.RequireAuth(comment.AddCommentHandler)) // Requires authentication

	// Admin routes
	http.HandleFunc("/admin/users", auth.RequireAuth(auth.AdminOnly(admin.UserManagementHandler)))     // Requires auth and admin privileges
	http.HandleFunc("/admin/dashboard", auth.RequireAuth(auth.AdminOnly(admin.AdminDashboardHandler))) // Requires auth and admin privileges

	// Like/Unlike routes
	http.HandleFunc("/like", like.LikeHandler)
	http.HandleFunc("/unlike", like.UnLikeHandler)

	// User profile route
	http.HandleFunc("/user/{id}", user.UserProfileHandler)

	// Error handler (for testing purposes)
	http.HandleFunc("/error", errors.ErrorHandler)

	// Determine the port to run the server on
	// Use the PORT environment variable if set, otherwise default to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Print the server's address to the console
	fmt.Println("Server is running on http://localhost:" + port)

	// Start the HTTP server
	// If there's an error starting the server, log it and exit the program
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
