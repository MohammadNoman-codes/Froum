package handlers

import (
	"html/template"
	"net/http"
)

func HomePageHandler(w http.ResponseWriter, r *http.Request) {
	// Only allow GET requests
	if r.Method != http.MethodGet {
		http.Error(w, "400: Bad Request", http.StatusBadRequest)
		return
	}

	// Serve error page if the request path is not "/homePage"
	if r.URL.Path != "/home" {
		t, err := template.ParseFiles("templates/error.html")
		if err != nil {
			http.Error(w, "500: internal server error", http.StatusInternalServerError)
			return
		}
		t.Execute(w, http.StatusNotFound)
		return
	}

	// Parse the homePage template
	t, err := template.ParseFiles("templates/homePage.html")
	if err != nil {
		http.Error(w, "500: internal server error", http.StatusInternalServerError)
		return
	}

	// Execute the homePage template
	err = t.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
