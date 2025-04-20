package internal

import (
	"context"
	"testing"
	"time"
	"fmt"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func getTestMongoURI() string {
	// Check environment variable first
	if uri := os.Getenv("TEST_MONGO_URI"); uri != "" {
		return uri
	}
	// Default to local MongoDB for testing
	return "mongodb://localhost:27017"
}

func waitForMongo(ctx context.Context, client *mongo.Client, t *testing.T, maxRetries int) error {
	var lastErr error
	for i := 0; i < maxRetries; i++ {
		if err := client.Ping(ctx, readpref.Primary()); err != nil {
			lastErr = err
			t.Logf("Attempt %d: MongoDB not ready, retrying... Error: %v", i+1, err)
			time.Sleep(2 * time.Second) // Wait longer between retries
			continue
		}
		return nil
	}
	return fmt.Errorf("failed to connect after %d retries: %v", maxRetries, lastErr)
}

func setupTestDB(t *testing.T) (*mongo.Client, func()) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	mongoURI := getTestMongoURI()
	t.Logf("Connecting to MongoDB at %s", mongoURI)

	clientOpts := options.Client().
		ApplyURI(mongoURI).
		SetServerSelectionTimeout(10 * time.Second).
		SetConnectTimeout(10 * time.Second)

	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		t.Fatalf("Failed to create MongoDB client: %v", err)
	}

	// Wait for MongoDB to be ready
	if err := waitForMongo(ctx, client, t, 10); err != nil {
		t.Fatalf("MongoDB is not available: %v\nMake sure MongoDB is running and accessible at %s", err, mongoURI)
	}

	t.Log("Successfully connected to MongoDB")

	// Clean up any existing data first
	db := client.Database(databaseName)
	if err := db.Collection(collectionName).Drop(ctx); err != nil {
		t.Logf("Warning: Failed to drop collection: %v", err)
	}

	// Return a cleanup function
	cleanup := func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := client.Database(databaseName).Collection(collectionName).Drop(ctx); err != nil {
			t.Logf("Warning: Failed to clean up test database: %v", err)
		}
		if err := client.Disconnect(ctx); err != nil {
			t.Logf("Warning: Failed to disconnect test client: %v", err)
		}
	}

	return client, cleanup
}

func TestCreatePost(t *testing.T) {
	client, cleanup := setupTestDB(t)
	defer cleanup()

	// Override the global Client for testing
	Client = client
	repo := NewPostRepository()

	ctx := context.Background()
	testPost := &Post{
		Title:   "Test Post",
		Content: "Test Content",
		Author:  "Test Author",
	}

	// Test creating a post
	err := repo.CreatePost(ctx, testPost)
	if err != nil {
		t.Errorf("Failed to create post: %v", err)
	}

	// Verify the post was created
	var savedPost Post
	err = repo.collection.FindOne(ctx, map[string]string{"title": "Test Post"}).Decode(&savedPost)
	if err != nil {
		t.Errorf("Failed to find created post: %v", err)
	}

	if savedPost.Title != testPost.Title {
		t.Errorf("Expected title %s, got %s", testPost.Title, savedPost.Title)
	}
}

func TestGetAllPosts(t *testing.T) {
	client, cleanup := setupTestDB(t)
	defer cleanup()

	Client = client
	repo := NewPostRepository()
	ctx := context.Background()

	// Create test posts
	testPosts := []Post{
		{Title: "Post 1", Content: "Content 1", Author: "Author 1", CreatedAt: time.Now()},
		{Title: "Post 2", Content: "Content 2", Author: "Author 2", CreatedAt: time.Now()},
	}

	for _, post := range testPosts {
		_, err := repo.collection.InsertOne(ctx, post)
		if err != nil {
			t.Fatalf("Failed to insert test post: %v", err)
		}
	}

	// Test getting all posts
	posts, err := repo.GetAllPosts(ctx)
	if err != nil {
		t.Errorf("Failed to get posts: %v", err)
	}

	if len(posts) != len(testPosts) {
		t.Errorf("Expected %d posts, got %d", len(testPosts), len(posts))
	}
}

func TestGetPostsByAuthor(t *testing.T) {
	client, cleanup := setupTestDB(t)
	defer cleanup()

	Client = client
	repo := NewPostRepository()
	ctx := context.Background()

	// Create test posts with different authors
	testPosts := []Post{
		{Title: "Post 1", Content: "Content 1", Author: "Author 1"},
		{Title: "Post 2", Content: "Content 2", Author: "Author 1"},
		{Title: "Post 3", Content: "Content 3", Author: "Author 2"},
	}

	for _, post := range testPosts {
		_, err := repo.collection.InsertOne(ctx, post)
		if err != nil {
			t.Fatalf("Failed to insert test post: %v", err)
		}
	}

	// Test getting posts by author
	posts, err := repo.GetPostsByAuthor(ctx, "Author 1")
	if err != nil {
		t.Errorf("Failed to get posts by author: %v", err)
	}

	if len(posts) != 2 {
		t.Errorf("Expected 2 posts for Author 1, got %d", len(posts))
	}
}

func TestUpdatePost(t *testing.T) {
	client, cleanup := setupTestDB(t)
	defer cleanup()

	Client = client
	repo := NewPostRepository()
	ctx := context.Background()

	// Create a test post
	originalPost := Post{
		Title:   "Original Title",
		Content: "Original Content",
		Author:  "Author",
	}

	_, err := repo.collection.InsertOne(ctx, originalPost)
	if err != nil {
		t.Fatalf("Failed to insert test post: %v", err)
	}

	// Update the post
	updatedPost := &Post{
		Title:   "Original Title",
		Content: "Updated Content",
		Author:  "Author",
	}

	err = repo.UpdatePost(ctx, "Original Title", updatedPost)
	if err != nil {
		t.Errorf("Failed to update post: %v", err)
	}

	// Verify the update
	var savedPost Post
	err = repo.collection.FindOne(ctx, map[string]string{"title": "Original Title"}).Decode(&savedPost)
	if err != nil {
		t.Errorf("Failed to find updated post: %v", err)
	}

	if savedPost.Content != "Updated Content" {
		t.Errorf("Expected content %s, got %s", "Updated Content", savedPost.Content)
	}
}

func TestDeletePost(t *testing.T) {
	client, cleanup := setupTestDB(t)
	defer cleanup()

	Client = client
	repo := NewPostRepository()
	ctx := context.Background()

	// Create a test post
	testPost := Post{
		Title:   "Test Post",
		Content: "Test Content",
		Author:  "Author",
	}

	_, err := repo.collection.InsertOne(ctx, testPost)
	if err != nil {
		t.Fatalf("Failed to insert test post: %v", err)
	}

	// Delete the post
	err = repo.DeletePost(ctx, "Test Post")
	if err != nil {
		t.Errorf("Failed to delete post: %v", err)
	}

	// Verify the deletion
	count, err := repo.collection.CountDocuments(ctx, map[string]string{"title": "Test Post"})
	if err != nil {
		t.Errorf("Failed to count documents: %v", err)
	}

	if count != 0 {
		t.Errorf("Expected 0 posts after deletion, got %d", count)
	}
} 