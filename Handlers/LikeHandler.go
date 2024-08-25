package handlers

import (
	"database/sql"
	"net/http"
)

func LikeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "405: Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	postID := r.FormValue("post_id")
	userID := 1 // Assuming user ID is 1, this should be fetched from session or context

	db, err := sql.Open("sqlite3", "storage/storage.db")
	if err != nil {
		http.Error(w, "500: Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	_, err = db.Exec("INSERT INTO likes (user_id, post_id) VALUES (?, ?)", userID, postID)
	if err != nil {
		http.Error(w, "500: Internal Server Error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/home", http.StatusSeeOther)
}
