package handlers

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"

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
	posts, err := models.GetAllPosts(h.DB)
	if err != nil {
		http.Error(w, "Failed to fetch posts", http.StatusInternalServerError)
		return
	}

	tmpl := template.Must(template.ParseFiles("web/templates/home.html"))
	tmpl.Execute(w, posts)
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
