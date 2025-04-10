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
		if r.Method == http.MethodGet {
			internal.ListPostsHandler(w, r)
		} else if r.Method == http.MethodPost {
			internal.AuthMiddleware(internal.CreatePostHandler)(w, r)
		}
	})

	log.Println("Post service running on :8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
