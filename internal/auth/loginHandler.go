package auth

import (
	"literary-lions-forum/internal/models"
	"literary-lions-forum/internal/utils"
	"literary-lions-forum/pkg/database"
	"literary-lions-forum/pkg/session"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

// LoginHandler handles both GET and POST requests for the login process.
// GET: Renders the login page.
// POST: Processes the login attempt.
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	// Handle GET request
	if r.Method == http.MethodGet {
		// Render the login page template
		utils.RenderTemplate(w, "login.html", models.PageData{
			Title: "Login - Literary Lions Forum",
			Page:  "login",
		})
		return
	}

	// Handle POST request
	if r.Method == http.MethodPost {
		// Parse the form data
		err := r.ParseForm()
		if err != nil {
			log.Printf("Error parsing form: %v", err)
			http.Error(w, "Bad Request", http.StatusBadRequest)
			utils.RenderErrorTemplate(w, err, http.StatusBadRequest, "Something went wrong")
			return
		}

		// Extract email and password from the form
		email := r.FormValue("email")
		password := r.FormValue("password")

		log.Printf("Attempting login for email: %s", email)

		// Check if email or password is empty
		if email == "" || password == "" {
			log.Println("Email or password is empty")
			utils.RenderTemplate(w, "login.html", models.PageData{
				Title: "Login - Literary Lions Forum",
				Page:  "login",
				Error: "Email and password are required",
			})
			return
		}

		// Variables to store user data from the database
		var dbPassword string
		var userID int
		var isAdmin bool

		// Query the database for user information
		err = database.DB.QueryRow("SELECT id, password, is_admin FROM users WHERE email = ?", email).Scan(&userID, &dbPassword, &isAdmin)
		if err != nil {
			log.Printf("Error querying user: %v", err)
			utils.RenderTemplate(w, "login.html", models.PageData{
				Title: "Login - Literary Lions Forum",
				Page:  "login",
				Error: "Invalid email or password",
			})
			return
		}

		// Compare the stored hashed password with the provided password
		err = bcrypt.CompareHashAndPassword([]byte(dbPassword), []byte(password))
		if err != nil {
			log.Printf("Password mismatch for user %d: %v", userID, err)
			utils.RenderTemplate(w, "login.html", models.PageData{
				Title: "Login - Literary Lions Forum",
				Page:  "login",
				Error: "Invalid email or password",
			})
			return
		}

		// Create a new session
		sess, err := session.GetSession(w, r)
		if err != nil {
			log.Printf("Error creating session: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			utils.RenderErrorTemplate(w, err, http.StatusInternalServerError, "Oh no, cookie can is empty :(")
			return
		}

		// Set user ID and admin status in the session
		session.SetUserID(sess, userID)
		session.SetIsAdmin(sess, isAdmin)

		// Log successful login
		log.Printf("Login successful for user ID: %d, Admin: %v", userID, isAdmin)

		// Redirect based on user role
		if isAdmin {
			http.Redirect(w, r, "/admin/dashboard", http.StatusSeeOther)
		} else {
			http.Redirect(w, r, "/home", http.StatusSeeOther)
		}
	}
}
