package main

import (
	"fmt"
	"literary-lions-forum/internal/auth"
	"literary-lions-forum/internal/category"
	"literary-lions-forum/internal/comment"
	"literary-lions-forum/internal/post"
	"literary-lions-forum/internal/utils"
	"literary-lions-forum/pkg/database"
	"literary-lions-forum/pkg/session"
	"log"
	"net/http"
	"os"
)

func requireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sess, err := session.GetSession(w, r)
		if err != nil || session.GetUserID(sess) == 0 {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	}
}

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

	database.VerifyDatabaseContents()

	// Set up routes
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	})
	http.HandleFunc("/home", requireAuth(auth.HomeHandler))
	http.HandleFunc("/register", auth.RegisterHandler)
	http.HandleFunc("/login", auth.LoginHandler)
	http.HandleFunc("/logout", auth.LogoutHandler)
	http.HandleFunc("/new-post", requireAuth(post.NewPostHandler))
	http.HandleFunc("/post/", requireAuth(post.PostDetailHandler))
	http.HandleFunc("/comment", requireAuth(comment.AddCommentHandler))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Println("Server is running on http://localhost:" + port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
