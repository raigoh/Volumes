package handlers

import (
	"fmt"
	"forum/models"
	"net/http"
	"strconv"
	"strings"
)

func LikeCommentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 4 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}
	commentID, err := strconv.Atoi(parts[2])
	if err != nil {
		http.Error(w, "Invalid comment ID", http.StatusBadRequest)
		return
	}
	err = models.LikeComment(commentID)
	if err != nil {
		http.Error(w, "Failed to like comment", http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "Liked comment %d", commentID)
}
