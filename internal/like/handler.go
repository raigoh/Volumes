package like

import (
	"encoding/json"
	"fmt"
	"literary-lions-forum/internal/utils"
	"literary-lions-forum/pkg/session"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// LikeHandler handles HTTP POST requests for liking or disliking a post or comment.
func LikeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	sess, err := session.GetSession(w, r)
	if err != nil {
		http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
		return
	}

	userID := session.GetUserID(sess)
	if userID == 0 {
		http.Error(w, "Unauthorized: User ID is 0", http.StatusUnauthorized)
		return
	}

	targetID, err := strconv.Atoi(r.FormValue("target_id"))
	if err != nil {
		http.Error(w, "Invalid target ID", http.StatusBadRequest)
		utils.RenderErrorTemplate(w, err, http.StatusBadRequest, "Server error. Database is not givin us what we want.. That lite picky database!")
		return
	}

	targetType := r.FormValue("target_type")
	if targetType != "post" && targetType != "comment" {
		http.Error(w, "Invalid target type", http.StatusBadRequest)
		utils.RenderErrorTemplate(w, err, http.StatusBadRequest, "Server error. Database is not givin us what we want.. That lite picky database!")
		return
	}

	isLike, err := strconv.ParseBool(r.FormValue("is_like"))
	if err != nil {
		http.Error(w, "Invalid is_like value", http.StatusBadRequest)
		utils.RenderErrorTemplate(w, err, http.StatusBadRequest, "Server error. Database is not givin us what we want.. That lite picky database!")
		return
	}

	err = AddLike(userID, targetID, targetType, isLike)
	if err != nil {
		log.Printf("Error processing like/dislike: %v", err)
		http.Error(w, fmt.Sprintf("Error processing like/dislike: %v", err), http.StatusInternalServerError)
		utils.RenderErrorTemplate(w, err, http.StatusInternalServerError, "Server error. IDK what to do with this data")
		return
	}

	// Get the referring page
	referer := r.Header.Get("Referer")
	if referer == "" {
		referer = "/"
	}

	// If it's a comment, add an anchor to the URL
	if targetType == "comment" {
		if !strings.Contains(referer, "#") {
			referer += "#comment-" + strconv.Itoa(targetID)
		}
	}

	// Redirect back to the referring page
	http.Redirect(w, r, referer, http.StatusSeeOther)
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
		utils.RenderErrorTemplate(w, err, http.StatusInternalServerError, "Server error. IDK what to do with this data")
		return
	}

	// Get the updated like and dislike counts
	likes, dislikes, err := GetLikesCount(targetID, targetType)
	if err != nil {
		http.Error(w, "Error getting like count", http.StatusInternalServerError)
		utils.RenderErrorTemplate(w, err, http.StatusInternalServerError, "Server error. IDK what to do with this data")
		return
	}

	// Send the updated counts as a JSON response
	json.NewEncoder(w).Encode(map[string]int{"likes": likes, "dislikes": dislikes})
}
