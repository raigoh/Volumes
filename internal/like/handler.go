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
// It ensures the request is valid, checks user authentication, and processes the like or dislike action.
func LikeHandler(w http.ResponseWriter, r *http.Request) {
	// Ensure the request method is POST. If not, return a "Method not allowed" error.
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Retrieve the user session to identify the user making the request.
	sess, err := session.GetSession(w, r)
	if err != nil {
		// If session retrieval fails, return an "Unauthorized" error.
		http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
		return
	}

	// Get the user ID from the session to ensure the user is logged in.
	userID := session.GetUserID(sess)
	if userID == 0 {
		// If the user ID is 0, they are not logged in, so return an "Unauthorized" error.
		http.Error(w, "Unauthorized: User ID is 0", http.StatusUnauthorized)
		return
	}

	// Parse the `target_id` from the form data. This is the ID of the post or comment to be liked/disliked.
	targetID, err := strconv.Atoi(r.FormValue("target_id"))
	if err != nil {
		// If the target ID is not a valid integer, return a "Bad Request" error and display a custom error message.
		http.Error(w, "Invalid target ID", http.StatusBadRequest)
		utils.RenderErrorTemplate(w, err, http.StatusBadRequest, "Server error. Database is not givin us what we want.. That lite picky database!")
		return
	}

	// Get the `target_type` from the form data (either "post" or "comment").
	targetType := r.FormValue("target_type")
	// Ensure that the target type is valid. If it's neither "post" nor "comment", return a "Bad Request" error.
	if targetType != "post" && targetType != "comment" {
		http.Error(w, "Invalid target type", http.StatusBadRequest)
		utils.RenderErrorTemplate(w, err, http.StatusBadRequest, "Server error. Database is not givin us what we want.. That lite picky database!")
		return
	}

	// Parse the `is_like` value from the form data to determine whether it's a like (true) or dislike (false).
	isLike, err := strconv.ParseBool(r.FormValue("is_like"))
	if err != nil {
		// If the `is_like` value is invalid (not true/false), return a "Bad Request" error.
		http.Error(w, "Invalid is_like value", http.StatusBadRequest)
		utils.RenderErrorTemplate(w, err, http.StatusBadRequest, "Server error. Database is not givin us what we want.. That lite picky database!")
		return
	}

	// Call the `AddLike` function to process the like or dislike for the specified user, target (post/comment), and target type.
	err = AddLike(userID, targetID, targetType, isLike)
	if err != nil {
		// Log the error and return a "Server Error" response if there is an issue with adding the like or dislike.
		log.Printf("Error processing like/dislike: %v", err)
		http.Error(w, fmt.Sprintf("Error processing like/dislike: %v", err), http.StatusInternalServerError)
		utils.RenderErrorTemplate(w, err, http.StatusInternalServerError, "Server error. IDK what to do with this data")
		return
	}

	// Retrieve the referring page (where the user came from), so the user can be redirected back after liking/disliking.
	referer := r.Header.Get("Referer")
	if referer == "" {
		// If no referrer is provided, redirect the user to the homepage ("/").
		referer = "/"
	}

	// Add an anchor to the URL for both posts and comments
	if strings.Contains(referer, "/all-posts") {
		referer += "#post-" + strconv.Itoa(targetID)
	} else if targetType == "comment" && !strings.Contains(referer, "#") {
		referer += "#comment-" + strconv.Itoa(targetID)
	}

	// Redirect the user back to the referring page or specific comment after processing the like/dislike action.
	http.Redirect(w, r, referer, http.StatusSeeOther)
}

// UnLikeHandler handles HTTP requests for removing a like or dislike from a post or comment.
// It checks user authentication and removes the like/dislike, then sends the updated like and dislike counts in response.
func UnLikeHandler(w http.ResponseWriter, r *http.Request) {
	// Retrieve the user session to verify that the user is logged in.
	sess, err := session.GetSession(w, r)
	if err != nil {
		// If session retrieval fails, return an "Unauthorized" error.
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get the user ID from the session to ensure the user is logged in.
	userID := session.GetUserID(sess)
	if userID == 0 {
		// If the user is not logged in, return an "Unauthorized" error.
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Parse the `target_id` from the form data to identify the post or comment from which the like/dislike is to be removed.
	targetID, _ := strconv.Atoi(r.FormValue("target_id"))
	// Get the `target_type` from the form data (either "post" or "comment").
	targetType := r.FormValue("target_type")

	// Call the `RemoveLike` function to remove the like/dislike for the given user, target, and target type.
	err = RemoveLike(userID, targetID, targetType)
	if err != nil {
		// If there is an error in removing the like/dislike, return a "Server Error" and log the error.
		http.Error(w, "Error removing like", http.StatusInternalServerError)
		utils.RenderErrorTemplate(w, err, http.StatusInternalServerError, "Server error. IDK what to do with this data")
		return
	}

	// After removing the like/dislike, retrieve the updated like and dislike counts for the target.
	likes, dislikes, err := GetLikesCount(targetID, targetType)
	if err != nil {
		// If there is an error fetching the like/dislike count, return a "Server Error".
		http.Error(w, "Error getting like count", http.StatusInternalServerError)
		utils.RenderErrorTemplate(w, err, http.StatusInternalServerError, "Server error. IDK what to do with this data")
		return
	}

	// Encode the updated like and dislike counts as a JSON object and send them as a response.
	json.NewEncoder(w).Encode(map[string]int{"likes": likes, "dislikes": dislikes})
}
