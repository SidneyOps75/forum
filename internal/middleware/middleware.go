package middleware

import (
	"database/sql"
	"net/http"
	"time"
)

// AuthMiddleware ensures the user is authenticated
func AuthMiddleware(db *sql.DB) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("session")
			if err != nil {
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}

			var expiresAt time.Time
			err = db.QueryRow(
				"SELECT expires_at FROM sessions WHERE session_token = ?",
				cookie.Value,
			).Scan(&expiresAt)

			if err != nil || time.Now().After(expiresAt) {
				http.SetCookie(w, &http.Cookie{
					Name:   "session",
					Value:  "",
					Path:   "/",
					MaxAge: -1,
				})
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}

			next(w, r)
		}
	}
}
