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
	"strconv"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	// Set headers to prevent caching
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
		}
	}

	categoryID, _ := strconv.Atoi(r.URL.Query().Get("category"))

	// Use the logged-in user's ID for filtering if no specific user is selected
	filterUserID := userID
	if r.URL.Query().Get("user") != "" {
		filterUserID, _ = strconv.Atoi(r.URL.Query().Get("user"))
	}

	likedOnly, _ := strconv.ParseBool(r.URL.Query().Get("liked"))

	// Only apply the likedOnly filter if the user is logged in
	if userID == 0 {
		likedOnly = false
	}

	posts, err := post.GetFilteredPosts(categoryID, filterUserID, likedOnly, 10)
	if err != nil {
		log.Printf("Error fetching filtered posts: %v", err)
		posts = []models.Post{}
	}

	// Fetch fresh like/dislike counts for each post
	for i, p := range posts {
		likes, dislikes, err := like.GetLikesCount(p.ID, "post")
		if err != nil {
			log.Printf("Error fetching like counts for post %d: %v", p.ID, err)
		} else {
			posts[i].Likes = likes
			posts[i].Dislikes = dislikes
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

	// Fetch all users
	allUsers, err := database.GetAllUsers()
	if err != nil {
		log.Printf("Error fetching all users: %v", err)
		allUsers = []models.User{}
	}

	pageData := models.PageData{
		Title: "Home - Literary Lions Forum",
		Page:  "home",
		User:  user,
		Data: map[string]interface{}{
			"Posts":             posts,
			"PopularCategories": popularCategories,
			"AllCategories":     allCategories,
			"SelectedCategory":  categoryID,
			"FilterUserID":      filterUserID,
			"LikedOnly":         likedOnly,
		},
		Users: allUsers,
	}

	utils.RenderTemplate(w, "home.html", pageData)
}
