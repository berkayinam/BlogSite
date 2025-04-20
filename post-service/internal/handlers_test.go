package internal

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
)

func setupTestHandler(t *testing.T) (*PostRepository, func()) {
	client, cleanup := setupTestDB(t)
	Client = client
	
	// Create a new repository instance for testing
	repo := NewPostRepository()
	postRepo = repo // Set the global repository

	// Clean up any existing data
	ctx := context.Background()
	if err := repo.collection.Drop(ctx); err != nil {
		t.Logf("Warning: Failed to clean up existing data: %v", err)
	}

	return repo, func() {
		postRepo = nil // Clear the global repository
		cleanup()      // Clean up the database
	}
}

func TestCreatePostHandler(t *testing.T) {
	_, cleanup := setupTestHandler(t)
	defer cleanup()

	tests := []struct {
		name           string
		input         Post
		expectedCode  int
		withAuth      bool
	}{
		{
			name: "Valid Post",
			input: Post{
				Title:   "Test Post",
				Content: "Test Content",
			},
			expectedCode: http.StatusCreated,
			withAuth:     true,
		},
		{
			name: "No Auth",
			input: Post{
				Title:   "Test Post",
				Content: "Test Content",
			},
			expectedCode: http.StatusUnauthorized,
			withAuth:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.input)
			req := httptest.NewRequest(http.MethodPost, "/posts", bytes.NewBuffer(body))
			
			if tt.withAuth {
				req.Header.Set("username", "testuser")
			}
			
			w := httptest.NewRecorder()
			
			if tt.withAuth {
				CreatePostHandler(w, req)
			} else {
				AuthMiddleware(CreatePostHandler)(w, req)
			}

			if w.Code != tt.expectedCode {
				t.Errorf("Expected status code %d, got %d", tt.expectedCode, w.Code)
			}
		})
	}
}

func TestListPostsHandler(t *testing.T) {
	repo, cleanup := setupTestHandler(t)
	defer cleanup()

	// Create some test posts
	testPosts := []Post{
		{Title: "Post 1", Content: "Content 1", Author: "Author 1", CreatedAt: time.Now()},
		{Title: "Post 2", Content: "Content 2", Author: "Author 2", CreatedAt: time.Now()},
	}

	ctx := context.Background()
	for _, post := range testPosts {
		err := repo.CreatePost(ctx, &post)
		if err != nil {
			t.Fatalf("Failed to insert test post: %v", err)
		}
	}

	req := httptest.NewRequest(http.MethodGet, "/posts", nil)
	w := httptest.NewRecorder()

	ListPostsHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var response []Post
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if len(response) != len(testPosts) {
		t.Errorf("Expected %d posts, got %d", len(testPosts), len(response))
	}
}

func TestGetPostsByAuthorHandler(t *testing.T) {
	repo, cleanup := setupTestHandler(t)
	defer cleanup()

	// Create test posts
	testPosts := []Post{
		{Title: "Post 1", Content: "Content 1", Author: "Author 1"},
		{Title: "Post 2", Content: "Content 2", Author: "Author 1"},
		{Title: "Post 3", Content: "Content 3", Author: "Author 2"},
	}

	ctx := context.Background()
	for _, post := range testPosts {
		err := repo.CreatePost(ctx, &post)
		if err != nil {
			t.Fatalf("Failed to insert test post: %v", err)
		}
	}

	tests := []struct {
		name         string
		author       string
		expectedCode int
		postCount    int
	}{
		{
			name:         "Valid Author",
			author:       "Author 1",
			expectedCode: http.StatusOK,
			postCount:    2,
		},
		{
			name:         "Missing Author",
			author:       "",
			expectedCode: http.StatusBadRequest,
			postCount:    0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reqURL := "/posts/author"
			if tt.author != "" {
				params := url.Values{}
				params.Add("author", tt.author)
				reqURL = reqURL + "?" + params.Encode()
			}

			req := httptest.NewRequest(http.MethodGet, reqURL, nil)
			w := httptest.NewRecorder()

			GetPostsByAuthorHandler(w, req)

			if w.Code != tt.expectedCode {
				t.Errorf("Expected status code %d, got %d", tt.expectedCode, w.Code)
			}

			if tt.expectedCode == http.StatusOK {
				var response []Post
				if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
					t.Fatalf("Failed to decode response: %v", err)
				}

				if len(response) != tt.postCount {
					t.Errorf("Expected %d posts, got %d", tt.postCount, len(response))
				}

				// Verify all posts are from the requested author
				for _, post := range response {
					if post.Author != tt.author {
						t.Errorf("Expected author %s, got %s", tt.author, post.Author)
					}
				}
			}
		})
	}
}

func TestUpdatePostHandler(t *testing.T) {
	repo, cleanup := setupTestHandler(t)
	defer cleanup()

	// Create initial post
	initialPost := Post{
		Title:   "Initial Post",
		Content: "Initial Content",
		Author:  "Author 1",
	}

	ctx := context.Background()
	err := repo.CreatePost(ctx, &initialPost)
	if err != nil {
		t.Fatalf("Failed to insert test post: %v", err)
	}

	tests := []struct {
		name         string
		title        string
		input        Post
		expectedCode int
		withAuth     bool
	}{
		{
			name:  "Valid Update",
			title: "Initial Post",
			input: Post{
				Title:   "Initial Post",
				Content: "Updated Content",
				Author:  "Author 1",
			},
			expectedCode: http.StatusOK,
			withAuth:     true,
		},
		{
			name:  "No Auth",
			title: "Initial Post",
			input: Post{
				Title:   "Initial Post",
				Content: "Updated Content",
				Author:  "Author 1",
			},
			expectedCode: http.StatusUnauthorized,
			withAuth:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.input)
			params := url.Values{}
			params.Add("title", tt.title)
			reqURL := "/posts/manage?" + params.Encode()

			req := httptest.NewRequest(http.MethodPut, reqURL, bytes.NewBuffer(body))
			
			if tt.withAuth {
				req.Header.Set("username", "testuser")
			}
			
			w := httptest.NewRecorder()
			
			if tt.withAuth {
				UpdatePostHandler(w, req)
			} else {
				AuthMiddleware(UpdatePostHandler)(w, req)
			}

			if w.Code != tt.expectedCode {
				t.Errorf("Expected status code %d, got %d", tt.expectedCode, w.Code)
			}

			if tt.expectedCode == http.StatusOK {
				// Verify the update
				post, err := repo.GetPostsByAuthor(ctx, tt.input.Author)
				if err != nil {
					t.Errorf("Failed to verify update: %v", err)
				}
				if len(post) == 0 {
					t.Error("Post not found after update")
				} else if post[0].Content != tt.input.Content {
					t.Errorf("Expected content %s, got %s", tt.input.Content, post[0].Content)
				}
			}
		})
	}
}

func TestDeletePostHandler(t *testing.T) {
	repo, cleanup := setupTestHandler(t)
	defer cleanup()

	// Create a post to delete
	testPost := Post{
		Title:   "Test Post",
		Content: "Test Content",
		Author:  "Author",
	}

	ctx := context.Background()
	err := repo.CreatePost(ctx, &testPost)
	if err != nil {
		t.Fatalf("Failed to insert test post: %v", err)
	}

	tests := []struct {
		name         string
		title        string
		expectedCode int
		withAuth     bool
	}{
		{
			name:         "Valid Delete",
			title:        "Test Post",
			expectedCode: http.StatusOK,
			withAuth:     true,
		},
		{
			name:         "No Auth",
			title:        "Test Post",
			expectedCode: http.StatusUnauthorized,
			withAuth:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params := url.Values{}
			params.Add("title", tt.title)
			reqURL := "/posts/manage?" + params.Encode()

			req := httptest.NewRequest(http.MethodDelete, reqURL, nil)
			
			if tt.withAuth {
				req.Header.Set("username", "testuser")
			}
			
			w := httptest.NewRecorder()
			
			if tt.withAuth {
				DeletePostHandler(w, req)
			} else {
				AuthMiddleware(DeletePostHandler)(w, req)
			}

			if w.Code != tt.expectedCode {
				t.Errorf("Expected status code %d, got %d", tt.expectedCode, w.Code)
			}

			if tt.expectedCode == http.StatusOK {
				// Verify the deletion
				posts, err := repo.GetAllPosts(ctx)
				if err != nil {
					t.Errorf("Failed to verify deletion: %v", err)
				}
				if len(posts) > 0 {
					t.Error("Post still exists after deletion")
				}
			}
		})
	}
} 