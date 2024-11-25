package comment

import (
	"literary-lions-forum/internal/utils"
	"literary-lions-forum/pkg/database"
	"literary-lions-forum/pkg/session"
	"log"
	"net/http"
	"strconv"
)

// CreateCommentHandler handles the HTTP request for creating a new comment.
// It performs authentication, validation, and directly inserts the comment into the database.
func CreateCommentHandler(w http.ResponseWriter, r *http.Request) {
	// Simulate an intentional error to test error handling
	// panic("intentional error for testing")

	// Ensure the request method is POST
	if r.Method != http.MethodPost {

		//http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		utils.RenderErrorTemplate(w, nil, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Get the user's session
	sess, err := session.GetSession(w, r)
	if err != nil {
		//http.Error(w, "Unauthorized", http.StatusUnauthorized)
		utils.RenderErrorTemplate(w, nil, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Get the user ID from the session
	userID := session.GetUserID(sess)
	if userID == 0 {
		//http.Error(w, "Unauthorized", http.StatusUnauthorized)
		utils.RenderErrorTemplate(w, nil, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Parse and validate the post ID from the form data
	postID, err := strconv.Atoi(r.FormValue("post_id"))
	if err != nil {
		//http.Error(w, "Invalid post ID", http.StatusBadRequest)
		utils.RenderErrorTemplate(w, err, http.StatusBadRequest, "Invalid post ID")
		return
	}

	// Get and validate the comment content
	content := r.FormValue("content")
	if content == "" {
		//http.Error(w, "Comment content cannot be empty", http.StatusBadRequest)
		utils.RenderErrorTemplate(w, err, http.StatusBadRequest, "Comment content cannot be empty")
		return
	}

	// Insert the comment directly into the database
	_, err = database.DB.Exec("INSERT INTO comments (post_id, user_id, content) VALUES (?, ?, ?)", postID, userID, content)
	if err != nil {
		log.Printf("Error creating comment: %v", err)
		//http.Error(w, "Error creating comment", http.StatusInternalServerError)
		utils.RenderErrorTemplate(w, err, http.StatusInternalServerError, "Error creating comment")
		return
	}

	// Redirect to the post page after successful comment creation
	http.Redirect(w, r, "/post/"+strconv.Itoa(postID), http.StatusSeeOther)
}

// AddCommentHandler handles the HTTP request for adding a new comment.
// It performs authentication, validation, and uses the CreateComment function to add the comment.
func AddCommentHandler(w http.ResponseWriter, r *http.Request) {
	// Ensure the request method is POST
	if r.Method != http.MethodPost {
		//http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		utils.RenderErrorTemplate(w, nil, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Get the user's session
	sess, err := session.GetSession(w, r)
	if err != nil {
		//http.Error(w, "Unauthorized", http.StatusUnauthorized)
		utils.RenderErrorTemplate(w, err, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Get the user ID from the session
	userID := session.GetUserID(sess)
	if userID == 0 {
		//http.Error(w, "Unauthorized", http.StatusUnauthorized)
		utils.RenderErrorTemplate(w, nil, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Parse and validate the post ID from the form data
	postID, err := strconv.Atoi(r.FormValue("post_id"))
	if err != nil {
		//http.Error(w, "Invalid post ID", http.StatusBadRequest)
		utils.RenderErrorTemplate(w, err, http.StatusBadRequest, "Invalid post ID")
		return
	}

	// Get and validate the comment content
	content := r.FormValue("content")
	if content == "" {
		//http.Error(w, "Comment content cannot be empty", http.StatusBadRequest)
		utils.RenderErrorTemplate(w, nil, http.StatusBadRequest, "Comment content cannot be empty")
		return
	}

	// Use the CreateComment function to add the comment
	err = CreateComment(userID, postID, content)
	if err != nil {
		log.Printf("Error creating comment: %v", err)
		//http.Error(w, "Error creating comment", http.StatusInternalServerError)
		utils.RenderErrorTemplate(w, err, http.StatusInternalServerError, "Error creating comment")
		return
	}

	// Redirect to the post page after successful comment creation
	http.Redirect(w, r, "/post/"+strconv.Itoa(postID), http.StatusSeeOther)
}
