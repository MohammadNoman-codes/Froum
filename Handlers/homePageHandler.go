package handlers

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type Post struct {
	ID            int
	Title         string
	Content       string
	Category      string
	LikesCount    int
	DislikesCount int
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
		rows, err = db.Query(`
			SELECT posts.id, posts.title, posts.content, 
			       (SELECT COUNT(*) FROM likes WHERE post_id = posts.id) AS likes_count, 
			       (SELECT COUNT(*) FROM dislikes WHERE post_id = posts.id) AS dislikes_count
			FROM posts`)
	} else {
		rows, err = db.Query(`
			SELECT posts.id, posts.title, posts.content, 
			       (SELECT COUNT(*) FROM likes WHERE post_id = posts.id) AS likes_count, 
			       (SELECT COUNT(*) FROM dislikes WHERE post_id = posts.id) AS dislikes_count
			FROM posts WHERE category = ?`, category)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %v", err)
	}
	defer rows.Close()

	posts := []Post{}
	for rows.Next() {
		var post Post
		err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.LikesCount, &post.DislikesCount)
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
