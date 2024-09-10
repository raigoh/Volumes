package auth

import (
	admin "literary-lions-forum/internal/Admin"
	"literary-lions-forum/internal/utils"
	"literary-lions-forum/pkg/session"
	"net/http"
)

// RequireAuth is a middleware that ensures a user is authenticated before accessing a route.
// If the user is not authenticated, they are redirected to the login page.
func RequireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Attempt to get the user's session
		sess, err := session.GetSession(w, r)
		// If there's an error retrieving the session or the user ID is 0 (indicating no user),
		// redirect to the login page
		if err != nil || session.GetUserID(sess) == 0 {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		// If the user is authenticated, proceed to the next handler
		next.ServeHTTP(w, r)
	}
}

// RequireAdmin is a middleware that ensures a user is both authenticated and an admin before accessing a route.
// If the user is not authenticated or not an admin, they are denied access.
func RequireAdmin(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Attempt to get the user's session
		sess, err := session.GetSession(w, r)
		// If there's an error retrieving the session or the user ID is 0 (indicating no user),
		// redirect to the login page
		if err != nil || session.GetUserID(sess) == 0 {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		// Check if the user is an admin
		isAdmin, err := admin.IsUserAdmin(session.GetUserID(sess))
		// If there's an error checking admin status or the user is not an admin,
		// return a Forbidden error
		if err != nil || !isAdmin {
			//http.Error(w, "Unauthorized", http.StatusForbidden)
			utils.RenderErrorTemplate(w, err, http.StatusUnauthorized, "You shuold not be here, unauthorized access")
			return
		}
		// If the user is authenticated and an admin, proceed to the next handler
		next.ServeHTTP(w, r)
	}
}

// AdminOnly is a middleware that ensures a user is an admin before accessing a route.
// It assumes the user is already authenticated and checks only for admin status.
func AdminOnly(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Attempt to get the user's session
		sess, err := session.GetSession(w, r)
		// If there's an error retrieving the session, return an Unauthorized error
		if err != nil {
			//http.Error(w, "Unauthorized", http.StatusUnauthorized)
			utils.RenderErrorTemplate(w, err, http.StatusUnauthorized, "You shuold not be here, unauthorized access")
			return
		}
		// Check if the user is an admin using the session data
		if !session.GetIsAdmin(sess) {
			//http.Error(w, "Unauthorized", http.StatusUnauthorized)
			utils.RenderErrorTemplate(w, err, http.StatusUnauthorized, "You shuold not be here, unauthorized access")
			return
		}
		// If the user is an admin, proceed to the next handler
		next.ServeHTTP(w, r)
	}
}
