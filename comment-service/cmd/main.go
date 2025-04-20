package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"

	"comment-service/internal"
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

	// Initialize handler
	handler := internal.NewCommentHandler(commentRepo, postLikeRepo, friendshipRepo)

	// Create router
	r := mux.NewRouter()

	// Comment routes
	r.HandleFunc("/posts/{postId}/comments", internal.AuthMiddleware(handler.CreateComment)).Methods("POST")
	r.HandleFunc("/posts/{postId}/comments", handler.GetCommentsByPost).Methods("GET")
	r.HandleFunc("/comments/{id}", internal.AuthMiddleware(handler.UpdateComment)).Methods("PUT")
	r.HandleFunc("/comments/{id}", internal.AuthMiddleware(handler.DeleteComment)).Methods("DELETE")
	r.HandleFunc("/comments/{id}/like", internal.AuthMiddleware(handler.LikeComment)).Methods("POST")
	r.HandleFunc("/comments/{id}/unlike", internal.AuthMiddleware(handler.UnlikeComment)).Methods("POST")

	// Post like routes
	r.HandleFunc("/posts/{postId}/like", internal.AuthMiddleware(handler.LikePost)).Methods("POST")
	r.HandleFunc("/posts/{postId}/unlike", internal.AuthMiddleware(handler.UnlikePost)).Methods("POST")
	r.HandleFunc("/posts/{postId}/likes", handler.GetPostLikes).Methods("GET")

	// Friendship routes
	r.HandleFunc("/friends/request", internal.AuthMiddleware(handler.SendFriendRequest)).Methods("POST")
	r.HandleFunc("/friends/requests", internal.AuthMiddleware(handler.GetFriendRequests)).Methods("GET")
	r.HandleFunc("/friends/requests/{id}/accept", internal.AuthMiddleware(handler.AcceptFriendRequest)).Methods("POST")
	r.HandleFunc("/friends", internal.AuthMiddleware(handler.GetFriends)).Methods("GET")

	// Health check
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}).Methods("GET")

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