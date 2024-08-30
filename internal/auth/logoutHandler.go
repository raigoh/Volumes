package auth

import (
	"literary-lions-forum/pkg/session"
	"net/http"
)

// LogoutHandler handles the user logout process.
// It destroys the user's session and redirects them to the home page.
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	// Destroy the session associated with the current request
	// This typically involves:
	// 1. Invalidating the session ID
	// 2. Removing any server-side session data
	// 3. Instructing the client to remove or expire the session cookie
	session.DestroySession(w, r)

	// Redirect the user to the home page after successful logout
	// http.StatusSeeOther (303) is used for redirecting after POST requests
	// It instructs the client to make a GET request to the specified location
	http.Redirect(w, r, "/home", http.StatusSeeOther)
}
