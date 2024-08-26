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

	// Public routes
	http.HandleFunc("/", handlers.IndexHandler)
	http.HandleFunc("/signup", handlers.SignUpHandler) // Register SignUpHandler
	http.HandleFunc("/signin", handlers.SignInHandler)

	// Protected routes - wrapped with SessionMiddleware
	protectedMux := http.NewServeMux()
	protectedMux.Handle("/home", handlers.SessionMiddleware(http.HandlerFunc(handlers.HomePageHandler)))
	protectedMux.Handle("/addpost", handlers.SessionMiddleware(http.HandlerFunc(handlers.AddPostHandler)))
	protectedMux.Handle("/logout", handlers.SessionMiddleware(http.HandlerFunc(handlers.LogoutHandler)))
	protectedMux.Handle("/like", handlers.SessionMiddleware(http.HandlerFunc(handlers.LikeHandler)))
	protectedMux.Handle("/dislike", handlers.SessionMiddleware(http.HandlerFunc(handlers.DislikeHandler)))
	protectedMux.Handle("/undislike", handlers.SessionMiddleware(http.HandlerFunc(handlers.UndislikeHandler)))
	protectedMux.Handle("/unlike", handlers.SessionMiddleware(http.HandlerFunc(handlers.UnlikeHandler)))
	protectedMux.Handle("/comments", handlers.SessionMiddleware(http.HandlerFunc(handlers.CommentsHandler)))
	protectedMux.Handle("/addcomment", handlers.SessionMiddleware(http.HandlerFunc(handlers.AddCommentHandler)))

	// Use the protected mux for routes that require authentication
	http.Handle("/home", protectedMux)
	http.Handle("/addpost", protectedMux)
	http.Handle("/logout", protectedMux)
	http.Handle("/like", protectedMux)
	http.Handle("/dislike", protectedMux)
	http.Handle("/undislike", protectedMux)
	http.Handle("/unlike", protectedMux)
	http.Handle("/comments", protectedMux)
	http.Handle("/addcomment", protectedMux)

	// Launch the server
	fmt.Println("Server listening on port http://localhost:1703")

	err = http.ListenAndServe(":1703", nil)
	if err != nil {
		log.Fatal("ListenAndServe", err)
	}
}
