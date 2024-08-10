package like

import (
	"literary-lions-forum/pkg/database"
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
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	userID := session.GetUserID(sess)
	if userID == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	postID, _ := strconv.Atoi(r.FormValue("post_id"))
	commentID, _ := strconv.Atoi(r.FormValue("comment_id"))
	isLike, _ := strconv.ParseBool(r.FormValue("is_like"))

	if commentID != 0 {
		_, err = database.DB.Exec("INSERT OR REPLACE INTO likes (user_id, comment_id, is_like) VALUES (?, ?, ?)", userID, commentID, isLike)
	} else if postID != 0 {
		_, err = database.DB.Exec("INSERT OR REPLACE INTO likes (user_id, post_id, is_like) VALUES (?, ?, ?)", userID, postID, isLike)
	} else {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if err != nil {
		log.Printf("Error processing like/dislike: %v", err)
		http.Error(w, "Error processing like/dislike", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, r.Header.Get("Referer"), http.StatusSeeOther)
}
