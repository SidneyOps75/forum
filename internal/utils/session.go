package utils

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
)

func GenerateSessionToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func CreateSession(db *sql.DB, userID int64) (string, error) {
	token, err := GenerateSessionToken()
	if err != nil {
		return "", err
	}

	_, err = db.Exec(
		"INSERT INTO sessions (user_id, session_token, expires_at) VALUES (?, ?, datetime('now', '+24 hours'))",
		userID, token,
	)
	if err != nil {
		return "", err
	}

	return token, nil
}
