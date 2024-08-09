package main

import (
	"fmt"
	"forum/database"
	"forum/handlers"
	"forum/models"
	"html/template"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

var templates *template.Template

func init() {
	templates = template.Must(template.ParseGlob("templates/*.html"))
}

func main() {
	// Serve static files
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	log.Println("Starting the application...")

	err := models.InitDB("./literary_lions.db")
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	log.Println("Database initialized successfully")

	// Initialize the database tables
	err = database.InitializeTables(models.GetDB())
	if err != nil {
		log.Fatalf("Failed to initialize database tables: %v", err)
	}
	log.Println("Database tables initialized successfully")

	db := models.GetDB() // Get the DB instance

	http.HandleFunc("/", handlers.HomeHandler(db))
	http.HandleFunc("/create-post", handlers.CreatePostHandler(db))
	http.HandleFunc("/login", handlers.LoginHandler)
	http.HandleFunc("/register", handlers.RegisterHandler)
	http.HandleFunc("/post/", handlers.PostHandler)
	http.HandleFunc("/post/*/comment", handlers.CommentHandler)
	http.HandleFunc("/post/*/like", handlers.LikePostHandler)
	http.HandleFunc("/post/*/dislike", handlers.DislikePostHandler)
	http.HandleFunc("/comment/*/like", handlers.LikeCommentHandler)
	http.HandleFunc("/comment/*/dislike", handlers.DislikeCommentHandler)

	log.Println("Routes registered successfully")

	fmt.Println("Server starting on :8080")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
