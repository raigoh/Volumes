package main

import (
	"fmt"
	admin "literary-lions-forum/internal/Admin"
	"literary-lions-forum/internal/auth"
	"literary-lions-forum/internal/category"
	"literary-lions-forum/internal/comment"
	"literary-lions-forum/internal/models"
	"literary-lions-forum/internal/post"
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
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	})
	http.HandleFunc("/home", auth.RequireAuth(auth.HomeHandler))
	http.HandleFunc("/register", auth.RegisterHandler)
	http.HandleFunc("/login", auth.LoginHandler)
	http.HandleFunc("/logout", auth.LogoutHandler)
	http.HandleFunc("/new-post", auth.RequireAuth(post.NewPostHandler))
	http.HandleFunc("/post/", auth.RequireAuth(post.PostDetailHandler))
	http.HandleFunc("/comment", auth.RequireAuth(comment.AddCommentHandler))
	http.HandleFunc("/admin/users", auth.RequireAuth(auth.AdminOnly(admin.UserManagementHandler)))
	http.HandleFunc("/admin/dashboard", auth.RequireAuth(auth.AdminOnly(adminDashboardHandler)))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Println("Server is running on http://localhost:" + port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}

func adminDashboardHandler(w http.ResponseWriter, r *http.Request) {
	utils.RenderTemplate(w, "admin-dashboard.html", models.PageData{
		Title: "Admin Dashboard - Literary Lions Forum",
		Page:  "admin-dashboard",
	})
}
