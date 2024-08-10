package main

import (
	"fmt"
	"literary-lions-forum/internal/auth"
	"literary-lions-forum/internal/models"
	"literary-lions-forum/internal/utils"
	"literary-lions-forum/pkg/database"
	"log"
	"net/http"
	"os"
)

// Handlers for different routes
func homeHandler(w http.ResponseWriter, r *http.Request) {
	// Redirect to login page
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	utils.RenderTemplate(w, "base.html", models.PageData{
		Title: "Login - Literary Lions Forum",
		Page:  "login",
	})
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	utils.RenderTemplate(w, "base.html", models.PageData{
		Title: "Register - Literary Lions Forum",
		Page:  "register",
	})
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

	// Set up your routes
	http.HandleFunc("/", homeHandler) // Handle root path
	http.HandleFunc("/register", auth.RegisterHandler)
	http.HandleFunc("/login", auth.LoginHandler)
	http.HandleFunc("/logout", auth.LogoutHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Println("Server is running on http://localhost:" + port)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
