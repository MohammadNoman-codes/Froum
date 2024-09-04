package handlers

import (
	"fmt"
	"net/http"

	"forum/models"
)

func SignUpHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	email := strings.trimSpace(r.FormValue("email"))
	username :=strings.trimSpace( r.FormValue("username"))
	password := strings.trimSpace( r.FormValue("password"))

	if email || username || password === "" {
		http.Error (w, "can not have empty fields please enter the data into it", http.statusBadRequest)
		return
	}

	err := models.CreateUser(email, username, password)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create user: %v", err), http.StatusInternalServerError)
		return
	}

	// Redirect to index.html after successful registration
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
