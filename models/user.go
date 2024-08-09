package models

import (
	"database/sql"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// User represents the structure of a user in the database.
type User struct {
	ID        int
	Username  string
	Email     string
	Password  string
	CreatedAt time.Time
}

// CreateUser inserts a new user into the database with a hashed password.
func CreateUser(db *sql.DB, user *User) error {
	// Hash the user's password before storing it
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	query := `INSERT INTO users (username, email, password, created_at) VALUES (?, ?, ?, ?)`
	_, err = db.Exec(query, user.Username, user.Email, hashedPassword, time.Now())
	return err
}

// GetUserByEmail retrieves a user from the database by email.
func GetUserByEmail(db *sql.DB, email string) (*User, error) {
	user := &User{}
	query := `SELECT id, username, email, password, created_at FROM users WHERE email = ?`
	err := db.QueryRow(query, email).Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil // No user found
	}
	return user, err
}

// GetUserByID retrieves a user from the database by ID.
func GetUserByID(db *sql.DB, id int) (*User, error) {
	user := &User{}
	query := `SELECT id, username, email, created_at FROM users WHERE id = ?`
	err := db.QueryRow(query, id).Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil // No user found
	}
	return user, err
}

// CheckPassword compares the provided password with the stored hashed password.
func (user *User) CheckPassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
}
