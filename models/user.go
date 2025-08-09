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
