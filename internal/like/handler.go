package like

import (
	"encoding/json"
	"fmt"
	"literary-lions-forum/pkg/session"
	"log"
	"net/http"
	"strconv"
)

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
		return
	}

	targetType := r.FormValue("target_type")
	if targetType != "post" && targetType != "comment" {
		http.Error(w, "Invalid target type", http.StatusBadRequest)
		return
	}

	isLike, err := strconv.ParseBool(r.FormValue("is_like"))
	if err != nil {
		http.Error(w, "Invalid is_like value", http.StatusBadRequest)
		return
	}

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

func UnLikeHandler(w http.ResponseWriter, r *http.Request) {
	sess, err := session.GetSession(w, r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	userID := session.GetUserID(sess)
	if userID == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	targetID, _ := strconv.Atoi(r.FormValue("target_id"))
	targetType := r.FormValue("target_type")

	err = RemoveLike(userID, targetID, targetType)
	if err != nil {
		http.Error(w, "Error removing like", http.StatusInternalServerError)
		return
	}

	likes, dislikes, err := GetLikesCount(targetID, targetType)
	if err != nil {
		http.Error(w, "Error getting like count", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]int{"likes": likes, "dislikes": dislikes})
}
