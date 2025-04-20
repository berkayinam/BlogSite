package internal

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
	"log"
	"sync"
)

var (
	postRepo *PostRepository
	repoOnce sync.Once
)

func initializeRepo() {
	repoOnce.Do(func() {
		postRepo = NewPostRepository()
	})
}

func CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	initializeRepo()
	
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	var post Post
	if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	post.Author = r.Header.Get("username")
	
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := postRepo.CreatePost(ctx, &post); err != nil {
		log.Printf("Failed to create post: %v", err)
		http.Error(w, "Failed to save post", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Post created"})
}

func ListPostsHandler(w http.ResponseWriter, r *http.Request) {
	initializeRepo()
	
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	posts, err := postRepo.GetAllPosts(ctx)
	if err != nil {
		log.Printf("Failed to get posts: %v", err)
		http.Error(w, "Failed to get posts", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}

func GetPostsByAuthorHandler(w http.ResponseWriter, r *http.Request) {
	initializeRepo()
	
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET allowed", http.StatusMethodNotAllowed)
		return
	}

	author := r.URL.Query().Get("author")
	if author == "" {
		http.Error(w, "Author parameter is required", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	posts, err := postRepo.GetPostsByAuthor(ctx, author)
	if err != nil {
		log.Printf("Failed to get posts by author: %v", err)
		http.Error(w, "Failed to get posts", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}

func UpdatePostHandler(w http.ResponseWriter, r *http.Request) {
	initializeRepo()
	
	if r.Method != http.MethodPut {
		http.Error(w, "Only PUT allowed", http.StatusMethodNotAllowed)
		return
	}

	title := r.URL.Query().Get("title")
	if title == "" {
		http.Error(w, "Title parameter is required", http.StatusBadRequest)
		return
	}

	var post Post
	if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := postRepo.UpdatePost(ctx, title, &post); err != nil {
		log.Printf("Failed to update post: %v", err)
		http.Error(w, "Failed to update post", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Post updated"})
}

func DeletePostHandler(w http.ResponseWriter, r *http.Request) {
	initializeRepo()
	
	if r.Method != http.MethodDelete {
		http.Error(w, "Only DELETE allowed", http.StatusMethodNotAllowed)
		return
	}

	title := r.URL.Query().Get("title")
	if title == "" {
		http.Error(w, "Title parameter is required", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := postRepo.DeletePost(ctx, title); err != nil {
		log.Printf("Failed to delete post: %v", err)
		http.Error(w, "Failed to delete post", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Post deleted"})
}

func SearchPostsHandler(w http.ResponseWriter, r *http.Request) {
	initializeRepo()
	
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET allowed", http.StatusMethodNotAllowed)
		return
	}

	query := r.URL.Query().Get("q")
	if query == "" {
		http.Error(w, "Search query parameter 'q' is required", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	posts, err := postRepo.SearchPosts(ctx, query)
	if err != nil {
		log.Printf("Failed to search posts: %v", err)
		http.Error(w, "Failed to search posts", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}
