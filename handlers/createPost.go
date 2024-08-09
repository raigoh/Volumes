package handlers

import (
	"database/sql"
	"encoding/json"
	"forum/models"
	"net/http"
	"time"
)

func CreatePostHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var post models.Post
		err := json.NewDecoder(r.Body).Decode(&post)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Set the current timestamp
		post.Timestamp = time.Now()

		// Call CreatePost with the database connection and the post object
		err = models.CreatePost(db, &post)
		if err != nil {
			http.Error(w, "Failed to create post", http.StatusInternalServerError)
			return
		}

		// Respond with the created post
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(post)
	}
}
