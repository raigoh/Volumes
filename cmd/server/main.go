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
	if err := utils.InitTemplates(); err != nil {
		log.Fatal("Failed to initialize templates:", err)
	}

	// Serve static files
	fs := http.FileServer(http.Dir("web/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Initialize database
	err := database.InitDB()
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer database.DB.Close()

	// Insert initial categories
	err = category.InsertInitialCategories()
	if err != nil {
		log.Printf("Failed to insert initial categories: %v", err)
	}

	// Set up routes
	http.HandleFunc("/", home.HomeHandler)
	http.HandleFunc("/register", auth.RegisterHandler)
	http.HandleFunc("/login", auth.LoginHandler)
	http.HandleFunc("/logout", auth.LogoutHandler)
	http.HandleFunc("/new-post", auth.RequireAuth(post.NewPostHandler))
	http.HandleFunc("/post/", post.PostDetailHandler)
	http.HandleFunc("/comment", auth.RequireAuth(comment.AddCommentHandler))
	http.HandleFunc("/admin/users", auth.RequireAuth(auth.AdminOnly(admin.UserManagementHandler)))
	http.HandleFunc("/admin/dashboard", auth.RequireAuth(auth.AdminOnly(admin.AdminDashboardHandler)))
	http.HandleFunc("/like", like.LikeHandler)
	http.HandleFunc("/unlike", like.UnLikeHandler)
	http.HandleFunc("/user/{id}", user.UserProfileHandler)

	// This handler is only for testing purposes
	http.HandleFunc("/error", errors.ErrorHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Println("Server is running on http://localhost:" + port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
