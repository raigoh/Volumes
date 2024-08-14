package auth

import (
	"literary-lions-forum/internal/category"
	"literary-lions-forum/internal/models"
	"literary-lions-forum/internal/post"
	"literary-lions-forum/internal/utils"
	"literary-lions-forum/pkg/database"
	"literary-lions-forum/pkg/session"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		utils.RenderTemplate(w, "register.html", models.PageData{
			Title: "Register - Literary Lions Forum",
			Page:  "register",
		})
		return
	}

	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		email := r.FormValue("email")
		password := r.FormValue("password")

		if username == "" || email == "" || password == "" {
			utils.RenderTemplate(w, "register.html", models.PageData{
				Title: "Register - Literary Lions Forum",
				Page:  "register",
				Error: "All fields are required",
			})
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("Error hashing password: %v", err)
			utils.RenderTemplate(w, "register.html", models.PageData{
				Title: "Register - Literary Lions Forum",
				Page:  "register",
				Error: "Error creating user",
			})
			return
		}

		_, err = database.DB.Exec("INSERT INTO users (username, email, password) VALUES (?, ?, ?)", username, email, hashedPassword)
		if err != nil {
			log.Printf("Error inserting user into database: %v", err)
			utils.RenderTemplate(w, "register.html", models.PageData{
				Title: "Register - Literary Lions Forum",
				Page:  "register",
				Error: "Error creating user: " + err.Error(),
			})
			return
		}

		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}

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
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	sess, err := session.GetSession(w, r)
	if err != nil || session.GetUserID(sess) == 0 {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	latestPosts, err := post.GetLatestPosts(5)
	if err != nil {
		log.Printf("Error fetching latest posts: %v", err)
		latestPosts = []models.Post{}
	}

	popularCategories, err := category.GetPopularCategories(5)
	if err != nil {
		log.Printf("Error fetching popular categories: %v", err)
		popularCategories = []models.Category{}
	}

	pageData := models.PageData{
		Title: "Home - Literary Lions Forum",
		Page:  "home",
		Data: map[string]interface{}{
			"LatestPosts":       latestPosts,
			"PopularCategories": popularCategories,
		},
	}

	utils.RenderTemplate(w, "home.html", pageData)
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	// Destroy the session using the custom session package
	session.DestroySession(w, r)
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
