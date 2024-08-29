package auth

import (
	"literary-lions-forum/pkg/session"
	"net/http"
)

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	// Destroy the session using the custom session package
	session.DestroySession(w, r)
	http.Redirect(w, r, "/home", http.StatusSeeOther)
}
