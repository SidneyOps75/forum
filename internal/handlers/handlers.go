package handlers

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"forum/internal/models"

	"golang.org/x/crypto/bcrypt"
)

type Handler struct {
	DB *sql.DB
}

func NewHandler(db *sql.DB) *Handler {
	return &Handler{DB: db}
}

// Home displays the homepage with all posts
func (h *Handler) Home(w http.ResponseWriter, r *http.Request) {
	// Debug information
	dir, _ := os.Getwd()
	log.Printf("Current working directory: %s", dir)

	// Get the project root directory (one level up from cmd)
	projectRoot := filepath.Join(dir, "..")

	// Use absolute paths for templates
	baseTemplate := filepath.Join(projectRoot, "web", "templates", "base.html")
	homeTemplate := filepath.Join(projectRoot, "web", "templates", "home.html")

	if _, err := os.Stat(baseTemplate); os.IsNotExist(err) {
		log.Printf("Template file does not exist at: %s", baseTemplate)
	} else {
		log.Printf("Template file found at: %s", baseTemplate)
	}

	// Get posts
	posts, err := models.GetAllPosts(h.DB)
	if err != nil {
		log.Printf("Error fetching posts: %v", err)
		http.Error(w, "Failed to fetch posts", http.StatusInternalServerError)
		return
	}

	// Get categories for the filter
	categories, err := models.GetAllCategories(h.DB)
	if err != nil {
		log.Printf("Error fetching categories: %v", err)
		http.Error(w, "Failed to fetch categories", http.StatusInternalServerError)
		return
	}

	// Get current user if logged in
	var user *models.User
	if cookie, err := r.Cookie("session"); err == nil {
		user, _ = models.GetUserBySession(h.DB, cookie.Value)
	}

	// Prepare template data
	data := struct {
		Posts      []models.Post
		Categories []models.Category
		User       *models.User
	}{
		Posts:      posts,
		Categories: categories,
		User:       user,
	}

	// Parse and execute template
	tmpl, err := template.ParseFiles(baseTemplate, homeTemplate)
	if err != nil {
		log.Printf("Template parsing error: %v", err)
		http.Error(w, fmt.Sprintf("Template error: %v", err), http.StatusInternalServerError)
		return
	}

	err = tmpl.ExecuteTemplate(w, "base", data)
	if err != nil {
		log.Printf("Template execution error: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// Register handles user registration
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		email := r.FormValue("email")
		password := r.FormValue("password")

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Failed to hash password", http.StatusInternalServerError)
			return
		}

		_, err = h.DB.Exec("INSERT INTO users (username, email, password) VALUES (?, ?, ?)", username, email, hashedPassword)
		if err != nil {
			http.Error(w, "Username or email already taken", http.StatusBadRequest)
			return
		}

		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	tmpl := template.Must(template.ParseFiles("web/templates/register.html"))
	tmpl.Execute(w, nil)
}

// Login handles user login
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		email := r.FormValue("email")
		password := r.FormValue("password")

		var user models.User
		err := h.DB.QueryRow("SELECT id, username, password FROM users WHERE email = ?", email).Scan(&user.ID, &user.Username, &user.Password)
		if err != nil {
			http.Error(w, "Invalid email or password", http.StatusUnauthorized)
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
		if err != nil {
			http.Error(w, "Invalid email or password", http.StatusUnauthorized)
			return
		}

		// Set session cookie
		http.SetCookie(w, &http.Cookie{
			Name:  "session",
			Value: fmt.Sprintf("%d", user.ID),
			Path:  "/",
		})

		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	tmpl := template.Must(template.ParseFiles("web/templates/login.html"))
	tmpl.Execute(w, nil)
}

// Logout handles user logout
func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:   "session",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
