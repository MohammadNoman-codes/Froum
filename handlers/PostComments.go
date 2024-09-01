package handlers

import (
	"database/sql"
	"fmt"
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
	Username  string
	CreatedAt time.Time
	CLikes    int
	CDislikes int
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
		http.Error(w, "500: Internal Server Error first", http.StatusInternalServerError)
		return
	}

	// Fetch all comments for the post with the username
	rows, err := db.Query(`
    SELECT comments.id, comments.content, comments.user_id, users.username, comments.post_id, comments.created_at, 
    (SELECT IFNULL(COUNT(*), 0) FROM likes WHERE comment_id = comments.id) as CLikes,
    (SELECT IFNULL(COUNT(*), 0) FROM dislikes WHERE comment_id = comments.id) as CDislikes
    FROM comments 
    JOIN users ON comments.user_id = users.id 
    WHERE comments.post_id = ?`, postID)

	if err != nil {
		http.Error(w, "500: Internal Server Error Second ", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var comments []Comment
	for rows.Next() {
		var comment Comment
		comment.CLikes = 0    // Default to 0
		comment.CDislikes = 0 // Default to 0
		if err := rows.Scan(&comment.ID, &comment.Content, &comment.UserID, &comment.Username, &comment.PostID, &comment.CreatedAt, &comment.CLikes, &comment.CDislikes); err != nil {
			fmt.Println(err)
			http.Error(w, "500: Internal Server Error", http.StatusInternalServerError)
			return
		}
		comments = append(comments, comment)
	}

	postIDInt, err := strconv.Atoi(postID)
	if err != nil {
		fmt.Println(err, "yusuf")

		http.Error(w, "500: Internal Server Error No this is the error  4", http.StatusInternalServerError)
		return
	}

	postWithComments := PostWithComments{
		PostID:    postIDInt,
		PostTitle: postTitle,
		Comments:  comments,
	}

	tmpl, err := template.ParseFiles("./templates/comments.html")
	if err != nil {
		fmt.Println(err, "noora")

		http.Error(w, "500: Internal Server Error fifth", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, postWithComments)
	if err != nil {
		fmt.Println(err, "error in Excte")

		http.Error(w, "500: Internal Server Error sixth", http.StatusInternalServerError)
		return
	}
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
