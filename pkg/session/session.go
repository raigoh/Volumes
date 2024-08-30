package session

import (
	"crypto/rand"
	"encoding/base64"
	"literary-lions-forum/internal/models"
	"literary-lions-forum/pkg/database"
	"net/http"
	"sync"
	"time"
)

// sessionStore is a thread-safe map to store active sessions
var sessionStore = struct {
	sync.RWMutex
	sessions map[string]*models.Session
}{sessions: make(map[string]*models.Session)}

const (
	sessionCookieName = "session_id"
	sessionDuration   = 24 * time.Hour
)

// generateSessionID creates a unique session ID
func generateSessionID() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

// GetSession retrieves or creates a session for the current request
func GetSession(w http.ResponseWriter, r *http.Request) (*models.Session, error) {
	cookie, err := r.Cookie(sessionCookieName)
	if err != nil || cookie.Value == "" {
		return createSession(w)
	}

	sessionStore.RLock()
	session, exists := sessionStore.sessions[cookie.Value]
	sessionStore.RUnlock()

	if !exists || time.Now().After(session.ExpiresAt) {
		if exists {
			DestroySession(w, r)
		}
		return createSession(w)
	}

	return session, nil
}

// createSession generates a new session and sets a cookie
func createSession(w http.ResponseWriter) (*models.Session, error) {
	sessionID := generateSessionID()
	session := &models.Session{
		ID:        sessionID,
		Data:      make(map[string]interface{}),
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(sessionDuration),
	}

	sessionStore.Lock()
	sessionStore.sessions[sessionID] = session
	sessionStore.Unlock()

	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    sessionID,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Expires:  session.ExpiresAt,
	})

	return session, nil
}

// DestroySession removes the session and clears the session cookie
func DestroySession(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(sessionCookieName)
	if err == nil && cookie.Value != "" {
		sessionStore.Lock()
		delete(sessionStore.sessions, cookie.Value)
		sessionStore.Unlock()

		http.SetCookie(w, &http.Cookie{
			Name:     sessionCookieName,
			Value:    "",
			Path:     "/",
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteStrictMode,
			MaxAge:   -1,
		})
	}
}

// SetUserID stores the user ID in the session
func SetUserID(session *models.Session, userID int) {
	session.UserID = userID
}

// GetUserID retrieves the user ID from the session
func GetUserID(session *models.Session) int {
	return session.UserID
}

// SetIsAdmin stores the admin status in the session
func SetIsAdmin(session *models.Session, isAdmin bool) {
	session.Data["isAdmin"] = isAdmin
}

// GetIsAdmin retrieves the admin status from the session
func GetIsAdmin(session *models.Session) bool {
	isAdmin, ok := session.Data["isAdmin"].(bool)
	return ok && isAdmin
}

// GetUserByID retrieves a user from the database by their ID
func GetUserByID(userID int) (*models.User, error) {
	var user models.User
	err := database.DB.QueryRow("SELECT id, username, email FROM users WHERE id = ?", userID).Scan(&user.ID, &user.Username, &user.Email)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
