package main

import (
	"fmt"
	"log"
	"net/http"

	"forum/internal/database"
	"forum/internal/handlers"
)

func main() {
	// Initialize the database
	db, err := database.InitDB("forum.db")
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Initialize handlers with the database
	handler := handlers.NewHandler(db)

	// Set up routes
	http.HandleFunc("/", handler.Home)
	http.HandleFunc("/register", handler.Register)
	http.HandleFunc("/login", handler.Login)
	http.HandleFunc("/logout", handler.Logout)
	// http.HandleFunc("/post/create", middleware.AuthMiddleware(handler.CreatePost))
	// http.HandleFunc("/post/", handler.ViewPost)
	// http.HandleFunc("/comment", middleware.AuthMiddleware(handler.CreateComment))
	// http.HandleFunc("/like", middleware.AuthMiddleware(handler.LikePostOrComment))
	// http.HandleFunc("/filter", handler.FilterPosts)

	// Serve static files
	fs := http.FileServer(http.Dir("web/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Start the server
	fmt.Println("Server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
