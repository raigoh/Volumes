package session

import (
	"crypto/rand"
	"encoding/base64"
	"literary-lions-forum/internal/models"
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

	// Extend session expiration
	session.ExpiresAt = time.Now().Add(sessionDuration)
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
}

func GetUserID(session *models.Session) int {
	return session.UserID
}
