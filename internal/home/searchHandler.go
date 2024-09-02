package home

import (
	"literary-lions-forum/internal/models"
	"literary-lions-forum/internal/post"
	"literary-lions-forum/internal/utils"
	"literary-lions-forum/pkg/session"
	"log"
	"net/http"
)

func SearchHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")
	log.Printf("Searching for: %s", query)

	sess, _ := session.GetSession(w, r)
	userID := session.GetUserID(sess)

	var user *models.User
	if userID != 0 {
		var err error
		user, err = session.GetUserByID(userID)
		if err != nil {
			log.Printf("Error fetching user: %v", err)
		}
	}

	posts, err := post.SearchPosts(query, 10)
	if err != nil {
		log.Printf("Error searching posts: %v", err)
		posts = []models.Post{}
	}
	log.Printf("Found %d posts", len(posts))

	pageData := models.PageData{
		Title: "Search Results - Literary Lions Forum",
		Page:  "search",
		User:  user,
		Data: map[string]interface{}{
			"Posts":       posts,
			"SearchQuery": query,
		},
	}

	utils.RenderTemplate(w, "search.html", pageData)
}
