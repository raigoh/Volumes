package session

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"sync"
)

var sessionStore = struct {
	sync.RWMutex
	sessions map[string]map[string]interface{}
}{sessions: make(map[string]map[string]interface{})}

// Generate a random session ID
func generateSessionID() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

// Create or retrieve a session
func GetSession(w http.ResponseWriter, r *http.Request) map[string]interface{} {
	cookie, err := r.Cookie("session_id")
	if err != nil || cookie.Value == "" {
		sessionID := generateSessionID()
		http.SetCookie(w, &http.Cookie{
			Name:  "session_id",
			Value: sessionID,
			Path:  "/",
		})

		sessionStore.Lock()
		sessionStore.sessions[sessionID] = make(map[string]interface{})
		sessionStore.Unlock()

		return sessionStore.sessions[sessionID]
	}

	sessionStore.RLock()
	session, exists := sessionStore.sessions[cookie.Value]
	sessionStore.RUnlock()

	if !exists {
		sessionStore.Lock()
		session = make(map[string]interface{})
		sessionStore.sessions[cookie.Value] = session
		sessionStore.Unlock()
	}

	return session
}

// Destroy a session
func DestroySession(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err == nil && cookie.Value != "" {
		sessionStore.Lock()
		delete(sessionStore.sessions, cookie.Value)
		sessionStore.Unlock()

		http.SetCookie(w, &http.Cookie{
			Name:   "session_id",
			Value:  "",
			Path:   "/",
			MaxAge: -1, // Delete the cookie
		})
	}
}
