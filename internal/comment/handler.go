package comment

import (
	"literary-lions-forum/pkg/database"
	"literary-lions-forum/pkg/session"
	"net/http"
	"strconv"
)

func CreateCommentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		// Get the session using the custom session management package
		sessionData := session.GetSession(w, r)
		userID := sessionData["user_id"]

		postID, _ := strconv.Atoi(r.FormValue("post_id"))
		content := r.FormValue("content")

		_, err := database.DB.Exec("INSERT INTO comments (post_id, user_id, content) VALUES (?, ?, ?)", postID, userID, content)
		if err != nil {
			http.Error(w, "Error creating comment", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/post/"+strconv.Itoa(postID), http.StatusSeeOther)
	}
}
