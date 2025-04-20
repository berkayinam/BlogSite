package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"

	"myblog-project/comment-service/internal"
)

func main() {
	// Connect to MongoDB
	if err := internal.ConnectToMongo(); err != nil {
		log.Fatal(err)
	}

	// Initialize repositories
	commentRepo := internal.NewCommentRepository()
	friendshipRepo := internal.NewFriendshipRepository()
	postLikeRepo := internal.NewPostLikeRepository()

	// Create router
	r := mux.NewRouter()

	// Comment routes
	r.HandleFunc("/posts/{postId}/comments", internal.AuthMiddleware(commentRepo.CreateComment)).Methods("POST")
	r.HandleFunc("/posts/{postId}/comments", commentRepo.GetCommentsByPost).Methods("GET")
	r.HandleFunc("/comments/{id}", internal.AuthMiddleware(commentRepo.UpdateComment)).Methods("PUT")
	r.HandleFunc("/comments/{id}", internal.AuthMiddleware(commentRepo.DeleteComment)).Methods("DELETE")
	r.HandleFunc("/comments/{id}/like", internal.AuthMiddleware(commentRepo.LikeComment)).Methods("POST")
	r.HandleFunc("/comments/{id}/unlike", internal.AuthMiddleware(commentRepo.UnlikeComment)).Methods("POST")

	// Post like routes
	r.HandleFunc("/posts/{postId}/like", internal.AuthMiddleware(postLikeRepo.LikePost)).Methods("POST")
	r.HandleFunc("/posts/{postId}/unlike", internal.AuthMiddleware(postLikeRepo.UnlikePost)).Methods("POST")
	r.HandleFunc("/posts/{postId}/likes", postLikeRepo.GetPostLikes).Methods("GET")

	// Friendship routes
	r.HandleFunc("/friends/request", internal.AuthMiddleware(friendshipRepo.SendFriendRequest)).Methods("POST")
	r.HandleFunc("/friends/requests", internal.AuthMiddleware(friendshipRepo.GetFriendRequests)).Methods("GET")
	r.HandleFunc("/friends/requests/{id}/accept", internal.AuthMiddleware(friendshipRepo.AcceptFriendRequest)).Methods("POST")
	r.HandleFunc("/friends", internal.AuthMiddleware(friendshipRepo.GetFriends)).Methods("GET")

	// CORS configuration
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{
			"Accept",
			"Authorization",
			"Content-Type",
			"X-CSRF-Token",
		},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:          300,
	})

	// Server configuration
	port := os.Getenv("PORT")
	if port == "" {
		port = "8083" // Default port for comment service
	}

	srv := &http.Server{
		Handler:      c.Handler(r),
		Addr:         ":" + port,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Printf("Comment service starting on port %s", port)
	log.Fatal(srv.ListenAndServe())
} 