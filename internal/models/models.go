package models

package models

import (
	"time"
)

type User struct {
	ID        int64
	Username  string
	Email     string
	Password  string // Hashed password
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Session struct {
	ID        string // UUID
	UserID    int64
	ExpiresAt time.Time
}

type Category struct {
	ID   int64
	Name string
}

type Post struct {
	ID        int64
	UserID    int64
	Title     string
	Content   string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Comment struct {
	ID        int64
	PostID    int64
	UserID    int64
	Content   string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Like struct {
	ID         int64
	UserID     int64
	TargetID   int64
	TargetType string // "post" or "comment"
	IsLike     bool   // true for like, false for dislike
	CreatedAt  time.Time
}

type PostCategory struct {
	PostID     int64
	CategoryID int64
}