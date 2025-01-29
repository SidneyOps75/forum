package models

import (
	"database/sql"
	"time"
)

type User struct {
	ID       int
	Username string
	Email    string
	Password string
}

type Post struct {
	ID         int
	UserID     int
	Title      string
	Content    string
	CreatedAt  time.Time
	Categories []Category
	Likes      int
	Dislikes   int
	Username   string // for display purposes
}

type Comment struct {
	ID        int
	UserID    int
	PostID    int
	Content   string
	CreatedAt time.Time
	Likes     int
	Dislikes  int
}

type Category struct {
	ID   int
	Name string
}

type Session struct {
	ID           int
	UserID       int
	SessionToken string
	ExpiresAt    time.Time
}

// GetAllPosts fetches all posts from the database
func GetAllPosts(db *sql.DB) ([]Post, error) {
	rows, err := db.Query("SELECT id, user_id, title, content, created_at FROM posts")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var post Post
		err := rows.Scan(&post.ID, &post.UserID, &post.Title, &post.Content, &post.CreatedAt)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	return posts, nil
}

// GetAllCategories fetches all categories from the database
func GetAllCategories(db *sql.DB) ([]Category, error) {
	rows, err := db.Query("SELECT id, name FROM categories")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []Category
	for rows.Next() {
		var category Category
		err := rows.Scan(&category.ID, &category.Name)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}

	return categories, nil
}

// GetUserBySession retrieves a user based on their session token
func GetUserBySession(db *sql.DB, sessionToken string) (*User, error) {
	var user User
	err := db.QueryRow(`
		SELECT u.id, u.username, u.email 
		FROM users u 
		JOIN sessions s ON u.id = s.user_id 
		WHERE s.session_token = ? AND s.expires_at > datetime('now')`,
		sessionToken,
	).Scan(&user.ID, &user.Username, &user.Email)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &user, nil
}
