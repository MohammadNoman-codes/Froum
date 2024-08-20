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
	ID       int
	Title    string
	Content  string
	Category string
}

// FetchPosts fetches posts based on the selected category.
func FetchPosts(category string) ([]Post, error) {
	db, err := sql.Open("sqlite3", "storage/storage.db")
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}
	defer db.Close()

	var rows *sql.Rows
	if category == "" {
		rows, err = db.Query("SELECT id, title, content FROM posts")
	} else {
		rows, err = db.Query("SELECT id, title, content FROM posts WHERE category = ?", category)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %v", err)
	}
	defer rows.Close()

	posts := []Post{}
	for rows.Next() {
		var post Post
		err := rows.Scan(&post.ID, &post.Title, &post.Content)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			return nil, fmt.Errorf("failed to scan row: %v", err)
		}
		posts = append(posts, post)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Row iteration error: %v", err)
		return nil, fmt.Errorf("row iteration error: %v", err)
	}

	return posts, nil
}

func HomePageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "405: Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	if r.URL.Path != "/home" {
		t, err := template.ParseFiles("templates/error.html")
		if err != nil {
			http.Error(w, "500: Internal Server Error", http.StatusInternalServerError)
			return
		}
		t.Execute(w, nil)
		return
	}

	category := r.URL.Query().Get("category")

	t, err := template.ParseFiles("templates/homePage.html")
	if err != nil {
		http.Error(w, "500: Internal Server Error", http.StatusInternalServerError)
		return
	}

	posts, err := FetchPosts(category)
	if err != nil {
		http.Error(w, "500: Internal Server Error", http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"Posts":            posts,
		"SelectedCategory": category,
	}

	if len(posts) == 0 {
		data["NoPosts"] = true
	}

	err = t.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
