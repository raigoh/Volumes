package post

import (
	"literary-lions-forum/internal/models"
	"literary-lions-forum/internal/utils"
	"literary-lions-forum/pkg/database"
	"literary-lions-forum/pkg/session"
	"net/http"
)

func CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		utils.RenderTemplate(w, "base.html", models.PageData{
			Title: "Create Post - Literary Lions Forum",
			Page:  "create_post",
		})
		return
	}

	if r.Method == http.MethodPost {
		// Get the session using the custom session management package
		sessionData := session.GetSession(w, r)
		userID := sessionData["user_id"]

		title := r.FormValue("title")
		content := r.FormValue("content")

		_, err := database.DB.Exec("INSERT INTO posts (user_id, title, content) VALUES (?, ?, ?)", userID, title, content)
		if err != nil {
			http.Error(w, "Error creating post", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}
