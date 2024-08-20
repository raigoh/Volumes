package models

import (
	"time"
)

type User struct {
	ID        string
	Username  string
	Email     string
	Password  string // Hashed password
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Session struct {
	ID        string
	UserID    int
	Data      map[string]interface{}
	CreatedAt time.Time
	ExpiresAt time.Time
}

type Category struct {
	ID   int
	Name string
}

type Post struct {
	ID        int
	UserID    int
	Title     string
	Content   string
	CreatedAt time.Time
	UpdatedAt time.Time
	User      User
	Category  Category
}

type Comment struct {
	ID        int
	PostID    int
	UserID    int
	Content   string
	CreatedAt time.Time
	UpdatedAt time.Time
	User      User
}

type Like struct {
	ID         int
	UserID     int
	TargetID   int
	TargetType string // "post" or "comment"
	IsLike     bool   // true for like, false for dislike
	CreatedAt  time.Time
}

type PostCategory struct {
	PostID     int
	CategoryID int
}

// PageData holds data to be passed to the templates
type PageData struct {
	Title          string
	Page           string
	Error          string
	Data           map[string]interface{}
	User           *User
	Post           *Post
	Comments       []Comment
	Categories     []Category
	TotalUsers     int
	TotalPosts     int
	TotalComments  int
	ActiveUsers    int
	RecentActivity []Activity
	Users          []User
}

type Activity struct {
	Type     string
	Username string
	Content  string
	Date     string
}
