package handlers

import (
	"database/sql"
	"forum/models"
	"net/http"
)

func CommentUnlikeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "405: Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	commentID := r.FormValue("comment_id")
	userID, err := models.GetUserIDFromSession(r) // Get the logged-in user ID
	if err != nil {
		http.Error(w, "Please log in to perform this action", http.StatusUnauthorized)
		return
	}

	db, err := sql.Open("sqlite3", "storage/storage.db")
	if err != nil {
		http.Error(w, "500: Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	_, err = db.Exec("DELETE FROM comment_likes WHERE user_id = ? AND comment_id = ?", userID, commentID)
	if err != nil {
		http.Error(w, "500: Internal Server Error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/home", http.StatusSeeOther)
}
