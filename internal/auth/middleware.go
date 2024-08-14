package auth

import (
	"literary-lions-forum/internal/user"
	"literary-lions-forum/pkg/session"
	"net/http"
)

func RequireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sess, err := session.GetSession(w, r)
		if err != nil || session.GetUserID(sess) == 0 {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	}
}

func RequireAdmin(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sess, err := session.GetSession(w, r)
		if err != nil || session.GetUserID(sess) == 0 {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		// Check if the user is an admin (you'll need to implement this logic)
		isAdmin, err := user.IsUserAdmin(session.GetUserID(sess))
		if err != nil || !isAdmin {
			http.Error(w, "Unauthorized", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	}
}
