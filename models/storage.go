package models

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"net/http"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func SetupDatabase() error {
	// Open SQLite database
	db, err := sql.Open("sqlite3", "./storage/storage.db")
	if err != nil {
		return err
	}
	defer db.Close()

	// Read setup.sql file
	setupSQL, err := os.ReadFile("./storage/setup.sql")
	if err != nil {
		return err
	}

	// Execute SQL statements
	_, err = db.Exec(string(setupSQL))
	if err != nil {
		return err
	}

	return nil
}

func CreateUser(email, username, password string) error {
	db, err := sql.Open("sqlite3", "./storage/storage.db")
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare("INSERT INTO users (email, username, password) VALUES (?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(email, username, password)
	if err != nil {
		return err
	}

	return nil
}

func AuthenticateUser(email, password string, w http.ResponseWriter) (bool, error) {
	db, err := sql.Open("sqlite3", "./storage/storage.db")
	if err != nil {
		return false, err
	}
	defer db.Close()

	var userID int
	var storedPassword string
	err = db.QueryRow("SELECT id, password FROM users WHERE email = ?", email).Scan(&userID, &storedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil // User not found
		}
		return false, err
	}

	// Compare passwords (in a real application, use bcrypt or similar for secure password storage)
	if storedPassword != password {
		return false, nil // Passwords do not match
	}

	// Generate a session ID
	sessionID := generateSessionID()

	// Set the session ID as a cookie
	cookie := &http.Cookie{
		Name:    "session_id",
		Value:   sessionID,
		Expires: time.Now().Add(24 * time.Hour), // Cookie expires in 24 hours
	}
	http.SetCookie(w, cookie)

	// Store the session in the database
	_, err = db.Exec("INSERT INTO sessions (session_id, user_id, expires_at) VALUES (?, ?, ?)", sessionID, userID, time.Now().Add(24*time.Hour))
	if err != nil {
		return false, err
	}

	return true, nil // Authentication and session creation successful
}

func generateSessionID() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}
