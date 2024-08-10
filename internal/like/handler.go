package like

import (
	"literary-lions-forum/pkg/database"
	"literary-lions-forum/pkg/session"
	"net/http"
	"strconv"
)

func LikeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		// Get the session using the custom session management package
		sessionData := session.GetSession(w, r)
		userID := sessionData["user_id"]

		postID, _ := strconv.Atoi(r.FormValue("post_id"))
		commentID, _ := strconv.Atoi(r.FormValue("comment_id"))
		isLike, _ := strconv.ParseBool(r.FormValue("is_like"))

		var err error
		if commentID != 0 {
			_, err = database.DB.Exec("INSERT OR REPLACE INTO likes (user_id, comment_id, is_like) VALUES (?, ?, ?)", userID, commentID, isLike)
		} else {
			_, err = database.DB.Exec("INSERT OR REPLACE INTO likes (user_id, post_id, is_like) VALUES (?, ?, ?)", userID, postID, isLike)
		}

		if err != nil {
			http.Error(w, "Error processing like/dislike", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, r.Header.Get("Referer"), http.StatusSeeOther)
	}
}
