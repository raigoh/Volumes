package post

import (
	"literary-lions-forum/internal/category"
	"literary-lions-forum/internal/like"
	"literary-lions-forum/internal/models"
	"literary-lions-forum/internal/utils"
	"literary-lions-forum/pkg/database"
	"literary-lions-forum/pkg/session"
	"log"
	"net/http"
	"strconv"
)

func AllPostsHandler(w http.ResponseWriter, r *http.Request) {
	sess, _ := session.GetSession(w, r)
	userID := session.GetUserID(sess)

	var user *models.User
	if userID != 0 {
		var err error
		user, err = session.GetUserByID(userID)
		if err != nil {
			utils.RenderErrorTemplate(w, err, http.StatusInternalServerError, "Error fetching user data")
			return
		}
	}

	var categoryID int
	var filterUserID int
	var likedOnly bool
	var searchQuery string

	if r.Method == http.MethodPost {
		r.ParseForm()
		searchQuery = r.FormValue("query")
		categoryID, _ = strconv.Atoi(r.FormValue("category"))
		filterUserID, _ = strconv.Atoi(r.FormValue("user"))
		likedOnly = r.FormValue("liked") == "true"
	} else {
		categoryID, _ = strconv.Atoi(r.URL.Query().Get("category"))
		filterUserID, _ = strconv.Atoi(r.URL.Query().Get("user"))
		likedOnly, _ = strconv.ParseBool(r.URL.Query().Get("liked"))
		searchQuery = r.URL.Query().Get("query")
	}

	if userID == 0 {
		likedOnly = false
	}

	var posts []models.Post
	var err error

	if searchQuery != "" {
		posts, err = EnhancedSearch(searchQuery, 10)
	} else {
		posts, err = GetFilteredPosts(categoryID, filterUserID, likedOnly, 10)
	}

	if err != nil {
		log.Printf("Error fetching posts: %v", err)
		posts = []models.Post{}
	}

	for i, p := range posts {
		likes, dislikes, err := like.GetLikesCount(p.ID, "post")
		if err != nil {
			log.Printf("Error fetching like counts for post %d: %v", p.ID, err)
			utils.RenderErrorTemplate(w, err, http.StatusInternalServerError, "Server error, who would like this even?!?")
		} else {
			posts[i].Likes = likes
			posts[i].Dislikes = dislikes
		}
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

	data := models.PageData{
		Title: "All Posts - Literary Lions Forum",
		Page:  "all-posts",
		User:  user,
		Data: map[string]interface{}{
			"Posts":            posts,
			"AllCategories":    allCategories,
			"SelectedCategory": categoryID,
			"FilterUserID":     filterUserID,
			"LikedOnly":        likedOnly,
			"SearchQuery":      searchQuery,
		},
		Users: allUsers,
	}

	utils.RenderTemplate(w, "all-posts.html", data)
}
