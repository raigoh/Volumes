package main

import (
	"literary-lions-forum/internal/models"
	"literary-lions-forum/internal/utils"
	"log"
	"net/http"
)

// Handlers for different routes
func homeHandler(w http.ResponseWriter, r *http.Request) {
	utils.RenderTemplate(w, "home.html", models.PageData{Title: "Home - Literary Lions Forum"})
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	utils.RenderTemplate(w, "login.html", models.PageData{Title: "Login - Literary Lions Forum"})
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	utils.RenderTemplate(w, "register.html", models.PageData{Title: "Register - Literary Lions Forum"})
}

func main() {
	// Initialize templates
	if err := utils.InitTemplates(); err != nil {
		log.Fatal("Failed to initialize templates:", err)
	}

	// Serve static files
	fs := http.FileServer(http.Dir("web/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Set up your routes
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/register", registerHandler)

	log.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
