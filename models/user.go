package models

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"fmt"
	"time"
)

type User struct {
	Username string
	IsAdmin  bool
	Password []byte
}

func GetUser(username string) (User, error) {
	var u User
	err := db.QueryRow("Select * FROM users WHERE username=?", username).
		Scan(&u.Username, &u.IsAdmin, &u.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return User{}, fmt.Errorf("user not found! %w", err)
		}
		return User{}, fmt.Errorf("query failed! %w", err)
	}
	return u, nil
}

func GenerateSessionToken(username string) (string, error) {
	b := make([]byte, 32) // 256 bit
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	token := base64.URLEncoding.EncodeToString(b)

	_, err = db.Exec("INSERT INTO sessions (token, username, expired_at) VALUES (?, ?, ?);", token, username, time.Now().Add(8*time.Hour))
	if err != nil {
		return "", fmt.Errorf("error inserting session token! %w", err)
	}
	return token, nil
}

func ValidateSession(token string) (*User, error) {
	var username string

	// Check if session exists and is not expired
	err := db.QueryRow(
		`SELECT username FROM sessions 
         WHERE token = ? AND expired_at > CURRENT_TIMESTAMP`,
		token).Scan(&username)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("invalid session")
		}
		return nil, fmt.Errorf("session validation failed: %w", err)
	}

	// Get the full user details
	user, err := GetUser(username)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	return &user, nil
}

func DeleteSession(token string) error {
	_, err := db.Exec(`DELETE FROM sessions WHERE token = ?`, token)
	return err
}
