package like

import (
	"encoding/json"
	"fmt"
	"literary-lions-forum/pkg/session"
	"log"
	"net/http"
	"strconv"
)

// LikeHandler handles HTTP POST requests for liking or disliking a post or comment.
func LikeHandler(w http.ResponseWriter, r *http.Request) {
	// Ensure the request method is POST
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get the user's session
	sess, err := session.GetSession(w, r)
	if err != nil {
		http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
		return
	}

	// Get the user ID from the session
	userID := session.GetUserID(sess)
	if userID == 0 {
		http.Error(w, "Unauthorized: User ID is 0", http.StatusUnauthorized)
		return
	}

	// Parse and validate the target ID (post or comment ID)
	targetID, err := strconv.Atoi(r.FormValue("target_id"))
	if err != nil {
		http.Error(w, "Invalid target ID", http.StatusBadRequest)
		return
	}

	// Validate the target type (must be either "post" or "comment")
	targetType := r.FormValue("target_type")
	if targetType != "post" && targetType != "comment" {
		http.Error(w, "Invalid target type", http.StatusBadRequest)
		return
	}

	// Parse and validate the is_like boolean
	isLike, err := strconv.ParseBool(r.FormValue("is_like"))
	if err != nil {
		http.Error(w, "Invalid is_like value", http.StatusBadRequest)
		return
	}

	// Process the like/dislike action
	err = AddLike(userID, targetID, targetType, isLike)
	if err != nil {
		log.Printf("Error processing like/dislike: %v", err)
		http.Error(w, fmt.Sprintf("Error processing like/dislike: %v", err), http.StatusInternalServerError)
		return
	}

	// Redirect back to the post detail page with an anchor to the specific comment
	postID := r.FormValue("post_id")
	http.Redirect(w, r, fmt.Sprintf("/post/%s#comment-%d", postID, targetID), http.StatusSeeOther)
}

// UnLikeHandler handles HTTP requests for removing a like or dislike from a post or comment.
func UnLikeHandler(w http.ResponseWriter, r *http.Request) {
	// Get the user's session
	sess, err := session.GetSession(w, r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get the user ID from the session
	userID := session.GetUserID(sess)
	if userID == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Parse the target ID and type
	targetID, _ := strconv.Atoi(r.FormValue("target_id"))
	targetType := r.FormValue("target_type")

	// Remove the like/dislike
	err = RemoveLike(userID, targetID, targetType)
	if err != nil {
		http.Error(w, "Error removing like", http.StatusInternalServerError)
		return
	}

	// Get the updated like and dislike counts
	likes, dislikes, err := GetLikesCount(targetID, targetType)
	if err != nil {
		http.Error(w, "Error getting like count", http.StatusInternalServerError)
		return
	}

	// Send the updated counts as a JSON response
	json.NewEncoder(w).Encode(map[string]int{"likes": likes, "dislikes": dislikes})
}
