package handlers

import (
	"database/sql"
	"fmt"
	"forum/models"
	"net/http"
	"strconv"
	"text/template"
)

var templates *template.Template

// Assuming you have a function to get the user from the session/cookie.
func getUserFromSession(db *sql.DB, r *http.Request) (*models.User, error) {
	// Retrieve the user ID from the session or cookie
	cookie, err := r.Cookie("user_id")
	if err != nil {
		return nil, nil // No user logged in
	}

	userID, err := strconv.Atoi(cookie.Value)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID in cookie: %v", err)
	}

	// Retrieve the user from the database
	user, err := models.GetUserByID(db, userID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func HomeHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the logged-in user, if any
		user, err := getUserFromSession(db, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		data := struct {
			Title string
			User  *models.User
		}{
			Title: "Home",
			User:  user,
		}

		// Render the template with the data
		err = templates.ExecuteTemplate(w, "layout.html", data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
