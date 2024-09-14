package models

import (
	"time"
)

// User represents a user in the system
type User struct {
	ID        string    // Unique identifier for the user
	Username  string    // User's chosen username
	Email     string    // User's email address
	Password  string    // Hashed password for security
	CreatedAt time.Time // Timestamp of when the user account was created
	UpdatedAt time.Time // Timestamp of when the user account was last updated
}

// Session represents a user session
type Session struct {
	ID        string                 // Unique identifier for the session
	UserID    int                    // ID of the user this session belongs to
	Data      map[string]interface{} // Additional session data
	CreatedAt time.Time              // Timestamp of when the session was created
	ExpiresAt time.Time              // Timestamp of when the session expires
}

// Category represents a post category
type Category struct {
	ID   int    // Unique identifier for the category
	Name string // Name of the category
}

// Post represents a forum post
type Post struct {
	ID        int       // Unique identifier for the post
	User      User      // User who created the post
	Title     string    // Title of the post
	Content   string    // Content of the post
	CreatedAt time.Time // Timestamp of when the post was created
	Category  Category  // Category the post belongs to
	Likes     int       // Number of likes the post has received
	Dislikes  int       // Number of dislikes the post has received
}

// Comment represents a comment on a post
type Comment struct {
	ID        int       // Unique identifier for the comment
	PostID    int       // ID of the post this comment belongs to
	UserID    int       // ID of the user who made the comment
	Likes     int       // Number of likes the comment has received
	Dislikes  int       // Number of dislikes the comment has received
	Content   string    // Content of the comment
	CreatedAt time.Time // Timestamp of when the comment was created
	UpdatedAt time.Time // Timestamp of when the comment was last updated
	User      User      // User who made the comment
}

// Like represents a like or dislike on a post or comment
type Like struct {
	ID         int       // Unique identifier for the like
	UserID     int       // ID of the user who made the like/dislike
	TargetID   int       // ID of the target (post or comment)
	TargetType string    // Type of the target ("post" or "comment")
	IsLike     bool      // true for like, false for dislike
	CreatedAt  time.Time // Timestamp of when the like was created
}

// PostCategory represents the many-to-many relationship between posts and categories
type PostCategory struct {
	PostID     int // ID of the post
	CategoryID int // ID of the category
}

// PageData holds data to be passed to the templates for rendering
type PageData struct {
	Title          string                 // Title of the page
	Page           string                 // Name of the current page
	Error          string                 // Error message, if any
	Data           map[string]interface{} // Additional data for the page
	User           *User                  // Current logged-in user
	Post           *Post                  // Post data, if applicable
	Comments       []Comment              // List of comments, if applicable
	Categories     []Category             // List of categories
	TotalUsers     int                    // Total number of users in the system
	TotalPosts     int                    // Total number of posts in the system
	TotalComments  int                    // Total number of comments in the system
	ActiveUsers    int                    // Number of currently active users
	RecentActivity []Activity             // List of recent activities
	Users          []User                 // List of users, if applicable
	IsAdmin        bool
}

// Activity represents a recent activity in the forum
type Activity struct {
	Type     string // Type of activity (e.g., "new_post", "new_comment")
	Username string // Username of the user who performed the activity
	Content  string // Content related to the activity
	Date     string // Date of the activity
}
