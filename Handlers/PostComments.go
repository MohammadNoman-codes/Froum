package handlers

import (
	"database/sql"
	"forum/models"
	"html/template"
	"net/http"
	"strconv"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Comment struct {
	ID        int
	Content   string
	UserID    int
	PostID    int
	CreatedAt time.Time
}

type PostWithComments struct {
	PostID    int
	PostTitle string
	Comments  []Comment
	UserID    int
}

// Fetch comments and display the comments page
func CommentsHandler(w http.ResponseWriter, r *http.Request) {
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

	// Fetch the post title (optional, for display)
	var postTitle string
	err = db.QueryRow("SELECT title FROM posts WHERE id = ?", postID).Scan(&postTitle)
	if err != nil {
		http.Error(w, "500: Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Fetch all comments for the post
	rows, err := db.Query("SELECT id, content, user_id, post_id, created_at FROM comments WHERE post_id = ?", postID)
	if err != nil {
		http.Error(w, "500: Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var comments []Comment
	for rows.Next() {
		var comment Comment
		if err := rows.Scan(&comment.ID, &comment.Content, &comment.UserID, &comment.PostID, &comment.CreatedAt); err != nil {
			http.Error(w, "500: Internal Server Error", http.StatusInternalServerError)
			return
		}
		comments = append(comments, comment)
	}

	postIDInt, err := strconv.Atoi(postID)
	if err != nil {
		http.Error(w, "500: Internal Server Error", http.StatusInternalServerError)
		return
	}

	postWithComments := PostWithComments{
		PostID:    postIDInt,
		PostTitle: postTitle,
		Comments:  comments,
	}

	tmpl, err := template.ParseFiles("./templates/comments.html")
	if err != nil {
		http.Error(w, "500: Internal Server Error", http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, postWithComments)
}

// Add a new comment to the database
func AddCommentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/home", http.StatusSeeOther)
		return
	}

	userID, err := models.GetUserIDFromSession(r)
	if err != nil {
		http.Error(w, "403: Forbidden", http.StatusForbidden)
		return
	}
	postID := r.FormValue("post_id")
	content := r.FormValue("content")

	db, err := sql.Open("sqlite3", "./storage/storage.db")
	if err != nil {
		http.Error(w, "500: Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	_, err = db.Exec("INSERT INTO comments (content, user_id, post_id, created_at) VALUES (?, ?, ?, ?)", content, userID, postID, time.Now())
	if err != nil {
		http.Error(w, "500: Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Redirect back to the comments page
	http.Redirect(w, r, "/comments?post_id="+postID, http.StatusSeeOther)
}
