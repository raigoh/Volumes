package handlers

import (
	"encoding/json"
	"fmt"
	"forum/models"
	"net/http"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// TODO: Authenticate user
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "User %s logged in successfully", user.Username)
}
