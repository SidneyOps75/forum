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
	ID        int
	UserID    int
	Title     string
	Content   string
	CreatedAt time.Time
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
