package handlers

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type GuestComment struct {
	ID        int
	Content   string
	UserID    int
	PostID    int
	Username  string
	CreatedAt time.Time
}

type GuestPostWithComments struct {
	PostID    int
	PostTitle string
	Comments  []GuestComment
}

// GuestCommentsHandler displays the comments for a post to guests without allowing new comments to be added
func GuestCommentsHandler(w http.ResponseWriter, r *http.Request) {
	postID := r.URL.Query().Get("post_id")
	if postID == "" {
		http.Error(w, "400: Bad Request", http.StatusBadRequest)
		return
	}

	db, err := sql.Open("sqlite3", "./storage/storage.db")
	if err != nil {
		http.Error(w, "500: Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	var postTitle string
	err = db.QueryRow("SELECT title FROM posts WHERE id = ?", postID).Scan(&postTitle)
	if err != nil {
		http.Error(w, "500: Internal Server Error", http.StatusInternalServerError)
		return
	}

	rows, err := db.Query(`
        SELECT comments.id, comments.content, comments.user_id, users.username, comments.post_id, comments.created_at 
        FROM comments 
        JOIN users ON comments.user_id = users.id 
        WHERE comments.post_id = ?`, postID)
	if err != nil {
		http.Error(w, "500: Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var comments []GuestComment
	for rows.Next() {
		var comment GuestComment
		if err := rows.Scan(&comment.ID, &comment.Content, &comment.UserID, &comment.Username, &comment.PostID, &comment.CreatedAt); err != nil {
			fmt.Println("Error scanning row:", err)
			http.Error(w, "500: Internal Server Error", http.StatusInternalServerError)
			return
		}
		comments = append(comments, comment)
	}

	postIDInt, err := strconv.Atoi(postID)
	if err != nil {
		fmt.Println("Invalid Post ID Conversion:", err)
		http.Error(w, "500: Internal Server Error", http.StatusInternalServerError)
		return
	}

	postWithComments := GuestPostWithComments{
		PostID:    postIDInt,
		PostTitle: postTitle,
		Comments:  comments,
	}

	tmpl, err := template.ParseFiles("./templates/guestcomments.html")
	if err != nil {
		fmt.Println("Error parsing template:", err)
		http.Error(w, "500: Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, postWithComments)
	if err != nil {
		fmt.Println("Error executing template:", err)
		http.Error(w, "500: Internal Server Error", http.StatusInternalServerError)
		return
	}
}
