package auth

import (
	"literary-lions-forum/internal/models"
	"literary-lions-forum/internal/utils"
	"literary-lions-forum/pkg/database"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

// RegisterHandler handles both GET and POST requests for user registration.
// GET: Renders the registration form.
// POST: Processes the registration attempt.
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	// Simulate an intentional error to test error handling
	// panic("intentional error for testing")

	// Handle GET request
	if r.Method == http.MethodGet {
		// Render the registration page template
		utils.RenderTemplate(w, "register.html", models.PageData{
			Title: "Register - Literary Lions Forum",
			Page:  "register",
		})
		return
	}

	// Handle POST request
	if r.Method == http.MethodPost {
		// Extract form values
		username := r.FormValue("username")
		email := r.FormValue("email")
		password := r.FormValue("password")

		// Validate that all required fields are provided
		if username == "" || email == "" || password == "" {
			utils.RenderTemplate(w, "register.html", models.PageData{
				Title: "Register - Literary Lions Forum",
				Page:  "register",
				Error: "All fields are required",
			})
			return
		}

		// Hash the password using bcrypt
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			// Log the error and render the registration page with an error message
			log.Printf("Error hashing password: %v", err)
			utils.RenderTemplate(w, "register.html", models.PageData{
				Title: "Register - Literary Lions Forum",
				Page:  "register",
				Error: "Error creating user",
			})
			return
		}

		// Insert the new user into the database
		_, err = database.DB.Exec("INSERT INTO users (username, email, password) VALUES (?, ?, ?)", username, email, hashedPassword)
		if err != nil {
			// Log the error and render the registration page with an error message
			log.Printf("Error inserting user into database: %v", err)
			utils.RenderTemplate(w, "register.html", models.PageData{
				Title: "Register - Literary Lions Forum",
				Page:  "register",
				Error: "Error creating user: " + err.Error(),
			})
			return
		}

		// Redirect to the login page after successful registration
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}
