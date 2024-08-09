package handlers

import (
	"fmt"
	"forum/models"
	"net/http"
	"strconv"
	"strings"
)

func DislikePostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 4 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}
	postID, err := strconv.Atoi(parts[2])
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}
	err = models.DislikePost(postID)
	if err != nil {
		http.Error(w, "Failed to dislike post", http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "Disliked post %d", postID)
}
