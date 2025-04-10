package internal

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
	"log"
)

func CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	var post Post
	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	post.Author = r.Header.Get("username")
	post.CreatedAt = time.Now()

	collection := Client.Database("postdb").Collection("posts")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = collection.InsertOne(ctx, post)
	if err != nil {
		log.Printf("errrrrr %v", err)
		http.Error(w, "Failed to save post", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Post created"})
}

func ListPostsHandler(w http.ResponseWriter, r *http.Request) {
	collection := Client.Database("postdb").Collection("posts")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, map[string]interface{}{})
	if err != nil {
		http.Error(w, "Failed to get posts", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	var posts []Post
	for cursor.Next(ctx) {
		var post Post
		cursor.Decode(&post)
		posts = append(posts, post)
	}

	json.NewEncoder(w).Encode(posts)
}
