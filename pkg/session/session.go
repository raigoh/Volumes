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

var sessionStore = struct {
	sync.RWMutex
	sessions map[string]*models.Session
}{sessions: make(map[string]*models.Session)}

const (
	sessionCookieName = "session_id"
	sessionDuration   = 24 * time.Hour
)

func generateSessionID() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

func GetSession(w http.ResponseWriter, r *http.Request) (*models.Session, error) {
	cookie, err := r.Cookie(sessionCookieName)
	if err != nil || cookie.Value == "" {
		// log.Println("No session cookie found, creating new session")
		return createSession(w)
	}

	sessionStore.RLock()
	session, exists := sessionStore.sessions[cookie.Value]
	sessionStore.RUnlock()

	if !exists || time.Now().After(session.ExpiresAt) {
		// log.Println("Session doesn't exist or has expired, creating new session")
		if exists {
			DestroySession(w, r)
		}
		return createSession(w)
	}

	// log.Printf("Retrieved existing session for user ID: %d", session.UserID)
	return session, nil
}

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

func SetUserID(session *models.Session, userID int) {
	session.UserID = userID
	// log.Printf("Set user ID in session: %d", userID)
}

func GetUserID(session *models.Session) int {
	// log.Printf("Getting user ID from session: %d", session.UserID)
	return session.UserID
}

func SetIsAdmin(session *models.Session, isAdmin bool) {
	session.Data["isAdmin"] = isAdmin
}

func GetIsAdmin(session *models.Session) bool {
	isAdmin, ok := session.Data["isAdmin"].(bool)
	return ok && isAdmin
}

func GetUserByID(userID int) (*models.User, error) {
	var user models.User
	err := database.DB.QueryRow("SELECT id, username, email FROM users WHERE id = ?", userID).Scan(&user.ID, &user.Username, &user.Email)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
