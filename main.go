package main

import (
	"fmt"
	"log"
	"net/http"

	"forum/handlers"
	"forum/models"
)

func main() {
	// Setup the database
	err := models.SetupDatabase()
	if err != nil {
		log.Fatalf("Failed to setup database: %v", err)
	}

	// Serve static files (CSS and templates)
	http.Handle("/templates/", http.StripPrefix("/templates/", http.FileServer(http.Dir("./templates/"))))
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("./css/"))))

	// Register handlers for specific paths
	http.HandleFunc("/", handlers.IndexHandler)
	http.HandleFunc("/home", handlers.HomePageHandler)
	http.HandleFunc("/signup", handlers.SignUpHandler) // Register SignUpHandler
	http.HandleFunc("/signin", handlers.SignInHandler)

	// Launch the server and open the default web browser
	fmt.Println("Server listening on port http://localhost:1703")

	err = http.ListenAndServe(":1703", nil)
	if err != nil {
		log.Fatal("ListenAndServe", err)
	}
}
