package post

import (
	"literary-lions-forum/internal/category"
	"literary-lions-forum/internal/comment"
	"literary-lions-forum/internal/models"
	"literary-lions-forum/internal/utils"
	"literary-lions-forum/pkg/database"
	"literary-lions-forum/pkg/session"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// CreatePostHandler handles both GET and POST requests for creating a new post
func CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	// Handle GET request: render the new post form
	if r.Method == http.MethodGet {
		utils.RenderTemplate(w, "new-post.html", models.PageData{
			Title: "Create Post - Literary Lions Forum",
			Page:  "create_post",
		})
		return
	}

	// Ensure the request method is POST
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get the user's session and check authorization
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

	// Get form values
	title := r.FormValue("title")
	content := r.FormValue("content")

	// Validate form input
	if title == "" || content == "" {
		http.Error(w, "Title and content cannot be empty", http.StatusBadRequest)
		return
	}

	// Insert the new post into the database
	_, err = database.DB.Exec("INSERT INTO posts (user_id, title, content) VALUES (?, ?, ?)", userID, title, content)
	if err != nil {
		log.Printf("Error creating post: %v", err)
		http.Error(w, "Error creating post", http.StatusInternalServerError)
		return
	}

	// Redirect to the home page after successful post creation
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// NewPostHandler handles both GET and POST requests for creating a new post (with category selection)
func NewPostHandler(w http.ResponseWriter, r *http.Request) {
	// Get the user's session and check authorization
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

	// Handle GET request: render the new post form with categories
	if r.Method == http.MethodGet {
		categories, err := category.GetCategories()
		if err != nil {
			log.Printf("Error fetching categories: %v", err)
			http.Error(w, "Error fetching categories", http.StatusInternalServerError)
			return
		}

		utils.RenderTemplate(w, "new-post.html", models.PageData{
			Title:      "Create New Post - Literary Lions Forum",
			Page:       "new-post",
			Categories: categories,
		})
		return
	}

	// Handle POST request: process the new post submission
	if r.Method == http.MethodPost {
		// Get form values
		title := r.FormValue("title")
		content := r.FormValue("content")
		categoryID, err := strconv.Atoi(r.FormValue("category"))
		if err != nil {
			http.Error(w, "Invalid category", http.StatusBadRequest)
			return
		}

		// Validate form input
		if title == "" || content == "" {
			http.Error(w, "Title and content cannot be empty", http.StatusBadRequest)
			return
		}

		// Create the new post
		postID, err := CreatePost(userID, categoryID, title, content)
		if err != nil {
			log.Printf("Error creating post: %v", err)
			http.Error(w, "Error creating post", http.StatusInternalServerError)
			return
		}

		// Redirect to the newly created post's detail page
		http.Redirect(w, r, "/post/"+strconv.Itoa(postID), http.StatusSeeOther)
		return
	}

	// Handle any other HTTP methods
	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}

// PostDetailHandler handles GET requests for viewing a specific post and its comments
func PostDetailHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the post ID from the URL
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 {
		http.NotFound(w, r)
		return
	}
	postID, err := strconv.Atoi(parts[2])
	if err != nil {
		http.NotFound(w, r)
		return
	}

	// Fetch the post details
	post, err := GetPostByID(postID)
	if err != nil {
		log.Printf("Error fetching post: %v", err)
		http.Error(w, "Error fetching post", http.StatusInternalServerError)
		return
	}

	// Fetch comments for the post
	comments, err := comment.GetCommentsByPostID(postID)
	if err != nil {
		log.Printf("Error fetching comments: %v", err)
		http.Error(w, "Error fetching comments", http.StatusInternalServerError)
		return
	}

	// Render the post detail template with post and comments data
	utils.RenderTemplate(w, "post-detail.html", models.PageData{
		Title:    post.Title + " - Literary Lions Forum",
		Page:     "post-detail",
		Post:     post,
		Comments: comments,
	})
}

// PostListHandler handles HTTP requests for the post list page
func PostListHandler(w http.ResponseWriter, r *http.Request) {
	// Fetch posts from the database
	posts, err := database.GetPosts()
	if err != nil {
		http.Error(w, "Error fetching posts", http.StatusInternalServerError)
		return
	}

	// Prepare data for the template
	data := models.PageData{
		Title: "Literary Lions Forum",
		Page:  "home",
		Data: map[string]interface{}{
			"Posts": posts,
		},
	}

	// Render the template
	utils.RenderTemplate(w, "home.html", data)
}
