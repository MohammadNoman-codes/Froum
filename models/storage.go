package models

import (
	"database/sql"
	"os"

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

func AuthenticateUser(email, password string) (bool, error) {
	db, err := sql.Open("sqlite3", "./storage/storage.db")
	if err != nil {
		return false, err
	}
	defer db.Close()

	var storedPassword string
	err = db.QueryRow("SELECT password FROM users WHERE email = ?", email).Scan(&storedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil // User not found
		}
		return false, err
	}

	// Compare passwords (in a real application, use bcrypt or similar for secure password storage)
	if storedPassword == password {
		return true, nil // Authentication successful
	}

	return false, nil // Passwords do not match
}
