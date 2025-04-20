package main

import (
	"log"
	"fmt"
	"os"
	"net/http"
	"post-service/internal"
)

func main() {
	internal.ConnectMongo()
	fmt.Println("JWT_KEY:", os.Getenv("JWT_SECRET")) // kontrol için

	http.HandleFunc("/posts", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			internal.ListPostsHandler(w, r)
		case http.MethodPost:
			internal.AuthMiddleware(internal.CreatePostHandler)(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Author's posts endpoint
	http.HandleFunc("/posts/author", internal.GetPostsByAuthorHandler)

	// Search posts endpoint
	http.HandleFunc("/posts/search", internal.SearchPostsHandler)

	// Post management endpoint (update and delete)
	http.HandleFunc("/posts/manage", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPut:
			internal.AuthMiddleware(internal.UpdatePostHandler)(w, r)
		case http.MethodDelete:
			internal.AuthMiddleware(internal.DeletePostHandler)(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	log.Println("Post service running on port 8082")
	log.Fatal(http.ListenAndServe(":8082", nil))
}
