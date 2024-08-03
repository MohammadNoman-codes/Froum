package handlers

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

type Post struct {
	ID      int
	Title   string
	Content string
}

func FetchAllPosts() ([]Post, error) {
	// Connect to the database
	db, err := sql.Open("sqlite3", "storage/storage.db")
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}
	defer db.Close()

	// Query the database to fetch all posts
	rows, err := db.Query("SELECT id, title, content FROM posts")
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %v", err)
	}
	defer rows.Close()

	// Create an empty slice to store the posts
	posts := []Post{}

	// Iterate over the rows and populate the posts slice
	for rows.Next() {
		var post Post
		err := rows.Scan(&post.ID, &post.Title, &post.Content)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			return nil, fmt.Errorf("failed to scan row: %v", err)
		}
		posts = append(posts, post)

	}

	// Check for any errors during iteration
	if err = rows.Err(); err != nil {
		log.Printf("Row iteration error: %v", err)
		return nil, fmt.Errorf("row iteration error: %v", err)
	}

	// Return the posts slice
	return posts, nil
}

func HomePageHandler(w http.ResponseWriter, r *http.Request) {
	// Only allow GET requests
	if r.Method != http.MethodGet {
		http.Error(w, "405: Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Serve error page if the request path is not "/home"
	if r.URL.Path != "/home" {
		t, err := template.ParseFiles("templates/error.html")
		if err != nil {
			http.Error(w, "500: Internal Server Error", http.StatusInternalServerError)
			return
		}
		t.Execute(w, nil)
		return
	}

	// Parse the homePage template
	t, err := template.ParseFiles("templates/homePage.html")
	if err != nil {
		http.Error(w, "500: Internal Server Error (This is the Parse File)", http.StatusInternalServerError)
		return
	}

	Posts, err := FetchAllPosts()
	if err != nil {
		http.Error(w, "500: Internal Server Error (This is because of Fetching all Posts)", http.StatusInternalServerError)
		return
	}

	if len(Posts) == 0 {
		// No posts found, pass a message or a flag to the template
		err = t.Execute(w, map[string]interface{}{
			"NoPosts": true,
		})
	} else {
		// Execute the homePage template with the posts
		err = t.Execute(w, map[string]interface{}{
			"Posts": Posts,
		})
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
