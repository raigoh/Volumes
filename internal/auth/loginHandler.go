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

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		utils.RenderTemplate(w, "login.html", models.PageData{
			Title: "Login - Literary Lions Forum",
			Page:  "login",
		})
		return
	}

	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			log.Printf("Error parsing form: %v", err)
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		email := r.FormValue("email")
		password := r.FormValue("password")

		log.Printf("Attempting login for email: %s", email)

		if email == "" || password == "" {
			log.Println("Email or password is empty")
			utils.RenderTemplate(w, "login.html", models.PageData{
				Title: "Login - Literary Lions Forum",
				Page:  "login",
				Error: "Email and password are required",
			})
			return
		}

		var dbPassword string
		var userID int
		var isAdmin bool
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

		// Create session
		sess, err := session.GetSession(w, r)
		if err != nil {
			log.Printf("Error creating session: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		session.SetUserID(sess, userID)
		session.SetIsAdmin(sess, isAdmin)

		// After successful login
		log.Printf("Login successful for user ID: %d, Admin: %v", userID, isAdmin)

		if isAdmin {
			http.Redirect(w, r, "/admin/dashboard", http.StatusSeeOther)
		} else {
			http.Redirect(w, r, "/home", http.StatusSeeOther)
		}
	}
}
