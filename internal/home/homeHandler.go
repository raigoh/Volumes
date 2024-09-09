package home

import (
	"literary-lions-forum/internal/category"
	"literary-lions-forum/internal/like"
	"literary-lions-forum/internal/models"
	"literary-lions-forum/internal/post"
	"literary-lions-forum/internal/utils"
	"literary-lions-forum/pkg/database"
	"literary-lions-forum/pkg/session"
	"log"
	"net/http"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")

	sess, _ := session.GetSession(w, r)
	userID := session.GetUserID(sess)

	var user *models.User
	if userID != 0 {
		var err error
		user, err = session.GetUserByID(userID)
		if err != nil {
			log.Printf("Error fetching user: %v", err)
			utils.RenderErrorTemplate(w, err, http.StatusInternalServerError, "Server error, the database is being picky today.. We are still training him")
			return
		}
	}

	popularCategories, err := category.GetPopularCategories(5)
	if err != nil {
		log.Printf("Error fetching popular categories: %v", err)
		popularCategories = []models.Category{}
	}

	allCategories, err := category.GetCategories()
	if err != nil {
		log.Printf("Error fetching all categories: %v", err)
		allCategories = []models.Category{}
	}

	allUsers, err := database.GetAllUsers()
	if err != nil {
		log.Printf("Error fetching all users: %v", err)
		allUsers = []models.User{}
	}

	// Fetch the latest posts
	latestPosts, err := post.GetLatestPosts(5)
	if err != nil {
		log.Printf("Error fetching latest posts: %v", err)
		latestPosts = []models.Post{}
	}

	// Update like and dislike counts for each post
	for i, p := range latestPosts {
		likes, dislikes, err := like.GetLikesCount(p.ID, "post")
		if err != nil {
			log.Printf("Error fetching like counts for post %d: %v", p.ID, err)
		} else {
			latestPosts[i].Likes = likes
			latestPosts[i].Dislikes = dislikes
		}
	}

	pageData := models.PageData{
		Title: "Home - Literary Lions Forum",
		Page:  "home",
		User:  user,
		Data: map[string]interface{}{
			"PopularCategories": popularCategories,
			"AllCategories":     allCategories,
			"CurrentUserID":     userID,
			"Posts":             latestPosts,
		},
		Users: allUsers,
	}

	utils.RenderTemplate(w, "home.html", pageData)
}
