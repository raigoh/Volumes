package handlers

import (
	"database/sql"
	"encoding/json"
	"forum/models"
	"net/http"
	"strconv"
	"strings"
)

// Assuming db is available globally or passed somehow
var db *sql.DB

func PostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 3 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}
	postID, err := strconv.Atoi(parts[2])
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}
	post, err := models.GetPostByID(db, postID) // Pass db to GetPostByID
	if err != nil {
		http.Error(w, "Failed to fetch post", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(post)
}
